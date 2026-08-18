package cmd

import (
	"fmt"

	"github.com/devxdh/gitingo/pkg/repo"
	"github.com/spf13/cobra"
)

var commitMsgStr string

var commitCmd = &cobra.Command{
	Use:   "commit -m <message>",
	Short: "Record changes to the repository",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if commitMsgStr == "" {
			return fmt.Errorf("fatal: commit message cannot be empty (use -m <message>)")
		}

		_, err := repo.Commit(".", commitMsgStr)
		return err
	},
}

func init() {
	commitCmd.Flags().StringVarP(&commitMsgStr, "message", "m", "", "Commit message")
	commitCmd.MarkFlagRequired("message")
	rootCmd.AddCommand(commitCmd)
}
