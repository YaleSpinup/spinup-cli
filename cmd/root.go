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
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/YaleSpinup/spinup/pkg/cas"
	"github.com/YaleSpinup/spinup/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/publicsuffix"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	spinupURL     string
	spinupUser    string
	spinupPass    string
	debug         bool
	verbose       bool
	SpinupClient  *spinup.Client
	spinupSpaceID string
)

// rootCmd represents the base command when called without any subcommands
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

		log.Debug("running root level prerun")

		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			return err
		}

		httpClient := &http.Client{
			Jar:     jar,
			Timeout: 15 * time.Second,
		}
		err = cas.Auth(spinupUser, spinupPass, spinupURL+"/login", httpClient)
		if err != nil {
			return err
		}

		s, err := spinup.New(spinupURL, httpClient)
		if err != nil {
			return err
		}

		SpinupClient = s
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.spinup.yaml)")
	rootCmd.PersistentFlags().StringVarP(&spinupURL, "url", "", "", "The base url for Spinup")
	rootCmd.PersistentFlags().StringVarP(&spinupUser, "username", "u", "", "Spinup username")
	rootCmd.PersistentFlags().StringVarP(&spinupPass, "password", "p", "", "Spinup password")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringVarP(&spinupSpaceID, "space", "s", "", "Space ID")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
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

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
