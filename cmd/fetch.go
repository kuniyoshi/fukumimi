package cmd

import (
	"fmt"

	"github.com/kuniyoshi/fukumimi/internal/fetcher"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch radio show episodes",
	Long:  `Fetch all radio show episodes from the fan club website and output as markdown.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Remove debug messages - output only markdown
		f := fetcher.New()
		episodes, err := f.FetchEpisodes()
		if err != nil {
			return fmt.Errorf("failed to fetch episodes: %w", err)
		}

		// Output markdown to STDOUT
		markdown := f.GenerateMarkdown(episodes)
		fmt.Print(markdown)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
