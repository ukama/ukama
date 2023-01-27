package server

import (
	"context"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	mocks "github.com/ukama/ukama/systems/ukama-agent/asr/mocks"
	pb "github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/client"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var Org = "40987edb-ebb6-4f84-a27c-99db7c136127"

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

func TestAsr_Read(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	gutiRepo := &mocks.GutiRepo{}
	factory := &mocks.Factory{}
	pcrf := &mocks.PolicyControl{}
	network := &mocks.Network{}
	mbC := &cmocks.MsgBusServiceClient{}

	t.Run("ReadByICCID", func(t *testing.T) {

		reqPb := pb.ReadReq{
			Id: &pb.ReadReq_Iccid{
				Iccid: "0123456789012345678912",
			},
		}

		asrRepo.On("GetByIccid", reqPb.GetIccid()).Return(&sub, nil).Once()

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, pcrf, Org, mbC)
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

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, pcrf, Org, mbC)
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
	pcrf := &mocks.PolicyControl{}
	network := &mocks.Network{}
	mbC := &cmocks.MsgBusServiceClient{}

	reqPb := pb.UpdatePackageReq{
		Iccid:     "0123456789012345678912",
		PackageId: "40987edb-ebb6-4f84-a27c-99db7c136127",
	}

	pId, err := uuid.FromString(reqPb.PackageId)
	assert.NoError(t, err)

	req := client.PolicyControlSimPackageUpdate{
		Imsi:      sub.Imsi,
		PackageId: pId,
	}

	asrRepo.On("GetByIccid", reqPb.GetIccid()).Return(&sub, nil).Once()
	pcrf.On("UpdateSim", req).Return(nil).Once()
	asrRepo.On("UpdatePackage", sub.Imsi, pId).Return(nil).Once()
	mbC.On("PublishRequest", "event.cloud.asr.activesubscriber.update", mock.Anything).Return(nil).Once()

	s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, pcrf, Org, mbC)
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
	pcrf := &mocks.PolicyControl{}
	network := &mocks.Network{}
	mbC := &cmocks.MsgBusServiceClient{}

	t.Run("ActivateByICCID", func(t *testing.T) {

		reqPb := pb.ActivateReq{
			Network:   "40987edb-ebb6-4f84-a27c-99db7c136127",
			Iccid:     "0123456789012345678912",
			PackageId: "40987edb-ebb6-4f84-a27c-99db7c136300",
		}

		pId, err := uuid.FromString(reqPb.PackageId)
		assert.NoError(t, err)

		nId, err := uuid.FromString(reqPb.Network)
		assert.NoError(t, err)

		pReq := client.PolicyControlSimInfo{
			Imsi:      sim.Imsi,
			Iccid:     sim.Iccid,
			PackageId: pId,
			NetworkId: nId,
			Visitor:   false, // We will using this flag on roaming in VLR
		}

		asr := &db.Asr{
			Iccid:          reqPb.Iccid,
			Imsi:           sim.Imsi,
			Op:             sim.Op,
			Key:            sim.Key,
			Amf:            sim.Amf,
			AlgoType:       sim.AlgoType,
			UeDlAmbrBps:    sim.UeDlAmbrBps,
			UeUlAmbrBps:    sim.UeUlAmbrBps,
			Sqn:            uint64(sim.Sqn),
			CsgIdPrsent:    sim.CsgIdPrsent,
			CsgId:          sim.CsgId,
			DefaultApnName: sim.DefaultApnName,
			PackageId:      pId,
			NetworkID:      nId,
		}

		network.On("ValidateNetwork", reqPb.Network, Org).Return(nil).Once()
		factory.On("ReadSimCardInfo", reqPb.Iccid).Return(&sim, nil).Once()
		pcrf.On("AddSim", pReq).Return(nil).Once()
		asrRepo.On("Add", asr).Return(nil).Once()
		mbC.On("PublishRequest", "event.cloud.asr.activesubscriber.create", mock.Anything).Return(nil).Once()

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, pcrf, Org, mbC)
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
	pcrf := &mocks.PolicyControl{}
	network := &mocks.Network{}
	mbC := &cmocks.MsgBusServiceClient{}

	t.Run("InactivateByICCID", func(t *testing.T) {

		reqPb := pb.InactivateReq{
			Id: &pb.InactivateReq_Iccid{
				Iccid: "0123456789012345678912",
			},
		}

		asrRepo.On("GetByIccid", reqPb.GetIccid()).Return(&sub, nil).Once()
		pcrf.On("DeleteSim", sub.Imsi).Return(nil).Once()
		asrRepo.On("Delete", sub.Imsi).Return(nil).Once()
		mbC.On("PublishRequest", "event.cloud.asr.activesubscriber.delete", mock.Anything).Return(nil).Once()

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, pcrf, Org, mbC)
		assert.NoError(t, err)

		hs, err := s.Inactivate(context.TODO(), &reqPb)
		assert.NoError(t, err)

		assert.NotNil(t, hs)

		asrRepo.AssertExpectations(t)
		gutiRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("InactivateByImsi", func(t *testing.T) {

		reqPb := pb.InactivateReq{
			Id: &pb.InactivateReq_Imsi{
				Imsi: "012345678912345",
			},
		}

		asrRepo.On("GetByImsi", reqPb.GetImsi()).Return(&sub, nil).Once()
		pcrf.On("DeleteSim", sub.Imsi).Return(nil).Once()
		asrRepo.On("Delete", sub.Imsi).Return(nil).Once()
		mbC.On("PublishRequest", "event.cloud.asr.activesubscriber.delete", mock.Anything).Return(nil).Once()

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, pcrf, Org, mbC)
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
	pcrf := &mocks.PolicyControl{}
	network := &mocks.Network{}
	mbC := &cmocks.MsgBusServiceClient{}
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

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, pcrf, Org, mbC)
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
	pcrf := &mocks.PolicyControl{}
	network := &mocks.Network{}
	mbC := &cmocks.MsgBusServiceClient{}

	t.Run("Update", func(t *testing.T) {

		reqPb := pb.UpdateTaiReq{
			Imsi:      sub.Imsi,
			PlmnId:    tai.PlmnId,
			Tac:       tai.Tac,
			UpdatedAt: uint32(guti.DeviceUpdatedAt.Unix()),
		}

		asrRepo.On("GetByImsi", reqPb.GetImsi()).Return(&sub, nil).Once()
		asrRepo.On("UpdateTai", sub.Imsi, tai).Return(nil).Once()

		s, err := NewAsrRecordServer(asrRepo, gutiRepo, factory, network, pcrf, Org, mbC)
		assert.NoError(t, err)

		hs, err := s.UpdateTai(context.TODO(), &reqPb)
		assert.NoError(t, err)

		assert.NotNil(t, hs)

		asrRepo.AssertExpectations(t)
		gutiRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})
}
