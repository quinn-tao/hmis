package cmd

import (
	"github.com/quinn-tao/hmis/v1/internal/profile"
	"github.com/spf13/cobra"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Managing budgeting profiles",
	Run: func(cmd *cobra.Command, args []string) {
        profile.Dump()
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
}

