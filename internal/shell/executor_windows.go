//go:build windows

package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// findExecutable searches for an executable in PATH with Windows-specific extensions
func (s *Shell) findExecutable(command string) (string, error) {
	// If command already has an extension, try it directly
	if strings.Contains(command, ".") {
		if _, err := os.Stat(command); err == nil {
			return command, nil
		}
		
		// Try absolute path
		if filepath.IsAbs(command) {
			return "", fmt.Errorf("executable not found: %s", command)
		}
	}
	
	// Get PATH environment variable
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return "", fmt.Errorf("PATH environment variable not set")
	}
	
	// Windows executable extensions to try
	extensions := []string{".exe", ".com", ".bat", ".cmd", ".ps1"}
	
	// If command already has extension, only try that
	if strings.Contains(command, ".") {
		extensions = []string{""}
	}
	
	// Split PATH and search each directory
	pathDirs := strings.Split(pathEnv, ";")
	
	for _, dir := range pathDirs {
		if dir == "" {
			continue
		}
		
		for _, ext := range extensions {
			candidate := filepath.Join(dir, command+ext)
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				// Check if file is executable (basic check for Windows)
				return candidate, nil
			}
		}
	}
	
	return "", fmt.Errorf("command not found: %s", command)
}

// executeExternalPlatform executes external commands on Windows with proper PATH resolution
func (s *Shell) executeExternalPlatform(cmd *Command) (*ExecutionResult, error) {
	// Try to find the executable
	executable, err := s.findExecutable(cmd.Name)
	if err != nil {
		// Return a helpful error message
		return &ExecutionResult{
			Error:    s.getCommandNotFoundMessage(cmd.Name),
			ExitCode: 127, // Standard "command not found" exit code
		}, fmt.Errorf("command not found: %s", cmd.Name)
	}
	
	// Create command with found executable
	var execCmd *exec.Cmd
	
	// Check file extension to determine how to execute
	ext := strings.ToLower(filepath.Ext(executable))
	
	switch ext {
	case ".ps1":
		// PowerShell script
		execCmd = exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", executable)
		if len(cmd.Args) > 0 {
			execCmd.Args = append(execCmd.Args, cmd.Args...)
		}
	case ".bat", ".cmd":
		// Batch file
		execCmd = exec.Command("cmd.exe", "/C", executable)
		if len(cmd.Args) > 0 {
			execCmd.Args = append(execCmd.Args, cmd.Args...)
		}
	default:
		// Regular executable
		execCmd = exec.Command(executable, cmd.Args...)
	}
	
	// Set working directory
	execCmd.Dir = s.currentDir
	
	// Set environment
	execCmd.Env = s.buildEnvSlice()
	
	// Handle input/output
	var outputBuilder strings.Builder
	var errorBuilder strings.Builder
	
	// For now, ignore custom Input/Output
	execCmd.Stdout = &outputBuilder
	
	execCmd.Stderr = &errorBuilder
	
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
	
	// Execute command
	err = execCmd.Run()
	exitCode := 0
	
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			} else {
				exitCode = 1
			}
		} else {
			return &ExecutionResult{
				Error:    err.Error(),
				ExitCode: 1,
			}, err
		}
	}
	
	return &ExecutionResult{
		Output:   outputBuilder.String(),
		Error:    errorBuilder.String(),
		ExitCode: exitCode,
	}, nil
}

// getCommandNotFoundMessage returns a user-friendly error message for Windows
func (s *Shell) getCommandNotFoundMessage(command string) string {
	suggestions := []string{}
	
	// Check for common typos or similar commands
	commonCommands := map[string]string{
		"ls":    "dir or Get-ChildItem",
		"cat":   "type or Get-Content", 
		"grep":  "findstr or Select-String",
		"ps":    "tasklist or Get-Process",
		"kill":  "taskkill or Stop-Process",
		"which": "where or Get-Command",
		"curl":  "Invoke-WebRequest",
		"wget":  "Invoke-WebRequest",
		"tail":  "Get-Content -Tail",
		"head":  "Get-Content -Head",
	}
	
	if windowsEquivalent, exists := commonCommands[command]; exists {
		suggestions = append(suggestions, fmt.Sprintf("Try: %s", windowsEquivalent))
	}
	
	// Check if it might be a PowerShell command
	if strings.Contains(command, "-") || 
	   strings.HasPrefix(command, "Get-") ||
	   strings.HasPrefix(command, "Set-") {
		suggestions = append(suggestions, "This looks like a PowerShell command - it should work in Lucien")
	}
	
	// Check PATH
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		suggestions = append(suggestions, "PATH environment variable is not set")
	}
	
	baseMsg := fmt.Sprintf("Command not found: '%s'", command)
	
	if len(suggestions) > 0 {
		return baseMsg + "\n" + strings.Join(suggestions, "\n")
	}
	
	return baseMsg + "\nTry 'help' for available commands or check your PATH"
}