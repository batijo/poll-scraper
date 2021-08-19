package main

import (
	"log"

	"github.com/batijo/poll-scraper/scraper"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}
}

func main() {
	scraper.Scrape()
}
