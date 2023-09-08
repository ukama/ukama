package registry

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v4"
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
	UserId         string
	Name           string
	Email          string
	Phone          string
	OwnerId        string
	MemberId       string
	OrgId          string
	OrgName        string
	NetName        string
	NetworkId      string
	RegistryClient *RegistryClient
	Nuc            *nucleus.NucleusClient
	Host           string
	MbHost         string

	// API requests
	reqAddUser      napi.AddUserRequest
	reqAddMember    api.MemberRequest
	reqGetMember    api.GetMemberRequest
	reqUpdateMember api.UpdateMemberRequest
	reqAddNetwork   api.AddNetworkRequest
	reqGetNetwork   api.GetNetworkRequest
	reqGetNetworks  api.GetNetworksRequest
}

func InitializeData() *RegistryData {
	config = pkg.NewConfig()
	d := &RegistryData{}

	d.OrgId = config.OrgId
	d.OrgName = config.OrgName
	d.Name = strings.ToLower(faker.FirstName())
	d.Email = strings.ToLower(faker.Email())
	d.Phone = strings.ToLower(faker.Phonenumber())
	d.AuthId = strings.ToLower(faker.UUIDHyphenated())

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

var TC_registry_get_members = &test.TestCase{
	Name:        "Get members",
	Description: "Get members of default org",
	Data:        &mempb.GetMembersResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required

		a := tc.GetWorkflowData().(*RegistryData)

		tc.SaveWorkflowData(a)
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
		tc.Data, err = td.RegistryClient.GetMembers()

		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*mempb.GetMembersResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			for _, member := range resp.Members {
				if member.OrgId == data.OrgId {
					check = true
				}
			}
		}
		return check, nil
	},
}

var TC_registry_get_member = &test.TestCase{
	Name:        "Get member",
	Description: "Get member from default org",
	Data:        &mempb.MemberResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a, ok := tc.GetWorkflowData().(*RegistryData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		uresp, err := a.Nuc.AddUser(a.reqAddUser)
		if err != nil {
			return err
		} else {
			a.UserId = uresp.User.Id
		}

		a.reqGetMember = api.GetMemberRequest{
			UserUuid: a.UserId,
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

		resp := tc.GetData().(*mempb.MemberResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			if data.reqGetMember.UserUuid == resp.Member.UserId &&
				resp.Member.OrgId != "" {
				check = true
			}
		}
		return check, nil
	},
}

var TC_registry_update_member = &test.TestCase{
	Name:        "Update member",
	Description: "Update member from default org",
	Data:        &mempb.MemberResponse{},

	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a, ok := tc.GetWorkflowData().(*RegistryData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a.NetName = strings.ToLower(faker.FirstName()) + "-net"
		a.reqUpdateMember = api.UpdateMemberRequest{
			UserUuid:      a.UserId,
			IsDeactivated: false,
			Role:          "Vendor",
		}

		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		// Test Case
		var err error

		a, ok := tc.GetWorkflowData().(*RegistryData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		// make sure the owner or admin is the request executor
		if ok {
			err = a.RegistryClient.UpdateMember(a.reqUpdateMember)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},
}

var TC_registry_add_network = &test.TestCase{
	Name:        "Add network",
	Description: "Add network to an organization",
	Data:        &netpb.AddResponse{},

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

		log.Debugf("Setting up watcher for %s", tc.Name)
		tc.Watcher = utils.SetupWatcher(a.MbHost,
			[]string{"event.cloud.network.network.add"})

		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		// Test Case
		var err error

		a, ok := tc.GetWorkflowData().(*RegistryData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		// make sure the owner or admin is the request executor
		if ok {
			tc.Data, err = a.RegistryClient.AddNetwork(a.reqAddNetwork)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {

		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*netpb.AddResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			data.NetworkId = resp.Network.Id
			if data.reqAddNetwork.NetName == resp.Network.Name &&
				data.OrgId == resp.Network.OrgId {
				check = true
			}
		}
		return check, nil
	},
}

var TC_registry_get_networks = &test.TestCase{
	Name:        "Get networks",
	Description: "Get networks of default org",
	Data:        &netpb.GetByOrgResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a := tc.GetWorkflowData().(*RegistryData)
		a.reqGetNetworks = api.GetNetworksRequest{
			OrgUuid: a.OrgId,
		}
		tc.SaveWorkflowData(a)
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

		tc.Data, err = td.RegistryClient.GetNetworks(td.reqGetNetworks)

		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*netpb.GetByOrgResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			if resp.OrgId == data.OrgId {
				check = true
			}
		}
		return check, nil
	},
}

var TC_registry_get_network = &test.TestCase{
	Name:        "Get network",
	Description: "Get network by id",
	Data:        &netpb.GetResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a := tc.GetWorkflowData().(*RegistryData)
		a.reqGetNetwork = api.GetNetworkRequest{
			NetworkId: a.NetworkId,
		}
		tc.SaveWorkflowData(a)
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

		tc.Data, err = td.RegistryClient.GetNetwork(td.reqGetNetwork)

		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*netpb.GetResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			if data.OrgId == resp.Network.OrgId {
				check = true
			}
		}
		return check, nil
	},
}
