package main

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"fmt"
	"log"
	"net/http"
	"net/url"

	_ "github.com/lib/pq"

	"github.com/TookenOrg/tooken-services/internal/api/handlers"
	"github.com/TookenOrg/tooken-services/internal/api/router"
	authUtils "github.com/TookenOrg/tooken-services/internal/auth/utils"
	blkGlobals "github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/services"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config ./api/oapi-codegen.yaml ./api/openapi.yaml

func main() {

	var err error

	// init logs
	logger.Init(true)

	// Without its own secret the server would trust tokens anyone can sign:
	// refuse to start rather than fall back on a default.
	if err := authUtils.SetSecret(os.Getenv("JWT_SECRET")); err != nil {
		log.Fatalf("Error on JWT_SECRET: %v (generate one with: openssl rand -base64 48)", err)
	}

	// init db
	globals.DB, err = setupDatabase()
	if err != nil {
		log.Fatalf("Error on SetupDatabase: %s", err.Error())
	}
	defer globals.DB.Close()

	// init redis

	// init rabbitMQ

	// Initialize blockchain
	blkGlobals.EthClient, err = services.SetupEthClient()
	if err != nil {
		log.Fatalf("Error on Setup blockchain: %s", err.Error())
	}

	startServer()
}

func startServer() {

	h := handlers.NewHandler()
	engine, err := router.New(h, "./api/openapi.yaml", getCorsOrigins())
	if err != nil {
		log.Fatalf("Error during router initialization: %v", err)
	}

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// The first SIGINT or SIGTERM starts a graceful shutdown. Once it has
	// begun, the default handling is restored: a second signal kills the
	// process, for whoever does not want to wait.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	context.AfterFunc(ctx, stop)

	logger.LogInfo("Starting server on %s", ln.Addr())
	if err := serve(ctx, &http.Server{Handler: engine}, ln, h, defaultStopTimeouts); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}

// stopTimeouts bounds the two waits of a shutdown.
type stopTimeouts struct {
	// requests is how long the requests in flight get to finish.
	requests time.Duration
	// background is how long the jobs they started get to finish. A delivery
	// may wait up to a minute for the node to take its transaction, then for
	// the transaction to be mined.
	background time.Duration
}

var defaultStopTimeouts = stopTimeouts{
	requests:   15 * time.Second,
	background: 2 * time.Minute,
}

// backgroundJobs is the part of the handler a shutdown waits on.
type backgroundJobs interface {
	CloseAndWait()
}

// serve answers on ln until ctx is cancelled, then shuts down in order:
//
//  1. it stops accepting connections and lets the requests in flight finish;
//  2. it closes the connections of those that do not finish in time, so their
//     clients learn the outcome is unknown instead of reading an answer;
//  3. it refuses new background jobs and waits for the running ones.
//
// The last step is the point of it all. Those jobs are detached from the
// request that started them and keep writing to the chain and the database:
// exiting under a delivery would cut it between the mint and the update of
// the order. That is why a listener that fails, possibly after hours of
// service, goes through the same steps before its error is returned.
func serve(ctx context.Context, srv *http.Server, ln net.Listener, jobs backgroundJobs, timeouts stopTimeouts) error {
	served := make(chan error, 1)
	go func() { served <- srv.Serve(ln) }()

	var failure error
	select {
	case failure = <-served:
		logger.LogWarn("⚠️ The server can no longer accept connections: %v", failure)
	case <-ctx.Done():
	}

	logger.LogInfo("🛑 Shutting down: waiting up to %s for the requests in flight", timeouts.requests)
	stopCtx, cancel := context.WithTimeout(context.Background(), timeouts.requests)
	defer cancel()
	if err := srv.Shutdown(stopCtx); err != nil {
		logger.LogWarn("⚠️ Requests still running after %s, closing their connections: %v", timeouts.requests, err)
		_ = srv.Close()
	}
	if failure == nil {
		<-served
	}

	return errors.Join(failure, waitForJobs(jobs, timeouts.background))
}

// waitForJobs refuses new background jobs and waits at most timeout for the
// running ones.
func waitForJobs(jobs backgroundJobs, timeout time.Duration) error {
	logger.LogInfo("⏳ Waiting up to %s for the background jobs", timeout)
	done := make(chan struct{})
	go func() {
		jobs.CloseAndWait()
		close(done)
	}()

	select {
	case <-done:
		logger.LogInfo("✅ Background jobs finished, server stopped")
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("background jobs still running after %s, exiting under them", timeout)
	}
}

func getCorsOrigins() []string {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		return []string{}
	}
	return strings.Split(origins, ",")
}

func setupDatabase() (dbClient *sql.DB, err error) {

	serviceURI := os.Getenv("DATABASE_URL")
	conn, _ := url.Parse(serviceURI)
	db, err := sql.Open("postgres", conn.String())
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT version()")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var result string
		err = rows.Scan(&result)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Version: %s\n", result)
	}

	logger.LogInfo("✅ Connected to PostgreSQL!")

	var (
		dbName string
		user   string
		addr   string
	)
	db.QueryRow(`SELECT current_database(), current_user, inet_server_addr()`).Scan(&dbName, &user, &addr)

	logger.LogInfo("🗄️ DB=%s | USER=%s | HOST=%s", dbName, user, addr)

	return db, nil
}
