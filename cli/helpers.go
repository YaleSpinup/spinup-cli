package cli

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
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
	defer timeTrack(time.Now(), "initClient()")

	if err := validateToken(spinupToken); err != nil {
		return err
	}

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

// findResourceInSpaces returns the space for the given resource, searching the spaces passed in the space list
func findResourceInSpaces(name string, spaces []string) (string, error) {
	log.Debugf("finding %s in spaces %+v", name, spaces)

	for _, s := range spinupSpaces {
		log.Debugf("listing resources for space %s", s)

		resources, err := SpinupClient.Resources(s)
		if err != nil {
			return "", err
		}

		for _, r := range resources {
			if r.Name == name {
				return s, nil
			}
		}
	}

	return "", fmt.Errorf("resource %s not found in any spaces", name)
}

// ingStatus prints the basic information about a resource and returns.
func ingStatus(resource *spinup.Resource) error {
	out, err := json.MarshalIndent(struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Status  string `json:"status"`
		SpaceID string `json:"space_id"`
	}{
		ID:      resource.ID.String(),
		Name:    resource.Name,
		Status:  resource.Status,
		SpaceID: resource.SpaceID.String(),
	}, "", "  ")
	if err != nil {
		return err
	}

	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	f.Write(out)

	return nil
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
	var output []byte
	o, ok := out.([]byte)
	if ok {
		output = o
	} else {
		var err error
		output, err = json.MarshalIndent(out, "", "  ")
		if err != nil {
			return err
		}
	}

	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	if _, err := f.Write(output); err != nil {
		return err
	}

	return nil
}

func validateToken(tokenString string) error {
	defer timeTrack(time.Now(), "validateToken()")

	if tokenString == "" {
		log.Debug("no token provided")
		return nil
	}

	log.Debugf("validating token: %s", tokenString)

	parts := strings.SplitN(tokenString, ".", 3)
	if l := len(parts); l != 3 {
		return fmt.Errorf("invalid token, unexpected number of parts (%d)", l)
	}

	rawPayload, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("invalid token, unable to decode payload: %s", err)
	}

	log.Debugf("got payload: %s", string(rawPayload))

	var payload map[string]interface{}
	if err := json.Unmarshal(rawPayload, &payload); err != nil {
		return fmt.Errorf("invalid token, unable to unmarshal payload: %s", err)
	}

	log.Debugf("unmarshalled payload: %+v", payload)

	exp, ok := payload["exp"].(float64)
	if !ok {
		return fmt.Errorf("invalid token, unable to parse expiration: %v", payload["exp"])
	}

	expirationTime := time.Unix(int64(exp), 0)
	if time.Now().After(expirationTime) {
		return fmt.Errorf("token is expired (%s)", expirationTime)
	}

	nbf, ok := payload["nbf"].(float64)
	if !ok {
		return fmt.Errorf("invalid token, unable to parse notbefore: %v", payload["nbf"])
	}

	notbeforeTime := time.Unix(int64(nbf), 0)
	if time.Now().Before(notbeforeTime) {
		return fmt.Errorf("token is not valid yet (%s)", notbeforeTime)
	}

	log.Debugf("token is valid (not before: %s, not after: %s", notbeforeTime, expirationTime)

	return nil
}

// timeTrack logs the time since the passed time
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Infof("%s took %s", name, elapsed)
}
