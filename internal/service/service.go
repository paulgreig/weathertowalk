// Package service provides shared weather fetch and recommendation logic used by
// the CLI and the gomobile Android binding.
package service

import (
	"fmt"

	"github.com/paulgreig/weathertowalk/internal/algorithm"
	"github.com/paulgreig/weathertowalk/internal/config"
	"github.com/paulgreig/weathertowalk/internal/location"
	"github.com/paulgreig/weathertowalk/internal/ui"
	"github.com/paulgreig/weathertowalk/internal/weather"
)

const defaultAPIBaseURL = "https://api.openweathermap.org/data/2.5"

// RunWalkReport fetches weather for the given coordinates and returns the same
// text summary used by the CLI (terminal UI).
func RunWalkReport(apiKey, baseURL string, lat, lon float64) (string, error) {
	if !config.ValidateAPIKey(apiKey) {
		return "", fmt.Errorf("invalid API key: must be 32 hex characters")
	}
	if baseURL == "" {
		baseURL = defaultAPIBaseURL
	}
	loc, err := location.NewLocation(lat, lon)
	if err != nil {
		return "", err
	}
	client := weather.NewClient(apiKey, baseURL)
	data, err := client.FetchWeather(loc)
	if err != nil {
		return "", err
	}
	prefs := config.DefaultPreferences()
	recs := algorithm.FindOptimalWindows(data.Forecast, prefs)
	return ui.RenderWeatherSummary(data, recs), nil
}
