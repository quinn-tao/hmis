/*
Copyright © 2024 Quinn Tao t.quinn.t.dev@gmail.com

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/quinn-tao/hmis/v1/internal/profile"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// Main application command 
var rootCmd = &cobra.Command{
	Use:   "hmis",
	Short: "(H)ow (M)uch (I) (S)pent?",
    Long: "hmis is a command-line tool for managing personal budgetting and expenses",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hmis.yaml)")
	rootCmd.Flags().BoolP("verbose", "v", false, "Enable application tracing")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		cfgDir, err := os.UserConfigDir()
		cobra.CheckErr(err)

        hmisDir := "hmis" 
        
        // Application config 
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".hmis")
        
        // Application config default values
		viper.SetDefault("profile.dir", path.Join(cfgDir, hmisDir, "profile"))
		viper.SetDefault("profile.name", "default")
	}

    // read in environment variables that match
	viper.AutomaticEnv() 

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
    
    profile.LoadProfile()

}
