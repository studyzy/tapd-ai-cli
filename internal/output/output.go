// Package output 提供 CLI 的 JSON 输出格式化工具函数
package output

import (
	"encoding/json"
	"io"

	"github.com/studyzy/tapd-ai-cli/internal/model"
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
