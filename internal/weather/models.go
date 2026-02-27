package weather

import (
	"time"

	"github.com/paulgreig/weathertowalk/internal/location"
)

// WeatherData is the aggregate weather information used by the app.
type WeatherData struct {
	Current   CurrentWeather
	Forecast  []ForecastSlot
	Location  location.Location
	Timestamp time.Time
}

// CurrentWeather represents current conditions at a location.
type CurrentWeather struct {
	Temperature   float64
	Humidity      float64
	Precipitation float64
	AirQuality    int
	Conditions    string
	WindSpeed     float64
}

// ForecastSlot is a single forecast entry, typically a 3-hour window.
type ForecastSlot struct {
	Time             time.Time
	Temperature      float64
	Humidity         float64
	Precipitation    float64
	AirQuality       int
	Conditions       string
	WalkabilityScore int
}

