//go:build !excludetest

package main

import (
	"database/sql"
	"encoding/json"
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
var store = session.New()

const newsUpdateInterval = 10 * time.Minute

func main() {
	setLogLevel()

	app := initApp(html.New("./views", ".html"))

	dbconn := initDatabase(app)
	defer dbconn.Close()

	// Routes
	app.Get("/", authRequired, controllers.Indexpage)

	app.Get("/login", controllers.LoginPage)
	app.Post("/login", login)
	app.Delete("/logout", authRequired, logout)

	app.Get("/newssources/:ID", authRequired, controllers.NewssourcePage)

	adminGroup := app.Group("/admin", authRequired, adminRequired)
	adminGroup.Get("/", controllers.AdminIndexPage)
	adminGroup.Get("/newssources/add", controllers.AdminAddNewssourcePage)
	adminGroup.Get("/newssources/edit/:ID", controllers.AdminEditNewssourcePage)
	adminGroup.Post("/newssources", authRequired, controllers.AdminAddNewssource)
	adminGroup.Put("/newssources", authRequired, controllers.AdminEditNewssource)
	adminGroup.Delete("/newssources/:ID", authRequired, controllers.AdminDeleteNewssource)

	app.Get("/article/:ID", authRequired, controllers.ArticlePage)

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

			log.Error("An error occurred: ", "error", err)

			return c.Status(code).SendString("An error occurred")
		},
	})

	app.Use(requestLogger())
	app.Use(csrf.New(csrf.Config{ContextKey: "csrf"}))
	app.Use(favicon.New(favicon.Config{File: "./favicon.ico", URL: "/favicon.ico"}))

	app.Use(pprof.New())

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

		err := c.Next()

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

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

func login(c *fiber.Ctx) error {
	var creds Credentials
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// Load user credentials from JSON file
	users, err := fetchUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot open users file"})
	}

	if user, ok := users[creds.Username]; ok && user.Password == creds.Password {
		sess, err := store.Get(c)
		if err != nil {
			return err
		}
		sess.Set("authenticated", true)
		sess.Set("username", creds.Username)
		sess.Set("isAdmin", user.IsAdmin)
		if err := sess.Save(); err != nil {
			return err
		}

		c.Set("HX-Redirect", "/")
		return c.SendStatus(fiber.StatusSeeOther) // 303 See Other
	}

	c.Set("HX-Refresh", "true")
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
}

func logout(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	if err := sess.Destroy(); err != nil {
		return err
	}

	c.Set("HX-Redirect", "/login")
	return c.SendStatus(fiber.StatusSeeOther) // 303 See Other
}

func fetchUsers() (map[string]User, error) {
	file, err := os.Open("users.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users map[string]User
	if err := json.NewDecoder(file).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

func authRequired(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	if auth, ok := sess.Get("authenticated").(bool); !ok || !auth {
		return c.Redirect("/login")
	}

	c.Locals("is_authenticated", true)

	isAdmin, _ := sess.Get("isAdmin").(bool)

	log.WithFields(logrus.Fields{
		"username": sess.Get("username"),
		"isAdmin":  isAdmin,
	}).Info("auth check")
	c.Locals("is_admin", isAdmin)

	return c.Next()
}

func adminRequired(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return err
	}

	if isAdmin, ok := sess.Get("isAdmin").(bool); !ok || !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}

	return c.Next()
}
