package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/axliupore/mk/internal"
)

var cpCmd = &cobra.Command{
	Use:   "cp <alias>",
	Short: "Copy a secret to clipboard",
	Long:  "Copy the API key to the system clipboard without printing it to the terminal.",
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

		if err := clipboard.WriteAll(key); err != nil {
			return fmt.Errorf(
				"failed to copy to clipboard: %w\n"+
					"On Linux, please install xclip: sudo apt install xclip",
				err,
			)
		}

		fmt.Println(internal.Successf("Copied %s to clipboard", internal.Alias(alias)))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cpCmd)
}
