package cmd

import (
	"fmt"

	"github.com/kuniyoshi/fukumimi/internal/auth"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the fan club website",
	Long:  `Authenticate with the fan club website and store session cookies locally.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting login process...")
		
		client := auth.NewClient()
		if err := client.Login(); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}
		
		fmt.Println("Login successful!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}