/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	_ "embed"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
	"gopkg.in/yaml.v3"
)

//go:embed default-metrics.yaml
var defaultMetricsYaml []byte

type NameUpdate struct {
	Required bool   `json:"required" default:"false"`
	Slice    string `json:"slice" default:""`
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Server            rest.HttpConfig
	Services          GrpcEndpoints  `mapstructure:"services"`
	HttpServices      HttpEndpoints  `mapstructure:"httpServices"`
	MetricsServer     config.Metrics `mapstructure:"metrics"`
	MetricsStore      string         `default:"http://localhost:8080"`
	Auth              *config.Auth   `mapstructure:"auth"`
	MetricsConfig     *MetricsConfig
	OrgName           string
	Period            time.Duration `default:"5s"`
	MetricsKeyMapFile string        `default:"default-metrics.yaml"`
	Http              HttpServices
}

type HttpServices struct {
	InitClient string `default:"api-gateway-init:8080"`
}

type Metric struct {
	NeedRate bool   `json:"needRate" yaml:"needRate"`
	Metric   string `json:"metric" yaml:"metric"`
	// Range vector duration used in Rate func https://prometheus.io/docs/prometheus/latest/querying/basics/#time-durations
	// if NeedRate is false then this field is ignored
	// Example: 1d or 5h, or 30s
	RateInterval string `json:"rateInterval" yaml:"rateInterval"`
}

type MetricsConfig struct {
	Metrics             map[string]Metric  `json:"metrics"`
	MetricsServer       string             `default:"http://localhost:9090"`
	Timeout             time.Duration
	DefaultRateInterval string
}

func loadMetricsFromFile(path string) (map[string]Metric, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return loadMetricsFromBytes(data)
}

type GrpcEndpoints struct {
	Timeout   time.Duration
	Exporter  string
	Sanitizer string
}

type HttpEndpoints struct {
	Timeout     time.Duration
	NodeMetrics string
}

func NewConfig() *Config {
	defaultCors := cors.DefaultConfig()
	defaultCors.AllowWildcard = true
	defaultCors.AllowOrigins = []string{"http://localhost", "https://localhost"}

	return &Config{
		BaseConfig: config.BaseConfig{
			DebugMode: false,
		},

		Services: GrpcEndpoints{
			Timeout:   3 * time.Second,
			Exporter:  "0.0.0.0:9090",
			Sanitizer: "sanitizer:9090",
		},
		HttpServices: HttpEndpoints{
			Timeout:     3 * time.Second,
			NodeMetrics: "http://localhost",
		},

		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},
		MetricsServer: *config.DefaultMetrics(),

		MetricsConfig: &MetricsConfig{
			Metrics:             make(map[string]Metric),
			Timeout:             time.Second * 5,
			DefaultRateInterval: "5m",
		},
		Auth:   config.LoadAuthHostConfig("auth"),
		Period: time.Second * 5,
	}
}

// loadDefaultMetrics loads the embedded default-metrics.yaml into a map. Used when no override file is set or loadable.
func loadDefaultMetrics() (map[string]Metric, error) {
	m, err := loadMetricsFromBytes(defaultMetricsYaml)
	if err != nil {
		return nil, err
	}
	if m == nil {
		m = make(map[string]Metric)
	}
	return m, nil
}

// loadMetricsFromBytes parses YAML (flat map or nested metrics-gateway format) into map[string]Metric.
func loadMetricsFromBytes(data []byte) (map[string]Metric, error) {
	var m map[string]Metric
	if err := yaml.Unmarshal(data, &m); err == nil && len(m) > 0 {
		return m, nil
	}
	var nested struct {
		MetricsGateway struct {
			APIGatewayConfig struct {
				MetricsConfig struct {
					Metrics map[string]Metric `yaml:"metrics"`
				} `yaml:"metricsConfig"`
			} `yaml:"apiGatewayConfig"`
		} `yaml:"metrics-gateway"`
	}
	if err := yaml.Unmarshal(data, &nested); err != nil {
		return nil, err
	}
	m = nested.MetricsGateway.APIGatewayConfig.MetricsConfig.Metrics
	if m == nil {
		m = make(map[string]Metric)
	}
	return m, nil
}

// ApplyMetricsFromEnvOverride sets MetricsConfig.Metrics from embedded default-metrics.yaml, then overlays
// the file at MetricsKeyMapFile (env METRICS_KEY_MAP_FILE) if set and loadable. Call after LoadConfig.
// If the override file is missing or invalid, defaults are still used.
func ApplyMetricsFromEnvOverride(c *Config) {
	if c == nil {
		return
	}
	if c.MetricsConfig == nil {
		c.MetricsConfig = &MetricsConfig{}
	}
	if c.MetricsConfig.Metrics == nil {
		c.MetricsConfig.Metrics = make(map[string]Metric)
	}
	// 1) Load defaults from embedded default-metrics.yaml
	defaults, err := loadDefaultMetrics()
	if err != nil {
		logrus.Warnf("failed to load embedded default metrics: %v", err)
	} else {
		for k, v := range defaults {
			c.MetricsConfig.Metrics[k] = v
		}
	}
	// 2) Overlay from MetricsKeyMapFile if set and loadable.
	// Prefer METRICS_KEY_MAP_FILE env (used by Helm) since viper may not bind METRICS_KEY_MAP_FILE to MetricsKeyMapFile.
	path := os.Getenv("METRICS_KEY_MAP_FILE")
	if path == "" {
		path = c.MetricsKeyMapFile
	}
	if path == "" {
		return
	}
	override, err := loadMetricsFromFile(path)
	if err != nil {
		logrus.Warnf("failed to load metrics from %s (using defaults): %v", path, err)
		return
	}
	logrus.Infof("loaded metrics override from %s (%d keys)", path, len(override))
	for k, v := range override {
		c.MetricsConfig.Metrics[k] = v
	}
}
