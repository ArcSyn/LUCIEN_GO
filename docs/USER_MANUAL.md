# Lucien CLI User Manual

## Table of Contents

1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [Core Commands](#core-commands)
4. [Operator Chaining](#operator-chaining)
5. [Security Modes](#security-modes)
6. [Variables and Environment](#variables-and-environment)
7. [Aliases](#aliases)
8. [History Management](#history-management)
9. [Job Control](#job-control)
10. [Batch Processing](#batch-processing)
11. [Configuration](#configuration)
12. [Troubleshooting](#troubleshooting)

## Introduction

Lucien CLI is a modern command-line shell designed with security and usability in mind. It provides advanced parsing capabilities, intelligent operator chaining, and built-in security controls while maintaining familiar shell semantics.

### Key Benefits

- **Security First**: Built-in command injection protection and configurable security modes
- **Intelligent Parsing**: Quote-aware parsing that handles complex command structures correctly
- **Cross-Platform**: Works consistently across Windows, macOS, and Linux
- **Production Ready**: Robust error handling and comprehensive testing

## Getting Started

### Installation

Choose your installation method:

**Build from Source:**
```bash
git clone https://github.com/ArcSyn/LUCIEN_GO.git
cd LUCIEN_GO
go build -o lucien ./cmd/lucien
```

**Download Binary:**
Download the appropriate binary for your platform from the [releases page](https://github.com/ArcSyn/LUCIEN_GO/releases).

### First Launch

#### Interactive Mode
```bash
./lucien
```

This launches the full interactive shell with the cyberpunk-styled TUI interface.

#### Batch Mode
```bash
echo "pwd" | ./lucien --batch
```

This runs commands non-interactively, perfect for scripts and automation.

## Core Commands

### Built-in Commands

| Command | Syntax | Description | Example |
|---------|--------|-------------|---------|
| `pwd` | `pwd` | Show current working directory | `pwd` |
| `home` | `home` | Navigate to home directory | `home` |
| `clear` | `clear` | Clear the screen | `clear` |
| `history` | `history` | Show command history | `history` |
| `jobs` | `jobs` | List active background jobs | `jobs` |
| `env` | `env` | Display environment variables | `env` |
| `export` | `export VAR=value` | Set environment variable | `export PATH=/usr/bin` |
| `alias` | `alias name='command'` | Create command alias | `alias ll='ls -la'` |
| `:secure` | `:secure [strict\|permissive]` | Control security mode | `:secure strict` |

### Navigation Commands

```bash
# Show current directory
lucien> pwd
/home/user

# Go to home directory
lucien> home
/home/user

# Change directory (if cd is available)
lucien> cd /tmp
lucien> pwd
/tmp
```

### System Information

```bash
# List environment variables
lucien> env
PATH=/usr/local/bin:/usr/bin:/bin
HOME=/home/user
SHELL=/bin/bash
...

# Show command history
lucien> history
   1  pwd
   2  echo hello
   3  home
   4  history
```

## Operator Chaining

Lucien supports sophisticated command chaining with proper precedence and short-circuit evaluation.

### Supported Operators

| Operator | Name | Description | Behavior |
|----------|------|-------------|----------|
| `&&` | AND | Execute next command only if previous succeeded | Short-circuits on failure |
| `\|\|` | OR | Execute next command only if previous failed | Short-circuits on success |
| `;` | Sequence | Execute commands in sequence regardless of status | Always continues |
| `\|` | Pipe | Pass output of first command to second | Stream processing |
| `&` | Background | Run command in background | Non-blocking execution |

### Chaining Examples

#### AND Operator (&&)
```bash
# Both commands run if first succeeds
lucien> echo "step1" && echo "step2"
step1
step2

# Second command doesn't run if first fails
lucien> false && echo "won't see this"
# No output from echo
```

#### OR Operator (||)
```bash
# Second command runs only if first fails
lucien> echo "success" || echo "backup"
success

# Backup command runs if first fails
lucien> false || echo "backup executed"
backup executed
```

#### Sequence Operator (;)
```bash
# Both commands run regardless of status
lucien> echo "first" ; echo "second"
first
second

# Both run even if first fails
lucien> false ; echo "still runs"
still runs
```

#### Complex Chaining
```bash
# Combine multiple operators
lucien> echo "start" && echo "middle" || echo "fallback" ; echo "end"
start
middle
end
```

### Quote Handling

Operators inside quotes are treated as literal text:

```bash
# Operators inside quotes are literal
lucien> echo "command && operator"
command && operator

# Mixed quotes and operators
lucien> echo "literal &&" && echo "actual operator"
literal &&
actual operator
```

## Security Modes

Lucien provides two security modes to balance usability and safety.

### Permissive Mode (Default)

```bash
lucien> :secure
Security mode: permissive
```

In permissive mode:
- All standard shell operations are allowed
- Operators work normally with any command
- Suitable for interactive development use

### Strict Mode

```bash
lucien> :secure strict
Security mode set to strict
```

In strict mode:
- Only whitelisted commands can use operators
- Dangerous command patterns are blocked
- Enhanced protection against command injection

### Whitelisted Commands

These commands are considered safe and can use operators in strict mode:
- `echo` - Text output
- `pwd` - Directory display
- `ls` - File listing
- `cd` - Directory change
- `clear` - Screen clear
- `home` - Home navigation
- `history` - Command history
- `jobs` - Job listing
- `env` - Environment display
- `export` - Variable setting
- `alias` - Alias creation

### Security Examples

```bash
# Safe in both modes
lucien> echo "hello" && pwd
hello
/home/user

# May be restricted in strict mode depending on implementation
lucien> rm file.txt && echo "deleted"
# Could be blocked in strict mode for safety
```

## Variables and Environment

### Environment Variables

#### Viewing Variables
```bash
# Show all environment variables
lucien> env
PATH=/usr/local/bin:/usr/bin
HOME=/home/user
USER=username

# Show specific variable (if supported)
lucien> echo $HOME
/home/user
```

#### Setting Variables
```bash
# Set environment variable
lucien> export MYVAR=hello
lucien> echo $MYVAR
hello

# Set with spaces
lucien> export MESSAGE="hello world"
lucien> echo $MESSAGE
hello world
```

### Variable Expansion

Lucien supports multiple variable expansion formats:

#### Unix Style
```bash
lucien> export NAME=Lucien
lucien> echo $NAME
Lucien
lucien> echo ${NAME}_CLI
Lucien_CLI
```

#### Windows Style
```bash
lucien> echo %HOME%
/home/user
```

### Tilde Expansion

```bash
# Expand to home directory
lucien> echo ~
/home/user

# Expand in paths
lucien> echo ~/documents
/home/user/documents
```

## Aliases

### Creating Aliases

```bash
# Simple alias
lucien> alias ll='ls -la'

# Alias with multiple commands
lucien> alias status='pwd && ls'

# Alias with quotes
lucien> alias greet='echo "Hello World"'
```

### Using Aliases

```bash
# Use the alias
lucien> ll
# Executes: ls -la

lucien> status
/home/user
file1.txt file2.txt

lucien> greet
Hello World
```

### Managing Aliases

```bash
# List aliases (if supported)
lucien> alias
ll='ls -la'
status='pwd && ls'
greet='echo "Hello World"'

# Remove alias (if unalias command exists)
lucien> unalias ll
```

## History Management

### Viewing History

```bash
lucien> history
   1  pwd
   2  echo hello
   3  home
   4  export MYVAR=test
   5  history
```

### History Features

- Commands are automatically saved to `~/.lucien/history`
- History persists between sessions
- Up/down arrows navigate history in interactive mode
- History expansion may be available (e.g., `!!`, `!n`)

## Job Control

### Background Jobs

```bash
# Run command in background
lucien> sleep 30 &
[1] 12345

# List active jobs
lucien> jobs
[1] Running    sleep 30 &
```

### Job Management

```bash
# List jobs
lucien> jobs
No active jobs

# With active jobs
lucien> jobs
[1] Running    sleep 30 &
[2] Running    ping example.com &
```

## Batch Processing

### Command Line Usage

```bash
# Single command
echo "pwd" | lucien --batch

# Multiple commands
echo -e "pwd\necho hello\nhome" | lucien --batch

# From file
cat commands.txt | lucien --batch
```

### Script Integration

```bash
#!/bin/bash
# Generate commands and pipe to Lucien
{
    echo "echo 'Starting batch process'"
    echo "pwd"
    echo "echo 'Process complete'"
} | lucien --batch
```

### Batch Mode Features

- Non-interactive execution
- Perfect for automation
- Preserves all security and parsing features
- Exit codes properly propagated

## Configuration

### Configuration Files

Lucien looks for configuration in:
- **Linux/macOS**: `~/.lucien/`
- **Windows**: `%USERPROFILE%\.lucien\`

### File Locations

| File | Purpose | Status |
|------|---------|--------|
| `history` | Command history | Active |
| `config.toml` | Main configuration | Future |
| `aliases` | User aliases | Future |

### Current Configuration

Currently, configuration is primarily through command-line flags:

```bash
# Enable safe mode
lucien --safe-mode

# Custom SSH port
lucien --ssh --port 2222

# Custom config file
lucien --config ~/.lucien/custom.toml
```

## Troubleshooting

### Common Issues

#### "Command not found" Errors

```bash
lucien> invalidcommand
Error: command not found: invalidcommand
```

**Solution**: Check spelling and ensure the command exists in your PATH.

#### Operator Not Working as Expected

```bash
# If operators aren't chaining properly
lucien> echo "first" && echo "second"
first
# Missing "second"
```

**Solutions**:
1. Check if you're in strict security mode: `:secure`
2. Verify quote handling: ensure operators aren't inside quotes
3. Check command exit status

#### Batch Mode Not Working

```bash
# If piped input goes to TUI instead of batch
echo "pwd" | lucien
# Shows TUI instead of output
```

**Solution**: Use the `--batch` flag explicitly:
```bash
echo "pwd" | lucien --batch
```

#### Permission Denied

```bash
lucien> some-restricted-command
Error: permission denied or command blocked
```

**Solutions**:
1. Check security mode: `:secure`
2. Switch to permissive if needed: `:secure permissive`
3. Verify file permissions

#### History Not Persisting

**Check**:
1. Directory permissions on `~/.lucien/`
2. Disk space availability
3. File system support for the history file

### Debug Information

Enable verbose output:
```bash
# Set debug environment variable
LUCIEN_DEBUG=1 lucien

# Or use debug flag if available
lucien --debug
```

### Getting Help

1. **Check Command Reference**: Review this manual
2. **Built-in Help**: Try `:help` in interactive mode
3. **GitHub Issues**: Report bugs at https://github.com/ArcSyn/LUCIEN_GO/issues
4. **Discussions**: Ask questions at https://github.com/ArcSyn/LUCIEN_GO/discussions

### Performance Tips

1. **Use Batch Mode for Scripts**: Non-interactive mode is faster
2. **Limit History Size**: Large history files can slow startup
3. **Disable Unused Features**: Use minimal configuration for better performance

---

## Quick Reference

### Essential Commands
```bash
pwd              # Current directory
home             # Go home
clear            # Clear screen
history          # Command history
jobs             # Background jobs
env              # Environment variables
:secure          # Check security mode
:secure strict   # Enable strict mode
```

### Operators
```bash
cmd1 && cmd2     # Run cmd2 if cmd1 succeeds
cmd1 || cmd2     # Run cmd2 if cmd1 fails
cmd1 ; cmd2      # Run both regardless
cmd1 | cmd2      # Pipe cmd1 output to cmd2
cmd1 &           # Run cmd1 in background
```

### Variables
```bash
export VAR=value # Set variable
echo $VAR        # Use variable
echo ~           # Home directory
```

### Aliases
```bash
alias name='cmd' # Create alias
name             # Use alias
```

This manual covers all currently implemented and tested features of Lucien CLI v1.0.0.