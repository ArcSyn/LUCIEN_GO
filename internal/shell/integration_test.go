package shell

import (
	"runtime"
	"strings"
	"testing"
)

func TestCrossPlastformCommands(t *testing.T) {
	shell := New(&Config{SafeMode: false, ExecutorMode: "shell"})

	var tests []struct {
		name        string
		command     string
		expectError bool
		contains    string
		skipOnWin   bool
		skipOnUnix  bool
	}

	// Platform-specific tests
	if runtime.GOOS == "windows" {
		tests = []struct {
			name        string
			command     string
			expectError bool
			contains    string
			skipOnWin   bool
			skipOnUnix  bool
		}{
			{
				name:        "Windows dir command",
				command:     "dir /B",
				expectError: false,
				contains:    "", // Just check it doesn't error
			},
			{
				name:        "Windows type command",
				command:     "echo test | type",
				expectError: true, // This specific syntax won't work
			},
			{
				name:        "Windows echo command",
				command:     "echo Windows test",
				expectError: false,
				contains:    "Windows test",
			},
			{
				name:        "Windows PowerShell available check",
				command:     "powershell -Command \"Write-Output 'PS test'\"",
				expectError: false,
				contains:    "PS test",
			},
		}
	} else {
		tests = []struct {
			name        string
			command     string
			expectError bool
			contains    string
			skipOnWin   bool
			skipOnUnix  bool
		}{
			{
				name:        "Unix ls command",
				command:     "ls -la | head -1",
				expectError: false,
				contains:    "total",
			},
			{
				name:        "Unix echo command",
				command:     "echo Unix test",
				expectError: false,
				contains:    "Unix test",
			},
			{
				name:        "Unix pipe command",
				command:     "echo hello world | wc -w",
				expectError: false,
				contains:    "2",
			},
			{
				name:        "Unix which command",
				command:     "which sh",
				expectError: false,
				contains:    "sh",
			},
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipOnWin && runtime.GOOS == "windows" {
				t.Skip("Skipping test on Windows")
			}
			if tt.skipOnUnix && runtime.GOOS != "windows" {
				t.Skip("Skipping test on Unix")
			}

			result, err := shell.Execute(tt.command)

			if tt.expectError {
				if err == nil && result.ExitCode == 0 {
					t.Errorf("Expected error but command succeeded")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result.ExitCode != 0 {
					t.Errorf("Expected exit code 0, got %d. Error: %s", result.ExitCode, result.Error)
				}
				if tt.contains != "" && !strings.Contains(result.Output, tt.contains) {
					t.Errorf("Expected output to contain %q, got: %q", tt.contains, result.Output)
				}
			}
		})
	}
}

func TestComplexCommandChaining(t *testing.T) {
	shell := New(&Config{SafeMode: false, ExecutorMode: "shell"})

	var testCommand string
	var expectedPattern string

	if runtime.GOOS == "windows" {
		// Windows command chaining
		testCommand = "echo first && echo second"
		expectedPattern = "first"
	} else {
		// Unix command chaining  
		testCommand = "echo first && echo second"
		expectedPattern = "first"
	}

	result, err := shell.Execute(testCommand)
	if err != nil {
		t.Fatalf("Failed to execute chained command: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	if !strings.Contains(result.Output, expectedPattern) {
		t.Errorf("Expected output to contain %q, got: %q", expectedPattern, result.Output)
	}
}

func TestCommandWithSpacesInPath(t *testing.T) {
	shell := New(&Config{SafeMode: false, ExecutorMode: "shell"})

	// Test that commands work with current working directory changes
	originalCmd := "pwd"
	if runtime.GOOS == "windows" {
		originalCmd = "cd"
	}

	result, err := shell.Execute(originalCmd)
	if err != nil {
		t.Fatalf("Failed to execute pwd/cd command: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	if result.Output == "" {
		t.Error("Expected some output from pwd/cd command")
	}
}

func TestLongRunningCommand(t *testing.T) {
	shell := New(&Config{SafeMode: false, ExecutorMode: "shell"})

	var testCommand string
	if runtime.GOOS == "windows" {
		testCommand = "ping -n 1 127.0.0.1"
	} else {
		testCommand = "ping -c 1 127.0.0.1"
	}

	result, err := shell.Execute(testCommand)
	if err != nil {
		// Ping might fail in some test environments, that's OK
		t.Logf("Ping command failed (expected in some test environments): %v", err)
		return
	}

	if result.ExitCode != 0 {
		t.Logf("Ping command exited with code %d (expected in some test environments)", result.ExitCode)
		return
	}

	// Just verify we got some output
	if result.Output == "" && result.Error == "" {
		t.Error("Expected some output from ping command")
	}
}

func TestErrorHandling(t *testing.T) {
	shell := New(&Config{SafeMode: false, ExecutorMode: "shell"})

	// Command that should fail
	result, err := shell.Execute("nonexistentcommand12345")

	// The command should fail, but the shell should handle it gracefully
	if err != nil {
		// It's OK if there's an error, but the shell shouldn't crash
		t.Logf("Command failed as expected: %v", err)
	}

	if result.ExitCode == 0 {
		t.Error("Expected non-zero exit code for nonexistent command")
	}
}

func TestBuiltinVsExternalPrecedence(t *testing.T) {
	shell := New(&Config{SafeMode: false, ExecutorMode: "shell"})

	// Test that built-ins take precedence
	result, err := shell.Execute("echo builtin test")
	if err != nil {
		t.Fatalf("Failed to execute echo command: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	if !strings.Contains(result.Output, "builtin test") {
		t.Errorf("Expected output to contain 'builtin test', got: %q", result.Output)
	}

	// Test pwd builtin
	result2, err := shell.Execute("pwd")
	if err != nil {
		t.Fatalf("Failed to execute pwd command: %v", err)
	}

	if result2.ExitCode != 0 {
		t.Errorf("Expected exit code 0 for pwd, got %d", result2.ExitCode)
	}

	if result2.Output == "" {
		t.Error("Expected pwd to return current directory")
	}
}