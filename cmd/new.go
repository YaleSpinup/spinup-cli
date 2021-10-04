package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	secretName        string
	secretValue       string
	secretValueFrom   string
	secretDescription string
)

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.AddCommand(newSpaceCmd)
	newCmd.AddCommand(newSecretCmd)
	newSecretCmd.PersistentFlags().StringVar(&secretName, "name", "", "The name of your secret")
	newSecretCmd.PersistentFlags().StringVar(&secretValue, "value", "", "The value of your secret")
	newSecretCmd.PersistentFlags().StringVar(&secretDescription, "description", "", "A short description for your secret (optional)")
	newSecretCmd.PersistentFlags().StringVar(&secretValueFrom, "from", "", "A file containing your secret value")
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new resources",
}

var newSpaceCmd = &cobra.Command{
	Use:   "space",
	Short: "Command to create a space",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating spaces is not currently supported from the CLI, please use the web interface.")
	},
}

var newSecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Command to create a secret in a space",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if secretName == "" {
			return errors.New("a secret name is required")
		}

		if secretValue == "" && secretValueFrom == "" {
			return errors.New("a secret value or file is required")
		}

		if cmd.Flags().Changed("from") {
			secretPath := filepath.Clean(secretValueFrom)
			f, err := os.Open(secretPath)
			if err != nil {
				return err
			}
			defer f.Close()

			body, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}

			log.Debugf("size of body is %d byets", len(body))

			if len(body) > 4000 {
				return errors.New("File size is greater than 4KB")
			}

			secretValue = string(body)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("creating secret %s:%s", secretName, secretValue)
		return nil
	},
}
