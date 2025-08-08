package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// Manager provides sandboxed execution capabilities
type Manager struct {
	config *Config
}

// Config holds sandbox configuration
type Config struct {
	TimeoutDuration time.Duration
	MemoryLimitMB   int64
	AllowNetworking bool
	AllowFileWrite  bool
	TempDirOnly     bool
}

// New creates a new sandbox manager
func New() (*Manager, error) {
	defaultConfig := &Config{
		TimeoutDuration: 30 * time.Second,
		MemoryLimitMB:   256,
		AllowNetworking: false,
		AllowFileWrite:  true,
		TempDirOnly:     false,
	}

	return &Manager{
		config: defaultConfig,
	}, nil
}

// ExecutionResult holds command execution results
type ExecutionResult struct {
	Output   string
	Error    string
	ExitCode int
	Duration time.Duration
}

// Execute runs a command in a sandboxed environment
func (m *Manager) Execute(cmd *exec.Cmd) (*ExecutionResult, error) {
	start := time.Now()

	// Comprehensive security validation before execution
	if err := m.validateExecution(cmd); err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("Execution validation failed: %v", err),
			ExitCode: 1,
			Duration: time.Since(start),
		}, err
	}

	// Apply command whitelist checking
	if err := m.enforceCommandWhitelist(cmd); err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("Command not whitelisted: %v", err),
			ExitCode: 1,
			Duration: time.Since(start),
		}, err
	}

	// Platform-specific sandboxing with enhanced security
	if err := m.applySandbox(cmd); err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("Failed to apply sandbox: %v", err),
			ExitCode: 1,
			Duration: time.Since(start),
		}, err
	}

	// Apply resource limits
	if err := m.applyResourceLimits(cmd); err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("Failed to apply resource limits: %v", err),
			ExitCode: 1,
			Duration: time.Since(start),
		}, err
	}

	// Set up timeout with process monitoring
	var timeoutCh <-chan time.Time
	if m.config.TimeoutDuration > 0 {
		timer := time.NewTimer(m.config.TimeoutDuration)
		defer timer.Stop()
		timeoutCh = timer.C
	}

	// Execute with monitoring
	result, err := m.executeWithMonitoring(cmd, timeoutCh, start)
	return result, err
}

// validateExecution performs comprehensive pre-execution validation
func (m *Manager) validateExecution(cmd *exec.Cmd) error {
	if cmd == nil {
		return fmt.Errorf("command is nil")
	}

	// Validate command path
	if err := m.ValidateCommand(cmd.Path, cmd.Args); err != nil {
		return fmt.Errorf("command validation failed: %v", err)
	}

	// Check environment variables for dangerous values
	for _, env := range cmd.Env {
		if strings.Contains(env, "LD_PRELOAD=") ||
		   strings.Contains(env, "LD_LIBRARY_PATH=") ||
		   strings.Contains(env, "DYLD_") {
			return fmt.Errorf("dangerous environment variable detected: %s", env)
		}
	}

	// Validate working directory
	if cmd.Dir != "" {
		if err := m.validateWorkingDirectory(cmd.Dir); err != nil {
			return fmt.Errorf("working directory validation failed: %v", err)
		}
	}

	return nil
}

// enforceCommandWhitelist ensures only whitelisted commands can be executed
func (m *Manager) enforceCommandWhitelist(cmd *exec.Cmd) error {
	// Define whitelist of safe commands
	whitelist := map[string]bool{
		// Basic file operations
		"cat":    true,
		"head":   true,
		"tail":   true,
		"grep":   true,
		"sed":    true,
		"awk":    true,
		"cut":    true,
		"sort":   true,
		"uniq":   true,
		"wc":     true,
		
		// Directory operations
		"ls":     true,
		"dir":    true, // Windows
		"pwd":    true,
		"cd":     true,
		
		// Text processing
		"echo":   true,
		"printf": true,
		
		// Network (limited)
		"curl":   true,
		"wget":   true,
		"ping":   true,
		
		// Development tools (safe ones)
		"git":    true,
		"node":   true,
		"python": true,
		"python3": true,
		"go":     true,
		"npm":    true,
		
		// System information (read-only)
		"uname":  true,
		"whoami": true,
		"id":     true,
		"date":   true,
	}

	// Extract command name from path
	cmdName := cmd.Path
	if idx := strings.LastIndex(cmdName, "/"); idx != -1 {
		cmdName = cmdName[idx+1:]
	}
	if idx := strings.LastIndex(cmdName, "\\"); idx != -1 {
		cmdName = cmdName[idx+1:]
	}

	// Remove file extensions for Windows
	if strings.HasSuffix(cmdName, ".exe") {
		cmdName = cmdName[:len(cmdName)-4]
	}

	if !whitelist[cmdName] {
		return fmt.Errorf("command '%s' is not in the whitelist", cmdName)
	}

	return nil
}

