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
			fmt.Printf(r.RenderShortOutput())
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Choose your ship: ")
		choose, _ := reader.ReadString('\n')
		idx, err := strconv.Atoi(strings.TrimSuffix(choose, "\n"))
		if len(resourcesList)+1 <= idx || idx < 1 || err != nil {
			return nil, fmt.Errorf("Unknown choose %s", choose)
		}
		return resourcesList[idx-1], nil
	} else {
		return nil, fmt.Errorf("No possible elements to display")
	}
}
