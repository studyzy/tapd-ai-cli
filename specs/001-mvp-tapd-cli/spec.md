# 功能规范: tapd-ai-cli MVP

**功能分支**: `001-mvp-tapd-cli`
**创建时间**: 2026-04-16
**状态**: 草稿
**输入**: 用户描述: "参考 docs/requirement.md 和 TAPD MCP 功能，实现 tapd-ai-cli 的 MVP 版本"

## 用户场景与测试 *(必填)*

### 用户故事 1 - 认证与凭据管理 (优先级: P1)

AI Agent 需要通过多种方式提供 TAPD 认证信息，系统支持两种认证模式（Access Token 和 API User/Password），按优先级从多个来源自动解析凭据，无需每次重复传入。

认证模式：
- **Access Token（推荐）**: 使用 `TAPD_ACCESS_TOKEN`，HTTP 请求头为 `Authorization: Bearer <token>`
- **API User/Password**: 使用 `api_user` + `api_password`，HTTP 请求头为 `Authorization: Basic <base64>`

凭据解析优先级（从高到低，Token 优先于 User/Password）：
1. 命令行参数 `--access-token` 或 `--api-user` / `--api-password`
2. 环境变量 `TAPD_ACCESS_TOKEN` 或 `TAPD_API_USER` / `TAPD_API_PASSWORD`
3. 当前目录下的 `.tapd.json` 配置文件
4. 用户主目录下的 `~/.tapd.json` 配置文件

在同一来源层级中，access_token 优先于 api_user/api_password。

**优先级原因**: 认证是所有 API 操作的前提条件，不完成认证则其他功能无法使用。

**独立测试**: 可以通过环境变量、本地配置文件、auth login 命令分别验证凭据是否被正确识别。

**验收场景**:

1. **给定** 用户有有效的 TAPD API 凭据, **当** 执行 `tapd auth login --api-user xxx --api-password yyy` 时, **那么** 凭据被持久化到 `~/.tapd.json`，命令输出 `{"success":true}`
2. **给定** 用户执行 `tapd auth login --access-token <token>` 时, **那么** token 被持久化到 `~/.tapd.json`，命令输出 `{"success":true}`
3. **给定** 用户执行 `tapd auth login --access-token <token> --local` 时, **那么** token 被持久化到当前目录的 `.tapd.json`
4. **给定** 环境变量 `TAPD_ACCESS_TOKEN` 已设置, **当** 执行任意需要认证的命令时, **那么** 系统自动使用 Bearer Token 完成认证
5. **给定** 环境变量同时设置了 `TAPD_ACCESS_TOKEN` 和 `TAPD_API_USER`/`TAPD_API_PASSWORD`, **当** 执行命令时, **那么** 系统优先使用 access_token
6. **给定** 当前目录存在 `.tapd.json` 且 `~/.tapd.json` 也存在, **当** 两者的凭据不同时, **那么** 系统优先使用当前目录的 `.tapd.json`
7. **给定** 用户未提供任何凭据来源, **当** 执行需要认证的命令时, **那么** stderr 输出凭据缺失的提示信息，列出所有支持的凭据配置方式，退出码为 1

---

### 用户故事 2 - 查看参与的项目列表并切换工作区 (优先级: P1)

AI Agent 需要查看当前用户参与的所有 TAPD 项目，并选择一个作为当前工作区，避免后续每条命令都重复指定 workspace_id。

**优先级原因**: 工作区选择是操作需求、缺陷等实体的前置步骤。

**独立测试**: 可以通过列出项目、切换工作区后验证配置文件中 workspace_id 是否正确写入。

**验收场景**:

1. **给定** 用户已登录, **当** 执行 `tapd workspace list` 时, **那么** 输出用户参与的项目列表（JSON 格式，包含 id、name、status 等字段），自动过滤掉 category 为 organization 的条目
2. **给定** 用户已登录, **当** 执行 `tapd workspace switch <id>` 时, **那么** workspace_id 被写入当前目录的 `.tapd.json`（不存在则自动创建），输出 `{"success":true,"workspace_id":"<id>"}`
3. **给定** 用户已登录, **当** 执行 `tapd workspace info` 时, **那么** 输出当前工作区的详细信息（id、名称、状态、创建时间等）

---

### 用户故事 3 - 需求与任务管理 (优先级: P1)

AI Agent 需要查询、创建和更新 TAPD 项目中的需求（stories）和任务（tasks），这是与 TAPD 平台交互的核心场景。

**优先级原因**: 需求和任务管理是 TAPD 最核心的功能，也是 AI Agent 最常用的操作。

**独立测试**: 可以通过查询需求列表、创建新需求、更新需求状态来完整验证。

**验收场景**:

