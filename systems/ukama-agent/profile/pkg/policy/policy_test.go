package policy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/ukama-agent/profile/mocks"
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

func TestPolicy_DataCapCheck(t *testing.T) {
	t.Run("DataCapChecks", func(t *testing.T) {
		valid := DataCapCheck(profile)
		assert.Equal(t, true, valid)

		newProfile := profile
		newProfile.ConsumedDataBytes = 1024000
		valid = DataCapCheck(newProfile)
		assert.Equal(t, false, valid)

		newProfile.ConsumedDataBytes = 2024000
		valid = DataCapCheck(newProfile)
		assert.Equal(t, false, valid)

	})
}

func TestPolicy_AllowedTimeOfService(t *testing.T) {
	t.Run("AllowedTimeOfService", func(t *testing.T) {
		valid := AllowedTimeOfServiceCheck(profile)
		assert.Equal(t, true, valid)

		newProfile := profile
		newProfile.AllowedTimeOfService = 0
		valid = AllowedTimeOfServiceCheck(newProfile)
		assert.Equal(t, false, valid)

	})
}

func TestPolicy_RemoveProfile(t *testing.T) {
	profileRepo := &mocks.ProfileRepo{}
	mbC := &cmocks.MsgBusServiceClient{}

	t.Run("RemoveProfile", func(t *testing.T) {

		reqPb := pb.RemoveReq{
			Id: &pb.RemoveReq_Imsi{
				Imsi: Imsi,
			},
		}

		profileRepo.On("Delete", reqPb.GetImsi(), db.POLICY_FAILURE).Return(nil).Once()
		mbC.On("PublishRequest", "event.cloud.profile.policy.delete", mock.Anything).Return(nil).Once()
		mbC.On("PublishToNodeFeeder", "event.cloud.profile.policy.node-feed", "*", Org, NodePolicyPath, "DELETE", mock.Anything).Return(nil).Once()

		pc := NewPolicyController(profileRepo, Org, mbC, NodePolicyPath, MonitoringPeriod)
		assert.NotNil(t, pc)

		err, state := RemoveProfile(pc, profile)
		assert.NoError(t, err)
		assert.Equal(t, true, state)

		profileRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})
}
