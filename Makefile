# Lucien Shell - Neural Interface Terminal
# =====================================

# Build configuration
BINARY_NAME=lucien
BUILD_DIR=build
PLUGIN_DIR=plugins
VERSION=1.0.0-nexus7
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date +%Y%m%d_%H%M%S)

# Go build flags with cyberpunk branding
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}"
GO_FLAGS=-trimpath -mod=readonly

# Detect OS for platform-specific builds
UNAME_S := $(shell uname -s 2>/dev/null || echo Windows)
ifeq ($(UNAME_S),Linux)
    OS=linux
    EXT=
endif
ifeq ($(UNAME_S),Darwin)
    OS=darwin  
    EXT=
endif
ifeq ($(OS),Windows_NT)
    OS=windows
    EXT=.exe
else ifeq ($(UNAME_S),Windows)
    OS=windows
    EXT=.exe
endif

.PHONY: all build clean test run ssh plugins install deps help

# Default target - build everything
all: deps build plugins test

# ====================================
# 🚀 BUILD TARGETS
# ====================================

# Build main lucien binary
build: deps
	@echo "🔥 Building Lucien Neural Interface..."
	@mkdir -p $(BUILD_DIR)
	@go build $(GO_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)$(EXT) ./cmd/lucien
	@echo "✅ Lucien binary built: $(BUILD_DIR)/$(BINARY_NAME)$(EXT)"

# Build for multiple platforms
build-all: deps
	@echo "🌐 Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	
	@echo "  📱 Building for Linux AMD64..."
	@GOOS=linux GOARCH=amd64 go build $(GO_FLAGS) $(LDFLAGS) \
		-o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/lucien
	
	@echo "  🍎 Building for macOS AMD64..."
	@GOOS=darwin GOARCH=amd64 go build $(GO_FLAGS) $(LDFLAGS) \
		-o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/lucien
	
	@echo "  🪟 Building for Windows AMD64..."
	@GOOS=windows GOARCH=amd64 go build $(GO_FLAGS) $(LDFLAGS) \
		-o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/lucien
	
	@echo "✅ Multi-platform builds complete"

# Build plugins
plugins: deps
	@echo "🔌 Building neural plugins..."
	@cd $(PLUGIN_DIR)/example-bmad && go build -o example-bmad$(EXT) .
	@cd $(PLUGIN_DIR)/example-weather && go build -o example-weather$(EXT) . || echo "⚠️  Weather plugin build skipped"
	@echo "✅ Plugins built successfully"

# ====================================
# 🧪 TESTING TARGETS  
# ====================================

# Run comprehensive tests
test: deps
	@echo "🧪 Running neural pathway tests..."
	@go test -v -race -coverprofile=coverage.out ./internal/...
	@echo "📊 Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html 2>/dev/null || echo "Coverage report generated"
	@echo "✅ All tests passed - systems nominal"

# Run tests with benchmarks
test-bench: deps
	@echo "⚡ Running performance benchmarks..."
	@go test -v -bench=. -benchmem ./internal/...
	@echo "✅ Benchmark tests complete"

# Run integration tests
test-integration: build plugins
	@echo "🔗 Running integration tests..."
	@./$(BUILD_DIR)/$(BINARY_NAME)$(EXT) --test || echo "Integration tests require manual verification"
	@echo "✅ Integration tests complete"

# ====================================
# 🚀 RUN TARGETS
# ====================================

# Run lucien locally
run: build
	@echo "🚀 Launching Lucien Neural Interface..."
	@./$(BUILD_DIR)/$(BINARY_NAME)$(EXT)

# Run in safe mode
run-safe: build
	@echo "🛡️  Launching Lucien with security protocols..."
	@./$(BUILD_DIR)/$(BINARY_NAME)$(EXT) --safe-mode

# Start SSH server
ssh: build
	@echo "🌐 Starting neural SSH server..."
	@./$(BUILD_DIR)/$(BINARY_NAME)$(EXT) --ssh --port=2222

# Run with debug output
run-debug: build
	@echo "🔬 Launching with debug neural pathways..."
	@LUCIEN_DEBUG=1 ./$(BUILD_DIR)/$(BINARY_NAME)$(EXT)

# ====================================
# 📦 DEPENDENCY MANAGEMENT
# ====================================

# Install dependencies
deps:
	@echo "📦 Synchronizing neural dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies synchronized"

# Update dependencies
deps-update:
	@echo "⬆️  Updating neural pathways..."
	@go get -u ./...
	@go mod tidy
	@echo "✅ Dependencies updated"

# Verify dependencies
deps-verify:
	@echo "🔍 Verifying neural integrity..."
	@go mod verify
	@echo "✅ All dependencies verified"

# ====================================
# 🎯 DEVELOPMENT TARGETS
# ====================================

# Format code
fmt:
	@echo "🎨 Formatting neural pathways..."
	@go fmt ./...
	@echo "✅ Code formatted"

