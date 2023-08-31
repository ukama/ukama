package subscriber

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker/v4"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"

	napi "github.com/ukama/ukama/systems/nucleus/api-gateway/pkg/rest"
	rapi "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"
	api "github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/rest"
	rpb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	mpb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	smutil "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/utils"
	spb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"github.com/ukama/ukama/testing/integration/pkg"
	"github.com/ukama/ukama/testing/integration/pkg/dataplan"
	"github.com/ukama/ukama/testing/integration/pkg/registry"
	"github.com/ukama/ukama/testing/integration/pkg/test"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
)

const MAX_POOL = 5

var config *pkg.Config

type InitData struct {
	Sys             *SubscriberClient
	Reg             *registry.RegistryClient
	Host            string
	RegHost         string
	SimType         string `default:"test"`
	ICCID           []string
	SimToken        []string
	MbHost          string
	SubscriberId    string
	OrgId           string
	OrgName         string
	NetworkId       string
	NetworkName     string
	UserId          string
	PackageId       string
	AddPackageId    string
	ActivePackageId string
	EncKey          string
	SimId           string
	SimStatus       string
	AllocatedSim    *spb.Sim

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
	reqRemovePkgFromSimReq       api.RemovePkgFromSimReq
	reqAllocateSimReq            api.AllocateSimReq
	reqActivateDeactivateSimReq  api.ActivateDeactivateSimReq
	reqSetActivePackageForSimReq api.SetActivePackageForSimReq
	reqAddOrgRequest             napi.AddOrgRequest
	reqAddNetworkRequest         rapi.AddNetworkRequest
	reqAddUserRequest            napi.AddUserRequest

	/* API responses */

	/* Dependencies on other system */
	//wfReg      *test.Workflow
	wfDataPlan *test.Workflow
}

func InitializeData() *InitData {
	config = pkg.NewConfig()

	d := &InitData{}

	d.ICCID = make([]string, MAX_POOL)
	//d.SimToken = make([]string, MAX_POOL)
	d.Host = config.System.Subscriber
	d.RegHost = config.System.Registry
	d.MbHost = config.System.MessageBus

	d.Sys = NewSubscriberClient(d.Host)
	d.Reg = registry.NewRegistryClient(d.RegHost)
	d.EncKey = "the-key-has-to-be-32-bytes-long!"
	d.SimType = "test"

	d.reqAddOrgRequest = napi.AddOrgRequest{
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

	d.reqAddUserRequest = napi.AddUserRequest{
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
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
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
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
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
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
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
		log.Tracef("Setting up watcher for %s", tc.Name)
		tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.simmanager.sim.allocate"})
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
				resp.Sim.Package != nil &&
				d.reqAllocateSimReq.PackageId == resp.Sim.Package.PackageId &&
				d.reqAllocateSimReq.SimType == resp.Sim.Type &&
				d.reqAllocateSimReq.NetworkId == resp.Sim.NetworkId &&
				tc.Watcher.Expections() {
				check = true
			}
		}

		return check, nil
	},

	ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Here we save any data required to be saved from the test case
		Cleanup any test specific data
		*/
		resp := tc.GetData().(*mpb.AllocateSimResponse)
		a := tc.GetWorkflowData().(*InitData)
		a.SimId = resp.Sim.Id
		a.ActivePackageId = resp.Sim.Package.PackageId
		tc.SaveWorkflowData(a)

		tc.Watcher.Stop()
		return nil
	},
}

var TC_registry_get_subscriber = &test.TestCase{
	Name:        "Get Subscriber ",
	Description: "Get subscriber",
	Data:        &rpb.GetSubscriberResponse{},

	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		a.reqSubscriberGetReq.SubscriberId = a.SubscriberId
		tc.SaveWorkflowData(a)
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Sys.SubscriberRegistryGetSusbscriber(a.reqSubscriberGetReq)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		/* Check for possible failures during test case */
		check := false

		resp := tc.GetData().(*rpb.GetSubscriberResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if data.reqSubscriberGetReq.SubscriberId == resp.Subscriber.SubscriberId {
				check = true
			}
		}

		return check, nil
	},
}

