package tests

import (
	"strings"
	"testing"
)

func TestCommandSeparators(t *testing.T) {
	s := newTestShell()
	
	// Test semicolon separator
	result, err := s.Execute("echo first; echo second")
	if err != nil {
		t.Fatalf("Semicolon separator failed: %v", err)
	}
	
	if !strings.Contains(result.Output, "first") || !strings.Contains(result.Output, "second") {
		t.Fatalf("Both commands should execute: %s", result.Output)
	}
}

func TestAndOperator(t *testing.T) {
	s := newTestShell()
	
	// Test && operator with success
	result, err := s.Execute("echo success && echo continued")
	if err != nil {
		t.Fatalf("AND operator failed: %v", err)
	}
	
	if !strings.Contains(result.Output, "success") || !strings.Contains(result.Output, "continued") {
		t.Fatalf("Both commands should execute on success: %s", result.Output)
	}
}

func TestAndOperatorFailure(t *testing.T) {
	s := newTestShell()
	
	// Test && operator with failure (false command doesn't exist, use exit 1)
	result, err := s.Execute("exit 1 && echo should_not_run")
	if err != nil {
		t.Fatalf("AND operator test failed: %v", err)
	}
	
	if strings.Contains(result.Output, "should_not_run") {
		t.Fatalf("Second command should not run after failure")
	}
}

func TestOrOperator(t *testing.T) {
	s := newTestShell()
	
	// Test || operator with failure
	result, err := s.Execute("exit 1 || echo backup")
	if err != nil {
		t.Fatalf("OR operator failed: %v", err)
	}
	
	if !strings.Contains(result.Output, "backup") {
		t.Fatalf("Second command should run after failure: %s", result.Output)
	}
}

func TestOrOperatorSuccess(t *testing.T) {
	s := newTestShell()
	
	// Test || operator with success
	result, err := s.Execute("echo success || echo should_not_run")
	if err != nil {
		t.Fatalf("OR operator success test failed: %v", err)
	}
	
	if strings.Contains(result.Output, "should_not_run") {
		t.Fatalf("Second command should not run after success")
	}
	
	if !strings.Contains(result.Output, "success") {
		t.Fatalf("First command should run: %s", result.Output)
	}
}

func TestComplexCommandChain(t *testing.T) {
	s := newTestShell()
	
	// Test complex chain: cmd1 && cmd2 || cmd3; cmd4
	result, err := s.Execute("echo step1 && echo step2 || echo backup; echo final")
	if err != nil {
		t.Fatalf("Complex command chain failed: %v", err)
	}
	
	// Should see step1, step2, and final (backup should not run)
	if !strings.Contains(result.Output, "step1") {
		t.Fatalf("Step1 should execute")
	}
	if !strings.Contains(result.Output, "step2") {
		t.Fatalf("Step2 should execute after step1 success")
	}
	if !strings.Contains(result.Output, "final") {
		t.Fatalf("Final should execute (semicolon separator)")
	}
	if strings.Contains(result.Output, "backup") {
		t.Fatalf("Backup should not execute after step2 success")
	}
}

func TestQuotedStrings(t *testing.T) {
	s := newTestShell()
	
	// Test single quotes
	result, err := s.Execute("echo 'hello world'")
	if err != nil {
		t.Fatalf("Single quoted string failed: %v", err)
	}
	
	if !strings.Contains(result.Output, "hello world") {
		t.Fatalf("Single quoted content should be preserved: %s", result.Output)
	}
}

func TestDoubleQuotedStrings(t *testing.T) {
	s := newTestShell()
	
	// Test double quotes
	result, err := s.Execute(`echo "hello world"`)
	if err != nil {
		t.Fatalf("Double quoted string failed: %v", err)
	}
	
	if !strings.Contains(result.Output, "hello world") {
		t.Fatalf("Double quoted content should be preserved: %s", result.Output)
	}
}

func TestNestedQuotes(t *testing.T) {
	s := newTestShell()
	
	// Test nested quotes
	result, err := s.Execute(`echo "It's a 'test'"`)
	if err != nil {
		t.Fatalf("Nested quotes failed: %v", err)
	}
	
	if !strings.Contains(result.Output, "It's a 'test'") {
		t.Fatalf("Nested quotes should be preserved: %s", result.Output)
	}
}

