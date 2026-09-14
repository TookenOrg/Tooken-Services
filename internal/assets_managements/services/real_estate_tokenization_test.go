package services

import (
	"context"
	"errors"
	"testing"
	"unicode/utf8"

	"github.com/TookenOrg/tooken-services/internal/api/server"
)

func TestDefineTokenName(t *testing.T) {
	tests := []struct {
		reName   string
		reID     int
		expected string
	}{
		{"Sample Real Estate", 123, "Tooken Sample Real Estate #123"},
		{"Another Real Estate", 456, "Tooken Another Real Estate #456"},
		{"Another Real Estate with a very long name and accents like this ééççç)ààaa and éééé more text to test truncation", 456, "Tooken Another Real Estate with a very long name and accents like this ééççç)ààaa and é #456"},
		{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaéééééééééééééééééééé", 456, "Tooken aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaééé #456"},
		{"👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍❤👍👍", 456, "Tooken 👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍👍❤ #456"},
		{"   ", 456, "Tooken Unnamed #456"},
	}
	for _, tt := range tests {
		result := defineTokenName(tt.reName, tt.reID)
		if !utf8.ValidString(result) {
			t.Errorf("defineTokenName(%q, %d) produced an invalid UTF-8 string: %q", tt.reName, tt.reID, result)
		}
		if result != tt.expected {
			t.Errorf("defineTokenName(%q, %d) = %q; want %q", tt.reName, tt.reID, result, tt.expected)
		}
	}
}

func TestDefineSymbol(t *testing.T) {
	realEstateID := 123
	expected := "TKN123"
	result := defineSymbol(realEstateID)
	if result != expected {
		t.Errorf("defineSymbol(%d) = %q; want %q", realEstateID, result, expected)
	}
}

func TestDefineSalt(t *testing.T) {
	realEstateID := 123
	expected := "re-123"
	result := defineSalt(realEstateID)
	if result != expected {
		t.Errorf("defineSalt(%d) = %q; want %q", realEstateID, result, expected)
	}
}

// fakeTokenDeployer stands in for the blockchain module. It records what it was
// asked to deploy so a test can assert on it, and it never touches a chain: the
// whole point of declaring TokenDeployer on this side is that the publication
// rule can be exercised without an RPC node, a private key or a factory.
type fakeTokenDeployer struct {
	calls   int
	gotReq  server.CreateTokenRequest
	gotSalt string
	tokenID int
	err     error

	// The lookup is recorded apart from the deployment on purpose. A single
	// shared field would be written twice on the deploy path, the second call
	// silently overwriting the first — and code that searched under one salt
	// while deploying under another would then look perfectly healthy.
	bySalt        *server.TokenInfos
	saltErr       error
	saltCalls     int
	gotLookupSalt string
}

var _ TokenDeployer = (*fakeTokenDeployer)(nil)

func (f *fakeTokenDeployer) CreateTokenWithSalt(_ context.Context, req server.CreateTokenRequest, salt string) (int, server.TokenInfos, error) {
	f.calls++
	f.gotReq = req
	f.gotSalt = salt

	if f.err != nil {
		return 0, server.TokenInfos{}, f.err
	}

	return f.tokenID, server.TokenInfos{}, nil
}

// GetTokenBySalt answers "nothing deployed under that salt" unless a test says
// otherwise. That is the ordinary case — a first publication — so it has to be
// the zero value; a fake that always found something would hide the deployment
// path entirely.
func (f *fakeTokenDeployer) GetTokenBySalt(_ context.Context, salt string, _ bool) (*server.TokenInfos, error) {
	f.saltCalls++
	f.gotLookupSalt = salt

	if f.saltErr != nil {
		return nil, f.saltErr
	}

	return f.bySalt, nil
}

// An asset that already carries a token must never reach the chain again: the
// salt is spent, and the binding investors hold is the one that counts. A call
// count of zero is the only thing that actually proves it.
func TestTokenForPublicationKeepsAnExistingToken(t *testing.T) {
	deployer := &fakeTokenDeployer{tokenID: 99}
	s := NewService(deployer)
	existing := 7

	got, err := s.tokenForPublication(context.Background(), 12, "Villa Belair", &existing)
	if err != nil {
		t.Fatalf("tokenForPublication returned an unexpected error: %v", err)
	}
	if got != existing {
		t.Errorf("tokenForPublication = %d; want the token the asset already carries (%d)", got, existing)
	}
	if deployer.calls != 0 {
		t.Errorf("deployer called %d time(s); an already tokenised asset must not be deployed again", deployer.calls)
	}
	// Resolving the salt would be a database round-trip spent to learn what the
	// asset already states. The shortcut is only real if nothing is asked.
	if deployer.saltCalls != 0 {
		t.Errorf("lookup called %d time(s); an asset that already carries a token has nothing to resolve", deployer.saltCalls)
	}
}

func TestTokenForPublicationDeploysWhenTheAssetHasNoToken(t *testing.T) {
	deployer := &fakeTokenDeployer{tokenID: 42}
	s := NewService(deployer)

	got, err := s.tokenForPublication(context.Background(), 12, "Villa Belair", nil)
	if err != nil {
		t.Fatalf("tokenForPublication returned an unexpected error: %v", err)
	}
	if got != deployer.tokenID {
		t.Errorf("tokenForPublication = %d; want the freshly deployed token (%d)", got, deployer.tokenID)
	}
	if deployer.calls != 1 {
		t.Fatalf("deployer called %d time(s); want exactly 1", deployer.calls)
	}
	// Deploying without resolving first is what leaves an interrupted attempt
	// unclaimed, so the lookup has to have happened — once.
	if deployer.saltCalls != 1 {
		t.Fatalf("lookup called %d time(s); want exactly 1 before deploying", deployer.saltCalls)
	}

	// The salt is what makes the deployment replayable and collision-free, so it
	// is checked apart from the rest of the payload.
	if want := "re-12"; deployer.gotSalt != want {
		t.Errorf("salt = %q; want %q", deployer.gotSalt, want)
	}
	// The whole recovery rests on these two being the same string. A suite
	// deployed under one salt and searched for under another could never be
	// adopted, and no assertion on a hardcoded value would reveal it.
	if deployer.gotLookupSalt != deployer.gotSalt {
		t.Errorf("looked up %q but deployed under %q; a replay would never find the orphan",
			deployer.gotLookupSalt, deployer.gotSalt)
	}
	if want := "Tooken Villa Belair #12"; deployer.gotReq.TokenName != want {
		t.Errorf("token name = %q; want %q", deployer.gotReq.TokenName, want)
	}
	if want := "TKN12"; deployer.gotReq.Symbol != want {
		t.Errorf("symbol = %q; want %q", deployer.gotReq.Symbol, want)
	}
	// Zero decimals would forbid ever selling a fraction of a share, and the
	// choice is frozen at deployment: the platform imposes it rather than
	// leaving it to whoever fills a form.
	if deployer.gotReq.NbDecimal != tokenDecimals {
		t.Errorf("decimals = %d; want %d", deployer.gotReq.NbDecimal, tokenDecimals)
	}
}

// A chain failure must leave nothing behind for the caller to bind, so that the
// asset stays a draft rather than becoming an offer backed by nothing.
func TestTokenForPublicationPropagatesADeploymentFailure(t *testing.T) {
	wantErr := errors.New("chain is unreachable")
	deployer := &fakeTokenDeployer{tokenID: 42, err: wantErr}
	s := NewService(deployer)

	got, err := s.tokenForPublication(context.Background(), 12, "Villa Belair", nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("tokenForPublication error = %v; want %v", err, wantErr)
	}
	if got != 0 {
		t.Errorf("tokenForPublication = %d on failure; want 0 so that no asset can be bound", got)
	}
}

// An attempt interrupted after the deployment leaves a suite on-chain, and a
// blk.token row, that nothing points at. Recognising it on the next try is what
// turns that failure into something a replay repairs, instead of an asset that
// can never be published again.
func TestTokenForPublicationAdoptsAnInterruptedDeployment(t *testing.T) {
	const orphan = 7
	deployer := &fakeTokenDeployer{
		tokenID: 42,
		bySalt:  &server.TokenInfos{Id: orphan},
	}
	s := NewService(deployer)

	got, err := s.tokenForPublication(context.Background(), 12, "Villa Belair", nil)
	if err != nil {
		t.Fatalf("tokenForPublication returned an unexpected error: %v", err)
	}
	if got != orphan {
		t.Errorf("tokenForPublication = %d; want the token left behind by the interrupted attempt (%d)", got, orphan)
	}
	// The suite is already on-chain and its salt is spent: a second deployment
	// would be refused by the factory, and the call count is the only thing that
	// proves none was attempted.
	if deployer.calls != 0 {
		t.Fatalf("deployer called %d time(s); the suite is already deployed", deployer.calls)
	}
	if want := "re-12"; deployer.gotLookupSalt != want {
		t.Errorf("looked up %q; want %q, derived from the asset identifier", deployer.gotLookupSalt, want)
	}
}

// The salt comes from the asset identifier, the token name from its title — and
// the title can be edited between the failed attempt and the replay. Adoption has
// to survive that, which is exactly what a lookup by name would not: it is the
// reason the salt, and not the name, is the key.
func TestTokenForPublicationAdoptsEvenAfterTheTitleChanged(t *testing.T) {
	const orphan = 7
	deployer := &fakeTokenDeployer{
		tokenID: 42,
		bySalt:  &server.TokenInfos{Id: orphan, TokenName: "Tooken Villa Belair #12"},
	}
	s := NewService(deployer)

	got, err := s.tokenForPublication(context.Background(), 12, "Renamed Since The Failure", nil)
	if err != nil {
		t.Fatalf("tokenForPublication returned an unexpected error: %v", err)
	}
	if got != orphan {
		t.Errorf("tokenForPublication = %d; want %d — a changed title must not hide the deployment", got, orphan)
	}
	if deployer.calls != 0 {
		t.Fatalf("deployer called %d time(s); renaming an asset does not entitle it to a second suite", deployer.calls)
	}
	// The salt is unchanged because the identifier is, which is the whole point.
	if want := "re-12"; deployer.gotLookupSalt != want {
		t.Errorf("looked up %q; want %q despite the new title", deployer.gotLookupSalt, want)
	}
}

// A database that cannot answer is not the same as a database answering "nothing
// was deployed". Treating the two alike is how a second suite gets created for an
// asset that already has one — the very orphan this lookup exists to prevent.
func TestTokenForPublicationDoesNotDeployWhenTheLookupFails(t *testing.T) {
	wantErr := errors.New("database is unreachable")
	deployer := &fakeTokenDeployer{tokenID: 42, saltErr: wantErr}
	s := NewService(deployer)

	got, err := s.tokenForPublication(context.Background(), 12, "Villa Belair", nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("tokenForPublication error = %v; want %v", err, wantErr)
	}
	if got != 0 {
		t.Errorf("tokenForPublication = %d on failure; want 0 so that no asset can be bound", got)
	}
	if deployer.calls != 0 {
		t.Fatalf("deployer called %d time(s); a failed lookup must not be read as 'nothing deployed'", deployer.calls)
	}
}
