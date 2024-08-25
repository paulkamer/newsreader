//go:build !excludetest

package controllers

import (
	"newsreader/config"
	"newsreader/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func ArticlePage(c *fiber.Ctx) error {
	appconfig := c.Locals("appconfig").(*config.AppConfig)

	id, err := uuid.Parse(c.Params("ID"))
	if err != nil {
		log.Error("Failed to parse UUID", "error", err)
		return fiber.ErrBadRequest
	}

	article, err := repositories.FetchArticle(appconfig.DB, id)
	if err != nil {
		log.Error("Failed to fetch article", "error", err)
		// return fiber.ErrBadRequest
	}

	newssource, err := repositories.FetchNewssource(appconfig.DB, article.Source)
	if err != nil {
		log.Error("Failed to fetch newssource for article", "error", err)
		return fiber.ErrBadRequest
	}

	return c.Render("article", fiber.Map{
		"Newssource": newssource,
		"Article":    article,
	})
}
