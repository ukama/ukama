package init

import (
	"fmt"
	"strings"
	"testing"

	"context"

	"github.com/bxcodec/faker/v4"
	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/ukama"
	api "github.com/ukama/ukama/systems/init/api-gateway/pkg/rest"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"github.com/ukama/ukama/testing/integration/pkg/test"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
)

type InitData struct {
	OrgName                  string
	OrgIP                    string
	OrgCerts                 string
	SysName, SysIP, SysCerts string
	SysPort                  int32
	NodeId                   ukama.NodeID
	NodeIP, NodeCerts        string
	Init                     *InitSys
	Host                     string
	ROrgIp                   string

	/* API requests */
	reqAddOrg    api.AddOrgRequest
	reqAddSystem api.AddSystemRequest
	reqAddNode   api.AddNodeRequest
	reqGetNode   api.GetNodeRequest
	reqGetSystem api.GetSystemRequest

	/* API responses */

}

func InitializeData() *InitData {
	d := &InitData{}
	d.OrgName = strings.ToLower(faker.FirstName()) + "_org"
	d.OrgIP = utils.RandomIPv4()
	d.OrgCerts = utils.RandomBase64String(2048)
	d.SysName = strings.ToLower(faker.FirstName()) + "_sys"
	d.SysIP = utils.RandomIPv4()
	d.SysCerts = utils.RandomBase64String(2048)
	d.SysPort = int32(utils.RandomPort())
	d.NodeIP = utils.RandomIPv4()
	d.NodeCerts = utils.RandomBase64String(2048)
	d.NodeId = ukama.NewVirtualHomeNodeId()
	d.Host = "http://localhost:8071"
	d.Init = NewInitSys(d.Host)

	d.reqAddOrg = api.AddOrgRequest{
		OrgName:     d.OrgName,
		Ip:          d.OrgIP,
		Certificate: d.OrgCerts,
	}

	d.reqAddSystem = api.AddSystemRequest{
		OrgName:     d.OrgName,
		SysName:     d.SysName,
		Ip:          d.SysIP,
		Certificate: d.SysCerts,
		Port:        int32(d.SysPort),
	}

	d.reqAddNode = api.AddNodeRequest{
		OrgName: d.OrgName,
		NodeId:  d.NodeId.String(),
	}

	d.reqGetNode = api.GetNodeRequest{
		OrgName: d.OrgName,
		NodeId:  d.NodeId.String(),
	}

	d.reqGetSystem = api.GetSystemRequest{
		OrgName: d.OrgName,
		SysName: d.SysName,
	}

	return d
}

func TestWorkflow_InitSystem(t *testing.T) {

	w := test.NewWorkflow("init_workflow_1", "Adding a system and getting its credentials")

	w.SetUpFxn = func(ctx context.Context, w *test.Workflow) error {
		log.SetLevel(log.DebugLevel)
		log.Infof("Initilizing Data for %s.", w.String())
		w.Data = InitializeData()

		log.Debugf("Workflow Data : %+v", w.Data)
		return nil
	}

	w.RegisterTestCase(&test.TestCase{

		Name:        "Add Organization.",
		Description: "Add organization to Init System",
		Data:        &pb.AddOrgResponse{},
		Workflow:    w,
		SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			log.Debugf("Setting up watcher for %s", tc.String())
			tc.Watcher = utils.SetupWatcher([]string{"event.cloud.lookup.organization.create"})
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*InitData)
			if ok {
				tc.Data, err = a.Init.InitAddOrg(a.reqAddOrg)
				log.Debugf("Expected: \n %v \n Actual:  type %T data: %v\n", a.reqAddOrg, tc.Data, tc.Data)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false
			d := tc.GetWorkflowData().(*InitData)
			resp := tc.GetData().(*pb.AddOrgResponse)
			if assert.NotNil(t, resp) {
				assert.Equal(t, d.OrgName, resp.OrgName)
				assert.Equal(t, d.OrgIP, utils.IPv4CIDRToStringNotation(resp.Ip))
				assert.Equal(t, d.OrgCerts, resp.Certificate)
				assert.Equal(t, true, tc.Watcher.Expections())
				check = true
			}

			return check, nil
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Here we save any data required to be saved from the test case
			Cleanup any test specific data
			*/
			resp := tc.GetData().(*pb.AddOrgResponse)
			a := tc.GetWorkflowData().(*InitData)
			a.ROrgIp = resp.Ip
			tc.SaveWorkflowData(a)
			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()
			return nil
		},
	})

	err := w.Run(context.Background())
	assert.NoError(t, err)

	log.Infof("Workflow Status: \t %s", w.String())
}
