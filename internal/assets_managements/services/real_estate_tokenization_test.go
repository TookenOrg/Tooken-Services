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

	// The salt is what makes the deployment replayable and collision-free, so it
	// is checked apart from the rest of the payload.
	if want := "re-12"; deployer.gotSalt != want {
		t.Errorf("salt = %q; want %q", deployer.gotSalt, want)
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
