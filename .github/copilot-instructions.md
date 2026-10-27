# Copilot instructions for Tooken-Services

Go REST API backend for tokenizing real-world assets (real estate) as ERC-3643 /
T-REX permissioned security tokens. It exposes an HTTP API (Gin), persists to
PostgreSQL (raw SQL), and drives Ethereum smart contracts via go-ethereum.

## Build, test, lint

- Build everything: `go build ./...`
- Run the server: `go run .` (entry point is `main.go` at the repo root; listens on `:8080`)
- Run all tests with the race detector (matches CI): `go test -race -timeout 5m ./...`
- Run a single test: `go test -run '^TestName$' ./internal/path/to/pkg`
- Vet: `go vet ./...`; format with `gofmt -w` before committing.
- **Unit tests** live next to the code (`internal/blockchain/{utils,services}/*_test.go`)
  and run in CI via `go test ./...` (no chain/DB needed) — e.g. `ConvertFloatToWei`,
  `buildTokenDetails`, the ONCHAINID claim signature/key derivation.
- **Integration tests** are gated by the `integration` build tag (skipped by default
  CI) and run against a local EVM. A Hardhat harness lives in `tools/hardhat/`:
  `cd tools/hardhat && npm install && npx hardhat node`, then in another terminal
  `go test -tags integration ./internal/blockchain/...`. Override with `ETH_TEST_RPC`
  / `PRIVATE_KEY`. See `internal/blockchain/services/integration_test.go` (deploy +
  tx helpers) and `e2e_test.go` (full T-REX lifecycle: deploy token → register a KYC
  identity → mint succeeds; mint to a non-KYC wallet reverts). Once
  `TEST_DATABASE_URL` or `ETH_TEST_RPC` is set, an unreachable database or node
  **fails** the chain/delivery integration tests instead of skipping them; run
  the full suite with `-p 1` (packages share the test database and the signer's
  nonce).
- Integration helpers live in **`internal/testenv`** (test-only):
  `ScratchDatabase(t, "tooken_it_<pkg>")` builds a private database from
  `tools/db/baseline.sql`, `DialNode(ctx, t)` dials the node. Use them instead of a
  new harness; `go list -deps . | grep internal/testenv` must stay empty.
- **End-to-end** tests (`internal/e2e`) drive the real router (`router.New`)
  over HTTP against the scratch database and Hardhat:
  `TEST_DATABASE_URL=… go test -race -tags integration -count=1 -timeout 20m ./internal/e2e/`
  (about 2 min). Background work started by a handler goes through
  `h.startBackground`; tests call `Handler.WaitBackground()` before reading its
  effects.
- `golangci-lint` is referenced in CI but **commented out / disabled** until the
  toolchain catches up to Go 1.25.3 — do not assume it runs.

## OpenAPI-first workflow (important)

`api/openapi.yaml` is the source of truth for the HTTP API. The Gin
`ServerInterface`, request/response models, and route registration are generated
into `internal/api/server/server.gen.go`.

- Regenerate after editing the spec: `go generate ./...` (runs `oapi-codegen`
  with `api/oapi-codegen.yaml`; the directive lives in `main.go`).
- **Never hand-edit `internal/api/server/server.gen.go`** (`DO NOT EDIT`).
- Adding/changing an endpoint = edit `openapi.yaml` → regenerate → implement the
  method on `Handler` in `internal/api/handlers/`. `Handler` must satisfy
  `server.ServerInterface`; a missing method breaks the build.

## Architecture & layering

Per-domain packages under `internal/` (`auth`, `assets_managements`,
`blockchain`, `orders`, `payments`, `users`) each split into the same layers:

`handlers` (in `internal/api/handlers`) → `services` → `database` (and, for
blockchain, on-chain `services`/`contracts`).

- Handlers are thin: bind/validate the request, call a service, map the result to
  a generated `server.*Response`. They live in `internal/api/handlers` and hang
  off a single `Handler` struct that holds one service per domain
  (`handlers.go` → `NewHandler()`). Routes, CORS and the spec-driven auth are
  wired in one place, `internal/api/router` → `router.New`, which `main.go`
  and the e2e tests both call.
- Services are stateless structs created via `NewService()`; business logic and
  orchestration live here.
- `database` packages hold raw SQL using the global `globals.DB` (`*sql.DB`).
  Queries reference Postgres schema prefixes: `usr.`, `blk.`, `iss.`, `ass.`,
  `rel.`. Row structs are named `...DTO`.

## Auth (route protection is driven by the spec)

- JWT, HS256, via `internal/auth/utils/jwt.go`. ⚠️ `JwtKey` is hardcoded
  (`"secret_key_to_change"`) — treat as a placeholder, never a real secret.
- `router.New` calls `middleware.InitAuth("./api/openapi.yaml")`, which parses the spec and
  marks every operation carrying `security: - bearerAuth: []` as protected.
  **To require auth on a route, add that `security` block in `openapi.yaml`** —
  there is no per-handler auth code.
