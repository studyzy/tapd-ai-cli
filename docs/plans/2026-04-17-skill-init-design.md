# `tapd skill init` 命令设计文档

## 概述

新增 `tapd skill init` 子命令，用于在当前项目目录下为各 AI Coding 工具生成 TAPD CLI 的 SKILL.md 文件。SKILL.md 告诉 AI Coding 工具如何使用 tapd CLI，包括安装方式、认证配置和完整的命令参考。

命令参考部分从 Cobra 命令树动态生成，确保与当前二进制版本的命令始终一致。

## 命令接口

```
tapd skill init
```

- 无需 TAPD API 认证（纯本地操作）
- 无需 `--workspace-id`
- 交互式多选菜单，无额外标志

## 支持的 AI Coding 工具

| 工具名称 | 检测目录 | SKILL.md 生成路径 |
|---------|---------|-----------------|
| Claude Code | `.claude/` | `.claude/skills/tapd/SKILL.md` |
| CodeBuddy | `.codebuddy/` | `.codebuddy/skills/tapd/SKILL.md` |
| Cursor | `.cursor/` | `.cursor/skills/tapd/SKILL.md` |
| Windsurf | `.windsurf/` | `.windsurf/skills/tapd/SKILL.md` |
| Trae | `.trae/` | `.trae/skills/tapd/SKILL.md` |
| Codex | `.codex/` | `.codex/skills/tapd/SKILL.md` |
| Gemini CLI | `.gemini/` | `.gemini/skills/tapd/SKILL.md` |
| Cline | `.cline/` | `.cline/skills/tapd/SKILL.md` |
| Roo Code | `.roo/` | `.roo/skills/tapd/SKILL.md` |
| Augment | `.augment/` | `.augment/skills/tapd/SKILL.md` |

## 交互流程

1. 扫描当前目录，检测上述 10 个工具的配置目录是否存在
2. 显示多选菜单，已存在目录的工具默认选中
3. 用户通过输入编号切换选中状态，输入 `y` 确认，输入 `q` 取消
4. 为每个选中的工具生成 `skills/tapd/SKILL.md`
5. 输出成功结果（JSON 格式，列出已生成的文件路径）

交互示例：

```
Select AI coding tools to initialize TAPD skill:
  [x] 1. Claude Code (.claude/)
  [x] 2. CodeBuddy (.codebuddy/)
  [ ] 3. Cursor (.cursor/)
  [ ] 4. Windsurf (.windsurf/)
  [ ] 5. Trae (.trae/)
  [ ] 6. Codex (.codex/)
  [ ] 7. Gemini CLI (.gemini/)
  [ ] 8. Cline (.cline/)
  [ ] 9. Roo Code (.roo/)
  [ ] 10. Augment (.augment/)
Toggle: enter number | Confirm: y | Cancel: q
> 3
  [x] 1. Claude Code (.claude/)
  [x] 2. CodeBuddy (.codebuddy/)
  [x] 3. Cursor (.cursor/)
  ...
> y
Generated TAPD skill for 3 tools:
  .claude/skills/tapd/SKILL.md
  .codebuddy/skills/tapd/SKILL.md
  .cursor/skills/tapd/SKILL.md
```

## SKILL.md 生成方式

### 固定部分（嵌入模板）

通过 `//go:embed skill_template.md` 嵌入到二进制中，包含：

- frontmatter（name、description）
- 标题和简介
- 安装说明
- 认证说明和凭据优先级
- 全局标志表格
- 占位符 `{{COMMAND_REFERENCE}}`（运行时替换）
- 典型 Agent 工作流示例

### 动态部分（命令参考）

复用 `help.go` 中已有的函数：

- `buildSpecLines(rootCmd)` — 遍历命令树生成命令参考行
- `commandToLine()` — 将单个 Cobra 命令格式化为紧凑参考文本

新增函数 `buildCommandReferenceMarkdown(rootCmd) string`：

1. 调用 `buildSpecLines(rootCmd)` 获取 `[]specLine`
2. 按 group 分组，每组生成一个 `### group` 标题 + ` ```bash ` 代码块
3. 拼接为完整的 Markdown 命令参考文本

### 生成流程

```
skill_template.md (embed) + buildCommandReferenceMarkdown(rootCmd)
    ↓ strings.Replace("{{COMMAND_REFERENCE}}", dynamicContent)
    ↓ os.MkdirAll + os.WriteFile
    → .xxx/skills/tapd/SKILL.md
