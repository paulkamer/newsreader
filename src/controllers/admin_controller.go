package controllers

import (
	"log"
	"newsreader/db"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AdminIndexPage(c *fiber.Ctx) error {
	newssources, _ := db.ListNewssources()

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
	guid, parse_err := uuid.Parse(c.Params("ID"))

	if parse_err != nil {
		return fiber.ErrBadRequest
	}

	newssource, err := db.FetchNewssource(guid)
	if err != nil {
		log.Printf("Failed to fetch newssource %s: %v", guid, err)
	}

	return c.Render("admin/newssources/edit", fiber.Map{
		"IsAdmin":    true,
		"Newssource": newssource,
		"Options":    []string{string(db.URGENT), string(db.HIGH), string(db.MED), string(db.LOW)},
	})
}

func AdminDeleteNewssource(c *fiber.Ctx) error {
	guid, parse_err := uuid.Parse(c.Params("ID"))

	if parse_err != nil {
		return fiber.ErrBadRequest
	}

	err := db.DeleteNewssource(guid)
	if err == nil {
		return c.SendString(c.Params("ID"))
	}

	return fiber.NewError(fiber.StatusInternalServerError, "Deleting failed")
}
