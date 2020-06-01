package cmd

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var includeCost bool

func init() {
	getCmd.AddCommand(getSpaceCmd)
	getSpaceCmd.PersistentFlags().BoolVarP(&includeCost, "cost", "c", false, "Query for cost (where available)")
}

var getSpaceCmd = &cobra.Command{
	Use:   "space",
	Short: "Get details about your space(s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		spaceIds, err := parseSpaceInput(args)
		if err != nil {
			return err
		}

		log.Debugf("getting space(s) '%+v'", spaceIds)

		output := map[string]*spinup.Space{}
		for _, s := range spaceIds {
			params := map[string]string{"id": s}
			space := &spinup.GetSpace{}
			if err := SpinupClient.GetResource(params, space); err != nil {
				return err
			}

			if includeCost {
				cost := &spinup.SpaceCost{}
				if err := SpinupClient.GetResource(params, cost); err != nil {
					return err
				}
				space.Space.Cost = cost
			}

			output[space.Space.Name] = space.Space
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
