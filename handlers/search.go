package handlers

import (
	"net/url"
	"strconv"

	"github.com/fiure-cms/fiure-web/internal/fcore"
	"github.com/fiure-cms/fiure-web/internal/models"
	"github.com/fiure-cms/fiure-web/loggers"
	"github.com/fiure-cms/fiure-web/services"
	"github.com/gofiber/fiber/v2"
)

// Search
func SearchPostSuggest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Bucket Params
		bucketParam := fcore.BUCKET_POST

		// Query with Page
		query := fcore.SearchTextLimit(c.Params("query"), fcore.SearchTextLimiter)
		query, err := url.QueryUnescape(query)
		if err != nil {
			loggers.Sugar.With("error", err).Error("QueryUnescape error")
		}

		// Suggest List
		results, err := services.Sm.Suggest(fcore.COLLECTION_POST_MODELS, bucketParam, query, 5)
		if err != nil {
			loggers.Sugar.With("error", err).Error("search client error")
		}

		loggers.Sugar.With("results", results, "query", query, "bucket", bucketParam).Info("Suggest")

		suggestions := []models.SuggestionItem{}
		if len(results) > 0 {
			for _, item := range results {
				suggestions = append(suggestions, models.SuggestionItem{
					Label: item,
					Value: item,
				})
			}
		}

		return c.JSON(fiber.Map{
			"result": suggestions,
		})
	}
}

func SearchPost() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Bucket Params
		bucketParam := fcore.BUCKET_POST // c.Cookies("sbucket")

		// Query with Page
		query := fcore.SearchTextLimit(c.Params("query"), fcore.SearchTextLimiter)
		paged := c.Params("paged", "1")
		pagedInt, _ := strconv.Atoi(paged)
		limit := 18
		offset := (fcore.Max(1, pagedInt) - 1) * limit

		//loggers.Sugar.Info(query)
		if c.Request().Header.IsPost() {
			query = fcore.SearchTextLimit(c.FormValue("query"), fcore.SearchTextLimiter)

			return c.Redirect("/s/"+query, 301)
		}

		// SonicSearch Query
		query, err := url.QueryUnescape(query)
		if err != nil {
			loggers.Sugar.With("error", err).Error("Query Unescape error")
		}
		loggers.Sugar.Info(query, bucketParam)

		// SonicSearch Result
		results, err := services.Sm.Search(fcore.COLLECTION_POST_MODELS, bucketParam, query, limit, offset)
		if err != nil {
			loggers.Sugar.With("error", err).Error("sonic search error")
		}

		if len(results) == 0 && pagedInt > 1 {
			c.Redirect("/")
		}

		// Collect PostItem
		items := services.Bbm.GetResultFromUIDS(results)

		return c.Render("search", fiber.Map{
			"Title":    generateMetaTitle("term", query),
			"Social":   generateMetaData("search", "", query),
			"Query":    query,
			"Total":    len(items),
			"NextPage": pagedInt + 1,
			"Result":   items,
		})
	}
}
