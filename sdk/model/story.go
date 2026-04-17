// Package model 中的 story.go 定义了 TAPD 需求数据模型
package model

import "encoding/json"

// Story 表示 TAPD 需求/工作项，字段覆盖 TAPD API 返回的所有常用字段
// 自定义字段（custom_field_*、custom_plan_field_*）通过 CustomFields map 保留，不会丢失
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/story/story.html
type Story struct {
	// 基本信息
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	WorkspaceID    string `json:"workspace_id,omitempty"`
	Status         string `json:"status,omitempty"`
	Step           string `json:"step,omitempty"`
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
	Progress       string `json:"progress,omitempty"`
	ProgressManual string `json:"progress_manual,omitempty"`
	TechRisk       string `json:"tech_risk,omitempty"`
	TestFocus      string `json:"test_focus,omitempty"`
	IsArchived     string `json:"is_archived,omitempty"`

	// 保密
	SecretRootID string `json:"secret_root_id,omitempty"`

	// 附加信息
	URL string `json:"url,omitempty"`

	// 自定义字段，key 为 custom_field_one、custom_field_9 等
	CustomFields map[string]string `json:"-"`
}

// UnmarshalJSON 自定义反序列化，在解析标准字段的同时收集 custom_field_* 和 custom_plan_field_* 字段
func (s *Story) UnmarshalJSON(data []byte) error {
	type Alias Story
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*s = Story(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	s.CustomFields = ExtractCustomFields(raw)
	return nil
}

// MarshalJSON 自定义序列化，将 CustomFields 中的键值对合并到输出 JSON
func (s Story) MarshalJSON() ([]byte, error) {
	type Alias Story
	b, err := json.Marshal(Alias(s))
	if err != nil {
		return nil, err
	}
	if len(s.CustomFields) == 0 {
		return b, nil
	}

	// 将自定义字段合并到 JSON 对象中
	var base map[string]json.RawMessage
	if err := json.Unmarshal(b, &base); err != nil {
		return nil, err
	}
	for k, v := range s.CustomFields {
		raw, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		base[k] = raw
	}
	return json.Marshal(base)
}

// ListStoriesRequest 查询需求列表的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/story/get_stories.html
type ListStoriesRequest struct {
	WorkspaceID   string // 必填：项目 ID
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
	WorkspaceID    string            // 必填：项目 ID
	Name           string            // 必填：标题
	Description    string            // 可选：详细描述
	PriorityLabel  string            // 可选：优先级
	Owner          string            // 可选：处理人
	Creator        string            // 可选：创建人
	Developer      string            // 可选：开发人员
	CC             string            // 可选：抄送人
	IterationID    string            // 可选：迭代 ID
	ParentID       string            // 可选：父需求 ID（创建子需求时使用）
	CategoryID     string            // 可选：需求分类 ID
	Type           string            // 可选：需求类型
	Source         string            // 可选：需求来源
	Begin          string            // 可选：预计开始日期
	Due            string            // 可选：预计结束日期
	Label          string            // 可选：标签
	TemplatedID    string            // 可选：模板 ID
	WorkitemTypeID string            // 可选：需求类别 ID
	CustomFields   map[string]string // 可选：自定义字段，key 如 custom_field_one、custom_field_9
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
	setOptional(params, "developer", r.Developer)
	setOptional(params, "cc", r.CC)
	setOptional(params, "iteration_id", r.IterationID)
	setOptional(params, "parent_id", r.ParentID)
	setOptional(params, "category_id", r.CategoryID)
	setOptional(params, "type", r.Type)
	setOptional(params, "source", r.Source)
	setOptional(params, "begin", r.Begin)
	setOptional(params, "due", r.Due)
	setOptional(params, "label", r.Label)
	setOptional(params, "templated_id", r.TemplatedID)
	setOptional(params, "workitem_type_id", r.WorkitemTypeID)
	MergeCustomFields(params, r.CustomFields)
	return params
}

// UpdateStoryRequest 更新需求的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/story/update_story.html
type UpdateStoryRequest struct {
	WorkspaceID   string            // 必填：项目 ID
	ID            string            // 必填：需求 ID
	Name          string            // 可选：标题
	Status        string            // 可选：状态
	VStatus       string            // 可选：中文状态名
	PriorityLabel string            // 可选：优先级
	Owner         string            // 可选：处理人
	CurrentUser   string            // 可选：变更人
	Developer     string            // 可选：开发人员
	CC            string            // 可选：抄送人
	Description   string            // 可选：详细描述
	IterationID   string            // 可选：迭代 ID
	CategoryID    string            // 可选：需求分类 ID
	Begin         string            // 可选：预计开始日期
	Due           string            // 可选：预计结束日期
	Label         string            // 可选：标签
	Type          string            // 可选：需求类型
	Source        string            // 可选：需求来源
	CustomFields  map[string]string // 可选：自定义字段，key 如 custom_field_one、custom_field_9
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
	setOptional(params, "developer", r.Developer)
	setOptional(params, "cc", r.CC)
	setOptional(params, "description", r.Description)
	setOptional(params, "iteration_id", r.IterationID)
	setOptional(params, "category_id", r.CategoryID)
	setOptional(params, "begin", r.Begin)
	setOptional(params, "due", r.Due)
	setOptional(params, "label", r.Label)
	setOptional(params, "type", r.Type)
	setOptional(params, "source", r.Source)
	MergeCustomFields(params, r.CustomFields)
	return params
}

// CountStoriesRequest 查询需求数量的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/story/get_stories_count.html
type CountStoriesRequest struct {
	WorkspaceID string // 必填：项目 ID
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
