package main

import (
	"github.com/ukama/ukamaX/cloud/net/cmd/version"
	"github.com/ukama/ukamaX/cloud/net/pkg/listener"
	ccmd "github.com/ukama/ukamaX/common/cmd"
	"github.com/ukama/ukamaX/common/config"
	"os"
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
