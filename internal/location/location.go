package location

import (
	"fmt"
)

// Location represents geographic coordinates.
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// ValidateCoordinates returns an error if lat or lon are out of valid range.
// Lat must be in [-90, 90], lon in [-180, 180].
func ValidateCoordinates(lat, lon float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude %v out of range [-90, 90]", lat)
	}
	if lon < -180 || lon > 180 {
		return fmt.Errorf("longitude %v out of range [-180, 180]", lon)
	}
	return nil
}

// String returns a short representation of the location.
func (l Location) String() string {
	return fmt.Sprintf("%.4f, %.4f", l.Lat, l.Lon)
}

// NewLocation validates lat/lon and returns a Location. Use for manual input.
func NewLocation(lat, lon float64) (Location, error) {
	if err := ValidateCoordinates(lat, lon); err != nil {
		return Location{}, err
	}
	return Location{Lat: lat, Lon: lon}, nil
}
