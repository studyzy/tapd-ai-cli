# 任务: tapd-ai-cli MVP

**输入**: 来自 `/specs/001-mvp-tapd-cli/` 的设计文档
**前置条件**: plan.md(必需), spec.md(必需), research.md, data-model.md, contracts/cli-commands.md, quickstart.md

**测试**: 项目章程要求核心代码必须有 UT，覆盖率 ≥60%。测试任务包含在各阶段中。

**组织结构**: 任务按用户故事分组，以便每个故事能够独立实施和测试。

## 格式: `[ID] [P?] [Story] 描述`
- **[P]**: 可以并行运行(不同文件, 无依赖关系)
- **[Story]**: 此任务属于哪个用户故事(US1, US2, US3, US4, US5, US6)
- 在描述中包含确切的文件路径

## 路径约定
- 仓库根目录下的标准 Go 项目布局: `cmd/`, `internal/`, `main.go`

---

## 阶段 1: 设置(项目初始化)

**目的**: Go 项目初始化和基本结构

- [x] T001 运行 `go mod init github.com/studyzy/tapd-ai-cli` 初始化 Go module，创建 go.mod
- [x] T002 创建 main.go 入口文件，调用 `cmd.Execute()`
- [x] T003 [P] 创建 Makefile，包含 build、test、lint、coverage 目标
- [x] T004 [P] 创建 README.md，包含项目简介、安装方式、认证配置和基本用法（中文）

---

## 阶段 2: 基础(阻塞前置条件)

**目的**: 所有用户故事依赖的核心基础设施

**⚠️ 关键**: 在此阶段完成之前, 无法开始任何用户故事工作

- [x] T005 在 internal/model/model.go 中定义所有数据模型结构体（Config, Workspace, Story, Task, Bug, Iteration, ListResponse, ErrorResponse），所有 JSON 字段使用 `omitempty` 标签，导出符号添加中文注释
- [x] T006 在 internal/output/output.go 中实现 JSON 输出工具函数：PrintJSON（标准/紧凑模式）、PrintError（错误输出到 stderr，含 error code + message + hint）、PrintSuccess（成功响应），以及退出码常量定义（0=成功，1=认证错误，2=未找到，3=参数错误，4=API错误）
- [x] T007 [P] 在 internal/output/output_test.go 中为 output 包编写单元测试：测试 PrintJSON 紧凑/非紧凑模式、PrintError 输出到 stderr、omitempty 行为、退出码常量
- [x] T008 在 internal/config/config.go 中实现配置管理：LoadConfig（按优先级链加载凭据：CLI flags > 环境变量 > ./.tapd.json > ~/.tapd.json，同层级 access_token 优先于 api_user/api_password）、SaveConfig（写入指定路径的 .tapd.json）、GetWorkspaceID（从配置或 --workspace-id 标志获取）
- [x] T009 [P] 在 internal/config/config_test.go 中为 config 包编写单元测试：测试优先级链（环境变量覆盖文件、当前目录覆盖主目录、token 优先于 user/password）、SaveConfig 写入和读取、缺失凭据时的错误信息
- [x] T010 在 internal/client/client.go 中实现 TAPD API HTTP 客户端：Client 结构体（baseURL, httpClient, authHeader）、NewClient（根据 access_token 或 api_user/api_password 构建 Authorization 头：Bearer 或 Basic）、doRequest 通用请求方法（GET/POST，处理 TAPD 响应包装格式 `{"status":1,"data":[...],"info":"success"}`，解包实体类型 wrapper key）、错误码映射（401→退出码1，404→退出码2，422→退出码3，429/500/502→退出码4）
- [x] T011 [P] 在 internal/client/client_test.go 中为 client 基础方法编写单元测试：使用 httptest.NewServer 模拟 TAPD API，测试 Bearer Token 认证头、Basic Auth 认证头、成功响应解析、各 HTTP 错误码的错误映射
- [x] T012 在 cmd/root.go 中实现根命令：注册全局 PersistentFlags（--workspace-id, --compact, --access-token, --api-user, --api-password），在 PersistentPreRun 中初始化 config 和 client，注册所有子命令

