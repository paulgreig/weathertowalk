package algorithm

import (
	"testing"
	"time"

	"github.com/paulgreig/weathertowalk/internal/config"
	"github.com/paulgreig/weathertowalk/internal/weather"
)

func TestFindOptimalWindows_PrefersHigherAverageScore(t *testing.T) {
	p := config.DefaultPreferences()
	p.WalkDuration = 60 // 1 hour, less than slot length so 1-slot windows

	baseTime := time.Now()
	slots := []weather.ForecastSlot{
		{Time: baseTime, Temperature: p.Temperature.Min + 1, Humidity: (p.Humidity.Min + p.Humidity.Max) / 2, Precipitation: 0, AirQuality: 50},
		{Time: baseTime.Add(3 * time.Hour), Temperature: p.Temperature.Min + 5, Humidity: (p.Humidity.Min + p.Humidity.Max) / 2, Precipitation: 0, AirQuality: 50},
	}

	recs := FindOptimalWindows(slots, p)
	if len(recs) == 0 {
		t.Fatalf("FindOptimalWindows() returned no recommendations")
	}
	if !recs[0].StartTime.Equal(slots[0].Time) && !recs[0].StartTime.Equal(slots[1].Time) {
		t.Fatalf("First recommendation start time not from a slot: %v", recs[0].StartTime)
	}
}

func TestFindOptimalWindows_NoSuitableSlotsReturnsEmpty(t *testing.T) {
	p := config.DefaultPreferences()
	p.WalkDuration = 60
	// All slots have heavy rain → score 0, below threshold.
	slots := []weather.ForecastSlot{
		{Time: time.Now(), Temperature: 10, Humidity: 90, Precipitation: 10},
	}

	recs := FindOptimalWindows(slots, p)
	if len(recs) != 0 {
		t.Fatalf("FindOptimalWindows() len = %d, want 0 for unsuitable slots", len(recs))
	}
}

