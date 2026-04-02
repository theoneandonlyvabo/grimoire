package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theoneandonlyvabo/grimoire/ui"
)

var carveCmd = &cobra.Command{
	Use:   "carve",
	Short: "Carve into your Grimoire — edit your technical notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ui.RunCarve()
	},
}

func init() {
	rootCmd.AddCommand(carveCmd)
}