**检查点**: 基础就绪 — `go build` 成功，`go test ./internal/...` 全部通过

---

## 阶段 3: 用户故事 1 - 认证与凭据管理 (优先级: P1) 🎯 MVP

**目标**: 用户可以通过 Access Token 或 API User/Password 登录，凭据持久化到 .tapd.json

**独立测试**: 执行 `tapd auth login --access-token <token>` 后检查 ~/.tapd.json 内容正确；设置环境变量后直接调用其他命令可自动认证

### 实施

- [x] T013 [US1] 在 cmd/auth.go 中实现 `tapd auth` 父命令和 `tapd auth login` 子命令：支持 --access-token 或 --api-user/--api-password（二选一），支持 --local 标志控制写入路径，login 成功时调用 `GET /quickstart/testauth` 验证凭据有效性，输出 `{"success":true}` 或错误信息
- [x] T014 [P] [US1] 在 cmd/auth_test.go 中为 auth 命令编写单元测试：测试 --access-token 登录、--api-user/--api-password 登录、--local 标志写入当前目录、缺少参数时的错误提示、凭据无效时的错误处理

**检查点**: `tapd auth login` 命令可正常工作，凭据正确持久化

---

## 阶段 4: 用户故事 2 - 项目列表与工作区切换 (优先级: P1)

**目标**: 用户可以列出参与的项目、切换工作区、查看工作区信息

**独立测试**: 执行 `tapd workspace list` 返回项目列表，`tapd workspace switch <id>` 后 ./.tapd.json 中 workspace_id 正确写入

### 实施

- [x] T015 [US2] 在 internal/client/workspace.go 中实现工作区相关 API 方法：ListWorkspaces（GET /workspaces/projects，过滤 category=organization 条目）、GetWorkspaceInfo（GET /workspaces/show）
- [x] T016 [P] [US2] 在 internal/client/workspace_test.go 中为工作区 API 方法编写单元测试：使用 httptest 模拟响应，测试列表过滤、详情获取、错误处理
- [x] T017 [US2] 在 cmd/workspace.go 中实现 `tapd workspace` 父命令及 list、switch、info 三个子命令：list 调用 ListWorkspaces 并格式化输出，switch 将 workspace_id 写入当前目录 .tapd.json（不存在则创建），info 调用 GetWorkspaceInfo 输出详情

**检查点**: 工作区命令完整可用，`tapd workspace switch` 正确写入当前目录配置

---

## 阶段 5: 用户故事 3 - 需求与任务管理 (优先级: P1)

**目标**: 用户可以对需求和任务执行 list/show/create/update/count 操作

**独立测试**: 执行 `tapd story list` 返回精简列表，`tapd story create --name "test"` 返回成功 JSON 含 id 和 url

### 实施

- [x] T018 [US3] 在 internal/client/story.go 中实现需求/任务共用的 API 方法：ListStories（GET /stories 或 /tasks，支持 status/owner/iteration_id/limit/page 参数，返回精简字段列表）、GetStory（GET /stories 或 /tasks，按 id 查询，description 经 HTML→Markdown 转换）、CreateStory（POST /stories 或 /tasks）、UpdateStory（POST /stories 或 /tasks）、CountStories（GET /stories/count 或 /tasks/count），所有方法通过 entityType 参数区分 stories/tasks
- [x] T019 [P] [US3] 安装 `github.com/JohannesKaufmann/html-to-markdown/v2` 依赖，在 internal/client/story.go 的 GetStory 方法中实现 description 字段的 HTML→Markdown 转换逻辑
- [x] T020 [P] [US3] 在 internal/client/story_test.go 中为需求/任务 API 方法编写单元测试：测试 list 精简字段、show 详情含 HTML→Markdown 转换、create 成功响应、update 成功响应、count 响应、entity_type=tasks 时的行为
- [x] T021 [US3] 在 cmd/story.go 中实现 `tapd story` 父命令及 list、show、create、update、count 五个子命令：list 支持 --status/--owner/--iteration-id/--limit/--page 标志，show 接受 story_id 位置参数，create 支持 --name(必需)/--description/--owner/--priority/--iteration-id 标志，update 接受 story_id 位置参数及 --name/--status/--owner/--priority 标志，count 支持 --status 标志，列表输出包含 total/has_more 分页信息
- [x] T022 [US3] 在 cmd/task.go 中实现 `tapd task` 父命令及 list、show、create、update、count 五个子命令：复用 story 命令的逻辑，entity_type 固定为 "tasks"，create 额外支持 --story-id 标志关联需求

