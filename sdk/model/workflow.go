// Package model 中的 workflow.go 定义了 TAPD 工作流相关数据模型
package model

// WorkflowTransition 表示工作流状态流转规则
type WorkflowTransition struct {
	StatusFrom string `json:"status_from,omitempty"`
	StatusTo   string `json:"status_to,omitempty"`
}

// WorkflowStatusMap 表示工作流状态的中英文映射
type WorkflowStatusMap struct {
	EnglishName string `json:"english_name,omitempty"`
	ChineseName string `json:"chinese_name,omitempty"`
}
