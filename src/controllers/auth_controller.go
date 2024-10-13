package controllers

import "github.com/gofiber/fiber/v2"

type User struct {
	Username string
	Password string
}

var users = []User{
	{Username: "user1", Password: "password1"},
	{Username: "user2", Password: "password2"},
}

func LoginPage(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

func HandleLogin(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	for _, user := range users {
		if user.Username == username && user.Password == password {

			store.Set(c, "user", user.Username)

			c.Set("HX-Redirect", "/")
			return c.SendStatus(fiber.StatusSeeOther) // 303 See Other
		}
	}

	return c.SendString("Invalid username or password")
}

func HandleLogout(c *fiber.Ctx) error {
	return c.SendString("Logged out!")
}
