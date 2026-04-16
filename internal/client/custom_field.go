// Package client 中的 custom_field.go 封装了 TAPD 自定义字段和需求字段相关 API
package client

import "encoding/json"

// GetCustomFields 获取自定义字段配置
// entityType 取值：stories, tasks, iterations, tcases
func (c *Client) GetCustomFields(params map[string]string) (json.RawMessage, error) {
	entityType := params["entity_type"]
	delete(params, "entity_type")
	return c.doGet("/"+entityType+"/custom_fields_settings", params)
}

// GetStoryFieldsLabel 获取需求所有字段的中英文名
func (c *Client) GetStoryFieldsLabel(params map[string]string) (json.RawMessage, error) {
	return c.doGet("/stories/get_fields_lable", params)
}

// GetStoryFieldsInfo 获取需求所有字段及候选值
func (c *Client) GetStoryFieldsInfo(params map[string]string) (json.RawMessage, error) {
	return c.doGet("/stories/get_fields_info", params)
}

// GetWorkitemTypes 获取需求类别列表
func (c *Client) GetWorkitemTypes(params map[string]string) (json.RawMessage, error) {
	return c.doGet("/workitem_types", params)
}
