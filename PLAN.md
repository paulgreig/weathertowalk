# Weather to Walk - Development Plan

## Project Overview

A terminal-based Go application that helps users find the optimal time to walk based on current and forecasted weather conditions, considering their personal preferences for temperature, humidity, precipitation, and air quality.

## Core Requirements

### Phase 1: MVP (Minimum Viable Product)

1. **Weather Data Integration**
   - Fetch current weather and 3-hour forecast from a public weather API
   - Support for location-based queries (user's local area)
   - Handle API errors gracefully

2. **User Preferences**
   - Temperature range (min/max in Celsius or Fahrenheit)
   - Humidity range (min/max percentage)
   - Precipitation tolerance (none, light, moderate, any)
   - Air quality preferences (good, moderate, acceptable thresholds)

3. **Best Fit Algorithm**
   - Analyze 3-hour forecast against user preferences
   - Calculate a "walkability score" for each time slot
   - Identify optimal time window(s) for walking
   - Consider user-specified walk duration

4. **Terminal User Interface (TUI)**
   - Display current weather conditions
   - Show 3-hour forecast with walkability scores
   - Highlight recommended walking times
   - Allow user to configure preferences interactively

### Phase 2: Future Enhancements

1. **Route Planning**
   - Plot walking routes from user's location
   - Find flattest route (elevation-based routing)
   - Minimize road crossings
   - Integration with mapping APIs

2. **Additional Features**
   - Save user preferences
   - Historical weather data analysis
   - Multiple location support
   - Notifications for optimal walking windows

## Architecture

### Package Structure

```
weathertowalk/
├── cmd/
│   └── weathertowalk/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                   # Configuration management
│   │   ├── config.go
│   │   └── preferences.go
│   ├── weather/                  # Weather API integration
│   │   ├── client.go            # API client
│   │   ├── models.go            # Weather data models
│   │   └── parser.go            # Response parsing
│   ├── algorithm/                # Best fit algorithm
│   │   ├── scorer.go            # Walkability scoring
│   │   └── optimizer.go         # Time window optimization
│   ├── location/                 # Location services
│   │   └── geocoder.go          # Location detection/resolution
│   └── ui/                       # Terminal UI
│       ├── tui.go               # Main TUI controller
│       ├── components/          # UI components
│       │   ├── weather_display.go
│       │   ├── preferences_form.go
│       │   └── recommendations.go
│       └── styles.go            # UI styling
└── pkg/
    └── weather/                  # Reusable weather types (if needed)
```

## Technology Stack

### Core
- **Go 1.22+**: Main language
- **Cobra**: CLI framework (optional, for command-line args)
- **Bubble Tea**: Terminal UI framework (https://github.com/charmbracelet/bubbletea)
- **Lip Gloss**: Terminal styling (https://github.com/charmbracelet/lipgloss)

### External APIs
- **OpenWeatherMap API**: Primary weather data source
  - Free tier: 1,000 calls/day
  - Provides current weather + 3-hour forecast
  - Includes temperature, humidity, precipitation, air quality
  - Alternative: WeatherAPI.com, MeteoSource

### Dependencies
```go
// UI
github.com/charmbracelet/bubbletea v0.25.0
github.com/charmbracelet/lipgloss v0.9.1

// HTTP Client
net/http (standard library)

// Configuration
gopkg.in/yaml.v3 or encoding/json (standard library)

// Location (if needed)
// May use IP-based geolocation or require user input
```

## Implementation Phases

### Phase 1.1: Foundation & Configuration (Week 1)

**Goal**: Set up project structure and configuration management

**Tasks**:
1. Create configuration package
   - Load API keys from environment/config file
   - Default user preferences structure
   - Configuration validation

2. Create location detection
   - IP-based geolocation or user input
   - Location validation

**Deliverables**:
- `internal/config/config.go` - Configuration loading
- `internal/config/preferences.go` - User preferences model
- `internal/location/geocoder.go` - Location resolution
- Tests for all components

**Acceptance Criteria**:
- Can load API key from environment variable
- Can store/load user preferences
- Can resolve user location (IP or manual)

---

### Phase 1.2: Weather API Integration (Week 1-2)

**Goal**: Integrate with weather API and parse responses

**Tasks**:
1. Create weather API client
   - HTTP client with error handling
   - Request/response handling
   - Rate limiting consideration

2. Define weather data models
   - Current weather structure
   - Forecast structure (3-hour intervals)
   - Air quality data structure

3. Implement response parsing
   - JSON unmarshaling
   - Data validation
   - Error handling

**Deliverables**:
- `internal/weather/client.go` - API client
- `internal/weather/models.go` - Data structures
- `internal/weather/parser.go` - Response parsing
- Tests with mocked API responses

**Acceptance Criteria**:
- Can fetch current weather for a location
- Can fetch 3-hour forecast
- Handles API errors gracefully
- Returns structured, validated data

**API Endpoints to Use**:
- OpenWeatherMap: `/weather` (current) + `/forecast` (3-hour)
- Include air quality if available in free tier

---

### Phase 1.3: Best Fit Algorithm (Week 2)

**Goal**: Implement algorithm to find optimal walking times

**Tasks**:
1. Create walkability scorer
   - Score each forecast time slot based on preferences
   - Weight factors: temperature, humidity, precipitation, air quality
   - Calculate composite score (0-100)

2. Implement time window optimizer
   - Consider user-specified walk duration
   - Find contiguous time windows that meet preferences
   - Rank windows by score
   - Handle edge cases (no suitable times, partial matches)

**Algorithm Design**:
```
For each time slot in 3-hour forecast:
  score = 0
  
  If temperature in preferred range:
    score += 30
  Else if close to range:
    score += 15 (with distance penalty)
  
  If humidity in preferred range:
    score += 25
  Else if close:
    score += 12 (with distance penalty)
  
  If precipitation acceptable:
    score += 25
  Else:
    score = 0 (precipitation is blocking)
  
  If air quality acceptable:
    score += 20
  Else:
    score -= 10 (penalty but not blocking)
  
  walkability_score = score
```

**Deliverables**:
- `internal/algorithm/scorer.go` - Walkability scoring
- `internal/algorithm/optimizer.go` - Time window optimization
- Comprehensive tests with various scenarios

**Acceptance Criteria**:
- Calculates accurate walkability scores
- Identifies best time windows
- Handles user-specified duration
- Returns ranked recommendations

---

### Phase 1.4: Terminal User Interface (Week 2-3)

**Goal**: Build interactive terminal interface

**Tasks**:
1. Set up Bubble Tea framework
   - Main application model
   - Update/View pattern
   - Event handling

2. Create UI components
   - Current weather display
   - 3-hour forecast table/list
   - Walkability scores visualization
   - Recommendations panel
   - Preferences form (interactive)

3. Implement styling
   - Color scheme
   - Layout management
   - Responsive to terminal size

**UI Layout**:
```
┌─────────────────────────────────────────┐
│  Weather to Walk                        │
├─────────────────────────────────────────┤
│  Current Weather                        │
│  🌤️  22°C  Humidity: 65%  AQI: Good   │
├─────────────────────────────────────────┤
│  3-Hour Forecast                        │
│  ┌──────┬──────┬──────┬──────┐         │
│  │ Now  │ +1h  │ +2h  │ +3h  │         │
│  ├──────┼──────┼──────┼──────┤         │
│  │ 22°C │ 24°C │ 23°C │ 21°C │         │
│  │ 65%  │ 70%  │ 68%  │ 62%  │         │
│  │ ⭐95 │ ⭐88 │ ⭐92 │ ⭐98 │         │
│  └──────┴──────┴──────┴──────┘         │
├─────────────────────────────────────────┤
│  🎯 Best Time to Walk: Now - +1h        │
│     Score: 95/100                       │
│     Duration: 30 minutes                │
├─────────────────────────────────────────┤
│  [P] Preferences  [R] Refresh  [Q] Quit│
└─────────────────────────────────────────┘
```

**Deliverables**:
- `internal/ui/tui.go` - Main TUI controller
- `internal/ui/components/weather_display.go`
- `internal/ui/components/preferences_form.go`
- `internal/ui/components/recommendations.go`
- `internal/ui/styles.go`
- Tests for UI components (where applicable)

**Acceptance Criteria**:
- Displays current weather clearly
- Shows 3-hour forecast with scores
- Highlights best walking times
- Allows preference configuration
- Handles terminal resize
- Responsive and intuitive

---

### Phase 1.5: Integration & Polish (Week 3)

**Goal**: Integrate all components and finalize MVP

**Tasks**:
1. Wire up all components
   - Connect weather client to UI
   - Connect algorithm to UI
   - Handle loading states
   - Error display in UI

2. Add CLI interface
   - Command-line arguments
   - Help text
   - Version info

3. Error handling & edge cases
   - Network failures
   - Invalid API responses
   - No suitable walking times
   - Invalid user preferences

4. Documentation
   - User guide
   - API setup instructions
   - Configuration examples

**Deliverables**:
- `cmd/weathertowalk/main.go` - Complete application
- Integration tests
- User documentation
- Example configuration files

**Acceptance Criteria**:
- Application runs end-to-end
- All error cases handled gracefully
- User can configure preferences
- Recommendations are accurate
- Code passes all linting and tests

---

## Testing Strategy

### Unit Tests
- Each package should have comprehensive unit tests
- Mock external dependencies (API calls)
- Test edge cases and error conditions
- Target: >80% code coverage

### Integration Tests
- Test API client with mock server
- Test algorithm with various weather scenarios
- Test UI components in isolation

### Test Data
- Create test fixtures for weather responses
- Various weather scenarios (rainy, sunny, extreme temps)
- Edge cases (no suitable times, API failures)

---

## Configuration

### Environment Variables
```bash
WEATHER_API_KEY=your_api_key_here
WEATHER_API_URL=https://api.openweathermap.org/data/2.5
```

### User Preferences File (YAML/JSON)
```yaml
preferences:
  temperature:
    min: 15  # Celsius
    max: 25
  humidity:
    min: 40  # Percentage
    max: 70
  precipitation:
    tolerance: "light"  # none, light, moderate, any
  air_quality:
    min_score: "moderate"  # good, moderate, acceptable
  walk_duration: 30  # minutes
```

---

## Future Enhancements (Phase 2)

### Route Planning
- Integration with OpenRouteService or similar
- Elevation data for flattest route
- Pedestrian-friendly routing (minimize road crossings)
- Map visualization in terminal (ASCII art or external viewer)

### Additional Features
- Persistent preference storage
- Multiple location profiles
- Weather history and trends
- Notifications (desktop/system)
- Export recommendations (JSON, CSV)

---

## Development Workflow

### TDD Approach
1. Write failing test for feature
2. Implement minimal code to pass
3. Refactor while keeping tests green
4. Repeat for next feature

### Code Quality Gates
- All tests must pass
- `go fmt` must be clean
- `golangci-lint` must pass
- Mutation Testing before merge
- Code review before merge

### Git Workflow
- Feature branches for each phase
- Descriptive commit messages
- Tests included in every commit


---

## Success Metrics

### MVP Success Criteria
- ✅ Fetches accurate weather data
- ✅ Calculates walkability scores correctly
- ✅ Identifies optimal walking times
- ✅ User-friendly terminal interface
- ✅ Handles errors gracefully
- ✅ All tests passing
- ✅ Code quality standards met

### User Experience
- Application starts in <2 seconds
- Weather data loads in <3 seconds
- UI is responsive and intuitive
- Clear recommendations provided

---

## Risk Mitigation

### API Limitations
- **Risk**: Free tier API rate limits
- **Mitigation**: Implement caching, request batching, consider paid tier if needed

### Location Detection
- **Risk**: IP-based location may be inaccurate
- **Mitigation**: Allow manual location input, support multiple methods

### Algorithm Accuracy
- **Risk**: Walkability scores may not match user expectations
- **Mitigation**: Allow user feedback, make weights configurable, iterate based on usage

---

## Timeline Estimate

- **Phase 1.1**: 2-3 days
- **Phase 1.2**: 3-4 days
- **Phase 1.3**: 3-4 days
- **Phase 1.4**: 4-5 days
- **Phase 1.5**: 2-3 days

**Total MVP**: ~2-3 weeks of focused development

---

## Next Steps

1. Set up API account (OpenWeatherMap or alternative)
2. Create initial test structure
3. Begin Phase 1.1: Foundation & Configuration
4. Follow TDD workflow throughout
