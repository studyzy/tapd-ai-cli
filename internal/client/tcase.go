// Package client 中的 tcase.go 封装了 TAPD 测试用例相关 API
package client

import (
	"encoding/json"
	"fmt"

	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// ListTCases 查询测试用例列表
func (c *Client) ListTCases(params map[string]string) ([]model.TCase, error) {
	data, err := c.doGet("/tcases", params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse tcase list: %w", err)
	}

	var results []model.TCase
	for _, item := range rawList {
		if raw, ok := item["Tcase"]; ok {
			var tc model.TCase
			if err := json.Unmarshal(raw, &tc); err == nil {
				results = append(results, tc)
			}
		}
	}
	return results, nil
}

// CountTCases 查询测试用例数量
func (c *Client) CountTCases(params map[string]string) (int, error) {
	data, err := c.doGet("/tcases/count", params)
	if err != nil {
		return 0, err
	}

	var result map[string]int
	if err := json.Unmarshal(data, &result); err != nil {
		return 0, fmt.Errorf("failed to parse tcase count: %w", err)
	}

	if count, ok := result["count"]; ok {
		return count, nil
	}
	return 0, nil
}

// CreateTCase 创建或更新测试用例
func (c *Client) CreateTCase(params map[string]string) (*model.TCase, error) {
	data, err := c.doPost("/tcases", params)
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse create tcase response: %w", err)
	}

	raw, ok := wrapper["Tcase"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var tc model.TCase
	if err := json.Unmarshal(raw, &tc); err != nil {
		return nil, fmt.Errorf("failed to parse created tcase: %w", err)
	}
	return &tc, nil
}

// BatchCreateTCases 批量创建测试用例
func (c *Client) BatchCreateTCases(params map[string]string) (json.RawMessage, error) {
	return c.doPost("/tcases/batch_save", params)
}
