package weather

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paulgreig/weathertowalk/internal/location"
)

// newClientForTest allows injecting base URL and HTTP client.
func newClientForTest(apiKey, baseURL string, httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		apiKey:     apiKey,
		baseURL:    baseURL,
	}
}

func TestClientFetchWeather_Success(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/data/2.5/weather", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("appid") == "" {
			t.Errorf("missing appid query param")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleCurrentJSON))
	})

	mux.HandleFunc("/data/2.5/forecast", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("appid") == "" {
			t.Errorf("missing appid query param")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleForecastJSON))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := newClientForTest("test-key", srv.URL+"/data/2.5", srv.Client())
	loc := location.Location{Lat: 51.5, Lon: -0.1}

	data, err := client.FetchWeather(loc)
	if err != nil {
		t.Fatalf("FetchWeather() err = %v, want nil", err)
	}
	if data.Location != loc {
		t.Errorf("Location = %+v, want %+v", data.Location, loc)
	}
	if data.Current.Conditions != "Clear" {
		t.Errorf("Current.Conditions = %q, want Clear", data.Current.Conditions)
	}
	if len(data.Forecast) != 2 {
		t.Fatalf("len(Forecast) = %d, want 2", len(data.Forecast))
	}
	if data.Timestamp.IsZero() {
		t.Error("Timestamp is zero, want non-zero")
	}
}

func TestClientFetchWeather_HTTPError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/data/2.5/weather", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	})
	mux.HandleFunc("/data/2.5/forecast", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(sampleForecastJSON))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := newClientForTest("test-key", srv.URL+"/data/2.5", srv.Client())
	loc := location.Location{Lat: 51.5, Lon: -0.1}

	_, err := client.FetchWeather(loc)
	if err == nil {
		t.Fatal("FetchWeather() err = nil, want error when /weather returns 400")
	}
}

func TestClientFetchWeather_ForecastError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/data/2.5/weather", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(sampleCurrentJSON))
	})
	mux.HandleFunc("/data/2.5/forecast", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := newClientForTest("test-key", srv.URL+"/data/2.5", srv.Client())
	loc := location.Location{Lat: 51.5, Lon: -0.1}

	_, err := client.FetchWeather(loc)
	if err == nil {
		t.Fatal("FetchWeather() err = nil, want error when /forecast returns 400")
	}
}

func TestClient_UsesReasonableTimeout(t *testing.T) {
	// Very basic check: ensure default HTTP client timeout is not zero.
	c := NewClient("test-key", "http://example.com")
	httpClient := c.httpClient
	if httpClient.Timeout == 0 {
		t.Error("expected non-zero HTTP client timeout")
	}
}
