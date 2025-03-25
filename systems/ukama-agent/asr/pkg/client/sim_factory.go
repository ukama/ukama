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

	"github.com/ukama/ukama/systems/common/rest"

	log "github.com/sirupsen/logrus"
)

type Factory interface {
	ReadSimCardInfo(Iccid string) (*SimCardInfo, error)
}

type factory struct {
	R *rest.RestClient
}

type ErrorMessage struct {
	Message string `json:"error"`
}

func NewFactoryClient(url string, debug bool) (*factory, error) {
	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		log.Errorf("Can't conncet to %s url.Error %s", url, err.Error())
		return nil, err
	}

	F := &factory{
		R: f,
	}

	return F, nil
}

func (f *factory) ReadSimCardInfo(Iccid string) (*SimCardInfo, error) {
	card := Sim{}
	errStatus := &ErrorMessage{}

	resp, err := f.R.C.R().
		SetError(errStatus).
		Get(f.R.URL.String() + "/v1/factory/simcards/" + Iccid)

	if err != nil {
		log.Errorf("Failed to send api request to Factory. Error %s", err.Error())
		return nil, err
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch sim card info.HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)
		return nil, fmt.Errorf("simcard request failure: %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &card)
	if err != nil {
		log.Tracef("Failed to desrialize sim card info. Error message is %s", err.Error())
		return nil, fmt.Errorf("simcard info deserailization failure: %s", err)
	}

	log.Infof("Sim card info: %+v", card.SimCardInfo)

	return card.SimCardInfo, nil
}