**检查点**: 需求和任务 CRUD 完整可用，`tapd story list` 和 `tapd task list` 均正确输出

---

## 阶段 6: 用户故事 4 - 缺陷管理 (优先级: P2)

**目标**: 用户可以对缺陷执行 list/show/create/update/count 操作

**独立测试**: 执行 `tapd bug list` 返回缺陷列表，`tapd bug create --title "test"` 返回成功 JSON

### 实施

- [x] T023 [US4] 在 internal/client/bug.go 中实现缺陷相关 API 方法：ListBugs（GET /bugs，支持 status/priority/severity/limit/page 参数）、GetBug（GET /bugs，按 id 查询，description 经 HTML→Markdown 转换）、CreateBug（POST /bugs）、UpdateBug（POST /bugs）、CountBugs（GET /bugs/count），注意 Bug 使用 title 而非 name 字段
- [x] T024 [P] [US4] 在 internal/client/bug_test.go 中为缺陷 API 方法编写单元测试：测试 list、show（含 HTML→Markdown）、create、update、count，以及 Bug 特有字段（severity/priority_label）
- [x] T025 [US4] 在 cmd/bug.go 中实现 `tapd bug` 父命令及 list、show、create、update、count 五个子命令：list 支持 --status/--priority/--severity/--limit/--page 标志，create 支持 --title(必需)/--description/--priority/--severity 标志，update 接受 bug_id 位置参数及 --title/--status/--priority/--severity 标志

**检查点**: 缺陷 CRUD 完整可用

---

## 阶段 7: 用户故事 5 - 迭代查询 (优先级: P2)

**目标**: 用户可以查询项目迭代列表

**独立测试**: 执行 `tapd iteration list` 返回迭代列表（id、name、startdate、enddate、status）

### 实施

- [x] T026 [US5] 在 internal/client/iteration.go 中实现迭代 API 方法：ListIterations（GET /iterations，支持 status 参数）
- [x] T027 [P] [US5] 在 internal/client/iteration_test.go 中为迭代 API 方法编写单元测试
- [x] T028 [US5] 在 cmd/iteration.go 中实现 `tapd iteration` 父命令及 list 子命令：支持 --status 标志

**检查点**: 迭代查询可用

---

## 阶段 8: 用户故事 6 - AI 自发现能力 (优先级: P3)

**目标**: `tapd spec` 命令输出 Tool Definition JSON，使 AI Agent 可自动发现所有命令

**独立测试**: 执行 `tapd spec` 输出合法 JSON，包含所有已注册命令的参数定义

### 实施

- [x] T029 [US6] 在 cmd/spec.go 中实现 `tapd spec` 命令：遍历 Cobra 命令树，为每个叶子命令生成 OpenAI/Anthropic 兼容的 Tool Definition（包含 name、description、parameters schema），输出 JSON 数组
- [x] T030 [P] [US6] 在 cmd/spec_test.go 中为 spec 命令编写单元测试：验证输出为合法 JSON 数组，每个 tool 定义包含必需字段（name, description, parameters），参数 schema 与实际命令标志一致

**检查点**: `tapd spec` 输出完整的 Tool Definition JSON

---

## 阶段 9: 完善与横切关注点

**目的**: 跨用户故事的质量改进

