package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zendesk/goship/color"
	"github.com/zendesk/goship/config"
	"github.com/zendesk/goship/version"
)

var (
	// ConfigFile stores path to config
	ConfigFile   string
	forceUncache bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:              "goship",
	Short:            "Find and connect to particular cloud resources",
	Long:             "This application helps find and connect to particular cloud resources",
	PersistentPreRun: initRootFlags,
}

func initRootFlags(cmd *cobra.Command, args []string) {

	err := viper.Unmarshal(&config.GlobalConfig)
	if err != nil {
		fmt.Printf("Error while unmarshalling config: %s", err.Error())
	}

	if forceUncache {
		config.GlobalConfig.CacheValidity = 0
	}

	version.CheckForNewVersion()
}

// Execute cobra
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&ConfigFile, "config", "", "config file (default is $HOME/.goship.yaml)")

	RootCmd.PersistentFlags().StringP("username", "u", "", "Username to use when logging into resources")
	RootCmd.PersistentFlags().BoolP("use-ec2-connect", "e", false, "Use EC2 Instance Connect to push your SSH key before connecting to the instance ")
	RootCmd.PersistentFlags().StringP("ec2-connect-key-path", "k", "~/.ssh/id_rsa.pub", "Path to public SSH key file used to connect via EC2 Instance Connect")
	RootCmd.PersistentFlags().StringP("ssh-command", "c", "", "command to be executed via SSH (applicable to ssh command only)")
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Be more verbose")
	RootCmd.PersistentFlags().BoolP("use-private-network", "p", false, "Use private resource identification")
	RootCmd.PersistentFlags().BoolP("use-dns", "d", false, "Use DNS instead of Resource IP")

	RootCmd.PersistentFlags().BoolVarP(&forceUncache, "uncache", "", false, "Drop any existing cache before obtaining resource list")
	RootCmd.PersistentFlags().StringP("cache-directory", "", "/tmp", "Cache directory (default is /tmp)")
	RootCmd.PersistentFlags().StringP("cache-file-prefix", "", "goship_cache_", "Cache file prefix")
	RootCmd.PersistentFlags().UintP("cache-validity", "t", 300, "Cache validity in seconds")

	_ = viper.BindPFlag("username", RootCmd.PersistentFlags().Lookup("username"))
	_ = viper.BindPFlag("use-ec2-connect", RootCmd.PersistentFlags().Lookup("use-ec2-connect"))
	_ = viper.BindPFlag("ec2-connect-key-path", RootCmd.PersistentFlags().Lookup("ec2-connect-key-path"))
	_ = viper.BindPFlag("ssh-command", RootCmd.PersistentFlags().Lookup("ssh-command"))
	_ = viper.BindPFlag("use_private_network", RootCmd.PersistentFlags().Lookup("use-private-network"))
	_ = viper.BindPFlag("use_dns", RootCmd.PersistentFlags().Lookup("use-dns"))
	_ = viper.BindPFlag("cache_directory", RootCmd.PersistentFlags().Lookup("cache-directory"))
	_ = viper.BindPFlag("cache_file_prefix", RootCmd.PersistentFlags().Lookup("cache-file-prefix"))
	_ = viper.BindPFlag("cache_validity", RootCmd.PersistentFlags().Lookup("cache-validity"))
	_ = viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(ConfigFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra-example" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".goship")
		ConfigFile = filepath.Join(home, ".goship.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if config.GlobalConfig.Verbose {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	}

	err := viper.ReadInConfig()
	if err != nil {
		color.PrintRed(fmt.Sprintf("Error while reading file %s\n", ConfigFile))
	}
}
