package shell

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestChangeDirectory(t *testing.T) {
	// Save original directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	shell := New(&Config{SafeMode: false, ExecutorMode: "internal"})

	tests := []struct {
		name        string
		args        []string
		expectError bool
		skipOnWin   bool
		skipOnUnix  bool
	}{
		{
			name:        "cd to home directory (no args)",
			args:        []string{},
			expectError: false,
		},
		{
			name:        "cd to root directory",
			args:        []string{"/"},
			expectError: false,
			skipOnWin:   true,
		},
		{
			name:        "cd to current directory",
			args:        []string{"."},
			expectError: false,
		},
		{
			name:        "cd to parent directory",
			args:        []string{".."},
			expectError: false,
		},
		{
			name:        "cd to tilde (home)",
			args:        []string{"~"},
			expectError: false,
		},
		{
			name:        "cd to nonexistent directory",
			args:        []string{"/nonexistent/directory/path"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip platform-specific tests
			if tt.skipOnWin && runtime.GOOS == "windows" {
				t.Skip("Skipping test on Windows")
			}
			if tt.skipOnUnix && runtime.GOOS != "windows" {
				t.Skip("Skipping test on Unix")
			}

			result, err := shell.changeDirectory(tt.args)
			
			if tt.expectError {
				if err == nil && result.ExitCode == 0 {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil || result.ExitCode != 0 {
					t.Errorf("Unexpected error: %v, result: %+v", err, result)
				}
			}
		})
	}
}

func TestChangeDirectoryWithQuotedPaths(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Quoted path test is primarily for Windows")
	}

	// Save original directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	shell := New(&Config{SafeMode: false, ExecutorMode: "internal"})

	// Create a test directory with spaces
	testDir := filepath.Join(os.TempDir(), "test directory with spaces")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "quoted path with spaces",
			args: []string{`"` + testDir + `"`},
		},
		{
			name: "single quoted path with spaces",
			args: []string{`'` + testDir + `'`},
		},
		{
			name: "unquoted path with spaces (should still work due to shell parsing)",
			args: []string{testDir},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := shell.changeDirectory(tt.args)
			
			if err != nil || result.ExitCode != 0 {
				t.Errorf("Failed to change to quoted directory: %v, result: %+v", err, result)
				return
			}

			// Verify we're in the right directory
			currentDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}

			// Clean both paths for comparison
			expectedDir := filepath.Clean(testDir)
			actualDir := filepath.Clean(currentDir)

			if !strings.EqualFold(expectedDir, actualDir) {
				t.Errorf("Expected to be in %q, but in %q", expectedDir, actualDir)
			}
		})
	}
}

func TestFirstArgAsPath(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedPath string
		expectError  bool
	}{
		{
			name:         "no arguments",
			args:         []string{},
			expectedPath: "",
			expectError:  true,
		},
		{
			name:         "simple path",
			args:         []string{"folder"},
			expectedPath: "folder",
			expectError:  false,
		},
		{
			name:         "double quoted path",
			args:         []string{`"C:\Program Files\Test"`},
			expectedPath: `C:\Program Files\Test`,
			expectError:  false,
		},
		{
			name:         "single quoted path",
			args:         []string{`'/usr/local/bin'`},
			expectedPath: `/usr/local/bin`,
			expectError:  false,
		},
		{
			name:         "path with both quote types",
			args:         []string{`"'mixed quotes'"`},
			expectedPath: `'mixed quotes'`,
			expectError:  false,
		},
		{
			name:         "empty quoted path",
			args:         []string{`""`},
			expectedPath: "",
			expectError:  true,
		},
		{
			name:         "path with internal quotes",
			args:         []string{`"path with \"internal\" quotes"`},
			expectedPath: `path with \"internal\" quotes`,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := FirstArgAsPath(tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if path != tt.expectedPath {
					t.Errorf("Expected path %q, got %q", tt.expectedPath, path)
				}
			}
		})
	}
}