package main

import (
	"github.com/num30/config"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/cloud/node/cmd/version"
	"github.com/ukama/ukama/services/cloud/node/pkg"
	"github.com/ukama/ukama/services/cloud/node/pkg/queue"
	ccmd "github.com/ukama/ukama/services/common/cmd"
	"github.com/ukama/ukama/services/common/metrics"
	"os"
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
	metrics.StartMetricsServer(&conf.Metrics)

	listener, err := queue.NewQueueListener(conf, serviceName, os.Getenv("POD_NAME"))
	if err != nil {
		logrus.Fatalf("Failed to create queue listener: %v", err)
	}
	listener.StartQueueListening()
}
