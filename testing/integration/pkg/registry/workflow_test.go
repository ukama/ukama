package registry

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/testing/integration/pkg"
	"github.com/ukama/ukama/testing/integration/pkg/nucleus"
	"github.com/ukama/ukama/testing/integration/pkg/test"
	"github.com/ukama/ukama/testing/integration/pkg/utils"

	log "github.com/sirupsen/logrus"
	napi "github.com/ukama/ukama/systems/nucleus/api-gateway/pkg/rest"
	api "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"
	mempb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
)

var config *pkg.Config

type RegistryData struct {
	AuthId         string
	Name           string
	Email          string
	Phone          string
	OwnerId        string
	MemberId       string
	OrgId          string
	OrgName        string
	NetName        string
	RegistryClient *RegistryClient
	Nuc            *nucleus.NucleusClient
	Host           string
	MbHost         string

	// API requests
	reqAddUser    napi.AddUserRequest
	reqAddMember  api.MemberRequest
	reqGetMember  api.GetMemberRequest
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
	config = pkg.NewConfig()
	d := &RegistryData{}

	d.Name = strings.ToLower(faker.FirstName())
	d.Email = strings.ToLower(faker.Email())
	d.Phone = strings.ToLower(faker.Phonenumber())
	d.AuthId = strings.ToLower(faker.UUIDHyphenated())
	d.OrgId = "8c6c2bec-5f90-4fee-8ffd-ee6456abf4fc"
	d.OrgName = "ukama-test-org"

	d.Host = config.System.Registry
	d.MbHost = config.System.MessageBus
	d.RegistryClient = NewRegistryClient(d.Host)
	d.Nuc = nucleus.NewNucleusClient(config.System.Nucleus)
	d.reqAddUser = napi.AddUserRequest{
		Name:   d.Name,
		Email:  d.Email,
		Phone:  d.Phone,
		AuthId: d.AuthId,
	}

	d.reqAddMember = api.MemberRequest{
		UserUuid: d.OwnerId,
		Role:     "admin",
	}

	d.reqAddNetwork = api.AddNetworkRequest{
		OrgName: d.OrgName,
	}

	return d
}

func TestWorkflow_RegistrySystem(t *testing.T) {
	w := test.NewWorkflow("Registry Workflows", "Various use cases whille adding registry items")

	w.SetUpFxn = func(t *testing.T, ctx context.Context, w *test.Workflow) error {
		log.Debugf("Initilizing Data for %s.", w.String())
		var err error

		d := InitializeData()

		w.Data = d

		log.Debugf("Workflow Data : %+v", w.Data)

		return err
	}

	w.RegisterTestCase(&test.TestCase{
		Name:        "Get member",
		Description: "Get member from default org",
		Data:        &mempb.MemberResponse{},
		Workflow:    w,
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			// Setup required for test case Initialize any
			// test specific data if required
			a, ok := tc.GetWorkflowData().(*RegistryData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			// we need to setup a new user not member of the org

			res, err := a.Nuc.AddUser(a.reqAddUser)

			if assert.NoError(t, err) {
				assert.NotNil(t, res)
				assert.Equal(t, a.reqAddUser.Name, res.User.Name)
				assert.Equal(t, a.reqAddUser.Email, res.User.Email)
				assert.Equal(t, a.reqAddUser.Phone, res.User.Phone)
			}

			a.reqGetMember = api.GetMemberRequest{
				UserUuid: res.User.Id,
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

			//  make sure the owner or admin is the request executor
			tc.Data, err = td.RegistryClient.GetMember(td.reqGetMember)

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

			tr, ok := tc.GetData().(*mempb.MemberResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return check, fmt.Errorf("invalid data type for Workflow data")
			}

			if assert.NotNil(t, tr) {
				// TODO: assert.Equal(t, d.MemberId, tr.Member.Uuid)
				assert.Equal(t, d.OrgId, tr.Member.OrgId)
				// assert.Equal(t, true, tc.Watcher.Expections())
				check = true
			}

			return check, nil
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			// Here we save any data required to be saved from the
			// test case Cleanup any test specific data

			resp, ok := tc.GetData().(*mempb.MemberResponse)
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
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
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

			// make sure the owner or admin is the request executor
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
