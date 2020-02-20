package spinup

import (
	"fmt"
	"strconv"
	"strings"
)

// ContainerService is a spinup container service
type ContainerService struct {
	ClusterArn   string
	CreatedAt    string
	DesiredCount int
	Events       []struct {
		CreateAt string
		ID       string
		Message  string
	}
	NetworkConfiguration struct {
		AwsvpcConfiguration struct {
			AssignPublicIp string
			SecurityGroups []string
			Subnets        []string
		}
	}
	PendingCount       int
	RoleArn            string
	RunningCount       int
	SchedulingStrategy string
	ServiceArn         string
	ServiceName        string
	ServiceRegistries  []struct {
		ContainerName string
		ContainerPort int
		Port          int
		RegistryArn   string
	}
	Status         string
	Tasks          []string
	TaskDefinition struct {
		CPU                  string
		ContainerDefinitions []struct {
			Image                 string
			Name                  string
			ContainerPortMappings []struct {
				ContainerPort int
				HostPort      int
				Protocol      int
			}

			Environment []struct {
				Name  string
				Value string
			}

			LogConfiguration struct {
				LogDriver     string
				Options       map[string]string
				SecretOptions []*Secret
			}
			Secrets []*Secret
		}
		Family            string
		Memory            string
		Revision          int
		Status            string
		TaskDefinitionArn string
		Volumes           []struct{}
	}
}

// GetEndpoint returns the endpoint to get details about a container service
func (c *ContainerService) GetEndpoint(id string) string {
	return BaseURL + ContainerURI + "/" + id
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
