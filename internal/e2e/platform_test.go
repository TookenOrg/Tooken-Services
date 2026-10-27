//go:build integration

package e2e

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/handlers"
	"github.com/TookenOrg/tooken-services/internal/api/router"
	"github.com/TookenOrg/tooken-services/internal/api/server"
	contracts "github.com/TookenOrg/tooken-services/internal/blockchain/contracts/bindings"
	chainglobals "github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	chainutils "github.com/TookenOrg/tooken-services/internal/blockchain/utils"
	appglobals "github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/testenv"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// pollTimeout bounds every wait on work the platform finishes in the
// background. A delivery on an auto-mining node takes well under a second;
// a minute leaves room for a slow machine without hiding a stuck order.
const pollTimeout = 60 * time.Second

// requestTimeout bounds a single HTTP call. The slowest synchronous one is a
// sign-up, whose bcrypt hash takes seconds under the race detector, longer
// when several investors sign up at once. Without a bound, a handler stuck on
// a lock would hang the run until `go test -timeout` dumps every goroutine,
// instead of failing on the request that never answered.
const requestTimeout = 2 * time.Minute

// account is somebody holding a session on the platform.
type account struct {
	id     int
	email  string
	token  string
	wallet common.Address // set once the investor is registered on chain
}

// platform is one running instance: the HTTP server, its database and the
// node it signs on.
type platform struct {
	srv     *httptest.Server
	db      *sql.DB
	eth     *ethclient.Client
	handler *handlers.Handler
	signer  common.Address
	admin   account
}

// reply is an HTTP answer, kept raw so that each caller decides what it
// expects of the body.
type reply struct {
	code int
	body []byte
}

// ledger is what a refused request must leave untouched: the orders on
// record, and the nonce of the key the platform signs with. A nonce that did
// not move is the proof that nothing was broadcast, whatever the logs say.
type ledger struct {
	orders int
	nonce  uint64
}

// assetSpec is what a manager types to list an asset.
type assetSpec struct {
	title          string
	totalShares    int
	pricePerShare  string
	entryFeeRate   string
	maxPerInvestor int
}

var emailSequence atomic.Int64

// bootPlatform starts the server the binary runs, on a database of its own
// and on the local node, and signs up its first administrator.
//
// Nothing is deployed yet: initialising the platform is part of what the
// end-to-end test proves.
func bootPlatform(t *testing.T) *platform {
	t.Helper()

	ctx := context.Background()

	// Registered first, so that its cleanup (dropping the database) runs last:
	// t.Cleanup is LIFO, and the background jobs must be waited for before the
	// database they write to disappears.
	db := testenv.ScratchDatabase(t, "tooken_it_e2e")
	logger.Init(true)
	client := testenv.DialNode(ctx, t)

	previousDB, previousClient := appglobals.DB, chainglobals.EthClient
	appglobals.DB, chainglobals.EthClient = db, client
	t.Setenv("PRIVATE_KEY", testenv.HardhatAccount0Key)

	previousWriter, previousMode := gin.DefaultWriter, gin.Mode()
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = io.Discard

	t.Cleanup(func() {
		gin.DefaultWriter = previousWriter
		gin.SetMode(previousMode)
		chainglobals.EthClient, appglobals.DB = previousClient, previousDB
	})

	h := handlers.NewHandler()
	engine, err := router.New(h, filepath.Join("..", "..", "api", "openapi.yaml"),
		[]string{"http://localhost:3000"})
	if err != nil {
		t.Fatalf("router.New: %v", err)
	}

	srv := httptest.NewServer(engine)
	t.Cleanup(h.WaitBackground)
	t.Cleanup(srv.Close)

	p := &platform{
		srv:     srv,
		db:      db,
		eth:     client,
		handler: h,
		signer:  chainutils.GetEthFrom(),
	}
	p.admin = p.promote(t, p.signUp(t, "admin"), "ADMIN")
	return p
}

// ---- HTTP -----------------------------------------------------------------

// send performs one request against the server. It takes no *testing.T so
// that it can run in a goroutine; the caller reports the error.
func (p *platform) send(method, path, token string, payload any, header http.Header) (reply, error) {
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return reply{}, fmt.Errorf("marshal %s %s: %w", method, path, err)
		}
		body = bytes.NewReader(raw)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, p.srv.URL+appglobals.BaseURL+path, body)
	if err != nil {
		return reply{}, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	for name, values := range header {
		for _, v := range values {
			req.Header.Add(name, v)
		}
	}

	resp, err := p.srv.Client().Do(req)
	if err != nil {
		return reply{}, fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return reply{}, fmt.Errorf("read %s %s: %w", method, path, err)
	}
	return reply{code: resp.StatusCode, body: raw}, nil
}

