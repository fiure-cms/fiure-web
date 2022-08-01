package middlewares

import (
	"github.com/fiure-cms/fiure-web/loggers"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	loggers.Sugar.With("err", err).Info("bi bok var")

	return ctx.Status(fiber.StatusInternalServerError).Render("502", fiber.Map{
		"Title":   "502 - Server Error",
		"Message": "Oppps Something wrong!",
	}, "layouts/error")
}
