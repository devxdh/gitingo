// Package cmd is the entrypoint of the project
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is populated at build time using -ldflags "-X github.com/devxdh/gig/cmd.Version=v1.0.0"
var Version = "v1.0.0"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gig",
	Short: "gig - A fast, lightweight Git VCS implementation written in Go",
	Long: `gig (Git-In-Go) is a standalone Version Control System (VCS) CLI written in Go.
It implements core Git internal objects (blobs, trees, commits), repository initialization,
object hashing, tree creation, commit logging, branch management, and checkout commands.`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Custom version template
	rootCmd.SetVersionTemplate("gig version {{.Version}}\n")
}
