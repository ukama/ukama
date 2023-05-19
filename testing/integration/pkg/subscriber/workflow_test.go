package subscriber

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v4"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/systems/common/ukama"
	api "github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/rest"
	spb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"github.com/ukama/ukama/testing/integration/pkg/test"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
)

const MAX_POOL = 5

type InitData struct {
	OrgName                  string
	OrgIP                    string
	OrgCerts                 string
	SysName, SysIP, SysCerts string
	SysPort                  int32
	NodeId                   ukama.NodeID
	NodeIP, NodeCerts        string
	Sys                      *SubscriberSys
	Host                     string
	ROrgIp                   string
	SimType                  string `default:"ukama_data"`
	ICCID                    [MAX_POOL]string
	MbHost                   string

	/* API requests */
	reqSimPoolUploadSimReq       api.SimPoolUploadSimReq
	reqSimPoolStatByTypeReq      api.SimPoolStatByTypeReq
	reqSimByIccidReq             api.SimByIccidReq
	reqSubscriberGetReq          api.SubscriberGetReq
	reqSubscriberAddReq          api.SubscriberAddReq
	reqSubscriberDeleteReq       api.SubscriberDeleteReq
	reqSubscriberUpdateReq       api.SubscriberUpdateReq
	reqGetSimsBySubReq           api.GetSimsBySubReq
	reqSimReq                    api.SimReq
	reqAddPkgToSimReq            api.AddPkgToSimReq
	reqAllocateSimReq            api.AllocateSimReq
	reqActivateDeactivateSimReq  api.ActivateDeactivateSimReq
	reqSetActivePackageForSimReq api.SetActivePackageForSimReq

	/* API responses */

}

func init() {
	log.SetLevel(log.TraceLevel)
	log.SetOutput(os.Stderr)
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
	d.Host = "http://192.168.0.22:8078"
	d.MbHost = "amqp://guest:guest@192.168.0.22:5672/"
	d.Sys = NewSubscriberSys(d.Host)
	d.SimType = "ukama_data"

	d.reqSimPoolUploadSimReq = api.SimPoolUploadSimReq{
		SimType: d.SimType,
		Data:    string(CreateSimPool(MAX_POOL, d.ICCID)),
	}

	d.reqSimPoolStatByTypeReq = api.SimPoolStatByTypeReq{
		SimType: SimType,
	}

	d.reqSimByIccidReq = api.SimByIccidReq{
		Iccid: d.ICCID[utils.RandomIntInRange(0, MAX_POOL-1)],
	}

	d.reqSubscriberGetReq = api.SubscriberGetReq{
		SubscriberId: uuid.NewV4().String(),
	}

	d.reqSubscriberAddReq = api.SubscriberAddReq{
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

	d.reqSubscriberDeleteReq = api.SubscriberDeleteReq{
		SubscriberId: "",
	}

	d.reqSubscriberUpdateReq = api.SubscriberUpdateReq{
		SubscriberId:          "",
		Email:                 faker.Email(),
		Phone:                 faker.Phonenumber(),
		ProofOfIdentification: "dl",
		IdSerial:              faker.UUIDDigit(),
		Address:               faker.Sentence(),
	}

	d.reqGetSimsBySubReq = api.GetSimsBySubReq{
		SubscriberId: "",
	}

	d.reqSimReq = api.SimReq{
		SimId: "",
	}

	d.reqAddPkgToSimReq = api.AddPkgToSimReq{
		SimId:     "",
		PackageId: "",
		//StartDate: "",
	}

	d.reqAllocateSimReq = api.AllocateSimReq{
		SubscriberId: "",
		SimToken:     "",
		PackageId:    "",
		NetworkId:    "",
		SimType:      "",
	}

	d.reqActivateDeactivateSimReq = api.ActivateDeactivateSimReq{
		SimId:  "",
		Status: "",
	}

	d.reqSetActivePackageForSimReq = api.SetActivePackageForSimReq{
		SimId:     "",
		PackageId: "",
	}
	return d
}

func CreateSimPool(count int, id [MAX_POOL]string) []byte {

	idx := 0
	str := "ICCID,MSISDN,SmDpAddress,ActivationCode,IsPhysical,QrCode"
	for count != 0 {
		id[idx] = fmt.Sprintf("891030000000%d%d", utils.RandomIntInRange(1000, 9999), utils.RandomIntInRange(10000, 99999))
		str = str + fmt.Sprintf("\n%s,%s,%s,%d,%t,%s", id[idx], faker.Phonenumber(), utils.RandomIPv4(), utils.RandomIntInRange(1000, 9999), false, faker.Word())
		count--
		idx++
	}
	pool := make([]byte, b64.StdEncoding.EncodedLen(len(str)))
	log.Tracef("Simpool: %s", str)
	b64.StdEncoding.Encode(pool, []byte(str))
	return pool
}

func TestWorkflow_SubscriberSystem(t *testing.T) {
	w := test.NewWorkflow("susbcriber_workflow_1", "Adding sims to sim pool")

	w.SetUpFxn = func(ctx context.Context, w *test.Workflow) error {
		log.Tracef("Initilizing Data for %s.", w.String())
		w.Data = InitializeData()

		log.Tracef("Workflow Data : %+v", w.Data)
		return nil
	}

	w.RegisterTestCase(&test.TestCase{
		Name:        "Add Sims to sim-pool",
		Description: "Add sims to sim-pool from the base64 encoded file",
		Data:        &spb.UploadResponse{},
		Workflow:    w,
		SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			a := tc.GetWorkflowData().(*InitData)
			log.Tracef("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.sim.sim.upload"})
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*InitData)
			if ok {
				tc.Data, err = a.Sys.SubscriberSimpoolUploadSims(a.reqSimPoolUploadSimReq)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false

			resp := tc.GetData().(*spb.UploadResponse)
			if assert.NotNil(t, resp) {
				data := tc.GetWorkflowData().(*InitData)
				assert.Equal(t, len(data.ICCID), len(resp.Iccid))
				assert.Equal(t, true, tc.Watcher.Expections())
				check = true
			}

			return check, nil
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Here we save any data required to be saved from the test case
			Cleanup any test specific data
			*/
			tc.Watcher.Stop()
			return nil
		},
	})

	err := w.Run(t, context.Background())
	assert.NoError(t, err)

	w.Status()

}
