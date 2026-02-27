package algorithm

import (
	"math"

	"github.com/paulgreig/weathertowalk/internal/config"
	"github.com/paulgreig/weathertowalk/internal/weather"
)

// ScoreSlot computes a walkability score (0–100) for a given forecast slot and preferences.
func ScoreSlot(slot weather.ForecastSlot, prefs *config.UserPreferences) int {
	if prefs == nil {
		return 0
	}
	score := 0

	// Temperature (30 points max)
	temp := slot.Temperature
	tMin, tMax := prefs.Temperature.Min, prefs.Temperature.Max
	if temp >= tMin && temp <= tMax {
		score += 30
	} else if temp >= tMin-5 && temp <= tMax+5 {
		distance := math.Min(math.Abs(temp-tMin), math.Abs(temp-tMax))
		partial := 30 - int(distance*3)
		if partial > 0 {
			score += partial
		}
	}

	// Humidity (25 points max)
	h := slot.Humidity
	hMin, hMax := prefs.Humidity.Min, prefs.Humidity.Max
	if h >= hMin && h <= hMax {
		score += 25
	} else if h >= hMin-10 && h <= hMax+10 {
		distance := math.Min(math.Abs(h-hMin), math.Abs(h-hMax))
		partial := 25 - int(distance*2)
		if partial > 0 {
			score += partial
		}
	}

	// Precipitation (25 points max, blocking)
	p := slot.Precipitation
	switch prefs.Precipitation {
	case config.ToleranceNone:
		if p == 0 {
			score += 25
		} else {
			return 0
		}
	case config.ToleranceLight:
		if p == 0 {
			score += 25
		} else if p < 0.5 {
			score += 20
		} else {
			return 0
		}
	case config.ToleranceModerate:
		if p == 0 {
			score += 25
		} else if p < 0.5 {
			score += 20
		} else if p < 2.0 {
			score += 15
		} else {
			return 0
		}
	case config.ToleranceAny:
		if p == 0 {
			score += 25
		} else if p < 0.5 {
			score += 20
		} else if p < 2.0 {
			score += 15
		} else {
			score += 10
		}
	}

	// Air Quality (20 points max)
	aq := slot.AirQuality
	switch prefs.AirQuality {
	case config.AQGood:
		if aq <= 50 {
			score += 20
		} else if aq <= 100 {
			score += 15
		} else {
			score -= 10
		}
	case config.AQModerate:
		if aq <= 50 {
			score += 20
		} else if aq <= 100 {
			score += 15
		} else if aq <= 150 {
			score += 10
		} else {
			score -= 10
		}
	case config.AQAcceptable:
		if aq <= 50 {
			score += 20
		} else if aq <= 100 {
			score += 15
		} else if aq <= 150 {
			score += 10
		} else {
			score += 5
		}
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	return score
}

