package router

import "github.com/gofiber/fiber/v2"

func MakeRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	v1.Put("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	v1.Delete("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
}
