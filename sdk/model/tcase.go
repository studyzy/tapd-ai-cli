// Package model 中的 tcase.go 定义了 TAPD 测试用例数据模型
package model

// TCase 表示 TAPD 测试用例
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/tcase/get_tcases.html
type TCase struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	WorkspaceID  string `json:"workspace_id,omitempty"`
	CategoryID   string `json:"category_id,omitempty"`
	Status       string `json:"status,omitempty"`
	Precondition string `json:"precondition,omitempty"`
	Steps        string `json:"steps,omitempty"`
	Expectation  string `json:"expectation,omitempty"`
	Type         string `json:"type,omitempty"`
	Priority     string `json:"priority,omitempty"`
	Creator      string `json:"creator,omitempty"`
	Modifier     string `json:"modifier,omitempty"`
	Created      string `json:"created,omitempty"`
	Modified     string `json:"modified,omitempty"`
	URL          string `json:"url,omitempty"`
}

// ListTCasesRequest 查询测试用例列表的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/tcase/get_tcases.html
type ListTCasesRequest struct {
	WorkspaceID  string // 必填：项目 ID
	ID           string // 可选：测试用例 ID，支持多 ID 查询
	Name         string // 可选：用例名称，支持模糊匹配
	Status       string // 可选：用例状态（normal|updating|abandon）
	CategoryID   string // 可选：用例目录 ID
	Priority     string // 可选：用例等级
	Type         string // 可选：用例类型
	Creator      string // 可选：创建人
	Modifier     string // 可选：最后修改人
	Created      string // 可选：创建时间，支持时间查询
	Modified     string // 可选：最后修改时间，支持时间查询
	Steps        string // 可选：用例步骤
	Precondition string // 可选：前置条件
	Expectation  string // 可选：预期结果
	Fields       string // 可选：返回字段列表，多个字段间以半角逗号隔开
	Limit        string // 可选：返回数量限制，默认为 30
	Page         string // 可选：页码，默认为 1
	Order        string // 可选：排序规则，如 created desc
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *ListTCasesRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "id", r.ID)
	setOptional(params, "name", r.Name)
	setOptional(params, "status", r.Status)
	setOptional(params, "category_id", r.CategoryID)
	setOptional(params, "priority", r.Priority)
	setOptional(params, "type", r.Type)
	setOptional(params, "creator", r.Creator)
	setOptional(params, "modifier", r.Modifier)
	setOptional(params, "created", r.Created)
	setOptional(params, "modified", r.Modified)
	setOptional(params, "steps", r.Steps)
	setOptional(params, "precondition", r.Precondition)
	setOptional(params, "expectation", r.Expectation)
	setOptional(params, "fields", r.Fields)
	setOptional(params, "limit", r.Limit)
	setOptional(params, "page", r.Page)
	setOptional(params, "order", r.Order)
	return params
}

// CountTCasesRequest 查询测试用例数量的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/tcase/get_tcases_count.html
type CountTCasesRequest struct {
	WorkspaceID        string // 必填：项目 ID
	ID                 string // 可选：测试用例 ID，支持多 ID 查询
	Name               string // 可选：用例名称，支持模糊匹配
	Status             string // 可选：用例状态（normal|updating|abandon）
	CategoryID         string // 可选：用例目录 ID
	Priority           string // 可选：用例等级
	Type               string // 可选：用例类型
	Creator            string // 可选：创建人
	Modifier           string // 可选：最后修改人
	Created            string // 可选：创建时间，支持时间查询
	Modified           string // 可选：最后修改时间，支持时间查询
	Steps              string // 可选：用例步骤
	Precondition       string // 可选：前置条件
	Expectation        string // 可选：预期结果
	TestPlanID         string // 可选：测试计划 ID，获取当前测试计划关联的测试用例数量
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CountTCasesRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "id", r.ID)
	setOptional(params, "name", r.Name)
	setOptional(params, "status", r.Status)
	setOptional(params, "category_id", r.CategoryID)
	setOptional(params, "priority", r.Priority)
	setOptional(params, "type", r.Type)
	setOptional(params, "creator", r.Creator)
	setOptional(params, "modifier", r.Modifier)
	setOptional(params, "created", r.Created)
	setOptional(params, "modified", r.Modified)
	setOptional(params, "steps", r.Steps)
	setOptional(params, "precondition", r.Precondition)
	setOptional(params, "expectation", r.Expectation)
	setOptional(params, "test_plan_id", r.TestPlanID)
	return params
}

// CreateTCaseRequest 创建测试用例的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/tcase/add_tcase.html
type CreateTCaseRequest struct {
	WorkspaceID  string // 必填：项目 ID
	Name         string // 必填：用例名称
	ID           string // 可选：测试用例 ID
	CategoryID   string // 可选：用例目录 ID
	Status       string // 可选：用例状态（normal|updating|abandon）
	Precondition string // 可选：前置条件
	Steps        string // 可选：用例步骤
	Expectation  string // 可选：预期结果
	Type         string // 可选：用例类型
	Priority     string // 可选：用例等级
	Creator      string // 可选：创建人
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CreateTCaseRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"name":         r.Name,
	}
	setOptional(params, "id", r.ID)
	setOptional(params, "category_id", r.CategoryID)
	setOptional(params, "status", r.Status)
	setOptional(params, "precondition", r.Precondition)
	setOptional(params, "steps", r.Steps)
	setOptional(params, "expectation", r.Expectation)
	setOptional(params, "type", r.Type)
	setOptional(params, "priority", r.Priority)
	setOptional(params, "creator", r.Creator)
	return params
}

// BatchCreateTCasesRequest 批量创建测试用例的请求参数
type BatchCreateTCasesRequest struct {
	WorkspaceID string // 必填：项目 ID
	Data        string // 必填：批量数据
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *BatchCreateTCasesRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "data", r.Data)
	return params
}
