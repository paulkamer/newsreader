//go:build !excludetest

package main

import (
	"database/sql"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/pprof"

	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/gofiber/fiber/v2/middleware/favicon"

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

	store := session.New()
	app := initApp(html.New("./views", ".html"))

	dbconn := initDatabase(app)
	defer dbconn.Close()

	// Routes
	app.Get("/", controllers.Indexpage)

	auth := app.Group("/auth")
	auth.Get("/login", controllers.LoginPage)
	auth.Post("/login", controllers.HandleLogin)
	auth.Post("/logout", controllers.HandleLogout)

	app.Get("/newssources/:ID", controllers.NewssourcePage)

	admin := app.Group("/admin")
	admin.Get("/", controllers.AdminIndexPage)
	admin.Get("/newssources/add", controllers.AdminAddNewssourcePage)
	admin.Get("/newssources/edit/:ID", controllers.AdminEditNewssourcePage)

	app.Post("/newssources", controllers.AdminAddNewssource)
	app.Put("/newssources", controllers.AdminEditNewssource)
	app.Delete("/newssources/:ID", controllers.AdminDeleteNewssource)

	app.Get("/article/:ID", controllers.ArticlePage)

	startNewsUpdateScheduler()

	log.Fatal(app.Listen(":3001"))
}

func setLogLevel() {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	logrusLevel, _ := logrus.ParseLevel(level)
	logrus.SetLevel(logrusLevel)
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

	app.Use(pprof.New())

	app.Use(favicon.New(favicon.Config{
		File: "./favicon.ico",
		URL:  "/favicon.ico",
	}))

	return app
}

func initDatabase(app *fiber.App) *sql.DB {
	dbConn, _ := db.InitDatabase(db.SQLiteType, db.SQLiteDataSource, db.MigrationsDir)

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
