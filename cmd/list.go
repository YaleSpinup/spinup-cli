package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var listSpaceCost bool
var showFailedResources bool

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listSpacesCmd)
	listCmd.AddCommand(listResourcesCmd)
	listCmd.AddCommand(listImagesCmd)
	listCmd.AddCommand(listSecretsCmd)

	listSpacesCmd.PersistentFlags().BoolVarP(&listSpaceCost, "cost", "c", false, "Query for the space cost")
	listResourcesCmd.PersistentFlags().BoolVar(&showFailedResources, "show-failed", false, "Also show failed resources")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List spinup objects",
}

var listSpacesCmd = &cobra.Command{
	Use:   "spaces",
	Short: "List spaces",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("Listing Spaces")

		spaces := spinup.Spaces{}
		if err := SpinupClient.GetResource(map[string]string{}, &spaces); err != nil {
			return err
		}

		if listSpaceCost {
			for _, s := range spaces.Spaces {
				spaceCost := &spinup.SpaceCost{}
				if err := SpinupClient.GetResource(map[string]string{"id": s.Id.String()}, spaceCost); err != nil {
					return err
				}

				s.Cost = spaceCost
			}
		}

		return formatOutput(spaces.Spaces)
	},
}

var listResourcesCmd = &cobra.Command{
	Use:   "resources [space space ...]",
	Short: "Lists the resources in your space(s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		spaces, err := parseSpaceInput(args)
		if err != nil {
			return err
		}

		if len(spaces) == 0 {
			return errors.New("at least one space is required")
		}

		output := []*spinup.Resource{}
		for _, s := range spaces {
			log.Debugf("listing resources for space %s", s)

			resources, err := SpinupClient.Resources(s)
			if err != nil {
				return err
			}

			for _, r := range resources {
				r.TypeName = r.Type.Name
				r.TypeCat = r.Type.Type
				r.TypeFlavor = r.Type.Flavor
				r.SizeID = nil
				r.IsA = ""
				r.Type = nil

				if showFailedResources || r.Status != "failed" {
					output = append(output, r)
				}
			}
		}

		j, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return err
		}

		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		f.Write(j)

		return nil
	},
}

var listImagesCmd = &cobra.Command{
	Use:   "images",
	Short: "List images in space",
	RunE: func(cmd *cobra.Command, args []string) error {
		spaces, err := parseSpaceInput(args)
		if err != nil {
			return err
		}

		if len(spaces) == 0 {
			return errors.New("space is required")
		}

		log.Debugf("listing Images for space %s", spaces)

		type ImageOutput struct {
			*spinup.Image
			OfferingName string `json:"offering_name"`
			Space        string `json:"space"`
		}

		output := []*ImageOutput{}
		for _, s := range spaces {
			images := spinup.Images{}
			if err := SpinupClient.GetResource(map[string]string{"id": s}, &images); err != nil {
				return err
			}

			for _, i := range []*spinup.Image(images) {
				oName := i.Offering.Name
				i.Offering = nil
				oOutput := ImageOutput{i, oName, s}
				output = append(output, &oOutput)
			}
		}

		j, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return err
		}

		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		f.Write(j)

		return nil
	},
}

var listSecretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "List secrets in space",
	RunE: func(cmd *cobra.Command, args []string) error {
		spaces, err := parseSpaceInput(args)
		if err != nil {
			return err
		}

		if len(spaces) == 0 {
			return errors.New("space is required")
		}

		log.Debugf("listing Secrets for space %s", spaces)

		type SecretOutput struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Space       string `json:"space"`
		}

		output := []*SecretOutput{}
		for _, s := range spaces {
			secrets, err := spaceSecrets(map[string]string{"space": s})
			if err != nil {
				return err
			}

			for _, secret := range secrets {
				output = append(output, &SecretOutput{
					Name:        secret.Name,
					Description: secret.Description,
					Space:       s,
				})
			}
		}

		j, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return err
		}

		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		f.Write(j)

		return nil
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
