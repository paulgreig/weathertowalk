# Quick Start Guide

## Project Summary

**Weather to Walk** is a terminal-based Go application that helps users find the optimal time to walk based on weather conditions and personal preferences.

## Key Features (MVP)

✅ Current weather display  
✅ 3-hour weather forecast  
✅ Personalized preferences (temperature, humidity, precipitation, air quality)  
✅ Best-fit algorithm to find optimal walking times  
✅ Terminal user interface (TUI)  
✅ Walk duration consideration  

## Architecture Overview

```
User (TUI) → Algorithm → Weather API
     ↓           ↓
Preferences  Scoring & Optimization
```

## Development Phases

1. **Foundation** - Configuration and location services
2. **Weather Integration** - API client and data parsing
3. **Algorithm** - Walkability scoring and optimization
4. **UI** - Terminal interface with Bubble Tea
5. **Polish** - Integration, error handling, documentation

## Technology Stack

- **Go 1.25.6+** - Core language
- **Standard library TUI** (line-based prompts) for now
- **OpenWeatherMap API** - Weather data source
- **TDD** - Test-driven development approach

## Running the App

```bash
export WEATHER_API_KEY=your_api_key_here
# Optional, defaults to London:
# export WEATHER_LAT=51.5074
# export WEATHER_LON=-0.1278

go run ./cmd/weathertowalk
```

You will see:
- A summary dashboard (current conditions, short forecast, recommendations)
- A prompt allowing you to refresh data or quit

## Getting Started

### 1. Set Up API Key

```bash
export WEATHER_API_KEY=your_api_key_here
```

Get your free API key from: https://openweathermap.org/api

### 2. Development Workflow

```bash
# Run tests
make test

# Format code
make fmt

# Run linter
make lint

# Run all checks
make check
```

### 3. Start Development

Begin with Phase 1.1 (Foundation & Configuration) following the TDD workflow:
1. Write failing test
2. Implement minimal code
3. Refactor
4. Repeat

## Project Structure

```
weathertowalk/
├── cmd/weathertowalk/     # Main application
├── internal/
│   ├── config/            # Configuration
│   ├── weather/           # Weather API
│   ├── algorithm/         # Best-fit algorithm
│   ├── location/          # Location services
│   └── ui/                # Terminal UI
├── .cursor/rules/         # Cursor IDE rules
└── docs/                  # Documentation
```

## Key Documents

- **PLAN.md** - Detailed development plan
- **ARCHITECTURE.md** - System architecture and design
- **DEVELOPMENT.md** - Development checklist
- **README.md** - Project overview and setup

## Next Steps

1. Review PLAN.md for detailed requirements
2. Review ARCHITECTURE.md for system design
3. Follow DEVELOPMENT.md checklist
4. Start with Phase 1.1: Foundation & Configuration

## Success Criteria

- ✅ Fetches accurate weather data
- ✅ Calculates walkability scores
- ✅ Identifies optimal walking times
- ✅ User-friendly terminal interface
- ✅ All tests passing
- ✅ Code quality standards met
