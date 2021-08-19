package main

import (
	"log"

	"github.com/batijo/poll-scraper/api/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}
}

func main() {
	app := fiber.New()
	app.Get("/", handlers.Data)
	app.Listen(":3000")
}
