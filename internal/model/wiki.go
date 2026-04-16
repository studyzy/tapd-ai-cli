// Package model 中的 wiki.go 定义了 TAPD Wiki 文档数据模型及请求参数结构体
package model

// Wiki 表示 TAPD Wiki 文档，字段覆盖 TAPD API 返回的所有常用字段
// 使用强类型结构体反序列化可自动过滤无用字段，节约 token
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/wiki/get_tapd_wikis.html
type Wiki struct {
	// 基本信息
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	WorkspaceID         string `json:"workspace_id,omitempty"`
	Description         string `json:"description,omitempty"`
	MarkdownDescription string `json:"markdown_description,omitempty"`
	IsRich              string `json:"is_rich,omitempty"`

	// 层级关系
	ParentWikiID string `json:"parent_wiki_id,omitempty"`

	// 人员相关
	Creator  string `json:"creator,omitempty"`
	Modifier string `json:"modifier,omitempty"`

	// 时间相关
	Created  string `json:"created,omitempty"`
	Modified string `json:"modified,omitempty"`

	// 统计信息
	ViewCount string `json:"view_count,omitempty"`
	Note      string `json:"note,omitempty"`

	// 附加信息
	URL string `json:"url,omitempty"`
}

// ListWikisRequest 查询 Wiki 列表的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/wiki/get_tapd_wikis.html
type ListWikisRequest struct {
	WorkspaceID string // 必填：项目 ID
	ID          string // 可选：Wiki ID
	Name        string // 可选：标题
	Creator     string // 可选：创建人
	Fields      string // 可选：返回字段列表
	Limit       string // 可选：返回数量限制
	Page        string // 可选：页码
	Order       string // 可选：排序规则
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *ListWikisRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "id", r.ID)
	setOptional(params, "name", r.Name)
	setOptional(params, "creator", r.Creator)
	setOptional(params, "fields", r.Fields)
	setOptional(params, "limit", r.Limit)
	setOptional(params, "page", r.Page)
	setOptional(params, "order", r.Order)
	return params
}

// CreateWikiRequest 创建 Wiki 的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/wiki/add_tapd_wiki.html
type CreateWikiRequest struct {
	WorkspaceID         string // 必填：项目 ID
	Name                string // 必填：标题
	Creator             string // 必填：创建人
	MarkdownDescription string // 可选：Markdown 格式内容
	Description         string // 可选：富文本格式内容
	Note                string // 可选：备注
	ParentWikiID        string // 可选：父 Wiki ID
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CreateWikiRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"name":         r.Name,
		"creator":      r.Creator,
	}
	setOptional(params, "markdown_description", r.MarkdownDescription)
	setOptional(params, "description", r.Description)
	setOptional(params, "note", r.Note)
	setOptional(params, "parent_wiki_id", r.ParentWikiID)
	return params
}

// UpdateWikiRequest 更新 Wiki 的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/wiki/update_tapd_wiki.html
type UpdateWikiRequest struct {
	WorkspaceID         string // 必填：项目 ID
	ID                  string // 必填：Wiki ID
	Name                string // 可选：新标题
	MarkdownDescription string // 可选：新 Markdown 内容
	Description         string // 可选：新富文本内容
	Note                string // 可选：新备注
	ParentWikiID        string // 可选：新父 Wiki ID
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *UpdateWikiRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"id":           r.ID,
	}
	setOptional(params, "name", r.Name)
	setOptional(params, "markdown_description", r.MarkdownDescription)
	setOptional(params, "description", r.Description)
	setOptional(params, "note", r.Note)
	setOptional(params, "parent_wiki_id", r.ParentWikiID)
	return params
}

// setOptional 当值非空时添加到参数 map
func setOptional(params map[string]string, key, value string) {
	if value != "" {
		params[key] = value
	}
}
