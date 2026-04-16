package client

import (
	"encoding/json"
	"fmt"

	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// ListIterations 查询迭代列表
func (c *Client) ListIterations(req *model.ListIterationsRequest) ([]model.Iteration, error) {
	data, err := c.doGet("/iterations", req.ToParams())
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

// CreateIteration 创建迭代
func (c *Client) CreateIteration(req *model.CreateIterationRequest) (*model.SuccessResponse, error) {
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

	var created model.Iteration
	if err := json.Unmarshal(raw, &created); err != nil {
		return nil, fmt.Errorf("failed to parse created iteration: %w", err)
	}

	wsID := req.WorkspaceID

	return &model.SuccessResponse{
		Success:     true,
		ID:          created.ID,
		WorkspaceID: wsID,
	}, nil
}

// UpdateIteration 更新迭代
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
