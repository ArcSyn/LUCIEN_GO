package tests

import (
	"strings"
	"testing"

	"github.com/ArcSyn/LucienCLI/internal/shell"
)

// newTestShell creates a shell instance for testing
func newTestShell() *shell.Shell {
	config := &shell.Config{
		SafeMode: false, // Disable security for tests
	}
	return shell.New(config)
}

func TestAliasQuoted(t *testing.T) {
	s := newTestShell()
	
	// Test quoted alias creation
	result, err := s.Execute(`alias g='git status'`)
	if err != nil {
		t.Fatalf("Failed to create alias: %v", err)
	}
	
	if result.ExitCode != 0 {
		t.Fatalf("Alias creation failed with exit code %d: %s", result.ExitCode, result.Error)
	}
	
	// Test alias execution (this will fail since git isn't available, but alias should be recognized)
	result, err = s.Execute("g")
	if err != nil && !strings.Contains(err.Error(), "command not found") {
		t.Fatalf("Unexpected error executing alias: %v", err)
	}
	
	// The alias should be expanded, so we should see git-related error, not g-related error
	if strings.Contains(result.Error, "'g'") {
		t.Fatalf("Alias was not expanded properly")
	}
}

func TestNestedAliasQuotes(t *testing.T) {
	s := newTestShell()
	
	// Test nested quotes in alias
	result, err := s.Execute(`alias s="echo 'hi'"`)
	if err != nil {
		t.Fatalf("Failed to create nested quote alias: %v", err)
	}
	
	if result.ExitCode != 0 {
		t.Fatalf("Nested quote alias creation failed: %s", result.Error)
	}
	
	// Execute the alias
	result, err = s.Execute("s")
	if err != nil {
		t.Fatalf("Failed to execute nested quote alias: %v", err)
	}
	
	if !strings.Contains(result.Output, "hi") {
		t.Fatalf("Nested quote alias output wrong: got '%s', expected to contain 'hi'", result.Output)
	}
}

func TestAliasEqualsFormat(t *testing.T) {
	s := newTestShell()
	
	// Test name=value format
	result, err := s.Execute("alias l=ls")
	if err != nil {
		t.Fatalf("Failed to create equals format alias: %v", err)
	}
	
	if result.ExitCode != 0 {
		t.Fatalf("Equals format alias creation failed: %s", result.Error)
	}
	
	// Test execution (will fail on command not found but alias should be recognized)
	result, err = s.Execute("l")
	if err != nil && !strings.Contains(err.Error(), "command not found") {
		t.Fatalf("Unexpected error executing equals format alias: %v", err)
	}
}

func TestUnalias(t *testing.T) {
	s := newTestShell()
	
	// Create an alias
	result, err := s.Execute("alias test=echo")
	if err != nil || result.ExitCode != 0 {
		t.Fatalf("Failed to create test alias: %v", err)
	}
	
	// Remove the alias
	result, err = s.Execute("unalias test")
	if err != nil {
		t.Fatalf("Failed to remove alias: %v", err)
	}
	
	if result.ExitCode != 0 {
		t.Fatalf("Unalias failed: %s", result.Error)
	}
	
	if !strings.Contains(result.Output, "Removed alias: test") {
		t.Fatalf("Unalias output incorrect: %s", result.Output)
	}
}

func TestUnaliasNotFound(t *testing.T) {
	s := newTestShell()
	
	// Try to remove non-existent alias
	result, err := s.Execute("unalias nonexistent")
	if err != nil {
		t.Fatalf("Unalias should handle non-existent aliases gracefully: %v", err)
	}
	
	if result.ExitCode != 1 {
		t.Fatalf("Unalias should return exit code 1 for non-existent alias")
	}
	
	if !strings.Contains(result.Error, "not found") {
		t.Fatalf("Unalias should report alias not found: %s", result.Error)
	}
}

func TestUnaliasTypoSuggestion(t *testing.T) {
	s := newTestShell()
	
	// Create an alias
	result, err := s.Execute("alias test=echo")
	if err != nil || result.ExitCode != 0 {
		t.Fatalf("Failed to create test alias: %v", err)
	}
	
	// Try to remove with typo (extra period)
	result, err = s.Execute("unalias test.")
	if err != nil {
		t.Fatalf("Unalias should handle typos gracefully: %v", err)
	}
	
	if result.ExitCode != 1 {
		t.Fatalf("Unalias should return exit code 1 for typo")
	}
	
	if !strings.Contains(result.Error, "Did you mean 'test'?") {
		t.Fatalf("Unalias should suggest correct alias name: %s", result.Error)
	}
}