```

## 文件变更清单

### 新增文件

| 文件 | 说明 |
|------|------|
| `internal/cmd/skill.go` | `skill` 父命令 + `skill init` 子命令实现 |
| `internal/cmd/skill_template.md` | SKILL.md 的固定部分模板（通过 embed 嵌入） |
| `internal/cmd/skill_test.go` | 单元测试 |

### 修改文件

| 文件 | 变更 |
|------|------|
| `internal/cmd/root.go` | `PersistentPreRunE` 中排除 `skill init`（跳过认证） |

## 实现细节

### 1. skill.go 核心结构

```go
package cmd

import (
    _ "embed"
    "os"
    "strings"
    // ...
)

//go:embed skill_template.md
var skillTemplate []byte

// aiTool 定义一个 AI Coding 工具的元信息
type aiTool struct {
    Name    string // 显示名称
    Dir     string // 配置目录名，如 ".claude"
}

// supportedTools 支持的 AI Coding 工具列表
var supportedTools = []aiTool{
    {Name: "Claude Code", Dir: ".claude"},
    {Name: "CodeBuddy", Dir: ".codebuddy"},
    {Name: "Cursor", Dir: ".cursor"},
    {Name: "Windsurf", Dir: ".windsurf"},
    {Name: "Trae", Dir: ".trae"},
    {Name: "Codex", Dir: ".codex"},
    {Name: "Gemini CLI", Dir: ".gemini"},
    {Name: "Cline", Dir: ".cline"},
    {Name: "Roo Code", Dir: ".roo"},
    {Name: "Augment", Dir: ".augment"},
}

var skillCmd = &cobra.Command{
    Use:   "skill",
    Short: "AI Coding 工具 Skill 管理",
}

var skillInitCmd = &cobra.Command{
    Use:   "init",
    Short: "为 AI Coding 工具生成 TAPD CLI 的 SKILL.md",
    RunE:  runSkillInit,
}
```

### 2. 认证跳过

在 `root.go` 的 `PersistentPreRunE` 中，`init` 命令与 `login` 一样跳过认证：

```go
if cmd.Name() == "login" || cmd.Name() == "init" {
    return nil
}
```

**注意**：当前判断逻辑是 `cmd.Name() == "login"`，按相同模式追加 `"init"` 即可。考虑到 `init` 名称较通用，可以进一步收紧判断条件为 `cmd.Parent() != nil && cmd.Parent().Name() == "skill"`，但现阶段只有 `skill init` 使用 `init` 名称，简单追加即可。如果未来有其他 `init` 子命令需要认证，再调整。

### 3. buildCommandReferenceMarkdown

```go
// buildCommandReferenceMarkdown 从命令树动态生成 Markdown 格式的命令参考
func buildCommandReferenceMarkdown(root *cobra.Command) string {
    lines := buildSpecLines(root)
    var b strings.Builder
    lastGroup := ""
    for _, l := range lines {
        if l.group != lastGroup {
            if lastGroup != "" {
                b.WriteString("```\n\n")
            }
            // 获取分组命令的 Short 描述
            b.WriteString("### " + l.group + " — " + getGroupShort(root, l.group) + "\n\n")
            b.WriteString("```bash\n")
            lastGroup = l.group
        }
        b.WriteString(l.text + "\n")
    }
    if lastGroup != "" {
        b.WriteString("```\n")
    }
    return b.String()
}
```

### 4. 交互逻辑

使用标准库 `bufio.Scanner` + `fmt.Print` 实现终端交互，不引入第三方依赖：

- 检测目录存在性：`os.Stat(dir)` 判断
- 切换选中：维护 `[]bool` 数组
- 输入解析：数字切换、`y` 确认、`q` 取消
- 每次输入后重新打印菜单

### 5. 文件写入

```go
for _, tool := range selectedTools {
    dir := filepath.Join(cwd, tool.Dir, "skills", "tapd")
    os.MkdirAll(dir, 0755)
    os.WriteFile(filepath.Join(dir, "SKILL.md"), content, 0644)
}
```

## 测试计划

### skill_test.go

| 测试函数 | 验证内容 |
|---------|---------|
| `TestBuildCommandReferenceMarkdown` | 从 rootCmd 生成的命令参考包含已知命令，格式正确 |
| `TestDetectExistingTools` | 目录存在/不存在时检测结果正确 |
| `TestGenerateSkillContent` | 模板占位符被正确替换，输出包含固定部分和动态命令参考 |
| `TestWriteSkillFiles` | 文件写入到正确路径，内容正确，中间目录自动创建 |

### 手动验证

```bash
# 在有 .claude/ 和 .codebuddy/ 的项目目录下执行
tapd skill init
# 验证：菜单中 Claude Code 和 CodeBuddy 默认选中
# 确认后检查生成的 SKILL.md 内容与命令参考一致
```

## 依赖

- 无新的第三方依赖
- 使用标准库：`embed`、`os`、`path/filepath`、`strings`、`bufio`、`fmt`
