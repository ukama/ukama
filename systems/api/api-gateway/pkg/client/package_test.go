/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	crest "github.com/ukama/ukama/systems/common/rest"
	cclient "github.com/ukama/ukama/systems/common/rest/client"
	cdplan "github.com/ukama/ukama/systems/common/rest/client/dataplan"
)

func TestCient_GetPackage(t *testing.T) {
	packageClient := &mocks.PackageClient{}

	packageId := uuid.NewV4()
	pkgName := "Monthly Data"

	p := client.NewPackageClientSet(packageClient)

	t.Run("PackageFoundAndStatusCompleted", func(t *testing.T) {
		packageClient.On("Get", packageId.String()).
			Return(&cdplan.PackageInfo{
				Id:         packageId.String(),
				Name:       pkgName,
				SyncStatus: ukama.StatusTypeCompleted.String(),
			}, nil).Once()

		pkgInfo, err := p.GetPackage(packageId.String())

		assert.NoError(t, err)

		assert.NotNil(t, pkgInfo)
		assert.Equal(t, pkgInfo.Id, packageId.String())
		assert.Equal(t, pkgInfo.Name, pkgName)
	})

	t.Run("PackageFoundAndStatusPending", func(t *testing.T) {
		packageClient.On("Get", packageId.String()).
			Return(&cdplan.PackageInfo{
				Id:         packageId.String(),
				Name:       pkgName,
				SyncStatus: ukama.StatusTypePending.String(),
			}, nil).Once()

		pkgInfo, err := p.GetPackage(packageId.String())

		assert.Error(t, err)
		assert.IsType(t, err, crest.HttpError{})
		assert.Contains(t, err.Error(), "partial")

		assert.NotNil(t, pkgInfo)
		assert.Equal(t, pkgInfo.Id, packageId.String())
		assert.Equal(t, pkgInfo.Name, pkgName)
	})

	t.Run("PackageFoundAndStatusFailed", func(t *testing.T) {
		packageClient.On("Get", packageId.String()).
			Return(&cdplan.PackageInfo{
				Id:         packageId.String(),
				Name:       pkgName,
				SyncStatus: ukama.StatusTypeFailed.String(),
			}, nil).Once()

		pkgInfo, err := p.GetPackage(packageId.String())

		assert.Error(t, err)
		assert.IsType(t, err, crest.HttpError{})
		assert.Contains(t, err.Error(), "invalid")

		assert.Nil(t, pkgInfo)
	})

	t.Run("PackageNotFound", func(t *testing.T) {
		packageClient.On("Get", packageId.String()).
			Return(nil, fmt.Errorf("GetNetwork failure: %w",
				&cclient.ErrorStatus{
					StatusCode: 404,
					RawError:   crest.ErrorResponse{Err: "not found"},
				})).Once()

		pkgInfo, err := p.GetPackage(packageId.String())

		assert.Error(t, err)
		assert.IsType(t, err, crest.HttpError{})
		assert.Contains(t, err.Error(), "404")

		assert.Nil(t, pkgInfo)
	})

	t.Run("PackageGetError", func(t *testing.T) {
		packageClient.On("Get", packageId.String()).
			Return(nil,
				fmt.Errorf("Some unexpected error")).Once()

		pkgInfo, err := p.GetPackage(packageId.String())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")

		assert.Nil(t, pkgInfo)
	})
}

func TestCient_AddPackage(t *testing.T) {
	pkgClient := &mocks.PackageClient{}

	pkgId := uuid.NewV4()
	pkgName := "Monthly Data"
	orgId := uuid.NewV4().String()
	ownerId := uuid.NewV4().String()
	from := "2023-04-01T00:00:00Z"
	to := "2023-04-01T00:00:00Z"
	baserateId := uuid.NewV4().String()
	voiceVolume := uint64(0)
	isActive := true
	dataVolume := uint64(1024)
	smsVolume := uint64(0)
	dataUnit := "MegaBytes"
	voiceUnit := "seconds"
	simType := "test"
	apn := "ukama.tel"
	markupValue := float64(0)
	markup := cdplan.PackageMarkup{
		PackageID:  pkgId.String(),
		BaseRateId: baserateId,
		Markup:     markupValue,
	}
	pType := "postpaid"
	duration := uint64(0)
	flatRate := false
	amount := float64(0)
	overdraft := float64(0)
	trafficPolicy := uint32(0)
	networks := []string{""}

	p := client.NewPackageClientSet(pkgClient)

	t.Run("PackageCreatedAndStatusUpdated", func(t *testing.T) {
		pkgClient.On("Add", cdplan.AddPackageRequest{
			Name:          pkgName,
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
			Markup:        markupValue,
			Type:          pType,
			Duration:      duration,
			Flatrate:      flatRate,
			Amount:        amount,
			Overdraft:     overdraft,
			TrafficPolicy: trafficPolicy,
			Networks:      networks,
		}).Return(&cdplan.PackageInfo{
			Id:            pkgId.String(),
			Name:          pkgName,
			OrgId:         orgId,
			OwnerId:       ownerId,
			From:          from,
			To:            to,
			BaserateId:    baserateId,
			IsActive:      isActive,
			SmsVolume:     smsVolume,
			VoiceVolume:   voiceVolume,
			DataVolume:    dataVolume,
			VoiceUnit:     voiceUnit,
			DataUnit:      dataUnit,
			SimType:       simType,
			Apn:           apn,
			Markup:        markup,
			Type:          pType,
			Duration:      duration,
			Flatrate:      flatRate,
			Amount:        amount,
			Overdraft:     overdraft,
			TrafficPolicy: trafficPolicy,
			Networks:      networks,
		}, nil).Once()

		pkgInfo, err := p.AddPackage(pkgName, orgId, ownerId, from, to, baserateId,
			isActive, flatRate, smsVolume, voiceVolume, dataVolume, voiceUnit, dataUnit,
			simType, apn, pType, duration, markupValue, amount, overdraft, trafficPolicy, networks)

		assert.NoError(t, err)

		assert.Equal(t, pkgInfo.Id, pkgId.String())
		assert.Equal(t, pkgInfo.Name, pkgName)
	})

	t.Run("PackageNotCreated", func(t *testing.T) {
		pkgClient.On("Add", cdplan.AddPackageRequest{
			Name:          pkgName,
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
			Markup:        markupValue,
			Type:          pType,
			Duration:      duration,
			Flatrate:      flatRate,
			Amount:        amount,
			Overdraft:     overdraft,
			TrafficPolicy: trafficPolicy,
			Networks:      networks,
		}).Return(nil, errors.New("some error")).Once()

		pkgInfo, err := p.AddPackage(pkgName, orgId, ownerId, from, to, baserateId,
			isActive, flatRate, smsVolume, voiceVolume, dataVolume, voiceUnit, dataUnit,
			simType, apn, pType, duration, markupValue, amount, overdraft, trafficPolicy, networks)

		assert.Contains(t, err.Error(), "error")
		assert.Nil(t, pkgInfo)
	})
}
