package main

import (
	"os"

	"github.com/num30/config"
	"github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/data-plan/base-rate/cmd/version"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/queue"
	// "github.com/ukama/ukama/systems/common/metrics"
)

const serviceName = pkg.ServiceName + "-listener"

func main() {
	ccmd.ProcessVersionArgument(serviceName, os.Args, version.Version)
	reader := config.NewConfReader(serviceName)
	conf := queue.QueueListenerConfig{}
	err := reader.Read(&conf)
	if err != nil {
		logrus.Fatalf("Failed to read config: %v", err)
	}
	// metrics.StartMetricsServer(&conf.Metrics)

	listener, err := queue.NewQueueListener(conf, serviceName, os.Getenv("POD_NAME"))
	if err != nil {
		logrus.Fatalf("Failed to create queue listener: %v", err)
	}

	err = listener.StartQueueListening()
	if err != nil {
		logrus.Fatalf("Failed to start queue listener: %v", err)
	}
}