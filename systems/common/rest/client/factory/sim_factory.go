/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package factory

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const SimFactoryEndpoint = "/v1/sims"

type SimFactoryClient interface {
	ReadSimCardInfo(iccid string) (*SimCardInfo, error)
}

type simFactoryClient struct {
	u *url.URL
	R *client.Resty
}

type ErrorMessage struct {
	Message string `json:"error"`
}

func NewSimFactoryClient(h string, options ...client.Option) *simFactoryClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &simFactoryClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (s *simFactoryClient) ReadSimCardInfo(iccid string) (*SimCardInfo, error) {
	log.Debugf("Reading sim card: %v", iccid)

	card := Sim{}

	resp, err := s.R.Get(s.u.String() + SimFactoryEndpoint + "/" + iccid)
	if err != nil {
		log.Errorf("Get sim card failure. error: %v", err)

		return nil, fmt.Errorf("getSimCard failure: %w", err)
	}

	log.Debugf("Unmarshaling resp from sim factory: %+v", resp)

	err = json.Unmarshal(resp.Body(), &card)
	if err != nil {
		log.Tracef("Failed to desrialize sim card info. Error message is %v", err)

		return nil, fmt.Errorf("simcard info deserailization failure: %w", err)
	}

	log.Infof("Sim card info: %+v", card.SimCardInfo)

	return card.SimCardInfo, nil
}

type SimCardInfo struct {
	Imsi           string `json:"imsi,omitempty"`
	Iccid          string `json:"iccid,omitempty"`
	Op             []byte `json:"op,omitempty"`
	Amf            []byte `json:"amf"`
	Key            []byte `json:"key,omitempty"`
	AlgoType       uint32 `json:"algo_type,omitempty"`
	UeDlAmbrBps    uint32 `json:"ue_dl_ambr_bps,omitempty"`
	UeUlAmbrBps    uint32 `json:"ue_ul_ambr_bps,omitempty"`
	Sqn            uint64 `json:"sqn,string,omitempty"`
	CsgIdPrsent    bool   `json:"c_sg_id_prsent,omitempty"`
	CsgId          uint32 `json:"csg_id,omitempty"`
	DefaultApnName string `json:"default_apn_name,omitempty"`
}

type Sim struct {
	SimCardInfo *SimCardInfo `json:"sim"`
}
