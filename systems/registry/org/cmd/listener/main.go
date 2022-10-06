package main

import (
	"os"

	sig "github.com/ukama/ukama/systems/common/signal"

	"github.com/num30/config"
	"github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/registry/org/cmd/version"
	"github.com/ukama/ukama/systems/registry/org/pkg"
	"github.com/ukama/ukama/systems/registry/org/pkg/queue"
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
	metrics.StartMetricsServer(conf.Metrics)

	listener, err := queue.NewQueueListener(conf, os.Getenv("POD_NAME"))
	if err != nil {
		logrus.Fatalf("Failed to create queue listener: %v", err)
	}
	sig.HandleSigterm(func() {
		listener.Close()
	})

	err = listener.StartQueueListening()
	if err != nil {
		logrus.Fatalf("Failed to start queue listener: %v", err)
	}
}
