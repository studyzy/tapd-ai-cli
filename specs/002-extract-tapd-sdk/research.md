# 研究报告: 提取 tapd-sdk 独立库

**分支**: `002-extract-tapd-sdk` | **日期**: 2026-04-17

## 决策 1: SDK 发布形式

**Decision**: 同仓库子目录（`sdk/` 目录 + 独立 `go.mod`），模块名 `github.com/studyzy/tapd-sdk-go`

**Rationale**:
- Go Modules 支持同仓库多模块（monorepo），通过 `replace` 指令在 CLI 的 `go.mod` 中引用本地路径，开发阶段无需发布即可联调
- 独立 `go.mod` 确保 SDK 的依赖图完全独立，第三方引入时不会拉入 `cobra`、`html-to-markdown` 等 CLI 专属依赖
- 相比独立仓库，同仓库子目录减少跨仓库同步负担，适合初期迭代

**Alternatives considered**:
- **独立仓库** (`github.com/studyzy/tapd-sdk-go`): 彻底隔离，但需要跨仓库 PR 联动，初期维护成本高
- **共用 `go.mod`** (internal 包改名): 最简单，但第三方无法引入 SDK（Go 规范禁止导入 `internal` 包）；即使改名也会拉入所有 CLI 依赖

---

## 决策 2: SDK 包结构（扁平 vs 分层）

**Decision**: 扁平结构——`sdk/` 根目录放客户端方法文件，`sdk/model/` 子包放数据模型

**Rationale**:
- 调用侧简洁：`sdk.NewClient(...)` 后直接 `client.ListStories(req)`，无需多层包导入
- `model` 子包分离数据结构，方便只使用模型定义（如做类型断言、序列化）而不实例化客户端
- 与现有 `internal/client` + `internal/model` 结构对应，迁移成本最低

**Alternatives considered**:
- **按资源类型分子包** (`tapd/story`, `tapd/bug`): 粒度过细，调用时需导入多个包，且各资源共享的 HTTP 基础层难以组织
- **单文件 SDK**: 所有代码放一个文件，不符合 Go 编码规范（文件 ≤ 800 行约束）

---

## 决策 3: HTML→Markdown 转换归属

**Decision**: HTML→Markdown 转换保留在 CLI 层（`internal/cmd/`），SDK 返回原始字符串

**Rationale**:
- SDK 职责是忠实传递 TAPD API 数据，格式转换是消费者（CLI/AI Agent）的关注点
- 避免 SDK 引入 `github.com/JohannesKaufmann/html-to-markdown` 依赖，保持 SDK 零非标准库依赖
- 不同消费者对 HTML 的处理需求不同（有些可能需要原始 HTML 或纯文本）

**Alternatives considered**:
- **SDK 内置转换**: 方便 CLI，但污染 SDK 依赖树，违背 SDK 最小依赖原则
- **可选转换（Option 模式）**: 过度设计，初期不需要

---

## 决策 4: 错误处理策略

**Decision**: SDK 返回结构化 `TAPDError`（含 HTTPStatus、ExitCode、Message），绝不调用 `os.Exit`

**Rationale**:
- CLI 工具可能调用 `os.Exit`，但 SDK 作为库绝对不能终止调用方进程
- 结构化错误允许调用方根据 HTTPStatus 或 ExitCode 做精细处理（如 401 触发重新认证）
- 现有 `TAPDError` 结构已满足需求，直接迁移即可

**Alternatives considered**:
- **纯字符串错误**: 丢失错误码信息，调用方无法区分 404 和 500
- **errors.Is/As 哨兵错误**: 适合枚举有限错误类型，但 TAPD 错误码多样，结构体更灵活

---

## 决策 5: CLI 对 SDK 的引用方式

**Decision**: CLI 的 `go.mod` 使用 `replace` 指令引用本地 `./sdk`，发布时改为版本引用

**Rationale**:
- 开发阶段无需发布 SDK 即可联调修改
- Go Modules 官方支持此模式，是 monorepo 标准做法
- CI 环境只需 `go work` 或保留 `replace` 指令即可正常构建

**Alternatives considered**:
- **git submodule**: 复杂度高，Go 生态不常用
- **vendor 模式**: 需要手动同步，维护成本高

---

## 技术风险与缓解措施

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| CLI 测试在迁移过程中破坏 | 高 | 先迁移，保持 internal/client 和 internal/model 的别名；所有测试通过后再删除 |
| SDK 意外引入 CLI 依赖 | 中 | `go mod tidy` + `go list -m all` 验证 SDK 依赖树 |
| 迁移遗漏部分 client 方法 | 中 | 逐文件对照迁移，迁移完成后用 grep 验证 internal/client 无遗留业务代码 |
