package cmd

import (
	"fmt"
	"os"

	"github.com/devxdh/gitingo/pkg/repo"
	"github.com/spf13/cobra"
)

var (
	prettyPrint bool
	showType    bool
)

var catFileCmd = &cobra.Command{
	Use:   "cat-file [flags] <object-hash>",
	Short: "Provide content or type of repository objects",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !prettyPrint && !showType {
			return fmt.Errorf("must specify either -p (pretty-print) or -t (type)")
		}

		hash := args[0]
		objType, content, err := repo.ReadObject(".", hash)
		if err != nil {
			return err
		}

		if showType {
			fmt.Println(objType)
			return nil
		}

		if prettyPrint {
			os.Stdout.Write(content)
		}

		return nil
	},
}

func init() {
	catFileCmd.Flags().BoolVarP(&prettyPrint, "pretty-print", "p", false, "Pretty-print the contents of <object>")
	catFileCmd.Flags().BoolVarP(&showType, "type", "t", false, "Show object type")
	rootCmd.AddCommand(catFileCmd)
}
