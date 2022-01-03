package main

import (
	"log"
	"os"

	"github.com/batijo/poll-scraper/server"
	"github.com/batijo/poll-scraper/utils/csv"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}
}

func main() {
	srv := server.New()
	if os.Getenv("PS_WRITE_TO_CSV") == "true" {
		err := csv.Create(os.Getenv("PS_CSV_PATH"))
		if err != nil {
			log.Fatalln(err)
		}
		err = csv.StartWriting()
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Println("Running on: ", os.Getenv("PS_IP")+":"+os.Getenv("PS_PORT"))
	log.Fatal(srv.App.Listen(os.Getenv("PS_IP") + ":" + os.Getenv("PS_PORT")))
}
