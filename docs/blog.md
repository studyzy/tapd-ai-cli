# tapd-ai-cli——专为 AI Agent 打造的 TAPD 命令行工具

当我们想要让 AI Agent 帮忙管理 TAPD 项目时，往往会面临一个尴尬的选择：用 MCP 接入吧，Agent 变得又慢又贵；不用吧，似乎没有更好的标准化方案。

本人最近在做 AI Coding（Vibe Coding）相关的工作，需要让 Claude Code、CodeBuddy 等 AI Agent 频繁地与 TAPD 交互——查需求、建任务、提缺陷。一开始我也是用 MCP 方案，但很快就被它的 **Token 消耗**和**响应延迟**劝退了。于是我开始思考：有没有一种更轻量、对 AI 更友好的接入方式？

答案是 **CLI**。

这篇文章记录了我开发 `tapd-ai-cli` 的设计思路和实测数据，希望对同样在做 Agent 工具集成的朋友有所帮助。

---

## 一、背景：MCP 的问题出在哪？

MCP（Model Context Protocol）是 Anthropic 推出的 AI 工具调用协议，一度被视为"标准答案"。但随着生产环境的深入应用，越来越多的开发者和团队开始反思它的代价。近年来，多家机构发布了 MCP 与 CLI 的对比基准测试，数据相当有说服力。

### 1.1 Token 膨胀——每轮对话都在交"Schema 税"

MCP 的核心问题在于：**每轮对话都需要将所有工具的 JSON Schema 全量注入上下文**。打个比方，这就像你去餐厅吃饭，服务员先把整本菜单（包括每道菜的配料表、烹饪方法、营养成分）念一遍，然后才问你"请问吃什么"。

这个"菜单税"到底有多重？我们来看几组真实数据：

**Scalekit** 团队用 75 次 benchmark 运行对比了 GitHub CLI（`gh`）和 GitHub MCP Server（43 个工具），结果如下：

| 任务 | CLI Token 消耗 | MCP Token 消耗 | 倍数差距 |
| :--- | :--- | :--- | :--- |
| 获取仓库语言和许可证 | 1,365 | 44,026 | **32x** |
| 获取 PR 详情 | 1,648 | 32,279 | **20x** |
| 获取仓库元数据 | 9,386 | 82,835 | **9x** |
| 按贡献者汇总 PR | 5,010 | 33,712 | **7x** |

换算成钱更直观：假设每月执行 10,000 次操作，CLI 月成本约 **$3.20**，MCP 约 **$55.20**，差距 **17 倍**。

微软 MVP **Jannik Reinhard** 在企业场景（Intune 设备管理）中做了类似测试，差距更夸张：用 Microsoft Graph MCP 导出 50 台非合规设备消耗约 **145,000 Tokens**，而用 `mgc` + `az` CLI 只需约 **4,150 Tokens**，差距达 **35 倍**。

他还指出一个容易被忽视的问题：GitHub MCP 服务器 93 个工具的 Schema 就要消耗约 55,000 Tokens，**直接占掉 GPT-4o 上下文窗口的一半**。如果企业同时接入多个 MCP 服务器，工具定义可能超过 150,000 Tokens，留给实际推理的空间所剩无几。

### 1.2 可靠性——MCP 的成功率只有 72%

Token 贵还不是最要命的，更让人头疼的是**不稳定**。

Scalekit 的测试显示，CLI 工具的执行成功率为 **100%**，而 MCP 因为 TCP 握手超时和复杂的参数映射，成功率仅为 **72%**——25 次运行中有 7 次失败。Jannik Reinhard 也观察到，MCP 方案在 3-4 次工具调用后，多步推理开始崩溃。

对于生产环境中的 Agent 来说，72% 的成功率是不可接受的。

### 1.3 CLI 是模型的"母语"

今年年初，OpenClaw 作者 Peter Steinberger（后加入 OpenAI）在社交媒体上直言：

> **"MCP were a mistake. Bash is better."**

这不仅仅是情绪表达。背后的逻辑是：大语言模型的训练数据中包含**数十亿行终端交互记录**，CLI 命令对模型来说是真正的"母语"。模型生成一条 `gh pr view 123` 命令比构造一个复杂的 MCP JSON 请求体要自然得多，出错概率也更低。

