package mobile

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/paulgreig/weathertowalk/internal/service"
)

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

func newMockServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/data/2.5/weather", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleCurrentJSON))
	})
	mux.HandleFunc("/data/2.5/forecast", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleForecastJSON))
	})
	return httptest.NewServer(mux)
}

func TestRunWalk_OutputMatchesService(t *testing.T) {
	srv := newMockServer()
	defer srv.Close()
	baseURL := srv.URL + "/data/2.5"

	mobileOut, err := RunWalk(testAPIKey, baseURL, 51.5, -0.1)
	if err != nil {
		t.Fatalf("RunWalk() err = %v, want nil", err)
	}

	serviceOut, err := service.RunWalkReport(testAPIKey, baseURL, 51.5, -0.1)
	if err != nil {
		t.Fatalf("service.RunWalkReport() err = %v, want nil", err)
	}

	if mobileOut != serviceOut {
		t.Errorf("output mismatch between mobile.RunWalk and service.RunWalkReport.\nmobile:\n%s\nservice:\n%s", mobileOut, serviceOut)
	}
}

func TestRunWalk_ContainsAllSections(t *testing.T) {
	srv := newMockServer()
	defer srv.Close()

	out, err := RunWalk(testAPIKey, srv.URL+"/data/2.5", 51.5, -0.1)
	if err != nil {
		t.Fatalf("RunWalk() err = %v, want nil", err)
	}

	sections := []string{
		"========================================",
		"Weather to Walk",
		"Location:",
		"As of :",
		"Current weather",
		"Conditions :",
		"Temp       :",
		"Humidity   :",
		"Precip     :",
		"Wind       :",
		"Forecast (next slots)",
		"Recommendations",
	}
	for _, s := range sections {
		if !strings.Contains(out, s) {
			t.Errorf("output missing section %q.\nOutput:\n%s", s, out)
		}
	}
}

func TestRunWalk_ContainsWeatherValues(t *testing.T) {
	srv := newMockServer()
	defer srv.Close()

	out, err := RunWalk(testAPIKey, srv.URL+"/data/2.5", 51.5, -0.1)
	if err != nil {
		t.Fatalf("RunWalk() err = %v, want nil", err)
	}

	values := []string{
		"51.5000, -0.1000",
		"Clear",
		"20.0°C",
		"60%",
		"0.0 mm",
		"3.2 m/s",
		"Clouds",
	}
	for _, v := range values {
		if !strings.Contains(out, v) {
			t.Errorf("output missing value %q.\nOutput:\n%s", v, out)
		}
	}
}

func TestRunWalk_InvalidAPIKey(t *testing.T) {
	_, err := RunWalk("bad", "http://example.com/data/2.5", 51.5, -0.1)
	if err == nil {
		t.Fatal("RunWalk() err = nil, want error for invalid API key")
	}
}

func TestRunWalk_InvalidCoordinates(t *testing.T) {
	_, err := RunWalk(testAPIKey, "http://example.com/data/2.5", 200, 0)
	if err == nil {
		t.Fatal("RunWalk() err = nil, want error for invalid latitude")
	}
}

func TestRunWalk_EmptyBaseURL_UsesDefault(t *testing.T) {
	_, err := RunWalk(testAPIKey, "", 51.5, -0.1)
	// The call will fail (no real API server) but should not panic.
	// The error should be a network error, not a configuration error.
	if err == nil {
		t.Skip("unexpectedly succeeded with default URL; may have network access")
	}
	if strings.Contains(err.Error(), "invalid API key") {
		t.Errorf("got API key error, expected network/fetch error: %v", err)
	}
}
