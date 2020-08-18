package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zendesk/goship/config"
	"github.com/zendesk/goship/utils"
	"os"
)

var pushkeyCmd = &cobra.Command{
	Use:    "push-key <search_keyword> [environment]",
	Short:  "Push key to instance via EC2 Instance Connect",
	Long:   "Push key to instance via EC2 Instance Connect",
	PreRun: setSSHArgsFunc,
	Args:   cobra.MinimumNArgs(1),
	Run:    pushkeyCmdFunc,
}

func init() {
	RootCmd.AddCommand(pushkeyCmd)
}

func pushkeyCmdFunc(cmd *cobra.Command, args []string) {

	cacheList := getCacheList()

	output := filterCacheList(&cacheList, cmd.Annotations)

	if len(output) == 0 {
		fmt.Printf("Could not find any matches for identifier %s", cmd.Annotations["resource"])
		os.Exit(1)
	}

	resource, err := utils.ChooseFromList(output)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	if config.GlobalConfig.UseEC2Connect {
		err := pushSSHKey(resource)
		if err != nil {
			fmt.Printf("Failed to push SSH key: %v", err)
			os.Exit(1)
		}
	}
}
