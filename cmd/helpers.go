package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
)

type ResourceSummary struct {
	ID       string `json:"id"`
	IP       string `json:"ip,omitempty"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Type     string `json:"type"`
	Flavor   string `json:"flavor"`
	Security string `json:"security"`
	SpaceID  string `json:"space_id,omitempty"`
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
		Timeout: 30 * time.Second,
	}

	s, err := spinup.New(spinupURL, httpClient, spinupToken)
	if err != nil {
		return err
	}

	SpinupClient = s

	return nil
}

// parseSpaceInput takes a list of space arguments and parses them into space ids, converting from
// names to IDs where necessary
func parseSpaceInput(args []string) ([]string, error) {
	log.Debugf("parsing space input args %+v", args)

	var spaceNames []string
	if len(args) > 0 {
		spaceNames = args
	} else if len(spinupSpaces) > 0 {
		spaceNames = spinupSpaces
	} else {
		return nil, errors.New("spaceid(s) or space name(s) required")
	}

	return spaceNames, nil
}

func newResourceSummary(resource *spinup.Resource, size spinup.Size, state string) *ResourceSummary {
	tryit := false
	if size.GetPrice() == "tryit" {
		tryit = true
	}

	return &ResourceSummary{
		ID:       resource.ID.String(),
		IP:       resource.IP,
		Name:     resource.Name,
		Status:   resource.Status,
		Type:     resource.Type.Name,
		Flavor:   resource.Type.Flavor,
		Security: resource.Type.Security,
		Size:     size.GetName(),
		// SpaceID:  resource.SpaceID.String(),
		Beta:  resource.Type.Beta.Bool(),
		TryIT: tryit,
		State: state,
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

// formatOutput prints the output as json or a string
func formatOutput(out interface{}) error {
	output, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}

	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	if _, err := f.Write(output); err != nil {
		return err
	}

	return nil
}
