# Agents

## Cursor Cloud specific instructions

### Project overview

Weather to Walk is a Go CLI application that fetches weather data from OpenWeatherMap and recommends optimal walking times. It has zero external Go dependencies (stdlib only). See `README.md` for full project structure and `QUICKSTART.md` for usage.

### Dev commands

All standard commands are in the `Makefile`: `make test`, `make fmt`, `make lint`, `make check` (runs all), `make test-coverage`. See `README.md` > Code Quality for details.

### Toolchain

- **Go 1.25.6** — installed via the Go toolchain wrapper; first invocation may download the toolchain.
- **golangci-lint v2** — installed at `~/go/bin/golangci-lint`. Ensure `$HOME/go/bin` is on `PATH` before running `make lint`.

### Gotchas

- The `.golangci.yml` was migrated to golangci-lint v2 format (requires `version: "2"` header, `gofmt`/`goimports` moved to `formatters` section, removed deprecated linters `golint`/`varcheck`/`gosimple`). Do not downgrade to v1.
- `make lint` returns exit code 1 due to pre-existing lint issues in the codebase (errcheck, gocritic, gosec, revive). These are not regressions.
- The app requires `WEATHER_API_KEY` (32-char hex, from OpenWeatherMap free tier). Without it, `config.Load()` fails immediately on startup.
- To run the app locally without a real API key, start a mock HTTP server and set `WEATHER_API_URL` to override the base URL (e.g. `http://127.0.0.1:8999/data/2.5`). Set `WEATHER_LAT`/`WEATHER_LON` to skip IP geolocation.
- There are no databases, Docker containers, or background services. The app is a standalone CLI that makes outbound HTTP calls only.
