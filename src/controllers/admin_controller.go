package controllers

import (
	"newsreader/config"
	"newsreader/models"
	"newsreader/repositories"

	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AdminIndexPage(c *fiber.Ctx) error {
	appconfig := c.Locals("appconfig").(*config.AppConfig)
	newssources, _ := repositories.ListNewssources(appconfig.DB)

	return c.Render("admin/index", fiber.Map{
		"Newssources": newssources,
		"IsAdmin":     true,
	})
}

func AdminAddNewssourcePage(c *fiber.Ctx) error {
	return c.Render("admin/newssources/add", fiber.Map{
		"IsAdmin": true,
	})
}

func AdminEditNewssourcePage(c *fiber.Ctx) error {
	appconfig := c.Locals("appconfig").(*config.AppConfig)
	guid, parse_err := uuid.Parse(c.Params("ID"))

	if parse_err != nil {
		return fiber.ErrBadRequest
	}

	newssource, err := repositories.FetchNewssource(appconfig.DB, guid)
	if err != nil {
		log.Errorf("Failed to fetch newssource %s: %v", guid, err)
	}

	return c.Render("admin/newssources/edit", fiber.Map{
		"IsAdmin":    true,
		"Newssource": newssource,
		"Options":    []string{string(models.URGENT), string(models.HIGH), string(models.MED), string(models.LOW)},
	})
}

func AdminDeleteNewssource(c *fiber.Ctx) error {
	appconfig := c.Locals("appconfig").(*config.AppConfig)
	guid, parse_err := uuid.Parse(c.Params("ID"))

	if parse_err != nil {
		return fiber.ErrBadRequest
	}

	err := repositories.DeleteNewssource(appconfig.DB, guid)
	if err == nil {
		return c.SendString(c.Params("ID"))
	}

	return fiber.NewError(fiber.StatusInternalServerError, "Deleting failed")
}
