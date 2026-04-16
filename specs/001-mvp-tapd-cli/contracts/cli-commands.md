# CLI 命令合同: tapd-ai-cli MVP

**格式**: `命令 [子命令] [参数] [标志]`
**全局标志**: 所有命令均支持以下标志：
- `--workspace-id <id>` — 覆盖本地配置的 workspace_id
- `--compact` — 输出紧凑 JSON（无缩进无多余空白）

**退出码约定**:
- `0` — 成功
- `1` — 认证错误
- `2` — 资源未找到
- `3` — 参数错误
- `4` — API 错误（网络/服务端）

---

## tapd auth login

持久化 API 凭据到本地配置文件。支持两种认证方式：Access Token（推荐）和 API User/Password。

```
tapd auth login --access-token <token> [--local]
tapd auth login --api-user <user> --api-password <password> [--local]
```

| 标志 | 必需 | 说明 |
|------|------|------|
| `--access-token` | 与 api-user/api-password 二选一 | TAPD Access Token（推荐） |
| `--api-user` | 与 access-token 二选一 | TAPD API 用户名 |
| `--api-password` | 与 api-user 配对 | TAPD API 密码 |
| `--local` | 否 | 写入当前目录 `.tapd.json`（默认写入 `~/.tapd.json`） |

**stdout（成功）**:
```json
{"success":true}
```

**stderr（参数缺失）**:
```json
{"error":"missing_parameter","message":"Provide --access-token or --api-user with --api-password","hint":"Usage: tapd auth login --access-token <token> OR --api-user <user> --api-password <password>"}
```

---

## tapd workspace list

列出当前用户参与的项目。

```
tapd workspace list
```

无额外标志。自动过滤 category=organization 的条目。

**stdout**:
```json
{"items":[{"id":"12345","name":"My Project","status":"active"}],"total":1}
```

---

## tapd workspace switch

切换当前工作区。workspace_id 始终写入当前目录的 `.tapd.json`（不存在则自动创建），使每个项目目录拥有独立的 workspace 配置。

```
tapd workspace switch <workspace_id>
```

| 参数 | 必需 | 说明 |
|------|------|------|
| `workspace_id` | 是 | 目标工作区 ID |

**stdout（成功）**:
```json
{"success":true,"workspace_id":"12345"}
```

---

## tapd workspace info

查看当前工作区详情。

```
tapd workspace info
```

**stdout**:
```json
{"id":"12345","name":"My Project","status":"active","creator":"user1","created":"2026-01-01"}
```

---

## tapd story list

查询需求列表（精简视图）。

```
tapd story list [--status <status>] [--owner <owner>] [--iteration-id <id>] [--limit <N>] [--page <N>]
```

| 标志 | 必需 | 默认值 | 说明 |
|------|------|--------|------|
| `--status` | 否 | — | 按状态筛选（中文状态名，如"已实现"） |
| `--owner` | 否 | — | 按处理人筛选 |
| `--iteration-id` | 否 | — | 按迭代 ID 筛选 |
| `--limit` | 否 | 10 | 返回数量限制 |
| `--page` | 否 | 1 | 页码 |

**stdout**:
```json
{"items":[{"id":"10001","name":"需求标题","status":"已实现","owner":"user1","modified":"2026-04-16"}],"total":25,"page":1,"limit":10,"has_more":true}
```

---

## tapd story show

查看需求详情（description 自动 HTML→Markdown）。

```
tapd story show <story_id>
```

| 参数 | 必需 | 说明 |
|------|------|------|
| `story_id` | 是 | 需求 ID |

**stdout**: 需求完整字段 JSON，description 为 Markdown 格式。

---

## tapd story create

创建需求。

```
tapd story create --name <title> [--description <desc>] [--owner <owner>] [--priority <priority>] [--iteration-id <id>]
```

| 标志 | 必需 | 说明 |
|------|------|------|
| `--name` | 是 | 需求标题 |
| `--description` | 否 | 描述内容 |
| `--owner` | 否 | 处理人 |
| `--priority` | 否 | 优先级（High/Middle/Low/Nice To Have） |
| `--iteration-id` | 否 | 关联迭代 ID |

**stdout（成功）**:
```json
{"success":true,"id":"10002","url":"https://www.tapd.cn/12345/prong/stories/view/10002"}
```

---

## tapd story update

更新需求。

```
tapd story update <story_id> [--name <title>] [--status <status>] [--owner <owner>] [--priority <priority>]
```

| 参数/标志 | 必需 | 说明 |
|-----------|------|------|
| `story_id` | 是 | 需求 ID |
| `--name` | 否 | 新标题 |
| `--status` | 否 | 新状态（中文状态名） |
| `--owner` | 否 | 新处理人 |
| `--priority` | 否 | 新优先级 |

**stdout**: 更新后的需求完整字段 JSON。

---

## tapd story count

查询需求数量。

```
tapd story count [--status <status>]
```

**stdout**:
```json
{"count":25}
```

---

## tapd task list / show / create / update / count

与 story 命令格式完全一致，内部 entity_type 自动设为 "tasks"。

区别：
- task 状态仅有 `open`、`progressing`、`done`
- create 支持额外标志 `--story-id <id>`（关联需求）

---

## tapd bug list

查询缺陷列表。

```
tapd bug list [--status <status>] [--priority <priority>] [--severity <severity>] [--limit <N>] [--page <N>]
```

| 标志 | 必需 | 默认值 | 说明 |
|------|------|--------|------|
| `--status` | 否 | — | 按状态筛选 |
| `--priority` | 否 | — | 按优先级（urgent/high/medium/low/insignificant） |
| `--severity` | 否 | — | 按严重程度（fatal/serious/normal/prompt/advice） |
| `--limit` | 否 | 10 | 返回数量 |
| `--page` | 否 | 1 | 页码 |

**stdout**: 同 story list 结构。

---

## tapd bug show

```
tapd bug show <bug_id>
```

**stdout**: 缺陷完整字段 JSON，description 为 Markdown 格式。

---

## tapd bug create

```
tapd bug create --title <title> [--description <desc>] [--priority <priority>] [--severity <severity>]
```

| 标志 | 必需 | 说明 |
|------|------|------|
| `--title` | 是 | 缺陷标题 |
| `--description` | 否 | 描述 |
| `--priority` | 否 | 优先级 |
| `--severity` | 否 | 严重程度 |

**stdout（成功）**:
```json
{"success":true,"id":"20001","url":"https://www.tapd.cn/12345/bugtrace/bugs/view/20001"}
```

---

## tapd bug update

```
tapd bug update <bug_id> [--title <title>] [--status <status>] [--priority <priority>] [--severity <severity>]
```

**stdout**: 更新后的缺陷完整字段 JSON。

---

## tapd bug count

```
tapd bug count [--status <status>]
```

**stdout**:
```json
{"count":12}
```

---

## tapd iteration list

查询迭代列表。

```
tapd iteration list [--status <status>]
```

| 标志 | 必需 | 说明 |
|------|------|------|
| `--status` | 否 | 按状态筛选（open/done） |

**stdout**:
```json
{"items":[{"id":"30001","name":"Sprint 1","status":"open","startdate":"2026-04-01","enddate":"2026-04-15"}],"total":3}
```

---

## tapd spec

输出 Tool Definition JSON（AI 自发现）。

```
tapd spec
```

无参数无标志。遍历 Cobra 命令树自动生成。

**stdout**: OpenAI/Anthropic 兼容的 Tool Definition JSON 数组。
