# 任务: TAPD URL 通用查询命令 + Wiki 子命令

**输入**: 来自 `/specs/001-tapd-url-query/` 的设计文档
**前置条件**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

## 格式: `[ID] [P?] [Story?] 描述`
- **[P]**: 可以并行运行（不同文件，无依赖关系）
- **[Story]**: 此任务属于哪个用户故事（US1~US5）
- 描述中包含确切文件路径

## 路径约定
本项目为 Go CLI 工具，源码根目录：
- 命令层：`internal/cmd/`
- 客户端层：`internal/client/`

---

## 阶段 1: 设置

**目的**: 确认现有项目结构就绪，无需新增依赖

- [x] T001 确认 `internal/cmd/` 和 `internal/client/` 目录存在，参考 `story.go`、`bug.go`、`task.go` 作为 wiki 命令实现模板

---

## 阶段 2: 基础（阻塞前置条件）

**目的**: 实现 URL 解析核心函数和 `url` 命令脚手架，所有用户故事分发逻辑依赖此阶段

**⚠️ 关键**: 此阶段完成前，无法开始任何用户故事实现

- [x] T002 在 `internal/cmd/url.go` 中创建文件骨架：定义 `urlCmd` cobra 命令（`Use: "url <url>"`，`Short: "根据 TAPD URL 查询对应条目详情"`），注册到 `rootCmd`，`RunE` 暂留空实现
- [x] T003 在 `internal/cmd/url.go` 中实现 `parseTAPDURL(rawURL string) (workspaceID, entityType, entityID string, err error)` 核心解析函数，支持以下所有格式：
  - 格式一：`dialog_preview_id={type}_{id}` 查询参数（story/bug/task）
  - 格式二：`/tapd_fe/{ws}/{type}/detail/{id}` 路径格式
  - 格式三：`/{ws}/markdown_wikis/show/#{id}` fragment 格式（wiki）
  - WorkspaceID 提取规则：`tapd_fe` 路径取 path[2]，直接路径取 path[1]
- [x] T004 [P] 在 `internal/cmd/url_test.go` 中编写 `parseTAPDURL` 的表格驱动单元测试，覆盖 7 种 URL 格式（Story 详情页、Story 列表预览、Bug 详情页、Bug 列表预览、Task 详情页、Task 看板预览、Wiki fragment）以及无效 URL 和不支持类型的错误场景

**检查点**: `parseTAPDURL` 单元测试全部通过，`go build ./...` 成功 → 可开始用户故事实现

---

## 阶段 3: 用户故事 1、2、3 — Story/Bug/Task URL 查询（优先级: P1）🎯 MVP

**目标**: 实现 Story、Bug、Task 三种类型的 URL 查询分发，复用现有 client 方法，零新增 client 代码

**独立测试**: 执行 `tapd url https://www.tapd.cn/tapd_fe/51081496/story/detail/1151081496001028684` 返回该需求的 JSON 详情；Bug/Task URL 同理

### 实现

- [x] T005 [P] [US1] 在 `internal/cmd/url.go` 的 `runURLQuery` 函数中实现 Story 分发：当 `entityType == "story"` 时，调用 `apiClient.GetStory(workspaceID, entityID, "stories")` 并通过 `output.PrintJSON` 输出结果
- [x] T006 [P] [US2] 在 `internal/cmd/url.go` 的 `runURLQuery` 函数中实现 Bug 分发：当 `entityType == "bug"` 时，调用 `apiClient.GetBug(workspaceID, entityID)` 并输出结果
- [x] T007 [P] [US3] 在 `internal/cmd/url.go` 的 `runURLQuery` 函数中实现 Task 分发：当 `entityType == "task"` 时，调用 `apiClient.GetStory(workspaceID, entityID, "tasks")` 并输出结果

**检查点**: `tapd url <story-url>`、`tapd url <bug-url>`、`tapd url <task-url>` 均返回正确 JSON → US1/US2/US3 可独立验证

---

## 阶段 4: 用户故事 4 — Wiki client + `tapd wiki` 子命令 + Wiki URL 查询（优先级: P2）

**目标**:
1. 新增 Wiki client 方法（`GetWiki`、`ListWikis`）
2. 新增 `tapd wiki` 顶级子命令（`list`、`show`），与 `tapd story`/`tapd bug`/`tapd task` 保持一致的使用体验
3. 实现 Wiki URL 查询分发（`tapd url <wiki-url>`）

