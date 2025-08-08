package completion

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CompletionType represents different types of completions
type CompletionType int

const (
	CommandCompletion CompletionType = iota
	FileCompletion
	DirectoryCompletion
	VariableCompletion
	AliasCompletion
	HistoryCompletion
)

// Suggestion represents a completion suggestion
type Suggestion struct {
	Text        string
	Type        CompletionType
	Description string
	Priority    int // Higher priority shows first
}

// Engine handles intelligent tab completion
type Engine struct {
	commands    map[string]bool
	builtins    map[string]bool
	aliases     map[string]string
	variables   map[string]string
	historyMgr  HistoryProvider // Interface to access history
}

// HistoryProvider interface for accessing command history
type HistoryProvider interface {
	GetRecent(n int) []HistoryEntry
	Search(query string, limit int) []HistoryEntry
}

// HistoryEntry represents a history entry for completion
type HistoryEntry struct {
	Command string
}

// New creates a new completion engine
func New() *Engine {
	engine := &Engine{
		commands:  make(map[string]bool),
		builtins:  make(map[string]bool),
		aliases:   make(map[string]string),
		variables: make(map[string]string),
	}
	
	// Initialize common commands
	engine.initializeCommands()
	
	return engine
}

// initializeCommands populates common system commands
func (e *Engine) initializeCommands() {
	// Common Unix/Linux commands
	commonCommands := []string{
		"ls", "cd", "pwd", "mkdir", "rmdir", "rm", "cp", "mv", "chmod", "chown",
		"grep", "find", "locate", "which", "whereis", "file", "stat", "du", "df",
		"ps", "top", "htop", "kill", "killall", "jobs", "bg", "fg", "nohup",
		"cat", "less", "more", "head", "tail", "tee", "sort", "uniq", "cut",
		"awk", "sed", "tr", "wc", "diff", "patch", "tar", "gzip", "gunzip",
		"curl", "wget", "ssh", "scp", "rsync", "ping", "netstat", "ss",
		"git", "svn", "hg", "make", "cmake", "gcc", "clang", "python", "node",
		"npm", "pip", "go", "rust", "cargo", "docker", "kubectl", "helm",
		"vim", "nano", "emacs", "code", "history", "alias", "export", "env",
		"echo", "printf", "date", "cal", "uptime", "uname", "whoami", "id",
	}
	
	for _, cmd := range commonCommands {
		e.commands[cmd] = true
	}
	
	// Built-in commands
	builtinCommands := []string{
		"cd", "set", "export", "alias", "history", "exit", "pwd", "echo", "help",
	}
	
	for _, builtin := range builtinCommands {
		e.builtins[builtin] = true
	}
}

// SetHistoryProvider sets the history provider for completion
func (e *Engine) SetHistoryProvider(provider HistoryProvider) {
	e.historyMgr = provider
}

// UpdateAliases updates the aliases used for completion
func (e *Engine) UpdateAliases(aliases map[string]string) {
	e.aliases = make(map[string]string)
	for k, v := range aliases {
		e.aliases[k] = v
	}
}

// UpdateVariables updates environment variables for completion
func (e *Engine) UpdateVariables(variables map[string]string) {
	e.variables = make(map[string]string)
	for k, v := range variables {
		e.variables[k] = v
	}
}

// Complete performs intelligent tab completion
func (e *Engine) Complete(input string, cursorPos int) []Suggestion {
	if cursorPos < 0 || cursorPos > len(input) {
		cursorPos = len(input)
	}
	
	// Extract the part of the command we're completing
	beforeCursor := input[:cursorPos]
	tokens := e.tokenize(beforeCursor)
	
	if len(tokens) == 0 {
		return e.getCommandSuggestions("")
	}
	
	lastToken := tokens[len(tokens)-1]
	
	// If we're completing the first token, it's a command
	if len(tokens) == 1 && !strings.HasSuffix(beforeCursor, " ") {
		return e.getCommandSuggestions(lastToken)
	}
	
	// If we have spaces after the command, we're completing arguments
	if len(tokens) >= 1 {
		command := tokens[0]
		
		// Handle special cases for specific commands
		switch command {
		case "cd":
			return e.getDirectorySuggestions(lastToken)
		case "ls", "cat", "less", "more", "head", "tail", "rm", "cp", "mv":
			return e.getFileSuggestions(lastToken)
		case "set", "export", "echo":
			if strings.HasPrefix(lastToken, "$") {
				return e.getVariableSuggestions(lastToken[1:])
			}
			return e.getFileSuggestions(lastToken)
		case "history":
			return e.getHistoryCommandSuggestions(lastToken)
		default:
			// Default to file completion for unknown commands
			return e.getFileSuggestions(lastToken)
		}
	}
	
	return []Suggestion{}
}

