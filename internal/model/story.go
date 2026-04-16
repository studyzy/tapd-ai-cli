// Package model 中的 story.go 定义了 TAPD 需求数据模型
package model

// Story 表示 TAPD 需求/工作项，字段覆盖 TAPD API 返回的所有常用字段
// 使用强类型结构体反序列化可自动过滤 custom_field_* 等无用字段，节约 token
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/story/story.html
type Story struct {
	// 基本信息
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	WorkspaceID    string `json:"workspace_id,omitempty"`
	Status         string `json:"status,omitempty"`
	Type           string `json:"type,omitempty"`
	Source         string `json:"source,omitempty"`
	Flows          string `json:"flows,omitempty"`
	CreatedFrom    string `json:"created_from,omitempty"`
	WorkitemTypeID string `json:"workitem_type_id,omitempty"`
	TemplatedID    string `json:"templated_id,omitempty"`

	// 优先级
	Priority      string `json:"priority,omitempty"`
	PriorityLabel string `json:"priority_label,omitempty"`
	BusinessValue string `json:"business_value,omitempty"`

	// 人员相关
	Owner     string `json:"owner,omitempty"`
	Creator   string `json:"creator,omitempty"`
	Developer string `json:"developer,omitempty"`
	CC        string `json:"cc,omitempty"`

	// 时间相关
	Created   string `json:"created,omitempty"`
	Modified  string `json:"modified,omitempty"`
	Completed string `json:"completed,omitempty"`
	Begin     string `json:"begin,omitempty"`
	Due       string `json:"due,omitempty"`

	// 关联与分类
	IterationID string `json:"iteration_id,omitempty"`
	Module      string `json:"module,omitempty"`
	Feature     string `json:"feature,omitempty"`
	Label       string `json:"label,omitempty"`
	CategoryID  string `json:"category_id,omitempty"`
	ParentID    string `json:"parent_id,omitempty"`
	ChildrenID  string `json:"children_id,omitempty"`
	AncestorID  string `json:"ancestor_id,omitempty"`
	Path        string `json:"path,omitempty"`
	Level       string `json:"level,omitempty"`
	ReleaseID   string `json:"release_id,omitempty"`
	BugID       string `json:"bug_id,omitempty"`
	Version     string `json:"version,omitempty"`

	// 规模与工时
	Size            string `json:"size,omitempty"`
	Effort          string `json:"effort,omitempty"`
	EffortCompleted string `json:"effort_completed,omitempty"`
	Remain          string `json:"remain,omitempty"`
	Exceed          string `json:"exceed,omitempty"`

	// 进度与风险
	Progress   string `json:"progress,omitempty"`
	TechRisk   string `json:"tech_risk,omitempty"`
	TestFocus  string `json:"test_focus,omitempty"`
	IsArchived string `json:"is_archived,omitempty"`

	// 附加信息
	URL string `json:"url,omitempty"`
}
