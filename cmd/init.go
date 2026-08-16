package cmd

import (
	"github.com/devxdh/gitingo/pkg/repo"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [target-file-path, defaults to current directory]",
	Short: "Initializes the a folder as repository and make .git folder",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetDir := "."
		if len(args) > 0 {
			targetDir = args[0]
		}
		repo.Init(targetDir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
