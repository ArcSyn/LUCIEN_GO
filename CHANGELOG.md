# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-08-08

### Added

#### 🔧 Core Implementation
- Security Guard System: Post-parsing validation with strict/permissive modes
- Full Operator Support: &&, ||, ;, |, & with precedence and short-circuiting
- Advanced Parser: Quote-aware parsing; treats operators inside quotes as literals
- Command Injection Protection: Whitelist-based validation + dangerous pattern detection
- Variable Expansion: $VAR, ${VAR}, %VAR% cross-platform
- Tilde Expansion: ~ and ~/ for home directory shortcuts

#### 🛡️ Security Features
- `:secure strict` → blocks risky chained commands unless whitelisted
- `:secure permissive` → normal shell behavior (default)
- Whitelisted builtins: echo, pwd, cd, ls, clear, home, etc.

#### 🔨 Enhanced Builtins
- `home` → go to platform-specific home dir
- `export` → set environment variables
- `env` → list environment variables
- `clear` → ANSI screen clear
- Aliases, history, job control

#### 🚀 Production Features
- `--batch` flag for non-interactive execution
- History persistence at `~/.lucien/history`
- Cross-platform: Windows, macOS, Linux
- Graceful error handling for invalid commands

#### ✅ Tested Commands & Examples
- `pwd`
- `echo test1 && echo test2`
- `echo success || echo backup`
- `echo 'operators && inside quotes'`
- `:secure`
- `:secure strict`
- `home`
- `set TESTVAR=hello` + `echo $TESTVAR`
- `alias ll='echo long listing'` + `ll`
- `history`
- `jobs`
- `env`

### Technical Details
- Built with Go 1.21+
- Cross-platform compatibility (Windows, macOS, Linux)
- Production-grade error handling and validation
- Comprehensive test coverage

---

## [Unreleased]

### Planned
- Plugin marketplace integration
- Enhanced configuration system
- Additional theme support
- Extended AI integration capabilities

---

**Note**: This is the first public release of Lucien CLI. All features listed above have been thoroughly tested and are production-ready.