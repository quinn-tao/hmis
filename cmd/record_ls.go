package cmd

import (
	"github.com/quinn-tao/hmis/v1/internal/chrono"
	"github.com/quinn-tao/hmis/v1/internal/db"
	"github.com/quinn-tao/hmis/v1/internal/debug"
	"github.com/quinn-tao/hmis/v1/internal/display"
	"github.com/quinn-tao/hmis/v1/internal/display/cli"
	"github.com/quinn-tao/hmis/v1/internal/profile"
	"github.com/quinn-tao/hmis/v1/internal/util"
	"github.com/spf13/cobra"
)

// recordCmd represents the record command
var recordLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List expense record",
	Run:   parseRecordLsArgs,
}

func parseRecordLsArgs(cmd *cobra.Command, args []string) {
    var searchOpt []db.RecordSearchOption

    // Parsing flags 
    parsedFromDate, err := chrono.ParseDate(recFlags.recFromDate)
    if err != nil {
        display.Errorf("Cannot parse end flag %v", recFlags.recFromDate)
    }
    searchOpt = append(searchOpt, db.RecordSearchWithFromDate(parsedFromDate))

    parsedToDate, err := chrono.ParseDate(recFlags.recToDate)
    if err != nil {
        display.Errorf("Cannot parse end flag %v", recFlags.recFromDate)
    }
    searchOpt = append(searchOpt, db.RecordSearchWithToDate(parsedToDate))
    
    if recFlags.recCategory != "" {
        parsedCategory := profile.Category{Name: recFlags.recCategory}
        searchOpt = append(searchOpt, db.RecordSearchWithCategory(parsedCategory))  
    }

    if recFlags.recNameRg != "" {
        searchOpt = append(searchOpt, db.RecordSearchWithNameRg(recFlags.recNameRg))
    }
    
    debug.Tracef("parsed start %v end %v", 1, 1)

	recs, err := db.GetRecords(searchOpt...)
	util.CheckErrorf(err, "Cannot get records from db:%v", recs)
    
	sum, err := db.GetSumRecord(searchOpt...)
	if err != nil {
		display.Errorf("Cannot get sum: %v", err)
	}

	tbl := cli.NewTable("Records",
		cli.Column{Name: "id", Required: true},
		cli.Column{Name: "amount", Required: true},
		cli.Column{Name: "name", Required: true},
		cli.Column{Name: "category", Required: true},
		cli.Column{Name: "date", Required: true})

	for row, rec := range recs {
		err := tbl.AppendRow(rec)
		if err != nil {
			display.Errorf("Cannot display row %v", row)
		}
	}

	tbl.AppendRowWithHighlight(sum)
	tbl.Render()
}

var recFlags struct {
    recCategory string
    recFromDate string 
    recToDate string
    recNameRg string 
}

const dateFmtHelpText = "Date format must be one of:\n" +
	"1. Special Identifiers:\n" +
	"  - 'present', 'today', and 'now' translates to today\n" +
	"  - 'mtd' and 'month' translates to first day of current month\n" +
	"  - 'ytd' and 'year' translates to first day of current year.\n" +
	"\n2. Date Formats:\n" +
    "Supports the standard date format parsing and a few extensions:\n" +
	"yyyy-mm-dd, yy-mm-dd, mm-dd, dd, mm/dd, mm/dd/yy, mm/dd/yyyy are\n" +
	"all valid formats. If a year is specified using 'yy', then it is\n" +
	"treated as '20mm' year.\n" +
	"\n3. Relative Dates\n" +
	"Supported formats are: Xd, Xm, Xy, X where X is a positive integer\n" +
	"Token Interpretation:\n" +
	"  - Xd, X - X days before current date\n" +
	"  - Xm - X months before current month\n" +
	"  - Xy - X years before current year\n"

func init() {
	recordLsCmd.Flags().StringVarP(&recFlags.recNameRg,
		"name", "n", "", 
        "Regex matching by name")
	recordLsCmd.Flags().StringVarP(&recFlags.recCategory,
		"category", "c", "", 
        "Specify category")

    recordLsCmd.Flags().StringVarP(&recFlags.recFromDate, 
        "start", "s", "month",
		"Specify recordus report start date\n" + dateFmtHelpText)
	recordLsCmd.Flags().StringVarP(&recFlags.recToDate, 
        "end", "e", "Present",
		"Specify recordus report end date")

	recordLsCmd.Flags().SortFlags = false
	recordCmd.AddCommand(recordLsCmd)
}
