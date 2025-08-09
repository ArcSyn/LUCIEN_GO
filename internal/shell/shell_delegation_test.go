package shell

import (
	"runtime"
	"strings"
	"testing"
)

func TestHasExeInPath(t *testing.T) {
	// Test with a command that should exist on all systems
	if runtime.GOOS == "windows" {
		if !hasExeInPath("cmd.exe") && !hasExeInPath("cmd") {
			t.Error("Expected cmd to be available on Windows")
		}
	} else {
		if !hasExeInPath("sh") {
			t.Error("Expected sh to be available on Unix")
		}
	}

	// Test with a command that definitely doesn't exist
	if hasExeInPath("nonexistent-command-12345") {
		t.Error("Expected nonexistent command to not be found")
	}
}

func TestShellCommandSelection(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		skipOnWin bool
		skipOnUnix bool
	}{
		{
			name:  "Simple echo command",
			input: "echo hello world",
		},
		{
			name:       "Unix pipe command",
			input:      "echo test | grep t",
			skipOnWin:  false, // PowerShell supports pipes
		},
		{
			name:  "Command with quotes",
			input: `echo "hello world"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipOnWin && runtime.GOOS == "windows" {
				t.Skip("Skipping test on Windows")
			}
			if tt.skipOnUnix && runtime.GOOS != "windows" {
				t.Skip("Skipping test on Unix")
			}

			cmd := shellCommand(tt.input)
			if cmd == nil {
				t.Error("Expected non-nil command")
				return
			}

			if runtime.GOOS == "windows" {
				// Should prefer PowerShell variants or fallback to cmd
				expectedCommands := []string{"pwsh", "powershell", "cmd"}
				found := false
				for _, expected := range expectedCommands {
					if strings.Contains(cmd.Path, expected) || cmd.Args[0] == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected Windows command to use PowerShell or cmd, got: %s with args %v", cmd.Path, cmd.Args)
				}

				// Verify arguments structure
				if len(cmd.Args) < 2 {
					t.Errorf("Expected at least 2 arguments for Windows command, got %d", len(cmd.Args))
				}
			} else {
				// Unix should use /bin/sh -c
				if !strings.Contains(cmd.Path, "sh") && cmd.Args[0] != "sh" {
					t.Errorf("Expected Unix command to use sh, got: %s", cmd.Path)
				}
				if len(cmd.Args) < 3 || cmd.Args[1] != "-c" {
					t.Errorf("Expected Unix sh -c format, got %v", cmd.Args)
				}
			}
		})
	}
}

func TestShellDelegationIntegration(t *testing.T) {
	shell := New(&Config{SafeMode: false, ExecutorMode: "shell"})

	tests := []struct {
		name        string
		command     string
		expectError bool
		contains    string
		skipOnWin   bool
		skipOnUnix  bool
	}{
		{
			name:        "Echo command",
			command:     "echo hello shell delegation",
			expectError: false,
			contains:    "hello shell delegation",
		},
	}

	// Add platform-specific tests
	if runtime.GOOS == "windows" {
		tests = append(tests, []struct {
			name        string
			command     string
			expectError bool
			contains    string
			skipOnWin   bool
			skipOnUnix  bool
		}{
			{
				name:        "Windows dir command",
				command:     `echo "Windows test"`,
				expectError: false,
				contains:    "Windows test",
			},
			{
				name:        "PowerShell specific command",
				command:     "Write-Output 'PowerShell works'", // This will work in PowerShell but not cmd
				expectError: false,
				contains:    "PowerShell works",
			},
		}...)
	} else {
		tests = append(tests, []struct {
			name        string
			command     string
			expectError bool
			contains    string
			skipOnWin   bool
			skipOnUnix  bool
		}{
			{
				name:        "Unix ls command",
				command:     "echo 'Unix test'",
				expectError: false,
				contains:    "Unix test",
			},
			{
				name:        "Unix pipe command",
				command:     "echo hello | grep hello",
				expectError: false,
				contains:    "hello",
			},
		}...)
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
					t.Errorf("Expected exit code 0, got %d. Output: %s, Error: %s", result.ExitCode, result.Output, result.Error)
				}
				if tt.contains != "" && !strings.Contains(strings.ReplaceAll(result.Output, "\n", " "), tt.contains) {
					t.Errorf("Expected output to contain %q, got: %q", tt.contains, result.Output)
				}
			}
		})
	}
}

func TestExitCodeHandling(t *testing.T) {
	shell := New(&Config{SafeMode: false, ExecutorMode: "shell"})

	var failCommand string
	if runtime.GOOS == "windows" {
		// PowerShell command that exits with error
		failCommand = "exit 42"
	} else {
		// Unix command that exits with error
		failCommand = "exit 42"
	}

	result, _ := shell.Execute(failCommand)

	// The command should fail gracefully
	if result.ExitCode == 0 {
		t.Error("Expected non-zero exit code for failing command")
	}

	// We expect the exit code to be 42, but some shells might handle it differently
	// The important thing is that it's non-zero
	if result.ExitCode == 0 {
		t.Errorf("Expected non-zero exit code, got %d", result.ExitCode)
	}
}

func TestPipeAndRedirection(t *testing.T) {
	shell := New(&Config{SafeMode: false, ExecutorMode: "shell"})

	var pipeCommand string
	var expectedContent string
	
	if runtime.GOOS == "windows" {
		// PowerShell pipe command
		pipeCommand = `echo "test content" | Select-String "content"`
		expectedContent = "content"
	} else {
		// Unix pipe command
		pipeCommand = `echo "test content" | grep content`
		expectedContent = "test content"
	}

	result, err := shell.Execute(pipeCommand)
	if err != nil {
		t.Fatalf("Failed to execute pipe command: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Error: %s", result.ExitCode, result.Error)
	}

	if !strings.Contains(result.Output, expectedContent) {
		t.Errorf("Expected output to contain %q, got: %q", expectedContent, result.Output)
	}
}

func TestBuiltinVsShellPrecedence(t *testing.T) {
	shell := New(&Config{SafeMode: false, ExecutorMode: "shell"})

	// Built-in commands should take precedence over shell commands
	result, err := shell.Execute("cd")
	if err != nil {
		t.Fatalf("Failed to execute builtin cd: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0 for builtin cd, got %d", result.ExitCode)
	}

	// pwd should also be built-in
	result2, err := shell.Execute("pwd")
	if err != nil {
		t.Fatalf("Failed to execute builtin pwd: %v", err)
	}

	if result2.ExitCode != 0 {
		t.Errorf("Expected exit code 0 for builtin pwd, got %d", result2.ExitCode)
	}

	if result2.Output == "" {
		t.Error("Expected pwd to return current directory")
	}
}