1. **给定** 用户已切换到某工作区, **当** 执行 `tapd story list [--status <status>] [--owner <owner>] [--limit N]` 时, **那么** 输出精简的需求列表（默认返回 id、name、status、owner、modified 字段），默认 limit 为 10，超出时提示剩余数量
2. **给定** 用户已切换到某工作区, **当** 执行 `tapd story show <id>` 时, **那么** 输出需求完整详情，description 字段中的 HTML 内容被转换为 Markdown
3. **给定** 用户已切换到某工作区, **当** 执行 `tapd story create --name <title> [--description <desc>] [--owner <owner>] [--priority <priority>]` 时, **那么** 创建需求并输出 `{"success":true,"id":"<id>","url":"<url>"}`
4. **给定** 用户已切换到某工作区, **当** 执行 `tapd story update <id> [--name <title>] [--status <status>] [--owner <owner>]` 时, **那么** 更新需求并输出更新后的数据
5. **给定** 用户已切换到某工作区, **当** 执行 `tapd story count [--status <status>]` 时, **那么** 输出符合条件的需求数量
6. **给定** 用户已切换到某工作区, **当** 执行 `tapd task list/show/create/update/count` 时, **那么** 与 story 命令行为一致，entity_type 参数自动设为 tasks

---

### 用户故事 4 - 缺陷管理 (优先级: P2)

AI Agent 需要查询、创建和更新 TAPD 项目中的缺陷（bugs），支持按状态、优先级、严重程度等条件筛选。

**优先级原因**: 缺陷管理是 TAPD 的高频使用场景，且与需求管理的操作模式高度相似，边际实现成本低。

**独立测试**: 可以通过查询缺陷列表、创建新缺陷、更新缺陷状态来完整验证。

**验收场景**:

1. **给定** 用户已切换到某工作区, **当** 执行 `tapd bug list [--status <status>] [--priority <priority>] [--severity <severity>] [--limit N]` 时, **那么** 输出精简的缺陷列表
2. **给定** 用户已切换到某工作区, **当** 执行 `tapd bug show <id>` 时, **那么** 输出缺陷完整详情
3. **给定** 用户已切换到某工作区, **当** 执行 `tapd bug create --title <title> [--description <desc>] [--priority <priority>] [--severity <severity>]` 时, **那么** 创建缺陷并输出 `{"success":true,"id":"<id>","url":"<url>"}`
4. **给定** 用户已切换到某工作区, **当** 执行 `tapd bug update <id> [--title <title>] [--status <status>]` 时, **那么** 更新缺陷并输出更新后的数据
5. **给定** 用户已切换到某工作区, **当** 执行 `tapd bug count [--status <status>]` 时, **那么** 输出符合条件的缺陷数量

---

### 用户故事 5 - 迭代查询 (优先级: P2)

AI Agent 需要查询 TAPD 项目中的迭代（iterations），以便了解项目进度和将需求/任务关联到对应迭代。

**优先级原因**: 迭代是组织需求和任务的重要维度，对理解项目上下文有帮助。

**独立测试**: 可以通过查询迭代列表并验证返回数据格式来独立验证。

**验收场景**:

1. **给定** 用户已切换到某工作区, **当** 执行 `tapd iteration list [--status <status>]` 时, **那么** 输出迭代列表（id、name、startdate、enddate、status 等字段）

---

### 用户故事 6 - AI 自发现能力 (优先级: P3)

AI Agent 需要通过一条命令获取该 CLI 工具的完整功能描述（Tool Definition JSON），以便自动理解可用的命令和参数，实现自发现能力。

**优先级原因**: 这是 AI-First 设计的关键特性，使 AI Agent 无需预置知识即可使用本工具。

**独立测试**: 可以通过执行 spec 命令并验证输出是否为合法的 Tool Definition JSON 来独立验证。

**验收场景**:

1. **给定** CLI 已安装, **当** 执行 `tapd spec` 时, **那么** 输出 OpenAI/Anthropic 兼容的 Tool Definition JSON，包含所有可用命令、参数及描述
2. **给定** CLI 已安装, **当** 工具定义被加载到 AI Agent 中时, **那么** AI Agent 能根据定义正确调用 CLI 命令

---

### 边界情况

- 当 API 凭据过期或无效时，所有命令返回明确的认证错误（退出码 1）并提示重新登录或检查环境变量
- 当多个凭据来源同时存在时，严格按优先级使用最高优先级的来源，不做合并
- 当 workspace_id 未设置且命令未通过 `--workspace-id` 传入时，提示用户先切换工作区
- 当查询结果为空时，返回空数组 `[]` 而非错误
- 当 TAPD API 返回 HTTP 错误（如 429 限流、500 服务端错误）时，输出结构化错误信息到 stderr
- 当需求/缺陷的 description 包含 HTML 标签时，正确转换为 Markdown
- 当使用 `--compact` 标志时，JSON 输出移除所有缩进和多余空白

