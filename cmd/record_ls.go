package cmd

import "github.com/spf13/cobra"

// recordCmd represents the record command
var recordLsCmd = &cobra.Command{
    Use:   "ls",
	Short: "List epense record",
	Run: parseRecordLsArgs,
}

func parseRecordLsArgs(cmd *cobra.Command, args []string) {
    // TODO: impl command detail
}

var recIdFlag int
var recNameFlag string 
var recCategoryFlag string 

func init() {
    recordLsCmd.Flags().IntVarP(&recIdFlag, "id", "i", -1, "Specify id")
    recordLsCmd.Flags().StringVarP(&recNameFlag, "name", "n", "", "Fuzzy matching by name")
    recordLsCmd.Flags().StringVarP(&recCategoryFlag, "category", "c", "", "Specify category")
	recordCmd.AddCommand(recordLsCmd)
}
