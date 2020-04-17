package resources

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/zendesk/goship/config"
	"io/ioutil"
)

// Ec2Instance represents resource for ec2 instances
type Ec2Instance struct {
	NativeObject ec2.Instance

	// Store AWS Profile name so we can recreate AWS Session when EC2 Instance Connect is enabled
	ProfileName string

	ShortOutputTemplate string
	LongOutputTemplate  string
}

// NewEc2Instance creates Ec2Instnace var with default render templates
func NewEc2Instance() *Ec2Instance {
	return &Ec2Instance{
		ShortOutputTemplate: DefaultShortOutputTemplate,
		LongOutputTemplate:  DefaultLongOutputTemplate,
	}
}

// Name returns resource name
func (i *Ec2Instance) Name() string {
	return i.GetTag("Name")
}

// ConnectIdentifier returns the identifier (eg. IP or DNS name) which will be used when connecting to instances
func (i *Ec2Instance) ConnectIdentifier(usePrivateID, useDNS bool) string {
	if usePrivateID || i.NativeObject.PublicDnsName == nil || i.NativeObject.PublicIpAddress == nil {
		if useDNS {
			return *i.NativeObject.PrivateDnsName
		}
		return *i.NativeObject.PrivateIpAddress
	}
	// else we're using public ID
	if useDNS {
		return *i.NativeObject.PublicDnsName
	}
	return *i.NativeObject.PublicIpAddress
}

// ResourceID returns unique resource ID
func (i *Ec2Instance) ResourceID() string {
	return *i.NativeObject.InstanceId
}

// GetTag returns tag value or empty string if not found
func (i *Ec2Instance) GetTag(tagName string) string {
	for _, tag := range i.NativeObject.Tags {
		if *tag.Key == tagName {
			return *tag.Value
		}
	}
	return ""
}

// GetZone returns Available Zone
func (i *Ec2Instance) GetZone() string {
	return *i.NativeObject.Placement.AvailabilityZone
}

// GetRegion returns Region
func (i *Ec2Instance) GetRegion() string {
	zone := i.GetZone()
	region := zone[:len(zone)-1]
	return region
}

// RenderShortOutput renders the list output for resource
func (i *Ec2Instance) RenderShortOutput() string {
	var tpl bytes.Buffer

	t := ParseTemplate(i.ShortOutputTemplate)
	//	fmt.Print("%v",i.NativeObject)
	err := t.Execute(&tpl, i.NativeObject)
	if err != nil {
		panic("ERROR WHILE PARSING TEMPLATE")
	}

	return tpl.String()
}

// RenderLongOutput renders the detailed output for resource
func (i *Ec2Instance) RenderLongOutput() string {
	var tpl bytes.Buffer
	t := ParseTemplate(i.LongOutputTemplate)
	err := t.Execute(&tpl, i.NativeObject)
	if err != nil {
		panic("ERROR WHILE PARSING TEMPLATE")
	}

	return tpl.String()
}

// String returns resource attributes as string
func (i *Ec2Instance) String() string {
	return i.NativeObject.String()
}

// SortKey returns sort key
func (i *Ec2Instance) SortKey() string {
	return i.GetTag("Name")
}

// PushSSHKey sends SSH key to EC2 Instance Connect
func (i *Ec2Instance) PushSSHKey(KeyPath string) error {
	key, err := ioutil.ReadFile(config.GlobalConfig.EC2ConnectKeyPath)

	if err != nil {
		return fmt.Errorf("failed to read SSH key from %s: %w", config.GlobalConfig.EC2ConnectKeyPath, err)
	}

	s, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           i.ProfileName,
		Config:            aws.Config{Region: aws.String(i.GetRegion())},
	})

	if err != nil {
		return fmt.Errorf("failed to create AWS Session: %w", err)
	}

	svc := ec2instanceconnect.New(s)
	input := &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: aws.String(i.GetZone()),
		InstanceId:       aws.String(i.ResourceID()),
		InstanceOSUser:   aws.String(config.GlobalConfig.LoginUsername),
		SSHPublicKey:     aws.String(string(key)),
	}
	_, err = svc.SendSSHPublicKey(input)

	return err
}
