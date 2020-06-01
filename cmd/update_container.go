package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var redeployContainerCmd bool

func init() {
	updateCmd.AddCommand(updateContainerCmd)
	updateContainerCmd.PersistentFlags().BoolVarP(&redeployContainerCmd, "redeploy", "r", false, "Redeploy with the current configuraiton.")
}

var updateContainerCmd = &cobra.Command{
	Use:   "container",
	Short: "Update a container resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("exactly 1 container service id is required")
		}

		resource := &spinup.Resource{}
		if err := SpinupClient.GetResource(map[string]string{"id": args[0]}, resource); err != nil {
			return err
		}

		var j []byte
		var err error
		switch {
		case resource.Status != "created":
			return fmt.Errorf("container must be in 'created' state, current state is %s", resource.Status)
		case redeployContainerCmd:
			if j, err = redeployContainer(resource); err != nil {
				return err
			}
		default:
			return errors.New("only redeployment is currently supported")
		}

		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		f.Write(j)

		return nil
	},
}

func redeployContainer(resource *spinup.Resource) ([]byte, error) {
	input, err := json.Marshal(map[string]bool{"only_redeploy": true})
	if err != nil {
		return []byte{}, err
	}

	log.Debugf("putting input: %s", string(input))

	info := &spinup.ContainerService{}
	if err = SpinupClient.PutResource(map[string]string{"id": resource.ID.String()}, input, info); err != nil {
		return []byte{}, err
	}

	return []byte("OK\n"), nil
}
