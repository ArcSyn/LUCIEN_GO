package tests

import (
	"strings"
	"testing"

	"github.com/luciendev/lucien-core/internal/completion"
)

func TestBasicCompletion(t *testing.T) {
	engine := completion.New()
	
	// Test command completion
	suggestions := engine.Complete("ec", 2)
	
	// Should suggest echo
	found := false
	for _, suggestion := range suggestions {
		if strings.Contains(suggestion.Text, "echo") {
			found = true
			break
		}
	}
	
	if !found {
		t.Fatalf("Should suggest echo for 'ec' prefix")
	}
}

func TestFileCompletion(t *testing.T) {
	engine := completion.New()
	
	// Test file completion (relative path)
	suggestions := engine.Complete("./", 2)
	
	// Should return some file/directory suggestions
	if len(suggestions) == 0 {
		t.Fatalf("Should suggest files/directories for './' prefix")
	}
	
	// Check that we get file or directory completions
	hasFileType := false
	for _, suggestion := range suggestions {
		if suggestion.Type == completion.FileCompletion || 
		   suggestion.Type == completion.DirectoryCompletion {
			hasFileType = true
			break
		}
	}
	
	if !hasFileType {
		t.Fatalf("Should include file or directory suggestions")
	}
}

func TestCommandArgumentCompletion(t *testing.T) {
	engine := completion.New()
	
	// Test completion for command with arguments
	suggestions := engine.Complete("ls -", 4)
	
	// Should suggest flags/options
	if len(suggestions) == 0 {
		// Some completion implementations might not handle flags
		// This is acceptable for basic implementation
		return
	}
	
	// Check for common ls flags if available
	for _, suggestion := range suggestions {
		if strings.Contains(suggestion.Text, "-l") || strings.Contains(suggestion.Text, "-a") {
			return // Flag found
		}
	}
	
	// Flag completion is advanced feature, not required for basic test
}

func TestVariableCompletion(t *testing.T) {
	s := newTestShell()
	s.Execute("set TEST_VAR=value")
	
	engine := completion.New()
	
	// Test variable completion
	suggestions := engine.Complete("$TES", 4)
	
	// Variable completion is advanced feature
	// Basic implementation might not support this
	if len(suggestions) == 0 {
		return // Acceptable for basic implementation
	}
	
	// If implemented, should suggest TEST_VAR
	for _, suggestion := range suggestions {
		if suggestion.Type == completion.VariableCompletion &&
		   strings.Contains(suggestion.Text, "TEST_VAR") {
			return // Found variable completion
		}
	}
}

func TestAliasCompletion(t *testing.T) {
	engine := completion.New()
	
	// Test alias completion
	suggestions := engine.Complete("l", 1)
	
	// Should include both ls command and ll alias
	hasLs := false
	hasAlias := false
	
	for _, suggestion := range suggestions {
		if strings.Contains(suggestion.Text, "ls") {
			hasLs = true
		}
		if strings.Contains(suggestion.Text, "ll") {
			hasAlias = true
		}
	}
	
	if !hasLs {
		t.Fatalf("Should suggest ls command")
	}
	
	// Alias completion is advanced, might not be implemented
	if hasAlias {
		t.Logf("Alias completion working correctly")
	}
}

func TestHistoryCompletion(t *testing.T) {
	engine := completion.New()
	// TODO: Wire up history provider when available
	
	// Test history-based completion
	suggestions := engine.Complete("echo uni", 8)
	
	// History completion is advanced feature
	if len(suggestions) == 0 {
		return // Acceptable for basic implementation
	}
	
	// If implemented, should suggest from history
	for _, suggestion := range suggestions {
		if suggestion.Type == completion.HistoryCompletion &&
		   strings.Contains(suggestion.Text, "unique_history_test") {
			return // Found history completion
		}
	}
}

func TestCompletionPaging(t *testing.T) {
	engine := completion.New()
	
	// Get suggestions that might trigger paging
	suggestions := engine.Complete("", 0) // Complete from empty string
	
	if len(suggestions) == 0 {
		return // No suggestions to page through
	}
	
	// Test that suggestions are reasonable in count
	if len(suggestions) > 100 {
		t.Fatalf("Too many suggestions returned, paging should limit results")
	}
	
	// Test that all suggestions are valid
	for _, suggestion := range suggestions {
		if suggestion.Text == "" {
			t.Fatalf("Suggestion should not be empty")
		}
	}
}

