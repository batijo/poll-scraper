package main

import (
	"log"
	"os"

	"github.com/batijo/poll-scraper/server"
	"github.com/batijo/poll-scraper/utils/file"
	"github.com/joho/godotenv"
)

func init() {
	logFile, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(logFile)
	if err := godotenv.Load(); err != nil {
		log.Fatalln(".env file not found")
	}
}

func main() {
	srv := server.New()
	err := file.InitFiles()
	if err != nil {
		log.Fatalln(err)
	}
	err = file.StartWriting()
	if err != nil {
		log.Fatalln(err)
	}
	if os.Getenv("PS_IP") == "" {
		os.Setenv("PS_IP", "localhost")
	}
	log.Println("Running on: ", os.Getenv("PS_IP")+":"+os.Getenv("PS_PORT"))
	log.Fatal(srv.App.Listen(os.Getenv("PS_IP") + ":" + os.Getenv("PS_PORT")))
}
