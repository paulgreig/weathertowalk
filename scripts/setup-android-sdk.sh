#!/usr/bin/env bash
# Install Android command-line SDK + API 34 + NDK for local / CI Linux hosts (not Android Studio).
# Requires: curl or wget, unzip, openjdk-17-jdk. After running, set ANDROID_HOME and run: gomobile init
set -euo pipefail

ANDROID_CLI_TOOLS_VERSION="${ANDROID_CLI_TOOLS_VERSION:-11076708}"
NDK_VERSION="${NDK_VERSION:-26.3.11579264}"
INSTALL_ROOT="${ANDROID_SDK_ROOT:-${ANDROID_HOME:-$HOME/android-sdk}}"

echo "Installing Android SDK under: $INSTALL_ROOT"

mkdir -p "$INSTALL_ROOT/cmdline-tools"
ZIP="/tmp/commandlinetools-linux-${ANDROID_CLI_TOOLS_VERSION}_latest.zip"
URL="https://dl.google.com/android/repository/commandlinetools-linux-${ANDROID_CLI_TOOLS_VERSION}_latest.zip"

if command -v wget >/dev/null 2>&1; then
  wget -q -O "$ZIP" "$URL"
elif command -v curl >/dev/null 2>&1; then
  curl -fsSL -o "$ZIP" "$URL"
else
  echo "Need wget or curl" >&2
  exit 1
fi

unzip -q -o "$ZIP" -d "$INSTALL_ROOT/cmdline-tools"
rm -f "$ZIP"
if [[ -d "$INSTALL_ROOT/cmdline-tools/cmdline-tools" ]]; then
  rm -rf "$INSTALL_ROOT/cmdline-tools/latest"
  mv "$INSTALL_ROOT/cmdline-tools/cmdline-tools" "$INSTALL_ROOT/cmdline-tools/latest"
fi

export ANDROID_SDK_ROOT="$INSTALL_ROOT"
export ANDROID_HOME="$INSTALL_ROOT"
export PATH="${PATH}:${ANDROID_HOME}/cmdline-tools/latest/bin:${ANDROID_HOME}/platform-tools"

yes | sdkmanager --licenses
sdkmanager \
  "platform-tools" \
  "platforms;android-34" \
  "build-tools;34.0.0" \
  "ndk;${NDK_VERSION}"

echo "Done. Add to your shell profile:"
echo "  export ANDROID_HOME=$INSTALL_ROOT"
echo "  export ANDROID_SDK_ROOT=$INSTALL_ROOT"
echo "  export ANDROID_NDK_HOME=$INSTALL_ROOT/ndk/${NDK_VERSION}"
echo "  export PATH=\"\$PATH:\$ANDROID_HOME/cmdline-tools/latest/bin:\$ANDROID_HOME/platform-tools\""
