package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theoneandonlyvabo/grimoire/ui"
)

var forgeCmd = &cobra.Command{
	Use:   "forge",
	Short: "Forge a new Grimoire into your codebase",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ui.RunForge()
	},
}

func init() {
	rootCmd.AddCommand(forgeCmd)
}
