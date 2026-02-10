/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	cenums "github.com/ukama/ukama/testing/common/enums"
	"github.com/ukama/ukama/testing/services/dummy/dnode/config"
)

type prometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Name    string `json:"__name__"`
				Env     string `json:"env"`
				Job     string `json:"job"`
				Network string `json:"network"`
				Org     string `json:"org"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func getMetricValue(query string) (int, error) {
	resp, err := http.Get("http://prometheus:9090/api/v1/query?query=" + query)
	if err != nil {
		return 0, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Warnf("Failed to close response body: %v", closeErr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result prometheusResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	if len(result.Data.Result) == 0 {
		return 0, nil
	}

	valueStr := result.Data.Result[0].Value[1].(string)
	return strconv.Atoi(valueStr)
}

var profile cenums.Profile
var scenario cenums.SCENARIOS

func Worker(id string, updateChan chan config.WMessage, initial config.WMessage) {
	kpis := initial.Kpis
	profile = initial.Profile
	scenario = initial.Scenario

	fmt.Printf("Coroutine %s started with: %d, %s\n", id, profile, scenario)

	cleanup := func() {
		fmt.Printf("Shutting down coroutine %s with scenario: %s\n", id, scenario)
		for _, kpi := range kpis.KPIs {
			kpi.KPI.Delete(prometheus.Labels{"node_id": id, "type": kpi.Type})
		}
	}

	for {
		time.Sleep(1 * time.Second)
		select {
		case msg, ok := <-updateChan:
			if !ok {
				fmt.Printf("Coroutine %s with Scenario is: %s, which leads to coroutine shutdown.", id, scenario)
				return
			}
			profile = msg.Profile
			scenario = msg.Scenario
			if scenario == cenums.SCENARIO_BACKHAUL_DOWN || scenario == cenums.SCENARIO_NODE_OFF {
				cleanup()
				fmt.Printf("Coroutine %s with Scenario is: %s, which leads to coroutine shutdown.", id, scenario)
				return
			}
			fmt.Printf("Coroutine %s updated args: %d, %s\n", id, profile, scenario)
		default:
		}

		fmt.Printf("Coroutine %s working with: %d, %s\n", id, profile, scenario)
		pushMetrics(kpis, id, scenario, profile)

	}
}

func pushMetrics(kpis config.NodeKPIs, nodeID string, scenario cenums.SCENARIOS, profile cenums.Profile) {
	values := make(map[string]float64)

	for _, kpi := range kpis.KPIs {
		labels := prometheus.Labels{"node_id": nodeID, "type": kpi.Type}
		switch kpi.Key {
		case "unit_uptime":
			kpi.KPI.With(labels).Inc()
			continue
		case "trx_lte_core_active_ue":
			if scenario == "node_rf_off" {
				values[kpi.Key] = 0
				kpi.KPI.With(labels).Set(values[kpi.Key])

			} else {
				count, err := getSubscriber()
				if err != nil {
					fmt.Printf("Error getting subscriber: %s\n", err)
					continue
				}
				values[kpi.Key] = float64(count)
				kpi.KPI.With(labels).Set(values[kpi.Key])
			}
			continue
		// TODO: Can handle different scenario cases here for different KPIs
		default:
			switch profile {
			case cenums.PROFILE_MIN:
				values[kpi.Key] = kpi.Min + rand.Float64()*(kpi.Normal-kpi.Min)*0.3
			case cenums.PROFILE_MAX:
				values[kpi.Key] = kpi.Normal + rand.Float64()*(kpi.Max-kpi.Normal)*0.3
			default:
				values[kpi.Key] = kpi.Min + rand.Float64()*(kpi.Normal-kpi.Min)*0.3
			}
			kpi.KPI.With(labels).Set(values[kpi.Key])
		}
	}
}

func getSubscriber() (int, error) {
	activeCount, err := getMetricValue("number_of_sims")
	if err != nil {
		return 0, err
	}

	return activeCount, nil
}
