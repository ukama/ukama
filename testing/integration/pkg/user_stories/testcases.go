package subscriber

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v4"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/testing/integration/pkg"
	"github.com/ukama/ukama/testing/integration/pkg/dataplan"
	"github.com/ukama/ukama/testing/integration/pkg/nucleus"
	"github.com/ukama/ukama/testing/integration/pkg/registry"
	"github.com/ukama/ukama/testing/integration/pkg/subscriber"
	"github.com/ukama/ukama/testing/integration/pkg/test"
	"github.com/ukama/ukama/testing/integration/pkg/utils"

	napi "github.com/ukama/ukama/systems/nucleus/api-gateway/pkg/rest"
	rapi "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"

	orgpb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"
	userpb "github.com/ukama/ukama/systems/nucleus/user/pb/gen"
	mempb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
)

var config *pkg.Config

type UserStoriesData struct {
	SubscriberClient *subscriber.SubscriberClient
	RegistryClient   *registry.RegistryClient
	NucleusClient    *nucleus.NucleusClient
	DataplanClient   *dataplan.DataplanClient

	MbHost string
	w      *test.Workflow

	OrgId      string
	OrgName    string
	OrgOwnerId string
	UserId     string
	UserAuthId string
	NetworkId  string

	reqGetOrg napi.GetOrgRequest
	reqAddOrg napi.AddOrgRequest

	reqAddUser       napi.AddUserRequest
	reqGetUser       napi.GetUserRequest
	reqWhoami        napi.GetUserRequest
	reqGetUserByAuth napi.GetUserByAuthIdRequest

	reqGetMember rapi.GetMemberRequest

	reqAddNetwork  rapi.AddNetworkRequest
	reqGetNetwork  rapi.GetNetworkRequest
	reqGetNetworks rapi.GetNetworksRequest
	reqGetSites    rapi.GetNetworkRequest

	reqAddNode         rapi.AddNodeRequest
	reqUpdateNode      rapi.UpdateNodeRequest
	reqUpdateNodeState rapi.UpdateNodeStateRequest
	reqAttachNode      rapi.AttachNodesRequest
	reqDetachNode      rapi.DetachNodeRequest
	reqAddNodeToSite   rapi.AddNodeToSiteRequest
	reqGetNodes        rapi.GetNodesRequest
	reqGetNode         rapi.GetNodeRequest
	reqGetNodeForSite  rapi.GetSiteNodesRequest
}

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stderr)
}

func InitializeData() *UserStoriesData {
	config = pkg.NewConfig()

	d := &UserStoriesData{}
	d.NucleusClient = nucleus.NewNucleusClient(config.System.Nucleus)
	d.RegistryClient = registry.NewRegistryClient(config.System.Registry)
	d.SubscriberClient = subscriber.NewSubscriberClient(config.System.Subscriber)
	d.DataplanClient = dataplan.NewDataplanClient(config.System.Dataplan)
	d.MbHost = config.System.MessageBus

	return d
}

func getWorkflowData(tc *test.TestCase) (*UserStoriesData, error) {
	a, ok := tc.GetWorkflowData().(*UserStoriesData)
	if !ok {
		log.Errorf("Invalid data type for Workflow data.")
		return nil, fmt.Errorf("invalid data type for Workflow data")
	}
	return a, nil
}

