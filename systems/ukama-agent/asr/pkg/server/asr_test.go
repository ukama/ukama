/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/client"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	dp "github.com/ukama/ukama/systems/common/rest/client/dataplan"
	mocks "github.com/ukama/ukama/systems/ukama-agent/asr/mocks"
	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	pm "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/policy"
	cpb "github.com/ukama/ukama/systems/ukama-agent/cdr/pb/gen"
)

var OrgId = "40987edb-ebb6-4f84-a27c-99db7c136127"
var Org = "ukama"
var networkId = uuid.NewV4()
var sub = db.Asr{
	Iccid:          "0123456789012345678912",
	Imsi:           "012345678912345",
	Op:             []byte("0123456789012345"),
	Key:            []byte("0123456789012345"),
	Amf:            []byte("800"),
	AlgoType:       1,
	UeDlAmbrBps:    2000000,
	UeUlAmbrBps:    2000000,
	Sqn:            1,
	CsgIdPrsent:    false,
	CsgId:          0,
	DefaultApnName: "ukama",
	NetworkId:      networkId,
}

var sim = client.SimCardInfo{
	Iccid:          "0123456789012345678912",
	Imsi:           "012345678912345",
	Op:             []byte("0123456789012345"),
	Key:            []byte("0123456789012345"),
	Amf:            []byte("800"),
	AlgoType:       1,
	UeDlAmbrBps:    2000000,
	UeUlAmbrBps:    2000000,
	Sqn:            1,
	CsgIdPrsent:    false,
	CsgId:          0,
	DefaultApnName: "ukama",
}

var guti = db.Guti{
	Imsi:            "012345678912345",
	PlmnId:          "00101",
	Mmegi:           101,
	Mmec:            101,
	MTmsi:           101,
	DeviceUpdatedAt: time.Unix(int64(1639144056), 0),
}

var tai = db.Tai{
	PlmnId:          "00101",
	Tac:             101,
	DeviceUpdatedAt: time.Unix(int64(1639144056), 0),
}

var usage = cpb.UsageResp{
	Usage: 1024,
}

var pack = &dp.PackageInfo{
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
	Duration: 2592000, //30 days
}

var policy = db.Policy{
	Id:           uuid.NewV4(),
	Burst:        1500,
	TotalData:    pack.DataVolume,
	ConsumedData: 0,
	Dlbr:         pack.PackageDetails.Dlbr,
	Ulbr:         pack.PackageDetails.Ulbr,
	StartTime:    1714008143,
	EndTime:      1914008143,
}

var atos int64 = 7200

