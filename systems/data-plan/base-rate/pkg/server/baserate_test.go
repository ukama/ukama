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
	"github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	mocks "github.com/ukama/ukama/systems/data-plan/base-rate/mocks"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const OrgName = "testorg"

// TestBaseRateService_UploadBaseRates tests the UploadBaseRates RPC
func TestBaseRateService_UploadBaseRates(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		mockRepo := &mocks.BaseRateRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewBaseRateServer(OrgName, mockRepo, msgbusClient)
		mockEffectiveAt := time.Now().Add(time.Hour * 24).Format(time.RFC3339)
		mockEndAt := time.Now().Add(time.Hour * 24 * 365 * 15).Format(time.RFC3339)
		mockFileUrl := "https://raw.githubusercontent.com/ukama/ukama/main/systems/data-plan/docs/template/template.csv"

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
		assert.NotNil(t, rateRes)
		assert.NotEmpty(t, rateRes.Rate)

		// Verify response structure
		firstRate := rateRes.Rate[0]
		assert.NotEmpty(t, firstRate.Uuid)
		assert.NotEmpty(t, firstRate.Country)
		assert.NotEmpty(t, firstRate.Provider)
		assert.Equal(t, "ukama_data", firstRate.SimType)
		assert.Equal(t, mockEffectiveAt, firstRate.EffectiveAt)
		assert.NotEmpty(t, firstRate.CreatedAt)
		assert.NotEmpty(t, firstRate.UpdatedAt)
	})

	t.Run("Validation errors", func(t *testing.T) {
		tests := []struct {
			name        string
			fileURL     string
			effectiveAt string
			endAt       string
			simType     string
			expectedErr string
		}{
			{
				name:        "Empty file URL",
				fileURL:     "",
				effectiveAt: time.Now().Add(time.Hour * 24).Format(time.RFC3339),
				endAt:       time.Now().Add(time.Hour * 24 * 365 * 15).Format(time.RFC3339),
				simType:     "ukama_data",
				expectedErr: "Please supply valid fileURL",
			},
			{
				name:        "Empty effective at",
				fileURL:     "https://example.com/file.csv",
				effectiveAt: "",
				endAt:       time.Now().Add(time.Hour * 24 * 365 * 15).Format(time.RFC3339),
				simType:     "ukama_data",
				expectedErr: "Please supply valid fileURL",
			},
			{
				name:        "Empty sim type",
				fileURL:     "https://example.com/file.csv",
				effectiveAt: time.Now().Add(time.Hour * 24).Format(time.RFC3339),
				endAt:       time.Now().Add(time.Hour * 24 * 365 * 15).Format(time.RFC3339),
				simType:     "",
				expectedErr: "invalid sim type",
			},
			{
				name:        "Invalid sim type",
				fileURL:     "https://example.com/file.csv",
				effectiveAt: time.Now().Add(time.Hour * 24).Format(time.RFC3339),
				endAt:       time.Now().Add(time.Hour * 24 * 365 * 15).Format(time.RFC3339),
				simType:     "invalid_sim_type",
				expectedErr: "invalid sim type",
			},
			{
				name:        "Sim type case mismatch",
				fileURL:     "https://example.com/file.csv",
				effectiveAt: time.Now().Add(time.Hour * 24).Format(time.RFC3339),
				endAt:       time.Now().Add(time.Hour * 24 * 365 * 15).Format(time.RFC3339),
				simType:     "UKAMA_DATA",
				expectedErr: "invalid sim type",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockRepo := &mocks.BaseRateRepo{}
				msgbusClient := &mbmocks.MsgBusServiceClient{}
				rateService := NewBaseRateServer(OrgName, mockRepo, msgbusClient)

				reqMock := &pb.UploadBaseRatesRequest{
					FileURL:     tt.fileURL,
					EffectiveAt: tt.effectiveAt,
					EndAt:       tt.endAt,
					SimType:     tt.simType,
				}

				_, err := rateService.UploadBaseRates(context.Background(), reqMock)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			})
		}
	})

	t.Run("Date validation errors", func(t *testing.T) {
		tests := []struct {
			name        string
			effectiveAt string
			endAt       string
			expectedErr string
		}{
			{
				name:        "Invalid effective at format",
				effectiveAt: "invalid-date-format",
				endAt:       time.Now().Add(time.Hour * 24 * 365 * 15).Format(time.RFC3339),
				expectedErr: "invalid date format, must be RFC3339 standard",
			},
			{
				name:        "Invalid end at format",
				effectiveAt: time.Now().Add(time.Hour * 24).Format(time.RFC3339),
				endAt:       "invalid-date-format",
				expectedErr: "invalid date format, must be RFC3339 standard",
			},
			{
				name:        "Effective at not in future",
				effectiveAt: time.Now().Add(-time.Hour * 24).Format(time.RFC3339),
				endAt:       time.Now().Add(time.Hour * 24 * 365 * 15).Format(time.RFC3339),
				expectedErr: "date is not in the future",
			},
			{
				name:        "End at not in future",
				effectiveAt: time.Now().Add(time.Hour * 24).Format(time.RFC3339),
				endAt:       time.Now().Add(-time.Hour * 24).Format(time.RFC3339),
				expectedErr: "date is not in the future",
			},
			{
				name:        "End at before effective at",
				effectiveAt: time.Now().Add(time.Hour * 24 * 2).Format(time.RFC3339),
				endAt:       time.Now().Add(time.Hour * 24).Format(time.RFC3339),
				expectedErr: "date is not after",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockRepo := &mocks.BaseRateRepo{}
				msgbusClient := &mbmocks.MsgBusServiceClient{}
				rateService := NewBaseRateServer(OrgName, mockRepo, msgbusClient)

				reqMock := &pb.UploadBaseRatesRequest{
					FileURL:     "https://example.com/file.csv",
					EffectiveAt: tt.effectiveAt,
					EndAt:       tt.endAt,
					SimType:     "ukama_data",
				}

				_, err := rateService.UploadBaseRates(context.Background(), reqMock)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			})
		}
	})

	t.Run("File fetch error", func(t *testing.T) {
		mockRepo := &mocks.BaseRateRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		rateService := NewBaseRateServer(OrgName, mockRepo, msgbusClient)

		reqMock := &pb.UploadBaseRatesRequest{
			FileURL:     "https://invalid-url-that-will-fail.com/file.csv",
			EffectiveAt: time.Now().Add(time.Hour * 24).Format(time.RFC3339),
			EndAt:       time.Now().Add(time.Hour * 24 * 365 * 15).Format(time.RFC3339),
			SimType:     "ukama_data",
		}

		_, err := rateService.UploadBaseRates(context.Background(), reqMock)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Get")
		assert.Contains(t, err.Error(), "no such host")
	})
}

