package cmd

import (
	"fmt"

	"github.com/devxdh/gitingo/pkg/repo"
	"github.com/spf13/cobra"
)

var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "Create a tree object from the current directory",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		treeHash, err := repo.WriteTree(".")
		if err != nil {
			return err
		}

		fmt.Println(treeHash)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
