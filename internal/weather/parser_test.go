package weather

import (
	"testing"
	"time"

	"github.com/paulgreig/weathertowalk/internal/location"
)

const sampleCurrentJSON = `{
  "coord": {"lat": 51.5, "lon": -0.1},
  "dt": 1700000000,
  "weather": [{"main": "Clear"}],
  "main": {
    "temp": 293.15,
    "humidity": 60
  },
  "wind": {"speed": 3.2},
  "rain": {"1h": 0.0}
}`

const sampleForecastJSON = `{
  "list": [
    {
      "dt": 1700003600,
      "weather": [{"main": "Clouds"}],
      "main": {
        "temp": 294.15,
        "humidity": 58
      },
      "pop": 0.1,
      "rain": {"3h": 0.1}
    },
    {
      "dt": 1700007200,
      "weather": [{"main": "Rain"}],
      "main": {
        "temp": 291.15,
        "humidity": 80
      },
      "pop": 0.8,
      "rain": {"3h": 2.0}
    }
  ],
  "city": {
    "coord": {"lat": 51.5, "lon": -0.1}
  }
}`

func TestParseCurrentWeather(t *testing.T) {
	cur, loc, err := ParseCurrent([]byte(sampleCurrentJSON))
	if err != nil {
		t.Fatalf("ParseCurrent() err = %v, want nil", err)
	}
	if loc != (location.Location{Lat: 51.5, Lon: -0.1}) {
		t.Errorf("Location = %+v, want 51.5,-0.1", loc)
	}
	if cur.Conditions != "Clear" {
		t.Errorf("Conditions = %q, want %q", cur.Conditions, "Clear")
	}
	// temp is converted from Kelvin to Celsius
	if cur.Temperature < 19 || cur.Temperature > 21 {
		t.Errorf("Temperature = %v, want around 20C", cur.Temperature)
	}
	if cur.Humidity != 60 {
		t.Errorf("Humidity = %v, want 60", cur.Humidity)
	}
	if cur.Precipitation != 0 {
		t.Errorf("Precipitation = %v, want 0", cur.Precipitation)
	}
	if cur.WindSpeed != 3.2 {
		t.Errorf("WindSpeed = %v, want 3.2", cur.WindSpeed)
	}
}

func TestParseForecast(t *testing.T) {
	slots, loc, err := ParseForecast([]byte(sampleForecastJSON))
	if err != nil {
		t.Fatalf("ParseForecast() err = %v, want nil", err)
	}
	if loc != (location.Location{Lat: 51.5, Lon: -0.1}) {
		t.Errorf("Location = %+v, want 51.5,-0.1", loc)
	}
	if len(slots) != 2 {
		t.Fatalf("len(slots) = %d, want 2", len(slots))
	}
	if slots[0].Conditions != "Clouds" {
		t.Errorf("slot[0].Conditions = %q, want Clouds", slots[0].Conditions)
	}
	if slots[1].Conditions != "Rain" {
		t.Errorf("slot[1].Conditions = %q, want Rain", slots[1].Conditions)
	}
	if slots[0].Time.After(slots[1].Time) {
		t.Errorf("slots not in chronological order: %v after %v", slots[0].Time, slots[1].Time)
	}
}

func TestToWeatherDataCombinesCurrentAndForecast(t *testing.T) {
	cur := CurrentWeather{
		Temperature:   20,
		Humidity:      60,
		Precipitation: 0,
		AirQuality:    50,
		Conditions:    "Clear",
		WindSpeed:     3,
	}
	loc := location.Location{Lat: 51.5, Lon: -0.1}
	now := time.Unix(1700000000, 0)
	forecast := []ForecastSlot{
		{Time: now.Add(1 * time.Hour), Temperature: 21},
	}

	data := ToWeatherData(cur, forecast, loc, now)
	if data.Location != loc {
		t.Errorf("Location = %+v, want %+v", data.Location, loc)
	}
	if len(data.Forecast) != 1 {
		t.Fatalf("Forecast len = %d, want 1", len(data.Forecast))
	}
	if !data.Timestamp.Equal(now) {
		t.Errorf("Timestamp = %v, want %v", data.Timestamp, now)
	}
}
