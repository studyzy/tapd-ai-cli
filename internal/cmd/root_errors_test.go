package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// newGuardTestTree 构造 root -> parent -> leaf 的独立命令树，避免污染全局 rootCmd
func newGuardTestTree() (*cobra.Command, *bool) {
	leafRan := false
	root := &cobra.Command{Use: "tapd", SilenceErrors: true, SilenceUsage: true}
	parent := &cobra.Command{Use: "bug", Short: "缺陷管理"}
	leaf := &cobra.Command{
		Use:  "update",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			leafRan = true
			return nil
		},
	}
	parent.AddCommand(leaf)
	root.AddCommand(parent)
	guardUnknownSubcommands(root)
	return root, &leafRan
}

func TestGuardUnknownSubcommands_UnknownChildReturnsError(t *testing.T) {
	root, _ := newGuardTestTree()
	root.SetArgs([]string{"bug", "comment", "add"})
	root.SetOut(&bytes.Buffer{})

	_, err := root.ExecuteC()
	if err == nil {
		t.Fatal("expected error for unknown subcommand, got nil")
	}
	if !strings.Contains(err.Error(), `unknown command "comment" for "tapd bug"`) {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGuardUnknownSubcommands_NoArgsShowsHelp(t *testing.T) {
	root, _ := newGuardTestTree()
	var out bytes.Buffer
	root.SetArgs([]string{"bug"})
	root.SetOut(&out)

	if _, err := root.ExecuteC(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !strings.Contains(out.String(), "缺陷管理") {
		t.Errorf("expected help output, got: %s", out.String())
	}
}

func TestGuardUnknownSubcommands_KnownChildStillRuns(t *testing.T) {
	root, leafRan := newGuardTestTree()
	root.SetArgs([]string{"bug", "update", "123"})

	if _, err := root.ExecuteC(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !*leafRan {
		t.Error("leaf command should have run")
	}
}

func TestPrintCommandError_WritesJSONWithHint(t *testing.T) {
	root, _ := newGuardTestTree()
	root.SetArgs([]string{"bug", "update", "--id=1"})

	cmd, err := root.ExecuteC()
	if err == nil {
		t.Fatal("expected unknown flag error, got nil")
	}
	var buf bytes.Buffer
	printCommandError(&buf, cmd, err)

	var resp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
		Hint    string `json:"hint"`
	}
	if jerr := json.Unmarshal(buf.Bytes(), &resp); jerr != nil {
		t.Fatalf("output is not valid JSON: %v, raw: %s", jerr, buf.String())
	}
	if resp.Error != "command_error" {
		t.Errorf("error code = %q, want command_error", resp.Error)
	}
	if resp.Message != "unknown flag: --id" {
		t.Errorf("message = %q, want %q", resp.Message, "unknown flag: --id")
	}
	if resp.Hint != "Run 'tapd bug update --help' for usage" {
		t.Errorf("hint = %q", resp.Hint)
	}
}

func TestPrintCommandError_NilCommand(t *testing.T) {
	var buf bytes.Buffer
	printCommandError(&buf, nil, errors.New("boom"))
	if !strings.Contains(buf.String(), "Run 'tapd <command> --help' for usage") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}
