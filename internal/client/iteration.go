package client

import (
	"encoding/json"
	"fmt"

	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// ListIterations 查询迭代列表
func (c *Client) ListIterations(params map[string]string) ([]model.Iteration, error) {
	data, err := c.doGet("/iterations", params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]model.Iteration
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse iteration list: %w", err)
	}

	var iterations []model.Iteration
	for _, item := range rawList {
		if iter, ok := item["Iteration"]; ok {
			iterations = append(iterations, iter)
		}
	}
	return iterations, nil
}
