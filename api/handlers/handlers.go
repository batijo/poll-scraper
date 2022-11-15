package handlers

import (
	"os"
	"strings"

	"github.com/batijo/poll-scraper/scraper"
	"github.com/gofiber/fiber/v2"
)

func Data(c *fiber.Ctx) error {
	links := strings.Split(os.Getenv("PS_LINKS"), " ")
	return c.JSON(scraper.ScrapeAll(links))
}
