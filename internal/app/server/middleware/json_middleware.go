package middleware

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func UseJSONMiddleware(app *fiber.App) {
	app.Use(func(ctx *fiber.Ctx) error {
		switch ctx.Method() {
		case fiber.MethodPost, fiber.MethodPut:
			if ctx.Is("json") {
				return ctx.Next()
			}

			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid format. Only JSON is allowed."})
		default:
			return ctx.Next()
		}
	})
}
