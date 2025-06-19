/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

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

const OrgName = "testorg"

func TestBaseRateService_UploadRates_Success(t *testing.T) {
	mockRepo := &mocks.BaseRateRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	rateService := NewBaseRateServer(OrgName, mockRepo, msgbusClient)
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
	msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.AnythingOfType("*events.EventBaserateUploaded")).Return(nil).Once()

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
	s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

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
	s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

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
	s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

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

func TestBaseRateService_GetBaseRatesForPeriod_Success(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	// Create test data with valid RFC3339 timestamps
	fromTime := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	toTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	mockFilters := &pb.GetBaseRatesByPeriodRequest{
		Country:  "Tycho crater",
		Provider: "ABC Tel",
		SimType:  "ukama_data",
		From:     fromTime,
		To:       toTime,
	}

	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

	// Parse the times for the mock expectation
	from, _ := time.Parse(time.RFC3339, fromTime)
	to, _ := time.Parse(time.RFC3339, toTime)

	expectedRates := []db.BaseRate{
		{
			X2g:         true,
			X3g:         true,
			Apn:         "Manual entry required",
			Country:     "Tycho crater",
			Data:        0.4,
			EffectiveAt: time.Now(),
			Imsi:        1,
			Lte:         true,
			Provider:    "ABC Tel",
			SimType:     db.SimTypeUkamaData,
			SmsMo:       0.1,
			SmsMt:       0.1,
			Vpmn:        "TTC",
		},
		{
			X2g:         false,
			X3g:         true,
			Apn:         "Another APN",
			Country:     "Tycho crater",
			Data:        0.5,
			EffectiveAt: time.Now().Add(time.Hour),
			Imsi:        2,
			Lte:         true,
			Provider:    "ABC Tel",
			SimType:     db.SimTypeUkamaData,
			SmsMo:       0.15,
			SmsMt:       0.15,
			Vpmn:        "TTC2",
		},
	}

	baseRateRepo.On("GetBaseRatesForPeriod", mockFilters.Country, mockFilters.Provider, from, to, db.ParseType("ukama_data")).Return(expectedRates, nil)

	rate, err := s.GetBaseRatesForPeriod(context.TODO(), mockFilters)

	assert.NoError(t, err)
	assert.NotNil(t, rate)
	assert.Len(t, rate.Rates, 2)
	assert.Equal(t, mockFilters.Country, rate.Rates[0].Country)
	assert.Equal(t, mockFilters.Provider, rate.Rates[0].Provider)
	assert.Equal(t, mockFilters.SimType, rate.Rates[0].SimType)
	assert.Equal(t, mockFilters.Country, rate.Rates[1].Country)
	assert.Equal(t, mockFilters.Provider, rate.Rates[1].Provider)
	assert.Equal(t, mockFilters.SimType, rate.Rates[1].SimType)

	baseRateRepo.AssertExpectations(t)
}

func TestBaseRateService_GetBaseRatesForPeriod_InvalidFromTime(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	mockFilters := &pb.GetBaseRatesByPeriodRequest{
		Country:  "Tycho crater",
		Provider: "ABC Tel",
		SimType:  "ukama_data",
		From:     "invalid-time-format",
		To:       time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}

	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

	rate, err := s.GetBaseRatesForPeriod(context.TODO(), mockFilters)

	assert.Error(t, err)
	assert.Nil(t, rate)
	assert.Contains(t, err.Error(), "invalid time format for from")
}

func TestBaseRateService_GetBaseRatesForPeriod_InvalidToTime(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	mockFilters := &pb.GetBaseRatesByPeriodRequest{
		Country:  "Tycho crater",
		Provider: "ABC Tel",
		SimType:  "ukama_data",
		From:     time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		To:       "invalid-time-format",
	}

	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

	rate, err := s.GetBaseRatesForPeriod(context.TODO(), mockFilters)

	assert.Error(t, err)
	assert.Nil(t, rate)
	assert.Contains(t, err.Error(), "invalid time format for to")
}

func TestBaseRateService_GetBaseRatesForPeriod_EmptyResult(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	fromTime := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	toTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	mockFilters := &pb.GetBaseRatesByPeriodRequest{
		Country:  "Non-existent Country",
		Provider: "Non-existent Provider",
		SimType:  "ukama_data",
		From:     fromTime,
		To:       toTime,
	}

	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

	// Parse the times for the mock expectation
	from, _ := time.Parse(time.RFC3339, fromTime)
	to, _ := time.Parse(time.RFC3339, toTime)

	baseRateRepo.On("GetBaseRatesForPeriod", mockFilters.Country, mockFilters.Provider, from, to, db.ParseType("ukama_data")).Return([]db.BaseRate{}, nil)

	rate, err := s.GetBaseRatesForPeriod(context.TODO(), mockFilters)

	assert.NoError(t, err)
	assert.NotNil(t, rate)
	assert.Len(t, rate.Rates, 0)

	baseRateRepo.AssertExpectations(t)
}