// tokenize splits the input into tokens, respecting quotes
func (e *Engine) tokenize(input string) []string {
	var tokens []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)
	
	for i := 0; i < len(input); i++ {
		char := input[i]
		
		switch {
		case char == '"' || char == '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
				quoteChar = 0
			} else {
				current.WriteByte(char)
			}
			
		case char == ' ' || char == '\t':
			if inQuotes {
				current.WriteByte(char)
			} else if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			
		default:
			current.WriteByte(char)
		}
	}
	
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	
	return tokens
}

// getCommandSuggestions returns command completion suggestions
func (e *Engine) getCommandSuggestions(prefix string) []Suggestion {
	var suggestions []Suggestion
	
	// Built-in commands (highest priority)
	for builtin := range e.builtins {
		if strings.HasPrefix(builtin, prefix) {
			suggestions = append(suggestions, Suggestion{
				Text:        builtin,
				Type:        CommandCompletion,
				Description: "built-in command",
				Priority:    100,
			})
		}
	}
	
	// Aliases (high priority)
	for alias := range e.aliases {
		if strings.HasPrefix(alias, prefix) {
			suggestions = append(suggestions, Suggestion{
				Text:        alias,
				Type:        AliasCompletion,
				Description: "alias: " + e.aliases[alias],
				Priority:    90,
			})
		}
	}
	
	// Common commands (medium priority)
	for command := range e.commands {
		if strings.HasPrefix(command, prefix) {
			suggestions = append(suggestions, Suggestion{
				Text:        command,
				Type:        CommandCompletion,
				Description: "command",
				Priority:    80,
			})
		}
	}
	
	// History-based commands (lower priority)
	if e.historyMgr != nil {
		historyEntries := e.historyMgr.GetRecent(100)
		seenCommands := make(map[string]bool)
		
		for _, entry := range historyEntries {
			if tokens := e.tokenize(entry.Command); len(tokens) > 0 {
				cmd := tokens[0]
				if strings.HasPrefix(cmd, prefix) && !seenCommands[cmd] {
					seenCommands[cmd] = true
					suggestions = append(suggestions, Suggestion{
						Text:        cmd,
						Type:        HistoryCompletion,
						Description: "from history",
						Priority:    70,
					})
				}
			}
		}
	}
	
	// Executable files in PATH (lowest priority)
	if pathSuggestions := e.getPathExecutables(prefix); len(pathSuggestions) > 0 {
		suggestions = append(suggestions, pathSuggestions...)
	}
	
	// Sort by priority (descending) then alphabetically
	sort.Slice(suggestions, func(i, j int) bool {
		if suggestions[i].Priority != suggestions[j].Priority {
			return suggestions[i].Priority > suggestions[j].Priority
		}
		return suggestions[i].Text < suggestions[j].Text
	})
	
	return suggestions
}

// getFileSuggestions returns file completion suggestions
func (e *Engine) getFileSuggestions(prefix string) []Suggestion {
	var suggestions []Suggestion
	
	// Determine the directory to search
	dir := "."
	pattern := prefix
	
	if strings.Contains(prefix, "/") || strings.Contains(prefix, "\\") {
		dir = filepath.Dir(prefix)
		pattern = filepath.Base(prefix)
		
		// Handle absolute paths and home directory expansion
		if strings.HasPrefix(prefix, "~/") {
			if homeDir, err := os.UserHomeDir(); err == nil {
				dir = filepath.Join(homeDir, strings.TrimPrefix(filepath.Dir(prefix), "~"))
			}
		}
	}
	
	// Read directory contents
	entries, err := os.ReadDir(dir)
	if err != nil {
		return suggestions
	}
	
	for _, entry := range entries {
		name := entry.Name()
		
		// Skip hidden files unless pattern starts with dot
		if strings.HasPrefix(name, ".") && !strings.HasPrefix(pattern, ".") {
			continue
		}
		
		// Check if name matches pattern
		if strings.HasPrefix(name, pattern) {
			suggestionType := FileCompletion
			description := "file"
			priority := 60
			
			if entry.IsDir() {
				suggestionType = DirectoryCompletion
				description = "directory"
				priority = 65
				name += "/"
			}
			
			// Reconstruct full path if needed
			fullPath := name
			if dir != "." {
				fullPath = filepath.Join(dir, name)
				if strings.HasPrefix(prefix, "~/") {
					fullPath = "~/" + strings.TrimPrefix(fullPath, filepath.Join(os.Getenv("HOME"), "/"))
				}
			}
			
			suggestions = append(suggestions, Suggestion{
				Text:        fullPath,
				Type:        suggestionType,
				Description: description,
				Priority:    priority,
			})
		}
	}
	
	// Sort directories first, then files, both alphabetically
	sort.Slice(suggestions, func(i, j int) bool {
		if suggestions[i].Type != suggestions[j].Type {
			return suggestions[i].Type == DirectoryCompletion
		}
		return suggestions[i].Text < suggestions[j].Text
	})
	
	return suggestions
}

