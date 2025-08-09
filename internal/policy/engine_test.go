package policy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPolicyEngineCreation(t *testing.T) {
	tmpDir := os.TempDir()
	policyDir := filepath.Join(tmpDir, "test_policies", time.Now().Format("20060102_150405"))
	defer os.RemoveAll(policyDir)

	engine, err := New(policyDir)
	if err != nil {
		t.Fatalf("Failed to create policy engine: %v", err)
	}

	if engine.policyDir != policyDir {
		t.Errorf("Expected policy directory %s, got %s", policyDir, engine.policyDir)
	}

	// Verify default policies were created
	expectedPolicies := []string{
		"deny_write_root.rego",
		"plugin_sandbox.rego", 
		"ai_safety.rego",
	}

	for _, policy := range expectedPolicies {
		policyPath := filepath.Join(policyDir, policy)
		if _, err := os.Stat(policyPath); os.IsNotExist(err) {
			t.Errorf("Default policy %s was not created", policy)
		}
	}
}

func TestCommandAuthorization(t *testing.T) {
	tmpDir := os.TempDir()
	policyDir := filepath.Join(tmpDir, "test_auth", time.Now().Format("20060102_150405"))
	defer os.RemoveAll(policyDir)

	engine, err := New(policyDir)
	if err != nil {
		t.Fatalf("Failed to create policy engine: %v", err)
	}

	tests := []struct {
		name        string
		action      string
		command     string
		args        []string
		shouldAllow bool
	}{
		{
			name:        "Safe echo command",
			action:      "execute",
			command:     "echo",
			args:        []string{"hello"},
			shouldAllow: true,
		},
		{
			name:        "Safe ls command",
			action:      "execute",
			command:     "ls",
			args:        []string{"-la"},
			shouldAllow: true,
		},
		{
			name:        "Dangerous rm command on root",
			action:      "execute", 
			command:     "rm",
			args:        []string{"-rf", "/"},
			shouldAllow: false,
		},
		{
			name:        "Dangerous dd command",
			action:      "execute",
			command:     "dd",
			args:        []string{"if=/dev/zero", "of=/dev/sda"},
			shouldAllow: false,
		},
		{
			name:        "Safe rm in home directory",
			action:      "execute",
			command:     "rm",
			args:        []string{"test.txt"},
			shouldAllow: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, err := engine.Authorize(tt.action, tt.command, tt.args)
			
			if err != nil && tt.shouldAllow {
				t.Errorf("Unexpected error for safe command: %v", err)
			}

			if allowed != tt.shouldAllow {
				t.Errorf("Expected authorization %v, got %v", tt.shouldAllow, allowed)
			}
		})
	}
}

func TestPluginAuthorization(t *testing.T) {
	tmpDir := os.TempDir()
	policyDir := filepath.Join(tmpDir, "test_plugin", time.Now().Format("20060102_150405"))
	defer os.RemoveAll(policyDir)

	engine, err := New(policyDir)
	if err != nil {
		t.Fatalf("Failed to create policy engine: %v", err)
	}

	homeDir, _ := os.UserHomeDir()
	
	tests := []struct {
		name        string
		pluginName  string
		action      string
		resource    string
		shouldAllow bool
	}{
		{
			name:        "Plugin read from home",
			pluginName:  "test-plugin",
			action:      "read",
			resource:    homeDir + "/test.txt",
			shouldAllow: true,
		},
		{
			name:        "Plugin write to plugin dir",
			pluginName:  "test-plugin", 
			action:      "write",
			resource:    homeDir + "/.lucien/plugins/test-plugin/data.txt",
			shouldAllow: true,
		},
		{
			name:        "Plugin read from root",
			pluginName:  "test-plugin",
			action:      "read",
			resource:    "/etc/passwd",
			shouldAllow: false,
		},
		{
			name:        "Plugin write outside sandbox",
			pluginName:  "test-plugin",
			action:      "write", 
			resource:    "/etc/hosts",
			shouldAllow: false,
		},
		{
			name:        "Plugin execute approved command",
			pluginName:  "test-plugin",
			action:      "execute",
			resource:    "echo",
			shouldAllow: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, err := engine.AuthorizePlugin(tt.pluginName, tt.action, tt.resource)
			
			if err != nil && tt.shouldAllow {
				t.Errorf("Unexpected error: %v", err)
			}

			if allowed != tt.shouldAllow {
				t.Errorf("Expected plugin authorization %v, got %v", tt.shouldAllow, allowed)
			}
		})
	}
}

func TestPolicyReloading(t *testing.T) {
	tmpDir := os.TempDir()
	policyDir := filepath.Join(tmpDir, "test_reload", time.Now().Format("20060102_150405"))
	defer os.RemoveAll(policyDir)

	engine, err := New(policyDir)
	if err != nil {
		t.Fatalf("Failed to create policy engine: %v", err)
	}

	initialPolicies := engine.ListPolicies()
	
	// Add a new policy file
	newPolicy := `
package lucien.test

default allow = true

deny {
    input.action == "test"
    input.command == "blocked"
}
`

	newPolicyPath := filepath.Join(policyDir, "test_policy.rego")
	if err := os.WriteFile(newPolicyPath, []byte(newPolicy), 0644); err != nil {
		t.Fatalf("Failed to write test policy: %v", err)
	}

	// Reload policies
	if err := engine.ReloadPolicies(); err != nil {
		t.Errorf("Failed to reload policies: %v", err)
	}

	reloadedPolicies := engine.ListPolicies()
	
	if len(reloadedPolicies) <= len(initialPolicies) {
		t.Error("New policy was not loaded after reload")
	}

	// Verify new policy is in the list
	found := false
	for _, policy := range reloadedPolicies {
		if policy == "test_policy" {
			found = true
			break
		}
	}

	if !found {
		t.Error("New test policy not found in loaded policies")
	}
}

