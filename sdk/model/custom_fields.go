// Package model 中的 custom_fields.go 提供自定义字段序列化/反序列化的通用工具函数
package model

import (
	"encoding/json"
	"strings"
)

// IsCustomField 判断 JSON key 是否为 TAPD 自定义字段
// 自定义字段包括 custom_field_* 和 custom_plan_field_* 两类
func IsCustomField(key string) bool {
	return strings.HasPrefix(key, "custom_field_") || strings.HasPrefix(key, "custom_plan_field_")
}

// ExtractCustomFields 从原始 JSON map 中提取所有自定义字段，返回非空值的键值对
func ExtractCustomFields(raw map[string]json.RawMessage) map[string]string {
	var result map[string]string
	for key, val := range raw {
		if !IsCustomField(key) {
			continue
		}
		var s string
		if err := json.Unmarshal(val, &s); err == nil && s != "" {
			if result == nil {
				result = make(map[string]string)
			}
			result[key] = s
		}
	}
	return result
}

// MergeCustomFields 将自定义字段 map 写入请求参数 map
func MergeCustomFields(params map[string]string, customFields map[string]string) {
	for k, v := range customFields {
		if v != "" {
			params[k] = v
		}
	}
}
