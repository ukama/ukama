/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package policy_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/ukama-agent/asr/mocks"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	dp "github.com/ukama/ukama/systems/common/rest/client/dataplan"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	db "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
	ip "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/policy"
)

const (
	Imsi    = "012345678912345"
	Iccid   = "123456789123456789123"
	Network = "db081ef5-a8ae-4a95-bff3-a7041d52bb9b"
	OrgId   = "abdc0cec-5553-46aa-b3a8-1e31b0ef58ad"
	OrgName = "ukama"
	Reroute = "10.10.10.10"

	Package        = "fab4f98d-2e82-47e8-adb5-e516346880d8"
	NodePolicyPath = "v1/epc/pcrf"
)

var (
	MonitoringPeriod time.Duration = 10 * time.Second

	pack = &dp.PackageInfo{
		Name:        "Monthly Data",
		OrgId:       uuid.NewV4().String(),
		OwnerId:     uuid.NewV4().String(),
		From:        "2023-04-01T00:00:00Z",
		To:          "2025-04-01T00:00:00Z",
		BaserateId:  uuid.NewV4().String(),
		VoiceVolume: 0,
		IsActive:    true,
		DataVolume:  1024000000,
		SmsVolume:   0,
		DataUnit:    "bytes",
		VoiceUnit:   "seconds",
		SimType:     "test",
		Apn:         "ukama.tel",
		PackageDetails: dp.PackageDetails{
			Dlbr: 15000,
			Ulbr: 2000,
			Apn:  "xyz",
		},
		Type:     "postpaid",
		Flatrate: false,
		Amount:   0,
		Duration: 2592000, //in minutes, 1800 days (which is roughly 4.93 years).
	}

	sub = db.Asr{
		Iccid:                "0123456789012345678912",
		Imsi:                 "012345678912345",
		Op:                   []byte("0123456789012345"),
		Key:                  []byte("0123456789012345"),
		Amf:                  []byte("800"),
		AlgoType:             1,
		UeDlAmbrBps:          2000000,
		UeUlAmbrBps:          2000000,
		Sqn:                  1,
		CsgIdPrsent:          false,
		CsgId:                0,
		DefaultApnName:       "ukama",
		NetworkId:            uuid.FromStringOrNil(Network),
		Policy:               policy,
		PackageId:            uuid.FromStringOrNil(Package),
		AllowedTimeOfService: 7200,
		LastStatusChangeAt:   time.Now(),
	}

	policy = db.Policy{
		Id:           uuid.NewV4(),
		Burst:        1500,
		TotalData:    pack.DataVolume,
		ConsumedData: 0,
		Dlbr:         pack.PackageDetails.Dlbr,
		Ulbr:         pack.PackageDetails.Ulbr,
		StartTime:    uint64(time.Now().Unix()),
		EndTime:      uint64(time.Now().Unix() + 10000000),
	}

	simInfo = ip.SimInfo{
		Imsi:      sub.Imsi,
		Iccid:     sub.Iccid,
		PackageId: sub.PackageId,
		NetworkId: sub.NetworkId,
		Visitor:   false,
		ID:        1,
	}
)

func TestController_NewPolicyController(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	mbC := &cmocks.MsgBusServiceClient{}

	pc := ip.NewPolicyController(asrRepo, mbC, OrgName, OrgId, Reroute, MonitoringPeriod, false)
	assert.NotNil(t, pc)
}

func TestController_StartStopPolicyController(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	mbC := &cmocks.MsgBusServiceClient{}

	pc := ip.NewPolicyController(asrRepo, mbC, OrgName, OrgId, Reroute, MonitoringPeriod, true)
	assert.NotNil(t, pc)
	lp := []db.Asr{sub}
	asrRepo.On("List").Return(lp, nil).Once()
	asrRepo.On("GetByImsi", Imsi).Return(&sub, nil).Once()

	time.Sleep(2 * time.Second)

	pc.StopPolicyControllerRoutine()

	time.Sleep(1 * time.Second)

}

func TestController_NewPolicy(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	mbC := &cmocks.MsgBusServiceClient{}

	pc := ip.NewPolicyController(asrRepo, mbC, OrgName, OrgId, Reroute, MonitoringPeriod, false)

	now := uint64(time.Now().Unix())

	p, err := pc.NewPolicy(ip.PolicyInput{
		TotalData: 1024000000,
		Dlbr:      15000,
		Ulbr:      2000,
		StartTime: now,
		EndTime:   now + 10000000,
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1024000000), p.TotalData)
	assert.Equal(t, uint64(15000), p.Dlbr)
	assert.Equal(t, uint64(2000), p.Ulbr)
	assert.Equal(t, now, p.StartTime)
	assert.Equal(t, now+10000000, p.EndTime)
	assert.Equal(t, uint64(0), p.ConsumedData)
}

func TestController_NewPolicy_RejectsZeroTotalData(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	mbC := &cmocks.MsgBusServiceClient{}

	pc := ip.NewPolicyController(asrRepo, mbC, OrgName, OrgId, Reroute, MonitoringPeriod, false)

	now := uint64(time.Now().Unix())

	_, err := pc.NewPolicy(ip.PolicyInput{
		TotalData: 0,
		StartTime: now,
		EndTime:   now + 100,
	})
	assert.Error(t, err)
}

func TestController_NewPolicy_RejectsEndBeforeStart(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	mbC := &cmocks.MsgBusServiceClient{}

	pc := ip.NewPolicyController(asrRepo, mbC, OrgName, OrgId, Reroute, MonitoringPeriod, false)

	now := uint64(time.Now().Unix())

	_, err := pc.NewPolicy(ip.PolicyInput{
		TotalData: 1000,
		StartTime: now,
		EndTime:   now - 1,
	})
	assert.Error(t, err)
}

func TestController_SyncProfile(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	mbC := &cmocks.MsgBusServiceClient{}

	pc := ip.NewPolicyController(asrRepo, mbC, OrgName, OrgId, Reroute, MonitoringPeriod, false)
	assert.NotNil(t, pc)

	mbC.On("PublishRequest", "event.cloud.local.ukama.ukamaagent.asr.policies.publish", mock.Anything).Return(nil).Once()
	mbC.On("PublishRequest", "event.cloud.local.ukama.ukamaagent.asr.activesubscriber.create", mock.Anything).Return(nil).Once()

	err := pc.SyncProfile(&simInfo, &sub, msgbus.ACTION_CRUD_CREATE, "activesubscriber", true)
	assert.NoError(t, err)
}