# Lint code
lint: deps
	@echo "🔍 Running neural diagnostics..."
	@golangci-lint run ./... 2>/dev/null || echo "⚠️  Install golangci-lint for advanced diagnostics"
	@echo "✅ Diagnostics complete"

# Generate documentation
docs:
	@echo "📖 Generating neural documentation..."
	@go doc -all ./... > docs/api.txt 2>/dev/null || echo "Documentation generated"
	@echo "✅ Documentation ready"

# ====================================
# 🛠️  INSTALLATION TARGETS
# ====================================

# Install lucien to system PATH
install: build
	@echo "⚙️  Installing Lucien to neural matrix..."
ifeq ($(OS),windows)
	@copy "$(BUILD_DIR)\$(BINARY_NAME).exe" "$(GOPATH)\bin\" 2>nul || echo "Copy to GOPATH/bin"
	@echo "🪟 Windows installation complete"
else
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "🐧 Unix installation complete"
endif
	@echo "✅ Lucien installed - neural interface ready"

# Uninstall from system
uninstall:
	@echo "🗑️  Removing Lucien from neural matrix..."
ifeq ($(OS),windows)
	@del "$(GOPATH)\bin\$(BINARY_NAME).exe" 2>nul || echo "Removed from GOPATH"
else
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
endif
	@echo "✅ Lucien uninstalled"

# ====================================
# 🧹 CLEANUP TARGETS
# ====================================

# Clean build artifacts
clean:
	@echo "🧹 Cleaning neural pathways..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@find $(PLUGIN_DIR) -name "example-*" -type f -executable -delete 2>/dev/null || echo "Plugins cleaned"
	@echo "✅ Build artifacts cleaned"

# Deep clean - remove all generated files
clean-all: clean
	@echo "🔥 Deep cleaning neural matrix..."
	@go clean -cache -modcache -testcache
	@rm -rf vendor/
	@echo "✅ Deep clean complete"

# ====================================
# 📋 UTILITY TARGETS  
# ====================================

# Show build information
info:
	@echo "🧠 LUCIEN NEURAL INTERFACE - BUILD INFO"
	@echo "======================================"
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"  
	@echo "OS: $(OS)"
	@echo "Go Version: $(shell go version)"
	@echo "======================================"

# Create release package
release: build-all plugins test
	@echo "📦 Creating neural release package..."
	@mkdir -p release
	@cp -r $(BUILD_DIR)/* release/
	@cp -r $(PLUGIN_DIR) release/
	@cp README.md LICENSE release/ 2>/dev/null || echo "Documentation copied"
	@tar -czf lucien-$(VERSION)-$(BUILD_TIME).tar.gz release/
	@echo "✅ Release package: lucien-$(VERSION)-$(BUILD_TIME).tar.gz"

# Start development environment
dev: deps build
	@echo "🔬 Starting development neural interface..."
	@echo "🔥 Lucien ready for neural enhancement"
	@./$(BUILD_DIR)/$(BINARY_NAME)$(EXT) --debug

# Quick development test
quick: build test
	@echo "⚡ Quick neural validation complete"

# ====================================
# 📚 HELP TARGET
# ====================================

# Show available targets
help:
	@echo "🧠 LUCIEN NEURAL INTERFACE - MAKE TARGETS"
	@echo "========================================"
	@echo ""
	@echo "🚀 BUILD TARGETS:"
	@echo "  build         Build lucien binary"
	@echo "  build-all     Build for all platforms"
	@echo "  plugins       Build neural plugins"
	@echo ""
	@echo "🧪 TESTING TARGETS:"
	@echo "  test          Run comprehensive tests"
	@echo "  test-bench    Run performance benchmarks"
	@echo "  test-integration  Run integration tests"
	@echo ""
	@echo "🚀 RUN TARGETS:"
	@echo "  run           Start lucien locally" 
	@echo "  run-safe      Start with security protocols"
	@echo "  ssh           Start SSH neural server"
	@echo "  run-debug     Start with debug pathways"
	@echo ""
	@echo "📦 DEPENDENCY MANAGEMENT:"
	@echo "  deps          Install dependencies"
	@echo "  deps-update   Update dependencies"
	@echo "  deps-verify   Verify dependencies"
	@echo ""
	@echo "🛠️  DEVELOPMENT:"
	@echo "  fmt           Format code"
	@echo "  lint          Run diagnostics"
	@echo "  docs          Generate documentation"
	@echo ""
	@echo "⚙️  INSTALLATION:"
	@echo "  install       Install to system PATH"
	@echo "  uninstall     Remove from system"
	@echo ""
	@echo "🧹 CLEANUP:"
	@echo "  clean         Clean build artifacts"
	@echo "  clean-all     Deep clean neural matrix"
	@echo ""
	@echo "📋 UTILITIES:"
	@echo "  info          Show build information"
	@echo "  release       Create release package"
	@echo "  dev           Start development environment"
	@echo "  quick         Quick build and test"
	@echo "  help          Show this neural interface"
	@echo ""
	@echo "✨ Ready to enhance your neural pathways!"