func TestBaseRateService_GetBaseRatesById(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		baseRateRepo := &mocks.BaseRateRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		mockCountry := "The lunar maria"
		rateID := uuid.NewV4()
		s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

		baseRateRepo.On("GetBaseRateById", rateID).Return(&db.BaseRate{
			Country: mockCountry,
		}, nil)

		rate, err := s.GetBaseRatesById(context.TODO(), &pb.GetBaseRatesByIdRequest{Uuid: rateID.String()})
		assert.NoError(t, err)
		assert.Equal(t, mockCountry, rate.Rate.Country)
		baseRateRepo.AssertExpectations(t)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		baseRateRepo := &mocks.BaseRateRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

		_, err := s.GetBaseRatesById(context.TODO(), &pb.GetBaseRatesByIdRequest{Uuid: "invalid-uuid"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("DatabaseError", func(t *testing.T) {
		baseRateRepo := &mocks.BaseRateRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		rateID := uuid.NewV4()
		s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

		baseRateRepo.On("GetBaseRateById", rateID).Return(nil, assert.AnError)

		_, err := s.GetBaseRatesById(context.TODO(), &pb.GetBaseRatesByIdRequest{Uuid: rateID.String()})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rpc error: code = Internal desc ")
		baseRateRepo.AssertExpectations(t)
	})
}

func TestBaseRateService_GetBaseRatesByCountry(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		mockFilters := &pb.GetBaseRatesByCountryRequest{
			Country:     "Tycho crater",
			Provider:    "ABC Tel",
			EffectiveAt: "2022-12-01T00:00:00Z",
			SimType:     "ukama_data",
		}

		baseRateRepo := &mocks.BaseRateRepo{}
		s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

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
				Provider:    "Multi Tel",
				SimType:     ukama.SimTypeUkamaData,
				SmsMo:       0.1,
				SmsMt:       0.1,
				Vpmn:        "TTC",
			},
		}

		baseRateRepo.On("GetBaseRatesByCountry", mockFilters.Country, mockFilters.Provider, ukama.SimTypeUkamaData).Return(expectedRates, nil)

		rate, err := s.GetBaseRatesByCountry(context.TODO(), mockFilters)
		assert.NoError(t, err)
		assert.Equal(t, mockFilters.Country, rate.Rates[0].Country)
		assert.Equal(t, mockFilters.SimType, rate.Rates[0].SimType)
		baseRateRepo.AssertExpectations(t)
	})

	t.Run("Multiple rates success", func(t *testing.T) {
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		mockFilters := &pb.GetBaseRatesByCountryRequest{
			Country:     "Multiple Country",
			Provider:    "Multiple Provider",
			EffectiveAt: "2022-12-01T00:00:00Z",
			SimType:     "ukama_data",
		}

		baseRateRepo := &mocks.BaseRateRepo{}
		s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

		expectedRates := []db.BaseRate{
			{
				X2g:         true,
				X3g:         true,
				Apn:         "APN 1",
				Country:     "Multiple Country",
				Data:        0.4,
				EffectiveAt: time.Now(),
				Imsi:        1,
				Lte:         true,
				Provider:    "Multiple Provider",
				SimType:     ukama.SimTypeUkamaData,
				SmsMo:       0.1,
				SmsMt:       0.1,
				Vpmn:        "VPMN1",
			},
			{
				X2g:         false,
				X3g:         true,
				Apn:         "APN 2",
				Country:     "Multiple Country",
				Data:        0.5,
				EffectiveAt: time.Now().Add(time.Hour),
				Imsi:        2,
				Lte:         true,
				Provider:    "Multiple Provider",
				SimType:     ukama.SimTypeUkamaData,
				SmsMo:       0.15,
				SmsMt:       0.15,
				Vpmn:        "VPMN2",
			},
		}

		baseRateRepo.On("GetBaseRatesByCountry", mockFilters.Country, mockFilters.Provider, ukama.SimTypeUkamaData).Return(expectedRates, nil)

		rate, err := s.GetBaseRatesByCountry(context.TODO(), mockFilters)
		assert.NoError(t, err)
		assert.Len(t, rate.Rates, 2)
		assert.Equal(t, mockFilters.Country, rate.Rates[0].Country)
		assert.Equal(t, mockFilters.Country, rate.Rates[1].Country)
		assert.Equal(t, mockFilters.Provider, rate.Rates[0].Provider)
		assert.Equal(t, mockFilters.Provider, rate.Rates[1].Provider)
		baseRateRepo.AssertExpectations(t)
	})

	t.Run("Empty result", func(t *testing.T) {
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		mockFilters := &pb.GetBaseRatesByCountryRequest{
			Country:     "Non-existent Country",
			Provider:    "Non-existent Provider",
			EffectiveAt: "2022-12-01T00:00:00Z",
			SimType:     "ukama_data",
		}

		baseRateRepo := &mocks.BaseRateRepo{}
		s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

		baseRateRepo.On("GetBaseRatesByCountry", mockFilters.Country, mockFilters.Provider, ukama.SimTypeUkamaData).Return([]db.BaseRate{}, nil)

		rate, err := s.GetBaseRatesByCountry(context.TODO(), mockFilters)
		assert.NoError(t, err)
		assert.NotNil(t, rate)
		assert.Len(t, rate.Rates, 0)
		baseRateRepo.AssertExpectations(t)
	})

	t.Run("Database error", func(t *testing.T) {
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		mockFilters := &pb.GetBaseRatesByCountryRequest{
			Country:     "Error Country",
			Provider:    "Error Provider",
			EffectiveAt: "2022-12-01T00:00:00Z",
			SimType:     "ukama_data",
		}

		baseRateRepo := &mocks.BaseRateRepo{}
		s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

		baseRateRepo.On("GetBaseRatesByCountry", mockFilters.Country, mockFilters.Provider, ukama.SimTypeUkamaData).Return(nil, assert.AnError)

		rate, err := s.GetBaseRatesByCountry(context.TODO(), mockFilters)
		assert.Error(t, err)
		assert.Nil(t, rate)
		assert.Contains(t, err.Error(), "rpc error: code = Internal desc ")
		baseRateRepo.AssertExpectations(t)
	})
}

