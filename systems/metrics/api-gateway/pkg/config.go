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
	Http              HttpEndpoints  `mapstructure:"http"`
	MetricsServer     config.Metrics `mapstructure:"metrics"`
	MetricsStore      string         `default:"http://localhost:8080"`
	Auth              *config.Auth   `mapstructure:"auth"`
	MetricsConfig     *MetricsConfig
	OrgName           string
	Period            time.Duration `default:"5s"`
	MetricsKeyMapFile string        `default:"default-metrics.yaml"`
}

type Metric struct {
	NeedRate bool   `json:"needRate" yaml:"needRate"`
	Metric   string `json:"metric" yaml:"metric"`
	// Range vector duration used in Rate func https://prometheus.io/docs/prometheus/latest/querying/basics/#time-durations
	// if NeedRate is false then this field is ignored
	// Example: 1d or 5h, or 30s
	RateInterval  string    `json:"rateInterval" yaml:"rateInterval"`
	Unit          string    `json:"unit,omitempty" yaml:"unit,omitempty"`
	Format        string    `json:"format,omitempty" yaml:"format,omitempty"`
	TickInterval  int       `json:"tickInterval,omitempty" yaml:"tickInterval,omitempty"`
	TickPositions []int     `json:"tickPositions,omitempty" yaml:"tickPositions,omitempty"`
	Threshold     Threshold `json:"threshold,omitempty" yaml:"threshold,omitempty"`
}

type Threshold struct {
	Min    float64 `json:"min" yaml:"min"`
	Normal float64 `json:"normal" yaml:"normal"`
	Max    float64 `json:"max" yaml:"max"`
}

type MetricsConfig struct {
	// nodeType (tnode/anode/cnode/system) → genericKey → Metric
	Metrics             map[string]map[string]Metric `json:"metrics"`
	MetricsServer       string                       `default:"http://localhost:9090"`
	Timeout             time.Duration
	DefaultRateInterval string
}

func loadMetricsFromFile(path string) (map[string]map[string]Metric, error) {
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
	Reasoning string
}

type HttpEndpoints struct {
	Timeout     time.Duration
	NodeMetrics string
	InitClient  string `default:"api-gateway-init:8080"`
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
			Exporter:  "exporter:9090",
			Sanitizer: "sanitizer:9090",
			Reasoning: "reasoning:9090",
		},
		Http: HttpEndpoints{
			Timeout:     3 * time.Second,
			NodeMetrics: "http://localhost",
		},

		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},
		MetricsServer: *config.DefaultMetrics(),

		MetricsConfig: &MetricsConfig{
			Metrics:             make(map[string]map[string]Metric),
			Timeout:             time.Second * 5,
			DefaultRateInterval: "5m",
		},
		Auth:   config.LoadAuthHostConfig("auth"),
		Period: time.Second * 5,
	}
}

// loadDefaultMetrics loads the embedded default-metrics.yaml into a nested map. Used when no override file is set or loadable.
func loadDefaultMetrics() (map[string]map[string]Metric, error) {
	m, err := loadMetricsFromBytes(defaultMetricsYaml)
	if err != nil {
		return nil, err
	}
	if m == nil {
		m = make(map[string]map[string]Metric)
	}
	return m, nil
}

// loadMetricsFromBytes parses YAML with structure nodeType -> genericKey -> Metric.
func loadMetricsFromBytes(data []byte) (map[string]map[string]Metric, error) {
	var m map[string]map[string]Metric
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	if m == nil {
		m = make(map[string]map[string]Metric)
	}
	return m, nil
}

// ApplyMetricsFromEnvOverride sets MetricsConfig.Metrics from embedded default-metrics.yaml, then overlays
// the file at MetricsKeyMapFile (env METRICS_KEY_MAP_FILE) if set and loadable. Call after LoadConfig.
// If the override file is missing or invalid, defaults are still used.
//
// Metrics is a two-level map: nodeType (tnode/anode/cnode/system) -> genericKey -> Metric.
// Override files must use the same nested structure; overlays are applied per node-type bucket.
func ApplyMetricsFromEnvOverride(c *Config) {
	if c == nil {
		return
	}
	if c.MetricsConfig == nil {
		c.MetricsConfig = &MetricsConfig{}
	}
	if c.MetricsConfig.Metrics == nil {
		c.MetricsConfig.Metrics = make(map[string]map[string]Metric)
	}
	// 1) Load defaults from embedded default-metrics.yaml
	defaults, err := loadDefaultMetrics()
	if err != nil {
		logrus.Warnf("failed to load embedded default metrics: %v", err)
	} else {
		for nodeType, keyMap := range defaults {
			if c.MetricsConfig.Metrics[nodeType] == nil {
				c.MetricsConfig.Metrics[nodeType] = make(map[string]Metric)
			}
			for k, v := range keyMap {
				c.MetricsConfig.Metrics[nodeType][k] = v
			}
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
	logrus.Infof("loaded metrics override from %s (%d node-type buckets)", path, len(override))
	for nodeType, keyMap := range override {
		if c.MetricsConfig.Metrics[nodeType] == nil {
			c.MetricsConfig.Metrics[nodeType] = make(map[string]Metric)
		}
		for k, v := range keyMap {
			c.MetricsConfig.Metrics[nodeType][k] = v
		}
	}
}
