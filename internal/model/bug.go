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

// ListBugsRequest 查询缺陷列表的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/bug/get_bugs.html
type ListBugsRequest struct {
	WorkspaceID   string // 必填：项目 ID
	ID            string // 可选：缺陷 ID
	Title         string // 可选：标题（支持模糊匹配）
	PriorityLabel string // 可选：优先级
	Severity      string // 可选：严重程度
	Status        string // 可选：状态
	VStatus       string // 可选：中文状态名
	CurrentOwner  string // 可选：处理人
	Reporter      string // 可选：创建人
	IterationID   string // 可选：迭代 ID
	Module        string // 可选：模块
	Label         string // 可选：标签
	Fields        string // 可选：返回字段列表
	Limit         string // 可选：返回数量限制
	Page          string // 可选：页码
	Order         string // 可选：排序规则
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *ListBugsRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "id", r.ID)
	setOptional(params, "title", r.Title)
	setOptional(params, "priority_label", r.PriorityLabel)
	setOptional(params, "severity", r.Severity)
	setOptional(params, "status", r.Status)
	setOptional(params, "v_status", r.VStatus)
	setOptional(params, "current_owner", r.CurrentOwner)
	setOptional(params, "reporter", r.Reporter)
	setOptional(params, "iteration_id", r.IterationID)
	setOptional(params, "module", r.Module)
	setOptional(params, "label", r.Label)
	setOptional(params, "fields", r.Fields)
	setOptional(params, "limit", r.Limit)
	setOptional(params, "page", r.Page)
	setOptional(params, "order", r.Order)
	return params
}

// CreateBugRequest 创建缺陷的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/bug/add_bug.html
type CreateBugRequest struct {
	WorkspaceID   string // 必填：项目 ID
	Title         string // 必填：缺陷标题
	PriorityLabel string // 可选：优先级
	Severity      string // 可选：严重程度
	Description   string // 可选：详细描述
	CurrentOwner  string // 可选：处理人
	Reporter      string // 可选：创建人
	DE            string // 可选：开发人员
	TE            string // 可选：测试人员
	Module        string // 可选：模块
	IterationID   string // 可选：迭代 ID
	Label         string // 可选：标签
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CreateBugRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"title":        r.Title,
	}
	setOptional(params, "priority_label", r.PriorityLabel)
	setOptional(params, "severity", r.Severity)
	setOptional(params, "description", r.Description)
	setOptional(params, "current_owner", r.CurrentOwner)
	setOptional(params, "reporter", r.Reporter)
	setOptional(params, "de", r.DE)
	setOptional(params, "te", r.TE)
	setOptional(params, "module", r.Module)
	setOptional(params, "iteration_id", r.IterationID)
	setOptional(params, "label", r.Label)
	return params
}

// UpdateBugRequest 更新缺陷的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/bug/update_bug.html
type UpdateBugRequest struct {
	WorkspaceID   string // 必填：项目 ID
	ID            string // 必填：缺陷 ID
	Title         string // 可选：标题
	PriorityLabel string // 可选：优先级
	Severity      string // 可选：严重程度
	Status        string // 可选：状态
	VStatus       string // 可选：中文状态名
	Description   string // 可选：详细描述
	CurrentOwner  string // 可选：处理人
	CurrentUser   string // 可选：变更人
	Module        string // 可选：模块
	Label         string // 可选：标签
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *UpdateBugRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"id":           r.ID,
	}
	setOptional(params, "title", r.Title)
	setOptional(params, "priority_label", r.PriorityLabel)
	setOptional(params, "severity", r.Severity)
	setOptional(params, "status", r.Status)
	setOptional(params, "v_status", r.VStatus)
	setOptional(params, "description", r.Description)
	setOptional(params, "current_owner", r.CurrentOwner)
	setOptional(params, "current_user", r.CurrentUser)
	setOptional(params, "module", r.Module)
	setOptional(params, "label", r.Label)
	return params
}

// CountBugsRequest 查询缺陷数量的请求参数
type CountBugsRequest struct {
	WorkspaceID string // 必填：项目 ID
	Status      string // 可选：状态
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CountBugsRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "status", r.Status)
	return params
}
