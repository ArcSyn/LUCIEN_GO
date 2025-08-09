package shell

import (
	"fmt"
	"strings"
)

// CommandType represents the type of command operation
type CommandType int

const (
	CommandSimple CommandType = iota
	CommandPipe                // |
	CommandAnd                 // &&
	CommandOr                  // ||
	CommandSequence            // ;
	CommandBackground          // &
)

// CommandChain represents a chain of commands with operators
type CommandChain struct {
	Commands  []Command
	Operators []string
	Types     []CommandType
}

// Len returns the number of commands in the chain
func (cc *CommandChain) Len() int {
	if cc == nil {
		return 0
	}
	return len(cc.Commands)
}

// At returns the command at the specified index
func (cc *CommandChain) At(i int) Command {
	return cc.Commands[i]
}

// Slice returns a copy of the commands slice
func (cc *CommandChain) Slice() []Command {
	if cc == nil {
		return nil
	}
	out := make([]Command, len(cc.Commands))
	copy(out, cc.Commands)
	return out
}

// parseCommandLineAdvanced parses a complete command line with all operators
func (s *Shell) parseCommandLineAdvanced(cmdLine string) (*CommandChain, error) {
	if strings.TrimSpace(cmdLine) == "" {
		return &CommandChain{}, nil
	}

	// Parse command line respecting quotes for operators
	commands, operators, types, err := s.splitCommandLineWithTypes(cmdLine)
	if err != nil {
		return nil, err
	}

	chain := &CommandChain{
		Commands:  commands,
		Operators: operators,
		Types:     types,
	}

	// Validate with security guard after parsing
	if err := s.securityGuard.ValidateCommandChain(chain); err != nil {
		return nil, err
	}

	return chain, nil
}

// splitCommandLineWithTypes splits command line by all operators while respecting quotes
func (s *Shell) splitCommandLineWithTypes(cmdLine string) ([]Command, []string, []CommandType, error) {
	var commands []Command
	var operators []string
	var types []CommandType
	var currentCommand strings.Builder
	
	inQuotes := false
	quoteChar := byte(0)
	escaped := false
	i := 0
	
	for i < len(cmdLine) {
		char := cmdLine[i]
		
		// Handle escape sequences
		if escaped {
			// For Windows paths, preserve backslashes when they're not escaping quotes
			if char == '"' || char == '\'' {
				// This is an escaped quote, add it literally
				currentCommand.WriteByte(char)
			} else {
				// This is a backslash followed by another character - preserve both
				currentCommand.WriteByte('\\')
				currentCommand.WriteByte(char)
			}
			escaped = false
			i++
			continue
		}
		
		if char == '\\' {
			// Check if next character is a quote - if so, this is an escape
			if i+1 < len(cmdLine) && (cmdLine[i+1] == '"' || cmdLine[i+1] == '\'') {
				escaped = true
				i++
				continue
			} else {
				// This is a regular backslash (like in Windows paths), keep it
				currentCommand.WriteByte(char)
				i++
				continue
			}
		}
		
		// Handle quotes - operators inside quotes are literal
		if char == '"' || char == '\'' {
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
				quoteChar = 0
			}
			currentCommand.WriteByte(char)
			i++
			continue
		}
		
		// If we're in quotes, treat everything as literal
		if inQuotes {
			currentCommand.WriteByte(char)
			i++
			continue
		}
		
		// Skip whitespace between commands/operators
		if char == ' ' || char == '\t' {
			if currentCommand.Len() == 0 {
				i++
				continue
			}
			// Add space to current command if it's not empty
			currentCommand.WriteByte(char)
			i++
			continue
		}
		
		// Check for two-character operators first
		if i < len(cmdLine)-1 {
			twoChar := cmdLine[i:i+2]
			var opType CommandType
			var isOperator bool
			
			switch twoChar {
			case "&&":
				opType = CommandAnd
				isOperator = true
			case "||":
				opType = CommandOr  
				isOperator = true
			}
			
			if isOperator {
				// Save current command
				if err := s.saveCurrentCommand(&currentCommand, &commands); err != nil {
					return nil, nil, nil, err
				}
				
				operators = append(operators, twoChar)
				types = append(types, opType)
				i += 2
				s.skipWhitespace(cmdLine, &i)
				continue
			}
		}
		
		// Check for single-character operators
		var opType CommandType
		var isOperator bool
		var opStr string
		
		switch char {
		case ';':
			opType = CommandSequence
			isOperator = true
			opStr = ";"
		case '|':
			opType = CommandPipe
			isOperator = true
			opStr = "|"
		case '&':
			opType = CommandBackground
			isOperator = true
			opStr = "&"
		}
		
		if isOperator {
			// Save current command
			if err := s.saveCurrentCommand(&currentCommand, &commands); err != nil {
				return nil, nil, nil, err
			}
			
			operators = append(operators, opStr)
			types = append(types, opType)
			i++
			s.skipWhitespace(cmdLine, &i)
			continue
		}
		
		// Regular character
		currentCommand.WriteByte(char)
		i++
	}
	
	// Add final command
	if currentCommand.Len() > 0 {
		if err := s.saveCurrentCommand(&currentCommand, &commands); err != nil {
			return nil, nil, nil, err
		}
	}
	
	return commands, operators, types, nil
}

