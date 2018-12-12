package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zendesk/goship/color"
	"github.com/zendesk/goship/config"
	"github.com/zendesk/goship/utils"
)

var sshCmd = &cobra.Command{
	Use:    "ssh <search_keyword> [environment]",
	Short:  "Connects to resources via SSH",
	Long:   `Connects to resources via SSH`,
	PreRun: setSSHArgsFunc,
	Args:   cobra.MinimumNArgs(1),
	Run:    sshCmdFunc,
}

func init() {
	RootCmd.AddCommand(sshCmd)
}

func setSSHArgsFunc(cmd *cobra.Command, args []string) {
	cmd.Annotations = map[string]string{
		"keyword": args[0],
	}
	if len(args) > 1 {
		cmd.Annotations["environment"] = args[1]
	}
}

func sshCmdFunc(cmd *cobra.Command, args []string) {

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

	sshCommandWithArgs := []string{
		config.GlobalConfig.SSHBinary,
		fmt.Sprintf(
			"%s@%s",
			config.GlobalConfig.LoginUsername,
			resource.ConnectIdentifier(config.GlobalConfig.UsePrivateNetwork, config.GlobalConfig.UseDNS),
		),
	}
	color.PrintGreen(fmt.Sprintf("Logging into %s (%s)\n",
		resource.Name(), resource.GetTag("environment")))

	if config.GlobalConfig.Verbose {
		color.PrintGreen(fmt.Sprintf("%s\n", sshCommandWithArgs))
	}

	env := os.Environ()
	comm := Command{
		Binary: config.GlobalConfig.SSHBinary,
		Cmd:    sshCommandWithArgs,
		Env:    env,
	}
	comm.Exec()
}
