/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// recordCmd represents the record command
var recordCmd = &cobra.Command{
	Use:   "rec [COMMAND]",
	Short: "Managing expense records",
    TraverseChildren: true, 
}

func init() {
	rootCmd.AddCommand(recordCmd)
}
