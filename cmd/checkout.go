package cmd

import (
	"github.com/devxdh/gitingo/pkg/repo"
	"github.com/spf13/cobra"
)

var createBranchFlag bool

var checkoutCmd = &cobra.Command{
	Use:   "checkout [-b] <branch-or-commit>",
	Short: "Switch branches or restore working tree files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		return repo.Checkout(".", target, createBranchFlag)
	},
}

func init() {
	checkoutCmd.Flags().BoolVarP(&createBranchFlag, "branch", "b", false, "Create and checkout a new branch")
	rootCmd.AddCommand(checkoutCmd)
}
