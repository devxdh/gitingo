package cmd

import (
	"fmt"

	"github.com/devxdh/gitingo/pkg/repo"
	"github.com/spf13/cobra"
)

var nameOnly bool

var lsTreeCmd = &cobra.Command{
	Use:   "ls-tree [flags] <tree-ish>",
	Short: "List the contents of a tree object",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		treeHash := args[0]
		entries, err := repo.LsTree(".", treeHash)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if nameOnly {
				fmt.Println(entry.Name)
			} else {
				objType := "blob"
				if entry.Mode == "040000" || entry.Mode == "40000" {
					objType = "tree"
				}
				fmt.Printf("%06s %s %s\t%s\n", entry.Mode, objType, entry.Hash, entry.Name)
			}
		}

		return nil
	},
}

func init() {
	lsTreeCmd.Flags().BoolVar(&nameOnly, "name-only", false, "List only filenames")
	rootCmd.AddCommand(lsTreeCmd)
}
