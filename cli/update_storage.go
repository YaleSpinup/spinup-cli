package cli

import (
	"errors"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
)

func updateStorage(params map[string]string, resource *spinup.Resource) error {
	return errors.New("storage updates are not currently supported")
}
