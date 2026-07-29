# tapd-ai-cli

面向 AI Agent 的 TAPD 命令行工具，通过 TAPD Open API 实现项目管理核心操作。

## 安装

### 方式一：直接下载二进制（推荐，无需 Go 环境）

从 [Releases](https://github.com/studyzy/tapd-ai-cli/releases) 下载对应平台的预编译包，或复制以下脚本一键安装：

**Linux / macOS**

```bash
# 自动检测系统与架构，下载并安装到 /usr/local/bin
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
VER=$(curl -s https://api.github.com/repos/studyzy/tapd-ai-cli/releases/latest | grep tag_name | cut -d '"' -f 4)
curl -sSLo /tmp/tapd.tar.gz "https://github.com/studyzy/tapd-ai-cli/releases/download/${VER}/tapd_${VER}_${OS}_${ARCH}.tar.gz"
sudo tar xzf /tmp/tapd.tar.gz -C /usr/local/bin tapd && rm /tmp/tapd.tar.gz
```

**Windows (PowerShell)**

```powershell
# 下载 zip 并解压到当前目录
$ver = (Invoke-RestMethod https://api.github.com/repos/studyzy/tapd-ai-cli/releases/latest).tag_name
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
Invoke-WebRequest -Uri "https://github.com/studyzy/tapd-ai-cli/releases/download/$ver/tapd_${ver}_windows_${arch}.zip" -OutFile tapd.zip
Expand-Archive tapd.zip -DestinationPath . ; Remove-Item tapd.zip
```

预编译包支持 **Linux / macOS / Windows**，**amd64 / arm64** 架构。

### 方式二：go install

```bash
go install github.com/studyzy/tapd-ai-cli/cmd/tapd@latest
```

### 方式三：从源码构建并安装

```bash
git clone git@github.com:studyzy/tapd-ai-cli.git
cd tapd-ai-cli
make install   # 编译并安装到 $GOPATH/bin
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

# 查询发布评审列表
tapd launch list

# 通过 URL 查询任意条目（需求/缺陷/任务/Wiki）
tapd url https://www.tapd.cn/tapd_fe/51081496/story/detail/1151081496001028684

# 查询 Wiki 文档列表
tapd wiki list

# 查看所有命令参考（AI 自发现）
tapd --help
```

## 高级过滤（--filter）

所有 `list` 命令支持 `--filter` 标志，可重复使用，直接透传 TAPD OpenAPI 的高级查询语法：

```bash
# 按名称模糊搜索
tapd story list --filter "name=LIKE<登录>"

# 按自定义字段精确匹配
tapd story list --filter "custom_field_one=EQ<高优先级>"

# 按时间范围查询
tapd bug list --filter "created=>2024-01-01" --filter "created=<2024-12-31"

# 组合多个过滤条件
tapd task list --owner zhangsan --filter "status=CONTAINS_OR<开发中|测试中>"

# 多人查询
tapd story list --filter "owner=USER_OR<张三|李四>"
```

支持的操作符：`LIKE`（模糊）、`EQ`（精确）、`NOT_EQ`（不等于）、`LIKE_OR`（多值模糊 OR）、`CONTAINS`（包含所有值 AND）、`CONTAINS_OR`（包含任一值 OR）、`USER_OR`（多人 OR）、`>`/`<`（时间/数值比较）、`~`（时间范围）、`<>`（不等于简写）、`|`（多值 OR）。

适用于所有标准字段和自定义字段（`custom_field_*`），可与已有标志（`--status`、`--owner` 等）组合使用。

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
├── launch    list | count | create | update <id> | templates | fields
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
