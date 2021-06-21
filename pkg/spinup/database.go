package spinup

import log "github.com/sirupsen/logrus"

type DatabaseInfo struct {
	Endpoint    string        `json:",omitempty"`
	DBClusters  []*DBCluster  `json:",omitempty"`
	DBInstances []*DBInstance `json:",omitempty"`
}

type DBCluster struct {
	AllocatedStorage           int64
	AutoMinorVersionUpgrade    bool
	AvailabilityZones          []string
	BackupRetentionPeriod      int64
	Capacity                   int64
	ClusterCreateTime          string
	DBClusterArn               string
	DBClusterIdentifier        string
	DbClusterResourceId        string
	EarliestRestorableTime     string
	Endpoint                   string
	Engine                     string
	EngineMode                 string
	EngineVersion              string
	HostedZoneId               string
	KmsKeyId                   string
	LatestRestorableTime       string
	MasterUsername             string
	MultiAZ                    bool
	PercentProgress            string
	Port                       int64
	PreferredBackupWindow      string
	PreferredMaintenanceWindow string
	ScalingConfigurationInfo   *DBScalingConfiguration
	Status                     string
	StorageEncrypted           bool
}

type DBScalingConfiguration struct {
	AutoPause             bool
	MaxCapacity           int64
	MinCapacity           int64
	SecondsUntilAutoPause int64
	TimeoutAction         string
}

type DBInstance struct {
	AllocatedStorage           int64
	AutoMinorVersionUpgrade    bool
	BackupRetentionPeriod      int64
	CACertificateIdentifier    string
	DBInstanceArn              string
	DBInstanceClass            string
	DBInstanceIdentifier       string
	DBInstanceStatus           string
	DBName                     string
	DbInstancePort             int64
	DbiResourceId              string
	Endpoint                   *DBInstanceEndpoint
	Engine                     string
	EngineVersion              string
	InstanceCreateTime         string
	Iops                       int64
	LatestRestorableTime       string
	LicenseModel               string
	ListenerEndpoint           string
	MasterUsername             string
	MultiAZ                    bool
	PreferredBackupWindow      string
	PreferredMaintenanceWindow string
	PubliclyAccessible         bool
	StorageEncrypted           bool
	StorageType                string
}

type DBInstanceEndpoint struct {
	Address      string
	HostedZoneId string
	Port         int64
}

// GetEndpoint gets the URL for server info
func (s *DatabaseInfo) GetEndpoint(params map[string]string) string {
	return BaseURL + SpaceURI + "/" + params["space"] + "/databases/" + params["name"]
}

// DatabaseSize is the size for a database satisfying the Size interface
type DatabaseSize struct {
	*BaseSize
}

// DatabaseSize returns a DatabaseSize as a Size
func (c *Client) DatabaseSize(id string) (*DatabaseSize, error) {
	size := &DatabaseSize{}
	if err := c.GetResource(map[string]string{"id": id}, size); err != nil {
		return nil, err
	}

	log.Debugf("returning database size %+v", size)

	return size, nil
}
