package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var containerEventsCmd bool
var containerTaskCmd bool

func init() {
	getCmd.AddCommand(getContainerCmd)
	getContainerCmd.PersistentFlags().BoolVar(&containerEventsCmd, "events", false, "Get container events")
	getContainerCmd.PersistentFlags().BoolVar(&containerTaskCmd, "tasks", false, "Get container tasks")
}

var getContainerCmd = &cobra.Command{
	Use:     "container [space]/[resource]",
	Short:   "Get a container service",
	PreRunE: getCmdPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("update container: %+v", args)

		status := getResource.Status
		if status != "created" && status != "creating" && status != "deleting" {
			return ingStatus(getResource)
		}

		var err error
		var out []byte
		switch {
		case detailedGetCmd:
			if out, err = containerDetails(getParams, getResource); err != nil {
				return err
			}
		case containerEventsCmd:
			if out, err = containerEvents(getParams, getResource); err != nil {
				return err
			}
		case containerTaskCmd:
			if out, err = containerTasks(getParams, getResource); err != nil {
				return err
			}
		default:
			if out, err = container(getParams, getResource); err != nil {
				return err
			}
		}

		return formatOutput(out)
	},
}

func container(params map[string]string, resource *spinup.Resource) ([]byte, error) {
	size, err := SpinupClient.ContainerSize(resource.SizeID.String())
	if err != nil {
		return []byte{}, err
	}

	info := &spinup.ContainerService{}
	if err = SpinupClient.GetResource(params, info); err != nil {
		return []byte{}, err
	}

	return json.MarshalIndent(newResourceSummary(resource, size, info.Status), "", "  ")
}

