package cmd

import (
	"github.com/quinn-tao/hmis/v1/config"
	"github.com/quinn-tao/hmis/v1/internal/db"
	"github.com/quinn-tao/hmis/v1/internal/display"
	"github.com/quinn-tao/hmis/v1/internal/display/cli"
	"github.com/quinn-tao/hmis/v1/internal/stat"
	"github.com/spf13/cobra"
)


var statCmd = &cobra.Command{
	Use:   "stat",
	Run:   parserStatCmd,
	Short: "Show current budget status",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		db.PersistorInit(config.StorageLocation())
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		db.PersistorClose()
	},
}

var statFlags struct {
	startDate string
	endDate   string
}

func parserStatCmd(cmd *cobra.Command, args []string) {
	st, err := stat.NewStat()
	if err != nil {
		display.Error(err.Error())
	}

	report := cli.NewReport("Expense Report")
	report.AddSection("Profile Settings")
	report.AddEntry("Budget", st.GetBudget())
	report.AddEntry("Mode", st.GetMode())
	report.AddSection("Status")
	report.AddEntry("Recorded Expenses", st.GetRecordedExp())
	report.AddEntry("Fixed Expenses", st.GetFixedExp())
	report.AddEntry("TotalExpenses", st.GetTotalExp())
	report.AddEntry("Remaining Budget", st.GetRemainingBudget())

	report.Render()
}

func init() {
	statCmd.Flags().SortFlags = false
	statCmd.Flags().BoolP("detail", "d", false, "Print detailed expense report.")
	
	rootCmd.AddCommand(statCmd)
}
