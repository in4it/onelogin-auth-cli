package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all roles and accounts available",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Roles:")
		for k, v := range config.Roles {
			fmt.Printf("[%d] %s\n", k, v)
		}

		fmt.Println("Accounts:")
		for k, v := range config.Accounts {
			fmt.Printf("[%d] %s\n", k, v.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
