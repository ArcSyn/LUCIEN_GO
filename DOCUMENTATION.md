# 🧠 Lucien CLI - Complete Documentation

## Table of Contents
1. [Overview](#overview)
2. [Revolutionary Features](#revolutionary-features)
3. [Installation & Setup](#installation--setup)
4. [Core Commands](#core-commands)
5. [Shell Replacement Features](#shell-replacement-features)
6. [AI & Intelligence System](#ai--intelligence-system)
7. [BMAD Method Agents](#bmad-method-agents)
8. [Configuration System](#configuration-system)
9. [PowerShell 7 Parity](#powershell-7-parity)
10. [Advanced Usage](#advanced-usage)
11. [Development & Extension](#development--extension)

---

## Overview

**Lucien CLI** is a next-generation, AI-enhanced shell replacement that combines the power of traditional command-line interfaces with advanced artificial intelligence capabilities. Built on the **BMAD METHOD** (Build, Manage, Analyze, Deploy), Lucien provides intelligent automation, predictive command execution, and real-time coding assistance.

### Key Capabilities
- **Full Shell Replacement** - Complete PowerShell/Bash functionality with pipes, redirects, variables
- **AI-Powered Intelligence** - Predictive command execution and contextual suggestions
- **Real-time Copilot** - AI pair programming with live code analysis
- **Adaptive Workflows** - Automated multi-step project orchestration
- **BMAD Agents** - Intelligent build, manage, analyze, and deploy agents
- **Persistent Memory** - Context-aware assistance that learns your patterns

---

## Revolutionary Features

### 1. 🧠 Intelligent Command Prediction & Auto-Execution
Learns your command patterns and predicts what you need before you finish typing.

```bash
# Lucien learns that after "git status" you usually run "git add ."
lucien> git status
# Suggests: git add . (confidence: 89%)

# Auto-executes safe commands when confidence is very high
lucien> gi[TAB] → git status (auto-executed)
```

**Commands:**
- `lucien ai predict [partial_input]` - Get command predictions
- `lucien ai suggestions` - Get context-aware suggestions
- `lucien ai completions <partial>` - Smart tab completions

### 2. 🤖 Real-time AI Pair Programming Copilot
Monitors your code in real-time and provides intelligent suggestions.

```bash
# Start monitoring your project
lucien ai start-copilot

# Get suggestions for current file
lucien ai copilot-suggestions main.py --line 42

# Analyze errors with AI
lucien ai analyze-error "ModuleNotFoundError: No module named 'requests'"
```

**Features:**
- **Live Code Analysis** - Detects issues as you code
- **Import Suggestions** - Auto-suggests missing imports
- **Error Resolution** - Provides specific fixes for errors
- **Best Practice Enforcement** - Suggests code improvements

### 3. 🚀 Adaptive Workflow Orchestration & Project Consciousness
Understands your entire project and orchestrates complex workflows automatically.

```bash
# Analyze your project
lucien ai analyze-project

# Get comprehensive insights
lucien ai project-insights

# Create adaptive workflows
lucien ai create-workflow "setup python project"
lucien ai execute-workflow workflow_123456789

# List active workflows
lucien ai list-workflows
```

**Capabilities:**
- **Deep Project Analysis** - Understands project type, dependencies, build systems
- **Intelligent Workflow Generation** - Creates optimal step sequences
- **Performance Learning** - Optimizes workflows based on execution history
- **Risk Assessment** - Identifies potential issues and suggests mitigations

---

## Installation & Setup

### Prerequisites
- Python 3.11+
- Git (for version control features)
- Windows, macOS, or Linux

### Installation
```bash
# Clone the repository
git clone https://github.com/your-username/lucien-cli.git
cd lucien-cli

# Install dependencies
pip install -r requirements.txt

# Test installation
python -m cli.main --help
```

### Setup as Default Shell
```bash
# Option 1: Windows Terminal integration
# Add to Windows Terminal settings.json:
{
  "commandline": "python C:/path/to/lucien-cli/cli/main.py interactive"
}

# Option 2: Direct shell mode
python -m cli.main interactive

# Option 3: Create system alias
# Add to ~/.bashrc or ~/.zshrc:
alias lucien="python /path/to/lucien-cli/cli/main.py"
```

---

## Core Commands

### Command Structure
```bash
lucien <command_group> <command> [options] [arguments]
```

### Main Command Groups
- **`shell`** - Shell operations and execution
- **`cast`** - Agent casting and AI operations
- **`ai`** - Advanced AI features and intelligence
- **`ps`** - PowerShell compatibility
- **`validate`** - System testing and validation
- **`summon`** - Agent creation and management
- **`daemon`** - Background processes

### Essential Commands
```bash
# Execute shell commands
lucien shell run "ls -la"
lucien exec "git status"  # Legacy support

# Interactive shell mode
lucien interactive

# Change directory with persistence
lucien shell cd ~/projects

# Manage environment variables
lucien shell env PATH "/new/path"
lucien shell alias ll "ls -la"

# Execute scripts
lucien shell script my_script.sh
```

---

## Shell Replacement Features

### Full Shell Functionality
Lucien provides complete shell replacement with all standard features:

#### Command Execution
```bash
# Basic commands
lucien> ls -la
lucien> cd ~/projects
lucien> pwd

# Complex commands with arguments
lucien> find . -name "*.py" -type f
lucien> grep -r "TODO" src/
```

#### Pipes and Redirects
```bash
# Pipes
lucien> ls | grep ".py"
lucien> ps aux | grep python

# Output redirection
lucien> echo "Hello World" > output.txt
lucien> ls -la >> log.txt

# Input redirection
lucien> sort < unsorted.txt
```

#### Variables and Environment
```bash
# Set variables
lucien> export PATH="/usr/local/bin:$PATH"
lucien> set MY_VAR=value

# Use variables
lucien> echo $PATH
lucien> cd $HOME/projects
```

#### Built-in Commands
- `cd` - Change directory with persistent state
- `pwd` - Print working directory
- `env` - Manage environment variables
- `alias` - Create command aliases
- `history` - Command history
- `exit` - Exit shell

#### Advanced Features
```bash
# Background processes
lucien> long_running_command &

# Command chaining
lucien> make && make test && make install

# Conditional execution
lucien> test -f file.txt && echo "File exists"
```

---

## AI & Intelligence System

### Command Prediction
The AI system learns your patterns and provides intelligent predictions:

```bash
# Get predictions for partial input
lucien ai predict "git"
# Output: git status (85%), git add . (72%), git commit (68%)

# Context-aware suggestions
lucien ai suggestions
# Output: Based on current Python project:
# 1. pip install -r requirements.txt
# 2. python -m pytest
# 3. git status
```

### Real-time Copilot
```bash
# Start AI monitoring
lucien ai start-copilot
lucien ai stop-copilot

# Get suggestions for current work
lucien ai copilot-suggestions main.py --line 42
lucien ai next-step
```

### Error Analysis
```bash
# Analyze error output
lucien ai analyze-error "ImportError: No module named 'requests'" --language python

# Output:
# Error Type: ImportError
# Confidence: 90%
# Suggested Fixes:
#   1. pip install requests
#   2. Add requests to requirements.txt
#   3. Check virtual environment activation
```

---

## BMAD Method Agents

### Build Agent
Handles project building, compilation, and dependency management.

```bash
# Cast build agent
lucien cast agent build "install dependencies"
lucien cast agent build "build python project"
lucien cast agent build "setup node environment"
```

**Capabilities:**
- Auto-detects project type (Python, Node.js, Go, Rust)
- Installs dependencies automatically
- Handles build processes
- Manages virtual environments

### Manage Agent
Environment and system management.

```bash
# Cast manage agent
lucien cast agent manage "show environment"
lucien cast agent manage "check system info"
lucien cast agent manage "manage dependencies"
```

**Capabilities:**
- Environment variable management
- System information display
- Dependency health checks
- Configuration management

### Analyze Agent
Code and project analysis.

```bash
# Cast analyze agent
lucien cast agent analyze "analyze project"
lucien cast agent analyze "count lines"
lucien cast agent analyze "find issues"
```

**Capabilities:**
- Project structure analysis
- Code quality assessment
- Line counting and metrics
- Issue detection (TODOs, security risks)
- Git repository analysis

### MCP Agent (Modular Command Pipeline)
AI-enhanced command processing with natural language support.

```bash
# Cast MCP agent
lucien cast agent mcp_agent "find all python files and count lines"
lucien cast agent mcp_agent "compress images in current folder"
lucien cast agent mcp_agent "analyze git history"
```

**Capabilities:**
- Natural language command interpretation
- AI-powered data filtering and processing
- Complex workflow execution
- Smart command suggestion

---

## Configuration System

### Configuration File
Lucien uses `~/.lucien/config.yaml` for persistent configuration:

```yaml
shell:
  prompt: "lucien> "
  history_size: 1000
  timeout: 30
  
aliases:
  ll: "ls -la"
  grep: "grep --color=auto"
  
environment:
  EDITOR: "nano"
  LUCIEN_HOME: "/home/user/.lucien"
  
agents:
  claude_api_key: ""
  default_model: "claude-3-sonnet-20240229"
  
ui:
  theme: "default"
  colors:
    primary: "blue"
    success: "green"
    error: "red"
```

### Configuration Commands
```bash
# Show current configuration
lucien shell env

# Set environment variables
lucien shell env EDITOR vim

# Create aliases
lucien shell alias ll "ls -la"

# Show command history
lucien shell history
```

---

## PowerShell 7 Parity

### Compatibility Features
Lucien provides comprehensive PowerShell 7 compatibility:

```bash
# Check parity status
lucien ps check-parity

# Translate PowerShell commands
lucien ps translate-command "Get-ChildItem -Recurse"
# Output: find . -type f

# PowerShell compatibility mode
lucien ps powershell-mode
PS > Get-Process
# Automatically translates to: ps aux

# Install PowerShell aliases
lucien ps install-powershell-aliases
```

### Command Mappings
| PowerShell | Lucien Equivalent |
|------------|-------------------|
| `Get-ChildItem` | `ls` |
| `Set-Location` | `cd` |
| `Get-Content` | `cat` |
| `Select-String` | `grep` |
| `Invoke-WebRequest` | `curl` |
| `Get-Process` | `ps` |
| `Stop-Process` | `kill` |

### Feature Parity Status
- ✅ Command Execution
- ✅ Pipeline Operations  
- ✅ Variables & Environment
- ✅ File Operations
- ✅ Process Management
- ✅ Network Operations
- ✅ Script Execution
- ✅ Module System (via Agents)
- ✅ Enhanced Pipeline (AI-powered)
- ⚠️ Remote Operations (via ssh/curl)
- ⚠️ Security Features (basic)

---

## Advanced Usage

### Interactive Shell Mode
```bash
# Start interactive mode
lucien interactive

# Features available in interactive mode:
lucien> # All shell commands
lucien> cast agent build "setup project"  # Agent casting
lucien> ai predict  # AI predictions
lucien> ps powershell-mode  # PowerShell compatibility
```

### Workflow Automation
```bash
# Create custom workflows
lucien ai create-workflow "deploy to production"
lucien ai create-workflow "run full test suite"

# Execute workflows
lucien ai execute-workflow workflow_123456

# Monitor workflow progress
lucien ai list-workflows
```

### Script Integration
```bash
# Execute shell scripts
lucien shell script setup.sh

# Execute Python scripts with context
lucien cast agent mcp_agent "run analysis.py with error handling"

# Batch processing
lucien ai create-workflow "process all data files"
```

### Memory and Learning
```bash
# Lucien automatically maintains memory in ~/.lucien/memory/
# - claude/default.md - Conversation history
# - agents/ - Agent-specific memory
# - sessions/ - Session contexts

# Memory is used for:
# - Command prediction
# - Context-aware suggestions  
# - Workflow optimization
# - Error pattern recognition
```

---

## Development & Extension

### Creating Custom Agents
```python
# agents/my_agent.py
from agents.agent_base import Agent

class My_agent(Agent):
    def run(self, input_text: str) -> str:
        return f"[MyAgent] Processing: {input_text}"
```

### Extending Commands
```python
# cli/commands/my_commands.py
import typer
app = typer.Typer()

@app.command()
def my_command():
    """My custom command"""
    pass
```

### Adding Intelligence
```python
# Extend the intelligence system
from core.intelligence import predictor

# Add custom predictions
predictor.learn_command("my_command", {"context": "custom"})
```

### Configuration Extension
```yaml
# Add to ~/.lucien/config.yaml
my_extension:
  enabled: true
  options:
    custom_setting: value
```

---

## System Requirements

### Minimum Requirements
- Python 3.11+
- 100MB disk space
- 256MB RAM

### Recommended Requirements
- Python 3.11+
- 500MB disk space
- 512MB RAM
- Git installed
- Modern terminal emulator

### Supported Platforms
- Windows 10/11
- macOS 10.15+
- Linux (Ubuntu 20.04+, CentOS 8+, Arch, etc.)

---

## Troubleshooting

### Common Issues
```bash
# Command not found
# Solution: Check if in interactive mode or use full path

# Unicode encoding errors
# Solution: Set environment variable
export PYTHONIOENCODING=utf-8

# Import errors
# Solution: Reinstall dependencies
pip install -r requirements.txt

# Permission errors
# Solution: Check file permissions
chmod +x /path/to/lucien-cli/cli/main.py
```

### Validation
```bash
# Run system validation
lucien validate system
lucien validate agents
lucien validate shell

# Check PowerShell parity
lucien ps check-parity
```

### Debug Mode
```bash
# Enable verbose output
export LUCIEN_DEBUG=1
python -m cli.main --help
```

---

## Performance

### Benchmarks
- Command execution: <100ms average
- AI prediction: <200ms average
- File analysis: <500ms for 1000 files
- Memory usage: ~50MB base, +10MB per active agent

### Optimization
- Commands are cached for faster execution
- AI predictions improve with usage
- Memory is automatically managed
- Background processes are optimized

---

## Security

### Security Features
- Sandboxed agent execution
- Configuration file protection
- Command validation
- Memory encryption (planned)

### Best Practices
- Keep configuration files secure
- Regularly update dependencies
- Use specific agent permissions
- Monitor command history

---

## Contributing

### Development Setup
```bash
git clone https://github.com/your-username/lucien-cli.git
cd lucien-cli
pip install -r requirements.txt
python -m cli.main validate system
```

### Testing
```bash
# Run validation suite
python -m cli.main validate system
python -m cli.main validate agents

# Test PowerShell parity
python -m cli.main ps check-parity
```

### Code Style
- Follow PEP 8
- Use type hints
- Document all functions
- Write tests for new features

---

## License

MIT License - see LICENSE file for details.

---

## Support

- **Documentation**: This file
- **Issues**: GitHub Issues
- **Discussions**: GitHub Discussions
- **Wiki**: GitHub Wiki

---

*Lucien CLI - Where Intelligence Meets Command Line*