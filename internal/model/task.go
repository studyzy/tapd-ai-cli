// Package model 中的 task.go 定义了 TAPD 任务数据模型
package model

// Task 表示 TAPD 任务，字段覆盖 TAPD API 返回的所有常用字段
// 使用强类型结构体反序列化可自动过滤 custom_field_* 等无用字段，节约 token
// 参考：https://open.tapd.cn/document/api-doc/API文档/api_reference/task/get_tasks.html
type Task struct {
	// 基本信息
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty"`
	Status      string `json:"status,omitempty"`
	CreatedFrom string `json:"created_from,omitempty"`

	// 优先级
	Priority      string `json:"priority,omitempty"`
	PriorityLabel string `json:"priority_label,omitempty"`

	// 人员相关
	Owner   string `json:"owner,omitempty"`
	Creator string `json:"creator,omitempty"`
	CC      string `json:"cc,omitempty"`

	// 时间相关
	Created   string `json:"created,omitempty"`
	Modified  string `json:"modified,omitempty"`
	Completed string `json:"completed,omitempty"`
	Begin     string `json:"begin,omitempty"`
	Due       string `json:"due,omitempty"`

	// 关联与分类
	StoryID     string `json:"story_id,omitempty"`
	IterationID string `json:"iteration_id,omitempty"`
	ReleaseID   string `json:"release_id,omitempty"`
	Label       string `json:"label,omitempty"`

	// 工时与进度
	Effort          string `json:"effort,omitempty"`
	EffortCompleted string `json:"effort_completed,omitempty"`
	Remain          string `json:"remain,omitempty"`
	Exceed          string `json:"exceed,omitempty"`
	Progress        string `json:"progress,omitempty"`
	HasAttachment   string `json:"has_attachment,omitempty"`

	// 附加信息
	URL string `json:"url,omitempty"`
}