func TestQuotedSeparators(t *testing.T) {
	s := newTestShell()
	
	// Test separators inside quotes should be literal
	result, err := s.Execute(`echo "cmd1 && cmd2"; echo after`)
	if err != nil {
		t.Fatalf("Quoted separators failed: %v", err)
	}
	
	// Should see literal && in output, not execute cmd1 and cmd2
	if !strings.Contains(result.Output, "cmd1 && cmd2") {
		t.Fatalf("Quoted separators should be literal: %s", result.Output)
	}
	
	if !strings.Contains(result.Output, "after") {
		t.Fatalf("Command after quoted string should execute: %s", result.Output)
	}
}

func TestEscapedCharacters(t *testing.T) {
	s := newTestShell()
	
	// Test escaped quotes
	result, err := s.Execute(`echo "He said \"hello\""`)
	if err != nil {
		t.Fatalf("Escaped quotes failed: %v", err)
	}
	
	if !strings.Contains(result.Output, `He said "hello"`) {
		t.Fatalf("Escaped quotes should be literal: %s", result.Output)
	}
}

func TestBackslashEscaping(t *testing.T) {
	s := newTestShell()
	
	// Test backslash escaping
	result, err := s.Execute("echo hello\\; echo world")
	if err != nil {
		t.Fatalf("Backslash escaping failed: %v", err)
	}
	
	// Should see literal semicolon, not command separation
	if !strings.Contains(result.Output, "hello;") {
		t.Fatalf("Escaped semicolon should be literal")
	}
	
	// "echo world" should not execute as separate command
	if strings.Contains(result.Output, "world") && !strings.Contains(result.Output, "echo world") {
		t.Fatalf("Escaped separator should not cause command separation")
	}
}

func TestEmptyCommands(t *testing.T) {
	s := newTestShell()
	
	// Test empty commands in chain
	result, err := s.Execute("echo first;; echo second")
	if err != nil {
		t.Fatalf("Empty commands should be handled: %v", err)
	}
	
	// Should execute first and second, ignoring empty command
	if !strings.Contains(result.Output, "first") {
		t.Fatalf("First command should execute")
	}
	if !strings.Contains(result.Output, "second") {
		t.Fatalf("Second command should execute despite empty command")
	}
}

func TestWhitespaceHandling(t *testing.T) {
	s := newTestShell()
	
	// Test whitespace around separators
	result, err := s.Execute("echo first  ;  echo second  &&  echo third")
	if err != nil {
		t.Fatalf("Whitespace handling failed: %v", err)
	}
	
	if !strings.Contains(result.Output, "first") {
		t.Fatalf("First command should execute")
	}
	if !strings.Contains(result.Output, "second") {
		t.Fatalf("Second command should execute")
	}
	if !strings.Contains(result.Output, "third") {
		t.Fatalf("Third command should execute")
	}
}

func TestOperatorPrecedence(t *testing.T) {
	s := newTestShell()
	
	// Test that && has higher precedence than ||
	// cmd1 || cmd2 && cmd3 should be: cmd1 || (cmd2 && cmd3)
	result, err := s.Execute("exit 1 || echo step2 && echo step3")
	if err != nil {
		t.Fatalf("Operator precedence test failed: %v", err)
	}
	
	// Should see both step2 and step3 (step2 runs due to ||, step3 runs due to && after step2 success)
	if !strings.Contains(result.Output, "step2") {
		t.Fatalf("Step2 should run after exit 1")
	}
	if !strings.Contains(result.Output, "step3") {
		t.Fatalf("Step3 should run after step2 success")
	}
}

func TestLongCommandChain(t *testing.T) {
	s := newTestShell()
	
	// Test long chain of commands
	longChain := "echo a; echo b && echo c || echo d; echo e && echo f && echo g || echo h"
	result, err := s.Execute(longChain)
	if err != nil {
		t.Fatalf("Long command chain failed: %v", err)
	}
	
	// Should see: a, b, c (d not run), e, f, g (h not run)
	expected := []string{"a", "b", "c", "e", "f", "g"}
	notExpected := []string{"d", "h"}
	
	for _, exp := range expected {
		if !strings.Contains(result.Output, exp) {
			t.Fatalf("Expected '%s' in output: %s", exp, result.Output)
		}
	}
	
	for _, notExp := range notExpected {
		if strings.Contains(result.Output, notExp) {
			t.Fatalf("Did not expect '%s' in output: %s", notExp, result.Output)
		}
	}
}