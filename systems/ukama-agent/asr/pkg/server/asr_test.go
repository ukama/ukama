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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/rest/client/factory"
	"github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	dp "github.com/ukama/ukama/systems/common/rest/client/dataplan"
	mocks "github.com/ukama/ukama/systems/ukama-agent/asr/mocks"
	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	pm "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/policy"
	cpb "github.com/ukama/ukama/systems/ukama-agent/cdr/pb/gen"
)

const (
	OrgId       = "40987edb-ebb6-4f84-a27c-99db7c136127"
	Org         = "ukama"
	Iccid       = "0123456789012345678912"
	Imsi        = "012345678912345"
	Atos  int64 = 7200
)

var (
	networkId = uuid.NewV4()
	sub       = db.Asr{
		Iccid:          Iccid,
		Imsi:           Imsi,
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

	Sim = factory.SimCardInfo{
		Iccid:          Iccid,
		Imsi:           Imsi,
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

	guti = db.Guti{
		Imsi:            Imsi,
		PlmnId:          "00101",
		Mmegi:           101,
		Mmec:            101,
		MTmsi:           101,
		DeviceUpdatedAt: time.Unix(int64(1639144056), 0),
	}

	tai = db.Tai{
		PlmnId:          "00101",
		Tac:             101,
		DeviceUpdatedAt: time.Unix(int64(1639144056), 0),
	}

	usage = cpb.UsageResp{
		Usage: 1024,
	}

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

	Policy = db.Policy{
		Id:           uuid.NewV4(),
		Burst:        1500,
		TotalData:    pack.DataVolume,
		ConsumedData: 0,
		Dlbr:         pack.PackageDetails.Dlbr,
		Ulbr:         pack.PackageDetails.Ulbr,
		StartTime:    1714008143,
		EndTime:      1914008143,
	}
)

func TestAsr_Read(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	gutiRepo := &mocks.GutiRepo{}

	factory := &cmocks.SimFactoryClient{}
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
		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, Atos)
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

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, Atos)
		assert.NoError(t, err)

		hs, err := s.Read(context.TODO(), &reqPb)
		assert.NoError(t, err)

		assert.NotNil(t, hs)

		asrRepo.AssertExpectations(t)
		gutiRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("ReadUlAmbrDoesNotUseDlAmbr", func(t *testing.T) {
		asymSub := sub
		asymSub.UeDlAmbrBps = 3000000
		asymSub.UeUlAmbrBps = 1000000

		reqPb := pb.ReadReq{
			Id: &pb.ReadReq_Imsi{
				Imsi: asymSub.Imsi,
			},
		}

		asrRepo.On("GetByImsi", reqPb.GetImsi()).Return(&asymSub, nil).Once()
		cdr.On("GetUsage", reqPb.GetImsi()).Return(&usage, nil).Once()

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, Atos)
		assert.NoError(t, err)

		hs, err := s.Read(context.TODO(), &reqPb)
		assert.NoError(t, err)

		assert.Equal(t, asymSub.UeDlAmbrBps, hs.Record.UeDlAmbrBps)
		assert.Equal(t, asymSub.UeUlAmbrBps, hs.Record.UeUlAmbrBps)

		asrRepo.AssertExpectations(t)
		gutiRepo.AssertExpectations(t)
	})

}

func TestAsr_UpdatePackage(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	gutiRepo := &mocks.GutiRepo{}

	factory := &cmocks.SimFactoryClient{}
	ctrl := &mocks.Controller{}
	network := &cmocks.NetworkClient{}
	mbC := &cmocks.MsgBusServiceClient{}
	cdr := &mocks.CDRService{}

	reqPb := pb.UpdatePackageReq{
		Iccid:        "0123456789012345678912",
		PackageId:    "40987edb-ebb6-4f84-a27c-99db7c136127",
		SimPackageId: "40987edb-ebb6-4f84-a27c-99db7c136300",
		TotalData:    1024000000,
		Dlbr:         15000,
		Ulbr:         2000,
		StartTime:    1700000000,
		EndTime:      1700100000,
	}

	pId, err := uuid.FromString(reqPb.PackageId)
	assert.NoError(t, err)

	spId, err := uuid.FromString(reqPb.SimPackageId)
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
	usub.Policy = Policy

	asrRepo.On("GetByIccid", reqPb.GetIccid()).Return(&sub, nil)
	ctrl.On("NewPolicy", pm.PolicyInput{
		TotalData: reqPb.TotalData,
		Dlbr:      reqPb.Dlbr,
		Ulbr:      reqPb.Ulbr,
		StartTime: reqPb.StartTime,
		EndTime:   reqPb.EndTime,
	}).Return(&Policy, nil).Once()
	asrRepo.On("UpdatePackage", sub.Imsi, pId, spId, &Policy).Return(nil).Once()
	ctrl.On("RunPolicyControl", sub.Imsi, false).Return(nil, false).Once()
	ctrl.On("SyncProfile", pcrfData, mock.Anything, msgbus.ACTION_CRUD_UPDATE, "activesubscriber", true).Return(nil, false).Once()

	s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, Atos)
	assert.NoError(t, err)

	hs, err := s.UpdatePackage(context.TODO(), &reqPb)
	assert.NoError(t, err)

	assert.NotNil(t, hs)

	asrRepo.AssertExpectations(t)
	gutiRepo.AssertExpectations(t)
	assert.NoError(t, err)

}