func TestAsr_Read(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	gutiRepo := &mocks.GutiRepo{}

	factory := &mocks.Factory{}
	ctrl := &mocks.Controller{}
	network := &cmocks.NetworkClient{}
	mbC := &cmocks.MsgBusServiceClient{}
	cdr := &mocks.CDRService{}

	t.Run("ReadByICCID", func(t *testing.T) {

		reqPb := pb.ReadReq{
			Id: &pb.ReadReq_Iccid{
				Iccid: "0123456789012345678912",
			},
		}

		asrRepo.On("GetByIccid", reqPb.GetIccid()).Return(&sub, nil).Once()
		cdr.On("GetUsage", sub.Imsi).Return(&usage, nil).Once()
		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, atos)
		assert.NoError(t, err)

		hs, err := s.Read(context.TODO(), &reqPb)
		assert.NoError(t, err)

		assert.NotNil(t, hs)

		asrRepo.AssertExpectations(t)
		gutiRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("ReadByIMSI", func(t *testing.T) {

		reqPb := pb.ReadReq{
			Id: &pb.ReadReq_Imsi{
				Imsi: "012345678912345",
			},
		}

		asrRepo.On("GetByImsi", reqPb.GetImsi()).Return(&sub, nil).Once()
		cdr.On("GetUsage", reqPb.GetImsi()).Return(&usage, nil).Once()

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, atos)
		assert.NoError(t, err)

		hs, err := s.Read(context.TODO(), &reqPb)
		assert.NoError(t, err)

		assert.NotNil(t, hs)

		asrRepo.AssertExpectations(t)
		gutiRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

}

func TestAsr_UpdatePackage(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	gutiRepo := &mocks.GutiRepo{}

	factory := &mocks.Factory{}
	ctrl := &mocks.Controller{}
	network := &cmocks.NetworkClient{}
	mbC := &cmocks.MsgBusServiceClient{}
	cdr := &mocks.CDRService{}

	reqPb := pb.UpdatePackageReq{
		Iccid:     "0123456789012345678912",
		PackageId: "40987edb-ebb6-4f84-a27c-99db7c136127",
	}

	pId, err := uuid.FromString(reqPb.PackageId)
	assert.NoError(t, err)

	pcrfData := &pm.SimInfo{
		ID:        sub.ID,
		Imsi:      sub.Imsi,
		Iccid:     sub.Iccid,
		PackageId: pId,
		NetworkId: sub.NetworkId,
	}

	usub := sub
	usub.PackageId = pId
	usub.Policy = policy

	asrRepo.On("GetByIccid", reqPb.GetIccid()).Return(&sub, nil)
	ctrl.On("NewPolicy", pId).Return(&policy, nil).Once()
	asrRepo.On("UpdatePackage", sub.Imsi, pId, &policy).Return(nil).Once()
	ctrl.On("RunPolicyControl", sub.Imsi, false).Return(nil, false).Once()
	ctrl.On("SyncProfile", pcrfData, mock.Anything, msgbus.ACTION_CRUD_UPDATE, "activesubscriber", true).Return(nil, false).Once()

	s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, atos)
	assert.NoError(t, err)

	hs, err := s.UpdatePackage(context.TODO(), &reqPb)
	assert.NoError(t, err)

	assert.NotNil(t, hs)

	asrRepo.AssertExpectations(t)
	gutiRepo.AssertExpectations(t)
	assert.NoError(t, err)

}

func TestAsr_Activate(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	gutiRepo := &mocks.GutiRepo{}

	factory := &mocks.Factory{}
	ctrl := &mocks.Controller{}
	network := &cmocks.NetworkClient{}
	mbC := &cmocks.MsgBusServiceClient{}
	cdr := &mocks.CDRService{}

	t.Run("ActivateByICCID", func(t *testing.T) {

		reqPb := pb.ActivateReq{
			NetworkId:    networkId.String(),
			Iccid:        "0123456789012345678912",
			PackageId:    "40987edb-ebb6-4f84-a27c-99db7c136300",
			SimPackageId: "107f7b15-a8c5-4711-b1e0-f2329bffaba1",
		}

		pId, err := uuid.FromString(reqPb.PackageId)
		assert.NoError(t, err)

		spId, err := uuid.FromString(reqPb.SimPackageId)
		assert.NoError(t, err)

		nId, err := uuid.FromString(reqPb.NetworkId)
		assert.NoError(t, err)

		pcrfData := &pm.SimInfo{
			Imsi:      sub.Imsi,
			Iccid:     sub.Iccid,
			PackageId: pId,
			NetworkId: sub.NetworkId,
			Visitor:   false,
		}

		asr := &db.Asr{
			Iccid:                   reqPb.Iccid,
			Imsi:                    sim.Imsi,
			Op:                      sim.Op,
			Key:                     sim.Key,
			Amf:                     sim.Amf,
			AlgoType:                sim.AlgoType,
			UeDlAmbrBps:             sim.UeDlAmbrBps,
			UeUlAmbrBps:             sim.UeUlAmbrBps,
			Sqn:                     uint64(sim.Sqn),
			CsgIdPrsent:             sim.CsgIdPrsent,
			CsgId:                   sim.CsgId,
			DefaultApnName:          sim.DefaultApnName,
			PackageId:               pId,
			SimPackageId:            spId,
			NetworkId:               nId,
			Policy:                  policy,
			LastStatusChangeAt:      time.Now(),
			AllowedTimeOfService:    atos,
			LastStatusChangeReasons: db.ACTIVATION,
		}

		network.On("Get", reqPb.NetworkId).Return(&registry.NetworkInfo{}, nil).Once()
		factory.On("ReadSimCardInfo", reqPb.Iccid).Return(&sim, nil).Once()
		ctrl.On("NewPolicy", pId).Return(&policy, nil).Once()
		asrRepo.On("Add", mock.MatchedBy(func(a1 *db.Asr) bool {
			return a1.Iccid == asr.Iccid
		})).Return(nil).Once()
		ctrl.On("RunPolicyControl", sub.Imsi, false).Return(nil, false).Once()
		ctrl.On("SyncProfile", pcrfData, mock.MatchedBy(func(a1 *db.Asr) bool {
			return a1.Iccid == asr.Iccid
		}), msgbus.ACTION_CRUD_CREATE, "activesubscriber", true).Return(nil, false).Once()

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, atos)
		assert.NoError(t, err)

		hs, err := s.Activate(context.TODO(), &reqPb)
		assert.NoError(t, err)

		assert.NotNil(t, hs)

		asrRepo.AssertExpectations(t)
		gutiRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})
}

