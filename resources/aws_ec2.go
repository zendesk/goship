package resources

import (
	"bytes"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Ec2Instance represents resource for ec2 instances
type Ec2Instance struct {
	NativeObject ec2.Instance

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
func (i *Ec2Instance) ConnectIdentifier(usePrivateID bool, useDNS bool) string {
	if usePrivateID {
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
