/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"errors"
	"net/http"

	"github.com/ukama/ukama/systems/common/types"

	log "github.com/sirupsen/logrus"
	crest "github.com/ukama/ukama/systems/common/rest"
	cclient "github.com/ukama/ukama/systems/common/rest/client"
)

const (
	failedRequestMsg  = "invalid content. request has failed"
	pendingRequestMsg = "partial content. request is still ongoing"
)

type Network interface {
	GetNetwork(string) (*cclient.NetworkInfo, error)
	CreateNetwork(string, string, []string, []string, float64, float64, uint32, bool) (*cclient.NetworkInfo, error)
}

type network struct {
	nc cclient.NetworkClient
}

func NewNetworkClientSet(ntwk cclient.NetworkClient) Network {
	n := &network{
		nc: ntwk,
	}

	return n
}

func (n *network) GetNetwork(id string) (*cclient.NetworkInfo, error) {
	net, err := n.nc.Get(id)
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	if net.SyncStatus == types.SyncStatusUnknown.String() || net.SyncStatus == types.SyncStatusFailed.String() {
		log.Error(failedRequestMsg)

		return nil, crest.HttpError{
			HttpCode: http.StatusUnprocessableEntity,
			Message:  failedRequestMsg,
		}
	}

	if net.SyncStatus == types.SyncStatusPending.String() {
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
	paymentLinks bool) (*cclient.NetworkInfo, error) {
	net, err := n.nc.Add(cclient.AddNetworkRequest{
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
		return nil, handleRestErrorStatus(err)
	}

	return net, nil
}

func handleRestErrorStatus(err error) error {
	e := cclient.ErrorStatus{}

	if errors.As(err, &e) {
		return crest.HttpError{
			HttpCode: e.StatusCode,
			Message:  err.Error(),
		}
	}

	return err
}
