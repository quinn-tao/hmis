package cmd

import (
	"github.com/quinn-tao/hmis/v1/internal/coins"
	"github.com/quinn-tao/hmis/v1/internal/db"
	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/display"
	"github.com/quinn-tao/hmis/v1/internal/profile"
	"github.com/quinn-tao/hmis/v1/internal/util"
	"github.com/spf13/cobra"
)

// recordCmd represents the record command
var recordAddCmd = &cobra.Command{
	Use:   "a [amount] [name] [category]",
	Short: "Add expense record",
	Args:  cobra.ExactArgs(3),
	Run:   parseRecordAddArgs,
}

func parseRecordAddArgs(cmd *cobra.Command, args []string) {
	recAmount, err := coins.NewFromString(args[0])
	if err != nil {
		display.Errorf("Error parsing amount:%v", err)
	}

	recName := args[1]
	recCategory := args[2]
	debug.Tracef("cmd: rec/add %v %v %v", recName, recAmount, recCategory)

	_, found := profile.FindCategory(recCategory)
	if !found {
		display.Errorf("Category %v not found", recCategory)
	}

	err = db.InsertRec(recAmount, recName, recCategory)
	util.CheckErrorf(err, "Error inserting record")
}

func init() {
	recordCmd.AddCommand(recordAddCmd)
}
