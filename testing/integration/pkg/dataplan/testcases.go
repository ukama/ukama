package dataplan

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bxcodec/faker/v4"
	log "github.com/sirupsen/logrus"

	api "github.com/ukama/ukama/systems/data-plan/api-gateway/pkg/rest"
	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	ppb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	rpb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	"github.com/ukama/ukama/testing/integration/pkg"
	"github.com/ukama/ukama/testing/integration/pkg/test"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
)

var config *pkg.Config

type InitData struct {
	Sys          *DataplanClient
	Host         string
	SimType      string `default:"test"`
	MbHost       string
	SubscriberId string
	BaseRateId   []string
	PackageId    string

	// This data is taken from the https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/template/template.csv */
	Providers  map[string][]string
	Countries  []string
	OwnerId    string
	OrgId      string
	BaserateId string
	Country    string
	Provider   string

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
	reqGetPackageByOrgRequest         api.GetPackageByOrgRequest
	reqPackagesRequest                api.PackagesRequest
	reqUpdatePackageRequest           api.UpdatePackageRequest
	ReqAddPackageRequest              api.AddPackageRequest

	/* API Responses */

}

func InitializeData() *InitData {

	config = pkg.NewConfig()

	d := &InitData{}
	d.Host = config.System.Dataplan
	d.MbHost = config.System.MessageBus
	d.Sys = NewDataplanClient(d.Host)
	d.SimType = "test"
	d.OrgId = config.OrgId
	d.OwnerId = config.OrgOwnerId
	d.reqUploadBaseRatesRequest = api.UploadBaseRatesRequest{
		EffectiveAt: utils.GenerateFutureDate(5 * time.Hour),
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
		OwnerId: d.OwnerId,
		Markup:  float64(utils.RandomInt(50)),
	}
	d.reqGetMarkupRequest = api.GetMarkupRequest{
		OwnerId: d.OwnerId,
	}
	d.reqGetMarkupHistoryRequest = api.GetMarkupHistoryRequest{
		OwnerId: d.OwnerId,
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

var TC_dp_add_baserate = &test.TestCase{
	Name:        "Adding base rate",
	Description: "Add base rate provided by third parties",
	Data:        &bpb.UploadBaseRatesResponse{},
	//Workflow:    w,

	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		log.Tracef("Setting up watcher for %s", tc.Name)
		tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.local.dataplan.rate.rate.upload"})
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
		if resp != nil {
			// if tc.Watcher.Expections() {
			check = true
			// } else {
			// 	log.Error("Expected events not found.")
			// }
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
		a.Country = resp.Rate[0].Country
		a.Provider = resp.Rate[0].Provider
		a.BaserateId = resp.Rate[0].Uuid
		tc.SaveWorkflowData(a)
		tc.Watcher.Stop()
		return nil
	},
}

var TC_dp_get_baserate_by_id = &test.TestCase{
	Name:        "Get Base rate",
	Description: "Get Base rate by Id",
	Data:        &bpb.GetBaseRatesByIdResponse{},
	//Workflow:    w,
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
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
		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if data.reqGetBaseRateRequest.RateId == resp.Rate.Uuid {
				check = true
			}
		}
		return check, nil
	},
}

var TC_dp_get_baserate_by_country = &test.TestCase{
	Name:        "Get Base rates for country",
	Description: "Get base rates for country",
	Data:        &bpb.GetBaseRatesResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)

		a.reqGetBaseRatesByCountryRequest = api.GetBaseRatesByCountryRequest{
			Country:  a.Country,
			Provider: a.Provider,
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
		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if data.reqGetBaseRatesByCountryRequest.Country == resp.Rates[0].Country &&
				data.reqGetBaseRatesByCountryRequest.Provider == resp.Rates[0].Provider {
				check = true
			}
		}

		return check, nil
	},
}

