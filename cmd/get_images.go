package cmd

import (
	"github.com/YaleSpinup/spinup-cli/pkg/spinup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	getCmd.AddCommand(getImagesCmd)
}

var getImagesCmd = &cobra.Command{
	Use:     "images [space]",
	Short:   "Get a list of images for a space",
	PreRunE: getSpaceLevelCmdPreRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("get images: %+v", args)

		type ImageOutput struct {
			*spinup.Image
			OfferingName string `json:"offering_name"`
			Space        string `json:"space"`
		}

		images := spinup.Images{}
		if err := SpinupClient.GetResource(getParams, &images); err != nil {
			return err
		}

		out := []*ImageOutput{}
		for _, i := range []*spinup.Image(images) {
			oName := i.Offering.Name
			i.Offering = nil
			oOutput := ImageOutput{i, oName, getParams["space"]}
			out = append(out, &oOutput)
		}

		return formatOutput(out)
	},
}
