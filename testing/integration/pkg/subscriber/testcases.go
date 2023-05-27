package subscriber

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"strings"

	"github.com/bxcodec/faker/v4"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	rapi "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"
	api "github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/rest"
	rpb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	mpb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	smutil "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/utils"
	spb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"github.com/ukama/ukama/testing/integration/pkg/registry"
	"github.com/ukama/ukama/testing/integration/pkg/test"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
)

const MAX_POOL = 5

type InitData struct {
	Sys          *SubscriberSys
	Reg          *registry.RegistrySys
	Host         string
	RegHost      string
	SimType      string `default:"ukama_data"`
	ICCID        []string
	SimToken     []string
	MbHost       string
	SubscriberId string
	OrgId        string
	OrgName      string
	NetworkId    string
	NetworkName  string
	UserId       string
	PackageId    string
	EncKey       string

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
	reqAddOrgRequest             rapi.AddOrgRequest
	reqAddNetworkRequest         rapi.AddNetworkRequest
	reqAddUserRequest            rapi.AddUserRequest

	/* API responses */

}

func InitializeData() *InitData {
	d := &InitData{}
	d.ICCID = make([]string, MAX_POOL)
	//d.SimToken = make([]string, MAX_POOL)
	d.Host = "http://192.168.0.23:8078"
	d.RegHost = "http://192.168.0.23:8075"
	d.MbHost = "amqp://guest:guest@192.168.0.23:5672/"
	d.Sys = NewSubscriberSys(d.Host)
	d.Reg = registry.NewRegistrySys(d.RegHost)
	d.EncKey = "the-key-has-to-be-32-bytes-long!"
	d.SimType = "ukama_data"

	d.reqAddOrgRequest = rapi.AddOrgRequest{
		OrgName:     strings.ToLower(faker.FirstName() + "-org"),
		Owner:       "",
		Certificate: utils.RandomBase64String(2048),
	}

	d.reqAddNetworkRequest = rapi.AddNetworkRequest{}

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
		OrgId:                 "",
		NetworkId:             "",
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

	d.reqAddUserRequest = rapi.AddUserRequest{
		Name:  d.reqSubscriberAddReq.FirstName,
		Email: d.reqSubscriberAddReq.Email,
		Phone: d.reqSubscriberAddReq.Phone,
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

var TC_simpool_upload = &test.TestCase{
	Name:        "Add Sims to sim-pool",
	Description: "Add sims to sim-pool from the base64 encoded file",
	Data:        &spb.UploadResponse{},
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

		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if len(data.ICCID) == len(resp.Iccid) &&
				tc.Watcher.Expections() {
				check = true
				for idx, i := range data.ICCID {
					if i != data.ICCID[idx] {
						check = false
						break
					} else {
						tok, err := smutil.GenerateTokenFromIccid(i, data.EncKey)
						if err == nil {
							data.SimToken = append(data.SimToken, tok)
						}
					}
					tc.SaveWorkflowData(data)
				}
			}

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
}

var TC_simpool_get_sim = &test.TestCase{
	Name:        "Get Sim from sim-pool",
	Description: "Get Sim from sim pool by ICCID",
	Data:        &spb.GetByIccidResponse{},

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
		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if data.reqSimByIccidReq.Iccid == resp.Sim.Iccid &&
				data.SimType == resp.Sim.SimType {
				check = true
			}
		}

		return check, nil
	},
}

var TC_simpool_get_stats = &test.TestCase{
	Name:        "Get stats from sim-pool",
	Description: "Get stats from sim pool by sim type",
	Data:        &spb.GetStatsResponse{},

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
		if resp != nil {
			check = true
		}

		return check, nil
	},
}

var TC_registry_add_subscriber = &test.TestCase{
	Name:        "Add Subscriber",
	Description: "Add subscriber to registry",
	Data:        &rpb.AddSubscriberResponse{},
	SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		a.reqSubscriberAddReq.NetworkId = a.NetworkId
		a.reqSubscriberAddReq.OrgId = a.OrgId
		tc.SaveWorkflowData(a)
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
		if resp != nil {
			log.Tracef("Resp data is %v", resp)
			d := tc.GetWorkflowData().(*InitData)
			if d.reqSubscriberAddReq.Email == resp.Subscriber.Email {
				check = true
			}
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
}

var TC_manager_allocate_sim = &test.TestCase{
	Name:        "Allocate sim",
	Description: "Allocating a sim to subscriber",
	Data:        &mpb.AllocateSimResponse{},
	SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		a.reqAllocateSimReq.NetworkId = a.NetworkId
		a.reqAllocateSimReq.PackageId = a.PackageId
		a.reqAllocateSimReq.SimType = a.SimType
		a.reqAllocateSimReq.SubscriberId = a.SubscriberId
		a.reqAllocateSimReq.SimToken = a.SimToken[utils.RandomInt(len(a.SimToken)-1)]

		tc.SaveWorkflowData(a)
		// log.Tracef("Setting up watcher for %s", tc.Name)
		// tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.sim.sim.upload"})
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Sys.SubscriberManagerAllocateSim(a.reqAllocateSimReq)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		/* Check for possible failures during test case */
		check := false

		resp := tc.GetData().(*mpb.AllocateSimResponse)
		if resp != nil {
			log.Tracef("Resp data is %v", resp)
			d := tc.GetWorkflowData().(*InitData)
			if d.reqAllocateSimReq.SubscriberId == resp.Sim.SubscriberId &&
				d.reqAllocateSimReq.PackageId == resp.Sim.Package.Id &&
				d.reqAllocateSimReq.SimType == resp.Sim.Type &&
				d.reqAllocateSimReq.NetworkId == resp.Sim.NetworkId {
				check = true
			}
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
}