// applyResourceLimits applies system resource limits to prevent abuse
func (m *Manager) applyResourceLimits(cmd *exec.Cmd) error {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	// Apply platform-specific resource limits
	switch runtime.GOOS {
	case "linux":
		return m.applyLinuxResourceLimits(cmd)
	case "windows":
		return m.applyWindowsResourceLimits(cmd)
	case "darwin":
		return m.applyDarwinResourceLimits(cmd)
	default:
		return m.applyGenericResourceLimits(cmd)
	}
}

// executeWithMonitoring executes command with comprehensive monitoring
func (m *Manager) executeWithMonitoring(cmd *exec.Cmd, timeoutCh <-chan time.Time, start time.Time) (*ExecutionResult, error) {
	// Start the command
	if err := cmd.Start(); err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("Failed to start command: %v", err),
			ExitCode: 1,
			Duration: time.Since(start),
		}, err
	}

	// Monitor execution
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// Wait for completion or timeout
	select {
	case err := <-done:
		duration := time.Since(start)
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
					Duration: duration,
				}, err
			}
		}

		return &ExecutionResult{
			ExitCode: exitCode,
			Duration: duration,
		}, nil

	case <-timeoutCh:
		// Kill the process on timeout
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return &ExecutionResult{
			Error:    "Command execution timed out",
			ExitCode: 124, // Standard timeout exit code
			Duration: time.Since(start),
		}, fmt.Errorf("execution timeout")
	}
}

// validateWorkingDirectory validates the working directory for security
func (m *Manager) validateWorkingDirectory(dir string) error {
	// Resolve absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("cannot resolve directory path: %v", err)
	}

	// Check for dangerous system directories
	dangerousPaths := []string{
		"/",
		"/etc",
		"/usr",
		"/var",
		"/sys",
		"/proc",
		"/boot",
		"C:\\Windows",
		"C:\\Program Files",
		"C:\\System32",
	}

	for _, dangerous := range dangerousPaths {
		if absDir == dangerous || strings.HasPrefix(absDir, dangerous+string(filepath.Separator)) {
			return fmt.Errorf("working directory '%s' is in a protected system path", absDir)
		}
	}

	// Ensure directory exists and is accessible
	info, err := os.Stat(absDir)
	if err != nil {
		return fmt.Errorf("cannot access directory: %v", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", absDir)
	}

	return nil
}

// Generic cross-platform resource limits
func (m *Manager) applyGenericResourceLimits(cmd *exec.Cmd) error {
	// Generic cross-platform limits
	return nil
}

func (m *Manager) applySandbox(cmd *exec.Cmd) error {
	switch runtime.GOOS {
	case "linux":
		return m.applyLinuxSandbox(cmd)
	case "windows":
		return m.applyWindowsSandbox(cmd)
	case "darwin":
		return m.applyDarwinSandbox(cmd)
	default:
		// For other platforms, apply basic restrictions
		return m.applyBasicRestrictions(cmd)
	}
}


// Basic cross-platform restrictions as fallback
func (m *Manager) applyBasicRestrictions(cmd *exec.Cmd) error {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	// Set minimal environment
	cmd.Env = m.createMinimalEnvironment()

	return nil
}

// createMinimalEnvironment creates a minimal, secure environment
func (m *Manager) createMinimalEnvironment() []string {
	env := []string{
		"PATH=/usr/local/bin:/usr/bin:/bin",
		"HOME=/tmp/lucien_sandbox",
		"USER=sandbox",
		"SHELL=/bin/sh",
		"TERM=dumb",
		"LC_ALL=C",
		"TZ=UTC",
	}

	// Platform-specific adjustments
	switch runtime.GOOS {
	case "windows":
		env = []string{
			"PATH=C:\\Windows\\System32",
			"TEMP=C:\\Temp\\lucien_sandbox",
			"TMP=C:\\Temp\\lucien_sandbox",
			"USERNAME=sandbox",
			"COMPUTERNAME=SANDBOX",
		}
	case "darwin":
		env = append(env, "TMPDIR=/tmp/lucien_sandbox")
	}

	return env
}

// SetConfig updates the sandbox configuration
func (m *Manager) SetConfig(config *Config) {
	m.config = config
}

