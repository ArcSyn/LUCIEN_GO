package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/luciendev/lucien-core/internal/ai"
	"github.com/luciendev/lucien-core/internal/plugin"
	"github.com/luciendev/lucien-core/internal/policy"
	"github.com/luciendev/lucien-core/internal/sandbox"
)

// ExecutionResult holds command execution results
type ExecutionResult struct {
	Output   string
	Error    string
	ExitCode int
	Duration time.Duration
}

// Config holds shell configuration
type Config struct {
	PolicyEngine *policy.Engine
	PluginMgr    *plugin.Manager
	SandboxMgr   *sandbox.Manager
	AIEngine     *ai.Engine
	SafeMode     bool
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

// Command represents a parsed command with pipes and redirects
type Command struct {
	Name      string
	Args      []string
	Input     io.Reader
	Output    io.Writer
	Error     io.Writer
	Pipes     []*Command
	Redirects map[string]string // ">" -> filename, "<" -> filename
}

// New creates a new shell instance
func New(config *Config) *Shell {
	homeDir, _ := os.UserHomeDir()
	
	shell := &Shell{
		config:     config,
		env:        make(map[string]string),
		aliases:    make(map[string]string),
		history:    []string{},
		currentDir: homeDir,
		builtins:   make(map[string]func([]string) (*ExecutionResult, error)),
	}

	// Copy environment variables
	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			shell.env[parts[0]] = parts[1]
		}
	}

	// Register built-in commands
	shell.registerBuiltins()

	return shell
}

func (s *Shell) registerBuiltins() {
	s.builtins["cd"] = s.changeDirectory
	s.builtins["set"] = s.setVariable
	s.builtins["alias"] = s.createAlias
	s.builtins["exit"] = s.exit
	s.builtins["history"] = s.showHistory
	s.builtins["pwd"] = s.printWorkingDirectory
	s.builtins["echo"] = s.echo
	s.builtins["export"] = s.exportVariable
}

// Execute runs a command string through the shell
func (s *Shell) Execute(cmdLine string) (*ExecutionResult, error) {
	start := time.Now()
	
	// Add to history
	s.history = append(s.history, cmdLine)
	
	// Parse command line
	commands, err := s.parseCommandLine(cmdLine)
	if err != nil {
		return &ExecutionResult{
			Error:    err.Error(),
			ExitCode: 1,
			Duration: time.Since(start),
		}, err
	}

	if len(commands) == 0 {
		return &ExecutionResult{Duration: time.Since(start)}, nil
	}

	// Execute command pipeline
	result, err := s.executePipeline(commands)
	result.Duration = time.Since(start)
	
	return result, err
}

func (s *Shell) parseCommandLine(cmdLine string) ([]*Command, error) {
	// Expand environment variables
	cmdLine = s.expandVariables(cmdLine)
	
	// Split by pipes
	pipeSegments := strings.Split(cmdLine, "|")
	commands := make([]*Command, 0, len(pipeSegments))

	for _, segment := range pipeSegments {
		cmd, err := s.parseCommand(strings.TrimSpace(segment))
		if err != nil {
			return nil, err
		}
		if cmd != nil { // Skip nil commands (empty segments)
			commands = append(commands, cmd)
		}
	}

	return commands, nil
}

func (s *Shell) parseCommand(cmdStr string) (*Command, error) {
	if cmdStr == "" || strings.TrimSpace(cmdStr) == "" {
		return nil, nil // Return nil without error for empty commands
	}

	// Tokenize respecting quotes and redirects
	tokens := s.tokenize(cmdStr)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no command specified")
	}

	cmd := &Command{
		Name:      tokens[0],
		Args:      []string{},
		Redirects: make(map[string]string),
	}

	// Process tokens for arguments and redirects
	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		
		switch token {
		case ">":
			if i+1 < len(tokens) {
				cmd.Redirects[">"] = tokens[i+1]
				i++ // Skip the filename
			}
		case "<":
			if i+1 < len(tokens) {
				cmd.Redirects["<"] = tokens[i+1]
				i++ // Skip the filename
			}
		case ">>":
			if i+1 < len(tokens) {
				cmd.Redirects[">>"] = tokens[i+1]
				i++ // Skip the filename
			}
		default:
			cmd.Args = append(cmd.Args, token)
		}
	}

	// Check if command is aliased
	if alias, exists := s.aliases[cmd.Name]; exists {
		aliasTokens := s.tokenize(alias)
		if len(aliasTokens) > 0 {
			cmd.Name = aliasTokens[0]
			cmd.Args = append(aliasTokens[1:], cmd.Args...)
		}
	}

	return cmd, nil
}

