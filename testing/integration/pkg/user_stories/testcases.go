package subscriber

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker/v4"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/testing/integration/pkg"
	"github.com/ukama/ukama/testing/integration/pkg/dataplan"
	"github.com/ukama/ukama/testing/integration/pkg/nucleus"
	"github.com/ukama/ukama/testing/integration/pkg/registry"
	"github.com/ukama/ukama/testing/integration/pkg/subscriber"
	"github.com/ukama/ukama/testing/integration/pkg/test"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/common/validation"
	dapi "github.com/ukama/ukama/systems/data-plan/api-gateway/pkg/rest"
	napi "github.com/ukama/ukama/systems/nucleus/api-gateway/pkg/rest"
	rapi "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"
	sapi "github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/rest"
	smutil "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/utils"

	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	ppb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	rpb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	orgpb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"
	userpb "github.com/ukama/ukama/systems/nucleus/user/pb/gen"
	invpb "github.com/ukama/ukama/systems/registry/invitation/pb/gen"
	mempb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	nodepb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	srpb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	smpb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	sppb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
)

var config *pkg.Config

const MAX_POOL = 5

type UserStoriesData struct {
	SubscriberClient *subscriber.SubscriberClient
	RegistryClient   *registry.RegistryClient
	NucleusClient    *nucleus.NucleusClient
	DataplanClient   *dataplan.DataplanClient

	MbHost string
	w      *test.Workflow

	OrgId         string
	OrgName       string
	OrgOwnerId    string
	UserId        string
	UserAuthId    string
	NetworkId     string
	SiteId        string
	NodeId        string
	lNodeId       string
	rNodeId       string
	simType       string
	country       string
	provider      string
	invitationId  string
	invitedUserId string
	subscriberId  string
	baserateId    string
	packageId     string
	spackageId    string
	ICCID         []string
	SimId         string
	SimPackageId  string

	reqGetOrg napi.GetOrgRequest
	reqAddOrg napi.AddOrgRequest

	reqAddUser       napi.AddUserRequest
	reqGetUser       napi.GetUserRequest
	reqWhoami        napi.GetUserRequest
	reqGetUserByAuth napi.GetUserByAuthIdRequest

	reqGetMember  rapi.GetMemberRequest
	reqAddMember  rapi.MemberRequest
	reqGetMembers rapi.GetMembersRequest

	reqAddNetwork  rapi.AddNetworkRequest
	reqGetNetwork  rapi.GetNetworkRequest
	reqGetNetworks rapi.GetNetworksRequest
	reqGetSites    rapi.GetNetworkRequest
	reqAddSite     rapi.AddSiteRequest

	reqAddNode           rapi.AddNodeRequest
	reqUpdateNode        rapi.UpdateNodeRequest
	reqUpdateNodeState   rapi.UpdateNodeStateRequest
	reqAttachNode        rapi.AttachNodesRequest
	reqDetachNode        rapi.DetachNodeRequest
	reqAddNodeToSite     rapi.AddNodeToSiteRequest
	reqGetNodes          rapi.GetNodesRequest
	reqGetNode           rapi.GetNodeRequest
	reqGetNodeForSite    rapi.GetSiteNodesRequest
	reqGetNodesByNetwork rapi.GetNetworkNodesRequest

	reqAddInvite        rapi.AddInvitationRequest
	reqUpdateInvitation rapi.UpdateInvitationRequest
	reqGetInvite        rapi.GetInvitationRequest
	reqGetInvites       rapi.GetInvitationByOrgRequest

	reqUploadBaseRates    dapi.UploadBaseRatesRequest
	reqGetBaseRates       dapi.GetBaseRatesByCountryRequest
	reqGetBaseRate        dapi.GetBaseRateRequest
	reqGetBaseratePackage dapi.GetBaseRatesForPeriodRequest
	reqGetBaserateHistory dapi.GetBaseRatesByCountryRequest
	reqGetBaseratesPeriod dapi.GetBaseRatesForPeriodRequest
	reqGetRateByUser      dapi.GetRateRequest

	reqAddPackage       dapi.AddPackageRequest
	reqGetPackage       dapi.PackagesRequest
	reqGetPackageDetail dapi.PackagesRequest
	reqGetPackageForOrg dapi.GetPackageByOrgRequest

	reqAddUserMarkup    dapi.SetMarkupRequest
	reqSetDefaultMarkup dapi.SetDefaultMarkupRequest
	reqGetMarkupForUser dapi.GetMarkupRequest
	reqDeleteMarkup     dapi.DeleteMarkupRequest

	reqSimPoolUploadSimReq  sapi.SimPoolUploadSimReq
	reqSimPoolStatByTypeReq sapi.SimPoolTypeReq
	reqSimByIccidReq        sapi.SimByIccidReq
	reqGetSimsByTypeReq     sapi.SimPoolTypeReq

	reqAddSubscriber sapi.SubscriberAddReq
	reqGetSubscriber sapi.SubscriberGetReq

	reqAllocateSim       sapi.AllocateSimReq
	reqAddpackageToSub   sapi.AddPkgToSimReq
	reqGetPackagesForSim sapi.SimReq
	reqGetSubscriberSims sapi.GetSimsBySubReq
	reqGetSimById        sapi.SimReq
	reqToggleSimState    sapi.ActivateDeactivateSimReq
	reqSetActivePackage  sapi.SetActivePackageForSimReq
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
	d.ICCID = make([]string, MAX_POOL)
	d.simType = "test"
	d.country = "The lunar maria"
	d.provider = "ABC Tel"
	d.OrgId = config.OrgId
	d.OrgName = config.OrgName
	d.OrgOwnerId = config.OrgOwnerId

	d.packageId = "29c7485b-cb87-43ed-8619-30c891cdd197"
	d.NetworkId = "eaa6b6e0-0ef4-4bd6-bb58-39964183a9b0"

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
			if err == nil && tc1.Org.Id == resp.Org.Id {
				return check1, fmt.Errorf("add org story failed on getOrg. Error %v", err)
			} else if err != nil {
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
				if !check2 {
					return check2, fmt.Errorf("add org story failed on whoami. Error %v", err)
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

		a, err := getWorkflowData(tc)
		if err != nil {
			return err
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

			if err == nil && tc1.User.Id == resp.User.Id {
				check1 = true
			} else if err != nil {
				return check1, fmt.Errorf("add user story failed on getUser. Error %v", err)
			}

			tc2, err := a.NucleusClient.GetUserByAuthId(a.reqGetUserByAuth)
			if err == nil && tc2.User.Id == resp.User.Id {
				check2 = true
			} else if err != nil {
				return check2, fmt.Errorf("add user story failed on getUserByAuth. Error %v", err)
			}

			tc3, err := a.NucleusClient.Whoami(a.reqWhoami)
			if err == nil && tc3.MemberOf[0].Id == a.OrgId {
				check3 = true
			} else if err != nil {
				return check3, fmt.Errorf("add user story failed on whoami. Error %v", err)
			}

			tc4, err := a.RegistryClient.GetMember(a.reqGetMember)
			if err == nil && tc4.Member.OrgId == a.OrgId {
				check4 = true
			} else if err != nil {
				return check4, fmt.Errorf("add user story failed on getMember. Error %v", err)
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

		a, err := getWorkflowData(tc)
		if err != nil {
			return err
		}

		a.UserId = resp.User.Id
		a.UserAuthId = resp.User.AuthId

		tc.SaveWorkflowData(a)
		log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
		tc.Watcher.Stop()

		return nil
	},
}

var Story_add_network = &test.TestCase{
	Name:        "Add network success case",
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
			UserUuid: a.OrgOwnerId,
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

		a, err := getWorkflowData(tc)
		if err != nil {
			return err
		}

		a.NetworkId = resp.Network.Id
		a.reqAddSite = rapi.AddSiteRequest{
			NetworkId: a.NetworkId,
			SiteName:  strings.ToLower(faker.FirstName()) + "-site",
		}
		res, er := a.RegistryClient.AddSite(a.reqAddSite)
		if er != nil {
			return er
		} else {
			a.SiteId = res.Site.Id
		}

		tc.SaveWorkflowData(a)
		log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
		tc.Watcher.Stop()

		return nil
	},
}

var Story_add_network_failed = &test.TestCase{
	Name:        "Add network failed case",
	Description: "Scenarios: Network name already exists",
	Data:        &netpb.AddResponse{},
	SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
		// Prepare the data for the test case
		a, err := getWorkflowData(tc)
		if err != nil {
			return err
		}

		log.Debugf("Setting up watcher for %s", tc.Name)
		tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.registry.network.create"})

		a.reqGetNetworks = rapi.GetNetworksRequest{
			OrgUuid: a.OrgId,
		}

		netResp, err := a.RegistryClient.GetNetworks(a.reqGetNetworks)
		if err != nil {
			return err
		}

		if netResp != nil && len(netResp.Networks) > 0 {
			netName := netResp.Networks[0].Name
			a.reqAddNetwork = rapi.AddNetworkRequest{
				OrgName: a.OrgName,
				NetName: netName,
			}
		} else {
			return fmt.Errorf("user is not an owner or admin of the org")
		}

		return nil
	},

	Fxn: func(ctx context.Context, tc *test.TestCase) error {
		// Test Case
		a, ok := getWorkflowData(tc)
		if ok != nil {
			return ok
		}
		_, err := a.RegistryClient.AddNetwork(a.reqAddNetwork)

		if err != nil && strings.Contains(err.Error(), "duplicate key value") {
			return nil
		}
		return fmt.Errorf("add network story duplicate name test failed. %v", nil)
	},
}

func Story_add_node(typ string) *test.TestCase {
	return &test.TestCase{
		Name:        "Add node",
		Description: "Add node",
		Data:        &nodepb.AddNodeResponse{},
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
				return err
			} else {
				var nId = ""
				if typ == "parent" {
					nId = utils.RandomGetNodeId("tnode")
				} else if typ == "left" || typ == "right" {
					nId = utils.RandomGetNodeId("anode")
				}
				a.reqAddNode = rapi.AddNodeRequest{
					NodeId: nId,
					Name:   strings.ToLower(faker.FirstName()) + "-node",
					OrgId:  orgResp.Org.Id,
					State:  "onboarded",
				}
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
			tc.Data, err = a.RegistryClient.AddNode(a.reqAddNode)
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			// Check for possible failures during user stories
			check1, check2 := false, false

			resp := tc.GetData().(*nodepb.AddNodeResponse)

			if resp != nil {
				a, ok := getWorkflowData(tc)
				if ok != nil {
					return false, ok
				}

				a.reqGetNode = rapi.GetNodeRequest{
					NodeId: resp.Node.Id,
				}
				a.reqGetNodes = rapi.GetNodesRequest{
					Free: true,
				}

				tc1, err := a.RegistryClient.GetNode(a.reqGetNode)
				if err != nil {
					return check1, fmt.Errorf("add node story failed on getNode. Error %v", err)
				} else if tc1.Node.Id == resp.Node.Id {
					check1 = true
				}

				tc2, err := a.RegistryClient.GetNodes(a.reqGetNodes)
				if err != nil {
					return check2, fmt.Errorf("add node story failed on getNodes. Error %v", err)
				} else {
					for _, node := range tc2.Nodes {
						if node.Id == resp.Node.Id {
							check2 = true
							break
						}
					}
				}
			}

			if check1 && check2 {
				return check1 && check2, nil
			} else {
				return false, fmt.Errorf("add node story failed. %v", nil)
			}
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			// Here we save any data required to be saved from the
			// test case Cleanup any test specific data

			resp, ok := tc.GetData().(*nodepb.AddNodeResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}

			if typ == "parent" {
				a.NodeId = resp.Node.Id
			} else if typ == "left" {
				a.lNodeId = resp.Node.Id
			} else if typ == "right" {
				a.rNodeId = resp.Node.Id
			}

			tc.SaveWorkflowData(a)
			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()

			return nil
		},
	}
}

