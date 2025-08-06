package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ComprehensiveShellTester performs systematic testing of Lucien shell
type ComprehensiveShellTester struct {
	binaryPath   string
	testDir      string
	logFile      *os.File
	testResults  []TestResult
}

type TestResult struct {
	Name     string
	Passed   bool
	Output   string
	Error    string
	Duration time.Duration
}

// NewTester creates a comprehensive shell tester
func NewTester(binaryPath string) (*ComprehensiveShellTester, error) {
	testDir := filepath.Join(os.TempDir(), "lucien_test_"+time.Now().Format("20060102_150405"))
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create test directory: %v", err)
	}

	logFile, err := os.Create(filepath.Join(testDir, "test_log.txt"))
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %v", err)
	}

	return &ComprehensiveShellTester{
		binaryPath:  binaryPath,
		testDir:     testDir,
		logFile:     logFile,
		testResults: []TestResult{},
	}, nil
}

// Cleanup removes test files and closes resources
func (t *ComprehensiveShellTester) Cleanup() {
	if t.logFile != nil {
		t.logFile.Close()
	}
	os.RemoveAll(t.testDir)
}

// RunAllTests executes comprehensive shell testing
func (t *ComprehensiveShellTester) RunAllTests() {
	fmt.Println("=== LUCIEN SHELL COMPREHENSIVE TESTING ===")
	t.logResult("Starting comprehensive shell tests", "", "")

	// Test 1: Basic Shell Commands and Built-ins
	t.testBasicCommands()
	
	// Test 2: Command parsing with pipes, redirects, and variables
	t.testCommandParsing()
	
	// Test 3: Environment variable handling
	t.testEnvironmentVariables()
	
	// Test 4: Alias system
	t.testAliasSystem()
	
	// Test 5: History management
	t.testHistoryManagement()
	
	// Test 6: Shell scripting capabilities
	t.testShellScripting()
	
	// Test 7: Error handling and edge cases
	t.testErrorHandling()
	
	// Test 8: Resource usage and performance
	t.testResourceUsage()
	
	// Test 9: Concurrent operations
	t.testConcurrentOperations()
	
	// Test 10: Plugin system
	t.testPluginSystem()
	
	// Generate final report
	t.generateReport()
}

// Test basic shell commands and built-ins
func (t *ComprehensiveShellTester) testBasicCommands() {
	fmt.Println("Testing basic shell commands and built-ins...")
	
	commands := []struct {
		name    string
		command string
		expectOutput bool
		expectError  bool
	}{
		{"pwd", "pwd", true, false},
		{"echo_simple", "echo hello world", true, false},
		{"echo_quoted", `echo "hello world"`, true, false},
		{"set_variable", "set TEST_VAR test_value", true, false},
		{"export_variable", "export TEST_VAR=exported_value", false, false},
		{"alias_create", "alias ll='ls -la'", true, false},
		{"alias_list", "alias", true, false},
		{"history_empty", "history", true, false},
	}

	for _, test := range commands {
		t.runSingleCommand(test.name, test.command, test.expectOutput, test.expectError)
	}
}

// Test command parsing features
func (t *ComprehensiveShellTester) testCommandParsing() {
	fmt.Println("Testing command parsing with pipes, redirects, and variables...")
	
	// Create test files
	testInputFile := filepath.Join(t.testDir, "input.txt")
	ioutil.WriteFile(testInputFile, []byte("line1\nline2\nline3\n"), 0644)
	
	commands := []struct {
		name    string
		command string
		expectOutput bool
		expectError  bool
	}{
		{"redirect_output", fmt.Sprintf("echo 'test output' > %s/output.txt", t.testDir), false, false},
		{"redirect_input", fmt.Sprintf("echo hello < %s/input.txt", testInputFile), true, false},
		{"pipe_simple", "echo hello | echo world", true, false},
		{"variable_expansion", "set MYVAR=test && echo $MYVAR", true, false},
		{"quoted_arguments", `echo "argument with spaces"`, true, false},
		{"multiple_redirects", fmt.Sprintf("echo test < %s/input.txt > %s/output2.txt", testInputFile, t.testDir), false, false},
	}

	for _, test := range commands {
		t.runSingleCommand(test.name, test.command, test.expectOutput, test.expectError)
	}
}

// Test environment variable handling
func (t *ComprehensiveShellTester) testEnvironmentVariables() {
	fmt.Println("Testing environment variable handling...")
	
	commands := []struct {
		name    string
		command string
		expectOutput bool
		expectError  bool
	}{
		{"set_env_var", "set TESTVAR=testvalue", true, false},
		{"export_env_var", "export EXPORTVAR=exportvalue", false, false},
		{"use_env_var", "echo $TESTVAR", true, false},
		{"use_braced_var", "echo ${TESTVAR}_suffix", true, false},
		{"list_exports", "export", true, false},
	}

	for _, test := range commands {
		t.runSingleCommand(test.name, test.command, test.expectOutput, test.expectError)
	}
}

