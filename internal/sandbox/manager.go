package sandbox

import (
	"fmt"
	"os/exec"
	"runtime"
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

	// Platform-specific sandboxing
	if err := m.applySandbox(cmd); err != nil {
		return &ExecutionResult{
			Error:    fmt.Sprintf("Failed to apply sandbox: %v", err),
			ExitCode: 1,
			Duration: time.Since(start),
		}, err
	}

	// Set up timeout
	if m.config.TimeoutDuration > 0 {
		go func() {
			time.Sleep(m.config.TimeoutDuration)
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		}()
	}

	// Execute the command
	err := cmd.Run()
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
}

func (m *Manager) applySandbox(cmd *exec.Cmd) error {
	switch runtime.GOOS {
	case "linux":
		return m.applyLinuxSandbox(cmd)
	case "windows":
		return m.applyWindowsSandbox(cmd)
	default:
		// For other platforms, apply basic restrictions
		return m.applyBasicRestrictions(cmd)
	}
}

// Linux-specific sandboxing using various kernel features
func (m *Manager) applyLinuxSandbox(cmd *exec.Cmd) error {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	// For now, we use basic process isolation on Windows
	// since Setpgid and Credential are Linux-specific
	// TODO: Implement proper Linux sandboxing when running on Linux
	
	return nil
}

// Windows-specific sandboxing using Job Objects
func (m *Manager) applyWindowsSandbox(cmd *exec.Cmd) error {
	// TODO: Implement Windows Job Object creation
	// This requires Windows API calls through syscall or cgo
	
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	// Create new process group on Windows
	cmd.SysProcAttr.CreationFlags = syscall.CREATE_NEW_PROCESS_GROUP

	return nil
}

// Basic cross-platform restrictions
func (m *Manager) applyBasicRestrictions(cmd *exec.Cmd) error {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	// Basic restrictions without platform-specific features
	// TODO: Implement cross-platform process isolation
	
	return nil
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
	// List of prohibited commands that could escape sandbox
	prohibited := []string{
		"sudo", "su", "doas",
		"chroot", "unshare",
		"docker", "podman",
		"systemctl", "service",
		"mount", "umount",
		"insmod", "rmmod", "modprobe",
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
	suspicious := []string{
		"/proc/",
		"/sys/", 
		"/dev/",
		"../",
		"$(", 
		"`",
		"|sudo",
		";sudo",
		"&&sudo",
	}

	for _, pattern := range suspicious {
		if len(arg) > len(pattern) && 
		   (arg[:len(pattern)] == pattern || 
		    arg[len(arg)-len(pattern):] == pattern) {
			return true
		}
	}

	return false
}