func Story_add_node_to_site(typ string) *test.TestCase {
	return &test.TestCase{
		Name:        "Add node to site",
		Description: "Add node to site",
		Data:        &nodepb.AddNodeToSiteResponse{},
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			// Prepare the data for the test case
			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}

			log.Debugf("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.registry.node.add.site"})

			var nId = ""
			if typ == "parent" {
				nId = a.NodeId
			} else if typ == "left" {
				nId = a.lNodeId
			} else if typ == "right" {
				nId = a.rNodeId
			}
			a.reqAddNodeToSite = rapi.AddNodeToSiteRequest{
				NodeId:    nId,
				SiteId:    a.SiteId,
				NetworkId: a.NetworkId,
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
			tc.Data, err = a.RegistryClient.AddToSite(a.reqAddNodeToSite)
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			// Check for possible failures during user stories
			check1, check2, check3 := false, false, false

			resp := tc.GetData().(*nodepb.AddNodeToSiteResponse)

			if resp != nil {
				a, ok := getWorkflowData(tc)
				if ok != nil {
					return false, ok
				}
				var nId = ""
				if typ == "parent" {
					nId = a.NodeId
				} else if typ == "left" {
					nId = a.lNodeId
				} else if typ == "right" {
					nId = a.rNodeId
				}

				a.reqGetNode = rapi.GetNodeRequest{
					NodeId: nId,
				}

				a.reqGetNodeForSite = rapi.GetSiteNodesRequest{
					SiteId: a.SiteId,
				}

				a.reqGetNodesByNetwork = rapi.GetNetworkNodesRequest{
					NetworkId: a.NetworkId,
				}

				tc1, err := a.RegistryClient.GetNode(a.reqGetNode)
				if err != nil {
					return check1, fmt.Errorf("add node to site story failed on getNode. Error %v", err)
				} else if tc1.Node.Site.SiteId == a.SiteId {
					check1 = true
				}

				tc2, err := a.RegistryClient.GetNodesForSite(a.reqGetNodeForSite)
				if err != nil {
					return check2, fmt.Errorf("add node to site story failed on getNodesForSite. Error %v", err)
				} else if tc2.SiteId == a.SiteId {
					check2 = true
				}

				tc3, err := a.RegistryClient.GetNodesByNetwork(a.reqGetNodesByNetwork)
				if err != nil {
					return check3, fmt.Errorf("add node to site story failed on getNodesByNetwork. Error %v", err)
				} else if tc3.NetworkId == a.NetworkId {
					check3 = true
				}
			}

			if check1 && check2 && check3 {
				return true, nil
			} else {
				return false, fmt.Errorf("add node to site story failed. %v", nil)
			}
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			// Here we save any data required to be saved from the
			// test case Cleanup any test specific data

			resp, ok := tc.GetData().(*nodepb.AddNodeToSiteResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}

			tc.SaveWorkflowData(a)
			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()

			return nil
		},
	}
}

