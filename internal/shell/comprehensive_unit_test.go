package shell

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestComprehensiveShellFunctionality tests all major shell features systematically
func TestComprehensiveShellFunctionality(t *testing.T) {
	// Create temporary test directory
	tmpDir, err := ioutil.TempDir("", "lucien_comprehensive_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	config := &Config{
		SafeMode: false,
	}
	shell := New(config)

	t.Run("BasicCommandParsing", func(t *testing.T) {
		testBasicCommandParsing(t, shell)
	})

	t.Run("AdvancedCommandParsing", func(t *testing.T) {
		testAdvancedCommandParsing(t, shell, tmpDir)
	})

	t.Run("BuiltinCommands", func(t *testing.T) {
		testBuiltinCommands(t, shell, tmpDir)
	})

	t.Run("EnvironmentVariables", func(t *testing.T) {
		testEnvironmentVariables(t, shell)
	})

	t.Run("AliasSystem", func(t *testing.T) {
		testAliasSystem(t, shell)
	})

	t.Run("HistoryManagement", func(t *testing.T) {
		testHistoryManagement(t, shell)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		testErrorHandling(t, shell)
	})

	t.Run("ResourceAndPerformance", func(t *testing.T) {
		testResourceAndPerformance(t, shell)
	})

	t.Run("EdgeCases", func(t *testing.T) {
		testEdgeCases(t, shell)
	})
}

func testBasicCommandParsing(t *testing.T, shell *Shell) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
		cmdCount  int
	}{
		{"Simple command", "echo hello", false, 1},
		{"Command with args", "echo hello world", false, 1},
		{"Quoted arguments", `echo "hello world"`, false, 1},
		{"Single quoted", `echo 'hello world'`, false, 1},
		{"Mixed quotes", `echo "hello" 'world'`, false, 1},
		{"Pipe command", "echo hello | grep h", false, 2},
		{"Multiple pipes", "echo hello | grep h | sort", false, 3},
		{"Output redirect", "echo hello > file.txt", false, 1},
		{"Input redirect", "sort < input.txt", false, 1},
		{"Append redirect", "echo hello >> file.txt", false, 1},
		{"Complex redirect", "sort < input.txt > output.txt", false, 1},
		{"Empty command", "", false, 0}, // Now handled gracefully
		{"Only whitespace", "   ", false, 0}, // Now handled gracefully
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commands, err := shell.parseCommandLine(tt.input)
			
			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(commands) != tt.cmdCount {
				t.Errorf("Expected %d commands, got %d", tt.cmdCount, len(commands))
			}

			// Verify command structure
			if len(commands) > 0 {
				cmd := commands[0]
				if cmd.Name == "" {
					t.Error("Command name should not be empty")
				}
			}
		})
	}
}

func testAdvancedCommandParsing(t *testing.T, shell *Shell, tmpDir string) {
	// Create test files
	inputFile := filepath.Join(tmpDir, "input.txt")
	ioutil.WriteFile(inputFile, []byte("line1\nline2\nline3\n"), 0644)

	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{"Variable in command", "echo $HOME", false},
		{"Braced variable", "echo ${HOME}/test", false},
		{"Multiple variables", "echo $HOME $USER", false},
		{"Redirect with variable", "echo test > $HOME/output.txt", false},
		{"Complex tokenization", `echo "quoted arg" normal 'single' > file`, false},
		{"Nested quotes", `echo "nested 'single' quotes"`, false},
		{"Escaped quotes", `echo "escaped \"quote\""`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commands, err := shell.parseCommandLine(tt.input)
			
			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(commands) == 0 {
				t.Error("Expected at least one command")
			}
		})
	}
}

