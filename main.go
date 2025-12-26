package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/TookenOrg/tooken-services/internal/api/handlers"
	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	swgui "github.com/swaggest/swgui/v5"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config ./api/oapi-codegen.yaml ./api/openapi.yaml

func main() {

	// var err error

	// init logs
	logger.Init(true)

	// init db
	// globals.DB, err = setupDatabase()
	// if err != nil {
	// 	log.Fatalf("Error on SetupDatabase: %s", err.Error())
	// }

	// init redis

	// init rabbitMQ

	// Initialize blockchain
	// globals.EthClient, err = services.SetupEthClient()
	// if err != nil {
	// 	log.Fatalf("Error on Setup blockchain: %s", err.Error())
	// }

	// services.SetGlobals()

	startServer()

}

func startServer() {

	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	handler := handlers.NewHandler()
	server.RegisterHandlers(apiV1, handler)

	// --- Swagger UI ---
	router.StaticFile("/openapi.yaml", "./api/openapi.yaml")
	ui := swgui.NewHandler("Tooken API Docs", "/openapi.yaml", "/docs")
	router.GET("/docs/*any", gin.WrapH(ui))

	logger.LogInfo("Starting server on %s", ":8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupDatabase() (dbClient *sql.DB, err error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return
	}

	logger.LogInfo("✅ Connected to PostgreSQL!")
	return db, nil
}
