package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/ukama/ukama/services/cloud/registry/cmd/version"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/cloud/registry/pkg/queue"
	ccmd "github.com/ukama/ukamaX/common/cmd"
)

const ServiceName = "registry-listener"
const DEFAUTL_GRPC_TIMEOUT = 5
const TIMEOUT_ENV_VAR_NAME = "GRPC_TIMEOUT_SECONDS"
const POD_NAME_ENV_VAR = "POD_NAME"

func main() {
	ccmd.ProcessVersionArgument(ServiceName, os.Args, version.Version)

	logrus.Info("Configure access to 'registry' and rabbitmq by setting REGISTRY and QUEUE environment variables")
	registry, ok := os.LookupEnv("REGISTRY")
	if !ok {
		registry = "registry:9090"
	}

	queueStr, ok := os.LookupEnv("QUEUE")
	if !ok {
		queueStr = "amqp://guest:guest@rabbitmq:5672"
	}
	logrus.Info("Configure grpc timeout by setting ", TIMEOUT_ENV_VAR_NAME, " environment variable")
	logrus.Info("Configure service id by setting  ", POD_NAME_ENV_VAR, " environment variable")

	logrus.Infof("Creating listener. Queue: %s. Registry: %s", queueStr[strings.LastIndex(queueStr, "@"):], registry)
	listener, err := queue.NewQueueListener(registry, queueStr, readTimeout(), os.Getenv(POD_NAME_ENV_VAR))
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

func readTimeout() int {
	timeOutVar, ext := os.LookupEnv(TIMEOUT_ENV_VAR_NAME)
	if ext {
		timeOut, err := strconv.Atoi(timeOutVar)
		if err != nil {
			logrus.Warningf("Error parsing timeout. Error: %v", err)
			return DEFAUTL_GRPC_TIMEOUT
		}
		return timeOut
	}
	return DEFAUTL_GRPC_TIMEOUT
}
