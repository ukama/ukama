/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package srvcrouter

import (
	"encoding/json"
	"net/url"

	"github.com/ukama/ukama/systems/common/config"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

const (
	RoutesExt   = "/routes"
	PatternExt  = "/pattern"
	ServicesExt = "/service"
	APIExt      = "/api"
	KeyPathExt  = "Path"
)

type QueryParams map[string]string

type ServiceRouter struct {
	C   *resty.Client
	Url *url.URL
}

func NewServiceRouter(path string) *ServiceRouter {
	url, _ := url.Parse(path)
	c := resty.New()
	c.SetDebug(true)
	rs := &ServiceRouter{
		C:   c,
		Url: url,
	}

	logrus.Tracef("Client created %+v for %s ", rs, rs.Url.String())
	return rs
}

func (r *ServiceRouter) RegisterService(apiIf config.ServiceApiIf) error {
	j, err := json.Marshal(apiIf)
	if err != nil {
		logrus.Errorf("Failed to encode service pattern into json. Error %v", err.Error())
		return err
	}

	logrus.Tracef("Requesting service router %s to add pattern %s for service.", (r.Url.String() + RoutesExt), string(j))
	resp, err := r.C.R().SetHeader("Content-Type", "application/json").SetBody(j).Post((r.Url.String() + RoutesExt))
	if err != nil {
		logrus.Errorf("Failed to register service to service router. Error %s", err.Error())
		return err
	}

	if resp.IsSuccess() {
		logrus.Infof("Service registered to service router. Response code %d", resp.StatusCode())
	} else {
		logrus.Errorf("Service failed to register itself to service router. Response code %d", resp.StatusCode())
	}

	return err
}
