// Package model 中的 attachment.go 定义了 TAPD 附件和图片相关数据模型
package model

// Attachment 表示 TAPD 附件
type Attachment struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	EntryID     string `json:"entry_id,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Description string `json:"description,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Created     string `json:"created,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty"`
	Owner       string `json:"owner,omitempty"`
	DownloadURL string `json:"download_url,omitempty"`
}

// ImageInfo 表示 TAPD 图片下载信息
type ImageInfo struct {
	Type        string `json:"type,omitempty"`
	Value       string `json:"value,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty"`
	Filename    string `json:"filename,omitempty"`
	DownloadURL string `json:"download_url,omitempty"`
}
