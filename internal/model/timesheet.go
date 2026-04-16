// Package model 中的 timesheet.go 定义了 TAPD 花费工时数据模型
package model

// Timesheet 表示 TAPD 花费工时记录
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/timesheet/
type Timesheet struct {
	ID         string `json:"id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	EntityID   string `json:"entity_id,omitempty"`
	Timespent  string `json:"timespent,omitempty"`
	Timeremain string `json:"timeremain,omitempty"`
	Spentdate  string `json:"spentdate,omitempty"`
	Owner      string `json:"owner,omitempty"`
	Memo       string `json:"memo,omitempty"`
	Created    string `json:"created,omitempty"`
	Modified   string `json:"modified,omitempty"`
	IsDelete   string `json:"is_delete,omitempty"`
}
