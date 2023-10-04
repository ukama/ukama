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
	invpb "github.com/ukama/ukama/systems/registry/invitation/pb/gen"
	mempb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	nodepb "github.com/ukama/ukama/systems/registry/node/pb/gen"
)

var config *pkg.Config

type RegistryData struct {
	AuthId         string
	UserId         string
	Name           string
	Email          string
	Phone          string
	OwnerId        string
	InviteId       string
	NodeId         string
	NodeName       string
	NodeState      string
	MemberId       string
	OrgId          string
	OrgName        string
	NetName        string
	rNodeId        string
	lNodeId        string
	NetworkId      string
	SiteName       string
	SiteId         string
	RegistryClient *RegistryClient
	Nuc            *nucleus.NucleusClient
	Host           string
	MbHost         string

	// API requests
	reqAddUser          napi.AddUserRequest
	reqAddMember        api.MemberRequest
	reqGetMember        api.GetMemberRequest
	reqUpdateMember     api.UpdateMemberRequest
	reqAddNetwork       api.AddNetworkRequest
	reqGetNetwork       api.GetNetworkRequest
	reqGetNetworks      api.GetNetworksRequest
	reqAddSite          api.AddSiteRequest
	reqGetSite          api.GetSiteRequest
	reqGetSites         api.GetNetworkRequest
	reqAddInvite        api.AddInvitationRequest
	reqUpdateInvitation api.UpdateInvitationRequest
	reqGetInvite        api.GetInvitationRequest
	reqGetInvites       api.GetInvitationByOrgRequest
	reqAddNode          api.AddNodeRequest
	reqUpdateNode       api.UpdateNodeRequest
	reqUpdateNodeState  api.UpdateNodeStateRequest
	reqAttachNode       api.AttachNodesRequest
	reqDetachNode       api.DetachNodeRequest
	reqAddNodeToSite    api.AddNodeToSiteRequest
	reqGetNodes         api.GetNodesRequest
	reqGetNode          api.GetNodeRequest
	reqGetNodeForSite   api.GetSiteNodesRequest
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

var TC_registry_add_site = &test.TestCase{
	Name:        "Add site",
	Description: "Add site to network",
	Data:        &netpb.AddSiteResponse{},

	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a, ok := tc.GetWorkflowData().(*RegistryData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a.SiteName = strings.ToLower(faker.FirstName()) + "-site"
		a.reqAddSite = api.AddSiteRequest{
			NetworkId: a.NetworkId,
			SiteName:  a.SiteName,
		}
		tc.SaveWorkflowData(a)
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
			tc.Data, err = a.RegistryClient.AddSite(a.reqAddSite)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {

		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*netpb.AddSiteResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			data.SiteId = resp.Site.Id
			data.SiteName = resp.Site.Name
			if data.SiteName == resp.Site.Name &&
				resp.Site.Id != "" {
				check = true
			}
		}
		return check, nil
	},
}

var TC_registry_get_site = &test.TestCase{
	Name:        "Get site",
	Description: "Get site by id",
	Data:        &netpb.GetSiteResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a := tc.GetWorkflowData().(*RegistryData)
		a.reqGetSite = api.GetSiteRequest{
			NetworkId: a.NetworkId,
			SiteName:  a.SiteName,
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

		tc.Data, err = td.RegistryClient.GetSite(td.reqGetSite)

		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*netpb.GetSiteResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			if data.NetworkId == resp.Site.NetworkId {
				check = true
			}
		}
		return check, nil
	},
}

var TC_registry_get_sites = &test.TestCase{
	Name:        "Get sites",
	Description: "Get sites by network",
	Data:        &netpb.GetSitesByNetworkResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a := tc.GetWorkflowData().(*RegistryData)
		a.reqGetSites = api.GetNetworkRequest{
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

		tc.Data, err = td.RegistryClient.GetSites(td.reqGetSites)

		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*netpb.GetSitesByNetworkResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			if data.NetworkId == resp.NetworkId {
				check = true
			}
		}
		return check, nil
	},
}

var TC_registry_add_invite = &test.TestCase{
	Name:        "Add invite",
	Description: "Add invite",
	Data:        &invpb.AddInvitationResponse{},

	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a, ok := tc.GetWorkflowData().(*RegistryData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a.reqAddInvite = api.AddInvitationRequest{
			Org:   a.OrgName,
			Name:  a.Name,
			Email: a.Email,
			Role:  "admin",
		}
		tc.SaveWorkflowData(a)
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
			tc.Data, err = a.RegistryClient.AddInvitations(a.reqAddInvite)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {

		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*invpb.AddInvitationResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			data.InviteId = resp.Invitation.Id

			if data.OrgName == resp.Invitation.Org &&
				resp.Invitation.Status == invpb.StatusType_Pending {
				check = true
			}
		}
		return check, nil
	},
}

