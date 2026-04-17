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

// ListCommentsRequest 查询评论列表的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/comment/get_comments.html
type ListCommentsRequest struct {
	WorkspaceID string // 必填：项目 ID
	EntryType   string // 可选：评论类型（stories|bug|tasks）
	EntryID     string // 可选：条目 ID
	Author      string // 可选：评论人
	Fields      string // 可选：返回字段列表
	Limit       string // 可选：返回数量限制
	Page        string // 可选：页码
	Order       string // 可选：排序规则
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *ListCommentsRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "entry_type", r.EntryType)
	setOptional(params, "entry_id", r.EntryID)
	setOptional(params, "author", r.Author)
	setOptional(params, "fields", r.Fields)
	setOptional(params, "limit", r.Limit)
	setOptional(params, "page", r.Page)
	setOptional(params, "order", r.Order)
	return params
}

// AddCommentRequest 添加评论的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/comment/add_comment.html
type AddCommentRequest struct {
	WorkspaceID string // 必填：项目 ID
	EntryType   string // 必填：评论类型（stories|bug|tasks）
	EntryID     string // 必填：条目 ID
	Description string // 必填：评论内容
	Author      string // 必填：评论人
	ReplyID     string // 可选：回复的评论 ID
	RootID      string // 可选：根评论 ID
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *AddCommentRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"entry_type":   r.EntryType,
		"entry_id":     r.EntryID,
		"description":  r.Description,
	}
	setOptional(params, "author", r.Author)
	setOptional(params, "reply_id", r.ReplyID)
	setOptional(params, "root_id", r.RootID)
	return params
}

// UpdateCommentRequest 更新评论的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/comment/update_comment.html
type UpdateCommentRequest struct {
	WorkspaceID   string // 必填：项目 ID
	ID            string // 必填：评论 ID
	Description   string // 必填：评论内容
	ChangeCreator string // 可选：变更人
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *UpdateCommentRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"id":           r.ID,
		"description":  r.Description,
	}
	setOptional(params, "change_creator", r.ChangeCreator)
	return params
}

// CountCommentsRequest 查询评论数量的请求参数
type CountCommentsRequest struct {
	WorkspaceID string // 必填：项目 ID
	EntryType   string // 可选：评论类型
	EntryID     string // 可选：条目 ID
	Author      string // 可选：评论人
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CountCommentsRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "entry_type", r.EntryType)
	setOptional(params, "entry_id", r.EntryID)
	setOptional(params, "author", r.Author)
	return params
}
