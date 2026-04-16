# tapd-ai-cli

面向 AI Agent 的 TAPD 命令行工具，通过 TAPD Open API 实现项目管理核心操作。

## 安装

```bash
go install github.com/studyzy/tapd-ai-cli/cmd/tapd@latest
```

或从源码构建：

```bash
git clone git@github.com:studyzy/tapd-ai-cli.git
cd tapd-ai-cli
make build
```

## 认证

支持两种认证方式：

### Access Token（推荐）

```bash
# 命令行登录
./tapd auth login --access-token <your_token>

# 或设置环境变量
export TAPD_ACCESS_TOKEN=<your_token>
```

### API User/Password

```bash
# 命令行登录
./tapd auth login --api-user <user> --api-password <password>

# 或设置环境变量
export TAPD_API_USER=<user>
export TAPD_API_PASSWORD=<password>
```

凭据也可以写入配置文件 `~/.tapd.json` 或当前目录的 `.tapd.json`。

## 基本用法

```bash
# 查看参与的项目
./tapd workspace list

# 切换工作区
./tapd workspace switch 12345

# 查询需求
./tapd story list

# 创建需求
./tapd story create --name "新功能需求"

# 查询缺陷
./tapd bug list

# 获取 Tool Definition（AI 自发现）
./tapd spec
```

## 全局标志

- `--workspace-id <id>` — 指定工作区 ID（覆盖本地配置）
- `--compact` — 输出紧凑 JSON（节省 token）

## 许可证

Apache License 2.0
