// Package model 中的 task.go 定义了 TAPD 任务数据模型
package model

// Task 表示 TAPD 任务，字段覆盖 TAPD API 返回的所有常用字段
// 使用强类型结构体反序列化可自动过滤 custom_field_* 等无用字段，节约 token
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/task/get_tasks.html
type Task struct {
	// 基本信息
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty"`
	Status      string `json:"status,omitempty"`
	CreatedFrom string `json:"created_from,omitempty"`

	// 优先级
	Priority      string `json:"priority,omitempty"`
	PriorityLabel string `json:"priority_label,omitempty"`

	// 人员相关
	Owner   string `json:"owner,omitempty"`
	Creator string `json:"creator,omitempty"`
	CC      string `json:"cc,omitempty"`

	// 时间相关
	Created   string `json:"created,omitempty"`
	Modified  string `json:"modified,omitempty"`
	Completed string `json:"completed,omitempty"`
	Begin     string `json:"begin,omitempty"`
	Due       string `json:"due,omitempty"`

	// 关联与分类
	StoryID     string `json:"story_id,omitempty"`
	IterationID string `json:"iteration_id,omitempty"`
	ReleaseID   string `json:"release_id,omitempty"`
	Label       string `json:"label,omitempty"`

	// 工时与进度
	Effort          string `json:"effort,omitempty"`
	EffortCompleted string `json:"effort_completed,omitempty"`
	Remain          string `json:"remain,omitempty"`
	Exceed          string `json:"exceed,omitempty"`
	Progress        string `json:"progress,omitempty"`
	HasAttachment   string `json:"has_attachment,omitempty"`

	// 附加信息
	URL string `json:"url,omitempty"`
}

// ListTasksRequest 查询任务列表的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/task/get_tasks.html
type ListTasksRequest struct {
	WorkspaceID   string // 必填：项目 ID
	ID            string // 可选：任务 ID
	Name          string // 可选：标题（支持模糊匹配）
	Status        string // 可选：状态（open/progressing/done）
	Owner         string // 可选：处理人
	Creator       string // 可选：创建人
	StoryID       string // 可选：关联需求 ID
	IterationID   string // 可选：迭代 ID
	PriorityLabel string // 可选：优先级
	Label         string // 可选：标签
	Fields        string // 可选：返回字段列表
	Limit         string // 可选：返回数量限制（默认 30，最大 200）
	Page          string // 可选：页码
	Order         string // 可选：排序规则
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *ListTasksRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "id", r.ID)
	setOptional(params, "name", r.Name)
	setOptional(params, "status", r.Status)
	setOptional(params, "owner", r.Owner)
	setOptional(params, "creator", r.Creator)
	setOptional(params, "story_id", r.StoryID)
	setOptional(params, "iteration_id", r.IterationID)
	setOptional(params, "priority_label", r.PriorityLabel)
	setOptional(params, "label", r.Label)
	setOptional(params, "fields", r.Fields)
	setOptional(params, "limit", r.Limit)
	setOptional(params, "page", r.Page)
	setOptional(params, "order", r.Order)
	return params
}

// CreateTaskRequest 创建任务的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/task/add_task.html
type CreateTaskRequest struct {
	WorkspaceID   string // 必填：项目 ID
	Name          string // 必填：任务标题
	Description   string // 可选：详细描述
	Owner         string // 可选：处理人
	Creator       string // 可选：创建人
	CC            string // 可选：抄送人
	Begin         string // 可选：预计开始日期
	Due           string // 可选：预计结束日期
	StoryID       string // 可选：关联需求 ID
	IterationID   string // 可选：迭代 ID
	PriorityLabel string // 可选：优先级（推荐使用）
	Effort        string // 可选：预估工时
	Label         string // 可选：标签
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CreateTaskRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"name":         r.Name,
	}
	setOptional(params, "description", r.Description)
	setOptional(params, "owner", r.Owner)
	setOptional(params, "creator", r.Creator)
	setOptional(params, "cc", r.CC)
	setOptional(params, "begin", r.Begin)
	setOptional(params, "due", r.Due)
	setOptional(params, "story_id", r.StoryID)
	setOptional(params, "iteration_id", r.IterationID)
	setOptional(params, "priority_label", r.PriorityLabel)
	setOptional(params, "effort", r.Effort)
	setOptional(params, "label", r.Label)
	return params
}

// UpdateTaskRequest 更新任务的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/task/update_task.html
type UpdateTaskRequest struct {
	WorkspaceID   string // 必填：项目 ID
	ID            string // 必填：任务 ID
	Name          string // 可选：任务标题
	Description   string // 可选：详细描述
	Status        string // 可选：状态（open/progressing/done）
	Owner         string // 可选：处理人
	CurrentUser   string // 可选：操作人
	CC            string // 可选：抄送人
	Begin         string // 可选：预计开始日期
	Due           string // 可选：预计结束日期
	StoryID       string // 可选：关联需求 ID
	IterationID   string // 可选：迭代 ID
	PriorityLabel string // 可选：优先级（推荐使用）
	Effort        string // 可选：预估工时
	Label         string // 可选：标签
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *UpdateTaskRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"id":           r.ID,
	}
	setOptional(params, "name", r.Name)
	setOptional(params, "description", r.Description)
	setOptional(params, "status", r.Status)
	setOptional(params, "owner", r.Owner)
	setOptional(params, "current_user", r.CurrentUser)
	setOptional(params, "cc", r.CC)
	setOptional(params, "begin", r.Begin)
	setOptional(params, "due", r.Due)
	setOptional(params, "story_id", r.StoryID)
	setOptional(params, "iteration_id", r.IterationID)
	setOptional(params, "priority_label", r.PriorityLabel)
	setOptional(params, "effort", r.Effort)
	setOptional(params, "label", r.Label)
	return params
}

// CountTasksRequest 查询任务数量的请求参数
type CountTasksRequest struct {
	WorkspaceID string // 必填：项目 ID
	Status      string // 可选：状态
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CountTasksRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "status", r.Status)
	return params
}
