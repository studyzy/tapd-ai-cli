# spec 命令输出格式重构：紧凑参考卡

## 背景

`tapd spec` 命令当前输出 OpenAI/Anthropic 兼容的 Tool Definition JSON 数组。但实际使用场景中，消费者是 AI Agent 通过阅读 spec 输出来了解可用命令，然后通过 Bash 调用 CLI。JSON Tool Definition 格式存在以下问题：

1. **Token 开销大** — JSON 结构性重复（`"type": "object"`、`"type": "string"` 等）消耗大量 token
2. **格式不匹配场景** — Tool Definition 是给 function calling API 注册用的，但这里只需要让 AI 理解"有哪些命令可用"
3. **可读性差** — 嵌套 JSON schema 对于"快速理解命令结构"的目标过于冗余

## 决策

将 `tapd spec` 输出从 JSON Tool Definition 数组改为**紧凑参考卡**纯文本格式。每个命令一行，使用标准 CLI 语法标记参数。

## 输出格式

### 语法约定

```
<arg>           — 必填位置参数
--flag=<val>    — 必填标志
[--flag]        — 可选标志
[--flag=default]— 可选标志，有默认值
```

### 输出结构

1. 标题行 + 全局标志说明
2. 按命令分组，组间空行分隔，组名用 `#` 注释
3. 每行一个完整可执行命令模板 + `#` 注释描述

### 示例

```
tapd - TAPD CLI for AI Agent
Global: [--workspace-id=<id>] [--json] [--pretty] [--no-comments]

# auth
tapd auth login [--access-token=<token>] [--api-user=<user>] [--api-password=<pwd>] [--local]  # 登录并持久化凭据

# workspace
tapd workspace list  # 列出参与的项目
tapd workspace switch <workspace_id>  # 切换当前工作区
tapd workspace info  # 查看当前工作区详情

# story
tapd story list [--status] [--owner] [--iteration-id] [--limit=10] [--page=1]  # 查询需求列表
tapd story show <story_id>  # 查看需求详情
tapd story create --name=<name> [--description] [--owner] [--priority] [--iteration-id]  # 创建需求
tapd story update <story_id> [--name] [--status] [--owner] [--priority]  # 更新需求
tapd story count [--status]  # 查询需求数量
tapd story todo [--limit=10] [--page=1]  # 查询当前用户待办需求
```

## 实现方案

### 改造 `internal/cmd/spec.go`

1. 删除 `toolDefinition`、`toolParameters`、`toolProperty` 结构体
2. `runSpec` 改为逐行输出文本
3. `walkCommands` 保留递归遍历，改为收集 `specLine` 结构（组名 + 命令文本）
4. `commandToTool` 改为 `commandToLine`，生成单行文本
5. 从 flag 的 usage 文本中检测"必需"/"必填"关键字，对必填 flag 使用 `--flag=<val>` 格式
6. 输出头部标题和全局标志说明
7. 全局认证标志（access-token、api-user、api-password）从全局标志列表中排除（与 spec 无关）

### 必填标志检测

当前代码未使用 `MarkFlagRequired`，必填性通过 RunE 内部检查。通过检测 flag Usage 文本中的"必需"/"必填"关键字来判断。

### 分组逻辑

遍历时记录每个命令的第一级父命令名作为组名。输出时在组切换处插入 `# group-name` 行和空行。

## Token 节省估算

以 story 组（6 个命令）为例：
- JSON 格式：~800-1000 字符
- 紧凑参考卡：~400-500 字符

全部 44 个命令预计从 ~7000 字符压缩到 ~3500 字符，节省约 50%。
