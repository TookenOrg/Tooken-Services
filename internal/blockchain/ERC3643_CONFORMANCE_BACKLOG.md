# ERC-3643 conformance backlog

All conformance deviations found during the ERC-3643 / T-REX review of
`internal/blockchain` have been addressed:

- Blocking deviations — **fixed**:
  - claim-signer key derived as `keccak256(abi.encode(address))` (ERC-734),
  - factory `TokenDetails.ONCHAINID` set to the zero address,
  - manual `TokenProxy` `_onchainID` set to the zero address.
- Quick wins — **applied**: identity-authority name literal, `decimals <= 18`
  validation, `AddManagementKeyToClaimIssuer` → `AddClaimSignerKeyToClaimIssuer`.
- Architecture refactor #1 — **done** (commit `c219bee`): `CreateToken` now goes
  through the shared `TREXFactory` and reuses a shared IRS / claim issuer / manager.

**No conformance items remain open.**

For related, non-conformance gaps (config tables only partially auto-populated,
mocked `GetWalletByUserId`, plaintext wallet keys, stubs, `SetGlobals` disabled),
see the "Known gaps & caveats" section of `internal/blockchain/README.md`.
