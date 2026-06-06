// go:generate runs the docs generator that writes docs/openapi.json and
// docs/API_CONTRACTS.md from the single source of truth in internal/apidocs.
// Run from the backend directory:
//
//	go generate ./...
//
//go:generate go run ../docs

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"prd-generator/internal/apidocs"
	"prd-generator/internal/config"
	"prd-generator/internal/db"
	"prd-generator/internal/handlers"
	"prd-generator/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Printf("Warning: Database connection failed: %v", err)
	}

	// Initialize database schema
	if err := db.InitSchema(conn); err != nil {
		log.Printf("Warning: Schema initialization failed: %v", err)
	}

	// Initialize AI service
	aiService := services.NewAIService(cfg.AnthropicAPIKey, cfg.OpenAIAPIKey)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "PRD Generator API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Initialize handlers
	handler := handlers.NewHandler(conn)
	generateHandler := handlers.NewGenerateHandler(conn, aiService)

	// ── Health check ────────────────────────────────────────────────────────
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "prd-generator-api",
		})
	})

	// ── API routes ──────────────────────────────────────────────────────────
	api := app.Group("/api")

	// Transcript routes
	api.Post("/transcript", handler.CreateTranscript)
	api.Get("/transcript/:id", handler.GetTranscript)

	// Generate routes
	api.Post("/generate/prd", generateHandler.GeneratePRD)
	api.Post("/generate/stories", generateHandler.GenerateStories)
	api.Post("/generate/requirements", generateHandler.GenerateRequirements)
	api.Post("/generate/wireframe", generateHandler.GenerateWireframe)

	// History routes
	api.Get("/history", handler.GetHistory)
	api.Get("/history/:id", handler.GetGeneration)

	// ── Documentation routes ─────────────────────────────────────────────────
	// These are generated LIVE from internal/apidocs so they are always in sync
	// with the codebase. Static files in docs/ can be regenerated with:
	//   go generate ./...

	// GET /docs/openapi.json — raw OpenAPI 3.0 spec
	app.Get("/docs/openapi.json", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.JSON(apidocs.OpenAPISpec())
	})

	// GET /docs/api-contracts.md — rendered markdown contract reference
	app.Get("/docs/api-contracts.md", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/markdown; charset=utf-8")
		return c.SendString(apidocs.Markdown())
	})

	// GET /docs — Swagger UI (loads spec from /docs/openapi.json)
	app.Get("/docs", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(swaggerUIHTML())
	})

	// ── Graceful shutdown ────────────────────────────────────────────────────
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		log.Println("Gracefully shutting down...")
		db.Close()
		_ = app.Shutdown()
	}()

	port := cfg.Port
	log.Printf("Starting server on port %s", port)
	log.Printf("Swagger UI:    http://localhost:%s/docs", port)
	log.Printf("OpenAPI JSON:  http://localhost:%s/docs/openapi.json", port)
	log.Printf("API Contracts: http://localhost:%s/docs/api-contracts.md", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// swaggerUIHTML returns a minimal self-contained Swagger UI page that loads
// the spec from /docs/openapi.json at runtime.
func swaggerUIHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>AI PRD Generator — API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <style>
    body { margin: 0; }
    #swagger-ui .topbar { background: #1a1a2e; }
    #swagger-ui .topbar .link { visibility: hidden; }
    #swagger-ui .topbar::before {
      content: "AI PRD Generator API";
      color: #fff;
      font-size: 1.1rem;
      font-weight: 600;
      padding: 0 1rem;
      display: flex;
      align-items: center;
    }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
    SwaggerUIBundle({
      url: "/docs/openapi.json",
      dom_id: "#swagger-ui",
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
      layout: "StandaloneLayout",
      deepLinking: true,
      defaultModelsExpandDepth: 1,
      defaultModelExpandDepth: 2,
    });
  </script>
</body>
</html>`
}
