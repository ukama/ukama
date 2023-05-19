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

	api "github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/rest"
	rpb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	spb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"github.com/ukama/ukama/testing/integration/pkg/test"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
)

const MAX_POOL = 5

type InitData struct {
	Sys          *SubscriberSys
	Host         string
	SimType      string `default:"ukama_data"`
	ICCID        []string
	MbHost       string
	SubscriberId string

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
	d.ICCID = make([]string, MAX_POOL)
	d.Host = "http://192.168.0.22:8078"
	d.MbHost = "amqp://guest:guest@192.168.0.22:5672/"
	d.Sys = NewSubscriberSys(d.Host)
	d.SimType = "ukama_data"

	d.reqSimPoolUploadSimReq = api.SimPoolUploadSimReq{
		SimType: d.SimType,
		Data:    string(CreateSimPool(MAX_POOL, &d.ICCID)),
	}

	d.reqSimPoolStatByTypeReq = api.SimPoolStatByTypeReq{
		SimType: d.SimType,
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
		Email:                 strings.ToLower(faker.FirstName() + "_" + faker.LastName() + "@gmail.com"),
		Phone:                 faker.Phonenumber(),
		Dob:                   utils.RandomPastDate(2000),
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

func CreateSimPool(count int, id *[]string) []byte {

	idx := 0
	str := "ICCID,MSISDN,SmDpAddress,ActivationCode,IsPhysical,QrCode"
	for count != 0 {
		(*id)[idx] = fmt.Sprintf("891030000000%d%d", utils.RandomIntInRange(1000, 9999), utils.RandomIntInRange(10000, 99999))
		str = str + fmt.Sprintf("\n%s,%s,%s,%d,%t,%s", (*id)[idx], faker.Phonenumber(), utils.RandomIPv4(), utils.RandomIntInRange(1000, 9999), false, faker.Word())
		count--
		idx++
	}
	pool := make([]byte, b64.StdEncoding.EncodedLen(len(str)))
	log.Tracef("Simpool: %s", str)
	b64.StdEncoding.Encode(pool, []byte(str))
	return pool
}

func TestWorkflow_SubscriberSystem(t *testing.T) {

	/* Sim pool */
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
				assert.Equal(t, data.ICCID, resp.Iccid)
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

	w.RegisterTestCase(&test.TestCase{
		Name:        "Get Sim from sim-pool",
		Description: "Get Sim from sim pool by ICCID",
		Data:        &spb.GetByIccidResponse{},
		Workflow:    w,

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*InitData)
			if ok {
				tc.Data, err = a.Sys.SubscriberSimpoolGetSimByICCID(a.reqSimByIccidReq)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false

			resp := tc.GetData().(*spb.GetByIccidResponse)
			if assert.NotNil(t, resp) {
				data := tc.GetWorkflowData().(*InitData)
				assert.Equal(t, data.reqSimByIccidReq.Iccid, resp.Sim.Iccid)
				assert.Equal(t, data.SimType, resp.Sim.SimType)
				check = true
			}

			return check, nil
		},
	})

	w.RegisterTestCase(&test.TestCase{
		Name:        "Get stats from sim-pool",
		Description: "Get stats from sim pool by sim type",
		Data:        &spb.GetStatsResponse{},
		Workflow:    w,

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*InitData)
			if ok {
				tc.Data, err = a.Sys.SubscriberSimpoolGetSimStats(a.reqSimPoolStatByTypeReq)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false

			resp := tc.GetData().(*spb.GetStatsResponse)
			if assert.NotNil(t, resp) {
				check = true
			}

			return check, nil
		},
	})

	/* subscriber */
	w.RegisterTestCase(&test.TestCase{
		Name:        "Add Subscriber",
		Description: "Add subscriber to registry",
		Data:        &rpb.AddSubscriberResponse{},
		Workflow:    w,
		SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			//a := tc.GetWorkflowData().(*InitData)
			// log.Tracef("Setting up watcher for %s", tc.Name)
			// tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.sim.sim.upload"})
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*InitData)
			if ok {
				tc.Data, err = a.Sys.SubscriberRegistryAddSusbscriber(a.reqSubscriberAddReq)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false

			resp := tc.GetData().(*rpb.AddSubscriberResponse)
			if assert.NotNil(t, resp) {
				d := tc.GetWorkflowData().(*InitData)
				assert.Equal(t, d.reqSubscriberAddReq.Email, resp.Subscriber.Email)
				assert.Equal(t, d.reqSubscriberAddReq.ProofOfIdentification, resp.Subscriber.ProofOfIdentification)
				assert.Equal(t, d.reqSubscriberAddReq.IdSerial, resp.Subscriber.IdSerial)
				assert.Equal(t, true, tc.Watcher.Expections())
				check = true
			}

			return check, nil
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Here we save any data required to be saved from the test case
			Cleanup any test specific data
			*/
			resp := tc.GetData().(*rpb.AddSubscriberResponse)
			a := tc.GetWorkflowData().(*InitData)
			a.SubscriberId = resp.Subscriber.SubscriberId
			tc.SaveWorkflowData(a)

			//tc.Watcher.Stop()
			return nil
		},
	})

	/* Run */
	err := w.Run(t, context.Background())
	assert.NoError(t, err)

	w.Status()

}
