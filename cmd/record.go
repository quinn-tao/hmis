/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// recordCmd represents the record command
var recordCmd = &cobra.Command{
	Use:   "$",
	Short: "Managing expense records",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("record called")
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)
}
