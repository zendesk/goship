package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zendesk/goship/color"
	"github.com/zendesk/goship/config"
	"github.com/zendesk/goship/utils"
	"io/ioutil"
	"strconv"

	"os"
	"path"
	"strings"
	"time"
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

	if config.GlobalConfig.UseSSM {
		fmt.Println("Using SSM")
		//Keep tmp pem key in cache directory
		currentTime := time.Now()
		tempKeyName := strconv.FormatInt(currentTime.UnixNano(), 10)
		pubKey := path.Join(config.GlobalConfig.CacheDirectory, tempKeyName+".pub")
		privKey := path.Join(config.GlobalConfig.CacheDirectory, tempKeyName)
		//here we need to generate SSH keys
		privKeyBlob, errKey := rsa.GenerateKey(rand.Reader, 4092)
		if errKey != nil {
			fmt.Printf("Error while generating key %s", errKey)
			os.Exit(1)
		}
		pubKeyBlob := privKeyBlob.PublicKey
		if errPubPEM := utils.SavePublicPEMKey(pubKey, &pubKeyBlob); errPubPEM != nil {
			fmt.Printf("Error while generating public key %s", errPubPEM)
			os.Exit(1)
		}
		defer func() {
			if err = utils.DeleteTempKey(pubKey); err != nil {
				fmt.Printf("Error while removing old key %s", err)
				os.Exit(1)
			}
		}()

		if errPrivPEM := utils.SavePrivPEMKey(privKey, privKeyBlob); errPrivPEM != nil {
			fmt.Printf("Error while saving private key %s", errPrivPEM)
			os.Exit(1)
		}
		defer func() {
			if err = utils.DeleteTempKey(privKey); err != nil {
				fmt.Printf("Error while removing old private key %s", err)
				os.Exit(1)
			}
		}()
		pubKeyData := []byte{}
		if pubKey != "" {
			pubKeyData, err = ioutil.ReadFile(pubKey)
			if err != nil {
				fmt.Printf("Error while removing old private key %s", err)
				os.Exit(1)
			}
		}
		proxyCmd, err := startSSH(resource, pubKeyData)
		if err != nil {
			fmt.Printf("Error while starting SSH session %s", err)
			os.Exit(1)
		}
		sshCommandWithArgs = append(sshCommandWithArgs, proxyCmd...)

		sshCommandWithArgs = append(sshCommandWithArgs, "-i", privKey)
		sshCommandWithArgs = append(sshCommandWithArgs, "-o", "IdentitiesOnly=yes")
	}
	env := os.Environ()

	comm := Command{
		Binary: config.GlobalConfig.SSHBinary,
		Cmd:    sshCommandWithArgs,
		Env:    env,
	}
	if config.GlobalConfig.Verbose {
		sshCommandWithArgs = append(sshCommandWithArgs, "-vvvv")
	}
	if config.GlobalConfig.Verbose {
		color.PrintGreen(fmt.Sprintf("%s\n", sshCommandWithArgs))
	}
	err = comm.Exec()
	if err != nil {
		fmt.Printf("Error while executing command %s", err)
		os.Exit(1)
	}
}
