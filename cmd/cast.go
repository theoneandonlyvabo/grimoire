package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theoneandonlyvabo/grimoire/ui"
)

var castCmd = &cobra.Command{
	Use:   "cast",
	Short: "Cast open your Grimoire — read your technical notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ui.RunCast()
	},
}

func init() {
	rootCmd.AddCommand(castCmd)
}