var Story_add_user = &test.TestCase{
	Name:        "Add User",
	Description: "After successful signup, addUser will be called to add the user to the system",
	Data:        &userpb.AddResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Prepare the data for the test case
		a, err := getWorkflowData(tc)
		if err != nil {
			return err
		}

		log.Debugf("Setting up watcher for %s", tc.Name)
		tc.Watcher = utils.SetupWatcher(a.MbHost,
			[]string{"event.cloud.users.user.add"})

		a.reqGetOrg = napi.GetOrgRequest{
			OrgName: config.OrgName,
		}

		res, err := a.NucleusClient.GetOrg(a.reqGetOrg)
		a.OrgId = res.Org.Id
		a.OrgName = res.Org.Name
		a.OrgOwnerId = res.Org.Owner

		a.reqAddUser = napi.AddUserRequest{
			Email:  strings.ToLower(faker.Email()),
			Name:   strings.ToLower(faker.FirstName()),
			Phone:  strings.ToLower(faker.Phonenumber()),
			AuthId: strings.ToLower(faker.UUIDHyphenated()),
		}

		return err
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		// Test Case
		var err error
		a, ok := getWorkflowData(tc)
		if ok != nil {
			return ok
		}
		tc.Data, err = a.NucleusClient.AddUser(a.reqAddUser)
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during user stories
		check1, check2, check3, check4 := false, false, false, false

		resp := tc.GetData().(*userpb.AddResponse)

		if resp != nil {
			a, ok := getWorkflowData(tc)
			if ok != nil {
				return false, ok
			}

			a.reqGetUser = napi.GetUserRequest{
				UserId: resp.User.Id,
			}
			a.reqGetUserByAuth = napi.GetUserByAuthIdRequest{
				AuthId: resp.User.AuthId,
			}
			a.reqWhoami = napi.GetUserRequest{
				UserId: resp.User.Id,
			}
			a.reqGetMember = rapi.GetMemberRequest{
				UserUuid: resp.User.Id,
			}

			tc1, err := a.NucleusClient.GetUser(a.reqGetUser)
			if err != nil {
				return check1, fmt.Errorf("add user story failed on getUser. Error %v", err)
			} else if tc1.User.Id == resp.User.Id {
				check1 = true
			}

			tc2, err := a.NucleusClient.GetUserByAuthId(a.reqGetUserByAuth)
			if err != nil {
				return check2, fmt.Errorf("add user story failed on getUserByAuth. Error %v", err)
			} else if tc2.User.Id == resp.User.Id {
				check2 = true
			}

			tc3, err := a.NucleusClient.Whoami(a.reqWhoami)
			if err != nil {
				return check3, fmt.Errorf("add user story failed on whoami. Error %v", err)
			} else if tc3.MemberOf[0].Id == a.OrgId {
				check3 = true
			}

			tc4, err := a.RegistryClient.GetMember(a.reqGetMember)
			if err != nil {
				return check4, fmt.Errorf("add user story failed on getMember. Error %v", err)
			} else if tc4.Member.OrgId == a.OrgId {
				check4 = true
			}
		}

		if check1 && check2 && check3 && check4 {
			return check1 && check2 && check3 && check4, nil
		} else {
			return false, fmt.Errorf("add user story failed. %v", nil)
		}
	},

	ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
		// Here we save any data required to be saved from the
		// test case Cleanup any test specific data

		resp, ok := tc.GetData().(*userpb.AddResponse)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a, ok := tc.GetWorkflowData().(*UserStoriesData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a.UserId = resp.User.Id
		a.UserAuthId = resp.User.AuthId

		tc.SaveWorkflowData(a)
		log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
		tc.Watcher.Stop()

		return nil
	},
}

var Story_add_org = &test.TestCase{
	Name:        "Add org",
	Description: "After successful signup, User can create org",
	Data:        &orgpb.AddResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Prepare the data for the test case
		a, err := getWorkflowData(tc)
		if err != nil {
			return err
		}

		log.Debugf("Setting up watcher for %s", tc.Name)
		tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.nucleus.org.create"})

		a.reqGetUser = napi.GetUserRequest{
			UserId: a.UserId,
		}

		res, err := a.NucleusClient.GetUser(a.reqGetUser)
		a.UserId = res.User.Id
		a.UserAuthId = res.User.AuthId

		a.reqAddOrg = napi.AddOrgRequest{
			Owner:       a.UserId,
			OrgName:     strings.ToLower(faker.FirstName()) + "-org",
			Certificate: "-----BEGIN CERTIFICATE-----",
		}

		return err
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		// Test Case
		var err error
		a, ok := getWorkflowData(tc)
		if ok != nil {
			return ok
		}
		tc.Data, err = a.NucleusClient.AddOrg(a.reqAddOrg)
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during user stories
		check1, check2, check3, check4 := false, false, false, false

		resp := tc.GetData().(*orgpb.AddResponse)

		if resp != nil {
			a, ok := getWorkflowData(tc)
			if ok != nil {
				return false, ok
			}

			a.reqGetOrg = napi.GetOrgRequest{
				OrgName: resp.Org.Name,
			}
			a.reqWhoami = napi.GetUserRequest{
				UserId: a.UserId,
			}
			a.reqGetMember = rapi.GetMemberRequest{
				UserUuid: a.UserId,
			}

			tc1, err := a.NucleusClient.GetOrg(a.reqGetOrg)
			if err != nil {
				return check1, fmt.Errorf("add org story failed on getOrg. Error %v", err)
			} else if tc1.Org.Id == resp.Org.Id {
				check1 = true
			}

			tc2, err := a.NucleusClient.Whoami(a.reqWhoami)
			if err != nil {
				return check2, fmt.Errorf("add org story failed on whoami. Error %v", err)
			} else {
				for _, org := range tc2.OwnerOf {
					if org.Id == resp.Org.Id {
						check2 = true
						break
					}
				}
			}

			tc3, err := a.RegistryClient.GetMembers()
			if err != nil {
				return check3, fmt.Errorf("add org story failed on getMembers. Error %v", err)
			} else {
				for _, member := range tc3.Members {
					if member.UserId == a.UserId {
						check3 = true
						break
					}
				}
			}

			tc4, err := a.RegistryClient.GetMember(a.reqGetMember)
			if err != nil {
				return check4, fmt.Errorf("add user story failed on getMember. Error %v", err)
			} else if tc4.Member.OrgId == a.OrgId {
				check4 = true
			}
		}

		if check1 && check2 && check3 && check4 {
			return check1 && check2 && check3 && check4, nil
		} else {
			return false, fmt.Errorf("add org story failed. %v", nil)
		}
	},

	ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
		// Here we save any data required to be saved from the
		// test case Cleanup any test specific data

		resp, ok := tc.GetData().(*orgpb.AddResponse)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a, ok := tc.GetWorkflowData().(*UserStoriesData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a.OrgId = resp.Org.Id
		a.OrgName = resp.Org.Name
		a.OrgOwnerId = resp.Org.Owner

		tc.SaveWorkflowData(a)
		log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
		tc.Watcher.Stop()

		return nil
	},
}

