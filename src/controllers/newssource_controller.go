package controllers

import (
	"log"
	"newsreader/db"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateNewssource(c *fiber.Ctx) error {
	var newssourceForm db.Newssource

	// TODO CSRF

	if err := c.BodyParser(&newssourceForm); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing form data")
	}

	// TODO validate values
	newssourceForm.ID = uuid.New()
	newssourceForm.UpdatePriority = db.MED
	newssourceForm.IsActive = true

	err := db.InsertNewssource(newssourceForm)
	if err != nil {
		log.Printf("Failed to insert newssource %s: %v", newssourceForm.Title, err)
	}

	c.Set("HX-Redirect", "/admin")
	return c.SendStatus(fiber.StatusNoContent) // 204 No Content
}

func EditNewssource(c *fiber.Ctx) error {
	var newssourceForm db.Newssource

	log.Printf("title %s", c.FormValue("title"))
	log.Printf("update_priority %s", c.FormValue("update_priority"))

	// TODO CSRF

	if err := c.BodyParser(&newssourceForm); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing form data")
	}

	// TODO validate values
	newssourceForm.UpdatePriority, _ = db.StringToUpdatePriority(c.FormValue("update_priority"))
	newssourceForm.IsActive = c.FormValue("is_active") == "1"

	err := db.UpdateNewssource(newssourceForm)
	if err != nil {
		log.Printf("Failed to update newssource %s: %v", newssourceForm.Title, err)
	}

	c.Set("HX-Redirect", "/admin")
	return c.SendStatus(fiber.StatusSeeOther) // 303 See Other
}
