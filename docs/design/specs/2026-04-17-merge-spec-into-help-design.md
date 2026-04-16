# 将 spec 子命令合并到 --help

## 背景

`tapd spec` 子命令输出一张紧凑的命令参考卡，供 AI Agent 发现所有可用命令及参数。但 `spec` 这个命令名不够直观——AI Agent 和人类用户都会自然地先尝试 `tapd --help`。

## 目标

- 将紧凑参考卡输出合并为 `tapd --help` 的默认输出
- 移除 `spec` 子命令
- 子命令的 `--help`（如 `tapd story --help`）保持 Cobra 默认行为不变

## 设计

### 变更清单

1. **`root.go`**：在 `init()` 中调用 `rootCmd.SetHelpFunc(customRootHelp)`，自定义函数判断：若当前命令是 rootCmd 则输出紧凑参考卡，否则 fallback 到 Cobra 默认 help。
2. **`spec.go`**：删除 `specCmd` 变量、`init()` 中的 `rootCmd.AddCommand(specCmd)`、`runSpec` 函数。保留所有工具函数（`buildSpecLines`、`walkSpecCommands`、`commandToLine`、`printSpecOutput` 等），供 help 函数复用。
3. **`root.go` 的 `PersistentPreRunE`**：移除 `cmd.Name() == "spec"` 的跳过判断。`--help` 在 Cobra 中天然在 PreRun 之前拦截，不会触发认证。
4. **`spec_test.go`**：删除 `TestSpecCommand_Exists`；修改 `TestBuildSpecLines` 移除对 `"tapd spec"` 的断言。其余辅助函数测试保留。

### 不变项

- 紧凑参考卡的输出格式和内容不变
- 子命令 `--help` 保持 Cobra 默认
- `buildSpecLines`、`walkSpecCommands`、`commandToLine`、`printSpecOutput` 等函数保留在 `spec.go` 中

### 效果

```
$ tapd --help
tapd - 面向 AI Agent 的 TAPD 命令行工具
Global: [--workspace-id=<id>] [--json] [--pretty] [--no-comments]

# auth
tapd auth login [--access-token] [--api-user] [--api-password] [--local]  # 登录并持久化凭据

# workspace
tapd workspace list  # 列出参与的项目
...

$ tapd story --help
(Cobra 默认格式，不变)
```

## 测试计划

- [ ] `go test ./...` 全部通过
- [ ] `tapd --help` 输出紧凑参考卡
- [ ] `tapd story --help` 保持 Cobra 默认格式
- [ ] `tapd spec` 返回 unknown command 错误
