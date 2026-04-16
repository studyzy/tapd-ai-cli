// Package model 中的 comment.go 定义了 TAPD 评论数据模型
package model

// Comment 表示 TAPD 评论，字段覆盖 TAPD API 返回的所有常用字段
// 使用强类型结构体反序列化可自动过滤无用字段，节约 token
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/comment/get_comments.html
type Comment struct {
	// 基本信息
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty"`

	// 关联信息
	EntryType string `json:"entry_type,omitempty"`
	EntryID   string `json:"entry_id,omitempty"`

	// 回复层级
	RootID  string `json:"root_id,omitempty"`
	ReplyID string `json:"reply_id,omitempty"`

	// 人员相关
	Author string `json:"author,omitempty"`

	// 时间相关
	Created  string `json:"created,omitempty"`
	Modified string `json:"modified,omitempty"`
}
