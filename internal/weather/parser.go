package weather

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/paulgreig/weathertowalk/internal/location"
)

// OpenWeatherMap current weather response (partial).
type owmCurrentResponse struct {
	Coord struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coord"`
	DT      int64 `json:"dt"`
	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity float64 `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Rain struct {
		OneH float64 `json:"1h"`
	} `json:"rain"`
}

// OpenWeatherMap forecast response (partial).
type owmForecastResponse struct {
	List []struct {
		DT      int64 `json:"dt"`
		Weather []struct {
			Main string `json:"main"`
		} `json:"weather"`
		Main struct {
			Temp     float64 `json:"temp"`
			Humidity float64 `json:"humidity"`
		} `json:"main"`
		Pop  float64 `json:"pop"`
		Rain struct {
			ThreeH float64 `json:"3h"`
		} `json:"rain"`
	} `json:"list"`
	City struct {
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
	} `json:"city"`
}

// ParseCurrent parses a current weather JSON payload into a CurrentWeather and location.
// Temperatures are converted from Kelvin to Celsius.
func ParseCurrent(data []byte) (CurrentWeather, location.Location, error) {
	var resp owmCurrentResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return CurrentWeather{}, location.Location{}, fmt.Errorf("parse current weather: %w", err)
	}
	loc := location.Location{Lat: resp.Coord.Lat, Lon: resp.Coord.Lon}
	cur := CurrentWeather{
		Temperature:   kelvinToCelsius(resp.Main.Temp),
		Humidity:      resp.Main.Humidity,
		Precipitation: resp.Rain.OneH,
		AirQuality:    0, // filled when AQ endpoint integrated
		Conditions:    firstWeatherMain(resp.Weather),
		WindSpeed:     resp.Wind.Speed,
	}
	return cur, loc, nil
}

// ParseForecast parses a forecast JSON payload into forecast slots and location.
// Temperatures are converted from Kelvin to Celsius.
func ParseForecast(data []byte) ([]ForecastSlot, location.Location, error) {
	var resp owmForecastResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, location.Location{}, fmt.Errorf("parse forecast: %w", err)
	}
	loc := location.Location{Lat: resp.City.Coord.Lat, Lon: resp.City.Coord.Lon}
	slots := make([]ForecastSlot, 0, len(resp.List))
	for _, item := range resp.List {
		slot := ForecastSlot{
			Time:          time.Unix(item.DT, 0),
			Temperature:   kelvinToCelsius(item.Main.Temp),
			Humidity:      item.Main.Humidity,
			Precipitation: item.Rain.ThreeH,
			AirQuality:    0,
			Conditions:    firstWeatherMain(item.Weather),
		}
		slots = append(slots, slot)
	}
	return slots, loc, nil
}

// ToWeatherData combines the parsed current and forecast data into a WeatherData struct.
func ToWeatherData(cur CurrentWeather, forecast []ForecastSlot, loc location.Location, ts time.Time) WeatherData {
	return WeatherData{
		Current:   cur,
		Forecast:  forecast,
		Location:  loc,
		Timestamp: ts,
	}
}

func kelvinToCelsius(k float64) float64 {
	return k - 273.15
}

func firstWeatherMain(list []struct{ Main string `json:"main"` }) string {
	if len(list) == 0 {
		return ""
	}
	return list[0].Main
}

