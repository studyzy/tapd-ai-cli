// Package client 封装了 TAPD API 的 HTTP 客户端，处理认证、请求构造和响应解析
package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/studyzy/tapd-ai-cli/internal/model"
)

const defaultBaseURL = "https://api.tapd.cn"

// TAPDError 表示 TAPD API 返回的错误
type TAPDError struct {
	HTTPStatus int
	ExitCode   int
	Message    string
}

// Error 返回错误描述信息
func (e *TAPDError) Error() string {
	return e.Message
}

// Client 是 TAPD API 的 HTTP 客户端封装
type Client struct {
	baseURL    string
	httpClient *http.Client
	authHeader string // "Bearer <token>" or "Basic <base64>"
	Nick       string // 当前用户昵称，Bearer Token 认证时自动获取
}

// NewClient 创建一个使用默认 baseURL 的 TAPD API 客户端。
// 如果 accessToken 非空则使用 Bearer 认证，否则使用 Basic 认证。
func NewClient(accessToken, apiUser, apiPassword string) *Client {
	return NewClientWithBaseURL(defaultBaseURL, accessToken, apiUser, apiPassword)
}

// NewClientWithBaseURL 创建一个指定 baseURL 的 TAPD API 客户端，适用于测试场景。
// 如果 accessToken 非空则使用 Bearer 认证，否则使用 Basic 认证。
func NewClientWithBaseURL(baseURL, accessToken, apiUser, apiPassword string) *Client {
	var authHeader string
	if accessToken != "" {
		authHeader = "Bearer " + accessToken
	} else {
		encoded := base64.StdEncoding.EncodeToString([]byte(apiUser + ":" + apiPassword))
		authHeader = "Basic " + encoded
	}
	c := &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		authHeader: authHeader,
	}

	return c
}

// FetchNick 通过 /users/info 接口获取当前用户昵称（Bearer Token 认证时使用）
func (c *Client) FetchNick() {
	data, err := c.doGet("/users/info", nil)
	if err != nil {
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return
	}

	if nick, ok := result["nick"].(string); ok {
		c.Nick = nick
	}
}

// TestAuth 通过调用 /quickstart/testauth 接口验证凭据是否有效。
// 返回 nil 表示认证成功，否则返回错误信息。
func (c *Client) TestAuth() error {
	_, err := c.doGet("/quickstart/testauth", nil)
	return err
}

// doGet 发送 GET 请求到指定端点，解析 TAPD 统一响应格式并返回 data 字段
func (c *Client) doGet(endpoint string, params map[string]string) (json.RawMessage, error) {
	reqURL, err := url.Parse(c.baseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	if len(params) > 0 {
		q := reqURL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		reqURL.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Content-Type", "application/json")

	return c.doRequest(req)
}

// doPost 发送 POST 请求到指定端点，使用 form-urlencoded 编码请求体
func (c *Client) doPost(endpoint string, body map[string]string) (json.RawMessage, error) {
	form := url.Values{}
	for k, v := range body {
		form.Set(k, v)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.doRequest(req)
}

// doRequest 执行 HTTP 请求并解析 TAPD 统一响应格式
func (c *Client) doRequest(req *http.Request) (json.RawMessage, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		snippet := string(bodyBytes)
		if len(snippet) > 200 {
			snippet = snippet[:200]
		}
		return nil, &TAPDError{
			HTTPStatus: resp.StatusCode,
			ExitCode:   c.mapHTTPError(resp.StatusCode),
			Message:    fmt.Sprintf("HTTP %d: %s", resp.StatusCode, snippet),
		}
	}

	var tapdResp model.TAPDResponse
	if err := json.Unmarshal(bodyBytes, &tapdResp); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	if tapdResp.Status != 1 {
		return nil, &TAPDError{
			HTTPStatus: resp.StatusCode,
			ExitCode:   4,
			Message:    fmt.Sprintf("TAPD API error: %s", tapdResp.Info),
		}
	}

	return tapdResp.Data, nil
}

// doPostJSON 发送 JSON 格式的 POST 请求到指定 URL（用于企业微信等外部 API）
func (c *Client) doPostJSON(url string, body []byte) error {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		snippet := string(bodyBytes)
		if len(snippet) > 200 {
			snippet = snippet[:200]
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, snippet)
	}
	return nil
}

// mapHTTPError 将 HTTP 状态码映射到 CLI 退出码
func (c *Client) mapHTTPError(statusCode int) int {
	switch statusCode {
	case 401:
		return 1
	case 404:
		return 2
	case 422:
		return 3
	case 429, 500, 502:
		return 4
	default:
		return 4
	}
}
