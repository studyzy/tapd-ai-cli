// Package model 中的 release.go 定义了 TAPD 发布计划数据模型
package model

// Release 表示 TAPD 发布计划
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/release/
type Release struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty"`
	Description string `json:"description,omitempty"`
	StartDate   string `json:"startdate,omitempty"`
	EndDate     string `json:"enddate,omitempty"`
	Status      string `json:"status,omitempty"`
	Creator     string `json:"creator,omitempty"`
	Created     string `json:"created,omitempty"`
	Modified    string `json:"modified,omitempty"`
}
