package server

import (
	"fmt"
	"log/slog"

	"github.com/allansbo/goapi/internal/app/server/middleware"
	"github.com/allansbo/goapi/internal/app/server/router"
	"github.com/allansbo/goapi/internal/domain/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

// Initialize is a function that initializes the server.
// It is responsible for setting up the server and starting the application
// by using the Fiber framework.
func Initialize() {
	app := fiber.New()

	app.Use(healthcheck.New())
	middleware.UseJSONMiddleware(app)
	router.MakeRoutes(app)

	slog.Info("Server running", "Port", usecase.AppConfig.AppPort)
	if err := app.Listen(fmt.Sprintf(":%s", usecase.AppConfig.AppPort)); err != nil {
		slog.Error("error on running server", "error", err)
	}

}
