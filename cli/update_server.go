package cli

import (
	"encoding/json"
	"errors"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverAction string

func init() {
	updateCmd.AddCommand(updateServerCmd)
	updateServerCmd.PersistentFlags().StringVar(&serverAction, "action", "", "Action to perform: start, stop, reboot, poweroff")
}

var updateServerCmd = &cobra.Command{
	Use:     "server [space]/[resource]",
	Short:   "Update a server service",
	PreRunE: updateCmdPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("update server: %+v", args)

		if serverAction == "" {
			return errors.New("--action flag is required (start, stop, reboot, poweroff)")
		}

		return updateServer(updateParams, serverAction)
	},
}

func updateServer(params map[string]string, action string) error {
	allowed := map[string]bool{"start": true, "stop": true, "reboot": true, "poweroff": true}
	if !allowed[action] {
		return errors.New("invalid action: " + action + " (allowed: start, stop, reboot, poweroff)")
	}

	input, err := json.Marshal(map[string]string{"action": action})
	if err != nil {
		return err
	}

	log.Debugf("putting input: %s", string(input))

	sa := &spinup.ServerAction{}
	return SpinupClient.PutResource(params, input, sa)
}