func TestBaseRateService_GetBaseRatesHistoryByCountry(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		mockFilters := &pb.GetBaseRatesByCountryRequest{
			Country:     "Tycho crater",
			Provider:    "ABC Tel",
			EffectiveAt: "2022-12-01T00:00:00Z",
			SimType:     "ukama_data",
		}

		baseRateRepo := &mocks.BaseRateRepo{}
		s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

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
				Provider:    "Multi Tel",
				SimType:     ukama.SimTypeUkamaData,
				SmsMo:       0.1,
				SmsMt:       0.1,
				Vpmn:        "TTC",
			},
		}

		baseRateRepo.On("GetBaseRatesHistoryByCountry", mockFilters.Country, mockFilters.Provider, ukama.SimTypeUkamaData).Return(expectedRates, nil)

		rate, err := s.GetBaseRatesHistoryByCountry(context.TODO(), mockFilters)
		assert.NoError(t, err)
		assert.Equal(t, mockFilters.Country, rate.Rates[0].Country)
		assert.Equal(t, mockFilters.SimType, rate.Rates[0].SimType)
		baseRateRepo.AssertExpectations(t)
	})

	t.Run("Multiple historical rates", func(t *testing.T) {
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		mockFilters := &pb.GetBaseRatesByCountryRequest{
			Country:     "Historical Country",
			Provider:    "Historical Provider",
			EffectiveAt: "2022-12-01T00:00:00Z",
			SimType:     "ukama_data",
		}

		baseRateRepo := &mocks.BaseRateRepo{}
		s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

		expectedRates := []db.BaseRate{
			{
				X2g:         true,
				X3g:         true,
				Apn:         "Historical APN 1",
				Country:     "Historical Country",
				Data:        0.3,
				EffectiveAt: time.Now().Add(-24 * time.Hour),
				Imsi:        1,
				Lte:         true,
				Provider:    "Historical Provider",
				SimType:     ukama.SimTypeUkamaData,
				SmsMo:       0.05,
				SmsMt:       0.05,
				Vpmn:        "HIST1",
			},
			{
				X2g:         true,
				X3g:         true,
				Apn:         "Historical APN 2",
				Country:     "Historical Country",
				Data:        0.4,
				EffectiveAt: time.Now().Add(-12 * time.Hour),
				Imsi:        2,
				Lte:         true,
				Provider:    "Historical Provider",
				SimType:     ukama.SimTypeUkamaData,
				SmsMo:       0.1,
				SmsMt:       0.1,
				Vpmn:        "HIST2",
			},
			{
				X2g:         true,
				X3g:         true,
				Apn:         "Historical APN 3",
				Country:     "Historical Country",
				Data:        0.5,
				EffectiveAt: time.Now(),
				Imsi:        3,
				Lte:         true,
				Provider:    "Historical Provider",
				SimType:     ukama.SimTypeUkamaData,
				SmsMo:       0.15,
				SmsMt:       0.15,
				Vpmn:        "HIST3",
			},
		}

		baseRateRepo.On("GetBaseRatesHistoryByCountry", mockFilters.Country, mockFilters.Provider, ukama.SimTypeUkamaData).Return(expectedRates, nil)

		rate, err := s.GetBaseRatesHistoryByCountry(context.TODO(), mockFilters)
		assert.NoError(t, err)
		assert.Len(t, rate.Rates, 3)
		assert.Equal(t, mockFilters.Country, rate.Rates[0].Country)
		assert.Equal(t, mockFilters.Country, rate.Rates[1].Country)
		assert.Equal(t, mockFilters.Country, rate.Rates[2].Country)
		assert.Equal(t, mockFilters.Provider, rate.Rates[0].Provider)
		assert.Equal(t, mockFilters.Provider, rate.Rates[1].Provider)
		assert.Equal(t, mockFilters.Provider, rate.Rates[2].Provider)
		baseRateRepo.AssertExpectations(t)
	})

	t.Run("Empty history result", func(t *testing.T) {
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		mockFilters := &pb.GetBaseRatesByCountryRequest{
			Country:     "No History Country",
			Provider:    "No History Provider",
			EffectiveAt: "2022-12-01T00:00:00Z",
			SimType:     "ukama_data",
		}

		baseRateRepo := &mocks.BaseRateRepo{}
		s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

		baseRateRepo.On("GetBaseRatesHistoryByCountry", mockFilters.Country, mockFilters.Provider, ukama.SimTypeUkamaData).Return([]db.BaseRate{}, nil)

		rate, err := s.GetBaseRatesHistoryByCountry(context.TODO(), mockFilters)
		assert.NoError(t, err)
		assert.NotNil(t, rate)
		assert.Len(t, rate.Rates, 0)
		baseRateRepo.AssertExpectations(t)
	})

	t.Run("Response structure validation", func(t *testing.T) {
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		mockFilters := &pb.GetBaseRatesByCountryRequest{
			Country:     "Structure Test Country",
			Provider:    "Structure Test Provider",
			EffectiveAt: "2022-12-01T00:00:00Z",
			SimType:     "ukama_data",
		}

		baseRateRepo := &mocks.BaseRateRepo{}
		s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

		expectedRates := []db.BaseRate{
			{
				X2g:         true,
				X3g:         false,
				X5g:         true,
				Apn:         "Test APN",
				Country:     "Structure Test Country",
				Data:        0.75,
				EffectiveAt: time.Now(),
				Imsi:        12345,
				Lte:         true,
				LteM:        false,
				Provider:    "Structure Test Provider",
				SimType:     ukama.SimTypeUkamaData,
				SmsMo:       0.25,
				SmsMt:       0.30,
				Vpmn:        "TEST_VPMN",
			},
		}

		baseRateRepo.On("GetBaseRatesHistoryByCountry", mockFilters.Country, mockFilters.Provider, ukama.SimTypeUkamaData).Return(expectedRates, nil)

		rate, err := s.GetBaseRatesHistoryByCountry(context.TODO(), mockFilters)
		assert.NoError(t, err)
		assert.NotNil(t, rate)
		assert.Len(t, rate.Rates, 1)

		resultRate := rate.Rates[0]
		assert.Equal(t, "Structure Test Country", resultRate.Country)
		assert.Equal(t, "Structure Test Provider", resultRate.Provider)
		assert.Equal(t, "ukama_data", resultRate.SimType)
		assert.Equal(t, "Test APN", resultRate.Apn)
		assert.Equal(t, "TEST_VPMN", resultRate.Vpmn)
		assert.Equal(t, int64(12345), resultRate.Imsi)
		assert.Equal(t, 0.75, resultRate.Data)
		assert.Equal(t, 0.25, resultRate.SmsMo)
		assert.Equal(t, 0.30, resultRate.SmsMt)
		assert.True(t, resultRate.X2G)
		assert.False(t, resultRate.X3G)
		assert.True(t, resultRate.X5G)
		assert.True(t, resultRate.Lte)
		assert.False(t, resultRate.LteM)
		assert.NotEmpty(t, resultRate.Uuid)
		assert.NotEmpty(t, resultRate.CreatedAt)
		assert.NotEmpty(t, resultRate.UpdatedAt)

		baseRateRepo.AssertExpectations(t)
	})

	t.Run("Database error", func(t *testing.T) {
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		mockFilters := &pb.GetBaseRatesByCountryRequest{
			Country:     "Error Country",
			Provider:    "Error Provider",
			EffectiveAt: "2022-12-01T00:00:00Z",
			SimType:     "ukama_data",
		}

		baseRateRepo := &mocks.BaseRateRepo{}
		s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

		baseRateRepo.On("GetBaseRatesHistoryByCountry", mockFilters.Country, mockFilters.Provider, ukama.SimTypeUkamaData).Return(nil, assert.AnError)

		rate, err := s.GetBaseRatesHistoryByCountry(context.TODO(), mockFilters)
		assert.Error(t, err)
		assert.Nil(t, rate)
		assert.Contains(t, err.Error(), "rpc error: code = Internal desc ")
		baseRateRepo.AssertExpectations(t)
	})
}

