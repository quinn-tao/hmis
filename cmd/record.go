/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/quinn-tao/hmis/v1/config"
	"github.com/quinn-tao/hmis/v1/internal/db"
	"github.com/spf13/cobra"
)

// recordCmd represents the record command
var recordCmd = &cobra.Command{
	Use:              "rec [COMMAND]",
	Short:            "Managing expense records",
	TraverseChildren: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		db.PersistorInit(config.StorageLocation())
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		db.PersistorClose()
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)
}
