package dataplan

import (
	json "encoding/json"
	"fmt"
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

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &bPb.UploadBaseRatesResponse{}

	resp, err := s.r.Post(s.u.String()+BASE_RATE+"/upload", b)
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

func (s *DataplanClient) DataPlanBaseRateGet(req api.GetBaseRateRequest) (*bPb.GetBaseRatesByIdResponse, error) {

	rsp := &bPb.GetBaseRatesByIdResponse{}

	resp, err := s.r.Get(s.u.String() + BASE_RATE + "/" + req.RateId)

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

	rsp := &rPb.UpdateDefaultMarkupResponse{}

	resp, err := s.r.Post(fmt.Sprintf("%s/%f/%s", s.u.String()+RATE, req.Markup, "/default"), nil)

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

func (s *DataplanClient) DataPlanGetDefaultMarkup(req api.GetDefaultMarkupRequest) (*rPb.GetDefaultMarkupResponse, error) {

	rsp := &rPb.GetDefaultMarkupResponse{}

	resp, err := s.r.Get(s.u.String() + RATE + "/default")

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

func (s *DataplanClient) DataPlanGetDefaultMarkupHistory(req api.GetDefaultMarkupHistoryRequest) (*rPb.GetDefaultMarkupHistoryResponse, error) {

	rsp := &rPb.GetDefaultMarkupHistoryResponse{}

	resp, err := s.r.Get(s.u.String() + RATE + "/default/history")

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

func (s *DataplanClient) DataPlanUpdateMarkup(req api.SetMarkupRequest) (*rPb.UpdateMarkupResponse, error) {

	rsp := &rPb.UpdateMarkupResponse{}

	resp, err := s.r.Post(fmt.Sprintf("%s/%f/%s", s.u.String()+RATE, req.Markup, "users/"+req.OwnerId), nil)

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

func (s *DataplanClient) DataPlanGetUserMarkup(req api.GetMarkupRequest) (*rPb.GetMarkupResponse, error) {

	rsp := &rPb.GetMarkupResponse{}

	resp, err := s.r.Get(s.u.String() + RATE + "/users/" + req.OwnerId)

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

func (s *DataplanClient) DataPlanGetMarkupHistory(req api.GetMarkupHistoryRequest) (*rPb.GetMarkupHistoryResponse, error) {

	rsp := &rPb.GetMarkupHistoryResponse{}

	resp, err := s.r.Get(s.u.String() + RATE + "/users/" + req.OwnerId + "/history")

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

func (s *DataplanClient) DataPlanGetRate(req api.GetRateRequest) (*rPb.GetRateResponse, error) {

	log.Tracef("GetRate request: %+v", req)

	rsp := &rPb.GetRateResponse{}

	url := fmt.Sprintf("%s/v1/rates/users/%s/rate", s.u.String(), req.OwnerId)
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
	log.Debugf("DataplanClient AddPackageRequest by  %+v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &pPb.AddPackageResponse{}

	resp, err := s.r.Post(s.u.String()+PACKAGE, b)
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

func (s *DataplanClient) DataPlanPackageGetByOrg(req api.GetPackageByOrgRequest) (*pPb.GetByOrgPackageResponse, error) {

	log.Tracef("GetPackgesForOrg %+v", req)
	rsp := &pPb.GetByOrgPackageResponse{}

	resp, err := s.r.Get(s.u.String() + PACKAGE + "/orgs/" + req.OrgId)

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

func (s *DataplanClient) DataPlanPackageGetById(req api.PackagesRequest) (*pPb.GetPackageResponse, error) {

	rsp := &pPb.GetPackageResponse{}

	resp, err := s.r.Get(s.u.String() + PACKAGE + "/" + req.Uuid)

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

func (s *DataplanClient) DataPlanPackageDetails(req api.PackagesRequest) (*pPb.GetPackageResponse, error) {

	rsp := &pPb.GetPackageResponse{}

	resp, err := s.r.Get(s.u.String() + PACKAGE + "/" + req.Uuid + "/details")

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

func (s *DataplanClient) DataPlanPackageUpdate(req api.UpdatePackageRequest) (*pPb.UpdatePackageResponse, error) {

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &pPb.UpdatePackageResponse{}

	resp, err := s.r.Patch(s.u.String()+PACKAGE+"/"+req.Uuid, b)
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

func (s *DataplanClient) DataPlanPackageDelete(req api.PackagesRequest) (*pPb.DeletePackageResponse, error) {

	rsp := &pPb.DeletePackageResponse{}

	resp, err := s.r.Delete(s.u.String() + PACKAGE + "/" + req.Uuid)
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
