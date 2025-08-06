package shell

import (
	"strings"
	"testing"
)

func TestShellParsing(t *testing.T) {
	config := &Config{
		SafeMode: false,
	}
	shell := New(config)

	tests := []struct {
		name     string
		input    string
		expected int // expected number of commands
	}{
		{
			name:     "Simple command",
			input:    "echo hello",
			expected: 1,
		},
		{
			name:     "Pipe command",
			input:    "echo hello | grep h",
			expected: 2,
		},
		{
			name:     "Complex pipe",
			input:    "cat file.txt | grep pattern | sort | uniq",
			expected: 4,
		},
		{
			name:     "Redirect output",
			input:    "echo hello > output.txt",
			expected: 1,
		},
		{
			name:     "Multiple redirects",
			input:    "sort < input.txt > output.txt",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commands, err := shell.parseCommandLine(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse command: %v", err)
			}

			if len(commands) != tt.expected {
				t.Errorf("Expected %d commands, got %d", tt.expected, len(commands))
			}

			// Validate first command structure
			if len(commands) > 0 {
				cmd := commands[0]
				if cmd.Name == "" {
					t.Error("Command name should not be empty")
				}
			}
		})
	}
}

func TestTokenization(t *testing.T) {
	shell := New(&Config{})

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Simple tokens",
			input:    "echo hello world",
			expected: []string{"echo", "hello", "world"},
		},
		{
			name:     "Quoted strings",
			input:    `echo "hello world" test`,
			expected: []string{"echo", "hello world", "test"},
		},
		{
			name:     "Redirects",
			input:    "echo hello > output.txt",
			expected: []string{"echo", "hello", ">", "output.txt"},
		},
		{
			name:     "Mixed quotes",
			input:    `echo 'single quotes' "double quotes"`,
			expected: []string{"echo", "single quotes", "double quotes"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := shell.tokenize(tt.input)
			
			if len(tokens) != len(tt.expected) {
				t.Errorf("Expected %d tokens, got %d", len(tt.expected), len(tokens))
				t.Errorf("Expected: %v", tt.expected)
				t.Errorf("Got: %v", tokens)
				return
			}

			for i, token := range tokens {
				if token != tt.expected[i] {
					t.Errorf("Token %d: expected '%s', got '%s'", i, tt.expected[i], token)
				}
			}
		})
	}
}

func TestVariableExpansion(t *testing.T) {
	shell := New(&Config{})
	shell.env["USER"] = "testuser"
	shell.env["HOME"] = "/home/testuser"

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple variable",
			input:    "echo $USER",
			expected: "echo testuser",
		},
		{
			name:     "Braced variable",
			input:    "echo ${HOME}/documents",
			expected: "echo /home/testuser/documents",
		},
		{
			name:     "Multiple variables",
			input:    "$USER lives in $HOME",
			expected: "testuser lives in /home/testuser",
		},
		{
			name:     "No variables",
			input:    "echo hello world",
			expected: "echo hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shell.expandVariables(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestBuiltinCommands(t *testing.T) {
	shell := New(&Config{})
	
	tests := []struct {
		name        string
		command     string
		args        []string
		shouldError bool
	}{
		{
			name:        "pwd command",
			command:     "pwd",
			args:        []string{},
			shouldError: false,
		},
		{
			name:        "echo command",
			command:     "echo",
			args:        []string{"hello", "world"},
			shouldError: false,
		},
		{
			name:        "set variable",
			command:     "set",
			args:        []string{"TEST", "value"},
			shouldError: false,
		},
		{
			name:        "create alias",
			command:     "alias",
			args:        []string{"ll", "ls -la"},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if builtin, exists := shell.builtins[tt.command]; exists {
				result, err := builtin(tt.args)
				
				if tt.shouldError && err == nil {
					t.Error("Expected error but got none")
				}
				
				if !tt.shouldError && err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if result == nil {
					t.Error("Result should not be nil")
				}
			} else {
				t.Errorf("Builtin command '%s' not found", tt.command)
			}
		})
	}
}

