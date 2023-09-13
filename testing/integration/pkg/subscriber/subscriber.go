package subscriber

import (
	"net/http"
	"net/url"

	api "github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/rest"
	rPb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	sPb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	pPb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
)

const VERSION = "/v1"
const SIM = "/sim/"
const SIMPOOL = "/simpool/"
const PACKAGES = "/packages/"
const SUBSCRIBER = "/subscriber/"

type SubscriberClient struct {
	u *url.URL
	r utils.Resty
}

func NewSubscriberClient(h string) *SubscriberClient {
	u, _ := url.Parse(h)
	return &SubscriberClient{
		u: u,
		r: *utils.NewResty(),
	}

}

// SIM POOL

func (s *SubscriberClient) SubscriberSimpoolUploadSims(req api.SimPoolUploadSimReq) (*pPb.UploadResponse, error) {
	url := s.u.String() + VERSION + SIMPOOL + "upload"
	rsp := &pPb.UploadResponse{}

	if err := s.r.SendRequest(http.MethodPut, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *SubscriberClient) SubscriberSimpoolGetSimStats(req api.SimPoolStatByTypeReq) (*pPb.GetStatsResponse, error) {
	url := s.u.String() + VERSION + SIMPOOL + "stats/" + req.SimType
	rsp := &pPb.GetStatsResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *SubscriberClient) SubscriberSimpoolGetSimByICCID(req api.SimByIccidReq) (*pPb.GetByIccidResponse, error) {
	url := s.u.String() + VERSION + SIMPOOL + "sim/" + req.Iccid
	rsp := &pPb.GetByIccidResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

// SUBSCRIBER

func (s *SubscriberClient) SubscriberRegistryGetSusbscriber(req api.SubscriberGetReq) (*rPb.GetSubscriberResponse, error) {
	url := s.u.String() + VERSION + SUBSCRIBER + req.SubscriberId
	rsp := &rPb.GetSubscriberResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *SubscriberClient) SubscriberRegistryAddSusbscriber(req api.SubscriberAddReq) (*rPb.AddSubscriberResponse, error) {
	url := s.u.String() + VERSION + SUBSCRIBER
	rsp := &rPb.AddSubscriberResponse{}

	if err := s.r.SendRequest(http.MethodPut, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *SubscriberClient) SubscriberRegistryDeleteSusbscriber(req api.SubscriberDeleteReq) (*rPb.DeleteSubscriberResponse, error) {
	url := s.u.String() + VERSION + SUBSCRIBER + req.SubscriberId
	rsp := &rPb.DeleteSubscriberResponse{}

	if err := s.r.SendRequest(http.MethodDelete, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *SubscriberClient) SubscriberRegistryUpdateSusbscriber(req api.SubscriberUpdateReq) (*rPb.UpdateSubscriberResponse, error) {
	url := s.u.String() + VERSION + SUBSCRIBER
	rsp := &rPb.UpdateSubscriberResponse{}

	if err := s.r.SendRequest(http.MethodPatch, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

// SIM MANAGER

func (s *SubscriberClient) SubscriberManagerGetSim(req api.SimReq) (*sPb.GetSimResponse, error) {
	url := s.u.String() + VERSION + SIM + req.SimId
	rsp := &sPb.GetSimResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *SubscriberClient) SubscriberManagerGetSubscriberSims(req api.GetSimsBySubReq) (*sPb.GetSimsBySubscriberResponse, error) {
	url := s.u.String() + VERSION + SIM + SUBSCRIBER + req.SubscriberId
	rsp := &sPb.GetSimsBySubscriberResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *SubscriberClient) SubscriberManagerGetPackageForSim(req api.SimReq) (*sPb.GetPackagesBySimResponse, error) {
	url := s.u.String() + VERSION + SIM + "packages/" + req.SimId
	rsp := &sPb.GetPackagesBySimResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *SubscriberClient) SubscriberManagerAddPackage(req api.AddPkgToSimReq) error {
	url := s.u.String() + VERSION + SIM + "package"

	if err := s.r.SendRequest(http.MethodPost, url, req, nil); err != nil {
		return err
	}

	return nil
}

func (s *SubscriberClient) SubscriberManagerAllocateSim(req api.AllocateSimReq) (*sPb.AllocateSimResponse, error) {
	url := s.u.String() + VERSION + SIM
	rsp := &sPb.AllocateSimResponse{}

	if err := s.r.SendRequest(http.MethodPost, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *SubscriberClient) SubscriberManagerUpdateSim(req api.ActivateDeactivateSimReq) (*sPb.ToggleSimStatusResponse, error) {
	url := s.u.String() + VERSION + SIM + req.SimId
	rsp := &sPb.ToggleSimStatusResponse{}

	if err := s.r.SendRequest(http.MethodPatch, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *SubscriberClient) SubscriberManagerActivatePackage(req api.SetActivePackageForSimReq) (*sPb.SetActivePackageResponse, error) {
	url := s.u.String() + VERSION + SIM + req.SimId + "/package/" + req.PackageId
	rsp := &sPb.SetActivePackageResponse{}

	if err := s.r.SendRequest(http.MethodPatch, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *SubscriberClient) SubscriberManagerDeletePackage(req api.RemovePkgFromSimReq) error {
	url := s.u.String() + VERSION + SIM + req.SimId + "/package/" + req.PackageId

	if err := s.r.SendRequest(http.MethodDelete, url, req, nil); err != nil {
		return err
	}

	return nil
}

func (s *SubscriberClient) SubscriberManagerDeleteSim(req api.RemovePkgFromSimReq) (*sPb.DeleteSimResponse, error) {
	url := s.u.String() + VERSION + SIM + req.SimId
	rsp := &sPb.DeleteSimResponse{}

	if err := s.r.SendRequest(http.MethodDelete, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}
