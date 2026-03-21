package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/paulgreig/weathertowalk/internal/algorithm"
	"github.com/paulgreig/weathertowalk/internal/location"
	"github.com/paulgreig/weathertowalk/internal/weather"
)

func TestRenderWeatherSummary_IncludesKeySections(t *testing.T) {
	loc := location.Location{Lat: 51.5, Lon: -0.1}
	now := time.Now().UTC()
	data := weather.WeatherData{
		Current: weather.CurrentWeather{
			Temperature:   20,
			Humidity:      60,
			Precipitation: 0,
			AirQuality:    42,
			Conditions:    "Clear",
			WindSpeed:     3.5,
		},
		Forecast: []weather.ForecastSlot{
			{Time: now.Add(1 * time.Hour), Temperature: 21, Humidity: 55, Conditions: "Clouds"},
		},
		Location:  loc,
		Timestamp: now,
	}

	recs := []algorithm.Recommendation{
		{
			StartTime:      now.Add(1 * time.Hour),
			EndTime:        now.Add(2 * time.Hour),
			Duration:       time.Hour,
			Score:          90,
			MinSlotScore:   90,
			WeatherSummary: "Clouds",
		},
	}

	out := RenderWeatherSummary(data, recs)

	for _, substr := range []string{
		"Weather to Walk",
		"Current weather",
		"Forecast (next slots)",
		"Recommendations",
		"20.0°C",
		"Humidity   : 60%",
		"AQI        : 42",
		"Score: 90",
	} {
		if !strings.Contains(out, substr) {
			t.Errorf("output missing %q.\nOutput:\n%s", substr, out)
		}
	}
}

func TestRenderWeatherSummary_ExactSectionOrder(t *testing.T) {
	loc := location.Location{Lat: 55.9533, Lon: -3.1883}
	ts := time.Date(2025, 6, 15, 14, 0, 0, 0, time.UTC)
	data := weather.WeatherData{
		Current: weather.CurrentWeather{
			Temperature:   18.5,
			Humidity:      72,
			Precipitation: 0.2,
			AirQuality:    0,
			Conditions:    "Drizzle",
			WindSpeed:     4.1,
		},
		Forecast: []weather.ForecastSlot{
			{Time: ts.Add(3 * time.Hour), Temperature: 19.0, Humidity: 65, Precipitation: 0, Conditions: "Clear"},
			{Time: ts.Add(6 * time.Hour), Temperature: 17.0, Humidity: 70, Precipitation: 0.5, Conditions: "Clouds"},
		},
		Location:  loc,
		Timestamp: ts,
	}
	recs := []algorithm.Recommendation{
		{
			StartTime:      ts.Add(3 * time.Hour),
			EndTime:        ts.Add(3*time.Hour + 30*time.Minute),
			Duration:       30 * time.Minute,
			Score:          85,
			MinSlotScore:   85,
			WeatherSummary: "Clear",
		},
	}

	out := RenderWeatherSummary(data, recs)

	ordered := []string{
		"========================================",
		"Weather to Walk",
		"========================================",
		"Location: 55.9533, -3.1883",
		"As of : 2025-06-15T14:00:00Z",
		"Current weather",
		"Conditions : Drizzle",
		"Temp       : 18.5°C",
		"Humidity   : 72%",
		"Precip     : 0.2 mm",
		"Wind       : 4.1 m/s",
		"Forecast (next slots)",
		"Clear",
		"Clouds",
		"Recommendations",
		"Score: 85",
		"Weather: Clear",
	}

	prev := 0
	for _, s := range ordered {
		idx := strings.Index(out[prev:], s)
		if idx < 0 {
			t.Errorf("output missing or out of order: %q (after position %d).\nOutput:\n%s", s, prev, out)
			continue
		}
		prev += idx + len(s)
	}
}

func TestRenderWeatherSummary_NoAQI_WhenZero(t *testing.T) {
	data := weather.WeatherData{
		Current: weather.CurrentWeather{
			Temperature:   22,
			Humidity:      50,
			Precipitation: 0,
			AirQuality:    0,
			Conditions:    "Sunny",
			WindSpeed:     2.0,
		},
		Location:  location.Location{Lat: 40.0, Lon: -74.0},
		Timestamp: time.Now().UTC(),
	}

	out := RenderWeatherSummary(data, nil)

	if strings.Contains(out, "AQI") {
		t.Errorf("AQI should not appear when AirQuality is 0.\nOutput:\n%s", out)
	}
}

func TestRenderWeatherSummary_AQI_WhenNonZero(t *testing.T) {
	data := weather.WeatherData{
		Current: weather.CurrentWeather{
			Temperature:   22,
			Humidity:      50,
			Precipitation: 0,
			AirQuality:    55,
			Conditions:    "Sunny",
			WindSpeed:     2.0,
		},
		Location:  location.Location{Lat: 40.0, Lon: -74.0},
		Timestamp: time.Now().UTC(),
	}

	out := RenderWeatherSummary(data, nil)

	if !strings.Contains(out, "AQI        : 55") {
		t.Errorf("expected AQI line for non-zero air quality.\nOutput:\n%s", out)
	}
}

