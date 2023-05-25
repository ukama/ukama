package pkg

import (
	"encoding/json"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	api "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	userpb "github.com/ukama/ukama/systems/registry/users/pb/gen"
	jsonpb "google.golang.org/protobuf/encoding/protojson"

	"github.com/ukama/ukama/testing/integration/pkg/util"
)

type RegistryClient struct {
	u *url.URL
	r util.Resty
}

func NewRegistryClient(h string) *RegistryClient {
	u, _ := url.Parse(h)

	return &RegistryClient{
		u: u,
		r: *util.NewResty(),
	}
}

func (s *RegistryClient) AddUser(req api.AddUserRequest) (*userpb.AddResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("AddUser: request marshal error. error: %w", err)
	}

	rsp := &userpb.AddResponse{}

	resp, err := s.r.Post(b, s.u.String()+"/v1/users")
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("AddUser failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("AddUser: response unmarshal error. error: %w", err)
	}

	return rsp, nil
}

func (s *RegistryClient) AddOrg(req api.AddOrgRequest) (*orgpb.AddResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("AddOrg: request marshal error. error: %w", err)
	}

	rsp := &orgpb.AddResponse{}

	resp, err := s.r.Post(b, s.u.String()+"/v1/orgs")
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("AddOrg failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("AddOrg: response unmarshal error. error: %w", err)
	}

	return rsp, nil
}

func (s *RegistryClient) GetOrg(req api.GetOrgRequest) (*orgpb.GetResponse, error) {
	rsp := &orgpb.GetResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/orgs/" + req.OrgName)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("GetOrg failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("GetOrg: response unmarshal error. error: %w", err)
	}

	return rsp, nil
}

func (s *RegistryClient) AddMember(req api.MemberRequest) (*orgpb.MemberResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("AddMember: request marshal error. error: %w", err)
	}

	rsp := &orgpb.MemberResponse{}

	resp, err := s.r.
		Post(b, s.u.String()+"/v1/orgs/"+req.OrgName+"/members")
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("AddMember failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("AddMember: response unmarshal error. error: %w", err)
	}

	return rsp, nil
}

func (s *RegistryClient) GetMember(req api.GetMemberRequest) (*orgpb.MemberResponse, error) {
	rsp := &orgpb.MemberResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/orgs/" + req.OrgName + "/members/" + req.UserUuid)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("GetMember failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("GetMember: response unmarshal error. error: %w", err)
	}

	return rsp, nil
}

func (s *RegistryClient) UpdateMember(req api.UpdateMemberRequest) error {
	b, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("UpdateMember: request marshal error. error: %w", err)
	}

	_, err = s.r.
		Patch(b, s.u.String()+"/v1/orgs/"+req.OrgName+"/members/"+req.UserUuid)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return fmt.Errorf("UpdateMember failure: %w", err)
	}

	return nil
}

func (s *RegistryClient) AddNetwork(req api.AddNetworkRequest) (*netpb.AddResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("AddNetwork: request marshal error. error: %w", err)
	}

	rsp := &netpb.AddResponse{}

	resp, err := s.r.Post(b, s.u.String()+"/v1/networks")
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("AddNetwork failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("AddNetwork: response unmarshal error. error: %w", err)
	}

	return rsp, nil
}

func (s *RegistryClient) GetNetwork(req api.GetNetworkRequest) (*netpb.GetResponse, error) {
	rsp := &netpb.GetResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/networks/" + req.NetworkId)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("GetNetwork failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("GetNetwork: response unmarshal error. error: %w", err)
	}

	return rsp, nil
}
