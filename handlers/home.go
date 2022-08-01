package handlers

import (
	"github.com/fiure-cms/fiure-web/loggers"
	"github.com/fiure-cms/fiure-web/services"
	"github.com/gofiber/fiber/v2"
)

// Homepage
func Home() fiber.Handler {
	return func(c *fiber.Ctx) error {

		lastid := c.Params("lastid")
		limit := 18

		// Get Lastest Post
		items, err := services.Bbm.PrevList(lastid, limit)
		if err != nil {
			loggers.Sugar.With("error", err).Error("bakbibu models list error")
		}

		if len(items) == 0 && lastid != "" {
			c.Redirect("/")
		}

		//loggers.Sugar.With("items", items).Info("items")

		return c.Render("home", fiber.Map{
			"Title":    generateMetaTitle("home", ""),
			"Social":   generateMetaData("home", "", ""),
			"Total":    len(items),
			"Result":   items,
			"NextPage": "",
		})
	}
}
