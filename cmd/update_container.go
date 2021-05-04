package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
)

var redeployContainerCmd bool

func init() {
	updateCmd.PersistentFlags().BoolVarP(&redeployContainerCmd, "redeploy", "r", false, "Redeploy with the current configuraiton.")
}

func updateContainer(params map[string]string, resource *spinup.Resource) error {
	var j []byte
	var err error
	switch {
	case resource.Status != "created":
		return fmt.Errorf("container must be in 'created' state, current state is %s", resource.Status)
	case redeployContainerCmd:
		if j, err = redeployContainer(params, resource); err != nil {
			return err
		}
	default:
		return errors.New("only redeployment is currently supported")
	}

	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	f.Write(j)

	return nil
}

func redeployContainer(params map[string]string, resource *spinup.Resource) ([]byte, error) {
	input, err := json.Marshal(map[string]bool{"only_redeploy": true})
	if err != nil {
		return []byte{}, err
	}

	log.Debugf("putting input: %s", string(input))

	info := &spinup.ContainerService{}
	if err = SpinupClient.PutResource(params, input, info); err != nil {
		return []byte{}, err
	}

	return []byte("OK\n"), nil
}