var TC_dp_get_baserate_by_period = &test.TestCase{
	Name:        "Get base rate for period",
	Description: "Get base rate for a period",
	Data:        &bpb.GetBaseRatesResponse{},
	//Workflow:    w,
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)

		a.reqGetBaseRatesForPeriodRequest = api.GetBaseRatesForPeriodRequest{
			Country:  a.Country,
			Provider: a.Provider,
			SimType:  a.SimType,
			From:     utils.GenerateFutureDate(24 * time.Hour),
			To:       utils.GenerateFutureDate(30 * 24 * time.Hour),
		}
		log.Info("DATE::: from: ", a.reqGetBaseRatesForPeriodRequest.From, "To:", a.reqGetBaseRatesForPeriodRequest.To)
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
		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if len(resp.Rates) > 0 &&
				data.reqGetBaseRatesForPeriodRequest.Country == resp.Rates[0].Country &&
				data.reqGetBaseRatesForPeriodRequest.Provider == resp.Rates[0].Provider {
				check = true
			}
		}

		return check, nil
	},
}

var TC_dp_add_markup = &test.TestCase{
	Name:        "Set Markup",
	Description: "Add markup rate fpr owner",
	Data:        &rpb.UpdateMarkupResponse{},
	//Workflow:    w,
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		log.Tracef("Setting up watcher for %s", tc.Name)
		tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.local.dataplan.markup.markup.set"})
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
		if resp != nil {
			// if tc.Watcher.Expections() {
			check = true
			// }
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

var TC_dp_get_markup = &test.TestCase{
	Name:        "Get markup",
	Description: "Get markup percentage for the owner",
	Data:        &rpb.GetMarkupResponse{},
	//Workflow:    w,

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
		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if data.reqGetMarkupRequest.OwnerId == resp.OwnerId {
				check = true
			}
		}

		return check, nil
	},
}

var TC_dp_get_rate = &test.TestCase{
	Name:        "Get rate for Owner's org",
	Description: "Get rate for a Owner's org",
	Data:        &rpb.GetRateResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		a.reqGetRateRequest = api.GetRateRequest{
			UserId:   a.OwnerId,
			Country:  a.Country,
			Provider: a.Provider,
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
		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if len(resp.Rates) > 0 &&
				data.reqGetRateRequest.Country == resp.Rates[0].Country &&
				data.reqGetRateRequest.Provider == resp.Rates[0].Provider {
				check = true
			}
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
				a.ReqAddPackageRequest.BaserateId = resp.Rates[0].Uuid
			}

			tc.SaveWorkflowData(a)
		}
		return nil
	},
}

var TC_dp_add_package = &test.TestCase{
	Name:        "Create a package",
	Description: "Create package",
	Data:        &ppb.AddPackageResponse{},
	//Workflow:    w,
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*InitData)
		log.Tracef("Setting up watcher for %s", tc.Name)

		a.ReqAddPackageRequest = api.AddPackageRequest{
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

		tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.local.dataplan.package.update"})
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*InitData)
		if ok {

			tc.Data, err = a.Sys.DataPlanPackageAdd(a.ReqAddPackageRequest)
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
		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if data.ReqAddPackageRequest.OrgId == resp.Package.OrgId &&
				resp.Package.Uuid != "" {
				check = true
			}

		}
		return check, nil
	},

	ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Here we save any data required to be saved from the test case
		Cleanup any test specific data
		*/
		tc.Watcher.Stop()

		if tc.State == test.StateTypePass {
			resp := tc.GetData().(*ppb.AddPackageResponse)
			a := tc.GetWorkflowData().(*InitData)
			a.PackageId = resp.Package.Uuid
			tc.SaveWorkflowData(a)
		}

		return nil
	},
}

var TC_dp_get_package_for_org = &test.TestCase{
	Name:        "Get packages for org",
	Description: "Get packages for the organization",
	Data:        &ppb.GetByOrgPackageResponse{},
	//Workflow:    w,
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
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
		if resp != nil {
			data := tc.GetWorkflowData().(*InitData)
			if len(resp.Packages) > 0 &&
				data.reqGetPackageByOrgRequest.OrgId == resp.Packages[0].OrgId {
				check = true
			}
		}

		return check, nil
	},
}
