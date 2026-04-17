# 任务: 提取 tapd-sdk-go 独立库

**输入**: 来自 `/specs/002-extract-tapd-sdk/` 的设计文档
**前置条件**: plan.md（必需）, spec.md（用户故事必需）, research.md, data-model.md, contracts/sdk-api.md, quickstart.md

**测试**: 本项目要求测试（spec.md SC-002 要求 100% 测试通过，覆盖率 ≥ 60%）。迁移现有测试至 SDK，CLI 层现有测试保持通过。

**组织结构**: 任务按用户故事分组，以便每个故事能够独立实施和测试。

## 格式: `[ID] [P?] [Story] 描述`
- **[P]**: 可以并行运行（不同文件，无依赖关系）
- **[Story]**: 此任务属于哪个用户故事（US1=SDK 可独立引入, US2=CLI 零回归, US3=SDK 独立认证, US4=清晰公开接口）
- 在描述中包含确切的文件路径

## 路径约定
- **SDK**: `sdk/`（根包 `package tapd`）、`sdk/model/`（数据模型子包）
- **CLI**: `internal/cmd/`、`internal/config/`、`internal/output/`
- **项目根**: `go.mod`

---

## 阶段 1: 设置（SDK 骨架与模块配置）

**目的**: 创建 SDK 目录结构和独立 go.mod，配置 CLI 主模块的本地依赖引用

- [x] T001 创建 SDK 目录结构 `sdk/` 和 `sdk/model/`
- [x] T002 创建 `sdk/go.mod`（module github.com/studyzy/tapd-sdk-go，go 1.24.12，无外部依赖）
- [x] T003 在主 `go.mod` 中添加 `require github.com/studyzy/tapd-sdk-go` 和 `replace github.com/studyzy/tapd-sdk-go => ./sdk` 指令

---

## 阶段 2: 基础（数据模型与客户端迁移）

**目的**: 将 `internal/model/` 和 `internal/client/` 的核心代码迁移到 SDK，这是所有用户故事的阻塞先决条件

**⚠️ 关键**: 在此阶段完成之前，无法开始任何用户故事工作

### 模型迁移（internal/model/ → sdk/model/）

- [x] T004 迁移通用模型到 `sdk/model/model.go`：Workspace、ListResponse、SuccessResponse、CountResponse、TAPDResponse（排除 Config 和 ErrorResponse，它们保留在 CLI 层）
- [x] T005 [P] 迁移 `internal/model/story.go` 到 `sdk/model/story.go`（Story、ListStoriesRequest、CreateStoryRequest、UpdateStoryRequest、CountStoriesRequest）
- [x] T006 [P] 迁移 `internal/model/bug.go` 到 `sdk/model/bug.go`（Bug 及其请求结构体）
- [x] T007 [P] 迁移 `internal/model/task.go` 到 `sdk/model/task.go`（Task 及其请求结构体）
- [x] T008 [P] 迁移 `internal/model/iteration.go` 到 `sdk/model/iteration.go`（Iteration 及其请求结构体）
- [x] T009 [P] 迁移 `internal/model/comment.go` 到 `sdk/model/comment.go`（Comment 及其请求结构体）
- [x] T010 [P] 迁移 `internal/model/wiki.go` 到 `sdk/model/wiki.go`（Wiki 及其请求结构体，含 setOptional 工具函数）
- [x] T011 [P] 迁移 `internal/model/attachment.go` 到 `sdk/model/attachment.go`（Attachment、ImageInfo 及请求结构体）
- [x] T012 [P] 迁移 `internal/model/tcase.go` 到 `sdk/model/tcase.go`（TCase 及其请求结构体）
- [x] T013 [P] 迁移 `internal/model/timesheet.go` 到 `sdk/model/timesheet.go`（Timesheet 及其请求结构体）
- [x] T014 [P] 迁移 `internal/model/category.go` 到 `sdk/model/category.go`
- [x] T015 [P] 迁移 `internal/model/workflow.go` 到 `sdk/model/workflow.go`
- [x] T016 [P] 迁移 `internal/model/release.go` 到 `sdk/model/release.go`
- [x] T017 [P] 迁移 `internal/model/request.go` 到 `sdk/model/request.go`（GetCustomFieldsRequest、WorkspaceIDRequest、WorkflowRequest、GetCommitMsgRequest、GetTodoRequest、GetRelatedBugsRequest、CreateRelationRequest）

### 客户端迁移（internal/client/ → sdk/）