func Story_attach_node() *test.TestCase {
	return &test.TestCase{
		Name:        "Attach nodes",
		Description: "Attch amplifier nodes with tower node",
		Data:        &nodepb.AttachNodesResponse{},
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			// Prepare the data for the test case
			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}

			log.Debugf("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.registry.node.attach"})

			a.reqAttachNode = rapi.AttachNodesRequest{
				ParentNode: a.NodeId,
				AmpNodeL:   a.lNodeId,
				AmpNodeR:   a.rNodeId,
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
			tc.Data, err = a.RegistryClient.AttachNode(a.reqAttachNode)
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			// Check for possible failures during user stories
			check1 := false

			resp := tc.GetData().(*nodepb.AttachNodesResponse)

			if resp != nil {
				a, ok := getWorkflowData(tc)
				if ok != nil {
					return false, ok
				}

				a.reqGetNode = rapi.GetNodeRequest{
					NodeId: a.NodeId,
				}

				tc1, err := a.RegistryClient.GetNode(a.reqGetNode)
				if err != nil {
					return check1, fmt.Errorf("attach node story failed on getNode. Error %v", err)
				} else if tc1.Node.Id == a.NodeId {
					for i, node := range tc1.Node.Attached {
						if (node.Id == a.lNodeId || node.Id == a.rNodeId) && len(tc1.Node.Attached)-1 == i {
							check1 = true
							break
						}
					}
				}
			}

			if check1 {
				return check1, nil
			} else {
				return false, fmt.Errorf("attach node story failed. %v", nil)
			}
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			// Here we save any data required to be saved from the
			// test case Cleanup any test specific data

			resp, ok := tc.GetData().(*nodepb.AttachNodesResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}

			tc.SaveWorkflowData(a)
			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()

			return nil
		},
	}
}

func Story_invite_add() *test.TestCase {
	return &test.TestCase{
		Name:        "Add invite",
		Description: "Invite user to org",
		Data:        &invpb.AddInvitationResponse{},
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			// Prepare the data for the test case
			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}

			log.Debugf("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.registry.invitation.add"})

			a.reqAddInvite = rapi.AddInvitationRequest{
				Org:   a.OrgName,
				Name:  strings.ToLower(faker.FirstName()) + "-invite",
				Email: strings.ToLower(faker.Email()),
				Role:  "USERS",
			}

			tc.SaveWorkflowData(a)
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			// Test Case
			var err error
			a, ok := getWorkflowData(tc)
			if ok != nil {
				return ok
			}
			tc.Data, err = a.RegistryClient.AddInvitations(a.reqAddInvite)
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			// Check for possible failures during user stories
			check1, check2 := false, false

			resp := tc.GetData().(*invpb.AddInvitationResponse)

			if resp != nil {
				a, ok := getWorkflowData(tc)
				if ok != nil {
					return false, ok
				}

				a.reqGetInvite = rapi.GetInvitationRequest{
					InvitationId: resp.Invitation.Id,
				}
				a.reqGetInvites = rapi.GetInvitationByOrgRequest{
					Org: a.OrgName,
				}

				tc1, err := a.RegistryClient.GetInvitation(a.reqGetInvite)
				timeObj := time.Unix(int64(tc1.Invitation.ExpireAt.Seconds), 0)
				expiry := timeObj.Format(time.RFC3339)
				if err == nil && tc1.Invitation.Id == resp.Invitation.Id &&
					tc1.Invitation.Role == resp.Invitation.Role &&
					validation.IsFutureDate(expiry) == nil &&
					tc1.Invitation.Status == invpb.StatusType_Pending &&
					tc1.Invitation.Org == a.OrgName {
					check1 = true
				} else {
					return check1, fmt.Errorf("add invite story failed on getInvite. Error %v", err)
				}

				tc2, err := a.RegistryClient.GetInvitationByOrg(a.reqGetInvites)
				if err == nil {
					for i, invite := range tc2.Invitations {
						if invite.Id == resp.Invitation.Id && invite.Org == a.OrgName {
							check2 = true
							break
						}
						if len(tc2.Invitations)-1 == i {
							check2 = false
						}
					}
				} else {
					return check1, fmt.Errorf("add invite story failed on getInvitations. Error %v", err)
				}
			}

			if check1 && check2 {
				return check1 && check2, nil
			} else {
				return false, fmt.Errorf("add invite story failed. %v", nil)
			}
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			resp, ok := tc.GetData().(*invpb.AddInvitationResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}
			a.invitationId = resp.Invitation.Id
			tc.SaveWorkflowData(a)
			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()

			return nil
		},
	}
}

