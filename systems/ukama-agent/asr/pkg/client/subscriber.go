/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

type Subscriber interface {
	GetSimDetails(iccid string, orgId string)
}

type subscriber struct {
	R *rest.RestClient
}

func NewSubscriberClient(url string, debug bool) (*subscriber, error) {

	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		logrus.Errorf("Can't conncet to %s url.Error %s", url, err.Error())
		return nil, err
	}

	S := &subscriber{
		R: f,
	}

	return S, nil
}

func (N *network) GetSimDetails(iccid string, orgId string) error {

	errStatus := &ErrorMessage{}

	network := NetworkInfo{}

	resp, err := N.R.C.R().
		SetError(errStatus).
		Get(N.R.URL.String() + "/v1/sim/" + iccid)

	if err != nil {
		logrus.Errorf("Failed to send api request to susbcriber system. Error %s", err.Error())
		return err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to fetch sim info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)
		return fmt.Errorf(" Sim Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &network)
	if err != nil {
		logrus.Tracef("Failed to desrialize network info. Error message is %s", err.Error())
		return fmt.Errorf("network info deserailization failure: %s", err)
	} else {
		logrus.Infof("Network Info: %+v", network)
	}

	if orgId != network.OrgId {
		logrus.Error("Missing network.")
		return fmt.Errorf("Network mismatch")
	}

	return nil
}
