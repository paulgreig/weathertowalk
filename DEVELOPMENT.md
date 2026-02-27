# Development Checklist

## Phase 1.1: Foundation & Configuration

### Configuration Package
- [x] Create `internal/config/config.go`
  - [x] Load API key from environment variable
  - [x] Load API key from config file (optional, `LoadFromFile`)
  - [x] Validate API key format
  - [x] Default configuration values
- [x] Create `internal/config/preferences.go`
  - [x] `UserPreferences` struct
  - [x] Load preferences from file
  - [x] Save preferences to file
  - [x] Validate preferences (ranges, values)
  - [x] Default preferences
- [x] Tests for config package
  - [x] Test environment variable loading
  - [x] Test file loading
  - [x] Test validation
  - [x] Test defaults

### Location Service
- [x] Create `internal/location/location.go`
  - [ ] IP-based geolocation (optional, deferred)
  - [x] Manual location input (`NewLocation`)
  - [x] Location validation (`ValidateCoordinates`)
  - [x] `Location` struct with `String()`
- [x] Tests for location service
  - [ ] Test IP geolocation (mocked) — deferred
  - [x] Test manual input
  - [x] Test validation

**Deliverable**: Configuration and location services working with tests

---

## Phase 1.2: Weather API Integration

### Weather Models
- [ ] Create `internal/weather/models.go`
  - [ ] `WeatherData` struct
  - [ ] `CurrentWeather` struct
  - [ ] `ForecastSlot` struct
  - [ ] `Location` struct
  - [ ] JSON tags for API response
- [ ] Tests for models
  - [ ] Test struct creation
  - [ ] Test JSON marshaling/unmarshaling

### Weather Client
- [ ] Create `internal/weather/client.go`
  - [ ] HTTP client setup
  - [ ] Current weather endpoint
  - [ ] Forecast endpoint (3-hour)
  - [ ] Error handling
  - [ ] Request timeout
  - [ ] Rate limiting awareness
- [ ] Tests for client
  - [ ] Mock HTTP server
  - [ ] Test successful requests
  - [ ] Test error cases (network, API errors)
  - [ ] Test timeout handling

### Weather Parser
- [ ] Create `internal/weather/parser.go`
  - [ ] Parse current weather response
  - [ ] Parse forecast response
  - [ ] Extract 3-hour intervals
  - [ ] Data validation
  - [ ] Error handling for malformed data
- [ ] Tests for parser
  - [ ] Test valid responses
  - [ ] Test missing fields
  - [ ] Test invalid data types
  - [ ] Test edge cases

**Deliverable**: Can fetch and parse weather data from API

---

## Phase 1.3: Best Fit Algorithm

### Walkability Scorer
- [ ] Create `internal/algorithm/scorer.go`
  - [ ] `ScoreTimeSlot` function
  - [ ] Temperature scoring (0-30 points)
  - [ ] Humidity scoring (0-25 points)
  - [ ] Precipitation scoring (0-25 points, blocking)
  - [ ] Air quality scoring (0-20 points)
  - [ ] Composite score calculation
- [ ] Tests for scorer
  - [ ] Test perfect conditions (score = 100)
  - [ ] Test each preference component
  - [ ] Test edge cases (extreme values)
  - [ ] Test blocking conditions (precipitation)
  - [ ] Test partial matches

### Time Window Optimizer
- [ ] Create `internal/algorithm/optimizer.go`
  - [ ] `FindOptimalWindows` function
  - [ ] Sliding window algorithm
  - [ ] Duration consideration
  - [ ] Ranking algorithm
  - [ ] Handle no suitable times
- [ ] Tests for optimizer
  - [ ] Test single optimal window
  - [ ] Test multiple windows
  - [ ] Test duration constraints
  - [ ] Test no suitable times
  - [ ] Test ranking order

**Deliverable**: Algorithm correctly identifies optimal walking times

---

## Phase 1.4: Terminal User Interface

### TUI Setup
- [ ] Add Bubble Tea dependency
  - [ ] `go get github.com/charmbracelet/bubbletea`
  - [ ] `go get github.com/charmbracelet/lipgloss`
