// Package cmd 定义了 tapd-ai-cli 的所有 Cobra 命令
package cmd

// Version 是当前程序版本，构建时通过 -ldflags 注入
// 格式示例：v0.1.0、v0.1.0-3-gabcdef
var Version = "dev"

func init() {
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("tapd version {{.Version}}\n")
}
