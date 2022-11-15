package handlers

import (
	"log"
	"os"
	"strings"

	"github.com/batijo/poll-scraper/models"
	"github.com/batijo/poll-scraper/scraper"
	"github.com/batijo/poll-scraper/utils"
	"github.com/gofiber/fiber/v2"
)

func Data(c *fiber.Ctx) error {
	links := strings.Split(os.Getenv("PS_LINKS"), " ")
	data := scraper.ScrapeAll(links)
	lines, err := utils.GetFilterLines("PS_FILTER_LINES")
	if err != nil {
		log.Println("ERROR: cannot parse filter lines")
		log.Println(err)
	}
	if len(lines) > 0 {
		data = models.FilterData(lines, data)
	}
	if os.Getenv("PS_ADD_SUM") == "true" {
		data = models.SumData(data)
	}
	return c.JSON(data)
}
