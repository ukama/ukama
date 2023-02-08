package server

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/data_plan/base_rate/mocks"
	"github.com/ukama/ukama/systems/data_plan/base_rate/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ukama/ukama/systems/data_plan/base_rate/pb/gen"
)

var mockCountry = "The lunar maria"
var mockSimTypeStr = "INTER_MNO_DATA"
var mockeEffectiveAt = time.Now().Add(time.Hour * 24 * 7 * time.Duration(4)).Format(time.RFC3339Nano)
var mockFileUrl = "https://raw.githubusercontent.com/ukama/ukama/main/systems/data_plan/docs/template/template.csv"

// Start UploadRates //
// Success case
func TestRateService_UploadRates_Success(t *testing.T) {
	mockRepo := &mocks.BaseRateRepo{}
	rateService := NewBaseRateServer(mockRepo)
	reqMock := &pb.UploadBaseRatesRequest{
		FileURL:     mockFileUrl,
		EffectiveAt: mockeEffectiveAt,
		SimType:     pb.SimType_INTER_MNO_DATA,
	}

	mockRepo.On("UploadBaseRates", mock.Anything).Return(nil)
	rateRes, err := rateService.UploadBaseRates(context.Background(), reqMock)
	assert.NoError(t, err)
	for i := range rateRes.Rate {
		assert.Equal(t, rateRes.Rate[i].EffectiveAt, reqMock.EffectiveAt)
		assert.Equal(t, rateRes.Rate[i].SimType, mockSimTypeStr)
	}
}

// Error case all empty args
func TestRateService_UploadRates_Error1(t *testing.T) {
	mockRepo := &mocks.BaseRateRepo{}
	rateService := NewBaseRateServer(mockRepo)

	reqMock := &pb.UploadBaseRatesRequest{
		FileURL:     "",
		EffectiveAt: "",
		SimType:     pb.SimType_INTER_MNO_DATA,
	}

	mockRepo.On("UploadBaseRates", mock.Anything).Return(status.Errorf(codes.InvalidArgument, "invalid argument"))
	failRes1, err := rateService.UploadBaseRates(context.Background(), reqMock)
	assert.Error(t, err)
	assert.Nil(t, failRes1)
}

// Error case invalid effectiveAt
func TestRateService_UploadRates_Error2(t *testing.T) {
	mockRepo := &mocks.BaseRateRepo{}
	rateService := NewBaseRateServer(mockRepo)

	reqMock := &pb.UploadBaseRatesRequest{
		FileURL:     mockFileUrl,
		EffectiveAt: time.Now().UTC().Format(time.RFC3339),
		SimType:     pb.SimType_INTER_MNO_DATA,
	}

	mockRepo.On("UploadBaseRates", mock.Anything).Return(status.Errorf(codes.InvalidArgument, "invalid argument"))
	failRes2, err := rateService.UploadBaseRates(context.Background(), reqMock)
	assert.Error(t, err)
	assert.Nil(t, failRes2)

}

// Error case invalid url
func TestRateService_UploadRates_Error3(t *testing.T) {
	mockRepo := &mocks.BaseRateRepo{}
	rateService := NewBaseRateServer(mockRepo)

	reqMock := &pb.UploadBaseRatesRequest{
		FileURL:     "https://example",
		EffectiveAt: mockeEffectiveAt,
		SimType:     pb.SimType_INTER_MNO_DATA,
	}

	mockRepo.On("UploadBaseRates", mock.Anything).Return(status.Errorf(codes.Internal, "internal error"))
	failRes3, err := rateService.UploadBaseRates(context.Background(), reqMock)
	assert.Error(t, err)
	assert.Nil(t, failRes3)
}

// End UploadRates //

// Start GetRate //
// Success case
func TestRateService_GetRate_Success(t *testing.T) {

	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(baseRateRepo)
	var uuid = uuid.New()
	baseRateRepo.On("GetBaseRate", uuid).Return(&db.Rate{
		Country: mockCountry,
	}, nil)
	rate, err := s.GetBaseRate(context.TODO(), &pb.GetBaseRateRequest{RateUuid: uuid.String()})
	assert.NoError(t, err)
	assert.Equal(t, mockCountry, rate.Rate.Country)
	baseRateRepo.AssertExpectations(t)

}

// Error case
func TestRateService_GetRate_Error(t *testing.T) {

	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(baseRateRepo)
	var uuid = uuid.New()
	baseRateRepo.On("GetBaseRate", uuid).Return(nil, status.Errorf(codes.NotFound, "record not found"))
	_rate, err := s.GetBaseRate(context.TODO(), &pb.GetBaseRateRequest{RateUuid: uuid.String()})
	assert.Error(t, err)
	assert.Nil(t, _rate)
}

// End GetRate //

// Start GetRates //
// Success case
func TestRateService_GetRates_Success(t *testing.T) {
	mockFilters := &pb.GetBaseRatesRequest{
		Country:     "Tycho crater",
		Provider:    "ABC Tel",
		EffectiveAt: "2022-12-01T00:00:00Z",
		SimType:     pb.SimType_INTER_MNO_DATA,
	}
	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(baseRateRepo)

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
			Sim_type:     "INTER_MNO_DATA",
			Sms_mo:       "$0.1",
			Sms_mt:       "$0.1",
			Vpmn:         "TTC"},
	}, nil)
	rate, err := s.GetBaseRates(context.TODO(), mockFilters)
	assert.NoError(t, err)
	assert.Equal(t, mockFilters.Country, rate.Rates[0].Country)
	assert.Equal(t, mockFilters.SimType.String(), rate.Rates[0].SimType)
	baseRateRepo.AssertExpectations(t)
}

// Error case
func TestRateService_GetRates_Error(t *testing.T) {

	mockFilters := &pb.GetBaseRatesRequest{
		Country:     "",
		Provider:    "",
		EffectiveAt: "",
		SimType:     pb.SimType_INTER_MNO_DATA,
	}
	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(baseRateRepo)

	baseRateRepo.On("GetBaseRates", mockFilters.Country, mockFilters.Provider, mockFilters.EffectiveAt, mockSimTypeStr).Return(nil, status.Errorf(codes.NotFound, "record not found"))
	_rate, err := s.GetBaseRates(context.TODO(), mockFilters)
	assert.Error(t, err)
	assert.Nil(t, _rate)
}

// End GetRates //
