# 📚 LUCIEN CLI - API DOCUMENTATION
**Version**: 1.0-alpha  
**Last Updated**: 2025-08-06  
**Status**: Development API Reference  

## 🎯 OVERVIEW

This document provides comprehensive API documentation for the Lucien CLI system, covering both the Go core engine and Python CLI layer interfaces.

**⚠️ IMPORTANT**: This API is under active development. Some interfaces may change before production release.

## 📋 TABLE OF CONTENTS

1. [Go Core API](#go-core-api)
2. [Shell Engine API](#shell-engine-api)
3. [Plugin System API](#plugin-system-api)
4. [Policy Engine API](#policy-engine-api)
5. [AI Engine API](#ai-engine-api)
6. [Python CLI API](#python-cli-api)
7. [RPC Interfaces](#rpc-interfaces)
8. [Configuration API](#configuration-api)
9. [Error Handling](#error-handling)
10. [Examples](#examples)

## 🛠️ GO CORE API

### Main Package Imports
```go
import (
    "github.com/luciendev/lucien-core/internal/shell"
    "github.com/luciendev/lucien-core/internal/plugin"
    "github.com/luciendev/lucien-core/internal/policy"
    "github.com/luciendev/lucien-core/internal/ai"
    "github.com/luciendev/lucien-core/internal/sandbox"
)
```

### Core Configuration Structure
```go
// Config holds shell configuration
type Config struct {
    PolicyEngine *policy.Engine
    PluginMgr    *plugin.Manager
    SandboxMgr   *sandbox.Manager
    AIEngine     *ai.Engine
    SafeMode     bool
}
```

## 🐚 SHELL ENGINE API

### Core Types

#### ExecutionResult
```go
// ExecutionResult holds command execution results
type ExecutionResult struct {
    Output   string        // Command standard output
    Error    string        // Error message if any
    ExitCode int           // Process exit code
    Duration time.Duration // Execution time (⚠️ Currently broken for builtins)
}
```

#### Shell Instance
```go
// Shell represents the core shell engine
type Shell struct {
    config      *Config                                    // Shell configuration
    env         map[string]string                         // Environment variables
    aliases     map[string]string                         // Command aliases
    history     []string                                  // Command history
    currentDir  string                                    // Current working directory
    builtins    map[string]func([]string) (*ExecutionResult, error) // Builtin commands
}
```

#### Command Structure
```go
// Command represents a parsed command with pipes and redirects
type Command struct {
    Name      string            // Command name
    Args      []string          // Command arguments
    Input     io.Reader         // Input reader
    Output    io.Writer         // Output writer
    Error     io.Writer         // Error writer
    Pipes     []*Command        // Piped commands
    Redirects map[string]string // Redirections: ">" -> filename, "<" -> filename
}
```

### Shell Methods

#### Constructor
```go
// New creates a new shell instance
func New(config *Config) *Shell

// Usage:
config := &shell.Config{
    SafeMode: true,
    PolicyEngine: policyEngine,
    PluginMgr: pluginManager,
}
sh := shell.New(config)
```

#### Core Execution
```go
// Execute runs a command string through the shell
func (s *Shell) Execute(cmdLine string) (*ExecutionResult, error)

// Usage:
result, err := sh.Execute("ls -la | grep .go")
if err != nil {
    log.Fatal(err)
}
fmt.Println(result.Output)
```

#### Variable Management
```go
// SetVariable sets a shell variable
func (s *Shell) SetVariable(name, value string) error

// GetVariable retrieves a shell variable
func (s *Shell) GetVariable(name string) (string, bool)

// Usage:
sh.SetVariable("PATH", "/usr/bin:/usr/local/bin")
path, exists := sh.GetVariable("PATH")
```

#### Alias Management
```go
// CreateAlias creates a command alias
func (s *Shell) CreateAlias(name, command string) error

// RemoveAlias removes an alias
func (s *Shell) RemoveAlias(name string) error

// ListAliases returns all defined aliases
func (s *Shell) ListAliases() map[string]string

// Usage:
sh.CreateAlias("ll", "ls -la")
sh.CreateAlias("grep-go", "grep -n -r --include='*.go'")
```

#### History Management
```go
// GetHistory returns command history
func (s *Shell) GetHistory() []string

// AddToHistory adds command to history
func (s *Shell) AddToHistory(cmd string)

// ClearHistory clears command history
func (s *Shell) ClearHistory()

// Usage:
history := sh.GetHistory()
for i, cmd := range history {
    fmt.Printf("%d: %s\n", i+1, cmd)
}
```

### Built-in Commands

#### Available Built-ins
| Command | Function | Purpose |
|---------|----------|---------|
| `cd` | `changeDirectory` | Change working directory |
| `set` | `setVariable` | Set shell variable |
| `export` | `exportVariable` | Export variable to environment |
| `alias` | `createAlias` | Create command alias |
| `history` | `showHistory` | Display command history |
| `pwd` | `printWorkingDirectory` | Print current directory |
| `echo` | `echo` | Display text |
| `exit` | `exit` | Exit shell |

#### Built-in Function Signatures
```go
// Built-in command signature
type BuiltinFunc func([]string) (*ExecutionResult, error)

// Example implementations:
func (s *Shell) changeDirectory(args []string) (*ExecutionResult, error)
func (s *Shell) setVariable(args []string) (*ExecutionResult, error)
func (s *Shell) createAlias(args []string) (*ExecutionResult, error)
```

## 🔌 PLUGIN SYSTEM API

### Plugin Interface

#### Core Plugin Interface
```go
// PluginInterface defines the interface that all Lucien plugins must implement
type PluginInterface interface {
    Execute(ctx context.Context, command string, args []string) (*Result, error)
    GetInfo() (*Info, error)
    Initialize(config map[string]interface{}) error
}
```

#### Plugin Result Structure
```go
// Result contains plugin execution results
type Result struct {
    Output   string                 `json:"output"`      // Plugin output
    Error    string                 `json:"error,omitempty"` // Error message
    ExitCode int                    `json:"exit_code"`   // Exit code
    Data     map[string]interface{} `json:"data,omitempty"` // Additional data
}
```

#### Plugin Info Structure
```go
// Info contains plugin metadata
type Info struct {
    Name         string            `json:"name"`         // Plugin name
    Version      string            `json:"version"`      // Plugin version
    Description  string            `json:"description"`  // Description
    Author       string            `json:"author"`       // Author name
    Capabilities []string          `json:"capabilities"` // Plugin capabilities
    Config       map[string]string `json:"config"`       // Configuration options
}
```

### Plugin Manifest

#### Manifest Structure
```go
// Manifest describes a plugin's configuration
type Manifest struct {
    Name         string            `json:"name"`         // Plugin name
    Version      string            `json:"version"`      // Version string
    Description  string            `json:"description"`  // Plugin description
    Author       string            `json:"author"`       // Author information
    Executable   string            `json:"executable"`   // Executable filename
    Capabilities []string          `json:"capabilities"` // Supported capabilities
    Config       map[string]string `json:"config"`       // Default configuration
    Dependencies []string          `json:"dependencies"` // Plugin dependencies
}
```

#### Example Manifest File
```json
{
    "name": "bmad",
    "version": "1.0.0",
    "description": "Build, Manage, Analyze, Deploy workflow plugin",
    "author": "Lucien Dev Team",
    "executable": "bmad.exe",
    "capabilities": [
        "build_management",
        "system_analysis",
        "deployment_automation"
    ],
    "config": {
        "timeout": "30s",
        "max_concurrent": "5"
    },
    "dependencies": []
}
```

### Plugin Manager

#### Manager Structure
```go
// Manager handles plugin lifecycle and execution
type Manager struct {
    pluginDir   string                      // Plugin directory path
    plugins     map[string]*LoadedPlugin    // Loaded plugins
    manifests   map[string]*Manifest        // Plugin manifests
    mu          sync.RWMutex               // Thread safety
}
```

#### Manager Methods
```go
// New creates a new plugin manager
func New(pluginDir string) *Manager

// LoadPlugins discovers and loads all plugins from the plugin directory
func (m *Manager) LoadPlugins() error

// LoadPlugin loads a specific plugin
func (m *Manager) LoadPlugin(name string) error

// UnloadPlugin unloads a specific plugin
func (m *Manager) UnloadPlugin(name string) error

// ExecutePlugin executes a plugin command
func (m *Manager) ExecutePlugin(name, command string, args []string) (*Result, error)

// ListPlugins returns information about all loaded plugins
func (m *Manager) ListPlugins() map[string]*Info

// GetPluginInfo returns information about a specific plugin
func (m *Manager) GetPluginInfo(name string) (*Info, error)
```

#### Usage Examples
```go
// Initialize plugin manager
manager := plugin.New("./plugins")

// Load all plugins
if err := manager.LoadPlugins(); err != nil {
    log.Fatal(err)
}

// Execute plugin command
result, err := manager.ExecutePlugin("bmad", "build", []string{"--verbose"})
if err != nil {
    log.Fatal(err)
}
fmt.Println(result.Output)

// List available plugins
plugins := manager.ListPlugins()
for name, info := range plugins {
    fmt.Printf("%s v%s: %s\n", name, info.Version, info.Description)
}
```

### Plugin RPC Interface

#### RPC Configuration
```go
// Handshake configuration for plugins
var Handshake = plugin.HandshakeConfig{
    ProtocolVersion:  1,
    MagicCookieKey:   "LUCIEN_PLUGIN",
    MagicCookieValue: "lucien_neural_interface",
}
```

#### RPC Client
```go
// Plugin RPC implementation
type PluginRPC struct {
    client *rpc.Client
}

// Execute calls the plugin's Execute method via RPC
func (p *PluginRPC) Execute(ctx context.Context, command string, args []string) (*Result, error)

// GetInfo calls the plugin's GetInfo method via RPC
func (p *PluginRPC) GetInfo() (*Info, error)
```

## 🛡️ POLICY ENGINE API

### Policy Types

#### Policy Engine Structure
```go
// Engine handles security policy evaluation
type Engine struct {
    rules    []PolicyRule     // Policy rules
    policies map[string]*Policy // Loaded policies
    mu       sync.RWMutex    // Thread safety
}
```

#### Policy Rule Structure
```go
// PolicyRule represents a security policy rule
type PolicyRule struct {
    Name        string       `json:"name"`        // Rule name
    Pattern     string       `json:"pattern"`     // Match pattern (regex)
    Action      ActionType   `json:"action"`      // Action to take
    Severity    SeverityLevel `json:"severity"`   // Risk severity
    Description string       `json:"description"` // Rule description
}
```

#### Policy Enums
```go
// ActionType defines policy actions
type ActionType int

const (
    Allow ActionType = iota  // Allow command execution
    Deny                     // Deny command execution
    Warn                     // Allow with warning
    Confirm                  // Require user confirmation
)

// SeverityLevel defines risk levels
type SeverityLevel int

const (
    Low SeverityLevel = iota    // Low risk
    Medium                      // Medium risk
    High                        // High risk
    Critical                    // Critical risk
)
```

#### Policy Evaluation Result
```go
// PolicyResult contains policy evaluation results
type PolicyResult struct {
    Allowed     bool          `json:"allowed"`      // Whether command is allowed
    Action      ActionType    `json:"action"`       // Action taken
    Severity    SeverityLevel `json:"severity"`     // Risk severity
    Message     string        `json:"message"`      // Result message
    MatchedRule string        `json:"matched_rule"` // Matched rule name
}
```

### Policy Engine Methods

#### Constructor and Configuration
```go
// New creates a new policy engine
func New() *Engine

// LoadPolicies loads policy rules from directory
func (e *Engine) LoadPolicies(dir string) error

// AddRule adds a single policy rule
func (e *Engine) AddRule(rule PolicyRule) error

// RemoveRule removes a policy rule
func (e *Engine) RemoveRule(name string) error
```

#### Policy Evaluation
```go
// Evaluate evaluates a command against all policy rules
func (e *Engine) Evaluate(command string) (*PolicyResult, error)

// EvaluateWithContext evaluates with additional context
func (e *Engine) EvaluateWithContext(command string, context map[string]interface{}) (*PolicyResult, error)

// Usage:
engine := policy.New()
engine.LoadPolicies("./policies")

result, err := engine.Evaluate("rm -rf /")
if err != nil {
    log.Fatal(err)
}

if !result.Allowed {
    fmt.Printf("Command blocked: %s\n", result.Message)
}
```

#### Policy Management
```go
// ListRules returns all loaded policy rules
func (e *Engine) ListRules() []PolicyRule

// GetRule returns a specific policy rule
func (e *Engine) GetRule(name string) (*PolicyRule, error)

// ValidateRule validates a policy rule
func (e *Engine) ValidateRule(rule PolicyRule) error
```

## 🧠 AI ENGINE API

### AI Engine Structure

#### Engine Configuration
```go
// Engine handles AI integration and query processing
type Engine struct {
    config   *Config      // AI configuration
    provider string       // AI provider (local, openai, anthropic)
    client   interface{}  // Provider-specific client
}

// Config holds AI engine configuration
type Config struct {
    Provider    string            `json:"provider"`     // AI provider
    ModelPath   string            `json:"model_path"`   // Local model path
    APIKey      string            `json:"api_key"`      // API key for cloud providers
    MaxTokens   int               `json:"max_tokens"`   // Maximum tokens
    Temperature float32           `json:"temperature"`  // Response randomness
    Timeout     time.Duration     `json:"timeout"`      // Request timeout
    Settings    map[string]string `json:"settings"`     // Additional settings
}
```

#### AI Response Structure
```go
// Response contains AI query results
type Response struct {
    Content   string            `json:"content"`    // AI response content
    Tokens    int               `json:"tokens"`     // Tokens used
    Model     string            `json:"model"`      // Model used
    Duration  time.Duration     `json:"duration"`   // Response time
    Metadata  map[string]interface{} `json:"metadata"` // Additional metadata
}
```

### AI Engine Methods

#### Constructor and Configuration
```go
// New creates a new AI engine
func New(config *Config) (*Engine, error)

// Initialize initializes the AI engine
func (e *Engine) Initialize() error

// SetProvider changes the AI provider
func (e *Engine) SetProvider(provider string) error
```

#### Query Interface
```go
// Query sends a query to the AI engine
func (e *Engine) Query(prompt string) (*Response, error)

// QueryWithContext sends a query with additional context
func (e *Engine) QueryWithContext(prompt string, context map[string]interface{}) (*Response, error)

// Usage:
config := &ai.Config{
    Provider: "local",
    ModelPath: "/models/llama-7b.gguf",
    MaxTokens: 1000,
    Temperature: 0.7,
}

engine, err := ai.New(config)
if err != nil {
    log.Fatal(err)
}

response, err := engine.Query("Explain how to optimize a Go application")
if err != nil {
    log.Fatal(err)
}
fmt.Println(response.Content)
```

#### Provider-Specific Methods
```go
// QueryLocal queries local AI model (⚠️ Currently placeholder)
func (e *Engine) QueryLocal(prompt string) (*Response, error)

// QueryOpenAI queries OpenAI API (⚠️ Not implemented)
func (e *Engine) QueryOpenAI(prompt string) (*Response, error)

// QueryAnthropic queries Anthropic Claude API (⚠️ Not implemented)
func (e *Engine) QueryAnthropic(prompt string) (*Response, error)
```

## 🐍 PYTHON CLI API

### Main CLI Application

#### CLI Structure
```python
import typer
from rich.console import Console
from pathlib import Path

app = typer.Typer()
console = Console()
```

#### Main Application Interface
```python
@app.callback()
def main(
    safe: bool = typer.Option(False, "--safe", help="Enable safe mode"),
    test: bool = typer.Option(False, "--test", help="Run in test mode"),
    debug: bool = typer.Option(False, "--debug", help="Enable debug output")
):
    """Lucien CLI - AI-Enhanced Shell Replacement"""
```

#### Interactive Mode
```python
@app.command()
def interactive():
    """Start Lucien in interactive shell mode"""
```

### UI Module

#### UI Functions
```python
from rich.console import Console
from rich.panel import Panel
from rich.text import Text

console = Console()

def spell_thinking(message: str):
    """Display thinking animation with cyberpunk styling"""
    
def spell_complete(result: str):
    """Display completion message with success styling"""
    
def spell_failed(error: str):
    """Display error message with failure styling"""
```

#### Usage Examples
```python
from lucien.ui import console, spell_complete, spell_failed

# Display status
console.print("[cyan]Initializing Lucien CLI...[/cyan]")

# Show success
spell_complete("System initialization complete")

# Show error
spell_failed("Configuration file not found")
```

### Safety Manager (Conceptual)

#### Safety Interface
```python
class SafetyManager:
    """Manages command safety validation"""
    
    def __init__(self, safe_mode: bool = True):
        self.safe_mode = safe_mode
        self.dangerous_patterns = [
            r'rm\s+-rf\s+/',
            r'del\s+/s',
            r'format\s+c:',
            # ... more patterns
        ]
    
    def validate_command(self, command: str) -> ValidationResult:
        """Validate command safety"""
        
    def is_dangerous(self, command: str) -> bool:
        """Check if command is dangerous"""
```

## 🔌 RPC INTERFACES

### Plugin RPC Protocol

#### RPC Method Definitions
```go
// Plugin RPC methods
type PluginRPCServer struct {
    Impl PluginInterface
}

// Execute RPC method
func (s *PluginRPCServer) Execute(args *ExecuteArgs, reply *Result) error

// GetInfo RPC method  
func (s *PluginRPCServer) GetInfo(args struct{}, reply *Info) error

// Initialize RPC method
func (s *PluginRPCServer) Initialize(args map[string]interface{}, reply *bool) error
```

#### RPC Argument Structures
```go
// ExecuteArgs contains RPC execute arguments
type ExecuteArgs struct {
    Command string   `json:"command"`  // Command to execute
    Args    []string `json:"args"`     // Command arguments
    Context map[string]interface{} `json:"context"` // Execution context
}

// RPC response wrapper
type RPCResponse struct {
    Result *Result `json:"result"`
    Error  string  `json:"error,omitempty"`
}
```

### Client-Server Communication

#### Plugin Client Setup
```go
// Plugin client configuration
clientConfig := &plugin.ClientConfig{
    HandshakeConfig: Handshake,
    Plugins: map[string]plugin.Plugin{
        "plugin": &PluginRPCPlugin{},
    },
    Cmd: exec.Command("./plugin-binary"),
}

client := plugin.NewClient(clientConfig)
rpcClient, err := client.Client()
```

## ⚙️ CONFIGURATION API

### Configuration Structure

#### Main Configuration
```go
// SystemConfig holds all system configuration
type SystemConfig struct {
    Shell    ShellConfig    `toml:"shell"`    // Shell configuration
    UI       UIConfig       `toml:"ui"`       // UI configuration
    Plugins  PluginConfig   `toml:"plugins"`  // Plugin configuration
    Security SecurityConfig `toml:"security"` // Security configuration
    AI       AIConfig       `toml:"ai"`       // AI configuration
}
```

#### Component Configurations
```go
// ShellConfig holds shell-specific configuration
type ShellConfig struct {
    SafeMode     bool   `toml:"safe_mode"`     // Enable safe mode
    HistorySize  int    `toml:"history_size"`  // Command history size
    AutoComplete bool   `toml:"auto_complete"` // Enable auto-completion
    DefaultDir   string `toml:"default_dir"`   // Default directory
}

// UIConfig holds UI configuration
type UIConfig struct {
    Theme        string `toml:"theme"`         // UI theme (nexus, synthwave, ghost)
    Animations   bool   `toml:"animations"`    // Enable animations
    Colors       bool   `toml:"colors"`        // Enable colors
    GlitchEffect bool   `toml:"glitch_effect"` // Enable glitch effects
}

// PluginConfig holds plugin configuration
type PluginConfig struct {
    AutoLoad       []string `toml:"auto_load"`       // Auto-load plugins
    PluginTimeout  int      `toml:"plugin_timeout"`  // Plugin timeout (seconds)
    MaxConcurrent  int      `toml:"max_concurrent"`  // Max concurrent plugins
    PluginDir      string   `toml:"plugin_dir"`      // Plugin directory
}
```

### Configuration Methods

#### Configuration Loading
```go
// LoadConfig loads configuration from file
func LoadConfig(filepath string) (*SystemConfig, error)

// SaveConfig saves configuration to file
func SaveConfig(config *SystemConfig, filepath string) error

// ValidateConfig validates configuration
func ValidateConfig(config *SystemConfig) error

// GetDefaultConfig returns default configuration
func GetDefaultConfig() *SystemConfig
```

#### Configuration Usage
```go
// Load configuration
config, err := LoadConfig("~/.lucien/config.toml")
if err != nil {
    config = GetDefaultConfig()
}

// Modify configuration
config.Shell.SafeMode = true
config.UI.Theme = "synthwave"

// Save configuration
if err := SaveConfig(config, "~/.lucien/config.toml"); err != nil {
    log.Fatal(err)
}
```

## ❌ ERROR HANDLING

### Error Types

#### Core Error Types
```go
// Error types used throughout the system
var (
    ErrCommandNotFound    = errors.New("command not found")
    ErrInvalidArguments   = errors.New("invalid arguments")
    ErrPermissionDenied   = errors.New("permission denied")
    ErrPluginNotFound     = errors.New("plugin not found")
    ErrPolicyViolation    = errors.New("policy violation")
    ErrConfigurationError = errors.New("configuration error")
)
```

#### Custom Error Structures
```go
// ExecutionError provides detailed execution error information
type ExecutionError struct {
    Command   string    // Failed command
    ExitCode  int       // Exit code
    Message   string    // Error message
    Timestamp time.Time // When error occurred
    Context   map[string]interface{} // Additional context
}

func (e *ExecutionError) Error() string {
    return fmt.Sprintf("command '%s' failed with exit code %d: %s", 
        e.Command, e.ExitCode, e.Message)
}
```

#### Error Wrapping
```go
// Wrap errors with additional context
func WrapError(err error, message string, args ...interface{}) error {
    return fmt.Errorf(message+": %w", append(args, err)...)
}

// Usage:
if err := someOperation(); err != nil {
    return WrapError(err, "failed to execute command '%s'", command)
}
```

### Error Response Patterns

#### Standard Error Response
```go
// Standard error response pattern
type ErrorResponse struct {
    Error   string            `json:"error"`             // Error message
    Code    string            `json:"code"`              // Error code
    Details map[string]interface{} `json:"details,omitempty"` // Error details
}

// Create error response
func NewErrorResponse(err error, code string) *ErrorResponse {
    return &ErrorResponse{
        Error: err.Error(),
        Code:  code,
        Details: make(map[string]interface{}),
    }
}
```

## 📝 EXAMPLES

### Complete Shell Integration Example

```go
package main

import (
    "log"
    "github.com/luciendev/lucien-core/internal/shell"
    "github.com/luciendev/lucien-core/internal/plugin"
    "github.com/luciendev/lucien-core/internal/policy"
)

func main() {
    // Initialize components
    policyEngine := policy.New()
    pluginManager := plugin.New("./plugins")
    
    // Load policies and plugins
    if err := policyEngine.LoadPolicies("./policies"); err != nil {
        log.Printf("Warning: could not load policies: %v", err)
    }
    
    if err := pluginManager.LoadPlugins(); err != nil {
        log.Printf("Warning: could not load plugins: %v", err)
    }
    
    // Create shell configuration
    config := &shell.Config{
        PolicyEngine: policyEngine,
        PluginMgr:    pluginManager,
        SafeMode:     true,
    }
    
    // Initialize shell
    sh := shell.New(config)
    
    // Set up some variables and aliases
    sh.SetVariable("PROJECT_DIR", "/home/user/projects")
    sh.CreateAlias("ll", "ls -la")
    sh.CreateAlias("grep-go", "grep -r --include='*.go'")
    
    // Execute commands
    commands := []string{
        "set GREETING 'Hello, World!'",
        "echo $GREETING",
        "ls -la | grep .go",
        "ll",
    }
    
    for _, cmd := range commands {
        result, err := sh.Execute(cmd)
        if err != nil {
            log.Printf("Error executing '%s': %v", cmd, err)
            continue
        }
        
        if result.Output != "" {
            log.Printf("Output: %s", result.Output)
        }
        
        log.Printf("Command '%s' completed in %v", cmd, result.Duration)
    }
}
```

### Plugin Development Example

```go
package main

import (
    "context"
    "fmt"
    "github.com/hashicorp/go-plugin"
    lucienPlugin "github.com/luciendev/lucien-core/internal/plugin"
)

// MyPlugin implements the PluginInterface
type MyPlugin struct{}

func (p *MyPlugin) Execute(ctx context.Context, command string, args []string) (*lucienPlugin.Result, error) {
    switch command {
    case "hello":
        return &lucienPlugin.Result{
            Output:   fmt.Sprintf("Hello from plugin! Args: %v", args),
            ExitCode: 0,
        }, nil
        
    case "version":
        return &lucienPlugin.Result{
            Output:   "MyPlugin v1.0.0",
            ExitCode: 0,
        }, nil
        
    default:
        return &lucienPlugin.Result{
            Error:    fmt.Sprintf("unknown command: %s", command),
            ExitCode: 1,
        }, nil
    }
}

func (p *MyPlugin) GetInfo() (*lucienPlugin.Info, error) {
    return &lucienPlugin.Info{
        Name:        "myplugin",
        Version:     "1.0.0",
        Description: "Example plugin for Lucien CLI",
        Author:      "Plugin Developer",
        Capabilities: []string{"hello", "version"},
    }, nil
}

func (p *MyPlugin) Initialize(config map[string]interface{}) error {
    // Plugin initialization logic
    return nil
}

func main() {
    plugin.Serve(&plugin.ServeConfig{
        HandshakeConfig: lucienPlugin.Handshake,
        Plugins: map[string]plugin.Plugin{
            "myplugin": &lucienPlugin.PluginRPCPlugin{Impl: &MyPlugin{}},
        },
    })
}
```

### Python CLI Integration Example

```python
#!/usr/bin/env python3
import typer
from rich.console import Console
from pathlib import Path
import subprocess
import json

app = typer.Typer()
console = Console()

@app.command()
def shell_command(
    command: str = typer.Argument(..., help="Command to execute"),
    safe: bool = typer.Option(True, "--safe/--no-safe", help="Use safe mode")
):
    """Execute a command through Lucien shell"""
    
    # Prepare Go binary execution
    cmd_args = ["./lucien.exe"]
    if safe:
        cmd_args.append("--safe")
    
    try:
        # Execute Go binary with command
        result = subprocess.run(
            cmd_args,
            input=command,
            capture_output=True,
            text=True,
            timeout=30
        )
        
        if result.returncode == 0:
            console.print(f"[green]✓[/green] {result.stdout}")
        else:
            console.print(f"[red]✗[/red] {result.stderr}")
            
    except subprocess.TimeoutExpired:
        console.print("[red]Command timed out[/red]")
    except Exception as e:
        console.print(f"[red]Error: {e}[/red]")

@app.command() 
def plugin_info():
    """Display information about loaded plugins"""
    
    try:
        result = subprocess.run(
            ["./lucien.exe", "--plugin-info"],
            capture_output=True,
            text=True
        )
        
        if result.returncode == 0:
            plugins = json.loads(result.stdout)
            
            console.print("[cyan]Loaded Plugins:[/cyan]")
            for name, info in plugins.items():
                console.print(f"  [bold]{name}[/bold] v{info['version']}")
                console.print(f"    {info['description']}")
                console.print(f"    Capabilities: {', '.join(info['capabilities'])}")
        else:
            console.print(f"[red]Error: {result.stderr}[/red]")
            
    except Exception as e:
        console.print(f"[red]Error: {e}[/red]")

if __name__ == "__main__":
    app()
```

---

## 🚧 API LIMITATIONS AND KNOWN ISSUES

### Current API Limitations

#### ⚠️ Known Issues
1. **Variable Expansion**: Undefined variables return literal text instead of empty strings
2. **Duration Tracking**: Builtin commands don't properly track execution duration
3. **Policy Engine**: Compilation errors prevent policy testing
4. **AI Integration**: Most AI methods are placeholder implementations
5. **Cross-Platform**: Limited non-Windows platform support

#### 🔄 Planned Improvements
1. Complete AI integration with actual model inference
2. Enhanced plugin security and sandboxing
3. Comprehensive policy rule engine
4. Full cross-platform compatibility
5. Performance optimizations

### Breaking Changes Warning

**This API is subject to change before version 1.0 release.** Major breaking changes may occur as:
- Security vulnerabilities are addressed
- Core functionality bugs are fixed
- Performance optimizations are implemented
- Feature completeness is achieved

---

**API Documentation Status**: 🔶 Development Version - Subject to Change

*This documentation reflects the current state of the Lucien CLI API as of 2025-08-06. Refer to the source code and test files for the most up-to-date interface definitions.*