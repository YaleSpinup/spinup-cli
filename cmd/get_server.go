package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	"github.com/spf13/cobra"
)

func init() {
	getCmd.AddCommand(getServerCmd)
}

var getServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Get details about a server resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("exactly 1 server id is required")
		}

		var j []byte
		if detailedGetCmd {
			var err error
			if j, err = serverDetails(args[0]); err != nil {
				return err
			}
		} else {
			var err error
			if j, err = server(args[0]); err != nil {
				return err
			}
		}

		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		f.Write(j)

		return nil
	},
}

func server(id string) ([]byte, error) {
	params := map[string]string{"id": id}
	resource := &spinup.Resource{}
	if err := SpinupClient.GetResource(params, resource); err != nil {
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

	size, err := SpinupClient.ServerSize(resource.SizeID.String())
	if err != nil {
		return []byte{}, err
	}

	info := &spinup.ServerInfo{}
	if err := SpinupClient.GetResource(params, info); err != nil {
		return []byte{}, err
	}

	return json.MarshalIndent(newResourceSummary(resource, size, info.State), "", "  ")
}

func serverDetails(id string) ([]byte, error) {
	params := map[string]string{"id": id}
	resource := &spinup.Resource{}
	if err := SpinupClient.GetResource(params, resource); err != nil {
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

	size, err := SpinupClient.ServerSize(resource.SizeID.String())
	if err != nil {
		return []byte{}, err
	}

	info := &spinup.ServerInfo{}
	if err := SpinupClient.GetResource(params, info); err != nil {
		return []byte{}, err
	}

	disks := spinup.Disks{}
	if err := SpinupClient.GetResource(params, &disks); err != nil {
		return []byte{}, err
	}

	snapshots := spinup.Snapshots{}
	if err := SpinupClient.GetResource(params, &snapshots); err != nil {
		return []byte{}, err
	}

	sgs := make([]string, 0, len(info.SecurityGroups))
	for _, s := range info.SecurityGroups {
		for k := range s {
			sgs = append(sgs, k)
		}
	}

	type InstanceDetails struct {
		ID               string
		IP               string
		Type             string
		Image            string
		Subnet           string
		SecurityGroups   []string
		AvailabilityZone string
		State            string
	}

	type InstanceVolume struct {
		spinup.Disk
		Snapshots []*spinup.Snapshot `json:"snapshots,omitempty"`
	}

	instanceDisks := []*InstanceVolume{}
	for _, d := range disks {
		volumeSnapshots := []*spinup.Snapshot{}

		for _, s := range snapshots {
			if s.VolumeID == d.ID {
				volumeSnapshots = append(volumeSnapshots, s)
			}
		}

		instanceDisks = append(instanceDisks, &InstanceVolume{
			spinup.Disk{
				ID:          d.ID,
				CreatedAt:   d.CreatedAt,
				Encrypted:   d.Encrypted,
				Size:        d.Size,
				VolumeType:  d.VolumeType,
				Attachments: d.Attachments,
			},
			volumeSnapshots,
		})
	}

	output := struct {
		*ResourceSummary
		InstanceDetails *InstanceDetails  `json:"instance_details"`
		Disks           []*InstanceVolume `json:"disks"`
	}{
		newResourceSummary(resource, size, resource.Status),
		&InstanceDetails{
			ID:               info.ID,
			IP:               info.IP,
			Type:             info.Type,
			Image:            info.Image,
			Subnet:           info.Subnet,
			SecurityGroups:   sgs,
			AvailabilityZone: info.AvailabilityZone,
			State:            info.State,
		},
		instanceDisks,
	}

	j, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}
