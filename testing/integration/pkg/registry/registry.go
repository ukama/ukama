package registry

import (
	"encoding/json"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	api "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"
	mempb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	nodepb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	jsonpb "google.golang.org/protobuf/encoding/protojson"

	"github.com/ukama/ukama/testing/integration/pkg/utils"
)

const VERSION = "/v1"
const MEMBERS = "/members/"
const NETWORKS = "/networks/"
const SITES = "/sites/"
const NODES = "/nodes/"

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
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &mempb.MemberResponse{}

	resp, err := s.r.
		Post(s.u.String()+"/v1/members", b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("AddMember failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) GetMember(req api.GetMemberRequest) (*mempb.MemberResponse, error) {

	rsp := &mempb.MemberResponse{}

	resp, err := s.r.Get(s.u.String() + VERSION + MEMBERS + req.UserUuid)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("GetUser failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) GetMembers() (*mempb.GetMembersResponse, error) {

	rsp := &mempb.GetMembersResponse{}

	resp, err := s.r.Get(s.u.String() + VERSION + MEMBERS)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("GetUser failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) UpdateMember(req api.UpdateMemberRequest) error {
	b, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	_, err = s.r.
		Patch(s.u.String()+VERSION+MEMBERS+req.UserUuid, b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return fmt.Errorf("UpdateMember failure: %w", err)
	}

	return nil
}

func (s *RegistryClient) AddNetwork(req api.AddNetworkRequest) (*netpb.AddResponse, error) {
	log.Debugf("Adding network: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &netpb.AddResponse{}

	resp, err := s.r.Post(s.u.String()+VERSION+NETWORKS, b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("AddNetwork failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) GetNetwork(req api.GetNetworkRequest) (*netpb.GetResponse, error) {
	rsp := &netpb.GetResponse{}

	resp, err := s.r.Get(s.u.String() + VERSION + NETWORKS + req.NetworkId)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) GetNetworks(req api.GetNetworksRequest) (*netpb.GetByOrgResponse, error) {
	rsp := &netpb.GetByOrgResponse{}

	resp, err := s.r.Get(s.u.String() + VERSION + NETWORKS)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) AddSite(req api.AddSiteRequest) (*netpb.AddSiteResponse, error) {
	log.Debugf("Adding site: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &netpb.AddSiteResponse{}

	resp, err := s.r.Post(s.u.String()+VERSION+NETWORKS+req.NetworkId+SITES+req.SiteName, b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("adding site failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) GetSites(req api.GetNetworkRequest) (*netpb.GetSitesByNetworkResponse, error) {
	rsp := &netpb.GetSitesByNetworkResponse{}

	resp, err := s.r.Get(s.u.String() + VERSION + NETWORKS + req.NetworkId + SITES)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) GetSite(req api.GetSiteRequest) (*netpb.GetSiteResponse, error) {
	rsp := &netpb.GetSiteResponse{}

	resp, err := s.r.Get(s.u.String() + VERSION + NETWORKS + req.NetworkId + SITES + req.SiteName)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) AddNode(req api.AddNodeRequest) (*nodepb.AddNodeResponse, error) {
	log.Debugf("Adding node: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &nodepb.AddNodeResponse{}

	resp, err := s.r.Post(s.u.String()+VERSION+NODES, b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("adding node failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) UpdateNode(req api.UpdateNodeRequest) (*nodepb.UpdateNodeResponse, error) {
	log.Debugf("Update node: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &nodepb.UpdateNodeResponse{}

	resp, err := s.r.Put(s.u.String()+VERSION+NODES+req.NodeId, b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("update node failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) UpdateNodeState(req api.UpdateNodeStateRequest) (*nodepb.UpdateNodeResponse, error) {
	log.Debugf("Update node state: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &nodepb.UpdateNodeResponse{}

	resp, err := s.r.Patch(s.u.String()+VERSION+NODES+req.NodeId, b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("update node state failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) AddToSite(req api.AddNodeToSiteRequest) (*nodepb.AddNodeToSiteResponse, error) {
	log.Debugf("Adding node to site: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &nodepb.AddNodeToSiteResponse{}

	resp, err := s.r.Post(s.u.String()+VERSION+NODES+req.NodeId+SITES, b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("adding node to site failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) ReleaseNodeFromSite(req api.ReleaseNodeFromSiteRequest) (*nodepb.ReleaseNodeFromSiteResponse, error) {
	log.Debugf("Release from site: %v", req)

	_, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &nodepb.ReleaseNodeFromSiteResponse{}

	resp, err := s.r.Delete(s.u.String() + VERSION + NODES + req.NodeId + SITES)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("release from site failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) AttachNode(req api.AttachNodesRequest) (*nodepb.AttachNodesResponse, error) {
	log.Debugf("Adding node: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &nodepb.AttachNodesResponse{}

	resp, err := s.r.Post(s.u.String()+VERSION+NODES+req.ParentNode+"/attach", b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("attach node failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) DetachNode(req api.DetachNodeRequest) (*nodepb.DetachNodeResponse, error) {
	log.Debugf("Detach node: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &nodepb.DetachNodeResponse{}

	resp, err := s.r.Delete(s.u.String()+VERSION+NODES+req.NodeId+"/attach", b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("attach node failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) GetNodes(req api.GetNodesRequest) (*nodepb.GetNodesResponse, error) {
	rsp := &nodepb.GetNodesResponse{}

	resp, err := s.r.Get(s.u.String() + VERSION + NODES)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) GetNode(req api.GetNodeRequest) (*nodepb.GetNodeResponse, error) {
	rsp := &nodepb.GetNodeResponse{}

	resp, err := s.r.Get(s.u.String() + VERSION + NODES + req.NodeId)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) GetNodesForSite(req api.GetSiteNodesRequest) (*nodepb.GetBySiteResponse, error) {
	rsp := &nodepb.GetBySiteResponse{}

	resp, err := s.r.Get(s.u.String() + VERSION + NODES + SITES + req.SiteId)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}
