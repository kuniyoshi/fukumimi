package cmd

import (
	"fmt"
	"os"

	"github.com/kuniyoshi/fukumimi/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fukumimi",
	Short: "A CLI tool for managing radio show episodes",
	Long: `fukumimi is a CLI tool that retrieves and tracks radio show episodes
from a fan club website. It manages local read/unread status for each
broadcast episode.`,
	Version: version.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
