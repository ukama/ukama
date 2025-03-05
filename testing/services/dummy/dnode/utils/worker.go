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
	cenums "github.com/ukama/ukama/testing/common/enums"
	"github.com/ukama/ukama/testing/services/dummy/dnode/config"
)

func Worker(id string, updateChan chan config.WMessage, initial config.WMessage) {
	count := 1.0
	kpis := initial.Kpis
	profile := initial.Profile
	scenario := initial.Scenario

	fmt.Printf("Coroutine %s started with: %d, %s\n", id, profile, scenario)

	for {
		count += 0.1
		time.Sleep(1 * time.Second)
		select {
		case msg, ok := <-updateChan:
			if !ok {
				fmt.Printf("Coroutine %s with Scenario is: %s, which leads to coroutine shutdown.", id, scenario)
				return
			}
			profile = msg.Profile
			scenario = msg.Scenario
			fmt.Printf("Coroutine %s updated args: %d, %s\n", id, profile, scenario)
		default:
		}

		fmt.Printf("Coroutine %s working with: %d, %s\n", id, profile, scenario)

		labels := prometheus.Labels{"nodeid": id}
		values := make(map[string]float64)

		for _, kpi := range kpis.KPIs {
			switch kpi.Key {
			case "unit_uptime":
				values[kpi.Key] = count
			// TODO: Can handle different scenario cases here
			default:
				switch profile {
				case cenums.PROFILE_MIN:
					values[kpi.Key] = kpi.Min + rand.Float64()*(kpi.Normal-kpi.Min)*0.1
				case cenums.PROFILE_MAX:
					values[kpi.Key] = kpi.Normal + rand.Float64()*(kpi.Max-kpi.Normal)*0.1
				default:
					values[kpi.Key] = kpi.Min + rand.Float64()*(kpi.Normal-kpi.Min)*0.1
				}
			}

			kpi.KPI.With(labels).Set(values[kpi.Key])
		}
	}
}
