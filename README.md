# Tooken Services

Go REST API that tokenizes real estate as ERC-3643 (T-REX) permissioned
security tokens: an HTTP API (Gin), PostgreSQL, and Ethereum smart contracts
driven through go-ethereum.

- API reference: [`api/openapi.yaml`](api/openapi.yaml), browsable at
  `http://localhost:8080/docs/` once the API runs.
- On-chain design: [`internal/blockchain/README.md`](internal/blockchain/README.md).

## Run the whole backend locally

All you need is Docker with Compose v2:

```bash
docker compose up -d --build          # PostgreSQL, a local EVM node and the API
```

The first `up` builds the image (a minute or two). Then:

| What | Where |
| --- | --- |
| API base URL | `http://localhost:8080/api/v1` |
| Health check | `http://localhost:8080/api/v1/health` |
| API docs (Swagger UI) | `http://localhost:8080/docs/` |
| PostgreSQL | `localhost:5433`, database, user and password `tooken` |
| EVM node (HTTP and WebSocket) | `http://localhost:8555`, `ws://localhost:8555`, chain id 31337 |

All ports are published on `127.0.0.1` only.

### First steps on an empty stack

The stack starts with the schema and its reference tables only: no account,
no platform, no asset. Nothing is loaded by default.

1. Create your account with `POST /api/v1/auth/sign-up`. It answers with a
   session: send `data.token` as `Authorization: Bearer <token>`.
2. Grant yourself a back-office role. No route does it, on purpose:

   ```bash
   docker compose exec postgres psql -U tooken -d tooken \
     -c "UPDATE usr.users SET role = 'ADMIN' WHERE email = 'you@example.com'"
   ```

   The role travels in the access token: sign in again, or call
   `POST /api/v1/auth/refresh`, to get a token that carries it.
3. Deploy the T-REX platform once, through the five `/trex/init/*` routes, in
   the order given in
   [`internal/blockchain/README.md`](internal/blockchain/README.md#4-platform-bootstrap-one-time-ordered).

### Configuration

Compose reads these variables from your shell or from a `.env` file next to
`docker-compose.yml`:

| Variable | Default | Purpose |
| --- | --- | --- |
| `API_PORT` | `8080` | Host port of the API |
| `POSTGRES_PORT` | `5433` | Host port of PostgreSQL |
| `EVM_PORT` | `8555` | Host port of the EVM node |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:3000` | Comma-separated origins allowed to call the API from a browser |
| `GIN_MODE` | `release` | `debug` also prints the route table and Gin's warnings at startup |
| `LOG_TRACE` | `false` | `true` adds trace logs |
| `ORDER_RESERVATION_TTL` | `15m` | How long an unpaid order holds its shares |
| `ETH_MIN_TIP_GWEI` | `1` | Minimum priority fee of the transactions |
| `ETH_TX_WAIT` | `5m` | How long to wait for a transaction to be mined |

For example, if port 8080 is taken on your machine:

```bash
API_PORT=18080 docker compose up -d
```

The API signs every transaction with the first development account of
Hardhat and Anvil. Its private key is public knowledge: it is only acceptable
on a local chain. Likewise, `JWT_SECRET` is fixed in `docker-compose.yml` to a
public value that only suits this stack.

### Everyday commands

```bash
docker compose ps                     # state of the services
docker compose logs -f api            # API logs (background jobs report errors only here)
docker compose stop                   # stop, keeping the data
docker compose start                  # start again
docker compose up -d --build          # after a git pull: rebuild and restart the API
docker compose down -v                # wipe everything (database AND chain)
```

### What survives a restart

The database and the chain live in two Docker volumes, `tooken-dev_postgres-data`
and `tooken-dev_evm-data`. Both survive `stop`, `restart` and `down`; only `down -v` deletes
them.

They must always be kept or deleted **together**: the database records the
addresses of contracts that exist only on that chain. A database without its
chain, or the other way round, breaks every on-chain operation. To start over,
always use `docker compose down -v`.

The EVM node is Anvil, from [Foundry](https://github.com/foundry-rs/foundry). It
writes its state to its volume every 10 seconds and when it stops. A crash or
a `kill -9` of the node can therefore lose the last few seconds of the chain,
while the database keeps them. If on-chain operations start failing after such
an event, reset with `docker compose down -v`.

The database is created from [`tools/db/baseline.sql`](tools/db/baseline.sql)
the first time its volume is empty, and marked as up to date for
[golang-migrate](tools/db/migrations/README.md). When a new migration lands,
either apply it with `migrate up` or reset the stack.

## Develop the API in Go

Run PostgreSQL and the node in Docker, and the API on your machine:

```bash
docker compose up -d postgres evm
set -a; . ./.env.example; set +a      # the server does not read .env files itself
export JWT_SECRET="$(openssl rand -base64 48)"
go run .                              # serves on :8080
```

`.env.example` holds working values for this setup, except `JWT_SECRET`: it is
left empty so that no environment inherits a known one, and the server refuses
to start without a secret of at least 32 bytes. A new secret logs every user
out. To keep your settings, copy the file to `.env` (ignored by git), fill it
in and source that file instead.

The API port is fixed to 8080 in this mode: stop the API container first if
it runs (`docker compose stop api`).

## Production

Besides the database, the node and the signer key, the server needs
`JWT_SECRET`: at least 32 random bytes (`openssl rand -base64 48`), with no
leading or trailing whitespace (a trailing newline is refused), kept in the
secret store of the host, never in the repository. It signs every session:
whoever holds it can sign an ADMIN token. The server refuses to start without
it. Replacing it logs every user out.

## Build and test

```bash
go build ./...
go vet ./...
go test -race -timeout 5m ./...
```

Integration tests run against a scratch database and a Hardhat node; see
[`internal/blockchain/README.md`](internal/blockchain/README.md).

The HTTP layer is generated from `api/openapi.yaml`: after editing the spec,
run `go generate ./...`.

## Troubleshooting

**The build fails with `x509: certificate signed by unknown authority`.**
A corporate proxy intercepts TLS, and the Go image does not trust its
certificate, so `go mod download` cannot reach `proxy.golang.org`. Either add
your company's root certificate to Docker's build environment, use the
Go development mode above, or feed the build with the modules of your machine
— the repository's Dockerfile stays unchanged:

```bash
# 1. Download the modules on the host, which trusts the proxy.
export GOMODCACHE="$HOME/.cache/tooken-gomod"
go mod download

# 2. A Go image that reads them from disk instead of the network.
docker build -t tooken-golang-offline -f - "$GOMODCACHE/cache/download" <<'EOF'
FROM golang:1.25.3-alpine
COPY . /goproxy
ENV GOPROXY=file:///goproxy GOSUMDB=off GOTOOLCHAIN=local
EOF

# 3. Build the API image on top of it, then start without building.
docker buildx build --load --target api -t tooken-api:dev \
  --build-context golang:1.25.3-alpine=docker-image://tooken-golang-offline .
docker compose up -d --no-build
unset GOMODCACHE
```

**`docker compose up` fails with `port is already allocated`.** Another
process uses one of the ports: set `API_PORT`, `POSTGRES_PORT` or `EVM_PORT`.

**A platform deployment never finishes.** Steps 1, 4 and 5 answer `202` and
run in the background: their errors are only in `docker compose logs api`. A
half-deployed platform cannot be resumed: `docker compose down -v`, then
`up -d` and deploy again.

**`api` is `unhealthy` or keeps restarting.** Read `docker compose logs api`:
the API exits when it cannot reach the database or the node at startup.
