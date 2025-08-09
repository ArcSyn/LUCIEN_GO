package shell

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
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
	SafeMode       bool
	ExecutorMode   string // "shell" or "internal"
	DisableHistory bool   // Skip loading history file (for tests)
}

// Command represents a parsed command with its arguments and redirections
type Command struct {
	Name      string
	Args      []string
	Input     interface{}
	Output    interface{}
	Redirects map[string]string
}

// Shell represents the core shell engine
type Shell struct {
	config        *Config
	aliases       map[string]string
	currentDir    string
	builtins      map[string]func([]string) (*ExecutionResult, error)
	history       []string
	jobs          map[int]*Job
	nextJobID     int
	variables     map[string]string
	historyFile   string
	securityGuard *SecurityGuard
	dispatcher    MessageDispatcher // For Bubble Tea integration
}

// Job represents a background job
type Job struct {
	ID      int
	Command string
	Status  string
	PID     int
}

// New creates a new shell instance
func New(config *Config) *Shell {
	if config == nil {
		config = &Config{
			SafeMode:     false,
			ExecutorMode: "shell",
		}
	}
	
	// Default executor mode if not set
	if config.ExecutorMode == "" {
		config.ExecutorMode = "shell"
	}

	s := &Shell{
		config:        config,
		aliases:       make(map[string]string),
		currentDir:    getCurrentDir(),
		history:       make([]string, 0),
		jobs:          make(map[int]*Job, 0),
		nextJobID:     1,
		variables:     make(map[string]string),
		securityGuard: NewSecurityGuard(),
	}

	s.historyFile = s.getHistoryFile()
	if !config.DisableHistory {
		s.loadHistory()
	}

	s.initBuiltins()
	s.initEnvironment()
	
	return s
}

// SetDispatcher sets the message dispatcher for Bubble Tea integration
func (s *Shell) SetDispatcher(dispatcher MessageDispatcher) {
	s.dispatcher = dispatcher
}

// getHistoryFile returns the path to the history file
func (s *Shell) getHistoryFile() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".lucien_history"
	}
	
	lucienDir := filepath.Join(homeDir, ".lucien")
	os.MkdirAll(lucienDir, 0755)
	
	return filepath.Join(lucienDir, "history")
}

// loadHistory loads command history from file
func (s *Shell) loadHistory() {
	data, err := os.ReadFile(s.historyFile)
	if err != nil {
		return // File doesn't exist or can't read
	}
	
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			s.history = append(s.history, line)
		}
	}
}

// saveHistory saves command history to file
func (s *Shell) saveHistory() {
	if len(s.history) == 0 {
		return
	}
	
	// Keep only last 1000 entries
	start := 0
	if len(s.history) > 1000 {
		start = len(s.history) - 1000
	}
	
	content := strings.Join(s.history[start:], "\n")
	os.WriteFile(s.historyFile, []byte(content), 0644)
}

// Execute runs a command line and returns the result
func (s *Shell) Execute(commandLine string) (*ExecutionResult, error) {
	start := time.Now()
	
	// Expand variables and tilde
	commandLine = s.expandVariables(commandLine)
	
	// Add to history (don't add duplicates)
	if len(s.history) == 0 || s.history[len(s.history)-1] != commandLine {
		s.history = append(s.history, commandLine)
		s.saveHistory()
	}
	
	// Expand history if needed
	expandedLine, err := s.expandHistory(commandLine)
	if err != nil {
		return &ExecutionResult{
			Error:    err.Error(),
			ExitCode: 1,
			Duration: time.Since(start),
		}, nil  // Return nil error - expansion failure is handled in ExecutionResult
	}
	
	// Parse command line with advanced parser (includes security validation)
	commandChain, err := s.parseCommandLineAdvanced(expandedLine)
	if err != nil {
		return &ExecutionResult{
			Error:    err.Error(),
			ExitCode: 1,
			Duration: time.Since(start),
		}, err
	}
	
	// Execute command chain
	result, err := s.executeCommandChain(*commandChain)
	if result != nil {
		result.Duration = time.Since(start)
	}
	
	return result, err
}

