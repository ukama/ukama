package registry_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/assert"
	pkg "github.com/ukama/ukama/testing/integration/pkg/registry"
	"github.com/ukama/ukama/testing/integration/pkg/test"
	"github.com/ukama/ukama/testing/integration/pkg/utils"

	log "github.com/sirupsen/logrus"
	api "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	userpb "github.com/ukama/ukama/systems/registry/users/pb/gen"
)

type RegistryData struct {
	Name           string
	Email          string
	Phone          string
	OwnerId        string
	MemberId       string
	OrgId          string
	OrgName        string
	NetName        string
	RegistryClient *pkg.RegistryClient
	Host           string
	MbHost         string

	// API requests
	reqAddUser    api.AddUserRequest
	reqAddOrg     api.AddOrgRequest
	reqAddMember  api.MemberRequest
	reqAddNetwork api.AddNetworkRequest
}

type TestData struct {
	Input   *RegistryData
	Results interface{}
}

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stderr)
}

func InitializeData() *RegistryData {
	d := &RegistryData{}

	d.Name = strings.ToLower(faker.FirstName())
	d.Email = strings.ToLower(faker.Email())
	d.Phone = strings.ToLower(faker.Phonenumber())

	d.OrgName = strings.ToLower(faker.FirstName()) + "-org"

	d.Host = "http://localhost:8082"
	d.RegistryClient = pkg.NewRegistryClient(d.Host)
	d.MbHost = "amqp://guest:guest@localhost:5672/"

	d.reqAddUser = api.AddUserRequest{
		Name:  d.Name,
		Email: d.Email,
		Phone: d.Phone,
	}

	d.reqAddOrg = api.AddOrgRequest{
		OrgName: d.OrgName,
	}

	d.reqAddMember = api.MemberRequest{
		OrgName: d.OrgName,
	}

	d.reqAddNetwork = api.AddNetworkRequest{
		OrgName: d.OrgName,
	}

	return d
}

func TestWorkflow_RegistrySystem(t *testing.T) {
	w := test.NewWorkflow("Registry Workflows", "Various use cases whille adding registry items")

	w.SetUpFxn = func(ctx context.Context, w *test.Workflow) error {
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
		SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
			// Setup required for test case Initialize any
			// test specific data if required
			a, ok := tc.GetWorkflowData().(*RegistryData)
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

			td, ok := tc.GetWorkflowData().(*RegistryData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			// // make sure the owner or admin is the request executor
			tc.Data, err = td.RegistryClient.AddUser(td.reqAddUser)

			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			// Check for possible failures during test case
			check := false

			d, ok := tc.GetWorkflowData().(*RegistryData)
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

			a, ok := tc.GetWorkflowData().(*RegistryData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a.OwnerId = resp.User.Uuid

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
		SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
			// Setup required for test case Initialize any
			// test specific data if required
			a, ok := tc.GetWorkflowData().(*RegistryData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a.reqAddOrg = api.AddOrgRequest{
				OrgName: a.OrgName,
				Owner:   a.OwnerId,
			}

			log.Debugf("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost,
				[]string{"event.cloud.org.org.add"})

			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			// Test Case
			var err error

			td, ok := tc.GetWorkflowData().(*RegistryData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			// // make sure the owner or admin is the request executor
			tc.Data, err = td.RegistryClient.AddOrg(td.reqAddOrg)

			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			// Check for possible failures during test case
			check := false

			d, ok := tc.GetWorkflowData().(*RegistryData)
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

			a, ok := tc.GetWorkflowData().(*RegistryData)
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

	w.RegisterTestCase(&test.TestCase{
		Name:        "Add member",
		Description: "Add a user to an organization",
		Data:        &orgpb.MemberResponse{},
		Workflow:    w,
		SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
			// Setup required for test case Initialize any
			// test specific data if required
			a, ok := tc.GetWorkflowData().(*RegistryData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			// we need to setup a new user not member of the org

			name := strings.ToLower(faker.FirstName())
			email := strings.ToLower(faker.Email())
			phone := strings.ToLower(faker.Phonenumber())

			res, err := a.RegistryClient.AddUser(api.AddUserRequest{
				Name:  name,
				Email: email,
				Phone: phone,
			})

			if assert.NoError(t, err) {
				assert.NotNil(t, res)
				assert.Equal(t, name, res.User.Name)
				assert.Equal(t, email, res.User.Email)
				assert.Equal(t, phone, res.User.Phone)
			}

			a.MemberId = res.User.Uuid

			a.reqAddMember = api.MemberRequest{
				OrgName:  a.OrgName,
				UserUuid: a.MemberId,
			}

			log.Debugf("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost,
				[]string{"event.cloud.org.member.add"})

			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			// Test Case
			var err error

			td, ok := tc.GetWorkflowData().(*RegistryData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			// // make sure the owner or admin is the request executor
			tc.Data, err = td.RegistryClient.AddMember(td.reqAddMember)

			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			// Check for possible failures during test case
			check := false

			d, ok := tc.GetWorkflowData().(*RegistryData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return check, fmt.Errorf("invalid data type for Workflow data")
			}

			tr, ok := tc.GetData().(*orgpb.MemberResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return check, fmt.Errorf("invalid data type for Workflow data")
			}

			if assert.NotNil(t, tr) {
				assert.Equal(t, d.MemberId, tr.Member.Uuid)
				assert.Equal(t, d.OrgId, tr.Member.OrgId)
				assert.Equal(t, true, tc.Watcher.Expections())
				check = true
			}

			return check, nil
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			// Here we save any data required to be saved from the
			// test case Cleanup any test specific data

			resp, ok := tc.GetData().(*orgpb.MemberResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, ok := tc.GetWorkflowData().(*RegistryData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()

			return nil
		},
	})

	w.RegisterTestCase(&test.TestCase{
		Name:        "Add network",
		Description: "Add network to an organization",
		Data:        &TestData{},
		Workflow:    w,
		SetUpFxn: func(ctx context.Context, tc *test.TestCase) error {
			// Setup required for test case Initialize any
			// test specific data if required
			a, ok := tc.GetWorkflowData().(*RegistryData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a.NetName = strings.ToLower(faker.FirstName()) + "-net"
			a.reqAddNetwork = api.AddNetworkRequest{
				OrgName: a.OrgName,
				NetName: a.NetName,
			}

			tc.Data = &TestData{Input: a}

			log.Debugf("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost,
				[]string{"event.cloud.network.network.add"})

			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			// Test Case
			var err error

			td, ok := tc.GetData().(*TestData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			ti := td.Input

			// // make sure the owner or admin is the request executor
			tr, err := ti.RegistryClient.AddNetwork(ti.reqAddNetwork)
			tc.Data = &TestData{Input: ti, Results: tr}

			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			// Check for possible failures during test case
			check := false

			td, ok := tc.GetData().(*TestData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return check, fmt.Errorf("invalid data type for Workflow data")
			}

			ti := td.Input

			tr, ok := td.Results.(*netpb.AddResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return check, fmt.Errorf("invalid data type for Workflow data")
			}

			if assert.NotNil(t, tr) {
				assert.Equal(t, ti.NetName, tr.Network.Name)
				assert.Equal(t, ti.OrgId, tr.Network.OrgId)
				assert.Equal(t, true, tc.Watcher.Expections())
				check = true
			}

			return check, nil
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			// Here we save any data required to be saved from the
			// test case Cleanup any test specific data

			resp, ok := tc.GetData().(*TestData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, ok := tc.GetWorkflowData().(*RegistryData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

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
