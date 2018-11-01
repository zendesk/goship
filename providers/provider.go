package providers

import (
	"github.com/zendesk/goship/resources"
)

// Provider defines interface which should be implemented by every provider
type Provider interface {
	Init() error
	Name() string
	GetResources() (resources.ResourceList, error)
}

// InitProvidersFromConfig returns list of initialized providers for particular types
func InitProvidersFromConfig(providersConfig map[string]interface{}) (p []Provider, err error) {

	for cloudVendor, cloudProvidersConfig := range providersConfig {
		switch cloudVendor {
		case "aws":
			for providerName, providerCfg := range cloudProvidersConfig.(map[string]interface{}) {
				switch providerName {
				case "ec2":
					for _, providerItem := range providerCfg.([]interface{}) {
						providersList, _ := InitAwsEc2ProvidersFromCfg(providerItem.(map[interface{}]interface{}))
						for _, provider := range providersList {
							p = append(p, provider)
						}

					}
				}
			}
		}
	}
	return p, err
}
