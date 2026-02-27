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