// expandVariables expands environment variables and tilde in the input
func (s *Shell) expandVariables(input string) string {
	// Handle tilde expansion first
	if strings.HasPrefix(input, "~") {
		if homeDir, err := os.UserHomeDir(); err == nil {
			if input == "~" {
				return homeDir
			} else if strings.HasPrefix(input, "~/") {
				return filepath.Join(homeDir, input[2:])
			}
		}
	}
	
	result := input
	
	// Create a combined variable map (shell variables + environment)
	allVars := make(map[string]string)
	
	// Add shell variables
	for k, v := range s.variables {
		allVars[k] = v
	}
	
	// Add common environment variables
	envVars := map[string]string{
		"HOME":        os.Getenv("HOME"),
		"USER":        os.Getenv("USER"),
		"USERPROFILE": os.Getenv("USERPROFILE"),
		"USERNAME":    os.Getenv("USERNAME"),
		"PATH":        os.Getenv("PATH"),
		"PWD":         s.currentDir,
	}
	
	for k, v := range envVars {
		if v != "" {
			allVars[k] = v
		}
	}
	
	// Expand ${VAR} format first (more specific)
	braceRegex := regexp.MustCompile(`\${([A-Za-z_][A-Za-z0-9_]*)}`)
	result = braceRegex.ReplaceAllStringFunc(result, func(match string) string {
		varName := match[2 : len(match)-1] // Remove ${ and }
		if value, exists := allVars[varName]; exists {
			return value
		}
		return ""
	})
	
	// Expand $VAR format (word boundary aware)
	dollarRegex := regexp.MustCompile(`\$([A-Za-z_][A-Za-z0-9_]*)`)
	result = dollarRegex.ReplaceAllStringFunc(result, func(match string) string {
		varName := match[1:] // Remove $
		if value, exists := allVars[varName]; exists {
			return value
		}
		return ""
	})
	
	// Windows style %VAR%
	percentRegex := regexp.MustCompile(`%([A-Za-z_][A-Za-z0-9_]*)%`)
	result = percentRegex.ReplaceAllStringFunc(result, func(match string) string {
		varName := match[1 : len(match)-1] // Remove % and %
		if value, exists := allVars[varName]; exists {
			return value
		}
		return ""
	})
	
	return result
}

// getCurrentDir gets the current working directory
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "/"
	}
	return dir
}

// initBuiltins initializes built-in commands
func (s *Shell) initBuiltins() {
	s.builtins = map[string]func([]string) (*ExecutionResult, error){
		"cd":      s.changeDirectory,
		"pwd":     s.builtinPwd,
		"echo":    s.builtinEcho,
		"set":     s.builtinSet,
		"unset":   s.builtinUnset,
		"export":  s.builtinExport,
		"alias":   s.builtinAlias,
		"unalias": s.builtinUnalias,
		"history": s.builtinHistory,
		"jobs":    s.builtinJobs,
		"fg":      s.builtinFg,
		"bg":      s.builtinBg,
		"disown":  s.builtinDisown,
		"suspend": s.builtinSuspend,
		"kill":    s.builtinKill,
		"exit":    s.builtinExit,
		"help":    s.builtinHelp,
		"clear":   s.builtinClear,
		"env":     s.builtinEnv,
		"home":    s.builtinHome,
	}
	
	// Add security toggle command
	s.builtins[":secure"] = s.builtinSecure
}

// initEnvironment initializes environment variables
func (s *Shell) initEnvironment() {
	// Copy environment variables
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			s.variables[parts[0]] = parts[1]
		}
	}
}

// expandHistory expands history references like !! and !n
func (s *Shell) expandHistory(input string) (string, error) {
	if !strings.Contains(input, "!") {
		return input, nil
	}
	
	// Simple implementation for !!, !n, and !prefix
	if input == "!!" {
		if len(s.history) < 2 {
			return "", fmt.Errorf("no previous command in history")
		}
		return s.history[len(s.history)-2], nil
	}
	
	// Handle !n (number)
	if strings.HasPrefix(input, "!") && len(input) > 1 {
		rest := input[1:]
		if num, err := strconv.Atoi(rest); err == nil {
			if num < 1 || num > len(s.history) {
				return "", fmt.Errorf("history entry %d not found", num)
			}
			return s.history[num-1], nil
		}
		
		// Handle !prefix
		for i := len(s.history) - 2; i >= 0; i-- {
			if strings.HasPrefix(s.history[i], rest) {
				return s.history[i], nil
			}
		}
		return "", fmt.Errorf("no history entry starting with '%s'", rest)
	}
	
	return input, nil
}

