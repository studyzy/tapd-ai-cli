# 快速入门: tapd-sdk-go

**分支**: `002-extract-tapd-sdk` | **日期**: 2026-04-17

## 概述

`tapd-sdk-go` 是 TAPD 平台的 Go SDK，提供对需求、缺陷、任务、迭代、评论、Wiki 等全部资源的类型安全访问。SDK 仅依赖 Go 标准库，可独立嵌入任何 Go 项目。

---

## 安装

```bash
go get github.com/studyzy/tapd-sdk-go
```

在开发阶段（与 tapd-ai-cli 同仓库），CLI 的 `go.mod` 使用本地引用：

```
replace github.com/studyzy/tapd-sdk-go => ./sdk
```

---

## 初始化客户端

SDK 支持两种认证方式：

```go
import "github.com/studyzy/tapd-sdk-go"

// 方式一：Bearer Token（个人访问令牌）
client := tapd.NewClient("your-access-token", "", "")
client.FetchNick() // 可选：获取当前用户昵称

// 方式二：Basic 认证（API 账号/密码）
client := tapd.NewClient("", "api-user", "api-password")

// 验证认证是否有效
if err := client.TestAuth(); err != nil {
    log.Fatal("认证失败:", err)
}
```

---

## 查询需求列表

```go
import "github.com/studyzy/tapd-sdk-go/model"

stories, err := client.ListStories(&model.ListStoriesRequest{
    WorkspaceID: "12345678",
    Status:      "doing",
    Limit:       "20",
    Page:        "1",
})
if err != nil {
    log.Fatal(err)
}
for _, s := range stories {
    fmt.Printf("[%s] %s\n", s.ID, s.Name)
}
```

---

## 获取单个需求

```go
story, err := client.GetStory("12345678", "story-id-001")
if err != nil {
    // 检查是否为 404
    var tapdErr *tapd.TAPDError
    if errors.As(err, &tapdErr) && tapdErr.ExitCode == 2 {
        fmt.Println("需求不存在")
        return
    }
    log.Fatal(err)
}
// story.Description 为原始 HTML，调用方自行处理格式转换
fmt.Println(story.Name, story.Status)
```

---

## 创建需求

```go
resp, err := client.CreateStory(&model.CreateStoryRequest{
    WorkspaceID:   "12345678",
    Name:          "支持暗黑模式",
    Description:   "用户希望在夜间使用时切换到暗黑主题",
    PriorityLabel: "High",
    Owner:         "zhangsan",
    IterationID:   "iter-001",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("创建成功，ID: %s，URL: %s\n", resp.ID, resp.URL)
```

---

## 错误处理

```go
import "errors"

_, err := client.GetBug("12345678", "nonexistent-id")
if err != nil {
    var tapdErr *tapd.TAPDError
    if errors.As(err, &tapdErr) {
        switch tapdErr.ExitCode {
        case 1:
            fmt.Println("认证失败，请检查凭据")
        case 2:
            fmt.Println("资源不存在")
        case 3:
            fmt.Println("请求参数错误:", tapdErr.Message)
        case 4:
            fmt.Println("服务端错误，请稍后重试:", tapdErr.Message)
        }
        return
    }
    // 网络错误等非 TAPD 错误
    log.Fatal("网络错误:", err)
}
```

---

## 在 CLI 中集成 SDK

CLI 层负责：读取配置文件 → 创建 SDK 客户端 → 调用 SDK 方法 → 格式化输出

```go
// internal/cmd/root.go（示意）
import (
    tapd "github.com/studyzy/tapd-sdk-go"
    "github.com/studyzy/tapd-ai-cli/internal/config"
)

func newClient() *tapd.Client {
    cfg := config.Load() // 读取 .tapd.json
    return tapd.NewClient(cfg.AccessToken, cfg.APIUser, cfg.APIPassword)
}
```

---

## 本地测试（使用 httptest）

```go
import "net/http/httptest"

server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"status":1,"data":[{"Story":{"id":"001","name":"测试需求"}}]}`))
}))
defer server.Close()

client := tapd.NewClientWithBaseURL(server.URL, "mock-token", "", "")
stories, err := client.ListStories(&model.ListStoriesRequest{WorkspaceID: "ws1"})
// 验证 stories[0].Name == "测试需求"
```
