package cmd

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information about the spinup-cli.",
	Run: func(cmd *cobra.Command, args []string) {

		output := struct {
			Version           string `json:"verison,omitempty"`
			VersionPrerelease string `json:"verisonPrerelease,omitempty"`
			BuildStamp        string `json:"buildStamp,omitempty"`
			GitHash           string `json:"gitHash,omitempty"`
		}{
			Version:           Version,
			VersionPrerelease: VersionPrerelease,
			BuildStamp:        BuildStamp,
			GitHash:           GitHash,
		}

		j, _ := json.MarshalIndent(output, "", "  ")
		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		f.Write(j)
	},
}
