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

func TestRateService_UploadRates_Success(t *testing.T) {
	mockRepo := &mocks.BaseRateRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}

	rateService := NewBaseRateServer(mockRepo,msgbusClient)
	var mockSimTypeStr = db.SimTypeInterMnoData
var mockeEffectiveAt = time.Now().Add(time.Hour * 24 * 7 * time.Duration(4)).Format(time.RFC3339Nano)
var mockFileUrl = "https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/template/template.csv"

	reqMock := &pb.UploadBaseRatesRequest{
		FileURL:     mockFileUrl,
		EffectiveAt: mockeEffectiveAt,
		SimType:    string("inter_mno_data"),
	}


	mockRepo.On("UploadBaseRates", mock.Anything).Return(nil)
		msgbusClient.On("PublishRequest", mock.Anything, reqMock).Return(nil).Once()

	rateRes, err := rateService.UploadBaseRates(context.Background(), reqMock)
	assert.NoError(t, err)
	for i := range rateRes.Rate {
		assert.Equal(t, rateRes.Rate[i].EffectiveAt, reqMock.EffectiveAt)
		 assert.Equal(t, rateRes.Rate[i].SimType, mockSimTypeStr.String())
	}
}
func TestRateService_GetRate_Success(t *testing.T) {

	baseRateRepo := &mocks.BaseRateRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	var mockCountry = "The lunar maria"
rateID:=uuid.NewV4()
	s := NewBaseRateServer(baseRateRepo,msgbusClient)

	baseRateRepo.On("GetBaseRate", rateID).Return(&db.Rate{
		Country: mockCountry,
	}, nil)
	rate, err := s.GetBaseRate(context.TODO(), &pb.GetBaseRateRequest{RateID:rateID.String()})
	assert.NoError(t, err)
	assert.Equal(t, mockCountry, rate.Rate.Country)
	baseRateRepo.AssertExpectations(t)

}
func TestRateService_GetRates_Success(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	mockFilters := &pb.GetBaseRatesRequest{
		Country:     "Tycho crater",
		Provider:    "ABC Tel",
		EffectiveAt: "2022-12-01T00:00:00Z",
		SimType:    db.SimTypeInterMnoData.String(),
	}
	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(baseRateRepo,msgbusClient)

	baseRateRepo.On("GetBaseRates", mockFilters.Country, mockFilters.Provider, mockFilters.EffectiveAt, db.SimTypeInterMnoData).Return([]db.Rate{
		{X2g: "2G",
			X3g:          "3G",
			Apn:          "Manual entry required",
			Country:      "Tycho crater",
			Data:         "$0.4",
			EffectiveAt: "2023-10-10",
			Imsi:         "1",
			Lte:          "LTE",
			Network:      "Multi Tel",
			SimType:     db.SimTypeInterMnoData,
			SmsMo:       "$0.1",
			SmsMt:       "$0.1",
			Vpmn:         "TTC"},
	}, nil)
	rate, err := s.GetBaseRates(context.TODO(), mockFilters)
	assert.NoError(t, err)
	assert.Equal(t, mockFilters.Country, rate.Rates[0].Country)
	assert.Equal(t, mockFilters.SimType, rate.Rates[0].SimType)
	baseRateRepo.AssertExpectations(t)
}
