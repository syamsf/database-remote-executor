package route

import (
	"database-remote-commander/internal/delivery/http"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	App                  *fiber.App
	HeaderAuthMiddleware fiber.Handler
	QueryController      *http.QueryController
}

func (c *Config) Setup() {
	c.SetupAuthRoute()
}

func (c *Config) SetupAuthRoute() {
	c.App.Use(c.HeaderAuthMiddleware)

	c.App.Post("/api/v1/query", c.QueryController.ExecQuery)
}
