package spinup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	BaseURL      = "http://localhost:8090"
	ContainerURI = "/api/v3/containers"
	ResourceURI  = "/api/v3/resources"
	SecretsURI   = "/api/v3/spaces"
	SizeURI      = "/api/v3/sizes"
	SpaceURI     = "/api/v3/spaces"
	StorageURI   = "/api/v3/storage"
)

// FlexInt is an int... or a string... or an int.... or...
type FlexInt int

// FlexBool is a bool... or a stirng... or an int... or...
type FlexBool bool

// Client is the spinup client
type Client struct {
	AuthToken  string
	CSRFToken  string
	HTTPClient *http.Client
}

// NameValue is the ubuquitous Name/Value struct
type NameValue struct {
	Name  string
	Value string
}

// NameValueFrom is a Name/ValueFrom struct
type NameValueFrom struct {
	Name      string
	ValueFrom string
}

// ResourceType is an interface for deteriming URLs
type ResourceType interface {
	GetEndpoint(params map[string]string) string
}

func New(spinupUrl string, client *http.Client, token string) (*Client, error) {
	u, err := url.Parse(spinupUrl)
	if err != nil {
		return nil, err
	}

	BaseURL = u.String()
	return &Client{
		AuthToken: token,
		// BaseURL:    u,
		HTTPClient: client,
	}, nil
}

// GetResource gets details about a resource and unmarshals them them into the passed
// ResourceType.  It first gets the resource endpoint by calling r.GetEndpoint(id) which
// is a function on the passed ResourceType interface.
func (c *Client) GetResource(params map[string]string, r ResourceType) error {
	defer timeTrack(time.Now(), "GetResource")

	endpoint := r.GetEndpoint(params)
	log.Infof("getting resource from endpoint: %s", endpoint)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed creating get resource request with params %+v: %s", params, err)
	}

	req.Header.Set("Content-Type", "application/json")

	if c.AuthToken != "" {
		log.Debugf("setting authorization bearer header")
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed getting resource with params %+v: %s", params, err)
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("error getting resource details: %s", res.Status)
	}

	log.Infof("got success response from api %s", res.Status)

	if log.GetLevel() == log.DebugLevel {
		dump, err := httputil.DumpResponse(res, true)
		if err != nil {
			log.Fatal(err)
		}

		log.Debugf("response: %s", string(dump))
	}

	for _, cookie := range res.Cookies() {
		if cookie.Name == "XSRF-TOKEN" {
			log.Debugf("found xsrf-token: %+v", cookie)

			decodedValue, err := url.QueryUnescape(cookie.Value)
			if err != nil {
				return err
			}

			log.Debugf("XSRF-TOKEN cookie %+v", decodedValue)
			c.CSRFToken = decodedValue
		}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed reading resource body: %s", err)
	}
	defer res.Body.Close()

	log.Debugf("read response body: %s", string(body))

	err = json.Unmarshal(body, r)
	if err != nil {
		return fmt.Errorf("failed unmarshalling resource body from json: %s", err)
	}

	log.Debugf("decoded output: %+v", r)

	return nil
}

// PutResources updates a resource
func (c *Client) PutResource(params map[string]string, input []byte, r ResourceType) error {
	defer timeTrack(time.Now(), "PutResource")

	endpoint := r.GetEndpoint(params)
	log.Infof("putting resource to endpoint: %s", endpoint)

	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewBuffer(input))
	if err != nil {
		return fmt.Errorf("failed creating update request with params %+v, %s: %s", params, string(input), err)
	}

	req.Header.Set("Content-Type", "application/json")

	if c.AuthToken != "" {
		log.Debugf("setting authorization bearer header")
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed updating resource with params %+v, %s: %s", params, string(input), err)
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("error updating resource: %s", res.Status)
	}

	log.Infof("got success response from api %s", res.Status)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed reading resource body: %s", err)
	}
	defer res.Body.Close()

	log.Debugf("got response body: %s", string(body))

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

func (nv *NameValue) String() string {
	log.Debugf("returning name/value as string: %v", *nv)
	return fmt.Sprintf("%s:%s", nv.Name, nv.Value)
}

func (nv *NameValueFrom) String() string {
	log.Debugf("returning name/value from as string: %v", *nv)
	return fmt.Sprintf("%s:%s", nv.Name, nv.ValueFrom)
}

// timeTrack logs the time since the passed time
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Infof("%s took %s", name, elapsed)
}
