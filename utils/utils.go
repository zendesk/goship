package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/zendesk/goship/resources"
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
		choose = strings.TrimSuffix(choose, "\n")

		if choosenIP(choose) {
			return findResourceByIP(choose, resourcesList)
		}

		idx, err := strconv.Atoi(choose)
		if len(resourcesList)+1 <= idx || idx < 1 || err != nil {
			return nil, fmt.Errorf("unknown choose %s", choose)
		}
		return resourcesList[idx-1], nil
	} else {
		return nil, fmt.Errorf("no possible elements to display")
	}
}

func choosenIP(id string) bool {
	ip := net.ParseIP(id)
	return ip != nil
}

func findResourceByIP(ip string, resourcesList []resources.Resource) (resources.Resource, error) {
	for _, v := range resourcesList {
		if v.ConnectIdentifier(true, false) == ip {
			return v, nil
		}
	}
	return nil, fmt.Errorf("could not find resource with specified IP")
}
