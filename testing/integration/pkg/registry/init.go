package pkg

import (
	"encoding/json"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	api "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"
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

func (s *RegistryClient) AddOrg(req api.AddOrgRequest) (*orgpb.AddResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &orgpb.AddResponse{}

	resp, err := s.r.Post(b, s.u.String()+"/v1/orgs")
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

func (s *RegistryClient) GetOrg(req api.GetOrgRequest) (*orgpb.GetResponse, error) {
	rsp := &orgpb.GetResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/orgs/" + req.OrgName)
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

func (s *RegistryClient) AddMember(req api.MemberRequest) (*orgpb.MemberResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &orgpb.MemberResponse{}

	resp, err := s.r.
		Post(b, s.u.String()+"/v1/orgs/"+req.OrgName+"/members")
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

func (s *RegistryClient) GetMember(req api.GetMemberRequest) (*orgpb.MemberResponse, error) {
	rsp := &orgpb.MemberResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/orgs/" + req.OrgName + "/members/" + req.UserUuid)

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

func (s *RegistryClient) UpdateMember(req api.UpdateMemberRequest) error {
	b, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	_, err = s.r.
		Patch(b, s.u.String()+"/v1/orgs/"+req.OrgName+"/members/"+req.UserUuid)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return err
	}

	return nil
}

func (s *RegistryClient) AddNetwork(req api.AddNetworkRequest) (*netpb.AddResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &netpb.AddResponse{}

	resp, err := s.r.Post(b, s.u.String()+"/v1/networks")
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
