# gomobile AAR

Place **`weathertowalk.aar`** here before building the Android app.

Generate it from the repository root (requires [Android SDK](https://developer.android.com/studio), `ANDROID_HOME`, and [Android NDK](https://developer.android.com/ndk)):

```bash
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init
cd mobile
gomobile bind -target=android -o ../app/libs/weathertowalk.aar .
```

Or use `make android-aar` from the repo root (see `Makefile`).

The Kotlin code imports **`mobile.Mobile`** and calls **`Mobile.runWalk(...)`**, matching the default Java package produced by `gomobile bind` for the `mobile` Go package.
