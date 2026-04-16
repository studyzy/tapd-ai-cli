// Package client 中的 attachment.go 封装了 TAPD 附件和图片相关 API
package client

import (
	"encoding/json"
	"fmt"

	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// GetImage 获取图片下载链接
func (c *Client) GetImage(params map[string]string) (*model.ImageInfo, error) {
	data, err := c.doGet("/files/get_image", params)
	if err != nil {
		return nil, err
	}

	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse image response: %w", err)
	}

	raw, ok := wrapper["Attachment"]
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var img model.ImageInfo
	if err := json.Unmarshal(raw, &img); err != nil {
		return nil, fmt.Errorf("failed to parse image info: %w", err)
	}
	return &img, nil
}

// GetAttachments 获取附件列表（含下载链接）
func (c *Client) GetAttachments(params map[string]string) ([]model.Attachment, error) {
	data, err := c.doGet("/attachments", params)
	if err != nil {
		return nil, err
	}

	var rawList []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, fmt.Errorf("failed to parse attachment list: %w", err)
	}

	var results []model.Attachment
	for _, item := range rawList {
		if raw, ok := item["Attachment"]; ok {
			var att model.Attachment
			if err := json.Unmarshal(raw, &att); err == nil {
				results = append(results, att)
			}
		}
	}
	return results, nil
}