func TestBaseRateService_GetBaseRatesForPackage_Success(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	// Create test data with valid RFC3339 timestamps
	fromTime := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	toTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	mockFilters := &pb.GetBaseRatesByPeriodRequest{
		Country:  "Tycho crater",
		Provider: "ABC Tel",
		SimType:  "ukama_data",
		From:     fromTime,
		To:       toTime,
	}

	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

	// Parse the times for the mock expectation
	from, _ := time.Parse(time.RFC3339, fromTime)
	to, _ := time.Parse(time.RFC3339, toTime)

	expectedRates := []db.BaseRate{
		{
			X2g:         true,
			X3g:         true,
			Apn:         "Package APN",
			Country:     "Tycho crater",
			Data:        0.4,
			EffectiveAt: time.Now(),
			Imsi:        1,
			Lte:         true,
			Provider:    "ABC Tel",
			SimType:     db.SimTypeUkamaData,
			SmsMo:       0.1,
			SmsMt:       0.1,
			Vpmn:        "TTC",
		},
		{
			X2g:         false,
			X3g:         true,
			Apn:         "Premium Package APN",
			Country:     "Tycho crater",
			Data:        0.8,
			EffectiveAt: time.Now().Add(time.Hour),
			Imsi:        2,
			Lte:         true,
			Provider:    "ABC Tel",
			SimType:     db.SimTypeUkamaData,
			SmsMo:       0.2,
			SmsMt:       0.2,
			Vpmn:        "TTC2",
		},
	}

	baseRateRepo.On("GetBaseRatesForPackage", mockFilters.Country, mockFilters.Provider, from, to, db.ParseType("ukama_data")).Return(expectedRates, nil)

	rate, err := s.GetBaseRatesForPackage(context.TODO(), mockFilters)

	assert.NoError(t, err)
	assert.NotNil(t, rate)
	assert.Len(t, rate.Rates, 2)
	assert.Equal(t, mockFilters.Country, rate.Rates[0].Country)
	assert.Equal(t, mockFilters.Provider, rate.Rates[0].Provider)
	assert.Equal(t, mockFilters.SimType, rate.Rates[0].SimType)
	assert.Equal(t, mockFilters.Country, rate.Rates[1].Country)
	assert.Equal(t, mockFilters.Provider, rate.Rates[1].Provider)
	assert.Equal(t, mockFilters.SimType, rate.Rates[1].SimType)

	baseRateRepo.AssertExpectations(t)
}

func TestBaseRateService_GetBaseRatesForPackage_InvalidFromTime(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	mockFilters := &pb.GetBaseRatesByPeriodRequest{
		Country:  "Tycho crater",
		Provider: "ABC Tel",
		SimType:  "ukama_data",
		From:     "invalid-time-format",
		To:       time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}

	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

	rate, err := s.GetBaseRatesForPackage(context.TODO(), mockFilters)

	assert.Error(t, err)
	assert.Nil(t, rate)
	assert.Contains(t, err.Error(), "invalid time format for from")
}

func TestBaseRateService_GetBaseRatesForPackage_InvalidToTime(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	mockFilters := &pb.GetBaseRatesByPeriodRequest{
		Country:  "Tycho crater",
		Provider: "ABC Tel",
		SimType:  "ukama_data",
		From:     time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		To:       "invalid-time-format",
	}

	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

	rate, err := s.GetBaseRatesForPackage(context.TODO(), mockFilters)

	assert.Error(t, err)
	assert.Nil(t, rate)
	assert.Contains(t, err.Error(), "invalid time format for to")
}

func TestBaseRateService_GetBaseRatesForPackage_EmptyResult(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	fromTime := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	toTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	mockFilters := &pb.GetBaseRatesByPeriodRequest{
		Country:  "Non-existent Country",
		Provider: "Non-existent Provider",
		SimType:  "ukama_data",
		From:     fromTime,
		To:       toTime,
	}

	baseRateRepo := &mocks.BaseRateRepo{}
	s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

	// Parse the times for the mock expectation
	from, _ := time.Parse(time.RFC3339, fromTime)
	to, _ := time.Parse(time.RFC3339, toTime)

	baseRateRepo.On("GetBaseRatesForPackage", mockFilters.Country, mockFilters.Provider, from, to, db.ParseType("ukama_data")).Return([]db.BaseRate{}, nil)

	rate, err := s.GetBaseRatesForPackage(context.TODO(), mockFilters)

	assert.NoError(t, err)
	assert.NotNil(t, rate)
	assert.Len(t, rate.Rates, 0)

	baseRateRepo.AssertExpectations(t)
}
