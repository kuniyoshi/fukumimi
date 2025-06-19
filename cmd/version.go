package cmd

import (
	"fmt"

	"github.com/kuniyoshi/fukumimi/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of fukumimi",
	Long:  `Print the version number of fukumimi.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("fukumimi version %s\n", version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}