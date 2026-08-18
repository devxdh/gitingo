package cmd

import (
	"fmt"

	"github.com/devxdh/gitingo/pkg/repo"
	"github.com/spf13/cobra"
)

var deleteBranchFlag string

var branchCmd = &cobra.Command{
	Use:   "branch [branch-name]",
	Short: "List, create, or delete branches",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if deleteBranchFlag != "" {
			return repo.DeleteBranch(".", deleteBranchFlag)
		}

		if len(args) == 1 {
			branchName := args[0]
			return repo.CreateBranch(".", branchName)
		}

		branches, err := repo.ListBranches(".")
		if err != nil {
			return err
		}

		for _, b := range branches {
			if b.IsActive {
				fmt.Printf("* %s\n", b.Name)
			} else {
				fmt.Printf("  %s\n", b.Name)
			}
		}

		return nil
	},
}

func init() {
	branchCmd.Flags().StringVarP(&deleteBranchFlag, "delete", "d", "", "Delete a branch")
	rootCmd.AddCommand(branchCmd)
}
