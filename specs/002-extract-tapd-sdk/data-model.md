# 数据模型: tapd-sdk

**分支**: `002-extract-tapd-sdk` | **日期**: 2026-04-17

## 核心实体

### Client（SDK 客户端）

SDK 的入口点，持有认证信息和 HTTP 连接配置。所有 API 操作通过 Client 方法发起。

| 字段 | 类型 | 说明 |
|------|------|------|
| baseURL | string | TAPD API 基础地址，默认 `https://api.tapd.cn` |
| authHeader | string | 认证头（Bearer Token 或 Basic 编码） |
| Nick | string | 当前用户昵称（Bearer Token 认证时自动获取） |

**创建方式**:
- `NewClient(accessToken, apiUser, apiPassword)` — 使用默认地址
- `NewClientWithBaseURL(baseURL, ...)` — 指定地址（测试用）

---

### TAPDError（结构化错误）

SDK 返回的所有错误类型，实现 `error` 接口。

| 字段 | 类型 | 说明 |
|------|------|------|
| HTTPStatus | int | HTTP 状态码（如 401、404、500） |
| ExitCode | int | CLI 退出码映射（1=未授权，2=未找到，3=参数错误，4=服务端错误） |
| Message | string | 人类可读错误描述（英文） |

---

### TAPDResponse（API 响应包装）

TAPD API 统一响应格式，SDK 内部使用，不暴露给调用方。

| 字段 | 类型 | 说明 |
|------|------|------|
| Status | int | 1 表示成功，其他表示错误 |
| Data | json.RawMessage | 实际数据负载 |
| Info | string | 错误时的描述信息 |

---

## 业务实体

### Story（需求）

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | string | 需求唯一标识 |
| Name | string | 需求标题 |
| Description | string | 详细描述（原始 HTML 字符串，CLI 层负责转 Markdown） |
| WorkspaceID | string | 所属项目 ID |
| Status | string | 状态（如 planning/developing/done） |
| Priority / PriorityLabel | string | 优先级 |
| Owner / Creator | string | 处理人 / 创建人 |
| IterationID | string | 所属迭代 ID |
| CategoryID | string | 需求分类 ID |
| Created / Modified | string | 创建/修改时间（ISO 格式） |
| URL | string | 需求详情页链接（SDK 自动拼接） |

**关联操作**: ListStories、GetStory、CreateStory、UpdateStory、CountStories

---

### Bug（缺陷）

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | string | 缺陷唯一标识 |
| Title | string | 缺陷标题 |
| Description | string | 详细描述（原始 HTML） |
| WorkspaceID | string | 所属项目 ID |
| Status | string | 状态 |
| Priority | string | 优先级 |
| Severity | string | 严重程度 |
| Reporter / CurrentOwner | string | 报告人 / 当前处理人 |
| IterationID | string | 所属迭代 ID |
| Created / Modified | string | 时间信息 |

**关联操作**: ListBugs、GetBug、CreateBug、UpdateBug、CountBugs

---

### Task（任务）

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | string | 任务唯一标识 |
| Name | string | 任务名称 |
| Description | string | 详细描述（原始 HTML） |
| WorkspaceID | string | 所属项目 ID |
| Status | string | 状态 |
| Priority | string | 优先级 |
| Owner / Creator | string | 处理人 / 创建人 |
| StoryID | string | 关联需求 ID |
| IterationID | string | 所属迭代 ID |
| Effort / EffortCompleted | string | 预估/已完成工时 |

**关联操作**: ListTasks、GetTask、CreateTask、UpdateTask、CountTasks

---

### Iteration（迭代）

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | string | 迭代唯一标识 |
| Name | string | 迭代名称 |
| WorkspaceID | string | 所属项目 ID |
| Status | string | 状态（planning/doing/done） |
| StartDate / EndDate | string | 开始/结束日期 |
| Goal | string | 迭代目标 |

**关联操作**: ListIterations、CreateIteration、UpdateIteration、CountIterations

---

### Comment（评论）

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | string | 评论唯一标识 |
| WorkspaceID | string | 所属项目 ID |
| EntryType | string | 关联实体类型（story/bug/task） |
| EntryID | string | 关联实体 ID |
| Description | string | 评论内容（原始 HTML） |
| Author | string | 评论人 |
| Created / Modified | string | 时间信息 |

**关联操作**: ListComments、AddComment、UpdateComment、CountComments

---

### Wiki（文档）

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | string | 文档唯一标识 |
| Title | string | 文档标题 |
| Detail | string | 文档内容（原始 HTML） |
| WorkspaceID | string | 所属项目 ID |
| Creator | string | 创建人 |
| ParentID | string | 父文档 ID（层级结构） |
| Created / Modified | string | 时间信息 |

**关联操作**: ListWikis、GetWiki、CreateWiki、UpdateWiki

---

### Attachment（附件）

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | string | 附件唯一标识 |
| WorkspaceID | string | 所属项目 ID |
| Filename | string | 文件名 |
| FileSize | string | 文件大小 |
| DownloadURL | string | 下载链接 |
| Uploader | string | 上传人 |

**关联操作**: GetAttachments、GetImage

---

### Timesheet（工时）

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | string | 工时记录唯一标识 |
| WorkspaceID | string | 所属项目 ID |
| EntityType | string | 关联实体类型 |
| EntityID | string | 关联实体 ID |
| Spentdate | string | 工时日期 |
| Timespent | string | 花费时间（小时） |
| Owner | string | 记录人 |

**关联操作**: ListTimesheets、AddTimesheet、UpdateTimesheet

---

### 通用响应实体

| 实体 | 用途 |
|------|------|
| SuccessResponse | 创建/更新操作结果（含 ID、URL） |
| CountResponse | 计数查询结果（含 Count） |
| Workspace | 项目/工作区信息 |
| Category | 需求/缺陷分类 |
| Release | 发布计划 |
| WorkflowStatusMap | 工作流状态映射 |
| WorkflowTransition | 工作流转换规则 |

---

## 实体关系

```
Workspace (1) ──── (*) Story
Workspace (1) ──── (*) Bug
Workspace (1) ──── (*) Task
Workspace (1) ──── (*) Iteration
Workspace (1) ──── (*) Wiki
Workspace (1) ──── (*) Release
Workspace (1) ──── (*) Category

Iteration (1) ──── (*) Story
Iteration (1) ──── (*) Bug
Iteration (1) ──── (*) Task

Story (1) ──── (*) Task
Story (1) ──── (*) Comment
Story (1) ──── (*) Attachment
Bug (1) ──── (*) Comment
Bug (1) ──── (*) Attachment
Task (1) ──── (*) Comment
Task (1) ──── (*) Timesheet

Story (*) ──── (*) Bug  [通过 Relation]
Story (*) ──── (*) Story [父子层级，ParentID]
```

---

## 通用请求参数模式

所有请求结构体实现 `ToParams() map[string]string` 方法，统一转换为 TAPD API 参数格式。

| 通用请求类型 | 适用场景 |
|-------------|---------|
| WorkspaceIDRequest | 仅需 workspace_id 的查询 |
| WorkflowRequest | 工作流相关查询 |
| GetTodoRequest | 用户待办查询（支持分页） |
| GetCommitMsgRequest | 源码提交关键字查询 |
| CreateRelationRequest | 创建实体关联关系 |
