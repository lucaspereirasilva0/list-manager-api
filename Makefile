# Variables
APP_NAME=lista-backend-api
MAIN_PATH=cmd/api/main.go
BUILD_DIR=build

# Go commands
GOCMD=go
GOPATH=$(shell go env GOPATH)

# Check if tools are installed
GOLANGCI_LINT=$(GOPATH)/bin/golangci-lint
GOIMPORTS=$(GOPATH)/bin/goimports

.PHONY: all build run test clean lint fmt deps help install-lint check-lint coverage

all: deps fmt lint test coverage build clean ## Update deps, format, lint, test, show coverage and build

build: ## Build the project
	@echo "Building..."
	$(GOCMD) build -o /dev/null $(MAIN_PATH)

run: ## Run the project
	@echo "Running..."
	$(GOCMD) run $(MAIN_PATH)

test: ## Run tests and generate coverage report
	@echo "Testing and generating coverage report..."
	@$(GOCMD) test -v -coverprofile=coverage.out ./...
	@grep -v "mock.go" coverage.out > coverage.tmp
	@mv coverage.tmp coverage.out

coverage: ## Show coverage report
	@echo "Generating coverage report..."
	@$(GOCMD) tool cover -func=coverage.out | grep total | awk '{print "Project coverage: " $$3}'

clean: ## Remove generated files
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out
	@rm -f coverage.tmp

check-lint: ## Check if linters are installed
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2; \
	}
	@command -v $(GOIMPORTS) >/dev/null 2>&1 || { \
		echo "Installing goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	}

lint: check-lint ## Run linter
	@echo "Linting..."
	$(GOLANGCI_LINT) run --timeout=5m ./...

fmt: check-lint ## Format Go files
	@echo "Formatting..."
	$(GOCMD) fmt ./...
	@echo "Running goimports..."
	$(GOIMPORTS) -w .

deps: ## Update dependencies
	@echo "Updating dependencies..."
	$(GOCMD) mod tidy

help: ## Show this help message
	@echo "Usage:"
	@echo "  make \033[36m<target>\033[0m"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'