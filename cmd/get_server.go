package cmd

import (
	"encoding/json"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	getCmd.AddCommand(getServerCmd)
}

var getServerCmd = &cobra.Command{
	Use:     "server [space][resource]",
	Short:   "Get a server service",
	PreRunE: getCmdPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("get server: %+v", args)

		status := getResource.Status
		if status != "created" && status != "creating" && status != "deleting" {
			return ingStatus(getResource)
		}

		var err error
		var out []byte
		switch {
		case detailedGetCmd:
			if out, err = serverDetails(getParams, getResource); err != nil {
				return err
			}
		default:
			if out, err = server(getParams, getResource); err != nil {
				return err
			}
		}

		return formatOutput(out)
	},
}

func server(params map[string]string, resource *spinup.Resource) ([]byte, error) {
	size, err := SpinupClient.ServerSize(resource.SizeID.String())
	if err != nil {
		return []byte{}, err
	}

	log.Debugf("collected server size: %+v", size)

	info := &spinup.ServerInfo{}
	if err := SpinupClient.GetResource(params, info); err != nil {
		return []byte{}, err
	}

	log.Debugf("collected server info: %+v", info)

	return json.MarshalIndent(newResourceSummary(resource, size, info.State), "", "  ")
}

func serverDetails(params map[string]string, resource *spinup.Resource) ([]byte, error) {
	size, err := SpinupClient.ServerSize(resource.SizeID.String())
	if err != nil {
		return []byte{}, err
	}

	log.Debugf("collected server size: %+v", size)

	info := &spinup.ServerInfo{}
	if err := SpinupClient.GetResource(params, info); err != nil {
		return []byte{}, err
	}

	log.Debugf("collected server info: %+v", info)

	disks := spinup.Disks{}
	if err := SpinupClient.GetResource(params, &disks); err != nil {
		return []byte{}, err
	}

	log.Debugf("collected server disks: %+v", disks)

	snapshots := spinup.Snapshots{}
	if err := SpinupClient.GetResource(params, &snapshots); err != nil {
		return []byte{}, err
	}

	log.Debugf("collected server snapshots: %+v", snapshots)

	sgs := make([]string, 0, len(info.SecurityGroups))
	for _, s := range info.SecurityGroups {
		for k := range s {
			sgs = append(sgs, k)
		}
	}

	type InstanceVolume struct {
		spinup.Disk
		Snapshots []*spinup.Snapshot `json:"snapshots,omitempty"`
	}

	type Details struct {
		AvailabilityZone string            `json:"availability_zone"`
		Disks            []*InstanceVolume `json:"disks"`
		ID               string            `json:"instance_id"`
		Image            string            `json:"image"`
		IP               string            `json:"ip"`
		SecurityGroups   []string          `json:"security_groups"`
		State            string            `json:"state"`
		Subnet           string            `json:"subnet"`
		InstanceType     string            `json:"instance_type"`
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
		Details *Details `json:"details"`
	}{
		newResourceSummary(resource, size, resource.Status),
		&Details{
			AvailabilityZone: info.AvailabilityZone,
			Disks:            instanceDisks,
			ID:               info.ID,
			Image:            info.Image,
			IP:               info.IP,
			SecurityGroups:   sgs,
			State:            info.State,
			Subnet:           info.Subnet,
			InstanceType:     info.Type,
		},
	}

	j, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}