// executeCommandChain executes a chain of commands with proper operator handling
func (s *Shell) executeCommandChain(chain CommandChain) (*ExecutionResult, error) {
	if len(chain.Commands) == 0 {
		return &ExecutionResult{ExitCode: 0}, nil
	}

	// Execute first command
	result, err := s.executeSingleCommand(chain.Commands[0])
	if err != nil {
		return result, err
	}

	// Process operators
	for i := 0; i < len(chain.Operators) && i+1 < len(chain.Commands); i++ {
		opType := chain.Types[i]
		nextCmd := chain.Commands[i+1]

		switch opType {
		case CommandAnd: // &&
			if result.ExitCode == 0 {
				nextResult, err := s.executeSingleCommand(nextCmd)
				if err != nil {
					return nextResult, err
				}
				result = s.combineResults(result, nextResult)
			}

		case CommandOr: // ||
			if result.ExitCode != 0 {
				nextResult, err := s.executeSingleCommand(nextCmd)
				if err != nil {
					return nextResult, err
				}
				result = s.combineResults(result, nextResult)
			}

		case CommandSequence: // ;
			nextResult, err := s.executeSingleCommand(nextCmd)
			if err != nil {
				return nextResult, err
			}
			result = s.combineResults(result, nextResult)

		case CommandPipe: // |
			// Simple pipe implementation - pass output as input
			nextCmd.Input = result.Output
			nextResult, err := s.executeSingleCommand(nextCmd)
			if err != nil {
				return nextResult, err
			}
			result = nextResult // Replace result with piped output

		case CommandBackground: // &
			// Start command in background
			go func(cmd Command) {
				s.executeSingleCommand(cmd)
			}(nextCmd)
			// Don't wait for background command
		}
	}

	return result, nil
}

// combineResults combines two execution results
func (s *Shell) combineResults(r1, r2 *ExecutionResult) *ExecutionResult {
	combined := &ExecutionResult{
		ExitCode: r2.ExitCode,
		Duration: r1.Duration + r2.Duration,
	}

	if r1.Output != "" && r2.Output != "" {
		combined.Output = r1.Output + "\n" + r2.Output
	} else if r1.Output != "" {
		combined.Output = r1.Output
	} else {
		combined.Output = r2.Output
	}

	if r1.Error != "" && r2.Error != "" {
		combined.Error = r1.Error + "\n" + r2.Error
	} else if r1.Error != "" {
		combined.Error = r1.Error
	} else {
		combined.Error = r2.Error
	}

	return combined
}


// executeSingleCommand executes a single command
func (s *Shell) executeSingleCommand(cmd Command) (*ExecutionResult, error) {
	// Expand aliases
	if alias, exists := s.aliases[cmd.Name]; exists {
		// Simple alias expansion
		parts := strings.Fields(alias)
		if len(parts) > 0 {
			cmd.Name = parts[0]
			cmd.Args = append(parts[1:], cmd.Args...)
		}
	}
	
	// Check built-ins first
	if builtin, exists := s.builtins[cmd.Name]; exists {
		return builtin(cmd.Args)
	}
	
	// Execute external command based on executor mode
	if s.config.ExecutorMode == "shell" {
		return s.executeExternalViaShell(&cmd)
	} else {
		return s.executeExternalPlatform(&cmd)
	}
}

// parseJobID parses job reference like %1, %+, %-
func (s *Shell) parseJobID(arg string) (int, error) {
	if !strings.HasPrefix(arg, "%") {
		return 0, fmt.Errorf("not a job reference")
	}
	
	ref := arg[1:]
	switch ref {
	case "+":
		// Current job (most recent)
		if len(s.jobs) == 0 {
			return 0, fmt.Errorf("no current job")
		}
		maxID := 0
		for id := range s.jobs {
			if id > maxID {
				maxID = id
			}
		}
		return maxID, nil
		
	case "-":
		// Previous job (second most recent)
		if len(s.jobs) < 2 {
			return 0, fmt.Errorf("no previous job")
		}
		// Find second highest ID
		var ids []int
		for id := range s.jobs {
			ids = append(ids, id)
		}
		if len(ids) < 2 {
			return 0, fmt.Errorf("no previous job")
		}
		// Sort and get second highest
		for i := 0; i < len(ids)-1; i++ {
			for j := i + 1; j < len(ids); j++ {
				if ids[i] < ids[j] {
					ids[i], ids[j] = ids[j], ids[i]
				}
			}
		}
		return ids[1], nil
		
	default:
		// Job number
		if jobID, err := strconv.Atoi(ref); err == nil {
			if _, exists := s.jobs[jobID]; exists {
				return jobID, nil
			}
			return 0, fmt.Errorf("no such job: %d", jobID)
		}
		
		// Job name prefix
		for id, job := range s.jobs {
			if strings.HasPrefix(job.Command, ref) {
				return id, nil
			}
		}
		return 0, fmt.Errorf("no job matching: %s", ref)
	}
}

