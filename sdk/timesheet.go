package tapd

import (
	"encoding/json"
	"fmt"

	"github.com/studyzy/tapd-sdk-go/model"
)

// ListTimesheets 查询花费工时列表
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/timesheet/get_timesheets.html
func (c *Client) ListTimesheets(req *model.ListTimesheetsRequest) ([]model.Timesheet, error) {
	data, err := c.doGet("/timesheets", req.ToParams())
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse timesheet list: %w", err)
	}

	var results []model.Timesheet
	for _, item := range rawList {
		if raw, ok := item["Timesheet"]; ok {
			var ts model.Timesheet
			if err := json.Unmarshal(raw, &ts); err == nil {
				results = append(results, ts)
			}
		}
	}
	return results, nil
}

// AddTimesheet 填写花费工时
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/timesheet/add_timesheet.html
func (c *Client) AddTimesheet(req *model.AddTimesheetRequest) (*model.Timesheet, error) {
	data, err := c.doPost("/timesheets", req.ToParams())
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse add timesheet response: %w", err)
	}

	raw, ok := wrapper["Timesheet"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var ts model.Timesheet
	if err := json.Unmarshal(raw, &ts); err != nil {
		return nil, fmt.Errorf("failed to parse created timesheet: %w", err)
	}
	return &ts, nil
}

// UpdateTimesheet 更新花费工时
// API 文档：https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/timesheet/update_timesheet.html
func (c *Client) UpdateTimesheet(req *model.UpdateTimesheetRequest) (*model.Timesheet, error) {
	data, err := c.doPost("/timesheets", req.ToParams())
	if err != nil {
		return nil, err
	}

	// 尝试直接解析
	var ts model.Timesheet
	if err := json.Unmarshal(data, &ts); err == nil && ts.ID != "" {
		return &ts, nil
	}

	// 兼容 "Timesheet" 包裹
	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse update timesheet response: %w", err)
	}

	raw, ok := wrapper["Timesheet"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	if err := json.Unmarshal(raw, &ts); err != nil {
		return nil, fmt.Errorf("failed to parse updated timesheet: %w", err)
	}
	return &ts, nil
}
