package spinup

import log "github.com/sirupsen/logrus"

// S3StorageInfo is the info about a S3 storage bucket
type S3StorageInfo struct {
	Empty bool
}

// S3StorageUsers is a list of storage users
type S3StorageUsers []*S3StorageUser

// S3StorageUser is a storage user
type S3StorageUser struct {
	Arn        string                    `json:"Arn"`
	Username   string                    `json:"UserName"`
	CreatedAt  string                    `json:"CreateDate"`
	LastUsed   string                    `json:"PasswordLastUsed"`
	AccessKeys []*S3StorageUserAccessKey `json:"AccessKeys"`
}

type S3StorageUserAccessKey struct {
	AccessKeyId string
	CreateDate  string
	Status      string
	UserName    string
}

// GetEndpoint returns the url for a storage resource
func (s *S3StorageInfo) GetEndpoint(params map[string]string) string {
	return BaseURL + SpaceURI + "/" + params["space"] + "/storage/" + params["name"]
}

// S3StorageSize is the size for a container satisfying the Size interface
type S3StorageSize struct {
	*BaseSize
}

// S3StorageSize returns S3StorageSize as a Size
func (c *Client) S3StorageSize(id string) (*S3StorageSize, error) {
	size := &S3StorageSize{}
	if err := c.GetResource(map[string]string{"id": id}, size); err != nil {
		return nil, err
	}

	log.Debugf("returing s3 storage size %+v", size)

	return size, nil
}

// GetEndpoint returns the URL for the list of users of a storage resource
func (s *S3StorageUsers) GetEndpoint(params map[string]string) string {
	return BaseURL + SpaceURI + "/" + params["space"] + "/storage/" + params["name"] + "/users"
}

// GetEndpoint returns the URL for the details about a user of a storage resource
func (s *S3StorageUser) GetEndpoint(params map[string]string) string {
	return BaseURL + SpaceURI + "/" + params["space"] + "/storage/" + params["name"] + "/users/" + params["username"]
}
