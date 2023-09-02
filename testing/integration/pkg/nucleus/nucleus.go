package nucleus

import (
	json "encoding/json"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	napi "github.com/ukama/ukama/systems/nucleus/api-gateway/pkg/rest"
	orgpb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"
	upb "github.com/ukama/ukama/systems/nucleus/user/pb/gen"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
)

const ORG = "/v1/orgs"
const USER = "/v1/users"

type NucleusClient struct {
	u *url.URL
	r utils.Resty
}

func NewNucleusClient(h string) *NucleusClient {
	u, _ := url.Parse(h)
	return &NucleusClient{
		u: u,
		r: *utils.NewResty(),
	}
}

func (s *NucleusClient) AddOrg(req napi.AddOrgRequest) (*orgpb.AddResponse, error) {
	log.Debugf("Adding org: %v", req)
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &orgpb.AddResponse{}

	resp, err := s.r.Post(s.u.String()+ORG, b)
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

func (s *NucleusClient) GetOrg(req napi.GetOrgRequest) (*orgpb.GetResponse, error) {
	rsp := &orgpb.GetResponse{}

	resp, err := s.r.Get(s.u.String() + ORG + req.OrgName)
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

func (s *NucleusClient) AddUser(req napi.AddUserRequest) (*upb.AddResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &upb.AddResponse{}

	resp, err := s.r.
		Post(s.u.String()+USER, b)
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
