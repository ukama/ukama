package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/uuid"
	mocks "github.com/ukama/ukama/systems/ukama-agent/profile/mocks"
	pb "github.com/ukama/ukama/systems/ukama-agent/profile/pb/gen"
	db "github.com/ukama/ukama/systems/ukama-agent/profile/pkg/db"
)

var Imsi = "012345678912345"
var Iccid = "123456789123456789123"
var Network = "db081ef5-a8ae-4a95-bff3-a7041d52bb9b"
var Org = "abdc0cec-5553-46aa-b3a8-1e31b0ef58ad"
var Package = "fab4f98d-2e82-47e8-adb5-e516346880d8"
var NodePolicyPath = "v1/epc/pcrf"
var MonitoringPeriod time.Duration = 10 * time.Second
var profile = db.Profile{
	Iccid:                   Iccid,
	Imsi:                    Imsi,
	UeDlBps:                 10000000,
	UeUlBps:                 1000000,
	ApnName:                 "ukama",
	AllowedTimeOfService:    2592000,
	TotalDataBytes:          1024000,
	ConsumedDataBytes:       0,
	NetworkId:               uuid.FromStringOrNil(Network),
	PackageId:               uuid.FromStringOrNil(Package),
	LastStatusChangeReasons: db.ACTIVATION,
	LastStatusChangeAt:      time.Now(),
}

var pack = db.PackageDetails{
	PackageId:            uuid.FromStringOrNil(Package),
	UeDlBps:              10000000,
	UeUlBps:              1000000,
	ApnName:              "ukama",
	AllowedTimeOfService: time.Second * 2592000,
	TotalDataBytes:       1024000,
	ConsumedDataBytes:    0,
	LastStatusChangeAt:   time.Now(),
}

func TestProfile_Read(t *testing.T) {
	profileRepo := &mocks.ProfileRepo{}
	mbC := &cmocks.MsgBusServiceClient{}

	t.Run("ReadByImsi", func(t *testing.T) {
		reqPb := pb.ReadReq{
			Id: &pb.ReadReq_Imsi{
				Imsi: Imsi,
			},
		}

		profileRepo.On("GetByImsi", reqPb.GetImsi()).Return(&profile, nil).Once()

		s, err := NewProfileServer(profileRepo, Org, mbC, NodePolicyPath, MonitoringPeriod)
		assert.NoError(t, err)

		p, err := s.Read(context.Background(), &reqPb)
		assert.NoError(t, err)

		profileRepo.AssertExpectations(t)
		assert.NoError(t, err)

		assert.NotNil(t, p)
		assert.EqualValues(t, p.Profile.GetImsi(), reqPb.GetImsi())

	})

	t.Run("ReadByIccid", func(t *testing.T) {
		reqPb := pb.ReadReq{
			Id: &pb.ReadReq_Iccid{
				Iccid: Iccid,
			},
		}

		profileRepo.On("GetByIccid", reqPb.GetIccid()).Return(&profile, nil).Once()

		s, err := NewProfileServer(profileRepo, Org, mbC, NodePolicyPath, MonitoringPeriod)
		assert.NoError(t, err)

		p, err := s.Read(context.Background(), &reqPb)
		assert.NoError(t, err)

		profileRepo.AssertExpectations(t)
		assert.NoError(t, err)

		assert.NotNil(t, p)
		assert.EqualValues(t, p.Profile.GetIccid(), reqPb.GetIccid())

	})
}

func TestProfile_Add(t *testing.T) {
	profileRepo := &mocks.ProfileRepo{}
	mbC := &cmocks.MsgBusServiceClient{}

	t.Run("Add_Success", func(t *testing.T) {
		reqPb := pb.AddReq{
			Profile: &pb.Profile{
				Iccid:   profile.Iccid,
				Imsi:    profile.Imsi,
				UeDlBps: profile.UeDlBps,
				UeUlBps: profile.UeUlBps,
				Apn: &pb.Apn{
					Name: profile.ApnName,
				},
				NetworkId:            profile.NetworkId.String(),
				PackageId:            profile.PackageId.String(),
				AllowedTimeOfService: profile.AllowedTimeOfService,
				TotalDataBytes:       profile.TotalDataBytes,
				ConsumedDataBytes:    profile.ConsumedDataBytes,
				LastChange:           db.ACTIVATION.String(),
				LastChangeAt:         profile.LastStatusChangeAt.Unix(),
			},
		}

		profileRepo.On("Add", mock.AnythingOfType("*db.Profile")).Return(nil).Once()
		profileRepo.On("GetByImsi", reqPb.Profile.GetImsi()).Return(&profile, nil).Once()
		mbC.On("PublishRequest", "event.cloud.profile.profile.create", mock.Anything).Return(nil).Once()
		mbC.On("PublishToNodeFeeder", "event.cloud.profile.server.node-feed", "*", Org, NodePolicyPath, "PUT", mock.Anything).Return(nil).Once()

		s, err := NewProfileServer(profileRepo, Org, mbC, NodePolicyPath, MonitoringPeriod)
		assert.NoError(t, err)

		_, err = s.Add(context.Background(), &reqPb)
		assert.NoError(t, err)

		profileRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})
}

