// Package model 中的 category.go 定义了 TAPD 需求分类数据模型
package model

// Category 表示 TAPD 需求分类
type Category struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ParentID    string `json:"parent_id,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty"`
	Creator     string `json:"creator,omitempty"`
	Modifier    string `json:"modifier,omitempty"`
	Created     string `json:"created,omitempty"`
	Modified    string `json:"modified,omitempty"`
}
