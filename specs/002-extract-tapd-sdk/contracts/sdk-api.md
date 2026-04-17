# SDK 公开 API 合约: tapd-sdk

**分支**: `002-extract-tapd-sdk` | **日期**: 2026-04-17
**模块**: `github.com/studyzy/tapd-sdk-go`

## 客户端初始化

```go
// NewClient 使用默认 TAPD API 地址创建客户端
// accessToken 非空时使用 Bearer 认证，否则使用 Basic 认证
func NewClient(accessToken, apiUser, apiPassword string) *Client

// NewClientWithBaseURL 使用指定地址创建客户端（测试场景使用）
func NewClientWithBaseURL(baseURL, accessToken, apiUser, apiPassword string) *Client

// FetchNick 通过 API 获取当前用户昵称（Bearer Token 认证时调用）
func (c *Client) FetchNick()

// TestAuth 验证认证凭据是否有效，返回 nil 表示成功
func (c *Client) TestAuth() error
```

---

## 需求（Story）

```go
func (c *Client) ListStories(req *model.ListStoriesRequest) ([]model.Story, error)
func (c *Client) GetStory(workspaceID, id string) (*model.Story, error)
func (c *Client) CreateStory(req *model.CreateStoryRequest) (*model.SuccessResponse, error)
func (c *Client) UpdateStory(req *model.UpdateStoryRequest) (*model.Story, error)
func (c *Client) CountStories(req *model.CountStoriesRequest) (int, error)
func (c *Client) GetStoryFieldsLabel(req *model.WorkspaceIDRequest) (json.RawMessage, error)
func (c *Client) GetStoryFieldsInfo(req *model.WorkspaceIDRequest) (json.RawMessage, error)
func (c *Client) GetTodoStories(req *model.GetTodoRequest) ([]model.Story, error)
```

---

## 缺陷（Bug）

```go
func (c *Client) ListBugs(req *model.ListBugsRequest) ([]model.Bug, error)
func (c *Client) GetBug(workspaceID, id string) (*model.Bug, error)
func (c *Client) CreateBug(req *model.CreateBugRequest) (*model.SuccessResponse, error)
func (c *Client) UpdateBug(req *model.UpdateBugRequest) (*model.Bug, error)
func (c *Client) CountBugs(req *model.CountBugsRequest) (int, error)
func (c *Client) GetTodoBugs(req *model.GetTodoRequest) ([]model.Bug, error)
func (c *Client) GetRelatedBugs(req *model.GetRelatedBugsRequest) (json.RawMessage, error)
```

---

## 任务（Task）

```go
func (c *Client) ListTasks(req *model.ListTasksRequest) ([]model.Task, error)
func (c *Client) GetTask(workspaceID, id string) (*model.Task, error)
func (c *Client) CreateTask(req *model.CreateTaskRequest) (*model.SuccessResponse, error)
func (c *Client) UpdateTask(req *model.UpdateTaskRequest) (*model.Task, error)
func (c *Client) CountTasks(req *model.CountTasksRequest) (int, error)
func (c *Client) GetTodoTasks(req *model.GetTodoRequest) ([]model.Task, error)
```

---

## 迭代（Iteration）

```go
func (c *Client) ListIterations(req *model.ListIterationsRequest) ([]model.Iteration, error)
func (c *Client) CreateIteration(req *model.CreateIterationRequest) (*model.SuccessResponse, error)
func (c *Client) UpdateIteration(req *model.UpdateIterationRequest) (*model.Iteration, error)
func (c *Client) CountIterations(req *model.CountIterationsRequest) (int, error)
```

---

## 评论（Comment）

```go
func (c *Client) ListComments(req *model.ListCommentsRequest) ([]model.Comment, error)
func (c *Client) AddComment(req *model.AddCommentRequest) (*model.Comment, error)
func (c *Client) UpdateComment(req *model.UpdateCommentRequest) (*model.Comment, error)
func (c *Client) CountComments(req *model.CountCommentsRequest) (int, error)
```

