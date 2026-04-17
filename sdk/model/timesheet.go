// Package model 中的 timesheet.go 定义了 TAPD 花费工时数据模型
package model

// Timesheet 表示 TAPD 花费工时记录
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/timesheet/
type Timesheet struct {
	ID          string `json:"id,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty"`
	EntityType  string `json:"entity_type,omitempty"`
	EntityID    string `json:"entity_id,omitempty"`
	Timespent   string `json:"timespent,omitempty"`
	Timeremain  string `json:"timeremain,omitempty"`
	Spentdate   string `json:"spentdate,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Memo        string `json:"memo,omitempty"`
	Created     string `json:"created,omitempty"`
	Modified    string `json:"modified,omitempty"`
	IsDelete    string `json:"is_delete,omitempty"`
}

// ListTimesheetsRequest 查询工时列表的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/timesheet/get_timesheets.html
type ListTimesheetsRequest struct {
	WorkspaceID                  string // 必填：项目 ID
	ID                           string // 可选：工时记录 ID，支持多 ID 查询
	EntityType                   string // 可选：对象类型（story/task/bug）
	EntityID                     string // 可选：对象 ID
	Timespent                    string // 可选：花费工时
	Spentdate                    string // 可选：花费日期，支持时间查询
	Modified                     string // 可选：最后修改时间，支持时间查询
	Owner                        string // 可选：花费创建人
	IncludeParentStoryTimesheet  string // 可选：值为 0 时不返回父需求的花费
	Created                      string // 可选：创建时间，支持时间查询
	Memo                         string // 可选：花费描述
	IsDelete                     string // 可选：是否已删除，默认 0 不返回已删除记录，1 返回已删除记录
	Fields                       string // 可选：返回字段列表，多个字段间以逗号隔开
	Limit                        string // 可选：返回数量限制，默认 30
	Page                         string // 可选：页码，默认 1
	Order                        string // 可选：排序规则，如 created desc
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *ListTimesheetsRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "id", r.ID)
	setOptional(params, "entity_type", r.EntityType)
	setOptional(params, "entity_id", r.EntityID)
	setOptional(params, "timespent", r.Timespent)
	setOptional(params, "spentdate", r.Spentdate)
	setOptional(params, "modified", r.Modified)
	setOptional(params, "owner", r.Owner)
	setOptional(params, "include_parent_story_timesheet", r.IncludeParentStoryTimesheet)
	setOptional(params, "created", r.Created)
	setOptional(params, "memo", r.Memo)
	setOptional(params, "is_delete", r.IsDelete)
	setOptional(params, "fields", r.Fields)
	setOptional(params, "limit", r.Limit)
	setOptional(params, "page", r.Page)
	setOptional(params, "order", r.Order)
	return params
}

// AddTimesheetRequest 填写工时的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/timesheet/add_timesheet.html
type AddTimesheetRequest struct {
	WorkspaceID string // 必填：项目 ID
	EntityType  string // 必填：对象类型（story/task/bug）
	EntityID    string // 必填：对象 ID
	Timespent   string // 必填：花费工时
	Owner       string // 必填：花费创建人
	Timeremain  string // 可选：剩余工时
	Spentdate   string // 可选：花费日期
	Memo        string // 可选：花费描述
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *AddTimesheetRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"entity_type":  r.EntityType,
		"entity_id":    r.EntityID,
		"timespent":    r.Timespent,
		"owner":        r.Owner,
	}
	setOptional(params, "timeremain", r.Timeremain)
	setOptional(params, "spentdate", r.Spentdate)
	setOptional(params, "memo", r.Memo)
	return params
}

// UpdateTimesheetRequest 更新工时的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/timesheet/update_timesheet.html
type UpdateTimesheetRequest struct {
	WorkspaceID string // 必填：项目 ID
	ID          string // 必填：工时记录 ID
	Timespent   string // 可选：花费工时
	Timeremain  string // 可选：剩余工时
	Memo        string // 可选：花费描述
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *UpdateTimesheetRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"id":           r.ID,
	}
	setOptional(params, "timespent", r.Timespent)
	setOptional(params, "timeremain", r.Timeremain)
	setOptional(params, "memo", r.Memo)
	return params
}
