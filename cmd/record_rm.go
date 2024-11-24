package cmd

import (
	"strconv"

	"github.com/quinn-tao/hmis/v1/internal/db"
	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/util"
	"github.com/spf13/cobra"
)

// recordCmd represents the record command
var recordRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove an expense record",
	Args:  cobra.ExactArgs(1),
	Run:   parseRecordRmArgs,
}

func parseRecordRmArgs(cmd *cobra.Command, args []string) {
	recId, err := strconv.Atoi(args[0])
	util.CheckErrorf(err, "Error parsing amount")
	debug.Tracef("cmd: rec/rm %v", recId)

	err = db.RemoveRec(recId)
	util.CheckErrorf(err, "Error removing rec")
}

func init() {
	recordCmd.AddCommand(recordRmCmd)
}
