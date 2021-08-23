package handlers

import (
	"os"

	"github.com/batijo/poll-scraper/models"
	"github.com/batijo/poll-scraper/scraper"
	"github.com/gofiber/fiber/v2"
)

func Data(c *fiber.Ctx) error {
	var data []models.Data
	if os.Getenv("PS_WITH_EQ") == "true" {
		data = scraper.ScrapeWithEquals()
	} else {
		data = scraper.Scrape()
	}
	return c.JSON(data)
}
