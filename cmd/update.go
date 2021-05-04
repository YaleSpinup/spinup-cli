package cmd

import (
	"errors"
	"fmt"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update [space] [resource]",
	Short: "Update a resource in a space",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("space name and resource name are required")
		}

		params := map[string]string{
			"space": args[0],
			"name":  args[1],
		}

		resource := &spinup.Resource{}
		if err := SpinupClient.GetResource(params, resource); err != nil {
			return err
		}

		switch resource.IsA {
		case "container":
			return updateContainer(params, resource)
		case "server":
			return updateServer(params, resource)
		case "storage":
			return updateStorage(params, resource)
		default:
			return fmt.Errorf("unrecognized type of resource '%s'", resource.IsA)
		}
	},
}
