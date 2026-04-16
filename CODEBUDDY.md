# CODEBUDDY.md

This file provides guidance to CodeBuddy Code when working with code in this repository.

## Project Overview

tapd-ai-cli is a Go CLI tool that interacts with the TAPD (Tencent Agile Product Development) platform via its API. The target users are **AI Agents** (Claude Code, GPT-4o), not humans. All output is optimized for minimal token usage.

## Tech Stack

- Go 1.24+, CLI framework: `spf13/cobra`, config: `spf13/viper`
- HTTP client: `net/http` (standard library)
- HTML to Markdown: `github.com/JohannesKaufmann/html-to-markdown/v2`
- Credentials: `TAPD_ACCESS_TOKEN` (Bearer, recommended) or `TAPD_API_USER`/`TAPD_API_PASSWORD` (Basic Auth); stored in `~/.tapd.json` or `./.tapd.json`
- License: Apache 2.0

## Build & Development Commands

```bash
# Build
go build -o tapd ./cmd/tapd/

# Run all tests
go test ./...

# Run a single test
go test ./path/to/package -run TestFunctionName

# Test coverage (must be >= 60%)
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Lint & format
gofmt -w .
goimports -w .
go vet ./...
```

## Architecture

### Planned Directory Structure

```
cmd/tapd/       # 入口 main.go (go install 目标)
internal/
  cmd/         # Cobra 命令定义 (root, auth, workspace, story, task, bug, iteration, spec)
  client/      # TAPD API HTTP 客户端封装 (认证、请求构造、响应解析)
  config/      # 多来源凭据管理 (CLI flags > env > ./.tapd.json > ~/.tapd.json)
  output/      # JSON 输出格式化 (默认紧凑、omitempty、列表截断)
  model/       # TAPD 数据模型 (Workspace, Story, Task, Bug, Iteration)
docs/          # 需求文档
```

### Command Tree

```
tapd
├── auth login --api-user <user> --api-password <pwd> [--local]
├── workspace list | switch <id> | info
├── story list | show <id> | create | update <id> | count
├── task list | show <id> | create | update <id> | count
├── bug list | show <id> | create | update <id> | count
├── iteration list
└── spec    # 输出 OpenAI/Anthropic 兼容的 Tool Definition JSON

Global flags: --workspace-id <id>, --pretty
```

### Key Design Patterns

1. **Command layer (internal/cmd/)**: Parameter parsing and output formatting only. MUST NOT construct HTTP requests directly.
2. **API client layer (internal/client/)**: Unified HTTP wrapper handling auth headers, error mapping. All TAPD API calls go through here.
3. **Execution flow per command**: Argument parsing -> API call -> Response transformation -> Formatted output.
4. **Output format**: Default compact JSON with `omitempty` on all structs. `--pretty` flag adds indentation for human reading. Lists truncated to 10 items by default with pagination hint.
5. **Error handling**: Distinct exit codes (0=success, 1=auth error, 2=not found, 3=param error, 4=API error). Errors to stderr with actionable hints.
6. **Credential resolution order**: CLI flags > env vars (`TAPD_ACCESS_TOKEN` or `TAPD_API_USER`/`TAPD_API_PASSWORD`) > `./.tapd.json` > `~/.tapd.json`. Within same level, access_token takes priority over api_user/api_password. Strict priority, no merging.

## Code Conventions

- All code comments and documentation in **Chinese** (中文)
- Error message strings and log output in **English** (for AI parsing)
- Every exported function, struct, and interface MUST have a Chinese doc comment
- Package comments in Chinese describing the package purpose
- Use `gofmt`/`goimports` formatting, camelCase naming per Go export rules
- No `panic` or `goto` in business logic
- Functions <= 80 lines, files <= 800 lines, nesting <= 4 levels
- Error as last return value, always handle or explicitly ignore
- `switch` statements must have `default` case

## Testing Requirements

- Every important exported function must have unit tests
- API client layer: use `net/http/httptest` for mock server testing
- Command layer: test argument parsing and output format
- Test files: `xxx_test.go`, test functions: `TestXxx`
- Coverage target: >= 60%

## Reference Documents

- Requirements spec: `docs/requirement.md`
- Project constitution: `.specify/memory/constitution.md`
- Implementation plan: `specs/001-mvp-tapd-cli/plan.md`
- CLI command contracts: `specs/001-mvp-tapd-cli/contracts/cli-commands.md`
- Data model: `specs/001-mvp-tapd-cli/data-model.md`
