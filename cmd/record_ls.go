package cmd

import (
	"fmt"

	"github.com/quinn-tao/hmis/v1/internal/db"
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
	recs, err := db.GetAllRec()
	util.CheckErrorf(err, "Cannot get records from db:%v", recs)
	for _, rec := range recs {
		// TODO: this should be in display module
		fmt.Println(rec)
	}
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
