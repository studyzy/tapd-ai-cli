# 数据模型: TAPD URL 通用查询命令

**功能**: `001-tapd-url-query` | **日期**: 2026-04-16（更新：扩展至 4 种类型）

## 核心实体

### ParsedTAPDURL（内部解析结果，不持久化）

| 字段 | 类型 | 说明 |
|------|------|------|
| WorkspaceID | string | 从 URL 路径中提取，如 `51081496` |
| EntityType | string | 枚举值：`story` / `bug` / `task` / `wiki` |
| EntityID | string | 条目 ID，如 `1151081496001028684` |

**验证规则**:
- WorkspaceID 必须非空且为纯数字
- EntityType 必须为 `story`、`bug`、`task`、`wiki` 之一
- EntityID 必须非空

**状态转换**: N/A（无状态，每次请求独立解析）

---

## URL 解析规则

### 解析优先级

```
1. 检查查询参数 dialog_preview_id={type}_{id}   → 格式 2
2. 检查路径段关键字                              → 格式 1
3. 检查 fragment (#)                             → Wiki 格式
```

### Story URL

**格式 1（详情页）**:
```
https://www.tapd.cn/tapd_fe/{workspace_id}/story/detail/{story_id}
                              ^^^^^^^^^^^^                ^^^^^^^^^^
```

**格式 2（列表页预览）**:
```
https://www.tapd.cn/tapd_fe/{workspace_id}/story/list?...&dialog_preview_id=story_{story_id}
                              ^^^^^^^^^^^^                                          ^^^^^^^^^^
```

### Bug URL

**格式 1（详情页）**:
```
https://www.tapd.cn/tapd_fe/{workspace_id}/bug/detail/{bug_id}
```

**格式 2（列表页预览）**:
```
https://www.tapd.cn/tapd_fe/{workspace_id}/bug/list?...&dialog_preview_id=bug_{bug_id}
```

### Task URL

**格式 1（详情页）**:
```
https://www.tapd.cn/tapd_fe/{workspace_id}/task/detail/{task_id}
```

**格式 2（看板页预览）**:
```
https://www.tapd.cn/{workspace_id}/prong/tasks?...&dialog_preview_id=task_{task_id}
                     ^^^^^^^^^^^^（无 tapd_fe 前缀）
```

### Wiki URL

**仅一种格式（fragment 携带 ID）**:
```
https://www.tapd.cn/{workspace_id}/markdown_wikis/show/#{wiki_id}
                     ^^^^^^^^^^^^                       ^^^^^^^^^^
                     path 第 1 段                       fragment
```

---

## 输出实体

`url` 命令的输出格式与对应类型的 `show` 命令完全一致：

| 类型 | 等价命令 | 输出包装键 |
|------|---------|-----------|
| story | `tapd story show <id>` | `Story` |
| bug | `tapd bug show <id>` | `Bug` |
| task | `tapd task show <id>` | `Task` |
| wiki | 新增（无现有 show 命令） | `TapdWiki` |

### Wiki 响应关键字段（来自 TAPD API）

| 字段 | 说明 |
|------|------|
| id | Wiki ID |
| name | 文档标题 |
| workspace_id | 工作区 ID |
| description | HTML 格式内容 |
| markdown_description | Markdown 格式内容（优先使用） |
| parent_wiki_id | 父文档 ID |
| creator | 创建人 |
| modifier | 最后修改人 |
| created | 创建时间 |
| modified | 最后修改时间 |
