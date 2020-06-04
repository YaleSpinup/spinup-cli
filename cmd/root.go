/*
Copyright © 2020 Yale University

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	Version           string
	VersionPrerelease string
	BuildStamp        string
	GitHash           string
	cfgFile           string
	spinupURL         string
	spinupUser        string
	spinupPass        string
	debug             bool
	verbose           bool
	SpinupClient      *spinup.Client
	spinupSpaceIDs    []string
)

// rootCmd represents the base command when called without any subcommands, it propogates the configuration items from the config file.
var rootCmd = &cobra.Command{
	Use:   "spinup ",
	Short: "A small CLI for interacting with Yale's Spinup service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			log.SetLevel(log.DebugLevel)
		} else if verbose {
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}

		spinupURL = viper.GetString("url")
		spinupUser = viper.GetString("username")
		spinupPass = viper.GetString("password")
		spinupSpaceIDs = viper.GetStringSlice("spaces")

		log.Debugf("command: %+v, args: %+v", cmd, args)

		called := cmd.CalledAs()
		if called != "version" && called != "help" {
			log.Debug("initializaing client from execute()")

			if err := initClient(); err != nil {
				log.Fatalf("failed to create client: %s", err)
			}
		}

		log.Debug("running root level prerun")

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	log.Debug("executing root command")
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute command: %s", err)
	}
}

func init() {
	log.Debug("binding flags to variables")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.spinup.yaml)")
	rootCmd.PersistentFlags().StringVarP(&spinupURL, "url", "", "", "The base url for Spinup")
	rootCmd.PersistentFlags().StringVarP(&spinupUser, "username", "u", "", "Spinup username")
	rootCmd.PersistentFlags().StringVarP(&spinupPass, "password", "p", "", "Spinup password")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringSliceVarP(&spinupSpaceIDs, "spaces", "s", nil, "Space ID")

	log.Debug("viper binding flags")
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("spaces", rootCmd.PersistentFlags().Lookup("spaces"))

	log.Debug("initializing configuration")
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		log.Debugf("viper setconfigfile %s", cfgFile)
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		log.Debug("finding default config file")

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".spinup" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".spinup")
	}

	viper.AutomaticEnv() // read in environment variables that match

	log.Debug("reading config file")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}
}
