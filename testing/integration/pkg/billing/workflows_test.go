package pkg_test

import (
	"context"
	"errors"
	"fmt"
	"os"
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

	"github.com/ukama/ukama/systems/common/uuid"
	dapi "github.com/ukama/ukama/systems/data-plan/api-gateway/pkg/rest"

	bilutil "github.com/ukama/ukama/systems/billing/invoice/pkg/util"

	billing "github.com/ukama/ukama/testing/integration/pkg/billing"
	dplan "github.com/ukama/ukama/testing/integration/pkg/dataplan"
)

var errTestFailure = errors.New("test failure")

type BillingData struct {
	OwnerId string
	OrgId   string
	SimType string `default:"ukama_data"`

	BillingClient *billing.BillingClient
	Host          string
	MbHost        string

	DataPlanClient *dplan.DataplanClient
	DplanHost      string
	PackageId      string
	BaseRateId     []string
	Country        string
	Provider       string

	// API requests
	reqUploadBaseRatesRequest       dapi.UploadBaseRatesRequest
	reqGetBaseRatesByCountryRequest dapi.GetBaseRatesByCountryRequest
	reqGetBaseRateRequest           dapi.GetBaseRateRequest
	reqSetDefaultMarkupRequest      dapi.SetDefaultMarkupRequest
	reqSetMarkupRequest             dapi.SetMarkupRequest
	reqAddPackageRequest            dapi.AddPackageRequest
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

	d.OwnerId = uuid.NewV4().String()
	d.OrgId = uuid.NewV4().String()
	// d.SimType = "ukama_data"
	d.SimType = "test"

	d.Host = "http://localhost:3000"
	d.MbHost = "amqp://guest:guest@localhost:5672/"

	d.BillingClient = billing.NewBillingClient(d.Host, serviceConfig.Key)

	d.DplanHost = "http://localhost:8080"
	d.DataPlanClient = dplan.NewDataplanClient(d.DplanHost)
	d.BaseRateId = make([]string, 8)
	d.Country = "The lunar maria"
	// d.Country = "Tycho crater"

	d.Provider = "ABC Tel"
	// d.Provider = "OWS Tel"

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

	return d
}

func TestWorkflow_BillingSystem(t *testing.T) {
	w := test.NewWorkflow("Billing Workflows", "Various use cases regarding the billing system")

	w.SetUpFxn = func(t *testing.T, ctx context.Context, w *test.Workflow) error {
		log.Debugf("Initilizing Data for %s.", w.String())
		var err error

		d := InitializeData()

		w.Data = d

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

			// Add base rates, ignore error for dupes if rates are already present
			abResp, err := a.DataPlanClient.DataPlanBaseRateUpload(a.reqUploadBaseRatesRequest)
			if assert.NoError(t, err) {
				assert.NotNil(t, abResp)
			}

			// Get one base rate
			a.reqGetBaseRateRequest = dapi.GetBaseRateRequest{
				RateId: a.BaseRateId[0],
			}

			// gbResp, err := a.DataPlanClient.DataPlanBaseRateGet(a.reqGetBaseRateRequest)
			gbResp, err := a.DataPlanClient.DataPlanBaseRateGetByCountry(a.reqGetBaseRatesByCountryRequest)
			if assert.NoError(t, err) {

				assert.NotNil(t, gbResp)

				if !assert.Equal(t, 1, len(gbResp.Rates)) {
					return fmt.Errorf("%w: setup failure while getting base rates", err)
				}
			}

			// Set markup
			a.reqSetDefaultMarkupRequest = dapi.SetDefaultMarkupRequest{
				Markup: float64(utils.RandomInt(50)),
			}

			mResp, err := a.DataPlanClient.DataPlanUpdateMarkup(a.reqSetMarkupRequest)
			if assert.NoError(t, err) {
				assert.NotNil(t, mResp)
			}

			// Add a new package
			a.reqAddPackageRequest = dapi.AddPackageRequest{
				OwnerId:    a.OwnerId,
				OrgId:      a.OrgId,
				Name:       faker.FirstName() + "-monthly-pack",
				SimType:    a.SimType,
				From:       utils.GenerateFutureDate(24 * time.Hour),
				To:         utils.GenerateFutureDate(30 * 24 * time.Hour),
				BaserateId: gbResp.Rates[0].Uuid,
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

	err := w.Run(t, context.Background())
	assert.NoError(t, err)

	w.Status()
}
