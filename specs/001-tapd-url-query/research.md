# 研究报告: TAPD URL 通用查询命令

**功能**: `001-tapd-url-query` | **日期**: 2026-04-16（更新：扩展至 4 种类型）

## 研究任务 1: TAPD URL 格式规律（全量）

### Decision: 双格式解析策略 — 路径格式 + 查询参数格式

**已确认的 URL 格式**:

| 类型 | 格式 1（详情页） | 格式 2（列表/看板页） |
|------|-----------------|----------------------|
| Story | `/tapd_fe/{ws}/story/detail/{id}` | `?dialog_preview_id=story_{id}` |
| Bug   | `/tapd_fe/{ws}/bug/detail/{id}`   | `?dialog_preview_id=bug_{id}` |
| Task  | `/tapd_fe/{ws}/task/detail/{id}`  | `?dialog_preview_id=task_{id}` |
| Wiki  | `/{ws}/markdown_wikis/show/#{id}` | N/A（仅 fragment 格式） |

**解析优先级**:
1. 检查 `dialog_preview_id` 查询参数（格式 2）
2. 检查 path 关键字（`story/detail`、`bug/detail`、`task/detail`、`markdown_wikis`）
3. 从路径提取工作区 ID：`tapd_fe` 路径取第 2 段，直接路径取第 1 段

**Rationale**: 两种格式均为用户实际使用中的真实 URL，必须全部支持。

**Alternatives considered**:
- 仅支持详情页格式：会丢失用户从列表页分享的 URL，体验差。

---

## 研究任务 2: Wiki API 规范

### Decision: 使用 `GET /tapd_wikis` 接口，以 `id` 参数查询

**API 端点**: `https://api.tapd.cn/tapd_wikis`

**关键参数**:
- `workspace_id`（必需）：工作区 ID
- `id`（可选）：Wiki 条目 ID
- `fields`（可选）：指定返回字段

**响应关键字段**: `id`, `name`, `workspace_id`, `description`, `markdown_description`, `parent_wiki_id`, `creator`, `modifier`, `created`, `modified`

**认证方式**: Basic Auth（与现有 client 一致）

**Rationale**: API 与现有 Story/Bug API 结构完全一致（统一的 TAPD API 风格），响应格式相同，可复用 `doGet` + JSON 解析模式。

**Alternatives considered**:
- 使用 `name` 参数查询：Wiki ID 更精确，不会因同名 Wiki 产生歧义。

---

## 研究任务 3: Task 查询 API 规范

### Decision: 复用 `GetStory(workspaceID, id, "tasks")` 方法

**API 端点（已验证）**: `https://api.tapd.cn/tasks`

**关键参数**:
- `workspace_id`（必需）：工作区 ID
- `id`（可选）：Task ID，支持多值

**响应包装键**: `Task`（与现有 `GetStory` 中 `entityType="tasks"` 的实现完全吻合）

**响应关键字段**: `id`, `name`, `workspace_id`, `status`（open/progressing/done）, `owner`, `priority_label`

**Rationale**: 现有 `GetStory` 通过 `entity_type` 参数区分 story/task（传 `"tasks"` 时调用 `/tasks` 端点，包装键为 `"Task"`）。官方 API 文档已确认端点和包装键，与现有实现完全一致。Task URL 解析出 ID 后直接调用此方法，零新增 client 代码。

---

## 研究任务 4: Wiki client 方法新增

### Decision: 新增 `GetWiki(workspaceID, id string)` 方法到 `internal/client/wiki.go`

**Rationale**: Wiki 使用独立的 `/tapd_wikis` 端点，与 Story/Bug/Task 不同，需要新文件。响应结构与 Bug 类似（单层包装键 `TapdWiki`）。

**新增文件**: `internal/client/wiki.go`

---

## 解决的 NEEDS CLARIFICATION 事项

无（规范阶段已无未澄清项）。

## 结论

实现方案清晰，风险低：

| 新增内容 | 说明 |
|---------|------|
| `internal/client/wiki.go` | 新增 `GetWiki` 方法，调用 `/tapd_wikis` |
| `internal/cmd/url.go` | URL 解析 + 四类型分发逻辑 |
| `internal/cmd/url_test.go` | 覆盖所有 URL 格式的解析单元测试 |