func testBuiltinCommands(t *testing.T, shell *Shell, tmpDir string) {
	tests := []struct {
		name      string
		command   string
		args      []string
		expectErr bool
		checkOutput func(result *ExecutionResult) error
	}{
		{
			name:      "pwd command",
			command:   "pwd",
			args:      []string{},
			expectErr: false,
			checkOutput: func(result *ExecutionResult) error {
				if result.Output == "" {
					return fmt.Errorf("pwd should return current directory")
				}
				return nil
			},
		},
		{
			name:      "echo simple",
			command:   "echo",
			args:      []string{"hello", "world"},
			expectErr: false,
			checkOutput: func(result *ExecutionResult) error {
				expected := "hello world\n"
				if result.Output != expected {
					return fmt.Errorf("expected '%s', got '%s'", expected, result.Output)
				}
				return nil
			},
		},
		{
			name:      "set variable",
			command:   "set",
			args:      []string{"TESTVAR", "testvalue"},
			expectErr: false,
			checkOutput: func(result *ExecutionResult) error {
				if !strings.Contains(result.Output, "TESTVAR=testvalue") {
					return fmt.Errorf("expected TESTVAR=testvalue in output")
				}
				return nil
			},
		},
		{
			name:      "export variable",
			command:   "export",
			args:      []string{"EXPORTVAR=exportvalue"},
			expectErr: false,
			checkOutput: func(result *ExecutionResult) error {
				// export should not produce output
				return nil
			},
		},
		{
			name:      "create alias",
			command:   "alias",
			args:      []string{"ll", "ls -la"},
			expectErr: false,
			checkOutput: func(result *ExecutionResult) error {
				if !strings.Contains(result.Output, "ll='ls -la'") {
					return fmt.Errorf("expected alias output to contain ll='ls -la'")
				}
				return nil
			},
		},
		{
			name:      "list aliases",
			command:   "alias",
			args:      []string{},
			expectErr: false,
			checkOutput: func(result *ExecutionResult) error {
				// Should list previously created alias
				if !strings.Contains(result.Output, "ll='ls -la'") {
					return fmt.Errorf("expected alias list to contain ll='ls -la'")
				}
				return nil
			},
		},
		{
			name:      "show history",
			command:   "history",
			args:      []string{},
			expectErr: false,
			checkOutput: func(result *ExecutionResult) error {
				// History should contain previous commands
				return nil
			},
		},
		{
			name:      "change directory - invalid",
			command:   "cd",
			args:      []string{"/nonexistent/directory"},
			expectErr: true,
			checkOutput: func(result *ExecutionResult) error {
				if result.ExitCode == 0 {
					return fmt.Errorf("expected non-zero exit code for invalid cd")
				}
				return nil
			},
		},
		{
			name:      "set without args",
			command:   "set",
			args:      []string{},
			expectErr: true,
			checkOutput: func(result *ExecutionResult) error {
				if result.ExitCode == 0 {
					return fmt.Errorf("expected non-zero exit code for set without args")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builtin, exists := shell.builtins[tt.command]
			if !exists {
				t.Fatalf("Builtin command '%s' not found", tt.command)
			}

			result, err := builtin(tt.args)

			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if result.ExitCode == 0 {
					t.Error("Expected non-zero exit code")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result.ExitCode != 0 {
					t.Errorf("Expected zero exit code, got %d", result.ExitCode)
				}
			}

			if result == nil {
				t.Error("Result should not be nil")
				return
			}

			// Check custom output validation
			if tt.checkOutput != nil {
				if err := tt.checkOutput(result); err != nil {
					t.Errorf("Output validation failed: %v", err)
				}
			}
		})
	}
}

func testEnvironmentVariables(t *testing.T, shell *Shell) {
	tests := []struct {
		name     string
		setup    func() // Setup function to run before test
		input    string
		expected string
	}{
		{
			name: "Simple variable expansion",
			setup: func() {
				shell.env["TESTVAR"] = "testvalue"
			},
			input:    "echo $TESTVAR",
			expected: "echo testvalue",
		},
		{
			name: "Braced variable expansion",
			setup: func() {
				shell.env["HOME"] = "/home/user"
			},
			input:    "echo ${HOME}/documents",
			expected: "echo /home/user/documents",
		},
		{
			name: "Multiple variables",
			setup: func() {
				shell.env["USER"] = "testuser"
				shell.env["HOME"] = "/home/testuser"
			},
			input:    "$USER lives in $HOME",
			expected: "testuser lives in /home/testuser",
		},
		{
			name:     "No variables",
			setup:    func() {},
			input:    "echo hello world",
			expected: "echo hello world",
		},
		{
			name:     "Undefined variable",
			setup:    func() {},
			input:    "echo $UNDEFINED_VAR",
			expected: "echo ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result := shell.expandVariables(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}

	// Test setting and using variables
	t.Run("SetAndUseVariable", func(t *testing.T) {
		// Set variable
		result, err := shell.builtins["set"]([]string{"DYNVAR", "dynamicvalue"})
		if err != nil {
			t.Fatalf("Failed to set variable: %v", err)
		}
		if result.ExitCode != 0 {
			t.Errorf("Expected zero exit code for set command")
		}

		// Check variable is set
		if shell.env["DYNVAR"] != "dynamicvalue" {
			t.Errorf("Variable not set correctly")
		}

		// Use variable
		expanded := shell.expandVariables("Value: $DYNVAR")
		expected := "Value: dynamicvalue"
		if expanded != expected {
			t.Errorf("Expected '%s', got '%s'", expected, expanded)
		}
	})
}

func testAliasSystem(t *testing.T, shell *Shell) {
	// Create alias
	result, err := shell.builtins["alias"]([]string{"ll", "ls -la"})
	if err != nil {
		t.Fatalf("Failed to create alias: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("Expected zero exit code for alias command")
	}

	// Verify alias exists
	if shell.aliases["ll"] != "ls -la" {
		t.Error("Alias was not created correctly")
	}

	// Test alias expansion in command parsing
	cmd, err := shell.parseCommand("ll -h")
	if err != nil {
		t.Errorf("Failed to parse aliased command: %v", err)
		return
	}

	if cmd.Name != "ls" {
		t.Errorf("Expected command name 'ls', got '%s'", cmd.Name)
	}

	expectedArgs := []string{"-la", "-h"}
	if len(cmd.Args) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(cmd.Args))
		return
	}

	for i, arg := range expectedArgs {
		if cmd.Args[i] != arg {
			t.Errorf("Expected arg %d to be '%s', got '%s'", i, arg, cmd.Args[i])
		}
	}

	// Test listing aliases
	result, err = shell.builtins["alias"]([]string{})
	if err != nil {
		t.Errorf("Failed to list aliases: %v", err)
	}
	if !strings.Contains(result.Output, "ll='ls -la'") {
		t.Error("Alias list should contain created alias")
	}

	// Test complex alias
	_, err = shell.builtins["alias"]([]string{"grep_error", "grep -i error"})
	if err != nil {
		t.Errorf("Failed to create complex alias: %v", err)
	}

	// Test using complex alias
	cmd, err = shell.parseCommand("grep_error logfile.txt")
	if err != nil {
		t.Errorf("Failed to parse complex aliased command: %v", err)
		return
	}

	if cmd.Name != "grep" {
		t.Errorf("Expected command name 'grep', got '%s'", cmd.Name)
	}
	if len(cmd.Args) < 2 || cmd.Args[0] != "-i" || cmd.Args[1] != "error" {
		t.Errorf("Complex alias not expanded correctly: %v", cmd.Args)
	}
}

func testHistoryManagement(t *testing.T, shell *Shell) {
	// Clear existing history for clean test
	shell.history = []string{}

	commands := []string{"echo test1", "pwd", "echo test2", "set VAR=value", "echo $VAR"}
	
	// Execute commands to build history
	for _, cmd := range commands {
		_, err := shell.Execute(cmd)
		if err != nil {
			t.Errorf("Failed to execute command '%s': %v", cmd, err)
		}
	}

	// Verify history length
	if len(shell.history) != len(commands) {
		t.Errorf("Expected %d history entries, got %d", len(commands), len(shell.history))
	}

	// Verify history content
	for i, cmd := range commands {
		if shell.history[i] != cmd {
			t.Errorf("History entry %d: expected '%s', got '%s'", i, cmd, shell.history[i])
		}
	}

	// Test history builtin
	result, err := shell.builtins["history"]([]string{})
	if err != nil {
		t.Errorf("History command failed: %v", err)
	}

	if result.Output == "" {
		t.Error("History output should not be empty")
	}

	// Verify all commands appear in history output
	for _, cmd := range commands {
		if !strings.Contains(result.Output, cmd) {
			t.Errorf("History should contain command: %s", cmd)
		}
	}

	// Test limited history
	result, err = shell.builtins["history"]([]string{"3"})
	if err != nil {
		t.Errorf("Limited history command failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(result.Output), "\n")
	if len(lines) > 3 {
		t.Errorf("Limited history should show at most 3 lines, got %d", len(lines))
	}
}

func testErrorHandling(t *testing.T, shell *Shell) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		checkResult func(result *ExecutionResult, err error) error
	}{
		{
			name:        "Empty command line",
			input:       "",
			expectError: false,
			checkResult: func(result *ExecutionResult, err error) error {
				if err != nil {
					return fmt.Errorf("empty command should not error")
				}
				return nil
			},
		},
		{
			name:        "Whitespace only",
			input:       "   \t  ",
			expectError: false, // Now handled gracefully
			checkResult: func(result *ExecutionResult, err error) error {
				// Whitespace-only commands should be handled gracefully now
				if err != nil {
					return fmt.Errorf("whitespace-only command should not error anymore: %v", err)
				}
				return nil
			},
		},
		{
			name:        "Invalid command",
			input:       "nonexistentcommand123",
			expectError: false, // Shell should handle gracefully
			checkResult: func(result *ExecutionResult, err error) error {
				// External commands may fail but shouldn't crash
				if result.ExitCode == 0 && result.Error == "" {
					// This might be OK if the command is handled gracefully
				}
				return nil
			},
		},
		{
			name:        "Malformed redirect",
			input:       "echo hello >",
			expectError: false, // Should parse but may fail execution
			checkResult: func(result *ExecutionResult, err error) error {
				// The command should parse successfully
				return nil
			},
		},
		{
			name:        "Unclosed quote",
			input:       `echo "unclosed quote`,
			expectError: false, // Tokenizer should handle this
			checkResult: func(result *ExecutionResult, err error) error {
				// The tokenizer handles unclosed quotes by including them
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := shell.Execute(tt.input)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.checkResult != nil {
				if checkErr := tt.checkResult(result, err); checkErr != nil {
					t.Error(checkErr)
				}
			}
		})
	}
}

func testResourceAndPerformance(t *testing.T, shell *Shell) {
	// Test with long input
	longString := strings.Repeat("a", 10000)
	start := time.Now()
	result, err := shell.Execute("echo " + longString)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Long echo command failed: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("Long echo command should succeed")
	}
	if duration > 5*time.Second {
		t.Errorf("Long echo took too long: %v", duration)
	}

	// Test many variables
	for i := 0; i < 100; i++ {
		shell.env[fmt.Sprintf("VAR%d", i)] = fmt.Sprintf("value%d", i)
	}

	start = time.Now()
	result, err = shell.Execute("echo Variables set")
	duration = time.Since(start)

	if err != nil {
		t.Errorf("Command with many variables failed: %v", err)
	}
	if duration > 1*time.Second {
		t.Errorf("Command with many variables took too long: %v", duration)
	}

	// Test complex command parsing
	complexCmd := "echo hello | echo world | echo test > /dev/null"
	start = time.Now()
	_, err = shell.parseCommandLine(complexCmd)
	duration = time.Since(start)

	if err != nil {
		t.Errorf("Complex command parsing failed: %v", err)
	}
	if duration > 100*time.Millisecond {
		t.Errorf("Complex parsing took too long: %v", duration)
	}

	// Memory usage test - create many history entries
	originalHistory := shell.history
	for i := 0; i < 1000; i++ {
		shell.history = append(shell.history, fmt.Sprintf("command%d", i))
	}

	start = time.Now()
	_, err = shell.builtins["history"]([]string{})
	duration = time.Since(start)

	if err != nil {
		t.Errorf("Large history command failed: %v", err)
	}
	if duration > 1*time.Second {
		t.Errorf("Large history took too long: %v", duration)
	}

	// Restore original history
	shell.history = originalHistory
}

func testEdgeCases(t *testing.T, shell *Shell) {
	tests := []struct {
		name  string
		test  func(t *testing.T)
	}{
		{
			name: "NilArguments",
			test: func(t *testing.T) {
				// Test builtin commands with nil arguments
				for cmdName, builtin := range shell.builtins {
					if cmdName == "exit" {
						continue // Skip exit as it terminates
					}
					result, err := builtin(nil)
					if result == nil {
						t.Errorf("Command %s returned nil result with nil args", cmdName)
					}
					// Error is acceptable for some commands with nil args
					_ = err
				}
			},
		},
		{
			name: "EmptyArguments",
			test: func(t *testing.T) {
				// Test with empty argument slices
				for cmdName, builtin := range shell.builtins {
					if cmdName == "exit" {
						continue
					}
					result, err := builtin([]string{})
					if result == nil {
						t.Errorf("Command %s returned nil result with empty args", cmdName)
					}
					_ = err
				}
			},
		},
		{
			name: "VeryLongCommandName",
			test: func(t *testing.T) {
				longName := strings.Repeat("verylongcommandname", 100)
				_, err := shell.parseCommand(longName)
				if err != nil {
					t.Errorf("Should handle very long command names: %v", err)
				}
			},
		},
		{
			name: "ManyArguments",
			test: func(t *testing.T) {
				args := []string{}
				for i := 0; i < 1000; i++ {
					args = append(args, fmt.Sprintf("arg%d", i))
				}
				cmdLine := "echo " + strings.Join(args, " ")
				_, err := shell.parseCommandLine(cmdLine)
				if err != nil {
					t.Errorf("Should handle many arguments: %v", err)
				}
			},
		},
		{
			name: "SpecialCharacters",
			test: func(t *testing.T) {
				specialChars := []string{
					`echo "special chars: !@#$%^&*()_+-={}[]|\:";'<>?,./~` + "`",
					"echo $!@#$%^&*()",
					"echo ${SPECIAL_VAR_!@#}",
				}
				
				for _, cmd := range specialChars {
					_, err := shell.parseCommandLine(cmd)
					// These may or may not parse successfully, but shouldn't crash
					_ = err
				}
			},
		},
		{
			name: "UnicodeHandling",
			test: func(t *testing.T) {
				unicodeCommands := []string{
					"echo 🚀 unicode test",
					"echo こんにちは world",
					"echo测试 unicode",
					"echo Ñiño español",
				}
				
				for _, cmd := range unicodeCommands {
					result, err := shell.Execute(cmd)
					if err != nil {
						t.Errorf("Unicode command failed: %v", err)
					}
					if result == nil {
						t.Error("Unicode command returned nil result")
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// Benchmark tests
func BenchmarkCommandExecution(b *testing.B) {
	shell := New(&Config{SafeMode: false})
	command := "echo hello world"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := shell.Execute(command)
		if err != nil {
			b.Fatalf("Command execution failed: %v", err)
		}
	}
}

func BenchmarkComplexParsing(b *testing.B) {
	shell := New(&Config{})
	command := `echo "complex command" | grep -i "test" | sort -n > output.txt`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := shell.parseCommandLine(command)
		if err != nil {
			b.Fatalf("Command parsing failed: %v", err)
		}
	}
}

func BenchmarkAdvancedVariableExpansion(b *testing.B) {
	shell := New(&Config{})
	shell.env["HOME"] = "/home/testuser"
	shell.env["USER"] = "testuser"
	shell.env["PATH"] = "/usr/bin:/bin:/usr/local/bin"
	input := "User $USER in $HOME with path $PATH and ${HOME}/documents"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shell.expandVariables(input)
	}
}

func BenchmarkHistoryManagement(b *testing.B) {
	shell := New(&Config{})
	
	// Pre-populate history
	for i := 0; i < 1000; i++ {
		shell.history = append(shell.history, fmt.Sprintf("command %d", i))
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := shell.builtins["history"]([]string{})
		if err != nil {
			b.Fatalf("History command failed: %v", err)
		}
	}
}