# 🧠 Lucien CLI

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/ArcSyn/LUCIEN_GO/workflows/CI/badge.svg)](https://github.com/ArcSyn/LUCIEN_GO/actions)

> A next-generation command-line shell with built-in security, advanced parsing, and intelligent command chaining.

Lucien CLI transforms your terminal experience with production-grade security controls, sophisticated operator handling, and cross-platform compatibility. Built for developers who need both power and protection.

## ✨ Features

### 🔧 Core Implementation
- **Security Guard System**: Post-parsing validation with strict/permissive modes
- **Full Operator Support**: `&&`, `||`, `;`, `|`, `&` with precedence and short-circuiting
- **Advanced Parser**: Quote-aware parsing; treats operators inside quotes as literals
- **Command Injection Protection**: Whitelist-based validation + dangerous pattern detection
- **Variable Expansion**: `$VAR`, `${VAR}`, `%VAR%` cross-platform
- **Tilde Expansion**: `~` and `~/` for home directory shortcuts

### 🛡️ Security Features
- `:secure strict` → blocks risky chained commands unless whitelisted
- `:secure permissive` → normal shell behavior (default)
- Whitelisted builtins: echo, pwd, cd, ls, clear, home, etc.

### 🔨 Enhanced Builtins
- `home` → go to platform-specific home dir
- `export` → set environment variables
- `env` → list environment variables
- `clear` → ANSI screen clear
- Aliases, history, job control

### 🚀 Production Features
- `--batch` flag for non-interactive execution
- History persistence at `~/.lucien/history`
- Cross-platform: Windows, macOS, Linux
- Graceful error handling for invalid commands

## 🚀 Quick Start

### Installation

#### Windows (PowerShell)
```powershell
# Download latest release
Invoke-WebRequest -Uri "https://github.com/ArcSyn/LUCIEN_GO/releases/latest/download/lucien-windows-amd64.exe" -OutFile "lucien.exe"
# Move to PATH
Move-Item lucien.exe "$env:USERPROFILE\bin\lucien.exe"
```

#### macOS
```bash
# Using Homebrew (coming soon)
# brew install arcsyn/tap/lucien

# Or download directly
curl -L "https://github.com/ArcSyn/LUCIEN_GO/releases/latest/download/lucien-darwin-amd64" -o lucien
chmod +x lucien
sudo mv lucien /usr/local/bin/
```

#### Linux
```bash
curl -L "https://github.com/ArcSyn/LUCIEN_GO/releases/latest/download/lucien-linux-amd64" -o lucien
chmod +x lucien
sudo mv lucien /usr/local/bin/
```

#### Build from Source
```bash
git clone https://github.com/ArcSyn/LUCIEN_GO.git
cd LUCIEN_GO
go build -o lucien ./cmd/lucien
```

### First Run

Start Lucien in interactive mode:
```bash
lucien
```

Or run commands in batch mode:
```bash
echo "pwd && echo 'Hello Lucien'" | lucien --batch
```

## 🎯 Usage Examples

### Basic Commands
```bash
lucien> pwd
/home/user

lucien> echo "Welcome to Lucien"
Welcome to Lucien

lucien> home
/home/user
```

### Operator Chaining
```bash
# Sequential execution
lucien> echo "step1" && echo "step2"
step1
step2

# Conditional execution
lucien> echo "success" || echo "backup"
success

# Command sequence
lucien> echo "first" ; echo "second"
first
second
```

### Quoted Operators
```bash
lucien> echo "operators && inside quotes"
operators && inside quotes
```

### Security Modes
```bash
# Check current security mode
lucien> :secure
Security mode: permissive

# Switch to strict mode
lucien> :secure strict
Security mode set to strict
```

### Variables and Aliases
```bash
# Set and use variables
lucien> set TESTVAR=hello
lucien> echo $TESTVAR
hello

# Create and use aliases
lucien> alias ll='echo long listing'
lucien> ll
long listing
```

### Batch Processing
```bash
# Process commands from file
cat commands.txt | lucien --batch

# Single command execution
echo "pwd" | lucien --batch
```

## 📚 Key Commands

| Command | Description |
|---------|-------------|
| `pwd` | Show current directory |
| `home` | Navigate to home directory |
| `clear` | Clear screen |
| `history` | Show command history |
| `jobs` | Show background jobs |
| `env` | List environment variables |
| `export VAR=value` | Set environment variable |
| `alias name='command'` | Create command alias |
| `:secure [strict\|permissive]` | Toggle security mode |

## ⚙️ Configuration

Lucien stores configuration and data in:

- **Linux/macOS**: `~/.lucien/`
- **Windows**: `%USERPROFILE%\.lucien\`

Files:
- `history` - Command history
- `config.toml` - Configuration settings (future)
- `aliases` - User-defined aliases (future)

## 🧪 Testing

Run the test suite:
```bash
go test ./...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Run integration tests:
```bash
go test ./tests/...
```

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

Quick start for contributors:
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./...`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## 📖 Documentation

- [User Manual](docs/USER_MANUAL.md) - Complete command reference
- [Contributing Guide](CONTRIBUTING.md) - How to contribute
- [Code of Conduct](CODE_OF_CONDUCT.md) - Community guidelines
- [Changelog](CHANGELOG.md) - Version history

## 🛡️ Security

Lucien CLI takes security seriously. If you discover a security vulnerability, please send an email to security@arcsyn.dev rather than opening a public issue.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🌟 Community

- 🐛 [Report bugs](https://github.com/ArcSyn/LUCIEN_GO/issues/new?template=bug_report.md)
- 💡 [Request features](https://github.com/ArcSyn/LUCIEN_GO/issues/new?template=feature_request.md)
- 💬 [Discussions](https://github.com/ArcSyn/LUCIEN_GO/discussions)
- 📧 [Email](mailto:hello@arcsyn.dev)

---

*Built with ⚡ by the ArcSyn team*