func TestBaseRateService_GetBaseRatesForPeriod(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		msgbusClient := &mbmocks.MsgBusServiceClient{}
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
				SimType:     ukama.SimTypeUkamaData,
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
				SimType:     ukama.SimTypeUkamaData,
				SmsMo:       0.15,
				SmsMt:       0.15,
				Vpmn:        "TTC2",
			},
		}

		baseRateRepo.On("GetBaseRatesForPeriod", mockFilters.Country, mockFilters.Provider, from, to, ukama.SimTypeUkamaData).Return(expectedRates, nil)

		rate, err := s.GetBaseRatesForPeriod(context.TODO(), mockFilters)
		assert.NoError(t, err)
		assert.NotNil(t, rate)
		assert.Len(t, rate.Rates, 2)
		assert.Equal(t, mockFilters.Country, rate.Rates[0].Country)
		assert.Equal(t, mockFilters.Provider, rate.Rates[0].Provider)
		assert.Equal(t, mockFilters.SimType, rate.Rates[0].SimType)
		baseRateRepo.AssertExpectations(t)
	})

	t.Run("Time validation errors", func(t *testing.T) {
		tests := []struct {
			name        string
			from        string
			to          string
			expectedErr string
		}{
			{
				name:        "Invalid from time",
				from:        "invalid-time-format",
				to:          time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				expectedErr: "invalid time format for from",
			},
			{
				name:        "Invalid to time",
				from:        time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				to:          "invalid-time-format",
				expectedErr: "invalid time format for to",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				msgbusClient := &mbmocks.MsgBusServiceClient{}
				mockFilters := &pb.GetBaseRatesByPeriodRequest{
					Country:  "Tycho crater",
					Provider: "ABC Tel",
					SimType:  "ukama_data",
					From:     tt.from,
					To:       tt.to,
				}

				baseRateRepo := &mocks.BaseRateRepo{}
				s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

				rate, err := s.GetBaseRatesForPeriod(context.TODO(), mockFilters)
				assert.Error(t, err)
				assert.Nil(t, rate)
				assert.Contains(t, err.Error(), tt.expectedErr)
			})
		}
	})

	t.Run("Empty result", func(t *testing.T) {
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

		from, _ := time.Parse(time.RFC3339, fromTime)
		to, _ := time.Parse(time.RFC3339, toTime)

		baseRateRepo.On("GetBaseRatesForPeriod", mockFilters.Country, mockFilters.Provider, from, to, ukama.SimTypeUkamaData).Return([]db.BaseRate{}, nil)

		rate, err := s.GetBaseRatesForPeriod(context.TODO(), mockFilters)
		assert.NoError(t, err)
		assert.NotNil(t, rate)
		assert.Len(t, rate.Rates, 0)
		baseRateRepo.AssertExpectations(t)
	})

	t.Run("Database error", func(t *testing.T) {
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		fromTime := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
		toTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

		mockFilters := &pb.GetBaseRatesByPeriodRequest{
			Country:  "Error Country",
			Provider: "Error Provider",
			SimType:  "ukama_data",
			From:     fromTime,
			To:       toTime,
		}

		baseRateRepo := &mocks.BaseRateRepo{}
		s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

		from, _ := time.Parse(time.RFC3339, fromTime)
		to, _ := time.Parse(time.RFC3339, toTime)

		baseRateRepo.On("GetBaseRatesForPeriod", mockFilters.Country, mockFilters.Provider, from, to, ukama.SimTypeUkamaData).Return(nil, assert.AnError)

		rate, err := s.GetBaseRatesForPeriod(context.TODO(), mockFilters)
		assert.Error(t, err)
		assert.Nil(t, rate)
		assert.Contains(t, err.Error(), "rpc error: code = Internal desc ")
		baseRateRepo.AssertExpectations(t)
	})

	t.Run("Empty time parameters", func(t *testing.T) {
		testCases := []struct {
			name        string
			from        string
			to          string
			expectedErr string
		}{
			{
				name:        "Empty from time",
				from:        "",
				to:          time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				expectedErr: "invalid time format for from",
			},
			{
				name:        "Empty to time",
				from:        time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				to:          "",
				expectedErr: "invalid time format for to",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				msgbusClient := &mbmocks.MsgBusServiceClient{}
				mockFilters := &pb.GetBaseRatesByPeriodRequest{
					Country:  "Test Country",
					Provider: "Test Provider",
					SimType:  "ukama_data",
					From:     tc.from,
					To:       tc.to,
				}

				baseRateRepo := &mocks.BaseRateRepo{}
				s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

				rate, err := s.GetBaseRatesForPeriod(context.TODO(), mockFilters)
				assert.Error(t, err)
				assert.Nil(t, rate)
				assert.Contains(t, err.Error(), tc.expectedErr)
			})
		}
	})
}