- [x] T018 迁移核心客户端到 `sdk/client.go`：Client 结构体、TAPDError、NewClient、NewClientWithBaseURL、FetchNick、TestAuth、doGet、doPost、doPostJSON、doRequest、mapHTTPError（包声明改为 `package tapd`，import 路径改为 `github.com/studyzy/tapd-sdk-go/model`）
- [x] T019 [P] 迁移 `internal/client/story.go` 到 `sdk/story.go`（移除 htmltomarkdown 导入和 GetStory 中的 HTML→Markdown 转换代码）
- [x] T020 [P] 迁移 `internal/client/bug.go` 到 `sdk/bug.go`（移除 htmltomarkdown 导入和 GetBug 中的 HTML→Markdown 转换代码）
- [x] T021 [P] 迁移 `internal/client/task.go` 到 `sdk/task.go`（移除 htmltomarkdown 导入和 GetTask 中的 HTML→Markdown 转换代码）
- [x] T022 [P] 迁移 `internal/client/wiki.go` 到 `sdk/wiki.go`（移除 htmltomarkdown 导入和 GetWiki 中的 HTML→Markdown 转换代码）
- [x] T023 [P] 迁移 `internal/client/comment.go` 到 `sdk/comment.go`（移除 htmltomarkdown 导入和 ListComments 中的 HTML→Markdown 转换代码）
- [x] T024 [P] 迁移 `internal/client/iteration.go` 到 `sdk/iteration.go`
- [x] T025 [P] 迁移 `internal/client/workspace.go` 到 `sdk/workspace.go`
- [x] T026 [P] 迁移 `internal/client/attachment.go` 到 `sdk/attachment.go`
- [x] T027 [P] 迁移 `internal/client/category.go` 到 `sdk/category.go`
- [x] T028 [P] 迁移 `internal/client/custom_field.go` 到 `sdk/custom_field.go`
- [x] T029 [P] 迁移 `internal/client/misc.go` 到 `sdk/misc.go`（含 GetCommitMsg、ListReleases、GetTodoStories/Tasks/Bugs、SendQiweiMessage）
- [x] T030 [P] 迁移 `internal/client/relation.go` 到 `sdk/relation.go`
- [x] T031 [P] 迁移 `internal/client/tcase.go` 到 `sdk/tcase.go`
- [x] T032 [P] 迁移 `internal/client/timesheet.go` 到 `sdk/timesheet.go`
- [x] T033 [P] 迁移 `internal/client/workflow.go` 到 `sdk/workflow.go`

### SDK 测试迁移

- [x] T034 迁移 `internal/client/client_test.go` 到 `sdk/client_test.go`（包声明改为 `package tapd`，调整 import 路径，TAPDError 相关断言改为包内引用）
- [x] T035 [P] 迁移 `internal/client/story_test.go` 到 `sdk/story_test.go`（HTML→Markdown 测试改为验证原始 HTML 保留）
- [x] T036 [P] 迁移 `internal/client/bug_test.go` 到 `sdk/bug_test.go`（同上调整 HTML 测试）
- [x] T037 [P] 迁移 `internal/client/task_test.go` 到 `sdk/task_test.go`（同上调整 HTML 测试）
- [x] T038 [P] 迁移 `internal/client/iteration_test.go` 到 `sdk/iteration_test.go`
- [x] T039 [P] 迁移 `internal/client/comment_test.go` 到 `sdk/comment_test.go`（HTML 测试调整为验证原始 HTML 保留）
- [x] T040 [P] 迁移 `internal/client/workspace_test.go` 到 `sdk/workspace_test.go`
- [x] T041 [P] 迁移 `internal/client/attachment_test.go` 到 `sdk/attachment_test.go`
- [x] T042 [P] 迁移 `internal/client/category_test.go` 到 `sdk/category_test.go`
- [x] T043 [P] 迁移 `internal/client/misc_test.go` 到 `sdk/misc_test.go`
- [x] T044 [P] 迁移 `internal/client/relation_test.go` 到 `sdk/relation_test.go`
- [x] T045 [P] 迁移 `internal/client/timesheet_test.go` 到 `sdk/timesheet_test.go`
- [x] T046 [P] 迁移 `internal/client/wiki_test.go` 到 `sdk/wiki_test.go`

### SDK 编译验证

- [x] T047 验证 SDK 独立编译：`cd sdk && go build ./...` 和 `go vet ./...`
- [x] T048 验证 SDK 测试通过：`cd sdk && go test ./...`
- [x] T049 验证 SDK 零非标准库依赖：`cd sdk && go list -m all` 确认仅标准库

