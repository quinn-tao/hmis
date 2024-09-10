package cmd

import (
	"github.com/quinn-tao/hmis/v1/internal/profile"
	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "Show current budget usage report",
	Run: parserStatCmd,
}

func parserStatCmd(cmd *cobra.Command, args []string) {
    profile.Dump()
}

func init() {
	rootCmd.AddCommand(statCmd)
}