## 需求 *(必填)*

### 功能需求

- **FR-001**: 系统必须支持两种认证模式：Access Token（Bearer）和 API User/Password（Basic Auth），同一层级中 access_token 优先于 api_user/api_password
- **FR-001a**: 凭据解析优先级：命令行参数 > 环境变量（`TAPD_ACCESS_TOKEN` 或 `TAPD_API_USER`/`TAPD_API_PASSWORD`）> 当前目录 `.tapd.json` > 用户主目录 `~/.tapd.json`
- **FR-001b**: `auth login` 命令支持 `--access-token` 或 `--api-user`/`--api-password` 两种方式，默认写入 `~/.tapd.json`，使用 `--local` 标志时写入当前目录的 `.tapd.json`
- **FR-002**: 系统必须支持查看用户参与的项目列表，自动过滤 category 为 organization 的条目
- **FR-003**: 系统必须支持切换当前工作区，将 workspace_id 写入当前目录的 `.tapd.json`（不存在则自动创建），使同目录下后续命令自动使用该 workspace
- **FR-004**: 系统必须支持查询当前工作区的详细信息
- **FR-005**: 系统必须支持对需求（stories）的列表查询、详情查看、创建、更新和计数操作
- **FR-006**: 系统必须支持对任务（tasks）的列表查询、详情查看、创建、更新和计数操作
- **FR-007**: 系统必须支持对缺陷（bugs）的列表查询、详情查看、创建、更新和计数操作
- **FR-008**: 系统必须支持查询迭代列表
- **FR-009**: 所有数据命令输出必须为 JSON 格式，struct 字段使用 `omitempty` 标签
- **FR-010**: 列表命令默认返回 10 条记录，超出时在输出中提示剩余数量和分页参数
- **FR-011**: 需求/缺陷的 description 字段中的 HTML 内容必须转换为 Markdown
- **FR-012**: 所有命令必须支持 `--workspace-id` 全局标志覆盖本地配置的 workspace_id
- **FR-013**: 系统必须支持 `--compact` 全局标志输出紧凑 JSON（无缩进无多余空白）
- **FR-014**: 错误信息必须输出到 stderr，包含可操作的修复建议
- **FR-015**: 系统必须使用明确的退出码（0=成功，1=认证错误，2=未找到，3=参数错误，4=API错误）
- **FR-016**: 系统必须支持 `spec` 命令输出 OpenAI/Anthropic 兼容的 Tool Definition JSON

### 关键实体

- **Workspace（工作区）**: TAPD 项目，包含 id、name、status、creator、created 等属性
- **Story（需求）**: TAPD 需求/工作项，包含 id、name、status、owner、priority、description、iteration_id 等属性。entity_type 为 "stories"
- **Task（任务）**: TAPD 任务，与 Story 共用 API（entity_type 为 "tasks"），状态仅有 open/progressing/done 三种
- **Bug（缺陷）**: TAPD 缺陷，包含 id、title、status、priority、severity、description 等属性
- **Iteration（迭代）**: TAPD 迭代，包含 id、name、startdate、enddate、status 等属性

## 成功标准 *(必填)*

### 可衡量的结果

- **SC-001**: AI Agent 可以通过 5 条以内的命令完成"登录 → 切换工作区 → 查询需求 → 创建需求 → 验证创建成功"的完整流程
- **SC-002**: 列表查询命令的 JSON 输出在紧凑模式下，单条记录的字段数不超过 10 个，确保 token 消耗最小化
- **SC-003**: 所有命令的 JSON 输出可被标准 JSON 解析器正确解析，无格式错误
- **SC-004**: 当发生错误时，stderr 输出的信息足以让 AI Agent 自主判断下一步操作（包含错误类型和建议动作）
- **SC-005**: `tapd spec` 输出的 Tool Definition 可被 AI Agent 直接加载并正确调用对应命令
- **SC-006**: CLI 功能覆盖 TAPD MCP Server 的核心操作（项目查看、需求 CRUD、任务 CRUD、缺陷 CRUD、迭代查询）

## 假设

- TAPD Open API 支持两种认证方式：Bearer Token（`Authorization: Bearer <token>`）和 Basic Auth（`Authorization: Basic <base64(user:password)>`）
- 用户在使用本工具前已在 TAPD 平台获取了 Access Token 或 API User/Password 凭据
- 单用户使用场景，本地配置文件无需处理并发读写
- `.tapd.json` 配置文件格式统一，无论在当前目录还是用户主目录，结构相同
- TAPD API 的返回格式稳定，字段不会频繁变更
- MVP 阶段不支持自定义字段（custom_field_*）查询，此功能留待后续版本
- MVP 阶段不支持 Wiki、工时、测试用例等高级模块
