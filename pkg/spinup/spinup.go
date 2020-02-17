package spinup

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

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

func New(spinupUrl string, client *http.Client) (*Client, error) {
	u, err := url.Parse(spinupUrl)
	if err != nil {
		return nil, err
	}

	BaseURL = u.String()
	return &Client{
		// BaseURL:    u,
		HTTPClient: client,
	}, nil
}

// Resource gets details about a resource and unmarshals them them into the passed ResourceType
func (c *Client) GetResource(id string, r ResourceType) error {
	res, err := c.HTTPClient.Get(r.GetEndpoint(id))
	if err != nil {
		return errors.Wrap(err, "failed getting resource "+id)
	}

	if res.StatusCode >= 400 {
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

func (fi *FlexInt) UnmarshalJSON(b []byte) error {
	if b[0] != '"' {
		return json.Unmarshal(b, (*int)(fi))
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*fi = FlexInt(i)
	return nil
}

func (fi *FlexInt) String() string {
	log.Debugf("converting flex int to string: %v", *fi)
	return strconv.Itoa(int(*fi))
}
