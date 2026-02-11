/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

import (
	"encoding/json"
	"os"
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig     `mapstructure:",squash"`
	EtcdHost          	 string
	DialTimeoutSecond 	 time.Duration
	Grpc                 *uconf.Grpc      `default:"{}"`
	Queue                *uconf.Queue     `default:"{}"`
	Timeout              time.Duration    `default:"3s"`
	MsgClient            *uconf.MsgClient `default:"{}"`
	SchedulerInterval    time.Duration    `default:"15s"`
	OrgName              string           `default:"ukama"`
	PrometheusHost       string           `default:"http://localhost:9079"`
	PrometheusInterval   int              `default:"15"`
	MetricsKeyMapFile    string           `default:"metric-key-map.json"`
	MetricKeyMap         *MetricKeyMap
	Service              *uconf.Service
	Http            	 HttpServices
}

type Metric struct {
	Name string `json:"name" yaml:"name"`
	Interval int `json:"interval" yaml:"interval"`
	Step int `json:"step" yaml:"step"`
	Category string `json:"category" yaml:"category"`
	Metric []MetricItem `json:"metric" yaml:"metric"`
}

type MetricItem struct {
	Key string `json:"key" yaml:"key"`
	Type string `json:"type" yaml:"type"`
}

type MetricKeyMap struct {
	Metrics []Metric `json:"metrics" yaml:"metrics"`
}

type HttpServices struct {
	InitClient string `default:"http://api-gateway-init:8080"`
}

func NewConfig(name string) *Config {
	return &Config{
		EtcdHost:          "localhost:2379",
		DialTimeoutSecond: 5 * time.Second,
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
		},
	}
}

func LoadMetricKeyMap(config *Config) (*MetricKeyMap, error) {
	metricKeyMap := &MetricKeyMap{}
	metricKeyMapFile := config.MetricsKeyMapFile
	bytes, err := os.ReadFile(metricKeyMapFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, metricKeyMap)
	if err != nil {
		return nil, err
	}
	return metricKeyMap, nil
}