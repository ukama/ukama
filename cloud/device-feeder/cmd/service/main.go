package main

import (
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/cloud/device-feeder/pkg/multipl"
	"os"

	"github.com/ukama/ukamaX/cloud/device-feeder/pkg"

	"github.com/ukama/ukamaX/cloud/device-feeder/cmd/version"

	ccmd "github.com/ukama/ukamaX/common/cmd"
	"github.com/ukama/ukamaX/common/config"
)

var serviceConfig *pkg.Config

const POD_NAME_ENV_VAR = "POD_NAME"

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	registryClient, err := multipl.NewRegistryClient(serviceConfig.Registry.Host, serviceConfig.Registry.TimeoutSeconds)
	if err != nil {
		logrus.Fatalf("Failed to create registry client: %v", err)
	}

	pub, err := multipl.NewQueuePublisher(serviceConfig.Queue.Uri)
	if err != nil {
		logrus.Fatalf("Failed to create publisher: %v", err)
	}

	m := multipl.NewRequestMultiplier(registryClient, pub)

	ipResolve, err := pkg.NewDeviceIpResolver(serviceConfig.Registry.Host, serviceConfig.Registry.TimeoutSeconds)
	if err != nil {
		logrus.Fatalf("Failed to create device ip resolver: %v", err)
	}

	exec := pkg.NewRequestExecutor(ipResolve, &serviceConfig.Device)

	listener, err := pkg.NewQueueListener(serviceConfig.Queue.Uri, os.Getenv(POD_NAME_ENV_VAR), m, exec, serviceConfig.Listener)
	if err != nil {
		logrus.WithError(err).Error("Error creating new listener")
		os.Exit(1)
	}

	logrus.Infof("Starting queue listening")
	err = listener.StartQueueListening()
	if err != nil {
		logrus.WithError(err).Error("Error starting queue listening")
		os.Exit(2)
	}
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, serviceConfig)
	pkg.IsDebugMode = serviceConfig.DebugMode
}
