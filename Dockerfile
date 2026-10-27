# Builds the API.
#
#   docker build -t tooken-api .

FROM golang:1.25.3-alpine AS build
WORKDIR /src
# Modules first: this layer is reused as long as go.mod and go.sum are unchanged.
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0

FROM build AS build-api
RUN go build -trimpath -ldflags="-s -w" -o /out/tooken .

# The base image already ships the CA bundle that the TLS clients need.
FROM alpine:3.22 AS runtime
RUN adduser -D -H -u 10001 tooken
WORKDIR /app
USER tooken

FROM runtime AS api
# The server reads its OpenAPI spec at ./api/openapi.yaml, relative to the
# working directory: route protection is derived from it.
COPY api/openapi.yaml /app/api/openapi.yaml
COPY --from=build-api /out/tooken /app/tooken
EXPOSE 8080
# Exec form: the server must receive SIGTERM itself to shut down gracefully.
CMD ["/app/tooken"]
