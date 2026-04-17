// Package cmd 中的 relation.go 实现了实体关联关系管理命令
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/studyzy/tapd-ai-cli/internal/output"
	"github.com/studyzy/tapd-sdk-go/model"
)

var (
	flagRelationStoryID    string
	flagRelationSourceType string
	flagRelationTargetType string
	flagRelationSourceID   string
	flagRelationTargetID   string
)

// relationCmd 是 relation 父命令
var relationCmd = &cobra.Command{
	Use:   "relation",
	Short: "关联关系管理",
}

var relationBugsCmd = &cobra.Command{
	Use:   "bugs",
	Short: "查询需求关联的缺陷",
	RunE:  runRelationBugs,
}

var relationCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建实体关联关系",
	RunE:  runRelationCreate,
}

func init() {
	// bugs 子命令
	relationBugsCmd.Flags().StringVar(&flagRelationStoryID, "story-id", "", "需求 ID（必需）")

	// create 子命令
	relationCreateCmd.Flags().StringVar(&flagRelationSourceType, "source-type", "", "源实体类型（story|bug|task，必需）")
	relationCreateCmd.Flags().StringVar(&flagRelationTargetType, "target-type", "", "目标实体类型（story|bug|task，必需）")
	relationCreateCmd.Flags().StringVar(&flagRelationSourceID, "source-id", "", "源实体 ID（必需）")
	relationCreateCmd.Flags().StringVar(&flagRelationTargetID, "target-id", "", "目标实体 ID（必需）")

	relationCmd.AddCommand(relationBugsCmd, relationCreateCmd)
	rootCmd.AddCommand(relationCmd)
}

func runRelationBugs(cmd *cobra.Command, args []string) error {
	if flagRelationStoryID == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--story-id is required",
			"Usage: tapd relation bugs --story-id <id>")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.GetRelatedBugsRequest{
		WorkspaceID: flagWorkspaceID,
		StoryID:     flagRelationStoryID,
	}

	data, err := apiClient.GetRelatedBugs(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, data, !flagPretty)
}

func runRelationCreate(cmd *cobra.Command, args []string) error {
	if flagRelationSourceType == "" || flagRelationTargetType == "" ||
		flagRelationSourceID == "" || flagRelationTargetID == "" {
		output.PrintError(os.Stderr, "missing_parameter",
			"--source-type, --target-type, --source-id and --target-id are required",
			"Usage: tapd relation create --source-type story --target-type bug --source-id <id> --target-id <id>")
		os.Exit(output.ExitParamError)
		return nil
	}

	req := &model.CreateRelationRequest{
		WorkspaceID: flagWorkspaceID,
		SourceType:  flagRelationSourceType,
		TargetType:  flagRelationTargetType,
		SourceID:    flagRelationSourceID,
		TargetID:    flagRelationTargetID,
	}

	data, err := apiClient.CreateRelation(req)
	if err != nil {
		output.PrintError(os.Stderr, "api_error", err.Error(), "")
		os.Exit(output.ExitAPIError)
		return nil
	}
	return output.PrintJSON(os.Stdout, data, !flagPretty)
}
