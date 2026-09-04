package services

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestGenerateSignatureAddClaim validates the ONCHAINID claim-signing path end to
// end (without a chain): the signature must recover to the signer derived from
// PRIVATE_KEY, and the ERC-734 key must be keccak256(abi.encode(address)) ==
// keccak256(left-pad-32(address)) — the conformance fix.
func TestGenerateSignatureAddClaim(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	t.Setenv("PRIVATE_KEY", common.Bytes2Hex(crypto.FromECDSA(key)))

	identity := common.HexToAddress("0xAbC0000000000000000000000000000000000123")

	res, err := generateSignatureAddClaim(identity, 7)
	if err != nil {
		t.Fatalf("generateSignatureAddClaim: %v", err)
	}
	if res == nil {
		t.Fatal("expected a signature result")
	}

	wantAddr := crypto.PubkeyToAddress(key.PublicKey)
	if res.RecoveredAddr != wantAddr {
		t.Errorf("recovered %s, want %s", res.RecoveredAddr.Hex(), wantAddr.Hex())
	}

	wantKey := crypto.Keccak256Hash(common.LeftPadBytes(wantAddr.Bytes(), 32))
	if res.Key != wantKey {
		t.Errorf("ONCHAINID key derivation mismatch: got %s, want %s", res.Key.Hex(), wantKey.Hex())
	}

	if len(res.SignatureBytes) != 65 {
		t.Errorf("signature length: got %d, want 65", len(res.SignatureBytes))
	}
	if res.V != 27 && res.V != 28 {
		t.Errorf("recovery id V should be 27 or 28, got %d", res.V)
	}
}
