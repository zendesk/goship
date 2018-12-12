package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zendesk/goship/version"
)

// helloCmd represents the hello command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints binary version.",
	Long:  `Prints binary version`,
	Run:   versionCmdFunc,
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func versionCmdFunc(cmd *cobra.Command, args []string) {
	fmt.Printf("Version: %s\n", version.VersionNumber)
}
