package cli

import (
	"errors"
	"strings"
	"time"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var (
	updateParams   = map[string]string{}
	updateResource = &spinup.Resource{}
)

func updateCmdPreRun(cmd *cobra.Command, args []string) error {
	defer timeTrack(time.Now(), "updateCmdPreRun()")

	if len(args) == 0 {
		return errors.New("space/resource required")
	}

	parts := strings.Split(args[0], "/")
	switch len(parts) {
	case 2:
		updateParams["space"] = parts[0]
		updateParams["name"] = parts[1]
	case 1:
		log.Debug("space not found in input, finding resource in default spaces")

		if len(spinupSpaces) == 0 {
			return errors.New("space not passed and no default spaces found")
		}

		space, err := findResourceInSpaces(parts[0], spinupSpaces)
		if err != nil {
			return err
		}

		updateParams["space"] = space
		updateParams["name"] = parts[0]
	default:
		return errors.New("space/resource required")
	}

	// set the global updateResource to the passed resource
	if err := SpinupClient.GetResource(updateParams, updateResource); err != nil {
		return err
	}

	return nil
}

var updateCmd = &cobra.Command{
	Use:   "update [type] [space]/[resource]",
	Short: "Update a resource in a space",
}
