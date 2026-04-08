# howdoi

AI-powered command-line assistant. Get shell command suggestions and explanations using Claude, GitHub Models, or Ollama.

## Install

```bash
go install github.com/WMcKibbin/howdoi@latest
```

Or download a binary from [Releases](https://github.com/WMcKibbin/howdoi/releases).

## Usage

### Suggest a command

```bash
howdoi suggest "find all large files over 100MB"
```

This opens an interactive menu where you can:

- **Execute** the suggested command
- **Revise** it with additional instructions
- **Explain** what the command does
- **Copy** it to your clipboard

### Explain a command

```bash
howdoi explain "tar -xzf archive.tar.gz"
```

### Configure

```bash
howdoi config
```

Sets up your preferred AI provider and credentials. Config is stored at `~/.config/howdoi/config.yaml`.

### Shell aliases

```bash
# Add to your .zshrc / .bashrc:
eval "$(howdoi alias)"
```

This creates shortcuts:

- `hdi` → `howdoi suggest`
- `hde` → `howdoi explain`

## Providers

| Provider             | Auth                          | Setup                                                                                |
| -------------------- | ----------------------------- | ------------------------------------------------------------------------------------ |
| **Claude** (default) | Claude CLI auth               | Install [Claude CLI](https://docs.anthropic.com/en/docs/claude-cli) and authenticate |
| **GitHub Models**    | GitHub PAT with `models:read` | `howdoi config` or set `GITHUB_TOKEN` env var                                        |
| **Ollama**           | None (local)                  | Install and run [Ollama](https://ollama.ai)                                          |

Use `--provider` to switch:

```bash
howdoi suggest --provider github "list docker containers"
howdoi suggest --provider ollama "disk usage by directory"
```

Use `--model` to override the default model:

```bash
howdoi suggest --provider github --model "openai/gpt-4.1" "compress a folder"
```

## License

MIT
