# 数据模型: tapd-ai-cli MVP

**来源**: spec.md 关键实体 + TAPD MCP Server 字段定义

## 实体

### Config（配置）

本地持久化的配置数据，存储于 `.tapd.json`。

凭据（api_user/api_password）通过 `auth login` 写入 `~/.tapd.json`（全局）或 `./.tapd.json`（`--local`）。
workspace_id 通过 `workspace switch` 始终写入当前目录的 `./.tapd.json`，使每个项目目录拥有独立的工作区配置。

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| access_token | string | 否 | TAPD Access Token（与 api_user/api_password 二选一，优先级更高） |
| api_user | string | 否 | TAPD API 用户名（与 access_token 二选一） |
| api_password | string | 否 | TAPD API 密码（与 access_token 二选一） |
| workspace_id | string | 否 | 当前工作区 ID（由 workspace switch 写入当前目录） |

**验证规则**:
- access_token 或 api_user+api_password 至少提供一组
- 如果同时存在，access_token 优先
- api_user 和 api_password 必须成对提供
- workspace_id 如果存在，必须为纯数字

---

### Workspace（工作区）

TAPD 项目，通过 `get_user_participant_projects` 和 `get_workspace_info` API 获取。

| 字段 | 类型 | 列表视图 | 详情视图 | 说明 |
|------|------|----------|----------|------|
| id | string | ✅ | ✅ | 项目 ID |
| name | string | ✅ | ✅ | 项目名称 |
| status | string | ✅ | ✅ | 项目状态 |
| category | string | — | ✅ | 类型（project/organization），列表需过滤 organization |
| creator | string | — | ✅ | 创建人 |
| created | string | — | ✅ | 创建时间 |

---

### Story（需求）

TAPD 需求/工作项，通过 stories API 操作，entity_type = "stories"。

| 字段 | 类型 | 列表视图 | 详情视图 | 创建时 | 更新时 | 说明 |
|------|------|----------|----------|--------|--------|------|
| id | string | ✅ | ✅ | — | 必需 | 需求 ID |
| name | string | ✅ | ✅ | 必需 | 可选 | 标题 |
| status | string | ✅ | ✅ | — | 可选 | 状态（使用 v_status 中文状态名） |
| owner | string | ✅ | ✅ | 可选 | 可选 | 处理人 |
| priority | string | — | ✅ | 可选 | 可选 | 优先级（High/Middle/Low/Nice To Have） |
| description | string | — | ✅ | 可选 | 可选 | 详细描述（HTML→Markdown） |
| iteration_id | string | — | ✅ | 可选 | — | 迭代 ID |
| modified | string | ✅ | ✅ | — | — | 最后修改时间 |
| url | string | — | ✅ | — | — | TAPD 链接（由客户端拼接） |

**状态转换**: 需求状态由 TAPD 工作流控制，通过 `v_status` 字段传入中文状态名。

---

### Task（任务）

TAPD 任务，与 Story 共用 API，entity_type = "tasks"。

| 字段 | 类型 | 列表视图 | 详情视图 | 创建时 | 更新时 | 说明 |
|------|------|----------|----------|--------|--------|------|
| id | string | ✅ | ✅ | — | 必需 | 任务 ID |
| name | string | ✅ | ✅ | 必需 | 可选 | 标题 |
| status | string | ✅ | ✅ | — | 可选 | 状态（仅 open/progressing/done） |
| owner | string | ✅ | ✅ | 可选 | 可选 | 处理人 |
| priority | string | — | ✅ | 可选 | 可选 | 优先级 |
| description | string | — | ✅ | 可选 | 可选 | 详细描述（HTML→Markdown） |
| story_id | string | — | ✅ | 可选 | — | 关联需求 ID |
| modified | string | ✅ | ✅ | — | — | 最后修改时间 |

---

### Bug（缺陷）

TAPD 缺陷，通过独立的 bugs API 操作。

| 字段 | 类型 | 列表视图 | 详情视图 | 创建时 | 更新时 | 说明 |
|------|------|----------|----------|--------|--------|------|
| id | string | ✅ | ✅ | — | 必需 | 缺陷 ID |
| title | string | ✅ | ✅ | 必需 | 可选 | 标题 |
| status | string | ✅ | ✅ | — | 可选 | 状态（使用 v_status） |
| priority | string | ✅ | ✅ | 可选 | 可选 | 优先级（urgent/high/medium/low/insignificant） |
| severity | string | — | ✅ | 可选 | 可选 | 严重程度（fatal/serious/normal/prompt/advice） |
| description | string | — | ✅ | 可选 | 可选 | 详细描述（HTML→Markdown） |
| current_owner | string | ✅ | ✅ | 可选 | — | 当前处理人 |
| reporter | string | — | ✅ | 可选 | — | 报告人 |
| modified | string | ✅ | ✅ | — | — | 最后修改时间 |
| url | string | — | ✅ | — | — | TAPD 链接 |

---

### Iteration（迭代）

TAPD 迭代，仅支持查询。

| 字段 | 类型 | 列表视图 | 说明 |
|------|------|----------|------|
| id | string | ✅ | 迭代 ID |
| name | string | ✅ | 迭代名称 |
| status | string | ✅ | 状态（open/done） |
| startdate | string | ✅ | 开始日期 |
| enddate | string | ✅ | 结束日期 |

## 关系

```
Workspace 1──* Story
Workspace 1──* Task
Workspace 1──* Bug
Workspace 1──* Iteration
Story 1──* Task (通过 story_id 关联)
Iteration 1──* Story (通过 iteration_id 关联)
```

## API 响应包装格式

TAPD API 统一返回以下 JSON 包装结构：

```json
{
  "status": 1,
  "data": [ ... ],
  "info": "success"
}
```

- `status`: 1 表示成功
- `data`: 实际数据（列表或单个对象）
- `info`: 状态信息

列表 API 的 data 中每个元素包装在实体类型键下：
```json
{"Story": {"id": "...", "name": "..."}}
```

## 输出格式约定

### 列表响应
```json
{
  "items": [...],
  "total": 25,
  "page": 1,
  "limit": 10,
  "has_more": true
}
```

### 创建/更新成功响应
```json
{
  "success": true,
  "id": "1234567890",
  "url": "https://www.tapd.cn/12345/prong/stories/view/1234567890"
}
```

### 错误响应（stderr）
```json
{
  "error": "authentication_failed",
  "message": "Invalid API credentials",
  "hint": "Run 'tapd auth login' or set TAPD_API_USER/TAPD_API_PASSWORD environment variables"
}
```
