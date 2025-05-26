/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"net/http"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/ukama"

	log "github.com/sirupsen/logrus"
	crest "github.com/ukama/ukama/systems/common/rest"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
)

const (
	failedRequestMsg  = "invalid content. request has failed"
	pendingRequestMsg = "partial content. request is still ongoing"
)

type Network interface {
	GetNetwork(string) (*creg.NetworkInfo, error)
	CreateNetwork(string, string, []string, []string, float64, float64, uint32, bool) (*creg.NetworkInfo, error)
}

type network struct {
	nc creg.NetworkClient
}

func NewNetworkClientSet(ntwk creg.NetworkClient) Network {
	n := &network{
		nc: ntwk,
	}

	return n
}

func (n *network) GetNetwork(id string) (*creg.NetworkInfo, error) {
	net, err := n.nc.Get(id)
	if err != nil {
		return nil, client.HandleRestErrorStatus(err)
	}

	if net.SyncStatus == ukama.StatusTypeUnknown.String() || net.SyncStatus == ukama.StatusTypeFailed.String() {
		log.Error(failedRequestMsg)

		return nil, crest.HttpError{
			HttpCode: http.StatusUnprocessableEntity,
			Message:  failedRequestMsg,
		}
	}

	if net.SyncStatus == ukama.StatusTypePending.String() || net.SyncStatus == ukama.StatusTypeProcessing.String() {
		log.Warn(pendingRequestMsg)

		return net, crest.HttpError{
			HttpCode: http.StatusPartialContent,
			Message:  pendingRequestMsg,
		}
	}

	return net, nil
}

func (n *network) CreateNetwork(orgName, NetworkName string, allowedCountries,
	allowedNetworks []string, budget, overdraft float64, trafficPolicy uint32,
	paymentLinks bool) (*creg.NetworkInfo, error) {
	net, err := n.nc.Add(creg.AddNetworkRequest{
		OrgName:          orgName,
		NetName:          NetworkName,
		AllowedCountries: allowedCountries,
		AllowedNetworks:  allowedNetworks,
		Budget:           budget,
		Overdraft:        overdraft,
		TrafficPolicy:    trafficPolicy,
		PaymentLinks:     paymentLinks,
	})
	if err != nil {
		return nil, client.HandleRestErrorStatus(err)
	}

	return net, nil
}
