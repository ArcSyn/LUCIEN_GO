package policy

import (
	"fmt"
	"os"
	"strings"
	"regexp"
)

// Engine represents a simplified policy engine for Lucien
type Engine struct {
	policyDir string
	rules     map[string]*PolicyRule
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
		policyDir: policyDir,
		rules:     make(map[string]*PolicyRule),
	}

	// Initialize with default security rules
	engine.initializeDefaultRules()

	return engine, nil
}

func (e *Engine) initializeDefaultRules() {
	// Security rules for dangerous commands
	e.rules["deny_dangerous"] = &PolicyRule{
		Name:        "deny_dangerous",
		Description: "Block dangerous system commands",
		Action:      "deny",
		Commands:    []string{"dd", "mkfs", "fdisk", "format", "sudo", "su", "rm -rf", "rmdir"},
		Conditions:  []string{"system_critical"},
	}

	// File system protection
	e.rules["protect_system_dirs"] = &PolicyRule{
		Name:        "protect_system_dirs",
		Description: "Protect system directories from modification",
		Action:      "deny",
		Commands:    []string{"rm", "rmdir", "chmod", "chown"},
		Args:        []string{"^/", "/etc", "/usr", "/var", "/sys", "/proc"},
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

	// AI safety
	e.rules["ai_safety"] = &PolicyRule{
		Name:        "ai_safety",
		Description: "AI cannot execute dangerous commands",
		Action:      "deny",
		Commands:    []string{"rm", "rmdir", "dd", "mkfs", "sudo", "su", "chmod", "chown"},
		Conditions:  []string{"ai_context"},
	}
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
		for _, arg := range args {
			for _, pattern := range rule.Args {
				if matched, _ := regexp.MatchString(pattern, arg); matched {
					return true
				}
			}
		}
		return false
	}

	return true
}

// AuthorizePlugin checks if a plugin action is allowed
func (e *Engine) AuthorizePlugin(pluginName, action, resource string) (bool, error) {
	// Check plugin-specific rules
	if rule, exists := e.rules["plugin_execute"]; exists {
		if rule.Action == "allow" {
			// Check if action is in allowed list
			for _, cmd := range rule.Commands {
				if cmd == action {
					return true, nil
				}
			}
		}
	}

	// Default deny for plugins
	return false, nil
}

func getCurrentDir() string {
	if cwd, err := os.Getwd(); err == nil {
		return cwd
	}
	return ""
}

// ReloadPolicies reloads all rules
func (e *Engine) ReloadPolicies() error {
	e.rules = make(map[string]*PolicyRule)
	e.initializeDefaultRules()
	return nil
}

// ListPolicies returns a list of loaded policy names
func (e *Engine) ListPolicies() []string {
	var names []string
	for name := range e.rules {
		names = append(names, name)
	}
	return names
}

// GetPolicyContent returns the description of a specific policy
func (e *Engine) GetPolicyContent(name string) (string, error) {
	if rule, exists := e.rules[name]; exists {
		return rule.Description, nil
	}
	return "", fmt.Errorf("policy %s not found", name)
}