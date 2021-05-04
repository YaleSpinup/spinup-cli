package spinup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Resource is a specific resource in the database, it represents an actual instance, container, s3 bucket, etc
type Resource struct {
	Admin     string    `json:"admin,omitempty"`
	CreatedAt string    `json:"created_at"`
	DeletedAt string    `json:"deleted_at,omitempty"`
	TypeName  string    `json:"type_name"`
	ID        *FlexInt  `json:"id"`
	IP        string    `json:"ip,omitempty"`
	IsA       string    `json:"is_a,omitempty"`
	Name      string    `json:"name"`
	ServerID  string    `json:"server_id,omitempty"`
	SizeID    *FlexInt  `json:"size_id,omitempty"`
	SpaceID   *FlexInt  `json:"-"`
	SpaceName string    `json:"space_name"`
	Space     *Space    `json:"-"`
	Status    string    `json:"status"`
	TypeID    *FlexInt  `json:"-"`
	Task      string    `json:"-"`
	Type      *Offering `json:"type,omitempty"`
	UpdatedAt string    `json:"updated_at,omitempty"`
}

// GetEndpoint returns the URL to get a resource
func (r *Resource) GetEndpoint(params map[string]string) string {
	return BaseURL + SpaceURI + "/" + params["space"] + "/resources/" + params["name"]
}

// Resources gets the resources from a space
func (c *Client) Resources(space string) ([]*Resource, error) {
	endpoint := BaseURL + SpaceURI + "/" + space
	log.Infof("getting resources from endpoint: %s", endpoint)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating get request for space %s: %s", space, err)
	}

	req.Header.Set("Content-Type", "application/json")

	if c.AuthToken != "" {
		log.Debugf("setting authorization bearer header")
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting space "+space)
	}

	if res.StatusCode >= 400 {
		return nil, errors.New("error getting space details: " + res.Status)
	}

	log.Infof("got success response from api %s", res.Status)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed reading space body")
	}
	defer res.Body.Close()

	log.Debugf("got response body: %s", string(body))

	output := new(struct {
		Resources []*Resource `json:"resources"`
	})
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, errors.Wrap(err, "failed unmarshalling resources body from json")
	}

	// add space name to return value
	for _, r := range output.Resources {
		r.SpaceName = space
	}

	log.Debugf("decoded output: %+v", output)
	return output.Resources, nil
}
