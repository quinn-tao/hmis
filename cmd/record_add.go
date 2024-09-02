package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// recordCmd represents the record command
var recordAddCmd = &cobra.Command{
	Use:   "a [amount] [name] [category] [FLAGS]",
	Short: "Add expense record",
    Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("????")
	},
}

func init() {
	recordCmd.AddCommand(recordAddCmd)
}
