package init

import (
	"strings"
	"testing"

	"github.com/bxcodec/faker/v4"
	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/ukama"
	api "github.com/ukama/ukama/systems/init/api-gateway/pkg/rest"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
)

var OrgName string
var OrgIP string
var OrgCerts string
var SysName, SysIP, SysCerts string
var SysPort int32
var NodeId ukama.NodeID
var NodeIP, NodeCerts string

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
}

func TestWorkflow_1(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	log.Infof("Starting test.")
	host := "http://localhost:8071"
	init := NewInitSys(host)

	reqAddOrg := api.AddOrgRequest{
		OrgName:     OrgName,
		Ip:          OrgIP,
		Certificate: OrgCerts,
	}

	reqAddSystem := api.AddSystemRequest{
		OrgName:     OrgName,
		SysName:     SysName,
		Ip:          SysIP,
		Certificate: SysCerts,
		Port:        int32(SysPort),
	}

	reqAddNode := api.AddNodeRequest{
		OrgName: OrgName,
		NodeId:  NodeId.String(),
	}

	reqGetNode := api.GetNodeRequest{
		OrgName: OrgName,
		NodeId:  NodeId.String(),
	}

	reqGetSystem := api.GetSystemRequest{
		OrgName: OrgName,
		SysName: SysName,
	}

	w := utils.SetupWatcher([]string{"event.cloud.lookup.organization.create", "event.cloud.lookup.node.create", "event.cloud.lookup.system.create"})
	resp, err := init.InitAddOrg(reqAddOrg)
	assert.NoError(t, err)
	log.Infof("Expected: \n %v \n Actual: %v\n", reqAddOrg, resp)
	if assert.NotNil(t, resp) {
		assert.Equal(t, OrgName, resp.OrgName)
		assert.Equal(t, OrgIP, utils.IPv4CIDRToStringNotation(resp.Ip))
		assert.Equal(t, OrgCerts, resp.Certificate)
	}

	sresp, err := init.InitAddSystem(reqAddSystem)
	assert.NoError(t, err)
	if assert.NotNil(t, sresp) {
		assert.Equal(t, SysName, sresp.SystemName)
		assert.Equal(t, SysIP, sresp.Ip)
		assert.Equal(t, SysPort, sresp.Port)
		assert.Equal(t, SysCerts, sresp.Certificate)
		assert.NotEmpty(t, sresp.SystemId)
	}

	nresp, err := init.InitAddNode(reqAddNode)
	assert.NoError(t, err)
	if assert.NotNil(t, nresp) {
		assert.Equal(t, NodeId.String(), nresp.NodeId)
		assert.Equal(t, OrgName, nresp.OrgName)
	}

	/* Getting system information */
	gSresp, err := init.InitGetSystem(reqGetSystem)
	assert.NoError(t, err)
	if assert.NotNil(t, gSresp) {
		assert.Equal(t, SysName, gSresp.SystemName)
		assert.Equal(t, SysIP, gSresp.Ip)
		assert.Equal(t, SysPort, gSresp.Port)
		assert.Equal(t, SysCerts, gSresp.Certificate)
		assert.NotEmpty(t, gSresp.SystemId)
	}

	/* Node bootstrapping */
	gNresp, err := init.InitGetNode(reqGetNode)
	assert.NoError(t, err)
	if assert.NotNil(t, gNresp) {
		assert.Equal(t, NodeId.String(), gNresp.NodeId)
		assert.Equal(t, OrgName, gNresp.OrgName)
	}

	assert.Equal(t, true, w.Expections())
}
