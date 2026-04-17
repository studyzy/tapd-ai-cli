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

// GetImageRequest 获取图片的请求参数
type GetImageRequest struct {
	WorkspaceID string // 必填：项目 ID
	ImagePath   string // 必填：图片路径
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *GetImageRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "image_path", r.ImagePath)
	return params
}

// GetAttachmentsRequest 获取附件列表的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/attachment/get_attachments.html
type GetAttachmentsRequest struct {
	WorkspaceID string // 必填：项目 ID
	Type        string // 可选：对象类型
	EntryID     string // 可选：条目 ID
	Limit       string // 可选：返回数量限制
	Page        string // 可选：页码
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *GetAttachmentsRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "type", r.Type)
	setOptional(params, "entry_id", r.EntryID)
	setOptional(params, "limit", r.Limit)
	setOptional(params, "page", r.Page)
	return params
}
