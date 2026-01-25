# poll-scraper

An application designed to collect and structure data from the [savitarna.mobi.lt](https://savitarna.mobi.lt) website with the purpose of enabling its integration into various graphic engines.

## Setup

- install [go](https://golang.org/dl/) at least 1.21 version
- copy `config.example.json` to `config.json` and edit it according to your needs
- open terminal or command line in project directory and type commands below:

```sh
go mod download
go run main.go
```

## Configuration

Create a `config.json` file in the project root (copy from `config.example.json`).

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `links` | array of strings | Yes | URLs to scrape data from |
| `port` | number | Yes | Server port (1-65535) |
| `ip` | string | No | Server IP address (defaults to "localhost") |
| `domains` | array of strings | No | Allowed CORS domains (empty allows all) |
| `with_eq` | boolean | No | If true, scrape from page with equal sign (default: false) |
| `filter_lines` | array of numbers | No | Filter scraped data by keeping only specified line numbers (1-based indexing). For example, `[1, 3, 5]` keeps only the 1st, 3rd, and 5th lines from the scraped data. Empty array keeps all lines |
| `add_lines` | array of strings | No | Additional values to add as lines |
| `add_sum` | boolean | No | If true, adds sum of all numeric values |
| `sum_symbols` | string | No | Symbol(s) to append after sum (e.g., "$") |
| `update_interval` | number | No | File update interval in milliseconds (default: 1000) |
| `write_to_csv` | boolean | No | If true, writes data to CSV file |
| `csv_path` | string | If write_to_csv | Path to CSV output file |
| `write_to_txt` | boolean | No | If true, writes data to TXT file (for Textus Live) |
| `txt_path` | string | If write_to_txt | Path to TXT output file |
| `dataset_name` | string | No | Dataset name for TXT output |

## Info

All breaking errors can be found in info.log file.

To create windows executable program run:

```sh
go build -o PollScraper.exe main.go
```
