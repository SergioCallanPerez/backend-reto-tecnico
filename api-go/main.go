package main

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/SergioCallanPerez/backend-reto-tecnico/api-go/internal/client"
	"github.com/SergioCallanPerez/backend-reto-tecnico/api-go/internal/handler"
	"github.com/SergioCallanPerez/backend-reto-tecnico/api-go/internal/logging"
	"github.com/SergioCallanPerez/backend-reto-tecnico/api-go/internal/service"
)

const defaultHTTPClientTimeoutMs = 2000

func main() {
	logging.Setup(os.Getenv("LOG_LEVEL"))

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		slog.Info("http_request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	nodeAPIURL := os.Getenv("NODE_API_URL")
	if nodeAPIURL == "" {
		nodeAPIURL = "http://localhost:3000"
	}
	timeoutMs := defaultHTTPClientTimeoutMs
	if v, err := strconv.Atoi(os.Getenv("HTTP_CLIENT_TIMEOUT_MS")); err == nil {
		timeoutMs = v
	}
	statsClient := client.NewStatsClient(nodeAPIURL, time.Duration(timeoutMs)*time.Millisecond)

	qrService := service.NewQRService()
	matrixHandler := handler.NewMatrixHandler(qrService, statsClient)

	rotateService := service.NewRotateService()
	rotateHandler := handler.NewRotateHandler(rotateService)

	v1 := app.Group("/api/v1")
	matrixHandler.RegisterRoutes(v1)
	rotateHandler.RegisterRoutes(v1)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := app.Listen(":" + port); err != nil {
		panic(err)
	}
}