// Test alias system functionality
func (t *ComprehensiveShellTester) testAliasSystem() {
	fmt.Println("Testing alias system...")
	
	commands := []struct {
		name    string
		command string
		expectOutput bool
		expectError  bool
	}{
		{"create_alias", "alias ll='ls -la'", true, false},
		{"create_complex_alias", `alias grep_error="grep -i error"`, true, false},
		{"list_aliases", "alias", true, false},
		{"use_alias", "ll", true, false}, // This might fail if ls is not available, but should still test alias expansion
	}

	for _, test := range commands {
		t.runSingleCommand(test.name, test.command, test.expectOutput, test.expectError)
	}
}

// Test command history management
func (t *ComprehensiveShellTester) testHistoryManagement() {
	fmt.Println("Testing command history...")
	
	// Execute several commands to build history
	t.runSingleCommand("history_build_1", "echo command1", true, false)
	t.runSingleCommand("history_build_2", "echo command2", true, false)
	t.runSingleCommand("history_build_3", "echo command3", true, false)
	
	// Test history commands
	commands := []struct {
		name    string
		command string
		expectOutput bool
		expectError  bool
	}{
		{"show_full_history", "history", true, false},
		{"show_last_5", "history 5", true, false},
	}

	for _, test := range commands {
		t.runSingleCommand(test.name, test.command, test.expectOutput, test.expectError)
	}
}

// Test shell scripting capabilities
func (t *ComprehensiveShellTester) testShellScripting() {
	fmt.Println("Testing shell scripting capabilities...")
	
	// Create a simple script
	scriptFile := filepath.Join(t.testDir, "test_script.sh")
	scriptContent := `#!/bin/bash
echo "Script start"
set SCRIPT_VAR=script_value
echo "Variable: $SCRIPT_VAR"
echo "Script end"
`
	ioutil.WriteFile(scriptFile, []byte(scriptContent), 0755)
	
	commands := []struct {
		name    string
		command string
		expectOutput bool
		expectError  bool
	}{
		{"multi_command", "echo first && echo second", true, false},
		{"variable_in_command", "set VAR=value && echo $VAR", true, false},
	}

	for _, test := range commands {
		t.runSingleCommand(test.name, test.command, test.expectOutput, test.expectError)
	}
}

// Test error handling and edge cases
func (t *ComprehensiveShellTester) testErrorHandling() {
	fmt.Println("Testing error handling and edge cases...")
	
	commands := []struct {
		name    string
		command string
		expectOutput bool
		expectError  bool
	}{
		{"invalid_command", "nonexistent_command", false, true},
		{"empty_command", "", false, false},
		{"only_spaces", "   ", false, false},
		{"invalid_redirect", "echo test > /invalid/path/file.txt", false, true},
		{"missing_quote", `echo "unterminated quote`, false, true},
		{"invalid_variable", "echo $NONEXISTENT_VAR", true, false}, // Should output empty
		{"cd_nonexistent", "cd /nonexistent/directory", false, true},
		{"set_without_args", "set", false, true},
		{"alias_without_args", "alias invalidalias", true, false}, // Should show empty or error gracefully
	}

	for _, test := range commands {
		t.runSingleCommand(test.name, test.command, test.expectOutput, test.expectError)
	}
}

// Test resource usage and performance
func (t *ComprehensiveShellTester) testResourceUsage() {
	fmt.Println("Testing resource usage and performance...")
	
	// Test commands that might stress the system
	commands := []struct {
		name    string
		command string
		expectOutput bool
		expectError  bool
	}{
		{"long_echo", "echo " + strings.Repeat("a", 1000), true, false},
		{"many_variables", "set VAR1=1 && set VAR2=2 && set VAR3=3 && set VAR4=4 && set VAR5=5", true, false},
		{"complex_pipe", "echo test | echo pipe1 | echo pipe2", true, false},
	}

	for _, test := range commands {
		start := time.Now()
		t.runSingleCommand(test.name, test.command, test.expectOutput, test.expectError)
		duration := time.Since(start)
		t.logResult(test.name+"_performance", fmt.Sprintf("Duration: %v", duration), "")
		
		// Flag performance issues
		if duration > 5*time.Second {
			t.logResult(test.name+"_PERFORMANCE_WARNING", fmt.Sprintf("Command took %v - may be too slow", duration), "")
		}
	}
}

// Test concurrent operations
func (t *ComprehensiveShellTester) testConcurrentOperations() {
	fmt.Println("Testing concurrent operations...")
	
	// This is limited since we're testing via CLI, but we can test rapid command execution
	commands := []string{
		"echo concurrent1",
		"echo concurrent2", 
		"echo concurrent3",
		"set CONC_VAR=concurrent",
		"echo $CONC_VAR",
	}

	start := time.Now()
	for i, cmd := range commands {
		t.runSingleCommand(fmt.Sprintf("concurrent_%d", i), cmd, true, false)
	}
	totalDuration := time.Since(start)
	
	t.logResult("concurrent_operations_total", fmt.Sprintf("Total time for %d commands: %v", len(commands), totalDuration), "")
}