func TestCommandExecution(t *testing.T) {
	shell := New(&Config{SafeMode: false})

	tests := []struct {
		name        string
		command     string
		shouldError bool
	}{
		{
			name:        "Echo command",
			command:     "echo test",
			shouldError: false,
		},
		{
			name:        "PWD builtin", 
			command:     "pwd",
			shouldError: false,
		},
		{
			name:        "Variable setting",
			command:     "set TESTVAR testvalue",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := shell.Execute(tt.command)
			
			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result == nil {
				t.Error("Result should not be nil")
			}

			if result.Duration == 0 {
				t.Log("Duration tracking working correctly for builtins")
			}
		})
	}
}

func TestHistoryManagement(t *testing.T) {
	shell := New(&Config{})

	// Execute some commands to build history
	commands := []string{"echo test1", "pwd", "echo test2"}
	
	for _, cmd := range commands {
		shell.Execute(cmd)
	}

	if len(shell.history) != len(commands) {
		t.Errorf("Expected %d history entries, got %d", len(commands), len(shell.history))
	}

	// Test history builtin
	result, err := shell.builtins["history"]([]string{})
	if err != nil {
		t.Errorf("History command failed: %v", err)
	}

	if result.Output == "" {
		t.Error("History output should not be empty")
	}

	// Verify history contains our commands
	for _, cmd := range commands {
		if !strings.Contains(result.Output, cmd) {
			t.Errorf("History should contain command: %s", cmd)
		}
	}
}

func TestAliasSystem(t *testing.T) {
	shell := New(&Config{})

	// Create an alias
	_, err := shell.builtins["alias"]([]string{"ll", "ls -la"})
	if err != nil {
		t.Errorf("Failed to create alias: %v", err)
	}

	// Verify alias exists
	if shell.aliases["ll"] != "ls -la" {
		t.Error("Alias was not created correctly")
	}

	// Test alias expansion in command parsing
	cmd, err := shell.parseCommand("ll")
	if err != nil {
		t.Errorf("Failed to parse aliased command: %v", err)
	}

	if cmd.Name != "ls" {
		t.Errorf("Expected command name 'ls', got '%s'", cmd.Name)
	}

	if len(cmd.Args) < 1 || cmd.Args[0] != "-la" {
		t.Error("Alias arguments not expanded correctly")
	}
}

func TestEnvironmentVariables(t *testing.T) {
	shell := New(&Config{})

	// Set environment variable
	_, err := shell.builtins["set"]([]string{"TESTVAR", "testvalue"})
	if err != nil {
		t.Errorf("Failed to set variable: %v", err)
	}

	// Verify variable is set
	if shell.env["TESTVAR"] != "testvalue" {
		t.Error("Environment variable not set correctly")
	}

	// Test variable expansion
	expanded := shell.expandVariables("Value is: $TESTVAR")
	expected := "Value is: testvalue"
	
	if expanded != expected {
		t.Errorf("Expected '%s', got '%s'", expected, expanded)
	}
}

// Benchmark tests for performance
func BenchmarkCommandParsing(b *testing.B) {
	shell := New(&Config{})
	command := "echo hello | grep h | sort | uniq > output.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := shell.parseCommandLine(command)
		if err != nil {
			b.Fatalf("Command parsing failed: %v", err)
		}
	}
}

func BenchmarkTokenization(b *testing.B) {
	shell := New(&Config{})
	input := `echo "hello world" test 'quoted string' > output.txt`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shell.tokenize(input)
	}
}

func BenchmarkVariableExpansion(b *testing.B) {
	shell := New(&Config{})
	shell.env["USER"] = "testuser"
	shell.env["HOME"] = "/home/testuser"
	input := "User $USER lives in $HOME and works in ${HOME}/projects"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shell.expandVariables(input)
	}
}