func containerDetails(params map[string]string, resource *spinup.Resource) ([]byte, error) {
	size, err := SpinupClient.ContainerSize(resource.SizeID.String())
	if err != nil {
		return []byte{}, err
	}

	info := &spinup.ContainerService{}
	if err = SpinupClient.GetResource(params, info); err != nil {
		return []byte{}, err
	}

	log.Debugf("collected container info %+v", info)

	spot := false
	for _, c := range info.CapacityProviderStrategy {
		if c.CapacityProvider == "FARGATE_SPOT" {
			spot = true
			break
		}
	}

	log.Debugf("container service spot: %t", spot)

	secrets, err := spaceSecrets(params)
	if err != nil {
		return []byte{}, err
	}

	log.Debugf("collected space secrets %+v", secrets)

	type ContainerDefinition struct {
		Auth         bool                          `json:"auth"`
		Environment  map[string]string             `json:"env,omitempty"`
		HealthCheck  *spinup.ContainerHealthCheck  `json:"healthcheck,omitempty"`
		Image        string                        `json:"image"`
		MountPoints  []*spinup.ContainerMountPoint `json:"mountpoints,omitempty"`
		Name         string                        `json:"name"`
		PortMappings []string                      `json:"portMappings,omitempty"`
		Secrets      map[string]string             `json:"secrets,omitempty"`
	}

	cdefs := make([]*ContainerDefinition, 0, len(info.TaskDefinition.ContainerDefinitions))
	for _, cdef := range info.TaskDefinition.ContainerDefinitions {
		auth := false
		if cdef.RepositoryCredentials.CredentialsParameter != "" {
			auth = true
		}

		env, err := mapNameValueArray(cdef.Environment)
		if err != nil {
			return []byte{}, err
		}

		cSecrets := make(map[string]string)
		if len(cdef.Secrets) > 0 {
			// map the secrets for the container def
			cSecrets, err = mapNameValueFromArray(cdef.Secrets)
			if err != nil {
				return []byte{}, err
			}

			// if the ARN matches the value of the container secret, override the value with the spinup secret name
			for _, s := range secrets {
				for k, v := range cSecrets {
					if v == s.ARN {
						cSecrets[k] = s.Name
					}
				}
			}
		}

		portMappings := []string{}
		if cdef.PortMappings != nil {
			for _, p := range cdef.PortMappings {
				portMappings = append(portMappings, fmt.Sprintf("%d/%s", p.ContainerPort, p.Protocol))
			}
		}

		cdefs = append(cdefs, &ContainerDefinition{
			Auth:         auth,
			Environment:  env,
			HealthCheck:  cdef.HealthCheck,
			Image:        cdef.Image,
			MountPoints:  cdef.MountPoints,
			Name:         cdef.Name,
			PortMappings: portMappings,
			Secrets:      cSecrets,
		})
	}

	type ContainerVolume struct {
		Name      string `json:"name"`
		Type      string `json:"type"`
		NfsVolume string `json:"nfs_volume,omitempty"`
	}

	volumes := make([]*ContainerVolume, 0, len(info.TaskDefinition.Volumes))
	for _, volume := range info.TaskDefinition.Volumes {
		v := ContainerVolume{
			Name: volume.Name,
		}

		v.Type = "persistent"
		if volume.Host != nil {
			v.Type = "ephemeral"
		}

		// TODO determine spinup resource instead of FileSystemId
		if volume.EfsVolumeConfiguration != nil {
			v.NfsVolume = volume.EfsVolumeConfiguration.FileSystemId
		}

		volumes = append(volumes, &v)
	}

	type Details struct {
		Containers   []*ContainerDefinition `json:"containers"`
		DesiredCount int64                  `json:"desiredCount"`
		Endpoint     string                 `json:"endpoint"`
		PendingCount int64                  `json:"pendingCount"`
		RunningCount int64                  `json:"runningCount"`
		Spot         bool                   `json:"spot"`
		Volumes      []*ContainerVolume     `json:"volumes"`
	}

	output := struct {
		*ResourceSummary
		Details *Details `json:"details"`
	}{
		newResourceSummary(resource, size, info.Status),
		&Details{
			Containers:   cdefs,
			DesiredCount: info.DesiredCount,
			Endpoint:     info.ServiceEndpoint,
			PendingCount: info.PendingCount,
			RunningCount: info.RunningCount,
			Spot:         spot,
			Volumes:      volumes,
		},
	}

	j, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

func containerEvents(params map[string]string, resource *spinup.Resource) ([]byte, error) {
	info := &spinup.ContainerService{}
	if err := SpinupClient.GetResource(params, info); err != nil {
		return []byte{}, err
	}

	log.Debugf("%+v", info)

	type Event struct {
		CreatedAt string `json:"createdAt"`
		Id        string `json:"id"`
		Message   string `json:"message"`
	}

	events := make([]*Event, 0, len(info.Events))
	for i := len(info.Events) - 1; i >= 0; i-- {
		e := info.Events[i]
		events = append(events, &Event{
			CreatedAt: e.CreatedAt,
			Id:        e.ID,
			Message:   e.Message,
		})
	}

	output := struct {
		Events []*Event `json:"events"`
	}{events}

	j, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

func containerTasks(params map[string]string, resource *spinup.Resource) ([]byte, error) {
	info := &spinup.ContainerService{}
	if err := SpinupClient.GetResource(params, info); err != nil {
		return []byte{}, err
	}

	log.Debugf("%+v", info)

	type Container struct {
		ExitCode     string `json:"exitCode"`
		HealthStatus string `json:"healthStatus"`
		Image        string `json:"image"`
		LastStatus   string `json:"lastStatus"`
		Name         string `json:"name"`
		Reason       string `json:"reason"`
	}

	type Task struct {
		AvailabilityZone string       `json:"availabilityZone"`
		CapacityProvider string       `json:"capacityProvider"`
		CPU              string       `json:"cpu"`
		CreatedAt        string       `json:"createdAt"`
		Id               string       `json:"id"`
		IpAddress        string       `json:"ipAddress"`
		LastStatus       string       `json:"lastStatus"`
		LaunchType       string       `json:"launchType"`
		Memory           string       `json:"memory"`
		PlatformVersion  string       `json:"platformVersion"`
		PullStartedAt    string       `json:"pullStartedAt"`
		PullStoppedAt    string       `json:"pullStoppedAt"`
		StopCode         string       `json:"stopCode"`
		StoppedAt        string       `json:"stoppedAt"`
		StoppedReason    string       `json:"stoppedReason"`
		StoppingAt       string       `json:"stoppingAt"`
		Containers       []*Container `json:"containers"`
		Version          int64        `json:"version"`
	}

	tasks := make([]*Task, 0, len(info.Tasks))
	for _, t := range info.Tasks {
		tid := strings.SplitN(t, "/", 2)
		params["taskId"] = tid[1]
		taskOut := &spinup.ContainerTask{}
		if err := SpinupClient.GetResource(params, taskOut); err != nil {
			return []byte{}, err
		}

		for _, task := range taskOut.Tasks {
			var ip string
			for _, a := range task.Attachments {
				if a.Type == "ElasticNetworkInterface" {
					for _, nv := range a.Details {
						if nv.Name == "privateIPv4Address" {
							ip = nv.Value
						}
					}
				}
			}

			containers := make([]*Container, 0, len(task.Containers))
			for _, c := range task.Containers {
				containers = append(containers, &Container{
					ExitCode:     c.ExitCode,
					HealthStatus: c.HealthStatus,
					Image:        c.Image,
					LastStatus:   c.LastStatus,
					Name:         c.Name,
					Reason:       c.Reason,
				})
			}

			tasks = append(tasks, &Task{
				AvailabilityZone: task.AvailabilityZone,
				CapacityProvider: task.CapacityProviderName,
				CPU:              task.Cpu,
				CreatedAt:        task.CreatedAt,
				Id:               tid[1],
				IpAddress:        ip,
				LastStatus:       task.LastStatus,
				LaunchType:       task.LaunchType,
				Memory:           task.Memory,
				PlatformVersion:  task.PlatformVersion,
				PullStartedAt:    task.PullStartedAt,
				PullStoppedAt:    task.PullStoppedAt,
				StopCode:         task.StopCode,
				StoppedAt:        task.StoppedAt,
				StoppedReason:    task.StoppedReason,
				StoppingAt:       task.StoppingAt,
				Containers:       containers,
				Version:          task.Version,
			})
		}
	}

	output := struct {
		Tasks []*Task `json:"tasks"`
	}{
		tasks,
	}

	j, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}
