// Package model 中的 iteration.go 定义了 TAPD 迭代数据模型
package model

// Iteration 表示 TAPD 迭代
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/iteration/get_iterations.html
type Iteration struct {
	// 基本信息
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	WorkspaceID    string `json:"workspace_id,omitempty"`
	Status         string `json:"status,omitempty"`
	EntityType     string `json:"entity_type,omitempty"`
	WorkitemTypeID string `json:"workitem_type_id,omitempty"`

	// 人员相关
	Creator string `json:"creator,omitempty"`
	Locker  string `json:"locker,omitempty"`

	// 时间相关
	StartDate string `json:"startdate,omitempty"`
	EndDate   string `json:"enddate,omitempty"`
	Created   string `json:"created,omitempty"`
	Modified  string `json:"modified,omitempty"`
	Completed string `json:"completed,omitempty"`

	// 关联与分类
	ReleaseID string `json:"release_id,omitempty"`
	ParentID  string `json:"parent_id,omitempty"`
	PlanAppID string `json:"plan_app_id,omitempty"`
	Label     string `json:"label,omitempty"`

	// 附加信息
	LockInfo string `json:"lock_info,omitempty"`
}

// ListIterationsRequest 查询迭代列表的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/iteration/get_iterations.html
type ListIterationsRequest struct {
	WorkspaceID    string // 必填：项目 ID
	ID             string // 可选：迭代 ID（支持多 ID 查询）
	Name           string // 可选：标题（支持模糊匹配）
	Description    string // 可选：详细描述
	StartDate      string // 可选：开始时间（支持时间查询）
	EndDate        string // 可选：结束时间（支持时间查询）
	WorkitemTypeID string // 可选：迭代类别
	PlanAppID      string // 可选：计划应用 ID
	Status         string // 可选：状态（open/done）
	Creator        string // 可选：创建人
	Created        string // 可选：创建时间（支持时间查询）
	Modified       string // 可选：最后修改时间（支持时间查询）
	Completed      string // 可选：完成时间
	Locker         string // 可选：锁定人
	Fields         string // 可选：返回字段列表
	Limit          string // 可选：返回数量限制
	Page           string // 可选：页码
	Order          string // 可选：排序规则
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *ListIterationsRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "id", r.ID)
	setOptional(params, "name", r.Name)
	setOptional(params, "description", r.Description)
	setOptional(params, "startdate", r.StartDate)
	setOptional(params, "enddate", r.EndDate)
	setOptional(params, "workitem_type_id", r.WorkitemTypeID)
	setOptional(params, "plan_app_id", r.PlanAppID)
	setOptional(params, "status", r.Status)
	setOptional(params, "creator", r.Creator)
	setOptional(params, "created", r.Created)
	setOptional(params, "modified", r.Modified)
	setOptional(params, "completed", r.Completed)
	setOptional(params, "locker", r.Locker)
	setOptional(params, "fields", r.Fields)
	setOptional(params, "limit", r.Limit)
	setOptional(params, "page", r.Page)
	setOptional(params, "order", r.Order)
	return params
}

// CreateIterationRequest 创建迭代的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/iteration/add_iteration.html
type CreateIterationRequest struct {
	WorkspaceID    string // 必填：项目 ID
	Name           string // 必填：标题
	StartDate      string // 必填：开始日期
	EndDate        string // 必填：结束日期
	Creator        string // 必填：创建人
	WorkitemTypeID string // 可选：迭代类别 ID
	PlanAppID      string // 可选：计划应用 ID
	EntityType     string // 可选：实体类型（iteration/release）
	ParentID       string // 可选：上层计划 ID
	Description    string // 可选：详细描述
	Status         string // 可选：状态（open/done）
	Label          string // 可选：标签（多个以竖线分隔）
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CreateIterationRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"name":         r.Name,
		"startdate":    r.StartDate,
		"enddate":      r.EndDate,
		"creator":      r.Creator,
	}
	setOptional(params, "workitem_type_id", r.WorkitemTypeID)
	setOptional(params, "plan_app_id", r.PlanAppID)
	setOptional(params, "entity_type", r.EntityType)
	setOptional(params, "parent_id", r.ParentID)
	setOptional(params, "description", r.Description)
	setOptional(params, "status", r.Status)
	setOptional(params, "label", r.Label)
	return params
}

// UpdateIterationRequest 更新迭代的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/iteration/update_iteration.html
type UpdateIterationRequest struct {
	WorkspaceID string // 必填：项目 ID
	ID          string // 必填：迭代 ID
	CurrentUser string // 必填：变更人
	Name        string // 可选：标题
	StartDate   string // 可选：开始日期
	EndDate     string // 可选：结束日期
	Description string // 可选：详细描述
	Status      string // 可选：状态（open/done）
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *UpdateIterationRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"id":           r.ID,
		"current_user": r.CurrentUser,
	}
	setOptional(params, "name", r.Name)
	setOptional(params, "startdate", r.StartDate)
	setOptional(params, "enddate", r.EndDate)
	setOptional(params, "description", r.Description)
	setOptional(params, "status", r.Status)
	return params
}

// CountIterationsRequest 查询迭代数量的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/iteration/get_iterations_count.html
type CountIterationsRequest struct {
	WorkspaceID    string // 必填：项目 ID
	ID             string // 可选：迭代 ID（支持多 ID 查询）
	Name           string // 可选：标题（支持模糊匹配）
	Description    string // 可选：详细描述
	StartDate      string // 可选：开始时间
	EndDate        string // 可选：结束时间
	WorkitemTypeID string // 可选：迭代类别
	PlanAppID      string // 可选：计划应用 ID
	Status         string // 可选：状态（open/done）
	Creator        string // 可选：创建人
	Created        string // 可选：创建时间（支持时间查询）
	Modified       string // 可选：最后修改时间（支持时间查询）
	Completed      string // 可选：完成时间
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CountIterationsRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "id", r.ID)
	setOptional(params, "name", r.Name)
	setOptional(params, "description", r.Description)
	setOptional(params, "startdate", r.StartDate)
	setOptional(params, "enddate", r.EndDate)
	setOptional(params, "workitem_type_id", r.WorkitemTypeID)
	setOptional(params, "plan_app_id", r.PlanAppID)
	setOptional(params, "status", r.Status)
	setOptional(params, "creator", r.Creator)
	setOptional(params, "created", r.Created)
	setOptional(params, "modified", r.Modified)
	setOptional(params, "completed", r.Completed)
	return params
}
