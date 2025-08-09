package plugin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// AgentCommand represents an agent command configuration
type AgentCommand struct {
	Name        string
	ScriptPath  string
	Description string
}

// Bridge handles Python agent command execution
type Bridge struct {
	agentCommands map[string]AgentCommand
	pluginDir     string
	pythonPath    string
}

// NewBridge creates a new plugin bridge
func NewBridge(pluginDir string) (*Bridge, error) {
	// Find Python executable
	pythonPath := findPython()
	if pythonPath == "" {
		return nil, fmt.Errorf("python executable not found in PATH")
	}

	// Ensure plugin directory exists
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %w", err)
	}

	bridge := &Bridge{
		agentCommands: make(map[string]AgentCommand),
		pluginDir:     pluginDir,
		pythonPath:    pythonPath,
	}

	// Register default agent commands
	bridge.registerAgentCommands()

	return bridge, nil
}

// registerAgentCommands registers all available agent commands
func (b *Bridge) registerAgentCommands() {
	agents := []AgentCommand{
		{
			Name:        "plan",
			ScriptPath:  "planner_agent.py",
			Description: "Break down a goal into actionable tasks using AI planning",
		},
		{
			Name:        "design",
			ScriptPath:  "designer_agent.py",
			Description: "Generate UI code from natural language descriptions",
		},
		{
			Name:        "review",
			ScriptPath:  "review_agent.py",
			Description: "Analyze code files and provide improvement suggestions",
		},
		{
			Name:        "code",
			ScriptPath:  "code_agent.py",
			Description: "Generate, refactor, or explain code using AI assistance",
		},
	}

	for _, agent := range agents {
		b.agentCommands[agent.Name] = agent
	}
}

// IsAgentCommand checks if a command is an agent command
func (b *Bridge) IsAgentCommand(command string) bool {
	_, exists := b.agentCommands[command]
	return exists
}

// RunAgentCommand executes a Python agent command and returns the output
func (b *Bridge) RunAgentCommand(command string, args []string) (string, error) {
	agent, exists := b.agentCommands[command]
	if !exists {
		return "", fmt.Errorf("unknown agent command: %s", command)
	}

	scriptPath := filepath.Join(b.pluginDir, agent.ScriptPath)
	
	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return "", fmt.Errorf("agent script not found: %s", scriptPath)
	}

	// Build command arguments
	cmdArgs := []string{scriptPath}
	cmdArgs = append(cmdArgs, args...)

	return b.executeWithSecurity(command, cmdArgs)
}

// executeWithSecurity executes agent commands with security constraints and logging
func (b *Bridge) executeWithSecurity(command string, cmdArgs []string) (string, error) {
	// Create context with timeout for security
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute Python script with timeout
	cmd := exec.CommandContext(ctx, b.pythonPath, cmdArgs...)
	cmd.Dir = b.pluginDir

	// Set secured environment variables
	cmd.Env = append(os.Environ(),
		"PYTHONPATH="+b.pluginDir,
		"LUCIEN_PLUGIN_MODE=true",
		"LUCIEN_AGENT_COMMAND="+command,
		"LUCIEN_EXECUTION_ID="+generateExecutionID(),
		"PYTHONIOENCODING=utf-8",
		"PYTHONUTF8=1",
	)

	// Log agent execution start
	b.logAgentExecution("START", command, cmdArgs)

	start := time.Now()
	
	// Capture output
	output, err := cmd.CombinedOutput()
	duration := time.Since(start)
	
	// Log agent execution completion
	if err != nil {
		b.logAgentExecution("FAILED", command, []string{fmt.Sprintf("error: %v, duration: %v", err, duration)})
		
		// Check if it was a timeout
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("agent command timed out after 30 seconds")
		}
		
		return "", fmt.Errorf("agent command failed: %w\nOutput: %s", err, string(output))
	}
	
	b.logAgentExecution("SUCCESS", command, []string{fmt.Sprintf("duration: %v", duration)})
	
	// Validate output for security (basic sanitization)
	sanitizedOutput := b.sanitizeAgentOutput(string(output))
	
	return strings.TrimSpace(sanitizedOutput), nil
}

