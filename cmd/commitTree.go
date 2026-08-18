package cmd

import (
	"fmt"

	"github.com/devxdh/gitingo/pkg/repo"
	"github.com/spf13/cobra"
)

var (
	parentCommit string
	commitMsg    string
)

var commitTreeCmd = &cobra.Command{
	Use:   "commit-tree <tree-hash> [flags]",
	Short: "Create a new commit object",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if commitMsg == "" {
			return fmt.Errorf("commit message cannot be empty (use -m <message>)")
		}

		treeHash := args[0]
		commitHash, err := repo.CommitTree(".", treeHash, parentCommit, commitMsg)
		if err != nil {
			return err
		}

		fmt.Println(commitHash)
		return nil
	},
}

func init() {
	commitTreeCmd.Flags().StringVarP(&commitMsg, "message", "m", "", "Commit message")
	commitTreeCmd.Flags().StringVarP(&parentCommit, "parent", "p", "", "Parent commit hash")
	commitTreeCmd.MarkFlagRequired("message")
	rootCmd.AddCommand(commitTreeCmd)
}
