# Dev environment: JDK 17, Android SDK (API 34), NDK — for gomobile + Gradle APK builds.
# Base image provides the `vscode` user used by Dev Containers.
FROM mcr.microsoft.com/devcontainers/base:jammy

USER root

RUN apt-get update && apt-get install -y --no-install-recommends \
    openjdk-17-jdk \
    unzip \
    wget \
    && rm -rf /var/lib/apt/lists/*

ENV ANDROID_SDK_ROOT=/opt/android-sdk
ENV ANDROID_HOME=/opt/android-sdk

# https://developer.android.com/studio#command-line-tools-only
ARG ANDROID_CLI_TOOLS_VERSION=11076708
RUN mkdir -p ${ANDROID_SDK_ROOT}/cmdline-tools \
    && wget -q "https://dl.google.com/android/repository/commandlinetools-linux-${ANDROID_CLI_TOOLS_VERSION}_latest.zip" -O /tmp/cmdline-tools.zip \
    && unzip -q /tmp/cmdline-tools.zip -d ${ANDROID_SDK_ROOT}/cmdline-tools \
    && mv ${ANDROID_SDK_ROOT}/cmdline-tools/cmdline-tools ${ANDROID_SDK_ROOT}/cmdline-tools/latest \
    && rm /tmp/cmdline-tools.zip

ENV PATH=${PATH}:${ANDROID_SDK_ROOT}/cmdline-tools/latest/bin:${ANDROID_SDK_ROOT}/platform-tools

# gomobile bind requires an NDK whose API metadata it accepts (21.4.x works; 26.x often fails with "unsupported API version").
ARG NDK_VERSION=21.4.7075529
RUN set -eux \
    && yes | sdkmanager --licenses \
    && sdkmanager \
        "platform-tools" \
        "platforms;android-34" \
        "build-tools;34.0.0" \
        "ndk;${NDK_VERSION}"

ENV ANDROID_NDK_HOME=${ANDROID_SDK_ROOT}/ndk/${NDK_VERSION}

RUN chown -R vscode:vscode ${ANDROID_SDK_ROOT}

USER vscode
