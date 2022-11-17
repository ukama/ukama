package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/data-plan/base-rate/mocks"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
	validations "github.com/ukama/ukama/systems/data-plan/base-rate/pkg/validations"

	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb"
)

var mockNetwork = "ABC Tel"
var mockCountry = "The lunar maria"
var mockSimType = "inter_mno_data"
var mockeEffectiveAt = "2022-12-01T00:00:00Z"
var mockFileUrl = "https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/template/template.csv"

func TestRateService_UploadRates(t *testing.T) {
	mockRepo := &mocks.BaseRateRepo{}
	mockRepo.On("UploadBaseRates", mock.Anything).Return(nil)

	rateService := NewBaseRateServer(mockRepo)

	rateReq := &pb.UploadBaseRatesRequest{
		FileURL:     mockFileUrl,
		EffectiveAt: mockeEffectiveAt,
		SimType:     validations.ReqStrTopb(mockSimType),
	}

	rateRes, err := rateService.UploadBaseRates(context.Background(), rateReq)
	assert.NoError(t, err)
	for i := range rateRes.Rate {
		assert.Equal(t, rateRes.Rate[i].EffectiveAt, rateReq.EffectiveAt)
		assert.Equal(t, rateRes.Rate[i].SimType, mockSimType)
	}
}

func TestRateService_GetRate(t *testing.T) {

	baseRateRepo := &mocks.BaseRateRepo{}
	baseRateRepo.On("GetBaseRate", uint64(1)).Return(&db.Rate{
		Country: mockCountry,
	}, nil)
	s := NewBaseRateServer(baseRateRepo)
	rate, err := s.GetBaseRate(context.TODO(), &pb.GetBaseRateRequest{RateId: uint64(1)})

	assert.NoError(t, err)
	assert.Equal(t, mockCountry, rate.Rate.Country)
	baseRateRepo.AssertExpectations(t)
}

func TestRateService_GetRates(t *testing.T) {
	var mockFilters = &pb.GetBaseRatesRequest{
		Country:     "Tycho crater",
		Provider:    "ABC Tel",
		EffectiveAt: "2022-12-01T00:00:00Z",
		SimType:     validations.ReqStrTopb("inter_mno_data"),
	}

	baseRateRepo := &mocks.BaseRateRepo{}
	baseRateRepo.On("GetBaseRates", mockFilters.Country, mockFilters.Provider, mockFilters.EffectiveAt, mockSimType).Return([]db.Rate{
		{X2g: "2G",
			X3g:          "3G",
			Apn:          "Manual entry required",
			Country:      "Tycho crater",
			Data:         "$0.4",
			Effective_at: "2023-10-10",
			Imsi:         "1",
			Lte:          "LTE",
			Network:      "Multi Tel",
			Sim_type:     "inter_mno_data",
			Sms_mo:       "$0.1",
			Sms_mt:       "$0.1",
			Vpmn:         "TTC"},
	}, nil)

	s := NewBaseRateServer(baseRateRepo)
	rate, err := s.GetBaseRates(context.TODO(), mockFilters)
	assert.NoError(t, err)
	assert.Equal(t, "Tycho crater", rate.Rates[0].Country)
	assert.Equal(t, "LTE", rate.Rates[0].Lte)
	baseRateRepo.AssertExpectations(t)
}
