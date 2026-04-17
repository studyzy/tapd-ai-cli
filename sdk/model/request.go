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
