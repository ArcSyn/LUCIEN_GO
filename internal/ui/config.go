package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config represents the Lucien CLI configuration
type Config struct {
	Shell     ShellConfig     `toml:"shell"`
	UI        UIConfig        `toml:"ui"`
	History   HistoryConfig   `toml:"history"`
	Completion CompletionConfig `toml:"completion"`
	AI        AIConfig        `toml:"ai"`
}

// ShellConfig holds shell-specific configuration
type ShellConfig struct {
	Prompt           string            `toml:"prompt"`
	SafeMode         bool              `toml:"safe_mode"`
	DefaultTheme     string            `toml:"default_theme"`
	Aliases          map[string]string `toml:"aliases"`
	Environment      map[string]string `toml:"environment"`
	ExecutionTimeout int               `toml:"execution_timeout"` // in seconds
	ExecutorMode     string            `toml:"executor_mode"`     // "shell" or "internal"
}

// UIConfig holds UI-specific configuration
type UIConfig struct {
	AnimatedStartup  bool   `toml:"animated_startup"`
	GlitchEffects    bool   `toml:"glitch_effects"`
	ColorSupport     string `toml:"color_support"` // "full", "256", "16", "none"
	FontSize         int    `toml:"font_size"`
	WindowWidth      int    `toml:"window_width"`
	WindowHeight     int    `toml:"window_height"`
}

// HistoryConfig holds history-specific configuration  
type HistoryConfig struct {
	Enabled     bool   `toml:"enabled"`
	MaxEntries  int    `toml:"max_entries"`
	SaveOnExit  bool   `toml:"save_on_exit"`
	FilePath    string `toml:"file_path"`
	SearchLimit int    `toml:"search_limit"`
}

// CompletionConfig holds tab completion configuration
type CompletionConfig struct {
	Enabled           bool `toml:"enabled"`
	SuggestionsPerPage int  `toml:"suggestions_per_page"`
	AutoCD            bool `toml:"auto_cd"`
	FuzzyMatching     bool `toml:"fuzzy_matching"`
	ShowDescriptions  bool `toml:"show_descriptions"`
}

