package cmd

import (
	"fmt"
	"os"

	"github.com/kuniyoshi/fukumimi/internal/merger"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge [filename]",
	Short: "Merge fetched episodes with local listened status",
	Long: `Merge fetched episodes from stdin with local file containing listened status.
The local file should contain episodes in the same format as fetch output,
with [x] marking listened episodes and [ ] marking unlistened ones.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]
		
		// Check if local file exists
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return fmt.Errorf("local file not found: %s", filename)
		}
		
		m := merger.New()
		
		// Read episodes from stdin
		newEpisodes, err := m.ReadEpisodesFromStdin()
		if err != nil {
			return fmt.Errorf("failed to read episodes from stdin: %w", err)
		}
		
		// Read local episodes with listened status
		localEpisodes, err := m.ReadEpisodesFromFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read local episodes: %w", err)
		}
		
		// Merge episodes preserving listened status
		merged := m.MergeEpisodes(newEpisodes, localEpisodes)
		
		// Output merged result
		fmt.Print(m.GenerateOutput(merged))
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}