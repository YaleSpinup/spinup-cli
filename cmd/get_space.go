package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"github.com/YaleSpinup/spinup/pkg/spinup"
	"github.com/spf13/cobra"
)

var includeCost bool

func init() {
	getCmd.AddCommand(getSpaceCmd)
	getCmd.PersistentFlags().BoolVarP(&includeCost, "cost", "c", false, "Query for cost (where available)")
}

var getSpaceCmd = &cobra.Command{
	Use:   "space",
	Short: "Commands to get details about a space",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("exactly 1 space id is required")
		}

		space := &spinup.Space{}
		err := SpinupClient.GetResource(args[0], space)
		if err != nil {
			return err
		}

		if includeCost {
			cost := &spinup.SpaceCost{}
			err := SpinupClient.GetResource(args[0], cost)
			if err != nil {
				return err
			}

			space.Cost = cost
		}

		j, err := json.MarshalIndent(space, "", "  ")
		if err != nil {
			return err
		}

		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		f.Write(j)

		return nil
	},
}