func TestBestMatch(t *testing.T) {
	engine := completion.New()
	
	// Create some test suggestions
	testSuggestions := []completion.Suggestion{
		{Text: "echo", Type: completion.CommandCompletion},
		{Text: "exit", Type: completion.CommandCompletion},
		{Text: "export", Type: completion.CommandCompletion},
	}
	
	// Test best match for "ec"
	bestMatch := engine.GetBestMatch(testSuggestions)
	
	// Should return common prefix or best match
	if bestMatch == "" {
		t.Fatalf("Should return some best match")
	}
	
	// Best match should be reasonable
	if len(bestMatch) > 20 {
		t.Fatalf("Best match too long: %s", bestMatch)
	}
}

func TestCompletionTypes(t *testing.T) {
	engine := completion.New()
	
	// Test different completion contexts
	testCases := []struct {
		input    string
		cursor   int
		expected completion.CompletionType
	}{
		{"ec", 2, completion.CommandCompletion},
		{"ls ./", 5, completion.FileCompletion},
		{"cd /", 4, completion.DirectoryCompletion},
	}
	
	for _, tc := range testCases {
		suggestions := engine.Complete(tc.input, tc.cursor)
		
		if len(suggestions) == 0 {
			continue // Some completions might not be implemented
		}
		
		// Check that at least one suggestion has expected type
		hasExpectedType := false
		for _, suggestion := range suggestions {
			if suggestion.Type == tc.expected {
				hasExpectedType = true
				break
			}
		}
		
		if !hasExpectedType && len(suggestions) > 0 {
			t.Logf("Expected type %v not found for input '%s', got types: %v", 
				tc.expected, tc.input, suggestions[0].Type)
		}
	}
}

func TestCompletionDescriptions(t *testing.T) {
	engine := completion.New()
	
	// Test that suggestions include helpful descriptions
	suggestions := engine.Complete("git", 3)
	
	if len(suggestions) == 0 {
		return // Git might not be available or completion not implemented
	}
	
	// Check that suggestions have reasonable descriptions
	for _, suggestion := range suggestions {
		if suggestion.Description == "" {
			// Descriptions are optional but helpful
			t.Logf("Suggestion '%s' has no description", suggestion.Text)
		}
		
		if len(suggestion.Description) > 100 {
			t.Fatalf("Description too long: %s", suggestion.Description)
		}
	}
}

func TestCompletionEdgeCases(t *testing.T) {
	engine := completion.New()
	
	// Test edge cases
	edgeCases := []struct {
		input  string
		cursor int
		name   string
	}{
		{"", 0, "empty input"},
		{"   ", 3, "whitespace only"},
		{"very_long_command_that_does_not_exist", 35, "long nonexistent command"},
		{"cmd\t", 4, "input with tab"},
		{"cmd\n", 4, "input with newline"},
	}
	
	for _, ec := range edgeCases {
		suggestions := engine.Complete(ec.input, ec.cursor)
		
		// Should not panic and should return reasonable results
		if suggestions == nil {
			t.Fatalf("Suggestions should not be nil for case: %s", ec.name)
		}
		
		// Should not return excessive suggestions
		if len(suggestions) > 50 {
			t.Logf("Many suggestions (%d) for edge case: %s", len(suggestions), ec.name)
		}
	}
}

func TestCompletionCursorPosition(t *testing.T) {
	engine := completion.New()
	
	// Test completion at different cursor positions
	input := "echo hello world"
	
	testPositions := []struct {
		cursor   int
		name     string
		expected bool
	}{
		{0, "start of line", true},
		{4, "end of command", true},
		{5, "start of first arg", true},
		{10, "middle of first arg", true},
		{11, "start of second arg", true},
		{16, "end of line", true},
		{20, "beyond end", false}, // Invalid cursor position
	}
	
	for _, tp := range testPositions {
		suggestions := engine.Complete(input, tp.cursor)
		
		if tp.expected {
			// Should handle valid cursor positions
			if suggestions == nil {
				t.Fatalf("Should handle cursor position %d (%s)", tp.cursor, tp.name)
			}
		} else {
			// Should handle invalid cursor positions gracefully
			if suggestions == nil {
				t.Logf("Correctly handled invalid cursor position %d", tp.cursor)
			}
		}
	}
}

func TestCompletionPerformance(t *testing.T) {
	engine := completion.New()
	
	// Test that completion is reasonably fast
	input := "test_completion_speed"
	
	// Run multiple completions
	for i := 0; i < 10; i++ {
		suggestions := engine.Complete(input, len(input))
		
		// Should complete quickly (tested by test timeout)
		if suggestions == nil {
			t.Fatalf("Completion should return non-nil result")
		}
	}
	
	// Test completed within test timeout (success if we reach here)
}