package history

import (
	"strings"
	"time"
)

var commands []string

// Add stores a command in the history.
func Add(cmd string) {
	commands = append(commands, cmd)
}

// All returns all commands in the history.
func All() []string {
	return commands
}

// Config represents configuration for the history manager
type Config struct {
	HistoryFile string
	MaxEntries  int
	AutoSave    bool
}

// Entry represents a history entry
type Entry struct {
	Command   string
	Timestamp time.Time
}

// Manager manages command history
type Manager struct {
	config  *Config
	entries []Entry
}

// New creates a new history manager
func New(config *Config) (*Manager, error) {
	return &Manager{
		config:  config,
		entries: make([]Entry, 0),
	}, nil
}

// GetRecent returns the n most recent history entries
func (m *Manager) GetRecent(n int) []Entry {
	if n <= 0 || len(m.entries) == 0 {
		return []Entry{}
	}
	
	start := len(m.entries) - n
	if start < 0 {
		start = 0
	}
	
	return m.entries[start:]
}

// Search searches for commands containing the query string
func (m *Manager) Search(query string, limit int) []Entry {
	var results []Entry
	
	for _, entry := range m.entries {
		if strings.Contains(entry.Command, query) {
			results = append(results, entry)
			if len(results) >= limit {
				break
			}
		}
	}
	
	return results
}