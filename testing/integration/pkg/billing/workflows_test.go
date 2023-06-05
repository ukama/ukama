package pkg_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/num30/config"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/testing/integration/pkg"
	"github.com/ukama/ukama/testing/integration/pkg/test"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"

	dapi "github.com/ukama/ukama/systems/data-plan/api-gateway/pkg/rest"
	rapi "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"
	sapi "github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/rest"

	bilutil "github.com/ukama/ukama/systems/billing/invoice/pkg/util"

	billing "github.com/ukama/ukama/testing/integration/pkg/billing"
	dplan "github.com/ukama/ukama/testing/integration/pkg/dataplan"
	reg "github.com/ukama/ukama/testing/integration/pkg/registry"
	subs "github.com/ukama/ukama/testing/integration/pkg/subscriber"
)

var errTestFailure = errors.New("test failure")

type BillingData struct {
	SimType string `default:"ukama_data"`

	BillingClient *billing.BillingClient
	Host          string
	MbHost        string

	DataPlanClient *dplan.DataplanClient
	DplanHost      string
	BaseRateId     []string
	BaserateId     string
	PackageId      string
	Country        string
	Provider       string

	SubscriberClient *subs.SubscriberClient
	SubsHost         string
	SubscriberId     string
	SubscriberName   string
	SubscriberEmail  string
	SubscriberPhone  string

	RegistryClient *reg.RegistryClient
	RegHost        string
	OwnerName      string
	OwnerEmail     string
	OwnerPhone     string
	OrgName        string
	OwnerId        string
	OrgId          string
	NetName        string
	NetworkId      string

	// API requests
	reqAddUserRequest               rapi.AddUserRequest
	reqAddOrgRequest                rapi.AddOrgRequest
	reqAddNetworkRequest            rapi.AddNetworkRequest
	reqUploadBaseRatesRequest       dapi.UploadBaseRatesRequest
	reqGetBaseRatesByCountryRequest dapi.GetBaseRatesByCountryRequest
	reqGetBaseRateRequest           dapi.GetBaseRateRequest
	reqSetDefaultMarkupRequest      dapi.SetDefaultMarkupRequest
	reqSetMarkupRequest             dapi.SetMarkupRequest
	reqAddPackageRequest            dapi.AddPackageRequest
	reqSubscriberAddReq             sapi.SubscriberAddReq
}

var serviceConfig = pkg.NewConfig()

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stderr)

	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatalf("Error reading config file. Error: %v", err)
	} else if serviceConfig.DebugMode {
		// output config in debug mode
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

}

func InitializeData() *BillingData {
	d := &BillingData{}

	d.SimType = "test"

	d.Host = "http://localhost:3000"
	d.MbHost = "amqp://guest:guest@localhost:5672/"

	d.BillingClient = billing.NewBillingClient(d.Host, serviceConfig.Key)

	d.DplanHost = "http://localhost:8080"
	d.DataPlanClient = dplan.NewDataplanClient(d.DplanHost)
	d.BaseRateId = make([]string, 8)
	d.Country = "The lunar maria"

	d.Provider = "ABC Tel"

	d.SubsHost = "http://localhost:8081"
	d.SubscriberClient = subs.NewSubscriberClient(d.SubsHost)
	d.SubscriberName = faker.FirstName()
	d.SubscriberEmail = strings.ToLower(faker.FirstName() + "_" + faker.LastName() + "@example.com")
	d.SubscriberPhone = faker.Phonenumber()

	d.RegHost = "http://localhost:8082"
	d.RegistryClient = reg.NewRegistryClient(d.RegHost)
	d.OwnerName = strings.ToLower(faker.FirstName())
	d.OwnerEmail = strings.ToLower(faker.Email())
	d.OwnerPhone = strings.ToLower(faker.Phonenumber())
	d.OrgName = strings.ToLower(faker.FirstName() + "-org")
	d.NetName = strings.ToLower(faker.FirstName()) + "-net"

	d.reqAddUserRequest = rapi.AddUserRequest{
		Name:  d.OwnerName,
		Email: d.OwnerEmail,
		Phone: d.OwnerPhone,
	}

	d.reqAddOrgRequest = rapi.AddOrgRequest{
		OrgName:     d.OrgName,
		Owner:       d.OwnerId,
		Certificate: utils.RandomBase64String(2048),
	}

	d.reqAddNetworkRequest = rapi.AddNetworkRequest{
		OrgName: d.OrgName,
		NetName: d.NetName,
	}

	d.reqUploadBaseRatesRequest = dapi.UploadBaseRatesRequest{
		EffectiveAt: utils.GenerateFutureDate(5 * time.Second),
		FileURL:     "https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/template/template.csv",
		EndAt:       utils.GenerateFutureDate(365 * 24 * time.Hour),
		SimType:     d.SimType,
	}

	d.reqGetBaseRatesByCountryRequest = dapi.GetBaseRatesByCountryRequest{
		Country:  d.Country,
		Provider: d.Provider,
		SimType:  d.SimType,
	}

	d.reqSetDefaultMarkupRequest = dapi.SetDefaultMarkupRequest{
		Markup: float64(utils.RandomInt(50)),
	}

	d.reqSetMarkupRequest = dapi.SetMarkupRequest{
		OwnerId: d.OwnerId,
		Markup:  float64(utils.RandomInt(50)),
	}

	d.reqSubscriberAddReq = sapi.SubscriberAddReq{
		FirstName:             d.SubscriberName,
		LastName:              faker.LastName(),
		Email:                 d.SubscriberEmail,
		Phone:                 d.SubscriberPhone,
		Dob:                   utils.RandomPastDate(2000),
		ProofOfIdentification: "passport",
		IdSerial:              faker.UUIDDigit(),
		Address:               faker.Sentence(),
		Gender:                "male",
		OrgId:                 "",
		NetworkId:             "",
	}

	return d
}

