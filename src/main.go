//go:build !excludetest

package main

import (
	"database/sql"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"

	"newsreader/config"
	"newsreader/controllers"
	"newsreader/db"
	"newsreader/jobs"

	"github.com/sirupsen/logrus"

	_ "net/http/pprof"
)

var log = logrus.New()

const newsUpdateInterval = 10 * time.Minute

func main() {
	setLogLevel()

	app := initApp(html.New("./views", ".html"))

	app.Use(pprof.New())

	dbconn := initDatabase(app)
	defer dbconn.Close()

	// Routes
	app.Get("/", controllers.Indexpage)

	app.Get("/newssources/:ID", controllers.NewssourcePage)

	app.Get("/admin", controllers.AdminIndexPage)
	app.Get("/admin/newssources/add", controllers.AdminAddNewssourcePage)
	app.Get("/admin/newssources/edit/:ID", controllers.AdminEditNewssourcePage)

	app.Post("/newssources", controllers.AddNewssource)
	app.Put("/newssources", controllers.EditNewssource)
	app.Delete("/newssources/:ID", controllers.AdminDeleteNewssource)

	app.Get("/article/:ID", controllers.ArticlePage)

	startNewsUpdateScheduler()

	log.Fatal(app.Listen(":3001"))
}

func setLogLevel() {
	level, _ := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	logrus.SetLevel(level)
}

func initApp(engine *html.Engine) *fiber.App {
	app := fiber.New(fiber.Config{
		Views:             engine,
		ViewsLayout:       "base",
		PassLocalsToViews: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			log.Error("An error occurred", "error", err)

			return c.Status(code).SendString("An error occurred")
		},
	})

	app.Use(requestLogger())

	app.Use(csrf.New(csrf.Config{
		ContextKey: "csrf",
	}))

	return app
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

func requestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next() // Proceed with the next middleware or handler

		log.WithFields(logrus.Fields{
			"method":     c.Method(),
			"path":       c.Path(),
			"status":     c.Response().StatusCode(),
			"latency_ns": time.Since(start).Nanoseconds(),
		}).Info("request")

		return err
	}
}

func startNewsUpdateScheduler() {
	listChan := make(chan uuid.UUID)

	go func() {
		go jobs.DetermineOutdatedNewssources(listChan) // Trigger immediately

		ticker := time.NewTicker(newsUpdateInterval)
		defer ticker.Stop()

		for range ticker.C {
			go jobs.DetermineOutdatedNewssources(listChan)
		}
	}()

	go func() {
		for id := range listChan {
			go jobs.FetchNews(id)
		}
	}()
}
