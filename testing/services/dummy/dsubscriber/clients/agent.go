/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package clients

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const CDREndpoint = "/v1/cdr"

type AddCDRRequest struct {
	Session       uint64 `json:"session"`
	Imsi          string `json:"imsi"`
	NodeId        string `json:"node_id"`
	Policy        string `json:"policy"`
	ApnName       string `json:"apn_name"`
	Ip            string `json:"ip"`
	StartTime     uint64 `json:"start_time"`
	EndTime       uint64 `json:"end_time"`
	LastUpdatedAt uint64 `json:"last_updated_at"`
	TxBytes       uint64 `json:"tx_bytes"`
	RxBytes       uint64 `json:"rx_bytes"`
	TotalBytes    uint64 `json:"total_bytes"`
}

type CDRClient interface {
	AddCDR(req AddCDRRequest) error
}

type cdrClient struct {
	u *url.URL
	R *client.Resty
}

func NewCDRClient(h string, options ...client.Option) *cdrClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &cdrClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (n *cdrClient) AddCDR(req AddCDRRequest) error {
	log.Debugf("Adding CDR: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("request marshal error. error: %w", err)
	}

	_, err = n.R.Post(n.u.String()+CDREndpoint+"/"+req.Imsi, b)
	if err != nil {
		log.Errorf("AddNetwork failure. error: %s", err.Error())

		return fmt.Errorf("AddNetwork failure: %w", err)
	}

	return nil
}
