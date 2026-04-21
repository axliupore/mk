package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/axliupore/mk/internal"
)

var rmCmd = &cobra.Command{
	Use:   "rm <alias>",
	Short: "Remove a secret",
	Long:  "Permanently delete the API key stored under the given alias.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]

		if err := internal.Delete(alias); err != nil {
			if internal.IsNotFound(err) {
				return fmt.Errorf("key %s not found", internal.Alias(alias))
			}
			return fmt.Errorf("failed to delete key: %w", err)
		}

		fmt.Println(internal.Successf("Removed %s", internal.Alias(alias)))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
