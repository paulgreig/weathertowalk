package service

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Valid 32-char hex OpenWeatherMap-style API key for tests.
const testAPIKey = "0123456789abcdef0123456789abcdef"

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

func TestRunWalkReport_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/data/2.5/weather", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleCurrentJSON))
	})
	mux.HandleFunc("/data/2.5/forecast", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleForecastJSON))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	out, err := RunWalkReport(testAPIKey, srv.URL+"/data/2.5", 51.5, -0.1)
	if err != nil {
		t.Fatalf("RunWalkReport() err = %v, want nil", err)
	}
	if !strings.Contains(out, "Weather to Walk") {
		t.Errorf("output missing title, got: %q", out)
	}
	if !strings.Contains(out, "Recommendations") {
		t.Errorf("output missing recommendations section")
	}
}

func TestRunWalkReport_InvalidAPIKey(t *testing.T) {
	_, err := RunWalkReport("short", "http://example.com/data/2.5", 51.5, -0.1)
	if err == nil {
		t.Fatal("RunWalkReport() err = nil, want error for invalid API key")
	}
}

func TestRunWalkReport_InvalidCoordinates(t *testing.T) {
	_, err := RunWalkReport(testAPIKey, "http://example.com/data/2.5", 200, 0)
	if err == nil {
		t.Fatal("RunWalkReport() err = nil, want error for invalid latitude")
	}
}
