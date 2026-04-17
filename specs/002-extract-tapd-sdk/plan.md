# 实施计划: 提取 tapd-sdk 独立库

**分支**: `002-extract-tapd-sdk` | **日期**: 2026-04-17 | **规范**: [spec.md](./spec.md)
**输入**: 来自 `/specs/002-extract-tapd-sdk/spec.md` 的功能规范

## 摘要

将 `internal/client`（TAPD HTTP 客户端）和 `internal/model`（数据模型与请求参数）从 CLI 内部迁移为独立的 `tapd-sdk-go` 包。CLI 重构为 SDK 的消费者，仅保留命令解析、配置管理和 AI 优化输出等 CLI 特有逻辑。SDK 以同仓库子目录方式发布（独立 `go.mod`），路径为 `sdk/`，模块名 `github.com/studyzy/tapd-sdk-go`。

## 技术背景

**语言/版本**: Go 1.24+
**主要依赖**:
- SDK: `net/http`（标准库）、`encoding/json`（标准库）
- CLI: `github.com/spf13/cobra`、`github.com/JohannesKaufmann/html-to-markdown/v2`（HTML→Markdown 转换保留在 CLI 层）
**存储**: N/A（SDK 无状态，CLI 使用 `.tapd.json` 本地配置文件）
**测试**: `go test`（标准库）+ `net/http/httptest`（Mock 服务器）
**目标平台**: Linux/macOS/Windows（CLI 工具）
**项目类型**: SDK 库 + CLI 工具（同仓库双模块）
**性能目标**: SDK 不引入额外延迟，网络 I/O 为瓶颈
**约束条件**: SDK 依赖项仅限标准库；CLI 现有命令输出格式不变；覆盖率 ≥ 60%

## 章程检查

*门控: 必须在阶段 0 研究前通过。阶段 1 设计后重新检查。*

| 章程原则 | 检查项 | 状态 | 说明 |
|---------|--------|------|------|
| I. API 优先 | SDK 封装所有 TAPD HTTP 调用，CLI 不直接构造 HTTP 请求 | ✅ 通过 | SDK 提供统一客户端，CLI 只调用 SDK 方法 |
| II. AI 优化输出 | CLI 继续输出 JSON，HTML→Markdown 转换保留在 CLI 层 | ✅ 通过 | SDK 返回原始结构体，CLI 负责格式化 |
| III. Go 编码规范 | 新 SDK 包遵循 gofmt/goimports/驼峰命名/error 最后返回 | ✅ 通过 | SDK panic 禁止，结构化错误返回 |
| IV. 测试纪律 | SDK 覆盖率 ≥ 60%，使用 httptest Mock | ✅ 通过 | 迁移现有测试 + 新增 SDK 独立测试 |
| V. 中文文档 | SDK 所有导出符号有中文注释，错误字符串用英文 | ✅ 通过 | 迁移现有注释风格 |

**门控结论**: 无违规，可进入阶段 0。

## 项目结构

### 文档（此功能）

```
specs/002-extract-tapd-sdk/
├── plan.md              # 此文件
├── research.md          # 阶段 0 输出
├── data-model.md        # 阶段 1 输出
├── quickstart.md        # 阶段 1 输出
├── contracts/           # 阶段 1 输出
│   └── sdk-api.md
└── tasks.md             # 阶段 2 输出（/speckit.tasks 命令）
```

### 源代码（重构后目录结构）

```
# 新增: tapd-sdk-go 子模块
sdk/                           # SDK 根目录（独立 go.mod）
├── go.mod                     # module github.com/studyzy/tapd-sdk-go
├── go.sum
├── client.go                  # Client 结构体、认证、HTTP 基础方法
├── error.go                   # TAPDError 类型
├── story.go                   # Story CRUD 方法
├── bug.go                     # Bug CRUD 方法
├── task.go                    # Task CRUD 方法
├── iteration.go               # Iteration 方法
├── comment.go                 # Comment 方法
├── wiki.go                    # Wiki 方法
├── attachment.go              # Attachment 方法
├── relation.go                # Relation 方法
├── timesheet.go               # Timesheet 方法
├── workflow.go                # Workflow 方法
├── workspace.go               # Workspace 方法
├── category.go                # Category 方法
├── custom_field.go            # CustomField 方法
├── misc.go                    # 杂项（Release、Todo、CommitMsg）
├── model/                     # 数据模型子包
│   ├── model.go               # 通用响应结构
│   ├── story.go               # Story 模型 + 请求参数
│   ├── bug.go
│   ├── task.go
│   ├── iteration.go
│   ├── comment.go
│   ├── wiki.go
│   ├── attachment.go
│   ├── relation.go
│   ├── timesheet.go
│   ├── workflow.go
│   ├── workspace.go
│   ├── category.go
│   ├── release.go
│   └── request.go             # 通用请求参数（WorkspaceIDRequest 等）
├── *_test.go                  # 与实现同包测试
└── go.sum

# 修改: CLI 重构（不改变外部接口）
internal/
├── cmd/                       # 不变（cobra 命令实现）
│   └── *.go                   # 改为导入 tapd-sdk-go 而非 internal/client
├── output/                    # 不变（输出格式化）
└── config/                    # 不变（配置文件管理）
# 删除: internal/client/ 和 internal/model/（内容迁移至 sdk/）

go.mod                         # 添加 tapd-sdk-go 本地路径依赖
```

**结构决策**: 采用同仓库子目录（`sdk/` 目录 + 独立 `go.mod`）。理由见 research.md。

## 复杂度跟踪

| 决策 | 为什么需要 | 拒绝更简单替代方案的原因 |
|------|-----------|------------------------|
| 独立 go.mod（`sdk/go.mod`） | SDK 需要可被第三方独立引入，无需依赖 CLI 的 cobra 等 | 若共用 go.mod，第三方引入 SDK 会拉入 cobra 等 CLI 依赖 |
| HTML→Markdown 留在 CLI | SDK 返回原始数据，格式转换是 CLI/AI 优化关注点 | 避免 SDK 引入 html-to-markdown 依赖，保持 SDK 最小依赖 |