var TC_manager_get_sim_by_subscriber = &test.TestCase{
	Name:        "Get Sim ",
	Description: "Get Sim by subscriber",
	Data:        &mpb.GetSimsBySubscriberResponse{},

	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		a.reqGetSimsBySubReq.SubscriberId = a.SubscriberId
		tc.SaveWorkflowData(a)
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Sys.SubscriberManagerGetSubscriberSims(a.reqGetSimsBySubReq)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		/* Check for possible failures during test case */
		check := false

		resp := tc.GetData().(*mpb.GetSimsBySubscriberResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if data.reqGetSimsBySubReq.SubscriberId == resp.SubscriberId &&
				len(resp.Sims) > 0 &&
				data.SimId == resp.Sims[0].Id &&
				data.SimStatus == resp.Sims[0].Status {
				check = true
			}
		}

		return check, nil
	},
}

var TC_manager_get_package_for_sim = &test.TestCase{
	Name:        "Get package for a sim ",
	Description: "Get package for a sim",
	Data:        &mpb.GetPackagesBySimResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		a.reqSimReq.SimId = a.SimId
		tc.SaveWorkflowData(a)
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Sys.SubscriberManagerGetPackageForSim(a.reqSimReq)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		/* Check for possible failures during test case */
		check := false

		resp := tc.GetData().(*mpb.GetPackagesBySimResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if data.reqSimReq.SimId == resp.SimId &&
				len(resp.Packages) > 0 && func(ps []*mpb.Package, id string) bool {
				for _, p := range ps {
					if id == p.PackageId {
						return true
					}
				}
				return false
			}(resp.Packages, data.PackageId) {
				check = true
			}
		}

		return check, nil
	},
}

var TC_manager_check_active_package_for_sim = &test.TestCase{
	Name:        "Check active package for a sim ",
	Description: "Check active package for a sim",
	Data:        &mpb.GetPackagesBySimResponse{},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Sys.SubscriberManagerGetPackageForSim(a.reqSimReq)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		/* Check for possible failures during test case */
		check := false

		resp := tc.GetData().(*mpb.GetPackagesBySimResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if data.reqSimReq.SimId == resp.SimId &&
				len(resp.Packages) > 0 && func(ps []*mpb.Package, id string) bool {
				for _, p := range ps {
					if id == p.PackageId && p.IsActive {
						return true
					}
				}
				return false
			}(resp.Packages, data.ActivePackageId) {
				check = true
			}
		}

		return check, nil
	},
}

var TC_manager_activate_sim = &test.TestCase{
	Name:        "Activate sim",
	Description: "Activate a sim of subscriber",
	Data:        &mpb.ToggleSimStatusResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		a.SimStatus = "active"
		a.reqActivateDeactivateSimReq.SimId = a.SimId
		a.reqActivateDeactivateSimReq.Status = a.SimStatus

		tc.SaveWorkflowData(a)
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Sys.SubscriberManagerUpdateSim(a.reqActivateDeactivateSimReq)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		/* Check for possible failures during test case */

		return true, nil
	},
}

var TC_manager_inactivate_sim = &test.TestCase{
	Name:        "Inactivate sim",
	Description: "Inactivate a sim of subscriber",
	Data:        &mpb.ToggleSimStatusResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		a.SimStatus = "inactive"
		a.reqActivateDeactivateSimReq.SimId = a.SimId
		a.reqActivateDeactivateSimReq.Status = a.SimStatus

		tc.SaveWorkflowData(a)

		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Sys.SubscriberManagerUpdateSim(a.reqActivateDeactivateSimReq)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		/* Check for possible failures during test case */
		return true, nil
	},
}

var TC_manager_get_sim = &test.TestCase{
	Name:        "Get Sim ",
	Description: "Get Sim by SimId",
	Data:        &mpb.GetSimResponse{},

	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		a.reqSimReq.SimId = a.SimId
		tc.SaveWorkflowData(a)
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Sys.SubscriberManagerGetSim(a.reqSimReq)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		/* Check for possible failures during test case */
		check := false

		resp := tc.GetData().(*mpb.GetSimResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if data.reqSimReq.SimId == resp.Sim.Id {
				check = true
			}
		}

		return check, nil
	},
}

