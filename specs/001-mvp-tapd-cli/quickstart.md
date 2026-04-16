# 快速开始: tapd-ai-cli

## 前提条件

- Go 1.24+ 已安装
- 已在 TAPD 平台获取 API 凭据（api_user / api_password）

## 构建

```bash
git clone git@github.com:studyzy/tapd-ai-cli.git
cd tapd-ai-cli
go build -o tapd .
```

## 认证

### 方式一：Access Token 登录（推荐）

```bash
# 凭据保存到 ~/.tapd.json
./tapd auth login --access-token your_token

# 或保存到当前目录 .tapd.json（项目级配置）
./tapd auth login --access-token your_token --local
```

### 方式二：API User/Password 登录

```bash
./tapd auth login --api-user your_user --api-password your_password

# 或保存到当前目录
./tapd auth login --api-user your_user --api-password your_password --local
```

### 方式三：环境变量

```bash
# Access Token（推荐）
export TAPD_ACCESS_TOKEN=your_token

# 或 API User/Password
export TAPD_API_USER=your_user
export TAPD_API_PASSWORD=your_password
```

### 方式四：配置文件

手动创建 `~/.tapd.json` 或 `./.tapd.json`：

```json
{
  "access_token": "your_token"
}
```

或使用 User/Password：

```json
{
  "api_user": "your_user",
  "api_password": "your_password"
}
  "api_password": "your_password",
  "workspace_id": "12345"
}
```

## 基本使用流程

```bash
# 1. 查看参与的项目
./tapd workspace list

# 2. 切换到目标工作区
./tapd workspace switch 12345

# 3. 查询需求列表
./tapd story list

# 4. 查看需求详情
./tapd story show 1000001

# 5. 创建新需求
./tapd story create --name "新功能需求" --description "需求描述" --priority High

# 6. 查询缺陷
./tapd bug list --status new --priority high

# 7. 创建缺陷
./tapd bug create --title "发现的问题" --severity serious
```

## AI Agent 集成

```bash
# 获取 Tool Definition JSON
./tapd spec

# AI Agent 可直接加载输出，自动发现所有可用命令
```

## 紧凑模式

```bash
# 输出紧凑 JSON（节省 token）
./tapd story list --compact
```

## 验证

成功完成以下流程即表示 MVP 功能正常：

1. `tapd auth login` → 输出 `{"success":true}`
2. `tapd workspace list` → 输出项目列表 JSON
3. `tapd workspace switch <id>` → 输出 `{"success":true,...}`
4. `tapd story list` → 输出需求列表 JSON
5. `tapd story create --name "test"` → 输出含 id 和 url 的成功 JSON
