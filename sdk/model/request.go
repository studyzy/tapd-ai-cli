// Package model 中的 request.go 定义了杂项 API 的请求参数结构体
package model

// GetCustomFieldsRequest 获取自定义字段配置的请求参数
type GetCustomFieldsRequest struct {
	WorkspaceID string // 必填：项目 ID
	EntityType  string // 必填：实体类型（stories/tasks/iterations/tcases）
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *GetCustomFieldsRequest) ToParams() map[string]string {
	return map[string]string{
		"workspace_id": r.WorkspaceID,
	}
}

// WorkspaceIDRequest 仅需项目 ID 的通用请求参数
// 适用于 GetStoryFieldsLabel、GetStoryFieldsInfo、GetWorkitemTypes、
// ListCategories、ListReleases 等仅需 workspace_id 的接口
type WorkspaceIDRequest struct {
	WorkspaceID string // 必填：项目 ID
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *WorkspaceIDRequest) ToParams() map[string]string {
	return map[string]string{
		"workspace_id": r.WorkspaceID,
	}
}

// WorkflowRequest 工作流相关接口的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/workflow/
type WorkflowRequest struct {
	WorkspaceID    string // 必填：项目 ID
	System         string // 必填：系统名（story/bug）
	WorkitemTypeID string // 可选：需求类别 ID
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *WorkflowRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "system", r.System)
	setOptional(params, "workitem_type_id", r.WorkitemTypeID)
	return params
}

// GetCommitMsgRequest 获取源码提交关键字的请求参数
type GetCommitMsgRequest struct {
	WorkspaceID string // 必填：项目 ID
	ObjectID    string // 必填：条目 ID
	Type        string // 必填：条目类型（story/task/bug）
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *GetCommitMsgRequest) ToParams() map[string]string {
	return map[string]string{
		"workspace_id": r.WorkspaceID,
		"object_id":    r.ObjectID,
		"type":         r.Type,
	}
}

// GetTodoRequest 获取用户待办事项的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/user/get_user_todo_story.html
type GetTodoRequest struct {
	WorkspaceID string // 必填：项目 ID
	EntityType  string // 必填：对象类型（story/bug/task）
	Limit       string // 可选：返回数量限制（默认 30，最大 200）
	Page        string // 可选：页码（默认 1）
	Order       string // 可选：排序规则
	Fields      string // 可选：返回字段列表（逗号分隔）
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *GetTodoRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "limit", r.Limit)
	setOptional(params, "page", r.Page)
	setOptional(params, "order", r.Order)
	setOptional(params, "fields", r.Fields)
	return params
}

// GetRelatedBugsRequest 获取需求关联缺陷的请求参数
type GetRelatedBugsRequest struct {
	WorkspaceID string // 必填：项目 ID
	StoryID     string // 必填：需求 ID
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *GetRelatedBugsRequest) ToParams() map[string]string {
	return map[string]string{
		"workspace_id": r.WorkspaceID,
		"story_id":     r.StoryID,
	}
}

// CreateRelationRequest 创建实体关联关系的请求参数
type CreateRelationRequest struct {
	WorkspaceID string // 必填：项目 ID
	SourceType  string // 必填：源实体类型（story/bug/task）
	TargetType  string // 必填：目标实体类型（story/bug/task）
	SourceID    string // 必填：源实体 ID
	TargetID    string // 必填：目标实体 ID
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CreateRelationRequest) ToParams() map[string]string {
	return map[string]string{
		"workspace_id": r.WorkspaceID,
		"source_type":  r.SourceType,
		"target_type":  r.TargetType,
		"source_id":    r.SourceID,
		"target_id":    r.TargetID,
	}
}

