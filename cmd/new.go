package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.AddCommand(newSpaceCmd)
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new resources",
}

var newSpaceCmd = &cobra.Command{
	Use:   "space",
	Short: "Command to create a space",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating spaces is not currently supported from the CLI, please use the web interface.")
	},
}
