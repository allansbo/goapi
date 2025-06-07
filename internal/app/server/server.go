package server

import (
	"fmt"
	"github.com/allansbo/goapi/internal/app/server/middleware"
	"github.com/allansbo/goapi/internal/app/server/router"
	"github.com/allansbo/goapi/internal/config"
	"github.com/gofiber/fiber/v2"
	"log"
)

func RunServer(cfg *config.EnvConfig) {
	app := fiber.New()

	middleware.UseJSONMiddleware(app)
	router.MakeRoutes(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", cfg.AppPort)))
}