func Story_invite_status_update() *test.TestCase {
	return &test.TestCase{
		Name:        "Update invite",
		Description: "Update invite status",
		Data:        &invpb.UpdateInvitationStatusResponse{},
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			// Prepare the data for the test case
			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}

			log.Debugf("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.registry.invitation.status.update"})

			a.reqUpdateInvitation = rapi.UpdateInvitationRequest{
				InvitationId: a.invitationId,
				Status:       invpb.StatusType_name[int32(invpb.StatusType_Accepted)],
			}

			tc.SaveWorkflowData(a)
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			// Test Case
			var err error
			a, ok := getWorkflowData(tc)
			if ok != nil {
				return ok
			}
			tc.Data, err = a.RegistryClient.UpdateInvitations(a.reqUpdateInvitation)
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			// Check for possible failures during user stories
			check1, check2 := false, false

			resp := tc.GetData().(*invpb.UpdateInvitationStatusResponse)

			if resp != nil {
				a, ok := getWorkflowData(tc)
				if ok != nil {
					return false, ok
				}

				a.reqGetInvite = rapi.GetInvitationRequest{
					InvitationId: a.invitationId,
				}

				tc1, err := a.RegistryClient.GetInvitation(a.reqGetInvite)
				timeObj := time.Unix(int64(tc1.Invitation.ExpireAt.Seconds), 0)
				expiry := timeObj.Format(time.RFC3339)
				if err == nil && tc1.Invitation.Id == resp.Id &&
					validation.IsFutureDate(expiry) == nil &&
					tc1.Invitation.Status == invpb.StatusType_Accepted {

					a.reqAddUser = napi.AddUserRequest{
						Name:   tc1.Invitation.Name,
						Email:  strings.ToLower(tc1.Invitation.Email),
						AuthId: faker.UUIDHyphenated(),
						Phone:  strings.ToLower(faker.Phonenumber()),
					}
					check1 = true
					tc2, err := a.NucleusClient.AddUser(a.reqAddUser)
					if err == nil && tc2.User.Email == tc1.Invitation.Email {
						a.invitedUserId = tc2.User.Id

						check2 = true
					} else {
						return check2, fmt.Errorf("update invitation status story failed on addUser. Error %v", err)
					}
				} else {
					return check1, fmt.Errorf("update invitation status story failed on getInvite. Error %v", err)
				}
			}

			if check1 && check2 {
				return true, nil
			} else {
				return false, fmt.Errorf("update invitation status story failed. %v", nil)
			}
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			resp, ok := tc.GetData().(*invpb.UpdateInvitationStatusResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}

			tc.SaveWorkflowData(a)
			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()

			return nil
		},
	}
}

func Story_member_add() *test.TestCase {
	return &test.TestCase{
		Name:        "Add member",
		Description: "Add member to org",
		Data:        &mempb.MemberResponse{},
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			// Prepare the data for the test case
			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}

			log.Debugf("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.registry.member.add"})

			a.reqAddMember = rapi.MemberRequest{
				UserUuid: a.invitedUserId,
			}

			tc.SaveWorkflowData(a)
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			// Test Case
			var err error
			a, ok := getWorkflowData(tc)
			if ok != nil {
				return ok
			}
			tc.Data, err = a.RegistryClient.AddMember(a.reqAddMember)
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			// Check for possible failures during user stories
			check1, check2 := false, false

			resp := tc.GetData().(*mempb.MemberResponse)

			if resp != nil {
				a, ok := getWorkflowData(tc)
				if ok != nil {
					return false, ok
				}

				a.reqGetMember = rapi.GetMemberRequest{
					UserUuid: a.invitedUserId,
				}

				tc1, err := a.RegistryClient.GetMember(a.reqGetMember)
				if err == nil && tc1.Member.UserId == resp.Member.UserId &&
					tc1.Member.OrgId == a.OrgId {
					check1 = true
				} else {
					return check1, fmt.Errorf("add member story failed on getMember. Error %v", err)
				}

				tc2, err := a.RegistryClient.GetMembers()

				if err == nil {
					for _, member := range tc2.Members {
						if member.UserId == a.invitedUserId && member.OrgId == a.OrgId {
							check2 = true
							break
						}
					}
				} else {
					return check2, fmt.Errorf("add member story failed on getMembers. Error %v", err)
				}
			}

			if check1 && check2 {
				return true, nil
			} else {
				return false, fmt.Errorf("add member story failed. %v", nil)
			}
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			resp, ok := tc.GetData().(*mempb.MemberResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}
			tc.SaveWorkflowData(a)
			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()

			return nil
		},
	}
}

