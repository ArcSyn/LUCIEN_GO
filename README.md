# 🧠 LUCIEN SHELL - NEXUS-7 NEURAL INTERFACE

```
██╗     ██╗   ██╗ ██████╗██╗███████╗███╗   ██╗    ███████╗██╗  ██╗███████╗██╗     ██╗     
██║     ██║   ██║██╔════╝██║██╔════╝████╗  ██║    ██╔════╝██║  ██║██╔════╝██║     ██║     
██║     ██║   ██║██║     ██║█████╗  ██╔██╗ ██║    ███████╗███████║█████╗  ██║     ██║     
██║     ██║   ██║██║     ██║██╔══╝  ██║╚██╗██║    ╚════██║██╔══██║██╔══╝  ██║     ██║     
███████╗╚██████╔╝╚██████╗██║███████╗██║ ╚████║    ███████║██║  ██║███████╗███████╗███████╗
╚══════╝ ╚═════╝  ╚═════╝╚═╝╚══════╝╚═╝  ╚═══╝    ╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝
```

**The world's first AI-enhanced shell with neural learning capabilities**

---

## 🎯 **WHAT IS LUCIEN?**

Lucien is a revolutionary command-line interface that combines the power of traditional shells with cutting-edge AI technology. Built with Go and the Bubble Tea TUI framework, it provides a cyberpunk-styled terminal experience with advanced security features and plugin architecture.

### 🌟 **Key Features**

- **🧠 Neural AI Integration** - Built-in AI assistance with local and cloud options
- **🛡️ Advanced Security** - OPA policy enforcement and sandboxed execution  
- **🔌 Plugin Architecture** - Extensible system with BMAD methodology plugins
- **🎨 Cyberpunk Aesthetic** - Matrix-inspired themes with glitch effects
- **⚡ High Performance** - Compiled Go binary with minimal resource usage
- **🌐 Remote Access** - SSH server with full TUI over network
- **🔒 Safe Mode** - Production-ready safety features and command filtering

---

## 🚀 **QUICK START**

### Prerequisites
- Go 1.22 or higher
- Git (for updates)
- SSH client (for remote access)

### Installation

```bash
# Clone the repository
git clone https://github.com/luciendev/lucien-core.git
cd lucien-core

# Build everything
make all

# Run Lucien locally
make run
```

### Alternative: One-Command Setup
```bash
# Build and run immediately
make build && ./build/lucien
```

---

## 🎮 **USAGE**

### Local Interface
```bash
# Start standard interface
lucien

# Start with security protocols
lucien --safe-mode

# Start with configuration file
lucien --config ~/.lucien/config.toml
```

### SSH Remote Access
```bash
# Start SSH server
lucien --ssh --port 2222

# Connect from another machine
ssh localhost -p 2222
```

### Interactive Commands
```bash
# AI assistance
:ai <query>

# Change visual theme
:theme nexus
:theme synthwave  
:theme ghost

# Activate glitch mode
:hack

# Clear terminal
:clear

# Show help
:help
```

---

## 🏗️ **ARCHITECTURE**

### Core Systems
```
lucien-core/
├── cmd/lucien/          # Main application entry point
├── internal/            # Core engine components
│   ├── ui/             # Bubble Tea TUI with cyberpunk themes
│   ├── shell/          # Advanced shell parser and executor
│   ├── plugin/         # Plugin management and RPC
│   ├── policy/         # OPA security policy engine
│   ├── sandbox/        # Process sandboxing (gVisor/Job Objects)
│   └── ai/             # AI integration (local/cloud)
└── plugins/            # Extensible plugin ecosystem
    ├── example-bmad/   # Build-Manage-Analyze-Deploy workflow
    └── example-weather/# Network-enabled weather plugin
```

### Technology Stack
- **UI Framework**: Bubble Tea + Lip Gloss + Bubbles (Charm Bracelet)
- **Remote Access**: Wish SSH middleware
- **Plugin Runtime**: HashiCorp go-plugin with RPC isolation
- **Security**: Open Policy Agent with Rego rules
- **Sandboxing**: gVisor (Linux) + Windows Job Objects
- **AI Engine**: llama.cpp integration + cloud API support

---

## 🔌 **PLUGIN SYSTEM**

### BMAD Plugin Example
```bash
# Execute BMAD workflow phases
lucien plugin bmad build     # Build phase operations
lucien plugin bmad manage    # System management
lucien plugin bmad analyze   # Code and security analysis  
lucien plugin bmad deploy    # Production deployment
lucien plugin bmad workflow  # Complete lifecycle
```

### Creating Custom Plugins
```go
type MyPlugin struct{}

func (p *MyPlugin) Execute(ctx context.Context, command string, args []string) (*Result, error) {
    return &Result{
        Output:   "Custom plugin response with cyberpunk flair",
        ExitCode: 0,
    }, nil
}
```

---

## 🛡️ **SECURITY FEATURES**

### OPA Policy Engine
```rego
# Example security policy
package lucien.security

# Block dangerous root operations
deny_write_root {
    input.action == "execute"
    input.command == "rm"
    contains(input.args[_], "/")
}
```

### Safe Mode Operations
- **Command Filtering** - Blocks dangerous operations automatically
- **Plugin Sandboxing** - Restricts plugin file system access
- **Resource Limits** - CPU/memory/time constraints
- **Audit Logging** - Comprehensive security event logging

---

## 🧠 **AI INTEGRATION**

### Local AI (llama.cpp)
```bash
export LLAMACPP_MODEL="/path/to/model.gguf"
lucien --config ai_local.toml
```

### Cloud AI APIs
```bash
export OPENAI_API_KEY="your-api-key"
export LUCIEN_AI_PROVIDER="openai"
lucien
```

