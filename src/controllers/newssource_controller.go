//go:build !excludetest

package controllers

import (
	"newsreader/config"
	"newsreader/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func NewssourcePage(c *fiber.Ctx) error {
	appconfig := c.Locals("appconfig").(*config.AppConfig)
	guid, parse_err := uuid.Parse(c.Params("ID"))
	if parse_err != nil {
		return fiber.ErrBadRequest
	}

	newssource, err := repositories.FetchNewssource(appconfig.DB, guid)

	if err != nil {
		log.Errorf("Failed to fetch newssource: %v", err)
		return fiber.ErrBadRequest
	}

	articles, err := repositories.ListArticles(appconfig.DB, guid)

	if err != nil {
		log.Errorf("Failed to fetch articles: %v", err)
		return fiber.ErrInternalServerError
	}

	return c.Render("feed", fiber.Map{
		"Newssource": &newssource,
		"Articles":   &articles,
	})
}