// CustomFieldConfig 表示 TAPD 自定义字段配置
// API 文档：https://open.tapd.cn/document/api-doc/API文档/api_reference/story/get_story_custom_fields_settings.html
type CustomFieldConfig struct {
	ID              string `json:"id,omitempty"`
	WorkspaceID     string `json:"workspace_id,omitempty"`
	AppID           string `json:"app_id,omitempty"`
	EntryType       string `json:"entry_type,omitempty"`        // 所属实体类型（story/bug/task 等）
	CustomField     string `json:"custom_field,omitempty"`      // 自定义字段标识（如 custom_field_17）
	Type            string `json:"type,omitempty"`              // 输入类型（text/select/cascade_radio 等）
	Name            string `json:"name,omitempty"`              // 显示名称
	Options         string `json:"options,omitempty"`           // 可选值（JSON 字符串，需二次解析）
	ExtraConfig     string `json:"extra_config,omitempty"`      // 额外配置
	Enabled         string `json:"enabled,omitempty"`           // 是否启用（1=启用，0=禁用）
	Creator         string `json:"creator,omitempty"`           // 创建人
	Created         string `json:"created,omitempty"`           // 创建时间
	Modified        string `json:"modified,omitempty"`          // 最后修改时间
	Freeze          string `json:"freeze,omitempty"`            // 是否冻结
	Sort            string `json:"sort,omitempty"`              // 显示排序
	Memo            string `json:"memo,omitempty"`              // 备注
	OpenExtensionID string `json:"open_extension_id,omitempty"` // 插件扩展字段标识
	IsOut           int    `json:"is_out,omitempty"`            // 已弃用
	IsUninstall     int    `json:"is_uninstall,omitempty"`      // 应用是否未安装（0=已安装，1=未安装）
	AppName         string `json:"app_name,omitempty"`          // 关联应用名称
}

// FieldInfo 表示需求字段的详细信息（含候选值）
// API 文档：https://open.tapd.cn/document/api-doc/API文档/api_reference/story/get_story_fields_info.html
type FieldInfo struct {
	Name         string        `json:"name,omitempty"`          // 字段英文名
	Label        string        `json:"label,omitempty"`         // 字段中文名
	HTMLType     string        `json:"html_type,omitempty"`     // 控件类型（input/select/rich_edit 等）
	Options      interface{}   `json:"options,omitempty"`       // 候选值（object 或 array，因字段类型而异）
	PureOptions  interface{}   `json:"pure_options,omitempty"`  // 详细候选值（含层级元信息）
	Readonly     int           `json:"readonly,omitempty"`      // 是否只读（0=可编辑，1=只读）
	ColorOptions []ColorOption `json:"color_options,omitempty"` // 带颜色的候选值（如优先级）
	UserOptions  interface{}   `json:"user_options,omitempty"`  // 用户选择器的用户列表
	Memo         string        `json:"memo,omitempty"`          // 自定义字段备注
	EnableColor  int           `json:"enable_color,omitempty"`  // 是否启用颜色标记
}

// ColorOption 表示带颜色的候选值（如优先级字段）
type ColorOption struct {
	Value string `json:"value,omitempty"` // 选项值
	Color string `json:"color,omitempty"` // 颜色（十六进制色值）
	Label string `json:"label,omitempty"` // 显示名称
}

// WorkitemType 表示 TAPD 需求类别
// API 文档：https://open.tapd.cn/document/api-doc/API文档/api_reference/story/get_workitem_types.html
type WorkitemType struct {
	ID             string `json:"id,omitempty"`
	WorkspaceID    string `json:"workspace_id,omitempty"`
	AppID          string `json:"app_id,omitempty"`
	EntityType     string `json:"entity_type,omitempty"`       // 类别别名
	Name           string `json:"name,omitempty"`              // 类别名称
	EnglishName    string `json:"english_name,omitempty"`      // 英文名称
	Status         string `json:"status,omitempty"`            // 状态（1=未完成，2=未启用，3=已启用）
	Color          string `json:"color,omitempty"`             // 颜色
	WorkflowID     string `json:"workflow_id,omitempty"`       // 关联工作流 ID
	ChildrenIDs    string `json:"children_ids,omitempty"`      // 允许的子需求类别
	ParentIDs      string `json:"parent_ids,omitempty"`        // 允许的父需求类别
	Icon           string `json:"icon,omitempty"`              // 图标路径
	IconSmall      string `json:"icon_small,omitempty"`        // 小图标路径
	Creator        string `json:"creator,omitempty"`           // 创建人
	Created        string `json:"created,omitempty"`           // 创建时间
	ModifiedBy     string `json:"modified_by,omitempty"`       // 最后修改人
	Modified       string `json:"modified,omitempty"`          // 最后修改时间
	IconViper      string `json:"icon_viper,omitempty"`        // 图标完整 URL
	IconSmallViper string `json:"icon_small_viper,omitempty"`  // 小图标完整 URL
}

// StoryBugRelation 表示需求与缺陷的关联关系
// API 文档：https://open.tapd.cn/document/api-doc/API文档/api_reference/story/get_story_related_bugs.html
type StoryBugRelation struct {
	WorkspaceID int    `json:"workspace_id,omitempty"` // 项目 ID
	StoryID     string `json:"story_id,omitempty"`     // 需求 ID
	BugID       string `json:"bug_id,omitempty"`       // 缺陷 ID
}
