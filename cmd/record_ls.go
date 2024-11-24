package cmd

import (
	"github.com/quinn-tao/hmis/v1/internal/db"
	"github.com/quinn-tao/hmis/v1/internal/display"
	"github.com/quinn-tao/hmis/v1/internal/display/cli"
	"github.com/quinn-tao/hmis/v1/internal/util"
	"github.com/spf13/cobra"
)

// recordCmd represents the record command
var recordLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List epense record",
	Run:   parseRecordLsArgs,
}

func parseRecordLsArgs(cmd *cobra.Command, args []string) {
	// TODO: [record ls] support filtering/fuzzy matching in names
	recs, err := db.GetAllRecords()
	util.CheckErrorf(err, "Cannot get records from db:%v", recs)

	sum, err := db.GetSumRecord()
	if err != nil {
		display.Errorf("Cannot get sum: %v", err)
	}

	tbl := cli.NewTable("Records",
		cli.Column{Name: "id", Required: true},
		cli.Column{Name: "amount", Required: true},
		cli.Column{Name: "name", Required: true},
		cli.Column{Name: "category", Required: true})

	for row, rec := range recs {
		err := tbl.AppendRow(rec)
		if err != nil {
			// TODO: failures could corrupt object
			display.Errorf("Cannot display row %v", row)
		}
	}

	tbl.AppendRowWithHighlight(sum)
	tbl.Render()
}

var recIdFlag int
var recNameFlag string
var recCategoryFlag string

func init() {
	recordLsCmd.Flags().IntVarP(&recIdFlag, "id", "i", -1, "Specify id")
	recordLsCmd.Flags().StringVarP(&recNameFlag,
		"name", "n", "", "Fuzzy matching by name")
	recordLsCmd.Flags().StringVarP(&recCategoryFlag,
		"category", "c", "", "Specify category")
	recordCmd.AddCommand(recordLsCmd)
}
