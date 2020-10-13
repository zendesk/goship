package utils

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/zendesk/goship/resources"
	"golang.org/x/crypto/ssh"
)

// ChooseFromList provides prompt to choose from list of resources
func ChooseFromList(resourcesList []resources.Resource) (resources.Resource, error) {
	if len(resourcesList) == 1 {
		return resourcesList[0], nil
	} else if len(resourcesList) > 1 {
		for i, r := range resourcesList {
			fmt.Printf("%d. ", i+1)
			fmt.Print(r.RenderShortOutput())
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Choose your ship: ")
		choose, _ := reader.ReadString('\n')
		idx, err := strconv.Atoi(strings.TrimSuffix(choose, "\n"))
		if len(resourcesList)+1 <= idx || idx < 1 || err != nil {
			return nil, fmt.Errorf("unknown choose %s", choose)
		}
		return resourcesList[idx-1], nil
	} else {
		return nil, fmt.Errorf("no possible elements to display")
	}
}

//Handle Temp PEM key creation
func SavePrivPEMKey(fileName string, key *rsa.PrivateKey) error {
	var privateKey = &pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(key),
	}

	// returns err
	return ioutil.WriteFile(fileName, pem.EncodeToMemory(privateKey), 0600)
}

func SavePublicPEMKey(fileName string, pubkey *rsa.PublicKey) error {
	pub, err := ssh.NewPublicKey(pubkey)
	if err != nil {
		return err
	}
	// returns err
	return ioutil.WriteFile(fileName, ssh.MarshalAuthorizedKey(pub), 0600)
}

func DeleteTempKey(keyPath string) error {
	_, err := os.Stat(keyPath)
	print("Removing key")
	print(keyPath)
	if err == nil {
		if errRm := os.Remove(keyPath); errRm != nil {
			return errRm
		}
	}
	return err
}
