# Weather to Walk — Android

This module is a minimal **Android Studio / Gradle** app that runs the same Go logic as the CLI via **gomobile** bindings.

## Dev Container

From the repository root, use **Dev Containers** (`.devcontainer/`) so **JDK 17**, **Android SDK**, **NDK**, **Go**, and **gomobile** are preinstalled. Then:

```bash
make android-apk
```

## Prerequisites (manual / Android Studio)

- **JDK 17**
- **Android SDK** (API 34) and **NDK** — install [Android Studio](https://developer.android.com/studio) or run `scripts/setup-android-sdk.sh` (Linux CLI tools), set `ANDROID_HOME` and **`ANDROID_NDK_HOME` to NDK 21.4.x** for `gomobile bind` (see `docs/APK_BUILD_VERIFICATION.md`)
- **gomobile**: `go install golang.org/x/mobile/cmd/gomobile@latest` then `gomobile init`

## Build the Go library (AAR)

From the repository root:

```bash
make android-aar
```

This produces `app/libs/weathertowalk.aar`, which Gradle consumes as a local dependency.

To build manually:

```bash
cd ../mobile
gomobile bind -target=android -o ../android/app/libs/weathertowalk.aar .
```

## Build the APK

From the repo root (after `make android-aar` or via `make android-apk`):

```bash
make android-apk
```

Or only Gradle, if the AAR is already built:

```bash
cd android
./gradlew assembleDebug
```

Output: `app/build/outputs/apk/debug/app-debug.apk`.

## How it works

- **`mobile` Go package** (`../mobile`) exports `RunWalk` for `gomobile bind`.
- **`WeatherWalkService`** is a foreground service (`dataSync`) that invokes **`mobile.Mobile.runWalk`** on a worker thread and broadcasts the result string to **`MainActivity`**.
- **Internet** permission is required for OpenWeatherMap calls.

## Java / Kotlin API

With the default `gomobile` settings (no `-javapkg`), the generated class is **`mobile.Mobile`** with a static **`runWalk`** method matching the Go signature.
