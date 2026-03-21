# Weather to Walk — Android

This module is a minimal **Android Studio / Gradle** app that runs the same Go logic as the CLI via **gomobile** bindings.

## Prerequisites

- **JDK 17**
- **Android SDK** (API 34) and **NDK** — install [Android Studio](https://developer.android.com/studio) or command-line tools, set `ANDROID_HOME`
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
