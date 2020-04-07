package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"
	"github.com/zendesk/goship/color"
	"github.com/zendesk/goship/config"
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
)

// default Providers configuration, for easier user onboarding
var defaultProviderConfig = `
providers:
    aws.ec2:
      - profile: default
        regions: 
        - us-east-1
`

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Runs configuration process for basic options",
	Long:  `Runs configuration process for basic options`,
	Run:   configureCmdFunc,
}

func init() {
	RootCmd.AddCommand(configureCmd)
}

func validateBool(s string) error {
	_, err := strconv.ParseBool(s)
	if err != nil {
		return fmt.Errorf("value is other than true/false")
	}
	return err
}

func validateUint(s string) error {
	_, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return fmt.Errorf("value is not valid non-negative integer")
	}
	return err
}

func validatePath(s string) (err error) {
	if _, err = os.Stat(s); os.IsNotExist(err) {
		return err
	}
	return err
}

func configureCmdFunc(cmd *cobra.Command, args []string) {

	if _, err := os.Stat(ConfigFile); !os.IsNotExist(err) {
		color.PrintYellow(fmt.Sprintf("Config file %s already exists. Refusing to reconfigure it. Please edit config file manually.\n", ConfigFile))
		os.Exit(1)
	}

	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	useEC2Connect, _ := ui.Ask("Use EC2 Instance Connect to connect to instances?", &input.Options{
		Required:     true,
		Loop:         true,
		Default:      RootCmd.PersistentFlags().Lookup("use-ec2-connect").DefValue,
		ValidateFunc: validateBool,
	})

	config.GlobalConfig.UseEC2Connect, _ = strconv.ParseBool(useEC2Connect)

	if config.GlobalConfig.UseEC2Connect {
		config.GlobalConfig.CacheDirectory, _ = ui.Ask("Path to SSH key file used to connect via EC2 Instance Connect", &input.Options{
			Required:     true,
			Default:      RootCmd.PersistentFlags().Lookup("ec2-connect-key-path").DefValue,
			Loop:         true,
			ValidateFunc: validatePath,
		})
	}

	config.GlobalConfig.LoginUsername, _ = ui.Ask("What username should be used when connecting to remote resources?", &input.Options{
		Required: true,
		Loop:     true,
	})

	usePrivateNetwork, _ := ui.Ask("Use private network when connecting to remote resources? [true/false]", &input.Options{
		Required:     false,
		Default:      RootCmd.PersistentFlags().Lookup("use-private-network").DefValue,
		Loop:         true,
		ValidateFunc: validateBool,
	})
	config.GlobalConfig.UsePrivateNetwork, _ = strconv.ParseBool(usePrivateNetwork)

	useDNS, _ := ui.Ask("Use DNS when connecting to remote resources? [true/false]", &input.Options{
		Required:     false,
		Default:      RootCmd.PersistentFlags().Lookup("use-dns").DefValue,
		Loop:         true,
		ValidateFunc: validateBool,
	})
	config.GlobalConfig.UseDNS, _ = strconv.ParseBool(useDNS)

	config.GlobalConfig.CacheDirectory, _ = ui.Ask("Cache directory where goship cache will be stored", &input.Options{
		Required:     false,
		Default:      RootCmd.PersistentFlags().Lookup("cache-directory").DefValue,
		Loop:         true,
		ValidateFunc: validatePath,
	})

	config.GlobalConfig.CacheFilePrefix, _ = ui.Ask("Cache files prefix", &input.Options{
		Required:     false,
		Default:      RootCmd.PersistentFlags().Lookup("cache-file-prefix").DefValue,
		Loop:         true,
		ValidateFunc: validatePath,
	})

	cacheValidity, _ := ui.Ask("Cache validity in seconds", &input.Options{
		Required:     false,
		Default:      RootCmd.PersistentFlags().Lookup("cache-validity").DefValue,
		Loop:         true,
		ValidateFunc: validateUint,
	})
	config.GlobalConfig.CacheValidity, _ = strconv.ParseUint(cacheValidity, 10, 0)

	verbose, _ := ui.Ask("Be verbose by default? [true/false]", &input.Options{
		Required:     false,
		Default:      RootCmd.PersistentFlags().Lookup("verbose").DefValue,
		Loop:         true,
		ValidateFunc: validateBool,
	})
	config.GlobalConfig.Verbose, _ = strconv.ParseBool(verbose)

	config.GlobalConfig.SSHBinary, _ = ui.Ask("Default ssh binary", &input.Options{
		Required:     false,
		Default:      config.GlobalConfig.SSHBinary,
		Loop:         true,
		ValidateFunc: validatePath,
	})

	config.GlobalConfig.ScpBinary, _ = ui.Ask("Default scp binary", &input.Options{
		Required:     false,
		Default:      config.GlobalConfig.ScpBinary,
		Loop:         true,
		ValidateFunc: validatePath,
	})

	c, _ := yaml.Marshal(config.GlobalConfig)
	c = append(c, defaultProviderConfig...)
	fmt.Printf("Config file: %s", viper.ConfigFileUsed())
	cacheFile, err := os.Create(ConfigFile)
	if err != nil {
		color.PrintRed(fmt.Sprintf("Error while creating config file %s: %s\n", ConfigFile, err.Error()))
		os.Exit(1)
	}
	defer cacheFile.Close()

	_, err = cacheFile.Write(c)
	if err != nil {
		color.PrintRed(fmt.Sprintf("Error while writing config file %s: %s\n", ConfigFile, err.Error()))
		os.Exit(1)
	}

	color.PrintGreen(fmt.Sprintf("Config file saved to `%s`. Please refer to documentation in order to customize cloud providers\n", ConfigFile))
}
