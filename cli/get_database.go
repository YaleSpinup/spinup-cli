package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	getCmd.AddCommand(getDatabaseCmd)
}

var getDatabaseCmd = &cobra.Command{
	Use:     "database [space]/[name]",
	Short:   "Get a container service",
	PreRunE: getCmdPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("get database: %+v", args)

		status := getResource.Status
		if status != "created" && status != "creating" && status != "deleting" {
			return ingStatus(getResource)
		}

		var err error
		var out []byte
		switch {
		case detailedGetCmd:
			if out, err = databaseDetails(getParams, getResource); err != nil {
				return err
			}
		default:
			if out, err = database(getParams, getResource); err != nil {
				return err
			}
		}

		return formatOutput(out)
	},
}

func database(params map[string]string, resource *spinup.Resource) ([]byte, error) {
	size, err := SpinupClient.DatabaseSize(resource.SizeID.String())
	if err != nil {
		return []byte{}, err
	}

	info := &spinup.DatabaseInfo{}
	if err := SpinupClient.GetResource(params, info); err != nil {
		return []byte{}, err
	}

	status := resource.Status
	if len(info.DBClusters) > 0 {
		status = info.DBClusters[0].Status
		if info.DBClusters[0].EngineMode == "serverless" && info.DBClusters[0].Capacity == 0 {
			status = "paused"
		}
	} else if len(info.DBInstances) > 0 {
		status = info.DBInstances[0].DBInstanceStatus
	}

	return json.MarshalIndent(newResourceSummary(resource, size, status), "", "  ")
}

func databaseDetails(params map[string]string, resource *spinup.Resource) ([]byte, error) {
	size, err := SpinupClient.DatabaseSize(resource.SizeID.String())
	if err != nil {
		return []byte{}, err
	}

	info := &spinup.DatabaseInfo{}
	if err := SpinupClient.GetResource(params, info); err != nil {
		return []byte{}, err
	}

	// I think we only ever have one cluster and instance (even in multi-az deployments)
	var cluster *spinup.DBCluster
	if len(info.DBClusters) > 0 {
		cluster = info.DBClusters[0]
	}

	var instance *spinup.DBInstance
	if len(info.DBInstances) > 0 {
		instance = info.DBInstances[0]
	}

	type DBDetails struct {
		AutoMinorVersionUpgrade bool   `json:"autoMinorVersionUpgrade"`
		CreatedAt               string `json:"createdAt"`
		Engine                  string `json:"engine"`
		EngineVersion           string `json:"engineVersion"`
		MasterUsername          string `json:"masterUsername"`
		BackupRetentionPeriod   int64  `json:"backupRetentionPeriod"`
		BackupWindow            string `json:"backupWindow"`
		MaintenanceWindow       string `json:"maintenanceWindow"`
		Port                    int64  `json:"port"`
		Endpoint                string `json:"endpoint"`
		Encrypted               bool   `json:"encrypted"`
	}

	type ServerlessDetails struct {
		*DBDetails
		AutoPauseEnabled       bool   `json:"autoPauseEnabled"`
		EarliestRestorableTime string `json:"earliestRestorableTime"`
		EngineMode             string `json:"engineMode"`
		LatestRestorableTime   string `json:"latestRestorableTime"`
		CurrentCapacity        int64  `json:"currentCapacity"`
		MaxCapacity            int64  `json:"maxCapacity"`
		MinCapacity            int64  `json:"minCapacity"`
		SecondsToAutoPause     int64  `json:"secondsToAutoPause"`
	}

	type ProvisionedDetails struct {
		*DBDetails
		EarliestRestorableTime string `json:"earliestRestorableTime"`
		EngineMode             string `json:"engineMode"`
		LatestRestorableTime   string `json:"latestRestorableTime"`
		Size                   string `json:"size"`
	}

	type InstanceDetails struct {
		*DBDetails
		AllocatedStorage int64  `json:"allocatedStorage"`
		MultiAZ          bool   `json:"multiAZ"`
		Size             string `json:"size"`
	}

	var status string
	var details interface{}
	if cluster != nil {
		status = cluster.Status

		dbd := &DBDetails{
			AutoMinorVersionUpgrade: cluster.AutoMinorVersionUpgrade,
			CreatedAt:               cluster.ClusterCreateTime,
			Engine:                  cluster.Engine,
			EngineVersion:           cluster.EngineVersion,
			MasterUsername:          cluster.MasterUsername,
			BackupRetentionPeriod:   cluster.BackupRetentionPeriod,
			BackupWindow:            cluster.PreferredBackupWindow,
			MaintenanceWindow:       cluster.PreferredMaintenanceWindow,
			Endpoint:                cluster.Endpoint,
			Encrypted:               cluster.StorageEncrypted,
			Port:                    cluster.Port,
		}

		if cluster.EngineMode == "serverless" {
			if cluster.Capacity == 0 {
				status = "paused"
			}

			details = &ServerlessDetails{
				dbd,
				cluster.ScalingConfigurationInfo.AutoPause,
				cluster.EarliestRestorableTime,
				cluster.EngineMode,
				cluster.LatestRestorableTime,
				cluster.Capacity,
				cluster.ScalingConfigurationInfo.MaxCapacity,
				cluster.ScalingConfigurationInfo.MinCapacity,
				cluster.ScalingConfigurationInfo.SecondsUntilAutoPause,
			}
		} else {
			s := size.Name
			if instance != nil {
				s = instance.DBInstanceClass
			}

			details = &ProvisionedDetails{
				dbd,
				cluster.EarliestRestorableTime,
				cluster.EngineMode,
				cluster.LatestRestorableTime,
				s,
			}
		}
	} else if instance != nil {
		status = instance.DBInstanceStatus

		dbd := &DBDetails{
			AutoMinorVersionUpgrade: instance.AutoMinorVersionUpgrade,
			CreatedAt:               instance.InstanceCreateTime,
			Endpoint:                instance.Endpoint.Address,
			Engine:                  instance.Engine,
			EngineVersion:           instance.EngineVersion,
			MasterUsername:          instance.MasterUsername,
			BackupRetentionPeriod:   instance.BackupRetentionPeriod,
			BackupWindow:            instance.PreferredBackupWindow,
			MaintenanceWindow:       instance.PreferredMaintenanceWindow,
			Encrypted:               instance.StorageEncrypted,
			Port:                    instance.Endpoint.Port,
		}

		details = &InstanceDetails{
			dbd,
			instance.AllocatedStorage,
			instance.MultiAZ,
			instance.DBInstanceClass,
		}

	} else {
		// else assume its a shared mysql database

		e := strings.SplitN(info.Endpoint, ":", 2)
		endpoint := e[0]

		var port int64
		if len(e) == 2 {
			p, err := strconv.ParseInt(e[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s as int: %s", e[1], err)
			}
			port = p
		}

		details = struct {
			Endpoint       string `json:"endpoint"`
			Port           int64  `json:"port"`
			MasterUsername string `json:"masterUsername"`
		}{Endpoint: endpoint, Port: port, MasterUsername: resource.Name}
	}

	output := struct {
		*ResourceSummary
		Details interface{} `json:"details"`
	}{
		newResourceSummary(resource, size, status),
		details,
	}

	j, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}
