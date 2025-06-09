package middleware

import (
	"net/http"

	"github.com/allansbo/goapi/internal/app/server/handler"
	"github.com/gofiber/fiber/v2"
)

// UseJSONMiddleware is a middleware that checks if the request is a JSON request
// and returns a 400 error if it is not. It is used to validate the request body.
func UseJSONMiddleware(app *fiber.App) {

	// Always that a user send data to the server,
	// we need to check if the request is a JSON request.
	// If it is not, we return a 400 error.
	app.Use(func(ctx *fiber.Ctx) error {
		switch ctx.Method() {
		case fiber.MethodPost, fiber.MethodPut, fiber.MethodPatch:
			if ctx.Is("json") {
				return ctx.Next()
			}

			return ctx.Status(http.StatusBadRequest).JSON(handler.GlobalErrorHandlerResp{
				Success: false,
				Message: "invalid format",
				Error:   "only json is allowed",
			})
		default:
			return ctx.Next()
		}
	})
}