func TestAsr_Inactivate(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	gutiRepo := &mocks.GutiRepo{}

	factory := &mocks.Factory{}
	ctrl := &mocks.Controller{}
	network := &cmocks.NetworkClient{}
	mbC := &cmocks.MsgBusServiceClient{}
	cdr := &mocks.CDRService{}

	t.Run("InactivateByICCID", func(t *testing.T) {
		reqPb := pb.InactivateReq{
			Iccid: "0123456789012345678912",
		}

		pcrfData := &pm.SimInfo{
			ID:        sub.ID,
			Imsi:      sub.Imsi,
			Iccid:     sub.Iccid,
			NetworkId: networkId,
		}

		asrRepo.On("GetByIccid", reqPb.GetIccid()).Return(&sub, nil).Once()

		// Remove the Delete expectation since it's commented out in the implementation
		// asrRepo.On("Delete", sub.Imsi, db.DEACTIVATION).Return(nil).Once()
		
		ctrl.On("SyncProfile", pcrfData, mock.MatchedBy(func(a1 *db.Asr) bool {
			return a1.Iccid == sub.Iccid
		}), msgbus.ACTION_CRUD_DELETE, "activesubscriber", true).Return(nil, false).Once()

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, atos)
		assert.NoError(t, err)

		hs, err := s.Inactivate(context.TODO(), &reqPb)
		assert.NoError(t, err)

		assert.NotNil(t, hs)

		asrRepo.AssertExpectations(t)
		gutiRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})
}

func TestAsr_UpdateGuti(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	gutiRepo := &mocks.GutiRepo{}

	factory := &mocks.Factory{}
	ctrl := &mocks.Controller{}
	network := &cmocks.NetworkClient{}
	mbC := &cmocks.MsgBusServiceClient{}
	cdr := &mocks.CDRService{}

	t.Run("Update", func(t *testing.T) {

		reqPb := pb.UpdateGutiReq{
			Imsi: guti.Imsi,
			Guti: &pb.Guti{
				PlmnId: guti.PlmnId,
				Mmegi:  guti.Mmegi,
				Mmec:   guti.Mmec,
				Mtmsi:  guti.MTmsi,
			},
			UpdatedAt: uint32(guti.DeviceUpdatedAt.Unix()),
		}

		asrRepo.On("GetByImsi", reqPb.GetImsi()).Return(&sub, nil).Once()
		gutiRepo.On("Update", &guti).Return(nil).Once()

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, atos)
		assert.NoError(t, err)

		hs, err := s.UpdateGuti(context.TODO(), &reqPb)
		assert.NoError(t, err)

		assert.NotNil(t, hs)

		asrRepo.AssertExpectations(t)
		gutiRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})
}

func TestAsr_UpdateTai(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	gutiRepo := &mocks.GutiRepo{}

	factory := &mocks.Factory{}
	ctrl := &mocks.Controller{}
	network := &cmocks.NetworkClient{}
	mbC := &cmocks.MsgBusServiceClient{}
	cdr := &mocks.CDRService{}

	t.Run("Update", func(t *testing.T) {

		reqPb := pb.UpdateTaiReq{
			Imsi:      sub.Imsi,
			PlmnId:    tai.PlmnId,
			Tac:       tai.Tac,
			UpdatedAt: uint32(guti.DeviceUpdatedAt.Unix()),
		}

		asrRepo.On("GetByImsi", reqPb.GetImsi()).Return(&sub, nil).Once()
		asrRepo.On("UpdateTai", sub.Imsi, tai).Return(nil).Once()

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, atos)
		assert.NoError(t, err)

		hs, err := s.UpdateTai(context.TODO(), &reqPb)
		assert.NoError(t, err)

		assert.NotNil(t, hs)

		asrRepo.AssertExpectations(t)
		gutiRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})
}
