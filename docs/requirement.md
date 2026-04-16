# 需求规格说明书：`tapd-ai-cli`

## 1. 项目概述
**目标**：使用 Golang 创建一个高性能、为 AI 优化的命令行界面（CLI），用于与 TAPD（腾讯敏捷产品开发）平台交互。
**目标用户**：AI 代理（如 Claude Code、GPT-4o），而非人类用户。
**设计理念**：
* **零噪音**：极简输出以节省 LLM Token。
* **无状态**：简单的请求-响应执行模式。
* **AI 优先**：结构化 JSON 输出、清晰的错误信息和机器可读的 Schema。

## 2. 技术栈
* **语言**：Golang 1.24+
* **CLI 框架**：`spf13/cobra`
* **API 通信**：`go-resty/resty` 或标准库 `net/http`
* **配置管理**：`spf13/viper`（用于 API 凭证）

## 3. 核心功能与命令
所有命令必须支持 `--json` 标志（尽可能默认为 `true`）和 `--project_id` 标志。

### A. 认证与工作区
* `tapd auth login --api-user <user> --api-key <key>`：在本地持久化凭证。
* `tapd workspace list`：列出可用的项目。
* `tapd workspace switch <id>`：在本地配置文件中设置"当前"工作区 ID，避免在每个命令中重复指定。

### B. 任务与故事管理（"精简"视图）
* `tapd task list [--status <status>] [--owner <me>]`：
    * **优化**：仅返回 `id`、`name`、`status` 和 `modified` 字段。
* `tapd task show <id>`：获取完整详情。
    * **优化**：使用 `html2md` 等库将 HTML 描述转换为 **Markdown**。
* `tapd task create --title <t> [--description <d>]`：创建新任务。

### C. AI 技能集成
* `tapd spec`：一个特殊命令，通过反射 Cobra 命令树输出 **工具定义 JSON**（兼容 OpenAI/Anthropic）。这允许 AI "自我发现" CLI 的能力。

## 4. AI 专属优化（关键）
1. **Token 精简**：
    * 从 JSON 输出中移除所有 `null` 字段（使用 `omitempty`）。
    * 使用 `--compact` 标志时，移除 JSON 中的缩进/空白。
    * 截断长列表（默认限制：10 条），并附加提示：`"更多条目可用，请使用 --limit。"`
2. **错误处理**：
    * 退出码必须区分明确（例如，`1` 表示认证错误，`2` 表示未找到）。
    * `Stderr` 应提供可操作的提示：`"错误：未找到任务 ID。您是指 ID 12345 吗？"`
3. **ID 映射（可选/第二阶段）**：
    * 维护一个本地缓存，将简单整数（`1`、`2`、`3`）映射到长 TAPD ID（最近访问的 10 个项目），以节省输入 Token。

## 5. 输出格式要求
* **标准输出**：数据命令输出原始、压缩的 JSON。
* **成功指示**：对于"创建/更新"操作，返回简单的 `{"success": true, "id": "..."}`。

---

## 参考代码
tapd mcp： /Users/devinzeng/Code/mcp-server-tapd  请参考这个代码实现MVP，MCP中的函数本cli都要实现，且输出都要符合上述规范。
api文档： 
- 使用必读 https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/%E4%BD%BF%E7%94%A8%E5%BF%85%E8%AF%BB.html
- API配置指引 https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/API%E9%85%8D%E7%BD%AE%E6%8C%87%E5%BC%95.html
- 研发协作API文档 https://open.tapd.cn/document/api-doc/API%E6%96%87%E6%A1%A3/api_reference/