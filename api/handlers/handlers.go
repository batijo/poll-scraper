package handlers

import (
	"github.com/batijo/poll-scraper/scraper"
	"github.com/gofiber/fiber/v2"
)

func Data(c *fiber.Ctx) error {
	return c.JSON(scraper.Scrape())
}
