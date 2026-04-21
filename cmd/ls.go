package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/axliupore/mk/internal"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all stored aliases",
	Long:  "Display all saved key aliases. Secret values are never shown.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		aliases, err := internal.List()
		if err != nil {
			return fmt.Errorf("failed to list keys: %w", err)
		}

		if len(aliases) == 0 {
			fmt.Println(internal.Dimf("  No keys stored."))
			return nil
		}

		// Sort for stable output
		sort.Strings(aliases)

		for _, alias := range aliases {
			fmt.Println(internal.ListItem(alias))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
