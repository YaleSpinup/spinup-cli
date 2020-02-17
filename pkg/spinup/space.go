package spinup

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Space holds details about a spinup space
type Space struct {
	Id             *FlexInt `json:"id"`
	Name           string   `json:"name"`
	Owner          string   `json:"owner"`
	Department     string   `json:"department"`
	Contact        string   `json:"contact"`
	QuestionaireID string   `json:"questid,omitempty"`
	SecurityGroup  string   `json:"sg,omitempty"`
	Security       string   `json:"security"`
	DataTypes      []struct {
		Id   *FlexInt
		Name string
	} `json:"data_types,omitempty"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at,omitempty"`
	DeletedAt string      `json:"deleted_at,omitempty"`
	Mine      bool        `json:"mine"`
	Resources []*Resource `json:"resources,omitempty"`
	Cost      *SpaceCost  `json:"cost,omitempty"`
}

// GetSpace is a space returned from a wonky endpoint
type GetSpace struct {
	Space *Space `json:"space"`
}

// Spaces is a list of spaces
type Spaces struct {
	Spaces []*Space `json:"spaces"`
}

// SoaceCost is the cost estimate for a space
type SpaceCost struct {
	Amount string
	Unit   string
	End    string
	Start  string
}

// GetEndpoint returns the endpoint to get the list of spaces
func (s *Spaces) GetEndpoint(_ string) string {
	return BaseURL + SpaceURI
}

// GetEndpoint returns the endpoint to get details about a space
func (s *Space) GetEndpoint(id string) string {
	return BaseURL + SpaceURI + "/" + id
}

// GetEndpoint returns the endpoint to get details about a space
func (s *GetSpace) GetEndpoint(id string) string {
	return BaseURL + SpaceURI + "/" + id
}

// GetEndpoint returns the endpoint to get cost of a space
func (s *SpaceCost) GetEndpoint(id string) string {
	return BaseURL + SpaceURI + "/" + id + "/cost"
}

// Resources gets the resources from a space
func (c *Client) Resources(id string) ([]*Resource, error) {
	res, err := c.HTTPClient.Get(BaseURL + SpaceURI + "/" + id)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting space "+id)
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

	log.Debugf("decoded output: %+v", output)

	// remove the type from the output
	for _, r := range output.Resources {
		r.Flavor = r.Type.Flavor
		r.Type = nil
	}

	return output.Resources, nil
}
