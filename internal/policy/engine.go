package policy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Engine represents a simplified policy engine for Lucien
type Engine struct {
	policyDir   string
	rules       map[string]*PolicyRule
	policyFiles map[string]string // Track .rego files on disk
}

// PolicyRule represents a simple policy rule
type PolicyRule struct {
	Name         string
	Description  string
	Action       string   // "allow" or "deny"
	Commands     []string // commands to match
	Args         []string // argument patterns to match
	Conditions   []string // additional conditions
}

// Decision represents a policy decision result
type Decision struct {
	Allow  bool              `json:"allow"`
	Reason string            `json:"reason,omitempty"`
	Meta   map[string]string `json:"meta,omitempty"`
}

// New creates a new policy engine
func New(policyDir string) (*Engine, error) {
	if err := os.MkdirAll(policyDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create policy directory: %w", err)
	}

	engine := &Engine{
		policyDir:   policyDir,
		rules:       make(map[string]*PolicyRule),
		policyFiles: make(map[string]string),
	}

	// Initialize with default security rules
	engine.initializeDefaultRules()
	
	// Create default policy files if they don't exist
	if err := engine.createDefaultPolicyFiles(); err != nil {
		return nil, fmt.Errorf("failed to create default policy files: %w", err)
	}
	
	// Load any existing policy files
	if err := engine.ReloadPolicies(); err != nil {
		return nil, fmt.Errorf("failed to load policy files: %w", err)
	}

	return engine, nil
}

func (e *Engine) initializeDefaultRules() {
	// Security rules for dangerous commands
	e.rules["deny_dangerous"] = &PolicyRule{
		Name:        "deny_dangerous",
		Description: "Block dangerous system commands",
		Action:      "deny",
		Commands:    []string{"dd", "mkfs", "fdisk", "format", "sudo", "su"},
		Conditions:  []string{"system_critical"},
	}

	// Special rule for dangerous rm operations
	e.rules["deny_dangerous_rm"] = &PolicyRule{
		Name:        "deny_dangerous_rm",
		Description: "Block dangerous rm operations",
		Action:      "deny",
		Commands:    []string{"rm"},
		Args:        []string{"/etc", "/usr", "/var", "/sys", "/proc", "/boot", "/root"},
		Conditions:  []string{"dangerous_paths"},
	}

	// Separate rule for root filesystem access
	e.rules["deny_root_access"] = &PolicyRule{
		Name:        "deny_root_access", 
		Description: "Block access to root filesystem",
		Action:      "deny",
		Commands:    []string{"rm", "rmdir"},
		Args:        []string{"/"}, // Only exact "/" path
		Conditions:  []string{"root_access"},
	}

	// File system protection - only block actual system directories
	e.rules["protect_system_dirs"] = &PolicyRule{
		Name:        "protect_system_dirs",
		Description: "Protect system directories from modification",
		Action:      "deny",
		Commands:    []string{"rm", "rmdir", "chmod", "chown"},
		Args:        []string{"/etc", "/usr", "/var", "/sys", "/proc", "/boot", "/root"},
		Conditions:  []string{"system_path"},
	}

	// Plugin sandboxing
	e.rules["plugin_execute"] = &PolicyRule{
		Name:        "plugin_execute",
		Description: "Allow only safe commands for plugins",
		Action:      "allow",
		Commands:    []string{"echo", "cat", "grep", "sed", "awk", "curl", "wget", "ls", "pwd"},
		Conditions:  []string{"plugin_context"},
	}

	// AI safety - only applies in AI context (not implemented yet)
	// Commented out until context-aware authorization is implemented
	// e.rules["ai_safety"] = &PolicyRule{
	//	Name:        "ai_safety", 
	//	Description: "AI cannot execute dangerous commands",
	//	Action:      "deny",
	//	Commands:    []string{"rm", "rmdir", "dd", "mkfs", "sudo", "su", "chmod", "chown"},
	//	Conditions:  []string{"ai_context"},
	// }
}

