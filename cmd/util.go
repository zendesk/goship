package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/zendesk/goship/cache"
	"github.com/zendesk/goship/color"
	"github.com/zendesk/goship/config"
	"github.com/zendesk/goship/providers"
	"github.com/zendesk/goship/resources"
)

// Command defines command to run
type Command struct {
	Binary string
	Cmd    []string
	Env    []string
}

// Exec executes defined command
func (c *Command) Exec() error {
	return syscall.Exec(
		c.Binary,
		c.Cmd,
		c.Env,
	)
}

// ScpCommand defines scp command
type ScpCommand struct {
	User      string
	Host      string
	HostPath  string
	LocalPath string
}

// CopyFromRemoteCmd defines command for copying from remote host
func (s *ScpCommand) CopyFromRemoteCmd() []string {
	return []string{fmt.Sprintf("%s@%s:%s", s.User, s.Host, s.HostPath), s.LocalPath}
}

// CopyToRemoteCmd defines command for copying to remote host
func (s *ScpCommand) CopyToRemoteCmd() []string {
	return []string{s.LocalPath, fmt.Sprintf("%s@%s:%s", s.User, s.Host, s.HostPath)}
}

func initCaches() (cacheList cache.GlobalCacheList) {
	providersList, err := providers.InitProvidersFromConfig(config.GlobalConfig.Providers)
	if err != nil {
		color.PrintRed(fmt.Sprintf("Error while initializing provider: %s", err.Error()))
	}

	for _, p := range providersList {
		switch p.Name() {
		case "aws.ec2":
			c := cache.NewAwsEc2Cache(*p.(*providers.AwsEc2Provider))
			cacheList = append(cacheList, c)
		}
	}

	return cacheList
}

func getCacheList() cache.GlobalCacheList {

	cacheList := initCaches()
	if len(cacheList) == 0 {
		color.PrintYellow("WARNING: No valid providers configured. Please refer to the documentation in order to configure it.\n")
	}
	refreshStart := time.Now()
	err := cacheList.RefreshInParallel(false)
	if err != nil {
		color.PrintRed(fmt.Sprintf("Error while getting cache list: %s", err.Error()))
	}
	refreshElapsed := time.Since(refreshStart)

	if config.GlobalConfig.Verbose {
		fmt.Printf("Cache operations total time: %s\n\n", refreshElapsed)
	}

	return cacheList
}

func filterCacheList(cacheList *cache.GlobalCacheList, criteria map[string]string) (output resources.ResourceList) {
	for _, c := range *cacheList {
		for _, r := range c.Resources() {
			if env, exists := criteria["environment"]; exists {
				if !strings.Contains(r.GetTag("environment"), env) {
					continue
				}
			}
			if strings.Contains(r.String(), criteria["keyword"]) {
				exists := false
				for _, o := range output {
					if r.ResourceID() == o.ResourceID() {
						exists = true
					}
				}
				if !exists {
					output = append(output, r)
				}
			}
		}
	}
	sort.Sort(output)
	return output
}

func validatePortNumber(portNumber string) error {
	_, err := strconv.Atoi(portNumber)
	if err != nil {
		return fmt.Errorf("unknown port value: %s", portNumber)
	}
	return nil
}

// We intentionally does not check IP address (if exists) in order to allow users
// to pass canonical hostanames like localhost:8080
func validateBindAddress(address string) error {
	l := strings.Split(address, ":")
	if len(l) == 1 {
		return validatePortNumber(l[0])
	} else if len(l) == 2 {
		return validatePortNumber(l[1])
	} else {
		return errors.New("too much sections for parameter")
	}
}

func formatProperAddressWithPort(s string, defaultHost string) string {
	tuple := strings.Split(s, ":")
	if len(tuple) == 1 {
		return fmt.Sprintf("%s:%s", defaultHost, s)
	}
	return s
}

func checkIfRemotePath(p string) bool {
	return strings.Contains(p, ":")
}

func parseScpURL(p string) (string, string) {
	return strings.Split(p, ":")[0], strings.Split(p, ":")[1]
}

func pushSSHKey(r resources.Resource) error {
	ec2, ok := r.(*resources.Ec2Instance)
	if !ok {
		return fmt.Errorf("pushing SSH key is supported only for EC2 instances")
	}
	color.PrintGreen(fmt.Sprintf("Sending SSH key to AWS SSM for %s (%s) in %s\n",
		r.Name(), r.GetTag("environment"), r.GetZone()))
	return ec2.PushSSHKey(config.GlobalConfig.EC2ConnectKeyPath)
}