func Story_upload_baserate() *test.TestCase {
	return &test.TestCase{
		Name:        "Adding base rate",
		Description: "Add base rate provided by third parties",
		Data:        &bpb.UploadBaseRatesResponse{},
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			a := tc.GetWorkflowData().(*UserStoriesData)
			log.Tracef("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.baserate.rate.update"})

			a.reqUploadBaseRates = dapi.UploadBaseRatesRequest{
				EffectiveAt: utils.GenerateUTCFutureDate(time.Second * 2),
				FileURL:     "https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/template/template.csv",
				EndAt:       utils.GenerateUTCFutureDate(365 * 24 * time.Hour),
				SimType:     a.simType,
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
			tc.Data, err = a.DataplanClient.DataPlanBaseRateUpload(a.reqUploadBaseRates)
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			check1, check2, check3, check4, check5 := false, false, false, false, false

			resp := tc.GetData().(*bpb.UploadBaseRatesResponse)

			if resp != nil {
				a, ok := getWorkflowData(tc)
				if ok != nil {
					return false, ok
				}

				a.reqGetBaseRates = dapi.GetBaseRatesByCountryRequest{
					Provider: a.provider,
					Country:  a.country,
					SimType:  a.simType,
				}

				tc1, err := a.DataplanClient.DataPlanBaseRateGetByCountry(a.reqGetBaseRates)
				if err == nil {
					for _, rate := range tc1.Rates {
						if rate.Country == a.country && rate.SimType == a.simType {
							check1 = true
							break
						}
					}
				}
				if err != nil || !check1 {
					return check1, fmt.Errorf("uploade baserate story failed on GetBaseRatesByCountryRequest. Error %v", err.Error())
				}

				rate := tc1.Rates[0]
				a.baserateId = rate.Uuid
				a.reqGetBaseRate = dapi.GetBaseRateRequest{
					RateId: rate.Uuid,
				}
				a.reqGetBaserateHistory = dapi.GetBaseRatesByCountryRequest{
					Country:     rate.Country,
					Provider:    rate.Provider,
					SimType:     rate.SimType,
					EffectiveAt: rate.EffectiveAt,
				}
				a.reqGetBaseratePackage = dapi.GetBaseRatesForPeriodRequest{
					Country:  rate.Country,
					Provider: rate.Provider,
					SimType:  rate.SimType,
					From:     utils.GenerateUTCFutureDate(time.Second * 5),
					To:       utils.GenerateUTCFutureDate(7 * 24 * time.Hour),
				}
				a.reqGetBaseratesPeriod = dapi.GetBaseRatesForPeriodRequest{
					Country:  rate.Country,
					Provider: rate.Provider,
					SimType:  rate.SimType,
					From:     utils.GenerateUTCFutureDate(7 * 24 * time.Hour),
					To:       utils.GenerateUTCFutureDate(7 * 24 * time.Hour),
				}

				tc2, err := a.DataplanClient.DataPlanBaseRateGet(a.reqGetBaseRate)
				if err == nil && tc2.Rate.Uuid == rate.Uuid {
					check2 = true
				} else {
					return check2, fmt.Errorf("uploade baserate story failed on GetBaserate. Error %v", err.Error())
				}

				tc3, err := a.DataplanClient.DataPlanBaseRateGetForPackage(a.reqGetBaseratePackage)
				if err == nil {
					for _, r := range tc3.Rates {
						if r.Uuid == rate.Uuid {
							check3 = true
							break
						}
					}
				} else {
					return check3, fmt.Errorf("uploade baserate story failed on BaseRateGetForPackage. Error %v", nil)
				}

				tc4, err := a.DataplanClient.DataPlanBaseRateGetByCountry(a.reqGetBaserateHistory)
				if err == nil {
					for _, r := range tc4.Rates {
						if r.Uuid == rate.Uuid {
							check4 = true
							break
						}
					}
				}
				if err != nil || !check4 {
					return check4, fmt.Errorf("uploade baserate story failed on BaseRateGetByCountry. Error %v", err.Error())
				}

				tc5, err := a.DataplanClient.DataPlanBaseRateGetByPeriod(a.reqGetBaseratesPeriod)
				if err == nil {
					for _, r := range tc5.Rates {
						if r.Uuid == rate.Uuid {
							check5 = true
							break
						}
					}
				}
				if err != nil || !check5 {
					return check5, fmt.Errorf("uploade baserate story failed on BaseRateGetByPeriod. Error %v", err.Error())
				}
			}

			return true, nil
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			resp, ok := tc.GetData().(*bpb.UploadBaseRatesResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}

			tc.SaveWorkflowData(a)
			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()

			return nil
		},
	}
}

func Story_markup() *test.TestCase {
	return &test.TestCase{
		Name:        "Add markup",
		Description: "Set default markup percentage, Set user markup percentage",
		Data:        &rpb.UpdateMarkupResponse{},
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			a := tc.GetWorkflowData().(*UserStoriesData)
			log.Tracef("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.markup.user.add"})

			a.reqGetMarkupForUser = dapi.GetMarkupRequest{
				OwnerId: a.OrgOwnerId,
			}
			resp, err := a.DataplanClient.DataPlanGetUserMarkup(a.reqGetMarkupForUser)
			if err == nil && resp.OwnerId == a.OrgOwnerId {
				a.reqDeleteMarkup = dapi.DeleteMarkupRequest{
					OwnerId: a.OrgOwnerId,
				}
				_, err := a.DataplanClient.DataPlanDeleteMarkup(a.reqDeleteMarkup)
				if err != nil {
					return fmt.Errorf("set markup story faild at SetUpFxn on DeleteMarkup. Error %v", err.Error())
				}
			}
			a.reqAddUserMarkup = dapi.SetMarkupRequest{
				OwnerId: a.OrgOwnerId,
				Markup:  10,
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
			tc.Data, err = a.DataplanClient.DataPlanUpdateMarkup(a.reqAddUserMarkup)
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			check1, check2 := false, false

			resp := tc.GetData().(*rpb.UpdateMarkupResponse)

			if resp != nil {
				a, ok := getWorkflowData(tc)
				if ok != nil {
					return false, ok
				}

				a.reqGetMarkupForUser = dapi.GetMarkupRequest{
					OwnerId: a.OrgOwnerId,
				}
				a.reqGetRateByUser = dapi.GetRateRequest{
					UserId:   a.OrgOwnerId,
					Country:  a.country,
					Provider: a.provider,
					SimType:  a.simType,
					To:       utils.GenerateUTCFutureDate(30 * 24 * time.Hour),
					From:     utils.GenerateUTCFutureDate(10 * time.Second),
				}

				tc1, err := a.DataplanClient.DataPlanGetUserMarkup(a.reqGetMarkupForUser)
				if err == nil && tc1.Markup == 10 {
					check1 = true
				}
				if err != nil || !check1 {
					return check1, fmt.Errorf("set markup story faild on GetMarkupForUser. Error %v", err.Error())
				}

				tc2, err := a.DataplanClient.DataPlanGetRate(a.reqGetRateByUser)
				if err == nil {
					for _, r := range tc2.Rates {
						if r.Uuid == a.baserateId {
							check2 = true
							break
						}
					}
				}
				if err != nil || !check2 {
					return check2, fmt.Errorf("set markup story failed on GetRateByUser. Error %v", err.Error())
				}
			}

			return true, nil
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			resp, ok := tc.GetData().(*rpb.UpdateMarkupResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}

			a.reqSetDefaultMarkup = dapi.SetDefaultMarkupRequest{
				Markup: 10,
			}

			// _, err = a.DataplanClient.DataPlanUpdateDefaultMarkup(a.reqSetDefaultMarkup)
			// if err != nil {
			// 	return fmt.Errorf("set markup story faild on setDefaultMarkup. Error %v", err.Error())
			// }

			tc.SaveWorkflowData(a)
			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()

			return nil
		},
	}
}

