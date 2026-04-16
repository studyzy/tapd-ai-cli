package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/studyzy/tapd-ai-cli/internal/config"
	"github.com/studyzy/tapd-ai-cli/internal/model"
)

// writeConfigFile 将 Config 写入指定路径的 JSON 文件
func writeConfigFile(t *testing.T, path string, cfg *model.Config) {
	t.Helper()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
}

// chdirTemp 切换到临时目录并在测试结束后恢复
func chdirTemp(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(orig)
	})
}

func TestLoadConfig_EnvAccessToken(t *testing.T) {
	tmp := t.TempDir()
	chdirTemp(t, tmp)
	t.Setenv("TAPD_ACCESS_TOKEN", "token123")
	t.Setenv("TAPD_API_USER", "")
	t.Setenv("TAPD_API_PASSWORD", "")
	t.Setenv("TAPD_WORKSPACE_ID", "")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AccessToken != "token123" {
		t.Errorf("expected access_token=token123, got %s", cfg.AccessToken)
	}
}

func TestLoadConfig_EnvAPIUserPassword(t *testing.T) {
	tmp := t.TempDir()
	chdirTemp(t, tmp)
	t.Setenv("TAPD_ACCESS_TOKEN", "")
	t.Setenv("TAPD_API_USER", "user1")
	t.Setenv("TAPD_API_PASSWORD", "pass1")
	t.Setenv("TAPD_WORKSPACE_ID", "")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIUser != "user1" || cfg.APIPassword != "pass1" {
		t.Errorf("expected api_user=user1, api_password=pass1, got %s/%s", cfg.APIUser, cfg.APIPassword)
	}
}

func TestLoadConfig_EnvAccessTokenOverAPIUser(t *testing.T) {
	tmp := t.TempDir()
	chdirTemp(t, tmp)
	t.Setenv("TAPD_ACCESS_TOKEN", "token_wins")
	t.Setenv("TAPD_API_USER", "user_loses")
	t.Setenv("TAPD_API_PASSWORD", "pass_loses")
	t.Setenv("TAPD_WORKSPACE_ID", "")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	// 同一来源中 access_token 和 api_user 都存在，两者都会被设置
	// 但 access_token 存在即表示使用 token 认证
	if cfg.AccessToken != "token_wins" {
		t.Errorf("expected access_token=token_wins, got %s", cfg.AccessToken)
	}
}

func TestLoadConfig_FromLocalFile(t *testing.T) {
	tmp := t.TempDir()
	chdirTemp(t, tmp)
	t.Setenv("TAPD_ACCESS_TOKEN", "")
	t.Setenv("TAPD_API_USER", "")
	t.Setenv("TAPD_API_PASSWORD", "")
	t.Setenv("TAPD_WORKSPACE_ID", "")

	writeConfigFile(t, filepath.Join(tmp, ".tapd.json"), &model.Config{
		AccessToken: "file_token",
		WorkspaceID: "ws123",
	})

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AccessToken != "file_token" {
		t.Errorf("expected access_token=file_token, got %s", cfg.AccessToken)
	}
	if cfg.WorkspaceID != "ws123" {
		t.Errorf("expected workspace_id=ws123, got %s", cfg.WorkspaceID)
	}
}

func TestLoadConfig_LocalFileOverHomeFile(t *testing.T) {
	tmp := t.TempDir()
	chdirTemp(t, tmp)
	t.Setenv("TAPD_ACCESS_TOKEN", "")
	t.Setenv("TAPD_API_USER", "")
	t.Setenv("TAPD_API_PASSWORD", "")
	t.Setenv("TAPD_WORKSPACE_ID", "")
	t.Setenv("HOME", tmp)

	// 模拟 home 目录下的配置
	writeConfigFile(t, filepath.Join(tmp, ".tapd.json"), &model.Config{
		AccessToken: "home_token",
		WorkspaceID: "home_ws",
	})

	// 创建子目录作为"当前目录"
	localDir := filepath.Join(tmp, "project")
	if err := os.MkdirAll(localDir, 0755); err != nil {
		t.Fatal(err)
	}
	chdirTemp(t, localDir)

	// 当前目录写入不同的配置
	writeConfigFile(t, filepath.Join(localDir, ".tapd.json"), &model.Config{
		APIUser:     "local_user",
		APIPassword: "local_pass",
		WorkspaceID: "local_ws",
	})

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	// 本地文件凭据应覆盖 home 文件
	if cfg.APIUser != "local_user" {
		t.Errorf("expected api_user=local_user, got %s", cfg.APIUser)
	}
	if cfg.WorkspaceID != "local_ws" {
		t.Errorf("expected workspace_id=local_ws, got %s", cfg.WorkspaceID)
	}
}

func TestSaveConfig(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "subdir", ".tapd.json")

	cfg := &model.Config{
		AccessToken: "save_token",
		WorkspaceID: "ws456",
	}
	if err := config.SaveConfig(cfg, path); err != nil {
		t.Fatal(err)
	}

	// 验证文件存在且内容正确
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var loaded model.Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatal(err)
	}
	if loaded.AccessToken != "save_token" || loaded.WorkspaceID != "ws456" {
		t.Errorf("unexpected config content: %+v", loaded)
	}

	// 验证文件权限
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("expected permission 0600, got %04o", perm)
	}
}

func TestSaveWorkspaceID_NewFile(t *testing.T) {
	tmp := t.TempDir()
	chdirTemp(t, tmp)

	if err := config.SaveWorkspaceID("new_ws"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, ".tapd.json"))
	if err != nil {
		t.Fatal(err)
	}
	var cfg model.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.WorkspaceID != "new_ws" {
		t.Errorf("expected workspace_id=new_ws, got %s", cfg.WorkspaceID)
	}
}

func TestSaveWorkspaceID_PreservesExisting(t *testing.T) {
	tmp := t.TempDir()
	chdirTemp(t, tmp)

	// 先写入已有凭据
	writeConfigFile(t, filepath.Join(tmp, ".tapd.json"), &model.Config{
		AccessToken: "keep_me",
		APIUser:     "keep_user",
	})

	// 只更新 workspace_id
	if err := config.SaveWorkspaceID("updated_ws"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, ".tapd.json"))
	if err != nil {
		t.Fatal(err)
	}
	var cfg model.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.AccessToken != "keep_me" {
		t.Errorf("expected access_token=keep_me, got %s", cfg.AccessToken)
	}
	if cfg.APIUser != "keep_user" {
		t.Errorf("expected api_user=keep_user, got %s", cfg.APIUser)
	}
	if cfg.WorkspaceID != "updated_ws" {
		t.Errorf("expected workspace_id=updated_ws, got %s", cfg.WorkspaceID)
	}
}

func TestGetConfigPath_Local(t *testing.T) {
	path, err := config.GetConfigPath(true)
	if err != nil {
		t.Fatal(err)
	}
	if path != ".tapd.json" {
		t.Errorf("expected .tapd.json, got %s", path)
	}
}

func TestGetConfigPath_Global(t *testing.T) {
	path, err := config.GetConfigPath(false)
	if err != nil {
		t.Fatal(err)
	}
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".tapd.json")
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}
