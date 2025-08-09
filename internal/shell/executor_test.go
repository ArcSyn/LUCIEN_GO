package shell

import (
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestExecutorModeConfiguration(t *testing.T) {
	// Test default configuration
	shell1 := New(nil)
	if shell1.config.ExecutorMode != "shell" {
		t.Errorf("Expected default executor mode 'shell', got %q", shell1.config.ExecutorMode)
	}

	// Test explicit configuration
	shell2 := New(&Config{SafeMode: true, ExecutorMode: "internal"})
	if shell2.config.ExecutorMode != "internal" {
		t.Errorf("Expected executor mode 'internal', got %q", shell2.config.ExecutorMode)
	}

	// Test empty configuration gets default
	shell3 := New(&Config{SafeMode: false})
	if shell3.config.ExecutorMode != "shell" {
		t.Errorf("Expected default executor mode 'shell' for empty config, got %q", shell3.config.ExecutorMode)
	}
}

func TestShellDelegation(t *testing.T) {
	shell := New(&Config{SafeMode: false, ExecutorMode: "shell"})

	var testCommand string
	var expectedOutput string

	if runtime.GOOS == "windows" {
		testCommand = "echo hello"
		expectedOutput = "hello"
	} else {
		testCommand = "echo hello"
		expectedOutput = "hello"
	}

	result, err := shell.Execute(testCommand)
	if err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	if !strings.Contains(result.Output, expectedOutput) {
		t.Errorf("Expected output to contain %q, got %q", expectedOutput, result.Output)
	}
}

func TestBuiltinCommandsStillWork(t *testing.T) {
	shell := New(&Config{SafeMode: false, ExecutorMode: "shell"})

	// Test built-in echo
	result, err := shell.Execute("echo test builtin")
	if err != nil {
		t.Fatalf("Failed to execute builtin command: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	if !strings.Contains(result.Output, "test builtin") {
		t.Errorf("Expected output to contain 'test builtin', got %q", result.Output)
	}
}

func TestExecutorModeFallback(t *testing.T) {
	// Test that internal mode still works
	shell := New(&Config{SafeMode: false, ExecutorMode: "internal"})

	// This should use the old executor system
	result, err := shell.Execute("echo fallback test")
	if err != nil {
		t.Fatalf("Failed to execute command in internal mode: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	if !strings.Contains(result.Output, "fallback test") {
		t.Errorf("Expected output to contain 'fallback test', got %q", result.Output)
	}
}

func TestProcessMessages(t *testing.T) {
	shell := New(&Config{SafeMode: false, ExecutorMode: "shell"})

	var receivedMessages []ProcessMessage
	dispatcher := func(msg ProcessMessage) {
		receivedMessages = append(receivedMessages, msg)
	}

	shell.SetDispatcher(dispatcher)

	// Execute a simple command that should generate messages
	var testCommand string
	if runtime.GOOS == "windows" {
		testCommand = "echo message test"
	} else {
		testCommand = "echo message test"
	}

	result, err := shell.Execute(testCommand)
	if err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	// Give some time for goroutines to complete
	time.Sleep(100 * time.Millisecond)

	// Check that we received some process messages
	if len(receivedMessages) == 0 {
		t.Error("Expected to receive process messages, got none")
	}

	// Verify message types
	var hasStarted, hasOutput, hasExited bool
	for _, msg := range receivedMessages {
		switch msg.(type) {
		case ProcessStartedMsg:
			hasStarted = true
		case ProcessStdoutMsg:
			hasOutput = true
		case ProcessExitedMsg:
			hasExited = true
		}
	}

	if !hasStarted {
		t.Error("Expected ProcessStartedMsg")
	}
	if !hasOutput {
		t.Error("Expected ProcessStdoutMsg")
	}
	if !hasExited {
		t.Error("Expected ProcessExitedMsg")
	}
}

func TestShellCommandCreation(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		expectedCmd string
		expectedArgs []string
	}{
		{
			name:    "simple command",
			command: "echo hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := shellCommand(tt.command)
			
			if runtime.GOOS == "windows" {
				if cmd.Path != "" && !strings.Contains(cmd.Path, "cmd") {
					t.Errorf("Expected Windows to use cmd, got %q", cmd.Path)
				}
				if len(cmd.Args) < 3 || cmd.Args[1] != "/C" {
					t.Errorf("Expected Windows cmd /C format, got %v", cmd.Args)
				}
			} else {
				if cmd.Path != "" && !strings.Contains(cmd.Path, "sh") {
					t.Errorf("Expected Unix to use sh, got %q", cmd.Path)
				}
				if len(cmd.Args) < 3 || cmd.Args[1] != "-c" {
					t.Errorf("Expected Unix sh -c format, got %v", cmd.Args)
				}
			}
		})
	}
}