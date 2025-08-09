# Lucien CLI Build System Makefile (Windows-safe)
.PHONY: all build clean test install bootstrap run dev help
.DEFAULT_GOAL := help

BINARY_NAME := lucien
MAIN_PATH := ./cmd/lucien
BUILD_DIR := ./build
VERSION := 1.0.0-nexus7
COMMIT := $(shell git rev-parse --short HEAD 2>NUL || echo unknown)
BUILD_TIME := $(shell powershell -Command "Get-Date -Format yyyyMMdd_HHmmss")

LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)
BUILD_FLAGS := -trimpath -mod=readonly -ldflags "$(LDFLAGS)"

ifeq ($(OS),Windows_NT)
    BINARY_EXT := .exe
    SHELL_CMD := powershell
    INSTALL_SCRIPT := scripts/install.ps1
else
    BINARY_EXT :=
    SHELL_CMD := bash
    INSTALL_SCRIPT := scripts/install.sh
endif

BINARY := $(BINARY_NAME)$(BINARY_EXT)
BUILD_BINARY := $(BUILD_DIR)/$(BINARY)

dev:
	@echo Running Lucien CLI in development mode...
	@go run $(MAIN_PATH)

run: build
	@echo Running built Lucien CLI...
	@$(BUILD_BINARY)

build:
	@echo Building Lucien CLI v$(VERSION)...
	@if not exist "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"
	@go build $(BUILD_FLAGS) -o $(BUILD_BINARY) $(MAIN_PATH)
	@echo Binary built: $(BUILD_BINARY)

build-all:
	@echo Building Lucien CLI for all platforms...
	@if not exist "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"
	@echo Building for Windows (amd64)...
	@GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo Building for Linux (amd64)...
	@GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@echo Building for Linux (arm64)...
	@GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	@echo Building for macOS (amd64)...
	@GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@echo Building for macOS (arm64)...
	@GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo All platform binaries built in $(BUILD_DIR)/

install: build
	@echo Installing Lucien CLI...
ifeq ($(OS),Windows_NT)
	@powershell -ExecutionPolicy Bypass -File scripts/install_simple.ps1
else
	@bash scripts/install.sh
endif
	@echo Installation complete! Run 'lucien --version' to verify

install-force: build
	@echo Force reinstalling Lucien CLI...
	@$(MAKE) uninstall 2>NUL || exit 0
	@$(MAKE) install
	@echo Force reinstall complete

install-dev: build
	@echo Installing development build...
	@if not exist "$(USERPROFILE)\.local\bin" mkdir "$(USERPROFILE)\.local\bin"
	@copy "$(BUILD_BINARY)" "$(USERPROFILE)\.local\bin\$(BINARY)" >NUL
	@echo Development build installed (no plugins/shortcuts)

uninstall:
	@echo Uninstalling Lucien CLI...
ifeq ($(OS),Windows_NT)
	@powershell -ExecutionPolicy Bypass -File scripts/uninstall.ps1
else
	@bash scripts/uninstall.sh
endif
	@echo Uninstallation complete

test:
	@echo Running tests...
	@go test -v ./...

test-agents: build
	@echo Running agent integration tests...
	@cd tests && LUCIEN_BINARY="../$(BUILD_BINARY)" go test -v -run TestAllAgents

test-shell: build
	@echo Running comprehensive shell tests...
	@cd tests && LUCIEN_BINARY="../$(BUILD_BINARY)" go test -v -run TestComprehensiveShell

test-all: build
	@echo Running all tests (unit + integration)...
	@go test -v ./...
	@cd tests && LUCIEN_BINARY="../$(BUILD_BINARY)" go test -v

benchmark:
	@echo Running benchmarks...
	@go test -bench=. -benchmem ./...

clean:
	@echo Cleaning build artifacts...
	@if exist "$(BUILD_DIR)" rmdir /s /q "$(BUILD_DIR)"
	@go clean
	@echo Clean complete

deps:
	@echo Updating dependencies...
	@go mod download
	@go mod tidy
	@echo Dependencies updated

fmt:
	@echo Formatting code...
	@go fmt ./...
	@echo Code formatted

vet:
	@echo Running go vet...
	@go vet ./...

version:
	@echo Lucien CLI Version: $(VERSION)
	@echo Commit: $(COMMIT)
	@echo Build Time: $(BUILD_TIME)
	@echo Go Version: $(shell go version)

help:
	@echo Lucien CLI Build System
	@echo
	@awk "BEGIN {FS = \":.*##\"; printf \"Usage: make \033[36m<target>\033[0m\n\n\"} /^[a-zA-Z_-]+:.*?##/ { printf \"  \033[36m%-15s\033[0m %s\n\", \$1, \$2 } /^##@/ { printf \"\n\033[1m%s\033[0m\n\", substr(\$0, 5) }" $(MAKEFILE_LIST)
	@echo
	@echo Examples:
	@echo   make build          # Build for current platform
	@echo   make test-all       # Run all tests
	@echo   make install        # Install using platform installer
	@echo   make release        # Build release for all platforms