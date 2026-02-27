.PHONY: test fmt lint vet clean help

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
