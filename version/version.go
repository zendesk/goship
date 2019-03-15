package version

import (
	"fmt"

	"github.com/tcnksm/go-latest"
	"github.com/zendesk/goship/color"
	"github.com/zendesk/goship/config"
)

var (
	// VersionNumber defines a version
	VersionNumber = "1.0.4"
)

// CheckForNewVersion checks for new version
func CheckForNewVersion() {

	githubTag := &latest.GithubTag{
		Owner:      "zendesk",
		Repository: "goship",
	}

	result, err := latest.Check(githubTag, VersionNumber)
	if err != nil {
		if config.GlobalConfig.Verbose {
			color.PrintYellow(fmt.Sprintf("Error while checking for the newest version: %s\n", err))
		}
		return
	}

	if result.Outdated && config.GlobalConfig.Verbose {
		color.PrintYellow(fmt.Sprintf("Newer version (%s) available! Checkout project repository to upgrade to the newest version.\n", result.Current))
	}
	return
}
