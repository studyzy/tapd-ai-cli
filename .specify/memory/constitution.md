<!--
同步影响报告
========================================
版本更改: 无(初始创建) → 1.0.0
修改的原则: 无(全部为新增)
添加的部分:
  - 核心原则 (5 项)
  - 技术栈与约束
  - 开发工作流
  - 治理
删除的部分: 无
需要更新的模板:
  - .specify/templates/plan-template.md ✅ 已检查(章程检查部分兼容)
  - .specify/templates/spec-template.md ✅ 已检查(无需修改)
  - .specify/templates/tasks-template.md ✅ 已检查(测试任务模式兼容)
后续 TODO: 无
========================================
-->
# tapd-ai-cli 项目章程

## 核心原则

### I. API 优先
本项目是一个 CLI 工具，所有功能 MUST 通过调用 TAPD 平台 API 实现。
API 交互层是项目的核心，所有命令 MUST 映射到具体的 TAPD API 端点。
HTTP 客户端封装 MUST 统一处理认证、错误响应和重试逻辑，
禁止在命令层直接构造 HTTP 请求。

### II. AI 优化输出
本工具的目标用户是 AI Agent（如 Claude Code、GPT-4o），而非人类用户。
所有命令输出 MUST 默认为结构化 JSON 格式。
输出 MUST 遵循 Zero Noise 原则：移除 null 字段（`omitempty`）、
支持 `--compact` 紧凑模式、长列表默认截断并提示分页。
错误信息 MUST 输出到 stderr 并包含可操作的修复建议。

### III. Go 编码规范
所有代码 MUST 符合 Golang 业界规范，以 Google Go 风格指南为基准。
具体要求：
- MUST 使用 `gofmt` 和 `goimports` 格式化代码
- MUST 使用驼峰命名，遵循 Go 导出规则
- MUST 正确处理 error，error 作为函数最后一个返回参数
- MUST 禁止在业务代码中使用 panic 和 goto
- SHOULD 保持函数不超过 80 行，文件不超过 800 行
- SHOULD 控制嵌套深度不超过 4 层
- SHOULD 使用 Go Modules 管理依赖

### IV. 测试纪律
核心和必要的代码 MUST 编写单元测试。
单元测试覆盖率 MUST 达到 60% 以上。
具体要求：
- 测试文件命名为 `xxx_test.go`
- 测试函数以 `Test` 开头
- 每个重要的可导出函数 MUST 有对应测试用例
- API 客户端层 SHOULD 使用 httptest 进行测试
- 命令层 SHOULD 测试参数解析和输出格式

### V. 中文文档与注释
所有文档和代码注释 MUST 使用中文。
具体要求：
- 每个导出的函数、结构体、接口 MUST 有中文注释
- 包注释 MUST 使用中文描述包的功能
- README 和其他文档 MUST 使用中文编写
- 错误信息字符串 SHOULD 使用英文（便于 AI Agent 解析）
- 日志输出 SHOULD 使用英文（便于机器处理）

## 技术栈与约束

- **语言**: Golang 1.24+
- **CLI 框架**: `spf13/cobra`
- **HTTP 客户端**: `go-resty/resty` 或标准 `net/http`
- **配置管理**: `spf13/viper`，凭据存储于 `~/.tapd/config.json`
- **项目类型**: 命令行工具（CLI）
- **架构模式**: 无状态请求-响应，Workspace ID 通过本地配置持久化
- **输出格式**: 默认 JSON，支持紧凑模式和人类可读模式
- **许可证**: Apache License 2.0

## 开发工作流

- 所有代码变更 MUST 通过 `go vet` 和 `golint` 静态检查
- 提交前 MUST 确保 `go test ./...` 全部通过
- 新增功能 MUST 同步编写对应的单元测试
- 代码审查 SHOULD 关注 API 调用的错误处理和边界情况
- 每个 TAPD API 命令的实现 SHOULD 遵循统一的模式：
  参数解析 → API 调用 → 响应转换 → 格式化输出

## 治理

本章程是 tapd-ai-cli 项目的最高指导文件，
所有开发实践 MUST 与章程原则保持一致。

**修改程序**：
- 章程修改 MUST 记录变更内容和理由
- 版本号遵循语义版本控制：MAJOR（原则删除或重定义）、
  MINOR（新原则或实质性扩展）、PATCH（措辞澄清或修正）
- 修改后 MUST 更新版本号和最后修订日期

**合规审查**：
- 实施计划 MUST 包含章程检查，确认不违反核心原则
- 代码审查 SHOULD 验证是否符合编码规范和测试纪律要求

**版本**: 1.0.0 | **批准日期**: 2026-04-16 | **最后修订**: 2026-04-16
