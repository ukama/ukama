package subscriber

import (
	"strings"
	"testing"

	b64 "encoding/base64"

	"github.com/bxcodec/faker/v4"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/ukama"
	api "github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/rest"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
)

var OrgName string
var OrgIP string
var OrgCerts string
var SysName, SysIP, SysCerts string
var SysPort int32
var NodeId ukama.NodeID
var NodeIP, NodeCerts string

var SimType, Iccid1, Iccid2, Iccid3 string

func init() {
	initializeData()
}

func initializeData() {
	OrgName = strings.ToLower(faker.FirstName()) + "_org"
	OrgIP = utils.RandomIPv4()
	log.Info(OrgIP)
	OrgCerts = utils.RandomBase64String(2048)
	SysName = strings.ToLower(faker.FirstName()) + "_sys"
	SysIP = utils.RandomIPv4()
	SysCerts = utils.RandomBase64String(2048)
	SysPort = int32(utils.RandomPort())
	NodeIP = utils.RandomIPv4()
	NodeCerts = utils.RandomBase64String(2048)
	NodeId = ukama.NewVirtualHomeNodeId()
	SimType = "test"
	Iccid1 = "8910300000003540855"
	Iccid2 = "8910300000003540845"
	Iccid3 = "8910300000003540835"

}

func TestWorkflow_1(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	log.Infof("Starting test.")
	host := "http://localhost:8071"
	subs := NewSubscriberSys(host)

	reqSimPoolUploadSimReq := api.SimPoolUploadSimReq{
		SimType: SimType,
		Data:    b64.StdEncoding.EncodeToString([]byte("https://github.com/ukama/ukama/blob/main/systems/subscriber/docs/template/SimPool.csv")),
	}

	reqSimPoolStatByTypeReq := api.SimPoolStatByTypeReq{
		SimType: SimType,
	}

	reqSimByIccidReq := api.SimByIccidReq{
		Iccid: "01234567890123456789",
	}

	reqSubscriberGetReq := api.SubscriberGetReq{
		SubscriberId: uuid.NewV4().String(),
	}

	reqSubscriberAddReq := api.SubscriberAddReq{
		FirstName:             faker.FirstName(),
		LastName:              faker.LastName(),
		Email:                 faker.Email(),
		Phone:                 faker.Phonenumber(),
		Dob:                   faker.TimeString(),
		ProofOfIdentification: "passport",
		IdSerial:              faker.UUIDDigit(),
		Address:               faker.Sentence(),
		Gender:                "male",
		OrgId:                 uuid.NewV4().String(),
		NetworkId:             uuid.NewV4().String(),
	}

	reqSubscriberDeleteReq := api.SubscriberDeleteReq{
		SubscriberId: "",
	}

	reqSubscriberUpdateReq := api.SubscriberUpdateReq{
		SubscriberId:          "",
		Email:                 faker.Email(),
		Phone:                 faker.Phonenumber(),
		ProofOfIdentification: "dl",
		IdSerial:              faker.UUIDDigit(),
		Address:               faker.Sentence(),
	}

	reqGetSimsBySubReq := api.GetSimsBySubReq{
		SubscriberId: "",
	}

	reqSimReq := api.SimReq{
		SimId: "",
	}

	reqAddPkgToSimReq := api.AddPkgToSimReq{
		SimId:     "",
		PackageId: "",
		//StartDate: "",
	}

	reqAllocateSimReq := api.AllocateSimReq{
		SubscriberId: "",
		SimToken:     "",
		PackageId:    "",
		NetworkId:    "",
		SimType:      "",
	}

	reqActivateDeactivateSimReq := api.ActivateDeactivateSimReq{
		SimId:  "",
		Status: "",
	}

	reqSetActivePackageForSimReq := api.SetActivePackageForSimReq{
		SimId:     "",
		PackageId: "",
	}

	w := utils.SetupWatcher([]string{"event.cloud.lookup.organization.create", "event.cloud.lookup.node.create", "event.cloud.lookup.system.create"})

	resp, err := subs.SubscriberSimpoolUploadSims(reqSimPoolUploadSimReq)
	assert.NoError(t, err)
	log.Infof("Expected: \n %v \n Actual: %v\n", reqSimPoolUploadSimReq, resp)
	if assert.NotNil(t, resp) {
		//assert.Equal(t, OrgName, resp.OrgName)
		//assert.Equal(t, OrgIP, utils.IPv4CIDRToStringNotation(resp.Ip))
		//assert.Equal(t, OrgCerts, resp.Certificate)
	}

	sGSSresp, err := subs.SubscriberSimpoolGetSimStats(reqSimPoolStatByTypeReq)
	assert.NoError(t, err)
	if assert.NotNil(t, sGSSresp) {

	}

	sGSBIresp, err := subs.SubscriberSimpoolGetSimByICCID(reqSimByIccidReq)
	assert.NoError(t, err)
	if assert.NotNil(t, sGSBIresp) {
		// assert.Equal(t, NodeId.String(), nresp.NodeId)
		// assert.Equal(t, OrgName, nresp.OrgName)
	}

	/* Getting system information */
	rSresp, err := subs.SubscriberRegistryAddSusbscriber(reqSubscriberAddReq)
	assert.NoError(t, err)
	if assert.NotNil(t, rSresp) {
		// assert.Equal(t, SysName, gSresp.SystemName)
		// assert.Equal(t, SysIP, gSresp.Ip)
		// assert.Equal(t, SysPort, gSresp.Port)
		// assert.Equal(t, SysCerts, gSresp.Certificate)
		// assert.NotEmpty(t, gSresp.SystemId)
	}

	/* Node bootstrapping */
	mASresp, err := subs.SubscriberManagerAllocateSim(reqAllocateSimReq)
	assert.NoError(t, err)
	if assert.NotNil(t, mASresp) {
		// assert.Equal(t, NodeId.String(), gNresp.NodeId)
		// assert.Equal(t, OrgName, gNresp.OrgName)
	}

	err = subs.SubscriberManagerAddPackage(reqAddPkgToSimReq)
	assert.NoError(t, err)

	rAPresp, err := subs.SubscriberManagerAcitvatePackage(reqSetActivePackageForSimReq)
	assert.NoError(t, err)
	if assert.NotNil(t, rAPresp) {
		// assert.Equal(t, NodeId.String(), gNresp.NodeId)
		// assert.Equal(t, OrgName, gNresp.OrgName)
	}

	rGSresp, err := subs.SubscriberManagerGetSubscriber(reqGetSimsBySubReq)
	assert.NoError(t, err)
	if assert.NotNil(t, rGSresp) {
		// assert.Equal(t, NodeId.String(), gNresp.NodeId)
		// assert.Equal(t, OrgName, gNresp.OrgName)
	}

	rDSresp, err := subs.SubscriberManagerDeleteSim(reqSimReq)
	assert.NoError(t, err)
	if assert.NotNil(t, rDSresp) {
		// assert.Equal(t, NodeId.String(), gNresp.NodeId)
		// assert.Equal(t, OrgName, gNresp.OrgName)
	}

	assert.Equal(t, true, w.Expections())
}
