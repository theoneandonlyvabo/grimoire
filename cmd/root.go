package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/theoneandonlyvabo/grimoire/ui"
)

var rootCmd = &cobra.Command{
	Use:   "grimoire",
	Short: "Technical memory for your codebase, straight from the terminal.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ui.StartMenu()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
