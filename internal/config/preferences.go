package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// UserPreferences holds user preferences for walking conditions.
type UserPreferences struct {
	Temperature   TemperatureRange       `json:"temperature"`
	Humidity      HumidityRange          `json:"humidity"`
	Precipitation PrecipitationTolerance `json:"precipitation"`
	AirQuality    AirQualityThreshold    `json:"air_quality"`
	WalkDuration  int                    `json:"walk_duration"` // minutes
}

// TemperatureRange defines preferred temperature bounds.
type TemperatureRange struct {
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
	Unit string  `json:"unit"` // "celsius" or "fahrenheit"
}

// HumidityRange defines preferred humidity bounds (percent).
type HumidityRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// PrecipitationTolerance indicates acceptable precipitation level.
type PrecipitationTolerance string

const (
	ToleranceNone     PrecipitationTolerance = "none"
	ToleranceLight    PrecipitationTolerance = "light"
	ToleranceModerate PrecipitationTolerance = "moderate"
	ToleranceAny      PrecipitationTolerance = "any"
)

// AirQualityThreshold indicates minimum acceptable air quality.
type AirQualityThreshold string

const (
	AQGood       AirQualityThreshold = "good"
	AQModerate   AirQualityThreshold = "moderate"
	AQAcceptable AirQualityThreshold = "acceptable"
)

// DefaultPreferences returns sensible default preferences.
func DefaultPreferences() *UserPreferences {
	return &UserPreferences{
		Temperature:   TemperatureRange{Min: 15, Max: 25, Unit: "celsius"},
		Humidity:      HumidityRange{Min: 40, Max: 70},
		Precipitation: ToleranceLight,
		AirQuality:    AQModerate,
		WalkDuration:  30,
	}
}

// ValidatePreferences checks that preferences are consistent and valid.
func ValidatePreferences(p *UserPreferences) error {
	if p == nil {
		return fmt.Errorf("preferences cannot be nil")
	}
	if p.Temperature.Min >= p.Temperature.Max {
		return fmt.Errorf("temperature min (%v) must be less than max (%v)", p.Temperature.Min, p.Temperature.Max)
	}
	if p.Humidity.Min >= p.Humidity.Max {
		return fmt.Errorf("humidity min (%v) must be less than max (%v)", p.Humidity.Min, p.Humidity.Max)
	}
	if p.WalkDuration <= 0 {
		return fmt.Errorf("walk_duration must be positive, got %d", p.WalkDuration)
	}
	unit := p.Temperature.Unit
	if unit != "celsius" && unit != "fahrenheit" {
		return fmt.Errorf("temperature unit must be celsius or fahrenheit, got %q", unit)
	}
	switch p.Precipitation {
	case ToleranceNone, ToleranceLight, ToleranceModerate, ToleranceAny:
	default:
		return fmt.Errorf("invalid precipitation tolerance %q", p.Precipitation)
	}
	switch p.AirQuality {
	case AQGood, AQModerate, AQAcceptable:
	default:
		return fmt.Errorf("invalid air_quality %q", p.AirQuality)
	}
	return nil
}

// LoadPreferences reads preferences from a JSON file.
func LoadPreferences(path string) (*UserPreferences, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read preferences file: %w", err)
	}
	var p UserPreferences
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parse preferences file: %w", err)
	}
	if err := ValidatePreferences(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

// SavePreferences writes preferences to a JSON file.
func SavePreferences(path string, p *UserPreferences) error {
	if err := ValidatePreferences(p); err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("encode preferences: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write preferences file: %w", err)
	}
	return nil
}