// AIConfig holds AI-specific configuration
type AIConfig struct {
	Enabled          bool   `toml:"enabled"`
	SuggestCommands  bool   `toml:"suggest_commands"`
	ConfidenceThreshold float64 `toml:"confidence_threshold"`
	ModelProvider    string `toml:"model_provider"`
	ApiKey          string `toml:"api_key"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	historyPath := filepath.Join(homeDir, ".lucien", "history")
	
	return &Config{
		Shell: ShellConfig{
			Prompt:           "lucien@nexus:~$ ",
			SafeMode:         true,
			DefaultTheme:     "nexus",
			Aliases:          make(map[string]string),
			Environment:      make(map[string]string),
			ExecutionTimeout: 30,
			ExecutorMode:     "shell",
		},
		UI: UIConfig{
			AnimatedStartup: true,
			GlitchEffects:   true,
			ColorSupport:    "full",
			FontSize:        12,
			WindowWidth:     120,
			WindowHeight:    40,
		},
		History: HistoryConfig{
			Enabled:     true,
			MaxEntries:  5000,
			SaveOnExit:  true,
			FilePath:    historyPath,
			SearchLimit: 100,
		},
		Completion: CompletionConfig{
			Enabled:           true,
			SuggestionsPerPage: 8,
			AutoCD:            true,
			FuzzyMatching:     true,
			ShowDescriptions:  true,
		},
		AI: AIConfig{
			Enabled:             true,
			SuggestCommands:     true,
			ConfidenceThreshold: 0.75,
			ModelProvider:       "local",
			ApiKey:              "",
		},
	}
}

// GetConfigPath returns the path to the configuration file with fallback directories
func GetConfigPath() (string, error) {
	// Primary config directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	
	primaryPath := filepath.Join(homeDir, ".lucien", "config.toml")
	
	// Check if primary directory is writable
	lucienDir := filepath.Join(homeDir, ".lucien")
	if err := os.MkdirAll(lucienDir, 0755); err == nil {
		// Test write access
		testFile := filepath.Join(lucienDir, ".write_test")
		if file, err := os.Create(testFile); err == nil {
			file.Close()
			os.Remove(testFile)
			return primaryPath, nil
		}
	}
	
	// Fallback directories
	var fallbackPaths []string
	
	switch runtime.GOOS {
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			fallbackPaths = append(fallbackPaths, filepath.Join(appData, "Lucien", "config.toml"))
		}
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			fallbackPaths = append(fallbackPaths, filepath.Join(localAppData, "Lucien", "config.toml"))
		}
	default:
		if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
			fallbackPaths = append(fallbackPaths, filepath.Join(xdgConfig, "lucien", "config.toml"))
		} else {
			fallbackPaths = append(fallbackPaths, filepath.Join(homeDir, ".config", "lucien", "config.toml"))
		}
		
		// Additional Unix fallbacks
		if tmpDir := os.Getenv("TMPDIR"); tmpDir != "" {
			fallbackPaths = append(fallbackPaths, filepath.Join(tmpDir, ".lucien", "config.toml"))
		}
		fallbackPaths = append(fallbackPaths, "/tmp/.lucien/config.toml")
	}
	
	// Try fallback paths
	for _, path := range fallbackPaths {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err == nil {
			// Test write access
			testFile := filepath.Join(dir, ".write_test")
			if file, err := os.Create(testFile); err == nil {
				file.Close()
				os.Remove(testFile)
				return path, nil
			}
		}
	}
	
	return "", fmt.Errorf("no writable config directory found")
}

// LoadConfig loads configuration from file or returns default config
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return DefaultConfig(), err
	}
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		config := DefaultConfig()
		if saveErr := SaveConfig(config); saveErr != nil {
			return config, fmt.Errorf("could not save default config: %v", saveErr)
		}
		return config, nil
	}
	
	// Load existing config
	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return DefaultConfig(), fmt.Errorf("could not parse config file: %v", err)
	}
	
	// Merge with defaults for missing fields
	defaultConfig := DefaultConfig()
	mergeConfigs(&config, defaultConfig)
	
	return &config, nil
}

// SaveConfig saves the configuration to file with debouncing
func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("could not create config directory: %v", err)
	}
	
	// Create temporary file for atomic write
	tempFile := configPath + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("could not create temp config file: %v", err)
	}
	defer file.Close()
	
	// Write TOML with comments
	if err := writeConfigWithComments(file, config); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("could not write config: %v", err)
	}
	
	// Atomic move
	if err := os.Rename(tempFile, configPath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("could not save config: %v", err)
	}
	
	return nil
}

// writeConfigWithComments writes the config with helpful comments
func writeConfigWithComments(file *os.File, config *Config) error {
	encoder := toml.NewEncoder(file)
	
	// Write header comment
	file.WriteString("# Lucien CLI Configuration File\n")
	file.WriteString("# This file is automatically generated but can be edited\n")
	file.WriteString("# Changes take effect on next shell restart\n\n")
	
	return encoder.Encode(config)
}

// mergeConfigs merges missing fields from default into target
func mergeConfigs(target, defaultConfig *Config) {
	// Simple field merging - in a real implementation you'd want more sophisticated merging
	if target.Shell.Prompt == "" {
		target.Shell.Prompt = defaultConfig.Shell.Prompt
	}
	if target.Shell.DefaultTheme == "" {
		target.Shell.DefaultTheme = defaultConfig.Shell.DefaultTheme
	}
	if target.Shell.ExecutionTimeout == 0 {
		target.Shell.ExecutionTimeout = defaultConfig.Shell.ExecutionTimeout
	}
	if target.UI.ColorSupport == "" {
		target.UI.ColorSupport = defaultConfig.UI.ColorSupport
	}
	if target.History.MaxEntries == 0 {
		target.History.MaxEntries = defaultConfig.History.MaxEntries
	}
	if target.Completion.SuggestionsPerPage == 0 {
		target.Completion.SuggestionsPerPage = defaultConfig.Completion.SuggestionsPerPage
	}
	if target.AI.ConfidenceThreshold == 0 {
		target.AI.ConfidenceThreshold = defaultConfig.AI.ConfidenceThreshold
	}
}

// SetConfigValue sets a configuration value by dot-notation path
func SetConfigValue(config *Config, key string, value string) error {
	parts := strings.Split(key, ".")
	if len(parts) < 2 {
		return fmt.Errorf("invalid config key format, use section.key")
	}
	
	section := parts[0]
	field := parts[1]
	
	switch section {
	case "shell":
		return setShellConfig(&config.Shell, field, value)
	case "ui":
		return setUIConfig(&config.UI, field, value)
	case "history":
		return setHistoryConfig(&config.History, field, value)
	case "completion":
		return setCompletionConfig(&config.Completion, field, value)
	case "ai":
		return setAIConfig(&config.AI, field, value)
	default:
		return fmt.Errorf("unknown config section: %s", section)
	}
}

func setShellConfig(config *ShellConfig, field, value string) error {
	switch field {
	case "prompt":
		config.Prompt = value
	case "safe_mode":
		if b, err := strconv.ParseBool(value); err == nil {
			config.SafeMode = b
		} else {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
	case "default_theme":
		config.DefaultTheme = value
	case "execution_timeout":
		if i, err := strconv.Atoi(value); err == nil {
			config.ExecutionTimeout = i
		} else {
			return fmt.Errorf("invalid integer value: %s", value)
		}
	case "executor_mode":
		if value == "shell" || value == "internal" {
			config.ExecutorMode = value
		} else {
			return fmt.Errorf("invalid executor mode: %s (must be 'shell' or 'internal')", value)
		}
	default:
		return fmt.Errorf("unknown shell config field: %s", field)
	}
	return nil
}

func setUIConfig(config *UIConfig, field, value string) error {
	switch field {
	case "animated_startup":
		if b, err := strconv.ParseBool(value); err == nil {
			config.AnimatedStartup = b
		} else {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
	case "glitch_effects":
		if b, err := strconv.ParseBool(value); err == nil {
			config.GlitchEffects = b
		} else {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
	case "color_support":
		config.ColorSupport = value
	case "font_size":
		if i, err := strconv.Atoi(value); err == nil {
			config.FontSize = i
		} else {
			return fmt.Errorf("invalid integer value: %s", value)
		}
	default:
		return fmt.Errorf("unknown UI config field: %s", field)
	}
	return nil
}

func setHistoryConfig(config *HistoryConfig, field, value string) error {
	switch field {
	case "enabled":
		if b, err := strconv.ParseBool(value); err == nil {
			config.Enabled = b
		} else {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
	case "max_entries":
		if i, err := strconv.Atoi(value); err == nil {
			config.MaxEntries = i
		} else {
			return fmt.Errorf("invalid integer value: %s", value)
		}
	case "save_on_exit":
		if b, err := strconv.ParseBool(value); err == nil {
			config.SaveOnExit = b
		} else {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
	case "file_path":
		config.FilePath = value
	default:
		return fmt.Errorf("unknown history config field: %s", field)
	}
	return nil
}

func setCompletionConfig(config *CompletionConfig, field, value string) error {
	switch field {
	case "enabled":
		if b, err := strconv.ParseBool(value); err == nil {
			config.Enabled = b
		} else {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
	case "suggestions_per_page":
		if i, err := strconv.Atoi(value); err == nil {
			config.SuggestionsPerPage = i
		} else {
			return fmt.Errorf("invalid integer value: %s", value)
		}
	case "auto_cd":
		if b, err := strconv.ParseBool(value); err == nil {
			config.AutoCD = b
		} else {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
	case "fuzzy_matching":
		if b, err := strconv.ParseBool(value); err == nil {
			config.FuzzyMatching = b
		} else {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
	default:
		return fmt.Errorf("unknown completion config field: %s", field)
	}
	return nil
}

func setAIConfig(config *AIConfig, field, value string) error {
	switch field {
	case "enabled":
		if b, err := strconv.ParseBool(value); err == nil {
			config.Enabled = b
		} else {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
	case "suggest_commands":
		if b, err := strconv.ParseBool(value); err == nil {
			config.SuggestCommands = b
		} else {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
	case "confidence_threshold":
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			config.ConfidenceThreshold = f
		} else {
			return fmt.Errorf("invalid float value: %s", value)
		}
	case "model_provider":
		config.ModelProvider = value
	case "api_key":
		config.ApiKey = value
	default:
		return fmt.Errorf("unknown AI config field: %s", field)
	}
	return nil
}