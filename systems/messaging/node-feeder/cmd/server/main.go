package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/num30/config"

	"github.com/ukama/ukama/systems/messaging/node-feeder/cmd/version"
	"github.com/ukama/ukama/systems/messaging/node-feeder/pkg"
	"github.com/ukama/ukama/systems/messaging/node-feeder/pkg/global"
	"github.com/ukama/ukama/systems/messaging/node-feeder/pkg/multipl"
	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/rest/client"
	ic "github.com/ukama/ukama/systems/common/rest/client/initclient"
)

var serviceConfig = pkg.NewConfig(global.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(global.ServiceName, os.Args, version.Version)
	initConfig()

	//registryClient := multipl.NewRegistryProvider(serviceConfig.Registry.Host, serviceConfig.Registry.TimeoutSeconds, serviceConfig.DebugMode)

	pub, err := multipl.NewQPub(serviceConfig.Queue.Uri, global.ServiceName, serviceConfig.Registry, os.Getenv(global.POD_NAME_ENV_VAR))
	if err != nil {
		log.Fatalf("Failed to create publisher: %v", err)
	}

	regUrl, err := ic.GetHostUrl(ic.NewInitClient(serviceConfig.Http.InitClient, client.WithDebug(serviceConfig.DebugMode)),
		ic.CreateHostString(serviceConfig.OrgName, "registry"), &serviceConfig.OrgName)
	if err != nil {
		log.Errorf("Failed to resolve registry address: %v", err)
	}

	m := multipl.NewRequestMultiplier(regUrl.String(), pub)

	ipResolve, err := pkg.NewNodeIpResolver(serviceConfig.Net, serviceConfig.TimeoutSeconds)
	if err != nil {
		log.Fatalf("Failed to create device ip resolver: %v", err)
	}

	exec := pkg.NewRequestExecutor(ipResolve, serviceConfig.DevicePort, serviceConfig.TimeoutSeconds)

	listener, err := pkg.NewQueueListener(global.ServiceName, serviceConfig.Queue.Uri, os.Getenv(global.POD_NAME_ENV_VAR), m, exec, serviceConfig.Listener)
	if err != nil {
		log.WithError(err).Error("Error creating new listener")
		os.Exit(1)
	}

	exposeMetrics()

	log.Infof("Starting queue listening")
	err = listener.StartQueueListening()
	if err != nil {
		log.WithError(err).Error("Error starting queue listening")
		os.Exit(2)
	}
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	err := config.NewConfReader(global.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if serviceConfig.DebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

	global.IsDebugMode = serviceConfig.DebugMode
}

func exposeMetrics() {
	if serviceConfig.Metrics.Enabled {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			log.Infof("Starting metrics server on port %d", serviceConfig.Metrics.Port)
			err := http.ListenAndServe(fmt.Sprintf(":%d", serviceConfig.Metrics.Port), nil)
			if err != nil {
				log.WithError(err).Error("Error starting metrics server")
			}
		}()
	}
}
