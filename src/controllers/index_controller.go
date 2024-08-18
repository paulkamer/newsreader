package controllers

import (
	"newsreader/config"
	"newsreader/repositories"

	"github.com/gofiber/fiber/v2"
)

func Indexpage(c *fiber.Ctx) error {
	appconfig := c.Locals("appconfig").(*config.AppConfig)

	newssources, _ := repositories.ListNewssources(appconfig.DB)

	return c.Render("index", fiber.Map{
		"Newssources": newssources,
		"IsAdmin":     false,
	})
}