// buildEnvSlice builds environment slice for exec
func (s *Shell) buildEnvSlice() []string {
	var env []string
	for key, value := range s.variables {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
}

// Built-in command implementations

func (s *Shell) builtinCd(args []string) (*ExecutionResult, error) {
	var targetDir string
	
	if len(args) == 0 {
		// Change to home directory
		home, err := os.UserHomeDir()
		if err != nil {
			return &ExecutionResult{
				Error:    "Cannot determine home directory",
				ExitCode: 1,
			}, nil
		}
		targetDir = home
	} else {
		targetDir = args[0]
	}
	
	// Expand ~ if present
	if strings.HasPrefix(targetDir, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return &ExecutionResult{
				Error:    "Cannot expand ~: " + err.Error(),
				ExitCode: 1,
			}, nil
		}
		targetDir = filepath.Join(home, targetDir[1:])
	}
	
	// Make it absolute
	if !filepath.IsAbs(targetDir) {
		targetDir = filepath.Join(s.currentDir, targetDir)
	}
	
	// Clean the path
	targetDir = filepath.Clean(targetDir)
	
	// Check if directory exists
	if info, err := os.Stat(targetDir); err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("cd: %s: No such file or directory", targetDir),
			ExitCode: 1,
		}, nil
	} else if !info.IsDir() {
		return &ExecutionResult{
			Error:    fmt.Sprintf("cd: %s: Not a directory", targetDir),
			ExitCode: 1,
		}, nil
	}
	
	// Change directory
	if err := os.Chdir(targetDir); err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("cd: %s", err.Error()),
			ExitCode: 1,
		}, nil
	}
	
	s.currentDir = targetDir
	
	return &ExecutionResult{ExitCode: 0}, nil
}

func (s *Shell) builtinPwd(args []string) (*ExecutionResult, error) {
	return &ExecutionResult{
		Output:   s.currentDir + "\n",
		ExitCode: 0,
	}, nil
}

func (s *Shell) builtinEcho(args []string) (*ExecutionResult, error) {
	output := strings.Join(args, " ") + "\n"
	return &ExecutionResult{
		Output:   output,
		ExitCode: 0,
	}, nil
}

func (s *Shell) builtinSet(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		// List all variables
		var output []string
		for key, value := range s.variables {
			output = append(output, fmt.Sprintf("%s=%s", key, value))
		}
		return &ExecutionResult{
			Output:   strings.Join(output, "\n") + "\n",
			ExitCode: 0,
		}, nil
	}
	
	// Handle traditional "set VAR value" format
	if len(args) == 2 && !strings.Contains(args[0], "=") {
		s.variables[args[0]] = args[1]
		return &ExecutionResult{
			Output:   fmt.Sprintf("%s=%s\n", args[0], args[1]),
			ExitCode: 0,
		}, nil
	}
	
	// Handle VAR=value format
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			s.variables[parts[0]] = parts[1]
		} else {
			return &ExecutionResult{
				Error:    fmt.Sprintf("set: invalid assignment: %s", arg),
				ExitCode: 1,
			}, nil
		}
	}
	
	return &ExecutionResult{ExitCode: 0}, nil
}

func (s *Shell) builtinUnset(args []string) (*ExecutionResult, error) {
	for _, arg := range args {
		delete(s.variables, arg)
	}
	return &ExecutionResult{ExitCode: 0}, nil
}

func (s *Shell) builtinAlias(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		// List all aliases
		var output []string
		for alias, command := range s.aliases {
			output = append(output, fmt.Sprintf("alias %s='%s'", alias, command))
		}
		return &ExecutionResult{
			Output:   strings.Join(output, "\n") + "\n",
			ExitCode: 0,
		}, nil
	}
	
	// Handle "alias name value" format
	if len(args) == 2 && !strings.Contains(args[0], "=") {
		name := args[0]
		value := args[1]
		s.aliases[name] = value
		return &ExecutionResult{
			Output:   fmt.Sprintf("alias %s='%s'\n", name, value),
			ExitCode: 0,
		}, nil
	}
	
	// Handle single argument (show alias or error)
	if len(args) == 1 && !strings.Contains(args[0], "=") {
		name := args[0]
		if command, exists := s.aliases[name]; exists {
			return &ExecutionResult{
				Output:   fmt.Sprintf("alias %s='%s'\n", name, command),
				ExitCode: 0,
			}, nil
		} else {
			return &ExecutionResult{
				Error:    fmt.Sprintf("alias: %s: not found", name),
				ExitCode: 1,
			}, nil
		}
	}
	
	// Handle NAME=value format  
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			name := parts[0]
			value := parts[1]
			
			// Remove quotes if present
			value = strings.Trim(value, "'\"")
			
			s.aliases[name] = value
		} else {
			return &ExecutionResult{
				Error:    fmt.Sprintf("alias: invalid format: %s", arg),
				ExitCode: 1,
			}, nil
		}
	}
	
	return &ExecutionResult{ExitCode: 0}, nil
}

