// Package model 定义了 tapd-ai-cli 使用的所有数据模型结构体
package model

import "encoding/json"

// Config 表示本地持久化的配置数据，存储于 .tapd.json
type Config struct {
	AccessToken string `json:"access_token,omitempty"`
	APIUser     string `json:"api_user,omitempty"`
	APIPassword string `json:"api_password,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty"`
}

// Workspace 表示 TAPD 项目/工作区
type Workspace struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Status   string `json:"status,omitempty"`
	Category string `json:"category,omitempty"`
	Creator  string `json:"creator,omitempty"`
	Created  string `json:"created,omitempty"`
}

// Iteration 表示 TAPD 迭代
type Iteration struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Status      string `json:"status,omitempty"`
	StartDate   string `json:"startdate,omitempty"`
	EndDate     string `json:"enddate,omitempty"`
	Description string `json:"description,omitempty"`
	Creator     string `json:"creator,omitempty"`
	Created     string `json:"created,omitempty"`
	Modified    string `json:"modified,omitempty"`
	Completed   string `json:"completed,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty"`
}

// ListIterationsRequest 查询迭代列表的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/iteration/get_iterations.html
type ListIterationsRequest struct {
	WorkspaceID string // 必填：项目 ID
	ID          string // 可选：迭代 ID
	Name        string // 可选：标题（支持模糊匹配）
	Status      string // 可选：状态（open/done）
	Fields      string // 可选：返回字段列表
	Limit       string // 可选：返回数量限制
	Page        string // 可选：页码
	Order       string // 可选：排序规则
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *ListIterationsRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "id", r.ID)
	setOptional(params, "name", r.Name)
	setOptional(params, "status", r.Status)
	setOptional(params, "fields", r.Fields)
	setOptional(params, "limit", r.Limit)
	setOptional(params, "page", r.Page)
	setOptional(params, "order", r.Order)
	return params
}

// CreateIterationRequest 创建迭代的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/iteration/add_iteration.html
type CreateIterationRequest struct {
	WorkspaceID string // 必填：项目 ID
	Name        string // 必填：标题
	StartDate   string // 必填：开始日期
	EndDate     string // 必填：结束日期
	Creator     string // 必填：创建人
	Description string // 可选：详细描述
	Status      string // 可选：状态（open/done）
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
	setOptional(params, "description", r.Description)
	setOptional(params, "status", r.Status)
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

// ListResponse 表示列表查询的通用响应结构
type ListResponse struct {
	Items   interface{} `json:"items"`
	Total   int         `json:"total,omitempty"`
	Page    int         `json:"page,omitempty"`
	Limit   int         `json:"limit,omitempty"`
	HasMore bool        `json:"has_more,omitempty"`
}

// SuccessResponse 表示创建/更新操作的成功响应
type SuccessResponse struct {
	Success     bool   `json:"success"`
	ID          string `json:"id,omitempty"`
	URL         string `json:"url,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty"`
}

// CountResponse 表示计数查询的响应
type CountResponse struct {
	Count int `json:"count"`
}

// ErrorResponse 表示输出到 stderr 的错误信息
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
}

// TAPDResponse 表示 TAPD API 的统一响应包装格式
type TAPDResponse struct {
	Status int             `json:"status"`
	Data   json.RawMessage `json:"data"`
	Info   string          `json:"info"`
}
