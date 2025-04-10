## Auto commit by ai with the all changes format.

![alt text](docs/image.png)

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

### Usage

Ensure that the required environment variables are set before running the application. You can override the defaults by setting them in your environment.

Example:
```sh
export OPENAI_API_KEY="your_api_key_here"
export AI_TYPE="openai"
```

Or using a different model with command line:
```sh
ai-git --model="deepseek" commit
```

Manual editing mode:
```sh
ai-git -m commit
```

Set a custom API endpoint:
```sh
export OPENAI_BASE_URL="https://your-custom-openai-endpoint.com/v1/chat/completions"
ai-git commit
```

Git-related code is referenced from https://github.com/jatinsandilya/mcp-server-auto-commit. Thanks to jatinsandilya!