func TestBaseRateService_GetBaseRatesForPackage(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		msgbusClient := &mbmocks.MsgBusServiceClient{}
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
				SimType:     ukama.SimTypeUkamaData,
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
				SimType:     ukama.SimTypeUkamaData,
				SmsMo:       0.2,
				SmsMt:       0.2,
				Vpmn:        "TTC2",
			},
		}

		baseRateRepo.On("GetBaseRatesForPackage", mockFilters.Country, mockFilters.Provider, from, to, ukama.SimTypeUkamaData).Return(expectedRates, nil)

		rate, err := s.GetBaseRatesForPackage(context.TODO(), mockFilters)
		assert.NoError(t, err)
		assert.NotNil(t, rate)
		assert.Len(t, rate.Rates, 2)
		assert.Equal(t, mockFilters.Country, rate.Rates[0].Country)
		assert.Equal(t, mockFilters.Provider, rate.Rates[0].Provider)
		assert.Equal(t, mockFilters.SimType, rate.Rates[0].SimType)
		baseRateRepo.AssertExpectations(t)
	})

	t.Run("Time validation errors", func(t *testing.T) {
		tests := []struct {
			name        string
			from        string
			to          string
			expectedErr string
		}{
			{
				name:        "Invalid from time",
				from:        "invalid-time-format",
				to:          time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				expectedErr: "invalid time format for from",
			},
			{
				name:        "Invalid to time",
				from:        time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				to:          "invalid-time-format",
				expectedErr: "invalid time format for to",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				msgbusClient := &mbmocks.MsgBusServiceClient{}
				mockFilters := &pb.GetBaseRatesByPeriodRequest{
					Country:  "Tycho crater",
					Provider: "ABC Tel",
					SimType:  "ukama_data",
					From:     tt.from,
					To:       tt.to,
				}

				baseRateRepo := &mocks.BaseRateRepo{}
				s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

				rate, err := s.GetBaseRatesForPackage(context.TODO(), mockFilters)
				assert.Error(t, err)
				assert.Nil(t, rate)
				assert.Contains(t, err.Error(), tt.expectedErr)
			})
		}
	})

	t.Run("Empty result", func(t *testing.T) {
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

		from, _ := time.Parse(time.RFC3339, fromTime)
		to, _ := time.Parse(time.RFC3339, toTime)

		baseRateRepo.On("GetBaseRatesForPackage", mockFilters.Country, mockFilters.Provider, from, to, ukama.SimTypeUkamaData).Return([]db.BaseRate{}, nil)

		rate, err := s.GetBaseRatesForPackage(context.TODO(), mockFilters)
		assert.NoError(t, err)
		assert.NotNil(t, rate)
		assert.Len(t, rate.Rates, 0)
		baseRateRepo.AssertExpectations(t)
	})

	t.Run("Invalid from time", func(t *testing.T) {
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
	})

	t.Run("Invalid to time", func(t *testing.T) {
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
	})

	t.Run("Database error", func(t *testing.T) {
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		fromTime := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
		toTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

		mockFilters := &pb.GetBaseRatesByPeriodRequest{
			Country:  "Error Country",
			Provider: "Error Provider",
			SimType:  "ukama_data",
			From:     fromTime,
			To:       toTime,
		}

		baseRateRepo := &mocks.BaseRateRepo{}
		s := NewBaseRateServer(OrgName, baseRateRepo, msgbusClient)

		from, _ := time.Parse(time.RFC3339, fromTime)
		to, _ := time.Parse(time.RFC3339, toTime)

		baseRateRepo.On("GetBaseRatesForPackage", mockFilters.Country, mockFilters.Provider, from, to, ukama.SimTypeUkamaData).Return(nil, assert.AnError)

		rate, err := s.GetBaseRatesForPackage(context.TODO(), mockFilters)
		assert.Error(t, err)
		assert.Nil(t, rate)
		assert.Contains(t, err.Error(), "rpc error: code = Internal desc ")
		baseRateRepo.AssertExpectations(t)
	})
}
