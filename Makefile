.PHONY: help fmt init upgrade install test run clean gamut

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'


GOFMT_FILES?=$$(find . -name '*.go')

fmt: ## Format Go code
	gofmt -w -l $(GOFMT_FILES)

init: ## Install Go
	brew install go@1.25

upgrade: ## Upgrade Go dependencies
	go get -u ./...
	go mod tidy

install: ## Install Go Go Gadget
	GOBIN=$(HOME)/go/bin/ go install -v ./...

test: ## Run tests
	go test -v ./...

run: ## Run Go Go Gadget
	go run main.go

clean: ## Clean Go Go Gadget
	rm -f go-go-gadget

gumut: ## Run the Gumut
	$(MAKE) fmt
	$(MAKE) init
	$(MAKE) upgrade
	$(MAKE) run
	$(MAKE) test
	$(MAKE) install
