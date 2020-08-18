package providers

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/zendesk/goship/color"
	"github.com/zendesk/goship/resources"
)

// AwsEc2Provider is an implementation of Provider interface for EC2
type AwsEc2Provider struct {
	AwsConfig      aws.Config
	AwsRegion      string
	AwsProfileName string
	AwsSession     *session.Session
}

// InitAwsEc2ProvidersFromCfg creates AwsEc2Provider list based on config sections
func InitAwsEc2ProvidersFromCfg(cfg map[interface{}]interface{}) (p []*AwsEc2Provider, err error) {
	profile := cfg["profile"].(string)

	for _, region := range cfg["regions"].([]interface{}) {
		provider := &AwsEc2Provider{
			AwsRegion:      region.(string),
			AwsProfileName: profile,
			AwsConfig: aws.Config{
				Region: aws.String(region.(string)),
			},
		}
		p = append(p, provider)
	}
	return p, err
}

// Name returns provider name
func (p *AwsEc2Provider) Name() string {
	return "aws.ec2"
}

// Init initializes provider necessary actions
func (p *AwsEc2Provider) Init() (err error) {
	p.AwsSession, err = session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           p.AwsProfileName,
		Config:            p.AwsConfig,
	})
	return err
}

// GetResources returns list of resources for particular provider
func (p *AwsEc2Provider) GetResources() (resourcesList resources.ResourceList, err error) {
	ec2Svc := ec2.New(p.AwsSession)
	resp, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		if strings.Contains(err.Error(), "UnauthorizedOperation: You are not authorized to perform this operation.") {
			color.PrintYellow(fmt.Sprintf("Skipping profile: %s, region: %s due to unsufficient privileges\n", p.AwsProfileName, p.AwsRegion))
		} else if strings.Contains(err.Error(), "is not authorized to perform: sts:AssumeRole") {
			color.PrintYellow(fmt.Sprintf("Skipping profile: %s, region: %s due to non-assumable role\n", p.AwsProfileName, p.AwsRegion))
		}
		color.PrintRed(fmt.Sprintf("Error while refreshing cache for profile: %s, region: %s):\n", p.AwsProfileName, p.AwsRegion))
		fmt.Printf("%s\n", err.Error())
	} else {
		for _, res := range resp.Reservations {
			for _, inst := range res.Instances {
				r := resources.NewEc2Instance()
				r.NativeObject = *inst
				r.ProfileName = p.AwsProfileName
				resourcesList = append(resourcesList, r)
			}
		}
	}
	return resourcesList, err
}
