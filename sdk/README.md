# tapd-sdk-go

[![Go](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](../LICENSE)

TAPD（腾讯敏捷产品研发平台）Open API 的 Go SDK，零外部依赖，仅使用标准库实现。

## API 文档来源

本 SDK 的设计基于 TAPD 官方开放平台 API 文档：

> <https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/>

TAPD Open API 共提供约 192 个接口，覆盖 20 个模块。本 SDK 封装了其中最常用的核心接口。

## 功能概览

SDK 当前覆盖 12 种 TAPD 资源类型，提供 **40 个导出方法**：

| 资源 | 方法数 | 支持操作 |
|------|--------|----------|
| **Workspace（项目）** | 2 | 列表、详情 |
| **Story（需求）** | 5 | 列表、详情、创建、更新、计数 |
| **Task（任务）** | 5 | 列表、详情、创建、更新、计数 |
| **Bug（缺陷）** | 5 | 列表、详情、创建、更新、计数 |
| **Iteration（迭代）** | 4 | 列表、创建、更新、计数 |
| **Comment（评论）** | 4 | 列表、添加、更新、计数 |
| **Wiki（文档）** | 4 | 列表、详情、创建、更新 |
| **TCase（测试用例）** | 4 | 列表、创建、批量创建、计数 |
| **Timesheet（工时）** | 3 | 列表、添加、更新 |
| **Attachment（附件）** | 2 | 附件列表、图片下载链接 |
| **Workflow（工作流）** | 3 | 状态流转、状态映射、结束状态 |
| **其他** | 若干 | 发布计划、需求分类、自定义字段、实体关联、待办事项、企业微信消息等 |

## 安装

```bash
go get github.com/studyzy/tapd-sdk-go
```

## 快速开始

```go
package main

import (
    "fmt"
    "log"

    tapd "github.com/studyzy/tapd-sdk-go"
    "github.com/studyzy/tapd-sdk-go/model"
)

func main() {
    // 方式一：使用 Access Token（推荐）
    client := tapd.NewClient("your-access-token", "", "")

    // 方式二：使用 Basic Auth
    // client := tapd.NewClient("", "api_user", "api_password")

    // 方式三：连接自定义 TAPD 站点
    // client := tapd.NewClientWithBaseURL(
    //     "https://api.my-tapd.com",   // API 地址
    //     "https://www.my-tapd.com",   // 前端页面地址
    //     "your-access-token", "", "",
    // )

    // 获取用户参与的项目列表
    workspaces, err := client.ListWorkspaces()
    if err != nil {
        log.Fatal(err)
    }
    for _, ws := range workspaces {
        fmt.Printf("项目: %s (ID: %s)\n", ws.Name, ws.ID)
    }

    // 查询需求列表
    stories, err := client.ListStories(&model.ListStoriesRequest{
        WorkspaceID: "12345678",
        Status:      "open",
        Limit:       10,
    })
    if err != nil {
        log.Fatal(err)
    }
    for _, s := range stories {
        fmt.Printf("需求: %s [%s]\n", s.Name, s.Status)
    }
}
```

## 架构设计

### 模块结构

```
sdk/
├── client.go          # Client 核心：认证、HTTP 封装、响应解析
├── story.go           # 需求 API
├── task.go            # 任务 API
├── bug.go             # 缺陷 API
├── iteration.go       # 迭代 API
├── comment.go         # 评论 API
├── wiki.go            # Wiki 文档 API
├── tcase.go           # 测试用例 API
├── timesheet.go       # 工时 API
├── attachment.go      # 附件/图片 API
├── workflow.go        # 工作流 API
├── category.go        # 需求分类 API
├── custom_field.go    # 自定义字段 API
├── relation.go        # 实体关联 API
├── misc.go            # 杂项：发布计划、待办事项、企业微信等
├── go.mod             # 独立模块，零外部依赖
└── model/             # 数据模型与请求参数
    ├── model.go       # 通用模型（TAPDResponse、ListResponse 等）
    ├── story.go       # Story 及其请求结构体
    ├── task.go        # Task 及其请求结构体
    ├── bug.go         # Bug 及其请求结构体
    ├── iteration.go   # Iteration 及其请求结构体
    ├── comment.go     # Comment 及其请求结构体
    ├── wiki.go        # Wiki 及其请求结构体
    ├── tcase.go       # TCase 及其请求结构体
    ├── timesheet.go   # Timesheet 及其请求结构体
    ├── attachment.go  # Attachment/ImageInfo 模型
    ├── category.go    # Category 模型
    ├── release.go     # Release 模型
    ├── workflow.go    # Workflow 模型
    └── request.go     # 通用请求结构体
```

### 设计原则

1. **零外部依赖**：仅使用 Go 标准库，可被任何项目无冲突引入。
2. **独立模块**：SDK 作为独立 Go module（`github.com/studyzy/tapd-sdk-go`），与 CLI 解耦，可被第三方项目直接引用。
3. **统一模式**：所有资源 API 遵循一致的调用模式：
   - 请求参数定义为 `model/` 下的结构体，通过 `ToParams()` 方法转为 `map[string]string`
   - Client 方法调用内部 `doGet`/`doPost` 发起 HTTP 请求
   - 自动解析 TAPD 统一响应格式 `{status, data, info}`
   - 自动处理 TAPD 特有的数据包裹格式（如 `[{"Story": {...}}]`）
4. **双模式认证**：支持 Bearer Token（推荐）和 Basic Auth 两种认证方式。
5. **可配置站点地址**：API 和前端 URL 均可自定义，支持连接非 `tapd.cn` 的 TAPD 部署。
6. **结构化错误**：返回 `TAPDError` 包含 HTTP 状态码、退出码和错误消息，便于调用方精确处理。

### 认证方式

| 方式 | 环境变量 | 说明 |
|------|----------|------|
| Bearer Token（推荐） | `TAPD_ACCESS_TOKEN` | 通过 TAPD 个人设置获取 |
| Basic Auth | `TAPD_API_USER` + `TAPD_API_PASSWORD` | API 账号密码 |

当同时提供两种凭据时，Bearer Token 优先。

### 自定义站点地址

SDK 默认连接 `https://api.tapd.cn`，如需连接其他 TAPD 部署，使用 `NewClientWithBaseURL`：

```go
client := tapd.NewClientWithBaseURL(
    "https://api.my-tapd.com",   // apiURL：API 请求地址
    "https://www.my-tapd.com",   // webURL：前端页面地址（用于生成条目链接）
    "your-access-token", "", "",
)

// 获取前端页面基础地址
fmt.Println(client.WebURL()) // https://www.my-tapd.com
```

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `apiURL` | `https://api.tapd.cn` | API 请求基础地址，传空字符串使用默认值 |
| `webURL` | `https://www.tapd.cn` | 前端页面基础地址，传空字符串使用默认值 |

### 错误处理

SDK 返回结构化的 `TAPDError`，包含以下退出码映射：

| HTTP 状态码 | 退出码 | 含义 |
|-------------|--------|------|
| 401 | 1 | 认证失败 |
| 404 | 2 | 资源未找到 |
| 422 | 3 | 参数错误 |
| 其他 | 4 | API 错误 |

```go
stories, err := client.ListStories(req)
if err != nil {
    if tapdErr, ok := err.(*tapd.TAPDError); ok {
        fmt.Printf("退出码: %d, HTTP: %d, 消息: %s\n",
            tapdErr.ExitCode, tapdErr.HTTPStatus, tapdErr.Message)
    }
}
```

## API 覆盖对照

下表对比 TAPD Open API 全量接口与本 SDK 的覆盖情况：

| TAPD 模块 | 官方接口数 | SDK 已覆盖 | 覆盖范围 |
|-----------|-----------|-----------|----------|
| 需求 (Story) | 33 | 5 | 基础 CRUD + 计数 |
| 缺陷 (Bug) | 21 | 5 | 基础 CRUD + 计数 |
| 任务 (Task) | 10 | 5 | 基础 CRUD + 计数 |
| 迭代 (Iteration) | 8 | 4 | 列表、创建、更新、计数 |
| 测试 (TCase) | 27 | 4 | 列表、创建、批量创建、计数 |
| Wiki | 11 | 4 | 基础 CRUD |
| 评论 (Comment) | 4 | 4 | 全部覆盖 |
| 工时 (Timesheet) | 5 | 3 | 列表、添加、更新 |
| 附件 (Attachment) | 3 | 2 | 附件列表、图片链接 |
| 项目 (Workspace) | 17 | 2 | 列表、详情 |
| 工作流 (Workflow) | 6 | 3 | 流转、状态映射、结束状态 |
| 发布 (Release) | 11 | 1 | 列表查询 |
| 配置 (Setting) | 21 | 4 | 自定义字段、字段信息、需求类别 |
| 其他 | - | 若干 | 待办事项、实体关联、源码提交关键字、企业微信消息 |

## 测试

SDK 使用 `net/http/httptest` 进行 mock server 测试：

```bash
cd sdk
go test ./...
```

## 许可证

[Apache License 2.0](../LICENSE)
