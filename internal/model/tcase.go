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
	WorkspaceID string // 必填：项目 ID
	ID          string // 可选：测试用例 ID
	Name        string // 可选：标题
	Status      string // 可选：状态
	CategoryID  string // 可选：目录 ID
	Priority    string // 可选：优先级
	Fields      string // 可选：返回字段列表
	Limit       string // 可选：返回数量限制
	Page        string // 可选：页码
	Order       string // 可选：排序规则
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
	setOptional(params, "fields", r.Fields)
	setOptional(params, "limit", r.Limit)
	setOptional(params, "page", r.Page)
	setOptional(params, "order", r.Order)
	return params
}

// CountTCasesRequest 查询测试用例数量的请求参数
type CountTCasesRequest struct {
	WorkspaceID string // 必填：项目 ID
	Status      string // 可选：状态
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CountTCasesRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
	}
	setOptional(params, "status", r.Status)
	return params
}

// CreateTCaseRequest 创建测试用例的请求参数
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/tcase/add_tcase.html
type CreateTCaseRequest struct {
	WorkspaceID  string // 必填：项目 ID
	Name         string // 必填：用例标题
	CategoryID   string // 可选：目录 ID
	Precondition string // 可选：前置条件
	Steps        string // 可选：用例步骤
	Expectation  string // 可选：预期结果
	Type         string // 可选：类型
	Priority     string // 可选：优先级
	Creator      string // 可选：创建人
}

// ToParams 将请求结构体转换为 TAPD API 参数 map
func (r *CreateTCaseRequest) ToParams() map[string]string {
	params := map[string]string{
		"workspace_id": r.WorkspaceID,
		"name":         r.Name,
	}
	setOptional(params, "category_id", r.CategoryID)
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
