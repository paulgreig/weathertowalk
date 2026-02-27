package weather

import (
	"testing"
	"time"

	"github.com/paulgreig/weathertowalk/internal/location"
)

func TestWeatherData_ModelShapes(t *testing.T) {
	loc := location.Location{Lat: 51.5, Lon: -0.1}
	now := time.Now()

	data := WeatherData{
		Current: CurrentWeather{
			Temperature:   20.5,
			Humidity:      60,
			Precipitation: 0.0,
			AirQuality:    42,
			Conditions:    "clear",
			WindSpeed:     3.2,
		},
		Forecast: []ForecastSlot{
			{
				Time:            now.Add(1 * time.Hour),
				Temperature:     21.0,
				Humidity:        58,
				Precipitation:   0.1,
				AirQuality:      45,
				Conditions:      "partly cloudy",
				WalkabilityScore: 0,
			},
		},
		Location:  loc,
		Timestamp: now,
	}

	if data.Location.Lat != loc.Lat || data.Location.Lon != loc.Lon {
		t.Errorf("Location mismatch: got %+v, want %+v", data.Location, loc)
	}
	if len(data.Forecast) != 1 {
		t.Fatalf("Forecast length = %d, want 1", len(data.Forecast))
	}
}

