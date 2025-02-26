package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	redeployContainerCmd bool
	scaleContainerCmd    int64
	imageTagCmd          string
)

func init() {
	updateCmd.AddCommand(updateContainerCmd)
	updateContainerCmd.PersistentFlags().BoolVarP(&redeployContainerCmd, "redeploy", "r", false, "Redeploy with the current configuration.")
	updateContainerCmd.PersistentFlags().Int64Var(&scaleContainerCmd, "scale", 0, "Scale the container service")
	updateContainerCmd.PersistentFlags().StringVar(&imageTagCmd, "tag", "", "Update the container image tag")
}

var updateContainerCmd = &cobra.Command{
	Use:     "container [name]",
	Short:   "Update a container service",
	PreRunE: updateCmdPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("update container: %+v", args)

		if updateResource == nil {
			return errors.New("no resource provided")
		}

		var j []byte
		var err error

		if cmd.Flags().Changed("tag") {
			if j, err = updateContainerImage(updateParams, updateResource, imageTagCmd, redeployContainerCmd); err != nil {
				return err
			}
		} else if cmd.Flags().Changed("scale") {
			if j, err = scaleContainer(updateParams, updateResource, scaleContainerCmd, redeployContainerCmd); err != nil {
				return err
			}
		} else if redeployContainerCmd {
			if j, err = redeployContainer(updateParams, updateResource); err != nil {
				return err
			}
		}

		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		f.Write(j)

		return nil
	},
}

func redeployContainer(params map[string]string, resource *spinup.Resource) ([]byte, error) {
	input, err := json.Marshal(map[string]bool{"only_redeploy": true})
	if err != nil {
		return []byte{}, err
	}

	log.Debugf("putting input: %s", string(input))

	info := &spinup.ContainerService{}
	if err = SpinupClient.PutResource(params, input, info); err != nil {
		return []byte{}, err
	}

	return []byte("OK\n"), nil
}

func scaleContainer(params map[string]string, resource *spinup.Resource, scale int64, force bool) ([]byte, error) {
	log.Infof("scaling container service to %d", scale)

	input, err := json.Marshal(spinup.ContainerServiceWrapperUpdateInput{
		Size: resource.SizeID,
		Service: &spinup.ContainerServiceUpdateInput{
			CapacityProviderStrategy: []*spinup.CapacityProviderStrategyInput{
				{
					Base:             1,
					CapacityProvider: "FARGATE_SPOT",
					Weight:           1,
				},
			},
			DesiredCount:    scale,
			PlatformVersion: "LATEST",
		},
		ForceRedeploy: force,
	})
	if err != nil {
		return []byte{}, err
	}

	log.Debugf("putting input: %s", string(input))

	info := &spinup.ContainerService{}
	if err = SpinupClient.PutResource(params, input, info); err != nil {
		return []byte{}, err
	}

	return []byte("OK\n"), nil
}

func updateContainerImage(params map[string]string, resource *spinup.Resource, tag string, force bool) ([]byte, error) {
	log.Infof("updating container service image tag to %s", tag)

	// Get the container service details
	info := &spinup.ContainerService{}
	if err := SpinupClient.GetResource(params, info); err != nil {
		return []byte{}, err
	}

	// Create a copy of the container service update input with the existing settings
	serviceInput := &spinup.ContainerServiceUpdateInput{
		CapacityProviderStrategy: []*spinup.CapacityProviderStrategyInput{
			{
				Base:             1,
				CapacityProvider: "FARGATE_SPOT",
				Weight:           1,
			},
		},
		DesiredCount:    info.DesiredCount,
		PlatformVersion: "LATEST",
		ImageTag:        tag,
	}

	input, err := json.Marshal(spinup.ContainerServiceWrapperUpdateInput{
		Size:         resource.SizeID,
		Service:      serviceInput,
		ForceRedeploy: force || true, // Force redeploy when updating the image tag
	})
	if err != nil {
		return []byte{}, err
	}

	log.Debugf("putting input: %s", string(input))

	if err = SpinupClient.PutResource(params, input, info); err != nil {
		return []byte{}, err
	}

	return []byte("OK\n"), nil
}