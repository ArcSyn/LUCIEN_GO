# 🛠️ LUCIEN CLI - TECHNICAL MANUAL
**Version**: 1.0-alpha  
**Last Updated**: 2025-08-06  
**Status**: Development Documentation  

## 📚 TABLE OF CONTENTS

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Installation Guide](#installation-guide)
4. [User Guide](#user-guide)
5. [Developer Guide](#developer-guide)
6. [API Reference](#api-reference)
7. [Plugin Development](#plugin-development)
8. [Configuration](#configuration)
9. [Troubleshooting](#troubleshooting)

## 🎯 OVERVIEW

### What is Lucien CLI?

Lucien CLI is an experimental AI-enhanced shell replacement written in Go with Python components. It combines traditional command-line functionality with modern plugin architecture and cyberpunk-themed user interface.

**⚠️ IMPORTANT**: Current version is NOT production-ready. This manual documents both implemented features and known limitations.

### Key Components
- **Go Core**: Shell engine, plugin system, security policies
- **Python CLI**: User interface layer and utility functions
- **Plugin System**: HashiCorp go-plugin based extensibility
- **Security Layer**: Policy engine with safe mode operation

## 🏗️ ARCHITECTURE

### System Architecture Diagram
```
┌─────────────────────────────────────────────────────────────────┐
│                          LUCIEN CLI                            │
├─────────────────────────────────────────────────────────────────┤
│ Python Layer (cli/main.py, lucien/ui.py)                       │
│ ├── User Interface (Rich-based TUI)                            │
│ ├── Command Dispatcher                                         │
│ └── Safety Manager                                             │
├─────────────────────────────────────────────────────────────────┤
│ Go Core (cmd/lucien/main.go)                                   │
│ ├── Shell Engine (internal/shell/)                            │
│ ├── Plugin Manager (internal/plugin/)                         │
│ ├── Policy Engine (internal/policy/)                          │
│ ├── Sandbox Manager (internal/sandbox/)                       │
│ └── AI Engine (internal/ai/)                                  │
├─────────────────────────────────────────────────────────────────┤
│ Plugin Ecosystem (plugins/)                                    │
│ └── BMAD Example Plugin                                        │
└─────────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

**Python Layer**:
- User interface and interaction
- Command parsing and initial validation
- Safety checks and dangerous command detection
- Rich-based terminal formatting

**Go Core**:
- Shell command execution
- Variable and alias management
- Plugin lifecycle management
- Security policy enforcement
- Process sandboxing

**Plugin System**:
- RPC-based plugin communication
- Isolated plugin execution
- Dynamic plugin loading/unloading

## 🔧 INSTALLATION GUIDE

### Prerequisites

**For Users**:
- Python 3.11 or higher
- Windows 10/11 (primary platform)
- 100MB free disk space

**For Developers**:
- Go 1.21 or higher
- Python 3.11+
- Git
- Make (optional but recommended)

### Installation Steps

#### Option 1: Binary Installation (Recommended)
```bash
# Download latest release
wget https://github.com/luciendev/lucien/releases/latest/lucien.exe

# Make executable (Windows)
# No additional steps needed

# Run Lucien
./lucien.exe
```

#### Option 2: Source Installation
```bash
# Clone repository
git clone https://github.com/luciendev/lucien-core.git
cd lucien-core

# Install Python dependencies
pip install -r requirements.txt

# Build Go components
make build

# Run from source
make run
```

#### Option 3: Development Setup
```bash
# Full development environment
git clone https://github.com/luciendev/lucien-core.git
cd lucien-core

# Setup Python virtual environment
python -m venv .venv
source .venv/bin/activate  # Linux/Mac
# OR
.venv\Scripts\activate     # Windows

# Install dependencies
pip install -r requirements.txt
go mod download

# Run tests
make test

# Start development server
make dev
```

## 👤 USER GUIDE

### Basic Usage

#### Starting Lucien
```bash
# Standard mode
lucien

# Safe mode (recommended)
lucien --safe

# Debug mode
lucien --debug

# Test mode (validation only)
lucien --test
```

#### Core Commands

**Built-in Shell Commands**:
```bash
# Directory operations
pwd                    # Print current directory
cd /path/to/directory  # Change directory

# Variables
set VAR value         # Set variable (current syntax)
export VAR            # Export to environment
echo $VAR             # Display variable value

# Aliases
alias ll='ls -la'     # Create alias
alias                 # List all aliases

# History
history               # Show command history
history 10            # Show last 10 commands

# System
exit                  # Exit shell
```

**Native Lucien Commands** (Python layer):
```bash
# AI and assistance
:claude <prompt>      # AI consultation (placeholder)
:mindmeld            # Learning system (placeholder)

# System operations
:exec <command>      # Enhanced command execution
:validate           # Run system tests
:ritual             # Memory system initialization

# Interface
:theme nexus        # Change visual theme
:theme synthwave    # Retro cyberpunk theme
:theme ghost        # Minimalist theme
```

#### Command Syntax

**Pipes and Redirection**:
```bash
# Pipes (working)
ls | grep txt

# Output redirection (working)
echo "hello" > file.txt
echo "world" >> file.txt

# Input redirection (working)
sort < input.txt

# Complex combinations (working)
cat file.txt | grep pattern > results.txt
```

**Variables and Expansion**:
```bash
# Variable setting (working)
set NAME "John Doe"
set PATH "/usr/bin:$PATH"

# Variable expansion (⚠️ PARTIAL - undefined vars broken)
echo "Hello $NAME"      # Works if NAME defined
echo "Value: $UNDEFINED" # ⚠️ Returns literal "$UNDEFINED"
```

#### Plugin Usage

```bash
# List available plugins
lucien plugin list

# Run BMAD plugin
lucien plugin bmad build
lucien plugin bmad manage
lucien plugin bmad analyze
lucien plugin bmad deploy
```

### Advanced Usage

#### Configuration Files

Create `~/.lucien/config.toml`:
```toml
[shell]
safe_mode = true
history_size = 1000

[ui]
theme = "synthwave"
glitch_effects = false

[plugins]
auto_load = ["bmad"]
plugin_dir = "~/.lucien/plugins"

[ai]
provider = "local"
model_path = "/path/to/model.gguf"
```

#### Environment Variables
```bash
export LUCIEN_DEBUG=1        # Enable debug mode
export LUCIEN_SAFE_MODE=1    # Force safe mode
export LUCIEN_CONFIG=/path   # Custom config location
```

## 💻 DEVELOPER GUIDE

### Development Environment Setup

#### Go Development
```bash
# Setup Go workspace
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Install tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Clone and setup
git clone https://github.com/luciendev/lucien-core.git
cd lucien-core
go mod tidy
```

#### Python Development
```bash
# Virtual environment
python -m venv .venv
source .venv/bin/activate

# Development dependencies
pip install -r requirements.txt
pip install -r requirements-dev.txt  # If exists

# Code formatting
pip install black isort
```

### Building from Source

#### Makefile Commands
```bash
# Build everything
make all

# Build Go binary only
make build

# Build with debugging
make build-debug

# Run tests
make test

# Clean build artifacts
make clean

# Format code
make fmt

# Lint code
make lint
```

#### Manual Build Process
```bash
# Go binary
cd cmd/lucien
go build -o ../../lucien.exe

# Python package (if needed)
python setup.py build
```

### Testing

#### Go Tests
```bash
# Unit tests
go test ./...

# Specific package tests
go test ./internal/shell/

# With coverage
go test -cover ./...

# Benchmark tests
go test -bench=. ./...
```

#### Python Tests
```bash
# Run pytest
pytest tests/

# With coverage
pytest --cov=lucien tests/
```

#### Integration Tests
```bash
# Full system test
make test-integration

# Shell functionality tests
go test ./internal/shell/ -tags=integration

# Plugin system tests
go test ./internal/plugin/ -tags=integration
```

## 📖 API REFERENCE

### Shell Engine API

#### Core Types
```go
// ExecutionResult holds command execution results
type ExecutionResult struct {
    Output   string        // Command output
    Error    string        // Error message
    ExitCode int          // Process exit code
    Duration time.Duration // Execution time
}

// Shell represents the core shell engine
type Shell struct {
    config      *Config
    env         map[string]string
    aliases     map[string]string
    history     []string
    currentDir  string
    builtins    map[string]func([]string) (*ExecutionResult, error)
}
```

#### Main Methods
```go
// Execute runs a command string through the shell
func (s *Shell) Execute(cmdLine string) (*ExecutionResult, error)

// SetVariable sets a shell variable
func (s *Shell) SetVariable(name, value string) error

// CreateAlias creates a command alias
func (s *Shell) CreateAlias(name, command string) error

// ChangeDirectory changes current working directory
func (s *Shell) ChangeDirectory(path string) error
```

### Plugin Interface

#### Plugin Contract
```go
type Plugin interface {
    Name() string
    Version() string
    Description() string
    Execute(ctx context.Context, args []string) (*PluginResult, error)
}

type PluginResult struct {
    Output   string
    Error    string
    ExitCode int
    Metadata map[string]interface{}
}
```

#### RPC Protocol
```go
// Plugin RPC interface
type PluginRPC struct {
    client rpc.Client
}

func (p *PluginRPC) Execute(args []string, reply *PluginResult) error
```

### Policy Engine API

#### Policy Types
```go
type PolicyEngine struct {
    rules []PolicyRule
}

type PolicyRule struct {
    Name        string
    Pattern     string
    Action      ActionType
    Severity    SeverityLevel
    Description string
}
```

#### Policy Methods
```go
func (p *PolicyEngine) Evaluate(command string) (*PolicyResult, error)
func (p *PolicyEngine) LoadPolicies(dir string) error
func (p *PolicyEngine) AddRule(rule PolicyRule) error
```

### Python CLI API

#### Core Classes
```python
class LucienCLI:
    """Main CLI application class"""
    
    def __init__(self, config: Config = None):
        self.config = config or Config()
        self.console = Console()
        
    def run_command(self, command: str) -> CommandResult:
        """Execute a command with safety checks"""
        
    def validate_command(self, command: str) -> ValidationResult:
        """Validate command safety"""
```

#### UI Components
```python
class LucienUI:
    """Rich-based user interface"""
    
    def spell_thinking(self, message: str):
        """Display thinking animation"""
        
    def spell_complete(self, result: Any):
        """Display completion message"""
        
    def spell_failed(self, error: str):
        """Display error message"""
```

## 🔌 PLUGIN DEVELOPMENT

### Creating a Plugin

#### Plugin Structure
```
my-plugin/
├── main.go          # Plugin implementation
├── manifest.json    # Plugin metadata
├── go.mod          # Go module file
└── README.md       # Plugin documentation
```

#### Basic Plugin Template
```go
package main

import (
    "context"
    "fmt"
    
    "github.com/hashicorp/go-plugin"
    "github.com/luciendev/lucien-core/internal/plugin"
)

type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "my-plugin"
}

func (p *MyPlugin) Version() string {
    return "1.0.0"
}

func (p *MyPlugin) Description() string {
    return "Example plugin for Lucien CLI"
}

func (p *MyPlugin) Execute(ctx context.Context, args []string) (*plugin.Result, error) {
    return &plugin.Result{
        Output:   fmt.Sprintf("Plugin executed with args: %v", args),
        ExitCode: 0,
    }, nil
}

func main() {
    plugin.Serve(&plugin.ServeConfig{
        HandshakeConfig: plugin.HandshakeConfig{
            ProtocolVersion:  1,
            MagicCookieKey:   "LUCIEN_PLUGIN",
            MagicCookieValue: "lucien_plugin_magic_cookie",
        },
        Plugins: map[string]plugin.Plugin{
            "my-plugin": &plugin.RPCPlugin{Impl: &MyPlugin{}},
        },
    })
}
```

#### Plugin Manifest
```json
{
    "name": "my-plugin",
    "version": "1.0.0",
    "description": "Example plugin for Lucien CLI",
    "author": "Your Name",
    "license": "MIT",
    "executable": "my-plugin.exe",
    "protocol": "rpc",
    "capabilities": [
        "command_execution",
        "file_access"
    ],
    "permissions": {
        "filesystem": "read-only",
        "network": "none"
    }
}
```

### Plugin Testing

#### Unit Tests
```go
func TestMyPlugin(t *testing.T) {
    plugin := &MyPlugin{}
    
    result, err := plugin.Execute(context.Background(), []string{"test"})
    assert.NoError(t, err)
    assert.Equal(t, 0, result.ExitCode)
    assert.Contains(t, result.Output, "test")
}
```

#### Integration Tests
```bash
# Test plugin loading
lucien plugin load ./my-plugin

# Test plugin execution
lucien plugin run my-plugin test-args

# Test plugin unloading
lucien plugin unload my-plugin
```

## ⚙️ CONFIGURATION

### Configuration Files

#### Main Configuration (~/.lucien/config.toml)
```toml
[shell]
safe_mode = true
history_size = 1000
auto_complete = true

[ui]
theme = "nexus"              # nexus, synthwave, ghost
glitch_effects = false
animations = true
color_scheme = "cyberpunk"

[plugins]
auto_load = ["bmad"]
plugin_timeout = 30
max_concurrent_plugins = 5

[security]
policy_engine = true
sandbox_mode = true
dangerous_commands = [
    "rm -rf",
    "format",
    "del /s"
]

[ai]
provider = "local"           # local, openai, anthropic
model_path = "/models/llama.gguf"
api_key_file = "~/.lucien/api_keys"
max_tokens = 1000

[performance]
command_timeout = 60
memory_limit = "512MB"
cpu_limit = 80
```

#### Policy Configuration (~/.lucien/policies/)
```rego
# security.rego
package lucien.security

# Block dangerous operations
deny_dangerous {
    input.command == "rm"
    input.args[_] == "-rf"
    input.args[_] == "/"
}

# Require confirmation for destructive commands
require_confirmation {
    dangerous_commands[input.command]
}

dangerous_commands := {
    "rm", "del", "format", "shutdown"
}
```

### Environment Variables

```bash
# Core settings
export LUCIEN_CONFIG=/path/to/config.toml
export LUCIEN_DEBUG=1
export LUCIEN_SAFE_MODE=1

# Plugin settings
export LUCIEN_PLUGIN_DIR=/path/to/plugins
export LUCIEN_PLUGIN_TIMEOUT=30

# AI settings
export OPENAI_API_KEY=your-key-here
export ANTHROPIC_API_KEY=your-key-here

# Security settings
export LUCIEN_POLICY_DIR=/path/to/policies
export LUCIEN_SANDBOX_MODE=1
```

## 🔍 TROUBLESHOOTING

### Common Issues

#### Build Issues
```bash
# Problem: Go build fails with dependency errors
# Solution: Update dependencies
go mod tidy
go mod download

# Problem: Python import errors
# Solution: Install requirements
pip install -r requirements.txt

# Problem: Make commands fail
# Solution: Check Make installation
make --version
```

#### Runtime Issues
```bash
# Problem: Plugin loading fails
# Solution: Check plugin manifest and permissions
lucien plugin validate ./my-plugin

# Problem: Commands hang or timeout
# Solution: Enable debug mode
LUCIEN_DEBUG=1 lucien --debug

# Problem: Variable expansion not working
# Solution: Known bug - use workaround
set VAR value  # Instead of VAR=value
```

#### Performance Issues
```bash
# Problem: Slow startup
# Solution: Check plugin loading
lucien --no-plugins

# Problem: High memory usage
# Solution: Limit plugins or reduce history
lucien --config minimal.toml

# Problem: Command execution slow
# Solution: Check sandbox settings
lucien --no-sandbox
```

### Debugging

#### Enable Debug Mode
```bash
# Environment variable
export LUCIEN_DEBUG=1

# Command line flag
lucien --debug

# Configuration file
[debug]
enabled = true
log_level = "trace"
```

#### Log Files
```bash
# Default log location
~/.lucien/logs/lucien.log

# Debug log
~/.lucien/logs/debug.log

# Plugin logs
~/.lucien/logs/plugins/
```

#### Debug Commands
```bash
# System information
:validate

# Plugin status
lucien plugin status

# Configuration dump
lucien config dump

# Health check
lucien health
```

---

## 🚧 KNOWN LIMITATIONS

### Critical Issues (Must Fix)
1. **Variable Expansion**: Undefined variables return literal text
2. **Duration Tracking**: Missing execution time measurement
3. **Command Syntax**: Standard `VAR=value` not supported
4. **Empty Commands**: Error instead of graceful handling

### Feature Limitations
1. **AI Integration**: Placeholder implementations only
2. **Script Execution**: No shell script file support
3. **Cross-Platform**: Windows-focused implementation
4. **Policy Engine**: Compilation errors prevent testing

### Performance Considerations
1. **Plugin Loading**: ~100ms per plugin
2. **Memory Usage**: Scales with history size
3. **Startup Time**: Affected by plugin count
4. **Command Response**: Generally <100ms

---

*This manual is a living document that reflects the current state of the Lucien CLI project. As issues are resolved and features are implemented, this documentation will be updated accordingly.*

**For the most current information, refer to the project repository and issue tracker.**