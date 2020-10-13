package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

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
		baseCommand = append(baseCommand, "-i", privKey)

		baseCommand = append(baseCommand, proxyCmd...)

	}
	baseCommand = append(baseCommand, config.GlobalConfig.ScpExtraParams...)

	if checkIfRemotePath(cmd.Annotations["copy_from"]) {
		baseCommand = append(baseCommand, command.CopyFromRemoteCmd()...)
	} else {
		baseCommand = append(baseCommand, command.CopyToRemoteCmd()...)
	}
	color.PrintGreen(fmt.Sprintf("Copying %s %s (%s)\n",
		cmd.Annotations["direction"], resource.Name(),
		resource.GetTag("environment")))

	env := os.Environ()
	comm := Command{
		Binary: config.GlobalConfig.ScpBinary,
		Cmd:    baseCommand,
		Env:    env,
	}
	if config.GlobalConfig.Verbose {
		color.PrintGreen(fmt.Sprintf("%s\n", baseCommand))
	}
	err = comm.Exec()
	if err != nil {
		fmt.Printf("Error while executing command %s", err)
		os.Exit(1)
	}
}