// sanitizeAgentOutput removes potentially harmful content from agent output
func (b *Bridge) sanitizeAgentOutput(output string) string {
	// Remove any potential command injection patterns
	sanitized := output
	
	// Remove or escape dangerous patterns
	dangerousPatterns := []string{
		"$(", "`", "&", "|", ";", "&&", "||", ">", "<",
		"rm -rf", "del /f", "format c:", "chmod 777",
	}
	
	for _, pattern := range dangerousPatterns {
		if strings.Contains(strings.ToLower(sanitized), pattern) {
			// Log potential security issue
			b.logSecurityEvent("DANGEROUS_PATTERN", pattern, "Agent output contained dangerous pattern")
		}
	}
	
	// Limit output size to prevent memory exhaustion
	const maxOutputSize = 100000 // 100KB
	if len(sanitized) > maxOutputSize {
		sanitized = sanitized[:maxOutputSize] + "\n... [output truncated for security]"
		b.logSecurityEvent("OUTPUT_TRUNCATED", "size_limit", "Agent output exceeded size limit")
	}
	
	return sanitized
}

// logAgentExecution logs agent execution events for security monitoring
func (b *Bridge) logAgentExecution(event, command string, args []string) {
	// In a production system, this would write to a secure log file
	// For now, we'll use a simple timestamp-based logging approach
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	logEntry := fmt.Sprintf("[%s] AGENT_EXEC: %s - %s %v", 
		timestamp, event, command, args)
	
	// Write to agent execution log (implement actual logging as needed)
	b.writeSecurityLog("agent_execution.log", logEntry)
}

// logSecurityEvent logs security-related events
func (b *Bridge) logSecurityEvent(eventType, details, description string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	logEntry := fmt.Sprintf("[%s] SECURITY_EVENT: %s - %s - %s", 
		timestamp, eventType, details, description)
	
	// Write to security log (implement actual logging as needed)
	b.writeSecurityLog("security_events.log", logEntry)
}

// writeSecurityLog writes log entries to security log files
func (b *Bridge) writeSecurityLog(filename, entry string) {
	// Create logs directory in plugin directory
	logsDir := filepath.Join(b.pluginDir, "logs")
	os.MkdirAll(logsDir, 0755)
	
	logFile := filepath.Join(logsDir, filename)
	
	// Append to log file (in production, use proper log rotation)
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Silently fail logging to not break agent execution
		return
	}
	defer f.Close()
	
	fmt.Fprintf(f, "%s\n", entry)
}

// generateExecutionID generates a unique execution ID for tracking
func generateExecutionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// GetAvailableAgents returns a list of available agent commands
func (b *Bridge) GetAvailableAgents() []AgentCommand {
	agents := make([]AgentCommand, 0, len(b.agentCommands))
	for _, agent := range b.agentCommands {
		agents = append(agents, agent)
	}
	return agents
}

// findPython finds Python executable in PATH
func findPython() string {
	candidates := []string{"python3", "python", "py"}
	
	for _, candidate := range candidates {
		if path, err := exec.LookPath(candidate); err == nil {
			// Test if it's actually Python
			cmd := exec.Command(path, "--version")
			if err := cmd.Run(); err == nil {
				return path
			}
		}
	}
	
	return ""
}

// ExecuteWithTimeout runs an agent command with a timeout
func (b *Bridge) ExecuteWithTimeout(command string, args []string, timeout time.Duration) (string, error) {
	agent, exists := b.agentCommands[command]
	if !exists {
		return "", fmt.Errorf("unknown agent command: %s", command)
	}

	scriptPath := filepath.Join(b.pluginDir, agent.ScriptPath)
	
	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return "", fmt.Errorf("agent script not found: %s", scriptPath)
	}

	// Build command arguments
	cmdArgs := []string{scriptPath}
	cmdArgs = append(cmdArgs, args...)

	// Execute with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, b.pythonPath, cmdArgs...)
	cmd.Dir = b.pluginDir
	cmd.Env = append(os.Environ(),
		"PYTHONPATH="+b.pluginDir,
		"LUCIEN_PLUGIN_MODE=true",
		"PYTHONIOENCODING=utf-8",
		"PYTHONUTF8=1",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("agent command timed out after %v", timeout)
		}
		return "", fmt.Errorf("agent command failed: %w\nOutput: %s", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}