**检查点**: SDK 独立编译通过、测试通过、零外部依赖 — 可以开始 CLI 重构

---

## 阶段 3: 用户故事 1 — 开发者通过 SDK 集成 TAPD（优先级: P1）🎯 MVP

**目标**: SDK 作为独立包可被第三方 Go 项目引入，无需 CLI 依赖

**独立测试**: 在 SDK 目录下运行 `go build ./...` 和 `go test ./...` 全部通过；`go list -m all` 不包含 cobra、html-to-markdown 等 CLI 依赖

> 此故事在阶段 2 完成后即已满足 — SDK 已经是独立可编译、可测试的模块。阶段 2 的 T047-T049 即为验证。

**检查点**: 用户故事 1 在阶段 2 完成后自动满足

---

## 阶段 4: 用户故事 2 — CLI 作为 SDK 封装层正常工作（优先级: P1）

**目标**: 重构 CLI 层改为消费 SDK，所有现有命令功能和输出与重构前完全一致

**独立测试**: `make test` 全部通过，`make build` 编译成功，`make coverage` 覆盖率 ≥ 60%

### CLI 层类型保留

- [x] T050 [P] [US2] 在 `internal/config/config.go` 中定义本地 Config 结构体（从 internal/model 迁出），移除对 `internal/model` 的 import
- [x] T051 [P] [US2] 在 `internal/output/output.go` 中定义本地 ErrorResponse 结构体，移除对 `internal/model` 的 import；同步更新 `internal/output/output_test.go`

### HTML→Markdown 转换移入 CLI 层

- [x] T052 [US2] 在 `internal/cmd/markdown.go` 中添加 `htmlToMarkdown` 函数（import `html-to-markdown/v2`，空字符串或转换失败时返回原值）

### CLI 命令文件 import 路径更新

- [x] T053 [US2] 更新 `internal/cmd/root.go`：`internal/client` → `github.com/studyzy/tapd-sdk-go`（别名 tapd）、`internal/model` → `github.com/studyzy/tapd-sdk-go/model`、`apiClient *client.Client` → `apiClient *tapd.Client`、`client.NewClient` → `tapd.NewClient`、`appConfig *model.Config` → `appConfig *config.Config`；在 `printComments` 中添加 `htmlToMarkdown` 调用
- [x] T054 [US2] 更新 `internal/cmd/auth.go`：import 路径调整，`client.NewClient` → `tapd.NewClient`
- [x] T055 [P] [US2] 更新 `internal/cmd/story.go`：`internal/model` → `github.com/studyzy/tapd-sdk-go/model`；在 `runStoryShow` 中添加 `story.Description = htmlToMarkdown(story.Description)`
- [x] T056 [P] [US2] 更新 `internal/cmd/bug.go`：import 路径调整；在 `runBugShow` 中添加 `bug.Description = htmlToMarkdown(bug.Description)`
- [x] T057 [P] [US2] 更新 `internal/cmd/task.go`：import 路径调整；在 `runTaskShow` 中添加 `task.Description = htmlToMarkdown(task.Description)`
- [x] T058 [P] [US2] 更新 `internal/cmd/wiki.go`：import 路径调整；在 `runWikiShow` 中添加 `wiki.Description = htmlToMarkdown(wiki.Description)`
- [x] T059 [P] [US2] 更新 `internal/cmd/iteration.go`：`internal/model` → `github.com/studyzy/tapd-sdk-go/model`
- [x] T060 [P] [US2] 更新 `internal/cmd/comment.go`：`internal/model` → `github.com/studyzy/tapd-sdk-go/model`
- [x] T061 [P] [US2] 更新 `internal/cmd/workspace.go`：import 路径调整（model 引用 SDK，Config 引用 config 包）
- [x] T062 [P] [US2] 更新 `internal/cmd/attachment.go`：`internal/model` → `github.com/studyzy/tapd-sdk-go/model`
- [x] T063 [P] [US2] 更新 `internal/cmd/category.go`：`internal/model` → `github.com/studyzy/tapd-sdk-go/model`
- [x] T064 [P] [US2] 更新 `internal/cmd/commit_msg.go`：`internal/model` → `github.com/studyzy/tapd-sdk-go/model`
- [x] T065 [P] [US2] 更新 `internal/cmd/custom_field.go`：`internal/model` → `github.com/studyzy/tapd-sdk-go/model`
- [x] T066 [P] [US2] 更新 `internal/cmd/tcase.go`：`internal/model` → `github.com/studyzy/tapd-sdk-go/model`
- [x] T067 [P] [US2] 更新 `internal/cmd/timesheet.go`：`internal/model` → `github.com/studyzy/tapd-sdk-go/model`
- [x] T068 [P] [US2] 更新 `internal/cmd/relation.go`：`internal/model` → `github.com/studyzy/tapd-sdk-go/model`
- [x] T069 [P] [US2] 更新 `internal/cmd/release.go`：`internal/model` → `github.com/studyzy/tapd-sdk-go/model`
- [x] T070 [P] [US2] 更新 `internal/cmd/workflow.go`：`internal/model` → `github.com/studyzy/tapd-sdk-go/model`
- [x] T071 [P] [US2] 更新 `internal/cmd/qiwei.go`：`internal/model` → `github.com/studyzy/tapd-sdk-go/model`（如有引用）
- [x] T072 [P] [US2] 更新 `internal/cmd/url.go`：`internal/model` → `github.com/studyzy/tapd-sdk-go/model`（如有引用）
- [x] T073 [P] [US2] 更新 `internal/cmd/id_expand.go`：`internal/model` → `github.com/studyzy/tapd-sdk-go/model`（如有引用）

