package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of got2s",
	Long:  `All software has versions. This is got2s's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("got2s v0.1")
	},
}

