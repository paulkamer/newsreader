package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"newsreader/config"
	"newsreader/controllers"
	"newsreader/db"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "base",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			log.Printf("An error occurred: %v\n", err)

			return c.Status(code).SendString("An error occurred")
		},
	})

	dbConn, _ := db.InitDatabase(db.SQLiteType, db.SQLiteDataSource)
	defer dbConn.Close()

	appconfig := &config.AppConfig{DB: dbConn}
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("appconfig", appconfig)
		return c.Next()
	})

	// Routes
	app.Get("/", controllers.Indexpage)

	app.Get("/admin", controllers.AdminIndexPage)
	app.Get("/admin/newssources/add", controllers.AdminAddNewssourcePage)
	app.Get("/admin/newssources/edit/:ID", controllers.AdminEditNewssourcePage)

	app.Post("/newssources", controllers.CreateNewssource)
	app.Put("/newssources", controllers.EditNewssource)
	app.Delete("/newssources/:ID", controllers.AdminDeleteNewssource)

	log.Fatal(app.Listen(":3001"))
}
