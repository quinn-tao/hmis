package cmd

import (
	"github.com/quinn-tao/hmis/v1/internal/profile"
	"github.com/spf13/cobra"
)

const dateFmtHelpText = "Date format must be one of:\n" +
    "1. Special Identifiers:\n" +
    "- 'present', 'today', and 'now' translates to today\n" +
    "- 'mtd' and 'month' translates to first day of current month\n" +
    "- 'ytd' and 'year' translates to first day of current year.\n" +
    "\n2. Date Formats:\n" +
    "Supports a bit more than the standard date format parsing\n" +
    "yyyy-mm-dd, yy-mm-dd, mm-dd, dd, mm/dd, mm/dd/yy, mm/dd/yyyy are\n" +
    "all valid formats. If a year is specified using 'yy', then it is\n" +
    "treated as '20mm' year.\n" +
    "\n3. Relative Dates\n" +
    "Supported formats are: Xd, Xm, Xy, X where X is an positive integer\n" +
    "Token Interpretation:\n" +
    "- Xd, X - X days before current date\n" +
    "- Xm - X months before current month\n" +
    "- Xy - X years before current year\n"

var statCmd = &cobra.Command{
	Use:   "stat",
	Run:   parserStatCmd,
    Short: "Show current budget status", 
}

func parserStatCmd(cmd *cobra.Command, args []string) {
	profile.Dump()
}

func init() {
    statCmd.Flags().SortFlags = false
	statCmd.Flags().BoolP("detail", "d", false, "Print detailed expense report.")
	statCmd.Flags().StringP("start", "s", "This Month",
		"Specify status report start date\n" + dateFmtHelpText)
	statCmd.Flags().StringP("end", "e", "Present",
		"Specify status report end date")

	rootCmd.AddCommand(statCmd)
}
