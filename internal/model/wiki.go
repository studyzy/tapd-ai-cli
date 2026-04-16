// Package model 中的 wiki.go 定义了 TAPD Wiki 文档数据模型
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
