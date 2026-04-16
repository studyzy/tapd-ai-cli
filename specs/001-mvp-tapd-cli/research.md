# 研究报告: tapd-ai-cli MVP

## R1: TAPD Open API 认证与请求模式

### Decision: 支持双认证模式（Bearer Token + Basic Auth）+ 标准 net/http 客户端

### Rationale:
- TAPD Open API 支持两种认证方式（参考 TAPD MCP Server 实现 `mcp-server-tapd/src/mcp_server_tapd/tapd.py:17-30`）:
  - **Access Token（推荐）**: `Authorization: Bearer <token>`，使用 `TAPD_ACCESS_TOKEN`
  - **API User/Password**: `Authorization: Basic <base64(user:password)>`，使用 `TAPD_API_USER` + `TAPD_API_PASSWORD`
- Access Token 认证时可调用 `GET /users/info` 获取当前用户 nick，用于创建/更新操作的 creator 字段
- 认证优先级: access_token > api_user/api_password（与 MCP Server 一致）
- API 基地址: `https://api.tapd.cn`
- 认证测试端点: `GET /quickstart/testauth` — 返回 `status=1` 表示认证成功，可用于 `auth login` 时验证凭据
- 请求格式: GET（query string）、POST（`application/x-www-form-urlencoded` 或 `application/json`）
- 响应格式统一为 JSON: `{"status": 1, "data": [...], "info": "success"}`
- 列表数据包装在实体类型键下: `data[].Story`、`data[].Bug`、`data[].Task`、`data[].Iteration`
- 分页: 默认 30 条/页，最大 200 条/页，使用 `limit` + `page` 参数，总数需调用专门的 count 接口
- 速率限制: 默认 60 req/min，超限返回 HTTP 429
- 高级查询语法: 枚举用 `|` 分隔，时间范围用 `~`，支持 `LIKE`/`EQ`/`NOT_EQ` 操作符
- 错误码: 401（认证失败）、404（资源不存在）、422（参数错误）、429（限流）、500/502（服务端错误）
- 标准 `net/http` 已完全满足需求，无需引入 `go-resty/resty`，减少依赖

### Alternatives considered:
- `go-resty/resty`: 功能丰富但对本项目来说过重，项目仅需 Basic Auth + JSON，标准库足够
- `net/http` + 自定义封装: **选择此方案**，轻量、无额外依赖、完全可控

---

## R2: HTML 转 Markdown 库

### Decision: 使用 `github.com/JohannesKaufmann/html-to-markdown/v2`

### Rationale:
- Go 生态中该领域的事实标准（3,600+ Stars）
- 2025 年仍活跃维护（v2.5.0 发布于 2025-11-30）
- API 极其简洁，最简场景一行代码: `htmltomarkdown.ConvertString(input)`
- 支持表格、删除线等插件扩展
- MIT 许可证，Goroutine 安全
- 所有备选方案要么已停止维护（lunny/html2md 2019 年归档），要么功能不匹配（jaytaylor/html2text 转纯文本而非 Markdown）

### Alternatives considered:
- `github.com/lunny/html2md`: 2019 年归档，已停止维护 7 年，不推荐
- `github.com/jaytaylor/html2text`: 转纯文本而非 Markdown，不符合需求
- `github.com/suntong/html2md`: 仅是 JohannesKaufmann 库的 CLI 封装

---

## R3: 配置管理与多来源凭据

### Decision: 使用 `spf13/viper` 管理配置，自定义优先级链

### Rationale:
- Viper 原生支持环境变量绑定、配置文件读取、多路径搜索
- 凭据优先级链: CLI flags > ENV vars > ./tapd.json > ~/.tapd.json
- Viper 的 `SetConfigName` + `AddConfigPath` 可以实现多路径搜索
- 但 Viper 默认合并所有来源，我们需要严格优先级（不合并），需要自定义逻辑:
  1. 先检查 CLI flags（Cobra PersistentFlags 绑定）
  2. 再检查 ENV（Viper AutomaticEnv + SetEnvPrefix）
  3. 再尝试读取 `./.tapd.json`
  4. 最后尝试读取 `~/.tapd.json`
- workspace_id 写入逻辑: `workspace switch` 始终写入当前目录的 `./.tapd.json`（不存在则创建），使每个项目目录有独立的 workspace 配置

