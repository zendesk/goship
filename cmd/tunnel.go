package cmd

import (
	"fmt"
	"os"

	"github.com/zendesk/goship/color"
	"github.com/zendesk/goship/config"
	"github.com/zendesk/goship/utils"
	"github.com/spf13/cobra"
)

// helloCmd represents the hello command
var tunnelCmd = &cobra.Command{
	Use:     "tunnel <Resource> <local_bind_address> <remote_bind_address>",
	Short:   "Creates SSH tunnel to remote resource",
	Example: "goship tunnel storm-master-1122334455 8080 127.0.0.1:8080",
	Long:    `This command creates SSH tunnel between localhost and remote resource`,
	PreRunE: validateTunnelCmdFunc,
	Args:    cobra.ExactArgs(3),
	Run:     tunnelCmdFunc,
}

func init() {
	RootCmd.AddCommand(tunnelCmd)
}

func validateTunnelCmdFunc(cmd *cobra.Command, args []string) error {
	cmd.Annotations = map[string]string{
		"keyword":             args[0],
		"local_bind_address":  args[1],
		"remote_bind_address": args[2],
	}
	if err := validateBindAddress(cmd.Annotations["local_bind_address"]); err != nil {
		return err
	}
	if err := validateBindAddress(cmd.Annotations["remote_bind_address"]); err != nil {
		return err
	}
	return nil
}

func tunnelCmdFunc(cmd *cobra.Command, args []string) {

	cacheList := getCacheList()
	output := filterCacheList(&cacheList, cmd.Annotations)

	if len(output) == 0 {
		fmt.Printf("Could not find any matches for identifier %s", cmd.Annotations["keyword"])
		os.Exit(1)
	}

	resource, err := utils.ChooseFromList(output)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	tunnelCommand := []string{
		config.GlobalConfig.SSHBinary,
		fmt.Sprintf(
			"%s@%s",
			config.GlobalConfig.LoginUsername,
			resource.ConnectIdentifier(config.GlobalConfig.UsePrivateNetwork, config.GlobalConfig.UseDNS),
		),
		"-N",
		"-L",
		fmt.Sprintf(
			"%s:%s",
			cmd.Annotations["local_bind_address"],
			cmd.Annotations["remote_bind_address"],
		),
	}
	color.PrintGreen(fmt.Sprintf("Creating ssh tunnel to %s (%s)\n",
		resource.Name(), resource.GetTag("environment")))

	if config.GlobalConfig.Verbose {
		color.PrintGreen(fmt.Sprintf("%s\n", tunnelCommand))
	}

	fmt.Printf("Tunnel created on %s\n", formatProperAddressWithPort(cmd.Annotations["local_bind_address"], "localhost"))
	env := os.Environ()
	comm := Command{
		Binary: config.GlobalConfig.SSHBinary,
		Cmd:    tunnelCommand,
		Env:    env,
	}
	comm.Exec()

}
