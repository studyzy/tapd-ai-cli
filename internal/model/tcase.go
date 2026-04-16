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
