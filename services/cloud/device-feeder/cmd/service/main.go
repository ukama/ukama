package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/cloud/device-feeder/pkg/global"
	"github.com/ukama/ukama/services/cloud/device-feeder/pkg/multipl"

	"github.com/ukama/ukama/services/cloud/device-feeder/pkg"

	"github.com/ukama/ukama/services/cloud/device-feeder/cmd/version"

	ccmd "github.com/ukama/ukama/services/common/cmd"
	"github.com/ukama/ukama/services/common/config"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(global.ServiceName, os.Args, version.Version)
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

	ipResolve, err := pkg.NewDeviceIpResolver(serviceConfig.Net.Host, serviceConfig.Registry.TimeoutSeconds)
	if err != nil {
		logrus.Fatalf("Failed to create device ip resolver: %v", err)
	}

	exec := pkg.NewRequestExecutor(ipResolve, &serviceConfig.Device)

	listener, err := pkg.NewQueueListener(serviceConfig.Queue.Uri, os.Getenv(global.POD_NAME_ENV_VAR), m, exec, serviceConfig.Listener)
	if err != nil {
		logrus.WithError(err).Error("Error creating new listener")
		os.Exit(1)
	}

	exposeMetrics()

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
	config.LoadConfig(global.ServiceName, serviceConfig)
	global.IsDebugMode = serviceConfig.DebugMode
	if global.IsDebugMode {
		logrus.Infof("Configuration: %+v", serviceConfig)
	}
}

func exposeMetrics() {
	if serviceConfig.Metrics.Enabled {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			logrus.Infof("Starting metrics server on port %d", serviceConfig.Metrics.Port)
			err := http.ListenAndServe(fmt.Sprintf(":%d", serviceConfig.Metrics.Port), nil)
			if err != nil {
				logrus.WithError(err).Error("Error starting metrics server")
			}
		}()

	}
}
