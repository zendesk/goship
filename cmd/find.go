package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// helloCmd represents the hello command
var findCmd = &cobra.Command{
	Use:    "find <search_keyword> [environment]",
	Short:  "Finds cloud resources by keyword",
	Long:   `This command searches for particular cloud resources matching particular criteria`,
	PreRun: setSSHArgsFunc,
	Args:   cobra.MinimumNArgs(1),
	Run:    findCmdFunc,
}

func init() {
	RootCmd.AddCommand(findCmd)

}

func findCmdFunc(cmd *cobra.Command, args []string) {
	cacheList := getCacheList()

	output := filterCacheList(&cacheList, cmd.Annotations)

	for _, r := range output {
		fmt.Printf(r.RenderLongOutput())
	}

}
