package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	getCmd.AddCommand(getContainerCmd)
}

var getContainerCmd = &cobra.Command{
	Use:   "container",
	Short: "Get details about a container resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("exactly 1 container service id is required")
		}

		var j []byte
		if detailedGetCmd {
			var err error
			if j, err = containerDetails(args[0]); err != nil {
				return err
			}
		} else {
			var err error
			if j, err = container(args[0]); err != nil {
				return err
			}
		}

		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		f.Write(j)

		return nil
	},
}

func container(id string) ([]byte, error) {
	resource := &spinup.Resource{}
	if err := SpinupClient.GetResource(map[string]string{"id": id}, resource); err != nil {
		return []byte{}, err
	}

	status := resource.Status
	if status != "created" && status != "creating" && status != "deleting" {
		return json.MarshalIndent(struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Status  string `json:"status"`
			SpaceID string `json:"space_id"`
		}{
			ID:      resource.ID.String(),
			Name:    resource.Name,
			Status:  resource.Status,
			SpaceID: resource.SpaceID.String(),
		}, "", "  ")
	}

	size, err := SpinupClient.ContainerSize(resource.SizeID.String())
	if err != nil {
		return []byte{}, err
	}

	// TODO change resource.Name to id once the API is changed to take ID
	info := &spinup.ContainerService{}
	if err = SpinupClient.GetResource(map[string]string{"id": resource.Name}, info); err != nil {
		return []byte{}, err
	}

	return resourceSummary(resource, size, info.Status)
}

func containerDetails(id string) ([]byte, error) {
	resource := &spinup.Resource{}
	if err := SpinupClient.GetResource(map[string]string{"id": id}, resource); err != nil {
		return []byte{}, err
	}

	status := resource.Status
	if status != "created" && status != "creating" && status != "deleting" {
		return json.MarshalIndent(struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Status  string `json:"status"`
			SpaceID string `json:"space_id"`
		}{
			ID:      resource.ID.String(),
			Name:    resource.Name,
			Status:  resource.Status,
			SpaceID: resource.SpaceID.String(),
		}, "", "  ")
	}

	size, err := SpinupClient.ContainerSize(resource.SizeID.String())
	if err != nil {
		return []byte{}, err
	}

	// TODO change resource.Name to id once the API is changed to take ID
	info := &spinup.ContainerService{}
	if err = SpinupClient.GetResource(map[string]string{"id": resource.Name}, info); err != nil {
		return []byte{}, err
	}

	log.Debugf("%+v", info)

	tryit := false
	if size.GetPrice() == "tryit" {
		tryit = true
	}

	output := struct {
		ID       string                   `json:"id"`
		Name     string                   `json:"name"`
		Status   string                   `json:"status"`
		Type     string                   `json:"type"`
		Flavor   string                   `json:"flavor"`
		Security string                   `json:"security"`
		Beta     bool                     `json:"beta"`
		Size     string                   `json:"size"`
		State    string                   `json:"state"`
		SpaceID  string                   `json:"space_id"`
		TryIT    bool                     `json:"tryit"`
		Info     *spinup.ContainerService `json:"info"`
	}{
		ID:       resource.ID.String(),
		Name:     resource.Name,
		Status:   resource.Status,
		Type:     resource.Type.Name,
		Flavor:   resource.Type.Flavor,
		Security: resource.Type.Security,
		SpaceID:  resource.SpaceID.String(),
		Size:     size.GetName(),
		Beta:     resource.Type.Beta.Bool(),
		TryIT:    tryit,
		State:    info.Status,
		Info:     info,
	}

	j, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}
