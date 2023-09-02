package nucleus

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/assert"
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

	d.reqAddOrg = napi.AddOrgRequest{
		OrgName: d.OrgName,
	}

	return d
}

func TestWorkflow_NucleusSystem(t *testing.T) {
	w := test.NewWorkflow("Nucleus Workflows", "Various use cases whille adding nucleus items")

	w.SetUpFxn = func(t *testing.T, ctx context.Context, w *test.Workflow) error {
		log.Debugf("Initilizing Data for %s.", w.String())
		var err error

		d := InitializeData()

		w.Data = d

		log.Debugf("Workflow Data : %+v", w.Data)

		return err
	}

	w.RegisterTestCase(&test.TestCase{
		Name:        "Add User",
		Description: "Register a new user to the system",
		Data:        &userpb.AddResponse{},
		Workflow:    w,
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

			d, ok := tc.GetWorkflowData().(*NucleusData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return check, fmt.Errorf("invalid data type for Workflow data")
			}

			tr, ok := tc.GetData().(*userpb.AddResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return check, fmt.Errorf("invalid data type for Workflow data")
			}

			if assert.NotNil(t, tr) {
				assert.Equal(t, d.Name, tr.User.Name)
				assert.Equal(t, d.Email, tr.User.Email)
				assert.Equal(t, d.Phone, tr.User.Phone)
				assert.Equal(t, true, tc.Watcher.Expections())
				check = true
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
	})

	w.RegisterTestCase(&test.TestCase{
		Name:        "Add org",
		Description: "Add an organization",
		Data:        &orgpb.AddResponse{},
		Workflow:    w,
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

			d, ok := tc.GetWorkflowData().(*NucleusData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return check, fmt.Errorf("invalid data type for Workflow data")
			}

			tr, ok := tc.GetData().(*orgpb.AddResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return check, fmt.Errorf("invalid data type for Workflow data")
			}

			if assert.NotNil(t, tr) {
				assert.Equal(t, d.OrgName, tr.Org.Name)
				assert.Equal(t, true, tc.Watcher.Expections())
				check = true
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
	})

	err := w.Run(t, context.Background())
	assert.NoError(t, err)

	w.Status()
}