var Story_add_network = &test.TestCase{
	Name:        "Add network",
	Description: "Add network to org",
	Data:        &netpb.AddResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Prepare the data for the test case
		a, err := getWorkflowData(tc)
		if err != nil {
			return err
		}

		log.Debugf("Setting up watcher for %s", tc.Name)
		tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.registry.network.create"})

		a.reqGetOrg = napi.GetOrgRequest{
			OrgName: a.OrgName,
		}
		a.reqGetMember = rapi.GetMemberRequest{
			UserUuid: a.UserId,
		}

		orgResp, err := a.NucleusClient.GetOrg(a.reqGetOrg)
		if err != nil {
			return err
		}

		memResp, err := a.RegistryClient.GetMember(a.reqGetMember)
		if err != nil {
			return err
		}

		if orgResp.Org.Id == memResp.Member.OrgId && (memResp.Member.Role == mempb.RoleType_OWNER || memResp.Member.Role == mempb.RoleType_ADMIN) {
			a.reqAddNetwork = rapi.AddNetworkRequest{
				OrgName: orgResp.Org.Name,
				NetName: strings.ToLower(faker.FirstName()) + "-network",
			}
		} else {
			return fmt.Errorf("user is not an owner or admin of the org")
		}

		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		// Test Case
		var err error
		a, ok := getWorkflowData(tc)
		if ok != nil {
			return ok
		}
		tc.Data, err = a.RegistryClient.AddNetwork(a.reqAddNetwork)
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during user stories
		check1, check2, check3 := false, false, false

		resp := tc.GetData().(*netpb.AddResponse)

		if resp != nil {
			a, ok := getWorkflowData(tc)
			if ok != nil {
				return false, ok
			}

			a.reqGetNetwork = rapi.GetNetworkRequest{
				NetworkId: resp.Network.Id,
			}
			a.reqGetNetworks = rapi.GetNetworksRequest{
				OrgUuid: a.OrgId,
			}
			a.reqGetSites = rapi.GetNetworkRequest{
				NetworkId: resp.Network.Id,
			}

			tc1, err := a.RegistryClient.GetNetwork(a.reqGetNetwork)
			if err != nil {
				return check1, fmt.Errorf("add network story failed on getNetwork. Error %v", err)
			} else if tc1.Network.Id == resp.Network.Id {
				check1 = true
			}

			tc2, err := a.RegistryClient.GetNetworks(a.reqGetNetworks)
			if err != nil {
				return check2, fmt.Errorf("add network story failed on getNetworks. Error %v", err)
			} else if tc2.OrgId == a.OrgId {
				for _, network := range tc2.Networks {
					if network.Id == resp.Network.Id {
						check2 = true
						break
					}
				}
			}

			tc3, err := a.RegistryClient.GetSites(a.reqGetSites)
			if err != nil {
				return check3, fmt.Errorf("add network story failed on getNetworks. Error %v", err)
			} else if tc3.NetworkId == resp.Network.Id && len(tc3.Sites) == 0 {
				check3 = true
			}

		}

		if check1 && check2 && check3 {
			return check1 && check2 && check3, nil
		} else {
			return false, fmt.Errorf("add network story failed. %v", nil)
		}
	},

	ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
		// Here we save any data required to be saved from the
		// test case Cleanup any test specific data

		resp, ok := tc.GetData().(*netpb.AddResponse)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a, ok := tc.GetWorkflowData().(*UserStoriesData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a.NetworkId = resp.Network.Id

		tc.SaveWorkflowData(a)
		log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
		tc.Watcher.Stop()

		return nil
	},
}