func (s *Shell) tokenize(input string) []string {
	var tokens []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(input); i++ {
		char := input[i]

		switch {
		case char == '"' || char == '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
				quoteChar = 0
			} else {
				current.WriteByte(char)
			}

		case char == ' ' || char == '\t':
			if inQuotes {
				current.WriteByte(char)
			} else if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

		case char == '>' || char == '<':
			if inQuotes {
				current.WriteByte(char)
			} else {
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}
				
				// Check for >> redirect
				if char == '>' && i+1 < len(input) && input[i+1] == '>' {
					tokens = append(tokens, ">>")
					i++ // Skip next >
				} else {
					tokens = append(tokens, string(char))
				}
			}

		default:
			current.WriteByte(char)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func (s *Shell) expandVariables(input string) string {
	result := input
	
	// Simple variable expansion for $VAR and ${VAR}
	for key, value := range s.env {
		result = strings.ReplaceAll(result, "$"+key, value)
		result = strings.ReplaceAll(result, "${"+key+"}", value)
	}
	
	// Handle undefined variables - replace with empty string
	// Pattern: $VARNAME where VARNAME contains letters, digits, underscore
	var i int
	for i < len(result) {
		if result[i] == '$' && i+1 < len(result) {
			start := i + 1
			end := start
			
			// Find end of variable name
			for end < len(result) && (result[end] >= 'A' && result[end] <= 'Z' ||
				result[end] >= 'a' && result[end] <= 'z' ||
				result[end] >= '0' && result[end] <= '9' ||
				result[end] == '_') {
				end++
			}
			
			if end > start {
				varName := result[start:end]
				// Only replace if variable is not defined
				if _, exists := s.env[varName]; !exists {
					result = result[:i] + result[end:]
					continue
				}
			}
		}
		i++
	}
	
	return result
}

func (s *Shell) executePipeline(commands []*Command) (*ExecutionResult, error) {
	if len(commands) == 1 {
		return s.executeCommand(commands[0])
	}

	// Set up pipes between commands
	var output strings.Builder
	var lastOutput io.Reader

	for i, cmd := range commands {
		cmd.Input = lastOutput

		if i < len(commands)-1 {
			// Create pipe for all but last command
			reader, writer := io.Pipe()
			cmd.Output = writer
			lastOutput = reader
		} else {
			// Last command outputs to our buffer
			cmd.Output = &output
		}

		// Execute command (in production, would run concurrently)
		result, err := s.executeCommand(cmd)
		if err != nil {
			return result, err
		}
	}

	return &ExecutionResult{
		Output:   output.String(),
		ExitCode: 0,
	}, nil
}

func (s *Shell) executeCommand(cmd *Command) (*ExecutionResult, error) {
	// Policy check if safe mode is enabled
	if s.config.SafeMode && s.config.PolicyEngine != nil {
		allowed, err := s.config.PolicyEngine.Authorize("execute", cmd.Name, cmd.Args)
		if err != nil {
			return &ExecutionResult{
				Error:    fmt.Sprintf("Policy check failed: %v", err),
				ExitCode: 1,
			}, err
		}
		if !allowed {
			return &ExecutionResult{
				Error:    fmt.Sprintf("Command '%s' blocked by security policy", cmd.Name),
				ExitCode: 1,
			}, fmt.Errorf("command blocked by policy")
		}
	}

	// Check if it's a builtin command
	if builtin, exists := s.builtins[cmd.Name]; exists {
		return builtin(cmd.Args)
	}

	// Execute external command
	return s.executeExternal(cmd)
}

func (s *Shell) executeExternal(cmd *Command) (*ExecutionResult, error) {
	var execCmd *exec.Cmd
	
	// Handle platform-specific command execution
	if runtime.GOOS == "windows" {
		execCmd = exec.Command("cmd", "/C", cmd.Name+" "+strings.Join(cmd.Args, " "))
	} else {
		execCmd = exec.Command(cmd.Name, cmd.Args...)
	}

	// Set working directory
	execCmd.Dir = s.currentDir

	// Set environment
	execCmd.Env = s.buildEnvSlice()

	// Handle redirects and pipes
	if cmd.Input != nil {
		execCmd.Stdin = cmd.Input
	}

	var output strings.Builder
	var errorOutput strings.Builder

	if cmd.Output != nil {
		execCmd.Stdout = cmd.Output
	} else {
		execCmd.Stdout = &output
	}
	execCmd.Stderr = &errorOutput

	// Handle file redirects
	if outFile, exists := cmd.Redirects[">"]; exists {
		file, err := os.Create(outFile)
		if err != nil {
			return &ExecutionResult{
				Error:    fmt.Sprintf("Cannot create output file: %v", err),
				ExitCode: 1,
			}, err
		}
		defer file.Close()
		execCmd.Stdout = file
	}

	if inFile, exists := cmd.Redirects["<"]; exists {
		file, err := os.Open(inFile)
		if err != nil {
			return &ExecutionResult{
				Error:    fmt.Sprintf("Cannot open input file: %v", err),
				ExitCode: 1,
			}, err
		}
		defer file.Close()
		execCmd.Stdin = file
	}

	// Execute with sandbox if configured
	if s.config.SandboxMgr != nil {
		sandboxResult, err := s.config.SandboxMgr.Execute(execCmd)
		if err != nil {
			return &ExecutionResult{
				Error:    err.Error(),
				ExitCode: 1,
			}, err
		}
		
		return &ExecutionResult{
			Output:   sandboxResult.Output,
			Error:    sandboxResult.Error,
			ExitCode: sandboxResult.ExitCode,
		}, nil
	}

	// Regular execution
	err := execCmd.Run()
	exitCode := 0

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		} else {
			return &ExecutionResult{
				Error:    err.Error(),
				ExitCode: 1,
			}, err
		}
	}

	return &ExecutionResult{
		Output:   output.String(),
		Error:    errorOutput.String(),
		ExitCode: exitCode,
	}, nil
}

