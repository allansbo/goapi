package router

import (
	"github.com/allansbo/goapi/internal/app/server/handler"
	"github.com/gofiber/fiber/v2"
)

// MakeRoutes is a function that makes the routes for the application.
// It is used to define the routes for the application.
func MakeRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/locations", handler.LocationsAdd)
	v1.Get("/locations/:id", handler.LocationsGet)
	v1.Patch("/locations/:id", handler.LocationsUpdate)
	v1.Delete("/locations/:id", handler.LocationsDelete)
}
