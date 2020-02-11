package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/YaleSpinup/spinup/pkg/spinup"
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
	resource := &spinup.Resource{}
	err := SpinupClient.GetResource(id, resource)
	if err != nil {
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
	err = SpinupClient.GetResource(id, info)
	if err != nil {
		return []byte{}, err
	}

	return resourceSummary(resource, size, info.State)
}

func serverDetails(id string) ([]byte, error) {
	resource := &spinup.Resource{}
	err := SpinupClient.GetResource(id, resource)
	if err != nil {
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
	err = SpinupClient.GetResource(id, info)
	if err != nil {
		return []byte{}, err
	}

	disks := spinup.Disks{}
	err = SpinupClient.GetResource(id, &disks)
	if err != nil {
		return []byte{}, err
	}

	snapshots := spinup.Snapshots{}
	err = SpinupClient.GetResource(id, &snapshots)
	if err != nil {
		return []byte{}, err
	}

	beta := false
	if b, err := strconv.Atoi(resource.Type.Beta); err != nil && b != 0 {
		beta = true
	}

	sgs := make([]string, 0, len(info.SecurityGroups))
	for _, s := range info.SecurityGroups {
		for k := range s {
			sgs = append(sgs, k)
		}
	}

	tryit := false
	if size.Price == "tryit" {
		tryit = true
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
		ID              string             `json:"id"`
		Name            string             `json:"name"`
		Status          string             `json:"status"`
		Type            string             `json:"type"`
		Flavor          string             `json:"flavor"`
		Security        string             `json:"security"`
		SpaceID         string             `json:"space_id"`
		Beta            bool               `json:"beta"`
		TryIT           bool               `json:"tryit"`
		InstanceDetails *InstanceDetails   `json:"instance_details"`
		Disks           []*InstanceVolume  `json:"disks"`
		Size            *spinup.ServerSize `json:"size"`
	}{
		ID:       resource.ID.String(),
		Name:     resource.Name,
		Status:   resource.Status,
		Type:     resource.Type.Name,
		Flavor:   resource.Type.Flavor,
		Security: resource.Type.Security,
		SpaceID:  resource.SpaceID.String(),
		Beta:     beta,
		TryIT:    tryit,
		InstanceDetails: &InstanceDetails{
			ID:               info.ID,
			IP:               info.IP,
			Type:             info.Type,
			Image:            info.Image,
			Subnet:           info.Subnet,
			SecurityGroups:   sgs,
			AvailabilityZone: info.AvailabilityZone,
			State:            info.State,
		},
		Disks: instanceDisks,
		Size:  size,
	}

	j, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}
