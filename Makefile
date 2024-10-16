.PHONY: clean run build install dep test lint format docker
# Set the bin path
PATHINSTBIN = $(abspath ./bin)
export PATH := $(PATHINSTBIN):$(PATH)

BIN_NAME					?= migration
DEFAULT_INSTALL_DIR			:= $(go env GOPATH)/bin
DEFAULT_ARCH				:= $(shell go env GOARCH)
DEFAULT_GOOS				:= $(shell go env GOOS)
ARCH						?= $(DEFAULT_ARCH)
GOOS						?= $(DEFAULT_GOOS)
INSTALL_DIR					?= $(DEFAULT_INSTALL_DIR)
.DEFAULT_GOAL 				:= run


VERSION   := $(shell git describe --tags || echo "v0.0.0")
VER_CUT   := $(shell echo $(VERSION) | cut -c2-)

# Dependency versions
GOLANGCI_VERSION   = latest

help:
	@echo "\nSpecify a subcommand:\n"
	@grep -hE '^[0-9a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[0;36m%-20s\033[m %s\n", $$1, $$2}'
	@echo ""


build: ## Build the migration binary
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(ARCH) \
		go build -o $(PATHINSTBIN)/$(BIN_NAME) ./cmd/$(BIN_NAME)


all: clean target

clean: ## Clean the project binaries
	@rm -rf $(PATHINSTBIN)
	

tidy:  ## tidy the go modules
	@go mod tidy

test: ## Run the all tests
	@go test ./...

lint: ## Run the linter
	@golangci-lint run --timeout 5m

format: ## Run the linter with fix
	@golangci-lint run --fix

tools-golangci-lint: ## Install golangci-lint
	@mkdir -p $(PATHINSTBIN)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | BINARY=golangci-lint bash -s -- ${GOLANGCI_VERSION}

tools: tools-golangci-lint  ## Install all tools