// getDirectorySuggestions returns directory-only completion suggestions
func (e *Engine) getDirectorySuggestions(prefix string) []Suggestion {
	suggestions := e.getFileSuggestions(prefix)
	
	// Filter to only include directories
	var dirSuggestions []Suggestion
	for _, suggestion := range suggestions {
		if suggestion.Type == DirectoryCompletion {
			dirSuggestions = append(dirSuggestions, suggestion)
		}
	}
	
	return dirSuggestions
}

// getVariableSuggestions returns environment variable suggestions
func (e *Engine) getVariableSuggestions(prefix string) []Suggestion {
	var suggestions []Suggestion
	
	for varName := range e.variables {
		if strings.HasPrefix(varName, prefix) {
			suggestions = append(suggestions, Suggestion{
				Text:        "$" + varName,
				Type:        VariableCompletion,
				Description: "variable: " + e.variables[varName],
				Priority:    75,
			})
		}
	}
	
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Text < suggestions[j].Text
	})
	
	return suggestions
}

// getHistoryCommandSuggestions returns suggestions for history command arguments
func (e *Engine) getHistoryCommandSuggestions(prefix string) []Suggestion {
	var suggestions []Suggestion
	
	historyCommands := []string{"clear", "stats", "search"}
	for _, cmd := range historyCommands {
		if strings.HasPrefix(cmd, prefix) {
			suggestions = append(suggestions, Suggestion{
				Text:        cmd,
				Type:        CommandCompletion,
				Description: "history subcommand",
				Priority:    85,
			})
		}
	}
	
	return suggestions
}

// getPathExecutables searches PATH for executable files
func (e *Engine) getPathExecutables(prefix string) []Suggestion {
	var suggestions []Suggestion
	
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return suggestions
	}
	
	pathDirs := strings.Split(pathEnv, string(os.PathListSeparator))
	seenExecutables := make(map[string]bool)
	
	for _, pathDir := range pathDirs {
		if pathDir == "" {
			continue
		}
		
		entries, err := os.ReadDir(pathDir)
		if err != nil {
			continue
		}
		
		for _, entry := range entries {
			name := entry.Name()
			
			// Skip if already seen or doesn't match prefix
			if seenExecutables[name] || !strings.HasPrefix(name, prefix) {
				continue
			}
			
			// Check if file is executable (simplified check)
			if !entry.IsDir() {
				if info, err := entry.Info(); err == nil {
					mode := info.Mode()
					if mode&0111 != 0 { // Has execute permission
						seenExecutables[name] = true
						suggestions = append(suggestions, Suggestion{
							Text:        name,
							Type:        CommandCompletion,
							Description: "executable",
							Priority:    50,
						})
					}
				}
			}
		}
	}
	
	return suggestions
}

// GetBestMatch returns the best single completion match if there's only one suggestion
// or a common prefix if multiple suggestions share one
func (e *Engine) GetBestMatch(suggestions []Suggestion) string {
	if len(suggestions) == 0 {
		return ""
	}
	
	if len(suggestions) == 1 {
		return suggestions[0].Text
	}
	
	// Find common prefix among all suggestions
	commonPrefix := suggestions[0].Text
	for _, suggestion := range suggestions[1:] {
		commonPrefix = longestCommonPrefix(commonPrefix, suggestion.Text)
		if commonPrefix == "" {
			break
		}
	}
	
	return commonPrefix
}

// longestCommonPrefix finds the longest common prefix between two strings
func longestCommonPrefix(a, b string) string {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	
	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return a[:i]
		}
	}
	
	return a[:minLen]
}