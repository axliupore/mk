package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/axliupore/mk/internal"
)

var getCmd = &cobra.Command{
	Use:   "get <alias>",
	Short: "Retrieve a secret",
	Long:  "Print the API key stored under the given alias to stdout.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]

		key, err := internal.Get(alias)
		if err != nil {
			if internal.IsNotFound(err) {
				return fmt.Errorf("key %s not found", internal.Alias(alias))
			}
			return fmt.Errorf("failed to get key: %w", err)
		}

		// Output only the key value for pipe compatibility.
		fmt.Println(key)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
