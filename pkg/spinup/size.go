package spinup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// BaseSize contains the basic size information for all sizes
type BaseSize struct {
	ID      *FlexInt `json:"id"`
	Name    string   `json:"name"`
	TypeID  *FlexInt `json:"type_id"`
	Value   string   `json:"value"`
	Details string   `json:"details"`
	Price   string   `json:"price"`
}

// Size is an interface that describes a Spinup size
type Size interface {
	GetName() string
	GetValue() string
	GetPrice() string
}

func (c *Client) Size(id string) (Size, error) {
	endpoint := BaseURL + SizeURI + "/" + id
	log.Infof("getting resource from endpoint: %s", endpoint)

	res, err := c.HTTPClient.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting size "+id)
	}

	if res.StatusCode >= 400 {
		msg := fmt.Sprintf("error getting size (ID: %s): %s", id, res.Status)
		return nil, errors.New(msg)
	}

	log.Infof("got success response from api %s", res.Status)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed reading size body")
	}
	defer res.Body.Close()

	log.Debugf("got response body: %s", string(body))

	size := BaseSize{}
	err = json.Unmarshal(body, &size)
	if err != nil {
		return nil, errors.Wrap(err, "failed unmarshalling size response body from json")
	}

	log.Debugf("decoded output: %+v", size)

	return &size, nil
}

func (s *BaseSize) GetEndpoint(params map[string]string) string {
	return BaseURL + SizeURI + "/" + params["id"]
}

func (s *BaseSize) GetName() string {
	return s.Name
}

func (s *BaseSize) GetValue() string {
	return s.Value
}

func (s *BaseSize) GetPrice() string {
	return s.Price
}
