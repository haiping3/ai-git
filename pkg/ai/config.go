package ai

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ModelType represents the type of AI model
type ModelType string

const (
	ModelOpenAI    ModelType = "openai"
	ModelOllama    ModelType = "ollama"
	ModelAnthropic ModelType = "anthropic"
)

// Config holds the configuration for AI models
type Config struct {
	Type      ModelType       `yaml:"type" json:"type"`
	OpenAI    OpenAIConfig    `yaml:"openai,omitempty" json:"openai,omitempty"`
	Ollama    OllamaConfig    `yaml:"ollama,omitempty" json:"ollama,omitempty"`
	Anthropic AnthropicConfig `yaml:"anthropic,omitempty" json:"anthropic,omitempty"`
}

// OpenAIConfig holds OpenAI-specific configuration
type OpenAIConfig struct {
	APIKey string `yaml:"api_key" json:"api_key"`
	Model  string `yaml:"model" json:"model"`
}

// OllamaConfig holds Ollama-specific configuration
type OllamaConfig struct {
	BaseURL string `yaml:"base_url" json:"base_url"`
	Model   string `yaml:"model" json:"model"`
}

// AnthropicConfig holds Anthropic-specific configuration
type AnthropicConfig struct {
	APIKey string `yaml:"api_key" json:"api_key"`
	Model  string `yaml:"model" json:"model"`
}

// LoadConfig loads the configuration from the specified file
func LoadConfig(configPath string) (*Config, error) {
	// Default config path is config.yaml in the current directory
	if configPath == "" {
		configPath = "config.yaml"
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// If not found in current directory, try user's home directory
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("could not find config file and unable to determine home directory: %w", err)
		}

		altPath := filepath.Join(home, ".ai-git", "config.yaml")
		if _, err := os.Stat(altPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found at %s or %s", configPath, altPath)
		}

		configPath = altPath
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Set default values if needed
	if config.Type == "" {
		config.Type = ModelOpenAI
	}

	// Validate model type
	switch config.Type {
	case ModelOpenAI:
		if config.OpenAI.Model == "" {
			config.OpenAI.Model = "gpt-3.5-turbo"
		}
	case ModelOllama:
		if config.Ollama.Model == "" {
			config.Ollama.Model = "llama2"
		}
		if config.Ollama.BaseURL == "" {
			config.Ollama.BaseURL = "http://localhost:11434"
		}
	case ModelAnthropic:
		if config.Anthropic.Model == "" {
			config.Anthropic.Model = "claude-3-opus-20240229"
		}
	default:
		return nil, fmt.Errorf("unsupported model type: %s", config.Type)
	}

	return &config, nil
}
