package cmd

import (
	"fmt"
	"os"

	"github.com/kuniyoshi/fukumimi/internal/merger"
	"github.com/kuniyoshi/fukumimi/internal/models"
	"github.com/spf13/cobra"
)

var (
	replaceFile bool
)

var mergeCmd = &cobra.Command{
	Use:   "merge [filename]",
	Short: "Merge fetched episodes with local listened status",
	Long: `Merge fetched episodes from stdin with local file containing listened status.
The local file should contain episodes in the same format as fetch output,
with [x] marking listened episodes and [ ] marking unlistened ones.

By default, the merged result is output to stdout. Use --replace to update
the local file in-place.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]
		
		// Check if local file exists
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			// If file doesn't exist and replace flag is set, we'll create it
			if !replaceFile {
				return fmt.Errorf("local file not found: %s", filename)
			}
		}
		
		m := merger.New()
		
		// Read episodes from stdin
		newEpisodes, err := m.ReadEpisodesFromStdin()
		if err != nil {
			return fmt.Errorf("failed to read episodes from stdin: %w", err)
		}
		
		// Read local episodes with listened status (if file exists)
		var localEpisodes []models.Episode
		if _, err := os.Stat(filename); err == nil {
			localEpisodes, err = m.ReadEpisodesFromFile(filename)
			if err != nil {
				return fmt.Errorf("failed to read local episodes: %w", err)
			}
		}
		
		// Merge episodes preserving listened status
		merged := m.MergeEpisodes(newEpisodes, localEpisodes)
		
		// Output merged result
		output := m.GenerateOutput(merged)
		
		if replaceFile {
			// Write to file
			if err := os.WriteFile(filename, []byte(output), 0644); err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}
		} else {
			// Output to stdout
			fmt.Print(output)
		}
		
		return nil
	},
}

func init() {
	mergeCmd.Flags().BoolVarP(&replaceFile, "replace", "r", false, "Replace the local file with merged content")
	rootCmd.AddCommand(mergeCmd)
}