package client

import (
	"encoding/json"
	"fmt"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// ListTasks 查询任务列表，返回强类型 Task 切片，自动过滤 custom_field 等无用字段
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

// GetTask 获取单个任务详情，description 字段自动从 HTML 转换为 Markdown
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

	// HTML 转 Markdown
	if task.Description != "" {
		md, err := htmltomarkdown.ConvertString(task.Description)
		if err == nil {
			task.Description = md
		}
	}

	task.URL = fmt.Sprintf("https://www.tapd.cn/%s/prong/tasks/view/%s", workspaceID, id)

	return &task, nil
}

// CreateTask 创建任务
func (c *Client) CreateTask(req *model.CreateTaskRequest) (*model.SuccessResponse, error) {
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

	var created model.Task
	if err := json.Unmarshal(raw, &created); err != nil {
		return nil, fmt.Errorf("failed to parse created task: %w", err)
	}

	return &model.SuccessResponse{
		Success: true,
		ID:      created.ID,
		URL:     fmt.Sprintf("https://www.tapd.cn/%s/prong/tasks/view/%s", req.WorkspaceID, created.ID),
	}, nil
}

// UpdateTask 更新任务
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
