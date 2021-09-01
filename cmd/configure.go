package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var show bool

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.Flags().BoolVar(&show, "show", false, "Display the current configurations")
}

var configureCmd = &cobra.Command{
	Use:     "configure",
	Aliases: []string{"config"},
	Short:   "Configure Spinup CLI",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("Configuring...")

		if show {
			for k, v := range viper.AllSettings() {
				fmt.Printf("%s:\t%+v\n", k, v)
			}
			return nil
		}

		var url, token, spaces string

		fmt.Printf("URL [%s]: ", spinupURL)
		fmt.Scanln(&url)
		if url == "" {
			url = spinupURL
		}
		viper.Set("url", url)

		fmt.Printf("Token [%s]: ", spinupToken)
		fmt.Scanln(&token)
		if token == "" {
			token = spinupToken
		}
		viper.Set("token", token)

		fmt.Printf("Spaces [%s]: ", strings.Join(spinupSpaces, ","))
		fmt.Scanln(&spaces)
		if spaces == "" {
			spaces = strings.Join(spinupSpaces, ",")
		}

		spaceNames := strings.Split(spaces, ",")
		if spaces != "" {
			viper.Set("spaces", spaceNames)
		} else {
			viper.Set("spaces", []string{})
		}

		log.Debugf("setting url %s, token %s, and spaces %+v", url, token, spaces)

		if err := viper.WriteConfig(); err != nil {
			return err
		}

		return nil
	},
}
