// Package config 管理 tapd-ai-cli 的凭据加载与持久化，支持四级优先链：
// CLI flags > 环境变量 > 当前目录 .tapd.json > 用户主目录 ~/.tapd.json
package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// LoadConfig 按优先级加载配置：环境变量 > ./.tapd.json > ~/.tapd.json
// 同一来源中 access_token 优先于 api_user/api_password
func LoadConfig() (*model.Config, error) {
	cfg := &model.Config{}

	// 尝试从 ~/.tapd.json 加载
	homePath, err := GetConfigPath(false)
	if err == nil {
		if homeCfg, e := readConfigFile(homePath); e != nil {
			return nil, e
		} else if homeCfg != nil {
			cfg = homeCfg
		}
	}

	// 尝试从 ./.tapd.json 加载（优先级高于 home）
	localPath, _ := GetConfigPath(true)
	if localCfg, e := readConfigFile(localPath); e != nil {
		return nil, e
	} else if localCfg != nil {
		// 当前目录有凭据则覆盖
		if localCfg.AccessToken != "" || localCfg.APIUser != "" {
			cfg.AccessToken = localCfg.AccessToken
			cfg.APIUser = localCfg.APIUser
			cfg.APIPassword = localCfg.APIPassword
		}
		if localCfg.WorkspaceID != "" {
			cfg.WorkspaceID = localCfg.WorkspaceID
		}
	}

	// 环境变量优先级最高
	envToken := os.Getenv("TAPD_ACCESS_TOKEN")
	envUser := os.Getenv("TAPD_API_USER")
	envPassword := os.Getenv("TAPD_API_PASSWORD")
	envWorkspace := os.Getenv("TAPD_WORKSPACE_ID")

	if envToken != "" || envUser != "" {
		cfg.AccessToken = envToken
		cfg.APIUser = envUser
		cfg.APIPassword = envPassword
	}
	if envWorkspace != "" {
		cfg.WorkspaceID = envWorkspace
	}

	return cfg, nil
}

// SaveConfig 将配置写入指定路径的 JSON 文件，自动创建父目录，文件权限 0600
func SaveConfig(cfg *model.Config, filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0600)
}

// GetConfigPath 返回配置文件路径：local=true 返回 ./.tapd.json，否则返回 ~/.tapd.json
func GetConfigPath(local bool) (string, error) {
	if local {
		return ".tapd.json", nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".tapd.json"), nil
}

// SaveWorkspaceID 将 workspace_id 保存到当前目录的 .tapd.json，保留已有的其他字段
func SaveWorkspaceID(workspaceID string) error {
	path, _ := GetConfigPath(true)
	cfg := &model.Config{}
	if existing, err := readConfigFile(path); err != nil {
		return err
	} else if existing != nil {
		cfg = existing
	}
	cfg.WorkspaceID = workspaceID
	return SaveConfig(cfg, path)
}

// readConfigFile 读取并解析指定路径的 .tapd.json 配置文件
// 文件不存在时返回 (nil, nil)，解析失败时返回错误
func readConfigFile(path string) (*model.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	cfg := &model.Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
