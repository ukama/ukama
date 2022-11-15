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
	mockRepo.On("UploadBaseRates", mock.Anything, mock.Anything).Return(nil)

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
	baseRateRepo := &mocks.BaseRateRepo{}
	baseRateRepo.On("GetBaseRates").Return(&db.Rate{
		Country:      mockCountry,
		Network:      mockNetwork,
		Effective_at: mockeEffectiveAt,
		Sim_type:     mockSimType,
	}, nil)

	s := NewBaseRateServer(baseRateRepo)
	rate, err := s.GetBaseRates(context.TODO(), &pb.GetBaseRatesRequest{Country: mockCountry, Provider: mockNetwork, EffectiveAt: mockeEffectiveAt, SimType: validations.ReqStrTopb(mockSimType)})
	assert.NoError(t, err)
	assert.Equal(t, mockCountry, rate.Rates[0].Country)
}
