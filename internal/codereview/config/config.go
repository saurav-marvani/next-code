package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the nextcode configuration
type Config struct {
	Version string                     `yaml:"version"`
	Rules   map[string]interface{}    `yaml:"rules"`
	Linters LinterConfig              `yaml:"linters"`
	Fixes   FixConfig                 `yaml:"fixes"`
	LLM     LLMConfig                 `yaml:"llm"`
	Slack   SlackConfig               `yaml:"slack"`
	Reports ReportConfig              `yaml:"reports"`
}

// LinterConfig configures linters
type LinterConfig struct {
	Enabled  []string               `yaml:"enabled"`
	Disabled []string               `yaml:"disabled"`
	Settings map[string]interface{} `yaml:"settings"`
}

// FixConfig configures auto-fixing
type FixConfig struct {
	Enabled     bool                   `yaml:"enabled"`
	Strategy    string                 `yaml:"strategy"` // "linter", "ai", "hybrid"
	AutoCommit  bool                   `yaml:"auto_commit"`
	Settings    map[string]interface{} `yaml:"settings"`
}

// LLMConfig configures LLM for AI fixes
type LLMConfig struct {
	Provider      string                 `yaml:"provider"` // "openai", "anthropic", "google"
	Model         string                 `yaml:"model"`
	Temperature   float32                `yaml:"temperature"`
	MaxTokens     int                    `yaml:"max_tokens"`
	CostPerToken  float32                `yaml:"cost_per_token"`
	Routing       map[string]string      `yaml:"routing"` // Task-specific model routing
}

// SlackConfig configures Slack integration
type SlackConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
	Channel    string `yaml:"channel"`
}

// ReportConfig configures report generation
type ReportConfig struct {
	Enabled bool                   `yaml:"enabled"`
	Format  []string               `yaml:"format"` // markdown, html, json
	Include map[string]interface{} `yaml:"include"`
}

// LoadConfig loads configuration from file
func LoadConfig(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// LoadConfigHierarchy loads config from multiple levels (global, project, directory)
func LoadConfigHierarchy(projectRoot string, targetDir string) (*Config, error) {
	// Start with defaults
	config := DefaultConfig()

	// Load global config
	globalConfig := filepath.Join(projectRoot, ".nextcode.yaml")
	if globalConfg, err := LoadConfig(globalConfig); err == nil {
		config = MergeConfigs(config, globalConfg)
	}

	// Load project-level config
	projectConfig := filepath.Join(projectRoot, "nextcode.yaml")
	if projConfig, err := LoadConfig(projectConfig); err == nil {
		config = MergeConfigs(config, projConfig)
	}

	// Load directory-specific config
	if targetDir != projectRoot {
		dirConfig := filepath.Join(targetDir, ".nextcode.yaml")
		if dConfig, err := LoadConfig(dirConfig); err == nil {
			config = MergeConfigs(config, dConfig)
		}
	}

	return config, nil
}

// SaveConfig saves configuration to file
func SaveConfig(path string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Version: "1.0",
		Rules:   make(map[string]interface{}),
		Linters: LinterConfig{
			Enabled: []string{"eslint", "pylint", "golangci-lint"},
			Settings: map[string]interface{}{
				"auto-detect": true,
			},
		},
		Fixes: FixConfig{
			Enabled:    true,
			Strategy:   "hybrid",
			AutoCommit: false,
		},
		LLM: LLMConfig{
			Provider:    "openai",
			Model:       "gpt-4",
			Temperature: 0.5,
			MaxTokens:   2000,
		},
		Reports: ReportConfig{
			Enabled: true,
			Format:  []string{"markdown"},
		},
	}
}

// MergeConfigs merges two configurations (right overrides left)
func MergeConfigs(left, right *Config) *Config {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}

	// Merge linters
	if len(right.Linters.Enabled) > 0 {
		left.Linters.Enabled = right.Linters.Enabled
	}
	for k, v := range right.Linters.Settings {
		left.Linters.Settings[k] = v
	}

	// Merge fixes
	left.Fixes = right.Fixes

	// Merge LLM
	if right.LLM.Provider != "" {
		left.LLM = right.LLM
	}

	// Merge reports
	if len(right.Reports.Format) > 0 {
		left.Reports = right.Reports
	}

	// Merge rules
	for k, v := range right.Rules {
		left.Rules[k] = v
	}

	return left
}
