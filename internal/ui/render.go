package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/paulgreig/weathertowalk/internal/algorithm"
	"github.com/paulgreig/weathertowalk/internal/weather"
)

// RenderWeatherSummary builds a textual summary of current weather, forecast, and recommendations.
// This is the primary terminal UI for the initial version of the app.
func RenderWeatherSummary(data weather.WeatherData, recs []algorithm.Recommendation) string {
	var b strings.Builder

	fmt.Fprintln(&b, "========================================")
	fmt.Fprintln(&b, "          Weather to Walk")
	fmt.Fprintln(&b, "========================================")
	fmt.Fprintf(&b, "Location: %.4f, %.4f\n", data.Location.Lat, data.Location.Lon)
	fmt.Fprintf(&b, "As of : %s\n\n", data.Timestamp.UTC().Format(time.RFC3339))

	// Current weather
	fmt.Fprintln(&b, "Current weather")
	fmt.Fprintln(&b, "----------------")
	fmt.Fprintf(&b, "Conditions : %s\n", data.Current.Conditions)
	fmt.Fprintf(&b, "Temp       : %.1f°C\n", data.Current.Temperature)
	fmt.Fprintf(&b, "Humidity   : %.0f%%\n", data.Current.Humidity)
	fmt.Fprintf(&b, "Precip     : %.1f mm\n", data.Current.Precipitation)
	if data.Current.AirQuality > 0 {
		fmt.Fprintf(&b, "AQI        : %d\n", data.Current.AirQuality)
	}
	fmt.Fprintf(&b, "Wind       : %.1f m/s\n", data.Current.WindSpeed)
	fmt.Fprintln(&b)

	// Forecast (show up to 4 slots)
	fmt.Fprintln(&b, "Forecast (next slots)")
	fmt.Fprintln(&b, "---------------------")
	maxSlots := 4
	if len(data.Forecast) < maxSlots {
		maxSlots = len(data.Forecast)
	}
	for i := 0; i < maxSlots; i++ {
		slot := data.Forecast[i]
		fmt.Fprintf(
			&b,
			"%s | %-12s | %5.1f°C | %3.0f%% humidity | %.1f mm\n",
			slot.Time.Format("15:04"),
			slot.Conditions,
			slot.Temperature,
			slot.Humidity,
			slot.Precipitation,
		)
	}
	if maxSlots == 0 {
		fmt.Fprintln(&b, "(no forecast data)")
	}
	fmt.Fprintln(&b)

	// Recommendations
	fmt.Fprintln(&b, "Recommendations")
	fmt.Fprintln(&b, "----------------")
	if len(recs) == 0 {
		fmt.Fprintln(&b, "No ideal walking windows found based on current preferences.")
	} else {
		for _, r := range recs {
			fmt.Fprintf(
				&b,
				"%s – %s | Duration: %dm | Score: %d | Weather: %s\n",
				r.StartTime.Format("15:04"),
				r.EndTime.Format("15:04"),
				int(r.Duration.Minutes()),
				r.Score,
				r.WeatherSummary,
			)
		}
	}

	return b.String()
}
