package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var containerEventsCmd bool
var containerTaskCmd bool

func init() {
	getCmd.AddCommand(getContainerCmd)
	getContainerCmd.PersistentFlags().BoolVarP(&containerEventsCmd, "events", "e", false, "Get container events")
	getContainerCmd.PersistentFlags().BoolVarP(&containerTaskCmd, "tasks", "t", false, "Get container tasks")
}

var getContainerCmd = &cobra.Command{
	Use:   "container",
	Short: "Get details about a container resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("exactly 1 container service id is required")
		}

		resource := &spinup.Resource{}
		if err := SpinupClient.GetResource(map[string]string{"id": args[0]}, resource); err != nil {
			return err
		}

		var j []byte
		var err error
		switch {
		case resource.Status != "created" && resource.Status != "creating" && resource.Status != "deleting":
			if j, err = json.MarshalIndent(struct {
				ID      string `json:"id"`
				Name    string `json:"name"`
				Status  string `json:"status"`
				SpaceID string `json:"space_id"`
			}{
				ID:      resource.ID.String(),
				Name:    resource.Name,
				Status:  resource.Status,
				SpaceID: resource.SpaceID.String(),
			}, "", "  "); err != nil {
				return err
			}
		case detailedGetCmd,
			containerEventsCmd:
			if j, err = containerEvents(resource); err != nil {
				return err
			}
		case containerTaskCmd:
			if j, err = containerTasks(resource); err != nil {
				return err
			}
		default:
			if j, err = container(resource); err != nil {
				return err
			}
		}

		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		f.Write(j)

		return nil
	},
}

func container(resource *spinup.Resource) ([]byte, error) {
	size, err := SpinupClient.ContainerSize(resource.SizeID.String())
	if err != nil {
		return []byte{}, err
	}

	// TODO change resource.Name to id once the API is changed to take ID
	info := &spinup.ContainerService{}
	if err = SpinupClient.GetResource(map[string]string{"id": resource.Name}, info); err != nil {
		return []byte{}, err
	}

	return json.MarshalIndent(newResourceSummary(resource, size, info.Status), "", "  ")
}

func containerDetails(resource *spinup.Resource) ([]byte, error) {
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

	type ContainerDefinition struct {
		Auth         bool              `json:"auth"`
		Image        string            `json:"image"`
		Name         string            `json:"name"`
		Environment  map[string]string `json:"env"`
		PortMappings []string          `json:"portMappings"`
		Secrets      map[string]string `json:"secrets"`
	}

	secrets, err := spaceSecrets(resource.SpaceID.String())
	if err != nil {
		return []byte{}, err
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
			Image:        cdef.Image,
			Name:         cdef.Name,
			PortMappings: portMappings,
			Secrets:      cSecrets,
		})
	}

	type Details struct {
		DesiredCount int64                  `json:"desiredCount"`
		Endpoint     string                 `json:"endpoint"`
		PendingCount int64                  `json:"pendingCount"`
		RunningCount int64                  `json:"runningCount"`
		Containers   []*ContainerDefinition `json:"containers"`
	}

	output := struct {
		*ResourceSummary
		Details *Details `json:"details"`
	}{
		newResourceSummary(resource, size, info.Status),
		&Details{
			DesiredCount: info.DesiredCount,
			Endpoint:     info.ServiceEndpoint,
			PendingCount: info.PendingCount,
			RunningCount: info.RunningCount,
			Containers:   cdefs,
		},
	}

	j, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

func containerEvents(resource *spinup.Resource) ([]byte, error) {
	// TODO change resource.Name to id once the API is changed to take ID
	info := &spinup.ContainerService{}
	if err := SpinupClient.GetResource(map[string]string{"id": resource.Name}, info); err != nil {
		return []byte{}, err
	}

	log.Debugf("%+v", info)

	type Event struct {
		CreatedAt string `json:"createdAt"`
		Id        string `json:"id"`
		Message   string `json:"message"`
	}

	events := make([]*Event, 0, len(info.Events))
	for _, e := range info.Events {
		events = append(events, &Event{
			CreatedAt: e.CreatedAt,
			Id:        e.ID,
			Message:   e.Message,
		})
	}

	j, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

func containerTasks(resource *spinup.Resource) ([]byte, error) {
	// TODO change resource.Name to id once the API is changed to take ID
	info := &spinup.ContainerService{}
	if err := SpinupClient.GetResource(map[string]string{"id": resource.Name}, info); err != nil {
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
		taskOut := &spinup.ContainerTask{}
		if err := SpinupClient.GetResource(map[string]string{"id": resource.Name, "taskId": tid[1]}, taskOut); err != nil {
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