func (p *platform) call(t *testing.T, method, path, token string, payload any) reply {
	t.Helper()

	r, err := p.send(method, path, token, payload, nil)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// expect fails the test unless the answer carries the wanted status.
func (r reply) expect(t *testing.T, want int, what string) reply {
	t.Helper()

	if r.code != want {
		t.Fatalf("%s: status %d, want %d. Body: %s", what, r.code, want, r.body)
	}
	return r
}

// decode reads the body as exactly one JSON document.
//
// Exactly one: a handler that writes an error and then carries on writes a
// second document after the first, and a client reading the first one only
// would never see that it was also told the opposite.
func decode[T any](t *testing.T, r reply) T {
	t.Helper()

	var v T
	dec := json.NewDecoder(bytes.NewReader(r.body))
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("decode %T: %v. Body: %s", v, err, r.body)
	}
	var extra json.RawMessage
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		t.Fatalf("the body holds more than one JSON document: %s", r.body)
	}
	return v
}

// message is the reason an APIResponse gives.
func (r reply) message(t *testing.T) string {
	t.Helper()
	return decode[server.APIResponse](t, r).Message
}

// ---- accounts -------------------------------------------------------------

const password = "correct horse battery staple"

// signUp creates a USER account and keeps the session the sign-up opens.
// The name is only a readable prefix: every email gets a sequence number, so
// one name may be used by several journeys.
//
// Signing in again would prove nothing more and cost a second bcrypt
// comparison, more than ten seconds each under the race detector.
func (p *platform) signUp(t *testing.T, name string) account {
	t.Helper()

	email := fmt.Sprintf("%s-%d@e2e.test", strings.ToLower(strings.ReplaceAll(name, " ", ".")),
		emailSequence.Add(1))

	r := p.call(t, http.MethodPost, "/auth/sign-up", "", server.SignUpRequest{
		Email:    openapi_types.Email(email),
		Password: password,
		FullName: name,
	}).expect(t, http.StatusCreated, "sign-up "+email)

	session := decode[server.SignUpResponse](t, r).Data
	if session == nil || session.Token == "" || session.User.Id == 0 {
		t.Fatalf("sign-up %s returned no session: %s", email, r.body)
	}
	return account{id: session.User.Id, email: email, token: session.Token}
}

func (p *platform) signIn(t *testing.T, email, password string) account {
	t.Helper()

	r := p.call(t, http.MethodPost, "/auth/sign-in", "", server.SignInRequest{
		Email:    openapi_types.Email(email),
		Password: password,
	}).expect(t, http.StatusOK, "sign-in "+email)

	session := decode[server.SignInResponse](t, r).Data
	if session == nil || session.Token == "" {
		t.Fatalf("sign-in %s returned no token: %s", email, r.body)
	}
	return account{id: session.User.Id, email: email, token: session.Token}
}

// together runs the steps at the same time and returns once all of them
// have, failing the test if one did. Onboarding is dominated by bcrypt, so
// the investors of a journey are onboarded side by side; their on-chain
// registrations then compete for the platform nonce, as they would in
// production.
func together(t *testing.T, name string, steps map[string]func(t *testing.T)) {
	t.Helper()

	if !t.Run(name, func(t *testing.T) {
		for step, run := range steps {
			t.Run(step, func(t *testing.T) {
				t.Parallel()
				run(t)
			})
		}
	}) {
		t.FailNow()
	}
}

// onboardInvestors onboards several investors at once.
func (p *platform) onboardInvestors(t *testing.T, names ...string) []account {
	t.Helper()

	investors := make([]account, len(names))
	steps := make(map[string]func(t *testing.T), len(names))
	for i, name := range names {
		// A repeated name would overwrite its step and leave an investor
		// empty, to fail much later on a request made with no token.
		if _, taken := steps[name]; taken {
			t.Fatalf("investor %q is listed twice", name)
		}
		steps[name] = func(t *testing.T) { investors[i] = p.onboardInvestor(t, name) }
	}
	together(t, "onboarding", steps)
	return investors
}

