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
	ModelDeepSeek  ModelType = "deepseek"
	ModelQwen      ModelType = "qwen"
)

// Config holds the configuration for AI models
type Config struct {
	Type      ModelType       `yaml:"type" json:"type"`
	OpenAI    OpenAIConfig    `yaml:"openai,omitempty" json:"openai,omitempty"`
	Ollama    OllamaConfig    `yaml:"ollama,omitempty" json:"ollama,omitempty"`
	Anthropic AnthropicConfig `yaml:"anthropic,omitempty" json:"anthropic,omitempty"`
	DeepSeek  DeepSeekConfig  `yaml:"deepseek,omitempty" json:"deepseek,omitempty"`
	Qwen      QwenConfig      `yaml:"qwen,omitempty" json:"qwen,omitempty"`
}

// OpenAIConfig holds OpenAI-specific configuration
type OpenAIConfig struct {
	APIKey  string `yaml:"api_key" json:"api_key"`
	Model   string `yaml:"model" json:"model"`
	BaseURL string `yaml:"base_url" json:"base_url"`
}

// OllamaConfig holds Ollama-specific configuration
type OllamaConfig struct {
	BaseURL string `yaml:"base_url" json:"base_url"`
	Model   string `yaml:"model" json:"model"`
}

// AnthropicConfig holds Anthropic-specific configuration
type AnthropicConfig struct {
	APIKey  string `yaml:"api_key" json:"api_key"`
	Model   string `yaml:"model" json:"model"`
	BaseURL string `yaml:"base_url" json:"base_url"`
}

// DeepSeekConfig holds DeepSeek-specific configuration
type DeepSeekConfig struct {
	APIKey  string `yaml:"api_key" json:"api_key"`
	Model   string `yaml:"model" json:"model"`
	BaseURL string `yaml:"base_url" json:"base_url"`
}

// QwenConfig holds Qwen-specific configuration
type QwenConfig struct {
	APIKey  string `yaml:"api_key" json:"api_key"`
	Model   string `yaml:"model" json:"model"`
	BaseURL string `yaml:"base_url" json:"base_url"`
}

// LoadConfig loads the configuration from the specified file
func LoadConfig() (*Config, error) {
	config := Config{
		Type: ModelType(getEnvWithDefault("AI_TYPE", "ollama")),
		OpenAI: OpenAIConfig{
			APIKey:  getEnvWithDefault("OPENAI_API_KEY", ""),
			Model:   getEnvWithDefault("OPENAI_MODEL", "gpt-3.5-turbo"),
			BaseURL: getEnvWithDefault("OPENAI_BASE_URL", "https://api.openai.com/v1/chat/completions"),
		},
		Ollama: OllamaConfig{
			BaseURL: getEnvWithDefault("OLLAMA_BASE_URL", "http://localhost:11434"),
			Model:   getEnvWithDefault("OLLAMA_MODEL", "qwen2.5:7b"),
		},
		Anthropic: AnthropicConfig{
			APIKey:  getEnvWithDefault("ANTHROPIC_API_KEY", ""),
			Model:   getEnvWithDefault("ANTHROPIC_MODEL", "claude-3-opus-20240229"),
			BaseURL: getEnvWithDefault("ANTHROPIC_BASE_URL", "https://api.anthropic.com/v1/messages1"),
		},
		DeepSeek: DeepSeekConfig{
			APIKey:  getEnvWithDefault("DEEPSEEK_API_KEY", ""),
			Model:   getEnvWithDefault("DEEPSEEK_MODEL", "deepseek-chat"),
			BaseURL: getEnvWithDefault("DEEPSEEK_BASE_URL", "https://api.deepseek.com/v1/chat/completions"),
		},
		Qwen: QwenConfig{
			APIKey:  getEnvWithDefault("QWEN_API_KEY", ""),
			Model:   getEnvWithDefault("QWEN_MODEL", "qwen-max"),
			BaseURL: getEnvWithDefault("QWEN_BASE_URL", "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"),
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
	case ModelDeepSeek:
		if config.DeepSeek.Model == "" {
			config.DeepSeek.Model = "deepseek-chat"
		}
	case ModelQwen:
		if config.Qwen.Model == "" {
			config.Qwen.Model = "qwen-max"
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
