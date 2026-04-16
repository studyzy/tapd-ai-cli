// Package cmd 中的 wiki.go 实现了 Wiki 文档管理命令
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/model"
	"github.com/studyzy/tapd-ai-cli/internal/output"
)

var flagWikiName string

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

func init() {
	wikiListCmd.Flags().IntVar(&flagLimit, "limit", 10, "返回数量限制")
	wikiListCmd.Flags().IntVar(&flagPage, "page", 1, "页码")
	wikiListCmd.Flags().StringVar(&flagWikiName, "name", "", "按标题筛选")

	wikiCmd.AddCommand(wikiListCmd, wikiShowCmd)
	rootCmd.AddCommand(wikiCmd)
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
	return output.PrintJSON(os.Stdout, wiki, !flagPretty)
}