func (s *Shell) buildEnvSlice() []string {
	var env []string
	for key, value := range s.env {
		env = append(env, key+"="+value)
	}
	return env
}

// Built-in command implementations
func (s *Shell) changeDirectory(args []string) (*ExecutionResult, error) {
	start := time.Now()
	var target string
	
	if len(args) == 0 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return &ExecutionResult{Error: "Cannot determine home directory", ExitCode: 1, Duration: time.Since(start)}, err
		}
		target = homeDir
	} else {
		target = args[0]
	}

	// Handle special cases
	if target == "-" {
		// cd - functionality (go to previous directory)
		if prev, exists := s.env["OLDPWD"]; exists {
			target = prev
		} else {
			return &ExecutionResult{Error: "OLDPWD not set", ExitCode: 1}, fmt.Errorf("no previous directory")
		}
	}

	// Expand ~ to home directory
	if strings.HasPrefix(target, "~/") {
		homeDir, _ := os.UserHomeDir()
		target = filepath.Join(homeDir, target[2:])
	}

	// Make path absolute
	absPath, err := filepath.Abs(target)
	if err != nil {
		return &ExecutionResult{Error: err.Error(), ExitCode: 1}, err
	}

	// Check if directory exists
	info, err := os.Stat(absPath)
	if err != nil {
		return &ExecutionResult{Error: fmt.Sprintf("%s: %v", target, err), ExitCode: 1}, err
	}

	if !info.IsDir() {
		return &ExecutionResult{Error: fmt.Sprintf("%s: not a directory", target), ExitCode: 1}, 
			fmt.Errorf("not a directory")
	}

	// Save old directory and change
	s.env["OLDPWD"] = s.currentDir
	s.currentDir = absPath
	s.env["PWD"] = absPath

	return &ExecutionResult{Output: "", Duration: time.Since(start)}, nil
}

