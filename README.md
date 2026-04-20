# tapd-ai-cli

面向 AI Agent 的 TAPD 命令行工具，通过 TAPD Open API 实现项目管理核心操作。

## 安装

### 方式一：go install（推荐）

```bash
go install github.com/studyzy/tapd-ai-cli/cmd/tapd@latest
```

### 方式二：从源码构建并安装

```bash
git clone git@github.com:studyzy/tapd-ai-cli.git
cd tapd-ai-cli
make install   # 编译并安装到 $GOPATH/bin
```

### 方式三：仅构建二进制

```bash
git clone git@github.com:studyzy/tapd-ai-cli.git
cd tapd-ai-cli
make build     # 在当前目录生成 ./tapd
```

## 认证

支持两种认证方式：

### Access Token（推荐）

```bash
# 命令行登录
tapd auth login --access-token <your_token>

# 或设置环境变量
export TAPD_ACCESS_TOKEN=<your_token>
```

### API User/Password

```bash
# 命令行登录
tapd auth login --api-user <user> --api-password <password>

# 或设置环境变量
export TAPD_API_USER=<user>
export TAPD_API_PASSWORD=<password>
```

凭据也可以写入配置文件 `~/.tapd.json` 或当前目录的 `.tapd.json`。

**凭据优先级**：CLI flags > 环境变量 > `./.tapd.json` > `~/.tapd.json`

### 自定义 TAPD 站点地址

如需连接非 `tapd.cn` 的 TAPD 部署，可通过环境变量或配置文件指定：

```bash
# 环境变量
export TAPD_API_BASE_URL=https://api.my-tapd.com
export TAPD_BASE_URL=https://www.my-tapd.com
```

或写入配置文件：

```json
{
  "access_token": "your-token",
  "api_base_url": "https://api.my-tapd.com",
  "base_url": "https://www.my-tapd.com"
}
```

| 配置项 | 环境变量 | JSON 字段 | 默认值 |
|--------|----------|-----------|--------|
| API 地址 | `TAPD_API_BASE_URL` | `api_base_url` | `https://api.tapd.cn` |
| 前端地址 | `TAPD_BASE_URL` | `base_url` | `https://www.tapd.cn` |

## 基本用法

```bash
# 查看参与的项目
tapd workspace list

# 切换工作区
tapd workspace switch 12345

# 查询需求列表
tapd story list

# 创建需求
tapd story create --name "新功能需求"

# 更新需求并切换迭代
tapd story update 10001 --iteration-id 12345

# 查询缺陷列表
tapd bug list

# 查询任务列表
tapd task list

# 查看迭代列表
tapd iteration list

# 通过 URL 查询任意条目（需求/缺陷/任务/Wiki）
tapd url https://www.tapd.cn/tapd_fe/51081496/story/detail/1151081496001028684

# 查询 Wiki 文档列表
tapd wiki list

# 查看所有命令参考（AI 自发现）
tapd --help
```

## 命令一览

```
tapd
├── auth      login --access-token <token> | --api-user <user> --api-password <pwd> [--local]
├── workspace list | switch <id> | info
├── story     list | show <id> | create | update <id> | count | todo
├── task      list | show <id> | create | update <id> | count | todo
├── bug       list | show <id> | create | update <id> | count | todo
├── wiki      list | show <id> | create | update
├── iteration list | create | update | count
├── comment   list | add | update | count
├── tcase     list | create | batch-create
├── timesheet list | add | update
├── workflow  transitions | status-map | last-steps
├── relation  bugs | create
├── skill     init
├── url       <tapd-url>
└── ...       release, attachment, image, category, custom-field, story-field, workitem-type, commit-msg, qiwei
```

## AI Coding 工具集成

`tapd skill init` 可一键为主流 AI Coding 工具生成 TAPD CLI 的 SKILL.md 指令文件：

```bash
tapd skill init
```

支持的工具：Claude Code、CodeBuddy、Cursor、Windsurf、Trae、Codex、Gemini CLI、Cline、Roo Code、Augment。

命令会自动检测当前目录下已有的工具配置文件夹并默认选中，交互式确认后生成 SKILL.md。生成的命令参考部分从当前 CLI 版本的命令树动态生成，始终保持同步。

## 全局标志

| 标志 | 说明 |
|------|------|
| `--workspace-id <id>` | 指定工作区 ID（覆盖本地配置） |
| `--pretty` | 输出格式化 JSON（带缩进，便于人类阅读；默认输出紧凑 JSON 以节省 token） |

## SDK

TAPD Go SDK 已独立为单独的模块，可直接引入使用：

```bash
go get github.com/studyzy/tapd-sdk-go@latest
```

详见 [tapd-sdk-go](https://github.com/studyzy/tapd-sdk-go)。

## 开发

```bash
make build      # 构建
make install    # 安装到 $GOPATH/bin
make test       # 运行测试
make coverage   # 测试覆盖率报告
make lint       # 代码检查
make fmt        # 代码格式化
make clean      # 清理构建产物
```

## 许可证

Apache License 2.0
