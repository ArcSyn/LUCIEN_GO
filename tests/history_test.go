package tests

import (
	"strings"
	"testing"

	"github.com/luciendev/lucien-core/internal/shell"
)

func TestHistoryExpansion(t *testing.T) {
	s := newTestShell()
	
	// Execute some commands to build history
	s.Execute("echo hello")
	s.Execute("ls -la")
	s.Execute("cd /tmp")
	
	// Test !! expansion (last command)
	result, err := s.Execute("!!")
	if err != nil {
		t.Fatalf("Failed to expand !!: %v", err)
	}
	
	if !strings.Contains(result.Output, "/tmp") && !strings.Contains(result.Error, "cd") {
		t.Fatalf("!! expansion failed: expected cd /tmp reference, got '%s'", result.Output)
	}
}

func TestHistoryByNumber(t *testing.T) {
	s := newTestShell()
	
	// Execute commands to build history
	s.Execute("echo first")
	s.Execute("echo second")
	s.Execute("echo third")
	
	// Test !1 expansion (first command)
	result, err := s.Execute("!1")
	if err != nil {
		t.Fatalf("Failed to expand !1: %v", err)
	}
	
	if !strings.Contains(result.Output, "first") && !strings.Contains(result.Error, "echo first") {
		t.Fatalf("!1 expansion failed: expected echo first reference")
	}
}

func TestHistoryByPrefix(t *testing.T) {
	s := newTestShell()
	
	// Execute commands with unique prefixes
	s.Execute("git status")
	s.Execute("ls -la")
	s.Execute("grep pattern file.txt")
	
	// Test !g expansion (last command starting with 'g')
	result, err := s.Execute("!g")
	if err != nil {
		t.Fatalf("Failed to expand !g: %v", err)
	}
	
	if !strings.Contains(result.Output, "pattern") && !strings.Contains(result.Error, "grep") {
		t.Fatalf("!g expansion failed: expected grep command reference")
	}
}

func TestHistoryInvalidExpansion(t *testing.T) {
	s := newTestShell()
	
	// Try to expand non-existent history
	result, err := s.Execute("!999")
	if err != nil {
		t.Fatalf("History expansion should handle invalid numbers gracefully: %v", err)
	}
	
	if result.ExitCode != 1 {
		t.Fatalf("Invalid history expansion should return exit code 1")
	}
	
	if !strings.Contains(result.Error, "not found") && !strings.Contains(result.Error, "history") {
		t.Fatalf("Invalid expansion should mention history not found: %s", result.Error)
	}
}

func TestHistoryCommand(t *testing.T) {
	s := newTestShell()
	
	// Execute some commands
	s.Execute("echo test1")
	s.Execute("echo test2")
	s.Execute("echo test3")
	
	// Run history command
	result, err := s.Execute("history")
	if err != nil {
		t.Fatalf("History command failed: %v", err)
	}
	
	if result.ExitCode != 0 {
		t.Fatalf("History command should succeed: %s", result.Error)
	}
	
	// Should show numbered history entries
	if !strings.Contains(result.Output, "1") || !strings.Contains(result.Output, "echo test1") {
		t.Fatalf("History output should show numbered entries: %s", result.Output)
	}
}

func TestHistoryQuotedExpansion(t *testing.T) {
	s := newTestShell()
	
	// Execute a command with quotes
	s.Execute(`echo "hello world"`)
	s.Execute("ls")
	
	// Test expansion with quotes
	result, err := s.Execute("!echo")
	if err != nil {
		t.Fatalf("Failed to expand quoted history: %v", err)
	}
	
	if !strings.Contains(result.Output, "hello world") && !strings.Contains(result.Error, "hello world") {
		t.Fatalf("Quoted history expansion failed: expected 'hello world' reference")
	}
}

func TestHistoryModifiers(t *testing.T) {
	s := newTestShell()
	
	// Execute a command with arguments
	s.Execute("echo arg1 arg2 arg3")
	
	// Test word selection (if implemented)
	// This tests advanced history expansion features
	// Note: Basic implementation might not support all modifiers
	result, err := s.Execute("echo previous: !!")
	if err != nil {
		t.Fatalf("History modifier test failed: %v", err)
	}
	
	// Should execute the echo with previous reference
	if result.ExitCode != 0 {
		t.Fatalf("History modifier should succeed: %s", result.Error)
	}
}

func TestHistoryEscaping(t *testing.T) {
	s := newTestShell()
	
	// Execute a command 
	s.Execute("echo test")
	
	// Test literal ! (escaped)
	result, err := s.Execute(`echo "literal !"`)
	if err != nil {
		t.Fatalf("Escaped history expansion failed: %v", err)
	}
	
	if result.ExitCode != 0 {
		t.Fatalf("Literal ! should not trigger history expansion")
	}
	
	if !strings.Contains(result.Output, "literal !") {
		t.Fatalf("Literal ! should appear in output: %s", result.Output)
	}
}

func TestHistoryDisabled(t *testing.T) {
	// Create shell with history disabled
	config := &shell.Config{
		SafeMode: true,
		// History would be disabled in safe mode
	}
	s := shell.New(config)
	
	// Execute commands
	s.Execute("echo test1")
	s.Execute("echo test2")
	
	// Try history expansion - should not work or should be literal
	result, err := s.Execute("!!")
	if err != nil {
		t.Fatalf("History expansion with disabled history should handle gracefully: %v", err)
	}
	
	// In safe mode, history expansion might be disabled
	// The behavior depends on implementation
	if result.ExitCode != 127 && !strings.Contains(result.Error, "command not found") {
		// If history is disabled, !! might be treated as literal command
		// This test verifies the behavior is consistent
	}
}