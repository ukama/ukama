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
	cdplan "github.com/ukama/ukama/systems/common/rest/client/dataplan"
)

type Package interface {
	GetPackage(string) (*cdplan.PackageInfo, error)
	AddPackage(string, string, string, string, string, string, bool, bool, uint64, uint64, uint64, string,
		string, string, string, string, uint64, float64, float64, float64, uint32, []string) (*cdplan.PackageInfo, error)
}

type datapackage struct {
	pc cdplan.PackageClient
}

func NewPackageClientSet(pkg cdplan.PackageClient) Package {
	p := &datapackage{
		pc: pkg,
	}

	return p
}

func (p *datapackage) GetPackage(id string) (*cdplan.PackageInfo, error) {
	pkg, err := p.pc.Get(id)
	if err != nil {
		return nil, client.HandleRestErrorStatus(err)
	}

	if pkg.SyncStatus == ukama.StatusTypeUnknown.String() || pkg.SyncStatus == ukama.StatusTypeFailed.String() {
		log.Error(failedRequestMsg)

		return nil, crest.HttpError{
			HttpCode: http.StatusUnprocessableEntity,
			Message:  failedRequestMsg,
		}
	}

	if pkg.SyncStatus == ukama.StatusTypePending.String() || pkg.SyncStatus == ukama.StatusTypeProcessing.String() {
		log.Warn(pendingRequestMsg)

		return pkg, crest.HttpError{
			HttpCode: http.StatusPartialContent,
			Message:  pendingRequestMsg,
		}
	}

	return pkg, nil
}

func (p *datapackage) AddPackage(name, orgId, ownerId, from, to, baserateId string,
	isActive, flatRate bool, smsVolume, voiceVolume, dataVolume uint64, voiceUnit, dataUnit,
	simType, apn, pType string, duration uint64, markup, amount, overdraft float64, trafficPolicy uint32,
	networks []string) (*cdplan.PackageInfo, error) {

	pkg, err := p.pc.Add(cdplan.AddPackageRequest{
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
		return nil, client.HandleRestErrorStatus(err)
	}

	return pkg, nil
}