var Story_add_node = &test.TestCase{
	Name:        "Add node",
	Description: "Add node",
	Data:        &netpb.AddResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Prepare the data for the test case
		a, err := getWorkflowData(tc)
		if err != nil {
			return err
		}

		log.Debugf("Setting up watcher for %s", tc.Name)
		tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.registry.node.add"})

		a.reqGetOrg = napi.GetOrgRequest{
			OrgName: a.OrgName,
		}
		orgResp, err := a.NucleusClient.GetOrg(a.reqGetOrg)
		if err != nil {
			return err
		}

		if err != nil {
			a.reqAddNode = rapi.AddNodeRequest{
				NodeId: utils.RandomGetNodeId("tnode"),
				Name:  strings.ToLower(faker.FirstName()) + "-node",
				OrgId: orgResp.Org.Id,
				State: "onboarded",
			}
		} else {
			return fmt.Errorf("user is not an owner or admin of the org")
		}

		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		// Test Case
		var err error
		a, ok := getWorkflowData(tc)
		if ok != nil {
			return ok
		}
		tc.Data, err = a.RegistryClient.AddNetwork(a.reqAddNetwork)
		return err
	},

	StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
		// Check for possible failures during user stories
		check1, check2, check3 := false, false, false

		resp := tc.GetData().(*netpb.AddResponse)

		if resp != nil {
			a, ok := getWorkflowData(tc)
			if ok != nil {
				return false, ok
			}

			a.reqGetNetwork = rapi.GetNetworkRequest{
				NetworkId: resp.Network.Id,
			}
			a.reqGetNetworks = rapi.GetNetworksRequest{
				OrgUuid: a.OrgId,
			}
			a.reqGetSites = rapi.GetNetworkRequest{
				NetworkId: resp.Network.Id,
			}

			tc1, err := a.RegistryClient.GetNetwork(a.reqGetNetwork)
			if err != nil {
				return check1, fmt.Errorf("add network story failed on getNetwork. Error %v", err)
			} else if tc1.Network.Id == resp.Network.Id {
				check1 = true
			}

			tc2, err := a.RegistryClient.GetNetworks(a.reqGetNetworks)
			if err != nil {
				return check2, fmt.Errorf("add network story failed on getNetworks. Error %v", err)
			} else if tc2.OrgId == a.OrgId {
				for _, network := range tc2.Networks {
					if network.Id == resp.Network.Id {
						check2 = true
						break
					}
				}
			}

			tc3, err := a.RegistryClient.GetSites(a.reqGetSites)
			if err != nil {
				return check3, fmt.Errorf("add network story failed on getNetworks. Error %v", err)
			} else if tc3.NetworkId == resp.Network.Id && len(tc3.Sites) == 0 {
				check3 = true
			}

		}

		if check1 && check2 && check3 {
			return check1 && check2 && check3, nil
		} else {
			return false, fmt.Errorf("add network story failed. %v", nil)
		}
	},

	ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
		// Here we save any data required to be saved from the
		// test case Cleanup any test specific data

		resp, ok := tc.GetData().(*netpb.AddResponse)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a, ok := tc.GetWorkflowData().(*UserStoriesData)
		if !ok {
			log.Errorf("Invalid data type for Workflow data.")

			return fmt.Errorf("invalid data type for Workflow data")
		}

		a.NetworkId = resp.Network.Id

		tc.SaveWorkflowData(a)
		log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
		tc.Watcher.Stop()

		return nil
	},
}
