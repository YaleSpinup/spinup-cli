package spinup

import (
	"fmt"
	"strconv"
	"strings"
)

type ContainerEvent struct {
	CreatedAt string
	ID        string
	Message   string
}

type ContainerHealthCheck struct {
	Command     []string
	Interval    int64
	Retries     int64
	StartPeriod int64
	Timeout     int64
}

type ContainerMountPoint struct {
	ContainerPath string
	ReadOnly      bool
	SourceVolume  string
}

type ContainerPortMapping struct {
	ContainerPort int64
	HostPort      int64
	Protocol      string
}

type Container struct {
	ContainerArn      string
	Cpu               string
	ExitCode          string
	HealthStatus      string
	Image             string
	LastStatus        string
	Memory            string
	MemoryReservation string
	Name              string
	NetworkBindings   []struct {
		BindIP        string
		ContainerPort int64
		HostPort      int64
		Protocol      string
	}
	NetworkInterfaces []struct {
		AttachmentId       string
		Ipv6Address        string
		PrivateIpv4Address string
	}
	Reason    string
	RuntimeId string
	TaskArn   string
}

type ContainerTask struct {
	Failures []string
	Tasks    []struct {
		AvailabilityZone string
		Attachments      []struct {
			Details []*NameValue
			Id      string
			Status  string
			Type    string
		}
		CapacityProviderName  string
		ClusterArn            string
		Connectivity          string
		ConnectivityAt        string
		Containers            []*Container
		Cpu                   string
		CreatedAt             string
		DesiredStatus         string
		ExecutionStoppedAt    string
		Group                 string
		HealthStatus          string
		InferenceAccelerators []struct {
			DeviceName string
			DeviceType string
		}
		LastStatus        string
		LaunchType        string
		Memory            string
		Overrides         interface{}
		PlatformVersion   string
		PullStartedAt     string
		PullStoppedAt     string
		StartedAt         string
		StartedBy         string
		StopCode          string
		StoppedAt         string
		StoppedReason     string
		StoppingAt        string
		Tags              []*NameValue
		TaskArn           string
		TaskDefinitionArn string
		Version           int64
	}
}

type ContainerDefinition struct {
	Command   []string
	CPU       int64
	DependsOn []struct {
		Condition     string
		ContainerName string
	}
	DisableNetworking bool
	DnsSearchDomains  []string
	DnsServers        []string
	DockerLabels      map[string]string
	EntryPoint        []string
	Environment       []*NameValue
	Essential         bool
	HealthCheck       *ContainerHealthCheck
	Image             string
	Interactive       bool
	Links             []string
	// LinuxParameters *LinuxParameter
	LogConfiguration struct {
		LogDriver     string
		Options       map[string]string
		SecretOptions []*NameValueFrom
	}
	Memory                 int64
	MemoryReservation      int64
	MountPoints            []*ContainerMountPoint
	Name                   string
	PortMappings           []*ContainerPortMapping
	Privileged             bool
	PseudoTerminal         bool
	ReadonlyRootFilesystem bool
	RepositoryCredentials  struct {
		CredentialsParameter string
	}
	Secrets      []*NameValueFrom
	StartTimeout int64
	StopTimeout  int64
	Ulimits      []struct {
		HardLimit int64
		Name      string
		SoftLimit int64
	}
	User        string
	VolumesFrom []struct {
		ReadOnly        bool
		SourceContainer string
	}
	WorkingDirectory string
}

// ContainerService is a spinup container service
type ContainerService struct {
	ClusterArn    string
	CreatedAt     string
	DesiredCount  int64
	Events        []*ContainerEvent
	LoadBalancers []struct {
		ContainerName    string
		ContainerPort    int64
		LoadBalancerName string
		TargetGroupArn   string
	}
	NetworkConfiguration struct {
		AwsvpcConfiguration struct {
			AssignPublicIp string
			SecurityGroups []string
			Subnets        []string
		}
	}
	PendingCount       int64
	RoleArn            string
	RunningCount       int64
	SchedulingStrategy string
	ServiceArn         string
	ServiceEndpoint    string
	ServiceName        string
	ServiceRegistries  []struct {
		ContainerName string
		ContainerPort int64
		Port          int64
		RegistryArn   string
	}
	Status         string
	Tasks          []string
	TaskDefinition struct {
		Compatibilities      []string
		CPU                  string
		ContainerDefinitions []*ContainerDefinition
		Family               string
		Memory               string
		Revision             int64
		Status               string
		TaskDefinitionArn    string
		Volumes              []struct{}
	}
}

// GetEndpoint returns the endpoint to get details about a container service
func (c *ContainerService) GetEndpoint(params map[string]string) string {
	return BaseURL + ContainerURI + "/" + params["id"]
}

// ContainerSize is the size for a container satisfying the Size interface
type ContainerSize struct {
	*BaseSize
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// ContainerSize returns ContainerSize as a Size
func (c *Client) ContainerSize(id string) (Size, error) {
	size, err := c.Size(id)
	if err != nil {
		return nil, err
	}

	cpu, mem := "", ""
	if size.GetValue() != "" {
		v := strings.SplitN(size.GetValue(), "-", 2)
		c, err := strconv.ParseFloat(v[0], 64)
		if err != nil {
			return nil, err
		}

		m, err := strconv.ParseFloat(v[1], 64)
		if err != nil {
			return nil, err
		}

		cpu = fmt.Sprintf("%0.00f vCPU", c/1024)
		mem = fmt.Sprintf("%0.00f GB", m/1024)
	}

	return &ContainerSize{size.(*BaseSize), cpu, mem}, nil
}

// GetEndpoint returns the endpoint to get details about a container service task
func (c *ContainerTask) GetEndpoint(params map[string]string) string {
	return BaseURL + ContainerURI + "/" + params["id"] + "/tasks/" + params["taskId"]
}
