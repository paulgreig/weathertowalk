# Mutation Testing Report

**Date**: 2025-01-28  
**Target**: `internal/config`, `internal/location`  
**Test suite**: `go test ./...`

## Summary

Mutation testing was run on config and location packages. One mutant **survived** initially; we added tests to kill it and to close related gaps.

## Tests Added (Improvements)

| Test | Purpose |
|------|---------|
| `TestLoadFromFile_DefaultAPIBaseURLWhenOmitted` | Asserts that when the config file omits `api_base_url`, `LoadFromFile` sets `APIBaseURL` to the default. **Kills survived mutant**: `if c.APIBaseURL == ""` → `!= ""`. |
| `TestValidatePreferences_NilReturnsError` | Asserts `ValidatePreferences(nil)` returns an error. Covers nil-check so mutations that remove it would be caught. |
| `TestSavePreferences_NilReturnsError` | Asserts `SavePreferences(path, nil)` returns an error. Ensures we never treat nil preferences as valid. |

## Survived Mutant (Now Killed)

- **Location**: `internal/config/config.go`, `LoadFromFile`
- **Mutation**: `if c.APIBaseURL == ""` → `if c.APIBaseURL != ""`
- **Effect**: When `api_base_url` is omitted from the config file, the default was no longer applied; `APIBaseURL` stayed `""`.
- **Gap**: No test asserted that `LoadFromFile` defaults `APIBaseURL` when it is missing from the file.
- **Fix**: `TestLoadFromFile_DefaultAPIBaseURLWhenOmitted` loads a config without `api_base_url` and asserts `APIBaseURL` equals the default URL.

## Other Mutants (All Killed)

- **ValidateAPIKey**: empty check `==`→`!=`, `\|\|`→`&&`, `len != 32`→`== 32`, hex-loop `&&`→`\|\|`
- **defaultOrEnv**: `v != ""`→`v == ""`
- **ValidatePreferences**: nil check, temp/humidity `>=`→`<`, walk duration `<=`→`>`
- **ValidateCoordinates**: `lat < -90`→`>=`, `lat > 90`→`<=`

## Recommendations

1. **Keep** `TestLoadFromFile_DefaultAPIBaseURLWhenOmitted` — it protects the APIBaseURL defaulting behaviour.
2. **Keep** `TestValidatePreferences_NilReturnsError` and `TestSavePreferences_NilReturnsError` — they lock down nil handling.
3. Re-run mutation testing periodically (e.g. after refactors) to catch new survivors.
