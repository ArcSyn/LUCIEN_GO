package env

import "os"

// Get returns the value of an environment variable, or a default if unset.
func Get(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

// Config represents configuration for the environment manager
type Config struct {
	PersistFile string
	AutoSave    bool
}

// Manager manages environment variables
type Manager struct {
	config *Config
}

// New creates a new environment manager
func New(config *Config) (*Manager, error) {
	return &Manager{
		config: config,
	}, nil
}