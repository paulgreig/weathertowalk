# APK build & run verification

This documents a successful **debug APK** build and **install on an Android emulator** (API 34, x86_64) in a Linux CI-style environment. A physical USB device was not used; `adb` targeted `emulator-5554`, which is the same install path used for real hardware.

## Artifacts

| Artifact | Path |
|----------|------|
| Demo video (title + UI screenshot + summary + raw screen capture) | `docs/weathertowalk-apk-demo.mp4` |
| Emulator screenshot (PNG) | `docs/weathertowalk-emulator-screenshot.png` |
| Raw `adb screenrecord` (low FPS on headless emulator) | `docs/weathertowalk-demo-raw-screen.mp4` |
| Debug APK | `android/app/build/outputs/apk/debug/app-debug.apk` (after `make android-apk`) |

## Toolchain notes

- **Gradle**: wrapper updated to **8.7** (AGP 8.5.2 requires Gradle ≥ 8.7).
- **gomobile**: current `gomobile bind` rejects **NDK 26** with `unsupported API version 16`; use **NDK 21.4.7075529** (side-by-side) and set `ANDROID_NDK_HOME` to that directory. The Dev Container and `scripts/setup-android-sdk.sh` use this version.
- **Go**: `go.mod` requires **1.25.6**; `make android-aar` runs gomobile with **`GOTOOLCHAIN=auto`** so the toolchain downloads when the host `go` is older.

## Commands run (high level)

1. Install SDK + NDK (`scripts/setup-android-sdk.sh` or Dev Container).
2. `go install golang.org/x/mobile/cmd/gomobile@latest && gomobile init`
3. `export ANDROID_NDK_HOME=$ANDROID_HOME/ndk/21.4.7075529` (if not using the default from the script).
4. `make android-apk` from repo root.
5. Start emulator or connect device; `adb install -r android/app/build/outputs/apk/debug/app-debug.apk`.
6. `adb shell am start -n com.paulgreig.weathertowalk/.MainActivity`

## Headless emulator caveat

Screenrecord on a **no-window** emulator with **software acceleration** produces **very low FPS**; the demo video therefore includes a **static `adb`** screenshot of the running UI and a text summary slide.