**独立测试**:
- `tapd wiki list` 返回当前工作区的 Wiki 文档列表 JSON
- `tapd wiki show 1151081496001001503` 返回指定 Wiki 文档完整内容
- `tapd url "https://www.tapd.cn/51081496/markdown_wikis/show/#1151081496001001503"` 返回该 Wiki 文档 JSON

### Wiki Client 层（阶段 4a）

- [x] T008 在 `internal/client/wiki.go` 中新建文件，实现以下两个方法（参考 `internal/client/bug.go` 实现模式）：
  - `ListWikis(params map[string]string) ([]map[string]interface{}, error)`：调用 `GET /tapd_wikis`，响应包装键为 `TapdWiki`，支持分页参数
  - `GetWiki(workspaceID, id string) (map[string]interface{}, error)`：调用 `GET /tapd_wikis`，参数 `workspace_id` 和 `id`，包装键 `TapdWiki`，`markdown_description` 字段保留原始 Markdown
- [x] T009 [P] 在 `internal/client/wiki_test.go` 中编写单元测试，使用 `httptest` mock 服务器，覆盖：`ListWikis` 成功返回列表、`GetWiki` 成功返回单条、`GetWiki` wiki not found 三个场景

### `tapd wiki` 命令层（阶段 4b，依赖 T008）

- [x] T010 [US4] 在 `internal/cmd/wiki.go` 中新建文件，定义以下命令（参考 `internal/cmd/task.go` 结构）：
  - `wikiCmd`：父命令，`Use: "wiki"`，`Short: "Wiki 文档管理"`
  - `wikiListCmd`：`Use: "list"`，`Short: "查询 Wiki 文档列表"`，支持 `--limit`（默认 10）、`--page`（默认 1）、`--name`（按标题筛选）标志
  - `wikiShowCmd`：`Use: "show <wiki_id>"`，`Short: "查看 Wiki 文档详情"`，`Args: cobra.ExactArgs(1)`
  - 在 `init()` 中将 `wikiCmd` 注册到 `rootCmd`
- [x] T011 [P] [US4] 在 `internal/cmd/wiki.go` 中实现 `runWikiList`：构造 `params`（`workspace_id`、可选 `name`、分页参数），调用 `apiClient.ListWikis(params)`，包装为 `model.ListResponse` 后输出
- [x] T012 [P] [US4] 在 `internal/cmd/wiki.go` 中实现 `runWikiShow`：调用 `apiClient.GetWiki(flagWorkspaceID, args[0])`，通过 `output.PrintJSON` 输出结果

### Wiki URL 分发（阶段 4c，依赖 T008）

- [x] T013 [US4] 在 `internal/cmd/url.go` 的 `runURLQuery` 函数中添加 Wiki 分发：当 `entityType == "wiki"` 时，调用 `apiClient.GetWiki(workspaceID, entityID)` 并输出结果

**检查点**: `tapd wiki list`、`tapd wiki show <id>`、`tapd url <wiki-url>` 均返回正确 JSON → US4 可独立验证

---

## 阶段 5: 用户故事 5 — 错误处理（优先级: P2）

**目标**: 无效 URL 和不支持类型给出清晰错误提示

**独立测试**: 执行 `tapd url "https://github.com/foo"` 输出 `invalid_tapd_url` 错误到 stderr，退出码为 3；执行不支持类型 URL 输出 `unsupported_entity_type` 错误并列出支持类型

### 实现

- [x] T014 [US5] 在 `internal/cmd/url.go` 的 `runURLQuery` 函数中添加错误分发：`parseTAPDURL` 返回错误时调用 `output.PrintError(os.Stderr, "invalid_tapd_url", ...)` 并 `os.Exit(output.ExitParamError)`；`entityType` 不在支持范围时输出 `unsupported_entity_type` 错误，提示 "supported types: story, bug, task, wiki"

**检查点**: 无效 URL 和不支持类型均返回正确错误码和错误消息 → US5 可独立验证

---

## 阶段 6: 完善与横切关注点

**目的**: 文档同步、代码质量验收

- [x] T015 [P] 在 `internal/cmd/spec.go` 中更新工具描述，将 `url` 和 `wiki` 命令加入 spec 输出（AI Agent 自发现）
- [x] T016 [P] 在 `README.md` 的"命令一览"表格和命令树中添加 `url` 和 `wiki` 命令条目
- [x] T017 运行 `make test` 确认全部测试通过，运行 `make build` 确认编译成功，检查覆盖率 ≥ 60%

