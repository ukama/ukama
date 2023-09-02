package nucleus

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v4"
	pkg "github.com/ukama/ukama/testing/integration/pkg"
	"github.com/ukama/ukama/testing/integration/pkg/test"
	"github.com/ukama/ukama/testing/integration/pkg/utils"

	log "github.com/sirupsen/logrus"
	napi "github.com/ukama/ukama/systems/nucleus/api-gateway/pkg/rest"
	orgpb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"
	userpb "github.com/ukama/ukama/systems/nucleus/user/pb/gen"
)

type NucleusData struct {
	Name          string
	Email         string
	Phone         string
	AuthId        string
	OwnerId       string
	OrgId         string
	OrgName       string
	NucleusClient *NucleusClient
	Host          string
	MbHost        string

	// API requests
	reqAddUser napi.AddUserRequest
	reqAddOrg  napi.AddOrgRequest
	reqGetUser napi.GetUserRequest
	reqGetOrg  napi.GetOrgRequest
}

type TestData struct {
	Input   *NucleusData
	Results interface{}
}

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stderr)
}

func InitializeData() *NucleusData {
	config := pkg.NewConfig()
	d := &NucleusData{}

	d.Name = strings.ToLower(faker.FirstName())
	d.Email = strings.ToLower(faker.Email())
	d.Phone = strings.ToLower(faker.Phonenumber())
	d.AuthId = strings.ToLower(faker.UUIDHyphenated())

	d.OrgName = strings.ToLower(faker.FirstName()) + "-org"

	d.Host = config.System.Nucleus
	d.NucleusClient = NewNucleusClient(d.Host)
	d.MbHost = config.System.MessageBus

	d.reqAddUser = napi.AddUserRequest{
		Name:   d.Name,
		Email:  d.Email,
		Phone:  d.Phone,
		AuthId: d.AuthId,
	}

	d.reqGetUser = napi.GetUserRequest{
		UserId: d.OwnerId,
	}

	d.reqAddOrg = napi.AddOrgRequest{
		OrgName: d.OrgName,
	}

	return d
}

var TC_nucleus_add_user = &test.TestCase{
	Name:        "Add User",
	Description: "Register a new user to the system",
	Data:        &userpb.AddResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a, ok := tc.GetWorkflowData().(*NucleusData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		log.Debugf("Setting up watcher for %s", tc.Name)
		tc.Watcher = utils.SetupWatcher(a.MbHost,
			[]string{"event.cloud.users.user.add"})

		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		// Test Case
		var err error

		td, ok := tc.GetWorkflowData().(*NucleusData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		// make sure the owner or admin is the request executor
		tc.Data, err = td.NucleusClient.AddUser(td.reqAddUser)

		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*userpb.AddResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*NucleusData)
			if data.reqAddUser.Email == resp.User.Email &&
				resp.User.Id != "" {
				check = true
			}
		}
		return check, nil
	},

	ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
		// Here we save any data required to be saved from the
		// test case Cleanup any test specific data

		resp, ok := tc.GetData().(*userpb.AddResponse)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a, ok := tc.GetWorkflowData().(*NucleusData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a.OwnerId = resp.User.Id

		tc.SaveWorkflowData(a)
		log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
		tc.Watcher.Stop()

		return nil
	},
}

var TC_nucleus_add_org = &test.TestCase{
	Name:        "Add org",
	Description: "Add an organization",
	Data:        &orgpb.AddResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a, ok := tc.GetWorkflowData().(*NucleusData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a.reqAddOrg = napi.AddOrgRequest{
			OrgName:     a.OrgName,
			Owner:       a.OwnerId,
			Certificate: "-----BEGIN CERTIFICATE-----",
		}

		log.Debugf("Setting up watcher for %s", tc.Name)
		tc.Watcher = utils.SetupWatcher(a.MbHost,
			[]string{"event.cloud.org.org.add"})

		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		// Test Case
		var err error

		td, ok := tc.GetWorkflowData().(*NucleusData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		// make sure the owner or admin is the request executor
		tc.Data, err = td.NucleusClient.AddOrg(td.reqAddOrg)

		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*orgpb.AddResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*NucleusData)
			if data.reqAddOrg.OrgName == resp.Org.Name &&
				resp.Org.Id != "" {
				check = true
			}
		}
		return check, nil
	},

	ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
		// Here we save any data required to be saved from the
		// test case Cleanup any test specific data

		resp, ok := tc.GetData().(*orgpb.AddResponse)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a, ok := tc.GetWorkflowData().(*NucleusData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a.OrgId = resp.Org.Id

		tc.SaveWorkflowData(a)
		log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
		tc.Watcher.Stop()

		return nil
	},
}

var TC_nucleus_get_org = &test.TestCase{
	Name:        "Get Org",
	Description: "Get Org by name",
	Data:        &orgpb.GetByNameRequest{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*NucleusData)
		a.reqGetOrg = napi.GetOrgRequest{
			OrgName: a.OrgName,
		}
		tc.SaveWorkflowData(a)
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*NucleusData)
		if ok {
			tc.Data, err = a.NucleusClient.GetOrg(a.reqGetOrg)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		/* Check for possible failures during test case */
		check := false

		resp := tc.GetData().(*orgpb.GetByNameRequest)
		if resp != nil {
			data := tc.GetWorkflowData().(*NucleusData)
			if data.reqGetOrg.OrgName == resp.Name {
				check = true
			}
		}

		return check, nil
	},
}

var TC_nucleus_get_user = &test.TestCase{
	Name:        "Get User",
	Description: "Get User By Id",
	Data:        &userpb.GetRequest{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		/* Setup required for test case
		Initialize any test specific data if required
		*/
		a := tc.GetWorkflowData().(*NucleusData)
		a.reqGetUser = napi.GetUserRequest{
			UserId: a.OwnerId,
		}
		tc.SaveWorkflowData(a)
		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		/* Test Case */
		var err error
		a, ok := tc.GetWorkflowData().(*NucleusData)
		if ok {
			tc.Data, err = a.NucleusClient.GetUser(a.reqGetUser)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		/* Check for possible failures during test case */
		check := false

		resp := tc.GetData().(*userpb.GetResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*NucleusData)
			if data.reqGetUser.UserId == resp.User.Id {
				check = true
			}
		}

		return check, nil
	},
}
