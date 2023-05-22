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

type DataPlanSys struct {
	u *url.URL
	r utils.Resty
}

func NewDataPlanSys(h string) *DataPlanSys {
	u, _ := url.Parse(h)
	return &DataPlanSys{
		u: u,
		r: *utils.NewResty(),
	}

}

func (s *DataPlanSys) DataPlanBaseRateUpload(req api.UploadBaseRatesRequest) (*bPb.UploadBaseRatesResponse, error) {

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

func (s *DataPlanSys) DataPlanBaseRateGet(req api.GetBaseRateRequest) (*bPb.GetBaseRatesByIdResponse, error) {

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

func (s *DataPlanSys) DataPlanBaseRateGetByCountry(req api.GetBaseRatesByCountryRequest) (*bPb.GetBaseRatesResponse, error) {

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &bPb.GetBaseRatesResponse{}

	resp, err := s.r.Post(s.u.String()+BASE_RATE+"/country/"+req.Country, b)

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

func (s *DataPlanSys) DataPlanBaseRateGetByPeriod(req api.GetBaseRatesForPeriodRequest) (*bPb.GetBaseRatesResponse, error) {

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &bPb.GetBaseRatesResponse{}

	resp, err := s.r.Post(s.u.String()+BASE_RATE+"/country/period", b)

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

func (s *DataPlanSys) DataPlanBaseRateGetForPackage(req api.GetBaseRatesForPeriodRequest) (*bPb.GetBaseRatesResponse, error) {

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &bPb.GetBaseRatesResponse{}

	resp, err := s.r.Post(s.u.String()+BASE_RATE+"/country/package", b)

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

func (s *DataPlanSys) DataPlanUpdateDefaultMarkup(req api.SetDefaultMarkupRequest) (*rPb.UpdateDefaultMarkupResponse, error) {

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

func (s *DataPlanSys) DataPlanGetDefaultMarkup(req api.GetDefaultMarkupRequest) (*rPb.GetDefaultMarkupResponse, error) {

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

func (s *DataPlanSys) DataPlanGetDefaultMarkupHistory(req api.GetDefaultMarkupHistoryRequest) (*rPb.GetDefaultMarkupHistoryResponse, error) {

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

func (s *DataPlanSys) DataPlanUpdateMarkup(req api.SetMarkupRequest) (*rPb.UpdateMarkupResponse, error) {

	rsp := &rPb.UpdateMarkupResponse{}

	resp, err := s.r.Post(fmt.Sprintf("%s/%f/%s/%s", s.u.String()+RATE, req.Markup, "/users/"+req.OwnerId), nil)

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

func (s *DataPlanSys) DataPlanGetDefaultUserMarkup(req api.GetMarkupRequest) (*rPb.GetMarkupResponse, error) {

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

func (s *DataPlanSys) DataPlanGetMarkupHistory(req api.GetMarkupHistoryRequest) (*rPb.GetMarkupHistoryResponse, error) {

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

func (s *DataPlanSys) DataPlanGetRate(req api.GetRateRequest) (*rPb.GetRateResponse, error) {

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	rsp := &rPb.GetRateResponse{}

	resp, err := s.r.Post(s.u.String()+"/v1/rates/users/"+req.OwnerId+"/rate", b)

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

func (s *DataPlanSys) DataPlanPackageAdd(req api.AddPackageRequest) (*pPb.AddPackageResponse, error) {

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

func (s *DataPlanSys) DataPlanPackageGetByOrg(req api.GetPackageByOrgRequest) (*pPb.GetByOrgPackageResponse, error) {

	// b, err := json.Marshal(req)
	// if err != nil {
	// 	return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	// }

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

func (s *DataPlanSys) DataPlanPackageGetById(req api.PackagesRequest) (*pPb.GetPackageResponse, error) {

	// b, err := json.Marshal(req)
	// if err != nil {
	// 	return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	// }

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

func (s *DataPlanSys) DataPlanPackageDetails(req api.PackagesRequest) (*pPb.GetPackageResponse, error) {

	// b, err := json.Marshal(req)
	// if err != nil {
	// 	return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	// }

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

func (s *DataPlanSys) DataPlanPackageUpdate(req api.UpdatePackageRequest) (*pPb.UpdatePackageResponse, error) {

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

func (s *DataPlanSys) DataPlanPackageDelete(req api.PackagesRequest) (*pPb.DeletePackageResponse, error) {

	// b, err := json.Marshal(req)
	// if err != nil {
	// 	return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	// }

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
