---
name: tapd
description: TAPD操作，涉及需求，缺陷，任务，Wiki等。tapd.cn快捷查询：tapd url <url> 
---

面向 AI Agent 的 TAPD 命令行工具。通过 `tapd` 命令与 TAPD 平台交互，所有输出针对最小 token 消耗优化。

## 安装

```bash
go install github.com/studyzy/tapd-ai-cli/cmd/tapd@latest
```

## 认证

```bash
# Access Token（推荐）
export TAPD_ACCESS_TOKEN=<your_token>

# 或交互式登录持久化凭据
tapd auth login
```

凭据优先级：CLI flags > 环境变量 > `./.tapd.json` > `~/.tapd.json`

## 命令参考

{{COMMAND_REFERENCE}}
