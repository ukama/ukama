package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/cloud/org-metrics/cmd/version"
	"github.com/ukama/ukama/services/cloud/org-metrics/pkg"
	"github.com/ukama/ukama/services/common/grpc"
	"github.com/ukama/ukama/services/common/metrics"
	"net/http"
	"os"

	"github.com/num30/config"
	reg "github.com/ukama/ukama/services/cloud/registry/pb/gen"
	ccmd "github.com/ukama/ukama/services/common/cmd"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()

	// those a service metirics, like number of requests, and note metrics that is serves
	metrics.StartMetricsServer(serviceConfig.Metrics)

	promReg := prometheus.NewRegistry()

	// create grpc connection
	// panics if connection fails
	conn := grpc.CreateGrpcConn(serviceConfig.Registry.GrpcService)

	regClient := reg.NewRegistryServiceClient(conn)

	promReg.MustRegister(pkg.NewMetricsCollector(regClient, serviceConfig.Registry.Timeout, serviceConfig.Registry.PollInterval))

	http.Handle("/", promhttp.HandlerFor(
		promReg,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))

	http.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}))

	logrus.Info("Starting server on ", serviceConfig.Server.Port)
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serviceConfig.Server.Port), nil))
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = &pkg.Config{}
	reader := config.NewConfReader(pkg.ServiceName)
	err := reader.Read(serviceConfig)
	if err != nil {
		logrus.Fatalf("Failed to read config: %v", err)
	}

	pkg.IsDebugMode = serviceConfig.DebugMode
}
