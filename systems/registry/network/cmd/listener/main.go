package main

import (
	"os"
	"time"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/registry/network/pkg"

	"github.com/ukama/ukama/systems/registry/network/cmd/version"

	rconf "github.com/num30/config"
	"github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/registry/network/pkg/queue"
)

const ServiceName = "network-listener"
const POD_NAME_ENV_VAR = "POD_NAME"

type QueueConfg struct {
	config.BaseConfig `mapstructure:",squash"`
	Queue             *config.Queue `default:"{}"`
	GrpcTimeout       time.Duration `default:"3s"`
	NetworkService    string        `default:"localhost:9090"`
}

func main() {
	ccmd.ProcessVersionArgument(ServiceName, os.Args, version.Version)

	config := &QueueConfg{}
	cr := rconf.NewConfReader(pkg.ServiceName + "-listener")
	err := cr.Read(config)
	if err != nil {
		logrus.Errorf("Error reading config: %v", err)
		return
	}
	pkg.IsDebugMode = config.DebugMode

	logrus.Infof("Creating listener. Queue: %s. Network: %s", config.Queue.SafeString(), config.NetworkService)
	listener, err := queue.NewQueueListener(config.NetworkService, config.Queue.Uri, config.GrpcTimeout, os.Getenv(POD_NAME_ENV_VAR))
	if err != nil {
		logrus.WithError(err).Error("Error starting listener")
		os.Exit(1)
	}
	logrus.Infof("Starting queue listening")
	err = listener.StartQueueListening()
	if err != nil {
		os.Exit(2)
	}
}
