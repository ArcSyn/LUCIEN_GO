//go:build !windows

package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// findExecutable searches for an executable in PATH on Unix-like systems
func (s *Shell) findExecutable(command string) (string, error) {
	// If command contains a slash, it's a path (relative or absolute)
	if strings.Contains(command, "/") {
		if filepath.IsAbs(command) {
			// Absolute path
			if info, err := os.Stat(command); err == nil && !info.IsDir() {
				// Check if file is executable
				if info.Mode()&0111 != 0 {
					return command, nil
				}
				return "", fmt.Errorf("file is not executable: %s", command)
			}
			return "", fmt.Errorf("file not found: %s", command)
		}
		
		// Relative path
		absPath, err := filepath.Abs(command)
		if err != nil {
			return "", err
		}
		
		if info, err := os.Stat(absPath); err == nil && !info.IsDir() {
			if info.Mode()&0111 != 0 {
				return absPath, nil
			}
			return "", fmt.Errorf("file is not executable: %s", command)
		}
		return "", fmt.Errorf("file not found: %s", command)
	}
	
	// Search in PATH
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return "", fmt.Errorf("PATH environment variable not set")
	}
	
	// Split PATH and search each directory
	pathDirs := strings.Split(pathEnv, ":")
	
	for _, dir := range pathDirs {
		if dir == "" {
			continue
		}
		
		candidate := filepath.Join(dir, command)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			// Check if file is executable
			if info.Mode()&0111 != 0 {
				return candidate, nil
			}
		}
	}
	
	return "", fmt.Errorf("command not found: %s", command)
}

// executeExternalPlatform executes external commands on Unix-like systems with proper PATH resolution
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
	execCmd := exec.Command(executable, cmd.Args...)
	
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

// getCommandNotFoundMessage returns a user-friendly error message for Unix systems
func (s *Shell) getCommandNotFoundMessage(command string) string {
	suggestions := []string{}
	
	// Check for common typos or missing packages
	packageSuggestions := map[string][]string{
		"git":    {"sudo apt install git", "brew install git", "yum install git"},
		"curl":   {"sudo apt install curl", "brew install curl", "yum install curl"},
		"wget":   {"sudo apt install wget", "brew install wget", "yum install wget"},
		"vim":    {"sudo apt install vim", "brew install vim", "yum install vim"},
		"emacs":  {"sudo apt install emacs", "brew install emacs", "yum install emacs"},
		"htop":   {"sudo apt install htop", "brew install htop", "yum install htop"},
		"tree":   {"sudo apt install tree", "brew install tree", "yum install tree"},
		"jq":     {"sudo apt install jq", "brew install jq", "yum install jq"},
		"python": {"sudo apt install python3", "brew install python3", "yum install python3"},
		"node":   {"sudo apt install nodejs", "brew install node", "yum install nodejs"},
		"docker": {"Install Docker from https://docker.com"},
	}
	
	if packages, exists := packageSuggestions[command]; exists {
		suggestions = append(suggestions, "This command might need to be installed:")
		suggestions = append(suggestions, packages...)
	}
	
	// Check for common alternatives
	alternatives := map[string]string{
		"dir":      "ls",
		"cls":      "clear", 
		"type":     "cat",
		"copy":     "cp",
		"move":     "mv",
		"del":      "rm",
		"md":       "mkdir",
		"rd":       "rmdir",
		"tasklist": "ps",
		"findstr":  "grep",
	}
	
	if alternative, exists := alternatives[command]; exists {
		suggestions = append(suggestions, fmt.Sprintf("Try the Unix equivalent: %s", alternative))
	}
	
	// Check PATH
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		suggestions = append(suggestions, "PATH environment variable is not set")
	} else {
		// Suggest checking PATH
		pathDirs := strings.Split(pathEnv, ":")
		suggestions = append(suggestions, fmt.Sprintf("Searched %d directories in PATH", len(pathDirs)))
	}
	
	baseMsg := fmt.Sprintf("Command not found: '%s'", command)
	
	if len(suggestions) > 0 {
		return baseMsg + "\n" + strings.Join(suggestions, "\n")
	}
	
	return baseMsg + "\nTry 'help' for available commands or check if the command is installed"
}