- `AutoAuthMiddleware` validates the `Bearer` token and stores `*CustomClaims`
  in the Gin context under `"user_claims"`. Read it with
  `middleware.GetUserClaims(gCtx)` (returns claims + ok).

## Blockchain conventions

- Contract bindings in `internal/blockchain/contracts/bindings/` are abigen-
  generated — do not edit by hand.
- Send every transaction through `utils.Transact(ctx, func(opts) (*types.Transaction, error))`
  (signs with the `PRIVATE_KEY` env key). It holds one lock for the key from
  nonce read to broadcast and bounds that step with a timeout; calling
  `GenerateTransactOpts` directly races other transactions for the nonce.
  The callback must broadcast one transaction and never call `Transact` itself.
  Wait for mining afterwards, outside it, with `utils.WaitDeployedTransaction` /
  `WaitTREXSuiteDeployment`.
- Long-running deploy/mint/burn operations run in the background and return
  `202 Accepted` immediately. Detach the context first
  (`context.WithoutCancel(reqCtx)` or `context.Background()`) so the work
  survives the HTTP request, and start it with `h.startBackground(what, job)`,
  never a bare `go`: on SIGINT/SIGTERM `main.go` `serve` stops the HTTP server,
  then `Handler.CloseAndWait()` refuses new jobs and waits for the running ones,
  so a delivery is not cut between the mint and the order update. A shutdown
  can take up to 15 s (requests) + 2 min (jobs): whatever stops the process
  must wait that long before killing it (Docker's default is 10 s; set
  `stop_grace_period` / `terminationGracePeriodSeconds` accordingly).
- Shared on-chain state (eth client, deployed contract addresses/instances) lives
  in `internal/blockchain/globals` and `internal/blockchain/services/globals.go`.
- **Token creation goes through the shared `TREXFactory`** (`services/token.go`
  `CreateToken` → `deployTREXSuite`), not hand-wired proxies. Target architecture:
  ONE shared factory, claim issuer, manager (the `PRIVATE_KEY` signer) and **IRS**
  (the shared investor whitelist); per-token IR/TIR/CTR/MC. Passing the shared IRS
  reuses one KYC across all tokens — the factory keeps IRS ownership, so reusing it
  across `deployTREXSuite` calls works. Build `TokenDetails` with `ONCHAINID = 0` so
  the factory creates the token identity (never an EOA).
- **ONCHAINID / ERC-734 keys** are `keccak256(abi.encode(address))` ==
  `keccak256(LeftPad32(address))`, **not** `keccak256(address[20])`. Applies to the
  claim-signer key (purpose 3) and management key (purpose 1).
- A token mint/transfer only succeeds if the recipient is
  `IdentityRegistry.isVerified` (registered **and** holding the required claims signed
  by a trusted issuer). Minting to a non-verified wallet reverts with
  `Identity is not verified.`.
- Singleton addresses (factory, claim issuer, module, shared IRS) are persisted to
  `blk.contract_role` via `database.InsertContractRole` and read with
  `GetContractRoleByName`; created tokens via `InsertToken`. Deep module docs +
  diagrams: **`internal/blockchain/README.md`** (section 9 covers the tests).

## Other conventions

- Logging goes through `pkg/logger` (slog wrapper): `logger.LogInfo/LogDebug/
  LogWarn`, and `logger.LogError`, which **both logs and returns an `error`**
  (handy for `return logger.LogError(...)`). Emoji-prefixed log messages
  (🚀, ✅, ⏳, 🆗) are the house style. `LOG_TRACE=true` enables `LogTrace`.
- Go named return values are used heavily; some service methods return an HTTP
  status code alongside the error (e.g. `SignUp` returns `(errCode int, ...)`).
- Globals: `globals.DB` and `globals.BaseURL` (`/api/v1`) in `internal/globals`.

## Environment

Config is read straight from `os.Getenv` — **there is no dotenv loader in the
code**, so variables must be exported in the process environment. The committed
VS Code `launch.json` injects `.env` (gitignored) only when debugging.

Required/used vars: `DATABASE_URL`, `WS_RPC_URL` (Ethereum WS RPC), `PRIVATE_KEY`
(signer), `CORS_ALLOWED_ORIGINS` (comma-separated), `LOG_TRACE`,
`ORDER_RESERVATION_TTL` (Go duration, default `15m`: how long an unpaid order
holds its shares).

## Git / PR conventions

- Branches: `feature/**`, `fix/**`, `hotfix/**` (these trigger CI alongside
  `main`/`develop`).
- PR titles follow Conventional Commits (`feat:`, `fix:`, `docs:`, …); target
  `develop` for integration and `main` for releases (see
  `.github/PULL_REQUEST_TEMPLATE.md`).

## Known gotcha

The `Dockerfile` builds `./cmd/server`, but the entry point is `main.go` at the
repo root (there is no `cmd/` directory), so the Docker build is currently stale.
Build/run against `main.go` until that is reconciled.