func (s *Shell) builtinUnalias(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		return &ExecutionResult{
			Error:    "unalias: usage: unalias name [name ...]",
			ExitCode: 1,
		}, nil
	}
	
	for _, name := range args {
		if _, exists := s.aliases[name]; exists {
			delete(s.aliases, name)
			// Return success message
			return &ExecutionResult{
				Output:   fmt.Sprintf("Removed alias: %s\n", name),
				ExitCode: 0,
			}, nil
		} else {
			// Check for typos
			var suggestion string
			for alias := range s.aliases {
				if levenshteinDistance(name, alias) <= 1 {
					suggestion = alias
					break
				}
			}
			
			if suggestion != "" {
				return &ExecutionResult{
					Error:    fmt.Sprintf("unalias: %s: not found. Did you mean '%s'?", name, suggestion),
					ExitCode: 1,
				}, nil
			} else {
				return &ExecutionResult{
					Error:    fmt.Sprintf("unalias: %s: not found", name),
					ExitCode: 1,
				}, nil
			}
		}
	}
	
	return &ExecutionResult{ExitCode: 0}, nil
}

func (s *Shell) builtinHistory(args []string) (*ExecutionResult, error) {
	var output []string
	for i, cmd := range s.history {
		output = append(output, fmt.Sprintf("%4d  %s", i+1, cmd))
	}
	return &ExecutionResult{
		Output:   strings.Join(output, "\n") + "\n",
		ExitCode: 0,
	}, nil
}

func (s *Shell) builtinJobs(args []string) (*ExecutionResult, error) {
	if len(s.jobs) == 0 {
		return &ExecutionResult{
			Output:   "No active jobs\n",
			ExitCode: 0,
		}, nil
	}
	
	var output []string
	for id, job := range s.jobs {
		output = append(output, fmt.Sprintf("[%d] %s %s", id, job.Status, job.Command))
	}
	
	return &ExecutionResult{
		Output:   strings.Join(output, "\n") + "\n",
		ExitCode: 0,
	}, nil
}

func (s *Shell) builtinFg(args []string) (*ExecutionResult, error) {
	var jobID int
	var err error
	
	if len(args) == 0 {
		// Use most recent job
		if len(s.jobs) == 0 {
			return &ExecutionResult{
				Error:    "fg: no current job",
				ExitCode: 1,
			}, nil
		}
		jobID = s.nextJobID - 1
	} else {
		jobID, err = s.parseJobID(args[0])
		if err != nil {
			return &ExecutionResult{
				Error:    "fg: " + err.Error(),
				ExitCode: 1,
			}, nil
		}
	}
	
	job, exists := s.jobs[jobID]
	if !exists {
		return &ExecutionResult{
			Error:    fmt.Sprintf("fg: job %d not found", jobID),
			ExitCode: 1,
		}, nil
	}
	
	return &ExecutionResult{
		Output:   fmt.Sprintf("Bringing job %d to foreground: %s\n", jobID, job.Command),
		ExitCode: 0,
	}, nil
}

func (s *Shell) builtinBg(args []string) (*ExecutionResult, error) {
	var jobID int
	var err error
	
	if len(args) == 0 {
		// Use most recent job
		if len(s.jobs) == 0 {
			return &ExecutionResult{
				Error:    "bg: no current job",
				ExitCode: 1,
			}, nil
		}
		jobID = s.nextJobID - 1
	} else {
		jobID, err = s.parseJobID(args[0])
		if err != nil {
			return &ExecutionResult{
				Error:    "bg: " + err.Error(),
				ExitCode: 1,
			}, nil
		}
	}
	
	job, exists := s.jobs[jobID]
	if !exists {
		return &ExecutionResult{
			Error:    fmt.Sprintf("bg: job %d not found", jobID),
			ExitCode: 1,
		}, nil
	}
	
	return &ExecutionResult{
		Output:   fmt.Sprintf("Sending job %d to background: %s\n", jobID, job.Command),
		ExitCode: 0,
	}, nil
}

