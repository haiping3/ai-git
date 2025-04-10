package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIRequest represents the request structure for OpenAI API
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"` // Always set to false
}

// OpenAIResponse represents the response structure from OpenAI API
type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// OllamaRequest represents the request structure for Ollama API
type OllamaRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"` // Always set to false
}

// OllamaResponse represents the response structure from Ollama API
type OllamaResponse struct {
	Message Message `json:"message"`
}

// AnthropicRequest represents the request structure for Anthropic API
type AnthropicRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// AnthropicResponse represents the response structure from Anthropic API
type AnthropicResponse struct {
	Content []Message `json:"content"`
}

// DeepSeekRequest represents the request structure for DeepSeek API
type DeepSeekRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"` // Always set to false
}

// DeepSeekResponse represents the response structure from DeepSeek API
type DeepSeekResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// QwenRequest represents the request structure for Qwen API
type QwenRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// QwenResponse represents the response structure from Qwen API
type QwenResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// GenerateCommitMessage generates a commit message using the configured AI model
func GenerateCommitMessage(prompt string, config Config) (string, error) {
	switch config.Type {
	case ModelOpenAI:
		return generateWithOpenAI(prompt, config.OpenAI)
	case ModelOllama:
		return generateWithOllama(prompt, config.Ollama)
	case ModelAnthropic:
		return generateWithAnthropic(prompt, config.Anthropic)
	case ModelDeepSeek:
		return generateWithDeepSeek(prompt, config.DeepSeek)
	case ModelQwen:
		return generateWithQwen(prompt, config.Qwen)
	default:
		return "", fmt.Errorf("unsupported model type: %s", config.Type)
	}
}

// GenerateBranchName generates a branch name using the configured AI model
func GenerateBranchName(prompt string, config Config) (string, error) {
	switch config.Type {
	case ModelOpenAI:
		return generateWithOpenAI(prompt, config.OpenAI)
	case ModelOllama:
		return generateWithOllama(prompt, config.Ollama)
	case ModelAnthropic:
		return generateWithAnthropic(prompt, config.Anthropic)
	case ModelDeepSeek:
		return generateWithDeepSeek(prompt, config.DeepSeek)
	case ModelQwen:
		return generateWithQwen(prompt, config.Qwen)
	default:
		return "", fmt.Errorf("unsupported model type: %s", config.Type)
	}
}

// generateWithOpenAI generates a commit message using OpenAI
func generateWithOpenAI(prompt string, config OpenAIConfig) (string, error) {
	apiKey := config.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return "", fmt.Errorf("OpenAI API key is not set")
		}
	}

	reqBody := OpenAIRequest{
		Model: config.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant that generates concise and descriptive git commit message based on the changes provided. Please generate shortly.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false, // Explicitly disable streaming
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1/chat/completions"
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", err
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

// generateWithOllama generates a commit message using Ollama
func generateWithOllama(prompt string, config OllamaConfig) (string, error) {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	reqBody := OllamaRequest{
		Model: config.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant that generates concise and descriptive git commit message based on the changes provided. Please generate shortly.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false, // Explicitly disable streaming
	}

	jsonData, err := json.Marshal(&reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/chat", baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", err
	}

	return ollamaResp.Message.Content, nil
}

// generateWithAnthropic generates a commit message using Anthropic
func generateWithAnthropic(prompt string, config AnthropicConfig) (string, error) {
	apiKey := config.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			return "", fmt.Errorf("Anthropic API key is not set")
		}
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1/messages"
	}

	// Create structure that matches Anthropic API for Claude
	reqBody := struct {
		Model     string    `json:"model"`
		MaxTokens int       `json:"max_tokens"`
		Messages  []Message `json:"messages"`
		Stream    bool      `json:"stream"` // Explicitly disable streaming
	}{
		Model:     config.Model,
		MaxTokens: 1000,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant that generates concise and descriptive git commit message based on the changes provided. Please generate shortly.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse response which has a different structure
	var anthropicResp struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(body, &anthropicResp); err != nil {
		return "", err
	}

	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("no response from Anthropic")
	}

	return anthropicResp.Content[0].Text, nil
}

// generateWithDeepSeek generates a commit message using DeepSeek
func generateWithDeepSeek(prompt string, config DeepSeekConfig) (string, error) {
	apiKey := config.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("DEEPSEEK_API_KEY")
		if apiKey == "" {
			return "", fmt.Errorf("DeepSeek API key is not set")
		}
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/v1/chat/completions"
	}

	reqBody := DeepSeekRequest{
		Model: config.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant that generates concise and descriptive git commit message based on the changes provided. Please generate shortly.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false, // Explicitly disable streaming
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var deepSeekResp DeepSeekResponse
	if err := json.Unmarshal(body, &deepSeekResp); err != nil {
		return "", err
	}

	if len(deepSeekResp.Choices) == 0 {
		return "", fmt.Errorf("no response from DeepSeek")
	}

	return deepSeekResp.Choices[0].Message.Content, nil
}

// generateWithQwen generates a commit message using Qwen
func generateWithQwen(prompt string, config QwenConfig) (string, error) {
	apiKey := config.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("QWEN_API_KEY")
		if apiKey == "" {
			return "", fmt.Errorf("Qwen API key is not set")
		}
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
	}

	reqBody := QwenRequest{
		Model: config.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant that generates concise and descriptive git commit message based on the changes provided. Please generate shortly.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var qwenResp QwenResponse
	if err := json.Unmarshal(body, &qwenResp); err != nil {
		return "", err
	}

	if len(qwenResp.Choices) == 0 {
		return "", fmt.Errorf("no response from Qwen")
	}

	content := qwenResp.Choices[0].Message.Content
	_, newContent, found := strings.Cut(content, "</think>\n\n")
	if found {
		return newContent, nil
	}
	return content, nil
}
