package tapd

import (
	"encoding/json"
	"fmt"

	"github.com/studyzy/tapd-sdk-go/model"
)

// ListIterations 查询迭代列表，返回强类型 Iteration 切片
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/iteration/get_iterations.html
func (c *Client) ListIterations(req *model.ListIterationsRequest) ([]model.Iteration, error) {
	data, err := c.doGet("/iterations", req.ToParams())
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse iteration list: %w", err)
	}

	var iterations []model.Iteration
	for _, item := range rawList {
		if raw, ok := item["Iteration"]; ok {
			var iter model.Iteration
			if err := json.Unmarshal(raw, &iter); err == nil {
				iterations = append(iterations, iter)
			}
		}
	}
	return iterations, nil
}

// CreateIteration 创建迭代，返回创建后的完整 Iteration 对象
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/iteration/add_iteration.html
func (c *Client) CreateIteration(req *model.CreateIterationRequest) (*model.Iteration, error) {
	data, err := c.doPost("/iterations", req.ToParams())
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse create iteration response: %w", err)
	}

	raw, ok := wrapper["Iteration"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var iteration model.Iteration
	if err := json.Unmarshal(raw, &iteration); err != nil {
		return nil, fmt.Errorf("failed to parse created iteration: %w", err)
	}

	return &iteration, nil
}

// UpdateIteration 更新迭代，返回更新后的完整 Iteration 对象
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/iteration/update_iteration.html
func (c *Client) UpdateIteration(req *model.UpdateIterationRequest) (*model.Iteration, error) {
	data, err := c.doPost("/iterations", req.ToParams())
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse update iteration response: %w", err)
	}

	raw, ok := wrapper["Iteration"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var iteration model.Iteration
	if err := json.Unmarshal(raw, &iteration); err != nil {
		return nil, fmt.Errorf("failed to parse updated iteration: %w", err)
	}

	return &iteration, nil
}

// CountIterations 查询迭代数量
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/iteration/get_iterations_count.html
func (c *Client) CountIterations(req *model.CountIterationsRequest) (int, error) {
	data, err := c.doGet("/iterations/count", req.ToParams())
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
