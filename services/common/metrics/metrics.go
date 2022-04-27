package metrics

import (
	"fmt"
	"net/http"

	"github.com/ukama/openIoR/services/common/config"

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