- [ ] Create `internal/ui/tui.go`
  - [ ] Main model struct
  - [ ] Init function
  - [ ] Update function (message handling)
  - [ ] View function (rendering)
  - [ ] Application lifecycle

### UI Components
- [ ] Create `internal/ui/components/weather_display.go`
  - [ ] Current weather display
  - [ ] Forecast table/list
  - [ ] Walkability score visualization
  - [ ] Weather icons/emojis
- [ ] Create `internal/ui/components/preferences_form.go`
  - [ ] Interactive form fields
  - [ ] Input validation
  - [ ] Save/cancel actions
- [ ] Create `internal/ui/components/recommendations.go`
  - [ ] Best time display
  - [ ] Score display
  - [ ] Duration display
  - [ ] Multiple recommendations
- [ ] Create `internal/ui/styles.go`
  - [ ] Color scheme
  - [ ] Layout helpers
  - [ ] Responsive sizing

### UI Integration
- [ ] Connect weather data to UI
- [ ] Connect algorithm to UI
- [ ] Handle loading states
- [ ] Handle error states
- [ ] Keyboard shortcuts (P for preferences, R for refresh, Q for quit)
- [ ] Auto-refresh option

**Deliverable**: Functional terminal UI displaying weather and recommendations

---

## Phase 1.5: Integration & Polish

### Main Application
- [ ] Create `cmd/weathertowalk/main.go`
  - [ ] Application entry point
  - [ ] Initialize all components
  - [ ] Error handling
  - [ ] Graceful shutdown
- [ ] CLI arguments (optional)
  - [ ] `--config` flag for config file
  - [ ] `--location` flag for manual location
  - [ ] `--help` flag
  - [ ] `--version` flag

### Error Handling
- [ ] Network failure handling
- [ ] API error handling
- [ ] Invalid response handling
- [ ] No suitable times message
- [ ] User-friendly error messages

### Documentation
- [ ] Update README.md with:
  - [ ] API setup instructions
  - [ ] Configuration examples
  - [ ] Usage examples
  - [ ] Troubleshooting
- [ ] Code comments
- [ ] Example config file

### Final Testing
- [ ] End-to-end test
- [ ] Manual testing of all features
- [ ] Test on different terminal sizes
- [ ] Test error scenarios
- [ ] Performance testing

### Code Quality
- [ ] All tests passing
- [ ] `go fmt ./...` clean
- [ ] `golangci-lint run ./...` passes
- [ ] Code coverage >80%
- [ ] No TODO comments in production code

**Deliverable**: Complete, polished MVP ready for use

---

## Testing Checklist (Apply to Each Phase)

### Unit Tests
- [ ] All exported functions have tests
- [ ] Edge cases covered
- [ ] Error cases covered
- [ ] Table-driven tests where appropriate
- [ ] Test helpers use `t.Helper()`

### Integration Tests
- [ ] Components work together
- [ ] Data flows correctly
- [ ] Error propagation works

### Test Quality
- [ ] Tests are readable
- [ ] Tests are maintainable
- [ ] Tests run quickly
- [ ] Tests are deterministic

---

## API Setup Checklist

### OpenWeatherMap (Recommended)
- [ ] Create account at https://openweathermap.org
- [ ] Get API key
- [ ] Verify free tier limits (1,000 calls/day)
- [ ] Test API key with curl/Postman
- [ ] Document API endpoints used

### Alternative APIs (if needed)
- [ ] WeatherAPI.com
- [ ] MeteoSource
- [ ] Compare features and limits

---

## Configuration Files to Create

### Example Config File
- [ ] Create `config.example.yaml`
- [ ] Document all options
- [ ] Provide sensible defaults

### Environment Template
- [ ] Create `.env.example`
- [ ] Document required variables
- [ ] Document optional variables

---

## Future Enhancements (Phase 2) - Not in MVP

- [ ] Route planning integration
- [ ] Elevation-based routing
- [ ] Minimize road crossings
- [ ] Map visualization
- [ ] Multiple location support
- [ ] Weather history
- [ ] Notifications
- [ ] Export functionality
