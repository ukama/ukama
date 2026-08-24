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
	"sync"
	"time"

	"github.com/ukama/ukama/systems/common/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var startOnce sync.Once

func StartMetricsServer(conf *config.Metrics) {
	if conf == nil || !conf.Enabled {
		return
	}

	startOnce.Do(func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())

		srv := &http.Server{
			Addr:              fmt.Sprintf(":%d", conf.Port),
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		}

		go func() {
			logrus.Infof("Starting metrics server on port %d", conf.Port)

			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logrus.WithError(err).Error("Error starting metrics server")
			}
		}()
	})
}
