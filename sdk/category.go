package tapd

import (
	"encoding/json"
	"fmt"

	"github.com/studyzy/tapd-sdk-go/model"
)

// ListCategories 查询需求分类列表
func (c *Client) ListCategories(params map[string]string) ([]model.Category, error) {
	data, err := c.doGet("/story_categories", params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse category list: %w", err)
	}

	var results []model.Category
	for _, item := range rawList {
		if raw, ok := item["Category"]; ok {
			var cat model.Category
			if err := json.Unmarshal(raw, &cat); err == nil {
				results = append(results, cat)
			}
		}
	}
	return results, nil
}