func (s *Shell) builtinKill(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		return &ExecutionResult{
			Error:    "kill: usage: kill [pid | %job]",
			ExitCode: 1,
		}, nil
	}
	
	for _, arg := range args {
		if strings.HasPrefix(arg, "%") {
			jobID, err := s.parseJobID(arg)
			if err != nil {
				return &ExecutionResult{
					Error:    "kill: " + err.Error(),
					ExitCode: 1,
				}, nil
			}
			
			job, exists := s.jobs[jobID]
			if !exists {
				return &ExecutionResult{
					Error:    fmt.Sprintf("kill: job %d not found", jobID),
					ExitCode: 1,
				}, nil
			}
			
			// Remove from jobs (simulate kill)
			delete(s.jobs, jobID)
			
			return &ExecutionResult{
				Output:   fmt.Sprintf("Killed job %d: %s\n", jobID, job.Command),
				ExitCode: 0,
			}, nil
		} else {
			// Handle PID (not implemented for safety)
			return &ExecutionResult{
				Error:    "kill: PID killing not implemented in safe mode",
				ExitCode: 1,
			}, nil
		}
	}
	
	return &ExecutionResult{ExitCode: 0}, nil
}

func (s *Shell) builtinExit(args []string) (*ExecutionResult, error) {
	exitCode := 0
	if len(args) > 0 {
		if code, err := strconv.Atoi(args[0]); err == nil {
			exitCode = code
		}
	}
	
	// In a real implementation, this would exit the program
	return &ExecutionResult{
		Output:   "Goodbye!\n",
		ExitCode: exitCode,
	}, nil
}

func (s *Shell) builtinHelp(args []string) (*ExecutionResult, error) {
	help := `Lucien Shell Built-in Commands:
  cd [dir]        Change directory
  pwd             Print working directory  
  echo [args]     Display arguments
  set [var=val]   Set environment variable
  unset [var]     Unset environment variable
  alias [name=cmd] Create command alias
  unalias [name]  Remove alias
  history         Show command history
  jobs            List active jobs
  fg [%job]       Bring job to foreground
  bg [%job]       Send job to background  
  kill [%job|pid] Terminate job or process
  exit [code]     Exit shell
  help            Show this help

Job Control:
  %1, %2, ...     Job by number
  %+              Current job
  %-              Previous job
  %name           Job by command name prefix

History:
  !!              Last command
  !n              Command number n
  !prefix         Last command starting with prefix
`
	
	return &ExecutionResult{
		Output:   help,
		ExitCode: 0,
	}, nil
}

func (s *Shell) builtinExport(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		// List all variables
		var output []string
		for key, value := range s.variables {
			output = append(output, fmt.Sprintf("export %s=%s", key, value))
		}
		return &ExecutionResult{
			Output:   strings.Join(output, "\n") + "\n",
			ExitCode: 0,
		}, nil
	}
	
	// Set and export variables
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			s.variables[parts[0]] = parts[1]
			os.Setenv(parts[0], parts[1])
		} else {
			return &ExecutionResult{
				Error:    fmt.Sprintf("export: invalid assignment: %s", arg),
				ExitCode: 1,
			}, nil
		}
	}
	
	return &ExecutionResult{ExitCode: 0}, nil
}

func (s *Shell) builtinEnv(args []string) (*ExecutionResult, error) {
	var output []string
	
	if len(args) == 0 {
		// Show all environment variables
		for key, value := range s.variables {
			output = append(output, fmt.Sprintf("%s=%s", key, value))
		}
		
		// Add system environment variables
		for _, env := range os.Environ() {
			if !strings.Contains(env, "=") {
				continue
			}
			parts := strings.SplitN(env, "=", 2)
			if _, exists := s.variables[parts[0]]; !exists {
				output = append(output, env)
			}
		}
	}
	
	return &ExecutionResult{
		Output:   strings.Join(output, "\n") + "\n",
		ExitCode: 0,
	}, nil
}

func (s *Shell) builtinClear(args []string) (*ExecutionResult, error) {
	// ANSI escape codes to clear screen
	return &ExecutionResult{
		Output:   "\033[2J\033[H",
		ExitCode: 0,
	}, nil
}

