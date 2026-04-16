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
