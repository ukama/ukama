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
	PrometheusInterval   int              `default:"60"`
	SchedulerInterval    time.Duration    `default:"60s"`
	OrgName              string           `default:"ukama"`
	PrometheusHost       string           `default:"http://localhost:9079"`
	MetricsKeyMapFile    string           `default:"metric-key-map.json"`
	FormatDecimalPoints  int              `default:"3"`
	MetricKeyMap         *MetricKeyMap
	Service              *uconf.Service
	Http            	 HttpServices
}

type Metric struct {
	Step             int         `json:"step" yaml:"step"`
	Key              string      `json:"key" yaml:"key"`
	Category         string      `json:"category" yaml:"category"`
	MetricType       string      `json:"metric_type" yaml:"metric_type"`
	TrendSensitivity float64     `json:"trend_sensitivity" yaml:"trend_sensitivity"`
	Thresholds       Thresholds  `json:"thresholds" yaml:"thresholds"`
	StateDirection   string      `json:"state_direction" yaml:"state_direction"`
}

type Thresholds struct {
	Min    float64 `json:"min" yaml:"min"`
	Medium float64 `json:"medium" yaml:"medium"`
	Max    float64 `json:"max" yaml:"max"`
	// For range direction
	LowWarning   float64 `json:"low_warning" yaml:"low_warning"`
	HighWarning  float64 `json:"high_warning" yaml:"high_warning"`
	LowCritical  float64 `json:"low_critical" yaml:"low_critical"`
	HighCritical float64 `json:"high_critical" yaml:"high_critical"`
}

type Metrics struct {
	Metrics []Metric `json:"metrics" yaml:"metrics"`
}

type MetricKeyMap map[string]Metrics

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
	metricKeyMap := new(MetricKeyMap)
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