### Alternatives considered:
- 纯手写 JSON 读写: 可行但需自行处理环境变量和优先级逻辑，Viper 提供了更好的抽象
- `github.com/knadh/koanf`: 更现代但生态不如 Viper，且 Cobra 与 Viper 集成更成熟

---

## R4: CLI 框架与命令组织

### Decision: 使用 `spf13/cobra`，按实体类型组织子命令

### Rationale:
- Cobra 是 Go CLI 事实标准，与 Viper 无缝集成
- 命令树结构:
  ```
  root (tapd)
  ├── auth → login
  ├── workspace → list, switch, info
  ├── story → list, show, create, update, count
  ├── task → list, show, create, update, count
  ├── bug → list, show, create, update, count
  ├── iteration → list
  └── spec
  ```
- 全局 PersistentFlags: `--workspace-id`, `--compact`
- Cobra 的 `TraverseChildren` 特性支持在父命令设置全局标志
- `spec` 命令可遍历 Cobra 命令树自动生成 Tool Definition

### Alternatives considered:
- `github.com/urfave/cli`: 功能足够但与 Viper 集成不如 Cobra
- 不使用框架: 工作量大且命令行解析代码难以维护

---

## R5: JSON 输出策略

### Decision: 自定义 output 包，基于 `encoding/json` + omitempty

### Rationale:
- Go 标准库 `encoding/json` 的 `omitempty` 标签直接满足"移除 null 字段"需求
- 紧凑模式: `json.Marshal`（默认紧凑），非紧凑模式: `json.MarshalIndent`
- 列表截断逻辑: 先通过 count API 获取总数，再分页获取，输出中包含 `total`、`has_more`
- 错误输出到 stderr: 使用 `fmt.Fprintf(os.Stderr, ...)` 输出 JSON 格式错误
- 退出码: 在 cmd 层统一通过 `os.Exit(code)` 处理

### Alternatives considered:
- 第三方 JSON 库（jsoniter、sonic）: 性能更好但 CLI 工具无需高吞吐，标准库足够

---

## R6: 测试策略

### Decision: httptest mock TAPD API + 表驱动测试

### Rationale:
- `internal/client/`: 使用 `httptest.NewServer` 模拟 TAPD API 响应，测试请求构造和响应解析
- `internal/config/`: 使用临时目录和临时文件测试配置读写和优先级链
- `internal/output/`: 纯函数测试，验证 JSON 格式、omitempty、compact 模式
- `cmd/`: 使用 Cobra 的 `Execute` 方法和 `bytes.Buffer` 捕获 stdout/stderr，测试端到端命令行为
- 覆盖率目标: ≥60%，优先覆盖 client 和 config 包

### Alternatives considered:
- 录制/回放（go-vcr）: 需要真实 API 调用录制，初始设置成本高
- 集成测试直连 TAPD: 需要真实凭据和测试项目，不适合 CI

---

## R7: TAPD API 实体操作模式

### Decision: Story/Task 共用一套 API 方法，Bug 独立一套

### Rationale:
- TAPD API 中 Story 和 Task 共用同一组端点（`get_stories_or_tasks`、`create_story_or_task` 等），通过 `entity_type` 参数区分
- Bug 使用独立的端点（`get_bugs`、`add_bug` 等），字段名也不同（如 Bug 用 `title` 而非 `name`）
- 在 client 层: story.go 和 task.go 可共享底层请求方法，通过 entity_type 参数区分
- 在 cmd 层: story 和 task 命令代码结构相似，可通过工厂函数减少重复

### TAPD API 端点映射:

| CLI 操作 | TAPD API 端点 | HTTP 方法 |
|----------|---------------|-----------|
| story/task list | `/stories` 或 `/tasks` | GET |
| story/task show | `/stories` 或 `/tasks` (id 参数) | GET |
| story/task create | `/stories` 或 `/tasks` | POST |
| story/task update | `/stories` 或 `/tasks` | POST |
| story/task count | `/stories/count` 或 `/tasks/count` | GET |
| bug list | `/bugs` | GET |
| bug show | `/bugs` (id 参数) | GET |
| bug create | `/bugs` | POST |
| bug update | `/bugs` | POST |
| bug count | `/bugs/count` | GET |
| workspace list | `/workspaces/user_participant_projects` | GET |
| workspace info | `/workspaces/get_workspace_info` | GET |
| iteration list | `/iterations` | GET |

所有端点的基地址: `https://api.tapd.cn`