// promote grants a role. No endpoint can, so it is written in SQL, the way an
// operator would; the account then signs in again because the role travels
// in the token.
func (p *platform) promote(t *testing.T, a account, role string) account {
	t.Helper()

	if _, err := p.db.Exec(`UPDATE usr.users SET role = $1 WHERE id = $2`, role, a.id); err != nil {
		t.Fatalf("promote %s to %s: %v", a.email, role, err)
	}
	promoted := p.signIn(t, a.email, password)
	promoted.wallet = a.wallet
	return promoted
}

// submitKyc files the investor's KYC and returns the id of the verification.
// The filing answers no id, so the back office's list is where it is found,
// exactly as an operator would find it.
func (p *platform) submitKyc(t *testing.T, a account, fullName string) int {
	t.Helper()

	p.call(t, http.MethodPost, "/users/kyc/verifications", a.token, server.KycVerificationRequest{
		DeclaredFullName:    fullName,
		DeclaredCountryCode: "FR",
	}).expect(t, http.StatusCreated, "KYC filing of "+a.email)

	r := p.call(t, http.MethodGet, "/users/kyc/verifications?userId="+strconv.Itoa(a.id)+"&status=submitted",
		p.admin.token, nil).expect(t, http.StatusOK, "KYC list of "+a.email)

	filed := decode[[]server.KycVerification](t, r)
	if len(filed) != 1 {
		t.Fatalf("%s must have exactly one KYC awaiting a decision, found %d: %s", a.email, len(filed), r.body)
	}
	return filed[0].Id
}

func (p *platform) approveKyc(t *testing.T, verificationID int) {
	t.Helper()

	p.call(t, http.MethodPost, fmt.Sprintf("/users/kyc/verifications/%d/approve", verificationID),
		p.admin.token, server.KycVerificationApproveRequest{
			ExpiresAt: time.Now().AddDate(1, 0, 0),
		}).expect(t, http.StatusOK, fmt.Sprintf("approval of KYC %d", verificationID))
}

func (p *platform) syncKyc(t *testing.T, verificationID int) server.KycOnChainSyncResponse {
	t.Helper()

	r := p.call(t, http.MethodPost, fmt.Sprintf("/users/kyc/verifications/%d/sync-onchain", verificationID),
		p.admin.token, nil).expect(t, http.StatusOK, fmt.Sprintf("on-chain sync of KYC %d", verificationID))
	return decode[server.KycOnChainSyncResponse](t, r)
}

// onboardInvestor takes a newcomer through the whole KYC journey: sign-up,
// filing, approval, registration on chain. The investor can then hold tokens.
func (p *platform) onboardInvestor(t *testing.T, fullName string) account {
	t.Helper()

	a := p.signUp(t, fullName)
	id := p.submitKyc(t, a, fullName)
	p.approveKyc(t, id)

	synced := p.syncKyc(t, id)
	if synced.KycStatus != "verified" {
		t.Fatalf("%s is %q after the on-chain sync, want verified", a.email, synced.KycStatus)
	}
	if !common.IsHexAddress(synced.WalletAddress) {
		t.Fatalf("%s got no wallet: %+v", a.email, synced)
	}
	a.wallet = common.HexToAddress(synced.WalletAddress)
	return a
}

// ---- catalogue ------------------------------------------------------------

// createAsset files a draft under an issuer of its own, so that no journey
// depends on the state another one left an issuer in.
func (p *platform) createAsset(t *testing.T, spec assetSpec) int {
	t.Helper()

	registration := fmt.Sprintf("B%d", emailSequence.Add(1))
	r := p.call(t, http.MethodPost, "/assets/issuers", p.admin.token, server.IssuerWriteRequest{
		Name:               spec.title + " SCI",
		LegalForm:          "SCI",
		CountryCode:        "FR",
		RegistrationNumber: &registration,
	}).expect(t, http.StatusCreated, "issuer of "+spec.title)
	issuer := decode[server.Issuer](t, r)

	payload := map[string]any{
		"title":          spec.title,
		"description":    "A listing created by the end-to-end test.",
		"estate_type_id": 1,
		"issuer_id":      issuer.Id,
		"address": map[string]any{
			"street": "12 avenue Belair", "postal_code": "06000", "city": "Nice", "country_code": "fr",
		},
		"configuration": map[string]any{
			"total_shares":              strconv.Itoa(spec.totalShares),
			"price_per_share":           spec.pricePerShare,
			"currency_code":             "eur",
			"yield":                     "4.25",
			"payment_frequency":         1,
			"payment_frequency_type_id": 1,
			"entry_fee_rate":            spec.entryFeeRate,
			"max_shares_per_investor":   strconv.Itoa(spec.maxPerInvestor),
		},
		"media": []map[string]any{{"url": "https://e2e.test/cover.jpg", "is_cover": true}},
	}

	r = p.call(t, http.MethodPost, "/assets/real-estates", p.admin.token, payload).
		expect(t, http.StatusCreated, "creation of "+spec.title)
	return decode[server.RealEstate](t, r).Id
}

