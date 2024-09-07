package cmd

import (
	"strconv"

	"github.com/quinn-tao/hmis/v1/internal/db"
	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/display"
	"github.com/spf13/cobra"
)

// recordCmd represents the record command
var recordRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove an expense record",
    Args: cobra.ExactArgs(1),
	Run: parseRecordRmArgs,
}

func parseRecordRmArgs(cmd *cobra.Command, args []string) {
    recId, err := strconv.Atoi(args[0])
    if err != nil {
        display.Error("Error parsing amount")
    }
    debug.TraceF("cmd: rec/rm %v", recId)

    err = db.RemoveRec(recId)    
    if err != nil {
        display.Error("Error removing rec")
    }
}

func init() {
	recordCmd.AddCommand(recordRmCmd)
}
