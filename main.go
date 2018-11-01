package main

import (
	"github.com/zendesk/goship/cmd"
	"github.com/zendesk/goship/version"
)

func main() {
	version.CheckForNewVersion()
	cmd.Execute()
}
