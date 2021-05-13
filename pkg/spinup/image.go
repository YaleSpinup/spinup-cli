package spinup

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// Image is a server image
type Image struct {
	Architecture string       `json:"architecture,omitempty"`
	CreatedAt    string       `json:"created_at,omitempty"`
	CreatedBy    string       `json:"created_by,omitempty"`
	Description  string       `json:"description,omitempty"`
	ID           string       `json:"id"`
	Name         string       `json:"name,omitempty"`
	ServerName   string       `json:"server_name,omitempty"`
	State        string       `json:"state,omitempty"`
	Status       string       `json:"status,omitempty"`
	Volumes      ImageVolumes `json:"volumes,omitempty"`
	Offering     *Offering    `json:"offering,omitempty"`
}

// Images is a list of server images
type Images []*Image

// GetEndpoint gets the endpoint UR for an image list
func (i *Images) GetEndpoint(params map[string]string) string {
	return BaseURL + SpaceURI + "/" + params["id"] + "/images"
}

type ImageVolumes map[string]*ImageVolume
type ImageVolume struct {
	DeleteOnTermination bool   `json:"delete_on_termination"`
	Encrypted           bool   `json:"encrypted"`
	ID                  string `json:"snapshot_id,omitempty"`
	Size                int    `json:"volume_size,omitempty"`
	Type                string `json:"volume_type,omitempty"`
}

func (iv *ImageVolumes) UnmarshalJSON(b []byte) error {
	if *iv == nil {
		*iv = ImageVolumes{}
	}
	imageVolumes := *iv

	var volMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &volMap); err != nil {
		return err
	}

	for k, v := range volMap {
		if string(v) == "[]" {
			log.Debugf("skipping empty volume %s", k)
			continue
		}

		vol := ImageVolume{}
		if err := json.Unmarshal(v, &vol); err != nil {
			return err
		}

		imageVolumes[k] = &vol
	}

	return nil
}
