// Package model 中的 bug.go 定义了 TAPD 缺陷数据模型
package model

import "encoding/json"

// Bug 表示 TAPD 缺陷，字段覆盖 TAPD API 返回的所有常用字段
// 自定义字段（custom_field_*、custom_plan_field_*）通过 CustomFields map 保留，不会丢失
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
	ReopenTime     string `json:"reopen_time,omitempty"`
	AuditTime      string `json:"audit_time,omitempty"`
	SuspendTime    string `json:"suspend_time,omitempty"`
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
	Milestone   string `json:"milestone,omitempty"`
	IssueID     string `json:"issue_id,omitempty"`

	// 版本相关
	VersionReport string `json:"version_report,omitempty"`
	VersionTest   string `json:"version_test,omitempty"`
	VersionFix    string `json:"version_fix,omitempty"`
	VersionClose  string `json:"version_close,omitempty"`

	// 基线相关
	BaselineFind  string `json:"baseline_find,omitempty"`
	BaselineJoin  string `json:"baseline_join,omitempty"`
	BaselineTest  string `json:"baseline_test,omitempty"`
	BaselineClose string `json:"baseline_close,omitempty"`

	// 测试相关
	OS               string `json:"os,omitempty"`
	Platform         string `json:"platform,omitempty"`
	TestMode         string `json:"testmode,omitempty"`
	TestPhase        string `json:"testphase,omitempty"`
	TestType         string `json:"testtype,omitempty"`
	OriginPhase      string `json:"originphase,omitempty"`
	SourcePhase      string `json:"sourcephase,omitempty"`
	Frequency        string `json:"frequency,omitempty"`
	RegressionNumber string `json:"regression_number,omitempty"`

	// 工时相关
	Effort          string `json:"effort,omitempty"`
	EffortCompleted string `json:"effort_completed,omitempty"`
	Remain          string `json:"remain,omitempty"`
	Exceed          string `json:"exceed,omitempty"`
	Estimate        string `json:"estimate,omitempty"`

	// 附加信息
	URL string `json:"url,omitempty"`

	// 自定义字段，key 为 custom_field_one、custom_field_9 等
	CustomFields map[string]string `json:"-"`
}

// UnmarshalJSON 自定义反序列化，在解析标准字段的同时收集 custom_field_* 和 custom_plan_field_* 字段
func (b *Bug) UnmarshalJSON(data []byte) error {
	type Alias Bug
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*b = Bug(alias)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	b.CustomFields = ExtractCustomFields(raw)
	return nil
}

// MarshalJSON 自定义序列化，将 CustomFields 中的键值对合并到输出 JSON
func (b Bug) MarshalJSON() ([]byte, error) {
	type Alias Bug
	data, err := json.Marshal(Alias(b))
	if err != nil {
		return nil, err
	}
	if len(b.CustomFields) == 0 {
		return data, nil
	}

	var base map[string]json.RawMessage
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}
	for k, v := range b.CustomFields {
		raw, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		base[k] = raw
	}
	return json.Marshal(base)
}

