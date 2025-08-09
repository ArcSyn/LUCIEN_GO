package shell

import (
	"fmt"
	"regexp"
	"strings"
)

// SecurityMode represents the security enforcement level
type SecurityMode int

const (
	SecurityModePermissive SecurityMode = iota
	SecurityModeStrict
)

// SecurityGuard handles command injection prevention and validation
type SecurityGuard struct {
	mode               SecurityMode
	whitelistedBuiltins []string
	dangerousPatterns  []*regexp.Regexp
}

// NewSecurityGuard creates a new security guard
func NewSecurityGuard() *SecurityGuard {
	return &SecurityGuard{
		mode: SecurityModePermissive,
		whitelistedBuiltins: []string{
			"echo", "pwd", "cd", "ls", "dir", "cat", "type", "help",
			"history", "alias", "unalias", "set", "unset", "export",
			"clear", "exit", "jobs", "fg", "bg", "kill",
		},
		dangerousPatterns: []*regexp.Regexp{
			regexp.MustCompile(`\$\([^)]*\)`),              // Command substitution $(...)
			regexp.MustCompile("`[^`]*`"),                  // Backtick command substitution
			regexp.MustCompile(`\|\s*\w+\s*>`),             // Pipe to file redirection
			regexp.MustCompile(`>\s*/dev/`),                // Writing to device files
			regexp.MustCompile(`>\s*&\d+`),                 // File descriptor manipulation
			regexp.MustCompile(`\|\s*sh\s*$`),              // Pipe to shell
			regexp.MustCompile(`\|\s*bash\s*$`),            // Pipe to bash
			regexp.MustCompile(`\|\s*powershell\s*$`),      // Pipe to PowerShell
			regexp.MustCompile(`;\s*rm\s+-rf`),             // Dangerous rm commands
			regexp.MustCompile(`;\s*del\s+/[sqf]`),         // Dangerous del commands
		},
	}
}

// SetMode sets the security mode
func (sg *SecurityGuard) SetMode(mode SecurityMode) {
	sg.mode = mode
}

// GetMode returns the current security mode
func (sg *SecurityGuard) GetMode() SecurityMode {
	return sg.mode
}

// IsWhitelistedBuiltin checks if a command is a whitelisted builtin
func (sg *SecurityGuard) IsWhitelistedBuiltin(command string) bool {
	for _, builtin := range sg.whitelistedBuiltins {
		if strings.EqualFold(command, builtin) {
			return true
		}
	}
	return false
}

// ValidateCommandChain validates a command chain after parsing
func (sg *SecurityGuard) ValidateCommandChain(chain *CommandChain) error {
	if sg.mode == SecurityModePermissive {
		return nil // Allow everything in permissive mode
	}

	// In strict mode, check each command and operator
	for i, cmd := range chain.Commands {
		// Allow whitelisted builtins to use any arguments
		if sg.IsWhitelistedBuiltin(cmd.Name) {
			continue
		}

		// Check for dangerous patterns in arguments
		fullCommand := cmd.Name + " " + strings.Join(cmd.Args, " ")
		for _, pattern := range sg.dangerousPatterns {
			if pattern.MatchString(fullCommand) {
				return fmt.Errorf("security guard: potentially dangerous pattern detected in strict mode: %s", pattern.String())
			}
		}

		// In strict mode, non-whitelisted commands with operators are blocked
		if i < len(chain.Operators) {
			op := chain.Operators[i]
			if op == "&&" || op == "||" || op == "|" {
				return fmt.Errorf("security guard: operator '%s' with non-whitelisted command '%s' blocked in strict mode", op, cmd.Name)
			}
		}
	}

	return nil
}

// ValidateCommand validates a single command for security issues
func (sg *SecurityGuard) ValidateCommand(cmd *Command) error {
	if sg.mode == SecurityModePermissive {
		return nil
	}

	// Allow whitelisted builtins
	if sg.IsWhitelistedBuiltin(cmd.Name) {
		return nil
	}

	// Check arguments for injection patterns
	for _, arg := range cmd.Args {
		for _, pattern := range sg.dangerousPatterns {
			if pattern.MatchString(arg) {
				return fmt.Errorf("security guard: dangerous pattern in argument blocked in strict mode: %s", arg)
			}
		}
	}

	return nil
}