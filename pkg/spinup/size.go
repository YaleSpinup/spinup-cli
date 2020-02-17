package spinup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type BaseSize struct {
	ID      *FlexInt `json:"id"`
	Name    string   `json:"name"`
	TypeID  string   `json:"type_id"`
	Value   string   `json:"value"`
	Details string   `json:"details"`
	Price   string   `json:"price"`
}

type Size interface {
	GetName() string
	GetValue() string
	GetPrice() string
}

func (c *Client) Size(id string) (Size, error) {
	res, err := c.HTTPClient.Get(BaseURL + SizeURI + "/" + id)
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

func (s *BaseSize) GetName() string {
	return s.Name
}

func (s *BaseSize) GetValue() string {
	return s.Value
}

func (s *BaseSize) GetPrice() string {
	return s.Price
}
