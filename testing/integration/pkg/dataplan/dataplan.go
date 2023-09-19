package dataplan

import (
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
	api "github.com/ukama/ukama/systems/data-plan/api-gateway/pkg/rest"
	bPb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	pPb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	rPb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
	pjson "google.golang.org/protobuf/encoding/protojson"
)

const BASE_RATE = "/v1/baserates"
const RATE = "/v1/markup"
const PACKAGE = "/v1/packages"

type DataplanClient struct {
	u *url.URL
	r utils.Resty
}

func NewDataplanClient(h string) *DataplanClient {
	u, _ := url.Parse(h)
	return &DataplanClient{
		u: u,
		r: *utils.NewResty(),
	}

}

func (s *DataplanClient) DataPlanBaseRateUpload(req api.UploadBaseRatesRequest) (*bPb.UploadBaseRatesResponse, error) {
	url := s.u.String() + BASE_RATE + "/upload"
	rsp := &bPb.UploadBaseRatesResponse{}

	if err := s.r.SendRequest(http.MethodPost, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *DataplanClient) DataPlanBaseRateGet(req api.GetBaseRateRequest) (*bPb.GetBaseRatesByIdResponse, error) {
	url := s.u.String() + BASE_RATE + "/" + req.RateId
	rsp := &bPb.GetBaseRatesByIdResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *DataplanClient) DataPlanBaseRateGetByCountry(req api.GetBaseRatesByCountryRequest) (*bPb.GetBaseRatesResponse, error) {

	rsp := &bPb.GetBaseRatesResponse{}
	q := fmt.Sprintf("country=%s&sim_type=%s", req.Country, req.SimType)
	if len(req.Provider) > 0 {
		q = fmt.Sprintf("%s&provider=%s", q, req.Provider)
	}

	resp, err := s.r.GetWithQuery(s.u.String()+BASE_RATE, q)

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

func (s *DataplanClient) DataPlanBaseRateGetByPeriod(req api.GetBaseRatesForPeriodRequest) (*bPb.GetBaseRatesResponse, error) {

	log.Debugf("DataplanClient GetBaseRate by period %v", req)

	rsp := &bPb.GetBaseRatesResponse{}

	q := fmt.Sprintf("country=%s&sim_type=%s&provider=%s&from=%s&to=%s", req.Country, req.SimType, req.Provider, req.From, req.To)
	resp, err := s.r.GetWithQuery(s.u.String()+BASE_RATE+"/period", q)

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

func (s *DataplanClient) DataPlanBaseRateGetForPackage(req api.GetBaseRatesForPeriodRequest) (*bPb.GetBaseRatesResponse, error) {

	rsp := &bPb.GetBaseRatesResponse{}

	q := fmt.Sprintf("country=%s&sim_type=%s&provider=%s&from=%s&to=%s", req.Country, req.SimType, req.Provider, req.From, req.To)
	resp, err := s.r.GetWithQuery(s.u.String()+"/package", q)

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

func (s *DataplanClient) DataPlanUpdateDefaultMarkup(req api.SetDefaultMarkupRequest) (*rPb.UpdateDefaultMarkupResponse, error) {
	url := fmt.Sprintf("%s/%f/%s", s.u.String()+RATE, req.Markup, "/default")
	rsp := &rPb.UpdateDefaultMarkupResponse{}

	if err := s.r.SendRequest(http.MethodPost, url, nil, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *DataplanClient) DataPlanGetDefaultMarkup(req api.GetDefaultMarkupRequest) (*rPb.GetDefaultMarkupResponse, error) {
	url := s.u.String() + RATE + "/default"
	rsp := &rPb.GetDefaultMarkupResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *DataplanClient) DataPlanGetDefaultMarkupHistory(req api.GetDefaultMarkupHistoryRequest) (*rPb.GetDefaultMarkupHistoryResponse, error) {
	url := s.u.String() + RATE + "/default/history"
	rsp := &rPb.GetDefaultMarkupHistoryResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *DataplanClient) DataPlanUpdateMarkup(req api.SetMarkupRequest) (*rPb.UpdateMarkupResponse, error) {
	url := fmt.Sprintf("%s/%f/%s", s.u.String()+RATE, req.Markup, "users/"+req.OwnerId)
	rsp := &rPb.UpdateMarkupResponse{}

	if err := s.r.SendRequest(http.MethodPost, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *DataplanClient) DataPlanGetUserMarkup(req api.GetMarkupRequest) (*rPb.GetMarkupResponse, error) {
	url := s.u.String() + RATE + "/users/" + req.OwnerId
	rsp := &rPb.GetMarkupResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *DataplanClient) DataPlanGetMarkupHistory(req api.GetMarkupHistoryRequest) (*rPb.GetMarkupHistoryResponse, error) {
	url := s.u.String() + RATE + "/users/" + req.OwnerId + "/history"
	rsp := &rPb.GetMarkupHistoryResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *DataplanClient) DataPlanGetRate(req api.GetRateRequest) (*rPb.GetRateResponse, error) {

	log.Tracef("GetRate request: %+v", req)

	rsp := &rPb.GetRateResponse{}

	url := fmt.Sprintf("%s/v1/rates/users/%s/rate", s.u.String(), req.UserId)
	q := fmt.Sprintf("country=%s&sim_type=%s&provider=%s&from=%s&to=%s", req.Country, req.SimType, req.Provider, req.From, req.To)
	resp, err := s.r.GetWithQuery(url, q)

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

func (s *DataplanClient) DataPlanPackageAdd(req api.AddPackageRequest) (*pPb.AddPackageResponse, error) {
	url := s.u.String() + PACKAGE
	rsp := &pPb.AddPackageResponse{}

	if err := s.r.SendRequest(http.MethodPost, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *DataplanClient) DataPlanPackageGetByOrg(req api.GetPackageByOrgRequest) (*pPb.GetByOrgPackageResponse, error) {
	url := s.u.String() + PACKAGE + "/orgs/" + req.OrgId
	rsp := &pPb.GetByOrgPackageResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *DataplanClient) DataPlanPackageGetById(req api.PackagesRequest) (*pPb.GetPackageResponse, error) {
	url := s.u.String() + PACKAGE + "/" + req.Uuid
	rsp := &pPb.GetPackageResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *DataplanClient) DataPlanPackageDetails(req api.PackagesRequest) (*pPb.GetPackageResponse, error) {
	url := s.u.String() + PACKAGE + "/" + req.Uuid + "/details"
	rsp := &pPb.GetPackageResponse{}

	if err := s.r.SendRequest(http.MethodGet, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *DataplanClient) DataPlanPackageUpdate(req api.UpdatePackageRequest) (*pPb.UpdatePackageResponse, error) {
	url := s.u.String() + PACKAGE + "/" + req.Uuid
	rsp := &pPb.UpdatePackageResponse{}

	if err := s.r.SendRequest(http.MethodPatch, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

func (s *DataplanClient) DataPlanPackageDelete(req api.PackagesRequest) (*pPb.DeletePackageResponse, error) {
	url := s.u.String() + PACKAGE + "/" + req.Uuid
	rsp := &pPb.DeletePackageResponse{}

	if err := s.r.SendRequest(http.MethodDelete, url, req, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}