### AI Commands
```bash
# Query the AI engine
:ai How do I optimize this Docker build?

# AI-powered command suggestions  
:ai Suggest commands for system monitoring

# Code analysis
:ai --context myfile.go Explain this function
```

---

## 🎨 **THEMES & CUSTOMIZATION**

### Available Themes
- **nexus** - Classic Matrix green-on-black
- **synthwave** - Retro cyberpunk pink/cyan
- **ghost** - Minimalist white-on-black

### Theme Switching
```bash
:theme nexus      # Switch to Matrix theme
:theme synthwave  # Switch to synthwave theme
:hack            # Toggle glitch effects
```

### Custom Configurations
```toml
# ~/.lucien/config.toml
[ui]
theme = "synthwave"
glitch_effects = true

[ai] 
provider = "local"
model_path = "/models/codellama-7b.gguf"

[security]
safe_mode = true
policy_dir = "~/.lucien/policies"
```

---

## 🧪 **TESTING**

### Run Test Suite
```bash
# Comprehensive testing
make test

# Performance benchmarks
make test-bench

# Integration tests
make test-integration
```

### Test Coverage
```bash
# Generate coverage report
make test
open coverage.html
```

---

## 📦 **BUILDING & DEPLOYMENT**

### Local Development
```bash
# Development environment
make dev

# Quick build and test
make quick

# Format and lint
make fmt lint
```

### Multi-Platform Builds
```bash
# Build for all platforms
make build-all

# Creates:
# - lucien-linux-amd64
# - lucien-darwin-amd64  
# - lucien-windows-amd64.exe
```

### Release Package
```bash
# Create complete release
make release

# Generates: lucien-1.0.0-nexus7-20240101_120000.tar.gz
```

---

## 🌐 **REMOTE ACCESS**

### SSH Server Setup
```bash
# Start SSH server
lucien --ssh --port 2222

# Generate host key automatically
# Server runs full TUI over SSH
```

### Remote Connection
```bash
# Connect to Lucien SSH server
ssh user@hostname -p 2222

# Full terminal interface available remotely
# All plugins and AI features accessible
```

---

## ⚙️ **INSTALLATION**

### System Installation
```bash
# Install to system PATH
make install

# Now available system-wide
lucien --help
```

### Package Managers
```bash
# Homebrew (macOS)
brew tap luciendev/tap
brew install lucien

# Chocolatey (Windows)
choco install lucien

# Snap (Linux)
snap install lucien
```

---

## 🤝 **CONTRIBUTING**

### Development Setup
```bash
git clone https://github.com/luciendev/lucien-core.git
cd lucien-core
make deps dev
```

### Plugin Development
```bash
# Create plugin template
mkdir plugins/my-plugin
cd plugins/my-plugin
# Implement PluginInterface
# Add manifest.json
```

### Testing & Quality
- Write comprehensive tests for all features
- Follow Go best practices and idioms
- Maintain cyberpunk aesthetic consistency
- Ensure cross-platform compatibility

---

## 📊 **PERFORMANCE**

### Benchmarks
- **Startup Time**: <100ms cold start
- **Memory Usage**: ~50MB base + 10MB per plugin
- **Command Execution**: <50ms average response
- **AI Query**: <200ms local, varies for cloud

### Resource Efficiency
- Compiled Go binary with minimal dependencies
- Efficient TUI rendering with Bubble Tea
- Smart caching for plugin operations
- Optional resource limits via sandbox

---

## 🆘 **TROUBLESHOOTING**

### Common Issues

**Build Failures**
```bash
# Update dependencies
make deps-update

# Clean rebuild
make clean all
```

**Plugin Errors**
```bash
# Check plugin manifest
cat plugins/my-plugin/manifest.json

# Rebuild plugins
make plugins
```

**SSH Connection Issues**
```bash
# Check SSH server status
lucien --ssh --port 2222 --debug

# Verify host key
ls -la .ssh/lucien_host_key*
```

### Debug Mode
```bash
# Enable verbose logging
LUCIEN_DEBUG=1 lucien --debug

# Check policy loading
lucien --safe-mode --debug
```

---

## 📚 **DOCUMENTATION**

### API Reference
```bash
# Generate documentation
make docs
open docs/api.txt
```

### Examples
- [Basic Usage](examples/basic.md)
- [Plugin Development](examples/plugins.md)  
- [Security Configuration](examples/security.md)
- [AI Integration](examples/ai.md)

---

## 📄 **LICENSE**

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🎯 **ROADMAP**

### Version 1.1
- [ ] Advanced AI model fine-tuning
- [ ] Plugin marketplace integration
- [ ] Enhanced security policies
- [ ] Mobile companion app

### Version 1.2  
- [ ] Distributed plugin execution
- [ ] Machine learning command optimization
- [ ] Advanced visualization dashboards
- [ ] Enterprise security features

---

## 🌟 **ACKNOWLEDGMENTS**

- **Charm Bracelet** - Bubble Tea TUI framework
- **HashiCorp** - go-plugin architecture
- **Open Policy Agent** - Security policy engine
- **gVisor** - Application sandbox technology

---

## 📞 **SUPPORT**

- **Issues**: [GitHub Issues](https://github.com/luciendev/lucien-core/issues)
- **Discussions**: [GitHub Discussions](https://github.com/luciendev/lucien-core/discussions)
- **Security**: security@luciendev.com
- **General**: hello@luciendev.com

---

**LUCIEN SHELL - WHERE INTELLIGENCE MEETS COMMAND LINE**

*The neural interface that learns, adapts, and evolves with you.*

```
▓▒░ NEURAL PATHWAYS OPTIMIZED ░▒▓
    SYSTEMS NOMINAL - READY FOR DEPLOYMENT
```