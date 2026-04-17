# 告别协议臃肿：为什么 CLI 才是 AI Agent 接入 TAPD 的最优解？

在 2026 年的 AI 开发者社区，一个曾被视为"标准答案"的协议正在遭遇严峻挑战。当 Anthropic 最初推出 **Model Context Protocol (MCP)** 时，业界曾寄希望于它能统一 AI 工具调用的接口。然而，随着生产环境的深入应用，开发者们发现了一个尴尬的现实：**MCP 太重了，重到让 Agent 变得迟钝且昂贵。**

### 一、效率之争：为什么 MCP 正在失宠？

最近，包括 **Scalekit** 和 **Jannik Reinhard** 在内的多家机构发布了针对 AI 工具调用协议的基准测试（Benchmark）。数据揭示了一个令人震惊的趋势：

* **Token 膨胀**：在执行类似"查询代码库详情"的任务时，**CLI 模式比 MCP 便宜了 10 到 32 倍**。原因很简单：MCP 需要在对话开始前将所有工具的 JSON Schema 注入上下文。接入 3 个 MCP 服务器可能瞬间烧掉 5 万个 Token，而 CLI 是**按需调用（On-demand）**。
* **可靠性黑洞**：Scalekit 的测试显示，CLI 工具的执行成功率为 **100%**，而 MCP 因为协议握手超时和复杂的参数映射，成功率仅为 **72%**。
* **渐进式披露（Progressive Disclosure）**：正如 OpenClaw 作者 Peter Steinberger 所言："Bash 才是 AI 的原生语言。"CLI 允许 Agent 通过 `--help` 逐步探测功能，而不是被一次性塞入数千行的协议文档。

正是基于这一趋势，我们开发了 **`tapd-ai-cli`**——一个专为 AI Agent 深度定制的、摒弃了冗余协议负担的 TAPD 命令行工具。

---

### 二、`tapd-ai-cli` 的核心设计思想

与传统为人设计的 CLI 不同，`tapd-ai-cli` 的每一行代码都为 **AI 的理解成本**和**运行效率**服务。以下是贯穿整个项目的五条核心设计原则：

#### 原则一：AI 友好的输出格式——每个 Token 都不浪费

这是 `tapd-ai-cli` 最核心的差异化设计。我们在输出层做了三级优化：

* **紧凑 JSON 为默认**：所有输出默认为单行、无缩进、带 `omitempty` 的 JSON。空字段、冗余时间戳、UI 渲染属性被自动剔除，只保留 AI 决策所需的关键 Fact。
* **Markdown 表格优先**：列表类数据自动转为 Markdown 表格输出。相比 JSON，表格格式节省约 **30%-50%** 的 Token，且 AI 对表格的语义对齐能力更强。仍可通过 `--json` 切换回 JSON。
* **HTML→Markdown 自动转换**：TAPD API 返回的富文本描述为 HTML，CLI 层自动转换为 Markdown。AI 无需处理 `<p>`、`<strong>` 等标签噪声，直接获得结构化的纯文本。
* **智能截断与分页提示**：列表默认只返回 10 条并附分页提示，避免 Agent 被海量数据淹没。

#### 原则二：渐进式功能发现——`--help` 就是最好的文档

Agent 不需要一次性加载所有 API 文档。`tapd-ai-cli` 的 `--help` 输出经过精心设计，是一张**紧凑的命令参考卡**：

```
tapd story list [--status=<用 workflow status-map 查询可用值>] [--owner] [--limit=10]
tapd bug create --title=<必需> [--priority=<urgent/high/medium/low/insignificant>]
```

每条命令一行，必选参数和可选参数一目了然，枚举值内联展示。Agent 可以逐层探测：先 `tapd --help` 了解命令树，再 `tapd story --help` 了解子命令，最后 `tapd story list --help` 了解具体参数——这就是**渐进式披露**。

#### 原则三：无状态原子操作——适配 Agent 的离散推理

Agent 的推理过程本质上是离散的，每次工具调用是一个独立的决策点。`tapd-ai-cli` 遵循 Unix 哲学，每个命令都是一次**自包含的原子操作**：

* 通过 `tapd workspace switch <id>` 切换项目后，工作区上下文自动写入 `.tapd.json`，后续所有命令自动关联，无需反复传递长串 ID。
* 认证凭据通过 `tapd auth login` 一次配置，自动沿优先级链查找（CLI flags > 环境变量 > `./.tapd.json` > `~/.tapd.json`）。
* 每条命令输出完整、自包含的结果，不依赖上一次调用的状态。

#### 原则四：标准管道组合——释放 Agent 的编排能力

AI 对 Shell 管道的理解超乎想象。CLI 模式天然支持 Unix 管道组合，让 Agent 可以在单次 Shell 调用中完成复杂编排：

```bash
# 查看当前迭代的所有待办需求
tapd story list --status=open --iteration-id=12345 --json | jq '.[].id'

# 给需求添加评论
echo "已完成代码审查" | tapd comment add --entry-type=stories --entry-id=100001

# 从文件创建 Wiki
tapd wiki create --name="API设计文档" --file=./design.md
```

