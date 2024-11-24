package cmd

import (
	"github.com/quinn-tao/hmis/v1/internal/profile"
	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "Show current budget usage report",
	Run:   parserStatCmd,
}

func parserStatCmd(cmd *cobra.Command, args []string) {
	profile.Dump()
}

func init() {
	statCmd.Flags().BoolP("detail", "d", false, "Print detailed expense report.")
	statCmd.Flags().StringP("start", "s", "This Month",
		"Specify status report start date")
	statCmd.Flags().StringP("end", "e", "Present",
		"Specify status report end date")

	rootCmd.AddCommand(statCmd)
}