var TC_registry_update_invite = &test.TestCase{
	Name:        "update invite",
	Description: "update invite",
	Data:        &invpb.UpdateInvitationStatusResponse{},

	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a, ok := tc.GetWorkflowData().(*RegistryData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a.reqUpdateInvitation = api.UpdateInvitationRequest{
			InvitationId: a.InviteId,
			Status:       "Accepted",
		}
		tc.SaveWorkflowData(a)
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
			tc.Data, err = a.RegistryClient.UpdateInvitations(a.reqUpdateInvitation)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {

		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*invpb.UpdateInvitationStatusResponse)
		if resp != nil {
			if resp.Status == invpb.StatusType_Accepted {
				check = true
			}
		}
		return check, nil
	},
}

var TC_registry_get_invite = &test.TestCase{
	Name:        "Get invite",
	Description: "Get invite by id",
	Data:        &invpb.GetInvitationResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a := tc.GetWorkflowData().(*RegistryData)
		a.reqGetInvite = api.GetInvitationRequest{
			InvitationId: a.InviteId,
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

		tc.Data, err = td.RegistryClient.GetInvitation(td.reqGetInvite)

		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*invpb.GetInvitationResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			if data.InviteId == resp.Invitation.Id {
				check = true
			}
		}
		return check, nil
	},
}

var TC_registry_get_invites = &test.TestCase{
	Name:        "Get invites",
	Description: "Get invites by org",
	Data:        &invpb.GetInvitationByOrgResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a := tc.GetWorkflowData().(*RegistryData)
		a.reqGetInvites = api.GetInvitationByOrgRequest{
			Org: a.OrgId,
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

		tc.Data, err = td.RegistryClient.GetInvitationByOrg(td.reqGetInvites)

		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*invpb.GetInvitationByOrgResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			for _, invite := range resp.Invitations {
				if invite.Id == data.InviteId {
					check = true
					break
				}
			}

		}
		return check, nil
	},
}

func TC_registry_add_node(typ string) *test.TestCase {
	return &test.TestCase{
		Name:        "Add node",
		Description: "Add Node",
		Data:        &nodepb.AddNodeResponse{},

		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			// Setup required for test case Initialize any
			// test specific data if required
			a, ok := tc.GetWorkflowData().(*RegistryData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}
			a.NodeName = "node" + strings.ToLower(faker.Word())
			var nId = ""
			if typ == "parent" {
				a.NodeId = utils.RandomGetNodeId("tnode")
				nId = a.NodeId
			} else if typ == "left" {
				a.lNodeId = utils.RandomGetNodeId("anode")
				nId = a.lNodeId
			} else if typ == "right" {
				a.rNodeId = utils.RandomGetNodeId("anode")
				nId = a.rNodeId
			}
			a.reqAddNode = api.AddNodeRequest{
				NodeId: nId,
				Name:   a.NodeName,
				OrgId:  a.OrgId,
				State:  "onboarded",
			}
			tc.SaveWorkflowData(a)
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
				tc.Data, err = a.RegistryClient.AddNode(a.reqAddNode)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {

			// Check for possible failures during test case
			check := false

			resp := tc.GetData().(*nodepb.AddNodeResponse)

			if resp != nil {
				data := tc.GetWorkflowData().(*RegistryData)
				if data.OrgId == resp.Node.OrgId {
					check = true
				}
			}
			return check, nil
		},
	}
}

var TC_registry_update_node = &test.TestCase{
	Name:        "Update node",
	Description: "Update Node",
	Data:        &nodepb.UpdateNodeResponse{},

	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a, ok := tc.GetWorkflowData().(*RegistryData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}
		a.NodeName = "update-node" + strings.ToLower(faker.Word())
		a.reqUpdateNode = api.UpdateNodeRequest{
			NodeId: a.NodeId,
			Name:   a.NodeName,
		}
		tc.SaveWorkflowData(a)
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
			tc.Data, err = a.RegistryClient.UpdateNode(a.reqUpdateNode)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {

		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*nodepb.UpdateNodeResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			if data.NodeName == resp.Node.Name {
				check = true
			}
		}
		return check, nil
	},
}

var TC_registry_update_node_state = &test.TestCase{
	Name:        "Update node state",
	Description: "Update Node State",
	Data:        &nodepb.UpdateNodeResponse{},

	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a, ok := tc.GetWorkflowData().(*RegistryData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}
		a.reqUpdateNodeState = api.UpdateNodeStateRequest{
			NodeId: a.NodeId,
			State:  "active",
		}
		tc.SaveWorkflowData(a)
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
			tc.Data, err = a.RegistryClient.UpdateNodeState(a.reqUpdateNodeState)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {

		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*nodepb.UpdateNodeResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			if data.NodeName == resp.Node.Name {
				check = true
			}
		}
		return check, nil
	},
}