func Story_package() *test.TestCase {
	return &test.TestCase{
		Name:        "Add package",
		Description: "Add package in an org",
		Data:        &ppb.AddPackageResponse{},
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			/* Setup required for test case
			Initialize any test specific data if required
			*/
			a := tc.GetWorkflowData().(*UserStoriesData)
			log.Tracef("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.package.create"})
			a.reqAddPackage = dapi.AddPackageRequest{
				OwnerId:    a.OrgOwnerId,
				OrgId:      a.OrgId,
				Name:       faker.FirstName() + "-monthly-pack",
				SimType:    a.simType,
				From:       utils.GenerateUTCFutureDate(24 * time.Hour),
				To:         utils.GenerateUTCFutureDate(30 * 24 * time.Hour),
				BaserateId: a.baserateId,
				SmsVolume:  100,
				DataVolume: 1024,
				DataUnit:   "MegaBytes",
				Type:       "postpaid",
				Active:     true,
				Flatrate:   false,
				Apn:        "ukama.tel",
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
			tc.Data, err = a.DataplanClient.DataPlanPackageAdd(a.reqAddPackage)
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			check1, check2, check3 := false, false, false

			resp := tc.GetData().(*ppb.AddPackageResponse)

			if resp != nil {
				a, ok := getWorkflowData(tc)
				if ok != nil {
					return false, ok
				}

				a.reqGetPackage = dapi.PackagesRequest{
					Uuid: resp.Package.Uuid,
				}

				a.reqGetPackageDetail = dapi.PackagesRequest{
					Uuid: resp.Package.Uuid,
				}

				a.reqGetPackageForOrg = dapi.GetPackageByOrgRequest{
					OrgId: a.OrgId,
				}

				tc1, err := a.DataplanClient.DataPlanPackageGetById(a.reqGetPackage)
				if err == nil && tc1.Package.Uuid == resp.Package.Uuid {
					check1 = true
				} else {
					return check1, fmt.Errorf("add package story failed on getPackage. Error %v", err)
				}

				tc2, err := a.DataplanClient.DataPlanPackageDetails(a.reqGetPackageDetail)
				if err == nil && tc2.Package.Uuid == tc2.Package.Uuid && tc2.Package.OrgId == a.OrgId {
					check2 = true
				} else {
					return check2, fmt.Errorf("add package story failed on getPackageDetatil. Error %v", err)
				}

				tc3, err := a.DataplanClient.DataPlanPackageGetByOrg(a.reqGetPackageForOrg)
				if err == nil {
					for _, p := range tc3.Packages {
						if p.Uuid == resp.Package.Uuid && p.OrgId == a.OrgId {
							check3 = true
							break
						}
					}
				} else {
					return check3, fmt.Errorf("add package story failed on getPackageForOrg. Error %v", err)
				}
			}

			if check1 && check2 && check3 {
				return true, nil
			} else {
				return false, fmt.Errorf("add package story failed. %v", nil)
			}
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			resp, ok := tc.GetData().(*ppb.AddPackageResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}
			a.packageId = resp.Package.Uuid
			tc.SaveWorkflowData(a)
			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()

			return nil
		},
	}
}

func Story_Simpool() *test.TestCase {
	return &test.TestCase{
		Name:        "Add sim pool",
		Description: "Add sim pool in an org",
		Data:        &sppb.UploadResponse{},
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			a := tc.GetWorkflowData().(*UserStoriesData)
			log.Tracef("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.simpool.upload"})
			a.reqSimPoolUploadSimReq = sapi.SimPoolUploadSimReq{
				SimType: a.simType,
				Data:    string(subscriber.CreateSimPool(MAX_POOL, &a.ICCID)),
			}

			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			var err error
			a, ok := getWorkflowData(tc)
			if ok != nil {
				return ok
			}
			tc.Data, err = a.SubscriberClient.SubscriberSimpoolUploadSims(a.reqSimPoolUploadSimReq)
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			check1, check2, check3 := false, false, false

			resp := tc.GetData().(*sppb.UploadResponse)

			if resp != nil {
				a, ok := getWorkflowData(tc)
				if ok != nil {
					return false, ok
				}

				a.reqSimByIccidReq = sapi.SimByIccidReq{
					Iccid: resp.Iccid[0],
				}

				a.reqSimPoolStatByTypeReq = sapi.SimPoolTypeReq{
					SimType: a.simType,
				}

				a.reqGetSimsByTypeReq = sapi.SimPoolTypeReq{
					SimType: a.simType,
				}

				tc1, err := a.SubscriberClient.SubscriberSimpoolGetSimByICCID(a.reqSimByIccidReq)
				if err == nil {
					for _, id := range resp.Iccid {
						if tc1.Sim.Iccid == id {
							check1 = true
							break
						}
					}
				} else {
					return check1, fmt.Errorf("upload sims story failed on getSimByIccid. Error %v", err)
				}

				tc2, err := a.SubscriberClient.SubscriberSimpoolGetSimStats(a.reqSimPoolStatByTypeReq)
				if err == nil && int(tc2.Total) >= len(resp.Iccid) {
					check2 = true
				} else {
					return check2, fmt.Errorf("upload sims story failed on getSimsStatsByType. Error %v", err)
				}

				tc3, err := a.SubscriberClient.SubscriberSimpoolGetSims(a.reqGetSimsByTypeReq)
				if err == nil {
					for _, sim := range tc3.Sims {
						if utils.Contains(resp.Iccid, sim.Iccid) {
							check3 = true
							break
						}
					}
				} else {
					return check3, fmt.Errorf("upload sims story failed on getSimsByType. Error %v", err)
				}
			}

			if check1 && check2 && check3 {
				return true, nil
			} else {
				return false, fmt.Errorf("upload sims story failed. %v", nil)
			}
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			resp, ok := tc.GetData().(*sppb.UploadResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}

			tc.SaveWorkflowData(a)
			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()

			return nil
		},
	}
}

