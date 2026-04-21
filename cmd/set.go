package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/axliupore/mk/internal"
)

var setCmd = &cobra.Command{
	Use:   "set <alias> <key>",
	Short: "Store a secret",
	Long:  "Save an API key to the system Keychain under the given alias.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]
		key := args[1]

		if err := internal.Set(alias, key); err != nil {
			return fmt.Errorf("failed to set key: %w", err)
		}

		fmt.Println(internal.Successf("Set %s", internal.Alias(alias)))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
