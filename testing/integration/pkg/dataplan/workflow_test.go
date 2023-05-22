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

	api "github.com/ukama/ukama/systems/data-plan/api-gateway/pkg/rest"
	//rpb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	//ppb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
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

	// This data is taken from the https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/template/template.csv */
	Providers map[string][]string
	Countries []string

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
	d.Host = "http://192.168.0.22:8078"
	d.MbHost = "amqp://guest:guest@192.168.0.22:5672/"
	d.Sys = NewDataPlanSys(d.Host)
	d.SimType = "ukama_data"
	d.reqUploadBaseRatesRequest = api.UploadBaseRatesRequest{
		EffectiveAt: utils.GenerateFutureDate(1 * time.Minute),
		FileURL:     "https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/template/template.csv",
		EndAt:       utils.GenerateFutureDate(365 * 24 * time.Hour),
		SimType:     d.SimType,
	}

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
		OwnerId: "",
		Markup:  float64(utils.RandomInt(50)),
	}
	d.reqGetMarkupRequest = api.GetMarkupRequest{}
	d.reqGetMarkupHistoryRequest = api.GetMarkupHistoryRequest{}
	d.reqGetRateRequest = api.GetRateRequest{
		OwnerId:  "",
		Country:  c,
		Provider: p[utils.RandomInt(len(p)-1)],
		SimType:  d.SimType,
		From:     utils.GenerateFutureDate(24 * time.Hour),
		To:       utils.GenerateFutureDate(30 * 24 * time.Hour),
	}

	d.reqAddPackageRequest = api.AddPackageRequest{
		OwnerId:    "",
		OrgId:      "",
		Name:       faker.FirstName(),
		SimType:    d.SimType,
		From:       utils.GenerateFutureDate(24 * time.Hour),
		To:         utils.GenerateFutureDate(30 * 24 * time.Hour),
		BaserateId: "",
		SmsVolume:  100,
		DataVolume: 1024,
		DataUnit:   "Mb",
		Type:       "postpaid",
		Active:     true,
	}

	d.reqGetPackageByOrgRequest = api.GetPackageByOrgRequest{
		OrgId: "",
	}

	d.reqPackagesRequest = api.PackagesRequest{
		Uuid: "",
	}

	d.reqUpdatePackageRequest = api.UpdatePackageRequest{
		Uuid:   "",
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
		Data:        &bpb.GetBaseRatesResponse{},
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

			resp := tc.GetData().(*bpb.GetBaseRatesResponse)
			if assert.NotNil(t, resp) {
				data := tc.GetWorkflowData().(*InitData)
				assert.Equal(t, data.reqGetBaseRateRequest.RateId, resp.Rates[0].Uuid)
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
		Description: "Get bae rate for a period",
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
				data := tc.GetWorkflowData().(*InitData)
				assert.Equal(t, data.reqGetBaseRatesByCountryRequest.Country, resp.Rates[0].Country)
				assert.Equal(t, data.reqGetBaseRatesByCountryRequest.Provider, resp.Rates[0].Provider)
				check = true
			}

			return check, nil
		},
	})

	/* Run */
	err := w.Run(t, context.Background())
	assert.NoError(t, err)

	w.Status()

}
