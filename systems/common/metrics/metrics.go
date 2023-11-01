/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package metrics

import (
	"fmt"
	"net/http"

	"github.com/ukama/ukama/systems/common/config"

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
