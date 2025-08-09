# 🏗️ LUCIEN CLI - BUILD AND DEPLOYMENT GUIDE
**Version**: 1.0-alpha  
**Last Updated**: 2025-08-06  
**Status**: Development Procedures  

## ⚠️ IMPORTANT SECURITY NOTICE

**THIS GUIDE IS FOR DEVELOPMENT PURPOSES ONLY**

The current version of Lucien CLI contains critical security vulnerabilities and is **NOT APPROVED FOR PRODUCTION DEPLOYMENT**. This guide documents build procedures for development and testing environments only.

**Production deployment is PROHIBITED until all security issues are resolved.**

## 📚 TABLE OF CONTENTS

1. [Prerequisites](#prerequisites)
2. [Development Environment Setup](#development-environment-setup)
3. [Building from Source](#building-from-source)
4. [Testing Procedures](#testing-procedures)
5. [Package Creation](#package-creation)
6. [Development Deployment](#development-deployment)
7. [CI/CD Pipeline](#cicd-pipeline)
8. [Security Considerations](#security-considerations)
9. [Troubleshooting](#troubleshooting)

## 🔧 PREREQUISITES

### System Requirements

#### Development Machine
- **OS**: Windows 10/11 (primary), Linux, macOS (experimental)
- **RAM**: 8GB minimum, 16GB recommended
- **Storage**: 5GB free space for build artifacts
- **Network**: Internet connection for dependency downloads

#### Software Dependencies
```bash
# Required
Go 1.21 or higher
Python 3.11 or higher
Git 2.30 or higher

# Recommended
Make (GNU Make 4.0+)
Docker (for containerized builds)
Visual Studio Code or similar IDE
```

#### Development Tools
```bash
# Go tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/godoc@latest

# Python tools
pip install black isort pytest coverage
pip install rich typer pyyaml  # Runtime dependencies
```

### Environment Validation

#### Pre-Build Checklist
```bash
# Verify Go installation
go version  # Should be 1.21+

# Verify Python installation  
python --version  # Should be 3.11+

# Verify Git installation
git --version  # Should be 2.30+

# Verify Make installation (optional)
make --version  # Should be 4.0+

# Check available disk space
df -h .  # Ensure 5GB+ available
```

## 🛠️ DEVELOPMENT ENVIRONMENT SETUP

### Initial Repository Setup

#### 1. Clone Repository
```bash
# Clone from repository
git clone https://github.com/luciendev/lucien-core.git
cd lucien-core

# Verify repository structure
ls -la
# Expected: cli/, cmd/, internal/, plugins/, etc.
```

#### 2. Go Environment Setup
```bash
# Initialize Go modules
go mod tidy
go mod download
go mod verify

# Verify Go environment
go env GOPATH
go env GOROOT
go env GOOS
go env GOARCH
```

#### 3. Python Environment Setup
```bash
# Create virtual environment
python -m venv .venv

# Activate virtual environment
# Windows:
.venv\Scripts\activate
# Linux/Mac:
source .venv/bin/activate

# Install Python dependencies
pip install -r requirements.txt

# Verify Python environment
python -c "import lucien; print('Python environment OK')"
```

### IDE Configuration

#### Visual Studio Code Setup
```json
// .vscode/settings.json
{
    "go.toolsManagement.checkForUpdates": "local",
    "go.useLanguageServer": true,
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "python.pythonPath": ".venv/bin/python",
    "python.formatting.provider": "black"
}
```

#### Go IDE Configuration
```json
// .vscode/launch.json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Lucien",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/lucien",
            "args": ["--debug", "--safe"],
            "env": {
                "LUCIEN_DEBUG": "1"
            }
        }
    ]
}
```

## 🔨 BUILDING FROM SOURCE

### Build Methods

#### Option 1: Makefile Build (Recommended)
```bash
# Clean previous builds
make clean

# Build all components
make all

# Build specific components
make build        # Go binary only
make python       # Python components only
make plugins      # Plugin system only

# Build with debugging
make build-debug

# Create release package
make release
```

#### Option 2: Manual Build Process
```bash
# Build Go components
cd cmd/lucien
go build -o ../../build/lucien.exe -v .
cd ../..

# Build plugins
cd plugins/example-bmad
go build -o ../../build/plugins/bmad.exe .
cd ../..

# Verify build
./build/lucien.exe --version
```

#### Option 3: Cross-Platform Build
```bash
# Build for multiple platforms
make build-all

# Manual cross-compilation
GOOS=windows GOARCH=amd64 go build -o build/lucien-windows-amd64.exe cmd/lucien/main.go
GOOS=linux GOARCH=amd64 go build -o build/lucien-linux-amd64 cmd/lucien/main.go
GOOS=darwin GOARCH=amd64 go build -o build/lucien-darwin-amd64 cmd/lucien/main.go
```

### Build Configuration

#### Makefile Overview
```makefile
# Key Make targets
build:          # Build main binary
build-debug:    # Build with debug symbols
build-all:      # Cross-platform build
test:           # Run test suite
clean:          # Clean build artifacts
fmt:            # Format source code
lint:           # Run linters
release:        # Create release package
```

#### Build Flags and Options
```bash
# Production build
go build -ldflags="-s -w" -o lucien.exe cmd/lucien/main.go

# Debug build
go build -gcflags="all=-N -l" -o lucien-debug.exe cmd/lucien/main.go

# Static linking (Linux)
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o lucien-static cmd/lucien/main.go
```

### Build Artifacts

#### Expected Output Structure
```
build/
├── lucien.exe                    # Main Windows binary
├── lucien                        # Main Linux binary
├── lucien.app/                   # Main macOS app bundle
├── plugins/
│   ├── bmad.exe                  # BMAD plugin
│   └── example-weather.exe       # Weather plugin
├── configs/
│   ├── default.toml              # Default configuration
│   └── safe-mode.toml            # Safe mode configuration
└── docs/
    ├── README.txt                # Basic documentation
    └── CHANGELOG.txt             # Version history
```

## 🧪 TESTING PROCEDURES

### Pre-Deployment Testing

#### Unit Testing
```bash
# Run all Go unit tests
go test ./... -v

# Run tests with coverage
go test ./... -cover -coverprofile=coverage.out

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html

# Run Python tests
pytest tests/ -v --cov=lucien
```

#### Integration Testing
```bash
# Shell functionality tests
go test ./internal/shell/ -tags=integration -v

# Plugin system tests
go test ./internal/plugin/ -tags=integration -v

# Full system integration test
make test-integration
```

#### Security Testing
```bash
# Security vulnerability scan
make security-scan

# Static analysis
golangci-lint run --enable-all

# Dependency vulnerability scan
go list -json -m all | nancy sleuth
```

#### Performance Testing
```bash
# Benchmark tests
go test -bench=. ./...

# Memory profiling
go test -memprofile=mem.prof ./internal/shell/
go tool pprof mem.prof

# CPU profiling
go test -cpuprofile=cpu.prof ./internal/shell/
go tool pprof cpu.prof
```

### Testing Checklist

#### ✅ Pre-Build Testing
- [ ] All unit tests pass
- [ ] Code formatting verified (`make fmt`)
- [ ] Linting passes (`make lint`)
- [ ] Dependencies updated (`go mod tidy`)
- [ ] No TODO/FIXME in critical paths

#### ✅ Post-Build Testing
- [ ] Binary executes successfully
- [ ] `--version` flag works
- [ ] `--help` flag displays correctly
- [ ] Basic shell commands functional
- [ ] Plugin loading works
- [ ] Safe mode activates properly

#### ✅ Integration Testing
- [ ] Command parsing works correctly
- [ ] Variable expansion functional (⚠️ Known issues)
- [ ] Plugin execution successful
- [ ] Error handling appropriate
- [ ] Cross-platform compatibility verified

## 📦 PACKAGE CREATION

### Development Package

#### Creating Distribution Package
```bash
# Create complete distribution
make package

# Manual package creation
mkdir -p dist/lucien-cli
cp build/lucien.exe dist/lucien-cli/
cp -r build/plugins dist/lucien-cli/
cp -r configs dist/lucien-cli/
cp README.md LICENSE dist/lucien-cli/
cd dist && zip -r lucien-cli-dev.zip lucien-cli/
```

#### Package Contents Verification
```bash
# Verify package contents
unzip -l lucien-cli-dev.zip

# Expected contents:
# lucien-cli/
# ├── lucien.exe
# ├── plugins/
# ├── configs/
# ├── README.md
# └── LICENSE
```

### Installer Creation

#### Windows Installer (NSIS)
```nsis
; installer.nsi
Name "Lucien CLI"
OutFile "lucien-cli-installer.exe"
InstallDir $PROGRAMFILES\LucienCLI

Section "Main"
    SetOutPath $INSTDIR
    File "build\lucien.exe"
    File /r "build\plugins"
    File /r "configs"
    
    WriteUninstaller $INSTDIR\uninstall.exe
SectionEnd
```

#### Linux Package (DEB)
```bash
# Create Debian package structure
mkdir -p lucien-cli-deb/DEBIAN
mkdir -p lucien-cli-deb/usr/bin
mkdir -p lucien-cli-deb/etc/lucien

# Copy files
cp build/lucien lucien-cli-deb/usr/bin/
cp -r configs lucien-cli-deb/etc/lucien/

# Create control file
cat > lucien-cli-deb/DEBIAN/control << EOF
Package: lucien-cli
Version: 1.0-alpha
Architecture: amd64
Maintainer: Lucien Dev Team
Description: AI-enhanced shell replacement
EOF

# Build package
dpkg-deb --build lucien-cli-deb
```

## 🚀 DEVELOPMENT DEPLOYMENT

### ⚠️ IMPORTANT: Development Only

**This section is for development and testing environments ONLY. Production deployment is not authorized due to security vulnerabilities.**

### Local Development Deployment

#### Direct Execution
```bash
# Run from build directory
./build/lucien.exe

# Run with safe mode
./build/lucien.exe --safe

# Run with debugging
LUCIEN_DEBUG=1 ./build/lucien.exe --debug
```

#### System Installation (Development)
```bash
# Windows (Run as Administrator)
copy build\lucien.exe C:\Windows\System32\
copy build\plugins\*.exe C:\ProgramData\Lucien\plugins\

# Linux (Development only)
sudo cp build/lucien /usr/local/bin/
sudo mkdir -p /etc/lucien/plugins
sudo cp build/plugins/* /etc/lucien/plugins/

# macOS (Development only)  
sudo cp build/lucien /usr/local/bin/
sudo mkdir -p /etc/lucien/plugins
sudo cp build/plugins/* /etc/lucien/plugins/
```

### Development Environment Configuration

#### Configuration File Setup
```bash
# Create user configuration directory
mkdir -p ~/.lucien

# Copy default configuration
cp configs/default.toml ~/.lucien/config.toml

# Edit configuration for development
cat > ~/.lucien/config.toml << EOF
[shell]
safe_mode = true
debug_mode = true

[plugins]
auto_load = ["bmad"]
plugin_dir = "~/.lucien/plugins"

[development]
log_level = "debug"
enable_profiling = true
EOF
```

#### Environment Variables
```bash
# Development environment setup
export LUCIEN_ENV=development
export LUCIEN_DEBUG=1
export LUCIEN_CONFIG=~/.lucien/config.toml
export LUCIEN_PLUGIN_DIR=~/.lucien/plugins
```

### Container Deployment (Development)

#### Docker Development Container
```dockerfile
# Dockerfile.dev
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN make build

FROM python:3.11-alpine
RUN apk add --no-cache bash
WORKDIR /app
COPY --from=builder /app/build/lucien .
COPY cli/ cli/
COPY lucien/ lucien/
COPY requirements.txt .
RUN pip install -r requirements.txt
CMD ["./lucien", "--safe"]
```

```bash
# Build development container
docker build -f Dockerfile.dev -t lucien-cli:dev .

# Run development container
docker run -it --rm lucien-cli:dev
```

## 🔄 CI/CD PIPELINE

### GitHub Actions Workflow

#### Build Pipeline
```yaml
# .github/workflows/build.yml
name: Build and Test
on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: 1.21
    - uses: actions/setup-python@v3
      with:
        python-version: 3.11
    
    - name: Install dependencies
      run: |
        go mod download
        pip install -r requirements.txt
    
    - name: Run tests
      run: make test
    
    - name: Build binary
      run: make build
    
    - name: Security scan
      run: make security-scan
    
    - name: Upload artifacts
      uses: actions/upload-artifact@v3
      with:
        name: lucien-cli
        path: build/
```

#### Security Pipeline
```yaml
# .github/workflows/security.yml
name: Security Scan
on: [push, pull_request]

jobs:
  security:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Run Gosec Security Scanner
      uses: securecodewarrior/github-action-gosec@master
    - name: Run Nancy Vulnerability Scanner
      run: |
        go list -json -m all | nancy sleuth
```

### Build Pipeline Stages

#### Stage 1: Source Validation
- Code formatting verification
- Lint checking
- Dependency validation
- Security scanning

#### Stage 2: Testing
- Unit test execution
- Integration test execution
- Performance benchmarking
- Coverage reporting

#### Stage 3: Build
- Cross-platform compilation
- Binary optimization
- Plugin compilation
- Asset bundling

#### Stage 4: Packaging
- Distribution package creation
- Installer generation
- Container image building
- Documentation generation

#### Stage 5: Deployment (Development Only)
- Development environment deployment
- Testing environment deployment
- Container registry publishing
- Artifact archiving

## 🔒 SECURITY CONSIDERATIONS

### Build Security

#### Secure Build Environment
```bash
# Verify build environment
echo "Verifying build security..."

# Check for malicious dependencies
go mod verify

# Scan for known vulnerabilities  
nancy sleuth < go.list

# Verify checksums
sha256sum build/lucien.exe > build/lucien.exe.sha256
```

#### Supply Chain Security
```bash
# Generate Software Bill of Materials (SBOM)
go list -json -m all > build/sbom.json

# Sign build artifacts (when ready for production)
# gpg --armor --detach-sign build/lucien.exe

# Verify signatures (when ready for production)
# gpg --verify build/lucien.exe.asc build/lucien.exe
```

### Deployment Security

#### Security Hardening Checklist
- [ ] All dependencies verified and up-to-date
- [ ] Security vulnerabilities scanned and addressed
- [ ] Build artifacts signed and verified
- [ ] Secure configuration defaults applied
- [ ] Runtime security measures enabled

#### Security Monitoring
```bash
# Enable security logging
export LUCIEN_SECURITY_LOG=1
export LUCIEN_AUDIT_LOG=/var/log/lucien-audit.log

# Monitor for security events
tail -f /var/log/lucien-audit.log | grep SECURITY
```

## 🔧 TROUBLESHOOTING

### Common Build Issues

#### Go Build Failures
```bash
# Problem: Module verification failed
# Solution: 
go clean -modcache
go mod download
go mod verify

# Problem: CGO errors on cross-compilation
# Solution:
CGO_ENABLED=0 go build cmd/lucien/main.go

# Problem: Missing dependencies
# Solution:
go mod tidy
go get -u all
```

#### Python Build Issues
```bash
# Problem: Missing Python dependencies
# Solution:
pip install --upgrade pip
pip install -r requirements.txt

# Problem: Virtual environment issues
# Solution:
rm -rf .venv
python -m venv .venv
source .venv/bin/activate  # Linux/Mac
.venv\Scripts\activate     # Windows
```

#### Plugin Build Failures
```bash
# Problem: Plugin RPC interface mismatch
# Solution:
cd plugins/example-bmad
go mod tidy
go build -o ../../build/plugins/bmad.exe .

# Problem: Plugin manifest errors
# Solution:
validate-json manifest.json
```

### Performance Issues

#### Slow Build Times
```bash
# Enable build caching
export GOCACHE=/tmp/go-build
export GOMODCACHE=/tmp/go-mod

# Parallel builds
make -j$(nproc) build

# Docker build optimization
docker build --build-arg BUILDKIT_INLINE_CACHE=1 .
```

#### Large Binary Size
```bash
# Build with size optimization
go build -ldflags="-s -w" cmd/lucien/main.go

# Strip debug information
strip build/lucien.exe

# Use UPX compression (with caution)
upx --best build/lucien.exe
```

### Runtime Issues

#### Startup Problems
```bash
# Check dependencies
ldd build/lucien  # Linux
otool -L build/lucien  # macOS

# Verify permissions
chmod +x build/lucien

# Check configuration
./build/lucien --validate-config
```

#### Plugin Loading Issues
```bash
# Debug plugin loading
LUCIEN_DEBUG=1 ./build/lucien --debug

# Verify plugin permissions
ls -la build/plugins/

# Test plugin manually
./build/plugins/bmad.exe --test
```

## 📋 DEPLOYMENT CHECKLIST

### Pre-Deployment Verification

#### ✅ Build Verification
- [ ] All components build successfully
- [ ] All tests pass
- [ ] Security scans complete
- [ ] Performance benchmarks acceptable
- [ ] Cross-platform compatibility verified

#### ✅ Package Verification  
- [ ] Distribution package created
- [ ] Package contents verified
- [ ] Installer tested
- [ ] Documentation included
- [ ] Checksums generated

#### ✅ Security Verification
- [ ] **CRITICAL**: All security vulnerabilities addressed ⚠️
- [ ] Build artifacts signed
- [ ] Dependencies verified
- [ ] Configuration hardened
- [ ] Audit logging enabled

#### ✅ Documentation Verification
- [ ] Installation guide updated
- [ ] User manual current
- [ ] API documentation complete
- [ ] Troubleshooting guide available
- [ ] Security guide provided

## 🚫 PRODUCTION DEPLOYMENT PROHIBITION

### ⛔ CRITICAL WARNING

**PRODUCTION DEPLOYMENT IS STRICTLY PROHIBITED**

The current version of Lucien CLI contains:
- 19 Critical security vulnerabilities
- Policy engine compilation failures  
- Incomplete sandbox implementation
- Multiple high-risk security issues

**Deployment authorization will only be granted after:**
1. All critical vulnerabilities resolved
2. Comprehensive security testing passed
3. Third-party security audit completed
4. Security clearance obtained from security team

### When Production Ready

#### Future Production Deployment (When Authorized)
```bash
# Production build (when security cleared)
make production-build

# Production configuration
cp configs/production.toml /etc/lucien/config.toml

# Production deployment
sudo systemctl enable lucien
sudo systemctl start lucien

# Monitor production deployment
journalctl -u lucien -f
```

---

**REMEMBER**: This is development documentation only. Production deployment requires security clearance and vulnerability remediation.

**Status**: 🔶 DEVELOPMENT ONLY - NOT PRODUCTION READY

*Next Steps: Address all security vulnerabilities identified in Security Assessment Summary before considering any production deployment.*