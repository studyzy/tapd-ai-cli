# CODEBUDDY.md

本文件为 CodeBuddy Code 在此仓库中工作时提供指导。

## 项目简介

tapd-ai-cli 是一个 Go CLI 工具，通过 TAPD（腾讯敏捷产品研发平台）Open API 与 TAPD 平台交互。目标用户是 **AI Agent**（Claude Code、GPT-4o 等），而非人类用户。所有输出均针对最小 token 消耗进行优化。

## 技术栈

- Go 1.24+，CLI 框架：`spf13/cobra`
- SDK：`github.com/studyzy/tapd-sdk-go`（独立仓库 https://github.com/studyzy/tapd-sdk-go ，仅依赖标准库）
- HTTP 客户端：`net/http`（标准库，封装在 SDK 中）
- HTML 转 Markdown：`github.com/JohannesKaufmann/html-to-markdown/v2`（CLI 层，SDK 返回原始 HTML）
- 认证方式：`TAPD_ACCESS_TOKEN`（Bearer，推荐）或 `TAPD_API_USER`/`TAPD_API_PASSWORD`（Basic Auth）；凭据存储于 `~/.tapd.json` 或 `./.tapd.json`
- 许可证：Apache 2.0

## 构建与开发命令

```bash
# 构建二进制
make build
# 等价于：go build -o tapd ./cmd/tapd/

# 安装到 $GOPATH/bin
make install
# 等价于：go install ./cmd/tapd/

# 运行所有测试
make test
# 等价于：go test ./...

# 运行单个测试
go test ./path/to/package -run TestFunctionName

# 测试覆盖率（须 >= 60%）
make coverage
# 等价于：
# go test ./... -coverprofile=coverage.out
# go tool cover -func=coverage.out

# 代码格式化与检查
make fmt    # gofmt + goimports
make lint   # go vet + goimports -l

# 清理产物
make clean
```

## 架构

### 目录结构

```
cmd/tapd/       # 入口 main.go（go install 目标）
internal/
  cmd/          # Cobra 命令定义，消费 SDK 接口
  config/       # 多来源凭据管理（CLI flags > env > ./.tapd.json > ~/.tapd.json）
  output/       # JSON/Markdown 输出格式化
docs/           # 需求文档
```

SDK 已独立为单独仓库：[github.com/studyzy/tapd-sdk-go](https://github.com/studyzy/tapd-sdk-go)，通过 `go get` 引入。本地开发联调可使用 `go.work`（已在 `.gitignore` 中忽略）。

### 命令树

```
tapd
├── auth login --api-user <user> --api-password <pwd> [--local]
├── workspace list | switch <id> | info
├── story list | show <id> | create | update <id> | count
├── task  list | show <id> | create | update <id> | count
├── bug   list | show <id> | create | update <id> | count
├── iteration list
└── iteration list

全局标志：--workspace-id <id>，--pretty
--help 输出紧凑命令参考卡（供 AI Agent 自发现）
```

### 核心设计原则

1. **SDK 层（github.com/studyzy/tapd-sdk-go）**：独立仓库和 Go 模块，封装所有 TAPD API HTTP 调用，返回原始数据（含 HTML），仅依赖标准库，可被第三方直接 `go get` 引入。
2. **命令层（internal/cmd/）**：只负责参数解析、HTML→Markdown 转换和输出格式化，通过 SDK 客户端调用 API。
3. **每条命令的执行流程**：参数解析 → SDK 方法调用 → HTML→Markdown 转换（如需）→ 格式化输出。
4. **输出格式**：默认所有结构体字段带 `omitempty` 的紧凑 JSON；`--pretty` 添加缩进便于人类阅读；列表默认截断为 10 条并附分页提示。
5. **错误处理**：明确的退出码（0=成功，1=认证错误，2=未找到，3=参数错误，4=API 错误）；错误信息输出至 stderr 并附可操作提示。
6. **凭据优先级**：CLI flags > 环境变量（`TAPD_ACCESS_TOKEN` 或 `TAPD_API_USER`/`TAPD_API_PASSWORD`）> `./.tapd.json` > `~/.tapd.json`。同级内 access_token 优先于 api_user/api_password，严格优先级，不合并。

## 代码规范

- 所有代码注释和文档使用**中文**
- 错误消息字符串和日志输出使用**英文**（便于 AI 解析）
- 每个导出的函数、结构体和接口必须有中文文档注释
- 包注释用中文描述包的用途
- 使用 `gofmt`/`goimports` 格式化，遵循 Go 导出规则的驼峰命名
- 业务逻辑中不使用 `panic` 或 `goto`
- 函数不超过 80 行，文件不超过 800 行，嵌套不超过 4 层
- 错误作为最后一个返回值，必须处理或显式忽略
- `switch` 语句必须有 `default` 分支

## 测试要求

- 每个重要的导出函数必须有单元测试
- SDK 层：独立仓库，使用 `net/http/httptest` 进行 mock server 测试（见 https://github.com/studyzy/tapd-sdk-go ）
- 命令层：测试参数解析和输出格式
- 测试文件：`xxx_test.go`，测试函数：`TestXxx`
- 覆盖率目标：>= 60%

## 参考文档

- 需求规格：`docs/requirement.md`
- 项目章程：`.specify/memory/constitution.md`
- MVP 实施计划：`specs/001-mvp-tapd-cli/plan.md`
- CLI 命令契约：`specs/001-mvp-tapd-cli/contracts/cli-commands.md`
- 数据模型：`specs/001-mvp-tapd-cli/data-model.md`
- SDK 提取规范：`specs/002-extract-tapd-sdk/spec.md`
- SDK API 合约：`specs/002-extract-tapd-sdk/contracts/sdk-api.md`
