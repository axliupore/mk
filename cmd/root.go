package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/axliupore/mk/internal"
)

// Version info injected by GoReleaser via ldflags at build time.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "mk",
	Short: "My Keys — Minimalist API key manager",
	Long: `mk (My Keys) — Two letters to manage all your API keys.
Securely stored in the system Keychain. Zero config, zero footprint.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate(fmt.Sprintf("mk %s (commit: %s, built: %s)\n", Version, Commit, Date))
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, internal.Errorf("  ✗  %s", err.Error()))
		os.Exit(1)
	}
}
