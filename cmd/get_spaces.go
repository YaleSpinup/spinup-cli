package cmd

import (
	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type SpacesOutput struct {
	*spinup.Space
	ResourceCount int `json:"resource_count"`
}

func init() {
	getCmd.AddCommand(getSpacesCmd)
	getSpacesCmd.PersistentFlags().BoolVarP(&includeCost, "cost", "c", false, "Query for cost (where available)")
}

var getSpacesCmd = &cobra.Command{
	Use:   "spaces",
	Short: "Get a list of your space(s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("getting all spaces")

		spaces := spinup.Spaces{}
		if err := SpinupClient.GetResource(map[string]string{}, &spaces); err != nil {
			return err
		}

		if includeCost {
			for _, s := range spaces.Spaces {
				spaceCost := &spinup.SpaceCosts{}
				if err := SpinupClient.GetResource(map[string]string{"id": s.Id.String()}, spaceCost); err != nil {
					return err
				}

				s.Cost = spaceCost
			}
		}

		out := []SpacesOutput{}
		for _, s := range spaces.Spaces {
			o := SpacesOutput{s, len(s.Resources)}
			o.Resources = nil
			out = append(out, o)
		}

		return formatOutput(out)
	},
}
