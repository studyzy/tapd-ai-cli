package client

import (
	"encoding/json"
	"fmt"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// ListStories 查询需求或任务列表，通过 params 中的 entity_type 区分 stories/tasks
func (c *Client) ListStories(params map[string]string) ([]map[string]interface{}, error) {
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

	var results []map[string]interface{}
	for _, item := range rawList {
		if raw, ok := item[wrapperKey]; ok {
			var obj map[string]interface{}
			if err := json.Unmarshal(raw, &obj); err == nil {
				results = append(results, obj)
			}
		}
	}
	return results, nil
}

// GetStory 获取单个需求或任务的详情，description 字段自动从 HTML 转换为 Markdown
func (c *Client) GetStory(workspaceID, id, entityType string) (map[string]interface{}, error) {
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

	var result map[string]interface{}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("failed to parse entity: %w", err)
	}

	// HTML 转 Markdown
	if desc, ok := result["description"].(string); ok && desc != "" {
		md, err := htmltomarkdown.ConvertString(desc)
		if err == nil {
			result["description"] = md
		}
	}

	// 拼接 URL
	if entityType == "tasks" {
		result["url"] = fmt.Sprintf("https://www.tapd.cn/%s/prong/tasks/view/%s", workspaceID, id)
	} else {
		result["url"] = fmt.Sprintf("https://www.tapd.cn/%s/prong/stories/view/%s", workspaceID, id)
	}

	return result, nil
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

	var created map[string]interface{}
	if err := json.Unmarshal(raw, &created); err != nil {
		return nil, fmt.Errorf("failed to parse created entity: %w", err)
	}

	id := fmt.Sprintf("%v", created["id"])
	wsID := params["workspace_id"]

	return &model.SuccessResponse{
		Success: true,
		ID:      id,
		URL:     fmt.Sprintf("https://www.tapd.cn/%s/%s/%s", wsID, urlPath, id),
	}, nil
}

// UpdateStory 更新需求或任务
func (c *Client) UpdateStory(params map[string]string, entityType string) (map[string]interface{}, error) {
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

	var result map[string]interface{}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("failed to parse updated entity: %w", err)
	}

	return result, nil
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
