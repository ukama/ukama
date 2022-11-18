package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/data-plan/base-rate/mocks"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb"
)

var mockCountry = "The lunar maria"
var mockSimTypeStr = "inter_mno_data"
var mockeEffectiveAt = "2022-12-01T00:00:00Z"
var mockFileUrl = "https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/template/template.csv"

func TestRateService_UploadRates(t *testing.T) {
	mockRepo := &mocks.BaseRateRepo{}
	rateService := NewBaseRateServer(mockRepo)
	reqMock := &pb.UploadBaseRatesRequest{
		FileURL:     mockFileUrl,
		EffectiveAt: mockeEffectiveAt,
		SimType:     pb.SimType_inter_mno_data,
	}
	failReqMock1 := &pb.UploadBaseRatesRequest{
		FileURL:     "",
		EffectiveAt: "",
		SimType:     pb.SimType_inter_mno_data,
	}
	failReqMock2 := &pb.UploadBaseRatesRequest{
		FileURL:     mockFileUrl,
		EffectiveAt: time.Now().UTC().Format(time.RFC3339),
		SimType:     pb.SimType_inter_mno_data,
	}
	failReqMock3 := &pb.UploadBaseRatesRequest{
		FileURL:     "https://example",
		EffectiveAt: mockeEffectiveAt,
		SimType:     pb.SimType_inter_mno_data,
	}

	//Success case
	mockRepo.On("UploadBaseRates", mock.Anything).Return(nil)
	rateRes, err := rateService.UploadBaseRates(context.Background(), reqMock)
	assert.NoError(t, err)
	for i := range rateRes.Rate {
		assert.Equal(t, rateRes.Rate[i].EffectiveAt, reqMock.EffectiveAt)
		assert.Equal(t, rateRes.Rate[i].SimType, mockSimTypeStr)
	}

	//Error case all empty args
	mockRepo.On("UploadBaseRates", mock.Anything).Return(status.Errorf(codes.InvalidArgument, "invalid argument"))
	failRes1, err := rateService.UploadBaseRates(context.Background(), failReqMock1)
	assert.Error(t, err)
	assert.Nil(t, failRes1)

	//Error case invalid effectiveAt
	mockRepo.On("UploadBaseRates", mock.Anything).Return(status.Errorf(codes.InvalidArgument, "invalid argument"))
	failRes2, err := rateService.UploadBaseRates(context.Background(), failReqMock2)
	assert.Error(t, err)
	assert.Nil(t, failRes2)

	//Error case invalid url
	mockRepo.On("UploadBaseRates", mock.Anything).Return(status.Errorf(codes.Internal, "internal error"))
	failRes3, err := rateService.UploadBaseRates(context.Background(), failReqMock3)
	assert.Error(t, err)
	assert.Nil(t, failRes3)
}

func TestRateService_GetRate(t *testing.T) {

	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(baseRateRepo)

	//Success case
	baseRateRepo.On("GetBaseRate", uint64(1)).Return(&db.Rate{
		Country: mockCountry,
	}, nil)
	rate, err := s.GetBaseRate(context.TODO(), &pb.GetBaseRateRequest{RateId: uint64(1)})
	assert.NoError(t, err)
	assert.Equal(t, mockCountry, rate.Rate.Country)
	baseRateRepo.AssertExpectations(t)

	//Error case
	baseRateRepo.On("GetBaseRate", uint64(0)).Return(nil, status.Errorf(codes.NotFound, "record not found"))
	_rate, err := s.GetBaseRate(context.TODO(), &pb.GetBaseRateRequest{RateId: uint64(0)})
	assert.Error(t, err)
	assert.Nil(t, _rate)
}

func TestRateService_GetRates(t *testing.T) {
	mockFilters := &pb.GetBaseRatesRequest{
		Country:     "Tycho crater",
		Provider:    "ABC Tel",
		EffectiveAt: "2022-12-01T00:00:00Z",
		SimType:     pb.SimType_inter_mno_data,
	}
	emptyMockFilters := &pb.GetBaseRatesRequest{
		Country:     "",
		Provider:    "",
		EffectiveAt: "",
		SimType:     pb.SimType_inter_mno_data,
	}
	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(baseRateRepo)

	// Success case
	baseRateRepo.On("GetBaseRates", mockFilters.Country, mockFilters.Provider, mockFilters.EffectiveAt, mockSimTypeStr).Return([]db.Rate{
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
	rate, err := s.GetBaseRates(context.TODO(), mockFilters)
	assert.NoError(t, err)
	assert.Equal(t, mockFilters.Country, rate.Rates[0].Country)
	assert.Equal(t, mockFilters.SimType.String(), rate.Rates[0].SimType)
	baseRateRepo.AssertExpectations(t)

	// Error case
	baseRateRepo.On("GetBaseRates", emptyMockFilters.Country, emptyMockFilters.Provider, emptyMockFilters.EffectiveAt, mockSimTypeStr).Return(nil, status.Errorf(codes.NotFound, "record not found"))
	_rate, err := s.GetBaseRates(context.TODO(), emptyMockFilters)
	assert.Error(t, err)
	assert.Nil(t, _rate)
}
