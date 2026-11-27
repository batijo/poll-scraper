# Poll Scraper

A desktop application designed to collect and structure data from the [savitarna.mobi.lt](https://savitarna.mobi.lt) website with the purpose of enabling its integration into various graphic engines.

---

## Features

### Scraper Control
Start and stop the scraper from the UI. The scraper fetches data from configured URLs at a set interval and writes results to output files. Use **Preview** to fetch data once without starting the continuous scraper.

### Link Login
Scrape pages that sit behind a login. Add an optional username and password to a URL in the add or edit dialog. When the site redirects to a login page, the scraper signs in, keeps the session cookies, and retries the request. Enter both fields or leave both empty. The same sign-in is used by Preview, the scraper, and the HTTP API. If the credentials are rejected, an error is logged and the URL returns no data.

### Filter Lines
Choose which scraped lines to include in the output. The filter modal shows all lines grouped by their source URL, so you can easily identify where each line comes from. Filters use stable indices tied to URL data only, so adding or removing custom lines won't shift your selections.

### Symbol Lines
Add a symbol (for example ` eur` or `%`) to any output line, scraped or custom. Set it per line in the filter modal. The app inserts a copy of that line right after the original, named `<name>_symbol`, with the symbol appended to the value. Only the original line counts toward the sum.

In `config.json`, `number` is the line's position in the output after filtering and custom lines, starting at 1. The inserted copies and the sum lines do not count toward this position:

```json
"symbol_lines": [{"number": 3, "symbol": " eur"}]
```

If you change the filters or hide a custom line, the positions shift. Saving the filter modal again recalculates them.

### Custom Lines
Add extra data rows with a name and value that get appended after the scraped data. Each custom line has a visibility toggle — hidden lines are excluded from output without deleting them.

### Sum
Automatically calculates the sum of all numeric values in the output. Optionally append a symbol (e.g. `$`, `€`) after the sum value. Adds one or two extra lines at the end: `sum` and optionally `sum_symbol`.

### Line Count Protection
If any URL starts returning a different number of lines than the first scrape, an error is logged and the scraper status becomes faulted. Optionally enable **Stop on URL line count change** in settings to automatically stop the scraper when this happens.

---

## Output

### TXT
Writes data in a structured format for Textus Live. Choose between **ANSI** (Windows-1252) or **UTF-8** encoding.

### CSV
Writes name/value pairs as CSV rows.

### HTTP API
JSON API. Any client can fetch the current data as a JSON array from the root endpoint. CORS domains can be restricted in settings.

---

## UI

### Status Tab
Shows scraper state, URL health (which URLs are returning data), line counts, output targets, and server status. Includes a live log viewer with level filtering (Error, Warn, Info, Debug) and color coding.

### Settings Tab
Update interval, debug mode, line count protection, server configuration, and output paths with encoding selection.

### Scraping Tab
URL list with optional login credentials, data parsing options, custom lines management, and line filter configuration.

### Status Indicator
| Color | Meaning |
|-------|---------|
| Green | Scraper is running, all URLs healthy |
| Gray  | Scraper is stopped |
| Red   | Error — URL failures, write errors, or line count changes |

---

## Setup

1. Install [Go](https://golang.org/dl/) 1.26+ and [Wails](https://wails.io/docs/gettingstarted/installation) v2.17+
2. Run: `wails dev`

A default `config.json` is created automatically on first launch if one doesn't exist.

---

## Build

```sh
wails build
```

---

## Notes

- All settings are configurable through the UI — the `config.json` file is managed automatically
- Errors are logged to `error.log`
- Link passwords are stored in plain text in `config.json`
- `links` in `config.json` are objects with a `url` and optional `name` and `password`, as in `config.example.json`. A config from before 2.0 lists links as plain strings and fails to load until you change each entry to `{"url": "..."}`