---

## 依赖关系与执行顺序

### 阶段依赖关系

- **阶段 1（设置）**: 无依赖，立即开始
- **阶段 2（基础）**: 依赖阶段 1 → **阻塞所有用户故事**
- **阶段 3（US1/US2/US3）**: 依赖阶段 2，三者内部可并行
- **阶段 4（US4）**: 依赖阶段 2；T008 完成后 T009/T010/T013 可并行；T011/T012 依赖 T010
- **阶段 5（US5）**: 依赖阶段 3 和阶段 4 完成
- **阶段 6（完善）**: 依赖所有用户故事完成

### 阶段 4 内部依赖

```
T008（wiki.go client）
  ├─ T009 [P]（client 测试，可与 T010 并行）
  ├─ T010（wiki cmd 骨架）
  │    ├─ T011 [P]（runWikiList）
  │    └─ T012 [P]（runWikiShow）
  └─ T013（url.go wiki 分发）
```

### 并行机会

- T003（parseTAPDURL）和 T008（wiki client）可并行开发（不同文件）
- T004（url_test.go）和 T009（wiki_test.go）可并行编写
- T005、T006、T007 逻辑独立，可并行编写
- T011、T012 不同函数，可并行实现
- T015、T016 完善任务可并行

---

## 并行执行示例

```bash
# 阶段 2 内并行
任务: "T003 实现 parseTAPDURL → internal/cmd/url.go"
任务: "T004 URL 解析单元测试 → internal/cmd/url_test.go"
任务: "T008 Wiki client → internal/client/wiki.go"  ← 可提前并行！

# 阶段 3+4 并行（基础完成后）
任务: "T005 Story 分发 → internal/cmd/url.go"
任务: "T006 Bug 分发 → internal/cmd/url.go"
任务: "T007 Task 分发 → internal/cmd/url.go"
任务: "T009 Wiki client 测试 → internal/client/wiki_test.go"
任务: "T010 tapd wiki 命令骨架 → internal/cmd/wiki.go"

# T010 完成后并行
任务: "T011 runWikiList → internal/cmd/wiki.go"
任务: "T012 runWikiShow → internal/cmd/wiki.go"
任务: "T013 Wiki URL 分发 → internal/cmd/url.go"
```

---

## 实施策略

### 仅 MVP（US1/US2/US3，P1 故事）

1. 完成阶段 1：设置
2. 完成阶段 2：基础（`parseTAPDURL` + 命令骨架）
3. 完成阶段 3：Story/Bug/Task URL 查询
4. **停止并验证**：三种类型 URL 均可正常查询
5. 可立即发布使用

### 完整交付

1. 设置 + 基础 → 解析器就绪
2. US1/US2/US3 → MVP，可发布
3. US4 → `tapd wiki` 子命令 + Wiki URL 查询
4. US5 → 错误处理完善
5. 完善阶段 → 文档同步，测试验收

---

## 新增文件清单

| 文件 | 说明 |
|------|------|
| `internal/client/wiki.go` | `ListWikis` + `GetWiki` 方法 |
| `internal/client/wiki_test.go` | Wiki client 单元测试 |
| `internal/cmd/url.go` | `url` 命令 + `parseTAPDURL` + 四类型分发 |
| `internal/cmd/url_test.go` | URL 解析表格驱动单元测试 |
| `internal/cmd/wiki.go` | `tapd wiki list/show` 命令 |

## 注意事项

- `parseTAPDURL` 是核心，建议优先完成并通过单元测试再开始分发逻辑
- `dialog_preview_id` 参数值带类型前缀（`story_`/`bug_`/`task_`），需去掉前缀提取纯数字 ID
- Wiki URL 的 ID 在 `#fragment` 中，需用 `u.Fragment` 而非路径或查询参数提取
- `tapd_fe` 路径的 workspaceID 在 path[2]，直接路径在 path[1]，注意区分
- `GetWiki` 响应包装键为 `TapdWiki`（参考 TAPD API 文档），参考 `GetBug` 实现
- `tapd wiki` 命令只需 `list` 和 `show`（只读操作），与 Wiki 是文档性质匹配
- 错误输出到 `os.Stderr`，正常输出到 `os.Stdout`，与现有命令保持一致
