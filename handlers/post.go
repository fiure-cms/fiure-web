package handlers

import (
	"github.com/fiure-cms/fiure-web/internal/fcore"
	"github.com/fiure-cms/fiure-web/internal/models"
	"github.com/fiure-cms/fiure-web/loggers"
	"github.com/fiure-cms/fiure-web/services"
	"github.com/gofiber/fiber/v2"
)

// Detail: Item
func PostDetail() fiber.Handler {
	return func(c *fiber.Ctx) error {

		shourlsid := c.Params("shourlsid")

		// Find Dream
		item := &models.Post{}
		_, err := services.Bbm.Get(shourlsid, item)
		if err != nil {
			loggers.Sugar.With("error", err).Error("redis error")
		}

		if item == nil {
			return c.Next()
		}

		if item.Status != fcore.StatusActive {
			return c.Next()
		}

		// Prev - Next Items
		prevResult, err := services.Bbm.PrevList(shourlsid, 1)
		if err != nil {
			loggers.Sugar.With("error", err).Error("livestore conn error")
		}

		nextResult, err := services.Bbm.List(shourlsid, 1)
		if err != nil {
			loggers.Sugar.With("error", err).Error("livestore conn error")
		}

		// Search: Category Releated Items
		query := ""
		if len(item.Cats) > 0 {
			query = item.Cats[0].Name
		}

		uids, err := services.Sm.Search(item.Collection, item.Bucket, query, 10, 0)
		if err != nil {
			loggers.Sugar.With("error", err).Error("sonic search error")
		}

		// Collect Post
		relatedItems := services.Bbm.GetResultFromUIDS(uids)

		return c.Render("item", fiber.Map{
			"Title": generateMetaTitle("single", item.Name),
			"Social": fiber.Map{
				"Type":        "single",
				"Description": fcore.GetTruncateText(item.Slug, 200),
				"Slug":        item.Slug + "/" + item.UID,
			},
			"Item": item,
			"ItemsNavi": fiber.Map{
				"Prev": prevResult,
				"Next": nextResult,
			},
			"RelatedItems": relatedItems,
		})
	}
}