Peter 不只是说说而已——他构建了 **MCPorter**，一个将 MCP 服务器转换为 CLI 的工具，用行动证明了自己的观点。他创建的 OpenClaw 项目在数周内获得了 190K GitHub Stars，其核心执行引擎正是基于 bash/CLI。

libgdx 作者 **Mario Zechner** 也做了 120 次评估运行的对比测试，得出了一个精辟的结论：**"工具设计比协议更重要。"** 如果从零开始，先做好 CLI——更简单、更通用；CLI 的输出还能通过管道过滤来提升 Token 效率，这是 MCP 做不到的。

### 1.4 行业趋势："CLI-first, MCP where needed"

值得一提的是，我并不是要说 MCP 一无是处。MCP 在需要多用户认证、有状态会话、企业级治理的场景中仍然有它的价值。但**行业共识正在形成**。正如 DeployHQ 在其对比分析中总结的：

> **CLI 覆盖 80% 的日常任务，MCP 用于 20% 需要有状态连接或无 CLI 替代的场景。**

正是基于这些观察和数据，我开发了 `tapd-ai-cli`。

---

## 二、核心设计思想

与传统为人设计的 CLI 不同，`tapd-ai-cli` 的每一行代码都为 **AI 的理解成本**和**运行效率**服务。以下是贯穿整个项目的五条核心设计原则。

### 2.1 AI 友好的输出格式——每个 Token 都不浪费

这是我认为最核心的差异化设计。我们在输出层做了几级优化：

- **紧凑输出为默认**：所有输出默认为 YAML 格式，带 `omitempty`。空字段、冗余时间戳、UI 渲染属性被自动剔除，只保留 AI 决策所需的关键信息。
- **Markdown 表格优先**：列表类数据自动转为 Markdown 表格输出。相比 JSON，表格格式节省约 **30%-50%** 的 Token，且 AI 对表格的语义对齐能力更强。仍可通过 `--json` 切换回 JSON。
- **HTML→Markdown 自动转换**：TAPD API 返回的富文本描述为 HTML，CLI 层自动转换为 Markdown。AI 无需处理 `<p>`、`<strong>` 等标签噪声，直接获得结构化的纯文本。
- **智能截断与分页提示**：列表默认只返回 10 条并附分页提示，避免 Agent 被海量数据淹没。

### 2.2 渐进式功能发现——`--help` 就是最好的文档

Agent 不需要一次性加载所有 API 文档。`tapd-ai-cli` 的 `--help` 输出经过精心设计，是一张**紧凑的命令参考卡**：

```
tapd story list [--status=<用 workflow status-map 查询可用值>] [--owner] [--limit=10]
tapd bug create --title=<必需> [--priority=<urgent/high/medium/low/insignificant>]
```

每条命令一行，必选参数和可选参数一目了然，枚举值内联展示。Agent 可以逐层探测：先 `tapd --help` 了解命令树，再 `tapd story --help` 了解子命令，最后 `tapd story list --help` 了解具体参数。这就是**渐进式披露（Progressive Disclosure）**，需要什么查什么，而不是一股脑全部塞进来。

### 2.3 无状态原子操作——适配 Agent 的离散推理

Agent 的推理过程本质上是离散的，每次工具调用是一个独立的决策点。`tapd-ai-cli` 遵循 Unix 哲学，每个命令都是一次**自包含的原子操作**：

- 通过 `tapd workspace switch <id>` 切换项目后，工作区上下文自动写入 `.tapd.json`，后续所有命令自动关联，无需反复传递长串 ID。
- 认证凭据通过 `tapd auth login` 一次配置，自动沿优先级链查找（CLI flags > 环境变量 > `./.tapd.json` > `~/.tapd.json`）。
- 每条命令输出完整、自包含的结果，不依赖上一次调用的状态。

### 2.4 标准管道组合——释放 Agent 的编排能力

AI 对 Shell 管道的理解超乎想象。CLI 模式天然支持 Unix 管道组合，让 Agent 可以在单次 Shell 调用中完成复杂编排。下面是几个典型的管道用法：

```bash
# 查看当前迭代的所有待办需求
tapd story list --status=open --iteration-id=12345 --json | jq '.[].id'

# 给需求添加评论
echo "已完成代码审查" | tapd comment add --entry-type=stories --entry-id=100001

# 从文件创建 Wiki
tapd wiki create --name="API设计文档" --file=./design.md
```