func TestPolicyContent(t *testing.T) {
	tmpDir := os.TempDir()
	policyDir := filepath.Join(tmpDir, "test_content", time.Now().Format("20060102_150405"))
	defer os.RemoveAll(policyDir)

	engine, err := New(policyDir)
	if err != nil {
		t.Fatalf("Failed to create policy engine: %v", err)
	}

	// Test getting content of default policy
	content, err := engine.GetPolicyContent("deny_write_root")
	if err != nil {
		t.Errorf("Failed to get policy content: %v", err)
	}

	if content == "" {
		t.Error("Policy content should not be empty")
	}

	// Verify content contains expected elements
	expectedElements := []string{
		"package lucien.security",
		"deny_write_root",
		"dangerous_commands",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Policy content should contain '%s'", element)
		}
	}
}

func TestPolicyWithEnvironmentVariables(t *testing.T) {
	tmpDir := os.TempDir()
	policyDir := filepath.Join(tmpDir, "test_env", time.Now().Format("20060102_150405"))
	defer os.RemoveAll(policyDir)

	// Set test environment
	originalHome := os.Getenv("HOME")
	testHome := "/test/home/user"
	os.Setenv("HOME", testHome)
	defer os.Setenv("HOME", originalHome)

	engine, err := New(policyDir)
	if err != nil {
		t.Fatalf("Failed to create policy engine: %v", err)
	}

	// Test command that should be allowed in home directory
	allowed, err := engine.Authorize("execute", "rm", []string{testHome + "/test.txt"})
	if err != nil {
		t.Errorf("Error checking home directory access: %v", err)
	}

	if !allowed {
		t.Error("Should allow rm in home directory")
	}

	// Test command that should be denied outside home
	denied, err := engine.Authorize("execute", "rm", []string{"/etc/passwd"})
	if err == nil {
		t.Error("Expected error for dangerous command")
	}

	if denied {
		t.Error("Should deny rm outside home directory")
	}
}

func TestCustomPolicyIntegration(t *testing.T) {
	tmpDir := os.TempDir()
	policyDir := filepath.Join(tmpDir, "test_custom", time.Now().Format("20060102_150405"))
	defer os.RemoveAll(policyDir)

	engine, err := New(policyDir)
	if err != nil {
		t.Fatalf("Failed to create policy engine: %v", err)
	}

	// Create custom policy
	customPolicy := `
package lucien.custom

# Block commands containing specific patterns
deny_custom {
    input.action == "execute"
    contains(input.command, "secret")
}

deny_custom {
    input.action == "execute"
    input.command == "custom_blocked"
}
`

	customPolicyPath := filepath.Join(policyDir, "custom_rules.rego")
	if err := os.WriteFile(customPolicyPath, []byte(customPolicy), 0644); err != nil {
		t.Fatalf("Failed to write custom policy: %v", err)
	}

	// Reload to pick up custom policy
	if err := engine.ReloadPolicies(); err != nil {
		t.Errorf("Failed to reload policies: %v", err)
	}

	// Test custom policy enforcement
	tests := []struct {
		name        string
		command     string
		args        []string
		shouldBlock bool
	}{
		{
			name:        "Normal command",
			command:     "echo",
			args:        []string{"hello"},
			shouldBlock: false,
		},
		{
			name:        "Blocked custom command",
			command:     "custom_blocked",
			args:        []string{},
			shouldBlock: true,
		},
		{
			name:        "Command with secret pattern",
			command:     "echo",
			args:        []string{"secret_data"},
			shouldBlock: false, // Pattern matching on command name, not args
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, err := engine.Authorize("execute", tt.command, tt.args)
			
			if tt.shouldBlock && (err == nil && allowed) {
				t.Error("Expected command to be blocked by custom policy")
			}
			
			if !tt.shouldBlock && err != nil {
				t.Errorf("Unexpected error for allowed command: %v", err)
			}
		})
	}
}

// Benchmark policy evaluation performance
func BenchmarkPolicyAuthorization(b *testing.B) {
	tmpDir := os.TempDir()
	policyDir := filepath.Join(tmpDir, "bench_policy", time.Now().Format("20060102_150405"))
	defer os.RemoveAll(policyDir)

	engine, err := New(policyDir)
	if err != nil {
		b.Fatalf("Failed to create policy engine: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Authorize("execute", "echo", []string{"hello"})
	}
}

func BenchmarkPluginAuthorization(b *testing.B) {
	tmpDir := os.TempDir()
	policyDir := filepath.Join(tmpDir, "bench_plugin", time.Now().Format("20060102_150405"))
	defer os.RemoveAll(policyDir)

	engine, err := New(policyDir)
	if err != nil {
		b.Fatalf("Failed to create policy engine: %v", err)
	}

	homeDir, _ := os.UserHomeDir()
	resource := homeDir + "/test.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.AuthorizePlugin("test-plugin", "read", resource)
	}
}