// createDefaultPolicyFiles creates default .rego policy files
func (e *Engine) createDefaultPolicyFiles() error {
	defaultPolicies := map[string]string{
		"deny_write_root.rego": `
package lucien.security

# Protect system directories from modification
deny_write_root {
    input.action == "execute"
    dangerous_commands := ["dd", "mkfs", "fdisk", "format", "rm", "rmdir"]
    input.command in dangerous_commands
    
    # Check if targeting root filesystem
    startswith(input.args[_], "/")
}

# Block dangerous system commands entirely
deny_dangerous {
    input.action == "execute"
    dangerous_commands := ["dd", "mkfs", "fdisk", "format"]
    input.command in dangerous_commands
}

# Protect system directories from modification
deny_system_dirs {
    input.action == "execute"  
    modify_commands := ["rm", "rmdir", "chmod", "chown"]
    input.command in modify_commands
    
    system_dirs := ["/etc", "/usr", "/var", "/sys", "/proc", "/boot"]
    startswith(input.args[_], system_dirs[_])
}
`,
		"plugin_sandbox.rego": `
package lucien.plugins

# Allow only safe commands for plugins in sandbox
allow_plugin_commands {
    input.action == "execute"
    input.context == "plugin"
    
    safe_commands := ["echo", "cat", "grep", "sed", "awk", "curl", "wget", "ls", "pwd", "wc", "sort", "uniq", "head", "tail"]
    input.command in safe_commands
}

# Deny plugin access to sensitive directories
deny_plugin_sensitive_access {
    input.action == "read"
    input.context == "plugin"
    
    sensitive_dirs := ["/etc", "/proc", "/sys", "/root", "/home"]
    startswith(input.path, sensitive_dirs[_])
}

# Allow plugin read access to designated directories
allow_plugin_read_access {
    input.action == "read"
    input.context == "plugin"
    
    allowed_dirs := ["/tmp", "/var/tmp"]
    startswith(input.path, allowed_dirs[_])
}
`,
		"ai_safety.rego": `
package lucien.ai

# AI cannot execute dangerous commands
deny_ai_dangerous {
    input.action == "execute"
    input.context == "ai"
    
    dangerous_commands := ["rm", "rmdir", "dd", "mkfs", "sudo", "su", "chmod", "chown", "kill", "killall", "reboot", "shutdown"]
    input.command in dangerous_commands
}

# AI can only execute safe, read-only commands
allow_ai_safe_commands {
    input.action == "execute"
    input.context == "ai"
    
    safe_commands := ["echo", "cat", "grep", "ls", "pwd", "wc", "sort", "uniq", "head", "tail", "find"]
    input.command in safe_commands
}

# Prevent AI from modifying critical configuration files
deny_ai_config_modification {
    input.action == "write"
    input.context == "ai"
    
    critical_files := ["/etc", "/root", "/home", ".ssh", ".profile", ".bashrc"]
    contains(input.path, critical_files[_])
}
`,
	}

	for filename, content := range defaultPolicies {
		policyPath := filepath.Join(e.policyDir, filename)
		
		// Only create if file doesn't exist
		if _, err := os.Stat(policyPath); os.IsNotExist(err) {
			if err := os.WriteFile(policyPath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to create policy file %s: %w", filename, err)
			}
		}
	}
	
	return nil
}

// Authorize checks if an action is allowed by policy using simple rules
func (e *Engine) Authorize(action, command string, args []string) (bool, error) {
	// Check each rule
	for _, rule := range e.rules {
		if e.ruleMatches(rule, action, command, args) {
			if rule.Action == "deny" {
				return false, fmt.Errorf("action denied by rule: %s", rule.Description)
			}
		}
	}

	// Default allow if no explicit denial
	return true, nil
}

func (e *Engine) ruleMatches(rule *PolicyRule, action, command string, args []string) bool {
	// Check if command matches
	commandMatches := false
	for _, ruleCmd := range rule.Commands {
		if ruleCmd == command || (strings.Contains(command, ruleCmd) && ruleCmd != command) {
			commandMatches = true
			break
		}
	}

	if !commandMatches {
		return false
	}

	// Check arguments if rule specifies them
	if len(rule.Args) > 0 {
		argMatches := false
		for _, arg := range args {
			for _, pattern := range rule.Args {
				// Use exact matching for root "/" to avoid blocking home directories
				if pattern == "/" {
					if arg == "/" {
						argMatches = true
						break
					}
				} else {
					// Use prefix matching for other path patterns
					if strings.HasPrefix(arg, pattern) {
						argMatches = true
						break
					}
				}
			}
			if argMatches {
				break
			}
		}
		return argMatches
	}

	return true
}