func TC_registry_add_node_to_site(typ string) *test.TestCase {
	return &test.TestCase{
		Name:        "Add node to site",
		Description: "Add Node to site",
		Data:        &nodepb.AddNodeToSiteResponse{},

		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			// Setup required for test case Initialize any
			// test specific data if required
			a, ok := tc.GetWorkflowData().(*RegistryData)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}
			var nodeId = ""
			if typ == "parent" {
				nodeId = a.NodeId
			} else if typ == "left" {
				nodeId = a.lNodeId
			} else if typ == "right" {
				nodeId = a.rNodeId
			}
			a.reqAddNodeToSite = api.AddNodeToSiteRequest{
				NodeId:    nodeId,
				SiteId:    a.SiteId,
				NetworkId: a.NetworkId,
			}
			tc.SaveWorkflowData(a)
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
				tc.Data, err = a.RegistryClient.AddToSite(a.reqAddNodeToSite)
			} else {
				log.Errorf("Invalid data type for Workflow data.")
				return fmt.Errorf("invalid data type for Workflow data")
			}
			return err
		}}
}

var TC_registry_attach_node = &test.TestCase{
	Name:        "Attach node",
	Description: "Attach node",
	Data:        &nodepb.AttachNodesResponse{},

	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a, ok := tc.GetWorkflowData().(*RegistryData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}
		a.reqAttachNode = api.AttachNodesRequest{
			ParentNode: a.NodeId,
			AmpNodeL:   a.lNodeId,
			AmpNodeR:   a.rNodeId,
		}
		tc.SaveWorkflowData(a)
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
			tc.Data, err = a.RegistryClient.AttachNode(a.reqAttachNode)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {

		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*nodepb.AttachNodesResponse)
		if resp != nil {
			check = true
		}
		return check, nil
	},
}

var TC_registry_detach_node = &test.TestCase{
	Name:        "Detach node",
	Description: "Detach node",
	Data:        &nodepb.DetachNodeResponse{},

	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a, ok := tc.GetWorkflowData().(*RegistryData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}
		a.reqDetachNode = api.DetachNodeRequest{
			NodeId: a.rNodeId,
		}
		tc.SaveWorkflowData(a)
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
			tc.Data, err = a.RegistryClient.DetachNode(a.reqDetachNode)
		} else {
			log.Errorf("Invalid data type for Workflow data.")
			return fmt.Errorf("invalid data type for Workflow data")
		}
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {

		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*nodepb.DetachNodeResponse)
		if resp != nil {
			check = true
		}
		return check, nil
	},
}

var TC_registry_get_nodes = &test.TestCase{
	Name:        "Get nodes",
	Description: "Get nodes",
	Data:        &nodepb.GetNodesResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a := tc.GetWorkflowData().(*RegistryData)
		a.reqGetNodes = api.GetNodesRequest{
			Free: true,
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

		tc.Data, err = td.RegistryClient.GetNodes(td.reqGetNodes)

		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*nodepb.GetNodesResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			for _, node := range resp.Nodes {
				check = true
				if node.Id == data.NodeId {
					break
				}
			}
		}
		return check, nil
	},
}

var TC_registry_get_node = &test.TestCase{
	Name:        "Get node",
	Description: "Get node",
	Data:        &nodepb.GetNodeResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a := tc.GetWorkflowData().(*RegistryData)
		a.reqGetNode = api.GetNodeRequest{
			NodeId: a.NodeId,
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

		tc.Data, err = td.RegistryClient.GetNode(td.reqGetNode)

		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*nodepb.GetNodeResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			if data.NodeId == resp.Node.Id {
				check = true
			}
		}
		return check, nil
	},
}

var TC_registry_get_nodes_by_site = &test.TestCase{
	Name:        "Get nodes for site",
	Description: "Get nodes for site",
	Data:        &nodepb.GetBySiteResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Setup required for test case Initialize any
		// test specific data if required
		a := tc.GetWorkflowData().(*RegistryData)
		a.reqGetNodeForSite = api.GetSiteNodesRequest{
			SiteId: a.SiteId,
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

		tc.Data, err = td.RegistryClient.GetNodesForSite(td.reqGetNodeForSite)

		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during test case
		check := false

		resp := tc.GetData().(*nodepb.GetBySiteResponse)
		if resp != nil {
			data := tc.GetWorkflowData().(*RegistryData)
			if data.SiteId == resp.SiteId {
				check = true
			}
		}
		return check, nil
	},
}