### CLI 测试文件更新

- [x] T074 [P] [US2] 更新 `internal/cmd/integration_test.go`：import 路径调整（client → tapd SDK，model → SDK model）
- [x] T075 [P] [US2] 更新其他 cmd 测试文件（`auth_test.go`、`helpers_test.go`、`help_test.go`、`id_expand_test.go`、`markdown_test.go`、`read_description_test.go`、`url_test.go`）中的 model import 路径

### 删除旧代码与依赖整理

- [x] T076 [US2] 删除 `internal/client/` 整个目录（29 个文件）
- [x] T077 [US2] 删除 `internal/model/` 整个目录（14 个文件）
- [x] T078 [US2] 运行 `go mod tidy`（主模块）和 `cd sdk && go mod tidy`（SDK 模块）

### CLI 编译与测试验证

- [x] T079 [US2] 验证 CLI 编译：`make build`
- [x] T080 [US2] 验证所有测试通过：`make test`
- [x] T081 [US2] 验证代码规范：`make lint`
- [x] T082 [US2] 验证覆盖率 ≥ 60%：`make coverage`

**检查点**: CLI 所有命令功能与重构前完全一致，零回归

---

## 阶段 5: 用户故事 3 — SDK 支持独立认证配置（优先级: P2）

**目标**: SDK 通过代码参数（而非配置文件）配置认证，可在服务端/CI/CD 等无文件系统场景使用

**独立测试**: SDK 测试全部使用代码参数初始化 Client，不依赖 .tapd.json 文件

> 此故事在阶段 2 迁移后已自动满足 — 现有 Client 的 NewClient(accessToken, apiUser, apiPassword) 构造函数本身就是纯代码参数方式，不读取任何配置文件。SDK 中无 config 包依赖。

**检查点**: 用户故事 3 在阶段 2 完成后自动满足

---

## 阶段 6: 用户故事 4 — SDK 提供清晰的公开接口（优先级: P2）

**目标**: SDK 暴露清晰、稳定的公开 API，包含完整数据模型定义和中文文档注释

**独立测试**: `go doc github.com/studyzy/tapd-sdk-go` 输出完整的包文档；所有导出符号有中文注释

> 此故事的大部分要求在迁移过程中已满足（现有代码已有中文注释）。仅需验证和补充。

- [x] T083 [P] [US4] 审查 `sdk/client.go` 所有导出符号的中文文档注释完整性
- [x] T084 [P] [US4] 审查 `sdk/model/*.go` 所有导出结构体和字段的中文文档注释完整性
- [x] T085 [P] [US4] 审查 `sdk/` 下所有资源文件（story.go、bug.go、task.go 等）的导出方法中文文档注释完整性
- [x] T086 [US4] 验证 `go doc` 输出覆盖所有资源类型（≥ 12 种：story、bug、task、iteration、comment、wiki、attachment、relation、timesheet、workflow、workspace、custom_field）

**检查点**: SDK 公开 API 文档完整，所有导出符号有中文注释

---

## 阶段 7: 完善与横切关注点

**目的**: 文档更新、最终验证

