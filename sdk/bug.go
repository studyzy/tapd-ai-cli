package tapd

import (
	"encoding/json"
	"fmt"

	"github.com/studyzy/tapd-sdk-go/model"
)

// ListBugs 查询缺陷列表，返回强类型 Bug 切片，自动过滤 custom_field 等无用字段
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/bug/get_bugs.html
func (c *Client) ListBugs(req *model.ListBugsRequest) ([]model.Bug, error) {
	data, err := c.doGet("/bugs", req.ToParams())
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse bug list: %w", err)
	}

	var results []model.Bug
	for _, item := range rawList {
		if raw, ok := item["Bug"]; ok {
			var bug model.Bug
			if err := json.Unmarshal(raw, &bug); err == nil {
				results = append(results, bug)
			}
		}
	}
	return results, nil
}

// GetBug 获取单个缺陷详情，description 字段保留原始 HTML
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/bug/get_bugs.html
func (c *Client) GetBug(workspaceID, id string) (*model.Bug, error) {
	params := map[string]string{
		"workspace_id": workspaceID,
		"id":           id,
	}

	data, err := c.doGet("/bugs", params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(rawList) == 0 {
		return nil, &TAPDError{ExitCode: 2, Message: fmt.Sprintf("bug %s not found", id)}
	}

	raw, ok := rawList[0]["Bug"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var bug model.Bug
	if err := json.Unmarshal(raw, &bug); err != nil {
		return nil, fmt.Errorf("failed to parse bug: %w", err)
	}

	bug.URL = fmt.Sprintf("%s/%s/bugtrace/bugs/view/%s", c.webURL, workspaceID, id)

	return &bug, nil
}

// CreateBug 创建缺陷，返回创建后的完整 Bug 对象
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/bug/add_bug.html
func (c *Client) CreateBug(req *model.CreateBugRequest) (*model.Bug, error) {
	data, err := c.doPost("/bugs", req.ToParams())
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse create response: %w", err)
	}

	raw, ok := wrapper["Bug"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var bug model.Bug
	if err := json.Unmarshal(raw, &bug); err != nil {
		return nil, fmt.Errorf("failed to parse created bug: %w", err)
	}

	bug.URL = fmt.Sprintf("%s/%s/bugtrace/bugs/view/%s", c.webURL, req.WorkspaceID, bug.ID)

	return &bug, nil
}

// UpdateBug 更新缺陷
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/bug/update_bug.html
func (c *Client) UpdateBug(req *model.UpdateBugRequest) (*model.Bug, error) {
	data, err := c.doPost("/bugs", req.ToParams())
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse update response: %w", err)
	}

	raw, ok := wrapper["Bug"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var bug model.Bug
	if err := json.Unmarshal(raw, &bug); err != nil {
		return nil, fmt.Errorf("failed to parse updated bug: %w", err)
	}

	return &bug, nil
}

// CountBugs 查询缺陷数量
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/bug/get_bugs_count.html
func (c *Client) CountBugs(req *model.CountBugsRequest) (int, error) {
	data, err := c.doGet("/bugs/count", req.ToParams())
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
