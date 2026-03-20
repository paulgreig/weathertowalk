// Package mobile provides gomobile bindings for Android. Build the AAR with:
//
//	go install golang.org/x/mobile/cmd/gomobile@latest
//	gomobile init
//	gomobile bind -target=android -o ../android/app/libs/weathertowalk.aar .
//
// See android/README.md for the full APK build.
package mobile

import "github.com/paulgreig/weathertowalk/internal/service"

// RunWalk fetches weather and walking recommendations for the given coordinates.
// baseURL may be empty to use the default OpenWeatherMap API base URL.
// Returns a non-nil error if the request fails or configuration is invalid.
func RunWalk(apiKey, baseURL string, lat, lon float64) (string, error) {
	return service.RunWalkReport(apiKey, baseURL, lat, lon)
}