// ListBugsRequest 查询缺陷列表的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/bug/get_bugs.html
type ListBugsRequest struct {
	// 必填参数
	WorkspaceID string // 必填：项目 ID

	// 基本筛选
	ID            string // 可选：缺陷 ID，支持多 ID 查询
	Title         string // 可选：标题（支持模糊匹配）
	Priority      string // 可选：优先级（建议使用 PriorityLabel）
	PriorityLabel string // 可选：优先级（推荐）
	Severity      string // 可选：严重程度，支持枚举查询
	Status        string // 可选：状态，支持不等于查询、枚举查询
	VStatus       string // 可选：中文状态名
	BugType       string // 可选：缺陷类型
	Source        string // 可选：缺陷根源，支持枚举查询
	Resolution    string // 可选：解决方法，支持枚举查询
	Description   string // 可选：详细描述（支持模糊匹配）

	// 人员筛选
	CurrentOwner string // 可选：处理人（支持模糊匹配）
	Reporter     string // 可选：创建人，支持多人员查询
	CC           string // 可选：抄送人
	Participator string // 可选：参与人，支持多人员查询
	TE           string // 可选：测试人员（支持模糊匹配）
	DE           string // 可选：开发人员（支持模糊匹配）
	Auditer      string // 可选：审核人
	Confirmer    string // 可选：验证人
	Fixer        string // 可选：修复人
	Closer       string // 可选：关闭人
	LastModify   string // 可选：最后修改人

	// 时间筛选
	Created        string // 可选：创建时间，支持时间查询
	Modified       string // 可选：最后修改时间，支持时间查询
	Resolved       string // 可选：解决时间，支持时间查询
	Closed         string // 可选：关闭时间，支持时间查询
	InProgressTime string // 可选：接受处理时间，支持时间查询
	VerifyTime     string // 可选：验证时间，支持时间查询
	RejectTime     string // 可选：拒绝时间，支持时间查询
	Begin          string // 可选：预计开始
	Due            string // 可选：预计结束
	Deadline       string // 可选：解决期限

	// 关联与分类筛选
	IterationID string // 可选：迭代 ID，支持枚举查询
	Module      string // 可选：模块，支持枚举查询
	Feature     string // 可选：特性
	ReleaseID   string // 可选：发布计划
	Label       string // 可选：标签，支持枚举查询

	// 版本筛选
	VersionReport string // 可选：发现版本，支持枚举查询
	VersionTest   string // 可选：验证版本
	VersionFix    string // 可选：合入版本
	VersionClose  string // 可选：关闭版本

	// 基线筛选
	BaselineFind  string // 可选：发现基线
	BaselineJoin  string // 可选：合入基线
	BaselineTest  string // 可选：验证基线
	BaselineClose string // 可选：关闭基线

	// 测试相关筛选
	OS          string // 可选：操作系统
	Size        string // 可选：规模
	Platform    string // 可选：软件平台
	TestMode    string // 可选：测试方式
	TestPhase   string // 可选：测试阶段
	TestType    string // 可选：测试类型
	Frequency   string // 可选：重现规律，支持枚举查询
	OriginPhase string // 可选：发现阶段
	SourcePhase string // 可选：引入阶段
	Estimate    string // 可选：预计解决时间

	// 分页与排序
	Fields string // 可选：返回字段列表，多个字段以逗号分隔
	Limit  string // 可选：返回数量限制，默认 30，最大 200
	Page   string // 可选：页码，默认 1
	Order  string // 可选：排序规则，如 created desc
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *ListBugsRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "id", r.ID)
	setOptional(params, "title", r.Title)
	setOptional(params, "priority", r.Priority)
	setOptional(params, "priority_label", r.PriorityLabel)
	setOptional(params, "severity", r.Severity)
	setOptional(params, "status", r.Status)
	setOptional(params, "v_status", r.VStatus)
	setOptional(params, "bugtype", r.BugType)
	setOptional(params, "source", r.Source)
	setOptional(params, "resolution", r.Resolution)
	setOptional(params, "description", r.Description)
	setOptional(params, "current_owner", r.CurrentOwner)
	setOptional(params, "reporter", r.Reporter)
	setOptional(params, "cc", r.CC)
	setOptional(params, "participator", r.Participator)
	setOptional(params, "te", r.TE)
	setOptional(params, "de", r.DE)
	setOptional(params, "auditer", r.Auditer)
	setOptional(params, "confirmer", r.Confirmer)
	setOptional(params, "fixer", r.Fixer)
	setOptional(params, "closer", r.Closer)
	setOptional(params, "lastmodify", r.LastModify)
	setOptional(params, "created", r.Created)
	setOptional(params, "modified", r.Modified)
	setOptional(params, "resolved", r.Resolved)
	setOptional(params, "closed", r.Closed)
	setOptional(params, "in_progress_time", r.InProgressTime)
	setOptional(params, "verify_time", r.VerifyTime)
	setOptional(params, "reject_time", r.RejectTime)
	setOptional(params, "begin", r.Begin)
	setOptional(params, "due", r.Due)
	setOptional(params, "deadline", r.Deadline)
	setOptional(params, "iteration_id", r.IterationID)
	setOptional(params, "module", r.Module)
	setOptional(params, "feature", r.Feature)
	setOptional(params, "release_id", r.ReleaseID)
	setOptional(params, "label", r.Label)
	setOptional(params, "version_report", r.VersionReport)
	setOptional(params, "version_test", r.VersionTest)
	setOptional(params, "version_fix", r.VersionFix)
	setOptional(params, "version_close", r.VersionClose)
	setOptional(params, "baseline_find", r.BaselineFind)
	setOptional(params, "baseline_join", r.BaselineJoin)
	setOptional(params, "baseline_test", r.BaselineTest)
	setOptional(params, "baseline_close", r.BaselineClose)
	setOptional(params, "os", r.OS)
	setOptional(params, "size", r.Size)
	setOptional(params, "platform", r.Platform)
	setOptional(params, "testmode", r.TestMode)
	setOptional(params, "testphase", r.TestPhase)
	setOptional(params, "testtype", r.TestType)
	setOptional(params, "frequency", r.Frequency)
	setOptional(params, "originphase", r.OriginPhase)
	setOptional(params, "sourcephase", r.SourcePhase)
	setOptional(params, "estimate", r.Estimate)
	setOptional(params, "fields", r.Fields)
	setOptional(params, "limit", r.Limit)
	setOptional(params, "page", r.Page)
	setOptional(params, "order", r.Order)
	return params
}

