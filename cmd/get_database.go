package cmd

import (
	"errors"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	"github.com/spf13/cobra"
)

/* Initialized the command in the getCmd thing. */
func init() {
	getCmd.AddCommand(getDatabaseCmd)
}

/* Defines the cmd via a bunch of flags
Using the cobra library */
var getDatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Get details about a database resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("exactly 1 container service id is required")
		}
		resource := &spinup.Resource{}
		if err := SpinupClient.GetResource(map[string]string{"id": args[0]}, resource); err != nil {
			return err
		}

		return nil
	},
}
