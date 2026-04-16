// Package output 提供 CLI 的输出格式化工具函数，支持 JSON 和 Markdown 两种格式
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/studyzy/tapd-ai-cli/internal/model"
	"gopkg.in/yaml.v3"
)

// 退出码常量，用于标识不同类型的命令执行结果
const (
	// ExitSuccess 表示命令执行成功
	ExitSuccess = 0
	// ExitAuthError 表示认证错误
	ExitAuthError = 1
	// ExitNotFound 表示资源未找到
	ExitNotFound = 2
	// ExitParamError 表示参数错误
	ExitParamError = 3
	// ExitAPIError 表示 API 调用错误
	ExitAPIError = 4
)

// PrintJSON 将数据以 JSON 格式写入 writer，支持紧凑和缩进两种模式
func PrintJSON(w io.Writer, data interface{}, compact bool) error {
	var b []byte
	var err error
	if compact {
		b, err = json.Marshal(data)
	} else {
		b, err = json.MarshalIndent(data, "", "  ")
	}
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = w.Write(b)
	return err
}

// PrintError 将错误信息以 JSON 格式写入 writer（始终使用紧凑模式）
func PrintError(w io.Writer, code string, message string, hint string) {
	resp := model.ErrorResponse{
		Error:   code,
		Message: message,
		Hint:    hint,
	}
	b, _ := json.Marshal(resp)
	b = append(b, '\n')
	w.Write(b)
}

// PrintSuccess 将成功响应以紧凑 JSON 格式写入 writer
func PrintSuccess(w io.Writer, resp interface{}) error {
	return PrintJSON(w, resp, true)
}

// PrintMarkdown 将数据以 YAML frontmatter + Markdown body 格式写入 writer
// bodyField 指定作为 Markdown 正文输出的字段名（JSON tag 名称，如 "description"）
// 其余非空字段输出为 YAML frontmatter
func PrintMarkdown(w io.Writer, data interface{}, bodyField string) error {
	meta, body, err := splitMarkdownFields(data, bodyField)
	if err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString("---\n")

	yamlBytes, err := yaml.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML frontmatter: %w", err)
	}
	sb.Write(yamlBytes)
	sb.WriteString("---\n")

	if body != "" {
		sb.WriteString("\n")
		sb.WriteString(body)
		if !strings.HasSuffix(body, "\n") {
			sb.WriteString("\n")
		}
	}

	_, err = w.Write([]byte(sb.String()))
	return err
}

// splitMarkdownFields 将结构体拆分为 frontmatter 的有序键值对和 body 字符串
// bodyField 是 JSON tag 名称，匹配到的字段值作为 body 返回，其余非空字段作为 meta 返回
func splitMarkdownFields(data interface{}, bodyField string) (meta yaml.Node, body string, err error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return meta, "", fmt.Errorf("PrintMarkdown only supports struct types, got %s", v.Kind())
	}

	t := v.Type()

	// 构建有序 YAML mapping node
	meta = yaml.Node{
		Kind: yaml.MappingNode,
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		// 获取 JSON tag 名称
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		tagName := strings.Split(jsonTag, ",")[0]

		// 跳过空值和整数类型的默认值（0、-1）
		strVal := fmt.Sprintf("%v", fieldVal.Interface())
		if strVal == "" {
			continue
		}
		if isIntegerKind(fieldVal.Kind()) && (strVal == "0" || strVal == "-1") {
			continue
		}
		// 字符串字段值为 "0" 或 "-1" 也视为无意义默认值，跳过
		if fieldVal.Kind() == reflect.String && (strVal == "0" || strVal == "-1") {
			continue
		}

		// body 字段单独提取
		if tagName == bodyField {
			body = strVal
			continue
		}

		// 添加到有序 mapping
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: tagName}
		valNode := &yaml.Node{Kind: yaml.ScalarNode, Value: strVal}
		meta.Content = append(meta.Content, keyNode, valNode)
	}

	return meta, body, nil
}

// isIntegerKind 判断是否为整数类型的 reflect.Kind
func isIntegerKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}