// CreateBugRequest 创建缺陷的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/bug/add_bug.html
type CreateBugRequest struct {
	// 必填参数
	WorkspaceID string // 必填：项目 ID
	Title       string // 必填：缺陷标题

	// 基本信息
	Description   string // 可选：详细描述
	Priority      string // 可选：优先级（建议使用 PriorityLabel）
	PriorityLabel string // 可选：优先级（推荐）
	Severity      string // 可选：严重程度
	BugType       string // 可选：缺陷类型
	Source        string // 可选：缺陷根源
	Resolution    string // 可选：解决方法

	// 人员相关
	CurrentOwner string // 可选：处理人
	Reporter     string // 可选：创建人
	CC           string // 可选：抄送人
	Participator string // 可选：参与人
	TE           string // 可选：测试人员
	DE           string // 可选：开发人员
	Auditer      string // 可选：审核人
	Confirmer    string // 可选：验证人
	Fixer        string // 可选：修复人
	Closer       string // 可选：关闭人
	LastModify   string // 可选：最后修改人

	// 时间相关
	InProgressTime string // 可选：接受处理时间
	VerifyTime     string // 可选：验证时间
	RejectTime     string // 可选：拒绝时间
	Begin          string // 可选：预计开始
	Due            string // 可选：预计结束
	Deadline       string // 可选：解决期限

	// 关联与分类
	IterationID string // 可选：迭代 ID
	Module      string // 可选：模块
	Feature     string // 可选：特性
	ReleaseID   string // 可选：发布计划
	Label       string // 可选：标签，多个以竖线分隔

	// 版本相关
	VersionReport string // 可选：发现版本
	VersionTest   string // 可选：验证版本
	VersionFix    string // 可选：合入版本
	VersionClose  string // 可选：关闭版本

	// 基线相关
	BaselineFind  string // 可选：发现基线
	BaselineJoin  string // 可选：合入基线
	BaselineTest  string // 可选：验证基线
	BaselineClose string // 可选：关闭基线

	// 测试相关
	OS          string // 可选：操作系统
	Size        string // 可选：规模
	Platform    string // 可选：软件平台
	TestMode    string // 可选：测试方式
	TestPhase   string // 可选：测试阶段
	TestType    string // 可选：测试类型
	Frequency   string // 可选：重现规律
	OriginPhase string // 可选：发现阶段
	SourcePhase string // 可选：引入阶段
	Estimate    string // 可选：预计解决时间

	// 工时相关
	Effort string // 可选：预估工时

	// 自定义字段
	CustomFields map[string]string // 可选：自定义字段，key 如 custom_field_one、custom_field_9
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CreateBugRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"title":        r.Title,
	}
	setOptional(params, "description", r.Description)
	setOptional(params, "priority", r.Priority)
	setOptional(params, "priority_label", r.PriorityLabel)
	setOptional(params, "severity", r.Severity)
	setOptional(params, "bugtype", r.BugType)
	setOptional(params, "source", r.Source)
	setOptional(params, "resolution", r.Resolution)
	setOptional(params, "current_owner", r.CurrentOwner)
	setOptional(params, "reporter", r.Reporter)
	setOptional(params, "cc", r.CC)
	setOptional(params, "participator", r.Participator)
	setOptional(params, "te", r.TE)
	setOptional(params, "de", r.DE)
	setOptional(params, "auditer", r.Auditer)
	setOptional(params, "confirmer", r.Confirmer)
	setOptional(params, "fixer", r.Fixer)
	setOptional(params, "closer", r.Closer)
	setOptional(params, "lastmodify", r.LastModify)
	setOptional(params, "in_progress_time", r.InProgressTime)
	setOptional(params, "verify_time", r.VerifyTime)
	setOptional(params, "reject_time", r.RejectTime)
	setOptional(params, "begin", r.Begin)
	setOptional(params, "due", r.Due)
	setOptional(params, "deadline", r.Deadline)
	setOptional(params, "iteration_id", r.IterationID)
	setOptional(params, "module", r.Module)
	setOptional(params, "feature", r.Feature)
	setOptional(params, "release_id", r.ReleaseID)
	setOptional(params, "label", r.Label)
	setOptional(params, "version_report", r.VersionReport)
	setOptional(params, "version_test", r.VersionTest)
	setOptional(params, "version_fix", r.VersionFix)
	setOptional(params, "version_close", r.VersionClose)
	setOptional(params, "baseline_find", r.BaselineFind)
	setOptional(params, "baseline_join", r.BaselineJoin)
	setOptional(params, "baseline_test", r.BaselineTest)
	setOptional(params, "baseline_close", r.BaselineClose)
	setOptional(params, "os", r.OS)
	setOptional(params, "size", r.Size)
	setOptional(params, "platform", r.Platform)
	setOptional(params, "testmode", r.TestMode)
	setOptional(params, "testphase", r.TestPhase)
	setOptional(params, "testtype", r.TestType)
	setOptional(params, "frequency", r.Frequency)
	setOptional(params, "originphase", r.OriginPhase)
	setOptional(params, "sourcephase", r.SourcePhase)
	setOptional(params, "estimate", r.Estimate)
	setOptional(params, "effort", r.Effort)
	MergeCustomFields(params, r.CustomFields)
	return params
}

