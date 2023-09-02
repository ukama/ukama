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

const VERSION = "/v1"
const ORG = "/orgs/"
const USER = "/users/"

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

	resp, err := s.r.Post(s.u.String()+VERSION+ORG, b)
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

	resp, err := s.r.Get(s.u.String() + VERSION + ORG + req.OrgName)
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

func (s *NucleusClient) AddUsrToOrg(req napi.UserOrgRequest) (*orgpb.RegisterUserResponse, error) {
	log.Debugf("Adding user to org: %v", req)
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &orgpb.RegisterUserResponse{}

	resp, err := s.r.Put(s.u.String()+VERSION+ORG+req.OrgId+USER+req.UserId, b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("AddUsrToOrg failure: %w", err)
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
		Post(s.u.String()+VERSION+USER, b)
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

func (s *NucleusClient) GetUser(req napi.GetUserRequest) (*upb.GetResponse, error) {
	rsp := &upb.GetResponse{}

	resp, err := s.r.Get(s.u.String() + VERSION + USER + req.UserId)
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

func (s *NucleusClient) GetUserByAuthId(req napi.GetUserByAuthIdRequest) (*upb.GetResponse, error) {
	rsp := &upb.GetResponse{}

	resp, err := s.r.Get(s.u.String() + VERSION + USER + req.AuthId)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("GetUserByAuthId failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *NucleusClient) Whoami(req napi.GetUserRequest) (*upb.WhoamiResponse, error) {
	rsp := &upb.WhoamiResponse{}

	resp, err := s.r.Get(s.u.String() + VERSION + USER + "whoami/" + req.UserId)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())

		return nil, fmt.Errorf("Whoami failure: %w", err)
	}

	err = jsonpb.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}