---

## 文档（Wiki）

```go
func (c *Client) ListWikis(req *model.ListWikisRequest) ([]model.Wiki, error)
func (c *Client) GetWiki(workspaceID, id string) (*model.Wiki, error)
func (c *Client) CreateWiki(req *model.CreateWikiRequest) (*model.SuccessResponse, error)
func (c *Client) UpdateWiki(req *model.UpdateWikiRequest) (*model.Wiki, error)
```

---

## 附件（Attachment）

```go
func (c *Client) GetAttachments(req *model.GetAttachmentsRequest) ([]model.Attachment, error)
func (c *Client) GetImage(req *model.GetImageRequest) (*model.ImageInfo, error)
```

---

## 工时（Timesheet）

```go
func (c *Client) ListTimesheets(req *model.ListTimesheetsRequest) ([]model.Timesheet, error)
func (c *Client) AddTimesheet(req *model.AddTimesheetRequest) (*model.Timesheet, error)
func (c *Client) UpdateTimesheet(req *model.UpdateTimesheetRequest) (*model.Timesheet, error)
```

---

## 测试用例（TCase）

```go
func (c *Client) ListTCases(req *model.ListTCasesRequest) ([]model.TCase, error)
func (c *Client) CreateTCase(req *model.CreateTCaseRequest) (*model.TCase, error)
func (c *Client) BatchCreateTCases(req *model.BatchCreateTCasesRequest) (json.RawMessage, error)
func (c *Client) CountTCases(req *model.CountTCasesRequest) (int, error)
```

---

## 关联关系（Relation）

```go
func (c *Client) CreateRelation(req *model.CreateRelationRequest) (json.RawMessage, error)
```

---

## 工作流（Workflow）

```go
func (c *Client) GetWorkflowStatusMap(req *model.WorkflowRequest) (json.RawMessage, error)
func (c *Client) GetWorkflowTransitions(req *model.WorkflowRequest) (json.RawMessage, error)
func (c *Client) GetWorkflowLastSteps(req *model.WorkflowRequest) (json.RawMessage, error)
func (c *Client) GetWorkitemTypes(req *model.WorkspaceIDRequest) (json.RawMessage, error)
```

---

## 工作区（Workspace）

```go
func (c *Client) ListWorkspaces() ([]model.Workspace, error)
func (c *Client) GetWorkspaceInfo(workspaceID string) (*model.Workspace, error)
```

---

## 分类（Category）

```go
func (c *Client) ListCategories(params map[string]string) ([]model.Category, error)
```

---

## 发布计划与杂项

```go
func (c *Client) ListReleases(req *model.WorkspaceIDRequest) ([]model.Release, error)
func (c *Client) GetCommitMsg(req *model.GetCommitMsgRequest) (json.RawMessage, error)
func (c *Client) GetCustomFields(req *model.GetCustomFieldsRequest) (json.RawMessage, error)
func (c *Client) SendQiweiMessage(webhookURL, msg string) error
```

---

## 错误类型

```go
// TAPDError 是 SDK 返回的结构化错误类型
type TAPDError struct {
    HTTPStatus int    // HTTP 状态码
    ExitCode   int    // 退出码（CLI 使用：1=未授权 2=未找到 3=参数错误 4=服务端错误）
    Message    string // 错误描述（英文）
}

func (e *TAPDError) Error() string
```

---

## API 设计约定

| 约定 | 说明 |
|------|------|
| 所有方法返回 `(T, error)` | SDK 绝不调用 `os.Exit`，错误通过返回值传递 |
| 请求参数使用结构体 | 通过 `ToParams()` 转为 API 参数 map |
| 列表方法返回 `[]T` | 空结果返回 `nil, nil`（非错误） |
| 单项方法返回 `*T` | 未找到时返回 `TAPDError{ExitCode: 2}` |
| 原始响应返回 `json.RawMessage` | 用于结构未完全定义的接口 |
| Description 字段保留原始 HTML | 调用方（CLI）负责格式转换 |
