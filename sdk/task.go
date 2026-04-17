package tapd

import (
	"encoding/json"
	"fmt"

	"github.com/studyzy/tapd-sdk-go/model"
)

// ListTasks 查询任务列表，返回强类型 Task 切片，自动过滤 custom_field 等无用字段
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/task/get_tasks.html
func (c *Client) ListTasks(req *model.ListTasksRequest) ([]model.Task, error) {
	data, err := c.doGet("/tasks", req.ToParams())
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse task list: %w", err)
	}

	var results []model.Task
	for _, item := range rawList {
		if raw, ok := item["Task"]; ok {
			var task model.Task
			if err := json.Unmarshal(raw, &task); err == nil {
				results = append(results, task)
			}
		}
	}
	return results, nil
}

// GetTask 获取单个任务详情，description 字段保留原始 HTML
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/task/get_tasks.html
func (c *Client) GetTask(workspaceID, id string) (*model.Task, error) {
	params := map[string]string{
		"workspace_id": workspaceID,
		"id":           id,
	}

	data, err := c.doGet("/tasks", params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(rawList) == 0 {
		return nil, &TAPDError{ExitCode: 2, Message: fmt.Sprintf("task %s not found", id)}
	}

	raw, ok := rawList[0]["Task"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var task model.Task
	if err := json.Unmarshal(raw, &task); err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)
	}

	task.URL = fmt.Sprintf("%s/%s/prong/tasks/view/%s", c.webURL, workspaceID, id)

	return &task, nil
}

// CreateTask 创建任务，返回创建后的完整 Task 对象
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/task/add_task.html
func (c *Client) CreateTask(req *model.CreateTaskRequest) (*model.Task, error) {
	data, err := c.doPost("/tasks", req.ToParams())
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse create response: %w", err)
	}

	raw, ok := wrapper["Task"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var task model.Task
	if err := json.Unmarshal(raw, &task); err != nil {
		return nil, fmt.Errorf("failed to parse created task: %w", err)
	}

	task.URL = fmt.Sprintf("%s/%s/prong/tasks/view/%s", c.webURL, req.WorkspaceID, task.ID)

	return &task, nil
}

// UpdateTask 更新任务
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/task/update_task.html
func (c *Client) UpdateTask(req *model.UpdateTaskRequest) (*model.Task, error) {
	data, err := c.doPost("/tasks", req.ToParams())
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse update response: %w", err)
	}

	raw, ok := wrapper["Task"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var task model.Task
	if err := json.Unmarshal(raw, &task); err != nil {
		return nil, fmt.Errorf("failed to parse updated task: %w", err)
	}

	return &task, nil
}

// CountTasks 查询任务数量
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/task/get_tasks_count.html
func (c *Client) CountTasks(req *model.CountTasksRequest) (int, error) {
	data, err := c.doGet("/tasks/count", req.ToParams())
	if err != nil {
		return 0, err
	}

	var result map[string]int
	if err := json.Unmarshal(data, &result); err != nil {
		return 0, fmt.Errorf("failed to parse count response: %w", err)
	}

	if count, ok := result["count"]; ok {
		return count, nil
	}
	return 0, nil
}
