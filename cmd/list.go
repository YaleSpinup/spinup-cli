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

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listSpacesCmd)
	listCmd.AddCommand(listResourcesCmd)
	listCmd.AddCommand(listImagesCmd)

	listCmd.PersistentFlags().BoolVarP(&listSpaceCost, "cost", "c", false, "Query for the space cost")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources",
}

var listSpacesCmd = &cobra.Command{
	Use:   "spaces",
	Short: "List spaces",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("Listing Spaces")

		spaces := spinup.Spaces{}
		err := SpinupClient.GetResource("", &spaces)
		if err != nil {
			return err
		}

		if listSpaceCost {
			for _, s := range spaces.Spaces {
				cost := &spinup.SpaceCost{}
				err := SpinupClient.GetResource(s.Id.String(), cost)
				if err != nil {
					return err
				}

				s.Cost = cost
			}
		}

		j, err := json.MarshalIndent(spaces.Spaces, "", "  ")
		if err != nil {
			return err
		}

		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		f.Write(j)

		return nil
	},
}

var listResourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "Lists the resources in space",
	RunE: func(cmd *cobra.Command, args []string) error {
		if spinupSpaceID == "" {
			return errors.New("space id is required")
		}

		log.Debugf("listing resources for space %s", spinupSpaceID)

		out, err := SpinupClient.Resources(spinupSpaceID)
		if err != nil {
			return err
		}

		j, err := json.MarshalIndent(out, "", "  ")
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
		if spinupSpaceID == "" {
			return errors.New("space id is required")
		}

		log.Debugf("listing Images for space %s", spinupSpaceID)

		j, err := listImages(spinupSpaceID)
		if err != nil {
			return err
		}

		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		f.Write(j)

		return nil
	},
}

func listImages(spaceId string) ([]byte, error) {
	images := spinup.Images{}
	err := SpinupClient.GetResource(spinupSpaceID, &images)
	if err != nil {
		return nil, err
	}

	type ImageOutput struct {
		*spinup.Image
		OfferingID string `json:"offering_id"`
	}

	output := []*ImageOutput{}
	for _, i := range []*spinup.Image(images) {
		oID := i.Offering.ID.String()
		i.Offering = nil
		output = append(output, &ImageOutput{i, oID})
	}

	j, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, err
	}

	return j, nil
}
