package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/YaleSpinup/spinup-cli/pkg/cas"
	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	"golang.org/x/net/publicsuffix"
)

type ResourceSummary struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Type     string `json:"type"`
	Flavor   string `json:"flavor"`
	Security string `json:"security"`
	SpaceID  string `json:"space_id"`
	Beta     bool   `json:"beta"`
	Size     string `json:"size"`
	TryIT    bool   `json:"tryit"`
	State    string `json:"state,omitempty"`
}

func initClient() error {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return err
	}

	httpClient := &http.Client{
		Jar:     jar,
		Timeout: 15 * time.Second,
	}

	if err = cas.Auth(spinupUser, spinupPass, spinupURL+"/login", httpClient); err != nil {
		return err
	}

	s, err := spinup.New(spinupURL, httpClient)
	if err != nil {
		return err
	}

	SpinupClient = s

	return nil
}

func parseSpaceInput(args []string) ([]string, error) {
	spaceIds := []string{}
	if len(args) > 0 {
		spaces := spinup.Spaces{}
		if err := SpinupClient.GetResource(map[string]string{}, &spaces); err != nil {
			return nil, err
		}

		for _, s := range spaces.Spaces {
			for _, arg := range args {
				if strings.EqualFold(s.Name, arg) {
					spaceIds = append(spaceIds, s.Id.String())
				}
			}
		}
	} else if len(spinupSpaceIDs) > 0 {
		spaceIds = spinupSpaceIDs
	} else {
		return nil, errors.New("spaceid(s) or space name(s) required")
	}

	return spaceIds, nil
}

func newResourceSummary(resource *spinup.Resource, size spinup.Size, state string) *ResourceSummary {
	tryit := false
	if size.GetPrice() == "tryit" {
		tryit = true
	}

	return &ResourceSummary{
		ID:       resource.ID.String(),
		Name:     resource.Name,
		Status:   resource.Status,
		Type:     resource.Type.Name,
		Flavor:   resource.Type.Flavor,
		Security: resource.Type.Security,
		Size:     size.GetName(),
		SpaceID:  resource.SpaceID.String(),
		Beta:     resource.Type.Beta.Bool(),
		TryIT:    tryit,
		State:    state,
	}
}

// mapNameValueArray maps the ubiquitous Name Value array into a key:value map
func mapNameValueArray(input []*spinup.NameValue) (map[string]string, error) {
	output := map[string]string{}
	for _, s := range input {
		if _, ok := output[s.Name]; ok {
			return nil, fmt.Errorf("name collision mapping name value: %s", s.Name)
		}
		output[s.Name] = s.Value
	}
	return output, nil
}

// mapNameValueArray maps the ubiquitous Name Value array into a key:value map
func mapNameValueFromArray(input []*spinup.NameValueFrom) (map[string]string, error) {
	output := map[string]string{}
	for _, s := range input {
		if _, ok := output[s.Name]; ok {
			return nil, fmt.Errorf("name collision mapping name valuefrom: %s", s.Name)
		}
		output[s.Name] = s.ValueFrom
	}
	return output, nil
}
