package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zendesk/goship/color"
	"github.com/zendesk/goship/config"
	"github.com/zendesk/goship/utils"
)

var scpCmd = &cobra.Command{
	Use:     "scp [<BOX>:]file [<BOX>:]file",
	Short:   "Copies files from/to remote resource in scp manner",
	Example: "goship scp kafka-production:~/file .\ngoship scp file.txt kafka-production:~/",
	Long:    `This command copies files from/to remote resources`,
	PreRunE: validateScpCmdFunc,
	Args:    cobra.ExactArgs(2),
	Run:     scpCmdFunc,
}

func validateScpCmdFunc(cmd *cobra.Command, args []string) error {
	cmd.Annotations = map[string]string{
		"copy_from": args[0],
		"copy_to":   args[1],
	}
	if checkIfRemotePath(cmd.Annotations["copy_from"]) {
		cmd.Annotations["keyword"], cmd.Annotations["remote_path"] = parseScpURL(cmd.Annotations["copy_from"])
		cmd.Annotations["local_path"] = cmd.Annotations["copy_to"]
		cmd.Annotations["direction"] = "from"
	} else if checkIfRemotePath(cmd.Annotations["copy_to"]) {
		cmd.Annotations["keyword"], cmd.Annotations["remote_path"] = parseScpURL(cmd.Annotations["copy_to"])
		cmd.Annotations["local_path"] = cmd.Annotations["copy_from"]
		cmd.Annotations["direction"] = "to"
	} else {
		return errors.New("none of paths is a remote path")
	}
	return nil
}

func init() {
	RootCmd.AddCommand(scpCmd)
}

func scpCmdFunc(cmd *cobra.Command, args []string) {

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

	command := ScpCommand{
		config.GlobalConfig.LoginUsername,
		resource.ConnectIdentifier(config.GlobalConfig.UsePrivateNetwork, config.GlobalConfig.UseDNS),
		cmd.Annotations["remote_path"],
		cmd.Annotations["local_path"],
	}
	baseCommand := []string{config.GlobalConfig.ScpBinary}
	baseCommand = append(baseCommand, config.GlobalConfig.ScpExtraParams...)

	if checkIfRemotePath(cmd.Annotations["copy_from"]) {
		baseCommand = append(baseCommand, command.CopyFromRemoteCmd()...)
	} else {
		baseCommand = append(baseCommand, command.CopyToRemoteCmd()...)
	}

	if config.GlobalConfig.UseEC2Connect {
		err := pushSSHKey(resource)
		if err != nil {
			fmt.Printf("Failed to push SSH key: %v", err)
			os.Exit(1)
		}
		baseCommand = append(baseCommand, "-i", sshPrivKeyPath(config.GlobalConfig.EC2ConnectKeyPath))
	}

	color.PrintGreen(fmt.Sprintf("Copying %s %s (%s)\n",
		cmd.Annotations["direction"], resource.Name(),
		resource.GetTag("environment")))

	if config.GlobalConfig.Verbose {
		color.PrintGreen(fmt.Sprintf("%s\n", baseCommand))
	}

	env := os.Environ()
	comm := Command{
		Binary: config.GlobalConfig.ScpBinary,
		Cmd:    baseCommand,
		Env:    env,
	}
	err = comm.Exec()
	if err != nil {
		fmt.Printf("Error while executing command %s", err)
		os.Exit(1)
	}
}
