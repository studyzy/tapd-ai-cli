// Package model 中的 story.go 定义了 TAPD 需求数据模型
package model

// Story 表示 TAPD 需求/工作项，字段覆盖 TAPD API 返回的所有常用字段
// 使用强类型结构体反序列化可自动过滤 custom_field_* 等无用字段，节约 token
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/story/story.html
type Story struct {
	// 基本信息
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	WorkspaceID    string `json:"workspace_id,omitempty"`
	Status         string `json:"status,omitempty"`
	Type           string `json:"type,omitempty"`
	Source         string `json:"source,omitempty"`
	Flows          string `json:"flows,omitempty"`
	CreatedFrom    string `json:"created_from,omitempty"`
	WorkitemTypeID string `json:"workitem_type_id,omitempty"`
	TemplatedID    string `json:"templated_id,omitempty"`

	// 优先级
	Priority      string `json:"priority,omitempty"`
	PriorityLabel string `json:"priority_label,omitempty"`
	BusinessValue string `json:"business_value,omitempty"`

	// 人员相关
	Owner     string `json:"owner,omitempty"`
	Creator   string `json:"creator,omitempty"`
	Developer string `json:"developer,omitempty"`
	CC        string `json:"cc,omitempty"`

	// 时间相关
	Created   string `json:"created,omitempty"`
	Modified  string `json:"modified,omitempty"`
	Completed string `json:"completed,omitempty"`
	Begin     string `json:"begin,omitempty"`
	Due       string `json:"due,omitempty"`

	// 关联与分类
	IterationID string `json:"iteration_id,omitempty"`
	Module      string `json:"module,omitempty"`
	Feature     string `json:"feature,omitempty"`
	Label       string `json:"label,omitempty"`
	CategoryID  string `json:"category_id,omitempty"`
	ParentID    string `json:"parent_id,omitempty"`
	ChildrenID  string `json:"children_id,omitempty"`
	AncestorID  string `json:"ancestor_id,omitempty"`
	Path        string `json:"path,omitempty"`
	Level       string `json:"level,omitempty"`
	ReleaseID   string `json:"release_id,omitempty"`
	BugID       string `json:"bug_id,omitempty"`
	Version     string `json:"version,omitempty"`

	// 规模与工时
	Size            string `json:"size,omitempty"`
	Effort          string `json:"effort,omitempty"`
	EffortCompleted string `json:"effort_completed,omitempty"`
	Remain          string `json:"remain,omitempty"`
	Exceed          string `json:"exceed,omitempty"`

	// 进度与风险
	Progress   string `json:"progress,omitempty"`
	TechRisk   string `json:"tech_risk,omitempty"`
	TestFocus  string `json:"test_focus,omitempty"`
	IsArchived string `json:"is_archived,omitempty"`

	// 附加信息
	URL string `json:"url,omitempty"`
}

// ListStoriesRequest 查询需求列表的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/story/get_stories.html
type ListStoriesRequest struct {
	WorkspaceID   string // 必填：项目 ID
	EntityType    string // 必填：实体类型（stories 或 tasks）
	ID            string // 可选：需求 ID
	Name          string // 可选：标题（支持模糊匹配）
	PriorityLabel string // 可选：优先级
	Status        string // 可选：状态
	Owner         string // 可选：处理人
	IterationID   string // 可选：迭代 ID
	CategoryID    string // 可选：需求分类
	Label         string // 可选：标签
	Fields        string // 可选：返回字段列表
	Limit         string // 可选：返回数量限制
	Page          string // 可选：页码
	Order         string // 可选：排序规则
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *ListStoriesRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "id", r.ID)
	setOptional(params, "name", r.Name)
	setOptional(params, "priority_label", r.PriorityLabel)
	setOptional(params, "status", r.Status)
	setOptional(params, "owner", r.Owner)
	setOptional(params, "iteration_id", r.IterationID)
	setOptional(params, "category_id", r.CategoryID)
	setOptional(params, "label", r.Label)
	setOptional(params, "fields", r.Fields)
	setOptional(params, "limit", r.Limit)
	setOptional(params, "page", r.Page)
	setOptional(params, "order", r.Order)
	return params
}

// CreateStoryRequest 创建需求的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/story/add_story.html
type CreateStoryRequest struct {
	WorkspaceID   string // 必填：项目 ID
	Name          string // 必填：标题
	EntityType    string // 必填：实体类型（stories 或 tasks）
	Description   string // 可选：详细描述
	PriorityLabel string // 可选：优先级
	Owner         string // 可选：处理人
	Creator       string // 可选：创建人
	IterationID   string // 可选：迭代 ID
	StoryID       string // 可选：关联需求 ID（仅 tasks 使用）
	ParentID      string // 可选：父需求 ID（创建子需求时使用）
	Label         string // 可选：标签
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CreateStoryRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"name":         r.Name,
	}
	setOptional(params, "description", r.Description)
	setOptional(params, "priority_label", r.PriorityLabel)
	setOptional(params, "owner", r.Owner)
	setOptional(params, "creator", r.Creator)
	setOptional(params, "iteration_id", r.IterationID)
	setOptional(params, "story_id", r.StoryID)
	setOptional(params, "parent_id", r.ParentID)
	setOptional(params, "label", r.Label)
	return params
}

// UpdateStoryRequest 更新需求的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/story/update_story.html
type UpdateStoryRequest struct {
	WorkspaceID   string // 必填：项目 ID
	ID            string // 必填：需求 ID
	EntityType    string // 必填：实体类型（stories 或 tasks）
	Name          string // 可选：标题
	Status        string // 可选：状态
	VStatus       string // 可选：中文状态名
	PriorityLabel string // 可选：优先级
	Owner         string // 可选：处理人
	CurrentUser   string // 可选：变更人
	Description   string // 可选：详细描述
	Label         string // 可选：标签
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *UpdateStoryRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"id":           r.ID,
	}
	setOptional(params, "name", r.Name)
	setOptional(params, "status", r.Status)
	setOptional(params, "v_status", r.VStatus)
	setOptional(params, "priority_label", r.PriorityLabel)
	setOptional(params, "owner", r.Owner)
	setOptional(params, "current_user", r.CurrentUser)
	setOptional(params, "description", r.Description)
	setOptional(params, "label", r.Label)
	return params
}

// CountStoriesRequest 查询需求数量的请求参数
type CountStoriesRequest struct {
	WorkspaceID string // 必填：项目 ID
	EntityType  string // 必填：实体类型（stories 或 tasks）
	Status      string // 可选：状态
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CountStoriesRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "status", r.Status)
	return params
}
