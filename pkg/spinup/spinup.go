package spinup

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	BaseURL      = "http://localhost:8090"
	SpaceURI     = "/api/v2/spaces"
	ResourceURI  = "/api/v2/resources"
	ServerURI    = "/api/v2/servers"
	SizeURI      = "/api/v2/sizes"
	ContainerURI = "/api/v2/containers"
	StorageURI   = "/api/v2/storage"
)

// FlexInt is an int... or a string... or an int.... or...
type FlexInt int

// Client is the spinup client
type Client struct {
	HTTPClient *http.Client
}

// ResourceType is an interface for deteriming URLs
type ResourceType interface {
	GetEndpoint(id string) string
}

// Resource gets details about a resource
func (c *Client) GetResource(id string, r ResourceType) error {
	res, err := c.HTTPClient.Get(r.GetEndpoint(id))
	if err != nil {
		return errors.Wrap(err, "failed getting resource "+id)
	}

	if res.StatusCode > 400 {
		return errors.New("error getting resource details: " + res.Status)
	}

	log.Infof("got success response from api %s", res.Status)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "failed reading resource body")
	}
	defer res.Body.Close()

	log.Debugf("go response body: %s", string(body))

	err = json.Unmarshal(body, r)
	if err != nil {
		return errors.Wrap(err, "failed unmarshalling resource body from json")
	}

	log.Debugf("decoded output: %+v", r)

	return nil
}
