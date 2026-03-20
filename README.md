# Weather to Walk

A Go application for determining optimal walking conditions based on weather data.

## Prerequisites

- Go 1.25.6 or later
- golangci-lint (for linting)

## Setup

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Install golangci-lint (if not already installed):
   ```bash
   # macOS
   brew install golangci-lint
   
   # Or using the install script
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
   ```

3. Configuration (Phase 1.1):
   - Set `WEATHER_API_KEY` (32-char OpenWeatherMap key), or use a config file.
   - Copy `config.example.json` and `preferences.example.json` as needed.
   - See `internal/config` for loading from env or file.

## Development Workflow

This project follows **Test Driven Development (TDD)**. All code changes must:

1. Write the test first (it should fail)
2. Run tests: `go test ./...`
3. Write minimal implementation to make tests pass
4. Refactor if needed
5. Ensure all tests pass

## Code Quality

Before committing, ensure:

```bash
# Format code
go fmt ./...

# Run tests
go test ./...

# Run linter
golangci-lint run ./...
```

All commands must pass before proposing changes.

## Project Structure

```
.
├── cmd/                 # Application entry points
│   └── weathertowalk/   # Main CLI/TUI
├── internal/            # Private application code
│   ├── algorithm/       # Scoring and window optimization
│   ├── config/          # Config and preferences
│   ├── location/        # Location utilities
│   ├── service/         # Shared fetch + recommendation logic (CLI + mobile)
│   ├── ui/              # Terminal UI rendering
│   └── weather/         # Weather models, parser, HTTP client
├── mobile/              # gomobile bindings for Android (see android/README.md)
├── android/             # Android app (APK) with foreground service
└── .cursor/             # Cursor IDE rules and skills
    ├── rules/           # Coding standards and conventions
    └── skills/          # Project-specific skills
```

## Android APK (foreground service)

The same weather logic runs on-device via **[gomobile](https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile)**:

### Dev Container (recommended)

Open the repo in **VS Code / Cursor** with Dev Containers and use **“Reopen in Container”**. The image includes **JDK 17**, **Android SDK (API 34)**, **NDK**, **Go 1.25.6**, **gomobile**, and **golangci-lint**. Then from the repo root:

```bash
make android-apk
```

Output: `android/app/build/outputs/apk/debug/app-debug.apk`.

### Manual setup

1. Install **JDK 17** and the **Android SDK + NDK** (API 34), set `ANDROID_HOME` (and `ANDROID_NDK_HOME` if needed). On Linux you can use `scripts/setup-android-sdk.sh` for a CLI-only SDK install.
2. `go install golang.org/x/mobile/cmd/gomobile@latest` and `gomobile init`.
3. `make android-aar` (writes `android/app/libs/weathertowalk.aar`), then `cd android && ./gradlew assembleDebug`, or run `make android-apk` from the repo root.

The app exposes a **foreground service** that calls the Go `RunWalk` binding on a background thread and shows the text report in the UI. Use a valid OpenWeatherMap API key and coordinates. See `android/README.md` for details.

## Running the app

```bash
export WEATHER_API_KEY=<your_openweathermap_api_key>
export WEATHER_LAT=<latitude>     # optional, default: London
export WEATHER_LON=<longitude>    # optional, default: London

go run ./cmd/weathertowalk
```

The app will:
- Fetch current weather and a short forecast
- Compute walkability recommendations based on preferences
- Render a terminal dashboard
- Prompt you to refresh or quit

## Cursor Rules

This project includes Cursor rules for:
- Go coding standards (formatting, naming, error handling)
- Test Driven Development workflow
- Linting requirements

Rules are automatically applied when working with Go files.
