package init

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
	api "github.com/ukama/ukama/systems/init/api-gateway/pkg/rest"
	lpb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"k8s.io/apimachinery/pkg/util/json"
)

type InitSys struct {
	Url *url.URL
	C   *resty.Client
}

func NewInitSys(host string) *InitSys {

	u, _ := url.Parse(host)
	c := resty.New()
	c.SetDebug(true)

	return &InitSys{
		Url: u,
		C:   c,
	}
}

func (s *InitSys) InitAddOrg(req api.AddOrgRequest) (*lpb.AddOrgResponse, error) {

	// req := api.AddOrgRequest{
	// 	OrgName:     faker.FirstName() + "_org",
	// 	Ip:          faker.IPv4(),
	// 	Certificate: util.RandomBase64(2048),
	// }

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &lpb.AddOrgResponse{}

	errStatus := &rest.ErrorResponse{}

	resp, err := s.C.R().
		SetError(errStatus).
		SetResult(rsp).
		SetBody(b).
		Put(s.Url.String() + "/v1/orgs/" + req.OrgName)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to perform operation HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return rsp, nil
}

func (s *InitSys) InitGetOrg(req api.GetOrgRequest) (*lpb.GetOrgResponse, error) {

	// req := api.AddOrgRequest{
	// 	OrgName:     faker.FirstName() + "_org",
	// 	Ip:          faker.IPv4(),
	// 	Certificate: util.RandomBase64(2048),
	// }
	rsp := &lpb.GetOrgResponse{}

	errStatus := &rest.ErrorResponse{}

	resp, err := s.C.R().
		SetError(errStatus).
		SetResult(rsp).
		Get(s.Url.String() + "/v1/orgs/" + req.OrgName)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to perform operation HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return rsp, nil
}

func (s *InitSys) InitAddSystem(req api.AddSystemRequest) (*lpb.AddSystemResponse, error) {

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &lpb.AddSystemResponse{}

	errStatus := &rest.ErrorResponse{}

	resp, err := s.C.R().
		SetError(errStatus).
		SetResult(rsp).
		SetBody(b).
		Put(s.Url.String() + "/v1/orgs/" + req.OrgName + "/systems/" + req.SysName)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to perform operation HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return rsp, nil
}

func (s *InitSys) InitGetSystem(req api.GetSystemRequest) (*lpb.GetSystemResponse, error) {

	rsp := &lpb.GetSystemResponse{}

	errStatus := &rest.ErrorResponse{}

	resp, err := s.C.R().
		SetError(errStatus).
		SetResult(rsp).
		Get(s.Url.String() + "/v1/orgs/" + req.OrgName + "/systems/" + req.SysName)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to perform operation HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return rsp, nil
}

func (s *InitSys) InitAddNode(req api.AddNodeRequest) (*lpb.AddNodeResponse, error) {

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &lpb.AddNodeResponse{}

	errStatus := &rest.ErrorResponse{}

	resp, err := s.C.R().
		SetError(errStatus).
		SetResult(rsp).
		SetBody(b).
		Put(s.Url.String() + "/v1/orgs/" + req.OrgName + "/nodes/" + req.NodeId)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to perform operation HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return rsp, nil
}

func (s *InitSys) InitGetNode(req api.GetNodeRequest) (*lpb.GetNodeResponse, error) {

	rsp := &lpb.GetNodeResponse{}

	errStatus := &rest.ErrorResponse{}

	resp, err := s.C.R().
		SetError(errStatus).
		SetResult(rsp).
		Get(s.Url.String() + "/v1/orgs/" + req.OrgName + "/nodes/" + req.NodeId)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to perform operation HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return rsp, nil
}