这种能力在 MCP 的请求-响应模式中很难实现——每一步都要等待协议往返。

#### 原则五：SDK 与 CLI 分层——一套投入、两种接入

项目采用清晰的两层架构：

* **SDK 层（[`tapd-sdk-go`](https://github.com/studyzy/tapd-sdk-go)）**：已发布为独立 Go 模块，零第三方依赖，封装所有 TAPD API 调用，返回强类型结构体。可通过 `go get github.com/studyzy/tapd-sdk-go@latest` 直接引入，也可以基于它构建 MCP Server。
* **CLI 层**：消费 SDK 接口，只负责参数解析、HTML→Markdown 转换和输出格式化。

这意味着如果未来需要 MCP 接入，只需在 SDK 之上再包一层薄适配，而非重写所有 API 调用逻辑。

---

### 三、实测对比：CLI vs MCP

为了验证效果，我们针对"在一个包含 20 个需求的迭代中提取待办事项并创建任务"这一典型场景进行了测试：

| 指标 | 传统 MCP 模式 | `tapd-ai-cli` 模式 | 提升/节省 |
| :--- | :--- | :--- | :--- |
| **首字响应延迟 (TTFT)** | ~4.5s (加载大量 Schema) | **~0.8s** | **82% 提速** |
| **单次任务消耗 Token** | 18,400 Tokens | **1,250 Tokens** | **14.7 倍节省** |
| **操作成功率 (一次性)** | 78% | **100%** | **显著提升** |
| **AI 解析负担** | 高（需处理复杂 JSON 层级） | **低 (精简 MD 表格)** | **响应更精准** |

---

### 四、快速上手

#### 安装

```bash
go install github.com/studyzy/tapd-ai-cli/cmd/tapd@latest
```

#### 认证配置

```bash
# 使用 Access Token（推荐）
export TAPD_ACCESS_TOKEN=your_token

# 或通过交互式登录持久化凭据
tapd auth login
```

#### 日常使用示例

```bash
# 切换到目标项目
tapd workspace switch 12345678

# 查看当前迭代的需求列表
tapd story list --status=open --iteration-id=100

# 查看需求详情（描述自动从 HTML 转为 Markdown）
tapd story show 1012345

# 创建任务并关联需求
tapd task create --name="实现登录接口" --story-id=1012345 --owner=dev1

# 提交缺陷
tapd bug create --title="登录页 CSS 错位" --severity=normal --priority=medium

# 查询工作流状态映射（了解有哪些可用状态值）
tapd workflow status-map --system=story --workitem-type-id=123

# 更新需求状态
tapd story update 1012345 --status=progressing

# 查看关联缺陷
tapd relation bugs --story-id=1012345

# 记录工时
tapd timesheet add --entity-type=task --entity-id=200001 --timespent=2 --owner=dev1
```

#### 在 AI Agent 中使用

将 `tapd` 加入 Agent 的可用工具列表后，Agent 会自动通过 `tapd --help` 发现所有功能。更推荐通过 **`tapd skill init`** 一键生成 SKILL.md 指令文件，让 AI Coding 工具主动理解如何使用 tapd CLI：

```bash
# 为 Claude Code、CodeBuddy、Cursor 等 10 种工具一键生成 SKILL.md
tapd skill init
```

命令会自动检测当前目录下已有的 AI Coding 工具配置文件夹并默认选中，交互式确认后将 SKILL.md 生成到对应工具的 `skills/tapd/SKILL.md` 路径下。生成的命令参考部分从当前 CLI 版本的命令树动态生成，始终与实际功能保持一致。

支持的工具：Claude Code、CodeBuddy、Cursor、Windsurf、Trae、Codex、Gemini CLI、Cline、Roo Code、Augment。

典型的 Agent 工作流：

1. `tapd workspace list` — 发现可用项目
2. `tapd workspace switch <id>` — 锁定目标项目
3. `tapd story list --status=open` — 获取待办需求
4. `tapd story show <id>` — 深入了解某个需求
5. `tapd task create --name=... --story-id=...` — 拆解并创建任务

---

### 五、结语

在 AI 编程（Vibe Coding）时代，我们不需要更复杂的协议，我们需要更直接的工具。`tapd-ai-cli` 不仅仅是一个命令行工具，它是对 **"AI 友好型基础设施"** 的一次实践——用最少的 Token 传递最多的信息，用最简单的接口释放最大的编排能力。

如果您正苦恼于 Agent 在 TAPD 项目管理中的迟钝和高额成本，欢迎尝试切换到 CLI 模式。

**项目主页**: [github.com/studyzy/tapd-ai-cli](https://github.com/studyzy/tapd-ai-cli)
**SDK 仓库**: [github.com/studyzy/tapd-sdk-go](https://github.com/studyzy/tapd-sdk-go)
