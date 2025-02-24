/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/ukama/ukama/testing/services/dummy/dnode/config"
)

func Worker(id string, updateChan chan config.WMessage, initial config.WMessage) {
	count := 1.0
	kpis := initial.Kpis
	profile := initial.Profile
	scenario := initial.Scenario

	fmt.Printf("Coroutine %s started with: %d, %s\n", id, profile, scenario)

	for {
		select {
		case msg := <-updateChan:
			profile = msg.Profile
			scenario = msg.Scenario
			fmt.Printf("Coroutine %s updated args: %d, %s\n", id, profile, scenario)
		default:
		}

		count += 0.1
		time.Sleep(1 * time.Second)
		fmt.Printf("Coroutine %s working with: %d, %s\n", id, profile, scenario)

		labels := prometheus.Labels{"nodeid": id}
		values := make(map[string]float64)

		for _, kpi := range kpis.KPIs {
			if kpi.Key == "unit_uptime" {
				values[kpi.Key] = count
			} else {
				switch profile {
				case config.PROFILE_MIN:
					values[kpi.Key] = kpi.Min + rand.Float64()*(kpi.Normal-kpi.Min)*0.1
				case config.PROFILE_MAX:
					values[kpi.Key] = kpi.Normal + rand.Float64()*(kpi.Max-kpi.Normal)*0.1
				default:
					values[kpi.Key] = kpi.Min + rand.Float64()*(kpi.Normal-kpi.Min)*0.1
				}
			}
			kpi.KPI.With(labels).Set(values[kpi.Key])
		}
	}
}
