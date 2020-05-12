package spinup

import (
	"fmt"
	"strconv"
	"strings"
)

// ServerInfo is the details about a server resource, filled in by fetching data from the backend APIs
type ServerInfo struct {
	ID               string                     `json:"id"`
	Name             string                     `json:"name"`
	Type             string                     `json:"type"`
	Image            string                     `json:"image"`
	IP               string                     `json:"ip"`
	Key              string                     `json:"key,omitempty"`
	Subnet           string                     `json:"subnet"`
	Tags             []map[string]string        `json:"tags,omitempty"`
	SecurityGroups   []map[string]string        `json:"sgs"`
	State            string                     `json:"state"`
	AvailabilityZone string                     `json:"az"`
	Platform         string                     `json:"platform,omitempty"`
	CreatedAt        string                     `json:"created_at,omitempty"`
	CreatedBy        string                     `json:"created_by,omitempty"`
	Volumes          map[string]*DiskAttachment `json:"volumes,omitempty"`
}

// DiskAttachment is the details about a disk/volumes attachment to an instance
type DiskAttachment struct {
	AttachTime          string `json:"attach_time"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
	Device              string `json:"device,omitempty"`
	InstanceID          string `json:"instance_id,omitempty"`
	State               string `json:"state,omitempty"`
	Status              string `json:"status,omitempty"`
}

// Disk is a volume
type Disk struct {
	ID          string          `json:"id"`
	CreatedAt   string          `json:"created_at"`
	Encrypted   bool            `json:"encrypted"`
	Size        int             `json:"size"`
	VolumeType  string          `json:"volume_type,omitempty"`
	Attachments *DiskAttachment `json:"attachments,omitempty"`
}

// Disksis a list of disks/volumes
type Disks []*Disk

// Snapshot is a snapshot of a volume
type Snapshot struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Encrypted bool   `json:"encrypted"`
	Progress  string `json:"progress,omitempty"`
	State     string `json:"state,omitempty"`
	Size      int    `json:"size,omitempty"`
	VolumeID  string `json:"volume_id,omitempty"`
}

// Snapshots is a list of snapshots
type Snapshots []*Snapshot

// ServerSize is the size for a server satisfying the Size interface
type ServerSize struct {
	*BaseSize
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// GetEndpoint gets the URL for server info
func (s *ServerInfo) GetEndpoint(params map[string]string) string {
	return BaseURL + ResourceURI + "/" + params["id"] + "/info"
}

// GetEndpoint gets the URL for server disks
func (s *Disks) GetEndpoint(params map[string]string) string {
	return BaseURL + ServerURI + "/" + params["id"] + "/disks"
}

// GetEndpoint gets the URL for server snapshots
func (s *Snapshots) GetEndpoint(params map[string]string) string {
	return BaseURL + ServerURI + "/" + params["id"] + "/snapshots"
}

// ServerSize returns a ServerSize as a Size
func (c *Client) ServerSize(id string) (*ServerSize, error) {
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

	return &ServerSize{size.(*BaseSize), cpu, mem}, nil
}
