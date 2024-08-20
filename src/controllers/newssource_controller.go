package controllers

import (
	"fmt"
	"log"
	"net/url"
	"newsreader/config"
	"newsreader/jobs"
	"newsreader/models"
	"newsreader/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func NewssourcePage(c *fiber.Ctx) error {
	appconfig := c.Locals("appconfig").(*config.AppConfig)
	guid, parse_err := uuid.Parse(c.Params("ID"))
	if parse_err != nil {
		return fiber.ErrBadRequest
	}

	newssource, err := repositories.FetchNewssource(appconfig.DB, guid)

	if err != nil {
		fmt.Printf("Failed to fetch newssource: %v", err)
		return fiber.ErrBadRequest
	}

	articles, err := repositories.ListArticles(appconfig.DB, guid)

	if err != nil {
		fmt.Printf("Failed to fetch articles: %v", err)
		return fiber.ErrInternalServerError
	}

	fmt.Printf("articles: %v", articles)

	return c.Render("feed", fiber.Map{
		"Newssource": newssource,
		"Articles":   articles,
	})
}

func AddNewssource(c *fiber.Ctx) error {
	appconfig := c.Locals("appconfig").(*config.AppConfig)
	var newssourceForm models.Newssource

	if err := c.BodyParser(&newssourceForm); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing form data")
	}

	// TODO validate values
	newssourceForm.ID = uuid.New()
	newssourceForm.UpdatePriority = models.MED
	newssourceForm.FeedType = models.RSS
	newssourceForm.IsActive = true

	if !validateNewssourceUrl(newssourceForm.Url) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid URL")
	}

	err := repositories.InsertNewssource(appconfig.DB, newssourceForm)
	if err != nil {
		log.Printf("Failed to insert newssource %s: %v", newssourceForm.Title, err)
	}

	go jobs.FetchNews(newssourceForm.ID)

	c.Set("HX-Redirect", "/admin")
	return c.SendStatus(fiber.StatusNoContent) // 204 No Content
}

func EditNewssource(c *fiber.Ctx) error {
	// TODO fetch news source from DB
	appconfig := c.Locals("appconfig").(*config.AppConfig)
	var newssourceForm models.Newssource

	if err := c.BodyParser(&newssourceForm); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing form data")
	}

	// TODO validate values
	newssourceForm.UpdatePriority, _ = models.StringToUpdatePriority(c.FormValue("update_priority"))
	newssourceForm.IsActive = c.FormValue("is_active") == "1"
	newssourceForm.FeedType = models.RSS

	if !validateNewssourceUrl(newssourceForm.Url) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid URL")
	}

	err := repositories.UpdateNewssource(appconfig.DB, newssourceForm)
	if err != nil {
		log.Printf("Failed to update newssource %s: %v", newssourceForm.Title, err)
	}

	c.Set("HX-Redirect", "/admin")
	return c.SendStatus(fiber.StatusSeeOther) // 303 See Other
}

func validateNewssourceUrl(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check if the scheme is https
	if parsedURL.Scheme != "https" {
		return false
	}

	// Check if the host is not empty
	if parsedURL.Host == "" {
		return false
	}

	return true
}