- [x] T087 更新 `CODEBUDDY.md` 中的目录结构，添加 `sdk/` 目录说明，更新架构描述
- [x] T088 运行 quickstart.md 中的代码示例验证 SDK 用法正确
- [x] T089 最终全量验证：`cd sdk && go build ./... && go test ./... && go vet ./...` 和 `make build && make test && make lint && make coverage`

---

## 依赖关系与执行顺序

### 阶段依赖关系

- **设置（阶段 1）**: 无依赖关系 — 可立即开始
- **基础（阶段 2）**: 依赖于设置完成 — 阻塞所有用户故事
- **用户故事 1（阶段 3）**: 阶段 2 完成后自动满足
- **用户故事 2（阶段 4）**: 依赖于阶段 2 完成 — CLI 重构的核心工作
- **用户故事 3（阶段 5）**: 阶段 2 完成后自动满足
- **用户故事 4（阶段 6）**: 依赖于阶段 2 完成 — 可与阶段 4 并行
- **完善（阶段 7）**: 依赖于阶段 4 和阶段 6 完成

### 用户故事依赖关系

- **US1（P1）**: 阶段 2 的 T047-T049 即为验证 — 无额外任务
- **US2（P1）**: 阶段 2 完成后可开始 — 核心重构工作量
- **US3（P2）**: 阶段 2 完成后自动满足 — 无额外任务
- **US4（P2）**: 阶段 2 完成后可开始 — 可与 US2 并行

### 每个阶段内部

- 模型迁移（T004-T017）可全部并行（T004 除外，它定义了其他模型依赖的通用类型）
- 客户端迁移（T018-T033）中 T018（核心客户端）必须先完成，其余可并行
- 测试迁移（T034-T046）中 T034（核心测试）建议先完成，其余可并行
- CLI 更新（T050-T075）中 T050-T053 必须先完成（定义本地类型和工具函数），其余可并行

### 并行机会

- 阶段 2 中：T005-T017 全部可并行（模型迁移，不同文件）
- 阶段 2 中：T019-T033 全部可并行（客户端资源文件迁移，需 T018 先完成）
- 阶段 2 中：T035-T046 全部可并行（测试迁移，需 T034 先完成）
- 阶段 4 中：T055-T073 全部可并行（cmd 文件 import 更新，不同文件）
- 阶段 4 与阶段 6 可并行（CLI 重构与文档审查独立）

---

## 并行示例: 阶段 2 模型迁移

```bash
# 一起启动所有模型迁移（T004 先完成后）:
任务: "迁移 internal/model/story.go 到 sdk/model/story.go"
任务: "迁移 internal/model/bug.go 到 sdk/model/bug.go"
任务: "迁移 internal/model/task.go 到 sdk/model/task.go"
# ... 其余模型文件同理
```

## 并行示例: 阶段 4 CLI 更新

```bash
# T050-T053 完成后，一起启动所有 cmd 文件更新:
任务: "更新 internal/cmd/story.go import 路径 + 添加 htmlToMarkdown"
任务: "更新 internal/cmd/bug.go import 路径 + 添加 htmlToMarkdown"
任务: "更新 internal/cmd/task.go import 路径"
# ... 其余 cmd 文件同理
```

---

## 实施策略

### 仅 MVP（用户故事 1 + 2）

1. 完成阶段 1: 设置（SDK 骨架）
2. 完成阶段 2: 基础（模型 + 客户端 + 测试迁移）
3. **验证 US1**: SDK 独立编译、测试通过、零外部依赖
4. 完成阶段 4: CLI 重构（US2）
5. **验证 US2**: `make test` 全部通过，零回归
6. 停止并验证：完整 CLI 功能与重构前一致

### 增量交付

1. 设置 + 基础 → SDK 可用（US1 + US3 满足）
2. CLI 重构 → 零回归（US2 满足）
3. 文档审查 → 接口清晰（US4 满足）
4. 完善 → 文档更新、最终验证

---

## 统计

- **总任务数**: 89
- **US1 任务数**: 0（阶段 2 完成后自动满足）
- **US2 任务数**: 33（T050-T082，CLI 重构核心）
- **US3 任务数**: 0（阶段 2 完成后自动满足）
- **US4 任务数**: 4（T083-T086，文档审查）
- **设置任务数**: 3（T001-T003）
- **基础任务数**: 46（T004-T049，模型+客户端+测试迁移+验证）
- **完善任务数**: 3（T087-T089）
- **并行机会**: 阶段 2 中 ~40 个任务可并行；阶段 4 中 ~20 个任务可并行
- **MVP 范围**: 阶段 1-4（US1 + US2），阶段 5-6 为 P2 增量
