package metrics

import (
	"fmt"
	"github.com/ukama/ukamaX/common/config"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func StartMetricsServer(conf *config.Metrics) {
	if conf.Enabled {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			logrus.Infof("Starting metrics server on port %d", conf.Port)
			err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
			if err != nil {
				logrus.WithError(err).Error("Error starting metrics server")
			}
		}()

	}
}
