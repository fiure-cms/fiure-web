package main

import (
	"github.com/fiure-cms/fiure-web/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {

	// Routes with Handlers

	// Sitemap: Static sitemap files in sitemap folder
	app.Static("/sitemap", "./static/sitemap", fiber.Static{
		Index: "index.xml",
	}).Name("fiure.sitemap")

	// Robot Txt
	app.Get("/robot.txt", handlers.RobotTxt()).Name("fiure.robot.txt")

	// Home
	app.Get("/:lastid?", handlers.Home()).Name("fiure.home")

	// Page
	app.Get("/p/:slug", handlers.PageDetail()).Name("fiure.page.detail")

	// Search
	app.Get("/suggest/:query", handlers.SearchPostSuggest()).Name("fiure.post.suggestion")
	app.Get("/s/:query/:paged?", handlers.SearchPost()).Name("fiure.post.search")
	app.Post("/s/:query", handlers.SearchPost())

	// Post
	app.Get("/:slug-:id", handlers.PostDetail()).Name("fiure.post.detail")
}