func TestAsr_CreateProfile(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	gutiRepo := &mocks.GutiRepo{}

	factory := &cmocks.SimFactoryClient{}
	ctrl := &mocks.Controller{}
	network := &cmocks.NetworkClient{}
	mbC := &cmocks.MsgBusServiceClient{}
	cdr := &mocks.CDRService{}

	t.Run("CreateProfileByICCID", func(t *testing.T) {
		reqPb := pb.CreateProfileReq{
			NetworkId:    networkId.String(),
			Iccid:        "0123456789012345678912",
			PackageId:    "40987edb-ebb6-4f84-a27c-99db7c136300",
			SimPackageId: "107f7b15-a8c5-4711-b1e0-f2329bffaba1",
			TotalData:    1024000000,
			Dlbr:         15000,
			Ulbr:         2000,
			StartTime:    1700000000,
			EndTime:      1700100000,
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
			Imsi:                    Sim.Imsi,
			Op:                      Sim.Op,
			Key:                     Sim.Key,
			Amf:                     Sim.Amf,
			AlgoType:                Sim.AlgoType,
			UeDlAmbrBps:             Sim.UeDlAmbrBps,
			UeUlAmbrBps:             Sim.UeUlAmbrBps,
			Sqn:                     uint64(Sim.Sqn),
			CsgIdPrsent:             Sim.CsgIdPrsent,
			CsgId:                   Sim.CsgId,
			DefaultApnName:          Sim.DefaultApnName,
			PackageId:               pId,
			SimPackageId:            spId,
			NetworkId:               nId,
			Policy:                  Policy,
			LastStatusChangeAt:      time.Now(),
			AllowedTimeOfService:    Atos,
			LastStatusChangeReasons: db.PROFILE_CREATION,
		}

		network.On("Get", reqPb.NetworkId).Return(&registry.NetworkInfo{}, nil).Once()
		factory.On("ReadSimCardInfo", reqPb.Iccid).Return(&Sim, nil).Once()
		ctrl.On("NewPolicy", pm.PolicyInput{
			TotalData: reqPb.TotalData,
			Dlbr:      reqPb.Dlbr,
			Ulbr:      reqPb.Ulbr,
			StartTime: reqPb.StartTime,
			EndTime:   reqPb.EndTime,
		}).Return(&Policy, nil).Once()
		asrRepo.On("Add", mock.MatchedBy(func(a1 *db.Asr) bool {
			return a1.Iccid == asr.Iccid
		})).Return(nil).Once()
		ctrl.On("RunPolicyControl", sub.Imsi, false).Return(nil, false).Once()
		ctrl.On("SyncProfile", pcrfData, mock.MatchedBy(func(a1 *db.Asr) bool {
			return a1.Iccid == asr.Iccid
		}), msgbus.ACTION_CRUD_CREATE, "activesubscriber", true).Return(nil, false).Once()

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, Atos)
		assert.NoError(t, err)

		hs, err := s.CreateProfile(context.TODO(), &reqPb)
		assert.NoError(t, err)

		assert.NotNil(t, hs)

		asrRepo.AssertExpectations(t)
		gutiRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("CreateProfileSucceedsDespitePcrfSyncFailure", func(t *testing.T) {
		reqPb := pb.CreateProfileReq{
			NetworkId:    networkId.String(),
			Iccid:        "0123456789012345678912",
			PackageId:    "40987edb-ebb6-4f84-a27c-99db7c136300",
			SimPackageId: "107f7b15-a8c5-4711-b1e0-f2329bffaba1",
			TotalData:    1024000000,
			Dlbr:         15000,
			Ulbr:         2000,
			StartTime:    1700000000,
			EndTime:      1700100000,
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
			Imsi:                    Sim.Imsi,
			Op:                      Sim.Op,
			Key:                     Sim.Key,
			Amf:                     Sim.Amf,
			AlgoType:                Sim.AlgoType,
			UeDlAmbrBps:             Sim.UeDlAmbrBps,
			UeUlAmbrBps:             Sim.UeUlAmbrBps,
			Sqn:                     uint64(Sim.Sqn),
			CsgIdPrsent:             Sim.CsgIdPrsent,
			CsgId:                   Sim.CsgId,
			DefaultApnName:          Sim.DefaultApnName,
			PackageId:               pId,
			SimPackageId:            spId,
			NetworkId:               nId,
			Policy:                  Policy,
			LastStatusChangeAt:      time.Now(),
			AllowedTimeOfService:    Atos,
			LastStatusChangeReasons: db.PROFILE_CREATION,
		}

		network.On("Get", reqPb.NetworkId).Return(&registry.NetworkInfo{}, nil).Once()
		factory.On("ReadSimCardInfo", reqPb.Iccid).Return(&Sim, nil).Once()
		ctrl.On("NewPolicy", pm.PolicyInput{
			TotalData: reqPb.TotalData,
			Dlbr:      reqPb.Dlbr,
			Ulbr:      reqPb.Ulbr,
			StartTime: reqPb.StartTime,
			EndTime:   reqPb.EndTime,
		}).Return(&Policy, nil).Once()
		asrRepo.On("Add", mock.MatchedBy(func(a1 *db.Asr) bool {
			return a1.Iccid == asr.Iccid
		})).Return(nil).Once()
		ctrl.On("RunPolicyControl", sub.Imsi, false).Return(nil, false).Once()
		ctrl.On("SyncProfile", pcrfData, mock.MatchedBy(func(a1 *db.Asr) bool {
			return a1.Iccid == asr.Iccid
		}), msgbus.ACTION_CRUD_CREATE, "activesubscriber", true).Return(fmt.Errorf("pcrf unreachable")).Once()

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, Atos)
		assert.NoError(t, err)

		hs, err := s.CreateProfile(context.TODO(), &reqPb)

		assert.NoError(t, err)
		assert.NotNil(t, hs)

		asrRepo.AssertExpectations(t)
		gutiRepo.AssertExpectations(t)
	})
}

