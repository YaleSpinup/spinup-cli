package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	redeployContainerCmd bool
	scaleContainerCmd    int64
	containerNameCmd     string
	containerTagCmd      string
)

func init() {
	updateCmd.AddCommand(updateContainerCmd)
	updateContainerCmd.PersistentFlags().BoolVarP(&redeployContainerCmd, "redeploy", "r", false, "Redeploy with the current configuration.")
	updateContainerCmd.PersistentFlags().Int64Var(&scaleContainerCmd, "scale", 0, "Scale the container service")
	updateContainerCmd.PersistentFlags().StringVar(&containerNameCmd, "container", "", "The name of the container to update")
	updateContainerCmd.PersistentFlags().StringVar(&containerTagCmd, "tag", "", "The new image tag for the container")
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

		// Check if container update flags are set
		if cmd.Flags().Changed("container") && cmd.Flags().Changed("tag") {
			if j, err = updateContainerImageTag(updateParams, updateResource, containerNameCmd, containerTagCmd, redeployContainerCmd); err != nil {
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
		} else if cmd.Flags().Changed("container") || cmd.Flags().Changed("tag") {
			return errors.New("both --container and --tag must be specified to update the container image")
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

func updateContainerImageTag(params map[string]string, resource *spinup.Resource, containerName, newTag string, forceRedeploy bool) ([]byte, error) {
	log.Infof("updating container %s image tag to %s", containerName, newTag)

	// Get the container service details
	info := &spinup.ContainerService{}
	if err := SpinupClient.GetResource(params, info); err != nil {
		return []byte{}, err
	}

	// Create a copy of the task definition
	taskDefinition := info.TaskDefinition
	containerUpdated := false

	// Loop through container definitions and update the specified container
	for i, container := range taskDefinition.ContainerDefinitions {
		if container.Name == containerName {
			// Parse the current image to get repository and update the tag
			imageParts := strings.Split(container.Image, ":")
			if len(imageParts) < 2 {
				return []byte{}, errors.New("current image format is not valid, expected repository:tag")
			}
			
			// Update the image with the new tag
			taskDefinition.ContainerDefinitions[i].Image = imageParts[0] + ":" + newTag
			containerUpdated = true
			break
		}
	}

	if !containerUpdated {
		return []byte{}, errors.New("container with name " + containerName + " not found in task definition")
	}

	// Create a wrapper for the update input
	updateWrapper := map[string]interface{}{
		"force_redeploy": forceRedeploy || true,
		"size_id":        resource.SizeID,
		"service": map[string]interface{}{
			"container_definitions": taskDefinition.ContainerDefinitions,
			"platform_version": "LATEST",
			"desired_count": info.DesiredCount,
		},
	}

	// Add the capacity provider strategy if it exists
	if len(info.CapacityProviderStrategy) > 0 {
		capProviders := make([]map[string]interface{}, 0, len(info.CapacityProviderStrategy))
		for _, cp := range info.CapacityProviderStrategy {
			capProviders = append(capProviders, map[string]interface{}{
				"base":              cp.Base,
				"capacity_provider": cp.CapacityProvider,
				"weight":            cp.Weight,
			})
		}
		updateWrapper["service"].(map[string]interface{})["capacity_provider_strategy"] = capProviders
	}

	input, err := json.Marshal(updateWrapper)
	if err != nil {
		return []byte{}, err
	}

	log.Debugf("putting input: %s", string(input))

	updatedInfo := &spinup.ContainerService{}
	if err = SpinupClient.PutResource(params, input, updatedInfo); err != nil {
		return []byte{}, err
	}

	return []byte("OK\n"), nil
}