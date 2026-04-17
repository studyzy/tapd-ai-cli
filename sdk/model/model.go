// Package model 定义了 tapd-sdk-go 的所有数据模型结构体
package model

import "encoding/json"

// Workspace 表示 TAPD 项目/工作区
type Workspace struct {
	ID                string `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	PrettyName        string `json:"pretty_name,omitempty"`
	Category          string `json:"category,omitempty"`
	Status            string `json:"status,omitempty"`
	Description       string `json:"description,omitempty"`
	BeginDate         string `json:"begin_date,omitempty"`
	EndDate           string `json:"end_date,omitempty"`
	Closed            string `json:"closed,omitempty"`
	ExternalOn        string `json:"external_on,omitempty"`
	Creator           string `json:"creator,omitempty"`
	Created           string `json:"created,omitempty"`
	ProductType       string `json:"product_type,omitempty"`
	PlatformType      string `json:"platform_type,omitempty"`
	IsSelfDevelopment string `json:"is_self_development,omitempty"`
	Objective         string `json:"objective,omitempty"`
	Secrecy           string `json:"secrecy,omitempty"`
	Schedule          string `json:"schedule,omitempty"`
	Milestone         string `json:"milestone,omitempty"`
	Risk              string `json:"risk,omitempty"`
	CompanyID         string `json:"company_id,omitempty"`
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

// TAPDResponse 表示 TAPD API 的统一响应包装格式
type TAPDResponse struct {
	Status int             `json:"status"`
	Data   json.RawMessage `json:"data"`
	Info   string          `json:"info"`
}