func TestWorkflow_BillingSystem(t *testing.T) {
	w := test.NewWorkflow("Billing Workflows", "Various use cases regarding the billing system")

	w.SetUpFxn = func(t *testing.T, ctx context.Context, w *test.Workflow) error {
		log.Debugf("Initilizing Data for %s.", w.String())
		var err error

		d := InitializeData()

		w.Data = d

		// Add new user
		aUserResp, err := d.RegistryClient.AddUser(d.reqAddUserRequest)
		if assert.NoError(t, err) {
			assert.NotNil(t, aUserResp)
			assert.Equal(t, d.OwnerName, aUserResp.User.Name)
			assert.Equal(t, d.OwnerEmail, aUserResp.User.Email)
			assert.Equal(t, d.OwnerPhone, aUserResp.User.Phone)
		}

		d.OwnerId = aUserResp.User.Uuid

		// Add new org
		d.reqAddOrgRequest.Owner = d.OwnerId

		aOrgResp, err := d.RegistryClient.AddOrg(d.reqAddOrgRequest)
		if assert.NoError(t, err) {
			assert.NotNil(t, aOrgResp)
			assert.Equal(t, d.OrgName, aOrgResp.Org.Name)
		}

		d.OrgId = aOrgResp.Org.Id

		// Add new network
		aNetResp, err := d.RegistryClient.AddNetwork(d.reqAddNetworkRequest)
		if assert.NoError(t, err) {
			assert.Equal(t, d.NetName, aNetResp.Network.Name)
			assert.Equal(t, d.OrgId, aNetResp.Network.OrgId)
		}

		d.NetworkId = aNetResp.Network.Id

		// Add base rates
		abResp, err := d.DataPlanClient.DataPlanBaseRateUpload(d.reqUploadBaseRatesRequest)
		if assert.NoError(t, err) {
			assert.NotNil(t, abResp)
		}

		// Get one base rate
		d.reqGetBaseRateRequest = dapi.GetBaseRateRequest{
			RateId: d.BaseRateId[0],
		}

		gbResp, err := d.DataPlanClient.DataPlanBaseRateGetByCountry(d.reqGetBaseRatesByCountryRequest)
		if assert.NoError(t, err) {
			assert.NotNil(t, gbResp)

			if !assert.Equal(t, 1, len(gbResp.Rates)) {
				return fmt.Errorf("%w: setup failure while getting base rates", err)
			}
		}

		d.BaserateId = gbResp.Rates[0].Uuid

		// Set markup
		d.reqSetDefaultMarkupRequest = dapi.SetDefaultMarkupRequest{
			Markup: float64(utils.RandomInt(50)),
		}

		d.reqSetMarkupRequest.OwnerId = d.OwnerId

		mResp, err := d.DataPlanClient.DataPlanUpdateMarkup(d.reqSetMarkupRequest)
		if assert.NoError(t, err) {
			assert.NotNil(t, mResp)
		}

		log.Debugf("Workflow Data : %+v", w.Data)

		return err
	}

	w.RegisterTestCase(&test.TestCase{
		Name:        "Add plan from package",
		Description: "Add a billing plan for a new data package",
		Data:        &bilutil.Plan{},
		Workflow:    w,
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			a := tc.GetWorkflowData().(*BillingData)
			log.Tracef("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.package.package.create"})

			// Add a new package
			a.reqAddPackageRequest = dapi.AddPackageRequest{
				OwnerId:    a.OwnerId,
				OrgId:      a.OrgId,
				Name:       faker.FirstName() + "-monthly-pack",
				SimType:    a.SimType,
				From:       utils.GenerateFutureDate(24 * time.Hour),
				To:         utils.GenerateFutureDate(30 * 24 * time.Hour),
				BaserateId: a.BaserateId,
				SmsVolume:  100,
				DataVolume: 1024,
				DataUnit:   "MegaBytes",
				Type:       "postpaid",
				Active:     true,
				Flatrate:   false,
				Apn:        "ukama.tel",
			}

			pResp, err := a.DataPlanClient.DataPlanPackageAdd(a.reqAddPackageRequest)
			if assert.NoError(t, err) {
				assert.NotNil(t, pResp)
				assert.Equal(t, a.reqAddPackageRequest.OrgId, pResp.Package.OrgId)
				assert.NotNil(t, pResp.Package.Uuid)
				assert.Equal(t, true, tc.Watcher.Expections())
			}

			a.PackageId = pResp.Package.Uuid

			return err
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*BillingData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}

			tc.Data, err = a.BillingClient.GetPlan(a.PackageId)

			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false

			a, ok := tc.GetWorkflowData().(*BillingData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return check, fmt.Errorf("invalid data type for Workflow data")
			}

			tr, ok := tc.GetData().(*bilutil.Plan)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return check, fmt.Errorf("invalid data type for Workflow data")
			}

			if assert.NotNil(t, tr) {
				assert.Equal(t, a.PackageId, tr.Code)
				check = true
			}

			return check, nil
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Here we save any data required to be saved from the test case
			Cleanup any test specific data
			*/
			tc.Watcher.Stop()

			assert.Equal(t, int(tc.State), int(test.StateTypePass))
			return nil
		},
	})

	w.RegisterTestCase(&test.TestCase{
		Name:        "Add customer from subscriber",
		Description: "Add a billing customer for a new subscriber",
		Data:        &bilutil.Customer{},
		Workflow:    w,
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			a := tc.GetWorkflowData().(*BillingData)
			log.Tracef("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.registry.subscriber.create"})

			// Add new subscriber
			a.reqSubscriberAddReq.NetworkId = a.NetworkId
			a.reqSubscriberAddReq.OrgId = a.OrgId

			addSub, err := a.SubscriberClient.SubscriberRegistryAddSusbscriber(a.reqSubscriberAddReq)
			if assert.NoError(t, err) {
				assert.NotNil(t, addSub)
				assert.Equal(t, a.reqSubscriberAddReq.OrgId, addSub.Subscriber.OrgId)
				assert.Equal(t, a.reqSubscriberAddReq.Email, addSub.Subscriber.Email)
				assert.NotNil(t, addSub.Subscriber.SubscriberId)
				assert.Equal(t, true, tc.Watcher.Expections())
			}

			a.SubscriberId = addSub.Subscriber.SubscriberId

			return err
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*BillingData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}

			tc.Data, err = a.BillingClient.GetCustomer(a.SubscriberId)

			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false

			a, ok := tc.GetWorkflowData().(*BillingData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return check, fmt.Errorf("invalid data type for Workflow data")
			}

			tr, ok := tc.GetData().(*bilutil.Customer)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return check, fmt.Errorf("invalid data type for Workflow data")
			}

			if assert.NotNil(t, tr) {
				assert.Equal(t, a.SubscriberId, tr.ExternalID)
				assert.Equal(t, a.SubscriberName, tr.Name)
				assert.Equal(t, a.SubscriberEmail, tr.Email)
				assert.Equal(t, a.SubscriberPhone, tr.Phone)
				check = true
			}

			return check, nil
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Here we save any data required to be saved from the test case
			Cleanup any test specific data
			*/
			tc.Watcher.Stop()

			assert.Equal(t, int(tc.State), int(test.StateTypePass))
			return nil
		},
	})

	err := w.Run(t, context.Background())
	assert.NoError(t, err)

	w.Status()
}
