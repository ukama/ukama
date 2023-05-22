package server

import (
	"context"
	"testing"
	"time"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	mocks "github.com/ukama/ukama/systems/data-plan/base-rate/mocks"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBaseRateService_UploadRates_Success(t *testing.T) {
	mockRepo := &mocks.BaseRateRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	rateService := NewBaseRateServer(mockRepo, msgbusClient)
	var mockSimTypeStr = string("ukama_data")
	var mockEffectiveAt = time.Now().Add(time.Hour * 24).Format(time.RFC3339)
	var mockEndAt = time.Now().Add(time.Hour * 24 * 365 * 15).Format(time.RFC3339)
	var mockFileUrl = "https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/template/template.csv"

	reqMock := &pb.UploadBaseRatesRequest{
		FileURL:     mockFileUrl,
		EffectiveAt: mockEffectiveAt,
		EndAt:       mockEndAt,
		SimType:     "ukama_data",
	}

	mockRepo.On("UploadBaseRates", mock.Anything).Return(nil)
	msgbusClient.On("PublishRequest", mock.Anything, reqMock).Return(nil).Once()

	rateRes, err := rateService.UploadBaseRates(context.Background(), reqMock)
	assert.NoError(t, err)
	for i := range rateRes.Rate {
		assert.Equal(t, rateRes.Rate[i].EffectiveAt, mockEffectiveAt)
		assert.Equal(t, rateRes.Rate[i].SimType, mockSimTypeStr)
	}
}

func TestBaseRateService_GetBaseRatesById(t *testing.T) {

	baseRateRepo := &mocks.BaseRateRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	var mockCountry = "The lunar maria"
	rateID := uuid.NewV4()
	s := NewBaseRateServer(baseRateRepo, msgbusClient)

	baseRateRepo.On("GetBaseRateById", rateID).Return(&db.BaseRate{
		Country: mockCountry,
	}, nil)
	rate, err := s.GetBaseRatesById(context.TODO(), &pb.GetBaseRatesByIdRequest{Uuid: rateID.String()})
	assert.NoError(t, err)
	assert.Equal(t, mockCountry, rate.Rate.Country)
	baseRateRepo.AssertExpectations(t)

}

func TestBaseRateService_GetBaseRatesByCountry(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	mockFilters := &pb.GetBaseRatesByCountryRequest{
		Country:     "Tycho crater",
		Provider:    "ABC Tel",
		EffectiveAt: "2022-12-01T00:00:00Z",
		SimType:     "ukama_data",
	}

	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(baseRateRepo, msgbusClient)

	baseRateRepo.On("GetBaseRatesByCountry", mockFilters.Country, mockFilters.Provider, db.ParseType("ukama_data")).Return([]db.BaseRate{
		{
			X2g:         true,
			X3g:         true,
			Apn:         "Manual entry required",
			Country:     "Tycho crater",
			Data:        0.4,
			EffectiveAt: time.Now(),
			Imsi:        1,
			Lte:         true,
			Provider:    "Multi Tel",
			SimType:     db.SimTypeUkamaData,
			SmsMo:       0.1,
			SmsMt:       0.1,
			Vpmn:        "TTC"},
	}, nil)
	rate, err := s.GetBaseRatesByCountry(context.TODO(), mockFilters)
	assert.NoError(t, err)
	assert.Equal(t, mockFilters.Country, rate.Rates[0].Country)
	assert.Equal(t, mockFilters.SimType, rate.Rates[0].SimType)
	baseRateRepo.AssertExpectations(t)
}

func TestBaseRateService_GetBaseRatesHistoryByCountry(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	mockFilters := &pb.GetBaseRatesByCountryRequest{
		Country:     "Tycho crater",
		Provider:    "ABC Tel",
		EffectiveAt: "2022-12-01T00:00:00Z",
		SimType:     "ukama_data",
	}

	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(baseRateRepo, msgbusClient)

	baseRateRepo.On("GetBaseRatesHistoryByCountry", mockFilters.Country, mockFilters.Provider, db.ParseType("ukama_data")).Return([]db.BaseRate{
		{
			X2g:         true,
			X3g:         true,
			Apn:         "Manual entry required",
			Country:     "Tycho crater",
			Data:        0.4,
			EffectiveAt: time.Now(),
			Imsi:        1,
			Lte:         true,
			Provider:    "Multi Tel",
			SimType:     db.SimTypeUkamaData,
			SmsMo:       0.1,
			SmsMt:       0.1,
			Vpmn:        "TTC"},
	}, nil)
	rate, err := s.GetBaseRatesHistoryByCountry(context.TODO(), mockFilters)
	assert.NoError(t, err)
	assert.Equal(t, mockFilters.Country, rate.Rates[0].Country)
	assert.Equal(t, mockFilters.SimType, rate.Rates[0].SimType)
	baseRateRepo.AssertExpectations(t)
}