func Story_Subscriber() *test.TestCase {
	return &test.TestCase{
		Name:        "Add subscriber",
		Description: "Add subscriber in an org",
		Data:        &srpb.AddSubscriberResponse{},
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			a := tc.GetWorkflowData().(*UserStoriesData)
			log.Tracef("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.subscriber.registry.add"})
			a.reqAddSubscriber = sapi.SubscriberAddReq{
				Dob:                   utils.GenerateRandomUTCPastDate(2005),
				Phone:                 strings.ToLower(faker.Phonenumber()),
				Email:                 strings.ToLower(faker.Email()),
				IdSerial:              faker.UUIDDigit(),
				FirstName:             faker.FirstName(),
				LastName:              faker.LastName(),
				Gender:                faker.Gender(),
				Address:               faker.Name(),
				NetworkId:             a.NetworkId,
				ProofOfIdentification: "passport",
				OrgId:                 a.OrgId,
			}

			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			var err error
			a, ok := getWorkflowData(tc)
			if ok != nil {
				return ok
			}
			tc.Data, err = a.SubscriberClient.SubscriberRegistryAddSusbscriber(a.reqAddSubscriber)
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			check1, check2 := false, false

			resp := tc.GetData().(*srpb.AddSubscriberResponse)

			if resp != nil {
				a, ok := getWorkflowData(tc)
				if ok != nil {
					return false, ok
				}

				a.reqGetSubscriber = sapi.SubscriberGetReq{
					SubscriberId: resp.Subscriber.SubscriberId,
				}

				tc1, err := a.SubscriberClient.SubscriberRegistryGetSusbscriber(a.reqGetSubscriber)
				if err == nil && tc1.Subscriber.OrgId == a.OrgId {
					check1 = true
				} else {
					return check1, fmt.Errorf("add subscriber story failed on getSubscriber. Error %v", err)
				}

				a.reqGetSubscriberSims = sapi.GetSimsBySubReq{
					SubscriberId: resp.Subscriber.SubscriberId,
				}
				tc2, err := a.SubscriberClient.SubscriberManagerGetSubscriberSims(a.reqGetSubscriberSims)
				if err == nil && len(tc2.Sims) == 0 {
					check2 = true
				} else {
					return check1, fmt.Errorf("allocate sim story failed on getSubscriberSims. Error %v", err)
				}
			}

			if check1 && check2 {
				return true, nil
			} else {
				return false, fmt.Errorf("add subscriber story failed. %v", nil)
			}
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			resp, ok := tc.GetData().(*srpb.AddSubscriberResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}
			a.subscriberId = resp.Subscriber.SubscriberId
			tc.SaveWorkflowData(a)
			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()

			return nil
		},
	}
}

func Story_Sim_Allocate() *test.TestCase {
	return &test.TestCase{
		Name:        "Allocate sim to subscriber",
		Description: "Allocate sim to subscriber in an org",
		Data:        &smpb.AllocateSimResponse{},
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			a := tc.GetWorkflowData().(*UserStoriesData)
			log.Tracef("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.subscriber.sim.allocate"})

			token, err := smutil.GenerateTokenFromIccid(a.ICCID[0], config.Key)
			if err != nil {
				return fmt.Errorf("failed to generate token from iccid %v", err)
			}
			a.reqAllocateSim = sapi.AllocateSimReq{
				SubscriberId: a.subscriberId,
				SimToken:     token,
				PackageId:    a.packageId,
				NetworkId:    a.NetworkId,
				SimType:      a.simType,
			}
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			var err error
			a, ok := getWorkflowData(tc)
			if ok != nil {
				return ok
			}
			tc.Data, err = a.SubscriberClient.SubscriberManagerAllocateSim(a.reqAllocateSim)
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			check1, check2 := false, false

			resp := tc.GetData().(*smpb.AllocateSimResponse)

			if resp != nil {
				a, ok := getWorkflowData(tc)
				if ok != nil {
					return false, ok
				}

				a.reqGetSubscriber = sapi.SubscriberGetReq{
					SubscriberId: resp.Sim.SubscriberId,
				}

				tc1, err := a.SubscriberClient.SubscriberRegistryGetSusbscriber(a.reqGetSubscriber)
				if err == nil && tc1.Subscriber.OrgId == a.OrgId {
					for _, sim := range tc1.Subscriber.Sim {
						if sim.Iccid == a.ICCID[0] {
							check1 = true
							break
						}
					}
				} else {
					return check1, fmt.Errorf("allocate sim story failed on getSubscriber. Error %v", err)
				}

				a.reqGetSubscriberSims = sapi.GetSimsBySubReq{
					SubscriberId: resp.Sim.SubscriberId,
				}
				tc2, err := a.SubscriberClient.SubscriberManagerGetSubscriberSims(a.reqGetSubscriberSims)
				if err == nil {
					for _, sim := range tc2.Sims {
						if sim.Iccid == a.ICCID[0] {
							check2 = true
							break
						}
					}
				} else {
					return check1, fmt.Errorf("allocate sim story failed on getSubscriberSims. Error %v", err)
				}
			}

			if check1 && check2 {
				return true, nil
			} else {
				return false, fmt.Errorf("allocate sim story failed. %v", nil)
			}
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			resp, ok := tc.GetData().(*smpb.AllocateSimResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}
			a.SimId = resp.Sim.Id
			tc.SaveWorkflowData(a)
			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()

			return nil
		},
	}
}

