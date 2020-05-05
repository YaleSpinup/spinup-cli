package spinup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

// FlexBool is a bool... or a stirng... or an int... or...
type FlexBool bool

// Client is the spinup client
type Client struct {
	HTTPClient *http.Client
}

// ResourceType is an interface for deteriming URLs
type ResourceType interface {
	GetEndpoint(params map[string]string) string
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

// Resource gets details about a resource and unmarshals them them into the passed
// ResourceType.  It first gets the resource endpoint by calling r.GetEndpoint(id) which
// is a function on the passed ResourceType interface.
func (c *Client) GetResource(params map[string]string, r ResourceType) error {
	res, err := c.HTTPClient.Get(r.GetEndpoint(params))
	if err != nil {
		return fmt.Errorf("failed getting resource with params %+v: %s", params, err)
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("error getting resource details: %s", res.Status)
	}

	log.Infof("got success response from api %s", res.Status)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed reading resource body: %s", err)
	}
	defer res.Body.Close()

	log.Debugf("go response body: %s", string(body))

	err = json.Unmarshal(body, r)
	if err != nil {
		return fmt.Errorf("failed unmarshalling resource body from json: %s", err)
	}

	log.Debugf("decoded output: %+v", r)

	return nil
}

func (fi *FlexInt) UnmarshalJSON(b []byte) error {
	// if b is not a string, it's an int
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

func (fb *FlexBool) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	sb, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	*fb = FlexBool(sb)
	return nil
}

func (fb *FlexBool) Bool() bool {
	log.Debugf("converting flex bool to bool: %v", *fb)
	return bool(*fb)
}

func (fb *FlexBool) String() string {
	log.Debugf("converting flex bool to string: %v", *fb)
	return strconv.FormatBool(bool(*fb))
}
