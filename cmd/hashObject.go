package cmd

import (
	"fmt"

	"github.com/devxdh/gitingo/pkg/repo"
	"github.com/spf13/cobra"
)

var writeFlag bool

var hashObjectCmd = &cobra.Command{
	Use:   "hashObject [flags] <file>",
	Short: "Compute object ID and optinally creates a blob from a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hash, err := repo.HashObject(".", args[0], writeFlag)
		if err != nil {
			return err
		}
		fmt.Println(hash)
		return nil
	},
}

func init() {
	hashObjectCmd.Flags().BoolVarP(&writeFlag, "write", "w", false, "Write the object into the object database")
	rootCmd.AddCommand(hashObjectCmd)
}
