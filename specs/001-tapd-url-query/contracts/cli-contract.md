# CLI 命令合同: tapd url

**功能**: `001-tapd-url-query` | **日期**: 2026-04-16（更新：扩展至 4 种类型）

## 命令签名

```
tapd url <url> [flags]
```

## 参数

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `url` | string | ✅ | TAPD 条目的完整 URL（支持详情页、列表预览页、Wiki fragment 格式） |

## 全局标志（继承自 root）

| 标志 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--pretty` | bool | false | 输出格式化 JSON（带缩进） |
| `--workspace-id` | string | 配置文件值 | 指定工作区 ID（URL 中的值优先） |
| `--access-token` | string | 配置文件值 | TAPD Access Token |
| `--api-user` | string | 配置文件值 | TAPD API 用户名 |
| `--api-password` | string | 配置文件值 | TAPD API 密码 |

## 支持的 URL 格式

| 类型 | 格式示例 |
|------|---------|
| Story 详情页 | `https://www.tapd.cn/tapd_fe/51081496/story/detail/1151081496001028684` |
| Story 列表预览 | `https://www.tapd.cn/tapd_fe/51081496/story/list?...&dialog_preview_id=story_1151081496001028684` |
| Bug 详情页 | `https://www.tapd.cn/tapd_fe/51081496/bug/detail/1151081496001016136` |
| Bug 列表预览 | `https://www.tapd.cn/tapd_fe/51081496/bug/list?...&dialog_preview_id=bug_1151081496001016136` |
| Task 详情页 | `https://www.tapd.cn/tapd_fe/51081496/task/detail/1151081496001028786` |
| Task 看板预览 | `https://www.tapd.cn/51081496/prong/tasks?...&dialog_preview_id=task_1151081496001028786` |
| Wiki | `https://www.tapd.cn/51081496/markdown_wikis/show/#1151081496001001503` |

## 输出格式

### 成功（Story）— 与 `tapd story show` 一致

```json
{"id":"1151081496001028684","name":"用户登录功能","status":"done","owner":"张三","priority":"High"}
```

### 成功（Bug）— 与 `tapd bug show` 一致

```json
{"id":"1151081496001016136","title":"登录页面崩溃","severity":"fatal","status":"opened","owner":"李四"}
```

### 成功（Task）— 与 `tapd task show` 一致

```json
{"id":"1151081496001028786","name":"实现登录接口","status":"open","owner":"王五"}
```

### 成功（Wiki）

```json
{"id":"1151081496001001503","name":"技术架构说明","markdown_description":"# 技术架构\n...","creator":"张三","modified":"2026-04-01 10:00:00"}
```

### 错误输出（stderr）

```json
{"error":"invalid_tapd_url","message":"not a valid TAPD URL","suggestion":"provide a TAPD URL like https://www.tapd.cn/tapd_fe/{workspace_id}/story/detail/{id}"}
```

```json
{"error":"unsupported_entity_type","message":"unsupported TAPD entity type: iteration","suggestion":"supported types: story, bug, task, wiki"}
```

## 退出码

| 退出码 | 含义 |
|--------|------|
| 0 | 成功 |
| 1 | 认证失败 |
| 2 | 资源不存在 |
| 3 | 参数错误（无效 URL、不支持类型） |
| 4 | API 或网络错误 |

## 使用示例

```bash
# 查询需求（详情页 URL）
tapd url https://www.tapd.cn/tapd_fe/51081496/story/detail/1151081496001028684

# 查询需求（列表页预览 URL）
tapd url "https://www.tapd.cn/tapd_fe/51081496/story/list?dialog_preview_id=story_1151081496001028684"

# 查询缺陷（格式化输出）
tapd url https://www.tapd.cn/tapd_fe/51081496/bug/detail/1151081496001016136 --pretty

# 查询任务（看板预览 URL）
tapd url "https://www.tapd.cn/51081496/prong/tasks?dialog_preview_id=task_1151081496001028786"

# 查询 Wiki
tapd url "https://www.tapd.cn/51081496/markdown_wikis/show/#1151081496001001503"
```

## 向后兼容性

- 新增顶级子命令，不影响现有任何命令
- Story/Bug/Task 输出格式与对应 `show` 命令完全一致
- Wiki 为新增类型，无现有等价命令
