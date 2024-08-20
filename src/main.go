package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/template/html/v2"

	"newsreader/config"
	"newsreader/controllers"
	"newsreader/db"
)

func main() {
	engine := html.New("./views", ".html")

	app := initApp(engine)
	dbconn := initDatabase(app)
	defer dbconn.Close()

	// CSRF Middleware
	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieSameSite: "Strict",
		Expiration:     8600,
		KeyGenerator:   utils.UUIDv4,
		ContextKey:     "csrf",
	}))

	// Routes
	app.Get("/", controllers.Indexpage)

	app.Get("/newssources/:ID", controllers.NewssourcePage)

	app.Get("/admin", controllers.AdminIndexPage)
	app.Get("/admin/newssources/add", controllers.AdminAddNewssourcePage)
	app.Get("/admin/newssources/edit/:ID", controllers.AdminEditNewssourcePage)

	app.Post("/newssources", controllers.AddNewssource)
	app.Put("/newssources", controllers.EditNewssource)
	app.Delete("/newssources/:ID", controllers.AdminDeleteNewssource)

	log.Fatal(app.Listen(":3001"))
}

func initApp(engine *html.Engine) *fiber.App {
	return fiber.New(fiber.Config{
		Views:             engine,
		ViewsLayout:       "base",
		PassLocalsToViews: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			log.Printf("An error occurred: %v\n", err)

			return c.Status(code).SendString("An error occurred")
		},
	})
}

func initDatabase(app *fiber.App) *sql.DB {
	dbConn, _ := db.InitDatabase(db.SQLiteType, db.SQLiteDataSource)

	appconfig := &config.AppConfig{DB: dbConn}
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("appconfig", appconfig)
		return c.Next()
	})

	return dbConn
}
