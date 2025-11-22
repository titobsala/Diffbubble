package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the user configuration
type Config struct {
	Theme       string      `yaml:"theme"`
	LineNumbers bool        `yaml:"line_numbers"`
	ContextMode string      `yaml:"context_mode"` // "focus" or "full"
	DiffMode    string      `yaml:"diff_mode"`    // "all", "staged", "unstaged"
	KeyBindings KeyBindings `yaml:"key_bindings,omitempty"`
}

// KeyBindings defines custom key bindings
type KeyBindings struct {
	Search            string `yaml:"search"`
	NextFile          string `yaml:"next_file"`
	PrevFile          string `yaml:"prev_file"`
	ToggleLineNumbers string `yaml:"toggle_line_numbers"`
	ToggleContext     string `yaml:"toggle_context"`
	Quit              string `yaml:"quit"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Theme:       "dark",
		LineNumbers: true,
		ContextMode: "focus",
		DiffMode:    "all",
		KeyBindings: DefaultKeyBindings(),
	}
}

// DefaultKeyBindings returns the default key bindings
func DefaultKeyBindings() KeyBindings {
	return KeyBindings{
		Search:            "/",
		NextFile:          "j",
		PrevFile:          "k",
		ToggleLineNumbers: "n",
		ToggleContext:     "c",
		Quit:              "q",
	}
}

// UserConfigPath returns the path to the user config file
func UserConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "diffbubble", "config.yaml")
}

// RepoConfigPath returns the path to the repository config file
func RepoConfigPath() string {
	return ".diffbubble.yml"
}

// Load loads configuration from disk, merging user and repo configs
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// Load user config (~/.config/diffbubble/config.yaml)
	if data, err := os.ReadFile(UserConfigPath()); err == nil {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			// If user config exists but is invalid, return error
			return nil, err
		}
	}

	// Override with repo config (.diffbubble.yml)
	if data, err := os.ReadFile(RepoConfigPath()); err == nil {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			// If repo config exists but is invalid, return error
			return nil, err
		}
	}

	return &cfg, nil
}

// Save saves the configuration to the user config file
func (c *Config) Save() error {
	configDir := filepath.Dir(UserConfigPath())
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(UserConfigPath(), data, 0644)
}

// Validate checks if the configuration values are valid
func (c *Config) Validate() error {
	// Validate theme (will be checked against ui.ValidateTheme in main)
	// Validate context mode
	if c.ContextMode != "focus" && c.ContextMode != "full" {
		c.ContextMode = "focus" // fallback to default
	}

	// Validate diff mode
	if c.DiffMode != "all" && c.DiffMode != "staged" && c.DiffMode != "unstaged" {
		c.DiffMode = "all" // fallback to default
	}

	return nil
}
