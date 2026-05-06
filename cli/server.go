package cli

import (
	"errors"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverParams = map[string]string{}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serverStartCmd)
	serverCmd.AddCommand(serverStopCmd)
	serverCmd.AddCommand(serverRebootCmd)
	serverCmd.AddCommand(serverPoweroffCmd)
}

func serverCmdPreRun(cmd *cobra.Command, args []string) error {
	defer timeTrack(time.Now(), "serverCmdPreRun()")

	if len(args) == 0 {
		return errors.New("space/resource required")
	}

	parts := strings.Split(args[0], "/")
	switch len(parts) {
	case 2:
		serverParams["space"] = parts[0]
		serverParams["name"] = parts[1]
	case 1:
		log.Debug("space not found in input, finding resource in default spaces")

		if len(spinupSpaces) == 0 {
			return errors.New("space not passed and no default spaces found")
		}

		space, err := findResourceInSpaces(parts[0], spinupSpaces)
		if err != nil {
			return err
		}

		serverParams["space"] = space
		serverParams["name"] = parts[0]
	default:
		return errors.New("space/resource required")
	}

	return nil
}

var serverCmd = &cobra.Command{
	Use:   "server [action] [space]/[resource]",
	Short: "Server power actions (start, stop, reboot, poweroff)",
}

var serverStartCmd = &cobra.Command{
	Use:     "start [space]/[resource]",
	Short:   "Start a server",
	PreRunE: serverCmdPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateServer(serverParams, "start")
	},
}

var serverStopCmd = &cobra.Command{
	Use:     "stop [space]/[resource]",
	Short:   "Stop a server",
	PreRunE: serverCmdPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateServer(serverParams, "stop")
	},
}

var serverRebootCmd = &cobra.Command{
	Use:     "reboot [space]/[resource]",
	Short:   "Reboot a server",
	PreRunE: serverCmdPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateServer(serverParams, "reboot")
	},
}

var serverPoweroffCmd = &cobra.Command{
	Use:     "poweroff [space]/[resource]",
	Short:   "Power off a server",
	PreRunE: serverCmdPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateServer(serverParams, "poweroff")
	},
}