func Story_add_sim_package() *test.TestCase {
	return &test.TestCase{
		Name:        "Add new package to sim",
		Description: "add new package to user sim",
		Data:        &smpb.AddPackageResponse{},
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			a := tc.GetWorkflowData().(*UserStoriesData)
			log.Tracef("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.subscriber.sim.package.add"})

			a.reqAddPackage = dapi.AddPackageRequest{
				OwnerId:    a.OrgOwnerId,
				OrgId:      a.OrgId,
				Name:       faker.FirstName() + "-monthly-pack",
				SimType:    a.simType,
				From:       utils.GenerateUTCFutureDate(24 * time.Hour),
				To:         utils.GenerateUTCFutureDate(30 * 24 * time.Hour),
				BaserateId: a.baserateId,
				SmsVolume:  100,
				DataVolume: 1024,
				DataUnit:   "MegaBytes",
				Type:       "postpaid",
				Active:     true,
				Flatrate:   false,
				Apn:        "ukama.tel",
			}

			resp, err := a.DataplanClient.DataPlanPackageAdd(a.reqAddPackage)
			if err != nil {
				return fmt.Errorf("add package story failed at addSimPkg. Error %v", err)
			}
			a.spackageId = resp.Package.Uuid

			a.reqAddpackageToSub = sapi.AddPkgToSimReq{
				SimId:     a.SimId,
				PackageId: resp.Package.Uuid,
				StartDate: timestamppb.New(time.Now().UTC().AddDate(0, 0, 1)),
			}
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			var err error
			a, ok := getWorkflowData(tc)
			if ok != nil {
				return ok
			}
			err = a.SubscriberClient.SubscriberManagerAddPackage(a.reqAddpackageToSub)
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			check1 := false

			a, ok := getWorkflowData(tc)
			if ok != nil {
				return false, ok
			}

			a.reqGetPackagesForSim = sapi.SimReq{
				SimId: a.SimId,
			}

			tc1, err := a.SubscriberClient.SubscriberManagerGetPackageForSim(a.reqGetPackagesForSim)
			if err == nil {
				for _, p := range tc1.Packages {
					if p.PackageId == a.spackageId {
						a.SimPackageId = p.Id
						check1 = true
						break
					}
				}
			} else {
				return check1, fmt.Errorf("add package for sim story failed on getSimPackages. Error %v", err)
			}

			if check1 {
				return true, nil
			} else {
				return false, fmt.Errorf("add package for sim story failed. %v", nil)
			}
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			resp, ok := tc.GetData().(*smpb.AddPackageResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}
			tc.SaveWorkflowData(a)
			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()

			return nil
		},
	}
}

func Story_activate_sim() *test.TestCase {
	return &test.TestCase{
		Name:        "Active package for sim",
		Description: "activate sim package failed case",
		Data:        &smpb.ToggleSimStatusResponse{},
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			a := tc.GetWorkflowData().(*UserStoriesData)
			log.Tracef("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.subscriber.sim.package.active"})

			a.reqSetActivePackage = sapi.SetActivePackageForSimReq{
				SimId:     a.SimId,
				PackageId: a.SimPackageId,
			}
			a.reqToggleSimState = sapi.ActivateDeactivateSimReq{
				SimId:  a.SimId,
				Status: "active",
			}
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			var err error
			a, ok := getWorkflowData(tc)
			if ok != nil {
				return ok
			}
			_, err = a.SubscriberClient.SubscriberManagerActivatePackage(a.reqSetActivePackage)
			if err != nil && strings.Contains(err.Error(), "sim's status is is inactive") {
				tc.Data, err = a.SubscriberClient.SubscriberManagerUpdateSim(a.reqToggleSimState)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("activate sim story failed at reqSetActivePackage. %v", nil)
			}
			return nil
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			check1 := false

			resp := tc.GetData().(*smpb.ToggleSimStatusResponse)

			if resp != nil {
				a, ok := getWorkflowData(tc)
				if ok != nil {
					return false, ok
				}

				a.reqGetSimById = sapi.SimReq{
					SimId: a.SimId,
				}

				tc1, err := a.SubscriberClient.SubscriberManagerGetSim(a.reqGetSimById)
				if err == nil && tc1.Sim.Status == "active" {
					check1 = true
				} else {
					return check1, fmt.Errorf("activate sim story failed on getSimById. Error %v", err)
				}
			}

			if check1 {
				return true, nil
			} else {
				return false, fmt.Errorf("activate sim story failed. %v", nil)
			}
		},

		ExitFxn: func(ctx context.Context, tc *test.TestCase) error {
			resp, ok := tc.GetData().(*smpb.ToggleSimStatusResponse)
			if !ok {
				log.Errorf("Invalid data type for Workflow data.")

				return fmt.Errorf("invalid data type for Workflow data")
			}

			a, err := getWorkflowData(tc)
			if err != nil {
				return err
			}
			tc.SaveWorkflowData(a)
			log.Debugf("Read resp Data %v \n Written data: %v", resp, a)
			tc.Watcher.Stop()

			return nil
		},
	}
}

func Story_active_sim_package() *test.TestCase {
	return &test.TestCase{
		Name:        "Active package for sim",
		Description: "activate sim package",
		Data:        &smpb.SetActivePackageResponse{},
		SetUpFxn: func(t *testing.T, ctx context.Context, tc *test.TestCase) error {
			a := tc.GetWorkflowData().(*UserStoriesData)
			log.Tracef("Setting up watcher for %s", tc.Name)
			tc.Watcher = utils.SetupWatcher(a.MbHost, []string{"event.cloud.subscriber.sim.package.active"})

			a.reqSetActivePackage = sapi.SetActivePackageForSimReq{
				SimId:     a.SimId,
				PackageId: a.spackageId,
			}
			return nil
		},

		Fxn: func(ctx context.Context, tc *test.TestCase) error {
			var err error
			a, ok := getWorkflowData(tc)
			if ok != nil {
				return ok
			}
			tc.Data, err = a.SubscriberClient.SubscriberManagerActivatePackage(a.reqSetActivePackage)
			return err
		},

		StateFxn: func(ctx context.Context, tc *test.TestCase) (bool, error) {
			check1 := false

			resp := tc.GetData().(*smpb.SetActivePackageResponse)

			if resp != nil {
				a, ok := getWorkflowData(tc)
				if ok != nil {
					return false, ok
				}

				a.reqGetPackagesForSim = sapi.SimReq{
					SimId: a.SimId,
				}

				tc1, err := a.SubscriberClient.SubscriberManagerGetPackageForSim(a.reqGetPackagesForSim)
				if err == nil {
					for _, p := range tc1.Packages {
						if p.PackageId == a.spackageId {
							check1 = true
							break
						}
					}
				} else {
					return check1, fmt.Errorf("activate sim package story failed on getSimPackages. Error %v", err)
				}
			}

			if check1 {
				return true, nil
			} else {
				return false, fmt.Errorf("activate sim package story failed. %v", nil)
			}
		},
	}
}
