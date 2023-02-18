package policy

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/ukama-agent/profile/mocks"
	db "github.com/ukama/ukama/systems/ukama-agent/profile/pkg/db"
)

var Imsi = "012345678912345"
var Iccid = "123456789123456789123"
var Network = "db081ef5-a8ae-4a95-bff3-a7041d52bb9b"
var Org = "abdc0cec-5553-46aa-b3a8-1e31b0ef58ad"
var Package = "fab4f98d-2e82-47e8-adb5-e516346880d8"
var NodePolicyPath = "v1/epc/pcrf"
var MonitoringPeriod time.Duration = 10 * time.Second

var Profile = db.Profile{
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

func TestController_NewPolicyController(t *testing.T) {
	profileRepo := &mocks.ProfileRepo{}
	mbC := &cmocks.MsgBusServiceClient{}
	pc := NewPolicyController(profileRepo, Org, mbC, NodePolicyPath, false, MonitoringPeriod)
	assert.NotNil(t, pc)
}

func TestController_StartStopPolicyController(t *testing.T) {
	lp := []db.Profile{Profile}
	profileRepo := &mocks.ProfileRepo{}
	mbC := &cmocks.MsgBusServiceClient{}
	pc := NewPolicyController(profileRepo, Org, mbC, NodePolicyPath, true, MonitoringPeriod)
	assert.NotNil(t, pc)
	profileRepo.On("List").Return(lp, nil).Once()
	profileRepo.On("GetByImsi", Imsi).Return(&Profile, nil).Once()

	time.Sleep(2 * time.Second)

	pc.StopPolicyRoutine()

	time.Sleep(1 * time.Second)

}

func TestController_syncProfile(t *testing.T) {

	profileRepo := &mocks.ProfileRepo{}
	mbC := &cmocks.MsgBusServiceClient{}
	pc := NewPolicyController(profileRepo, Org, mbC, NodePolicyPath, false, MonitoringPeriod)
	assert.NotNil(t, pc)

	mbC.On("PublishToNodeFeeder", "event.cloud.profile.policy.node-feed", "*", Org, NodePolicyPath, "DELETE", mock.Anything).Return(nil).Once()

	err := pc.syncProfile(http.MethodDelete, Profile)
	assert.NoError(t, err)

}

func TestController_publishEvent(t *testing.T) {

	profileRepo := &mocks.ProfileRepo{}
	mbC := &cmocks.MsgBusServiceClient{}
	pc := NewPolicyController(profileRepo, Org, mbC, NodePolicyPath, false, MonitoringPeriod)
	assert.NotNil(t, pc)

	mbC.On("PublishRequest", "event.cloud.profile.policy.delete", mock.Anything).Return(nil).Once()

	e := &epb.ProfileRemoved{
		Profile: &epb.Profile{
			Imsi:                 Profile.Imsi,
			Iccid:                Profile.Iccid,
			Network:              Profile.NetworkId.String(),
			Package:              Profile.PackageId.String(),
			Org:                  Org,
			AllowedTimeOfService: Profile.AllowedTimeOfService,
			TotalDataBytes:       Profile.TotalDataBytes,
		},
	}

	err := pc.publishEvent(msgbus.ACTION_CRUD_DELETE, "policy", e)
	assert.NoError(t, err)

}
