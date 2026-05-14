.PHONY: help fmt lint lint-fix init upgrade install test test-coverage run clean gamut

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'


GOFMT_FILES?=$$(find . -name '*.go')

fmt: ## Format Go code
	gofmt -w -l $(GOFMT_FILES)

lint: ## Run golangci-lint
	golangci-lint run ./...

lint-fix: ## Run golangci-lint and fix issues
	golangci-lint run ./... --fix

init: ## Install Go
	brew install go@1.26
	brew install golangci-lint

upgrade: ## Upgrade Go dependencies
	go get -u ./...
	go mod tidy

install: ## Install Go Go Gadget
	GOBIN=$(HOME)/go/bin/ go install -v ./...

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -coverprofile=coverage.out ./...

run: ## Run Go Go Gadget
	go run main.go

clean: ## Clean Go Go Gadget
	rm -f go-go-gadget

gamut: ## Run the Gamut
	$(MAKE) fmt
	$(MAKE) init
	$(MAKE) upgrade
	$(MAKE) lint
	$(MAKE) run
	$(MAKE) test
	$(MAKE) test-coverage
	$(MAKE) install
