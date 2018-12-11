package version

import (
	"github.com/zendesk/goship/color"
	"github.com/zendesk/goship/config"
)

var (
	// VersionNumber defines a version
	VersionNumber = "1.0.0"
)

// CheckForNewVersion checks for new version
func CheckForNewVersion() {

	if config.GlobalConfig.Verbose {
		color.PrintYellow("Checking for newest version temporarily disabled")
	}
	return
}
