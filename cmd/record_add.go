package cmd

import (
	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/display"
	"github.com/quinn-tao/hmis/v1/internal/util"
	"github.com/spf13/cobra"
)

// recordCmd represents the record command
var recordAddCmd = &cobra.Command{
	Use:   "a [amount] [name] [category]",
	Short: "Add expense record",
    Args: cobra.ExactArgs(3),
	Run: parseRecordAddArgs,
}

func parseRecordAddArgs(cmd *cobra.Command, args []string) {
    recName := args[0]
    recAmount, err := util.StringToCents(args[1])
    if err != nil {
        display.Error("Error parsing amount")
    }
    recCategory := args[2]
    debug.TraceF("cmd: rec/add %v %v %v", recName, recAmount, recCategory)
}

func init() {
	recordCmd.AddCommand(recordAddCmd)
}
