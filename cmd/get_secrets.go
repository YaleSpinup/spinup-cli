package cmd

import (
	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	getCmd.AddCommand(getSecretsCmd)
}

var getSecretsCmd = &cobra.Command{
	Use:     "secrets [space]",
	Short:   "Get a list of secrets for a space",
	PreRunE: getSpaceLevelCmdPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("get secrets: %+v", args)

		type SecretOutput struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Space       string `json:"space"`
		}

		out := []*SecretOutput{}
		secrets, err := spaceSecrets(getParams)
		if err != nil {
			return err
		}

		for _, secret := range secrets {
			out = append(out, &SecretOutput{
				Name:        secret.Name,
				Description: secret.Description,
				Space:       getParams["space"],
			})
		}
		return formatOutput(out)
	},
}

func spaceSecrets(params map[string]string) ([]*spinup.Secret, error) {
	// collect a list of secrets from the space
	secrets := &spinup.Secrets{}
	if err := SpinupClient.GetResource(params, secrets); err != nil {
		return nil, err
	}

	log.Debugf("got list of secrets in space %+v", secrets)

	// get details about each secret (necessary to map the ARN to the name)
	spaceSecrets := []*spinup.Secret{}
	for _, s := range *secrets {
		secret := &spinup.Secret{}
		if err := SpinupClient.GetResource(
			map[string]string{
				"space":      params["space"],
				"secretname": string(s),
			}, secret); err != nil {
			return nil, err
		}
		spaceSecrets = append(spaceSecrets, secret)
	}

	return spaceSecrets, nil
}