// publish deploys the asset's token and returns its address.
func (p *platform) publish(t *testing.T, assetID int) common.Address {
	t.Helper()

	r := p.call(t, http.MethodPost, fmt.Sprintf("/assets/real-estates/%d/publish", assetID), p.admin.token, nil).
		expect(t, http.StatusOK, fmt.Sprintf("publication of asset %d", assetID))

	asset := decode[server.RealEstate](t, r)
	if asset.Token == nil || asset.Token.Address == nil || !common.IsHexAddress(*asset.Token.Address) {
		t.Fatalf("asset %d was published without a token: %s", assetID, r.body)
	}
	return common.HexToAddress(*asset.Token.Address)
}

func (p *platform) openFundraising(t *testing.T, assetID int) {
	t.Helper()

	p.call(t, http.MethodPost, fmt.Sprintf("/assets/real-estates/%d/open-fundraising", assetID), p.admin.token, nil).
		expect(t, http.StatusOK, fmt.Sprintf("opening of asset %d", assetID))
}

// listAsset takes an asset from nothing to open for subscription.
func (p *platform) listAsset(t *testing.T, spec assetSpec) (int, common.Address) {
	t.Helper()

	id := p.createAsset(t, spec)
	token := p.publish(t, id)
	p.openFundraising(t, id)
	return id, token
}

// asset reads the asset as a manager sees it, whatever its status.
func (p *platform) asset(t *testing.T, assetID int) server.RealEstate {
	t.Helper()

	r := p.call(t, http.MethodGet, fmt.Sprintf("/assets/real-estates/%d", assetID), p.admin.token, nil).
		expect(t, http.StatusOK, fmt.Sprintf("read of asset %d", assetID))
	return decode[server.RealEstate](t, r)
}

func (p *platform) tokensSold(t *testing.T, assetID int) int {
	t.Helper()

	asset := p.asset(t, assetID)
	if asset.Progression == nil || asset.Progression.TokensSold == nil {
		return 0
	}
	return *asset.Progression.TokensSold
}

