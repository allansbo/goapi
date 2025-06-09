package middleware

import (
	"github.com/allansbo/goapi/internal/app/server/handler"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func UseJSONMiddleware(app *fiber.App) {
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
