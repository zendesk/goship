package config

import (
	"fmt"
	"os"

	"github.com/zendesk/goship/color"
)

// VariablesToCheck defines list of variables to check in env in order to print warning when set
var VariablesToCheck = []string{"AWS_PROFILE", "AWS_ACCESS_KEY_ID", "AWS_ACCESS_KEY", "AWS_SECRET_ACCESS_KEY", "AWS_SECRET_KEY"}

// CheckEnv checks and warns if particular variables form VariablesToCheck are set in env.
func CheckEnv() {
	for _, e := range VariablesToCheck {
		if len(os.Getenv(e)) > 0 {
			color.PrintYellow(fmt.Sprintf("Warning! %s env variable is set and will override credentials from ~/.aws/credentials file. This will affect results.\n", e))
		}
	}
}