var TC_manager_add_extra_package_to_sim = &test.TestCase{
	Name:        "Add an extra package to sim",
	Description: "Allocating multiple packages to a sim of subscriber",
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)

		/* Add a new package */
		err := func(t *testing.T) error {
			dTC := a.wfDataPlan.GetTestCase(dataplan.TC_dp_add_package.Name)
			if dTC == nil {
				return fmt.Errorf("%s", "invalid test case name")
			}

			d := dTC.GetWorkflowData().(*dataplan.InitData)
			d.ReqAddPackageRequest.Name = "Additonal-package"

			/* Add a package to Data plan first */
			err := a.wfDataPlan.ExecuteTestCase(t, ctx, dTC)
			if err != nil {
				log.Errorf("Adding addtional test package failed: %v", err)
				return err
			}

			a.AddPackageId = d.PackageId
			tc.SaveWorkflowData(a)
			return nil
		}(t)
		if err != nil {
			log.Errorf("Failed to add new pacakge.")
			return err
		}

		a.reqAddPkgToSimReq.PackageId = a.AddPackageId
		a.reqAddPkgToSimReq.SimId = a.SimId
		a.reqAddPkgToSimReq.StartDate = timestamppb.New(time.Now().Add(24 * time.Hour))
		a.ActivePackageId = a.AddPackageId

		tc.SaveWorkflowData(a)

		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			err = a.Sys.SubscriberManagerAddPackage(a.reqAddPkgToSimReq)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},
}

var TC_manager_get_multiple_package_for_sim = &test.TestCase{
	Name:        "Get multiple package for a sim ",
	Description: "Get multiple package for a sim",
	Data:        &mpb.GetPackagesBySimResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		a.reqSimReq.SimId = a.SimId
		tc.SaveWorkflowData(a)
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Sys.SubscriberManagerGetPackageForSim(a.reqSimReq)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		/* Check for possible failures during test case */
		check := false

		resp := tc.GetData().(*mpb.GetPackagesBySimResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if data.reqSimReq.SimId == resp.SimId &&
				len(resp.Packages) > 0 && func(ps []*mpb.Package, id string) bool {
				for _, p := range ps {
					if id == p.PackageId {
						return true
					}
				}
				return false
			}(resp.Packages, data.PackageId) &&
				func(ps []*mpb.Package, id string) bool {
					for _, p := range ps {
						if id == p.PackageId {
							return true
						}
					}
					return false
				}(resp.Packages, data.AddPackageId) {
				check = true
			}
		}

		return check, nil
	},
}

var TC_manager_set_active_package_for_sim = &test.TestCase{
	Name:        "Set active package for a sim ",
	Description: "Set active package for a sim",
	Data:        &mpb.SetActivePackageResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		a.reqSetActivePackageForSimReq.SimId = a.SimId
		a.reqSetActivePackageForSimReq.PackageId = a.ActivePackageId
		tc.SaveWorkflowData(a)
		log.Tracef("Setting up watcher for %s", tc.Name)
		tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.simmanager.sim.activepackage"})
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Sys.SubscriberManagerActivatePackage(a.reqSetActivePackageForSimReq)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		/* Check for possible failures during test case */
		check := false

		if tc.Watcher.Expections() {
			check = true
			a := tc.GetWorkflowData().(*InitData)
			a.ActivePackageId = a.reqSetActivePackageForSimReq.PackageId
			tc.SaveWorkflowData(a)
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

var TC_manager_set_delete_package_for_sim = &test.TestCase{
	Name:        "Delete package for a sim ",
	Description: "Delete package for a sim",

	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		a.reqRemovePkgFromSimReq.SimId = a.SimId
		a.reqRemovePkgFromSimReq.PackageId = a.PackageId
		tc.SaveWorkflowData(a)
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			err = a.Sys.SubscriberManagerDeletePackage(a.reqRemovePkgFromSimReq)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},
}

var TC_manager_get_package_after_removal_for_sim = &test.TestCase{
	Name:        "Get availaible packages for a sim ",
	Description: "Get availaible packages after removal of a pacakge for a sim",
	Data:        &mpb.GetPackagesBySimResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		a.reqSimReq.SimId = a.SimId
		tc.SaveWorkflowData(a)
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {
			tc.Data, err = a.Sys.SubscriberManagerGetPackageForSim(a.reqSimReq)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		/* Check for possible failures during test case */
		check := false

		resp := tc.GetData().(*mpb.GetPackagesBySimResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if data.reqSimReq.SimId == resp.SimId &&
				len(resp.Packages) > 0 && func(ps []*mpb.Package, id string) bool {
				for _, p := range ps {
					if id == p.PackageId {
						return false
					}
				}
				return true
			}(resp.Packages, data.PackageId) &&
				func(ps []*mpb.Package, id string) bool {
					for _, p := range ps {
						if id == p.PackageId {
							return true
						}
					}
					return false
				}(resp.Packages, data.AddPackageId) {
				check = true
			}
		}

		return check, nil
	},
}
