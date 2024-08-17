package controllers

import (
	"newsreader/db"

	"github.com/gofiber/fiber/v2"
)

func Indexpage(c *fiber.Ctx) error {
	newssources, _ := db.ListNewssources()

	return c.Render("index", fiber.Map{
		"Newssources": newssources,
		"IsAdmin":     false,
	})
}
