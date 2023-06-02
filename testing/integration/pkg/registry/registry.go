package registry

import (
	"encoding/json"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	api "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	upb "github.com/ukama/ukama/systems/registry/users/pb/gen"
	jsonpb "google.golang.org/protobuf/encoding/protojson"

	"github.com/ukama/ukama/testing/integration/pkg/utils"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
)

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

func (s *RegistryClient) AddOrg(req api.AddOrgRequest) (*orgpb.AddResponse, error) {
	log.Debugf("Adding org: %v", req)
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &orgpb.AddResponse{}

	resp, err := s.r.Post(s.u.String()+"/v1/orgs", b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("AddOrg failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
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
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *RegistryClient) AddMember(req api.MemberRequest) (*orgpb.MemberResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &orgpb.MemberResponse{}

	resp, err := s.r.
		Post(s.u.String()+"/v1/orgs/"+req.OrgName+"/members", b)
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

func (s *RegistryClient) GetMember(req api.GetMemberRequest) (*orgpb.MemberResponse, error) {
	rsp := &orgpb.MemberResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/orgs/" + req.OrgName + "/members/" + req.UserUuid)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("GetMember failure: %w", err)
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
		Patch(s.u.String()+"/v1/orgs/"+req.OrgName+"/members/"+req.UserUuid, b)
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

	resp, err := s.r.Post(s.u.String()+"/v1/networks", b)
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

	resp, err := s.r.Get(s.u.String() + "/v1/networks/" + req.NetworkId)
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

func (s *RegistryClient) AddUser(req api.AddUserRequest) (*upb.AddResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &upb.AddResponse{}

	resp, err := s.r.
		Post(s.u.String()+"/v1/users", b)
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
