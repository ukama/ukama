package dataplan

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bxcodec/faker/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	uuid "github.com/ukama/ukama/systems/common/uuid"

	api "github.com/ukama/ukama/systems/data-plan/api-gateway/pkg/rest"
	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	ppb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	rpb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	"github.com/ukama/ukama/testing/integration/pkg/test"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
)

type InitData struct {
	Sys          *DataPlanSys
	Host         string
	SimType      string `default:"ukama_data"`
	MbHost       string
	SubscriberId string
	BaseRateId   []string
	PackageId    string

	// This data is taken from the https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/template/template.csv */
	Providers map[string][]string
	Countries []string
	OwnerId   string
	OrgId     string

	/* API requests */
	reqUploadBaseRatesRequest         api.UploadBaseRatesRequest
	reqGetBaseRateRequest             api.GetBaseRateRequest
	reqGetBaseRatesByCountryRequest   api.GetBaseRatesByCountryRequest
	reqGetBaseRatesForPeriodRequest   api.GetBaseRatesForPeriodRequest
	reqSetDefaultMarkupRequest        api.SetDefaultMarkupRequest
	reqGetDefaultMarkupRequest        api.GetDefaultMarkupRequest
	reqGetDefaultMarkupHistoryRequest api.GetDefaultMarkupHistoryRequest
	reqSetMarkupRequest               api.SetMarkupRequest
	reqGetMarkupRequest               api.GetMarkupRequest
	reqGetMarkupHistoryRequest        api.GetMarkupHistoryRequest
	reqGetRateRequest                 api.GetRateRequest
	reqAddPackageRequest              api.AddPackageRequest
	reqGetPackageByOrgRequest         api.GetPackageByOrgRequest
	reqPackagesRequest                api.PackagesRequest
	reqUpdatePackageRequest           api.UpdatePackageRequest

	/* API Responses */

}

func init() {
	log.SetLevel(log.TraceLevel)
	log.SetOutput(os.Stderr)
}

func InitializeData() *InitData {
	d := &InitData{}
	d.Host = "http://192.168.0.22:8074"
	d.MbHost = "amqp://guest:guest@192.168.0.22:5672/"
	d.Sys = NewDataPlanSys(d.Host)
	d.SimType = "ukama_data"
	d.reqUploadBaseRatesRequest = api.UploadBaseRatesRequest{
		EffectiveAt: utils.GenerateFutureDate(5 * time.Second),
		FileURL:     "https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/template/template.csv",
		EndAt:       utils.GenerateFutureDate(365 * 24 * time.Hour),
		SimType:     d.SimType,
	}

	d.OwnerId = uuid.NewV4().String()
	d.OrgId = uuid.NewV4().String()
	d.BaseRateId = make([]string, 8)
	d.Countries = []string{"The Lunar Maria", "Montes Appenninus", "Tycho crater"}
	d.Providers = make(map[string][]string)
	d.Providers[d.Countries[0]] = []string{"ABC Tel", "Light Tel", "Eagle Tel"}
	d.Providers[d.Countries[1]] = []string{"Power Tel", "2D Tel"}
	d.Providers[d.Countries[2]] = []string{"Multi Tel", "Connect Tel", "OWS Tel"}

	/* Read this one form response of Upload Base rates */
	d.reqGetBaseRateRequest = api.GetBaseRateRequest{}

	/* Read Country name form the  Upload base rate */
	c := d.Countries[utils.RandomInt(2)]
	p := d.Providers[c]
	d.reqGetBaseRatesByCountryRequest = api.GetBaseRatesByCountryRequest{
		Country:  c,
		Provider: p[utils.RandomInt(len(p)-1)],
		SimType:  d.SimType,
	}

	d.reqGetBaseRatesForPeriodRequest = api.GetBaseRatesForPeriodRequest{
		Country:  c,
		Provider: p[utils.RandomInt(len(p)-1)],
		SimType:  d.SimType,
	}

	d.reqSetDefaultMarkupRequest = api.SetDefaultMarkupRequest{
		Markup: float64(utils.RandomInt(50)),
	}

	d.reqGetDefaultMarkupRequest = api.GetDefaultMarkupRequest{}
	d.reqGetDefaultMarkupHistoryRequest = api.GetDefaultMarkupHistoryRequest{}
	d.reqSetMarkupRequest = api.SetMarkupRequest{
		OwnerId: d.OwnerId,
		Markup:  float64(utils.RandomInt(50)),
	}
	d.reqGetMarkupRequest = api.GetMarkupRequest{
		OwnerId: d.OwnerId,
	}
	d.reqGetMarkupHistoryRequest = api.GetMarkupHistoryRequest{
		OwnerId: d.OwnerId,
	}

	d.reqGetRateRequest = api.GetRateRequest{
		OwnerId:  d.OwnerId,
		Country:  c,
		Provider: p[utils.RandomInt(len(p)-1)],
		SimType:  d.SimType,
		From:     utils.GenerateFutureDate(24 * time.Hour),
		To:       utils.GenerateFutureDate(30 * 24 * time.Hour),
	}

	d.reqAddPackageRequest = api.AddPackageRequest{
		OwnerId:    d.OwnerId,
		OrgId:      d.OrgId,
		Name:       faker.FirstName() + "-monthly-pack",
		SimType:    d.SimType,
		From:       utils.GenerateFutureDate(24 * time.Hour),
		To:         utils.GenerateFutureDate(30 * 24 * time.Hour),
		BaserateId: "",
		SmsVolume:  100,
		DataVolume: 1024,
		DataUnit:   "MegaBytes",
		Type:       "postpaid",
		Active:     true,
		Flatrate:   false,
		Apn:        "ukama.tel",
	}

	d.reqGetPackageByOrgRequest = api.GetPackageByOrgRequest{
		OrgId: d.OrgId,
	}

	d.reqPackagesRequest = api.PackagesRequest{
		Uuid: d.OwnerId,
	}

	d.reqUpdatePackageRequest = api.UpdatePackageRequest{
		Uuid:   d.OwnerId,
		Name:   faker.FirstName(),
		Active: false,
	}
	return d
}

