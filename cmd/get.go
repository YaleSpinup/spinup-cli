package cmd

import (
	"github.com/spf13/cobra"
)

var detailedGetCmd bool

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().BoolVarP(&detailedGetCmd, "details", "d", false, "Get detailed output about the resource")
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details about a resource",
}
