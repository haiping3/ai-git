## AI-Git: Git CLI with AI Capabilities

![alt text](docs/image.png)

AI-Git is a Git wrapper that enhances standard Git commands with AI capabilities. It can generate commit messages and branch names based on your changes, making your Git workflow more efficient.

## Features

- **AI-assisted Commit Messages**: Automatically generate meaningful commit messages based on your changes
- **AI-assisted Branch Names**: Create descriptive branch names based on your changes
- **Full Git Compatibility**: Works with all standard Git commands, falling back to native Git for unsupported commands

## Installation

```sh
# Clone the repository
git clone https://github.com/Codexiaoyi/ai-git.git

# Build and install
cd ai-git
go build -o ai-git .
sudo mv ai-git /usr/local/bin/
```

## Usage

### Basic Usage

AI-Git works as a drop-in replacement for Git. Simply use `ai-git` instead of `git` for any command:

```sh
# AI-assisted commit (automatically generates a commit message)
ai-git commit

# AI-assisted branch creation (automatically generates a branch name)
ai-git checkout -b

# Use any other Git command (falls back to native Git)
ai-git status
ai-git push
ai-git pull
ai-git log
```

### Manual Editing Mode

You can manually edit the AI-generated commit messages and branch names:

```sh
# Commit with manual editing of the message
ai-git -m commit

# Checkout with manual editing of the branch name
ai-git -b checkout
```

## Configuration

The application supports multiple AI models, including OpenAI, Ollama, Anthropic, DeepSeek, and Qwen. The configuration is managed using environment variables with default values.

### Environment Variables

| Variable Name          | Default Value                                                         | Description                          |
|------------------------|---------------------------------------------------------------------|--------------------------------------|
| `AI_TYPE`              | `ollama`                                                            | Specifies the AI model type to use (`openai`, `ollama`, `anthropic`, `deepseek`, `qwen`)  |
| `AI_GIT_EDITOR`        | `$EDITOR` or `vim`                                                  | Editor to use for manual editing mode |
| `OPENAI_API_KEY`       | `""`                                                                | OpenAI API key                      |
| `OPENAI_MODEL`         | `gpt-3.5-turbo`                                                     | OpenAI model to be used             |
| `OPENAI_BASE_URL`      | `https://api.openai.com/v1/chat/completions`                        | OpenAI API endpoint URL             |
| `OLLAMA_BASE_URL`      | `http://localhost:11434`                                            | Base URL for Ollama                 |
| `OLLAMA_MODEL`         | `qwen2.5:7b`                                                        | Ollama model to be used             |
| `ANTHROPIC_API_KEY`    | `""`                                                                | Anthropic API key                   |
| `ANTHROPIC_MODEL`      | `claude-3-opus-20240229`                                            | Anthropic model to be used          |
| `ANTHROPIC_BASE_URL`   | `https://api.anthropic.com/v1/messages`                             | Anthropic API endpoint URL          |
| `DEEPSEEK_API_KEY`     | `""`                                                                | DeepSeek API key                    |
| `DEEPSEEK_MODEL`       | `deepseek-chat`                                                     | DeepSeek model to be used           |
| `DEEPSEEK_BASE_URL`    | `https://api.deepseek.com/v1/chat/completions`                      | DeepSeek API endpoint URL           |
| `QWEN_API_KEY`         | `""`                                                                | Qwen API key                        |
| `QWEN_MODEL`           | `qwen-max`                                                          | Qwen model to be used               |
| `QWEN_BASE_URL`        | `https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation` | Qwen API endpoint URL |

### Configuration Examples

Set up your AI model type and API key:
```sh
export OPENAI_API_KEY="your_api_key_here"
export AI_TYPE="openai"
```

## How It Works

When you run an AI-Git command:

1. If it's a supported command (like `commit` or `checkout -b`), AI-Git uses AI to enhance the command's behavior
2. If it's an unsupported command, AI-Git passes it through to the native Git command

This means you can use AI-Git for your entire Git workflow without having to switch between different commands.

## License

[License details]

## Acknowledgements

Git-related code is referenced from https://github.com/jatinsandilya/mcp-server-auto-commit. Thanks to jatinsandilya!