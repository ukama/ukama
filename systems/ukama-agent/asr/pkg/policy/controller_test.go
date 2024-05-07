package policy_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/msgbus"
	dp "github.com/ukama/ukama/systems/common/rest/client/dataplan"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/ukama-agent/asr/mocks"
	db "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
	ip "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/policy"
)

var Imsi = "012345678912345"
var Iccid = "123456789123456789123"
var Network = "db081ef5-a8ae-4a95-bff3-a7041d52bb9b"
var OrgId = "abdc0cec-5553-46aa-b3a8-1e31b0ef58ad"
var OrgName = "ukama"
var dataplanHost = "localhost:8080"
var Reroute = "10.10.10.10"

var Package = "fab4f98d-2e82-47e8-adb5-e516346880d8"
var NodePolicyPath = "v1/epc/pcrf"
var MonitoringPeriod time.Duration = 10 * time.Second

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

var sub = db.Asr{
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

var policy = db.Policy{
	Id:           uuid.NewV4(),
	Burst:        1500,
	TotalData:    pack.DataVolume,
	ConsumedData: 0,
	Dlbr:         pack.PackageDetails.Dlbr,
	Ulbr:         pack.PackageDetails.Ulbr,
	StartTime:    uint64(time.Now().Unix()),
	EndTime:      uint64(time.Now().Unix() + 10000000),
}

var simInfo = ip.SimInfo{
	Imsi:      sub.Imsi,
	Iccid:     sub.Iccid,
	PackageId: sub.PackageId,
	NetworkId: sub.NetworkId,
	Visitor:   false,
	ID:        1,
}

func TestController_NewPolicyController(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	mbC := &cmocks.MsgBusServiceClient{}

	pc := ip.NewPolicyController(asrRepo, mbC, dataplanHost, OrgName, OrgId, Reroute, MonitoringPeriod, false)
	assert.NotNil(t, pc)
}

func TestController_StartStopPolicyController(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	mbC := &cmocks.MsgBusServiceClient{}

	pc := ip.NewPolicyController(asrRepo, mbC, dataplanHost, OrgName, OrgId, Reroute, MonitoringPeriod, true)
	assert.NotNil(t, pc)
	lp := []db.Asr{sub}
	asrRepo.On("List").Return(lp, nil).Once()
	asrRepo.On("GetByImsi", Imsi).Return(&sub, nil).Once()

	time.Sleep(2 * time.Second)

	pc.StopPolicyControllerRoutine()

	time.Sleep(1 * time.Second)

}

func TestController_SyncProfile(t *testing.T) {

	asrRepo := &mocks.AsrRecordRepo{}
	mbC := &cmocks.MsgBusServiceClient{}

	pc := ip.NewPolicyController(asrRepo, mbC, dataplanHost, OrgName, OrgId, Reroute, MonitoringPeriod, false)
	assert.NotNil(t, pc)

	mbC.On("PublishRequest", "request.cloud.local.ukama.ukamaagent.asr.nodefeeder.publish", mock.Anything).Return(nil).Once()
	mbC.On("PublishRequest", "event.cloud.local.ukama.ukamaagent.asr.activesubscriber.create", mock.Anything).Return(nil).Once()

	err := pc.SyncProfile(&simInfo, &sub, msgbus.ACTION_CRUD_CREATE, "activesubscriber", true)
	assert.NoError(t, err)

}
