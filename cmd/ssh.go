package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zendesk/goship/color"
	"github.com/zendesk/goship/config"
	"github.com/zendesk/goship/utils"
)

var sshCmd = &cobra.Command{
	Use:    "ssh <search_keyword> [environment] [--ssh-command \"<command to execute via ssh>\"]",
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

	if config.GlobalConfig.UseEC2Connect {
		err := pushSSHKey(resource)
		if err != nil {
			fmt.Printf("Failed to push SSH key: %v", err)
			os.Exit(1)
		}
	}

	baseCommand := []string{config.GlobalConfig.SSHBinary}
	baseCommand = append(baseCommand, config.GlobalConfig.SSHExtraParams...)

	sshCommandWithArgs := append(baseCommand, []string{
		fmt.Sprintf(
			"%s@%s",
			config.GlobalConfig.LoginUsername,
			resource.ConnectIdentifier(config.GlobalConfig.UsePrivateNetwork, config.GlobalConfig.UseDNS),
		),
	}...)

	sshCommand, err := RootCmd.PersistentFlags().GetString("ssh-command")
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	if sshCommand != "" {
		sshCommandWithArgs = append(sshCommandWithArgs, strings.Split(sshCommand, " ")...)
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
	err = comm.Exec()
	if err != nil {
		fmt.Printf("Error while executing command %s", err)
		os.Exit(1)
	}
}
