/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"encoding/json"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
)

const PackageEndpoint = "/v1/packages"

type PackageInfo struct {
	Id            string   `json:"uuid"`
	Name          string   `json:"name"`
	From          string   `json:"from" validation:"required"`
	To            string   `json:"to" validation:"required"`
	OrgId         string   `json:"org_id" validation:"required"`
	OwnerId       string   `json:"owner_id" validation:"required"`
	SimType       string   `json:"sim_type" validation:"required"`
	SmsVolume     int64    `json:"sms_volume" validation:"required"`
	VoiceVolume   int64    `json:"voice_volume" default:"0"`
	DataVolume    int64    `json:"data_volume" validation:"required"`
	VoiceUnit     string   `json:"voice_unit" validation:"required"`
	DataUnit      string   `json:"data_unit" validation:"required"`
	Type          string   `json:"type" validation:"required"`
	Flatrate      bool     `json:"flat_rate" default:"false"`
	Amount        float64  `json:"amount" default:"0.00"`
	Markup        float64  `json:"markup" default:"0.00"`
	Apn           string   `json:"apn" default:"ukama.tel"`
	BaserateId    string   `json:"baserate_id" validation:"required"`
	IsActive      bool     `json:"active"`
	Duration      uint64   `json:"duration,string"`
	Overdraft     float64  `json:"overdraft"`
	TrafficPolicy uint32   `json:"traffic_policy"`
	Networks      []string `json:"networks"`
	SyncStatus    string   `json:"sync_status,omitempty"`
}

type Package struct {
	PackageInfo *PackageInfo `json:"package"`
}

type AddPackageRequest struct {
	Name          string   `json:"name" validation:"required"`
	From          string   `json:"from" validation:"required"`
	To            string   `json:"to" validation:"required"`
	OrgId         string   `json:"org_id" validation:"required"`
	OwnerId       string   `json:"owner_id" validation:"required"`
	SimType       string   `json:"sim_type" validation:"required"`
	SmsVolume     int64    `json:"sms_volume" validation:"required"`
	VoiceVolume   int64    `json:"voice_volume" default:"0"`
	DataVolume    int64    `json:"data_volume" validation:"required"`
	VoiceUnit     string   `json:"voice_unit" validation:"required"`
	DataUnit      string   `json:"data_unit" validation:"required"`
	Duration      uint64   `json:"duration" validation:"required"`
	Type          string   `json:"type" validation:"required"`
	Flatrate      bool     `json:"flat_rate" default:"false"`
	Amount        float64  `json:"amount" default:"0.00"`
	Markup        float64  `json:"markup" default:"0.00"`
	Apn           string   `json:"apn" default:"ukama.tel"`
	Active        bool     `json:"active" validation:"required"`
	BaserateId    string   `json:"baserate_id" validation:"required"`
	Overdraft     float64  `json:"overdraft"`
	TrafficPolicy uint32   `json:"traffic_policy"`
	Networks      []string `json:"networks"`
}

type PackageClient interface {
	Get(Id string) (*PackageInfo, error)
	Add(req AddPackageRequest) (*PackageInfo, error)
}

type packageClient struct {
	u *url.URL
	R *Resty
}

func NewPackageClient(h string) *packageClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error %s", h, err.Error())
	}

	return &packageClient{
		u: u,
		R: NewResty(),
	}
}

// TODO check upstream returns payload
func (p *packageClient) Add(req AddPackageRequest) (*PackageInfo, error) {
	log.Debugf("Adding package: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	pkg := Package{}

	resp, err := p.R.Post(p.u.String()+PackageEndpoint, b)
	if err != nil {
		log.Errorf("AddPackage failure. error: %s", err.Error())

		return nil, fmt.Errorf("AddPackage failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &pkg)
	if err != nil {
		log.Tracef("Failed to deserialize package info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("package info deserailization failure: %w", err)
	}

	log.Infof("Package Info: %+v", pkg.PackageInfo)

	return pkg.PackageInfo, nil
}

func (p *packageClient) Get(id string) (*PackageInfo, error) {
	log.Debugf("Getting package: %v", id)

	pkg := Package{}

	resp, err := p.R.Get(p.u.String() + PackageEndpoint + "/" + id)
	if err != nil {
		log.Errorf("GetPackage failure. error %s", err.Error())

		return nil, fmt.Errorf("GetPackage failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &pkg)
	if err != nil {
		log.Tracef("Failed to deserialize org info. Error message is %s", err.Error())

		return nil, fmt.Errorf("package info deserailization failure: %w", err)
	}

	log.Infof("Package Info: %+v", pkg.PackageInfo)

	return pkg.PackageInfo, nil
}
