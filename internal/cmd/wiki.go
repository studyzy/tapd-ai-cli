// Package cmd 中的 wiki.go 实现了 Wiki 文档管理命令
package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/model"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

var (
	flagWikiName    string
	flagWikiContent string
	flagWikiFile    string
	flagWikiNote    string
	flagParentWiki  string
)

// wikiCmd 是 wiki 父命令
var wikiCmd = &cobra.Command{
	Use:   "wiki",
	Short: "Wiki 文档管理",
}

var wikiListCmd = &cobra.Command{
	Use:   "list",
	Short: "查询 Wiki 文档列表",
	RunE:  runWikiList,
}

var wikiShowCmd = &cobra.Command{
	Use:   "show <wiki_id>",
	Short: "查看 Wiki 文档详情",
	Args:  cobra.ExactArgs(1),
	RunE:  runWikiShow,
}

var wikiCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建 Wiki 文档",
	Long: `创建 Wiki 文档，内容支持三种输入方式：
  1. --content <markdown_text>  直接传入 Markdown 文本
  2. --file <path>              从本地 Markdown 文件读取内容
  3. echo "..." | tapd wiki create --name <title> --creator <user>  通过 stdin 管道输入`,
	RunE: runWikiCreate,
}

var wikiUpdateCmd = &cobra.Command{
	Use:   "update <wiki_id>",
	Short: "更新 Wiki 文档",
	Long: `更新 Wiki 文档，内容支持三种输入方式：
  1. --content <markdown_text>  直接传入 Markdown 文本
  2. --file <path>              从本地 Markdown 文件读取内容
  3. echo "..." | tapd wiki update <wiki_id>  通过 stdin 管道输入`,
	Args: cobra.ExactArgs(1),
	RunE: runWikiUpdate,
}

func init() {
	wikiListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	wikiListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")
	wikiListCmd.Flags().StringVar(&flagWikiName, "name", "", "按标题筛选")

	wikiCreateCmd.Flags().StringVar(&flagWikiName, "name", "", "Wiki 标题（必需）")
	wikiCreateCmd.Flags().StringVar(&flagCreator, "creator", "", "创建人（必需）")
	wikiCreateCmd.Flags().StringVar(&flagWikiContent, "content", "", "Wiki 内容（Markdown 格式）")
	wikiCreateCmd.Flags().StringVar(&flagWikiFile, "file", "", "从本地 Markdown 文件读取 Wiki 内容")
	wikiCreateCmd.Flags().StringVar(&flagWikiNote, "note", "", "备注")
	wikiCreateCmd.Flags().StringVar(&flagParentWiki, "parent-wiki-id", "", "父 Wiki ID")

	wikiUpdateCmd.Flags().StringVar(&flagWikiName, "name", "", "新标题")
	wikiUpdateCmd.Flags().StringVar(&flagWikiContent, "content", "", "新内容（Markdown 格式）")
	wikiUpdateCmd.Flags().StringVar(&flagWikiFile, "file", "", "从本地 Markdown 文件读取新内容")
	wikiUpdateCmd.Flags().StringVar(&flagWikiNote, "note", "", "新备注")
	wikiUpdateCmd.Flags().StringVar(&flagParentWiki, "parent-wiki-id", "", "新父 Wiki ID")

	wikiCmd.AddCommand(wikiListCmd, wikiShowCmd, wikiCreateCmd, wikiUpdateCmd)
	rootCmd.AddCommand(wikiCmd)
}

// readWikiContent 从 --content、--file 或 stdin 读取 Wiki 内容
// 优先级：--content > --file > stdin
func readWikiContent() (string, error) {
	if flagWikiContent != "" {
		return flagWikiContent, nil
	}
	if flagWikiFile != "" {
		data, err := os.ReadFile(flagWikiFile)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	// 尝试从 stdin 读取（仅当 stdin 不是终端时）
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", nil
}

func runWikiList(cmd *cobra.Command, args []string) error {
	params := map[string]string{
		"workspace_id": flagWorkspaceID,
	}
	addOptionalParam(params, "name", flagWikiName)
	addPaginationParams(params, flagLimit, flagPage)
	params["fields"] = "id,name,creator,modifier,modified"

	wikis, err := apiClient.ListWikis(params)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}

	resp := &model.ListResponse{
		Items:   wikis,
		Total:   len(wikis),
		Page:    flagPage,
		Limit:   flagLimit,
		HasMore: len(wikis) == flagLimit,
	}
	return output.PrintJSON(os.Stdout, resp, !flagPretty)
}

func runWikiShow(cmd *cobra.Command, args []string) error {
	wiki, err := apiClient.GetWiki(flagWorkspaceID, args[0])
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return printDetail(wiki, "markdown_description")
}

func runWikiCreate(cmd *cobra.Command, args []string) error {
	if flagWikiName == "" {
		output.PrintError(os.Stderr, "missing_parameter", "--name is required", "Usage: tapd wiki create --name <title> --creator <user> [--content <text> | --file <path> | stdin]")
		os.Exit(output.ExitParamError)
		return nil
	}
	if flagCreator == "" {
		output.PrintError(os.Stderr, "missing_parameter", "--creator is required", "Usage: tapd wiki create --name <title> --creator <user> [--content <text> | --file <path> | stdin]")
		os.Exit(output.ExitParamError)
		return nil
	}

	content, err := readWikiContent()
	if err != nil {
		output.PrintError(os.Stderr, "file_error", err.Error(), "Check that the file path is correct and readable")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.CreateWikiRequest{
		WorkspaceID:         flagWorkspaceID,
		Name:                flagWikiName,
		Creator:             flagCreator,
		MarkdownDescription: content,
		Note:                flagWikiNote,
		ParentWikiID:        flagParentWiki,
	}

	result, err := apiClient.CreateWiki(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}

func runWikiUpdate(cmd *cobra.Command, args []string) error {
	content, err := readWikiContent()
	if err != nil {
		output.PrintError(os.Stderr, "file_error", err.Error(), "Check that the file path is correct and readable")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.UpdateWikiRequest{
		WorkspaceID:         flagWorkspaceID,
		ID:                  args[0],
		Name:                flagWikiName,
		MarkdownDescription: content,
		Note:                flagWikiNote,
		ParentWikiID:        flagParentWiki,
	}

	result, err := apiClient.UpdateWiki(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, result, !flagPretty)
}
