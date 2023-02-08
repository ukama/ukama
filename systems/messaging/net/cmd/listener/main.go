package main

import (
	"os"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/messaging/net/cmd/version"
	"github.com/ukama/ukama/systems/messaging/net/pkg/listener"
)

var listenerConfig *listener.ListenerConfig

const ListenerName = "net-listener"

func main() {
	ccmd.ProcessVersionArgument(ListenerName, os.Args, version.Version)
	initConfig()
	listener.StartListener(listenerConfig)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	listenerConfig = listener.NewLiseterConfig()
	config.LoadConfig(ListenerName, listenerConfig)
}