// Test plugin loading and execution
func (t *ComprehensiveShellTester) testPluginSystem() {
	fmt.Println("Testing plugin system...")
	
	// Check if example plugin exists
	pluginPath := filepath.Join(filepath.Dir(t.binaryPath), "..", "plugins", "example-bmad", "example-bmad.exe")
	if _, err := os.Stat(pluginPath); err == nil {
		t.logResult("plugin_exists", "Example BMAD plugin found", "")
		
		// Try to test plugin functionality if there's a way to invoke it
		commands := []struct {
			name    string
			command string
			expectOutput bool
			expectError  bool
		}{
			// Note: We'd need to know how to invoke plugins through the shell
			// This is a placeholder - actual plugin invocation depends on shell design
		}

		for _, test := range commands {
			t.runSingleCommand(test.name, test.command, test.expectOutput, test.expectError)
		}
	} else {
		t.logResult("plugin_missing", "Example BMAD plugin not found", "Plugin system cannot be fully tested")
	}
}

// Run a single command and record results
func (t *ComprehensiveShellTester) runSingleCommand(name, command string, expectOutput, expectError bool) {
	start := time.Now()
	
	// Create command to run lucien shell with the test command
	cmd := exec.Command(t.binaryPath)
	cmd.Dir = t.testDir
	
	// Send command via stdin
	cmd.Stdin = strings.NewReader(command + "\nexit\n")
	
	output, err := cmd.CombinedOutput()
	duration := time.Since(start)
	
	result := TestResult{
		Name:     name,
		Passed:   true,
		Output:   string(output),
		Duration: duration,
	}
	
	if err != nil {
		result.Error = err.Error()
		if !expectError {
			result.Passed = false
			fmt.Printf("FAIL: %s - Unexpected error: %v\n", name, err)
		} else {
			fmt.Printf("PASS: %s - Expected error occurred\n", name)
		}
	} else {
		if expectError {
			result.Passed = false
			fmt.Printf("FAIL: %s - Expected error but command succeeded\n", name)
		} else {
			fmt.Printf("PASS: %s\n", name)
		}
	}
	
	if expectOutput && strings.TrimSpace(result.Output) == "" {
		result.Passed = false
		fmt.Printf("FAIL: %s - Expected output but got none\n", name)
	}
	
	t.testResults = append(t.testResults, result)
	t.logResult(name, result.Output, result.Error)
}

// Log test results
func (t *ComprehensiveShellTester) logResult(name, output, errorMsg string) {
	timestamp := time.Now().Format("15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, name)
	if output != "" {
		logEntry += fmt.Sprintf("  OUTPUT: %s\n", strings.ReplaceAll(output, "\n", "\\n"))
	}
	if errorMsg != "" {
		logEntry += fmt.Sprintf("  ERROR: %s\n", errorMsg)
	}
	logEntry += "\n"
	
	t.logFile.WriteString(logEntry)
	t.logFile.Sync()
}

// Generate final test report
func (t *ComprehensiveShellTester) generateReport() {
	fmt.Println("\n=== FINAL TEST REPORT ===")
	
	totalTests := len(t.testResults)
	passedTests := 0
	failedTests := 0
	
	for _, result := range t.testResults {
		if result.Passed {
			passedTests++
		} else {
			failedTests++
		}
	}
	
	fmt.Printf("Total Tests: %d\n", totalTests)
	fmt.Printf("Passed: %d\n", passedTests)
	fmt.Printf("Failed: %d\n", failedTests)
	fmt.Printf("Success Rate: %.2f%%\n", float64(passedTests)/float64(totalTests)*100)
	
	if failedTests > 0 {
		fmt.Println("\nFAILED TESTS:")
		for _, result := range t.testResults {
			if !result.Passed {
				fmt.Printf("- %s: %s\n", result.Name, result.Error)
			}
		}
	}
	
	// Write summary to log
	summary := fmt.Sprintf(`
=== TEST SUMMARY ===
Total: %d, Passed: %d, Failed: %d
Success Rate: %.2f%%
Test Duration: %v
Log File: %s
`, totalTests, passedTests, failedTests, 
		float64(passedTests)/float64(totalTests)*100,
		time.Since(time.Now()),
		t.logFile.Name())
	
	t.logFile.WriteString(summary)
	
	fmt.Printf("\nFull test log saved to: %s\n", t.logFile.Name())
}

// Main test execution
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run comprehensive_shell_test.go <path_to_lucien_binary>")
		os.Exit(1)
	}
	
	binaryPath := os.Args[1]
	if _, err := os.Stat(binaryPath); err != nil {
		fmt.Printf("Error: Lucien binary not found at %s\n", binaryPath)
		os.Exit(1)
	}
	
	tester, err := NewTester(binaryPath)
	if err != nil {
		fmt.Printf("Error creating tester: %v\n", err)
		os.Exit(1)
	}
	defer tester.Cleanup()
	
	tester.RunAllTests()
}