package client

import (
	"encoding/json"
	"fmt"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// ListStories 查询需求或任务列表，通过 params 中的 entity_type 区分 stories/tasks
// 返回强类型切片（[]model.Story 或 []model.Task），自动过滤 custom_field 等无用字段
func (c *Client) ListStories(params map[string]string) (interface{}, error) {
	entityType := params["entity_type"]
	delete(params, "entity_type")

	endpoint := "/stories"
	wrapperKey := "Story"
	if entityType == "tasks" {
		endpoint = "/tasks"
		wrapperKey = "Task"
	}

	data, err := c.doGet(endpoint, params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse list response: %w", err)
	}

	if entityType == "tasks" {
		var results []model.Task
		for _, item := range rawList {
			if raw, ok := item[wrapperKey]; ok {
				var task model.Task
				if err := json.Unmarshal(raw, &task); err == nil {
					results = append(results, task)
				}
			}
		}
		return results, nil
	}

	var results []model.Story
	for _, item := range rawList {
		if raw, ok := item[wrapperKey]; ok {
			var story model.Story
			if err := json.Unmarshal(raw, &story); err == nil {
				results = append(results, story)
			}
		}
	}
	return results, nil
}

// GetStory 获取单个需求或任务的详情，description 字段自动从 HTML 转换为 Markdown
// 返回强类型（*model.Story 或 *model.Task），自动过滤 custom_field 等无用字段
func (c *Client) GetStory(workspaceID, id, entityType string) (interface{}, error) {
	endpoint := "/stories"
	wrapperKey := "Story"
	if entityType == "tasks" {
		endpoint = "/tasks"
		wrapperKey = "Task"
	}

	params := map[string]string{
		"workspace_id": workspaceID,
		"id":           id,
	}

	data, err := c.doGet(endpoint, params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(rawList) == 0 {
		return nil, &TAPDError{ExitCode: 2, Message: fmt.Sprintf("%s %s not found", entityType, id)}
	}

	raw, ok := rawList[0][wrapperKey]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	if entityType == "tasks" {
		var task model.Task
		if err := json.Unmarshal(raw, &task); err != nil {
			return nil, fmt.Errorf("failed to parse task: %w", err)
		}
		if task.Description != "" {
			md, err := htmltomarkdown.ConvertString(task.Description)
			if err == nil {
				task.Description = md
			}
		}
		task.URL = fmt.Sprintf("https://www.tapd.cn/%s/prong/tasks/view/%s", workspaceID, id)
		return &task, nil
	}

	var story model.Story
	if err := json.Unmarshal(raw, &story); err != nil {
		return nil, fmt.Errorf("failed to parse story: %w", err)
	}
	if story.Description != "" {
		md, err := htmltomarkdown.ConvertString(story.Description)
		if err == nil {
			story.Description = md
		}
	}
	story.URL = fmt.Sprintf("https://www.tapd.cn/%s/prong/stories/view/%s", workspaceID, id)
	return &story, nil
}

// CreateStory 创建需求或任务
func (c *Client) CreateStory(params map[string]string, entityType string) (*model.SuccessResponse, error) {
	endpoint := "/stories"
	wrapperKey := "Story"
	urlPath := "prong/stories/view"
	if entityType == "tasks" {
		endpoint = "/tasks"
		wrapperKey = "Task"
		urlPath = "prong/tasks/view"
	}

	data, err := c.doPost(endpoint, params)
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse create response: %w", err)
	}

	raw, ok := wrapper[wrapperKey]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(raw, &created); err != nil {
		return nil, fmt.Errorf("failed to parse created entity: %w", err)
	}

	wsID := params["workspace_id"]

	return &model.SuccessResponse{
		Success: true,
		ID:      created.ID,
		URL:     fmt.Sprintf("https://www.tapd.cn/%s/%s/%s", wsID, urlPath, created.ID),
	}, nil
}

// UpdateStory 更新需求或任务，返回强类型（*model.Story 或 *model.Task）
func (c *Client) UpdateStory(params map[string]string, entityType string) (interface{}, error) {
	endpoint := "/stories"
	wrapperKey := "Story"
	if entityType == "tasks" {
		endpoint = "/tasks"
		wrapperKey = "Task"
	}

	data, err := c.doPost(endpoint, params)
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse update response: %w", err)
	}

	raw, ok := wrapper[wrapperKey]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	if entityType == "tasks" {
		var task model.Task
		if err := json.Unmarshal(raw, &task); err != nil {
			return nil, fmt.Errorf("failed to parse updated task: %w", err)
		}
		return &task, nil
	}

	var story model.Story
	if err := json.Unmarshal(raw, &story); err != nil {
		return nil, fmt.Errorf("failed to parse updated story: %w", err)
	}
	return &story, nil
}

// CountStories 查询需求或任务数量
func (c *Client) CountStories(params map[string]string) (int, error) {
	entityType := params["entity_type"]
	delete(params, "entity_type")

	endpoint := "/stories/count"
	if entityType == "tasks" {
		endpoint = "/tasks/count"
	}

	data, err := c.doGet(endpoint, params)
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
