package cmd

import (
	"fmt"

	"github.com/devxdh/gitingo/pkg/repo"
	"github.com/spf13/cobra"
)

var oneline bool

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit logs",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		commits, err := repo.Log(".")
		if err != nil {
			return err
		}

		for i, commit := range commits {
			if oneline {
				firstLine := commit.Message
				if idx := len(firstLine); idx > 0 {
					lines := fmt.Sprintf("%s", firstLine)
					fmt.Printf("%s %s\n", commit.Hash[:7], lines)
				}
			} else {
				fmt.Printf("commit %s\n", commit.Hash)
				fmt.Printf("Author: %s\n", commit.Author)
				fmt.Printf("Date:   %s\n\n", commit.Date.Format("Mon Jan 02 15:04:05 2006 -0700"))
				fmt.Printf("    %s\n", commit.Message)

				if i < len(commits)-1 {
					fmt.Println()
				}
			}
		}

		return nil
	},
}

func init() {
	logCmd.Flags().BoolVar(&oneline, "oneline", false, "Show commits in a single line format")
	rootCmd.AddCommand(logCmd)
}
