# poll-scraper

## Setup

- install [go](https://golang.org/dl/) at least 1.16 version
- edit `example.env` and then remove example and save it as `.env`
- open terminal or command line in project directory and type commands bellow:

```sh
go mod download
go run main.go
```

## Info

All logs and errors can be found in info.log file.

To create windows executable program run:

```sh
go build PollScraper.exe main.go
```
