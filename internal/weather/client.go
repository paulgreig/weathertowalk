package weather

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/paulgreig/weathertowalk/internal/location"
)

// Client calls the weather API to fetch current conditions and forecast.
type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

// NewClient creates a Client with a sensible default HTTP timeout.
// baseURL should be something like "https://api.openweathermap.org/data/2.5".
func NewClient(apiKey, baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
			},
		},
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// FetchWeather fetches current weather and forecast for the given location.
func (c *Client) FetchWeather(loc location.Location) (WeatherData, error) {
	cur, curLoc, err := c.fetchCurrent(loc)
	if err != nil {
		return WeatherData{}, err
	}
	forecast, fcLoc, err := c.fetchForecast(loc)
	if err != nil {
		return WeatherData{}, err
	}

	// Prefer forecast location if provided, otherwise current.
	finalLoc := curLoc
	if fcLoc != (location.Location{}) {
		finalLoc = fcLoc
	}
	return ToWeatherData(cur, forecast, finalLoc, time.Now().UTC()), nil
}

func (c *Client) fetchCurrent(loc location.Location) (CurrentWeather, location.Location, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return CurrentWeather{}, location.Location{}, fmt.Errorf("parse base url: %w", err)
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/weather"
	q := u.Query()
	q.Set("lat", fmt.Sprintf("%f", loc.Lat))
	q.Set("lon", fmt.Sprintf("%f", loc.Lon))
	q.Set("appid", c.apiKey)
	u.RawQuery = q.Encode()

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return CurrentWeather{}, location.Location{}, fmt.Errorf("GET /weather: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return CurrentWeather{}, location.Location{}, fmt.Errorf("GET /weather: status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return CurrentWeather{}, location.Location{}, fmt.Errorf("read /weather body: %w", err)
	}
	return ParseCurrent(data)
}

func (c *Client) fetchForecast(loc location.Location) ([]ForecastSlot, location.Location, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, location.Location{}, fmt.Errorf("parse base url: %w", err)
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/forecast"
	q := u.Query()
	q.Set("lat", fmt.Sprintf("%f", loc.Lat))
	q.Set("lon", fmt.Sprintf("%f", loc.Lon))
	q.Set("appid", c.apiKey)
	u.RawQuery = q.Encode()

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, location.Location{}, fmt.Errorf("GET /forecast: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, location.Location{}, fmt.Errorf("GET /forecast: status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, location.Location{}, fmt.Errorf("read /forecast body: %w", err)
	}
	return ParseForecast(data)
}
