package cmd

import (
	"github.com/quinn-tao/hmis/v1/internal/chrono"
	"github.com/quinn-tao/hmis/v1/internal/coins"
	"github.com/quinn-tao/hmis/v1/internal/db"
	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/display"
	"github.com/quinn-tao/hmis/v1/internal/profile"
	"github.com/quinn-tao/hmis/v1/internal/record"
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

var flags struct {
    datestr string 
}

func parseRecordAddArgs(cmd *cobra.Command, args []string) {
	recAmount, err := coins.NewFromString(args[0])
	if err != nil {
		display.Errorf("Error parsing amount:%v", err)
	}

	recName := args[1]
	recCategory := args[2]

	_, found := profile.FindCategory(recCategory)
	if !found {
		display.Errorf("Category %v not found", recCategory)
	}

    recDate := chrono.Today()
    if flags.datestr != "today" {
        date, err := chrono.ParseDate(flags.datestr)
        if err != nil {
            display.Error(err.Error()) 
        }
        recDate = date 
    }
    
    newRec := record.Record{
        Amount: recAmount,
        Name: recName,
        Category: recCategory,
        Date: recDate,
    }
    
    debug.Tracef("Adding %v", newRec)
	err = db.InsertRec(newRec)
	util.CheckErrorf(err, "Error inserting record")
}

func init() {
    recordAddCmd.Flags().StringVarP(&flags.datestr, "date", "d", "today",
        "Specify date of the record")
	recordCmd.AddCommand(recordAddCmd)
}