func (s *Shell) builtinHome(args []string) (*ExecutionResult, error) {
	var homeDir string
	var err error
	
	// Platform-specific home directory detection
	if runtime.GOOS == "windows" {
		homeDir = os.Getenv("USERPROFILE")
		if homeDir == "" {
			homeDir = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		}
	} else {
		homeDir = os.Getenv("HOME")
	}
	
	if homeDir == "" {
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return &ExecutionResult{
				Error:    "Cannot determine home directory",
				ExitCode: 1,
			}, nil
		}
	}
	
	// Change to home directory
	if err := os.Chdir(homeDir); err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("home: %s", err.Error()),
			ExitCode: 1,
		}, nil
	}
	
	s.currentDir = homeDir
	
	return &ExecutionResult{
		Output:   homeDir + "\n",
		ExitCode: 0,
	}, nil
}

func (s *Shell) builtinSecure(args []string) (*ExecutionResult, error) {
	if len(args) == 0 {
		// Show current mode
		mode := "permissive"
		if s.securityGuard.GetMode() == SecurityModeStrict {
			mode = "strict"
		}
		return &ExecutionResult{
			Output:   fmt.Sprintf("Security mode: %s\n", mode),
			ExitCode: 0,
		}, nil
	}
	
	switch strings.ToLower(args[0]) {
	case "strict":
		s.securityGuard.SetMode(SecurityModeStrict)
		return &ExecutionResult{
			Output:   "Security mode set to strict\n",
			ExitCode: 0,
		}, nil
	case "permissive":
		s.securityGuard.SetMode(SecurityModePermissive)
		return &ExecutionResult{
			Output:   "Security mode set to permissive\n",
			ExitCode: 0,
		}, nil
	default:
		return &ExecutionResult{
			Error:    "Usage: :secure [strict|permissive]\n",
			ExitCode: 1,
		}, nil
	}
}

// levenshteinDistance calculates edit distance for typo suggestions
func levenshteinDistance(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	
	matrix := make([][]int, len(a)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(b)+1)
		matrix[i][0] = i
	}
	
	for j := 0; j <= len(b); j++ {
		matrix[0][j] = j
	}
	
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}
	
	return matrix[len(a)][len(b)]
}

// ProcessMessage represents messages for Bubble Tea integration
type ProcessMessage interface {
	ProcessMessage()
}

// ProcessStartedMsg indicates a process has started
type ProcessStartedMsg struct {
	Cmd string
	PID int
	Err error
}

func (ProcessStartedMsg) ProcessMessage() {}

// ProcessStdoutMsg contains stdout output from a process
type ProcessStdoutMsg struct {
	Line string
}

func (ProcessStdoutMsg) ProcessMessage() {}

// ProcessStderrMsg contains stderr output from a process
type ProcessStderrMsg struct {
	Line string
}

func (ProcessStderrMsg) ProcessMessage() {}

// ProcessExitedMsg indicates a process has exited
type ProcessExitedMsg struct {
	Code int
	Err  error
}

func (ProcessExitedMsg) ProcessMessage() {}

// MessageDispatcher is a function type for dispatching messages to Bubble Tea
type MessageDispatcher func(ProcessMessage)

// hasExeInPath checks if an executable exists in PATH
func hasExeInPath(exeName string) bool {
	_, err := exec.LookPath(exeName)
	return err == nil
}

// shellCommand creates the appropriate shell command for the platform
func shellCommand(raw string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		// Prefer PowerShell 7 (pwsh) for best cross-platform experience
		if hasExeInPath("pwsh.exe") || hasExeInPath("pwsh") {
			return exec.Command("pwsh", "-NoLogo", "-NoProfile", "-NonInteractive", "-Command", raw)
		}
		// Fallback to Windows PowerShell 5.x
		if hasExeInPath("powershell.exe") || hasExeInPath("powershell") {
			return exec.Command("powershell", "-NoLogo", "-NoProfile", "-NonInteractive", "-Command", raw)
		}
		// Final fallback to cmd.exe
		return exec.Command("cmd", "/C", raw)
	}
	return exec.Command("/bin/sh", "-c", raw)
}

// exitCodeFromError extracts the exit code from an exec error
func exitCodeFromError(err error) int {
	if err == nil {
		return 0
	}
	
	if exitError, ok := err.(*exec.ExitError); ok {
		if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
		return 1
	}
	return 1
}

// streamOutput streams data from a pipe to the dispatcher
func streamOutput(pipe io.ReadCloser, dispatch func(string)) {
	if pipe == nil || dispatch == nil {
		return
	}
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		dispatch(scanner.Text())
	}
}

