/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

package pkg

import (
	"encoding/json"
	"os"
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc   `default:"{}"`
	Queue            *uconf.Queue  `default:"{}"`
	Timeout          time.Duration `default:"20s"`
	Service          *uconf.Service
	MsgClient        *uconf.MsgClient `default:"{}"`
	OrgName          string           `default:"ukama"`
	Lookup           string           `default:"lookup:9090"`
	Http             HttpServices
	DNSMap           []OrgDNS `mapstructure:"-"`
	MeshNamespace    string   `default:"messaging"`
}

type OrgDNS struct {
	OrgName string `mapstructure:"org_name" json:"org_name"`
	DNS     string `mapstructure:"dns" json:"dns"`
}

type HttpServices struct {
	InitClient string `default:"api-gateway-init:8080"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout:        7 * time.Second,
			ListenerRoutes: []string{},
		},
		DNSMap: []OrgDNS{
			{
				OrgName: "ukama",
				DNS:     "localhost",
			},
		},
	}
}

func (c *Config) ParseDNSMapFromEnv() error {
	dnsMapEnv := os.Getenv("DNS_MAP")
	if dnsMapEnv == "" {
		return nil
	}

	var dnsMap []OrgDNS
	if err := json.Unmarshal([]byte(dnsMapEnv), &dnsMap); err != nil {
		return err
	}

	c.DNSMap = dnsMap
	return nil
}

func (c *Config) ToDNSMap() map[string]string {
	dnsMap := make(map[string]string)
	for _, orgDNS := range c.DNSMap {
		dnsMap[orgDNS.OrgName] = orgDNS.DNS
	}
	return dnsMap
}