这种能力在 MCP 的请求-响应模式中很难实现——每一步都要等待协议往返。

### 2.5 SDK 与 CLI 分层——一套投入、两种接入

项目采用清晰的两层架构：

- **SDK 层（[`tapd-sdk-go`](https://github.com/studyzy/tapd-sdk-go)）**：已发布为独立 Go 模块，零第三方依赖，封装所有 TAPD API 调用，返回强类型结构体。可通过 `go get github.com/studyzy/tapd-sdk-go@latest` 直接引入，也可以基于它构建 MCP Server。
- **CLI 层**：消费 SDK 接口，只负责参数解析、HTML→Markdown 转换和输出格式化。

这意味着如果未来需要 MCP 接入，只需在 SDK 之上再包一层薄适配，而非重写所有 API 调用逻辑。

---

## 三、实测对比：CLI vs MCP 到底差多少？

光说原理不够，我们来看实际数据。

测试场景很简单：**查询一个 TAPD 需求单的详情**。CLI 模式下 Agent 需要两次调用——先 `tapd --help` 了解命令用法（约 7KB 的命令参考卡），再 `tapd story show <id>` 获取详情。MCP 模式只需一次调用，但每轮对话都要注入全部工具的 Schema。

### 3.1 单次查询 Token 消耗明细

| 消耗项 | MCP 模式 | CLI 模式 | 说明 |
| :--- | :--- | :--- | :--- |
| **工具 Schema 注入** | ~12,000 Tokens（30+ 工具 Schema 全量注入） | 0 | MCP 每轮对话的固定开销，CLI 无需额外 Schema |
| **第 1 次调用：了解用法** | — | ~2,100 Tokens（`tapd --help` 返回命令参考卡） | CLI 特有，同会话内只需一次 |
| **查询请求参数** | ~40 Tokens | ~20 Tokens | 差异不大 |
| **查询返回结果** | ~520 Tokens（JSON + 原始 HTML） | ~350 Tokens（YAML + Markdown） | CLI 自动将 HTML 转为 MD，并剔除空字段 |
| **单次查询合计** | **~12,560 Tokens** | **~2,480 Tokens** | **CLI 节省约 5 倍** |

**特别注意**：这里 CLI 的数据已经包含了 `tapd --help` 的开销。如果 Agent 在同一会话中多次查询 TAPD，`--help` 只需要调一次，后续查询只消耗约 370 Tokens。

### 3.2 多轮对话场景才是关键

单次查询节省 5 倍，这个数字看起来还行，但不算惊艳。真正拉开差距的是**多轮对话场景**。

我们来算一个更贴近实际的场景：一个 Agent 任务涉及 10 轮对话，其中有 3 次需要调用 TAPD。

| 指标 | MCP 模式 | CLI 模式 | 节省 |
| :--- | :--- | :--- | :--- |
| **Schema 累计开销** | 12,000 × 10 = **120,000 Tokens** | **0** | MCP 每轮都要注入，不管用不用 |
| **工具调用开销** | 560 × 3 = **1,680 Tokens** | 2,100 + 370 × 3 = **3,210 Tokens** | CLI 首次含 `--help`，后续只有命令本身 |
| **合计** | **~121,680 Tokens** | **~3,210 Tokens** | **CLI 节省约 38 倍** |

为什么差距突然从 5 倍跳到 38 倍？因为 MCP 的 Schema 注入是**每轮对话的固定税**——不管这一轮你有没有调用 TAPD 工具，那 12,000 Tokens 的 Schema 都要交。而 CLI 没有这个税。

> **核心结论**：单次查询时 CLI 节省约 5 倍，但随着对话轮数增加，差距会持续放大。在典型的多轮 Agent 工作流中，CLI 可节省 **10-40 倍** Token。

---

## 四、快速上手

### 4.1 安装

一行命令搞定：

```bash
go install github.com/studyzy/tapd-ai-cli/cmd/tapd@latest
```

### 4.2 认证配置

`tapd-ai-cli` 支持两种认证方式，下面的命令二选一：

```bash
# 方式一：使用 Access Token（推荐）
export TAPD_ACCESS_TOKEN=your_token

# 方式二：通过交互式登录持久化凭据
tapd auth login
```

### 4.3 日常使用示例

以下是一些常用场景的命令，复制即用：

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

### 4.4 在 AI Agent 中使用

将 `tapd` 加入 Agent 的可用工具列表后，Agent 会自动通过 `tapd --help` 发现所有功能。不过我更推荐通过 **`tapd skill init`** 一键生成 SKILL.md 指令文件，让 AI Coding 工具主动理解如何使用 tapd CLI：

```bash
# 为 Claude Code、CodeBuddy、Cursor 等 10 种工具一键生成 SKILL.md
tapd skill init
```

命令会自动检测当前目录下已有的 AI Coding 工具配置文件夹并默认选中，交互式确认后将 SKILL.md 生成到对应工具的 `skills/tapd/SKILL.md` 路径下。生成的命令参考部分从当前 CLI 版本的命令树动态生成，始终与实际功能保持一致。

目前支持的工具：Claude Code、CodeBuddy、Cursor、Windsurf、Trae、Codex、Gemini CLI、Cline、Roo Code、Augment。

接入之后，一个典型的 Agent 工作流是这样的：

1. `tapd workspace list` — 发现可用项目
2. `tapd workspace switch <id>` — 锁定目标项目
3. `tapd story list --status=open` — 获取待办需求
4. `tapd story show <id>` — 深入了解某个需求
5. `tapd task create --name=... --story-id=...` — 拆解并创建任务

---

## 五、总结与思考

回顾一下核心要点：

1. **MCP 的最大成本不在调用本身，而在 Schema 注入**。30+ 个工具的 JSON Schema 在每轮对话中全量注入，这是一笔无法回避的"固定税"。
2. **CLI 通过渐进式披露规避了这个问题**。Agent 按需调用 `--help`，只在需要时才获取命令说明，且同会话内不需要重复获取。
3. **输出格式的优化同样重要**。HTML→Markdown 转换、空字段剔除、智能截断，这些细节加起来又能再省 30%-50% 的 Token。
4. **差距随对话轮数放大**。单次查询节省 5 倍，多轮工作流可节省 10-40 倍。

在 AI 编程时代，我个人的一个感受是：我们不需要更复杂的协议，我们需要更直接的工具。`tapd-ai-cli` 是对"AI 友好型基础设施"的一次实践——用最少的 Token 传递最多的信息，用最简单的接口释放最大的编排能力。

如果你也在做 Agent 与项目管理平台的集成，或者正苦恼于 MCP 的高额 Token 成本，欢迎试试 CLI 模式。有任何问题或建议，也欢迎在 GitHub Issues 中讨论。

---

**参考资源：**

- 项目主页：[github.com/studyzy/tapd-ai-cli](https://github.com/studyzy/tapd-ai-cli)
- SDK 仓库：[github.com/studyzy/tapd-sdk-go](https://github.com/studyzy/tapd-sdk-go)

**延伸阅读：**

- Scalekit, *MCP vs CLI: Benchmarking AI Agent Cost & Reliability*, [scalekit.com/blog/mcp-vs-cli-use](https://www.scalekit.com/blog/mcp-vs-cli-use)（含 75 次 benchmark 运行数据和开源测试仓库）
- Jannik Reinhard, *Why CLI Tools Are Beating MCP for AI Agents*, [jannikreinhard.com/2026/02/22/why-cli-tools-are-beating-mcp-for-ai-agents](https://jannikreinhard.com/2026/02/22/why-cli-tools-are-beating-mcp-for-ai-agents/)（企业场景实测，含上下文窗口占用分析）
- Mario Zechner, *MCP vs CLI: Benchmarking Tools for Coding Agents*, [mariozechner.at/posts/2025-08-15-mcp-vs-cli](https://mariozechner.at/posts/2025-08-15-mcp-vs-cli/)（120 次评估运行，"工具设计比协议更重要"）
- Firecrawl, *MCP vs CLI for AI Agents: Which One Should You Use in 2026?*, [firecrawl.dev/blog/mcp-vs-cli](https://www.firecrawl.dev/blog/mcp-vs-cli)（综合分析，含 MCP 安全问题和未来优化方向）
- DeployHQ, *CLIs or MCP for Coding Agents? A Practical Comparison*, [deployhq.com/blog/clis-or-mcp-for-coding-agents-practical-comparison](https://www.deployhq.com/blog/clis-or-mcp-for-coding-agents-practical-comparison)（"CLI 覆盖 80% 日常任务"的实践建议）
