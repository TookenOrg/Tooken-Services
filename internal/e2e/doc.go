// Package e2e proves that the platform works as a whole: a real HTTP server,
// a real PostgreSQL database and a real EVM node, driven only through the
// public API, the way a front end and an operator would drive them.
//
// It holds no production code. Its tests carry the integration build tag, so
// a plain `go test ./...` compiles nothing here. To run them:
//
//	cd tools/hardhat && npx hardhat node          # in one terminal
//	TEST_DATABASE_URL="postgres://…" \
//	  go test -tags integration -count=1 ./internal/e2e/...
//
// The test builds its own database, tooken_it_e2e, from tools/db/baseline.sql
// and drops it at the end. It signs with the first Hardhat account, like every
// other chain test: run the packages one at a time (-p 1) when several of them
// share the node, or they race for the nonce of that account.
//
// The platform is initialised once, through the same routes an operator
// calls. Each journey then lists its own asset and onboards its own investors,
// so that a failure in one of them neither hides nor causes a failure in
// another. To play a single journey, keep the initialisation in the -run
// pattern, or the journey starts on a platform that was never set up:
//
//	go test -tags integration -count=1 -run 'TestEndToEnd/(initialises|race)' ./internal/e2e/...
//
// Some failures cannot be provoked through the API, because the API exists to
// prevent them: an investor removed from the whitelist on chain while the
// database still says verified, or a per-investor limit raised in the database
// above the one enforced by the token. Those are injected directly, in SQL or
// on chain, and each injection says what it simulates.
package e2e