// UpdateBugRequest 更新缺陷的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/bug/update_bug.html
type UpdateBugRequest struct {
	// 必填参数
	WorkspaceID string // 必填：项目 ID
	ID          string // 必填：缺陷 ID

	// 基本信息
	Title         string // 可选：标题
	Description   string // 可选：详细描述
	Priority      string // 可选：优先级（建议使用 PriorityLabel）
	PriorityLabel string // 可选：优先级（推荐）
	Severity      string // 可选：严重程度
	Status        string // 可选：状态
	VStatus       string // 可选：中文状态名
	BugType       string // 可选：缺陷类型
	Source        string // 可选：缺陷根源
	Resolution    string // 可选：解决方法

	// 人员相关
	CurrentOwner string // 可选：处理人
	CurrentUser  string // 可选：变更人
	Reporter     string // 可选：创建人
	CC           string // 可选：抄送人
	Participator string // 可选：参与人
	TE           string // 可选：测试人员
	DE           string // 可选：开发人员
	Auditer      string // 可选：审核人
	Confirmer    string // 可选：验证人
	Fixer        string // 可选：修复人
	Closer       string // 可选：关闭人
	LastModify   string // 可选：最后修改人

	// 时间相关
	InProgressTime string // 可选：接受处理时间
	VerifyTime     string // 可选：验证时间
	RejectTime     string // 可选：拒绝时间
	Begin          string // 可选：预计开始
	Due            string // 可选：预计结束
	Deadline       string // 可选：解决期限

	// 关联与分类
	IterationID string // 可选：迭代 ID
	Module      string // 可选：模块
	Feature     string // 可选：特性
	ReleaseID   string // 可选：发布计划
	Label       string // 可选：标签，多个以竖线分隔

	// 版本相关
	VersionReport string // 可选：发现版本
	VersionTest   string // 可选：验证版本
	VersionFix    string // 可选：合入版本
	VersionClose  string // 可选：关闭版本

	// 基线相关
	BaselineFind  string // 可选：发现基线
	BaselineJoin  string // 可选：合入基线
	BaselineTest  string // 可选：验证基线
	BaselineClose string // 可选：关闭基线

	// 测试相关
	OS          string // 可选：操作系统
	Size        string // 可选：规模
	Platform    string // 可选：软件平台
	TestMode    string // 可选：测试方式
	TestPhase   string // 可选：测试阶段
	TestType    string // 可选：测试类型
	Frequency   string // 可选：重现规律
	OriginPhase string // 可选：发现阶段
	SourcePhase string // 可选：引入阶段
	Estimate    string // 可选：预计解决时间

	// 工时相关
	Effort string // 可选：预估工时

	// 自定义字段
	CustomFields map[string]string // 可选：自定义字段，key 如 custom_field_one、custom_field_9
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *UpdateBugRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"id":           r.ID,
	}
	setOptional(params, "title", r.Title)
	setOptional(params, "description", r.Description)
	setOptional(params, "priority", r.Priority)
	setOptional(params, "priority_label", r.PriorityLabel)
	setOptional(params, "severity", r.Severity)
	setOptional(params, "status", r.Status)
	setOptional(params, "v_status", r.VStatus)
	setOptional(params, "bugtype", r.BugType)
	setOptional(params, "source", r.Source)
	setOptional(params, "resolution", r.Resolution)
	setOptional(params, "current_owner", r.CurrentOwner)
	setOptional(params, "current_user", r.CurrentUser)
	setOptional(params, "reporter", r.Reporter)
	setOptional(params, "cc", r.CC)
	setOptional(params, "participator", r.Participator)
	setOptional(params, "te", r.TE)
	setOptional(params, "de", r.DE)
	setOptional(params, "auditer", r.Auditer)
	setOptional(params, "confirmer", r.Confirmer)
	setOptional(params, "fixer", r.Fixer)
	setOptional(params, "closer", r.Closer)
	setOptional(params, "lastmodify", r.LastModify)
	setOptional(params, "in_progress_time", r.InProgressTime)
	setOptional(params, "verify_time", r.VerifyTime)
	setOptional(params, "reject_time", r.RejectTime)
	setOptional(params, "begin", r.Begin)
	setOptional(params, "due", r.Due)
	setOptional(params, "deadline", r.Deadline)
	setOptional(params, "iteration_id", r.IterationID)
	setOptional(params, "module", r.Module)
	setOptional(params, "feature", r.Feature)
	setOptional(params, "release_id", r.ReleaseID)
	setOptional(params, "label", r.Label)
	setOptional(params, "version_report", r.VersionReport)
	setOptional(params, "version_test", r.VersionTest)
	setOptional(params, "version_fix", r.VersionFix)
	setOptional(params, "version_close", r.VersionClose)
	setOptional(params, "baseline_find", r.BaselineFind)
	setOptional(params, "baseline_join", r.BaselineJoin)
	setOptional(params, "baseline_test", r.BaselineTest)
	setOptional(params, "baseline_close", r.BaselineClose)
	setOptional(params, "os", r.OS)
	setOptional(params, "size", r.Size)
	setOptional(params, "platform", r.Platform)
	setOptional(params, "testmode", r.TestMode)
	setOptional(params, "testphase", r.TestPhase)
	setOptional(params, "testtype", r.TestType)
	setOptional(params, "frequency", r.Frequency)
	setOptional(params, "originphase", r.OriginPhase)
	setOptional(params, "sourcephase", r.SourcePhase)
	setOptional(params, "estimate", r.Estimate)
	setOptional(params, "effort", r.Effort)
	MergeCustomFields(params, r.CustomFields)
	return params
}

