# Architecture Overview

## System Components

```
┌─────────────────────────────────────────────────────────┐
│                    User Interface (TUI)                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   Weather    │  │ Preferences  │  │Recommendations│ │
│  │   Display    │  │    Form      │  │    Panel      │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│              Application Core                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │   Weather    │  │  Algorithm    │  │  Location    │ │
│  │   Client     │  │  Optimizer    │  │  Service     │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│              External Services                          │
│  ┌──────────────┐                                      │
│  │  Weather API │  (OpenWeatherMap / WeatherAPI)       │
│  └──────────────┘                                      │
└─────────────────────────────────────────────────────────┘
```

## Data Flow

### 1. Application Startup
```
User runs app
    ↓
Load configuration (API key, preferences)
    ↓
Detect/resolve user location
    ↓
Initialize TUI
```

### 2. Weather Data Fetching
```
User requests weather (or auto-refresh)
    ↓
Location Service → Get coordinates
    ↓
Weather Client → API Request
    ↓
Weather Parser → Structured Data
    ↓
TUI → Display
```

### 3. Recommendation Generation
```
Weather Data Available
    ↓
Algorithm Scorer → Calculate scores for each time slot
    ↓
Algorithm Optimizer → Find best time windows
    ↓
TUI → Display recommendations
```

## Component Details

### Weather Client
- **Responsibility**: HTTP communication with weather API
- **Input**: Location (lat/lon or city name)
- **Output**: Raw API response
- **Error Handling**: Network errors, API errors, rate limits

### Weather Parser
- **Responsibility**: Transform API response to internal models
- **Input**: JSON response
- **Output**: `WeatherData` struct
- **Validation**: Required fields, data types, ranges

### Algorithm Scorer
- **Responsibility**: Calculate walkability score for a time slot
- **Input**: Weather conditions + User preferences
- **Output**: Score (0-100)
- **Logic**: Weighted scoring based on preferences

### Algorithm Optimizer
- **Responsibility**: Find optimal time windows
- **Input**: Scored time slots + Walk duration
- **Output**: Ranked recommendations
- **Logic**: Sliding window algorithm

### Location Service
- **Responsibility**: Resolve user location
- **Methods**: IP geolocation, manual input, saved location
- **Output**: Coordinates (lat/lon)

### TUI Controller
- **Responsibility**: Manage UI state and interactions
- **Components**: Weather display, preferences form, recommendations
- **Framework**: Bubble Tea (Model-Update-View pattern)

## Data Models

### WeatherData
```go
type WeatherData struct {
    Current    CurrentWeather
    Forecast   []ForecastSlot  // 3-hour intervals
    Location   Location
    Timestamp  time.Time
}

type CurrentWeather struct {
    Temperature     float64  // Celsius
    Humidity        float64  // Percentage
    Precipitation   float64  // mm
    AirQuality      int      // AQI
    Conditions      string   // "clear", "rain", etc.
    WindSpeed       float64  // m/s
}

type ForecastSlot struct {
    Time            time.Time
    Temperature     float64
    Humidity        float64
    Precipitation   float64
    AirQuality      int
    Conditions      string
    WalkabilityScore int  // Calculated
}
```

### UserPreferences
```go
type UserPreferences struct {
    Temperature     TemperatureRange
    Humidity        HumidityRange
    Precipitation   PrecipitationTolerance
    AirQuality      AirQualityThreshold
    WalkDuration    int  // minutes
}

type TemperatureRange struct {
    Min float64
    Max float64
    Unit string  // "celsius" or "fahrenheit"
}

type HumidityRange struct {
    Min float64  // percentage
    Max float64
}

type PrecipitationTolerance string
const (
    ToleranceNone     PrecipitationTolerance = "none"
    ToleranceLight    PrecipitationTolerance = "light"
    ToleranceModerate PrecipitationTolerance = "moderate"
    ToleranceAny      PrecipitationTolerance = "any"
)

type AirQualityThreshold string
const (
    AQGood       AirQualityThreshold = "good"
    AQModerate   AirQualityThreshold = "moderate"
    AQAcceptable AirQualityThreshold = "acceptable"
)
```

