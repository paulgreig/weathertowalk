package location

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// geoBaseURL is the IP geolocation endpoint. Overridable in tests.
var geoBaseURL = "http://ip-api.com/json/"

type ipAPIResponse struct {
	Status  string  `json:"status"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Message string  `json:"message"`
}

// DetectLocation attempts to determine the user's approximate location based on IP.
// It returns a validated Location or an error if detection fails.
func DetectLocation() (Location, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(geoBaseURL)
	if err != nil {
		return Location{}, fmt.Errorf("detect location: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Location{}, fmt.Errorf("detect location: status %d", resp.StatusCode)
	}

	var parsed ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return Location{}, fmt.Errorf("detect location: parse response: %w", err)
	}
	if parsed.Status != "success" {
		return Location{}, fmt.Errorf("detect location: %s", parsed.Message)
	}

	if err := ValidateCoordinates(parsed.Lat, parsed.Lon); err != nil {
		return Location{}, err
	}
	return Location{Lat: parsed.Lat, Lon: parsed.Lon}, nil
}
