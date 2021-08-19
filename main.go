package main

import (
	"log"
	"os"

	"github.com/batijo/poll-scraper/server"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}
}

func main() {
	srv := server.New()
	log.Fatal(srv.App.Listen(os.Getenv("PS_IP") + ":" + os.Getenv("PS_PORT")))
}
