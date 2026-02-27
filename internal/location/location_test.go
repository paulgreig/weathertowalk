package location

import (
	"testing"
)

func TestValidateCoordinates_Valid(t *testing.T) {
	tests := []struct {
		lat, lon float64
	}{
		{0, 0},
		{51.5074, -0.1278},
		{-33.8688, 151.2093},
		{90, 180},
		{-90, -180},
	}
	for _, tt := range tests {
		err := ValidateCoordinates(tt.lat, tt.lon)
		if err != nil {
			t.Errorf("ValidateCoordinates(%v, %v) err = %v, want nil", tt.lat, tt.lon, err)
		}
	}
}

func TestValidateCoordinates_InvalidLatitude(t *testing.T) {
	tests := []struct {
		lat, lon float64
	}{
		{91, 0},
		{-91, 0},
		{100, 10},
	}
	for _, tt := range tests {
		err := ValidateCoordinates(tt.lat, tt.lon)
		if err == nil {
			t.Errorf("ValidateCoordinates(%v, %v) err = nil, want error", tt.lat, tt.lon)
		}
	}
}

func TestValidateCoordinates_InvalidLongitude(t *testing.T) {
	tests := []struct {
		lat, lon float64
	}{
		{0, 181},
		{0, -181},
		{0, 200},
	}
	for _, tt := range tests {
		err := ValidateCoordinates(tt.lat, tt.lon)
		if err == nil {
			t.Errorf("ValidateCoordinates(%v, %v) err = nil, want error", tt.lat, tt.lon)
		}
	}
}

func TestLocation_String(t *testing.T) {
	loc := Location{Lat: 51.5, Lon: -0.1}
	s := loc.String()
	if s == "" {
		t.Error("Location.String() returned empty string")
	}
}

func TestNewLocation_Valid(t *testing.T) {
	loc, err := NewLocation(51.5, -0.1)
	if err != nil {
		t.Fatalf("NewLocation() err = %v, want nil", err)
	}
	if loc.Lat != 51.5 || loc.Lon != -0.1 {
		t.Errorf("NewLocation() = %+v, want Lat=51.5 Lon=-0.1", loc)
	}
}

func TestNewLocation_Invalid(t *testing.T) {
	_, err := NewLocation(100, 0)
	if err == nil {
		t.Fatal("NewLocation(100, 0) err = nil, want error")
	}
}
