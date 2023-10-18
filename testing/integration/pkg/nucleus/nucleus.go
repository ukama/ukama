package nucleus

import (
	"net/http"
	"net/url"

	napi "github.com/ukama/ukama/systems/nucleus/api-gateway/pkg/rest"
	orgpb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"
	upb "github.com/ukama/ukama/systems/nucleus/user/pb/gen"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
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

	url := s.u.String() + VERSION + ORG
	rsp := &orgpb.AddResponse{}

	if err := s.r.SendRequest(http.MethodPost, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *NucleusClient) GetOrg(req napi.GetOrgRequest) (*orgpb.GetResponse, error) {
	url := s.u.String() + VERSION + ORG + req.OrgName
	rsp := &orgpb.GetResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, nil, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *NucleusClient) AddUsrToOrg(req napi.UserOrgRequest) (*orgpb.RegisterUserResponse, error) {
	url := s.u.String() + VERSION + ORG + req.OrgId + USER + req.UserId
	rsp := &orgpb.RegisterUserResponse{}

	if err := s.r.SendRequest(http.MethodPut, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *NucleusClient) RemoveUsrFromOrg(req napi.UserOrgRequest) (*orgpb.RemoveOrgForUserResponse, error) {
	url := s.u.String() + VERSION + ORG + req.OrgId + USER + req.UserId
	rsp := &orgpb.RemoveOrgForUserResponse{}

	if err := s.r.SendRequest(http.MethodDelete, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *NucleusClient) AddUser(req napi.AddUserRequest) (*upb.AddResponse, error) {
	url := s.u.String() + VERSION + USER
	rsp := &upb.AddResponse{}

	if err := s.r.SendRequest(http.MethodPost, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *NucleusClient) GetUser(req napi.GetUserRequest) (*upb.GetResponse, error) {
	url := s.u.String() + VERSION + USER + req.UserId
	rsp := &upb.GetResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, nil, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *NucleusClient) GetUserByAuthId(req napi.GetUserByAuthIdRequest) (*upb.GetResponse, error) {
	url := s.u.String() + VERSION + USER + "auth/" + req.AuthId
	rsp := &upb.GetResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, nil, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *NucleusClient) Whoami(req napi.GetUserRequest) (*upb.WhoamiResponse, error) {
	url := s.u.String() + VERSION + USER + "whoami/" + req.UserId
	rsp := &upb.WhoamiResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, nil, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}
