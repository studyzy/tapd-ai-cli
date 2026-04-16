# 实施计划: tapd-ai-cli MVP

**分支**: `001-mvp-tapd-cli` | **日期**: 2026-04-16 | **规范**: [spec.md](./spec.md)
**输入**: 来自 `/specs/001-mvp-tapd-cli/spec.md` 的功能规范

## 摘要

构建 tapd-ai-cli 的 MVP 版本：一个面向 AI Agent 的 Go CLI 工具，通过 TAPD Open API 实现项目管理核心操作。MVP 覆盖认证（多来源凭据）、工作区管理、需求/任务 CRUD、缺陷 CRUD、迭代查询和 AI 自发现（spec 命令）。所有输出为精简 JSON 格式，支持紧凑模式。

## 技术背景

**语言/版本**: Go 1.24+
**主要依赖**:
- `github.com/spf13/cobra` — CLI 命令框架
- `github.com/spf13/viper` — 配置管理（多来源凭据解析）
- `net/http` — HTTP 客户端（标准库，轻量无额外依赖）
- `github.com/JohannesKaufmann/html-to-markdown/v2` — HTML 转 Markdown
**存储**: 本地 JSON 文件（`~/.tapd.json` 和 `./.tapd.json`）
**测试**: `go test`，API 层使用 `net/http/httptest`
**目标平台**: macOS / Linux / Windows（跨平台 CLI）
**项目类型**: CLI 工具
**性能目标**: N/A（CLI 工具，响应时间取决于 TAPD API）
**约束条件**: JSON 输出必须 `omitempty`，列表默认截断 10 条
**规模/范围**: 单用户 CLI，约 20 个子命令

## 章程检查

*门控: 必须在阶段 0 研究前通过. 阶段 1 设计后重新检查.*

| 原则 | 状态 | 说明 |
|------|------|------|
| I. API 优先 | ✅ 通过 | 所有功能通过 TAPD API 实现，HTTP 客户端统一封装在 `internal/client/` |
| II. AI 优化输出 | ✅ 通过 | 默认 JSON 输出，`omitempty`，`--compact` 模式，列表截断 |
| III. Go 编码规范 | ✅ 通过 | 使用 gofmt/goimports，标准 Go 项目结构 |
| IV. 测试纪律 | ✅ 通过 | 核心包（client, config, output）均需 UT，目标 ≥60% |
| V. 中文文档与注释 | ✅ 通过 | 所有导出符号使用中文注释，错误信息用英文 |

无章程违规，无需复杂度跟踪。

## 项目结构

### 文档(此功能)

```
specs/001-mvp-tapd-cli/
├── plan.md              # 此文件
├── research.md          # 阶段 0 输出
├── data-model.md        # 阶段 1 输出
├── quickstart.md        # 阶段 1 输出
├── contracts/           # 阶段 1 输出（CLI 命令合同）
└── tasks.md             # 阶段 2 输出 (/speckit.tasks)
```

### 源代码(仓库根目录)

```
tapd-ai-cli/
├── main.go                    # 入口
├── go.mod
├── go.sum
├── cmd/                       # Cobra 命令定义
│   ├── root.go                # 根命令，全局标志（--workspace-id, --compact）
│   ├── auth.go                # tapd auth login
│   ├── workspace.go           # tapd workspace list/switch/info
│   ├── story.go               # tapd story list/show/create/update/count
│   ├── task.go                # tapd task list/show/create/update/count
│   ├── bug.go                 # tapd bug list/show/create/update/count
│   ├── iteration.go           # tapd iteration list
│   └── spec.go                # tapd spec（Tool Definition JSON）
├── internal/
│   ├── client/                # TAPD API HTTP 客户端封装
│   │   ├── client.go          # Client 结构体、认证、请求/响应处理
│   │   ├── client_test.go
│   │   ├── story.go           # 需求相关 API 方法
│   │   ├── task.go            # 任务相关 API 方法
│   │   ├── bug.go             # 缺陷相关 API 方法
│   │   ├── workspace.go       # 工作区相关 API 方法
│   │   └── iteration.go       # 迭代相关 API 方法
│   ├── config/                # 配置管理（多来源凭据、workspace_id）
│   │   ├── config.go
│   │   └── config_test.go
│   ├── output/                # JSON 输出格式化（compact、omitempty、截断）
│   │   ├── output.go
│   │   └── output_test.go
│   └── model/                 # 数据模型（Workspace, Story, Bug, Iteration 等）
│       └── model.go
├── Makefile                   # 构建、测试、lint 命令
└── README.md                  # 项目说明
```

**结构决策**: 采用标准 Go CLI 项目布局。`cmd/` 放 Cobra 命令，`internal/` 放不对外暴露的业务逻辑。`internal/client/` 统一封装所有 TAPD API 调用，命令层只做参数解析和输出格式化。

## 复杂度跟踪

> 无章程违规，无需填写。
