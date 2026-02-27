package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultPreferences(t *testing.T) {
	p := DefaultPreferences()
	if p.Temperature.Min >= p.Temperature.Max {
		t.Errorf("Temperature Min %v >= Max %v", p.Temperature.Min, p.Temperature.Max)
	}
	if p.Humidity.Min >= p.Humidity.Max {
		t.Errorf("Humidity Min %v >= Max %v", p.Humidity.Min, p.Humidity.Max)
	}
	if p.WalkDuration <= 0 {
		t.Errorf("WalkDuration = %v, want positive", p.WalkDuration)
	}
	if p.Temperature.Unit != "celsius" && p.Temperature.Unit != "fahrenheit" {
		t.Errorf("Temperature.Unit = %q, want celsius or fahrenheit", p.Temperature.Unit)
	}
}

func TestValidatePreferences_Valid(t *testing.T) {
	p := DefaultPreferences()
	if err := ValidatePreferences(p); err != nil {
		t.Errorf("ValidatePreferences(default) err = %v, want nil", err)
	}
}

func TestValidatePreferences_InvalidTemperatureRange(t *testing.T) {
	p := DefaultPreferences()
	p.Temperature.Min = 30
	p.Temperature.Max = 20
	err := ValidatePreferences(p)
	if err == nil {
		t.Fatal("ValidatePreferences() err = nil, want error for Min > Max")
	}
}

func TestValidatePreferences_InvalidHumidityRange(t *testing.T) {
	p := DefaultPreferences()
	p.Humidity.Min = 90
	p.Humidity.Max = 50
	err := ValidatePreferences(p)
	if err == nil {
		t.Fatal("ValidatePreferences() err = nil, want error for Min > Max")
	}
}

func TestValidatePreferences_InvalidWalkDuration(t *testing.T) {
	p := DefaultPreferences()
	p.WalkDuration = 0
	err := ValidatePreferences(p)
	if err == nil {
		t.Fatal("ValidatePreferences() err = nil, want error for WalkDuration <= 0")
	}
}

func TestValidatePreferences_NilReturnsError(t *testing.T) {
	err := ValidatePreferences(nil)
	if err == nil {
		t.Fatal("ValidatePreferences(nil) err = nil, want error")
	}
}

func TestSavePreferences_NilReturnsError(t *testing.T) {
	err := SavePreferences(filepath.Join(t.TempDir(), "prefs.json"), nil)
	if err == nil {
		t.Fatal("SavePreferences(_, nil) err = nil, want error")
	}
}

func TestLoadPreferences_FileNotFound(t *testing.T) {
	_, err := LoadPreferences(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err == nil {
		t.Fatal("LoadPreferences() err = nil, want error for missing file")
	}
}

func TestLoadPreferences_ValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "prefs.json")
	content := `{
		"temperature": {"min": 10, "max": 25, "unit": "celsius"},
		"humidity": {"min": 30, "max": 70},
		"precipitation": "light",
		"air_quality": "moderate",
		"walk_duration": 45
	}`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write prefs file: %v", err)
	}

	p, err := LoadPreferences(path)
	if err != nil {
		t.Fatalf("LoadPreferences() err = %v, want nil", err)
	}
	if p.Temperature.Min != 10 || p.Temperature.Max != 25 {
		t.Errorf("Temperature = [%v, %v], want [10, 25]", p.Temperature.Min, p.Temperature.Max)
	}
	if p.Humidity.Min != 30 || p.Humidity.Max != 70 {
		t.Errorf("Humidity = [%v, %v], want [30, 70]", p.Humidity.Min, p.Humidity.Max)
	}
	if p.Precipitation != ToleranceLight {
		t.Errorf("Precipitation = %q, want light", p.Precipitation)
	}
	if p.AirQuality != AQModerate {
		t.Errorf("AirQuality = %q, want moderate", p.AirQuality)
	}
	if p.WalkDuration != 45 {
		t.Errorf("WalkDuration = %v, want 45", p.WalkDuration)
	}
}

func TestSavePreferences_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "prefs.json")
	p := DefaultPreferences()
	p.WalkDuration = 30

	if err := SavePreferences(path, p); err != nil {
		t.Fatalf("SavePreferences() err = %v, want nil", err)
	}
	loaded, err := LoadPreferences(path)
	if err != nil {
		t.Fatalf("LoadPreferences() after save err = %v, want nil", err)
	}
	if loaded.WalkDuration != 30 {
		t.Errorf("after round-trip WalkDuration = %v, want 30", loaded.WalkDuration)
	}
	if loaded.Temperature.Min != p.Temperature.Min || loaded.Temperature.Max != p.Temperature.Max {
		t.Error("Temperature changed after round-trip")
	}
}