func TestProfile_UpdateUsage(t *testing.T) {
	profileRepo := &mocks.ProfileRepo{}
	mbC := &cmocks.MsgBusServiceClient{}

	t.Run("UpdateUsage_Success", func(t *testing.T) {
		var usage uint64 = 1000
		reqPb := pb.UpdateUsageReq{
			Imsi:              profile.Imsi,
			ConsumedDataBytes: usage,
		}

		profileRepo.On("GetByImsi", reqPb.GetImsi()).Return(&profile, nil).Once()
		profileRepo.On("UpdateUsage", reqPb.GetImsi(), reqPb.GetConsumedDataBytes()).Return(nil).Once()
		profileRepo.On("GetByImsi", reqPb.GetImsi()).Return(&profile, nil).Once()

		s, err := NewProfileServer(profileRepo, Org, mbC, NodePolicyPath, MonitoringPeriod)
		assert.NoError(t, err)

		_, err = s.UpdateUsage(context.Background(), &reqPb)
		assert.NoError(t, err)

		profileRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("UpdateUsageWithDataBytesOver_Success", func(t *testing.T) {
		var usage uint64 = 1024000

		reqPb := pb.UpdateUsageReq{
			Imsi:              profile.Imsi,
			ConsumedDataBytes: usage,
		}

		updatedProfile := profile
		updatedProfile.ConsumedDataBytes = usage
		profileRepo.On("GetByImsi", reqPb.GetImsi()).Return(&profile, nil).Once()
		profileRepo.On("UpdateUsage", reqPb.GetImsi(), reqPb.GetConsumedDataBytes()).Return(nil).Once()
		profileRepo.On("GetByImsi", reqPb.GetImsi()).Return(&updatedProfile, nil).Once()
		profileRepo.On("Delete", reqPb.GetImsi(), db.POLICY_FAILURE).Return(nil).Once()
		mbC.On("PublishRequest", "event.cloud.profile.policy.delete", mock.Anything).Return(nil).Once()
		mbC.On("PublishToNodeFeeder", "event.cloud.profile.policy.node-feed", "*", Org, NodePolicyPath, "DELETE", mock.Anything).Return(nil).Once()

		s, err := NewProfileServer(profileRepo, Org, mbC, NodePolicyPath, MonitoringPeriod)
		assert.NoError(t, err)

		_, err = s.UpdateUsage(context.Background(), &reqPb)
		assert.NoError(t, err)

		profileRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("UpdateUsageWithAllowedServiceTimeOver_Success", func(t *testing.T) {
		var usage uint64 = 1024000

		reqPb := pb.UpdateUsageReq{
			Imsi:              profile.Imsi,
			ConsumedDataBytes: usage,
		}

		updatedProfile := profile
		updatedProfile.ConsumedDataBytes = usage
		updatedProfile.AllowedTimeOfService = 0
		profileRepo.On("GetByImsi", reqPb.GetImsi()).Return(&profile, nil).Once()
		profileRepo.On("UpdateUsage", reqPb.GetImsi(), reqPb.GetConsumedDataBytes()).Return(nil).Once()
		profileRepo.On("GetByImsi", reqPb.GetImsi()).Return(&updatedProfile, nil).Once()
		profileRepo.On("Delete", reqPb.GetImsi(), db.POLICY_FAILURE).Return(nil).Once()
		mbC.On("PublishRequest", "event.cloud.profile.policy.delete", mock.Anything).Return(nil).Once()
		mbC.On("PublishToNodeFeeder", "event.cloud.profile.policy.node-feed", "*", Org, NodePolicyPath, "DELETE", mock.Anything).Return(nil).Once()

		s, err := NewProfileServer(profileRepo, Org, mbC, NodePolicyPath, MonitoringPeriod)
		assert.NoError(t, err)

		_, err = s.UpdateUsage(context.Background(), &reqPb)
		assert.NoError(t, err)

		profileRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

}

func TestProfile_Remove(t *testing.T) {
	profileRepo := &mocks.ProfileRepo{}
	mbC := &cmocks.MsgBusServiceClient{}

	t.Run("RemoveByImsi_Success", func(t *testing.T) {
		reqPb := pb.RemoveReq{
			Id: &pb.RemoveReq_Imsi{
				Imsi: Imsi,
			},
		}

		profileRepo.On("GetByImsi", reqPb.GetImsi()).Return(&profile, nil).Once()
		profileRepo.On("Delete", reqPb.GetImsi(), db.DEACTIVATION).Return(nil).Once()
		mbC.On("PublishRequest", "event.cloud.profile.profile.delete", mock.Anything).Return(nil).Once()
		mbC.On("PublishToNodeFeeder", "event.cloud.profile.server.node-feed", "*", Org, NodePolicyPath, "DELETE", mock.Anything).Return(nil).Once()

		s, err := NewProfileServer(profileRepo, Org, mbC, NodePolicyPath, MonitoringPeriod)
		assert.NoError(t, err)

		_, err = s.Remove(context.Background(), &reqPb)
		assert.NoError(t, err)

		profileRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("RemoveByIccid_Success", func(t *testing.T) {
		reqPb := pb.RemoveReq{
			Id: &pb.RemoveReq_Iccid{
				Iccid: Iccid,
			},
		}

		profileRepo.On("GetByIccid", reqPb.GetIccid()).Return(&profile, nil).Once()
		profileRepo.On("Delete", profile.Imsi, db.DEACTIVATION).Return(nil).Once()
		mbC.On("PublishRequest", "event.cloud.profile.profile.delete", mock.Anything).Return(nil).Once()
		mbC.On("PublishToNodeFeeder", "event.cloud.profile.server.node-feed", "*", Org, NodePolicyPath, "DELETE", mock.Anything).Return(nil).Once()

		s, err := NewProfileServer(profileRepo, Org, mbC, NodePolicyPath, MonitoringPeriod)
		assert.NoError(t, err)

		_, err = s.Remove(context.Background(), &reqPb)
		assert.NoError(t, err)

		profileRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

}

func TestProfile_UpdatePackage(t *testing.T) {
	profileRepo := &mocks.ProfileRepo{}
	mbC := &cmocks.MsgBusServiceClient{}
	t.Run("UpdatePackage_Success", func(t *testing.T) {
		reqPb := pb.UpdatePackageReq{
			Iccid: Iccid,
			Package: &pb.Package{
				UeDlBps: profile.UeDlBps,
				UeUlBps: profile.UeUlBps,
				Apn: &pb.Apn{
					Name: profile.ApnName,
				},
				PackageId:            profile.PackageId.String(),
				AllowedTimeOfService: profile.AllowedTimeOfService,
				TotalDataBytes:       profile.TotalDataBytes,
				ConsumedDataBytes:    profile.ConsumedDataBytes,
			},
		}

		profileRepo.On("GetByIccid", reqPb.GetIccid()).Return(&profile, nil)
		profileRepo.On("UpdatePackage", profile.Imsi, mock.AnythingOfTypeArgument("db.PackageDetails")).Return(nil).Once()
		profileRepo.On("GetByImsi", profile.Imsi).Return(&profile, nil).Once()
		mbC.On("PublishRequest", "event.cloud.profile.profile.update", mock.Anything).Return(nil).Once()
		mbC.On("PublishToNodeFeeder", "event.cloud.profile.server.node-feed", "*", Org, NodePolicyPath, "PUT", mock.Anything).Return(nil).Once()

		s, err := NewProfileServer(profileRepo, Org, mbC, NodePolicyPath, MonitoringPeriod)
		assert.NoError(t, err)

		_, err = s.UpdatePackage(context.Background(), &reqPb)
		assert.NoError(t, err)

		profileRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})
}

func TestProfile_Sync(t *testing.T) {
	profileRepo := &mocks.ProfileRepo{}
	mbC := &cmocks.MsgBusServiceClient{}
	t.Run("Sync", func(t *testing.T) {
		reqPb := pb.SyncReq{
			Iccid: []string{Iccid},
		}

		profileRepo.On("GetByIccid", reqPb.Iccid[0]).Return(&profile, nil).Once()
		profileRepo.On("GetByImsi", profile.Imsi).Return(&profile, nil).Once()
		mbC.On("PublishToNodeFeeder", "event.cloud.profile.server.node-feed", "*", Org, NodePolicyPath, "PUT", mock.Anything).Return(nil).Once()

		s, err := NewProfileServer(profileRepo, Org, mbC, NodePolicyPath, MonitoringPeriod)
		assert.NoError(t, err)

		_, err = s.Sync(context.Background(), &reqPb)
		assert.NoError(t, err)

		profileRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})
}
