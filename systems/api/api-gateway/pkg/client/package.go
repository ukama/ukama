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

	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client/rest"
	"github.com/ukama/ukama/systems/common/types"

	log "github.com/sirupsen/logrus"
	crest "github.com/ukama/ukama/systems/common/rest"
)

type Package interface {
	GetPackage(string) (*rest.PackageInfo, error)
	AddPackage(string, string, string, string, string, string, bool, bool, int64, int64, int64, string,
		string, string, string, string, uint64, float64, float64, float64, uint32, []string) (*rest.PackageInfo, error)
}

type datapackage struct {
	pc rest.PackageClient
}

func NewPackageClientSet(pkg rest.PackageClient) Package {
	p := &datapackage{
		pc: pkg,
	}

	return p
}

func (p *datapackage) GetPackage(id string) (*rest.PackageInfo, error) {
	pkg, err := p.pc.Get(id)
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	if pkg.SyncStatus == types.SyncStatusUnknown.String() || pkg.SyncStatus == types.SyncStatusFailed.String() {
		log.Error(failedRequestMsg)

		return nil, crest.HttpError{
			HttpCode: http.StatusUnprocessableEntity,
			Message:  failedRequestMsg,
		}
	}

	if pkg.SyncStatus == types.SyncStatusPending.String() {
		log.Warn(pendingRequestMsg)

		return pkg, crest.HttpError{
			HttpCode: http.StatusPartialContent,
			Message:  pendingRequestMsg,
		}
	}

	return pkg, nil
}

func (p *datapackage) AddPackage(name, orgId, ownerId, from, to, baserateId string,
	isActive, flatRate bool, smsVolume, voiceVolume, dataVolume int64, voiceUnit, dataUnit,
	simType, apn, pType string, duration uint64, markup, amount, overdraft float64, trafficPolicy uint32,
	networks []string) (*rest.PackageInfo, error) {

	pkg, err := p.pc.Add(rest.AddPackageRequest{
		Name:          name,
		OrgId:         orgId,
		OwnerId:       ownerId,
		From:          from,
		To:            to,
		BaserateId:    baserateId,
		Active:        isActive,
		SmsVolume:     smsVolume,
		VoiceVolume:   voiceVolume,
		DataVolume:    dataVolume,
		VoiceUnit:     voiceUnit,
		DataUnit:      dataUnit,
		SimType:       simType,
		Apn:           apn,
		Markup:        markup,
		Type:          pType,
		Flatrate:      flatRate,
		Amount:        amount,
		Overdraft:     overdraft,
		TrafficPolicy: trafficPolicy,
		Networks:      networks,
	})
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	return pkg, nil
}