- [x] T031 [P] 运行 `go test ./... -coverprofile=coverage.out` 检查覆盖率，补充测试使 internal/ 包覆盖率达到 ≥60%
- [x] T032 [P] 运行 `go vet ./...` 和 `goimports -w .` 确保代码通过静态检查和格式化
- [x] T033 确认所有导出的函数、结构体、接口都有中文注释，包注释使用中文描述包功能
- [x] T034 运行 quickstart.md 中的验证流程，确认端到端完整链路可用

---

## 依赖关系与执行顺序

### 阶段依赖关系

- **设置(阶段 1)**: 无依赖关系 — 可立即开始
- **基础(阶段 2)**: 依赖于设置完成 — 阻塞所有用户故事
- **US1 认证(阶段 3)**: 依赖于基础(阶段 2)
- **US2 工作区(阶段 4)**: 依赖于 US1 完成（需要认证才能调用 API）
- **US3 需求/任务(阶段 5)**: 依赖于 US2 完成（需要 workspace_id）
- **US4 缺陷(阶段 6)**: 依赖于 US2 完成（需要 workspace_id），可与 US3 并行
- **US5 迭代(阶段 7)**: 依赖于 US2 完成（需要 workspace_id），可与 US3/US4 并行
- **US6 AI 自发现(阶段 8)**: 依赖于所有其他 US 完成（需遍历完整命令树）
- **完善(阶段 9)**: 依赖于所有用户故事完成

### 用户故事依赖关系

```
US1 (认证) ─→ US2 (工作区) ─→ US3 (需求/任务)
                             ├→ US4 (缺陷)       ← 可并行
                             └→ US5 (迭代)       ← 可并行
                                                    ↓
                                               US6 (spec)
```

### 每个用户故事内部

- 模型定义（阶段 2 统一完成）在 API 方法之前
- API 方法（internal/client/）在命令实现（cmd/）之前
- 测试可与 API 方法并行编写

### 并行机会

- T003, T004 可与 T001/T002 并行
- T007, T009, T011 测试任务可各自与实现并行
- 阶段 5（US3）、阶段 6（US4）、阶段 7（US5）可在 US2 完成后并行
- 阶段 9 的 T031、T032 可并行

---

## 并行示例: 阶段 2 基础

```bash
# 先完成有依赖的任务:
任务: T005 "定义数据模型"
任务: T006 "实现输出工具"
任务: T008 "实现配置管理"
任务: T010 "实现 HTTP 客户端"
任务: T012 "实现根命令"

# 测试任务可与对应实现并行:
任务: T007 "output 测试" (与 T006 并行)
任务: T009 "config 测试" (与 T008 并行)
任务: T011 "client 测试" (与 T010 并行)
```

## 并行示例: US2 完成后并行推进 US3 + US4 + US5

```bash
# 三个独立的用户故事可并行:
开发者 A: US3 (T018→T019→T020→T021→T022) — 需求/任务
开发者 B: US4 (T023→T024→T025) — 缺陷
开发者 C: US5 (T026→T027→T028) — 迭代
```

---

## 实现策略

### 仅 MVP(用户故事 1-3)

1. 完成阶段 1: 设置
2. 完成阶段 2: 基础
3. 完成阶段 3: US1 认证
4. 完成阶段 4: US2 工作区
5. 完成阶段 5: US3 需求/任务
6. **停止并验证**: 可以完整运行 `login → workspace switch → story list → story create`

### 增量交付

1. 设置 + 基础 → 基础就绪
2. US1 认证 → 可登录
3. US2 工作区 → 可切换项目
4. US3 需求/任务 → **核心 MVP 可用**
5. US4 缺陷 → 扩展能力
6. US5 迭代 → 扩展能力
7. US6 spec → AI 自发现
8. 完善 → 质量达标

---

## 注意事项

- [P] 任务 = 不同文件, 无依赖关系
- [Story] 标签将任务映射到特定用户故事以实现可追溯性
- 每个用户故事应该独立可完成和可测试
- 在每个任务或逻辑组后提交
- 在任何检查点停止以独立验证故事
- 避免: 模糊任务, 相同文件冲突, 破坏独立性的跨故事依赖
- 所有代码注释和文档使用中文，错误信息字符串使用英文
