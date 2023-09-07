package registry

import (
	"encoding/json"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	api "github.com/ukama/ukama/systems/registry/api-gateway/pkg/rest"
	mempb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	jsonpb "google.golang.org/protobuf/encoding/protojson"

	"github.com/ukama/ukama/testing/integration/pkg/utils"
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

	resp, err := s.r.Get(s.u.String() + "/v1/members/" + req.UserUuid)
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
		Patch(s.u.String()+"/v1/members/"+req.UserUuid, b)
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
