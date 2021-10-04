package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type CmdVersion struct {
	AppVersion string
	BuildTime  string
	GitCommit  string
	GitRef     string
}

var (
	Version *CmdVersion
	long    bool

	versionCmd = &cobra.Command{
		Use:     "version",
		Aliases: []string{"vers"},
		Short:   "Display version information",
		RunE: func(_ *cobra.Command, args []string) error {
			if long {
				fmt.Printf("Toker version: %s\nBuildtime: %s\nGitCommit: %s\n", Version.AppVersion, Version.BuildTime, Version.GitCommit)
				return nil
			}

			fmt.Printf("%s\n", Version.AppVersion)
			return nil
		},
	}
)

func init() {
	versionCmd.Flags().BoolVarP(&long, "long", "l", false, "get more verbose version information")
	rootCmd.AddCommand(versionCmd)
}
