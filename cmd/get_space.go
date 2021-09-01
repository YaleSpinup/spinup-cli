package cmd

import (
	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type SpaceOutput struct {
	*spinup.Space
	Resources []*spinup.Resource `json:",omitempty"`
}

var (
	includeCost         bool
	showFailedResources bool
	showResources       bool
)

func init() {
	getCmd.AddCommand(getSpaceCmd)
	getSpaceCmd.PersistentFlags().BoolVarP(&includeCost, "cost", "c", false, "Query for cost (where available)")
	getSpaceCmd.PersistentFlags().BoolVar(&showResources, "resources", false, "Show resources")
	getSpaceCmd.PersistentFlags().BoolVar(&showFailedResources, "failed", false, "Also show failed resources")
}

var getSpaceCmd = &cobra.Command{
	Use:   "space",
	Short: "Get details about your space(s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		spaces, err := parseSpaceInput(args)
		if err != nil {
			return err
		}

		log.Debugf("getting space(s) '%+v'", spaces)

		output := map[string]SpaceOutput{}
		for _, s := range spaces {
			params := map[string]string{"id": s}
			space := &spinup.GetSpace{}
			if err := SpinupClient.GetResource(params, space); err != nil {
				return err
			}

			if includeCost {
				cost := &spinup.SpaceCosts{}
				if err := SpinupClient.GetResource(params, cost); err != nil {
					return err
				}
				space.Space.Cost = cost
			}

			var resourcesOut []*spinup.Resource
			if showResources {
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
						resourcesOut = append(resourcesOut, r)
					}
				}
			}

			out := SpaceOutput{space.Space, resourcesOut}
			output[space.Space.Name] = out
		}

		return formatOutput(output)
	},
}