// saveCurrentCommand saves the current command buffer to the commands slice
func (s *Shell) saveCurrentCommand(buffer *strings.Builder, commands *[]Command) error {
	cmdStr := strings.TrimSpace(buffer.String())
	if cmdStr == "" {
		buffer.Reset()
		return nil
	}
	
	cmd, err := s.parseCommand(cmdStr)
	if err != nil {
		return err
	}
	
	if cmd.Name != "" {
		*commands = append(*commands, cmd)
	}
	
	buffer.Reset()
	return nil
}

// skipWhitespace skips whitespace characters and advances the index
func (s *Shell) skipWhitespace(cmdLine string, i *int) {
	for *i < len(cmdLine) && (cmdLine[*i] == ' ' || cmdLine[*i] == '\t') {
		(*i)++
	}
}

// parseCommand parses a single command into name and args
func (s *Shell) parseCommand(cmdStr string) (Command, error) {
	cmdStr = strings.TrimSpace(cmdStr)
	if cmdStr == "" {
		return Command{}, nil
	}
	
	// Tokenize respecting quotes
	tokens, err := s.tokenizeAdvanced(cmdStr)
	if err != nil {
		return Command{}, err
	}
	
	if len(tokens) == 0 {
		return Command{}, nil
	}
	
	// Handle alias expansion
	commandName := tokens[0]
	args := tokens[1:]
	
	if aliasValue, exists := s.aliases[commandName]; exists {
		// Parse the alias value and prepend to arguments
		aliasTokens, err := s.tokenizeAdvanced(aliasValue)
		if err == nil && len(aliasTokens) > 0 {
			commandName = aliasTokens[0]
			// Merge alias args with original args
			newArgs := make([]string, 0, len(aliasTokens)-1+len(args))
			newArgs = append(newArgs, aliasTokens[1:]...)
			newArgs = append(newArgs, args...)
			args = newArgs
		}
	}
	
	cmd := Command{
		Name:      commandName,
		Args:      args,
		Redirects: make(map[string]string),
	}
	
	// Validate individual command
	if err := s.securityGuard.ValidateCommand(&cmd); err != nil {
		return Command{}, err
	}
	
	return cmd, nil
}

// tokenizeAdvanced tokenizes a command string respecting quotes and escapes
func (s *Shell) tokenizeAdvanced(input string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)
	escaped := false
	
	for i := 0; i < len(input); i++ {
		char := input[i]
		
		if escaped {
			// For Windows paths, preserve backslashes when they're not escaping quotes
			if char == '"' || char == '\'' {
				// This is an escaped quote, add it literally
				current.WriteByte(char)
			} else {
				// This is a backslash followed by another character - preserve both
				current.WriteByte('\\')
				current.WriteByte(char)
			}
			escaped = false
			continue
		}
		
		if char == '\\' {
			// Check if next character is a quote - if so, this is an escape
			if i+1 < len(input) && (input[i+1] == '"' || input[i+1] == '\'') {
				escaped = true
				continue
			} else {
				// This is a regular backslash (like in Windows paths), keep it
				current.WriteByte(char)
			}
		} else if char == '"' || char == '\'' {
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
				quoteChar = 0
			} else {
				current.WriteByte(char)
			}
		} else if char == ' ' || char == '\t' {
			if inQuotes {
				current.WriteByte(char)
			} else {
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}
			}
		} else {
			current.WriteByte(char)
		}
	}
	
	if inQuotes {
		return nil, fmt.Errorf("unclosed quote in command")
	}
	
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	
	return tokens, nil
}

// FirstArgAsPath extracts the first argument as a path, handling quotes correctly
func FirstArgAsPath(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("no path argument provided")
	}
	
	path := args[0]
	
	// Remove surrounding quotes - only outer quotes, preserving inner ones
	if len(path) >= 2 {
		if (path[0] == '"' && path[len(path)-1] == '"') ||
		   (path[0] == '\'' && path[len(path)-1] == '\'') {
			path = path[1 : len(path)-1]
		}
	}
	
	if path == "" {
		return "", fmt.Errorf("empty path argument")
	}
	
	return path, nil
}