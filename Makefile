GOMOBILE := $(shell go env GOPATH)/bin/gomobile

.PHONY: test fmt lint vet clean help android-aar android-apk

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

fmt: ## Format code
	go fmt ./...

lint: ## Run linter
	golangci-lint run ./...

vet: ## Run go vet
	go vet ./...

check: fmt vet lint test ## Run all checks (format, vet, lint, test)

clean: ## Clean build artifacts
	go clean ./...
	rm -f coverage.out coverage.html

android-aar: ## Build android/app/libs/weathertowalk.aar via gomobile (needs ANDROID_HOME, NDK)
	@test -n "$$ANDROID_HOME" || (echo "Set ANDROID_HOME to your Android SDK path." && exit 1)
	@test -f "$(GOMOBILE)" || (echo "Install gomobile: go install golang.org/x/mobile/cmd/gomobile@latest && gomobile init" && exit 1)
	cd mobile && env GOTOOLCHAIN=auto $(GOMOBILE) bind -target=android -o ../android/app/libs/weathertowalk.aar .

android-apk: android-aar ## Build debug APK (android/app/build/outputs/apk/debug/app-debug.apk)
	cd android && ./gradlew assembleDebug
