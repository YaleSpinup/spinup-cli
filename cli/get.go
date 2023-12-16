package cli

import (
	"errors"
	"strings"
	"time"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var detailedGetCmd bool

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().BoolVarP(&detailedGetCmd, "details", "d", false, "Get detailed output about the resource")
}

var (
	getParams   = map[string]string{}
	getResource = &spinup.Resource{}
)

func getCmdPreRun(cmd *cobra.Command, args []string) error {
	defer timeTrack(time.Now(), "getCmdPreRun()")

	if len(args) == 0 {
		return errors.New("space/resource required")
	}

	parts := strings.Split(args[0], "/")
	switch len(parts) {
	case 2:
		getParams["space"] = parts[0]
		getParams["name"] = parts[1]
	case 1:
		log.Debug("space not found in input, finding resource in default spaces")

		if len(spinupSpaces) == 0 {
			return errors.New("space not passed and no default spaces found")
		}

		space, err := findResourceInSpaces(parts[0], spinupSpaces)
		if err != nil {
			return err
		}

		getParams["space"] = space
		getParams["name"] = parts[0]
	default:
		return errors.New("space/resource required")
	}

	// set the global getResource to the passed resource
	if err := SpinupClient.GetResource(getParams, getResource); err != nil {
		return err
	}

	return nil
}

func getSpaceLevelCmdPreRun(cmd *cobra.Command, args []string) error {
	defer timeTrack(time.Now(), "getSpaceLevelCmdPreRun()")

	if len(args) == 0 {
		return errors.New("space required")
	}

	getParams["space"] = args[0]

	return nil
}

var getCmd = &cobra.Command{
	Use:   "get [type] [space]/[resource]",
	Short: "Get information about a resource in a space",
}
