// Package model 中的 bug.go 定义了 TAPD 缺陷数据模型
package model

// Bug 表示 TAPD 缺陷，字段覆盖 TAPD API 返回的所有常用字段
// 使用强类型结构体反序列化可自动过滤 custom_field_* 等无用字段，节约 token
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/bug/bug.html
type Bug struct {
	// 基本信息
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty"`
	Status      string `json:"status,omitempty"`
	BugType     string `json:"bugtype,omitempty"`
	Source      string `json:"source,omitempty"`
	Resolution  string `json:"resolution,omitempty"`
	Flows       string `json:"flows,omitempty"`

	// 优先级与严重程度
	Priority      string `json:"priority,omitempty"`
	PriorityLabel string `json:"priority_label,omitempty"`
	Severity      string `json:"severity,omitempty"`

	// 人员相关
	CurrentOwner string `json:"current_owner,omitempty"`
	Reporter     string `json:"reporter,omitempty"`
	Fixer        string `json:"fixer,omitempty"`
	Closer       string `json:"closer,omitempty"`
	CC           string `json:"cc,omitempty"`
	Participator string `json:"participator,omitempty"`
	TE           string `json:"te,omitempty"`
	DE           string `json:"de,omitempty"`
	Auditer      string `json:"auditer,omitempty"`
	Confirmer    string `json:"confirmer,omitempty"`
	LastModify   string `json:"lastmodify,omitempty"`

	// 时间相关
	Created        string `json:"created,omitempty"`
	Modified       string `json:"modified,omitempty"`
	Resolved       string `json:"resolved,omitempty"`
	Closed         string `json:"closed,omitempty"`
	InProgressTime string `json:"in_progress_time,omitempty"`
	VerifyTime     string `json:"verify_time,omitempty"`
	RejectTime     string `json:"reject_time,omitempty"`
	Begin          string `json:"begin,omitempty"`
	Due            string `json:"due,omitempty"`
	Deadline       string `json:"deadline,omitempty"`

	// 关联与分类
	IterationID string `json:"iteration_id,omitempty"`
	Module      string `json:"module,omitempty"`
	Feature     string `json:"feature,omitempty"`
	StoryID     string `json:"story_id,omitempty"`
	ReleaseID   string `json:"release_id,omitempty"`
	Label       string `json:"label,omitempty"`
	CreatedFrom string `json:"created_from,omitempty"`

	// 版本相关
	VersionReport string `json:"version_report,omitempty"`
	VersionTest   string `json:"version_test,omitempty"`
	VersionFix    string `json:"version_fix,omitempty"`
	VersionClose  string `json:"version_close,omitempty"`

	// 测试相关
	OS          string `json:"os,omitempty"`
	Platform    string `json:"platform,omitempty"`
	TestMode    string `json:"testmode,omitempty"`
	TestPhase   string `json:"testphase,omitempty"`
	TestType    string `json:"testtype,omitempty"`
	OriginPhase string `json:"originphase,omitempty"`
	SourcePhase string `json:"sourcephase,omitempty"`
	Frequency   string `json:"frequency,omitempty"`

	// 工时相关
	Effort          string `json:"effort,omitempty"`
	EffortCompleted string `json:"effort_completed,omitempty"`
	Remain          string `json:"remain,omitempty"`
	Exceed          string `json:"exceed,omitempty"`
	Estimate        string `json:"estimate,omitempty"`

	// 附加信息
	URL string `json:"url,omitempty"`
}
