package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/paulgreig/weathertowalk/internal/algorithm"
	"github.com/paulgreig/weathertowalk/internal/config"
	"github.com/paulgreig/weathertowalk/internal/location"
	"github.com/paulgreig/weathertowalk/internal/ui"
	"github.com/paulgreig/weathertowalk/internal/weather"
)

func main() {
	// Load core configuration (API key & base URL).
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	// Load preferences (future: from file/flags). For now: defaults.
	prefs := config.DefaultPreferences()

	reader := bufio.NewReader(os.Stdin)

	for {
		// Determine location: from env (WEATHER_LAT/WEATHER_LON) or fallback.
		loc := resolveLocationFromEnv()

		client := weather.NewClient(cfg.APIKey, cfg.APIBaseURL)
		data, err := client.FetchWeather(loc)
		if err != nil {
			log.Fatalf("failed to fetch weather: %v", err)
		}

		recs := algorithm.FindOptimalWindows(data.Forecast, prefs)
		out := ui.RenderWeatherSummary(data, recs)
		fmt.Println(out)

		fmt.Print("\n[r]efresh, [q]uit: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "read error:", err)
			return
		}
		cmd := strings.TrimSpace(strings.ToLower(line))
		if cmd == "q" {
			return
		}
		if cmd != "r" && cmd != "" {
			fmt.Println("Unknown command. Type 'r' to refresh or 'q' to quit.")
		}
		// Loop continues on 'r' or empty input.
	}
}

func resolveLocationFromEnv() location.Location {
	latStr := os.Getenv("WEATHER_LAT")
	lonStr := os.Getenv("WEATHER_LON")
	if latStr == "" || lonStr == "" {
		// Attempt automatic IP-based geolocation first.
		if loc, err := location.DetectLocation(); err == nil {
			return loc
		}
		// Fallback example: London coordinates.
		return location.Location{Lat: 51.5074, Lon: -0.1278}
	}
	lat, err1 := strconv.ParseFloat(latStr, 64)
	lon, err2 := strconv.ParseFloat(lonStr, 64)
	if err1 != nil || err2 != nil {
		return location.Location{Lat: 51.5074, Lon: -0.1278}
	}
	if err := location.ValidateCoordinates(lat, lon); err != nil {
		return location.Location{Lat: 51.5074, Lon: -0.1278}
	}
	return location.Location{Lat: lat, Lon: lon}
}
