package policy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/ukama-agent/profile/mocks"
	pb "github.com/ukama/ukama/systems/ukama-agent/profile/pb/gen"
	db "github.com/ukama/ukama/systems/ukama-agent/profile/pkg/db"
)

func TestPolicy_DataCapCheck(t *testing.T) {
	t.Run("DataCapChecks", func(t *testing.T) {
		valid := DataCapCheck(Profile)
		assert.Equal(t, true, valid)

		newProfile := Profile
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
		valid := AllowedTimeOfServiceCheck(Profile)
		assert.Equal(t, true, valid)

		newProfile := Profile
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

		pc := NewPolicyController(profileRepo, Org, mbC, NodePolicyPath, false, MonitoringPeriod)
		assert.NotNil(t, pc)

		err, state := RemoveProfile(pc, Profile)
		assert.NoError(t, err)
		assert.Equal(t, true, state)

		profileRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})
}
