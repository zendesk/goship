package utils

import (
	"bufio"
	"fmt"
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
		choice, _ := reader.ReadString('\n')
		idx, err := strconv.Atoi(strings.TrimSuffix(choice, "\n"))
		if len(resourcesList)+1 <= idx || idx < 1 || err != nil {
			return nil, fmt.Errorf("unknown choice %s", choice)
		}
		return resourcesList[idx-1], nil
	} else {
		return nil, fmt.Errorf("no possible elements to display")
	}
}