// awaitAssetStatus polls the asset until it reaches status. The funded status
// in particular is written after the last order is CLOSED, by the same
// background job, so reading it right after CLOSED would race that write.
func (p *platform) awaitAssetStatus(t *testing.T, assetID int, status string) {
	t.Helper()

	var last string
	deadline := time.Now().Add(pollTimeout)
	for time.Now().Before(deadline) {
		if s := p.asset(t, assetID).Status; s != nil {
			last = *s
		}
		if last == status {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("asset %d is still %q after %s, want %q", assetID, last, pollTimeout, status)
}

// ---- orders ---------------------------------------------------------------

func orderPath(ref string) string {
	return "/assets/real-estate/issuance/orders/" + ref
}

// placeOrder sends the order the way the front end does, Idempotency-Key
// included. It takes no *testing.T so that investors can place orders from
// several goroutines at once.
func (p *platform) placeOrder(a account, assetID, quantity int, idempotencyKey string) (reply, error) {
	return p.send(http.MethodPost, "/assets/real-estate/issuance/orders", a.token,
		server.CreateIssuanceOrderRequest{RealEstateId: assetID, Quantity: quantity},
		http.Header{"Idempotency-Key": []string{idempotencyKey}})
}

// order places an order that must be accepted, and returns it.
func (p *platform) order(t *testing.T, a account, assetID, quantity int) server.IssuanceOrder {
	t.Helper()

	r, err := p.placeOrder(a, assetID, quantity, uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}
	r.expect(t, http.StatusCreated, fmt.Sprintf("order of %d shares of asset %d by %s", quantity, assetID, a.email))
	return createdOrder(t, r)
}

func createdOrder(t *testing.T, r reply) server.IssuanceOrder {
	t.Helper()

	created := decode[server.CreateIssuanceOrderResponse](t, r)
	if created.Data == nil {
		t.Fatalf("the order came back empty: %s", r.body)
	}
	return *created.Data
}

func (p *platform) confirmPayment(as account, orderRef string) (reply, error) {
	return p.send(http.MethodPost, orderPath(orderRef)+"/confirm-payment", as.token, nil, nil)
}

func (p *platform) fetchOrder(t *testing.T, a account, orderRef string) server.IssuanceOrder {
	t.Helper()

	r := p.call(t, http.MethodGet, orderPath(orderRef), a.token, nil).
		expect(t, http.StatusOK, "read of order "+orderRef)
	return decode[server.IssuanceOrder](t, r)
}

// awaitOrder polls the order, as its owner, until it reaches one of the
// final statuses. It then waits for the background job to return: the
// assertions that follow read the chain and the database, and must see
// everything that job wrote, not part of it.
func (p *platform) awaitOrder(t *testing.T, a account, orderRef string, final ...string) server.IssuanceOrder {
	t.Helper()

	var last server.IssuanceOrder
	deadline := time.Now().Add(pollTimeout)
	for time.Now().Before(deadline) {
		last = p.fetchOrder(t, a, orderRef)
		if slices.Contains(final, last.StatusCode) {
			p.handler.WaitBackground()
			return p.fetchOrder(t, a, orderRef)
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("order %s is still %s after %s, want one of %v. Last reason: %q",
		orderRef, last.StatusCode, pollTimeout, final, p.lastReason(t, orderRef))
	return last
}

// ---- database -------------------------------------------------------------

func (p *platform) count(t *testing.T, query string, args ...any) int {
	t.Helper()

	var n int
	if err := p.db.QueryRow(query, args...).Scan(&n); err != nil {
		t.Fatalf("%s: %v", query, err)
	}
	return n
}

// holdAsset takes the row lock on the asset, as a manager editing it would,
// and keeps it until release is called or the test ends.
//
// Every order on the asset then stalls on that row: in the lock it takes
// first, or, should it no longer take one, in the foreign key check of its
// insert, after it has read the stock. Releasing the row once they all wait
// sets them loose together, which turns a race left to the scheduler into one
// that happens every time.
func (p *platform) holdAsset(t *testing.T, assetID int) (release func()) {
	t.Helper()

	tx, err := p.db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(`SELECT 1 FROM ass.real_estate WHERE id = $1 FOR UPDATE`, assetID); err != nil {
		_ = tx.Rollback()
		t.Fatalf("lock asset %d: %v", assetID, err)
	}

	var once sync.Once
	release = func() { once.Do(func() { _ = tx.Rollback() }) }
	t.Cleanup(release)
	return release
}

// awaitLockWaiters reports whether n sessions of the database are waiting on
// a lock before the poll timeout.
func (p *platform) awaitLockWaiters(t *testing.T, n int) bool {
	t.Helper()

	deadline := time.Now().Add(pollTimeout)
	for time.Now().Before(deadline) {
		if p.count(t, `
			SELECT count(*) FROM pg_stat_activity
			WHERE datname = current_database() AND wait_event_type = 'Lock'`) >= n {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return false
}

// lastReason is the reason the latest step of the order's timeline gives.
func (p *platform) lastReason(t *testing.T, orderRef string) string {
	t.Helper()

	var reason sql.NullString
	if err := p.db.QueryRow(`
		SELECT h.reason
		FROM iss.issuance_order_status_history h
		JOIN iss.issuance_orders o ON o.id = h.order_id
		WHERE o.order_reference = $1
		ORDER BY h.created_at DESC, h.id DESC
		LIMIT 1`, orderRef).Scan(&reason); err != nil {
		t.Fatalf("timeline of %s: %v", orderRef, err)
	}
	return reason.String
}

// timeline lists the statuses the order went through, in order.
func (p *platform) timeline(t *testing.T, orderRef string) []int {
	t.Helper()

	rows, err := p.db.Query(`
		SELECT h.to_status_id
		FROM iss.issuance_order_status_history h
		JOIN iss.issuance_orders o ON o.id = h.order_id
		WHERE o.order_reference = $1
		ORDER BY h.created_at, h.id`, orderRef)
	if err != nil {
		t.Fatalf("timeline of %s: %v", orderRef, err)
	}
	defer rows.Close()

	var steps []int
	for rows.Next() {
		var s int
		if err := rows.Scan(&s); err != nil {
			t.Fatal(err)
		}
		steps = append(steps, s)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return steps
}

// ---- chain ----------------------------------------------------------------

func (p *platform) platformNonce(t *testing.T) uint64 {
	t.Helper()

	nonce, err := p.eth.PendingNonceAt(context.Background(), p.signer)
	if err != nil {
		t.Fatalf("PendingNonceAt: %v", err)
	}
	return nonce
}

func (p *platform) ledger(t *testing.T) ledger {
	t.Helper()

	return ledger{
		orders: p.count(t, `SELECT count(*) FROM iss.issuance_orders`),
		nonce:  p.platformNonce(t),
	}
}

// refused checks a refusal: the status, the reason when there is one to
// check, and that the platform neither recorded an order nor broadcast a
// transaction.
func (p *platform) refused(t *testing.T, before ledger, r reply, want int, reason error) {
	t.Helper()

	r.expect(t, want, "refusal")
	if reason != nil {
		if got := r.message(t); !strings.Contains(got, reason.Error()) {
			t.Fatalf("the refusal says %q, want it to contain %q", got, reason.Error())
		}
	}
	if after := p.ledger(t); after != before {
		t.Fatalf("a refused request changed the ledger: %+v before, %+v after", before, after)
	}
}

func (p *platform) token(t *testing.T, address common.Address) *contracts.Token {
	t.Helper()

	token, err := contracts.NewToken(address, p.eth)
	if err != nil {
		t.Fatalf("bind token %s: %v", address, err)
	}
	return token
}

func (p *platform) balanceOf(t *testing.T, token common.Address, wallet common.Address) *big.Int {
	t.Helper()

	balance, err := p.token(t, token).BalanceOf(&bind.CallOpts{}, wallet)
	if err != nil {
		t.Fatalf("balanceOf %s: %v", wallet, err)
	}
	return balance
}

func (p *platform) totalSupply(t *testing.T, token common.Address) *big.Int {
	t.Helper()

	supply, err := p.token(t, token).TotalSupply(&bind.CallOpts{})
	if err != nil {
		t.Fatalf("totalSupply of %s: %v", token, err)
	}
	return supply
}

// mints lists the Transfer events from the zero address the transaction
// emitted on token.
func (p *platform) mints(t *testing.T, token common.Address, txHash string) []*contracts.TokenTransfer {
	t.Helper()

	receipt, err := p.eth.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		t.Fatalf("receipt of %s: %v", txHash, err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Fatalf("transaction %s reverted", txHash)
	}

	binding := p.token(t, token)
	var minted []*contracts.TokenTransfer
	for _, l := range receipt.Logs {
		if l.Address != token {
			continue
		}
		if transfer, err := binding.ParseTransfer(*l); err == nil && transfer.From == (common.Address{}) {
			minted = append(minted, transfer)
		}
	}
	return minted
}

// strikeOff removes the wallet from the investor whitelist, signed by the
// platform, which is an agent of every IdentityRegistry it deploys. The
// whitelist is shared by all tokens, so the investor is struck off for all of
// them: only an investor dedicated to one journey may go through this.
func (p *platform) strikeOff(t *testing.T, token common.Address, wallet common.Address) {
	t.Helper()

	ctx := context.Background()
	irAddress, err := p.token(t, token).IdentityRegistry(&bind.CallOpts{})
	if err != nil {
		t.Fatalf("identityRegistry of %s: %v", token, err)
	}
	ir, err := contracts.NewIdentityRegistry(irAddress, p.eth)
	if err != nil {
		t.Fatalf("bind identity registry %s: %v", irAddress, err)
	}

	tx, err := chainutils.Transact(ctx, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return ir.DeleteIdentity(opts, wallet)
	})
	if err != nil {
		t.Fatalf("deleteIdentity(%s): %v", wallet, err)
	}
	if _, err := chainutils.WaitDeployedTransaction(ctx, tx, false); err != nil {
		t.Fatalf("deleteIdentity(%s) was not mined: %v", wallet, err)
	}

	verified, err := ir.IsVerified(&bind.CallOpts{}, wallet)
	if err != nil {
		t.Fatalf("isVerified(%s): %v", wallet, err)
	}
	if verified {
		t.Fatalf("%s is still verified after deleteIdentity", wallet)
	}
}

// shares converts a number of shares into token units: the tokens carry
// 18 decimals, one share being one whole token.
func shares(n int) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(n)), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}