// AuthorizePlugin checks if a plugin action is allowed
func (e *Engine) AuthorizePlugin(pluginName, action, resource string) (bool, error) {
	homeDir := getCurrentDir()
	
	switch action {
	case "execute":
		// Check if command is in allowed list
		if rule, exists := e.rules["plugin_execute"]; exists {
			for _, cmd := range rule.Commands {
				if cmd == resource {
					return true, nil
				}
			}
		}
		return false, fmt.Errorf("command %s not allowed for plugins", resource)
		
	case "read":
		// Allow reading from home directory and plugin-specific directories
		if strings.HasPrefix(resource, homeDir) {
			return true, nil
		}
		// Allow reading from /tmp
		if strings.HasPrefix(resource, "/tmp") {
			return true, nil
		}
		// Deny reading from sensitive system directories
		sensitiveDir := []string{"/etc", "/proc", "/sys", "/root"}
		for _, dir := range sensitiveDir {
			if strings.HasPrefix(resource, dir) {
				return false, fmt.Errorf("read access denied to %s", resource)
			}
		}
		return true, nil
		
	case "write":
		// Allow writing to plugin-specific directories and temp directories
		pluginDir := homeDir + "/.lucien/plugins/" + pluginName
		if strings.HasPrefix(resource, pluginDir) {
			return true, nil
		}
		if strings.HasPrefix(resource, "/tmp") || strings.HasPrefix(resource, "/var/tmp") {
			return true, nil
		}
		// Deny writing to system directories
		return false, fmt.Errorf("write access denied to %s", resource)
		
	default:
		return false, fmt.Errorf("unknown action: %s", action)
	}
}

func getCurrentDir() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return homeDir
	}
	return ""
}

// ReloadPolicies reloads all rules and scans for new .rego files
func (e *Engine) ReloadPolicies() error {
	// Clear existing rules
	e.rules = make(map[string]*PolicyRule)
	e.policyFiles = make(map[string]string)
	
	// Initialize default in-memory rules
	e.initializeDefaultRules()
	
	// Scan directory for .rego files
	entries, err := os.ReadDir(e.policyDir)
	if err != nil {
		return fmt.Errorf("failed to read policy directory: %w", err)
	}
	
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".rego") {
			policyPath := filepath.Join(e.policyDir, entry.Name())
			content, err := os.ReadFile(policyPath)
			if err == nil {
				// Store the file content
				baseName := strings.TrimSuffix(entry.Name(), ".rego")
				e.policyFiles[baseName] = string(content)
				
				// Parse simple rules from the policy content
				e.parseCustomPolicyRules(baseName, string(content))
			}
		}
	}
	
	return nil
}

// parseCustomPolicyRules parses simple deny rules from Rego content  
func (e *Engine) parseCustomPolicyRules(policyName, content string) {
	lines := strings.Split(content, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Look for simple deny rules like: input.command == "custom_blocked"
		if strings.Contains(line, "input.command ==") && strings.Contains(line, "\"") {
			// Extract the command name from the line
			start := strings.Index(line, "\"")
			end := strings.LastIndex(line, "\"")
			if start != -1 && end != -1 && end > start {
				command := line[start+1 : end]
				
				// Create a custom rule
				ruleName := policyName + "_" + command
				e.rules[ruleName] = &PolicyRule{
					Name:        ruleName,
					Description: "Custom policy rule from " + policyName,
					Action:      "deny",
					Commands:    []string{command},
					Conditions:  []string{"custom_policy"},
				}
			}
		}
	}
}

// ListPolicies returns a list of loaded policy names (both rules and files)
func (e *Engine) ListPolicies() []string {
	var names []string
	
	// Add in-memory rules
	for name := range e.rules {
		names = append(names, name)
	}
	
	// Add policy files
	for name := range e.policyFiles {
		names = append(names, name)
	}
	
	return names
}

// GetPolicyContent returns the description or content of a specific policy
func (e *Engine) GetPolicyContent(name string) (string, error) {
	// Check in-memory rules first
	if rule, exists := e.rules[name]; exists {
		return rule.Description, nil
	}
	
	// Check policy files
	if content, exists := e.policyFiles[name]; exists {
		return content, nil
	}
	
	return "", fmt.Errorf("policy %s not found", name)
}