func TestRenderWeatherSummary_NoForecastData(t *testing.T) {
	data := weather.WeatherData{
		Current: weather.CurrentWeather{
			Temperature: 15,
			Humidity:    80,
			Conditions:  "Rain",
			WindSpeed:   6.0,
		},
		Location:  location.Location{Lat: 48.8566, Lon: 2.3522},
		Timestamp: time.Now().UTC(),
	}

	out := RenderWeatherSummary(data, nil)

	if !strings.Contains(out, "(no forecast data)") {
		t.Errorf("expected '(no forecast data)' when forecast is empty.\nOutput:\n%s", out)
	}
	if !strings.Contains(out, "No ideal walking windows found") {
		t.Errorf("expected no-recommendations message.\nOutput:\n%s", out)
	}
}

func TestRenderWeatherSummary_MultipleRecommendations(t *testing.T) {
	ts := time.Date(2025, 7, 1, 10, 0, 0, 0, time.UTC)
	data := weather.WeatherData{
		Current: weather.CurrentWeather{
			Temperature: 22,
			Humidity:    55,
			Conditions:  "Clear",
			WindSpeed:   2.5,
		},
		Forecast: []weather.ForecastSlot{
			{Time: ts.Add(3 * time.Hour), Temperature: 23, Humidity: 50, Conditions: "Clear"},
			{Time: ts.Add(6 * time.Hour), Temperature: 24, Humidity: 45, Conditions: "Clear"},
		},
		Location:  location.Location{Lat: 34.0522, Lon: -118.2437},
		Timestamp: ts,
	}
	recs := []algorithm.Recommendation{
		{
			StartTime:      ts.Add(3 * time.Hour),
			EndTime:        ts.Add(3*time.Hour + 30*time.Minute),
			Duration:       30 * time.Minute,
			Score:          95,
			MinSlotScore:   95,
			WeatherSummary: "Clear",
		},
		{
			StartTime:      ts.Add(6 * time.Hour),
			EndTime:        ts.Add(6*time.Hour + 30*time.Minute),
			Duration:       30 * time.Minute,
			Score:          88,
			MinSlotScore:   88,
			WeatherSummary: "Clear",
		},
	}

	out := RenderWeatherSummary(data, recs)

	if !strings.Contains(out, "Score: 95") {
		t.Errorf("missing first recommendation score.\nOutput:\n%s", out)
	}
	if !strings.Contains(out, "Score: 88") {
		t.Errorf("missing second recommendation score.\nOutput:\n%s", out)
	}

	firstIdx := strings.Index(out, "Score: 95")
	secondIdx := strings.Index(out, "Score: 88")
	if firstIdx >= secondIdx {
		t.Errorf("recommendations should appear in order; score 95 at %d, score 88 at %d", firstIdx, secondIdx)
	}
}

func TestRenderWeatherSummary_DeterministicOutput(t *testing.T) {
	ts := time.Date(2025, 3, 20, 12, 0, 0, 0, time.UTC)
	data := weather.WeatherData{
		Current: weather.CurrentWeather{
			Temperature:   18,
			Humidity:      65,
			Precipitation: 0,
			AirQuality:    30,
			Conditions:    "Clouds",
			WindSpeed:     3.0,
		},
		Forecast: []weather.ForecastSlot{
			{Time: ts.Add(3 * time.Hour), Temperature: 19, Humidity: 60, Conditions: "Clear"},
		},
		Location:  location.Location{Lat: 51.5074, Lon: -0.1278},
		Timestamp: ts,
	}
	recs := []algorithm.Recommendation{
		{
			StartTime:      ts.Add(3 * time.Hour),
			EndTime:        ts.Add(3*time.Hour + 30*time.Minute),
			Duration:       30 * time.Minute,
			Score:          80,
			MinSlotScore:   80,
			WeatherSummary: "Clear",
		},
	}

	out1 := RenderWeatherSummary(data, recs)
	out2 := RenderWeatherSummary(data, recs)
	if out1 != out2 {
		t.Errorf("RenderWeatherSummary is not deterministic.\nFirst:\n%s\nSecond:\n%s", out1, out2)
	}
}

func TestRenderWeatherSummary_ForecastLimitsFourSlots(t *testing.T) {
	ts := time.Date(2025, 8, 1, 8, 0, 0, 0, time.UTC)
	slots := make([]weather.ForecastSlot, 6)
	for i := range slots {
		slots[i] = weather.ForecastSlot{
			Time:        ts.Add(time.Duration(i+1) * 3 * time.Hour),
			Temperature: float64(20 + i),
			Humidity:    50,
			Conditions:  "Clear",
		}
	}
	data := weather.WeatherData{
		Current: weather.CurrentWeather{
			Temperature: 19,
			Humidity:    55,
			Conditions:  "Clear",
			WindSpeed:   2.0,
		},
		Forecast:  slots,
		Location:  location.Location{Lat: 51.5, Lon: -0.1},
		Timestamp: ts,
	}

	out := RenderWeatherSummary(data, nil)

	if !strings.Contains(out, "20.0°C") {
		t.Error("missing slot 1 temperature")
	}
	if !strings.Contains(out, "23.0°C") {
		t.Error("missing slot 4 temperature")
	}
	if strings.Contains(out, "24.0°C") {
		t.Error("slot 5 should not appear (max 4 slots)")
	}
	if strings.Contains(out, "25.0°C") {
		t.Error("slot 6 should not appear (max 4 slots)")
	}
}
