package handlers

import (
	"github.com/fiure-cms/fiure-web/internal/fcore"
	"github.com/fiure-cms/fiure-web/internal/models"
	"github.com/fiure-cms/fiure-web/loggers"
	"github.com/fiure-cms/fiure-web/services"
	"github.com/gofiber/fiber/v2"
)

// Detail: Page
func PageDetail() fiber.Handler {
	return func(c *fiber.Ctx) error {

		slug := c.Params("slug")

		// Find Page
		item := &models.Page{}
		found, err := services.Pm.Get(slug, item)
		if err != nil {
			loggers.Sugar.With("error", err).Error("redis error")
		}

		if found == nil {
			return c.Next()
		}

		if item.Status == fcore.StatusPassive {
			return c.Next()
		}

		return c.Render("page", fiber.Map{
			"Title": generateMetaTitle("page", item.Name),
			"Social": fiber.Map{
				"Type":        "page",
				"Description": fcore.GetTruncateText(item.Content, 200),
				"Slug":        slug,
			},
			"Item": item,
		})
	}
}
