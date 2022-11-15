package server

import (
	"os"

	"github.com/batijo/poll-scraper/api/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Obj struct {
	*fiber.App
}

func (srv *Obj) registerApiRoutes() {
	srv.App.Get("/", handlers.Data)
}

func New() *Obj {
	srv := Obj{
		App: fiber.New(fiber.Config{
			GETOnly:               true,
			ServerHeader:          "poll-scraper",
			DisableStartupMessage: true,
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			},
		}),
	}
	srv.registerMiddleware()
	srv.registerApiRoutes()
	srv.statusNotFoundMiddleware()
	return &srv
}

func (srv *Obj) registerMiddleware() {
	origins := "*"
	if os.Getenv("PS_DOMAINS") != "" {
		origins = os.Getenv("PS_DOMAINS")
	}
	srv.Use(cors.New(cors.Config{
		AllowOrigins: origins,
		AllowMethods: "GET",
	}))
}

func (srv *Obj) statusNotFoundMiddleware() {
	srv.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(
			map[string]interface{}{"message": "page not found"},
		)
	})
}