### Recommendation
```go
type Recommendation struct {
    StartTime       time.Time
    EndTime         time.Time
    Duration        time.Duration
    Score           int  // 0-100
    Reason          string  // Why this time is good
    WeatherSummary  string
}
```

## Algorithm Details

### Walkability Scoring

**Formula**:
```
score = 0

// Temperature (30 points max)
if temp in [pref_min, pref_max]:
    score += 30
else if temp in [pref_min - 5, pref_max + 5]:
    distance = min(|temp - pref_min|, |temp - pref_max|)
    score += 30 - (distance * 3)  // Linear penalty
else:
    score += 0  // Too far from preference

// Humidity (25 points max)
if humidity in [pref_min, pref_max]:
    score += 25
else if humidity in [pref_min - 10, pref_max + 10]:
    distance = min(|humidity - pref_min|, |humidity - pref_max|)
    score += 25 - (distance * 2)
else:
    score += 0

// Precipitation (25 points max, blocking)
if precipitation == 0:
    score += 25
else if precipitation < 0.5 && tolerance >= "light":
    score += 20
else if precipitation < 2.0 && tolerance >= "moderate":
    score += 15
else if tolerance == "any":
    score += 10
else:
    score = 0  // Blocking condition

// Air Quality (20 points max)
if aqi <= 50:  // Good
    score += 20
else if aqi <= 100 && threshold >= "good":
    score += 15
else if aqi <= 150 && threshold >= "moderate":
    score += 10
else if threshold == "acceptable":
    score += 5
else:
    score -= 10  // Penalty but not blocking

walkability_score = min(100, max(0, score))
```

### Time Window Optimization

**Algorithm**:
1. Score all time slots in 3-hour forecast
2. Filter slots with score >= threshold (e.g., 70)
3. For each potential start time:
   - Check if duration fits within forecast window
   - Calculate average score for the window
   - Consider continuity (prefer contiguous good slots)
4. Rank windows by:
   - Average score (primary)
   - Minimum score in window (secondary)
   - Start time proximity (tertiary)
5. Return top N recommendations

## Error Handling Strategy

### API Errors
- **Network failures**: Retry with exponential backoff (max 3 retries)
- **Rate limiting**: Cache last successful response, show warning
- **Invalid responses**: Log error, show user-friendly message
- **Missing data**: Use defaults where possible, mark as incomplete

### User Input Errors
- **Invalid preferences**: Validate on input, show clear error messages
- **Invalid location**: Allow retry, suggest alternatives
- **No suitable times**: Explain why, suggest preference adjustments

### Application Errors
- **Panic recovery**: Catch panics, log, show error message, allow graceful exit
- **State corruption**: Reset to safe state, reload data

## Performance Considerations

### Caching
- Cache weather data for 10 minutes (API rate limits)
- Cache location resolution
- Cache user preferences

### Optimization
- Lazy loading of UI components
- Efficient string formatting for TUI
- Minimize API calls (batch if possible)

### Resource Usage
- Memory: Target <50MB for typical usage
- CPU: Minimal during idle, responsive during updates
- Network: Only when fetching weather data

## Security Considerations

### API Keys
- Never commit API keys to repository
- Load from environment variables or secure config file
- Validate API key format before use

### User Data
- Preferences stored locally only
- No personal data transmitted (except location for weather)
- Location can be approximate (city-level) for privacy

## Testing Architecture

### Unit Tests
- Mock external dependencies
- Test each component in isolation
- Use table-driven tests for algorithms

### Integration Tests
- Mock HTTP server for API testing
- Test data flow between components
- Test error propagation

### UI Tests
- Test UI models and state transitions
- Manual testing for visual components
- Test user interaction flows
