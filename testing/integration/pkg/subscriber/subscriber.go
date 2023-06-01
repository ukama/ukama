package subscriber

import (
	json "encoding/json"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	api "github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/rest"
	rPb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	sPb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	pPb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
	pjson "google.golang.org/protobuf/encoding/protojson"
)

type SubscriberSys struct {
	u *url.URL
	r utils.Resty
}

func NewSubscriberSys(h string) *SubscriberSys {
	u, _ := url.Parse(h)
	return &SubscriberSys{
		u: u,
		r: *utils.NewResty(),
	}

}

func (s *SubscriberSys) SubscriberSimpoolUploadSims(req api.SimPoolUploadSimReq) (*pPb.UploadResponse, error) {

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &pPb.UploadResponse{}

	resp, err := s.r.Put(s.u.String()+"/v1/simpool/"+"upload", b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = pjson.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *SubscriberSys) SubscriberSimpoolGetSimStats(req api.SimPoolStatByTypeReq) (*pPb.GetStatsResponse, error) {

	rsp := &pPb.GetStatsResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/simpool/stats/" + req.SimType)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = pjson.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *SubscriberSys) SubscriberSimpoolGetSimByICCID(req api.SimByIccidReq) (*pPb.GetByIccidResponse, error) {

	rsp := &pPb.GetByIccidResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/simpool/sim/" + req.Iccid)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = pjson.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *SubscriberSys) SubscriberRegistryGetSusbscriber(req api.SubscriberGetReq) (*rPb.GetSubscriberResponse, error) {

	rsp := &rPb.GetSubscriberResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/subscriber/" + req.SubscriberId)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = pjson.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *SubscriberSys) SubscriberRegistryAddSusbscriber(req api.SubscriberAddReq) (*rPb.AddSubscriberResponse, error) {
	log.Debugf("Adding Subscriber to subscriber registry: %+v", req)
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &rPb.AddSubscriberResponse{}

	resp, err := s.r.Put(s.u.String()+"/v1/subscriber", b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = pjson.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *SubscriberSys) SubscriberRegistryDeleteSusbscriber(req api.SubscriberDeleteReq) (*rPb.DeleteSubscriberResponse, error) {

	rsp := &rPb.DeleteSubscriberResponse{}

	_, err := s.r.Delete(s.u.String() + "/v1/subscriber/" + req.SubscriberId)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	return rsp, nil
}

func (s *SubscriberSys) SubscriberRegistryUpdateSusbscriber(req api.SubscriberUpdateReq) (*rPb.UpdateSubscriberResponse, error) {

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &rPb.UpdateSubscriberResponse{}

	_, err = s.r.Patch(s.u.String()+"/v1/subscriber", b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	return rsp, nil
}

func (s *SubscriberSys) SubscriberManagerGetSim(req api.SimReq) (*sPb.GetSimResponse, error) {

	rsp := &sPb.GetSimResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/sim/" + req.SimId)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = pjson.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *SubscriberSys) SubscriberManagerGetSubscriberSims(req api.GetSimsBySubReq) (*sPb.GetSimsBySubscriberResponse, error) {

	rsp := &sPb.GetSimsBySubscriberResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/sim/subscriber/" + req.SubscriberId)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = pjson.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *SubscriberSys) SubscriberManagerGetPackageForSim(req api.SimReq) (*sPb.GetPackagesBySimResponse, error) {

	rsp := &sPb.GetPackagesBySimResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/sim/packages/" + req.SimId)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = pjson.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *SubscriberSys) SubscriberManagerAddPackage(req api.AddPkgToSimReq) error {
	log.Tracef("Request to add pacakage: %+v", req)
	b, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	_, err = s.r.Post(s.u.String()+"/v1/sim/package", b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return err
	}

	return nil
}

func (s *SubscriberSys) SubscriberManagerAllocateSim(req api.AllocateSimReq) (*sPb.AllocateSimResponse, error) {
	log.Tracef("Allocate sim req %+v", req)
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &sPb.AllocateSimResponse{}

	resp, err := s.r.Post(s.u.String()+"/v1/sim/", b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = pjson.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *SubscriberSys) SubscriberManagerUpdateSim(req api.ActivateDeactivateSimReq) (*sPb.ToggleSimStatusResponse, error) {

	log.Tracef("SimStatus update request %+v", req)
	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &sPb.ToggleSimStatusResponse{}

	_, err = s.r.Patch(s.u.String()+"/v1/sim/"+req.SimId, b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	return rsp, nil
}

func (s *SubscriberSys) SubscriberManagerActivatePackage(req api.SetActivePackageForSimReq) (*sPb.SetActivePackageResponse, error) {

	log.Tracef("Package activation request: %+v", req)
	rsp := &sPb.SetActivePackageResponse{}

	_, err := s.r.Patch(s.u.String()+"/v1/sim/"+req.SimId+"/package/"+req.PackageId, nil)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	return rsp, nil
}

func (s *SubscriberSys) SubscriberManagerDeletePackage(req api.RemovePkgFromSimReq) error {

	log.Tracef("Package deletion request: %v", req)
	_, err := s.r.Delete(s.u.String() + "/v1/sim/" + req.SimId + "/package/" + req.PackageId)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return err
	}

	return nil
}

func (s *SubscriberSys) SubscriberManagerDeleteSim(req api.RemovePkgFromSimReq) (*sPb.DeleteSimResponse, error) {

	rsp := &sPb.DeleteSimResponse{}

	_, err := s.r.Delete(s.u.String() + "/v1/sim/" + req.SimId)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	return rsp, nil
}
