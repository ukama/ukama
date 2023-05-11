package init

import (
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	api "github.com/ukama/ukama/systems/init/api-gateway/pkg/rest"
	lpb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"github.com/ukama/ukama/testing/integration/pkg/util"
	"k8s.io/apimachinery/pkg/util/json"
)

type InitSys struct {
	u *url.URL
	r util.Resty
}

func NewInitSys(h string) *InitSys {
	u, _ := url.Parse(h)
	return &InitSys{
		u: u,
		r: *util.NewResty(),
	}

}
func (s *InitSys) InitAddOrg(req api.AddOrgRequest) (*lpb.AddOrgResponse, error) {

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &lpb.AddOrgResponse{}

	resp, err := s.r.Put(s.u.String()+"/v1/orgs/"+req.OrgName, b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *InitSys) InitGetOrg(req api.GetOrgRequest) (*lpb.GetOrgResponse, error) {

	rsp := &lpb.GetOrgResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/orgs/" + req.OrgName)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *InitSys) InitAddSystem(req api.AddSystemRequest) (*lpb.AddSystemResponse, error) {

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &lpb.AddSystemResponse{}

	resp, err := s.r.Put(s.u.String()+"/v1/orgs/"+req.OrgName+"/systems/"+req.SysName, b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *InitSys) InitGetSystem(req api.GetSystemRequest) (*lpb.GetSystemResponse, error) {

	rsp := &lpb.GetSystemResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/orgs/" + req.OrgName + "/systems/" + req.SysName)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *InitSys) InitAddNode(req api.AddNodeRequest) (*lpb.AddNodeResponse, error) {

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &lpb.AddNodeResponse{}

	resp, err := s.r.Put(s.u.String()+"/v1/orgs/"+req.OrgName+"/nodes/"+req.NodeId, b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *InitSys) InitGetNode(req api.GetNodeRequest) (*lpb.GetNodeResponse, error) {

	rsp := &lpb.GetNodeResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/orgs/" + req.OrgName + "/nodes/" + req.NodeId)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}
