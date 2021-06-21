package cmd

import (
	"errors"
	"fmt"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	"github.com/spf13/cobra"
)

var detailedGetCmd bool

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().BoolVarP(&detailedGetCmd, "details", "d", false, "Get detailed output about the resource")
}

var getCmd = &cobra.Command{
	Use:   "get [space] [resource]",
	Short: "Get information about a resource in a space",
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
			return getContainer(params, resource)
		case "db":
			return getDatabase(params, resource)
		case "server":
			return getServer(params, resource)
		case "storage":
			return getStorage(params, resource)
		default:
			return fmt.Errorf("unrecognized type of resource '%s'", resource.IsA)
		}
	},
}