func TestWorkflow_DataPlanSystem(t *testing.T) {

	/* Sim pool */
	w := test.NewWorkflow("dataplan_workflow_1", "Adding rates and packages")

	w.SetUpFxn = func(ctx context.Context, w *test.Workflow) error {
		log.Tracef("Initilizing Data for %s.", w.String())
		w.Data = InitializeData()

		log.Tracef("Workflow Data : %+v", w.Data)
		return nil
	}

	w.RegisterTestCase(&test.TestCase{
		Name:        "Adding base rate",
		Description: "Add base rate provided by third parties",
		Data:        &bpb.UploadBaseRatesResponse{},
		Workflow:    w,
		SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			a := tc.GetWorkflowData().(*InitData)
			log.Tracef("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.baserate.rate.update"})
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*InitData)
			if ok {
				tc.Data, err = a.Sys.DataPlanBaseRateUpload(a.reqUploadBaseRatesRequest)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false

			resp := tc.GetData().(*bpb.UploadBaseRatesResponse)
			if assert.NotNil(t, resp) {
				assert.Equal(t, true, tc.Watcher.Expections())
				check = true
			}

			return check, nil
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Here we save any data required to be saved from the test case
			Cleanup any test specific data
			*/
			resp := tc.GetData().(*bpb.UploadBaseRatesResponse)

			a := tc.GetWorkflowData().(*InitData)
			for _, r := range resp.Rate {
				a.BaseRateId = append(a.BaseRateId, r.Uuid)
			}

			tc.SaveWorkflowData(a)
			tc.Watcher.Stop()
			return nil
		},
	})

	w.RegisterTestCase(&test.TestCase{
		Name:        "Get Base rate",
		Description: "Get Base rate by Id",
		Data:        &bpb.GetBaseRatesByIdResponse{},
		Workflow:    w,
		SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			a := tc.GetWorkflowData().(*InitData)
			a.reqGetBaseRateRequest = api.GetBaseRateRequest{
				RateId: a.BaseRateId[len(a.BaseRateId)-1],
			}
			tc.SaveWorkflowData(a)
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*InitData)
			if ok {
				tc.Data, err = a.Sys.DataPlanBaseRateGet(a.reqGetBaseRateRequest)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false

			resp := tc.GetData().(*bpb.GetBaseRatesByIdResponse)
			if assert.NotNil(t, resp) {
				data := tc.GetWorkflowData().(*InitData)
				assert.Equal(t, data.reqGetBaseRateRequest.RateId, resp.Rate.Uuid)
				check = true
			}

			return check, nil
		},
	})

	w.RegisterTestCase(&test.TestCase{
		Name:        "Get Base rates for country",
		Description: "Get base rates for country",
		Data:        &bpb.GetBaseRatesResponse{},
		Workflow:    w,
		SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			a := tc.GetWorkflowData().(*InitData)
			c := a.Countries[len(a.Countries)-1]
			p := a.Providers[c]
			a.reqGetBaseRatesByCountryRequest = api.GetBaseRatesByCountryRequest{
				Country:  c,
				Provider: p[len(p)-1],
				SimType:  a.SimType,
			}
			tc.SaveWorkflowData(a)
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*InitData)
			if ok {
				tc.Data, err = a.Sys.DataPlanBaseRateGetByCountry(a.reqGetBaseRatesByCountryRequest)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false

			resp := tc.GetData().(*bpb.GetBaseRatesResponse)
			if assert.NotNil(t, resp) {
				data := tc.GetWorkflowData().(*InitData)
				assert.Equal(t, data.reqGetBaseRatesByCountryRequest.Country, resp.Rates[0].Country)
				assert.Equal(t, data.reqGetBaseRatesByCountryRequest.Provider, resp.Rates[0].Provider)
				check = true
			}

			return check, nil
		},
	})

	/* Get rates by Period */
	w.RegisterTestCase(&test.TestCase{
		Name:        "Get base rate for period",
		Description: "Get base rate for a period",
		Data:        &bpb.GetBaseRatesResponse{},
		Workflow:    w,
		SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			a := tc.GetWorkflowData().(*InitData)
			c := a.Countries[len(a.Countries)-1]
			p := a.Providers[c]
			a.reqGetBaseRatesForPeriodRequest = api.GetBaseRatesForPeriodRequest{
				Country:  c,
				Provider: p[len(p)-1],
				SimType:  a.SimType,
				From:     utils.GenerateFutureDate(24 * time.Hour),
				To:       utils.GenerateFutureDate(30 * 24 * time.Hour),
			}
			tc.SaveWorkflowData(a)
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*InitData)
			if ok {
				tc.Data, err = a.Sys.DataPlanBaseRateGetByPeriod(a.reqGetBaseRatesForPeriodRequest)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false

			resp := tc.GetData().(*bpb.GetBaseRatesResponse)
			if assert.NotNil(t, resp) {
				if assert.Condition(t, func() bool { return (len(resp.Rates) > 0) }, "somebase rates should be returned.") {
					data := tc.GetWorkflowData().(*InitData)
					assert.Equal(t, data.reqGetBaseRatesByCountryRequest.Country, resp.Rates[0].Country)
					assert.Equal(t, data.reqGetBaseRatesByCountryRequest.Provider, resp.Rates[0].Provider)
					check = true
				}
			}

			return check, nil
		},
	})

	// Add Mark ups
	w.RegisterTestCase(&test.TestCase{
		Name:        "Set Markup",
		Description: "Add markup rate fpr owner",
		Data:        &rpb.UpdateMarkupResponse{},
		Workflow:    w,
		SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			a := tc.GetWorkflowData().(*InitData)
			log.Tracef("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.rate.markup.update"})
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*InitData)
			if ok {
				tc.Data, err = a.Sys.DataPlanUpdateMarkup(a.reqSetMarkupRequest)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false

			resp := tc.GetData().(*rpb.UpdateMarkupResponse)
			if assert.NotNil(t, resp) {
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

	/* Get Mark up */
	w.RegisterTestCase(&test.TestCase{
		Name:        "Get markup",
		Description: "Get markup percentage for the owner",
		Data:        &rpb.GetMarkupResponse{},
		Workflow:    w,

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*InitData)
			if ok {
				tc.Data, err = a.Sys.DataPlanGetUserMarkup(a.reqGetMarkupRequest)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false

			resp := tc.GetData().(*rpb.GetMarkupResponse)
			if assert.NotNil(t, resp) {
				data := tc.GetWorkflowData().(*InitData)
				assert.Equal(t, data.reqGetMarkupRequest.OwnerId, resp.OwnerId)
				check = true
			}

			return check, nil
		},
	})

	/* Get rate */
	w.RegisterTestCase(&test.TestCase{
		Name:        "Get rate for Owner's org",
		Description: "Get rate for a Owner's org",
		Data:        &rpb.GetRateResponse{},
		Workflow:    w,
		SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			a := tc.GetWorkflowData().(*InitData)
			c := a.Countries[len(a.Countries)-1]
			p := a.Providers[c]
			a.reqGetRateRequest = api.GetRateRequest{
				OwnerId:  a.OwnerId,
				Country:  c,
				Provider: p[len(p)-1],
				SimType:  a.SimType,
				From:     utils.GenerateFutureDate(24 * time.Hour),
				To:       utils.GenerateFutureDate(30 * 24 * time.Hour),
			}
			tc.SaveWorkflowData(a)
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*InitData)
			if ok {
				tc.Data, err = a.Sys.DataPlanGetRate(a.reqGetRateRequest)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false

			resp := tc.GetData().(*rpb.GetRateResponse)
			if assert.NotNil(t, resp) {
				data := tc.GetWorkflowData().(*InitData)
				assert.Equal(t, data.reqGetBaseRatesByCountryRequest.Country, resp.Rates[0].Country)
				assert.Equal(t, data.reqGetBaseRatesByCountryRequest.Provider, resp.Rates[0].Provider)
				check = true
			}

			return check, nil
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Here we save any data required to be saved from the test case
			Cleanup any test specific data
			*/

			if tc.State == test.StateTypePass {
				resp := tc.GetData().(*rpb.GetRateResponse)

				a := tc.GetWorkflowData().(*InitData)
				if len(resp.Rates) > 0 {
					a.reqAddPackageRequest.BaserateId = resp.Rates[0].Uuid
				}

				tc.SaveWorkflowData(a)
			}
			return nil
		},
	})

	/* Add a package */
	w.RegisterTestCase(&test.TestCase{
		Name:        "Create a package",
		Description: "Cretae package",
		Data:        &ppb.AddPackageResponse{},
		Workflow:    w,
		SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			a := tc.GetWorkflowData().(*InitData)
			log.Tracef("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.package.package.create"})
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*InitData)
			if ok {
				tc.Data, err = a.Sys.DataPlanPackageAdd(a.reqAddPackageRequest)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false

			resp := tc.GetData().(*ppb.AddPackageResponse)
			if assert.NotNil(t, resp) {
				data := tc.GetWorkflowData().(*InitData)
				assert.Equal(t, data.reqAddPackageRequest.OrgId, resp.Package.OrgId)
				assert.NotNil(t, resp.Package.Uuid)
				assert.Equal(t, true, tc.Watcher.Expections())
				check = true
			}

			return check, nil
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Here we save any data required to be saved from the test case
			Cleanup any test specific data
			*/
			resp := tc.GetData().(*ppb.AddPackageResponse)

			a := tc.GetWorkflowData().(*InitData)
			a.PackageId = resp.Package.Uuid

			tc.SaveWorkflowData(a)
			tc.Watcher.Stop()
			return nil
		},
	})

	/* Get Packages */
	w.RegisterTestCase(&test.TestCase{
		Name:        "Get packages for org",
		Description: "Get packages for the organization",
		Data:        &ppb.GetByOrgPackageResponse{},
		Workflow:    w,
		SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			a := tc.GetWorkflowData().(*InitData)
			a.reqGetPackageByOrgRequest = api.GetPackageByOrgRequest{
				OrgId: a.OrgId,
			}
			tc.SaveWorkflowData(a)
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			/* Test Case */
			var err error
			a, ok := tc.GetWorkflowData().(*InitData)
			if ok {
				tc.Data, err = a.Sys.DataPlanPackageGetByOrg(a.reqGetPackageByOrgRequest)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			/* Check for possible failures during test case */
			check := false

			resp := tc.GetData().(*ppb.GetByOrgPackageResponse)
			if assert.NotNil(t, resp) {
				if assert.Condition(t, func() bool { return (len(resp.Packages) > 0) }) {
					data := tc.GetWorkflowData().(*InitData)
					assert.Equal(t, data.reqGetPackageByOrgRequest.OrgId, resp.Packages[0].OrgId)
					check = true
				}
			}

			return check, nil
		},
	})

	/* Run */
	err := w.Run(t, context.Background())
	assert.NoError(t, err)

	w.Status()

}