// runViaSystemShell executes a command via the system shell with Bubble Tea integration
func (s *Shell) runViaSystemShell(ctx context.Context, raw string, env []string, dispatch MessageDispatcher) (int, error) {
	cmd := shellCommand(raw)
	cmd.Dir = s.currentDir
	cmd.Env = env

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		if dispatch != nil {
			dispatch(ProcessStartedMsg{Cmd: raw, Err: err})
		}
		return 1, err
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		if dispatch != nil {
			dispatch(ProcessStartedMsg{Cmd: raw, Err: err})
		}
		return 1, err
	}

	if err := cmd.Start(); err != nil {
		if dispatch != nil {
			dispatch(ProcessStartedMsg{Cmd: raw, Err: err})
		}
		return 1, err
	}

	if dispatch != nil {
		dispatch(ProcessStartedMsg{Cmd: raw, PID: cmd.Process.Pid})
	}

	// Stream stdout in goroutine
	go streamOutput(stdout, func(line string) {
		if dispatch != nil {
			dispatch(ProcessStdoutMsg{Line: line})
		}
	})

	// Stream stderr in goroutine
	go streamOutput(stderr, func(line string) {
		if dispatch != nil {
			dispatch(ProcessStderrMsg{Line: line})
		}
	})

	err = cmd.Wait()
	exitCode := exitCodeFromError(err)

	if dispatch != nil {
		dispatch(ProcessExitedMsg{Code: exitCode, Err: err})
	}

	return exitCode, err
}

// executeExternalViaShell executes commands via the system shell (new default behavior)
func (s *Shell) executeExternalViaShell(cmd *Command) (*ExecutionResult, error) {
	// Reconstruct the original command line
	cmdLine := cmd.Name
	if len(cmd.Args) > 0 {
		cmdLine += " " + strings.Join(cmd.Args, " ")
	}
	
	var outputBuilder strings.Builder
	var errorBuilder strings.Builder
	
	// Create a simple dispatcher that captures output
	dispatch := func(msg ProcessMessage) {
		switch m := msg.(type) {
		case ProcessStdoutMsg:
			outputBuilder.WriteString(m.Line + "\n")
		case ProcessStderrMsg:
			errorBuilder.WriteString(m.Line + "\n")
		}
		
		// Also forward to the UI dispatcher if available
		if s.dispatcher != nil {
			s.dispatcher(msg)
		}
	}
	
	// Run via system shell
	exitCode, err := s.runViaSystemShell(context.Background(), cmdLine, s.buildEnvSlice(), dispatch)
	
	return &ExecutionResult{
		Output:   outputBuilder.String(),
		Error:    errorBuilder.String(),
		ExitCode: exitCode,
	}, err
}

// builtinDisown removes jobs from the jobs table
func (s *Shell) builtinDisown(args []string) (*ExecutionResult, error) {
	if len(s.jobs) == 0 {
		return &ExecutionResult{
			Error:    "disown: no jobs to disown",
			ExitCode: 1,
		}, nil
	}
	
	if len(args) == 0 {
		// Disown the most recent job
		for id := range s.jobs {
			delete(s.jobs, id)
			return &ExecutionResult{
				Output:   fmt.Sprintf("Job [%d] disowned\n", id),
				ExitCode: 0,
			}, nil
		}
	} else {
		// Disown specific job
		jobID, err := s.parseJobID(args[0])
		if err != nil {
			return &ExecutionResult{
				Error:    "disown: " + err.Error(),
				ExitCode: 1,
			}, nil
		}
		
		if _, exists := s.jobs[jobID]; !exists {
			return &ExecutionResult{
				Error:    fmt.Sprintf("disown: job %d not found", jobID),
				ExitCode: 1,
			}, nil
		}
		
		delete(s.jobs, jobID)
		return &ExecutionResult{
			Output:   fmt.Sprintf("Job [%d] disowned\n", jobID),
			ExitCode: 0,
		}, nil
	}
	
	return &ExecutionResult{ExitCode: 0}, nil
}

// builtinSuspend suspends the current shell (simulated for testing)
func (s *Shell) builtinSuspend(args []string) (*ExecutionResult, error) {
	// In a real shell, this would send SIGTSTP to the current process
	// For our shell, we'll just acknowledge the command
	return &ExecutionResult{
		Output:   "Shell suspend requested (simulated)\n",
		ExitCode: 0,
	}, nil
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}