func TestAsr_DeleteProfile(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	gutiRepo := &mocks.GutiRepo{}

	factory := &cmocks.SimFactoryClient{}
	ctrl := &mocks.Controller{}
	network := &cmocks.NetworkClient{}
	mbC := &cmocks.MsgBusServiceClient{}
	cdr := &mocks.CDRService{}

	t.Run("DeleteProfileByICCID", func(t *testing.T) {
		reqPb := pb.DeleteProfileReq{
			Iccid: "0123456789012345678912",
		}

		pcrfData := &pm.SimInfo{
			ID:        sub.ID,
			Imsi:      sub.Imsi,
			Iccid:     sub.Iccid,
			NetworkId: networkId,
		}

		asrRepo.On("GetByIccid", reqPb.GetIccid()).Return(&sub, nil).Once()

		asrRepo.On("Delete", sub.Imsi, db.PROFILE_DELETION).Return(nil).Once()
		ctrl.On("SyncProfile", pcrfData, mock.MatchedBy(func(a1 *db.Asr) bool {
			return a1.Iccid == sub.Iccid
		}), msgbus.ACTION_CRUD_DELETE, "activesubscriber", true).Return(nil, false).Once()

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, Atos)
		assert.NoError(t, err)

		hs, err := s.DeleteProfile(context.TODO(), &reqPb)
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

	factory := &cmocks.SimFactoryClient{}
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

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, Atos)
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

	factory := &cmocks.SimFactoryClient{}
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

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, Atos)
		assert.NoError(t, err)

		hs, err := s.UpdateTai(context.TODO(), &reqPb)
		assert.NoError(t, err)

		assert.NotNil(t, hs)

		asrRepo.AssertExpectations(t)
		gutiRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})
}

func TestAsr_UpdateAndSyncAsrProfileFromCdr(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	gutiRepo := &mocks.GutiRepo{}

	factory := &cmocks.SimFactoryClient{}
	ctrl := &mocks.Controller{}
	network := &cmocks.NetworkClient{}
	mbC := &cmocks.MsgBusServiceClient{}
	cdr := &mocks.CDRService{}

	usub := sub
	usub.Policy = Policy

	asrRepo.On("GetByImsi", usub.Imsi).Return(&usub, nil).Once()
	cdr.On("GetUsage", usub.Imsi).Return(&cpb.UsageResp{Usage: 1024}, nil).Once()
	asrRepo.On("UpdateConsumedData", usub.ID, uint64(1024)).Return(nil).Once()

	ctrl.On("RunPolicyControl", usub.Imsi, true).Return(nil, false).Once()
	ctrl.On("SyncProfile", mock.Anything, mock.Anything, msgbus.ACTION_CRUD_UPDATE, "activesubscriber", false).Return(nil).Once()

	s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, ctrl, cdr, OrgId, Org, mbC, Atos)
	assert.NoError(t, err)

	err = s.UpdateAndSyncAsrProfileFromCdr(usub.Imsi)
	assert.NoError(t, err)

	asrRepo.AssertExpectations(t)
	ctrl.AssertExpectations(t)
	cdr.AssertExpectations(t)
	asrRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}
