package server

import (
	"fmt"
	"github.com/allansbo/goapi/internal/app/server/middleware"
	"github.com/allansbo/goapi/internal/app/server/router"
	"github.com/allansbo/goapi/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"log/slog"
)

func Initialize(cfg *config.EnvConfig) {
	app := fiber.New()

	app.Use(healthcheck.New())
	middleware.UseJSONMiddleware(app)
	router.MakeRoutes(app)

	slog.Info("Server running", "Port", cfg.AppPort)
	if err := app.Listen(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
		slog.Error("error on running server", "error", err)
	}

}
