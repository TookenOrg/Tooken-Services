package main

import (
	"database/sql"
	"os"

	"fmt"
	"log"
	"net/url"

	_ "github.com/lib/pq"

	"github.com/TookenOrg/tooken-services/internal/api/handlers"
	"github.com/TookenOrg/tooken-services/internal/api/server"
	blkGlobals "github.com/TookenOrg/tooken-services/internal/blockchain/globals"
	"github.com/TookenOrg/tooken-services/internal/blockchain/services"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	swgui "github.com/swaggest/swgui/v5"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config ./api/oapi-codegen.yaml ./api/openapi.yaml

func main() {

	var err error

	// init logs
	logger.Init(true)

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
	return db, nil
}