func (s *Shell) setVariable(args []string) (*ExecutionResult, error) {
	start := time.Now()
	if len(args) == 0 {
		return &ExecutionResult{Error: "Usage: set VARIABLE value OR set VARIABLE=value", ExitCode: 1, Duration: time.Since(start)}, 
			fmt.Errorf("insufficient arguments")
	}

	// Handle both syntaxes: set VAR value AND set VAR=value
	if len(args) == 1 && strings.Contains(args[0], "=") {
		// Handle VAR=value syntax
		parts := strings.SplitN(args[0], "=", 2)
		if len(parts) == 2 {
			s.env[parts[0]] = parts[1]
			return &ExecutionResult{Output: fmt.Sprintf("%s=%s", parts[0], parts[1]), Duration: time.Since(start)}, nil
		}
	}
	
	if len(args) < 2 {
		return &ExecutionResult{Error: "Usage: set VARIABLE value OR set VARIABLE=value", ExitCode: 1, Duration: time.Since(start)}, 
			fmt.Errorf("insufficient arguments")
	}

	variable := args[0]
	value := strings.Join(args[1:], " ")
	s.env[variable] = value

	return &ExecutionResult{Output: fmt.Sprintf("%s=%s", variable, value), Duration: time.Since(start)}, nil
}

func (s *Shell) exportVariable(args []string) (*ExecutionResult, error) {
	start := time.Now()
	if len(args) == 0 {
		// Show all exported variables
		var output strings.Builder
		for key, value := range s.env {
			output.WriteString(fmt.Sprintf("export %s=%s\n", key, value))
		}
		return &ExecutionResult{Output: output.String(), Duration: time.Since(start)}, nil
	}

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			s.env[parts[0]] = parts[1]
			os.Setenv(parts[0], parts[1])
		} else {
			// Export existing variable
			if value, exists := s.env[parts[0]]; exists {
				os.Setenv(parts[0], value)
			}
		}
	}

	return &ExecutionResult{Duration: time.Since(start)}, nil
}

func (s *Shell) createAlias(args []string) (*ExecutionResult, error) {
	start := time.Now()
	if len(args) < 2 {
		// Show all aliases
		var output strings.Builder
		for name, command := range s.aliases {
			output.WriteString(fmt.Sprintf("alias %s='%s'\n", name, command))
		}
		return &ExecutionResult{Output: output.String(), Duration: time.Since(start)}, nil
	}

	aliasName := args[0]
	aliasCommand := strings.Join(args[1:], " ")
	s.aliases[aliasName] = aliasCommand

	return &ExecutionResult{Output: fmt.Sprintf("alias %s='%s'", aliasName, aliasCommand), Duration: time.Since(start)}, nil
}

func (s *Shell) printWorkingDirectory(args []string) (*ExecutionResult, error) {
	start := time.Now()
	return &ExecutionResult{Output: s.currentDir + "\n", Duration: time.Since(start)}, nil
}

func (s *Shell) echo(args []string) (*ExecutionResult, error) {
	start := time.Now()
	output := strings.Join(args, " ") + "\n"
	return &ExecutionResult{Output: output, Duration: time.Since(start)}, nil
}

func (s *Shell) showHistory(args []string) (*ExecutionResult, error) {
	start := time.Now()
	var output strings.Builder
	
	startIdx := 0
	if len(args) > 0 {
		if n, err := strconv.Atoi(args[0]); err == nil && n > 0 {
			startIdx = len(s.history) - n
			if startIdx < 0 {
				startIdx = 0
			}
		}
	}

	for i := startIdx; i < len(s.history); i++ {
		output.WriteString(fmt.Sprintf("%4d  %s\n", i+1, s.history[i]))
	}

	return &ExecutionResult{Output: output.String(), Duration: time.Since(start)}, nil
}

func (s *Shell) exit(args []string) (*ExecutionResult, error) {
	start := time.Now()
	exitCode := 0
	if len(args) > 0 {
		if code, err := strconv.Atoi(args[0]); err == nil {
			exitCode = code
		}
	}
	
	os.Exit(exitCode)
	return &ExecutionResult{Duration: time.Since(start)}, nil
}