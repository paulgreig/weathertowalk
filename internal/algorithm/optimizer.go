package algorithm

import (
	"sort"
	"time"

	"github.com/paulgreig/weathertowalk/internal/config"
	"github.com/paulgreig/weathertowalk/internal/weather"
)

// Recommendation describes a suggested walking window.
type Recommendation struct {
	StartTime      time.Time
	EndTime        time.Time
	Duration       time.Duration
	Score          int
	MinSlotScore   int
	WeatherSummary string
}

// FindOptimalWindows scores all slots and returns high-scoring recommendations.
// It uses a simple threshold on per-slot scores and returns individual-slot windows
// sized to the user's desired walk duration.
func FindOptimalWindows(slots []weather.ForecastSlot, prefs *config.UserPreferences) []Recommendation {
	if prefs == nil || len(slots) == 0 {
		return nil
	}
	const threshold = 70

	var recs []Recommendation
	for _, slot := range slots {
		score := ScoreSlot(slot, prefs)
		if score < threshold {
			continue
		}
		dur := time.Duration(prefs.WalkDuration) * time.Minute
		rec := Recommendation{
			StartTime:      slot.Time,
			EndTime:        slot.Time.Add(dur),
			Duration:       dur,
			Score:          score,
			MinSlotScore:   score,
			WeatherSummary: slot.Conditions,
		}
		recs = append(recs, rec)
	}

	// Sort recommendations by score desc, then start time asc.
	sort.Slice(recs, func(i, j int) bool {
		if recs[i].Score == recs[j].Score {
			return recs[i].StartTime.Before(recs[j].StartTime)
		}
		return recs[i].Score > recs[j].Score
	})

	// Limit to top N recommendations.
	const maxRecs = 5
	if len(recs) > maxRecs {
		recs = recs[:maxRecs]
	}
	return recs
}
