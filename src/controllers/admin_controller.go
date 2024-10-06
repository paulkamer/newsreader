//go:build !excludetest

package controllers

import (
	"newsreader/config"
	"newsreader/jobs"
	"newsreader/models"
	"newsreader/repositories"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var newssourcePriorityOptions = []string{string(models.URGENT), string(models.HIGH), string(models.MED), string(models.LOW)}

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
		"Options":    newssourcePriorityOptions,
	})
}

func AdminAddNewssource(c *fiber.Ctx) error {
	appconfig := c.Locals("appconfig").(*config.AppConfig)
	var newssourceForm models.Newssource

	if err := c.BodyParser(&newssourceForm); err != nil {
		log.Debugf("Error parsing form data: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing form data")
	}

	newssourceForm.ID = uuid.New()
	newssourceForm.UpdatePriority = models.MED
	newssourceForm.IsActive = true

	errors := validateForm(newssourceForm)
	if len(errors) > 0 {
		return c.Render("admin/newssource/_addform.html", fiber.Map{
			"Errors":     errors,
			"Newssource": &newssourceForm,
		}, "base_empty")
	}

	err := repositories.InsertNewssource(appconfig.DB, &newssourceForm)
	if err != nil {
		log.Errorf("Failed to insert newssource %s: %v", newssourceForm.Title, err)
	}

	go jobs.FetchNews(newssourceForm.ID)

	c.Set("HX-Redirect", "/admin")
	return c.SendStatus(fiber.StatusNoContent) // 204 No Content
}

func AdminEditNewssource(c *fiber.Ctx) error {
	appconfig := c.Locals("appconfig").(*config.AppConfig)
	var newssourceForm models.Newssource

	if err := c.BodyParser(&newssourceForm); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing form data")
	}

	orgNewssource, err := repositories.FetchNewssource(appconfig.DB, newssourceForm.ID)
	if err != nil {
		log.Errorf("Failed to fetch newssource: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Failed to fetch newssource")
	}

	errors := validateForm(newssourceForm)
	if len(errors) > 0 {
		return c.Render("admin/newssources/_editform", fiber.Map{
			"Errors":     errors,
			"Newssource": &newssourceForm,
			"Options":    newssourcePriorityOptions,
		}, "base_empty")
	}

	orgNewssource.Title = newssourceForm.Title
	orgNewssource.Url = newssourceForm.Url
	orgNewssource.UpdatePriority, _ = models.StringToUpdatePriority(c.FormValue("updatepriority"))
	orgNewssource.IsActive = c.FormValue("isactive") == "1"

	err = repositories.UpdateNewssource(appconfig.DB, &orgNewssource)
	if err != nil {
		log.Errorf("failed to update newssource %s: %v", orgNewssource.Title, err)
	}

	c.Set("HX-Redirect", "/admin")
	return c.SendStatus(fiber.StatusSeeOther) // 303 See Other
}

func AdminDeleteNewssource(c *fiber.Ctx) error {
	appconfig := c.Locals("appconfig").(*config.AppConfig)
	guid, parse_err := uuid.Parse(c.Params("ID"))
	if parse_err != nil {
		return fiber.ErrBadRequest
	}

	redirect := c.Query("redirect", "false")

	err := repositories.DeleteNewssource(appconfig.DB, guid)
	if err == nil {
		if redirect == "true" {
			c.Response().Header.Set("HX-redirect", "/admin")
		}

		return c.SendString(c.Params("ID"))
	}

	return fiber.NewError(fiber.StatusInternalServerError, "Deleting failed")
}

func validateForm(form models.Newssource) map[string]string {
	validatorInst := validator.New(validator.WithRequiredStructEnabled())

	err := validatorInst.Struct(form)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if len(validationErrors) > 0 {
			log.Debugf("Validation errors: %v", validationErrors)

			errors := make(map[string]string)
			for _, err := range err.(validator.ValidationErrors) {
				errors[err.Field()] = err.Tag()
			}

			return errors
		}
	}

	return nil
}
