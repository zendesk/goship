package resources

import (
	"bytes"
	"encoding/json"

	//"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"os"
	"os/exec"
	"strings"

	"strconv"

	"github.com/alessio/shellescape"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect/ec2instanceconnectiface"
)

// Ec2Instance represents resource for ec2 instances
type Ec2Instance struct {
	NativeObject ec2.Instance
	ProfileName  string

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

// ResourceID returns unique resource ID
func (i *Ec2Instance) PlacementAZ() string {
	return *i.NativeObject.Placement.AvailabilityZone
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

// GetZone returns Available Zone
func (i *Ec2Instance) GetZone() string {
	return *i.NativeObject.Placement.AvailabilityZone
}

// StartSSH will start a ssh proxy session for a chosed node
func (i *Ec2Instance) StartSSH(osUser string, portNumber uint64, publicKey []byte) ([]string, error) {
	zone := i.GetZone()
	region := zone[:len(zone)-1]
	s, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "sandbox",
		Config:            aws.Config{Region: aws.String(region)},
	})
	instanceId := i.ResourceID()
	placementAZ := i.PlacementAZ()

	if len(publicKey) > 0 {
		svc := ec2instanceconnect.New(s)
		err := uploadPublicKey(svc, publicKey, osUser, instanceId, placementAZ)
		if err != nil {
			return nil, err
		}
	}

	ssmClient := ssm.New(s)
	parameters := map[string][]*string{
		"portNumber": {aws.String(strconv.FormatUint(portNumber, 10))},
	}
	port := strconv.FormatUint(portNumber, 10)
	input := &ssm.StartSessionInput{
		DocumentName: aws.String("AWS-StartSSHSession"),
		Parameters:   map[string][]*string{"portNumber": []*string{&port}},
		Target:       &instanceId,
	}
	startSessionInput := &ssm.StartSessionInput{
		Parameters:   parameters,
		Target:       &instanceId,
		DocumentName: aws.String("AWS-StartSSHSession"),
	}
	output, err := ssmClient.StartSession(input)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	startSessionInputJSON, err := json.Marshal(startSessionInput)
	if err != nil {
		return nil, err
	}

	endpoint := ssmClient.Client.Endpoint
	//proxy := fmt.Sprintf("ProxyCommand=aws ssm start-session --target %s --document-name AWS-StartSSHSession --parameters 'portNumber=%s' --profile %s", instanceId, "22", "sandbox")
	//sshArgs := []string{"-o", proxy, "-vvvv"}
	return GetRunSessionPluginManager(string(payload), region, i.ProfileName, string(startSessionInputJSON), endpoint), err

}

func uploadPublicKey(client ec2instanceconnectiface.EC2InstanceConnectAPI, publicKey []byte, osUser, instanceID, availabilityZone string) error {
	pubKey := string(publicKey)
	out, err := client.SendSSHPublicKey(&ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: &availabilityZone,
		InstanceId:       &instanceID,
		InstanceOSUser:   &osUser,
		SSHPublicKey:     &pubKey,
	})
	if err != nil {

		return err
	}
	if !*out.Success {
		return fmt.Errorf("failed SendSSHPublicKey. RequestID: %s", *out.RequestId)
	}

	return nil
}
func TerminateSession(ssmClient ssmiface.SSMAPI, sessionID string) error {
	_, err := ssmClient.TerminateSession(&ssm.TerminateSessionInput{
		SessionId: &sessionID,
	})
	if err != nil {
		return err
	}
	return nil
}
func GetRunSessionPluginManager(payloadJSON, region, profile, inputJSON, endpoint string) []string {
	pluginName := "session-manager-plugin"
	// TODO allowing logging
	// https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html#configure-logs-linux
	// https://github.com/aws/aws-cli/blob/5f16b26/awscli/customizations/sessionmanager.py#L83-L89
	shell := exec.Command(pluginName, payloadJSON, region, "StartSession", profile, inputJSON, endpoint)
	shell.Stdout = os.Stdout
	shell.Stdin = os.Stdin
	shell.Stderr = os.Stderr
	//shell.Run()
	proxy := fmt.Sprintf("-oProxyCommand=%s", shellescape.QuoteCommand(strings.Split(shell.String(), " ")))

	sshArgs := []string{proxy}

	return sshArgs
}
