package ai

import (
	"fmt"
	"os"
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
func LoadConfig() (*Config, error) {
	config := Config{
		Type: ModelType(getEnvWithDefault("AI_TYPE", "ollama")),
		OpenAI: OpenAIConfig{
			APIKey: getEnvWithDefault("OPENAI_API_KEY", ""),
			Model:  getEnvWithDefault("OPENAI_MODEL", "gpt-3.5-turbo"),
		},
		Ollama: OllamaConfig{
			BaseURL: getEnvWithDefault("OLLAMA_BASE_URL", "http://localhost:11434"),
			Model:   getEnvWithDefault("OLLAMA_MODEL", "qwen2.5:7b"),
		},
		Anthropic: AnthropicConfig{
			APIKey: getEnvWithDefault("ANTHROPIC_API_KEY", ""),
			Model:  getEnvWithDefault("ANTHROPIC_MODEL", "claude-3-opus-20240229"),
		},
	}

	// Set default values if needed
	if config.Type == "" {
		config.Type = ModelOllama
	}

	// Validate model type
	switch config.Type {
	case ModelOpenAI:
		if config.OpenAI.Model == "" {
			config.OpenAI.Model = "gpt-3.5-turbo"
		}
	case ModelOllama:
		if config.Ollama.Model == "" {
			config.Ollama.Model = "qwen2.5:7b"
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

// Helper to get environment variable with default fallback
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
