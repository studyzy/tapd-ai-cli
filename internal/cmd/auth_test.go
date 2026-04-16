package cmd

import (
	"testing"
)

// TestAuthCommand_HasLogin 验证 authCmd 下注册了 login 子命令
func TestAuthCommand_HasLogin(t *testing.T) {
	found := false
	for _, sub := range authCmd.Commands() {
		if sub.Name() == "login" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("authCmd should have 'login' subcommand")
	}
}

// TestLoginCommand_Exists 验证 loginCmd 已注册在 authCmd 下且可通过 rootCmd 访问
func TestLoginCommand_Exists(t *testing.T) {
	// 验证 authCmd 注册在 rootCmd 下
	found := false
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == "auth" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("rootCmd should have 'auth' subcommand")
	}

	// 验证 loginCmd 的基本属性
	if loginCmd.Use != "login" {
		t.Errorf("loginCmd.Use = %q, want %q", loginCmd.Use, "login")
	}
	if loginCmd.Short == "" {
		t.Error("loginCmd.Short should not be empty")
	}
	if loginCmd.RunE == nil {
		t.Error("loginCmd.RunE should not be nil")
	}
}

// TestLoginCommand_Flags 验证 loginCmd 上定义了预期的标志
func TestLoginCommand_Flags(t *testing.T) {
	// --local 是 loginCmd 自身的 flag
	localFlag := loginCmd.Flags().Lookup("local")
	if localFlag == nil {
		t.Fatal("loginCmd should have --local flag")
	}
	if localFlag.DefValue != "false" {
		t.Errorf("--local default = %q, want %q", localFlag.DefValue, "false")
	}

	// --access-token, --api-user, --api-password 是 rootCmd 的 PersistentFlags，loginCmd 应可继承
	persistentFlags := []string{"access-token", "api-user", "api-password"}
	for _, name := range persistentFlags {
		f := loginCmd.InheritedFlags().Lookup(name)
		if f == nil {
			// 也检查 Flags()（合并后的结果）
			f = loginCmd.Flags().Lookup(name)
		}
		if f == nil {
			t.Errorf("loginCmd should inherit --%s flag from rootCmd", name)
		}
	}
}
