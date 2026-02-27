package algorithm

import (
	"testing"

	"github.com/paulgreig/weathertowalk/internal/config"
	"github.com/paulgreig/weathertowalk/internal/weather"
)

func TestScoreSlot_PerfectConditionsReturns100(t *testing.T) {
	p := config.DefaultPreferences()
	slot := weather.ForecastSlot{
		Temperature:   (p.Temperature.Min + p.Temperature.Max) / 2,
		Humidity:      (p.Humidity.Min + p.Humidity.Max) / 2,
		Precipitation: 0,
		AirQuality:    50,
		Conditions:    "Clear",
	}

	score := ScoreSlot(slot, p)
	if score != 100 {
		t.Fatalf("ScoreSlot(perfect) = %d, want 100", score)
	}
}

func TestScoreSlot_HeavyRainBlocksScore(t *testing.T) {
	p := config.DefaultPreferences()
	p.Precipitation = config.ToleranceLight
	slot := weather.ForecastSlot{
		Temperature:   (p.Temperature.Min + p.Temperature.Max) / 2,
		Humidity:      (p.Humidity.Min + p.Humidity.Max) / 2,
		Precipitation: 10.0, // heavy rain
		AirQuality:    50,
	}

	score := ScoreSlot(slot, p)
	if score != 0 {
		t.Fatalf("ScoreSlot(heavy rain) = %d, want 0 (blocking)", score)
	}
}

func TestScoreSlot_TemperaturePenaltyApplied(t *testing.T) {
	p := config.DefaultPreferences()
	// Slightly outside preferred range but within tolerance band.
	slot := weather.ForecastSlot{
		Temperature:   p.Temperature.Max + 3,
		Humidity:      (p.Humidity.Min + p.Humidity.Max) / 2,
		Precipitation: 0,
		AirQuality:    50,
	}

	score := ScoreSlot(slot, p)
	if score <= 0 || score >= 100 {
		t.Fatalf("ScoreSlot(temp slightly high) = %d, want between 1 and 99", score)
	}
}

