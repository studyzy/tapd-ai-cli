package tapd

import (
	"encoding/json"
	"fmt"

	"github.com/studyzy/tapd-sdk-go/model"
)

// ListStories 查询需求列表，返回强类型 Story 切片，自动过滤 custom_field 等无用字段
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/story/get_stories.html
func (c *Client) ListStories(req *model.ListStoriesRequest) ([]model.Story, error) {
	data, err := c.doGet("/stories", req.ToParams())
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse story list: %w", err)
	}

	var results []model.Story
	for _, item := range rawList {
		if raw, ok := item["Story"]; ok {
			var story model.Story
			if err := json.Unmarshal(raw, &story); err == nil {
				results = append(results, story)
			}
		}
	}
	return results, nil
}

// GetStory 获取单个需求详情，description 字段保留原始 HTML
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/story/get_stories.html
func (c *Client) GetStory(workspaceID, id string) (*model.Story, error) {
	params := map[string]string{
		"workspace_id": workspaceID,
		"id":           id,
	}

	data, err := c.doGet("/stories", params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(rawList) == 0 {
		return nil, &TAPDError{ExitCode: 2, Message: fmt.Sprintf("story %s not found", id)}
	}

	raw, ok := rawList[0]["Story"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var story model.Story
	if err := json.Unmarshal(raw, &story); err != nil {
		return nil, fmt.Errorf("failed to parse story: %w", err)
	}

	story.URL = fmt.Sprintf("%s/%s/prong/stories/view/%s", c.webURL, workspaceID, id)

	return &story, nil
}

// CreateStory 创建需求，返回创建后的完整 Story 对象
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/story/add_story.html
func (c *Client) CreateStory(req *model.CreateStoryRequest) (*model.Story, error) {
	data, err := c.doPost("/stories", req.ToParams())
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse create response: %w", err)
	}

	raw, ok := wrapper["Story"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var story model.Story
	if err := json.Unmarshal(raw, &story); err != nil {
		return nil, fmt.Errorf("failed to parse created story: %w", err)
	}

	story.URL = fmt.Sprintf("%s/%s/prong/stories/view/%s", c.webURL, req.WorkspaceID, story.ID)

	return &story, nil
}

// UpdateStory 更新需求
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/story/update_story.html
func (c *Client) UpdateStory(req *model.UpdateStoryRequest) (*model.Story, error) {
	data, err := c.doPost("/stories", req.ToParams())
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse update response: %w", err)
	}

	raw, ok := wrapper["Story"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var story model.Story
	if err := json.Unmarshal(raw, &story); err != nil {
		return nil, fmt.Errorf("failed to parse updated story: %w", err)
	}

	return &story, nil
}

// CountStories 查询需求数量
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/story/get_stories_count.html
func (c *Client) CountStories(req *model.CountStoriesRequest) (int, error) {
	data, err := c.doGet("/stories/count", req.ToParams())
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
