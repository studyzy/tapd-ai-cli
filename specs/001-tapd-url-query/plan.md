# 实施计划: TAPD URL 通用查询命令

**分支**: `001-tapd-url-query` | **日期**: 2026-04-16 | **规范**: [spec.md](./spec.md)
**输入**: 来自 `/specs/001-tapd-url-query/spec.md` 的功能规范

## 摘要

新增 `tapd url <url>` 顶级命令，支持 Story、Bug、Task、Wiki 四种 TAPD 条目类型，兼容详情页 URL、列表/看板预览 URL（含 `dialog_preview_id` 参数）和 Wiki fragment URL 三种格式。解析后复用现有 `GetStory`/`GetBug` 方法，Task 通过 `GetStory(..., "tasks")` 处理，Wiki 新增 `GetWiki` client 方法。

## 技术背景

- **语言/版本**: Go 1.24
- **主要依赖**: `spf13/cobra`（CLI 框架）、标准库 `net/url`（URL 解析）
- **存储**: N/A（无状态）
- **测试**: `go test`，表格驱动测试
- **目标平台**: Linux / macOS CLI
- **项目类型**: 命令行工具（CLI）
- **性能目标**: 单次查询响应 < 10 秒（网络正常）
- **约束条件**: 仅新增一个外部依赖（无），复用现有 client 方法为主
- **规模/范围**: 2 个新源文件（`cmd/url.go`、`client/wiki.go`）+ 2 个测试文件

## 章程检查

| 原则 | 符合？ | 说明 |
|------|--------|------|
| I. API 优先 | ✅ | 复用 `GetStory`/`GetBug`，新增 `GetWiki` 调用 `/tapd_wikis` 端点 |
| II. AI 优化输出 | ✅ | 输出 JSON，支持 `--pretty`，错误到 stderr，含可操作建议 |
| III. Go 编码规范 | ✅ | gofmt、驼峰命名、error 最后返回，函数 < 80 行 |
| IV. 测试纪律 | ✅ | URL 解析和 Wiki client 均有单元测试，目标覆盖率 ≥ 60% |
| V. 中文文档与注释 | ✅ | 所有导出符号和注释使用中文 |

## 项目结构

### 文档（此功能）

```
specs/001-tapd-url-query/
├── plan.md              # 此文件
├── research.md          # 阶段 0 输出
├── data-model.md        # 阶段 1 输出
├── contracts/
│   └── cli-contract.md  # 阶段 1 输出
└── tasks.md             # 阶段 2 输出（/speckit.tasks 生成）
```

### 源代码（仓库根目录）

```
internal/
├── client/
│   ├── wiki.go          # 新增：GetWiki 方法，调用 /tapd_wikis API
│   └── wiki_test.go     # 新增：Wiki client 单元测试
└── cmd/
    ├── url.go           # 新增：url 命令 + parseTAPDURL 解析函数
    └── url_test.go      # 新增：URL 解析表格驱动单元测试
```

**结构决策**: 采用单一项目结构，与现有 `story.go`/`bug.go` 和 `client/bug.go` 模式完全一致。

## URL 解析策略

```
输入 URL
  │
  ├─ 检查 dialog_preview_id 查询参数
  │    ├─ story_{id}  → type=story, id=...
  │    ├─ bug_{id}    → type=bug,   id=...
  │    └─ task_{id}   → type=task,  id=...
  │
  ├─ 检查路径关键字（tapd_fe 或直接路径）
  │    ├─ /story/detail/{id}     → type=story
  │    ├─ /bug/detail/{id}       → type=bug
  │    ├─ /task/detail/{id}      → type=task
  │    └─ /markdown_wikis/show/  → type=wiki, id=fragment
  │
  └─ 无法识别 → 错误：invalid_tapd_url

WorkspaceID 提取：
  - tapd_fe 路径: path[2]（/tapd_fe/{ws}/...）
  - 直接路径:     path[1]（/{ws}/...）
```
