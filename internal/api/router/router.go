// Package router builds the HTTP engine the server runs.
//
// It lives outside main so that the binary and the end-to-end tests serve the
// very same engine: same CORS policy, same authentication middleware, same
// routes, same documentation. A test that rebuilt its own router would keep
// passing while the real one drifted.
package router

import (
	"fmt"
	"time"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	"github.com/TookenOrg/tooken-services/internal/globals"
	"github.com/TookenOrg/tooken-services/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swgui "github.com/swaggest/swgui/v5"
)

// New returns the engine serving h under globals.BaseURL.
//
// specPath is the OpenAPI document: the authentication middleware reads from
// it which operations require a bearer token, and it is served as is at
// /openapi.yaml for the documentation UI at /docs.
//
// corsOrigins lists the origins a browser may call the API from. It cannot be
// empty, and each origin needs its scheme: the CORS middleware refuses such a
// configuration, and it does so by panicking, so it is checked here first and
// reported as an error.
func New(h server.ServerInterface, specPath string, corsOrigins []string) (*gin.Engine, error) {
	corsConfig := cors.Config{
		AllowOrigins: corsOrigins,
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	if err := corsConfig.Validate(); err != nil {
		return nil, fmt.Errorf("CORS configuration (CORS_ALLOWED_ORIGINS): %w", err)
	}

	if err := middleware.InitAuth(specPath); err != nil {
		return nil, fmt.Errorf("authentication initialization: %w", err)
	}

	router := gin.Default()
	router.Use(cors.New(corsConfig))

	apiV1 := router.Group(globals.BaseURL)
	apiV1.Use(middleware.AutoAuthMiddleware())
	server.RegisterHandlers(apiV1, h)

	router.StaticFile("/openapi.yaml", specPath)
	ui := swgui.NewHandler("Tooken API Docs", "/openapi.yaml", "/docs")
	router.GET("/docs/*any", gin.WrapH(ui))

	return router, nil
}