// CountBugsRequest 查询缺陷数量的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/bug/get_bugs_count.html
type CountBugsRequest struct {
	// 必填参数
	WorkspaceID string // 必填：项目 ID

	// 基本筛选
	ID            string // 可选：缺陷 ID，支持多 ID 查询
	Title         string // 可选：标题（支持模糊匹配）
	Priority      string // 可选：优先级（建议使用 PriorityLabel）
	PriorityLabel string // 可选：优先级（推荐）
	Severity      string // 可选：严重程度，支持枚举查询
	Status        string // 可选：状态，支持不等于查询、枚举查询
	BugType       string // 可选：缺陷类型
	Source        string // 可选：缺陷根源，支持枚举查询
	Resolution    string // 可选：解决方法，支持枚举查询
	Description   string // 可选：详细描述（支持模糊匹配）

	// 人员筛选
	CurrentOwner string // 可选：处理人（支持模糊匹配）
	Reporter     string // 可选：创建人
	CC           string // 可选：抄送人
	Participator string // 可选：参与人，支持多人员查询
	TE           string // 可选：测试人员（支持模糊匹配）
	DE           string // 可选：开发人员（支持模糊匹配）
	Auditer      string // 可选：审核人
	Confirmer    string // 可选：验证人
	Fixer        string // 可选：修复人
	Closer       string // 可选：关闭人
	LastModify   string // 可选：最后修改人

	// 时间筛选
	Created        string // 可选：创建时间，支持时间查询
	Modified       string // 可选：最后修改时间，支持时间查询
	Resolved       string // 可选：解决时间，支持时间查询
	Closed         string // 可选：关闭时间，支持时间查询
	InProgressTime string // 可选：接受处理时间，支持时间查询
	VerifyTime     string // 可选：验证时间，支持时间查询
	RejectTime     string // 可选：拒绝时间，支持时间查询
	Begin          string // 可选：预计开始
	Due            string // 可选：预计结束
	Deadline       string // 可选：解决期限

	// 关联与分类筛选
	IterationID string // 可选：迭代 ID
	Module      string // 可选：模块，支持枚举查询
	Feature     string // 可选：特性
	ReleaseID   string // 可选：发布计划
	Label       string // 可选：标签，支持枚举查询

	// 版本筛选
	VersionReport string // 可选：发现版本，支持枚举查询
	VersionTest   string // 可选：验证版本
	VersionFix    string // 可选：合入版本
	VersionClose  string // 可选：关闭版本

	// 基线筛选
	BaselineFind  string // 可选：发现基线
	BaselineJoin  string // 可选：合入基线
	BaselineTest  string // 可选：验证基线
	BaselineClose string // 可选：关闭基线

	// 测试相关筛选
	OS          string // 可选：操作系统
	Platform    string // 可选：软件平台
	TestMode    string // 可选：测试方式
	TestPhase   string // 可选：测试阶段
	TestType    string // 可选：测试类型
	Frequency   string // 可选：重现规律，支持枚举查询
	OriginPhase string // 可选：发现阶段
	SourcePhase string // 可选：引入阶段
	Estimate    string // 可选：预计解决时间

	// 工时筛选
	Effort string // 可选：预估工时
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CountBugsRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "id", r.ID)
	setOptional(params, "title", r.Title)
	setOptional(params, "priority", r.Priority)
	setOptional(params, "priority_label", r.PriorityLabel)
	setOptional(params, "severity", r.Severity)
	setOptional(params, "status", r.Status)
	setOptional(params, "bugtype", r.BugType)
	setOptional(params, "source", r.Source)
	setOptional(params, "resolution", r.Resolution)
	setOptional(params, "description", r.Description)
	setOptional(params, "current_owner", r.CurrentOwner)
	setOptional(params, "reporter", r.Reporter)
	setOptional(params, "cc", r.CC)
	setOptional(params, "participator", r.Participator)
	setOptional(params, "te", r.TE)
	setOptional(params, "de", r.DE)
	setOptional(params, "auditer", r.Auditer)
	setOptional(params, "confirmer", r.Confirmer)
	setOptional(params, "fixer", r.Fixer)
	setOptional(params, "closer", r.Closer)
	setOptional(params, "lastmodify", r.LastModify)
	setOptional(params, "created", r.Created)
	setOptional(params, "modified", r.Modified)
	setOptional(params, "resolved", r.Resolved)
	setOptional(params, "closed", r.Closed)
	setOptional(params, "in_progress_time", r.InProgressTime)
	setOptional(params, "verify_time", r.VerifyTime)
	setOptional(params, "reject_time", r.RejectTime)
	setOptional(params, "begin", r.Begin)
	setOptional(params, "due", r.Due)
	setOptional(params, "deadline", r.Deadline)
	setOptional(params, "iteration_id", r.IterationID)
	setOptional(params, "module", r.Module)
	setOptional(params, "feature", r.Feature)
	setOptional(params, "release_id", r.ReleaseID)
	setOptional(params, "label", r.Label)
	setOptional(params, "version_report", r.VersionReport)
	setOptional(params, "version_test", r.VersionTest)
	setOptional(params, "version_fix", r.VersionFix)
	setOptional(params, "version_close", r.VersionClose)
	setOptional(params, "baseline_find", r.BaselineFind)
	setOptional(params, "baseline_join", r.BaselineJoin)
	setOptional(params, "baseline_test", r.BaselineTest)
	setOptional(params, "baseline_close", r.BaselineClose)
	setOptional(params, "os", r.OS)
	setOptional(params, "platform", r.Platform)
	setOptional(params, "testmode", r.TestMode)
	setOptional(params, "testphase", r.TestPhase)
	setOptional(params, "testtype", r.TestType)
	setOptional(params, "frequency", r.Frequency)
	setOptional(params, "originphase", r.OriginPhase)
	setOptional(params, "sourcephase", r.SourcePhase)
	setOptional(params, "estimate", r.Estimate)
	setOptional(params, "effort", r.Effort)
	return params
}
