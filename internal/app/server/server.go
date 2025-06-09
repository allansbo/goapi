package server

import (
	"fmt"
	"log/slog"

	"github.com/allansbo/goapi/internal/app/server/middleware"
	"github.com/allansbo/goapi/internal/app/server/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

type AppServer struct {
	FiberApp *fiber.App
	appPort  string
}

func NewAppServer(appPort string) *AppServer {
	return &AppServer{
		FiberApp: fiber.New(),
		appPort:  appPort,
	}
}

// Start is a function that initializes the server.
// It is responsible for setting up the server and starting the application
// by using the Fiber framework.
func (s *AppServer) Start() {
	s.FiberApp.Use(healthcheck.New())
	middleware.UseJSONMiddleware(s.FiberApp)
	router.MakeRoutes(s.FiberApp)

	slog.Info("Server running", "Port", s.appPort)
	if err := s.FiberApp.Listen(fmt.Sprintf(":%s", s.appPort)); err != nil {
		slog.Error("error on running server", "error", err)
	}

}