// GetConfig returns the current sandbox configuration
func (m *Manager) GetConfig() *Config {
	return m.config
}

// IsSupported checks if sandboxing is supported on the current platform
func IsSupported() bool {
	switch runtime.GOOS {
	case "linux", "darwin", "windows":
		return true
	default:
		return false
	}
}

// GetSandboxInfo returns information about sandbox capabilities
func GetSandboxInfo() map[string]interface{} {
	info := make(map[string]interface{})
	
	info["platform"] = runtime.GOOS
	info["supported"] = IsSupported()
	
	switch runtime.GOOS {
	case "linux":
		info["features"] = []string{
			"process_isolation",
			"resource_limits", 
			"gvisor_ready", // TODO: detect actual gVisor availability
		}
	case "windows":
		info["features"] = []string{
			"job_objects",
			"process_isolation",
		}
	case "darwin":
		info["features"] = []string{
			"process_isolation",
			"sandbox_exec", // macOS sandbox-exec support
		}
	default:
		info["features"] = []string{
			"basic_isolation",
		}
	}
	
	return info
}

// CreateSecureEnvironment sets up a minimal, secure environment
func (m *Manager) CreateSecureEnvironment() map[string]string {
	secureEnv := map[string]string{
		"PATH":     "/usr/local/bin:/usr/bin:/bin",
		"HOME":     "/tmp/lucien_sandbox",
		"USER":     "lucien_sandbox",
		"SHELL":    "/bin/sh",
		"TERM":     "dumb",
		"TMPDIR":   "/tmp",
		"LC_ALL":   "C",
	}

	// Platform-specific adjustments
	switch runtime.GOOS {
	case "windows":
		secureEnv = map[string]string{
			"PATH":     "C:\\Windows\\System32",
			"TEMP":     "C:\\Temp\\lucien_sandbox",
			"TMP":      "C:\\Temp\\lucien_sandbox",
			"USERNAME": "lucien_sandbox",
		}
	}

	return secureEnv
}

// ValidateCommand checks if a command is safe to execute in sandbox
func (m *Manager) ValidateCommand(cmdName string, args []string) error {
	// List of prohibited commands that could be dangerous or escape sandbox
	prohibited := []string{
		// System administration
		"sudo", "su", "doas",
		"chroot", "unshare",
		"docker", "podman",
		"systemctl", "service",
		"mount", "umount",
		"insmod", "rmmod", "modprobe",
		
		// File operations
		"rm", "rmdir", "del", "erase",
		"format", "fdisk", "mkfs", "dd",
		"chmod", "chown", "chgrp",
		
		// Process control
		"kill", "killall", "pkill",
		"reboot", "shutdown", "halt",
		"init", "systemd",
	}

	for _, p := range prohibited {
		if cmdName == p {
			return fmt.Errorf("command '%s' is prohibited in sandbox", cmdName)
		}
	}

	// Check for suspicious argument patterns
	for _, arg := range args {
		if containsSuspiciousPattern(arg) {
			return fmt.Errorf("argument contains suspicious pattern: %s", arg)
		}
	}

	return nil
}

func containsSuspiciousPattern(arg string) bool {
	// Convert to lowercase for case-insensitive matching
	lowerArg := strings.ToLower(arg)
	
	// Check for path traversal patterns
	pathTraversalPatterns := []string{
		"../", "..\\", // Path traversal
		"/etc/", "/proc/", "/sys/", "/dev/", "/boot/", // Unix system directories
		"c:\\windows\\", "c:\\program files\\", "c:\\system32\\", // Windows system directories
		"\\windows\\", "\\system32\\", // Windows paths without drive
	}
	
	for _, pattern := range pathTraversalPatterns {
		if strings.Contains(lowerArg, strings.ToLower(pattern)) {
			return true
		}
	}
	
	// Check for absolute paths that could be dangerous
	if strings.HasPrefix(lowerArg, "/etc/") ||
	   strings.HasPrefix(lowerArg, "/proc/") ||
	   strings.HasPrefix(lowerArg, "/sys/") ||
	   strings.HasPrefix(lowerArg, "/dev/") ||
	   strings.HasPrefix(lowerArg, "c:\\windows\\") ||
	   strings.HasPrefix(lowerArg, "c:\\program files\\") {
		return true
	}
	
	// Check for command injection patterns
	injectionPatterns := []string{
		"$(", 
		"`",
		"|sudo",
		";sudo",
		"&&sudo",
	}
	
	for _, pattern := range injectionPatterns {
		if strings.Contains(arg, pattern) {
			return true
		}
	}

	return false
}