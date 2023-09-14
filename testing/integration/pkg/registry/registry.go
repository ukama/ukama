package registry

import (
	"net/http"
	"net/url"

	api "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"
	invpb "github.com/ukama/ukama/systems/registry/invitation/pb/gen"
	mempb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	nodepb "github.com/ukama/ukama/systems/registry/node/pb/gen"

	"github.com/ukama/ukama/testing/integration/pkg/utils"
)

const VERSION = "/v1"
const MEMBERS = "/members/"
const NETWORKS = "/networks/"
const SITES = "/sites/"
const NODES = "/nodes/"
const INVITATIONS = "/invitations/"

type RegistryClient struct {
	u *url.URL
	r utils.Resty
}

func NewRegistryClient(h string) *RegistryClient {
	u, _ := url.Parse(h)
	return &RegistryClient{
		u: u,
		r: *utils.NewResty(),
	}
}

func (s *RegistryClient) AddMember(req api.MemberRequest) (*mempb.MemberResponse, error) {
	url := s.u.String() + VERSION + MEMBERS
	rsp := &mempb.MemberResponse{}

	if err := s.r.SendRequest(http.MethodPost, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) GetMember(req api.GetMemberRequest) (*mempb.MemberResponse, error) {
	url := s.u.String() + VERSION + MEMBERS + req.UserUuid
	rsp := &mempb.MemberResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) GetMembers() (*mempb.GetMembersResponse, error) {
	url := s.u.String() + VERSION + MEMBERS
	rsp := &mempb.GetMembersResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, nil, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) UpdateMember(req api.UpdateMemberRequest) error {
	url := s.u.String() + VERSION + MEMBERS + req.UserUuid
	if err := s.r.SendRequest(http.MethodPatch, url, req, nil); err != nil {
		return err
	}

	return nil
}

func (s *RegistryClient) AddNetwork(req api.AddNetworkRequest) (*netpb.AddResponse, error) {
	url := s.u.String() + VERSION + NETWORKS
	rsp := &netpb.AddResponse{}

	if err := s.r.SendRequest(http.MethodPost, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) GetNetwork(req api.GetNetworkRequest) (*netpb.GetResponse, error) {
	url := s.u.String() + VERSION + NETWORKS + req.NetworkId
	rsp := &netpb.GetResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) GetNetworks(req api.GetNetworksRequest) (*netpb.GetByOrgResponse, error) {
	url := s.u.String() + VERSION + NETWORKS + "?org=" + req.OrgUuid
	rsp := &netpb.GetByOrgResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) AddSite(req api.AddSiteRequest) (*netpb.AddSiteResponse, error) {
	url := s.u.String() + VERSION + NETWORKS + req.NetworkId + SITES
	rsp := &netpb.AddSiteResponse{}

	if err := s.r.SendRequest(http.MethodPost, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) GetSites(req api.GetNetworkRequest) (*netpb.GetSitesByNetworkResponse, error) {
	url := s.u.String() + VERSION + NETWORKS + req.NetworkId + SITES
	rsp := &netpb.GetSitesByNetworkResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, nil, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) GetSite(req api.GetSiteRequest) (*netpb.GetSiteResponse, error) {
	url := s.u.String() + VERSION + NETWORKS + req.NetworkId + SITES + req.SiteName
	rsp := &netpb.GetSiteResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, nil, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) AddNode(req api.AddNodeRequest) (*nodepb.AddNodeResponse, error) {
	url := s.u.String() + VERSION + NODES
	rsp := &nodepb.AddNodeResponse{}

	if err := s.r.SendRequest(http.MethodPost, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) UpdateNode(req api.UpdateNodeRequest) (*nodepb.UpdateNodeResponse, error) {
	url := s.u.String() + VERSION + NODES + req.NodeId
	rsp := &nodepb.UpdateNodeResponse{}

	if err := s.r.SendRequest(http.MethodPut, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) UpdateNodeState(req api.UpdateNodeStateRequest) (*nodepb.UpdateNodeResponse, error) {
	url := s.u.String() + VERSION + NODES + req.NodeId
	rsp := &nodepb.UpdateNodeResponse{}

	if err := s.r.SendRequest(http.MethodPatch, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) AddToSite(req api.AddNodeToSiteRequest) (*nodepb.AddNodeToSiteResponse, error) {
	url := s.u.String() + VERSION + NODES + req.NodeId + SITES
	rsp := &nodepb.AddNodeToSiteResponse{}

	if err := s.r.SendRequest(http.MethodPost, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) ReleaseNodeFromSite(req api.ReleaseNodeFromSiteRequest) (*nodepb.ReleaseNodeFromSiteResponse, error) {
	url := s.u.String() + VERSION + NODES + req.NodeId + SITES
	rsp := &nodepb.ReleaseNodeFromSiteResponse{}

	if err := s.r.SendRequest(http.MethodDelete, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) AttachNode(req api.AttachNodesRequest) (*nodepb.AttachNodesResponse, error) {
	url := s.u.String() + VERSION + NODES + req.ParentNode + "/attach"
	rsp := &nodepb.AttachNodesResponse{}

	if err := s.r.SendRequest(http.MethodPost, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) DetachNode(req api.DetachNodeRequest) (*nodepb.DetachNodeResponse, error) {
	url := s.u.String() + VERSION + NODES + req.NodeId + "/attach"
	rsp := &nodepb.DetachNodeResponse{}
	if err := s.r.SendRequest(http.MethodDelete, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) GetNodes(req api.GetNodesRequest) (*nodepb.GetNodesResponse, error) {
	url := s.u.String() + VERSION + NODES
	rsp := &nodepb.GetNodesResponse{}
	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) GetNodesByNetwork(req api.GetNetworkNodesRequest) (*nodepb.GetByNetworkResponse, error) {
	url := s.u.String() + VERSION + NODES + "networks/" + req.NetworkId
	rsp := &nodepb.GetByNetworkResponse{}
	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) GetNode(req api.GetNodeRequest) (*nodepb.GetNodeResponse, error) {
	url := s.u.String() + VERSION + NODES + req.NodeId
	rsp := &nodepb.GetNodeResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) GetNodesForSite(req api.GetSiteNodesRequest) (*nodepb.GetBySiteResponse, error) {
	url := s.u.String() + VERSION + NODES + "sites/" + req.SiteId
	rsp := &nodepb.GetBySiteResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) AddInvitations(req api.AddInvitationRequest) (*invpb.AddInvitationResponse, error) {
	url := s.u.String() + VERSION + INVITATIONS
	rsp := &invpb.AddInvitationResponse{}

	if err := s.r.SendRequest(http.MethodPost, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) UpdateInvitations(req api.UpdateInvitationRequest) (*invpb.UpdateInvitationStatusResponse, error) {
	url := s.u.String() + VERSION + INVITATIONS
	rsp := &invpb.UpdateInvitationStatusResponse{}

	if err := s.r.SendRequest(http.MethodPost, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) GetInvitationByOrg(req api.GetInvitationByOrgRequest) (*invpb.GetInvitationByOrgResponse, error) {
	url := s.u.String() + VERSION + INVITATIONS + "orgs/" + req.Org
	rsp := &invpb.GetInvitationByOrgResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *RegistryClient) GetInvitation(req api.GetInvitationRequest) (*invpb.GetInvitationResponse, error) {
	url := s.u.String() + VERSION + INVITATIONS + req.InvitationId
	rsp := &invpb.GetInvitationResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}
