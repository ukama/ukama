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
	"errors"
	"testing"
	"time"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/common/validation"
	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	mocks "github.com/ukama/ukama/systems/data-plan/rate/mocks"
	pb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	"github.com/ukama/ukama/systems/data-plan/rate/pkg/db"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	splmocks "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen/mocks"
)

const OrgName = "testorg"

const (
	TestMarkup10 = 10.0
	TestMarkup5  = 5.0

	TestDataRate  = 0.0014
	TestSmsMoRate = 0.0100
	TestSmsMtRate = 0.0001
	TestImsi      = 1

	TestCountry  = "USA"
	TestProvider = "Ukama"
	TestSimType  = "ukama_data"

	TestApn  = "Manual entry required"
	TestVpmn = "TTC"

	TestFromDate = "2033-04-20T20:31:24-00:00"
	TestToDate   = "2043-04-20T20:31:24-00:00"

	ErrInvalidUUID        = "Error parsing UUID"
	ErrDatabaseConnection = "database connection failed"
	ErrDatabaseError      = "database error"
	ErrMessageBusError    = "message bus error"
	ErrRecordNotFound     = "record not found"
	ErrBaseRateClient     = "base rate client error"
	ErrBaseRates          = "base rates error"
	ErrNoValidBaseRates   = "no valid base rates found"
	ErrMarkupDatabase     = "markup database error"
	ErrBaseRateService    = "base rate service error"
	ErrInvalidDateFormat  = "invalid date format for"

	TestInvalidUUID = "invalid-uuid"
)

var (
	TestCreatedTime = time.Date(2021, 11, 12, 11, 45, 26, 371000000, time.UTC)
	TestUpdatedTime = time.Date(2022, 10, 12, 11, 45, 26, 371000000, time.UTC)
	TestDeletedTime = time.Date(2022, 11, 12, 11, 45, 26, 371000000, time.UTC)
	TestDeletedAt   = gorm.DeletedAt{Time: TestDeletedTime, Valid: true}
)

type testSetup struct {
	markupRepo     *mocks.MarkupsRepo
	defMarkupRepo  *mocks.DefaultMarkupRepo
	baserateSvc    *mocks.BaserateClientProvider
	msgbusClient   *mbmocks.MsgBusServiceClient
	rateService    *RateServer
	baserateClient *splmocks.BaseRatesServiceClient
}

func newTestSetup() *testSetup {
	markupRepo := &mocks.MarkupsRepo{}
	defMarkupRepo := &mocks.DefaultMarkupRepo{}
	baserateSvc := &mocks.BaserateClientProvider{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

	return &testSetup{
		markupRepo:    markupRepo,
		defMarkupRepo: defMarkupRepo,
		baserateSvc:   baserateSvc,
		msgbusClient:  msgbusClient,
		rateService:   rateService,
	}
}

func (ts *testSetup) setupBaseRateClient() *splmocks.BaseRatesServiceClient {
	baserateClient := &splmocks.BaseRatesServiceClient{}
	ts.baserateSvc.On("GetClient").Return(baserateClient, nil)
	ts.baserateClient = baserateClient
	return baserateClient
}

func (ts *testSetup) assertAllExpectations(t *testing.T) {
	ts.markupRepo.AssertExpectations(t)
	ts.defMarkupRepo.AssertExpectations(t)
	ts.baserateSvc.AssertExpectations(t)
	ts.msgbusClient.AssertExpectations(t)
	if ts.baserateClient != nil {
		ts.baserateClient.AssertExpectations(t)
	}
}

func TestRateService_GetMarkup(t *testing.T) {
	tests := []struct {
		name           string
		ownerId        string
		setupMocks     func(*testSetup, string)
		expectedError  bool
		expectedMarkup float64
		errorContains  string
	}{
		{
			name:    "MarkupforOwnerIdExists",
			ownerId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				markups := &db.Markups{
					OwnerId: ownerUUID,
					Markup:  TestMarkup10,
				}
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(markups, nil)
			},
			expectedError:  false,
			expectedMarkup: TestMarkup10,
		},
		{
			name:    "MarkupforOwnerIdDoesn'tExists",
			ownerId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				defMarkup := &db.DefaultMarkup{Markup: TestMarkup5}
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(nil, gorm.ErrRecordNotFound)
				ts.defMarkupRepo.On("GetDefaultMarkupRate").Return(defMarkup, nil)
			},
			expectedError:  false,
			expectedMarkup: TestMarkup5,
		},
		{
			name:    "InvalidUUID",
			ownerId: TestInvalidUUID,
			setupMocks: func(ts *testSetup, ownerId string) {
			},
			expectedError: true,
			errorContains: ErrInvalidUUID,
		},
		{
			name:    "DatabaseErrorWhenGettingMarkup",
			ownerId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				dbError := errors.New(ErrDatabaseConnection)
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(nil, dbError)
			},
			expectedError: true,
			errorContains: ErrDatabaseConnection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestSetup()
			tt.setupMocks(ts, tt.ownerId)

			req := &pb.GetMarkupRequest{OwnerId: tt.ownerId}
			rateRes, err := ts.rateService.GetMarkup(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, rateRes)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, rateRes)
				assert.Equal(t, tt.expectedMarkup, rateRes.Markup)
				assert.Equal(t, tt.ownerId, rateRes.OwnerId)
			}

			ts.assertAllExpectations(t)
		})
	}
}

func TestRateService_UpdateDefaultMarkup(t *testing.T) {
	tests := []struct {
		name          string
		markup        float64
		setupMocks    func(*testSetup)
		expectedError bool
		errorContains string
	}{
		{
			name:   "UpdateDefaultMarkupSuccess",
			markup: TestMarkup5,
			setupMocks: func(ts *testSetup) {
				ts.defMarkupRepo.On("UpdateDefaultMarkupRate", TestMarkup5).Return(nil)
				ts.msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "UpdateDefaultMarkup_DatabaseError",
			markup: TestMarkup5,
			setupMocks: func(ts *testSetup) {
				ts.defMarkupRepo.On("UpdateDefaultMarkupRate", TestMarkup5).Return(errors.New(ErrDatabaseError))
			},
			expectedError: true,
			errorContains: ErrDatabaseError,
		},
		{
			name:   "UpdateDefaultMarkup_MessageBusError",
			markup: TestMarkup5,
			setupMocks: func(ts *testSetup) {
				ts.defMarkupRepo.On("UpdateDefaultMarkupRate", TestMarkup5).Return(nil)
				ts.msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(errors.New(ErrMessageBusError))
			},
			expectedError: false,
		},
		{
			name:   "UpdateDefaultMarkup_RecordNotFound",
			markup: TestMarkup5,
			setupMocks: func(ts *testSetup) {
				ts.defMarkupRepo.On("UpdateDefaultMarkupRate", TestMarkup5).Return(gorm.ErrRecordNotFound)
			},
			expectedError: true,
			errorContains: ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestSetup()
			tt.setupMocks(ts)

			req := &pb.UpdateDefaultMarkupRequest{Markup: tt.markup}
			_, err := ts.rateService.UpdateDefaultMarkup(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			ts.assertAllExpectations(t)
		})
	}
}

func TestRateService_GetDefaultMarkup(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*testSetup)
		expectedError  bool
		expectedMarkup float64
		errorContains  string
	}{
		{
			name: "GetDefaultMarkupSuccess",
			setupMocks: func(ts *testSetup) {
				defMarkup := &db.DefaultMarkup{Markup: TestMarkup5}
				ts.defMarkupRepo.On("GetDefaultMarkupRate").Return(defMarkup, nil)
			},
			expectedError:  false,
			expectedMarkup: TestMarkup5,
		},
		{
			name: "GetDefaultMarkup_DatabaseError",
			setupMocks: func(ts *testSetup) {
				ts.defMarkupRepo.On("GetDefaultMarkupRate").Return(nil, errors.New(ErrDatabaseError))
			},
			expectedError: true,
			errorContains: ErrDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestSetup()
			tt.setupMocks(ts)

			req := &pb.GetDefaultMarkupRequest{}
			rateRes, err := ts.rateService.GetDefaultMarkup(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, rateRes)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMarkup, rateRes.Markup)
			}

			ts.assertAllExpectations(t)
		})
	}
}

func TestRateService_GetDefaultMarkupHistory(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(*testSetup)
		expectedError bool
		errorContains string
	}{
		{
			name: "GetDefaultMarkupHistorySuccess",
			setupMocks: func(ts *testSetup) {
				defMarkup := []db.DefaultMarkup{
					{
						Model: gorm.Model{
							ID:        1,
							CreatedAt: TestCreatedTime,
							DeletedAt: TestDeletedAt,
							UpdatedAt: TestUpdatedTime,
						},
						Markup: TestMarkup5,
					},
					{
						Model: gorm.Model{
							ID:        2,
							CreatedAt: TestCreatedTime,
							UpdatedAt: TestUpdatedTime,
						},
						Markup: TestMarkup5,
					},
				}
				ts.defMarkupRepo.On("GetDefaultMarkupRateHistory").Return(defMarkup, nil)
			},
			expectedError: false,
		},
		{
			name: "GetDefaultMarkupHistory_DatabaseError",
			setupMocks: func(ts *testSetup) {
				ts.defMarkupRepo.On("GetDefaultMarkupRateHistory").Return(nil, errors.New(ErrDatabaseError))
			},
			expectedError: true,
			errorContains: ErrDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestSetup()
			tt.setupMocks(ts)

			req := &pb.GetDefaultMarkupHistoryRequest{}
			rateRes, err := ts.rateService.GetDefaultMarkupHistory(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, rateRes)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, rateRes)
				for _, rate := range rateRes.MarkupRates {
					assert.Equal(t, TestMarkup5, rate.Markup)
					assert.Equal(t, TestCreatedTime.Format(time.RFC3339), rate.CreatedAt)
				}
			}

			ts.assertAllExpectations(t)
		})
	}
}

func TestRateService_GetRate(t *testing.T) {
	tests := []struct {
		name          string
		req           *pb.GetRateRequest
		setupMocks    func(*testSetup, *pb.GetRateRequest)
		expectedError bool
		errorContains string
	}{
		{
			name: "GetRate_Success",
			req: &pb.GetRateRequest{
				OwnerId:  uuid.NewV4().String(),
				Country:  TestCountry,
				Provider: TestProvider,
				SimType:  TestSimType,
				From:     TestFromDate,
				To:       TestToDate,
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				markups := &db.Markups{OwnerId: ownerId, Markup: TestMarkup10}
				ts.markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)

				baserateClient := ts.setupBaseRateClient()
				to, _ := validation.FromString(req.To)
				from, _ := validation.FromString(req.From)

				baserateClient.On("GetBaseRatesForPeriod", mock.Anything, &bpb.GetBaseRatesByPeriodRequest{
					Country:  req.Country,
					Provider: req.Provider,
					To:       to.Format(time.RFC3339),
					From:     from.Format(time.RFC3339),
					SimType:  req.SimType,
				}).Return(&bpb.GetBaseRatesResponse{
					Rates: []*bpb.Rate{
						{
							X2G:         true,
							X3G:         true,
							Apn:         TestApn,
							Country:     req.Country,
							Data:        TestDataRate,
							EffectiveAt: "2033-04-20T20:31:24+00:00",
							Imsi:        TestImsi,
							Lte:         true,
							Provider:    "Multi Tel",
							SimType:     req.SimType,
							SmsMo:       TestSmsMoRate,
							SmsMt:       TestSmsMtRate,
							Vpmn:        TestVpmn,
						},
					},
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "GetRate_InvalidUUID",
			req: &pb.GetRateRequest{
				OwnerId:  TestInvalidUUID,
				Country:  TestCountry,
				Provider: TestProvider,
				SimType:  TestSimType,
				From:     TestFromDate,
				To:       TestToDate,
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				// No mocks needed for invalid UUID
			},
			expectedError: true,
			errorContains: ErrInvalidUUID,
		},
		{
			name: "GetRate_InvalidDateFormatTo",
			req: &pb.GetRateRequest{
				OwnerId:  uuid.NewV4().String(),
				Country:  TestCountry,
				Provider: TestProvider,
				SimType:  TestSimType,
				From:     TestFromDate,
				To:       "invalid-date-format",
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				markups := &db.Markups{OwnerId: ownerId, Markup: TestMarkup10}
				ts.markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)
			},
			expectedError: true,
			errorContains: ErrInvalidDateFormat + " to",
		},
		{
			name: "GetRate_InvalidDateFormatFrom",
			req: &pb.GetRateRequest{
				OwnerId:  uuid.NewV4().String(),
				Country:  TestCountry,
				Provider: TestProvider,
				SimType:  TestSimType,
				From:     "invalid-date-format",
				To:       TestToDate,
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				markups := &db.Markups{OwnerId: ownerId, Markup: TestMarkup10}
				ts.markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)
			},
			expectedError: true,
			errorContains: ErrInvalidDateFormat + " from",
		},
		{
			name: "GetRate_BaseRateClientError",
			req: &pb.GetRateRequest{
				OwnerId:  uuid.NewV4().String(),
				Country:  TestCountry,
				Provider: TestProvider,
				SimType:  TestSimType,
				From:     TestFromDate,
				To:       TestToDate,
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				markups := &db.Markups{OwnerId: ownerId, Markup: TestMarkup10}
				ts.markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)
				ts.baserateSvc.On("GetClient").Return(nil, errors.New(ErrBaseRateClient))
			},
			expectedError: true,
			errorContains: ErrBaseRateClient,
		},
		{
			name: "GetRate_BaseRatesError",
			req: &pb.GetRateRequest{
				OwnerId:  uuid.NewV4().String(),
				Country:  TestCountry,
				Provider: TestProvider,
				SimType:  TestSimType,
				From:     TestFromDate,
				To:       TestToDate,
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				markups := &db.Markups{OwnerId: ownerId, Markup: TestMarkup10}
				ts.markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)

				baserateClient := ts.setupBaseRateClient()
				baserateClient.On("GetBaseRatesForPeriod", mock.Anything, mock.Anything).Return(nil, errors.New(ErrBaseRates))
			},
			expectedError: true,
			errorContains: ErrBaseRates,
		},
		{
			name: "GetRate_NoValidBaseRates",
			req: &pb.GetRateRequest{
				OwnerId:  uuid.NewV4().String(),
				Country:  TestCountry,
				Provider: TestProvider,
				SimType:  TestSimType,
				From:     TestFromDate,
				To:       TestToDate,
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				markups := &db.Markups{OwnerId: ownerId, Markup: TestMarkup10}
				ts.markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)

				baserateClient := ts.setupBaseRateClient()
				baserateClient.On("GetBaseRatesForPeriod", mock.Anything, mock.Anything).Return(&bpb.GetBaseRatesResponse{
					Rates: []*bpb.Rate{},
				}, nil)
			},
			expectedError: true,
			errorContains: ErrNoValidBaseRates,
		},
		{
			name: "GetRate_NilBaseRates",
			req: &pb.GetRateRequest{
				OwnerId:  uuid.NewV4().String(),
				Country:  TestCountry,
				Provider: TestProvider,
				SimType:  TestSimType,
				From:     TestFromDate,
				To:       TestToDate,
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				markups := &db.Markups{OwnerId: ownerId, Markup: TestMarkup10}
				ts.markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)

				baserateClient := ts.setupBaseRateClient()
				baserateClient.On("GetBaseRatesForPeriod", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedError: true,
			errorContains: ErrNoValidBaseRates,
		},
		{
			name: "GetRate_UserMarkupError",
			req: &pb.GetRateRequest{
				OwnerId:  uuid.NewV4().String(),
				Country:  TestCountry,
				Provider: TestProvider,
				SimType:  TestSimType,
				From:     TestFromDate,
				To:       TestToDate,
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				ts.markupRepo.On("GetMarkupRate", ownerId).Return(nil, errors.New(ErrMarkupDatabase))
			},
			expectedError: true,
			errorContains: ErrMarkupDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestSetup()
			tt.setupMocks(ts, tt.req)

			rateRes, err := ts.rateService.GetRate(context.Background(), tt.req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, rateRes)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, rateRes)
				for _, r := range rateRes.Rates {
					assert.Equal(t, tt.req.Country, r.Country)
					assert.Equal(t, tt.req.SimType, r.SimType)
				}
			}

			ts.assertAllExpectations(t)
		})
	}
}

func TestRateService_UpdateMarkup(t *testing.T) {
	tests := []struct {
		name          string
		ownerId       string
		markup        float64
		setupMocks    func(*testSetup, string, float64)
		expectedError bool
		errorContains string
	}{
		{
			name:    "UpdateMarkupSuccess",
			ownerId: uuid.NewV4().String(),
			markup:  TestMarkup10,
			setupMocks: func(ts *testSetup, ownerId string, markup float64) {
				ownerUUID, _ := uuid.FromString(ownerId)
				ts.markupRepo.On("UpdateMarkupRate", ownerUUID, markup).Return(nil)
				ts.msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "UpdateMarkup_InvalidUUID",
			ownerId: TestInvalidUUID,
			markup:  TestMarkup10,
			setupMocks: func(ts *testSetup, ownerId string, markup float64) {
				// No mocks needed for invalid UUID
			},
			expectedError: true,
			errorContains: ErrInvalidUUID,
		},
		{
			name:    "UpdateMarkup_DatabaseError",
			ownerId: uuid.NewV4().String(),
			markup:  TestMarkup10,
			setupMocks: func(ts *testSetup, ownerId string, markup float64) {
				ownerUUID, _ := uuid.FromString(ownerId)
				ts.markupRepo.On("UpdateMarkupRate", ownerUUID, markup).Return(errors.New(ErrDatabaseError))
			},
			expectedError: true,
			errorContains: ErrDatabaseError,
		},
		{
			name:    "UpdateMarkup_MessageBusError",
			ownerId: uuid.NewV4().String(),
			markup:  TestMarkup10,
			setupMocks: func(ts *testSetup, ownerId string, markup float64) {
				ownerUUID, _ := uuid.FromString(ownerId)
				ts.markupRepo.On("UpdateMarkupRate", ownerUUID, markup).Return(nil)
				ts.msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(errors.New(ErrMessageBusError))
			},
			expectedError: false, // Message bus errors are logged but don't fail the operation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestSetup()
			tt.setupMocks(ts, tt.ownerId, tt.markup)

			req := &pb.UpdateMarkupRequest{
				OwnerId: tt.ownerId,
				Markup:  tt.markup,
			}
			_, err := ts.rateService.UpdateMarkup(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			ts.assertAllExpectations(t)
		})
	}
}

func TestRateService_DeleteMarkup(t *testing.T) {
	tests := []struct {
		name          string
		ownerId       string
		setupMocks    func(*testSetup, string)
		expectedError bool
		errorContains string
	}{
		{
			name:    "DeleteMarkupSuccess",
			ownerId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				ts.markupRepo.On("DeleteMarkupRate", ownerUUID).Return(nil)
				ts.msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "DeleteMarkup_InvalidUUID",
			ownerId: TestInvalidUUID,
			setupMocks: func(ts *testSetup, ownerId string) {
				// No mocks needed for invalid UUID
			},
			expectedError: true,
			errorContains: ErrInvalidUUID,
		},
		{
			name:    "DeleteMarkup_DatabaseError",
			ownerId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				ts.markupRepo.On("DeleteMarkupRate", ownerUUID).Return(errors.New(ErrDatabaseError))
			},
			expectedError: true,
			errorContains: ErrDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestSetup()
			tt.setupMocks(ts, tt.ownerId)

			req := &pb.DeleteMarkupRequest{OwnerId: tt.ownerId}
			_, err := ts.rateService.DeleteMarkup(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			ts.assertAllExpectations(t)
		})
	}
}

func TestRateService_GetMarkupVal(t *testing.T) {
	t.Run("GetMarkupSuccess", func(t *testing.T) {
		ts := newTestSetup()
		ownerId := uuid.NewV4()
		markup := &db.Markups{
			OwnerId: ownerId,
			Markup:  TestMarkup10,
		}

		ts.markupRepo.On("GetMarkupRate", ownerId).Return(markup, nil)

		req := &pb.GetMarkupRequest{OwnerId: ownerId.String()}
		rateRes, err := ts.rateService.GetMarkup(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, rateRes)
		assert.Equal(t, markup.Markup, rateRes.Markup)
		assert.Equal(t, markup.OwnerId.String(), rateRes.OwnerId)

		ts.assertAllExpectations(t)
	})
}

func TestRateService_GetMarkupHistory(t *testing.T) {
	tests := []struct {
		name          string
		ownerId       string
		setupMocks    func(*testSetup, string)
		expectedError bool
		errorContains string
	}{
		{
			name:    "GetMarkupHistorySuccess",
			ownerId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				markup := []db.Markups{
					{
						Model: gorm.Model{
							ID:        1,
							CreatedAt: TestCreatedTime,
							DeletedAt: TestDeletedAt,
							UpdatedAt: TestUpdatedTime,
						},
						OwnerId: ownerUUID,
						Markup:  TestMarkup5,
					},
					{
						Model: gorm.Model{
							ID:        2,
							CreatedAt: TestCreatedTime,
							UpdatedAt: TestUpdatedTime,
						},
						OwnerId: ownerUUID,
						Markup:  TestMarkup10,
					},
				}
				ts.markupRepo.On("GetMarkupRateHistory", ownerUUID).Return(markup, nil)
			},
			expectedError: false,
		},
		{
			name:    "GetMarkupHistory_InvalidUUID",
			ownerId: TestInvalidUUID,
			setupMocks: func(ts *testSetup, ownerId string) {
				// No mocks needed for invalid UUID
			},
			expectedError: true,
			errorContains: ErrInvalidUUID,
		},
		{
			name:    "GetMarkupHistory_DatabaseError",
			ownerId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				ts.markupRepo.On("GetMarkupRateHistory", ownerUUID).Return(nil, errors.New(ErrDatabaseError))
			},
			expectedError: true,
			errorContains: ErrDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestSetup()
			tt.setupMocks(ts, tt.ownerId)

			req := &pb.GetMarkupHistoryRequest{OwnerId: tt.ownerId}
			rateRes, err := ts.rateService.GetMarkupHistory(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, rateRes)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, rateRes)
				for i, rate := range rateRes.MarkupRates {
					expectedMarkup := TestMarkup5 + float64(i)*TestMarkup5 // 5, 10
					assert.Equal(t, expectedMarkup, rate.Markup)
					assert.Equal(t, TestCreatedTime.Format(time.RFC3339), rate.CreatedAt)
				}
			}

			ts.assertAllExpectations(t)
		})
	}
}

func TestRateService_GetRateById(t *testing.T) {
	tests := []struct {
		name          string
		ownerId       string
		baseRateId    string
		setupMocks    func(*testSetup, string, string)
		expectedError bool
		errorContains string
	}{
		{
			name:       "GetRateById_Success",
			ownerId:    uuid.NewV4().String(),
			baseRateId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string, baseRateId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				baseRateUUID, _ := uuid.FromString(baseRateId)

				markups := &db.Markups{OwnerId: ownerUUID, Markup: TestMarkup10}
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(markups, nil)

				baserateClient := ts.setupBaseRateClient()

				baseRate := &bpb.Rate{
					X2G:         true,
					X3G:         true,
					Apn:         TestApn,
					Country:     TestCountry,
					Data:        TestDataRate,
					EffectiveAt: "2033-04-20T20:31:24+00:00",
					Imsi:        TestImsi,
					Lte:         true,
					Provider:    "Multi Tel",
					SimType:     TestSimType,
					SmsMo:       TestSmsMoRate,
					SmsMt:       TestSmsMtRate,
					Vpmn:        TestVpmn,
				}

				baserateClient.On("GetBaseRatesById", mock.Anything, &bpb.GetBaseRatesByIdRequest{
					Uuid: baseRateUUID.String(),
				}).Return(&bpb.GetBaseRatesByIdResponse{
					Rate: baseRate,
				}, nil)
			},
			expectedError: false,
		},
		{
			name:       "GetRateById_InvalidOwnerId",
			ownerId:    TestInvalidUUID,
			baseRateId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string, baseRateId string) {
				// No mocks needed for invalid UUID
			},
			expectedError: true,
			errorContains: ErrInvalidUUID,
		},
		{
			name:       "GetRateById_BaseRateClientError",
			ownerId:    uuid.NewV4().String(),
			baseRateId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string, baseRateId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				markups := &db.Markups{OwnerId: ownerUUID, Markup: TestMarkup10}
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(markups, nil)
				ts.baserateSvc.On("GetClient").Return(nil, errors.New(ErrBaseRateClient))
			},
			expectedError: true,
			errorContains: ErrBaseRateClient,
		},
		{
			name:       "GetRateById_BaseRateServiceError",
			ownerId:    uuid.NewV4().String(),
			baseRateId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string, baseRateId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				baseRateUUID, _ := uuid.FromString(baseRateId)

				markups := &db.Markups{OwnerId: ownerUUID, Markup: TestMarkup10}
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(markups, nil)

				baserateClient := ts.setupBaseRateClient()
				baserateClient.On("GetBaseRatesById", mock.Anything, &bpb.GetBaseRatesByIdRequest{
					Uuid: baseRateUUID.String(),
				}).Return(nil, errors.New(ErrBaseRateService))
			},
			expectedError: true,
			errorContains: ErrBaseRateService,
		},
		{
			name:       "GetRateById_BaseRateNotFound",
			ownerId:    uuid.NewV4().String(),
			baseRateId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string, baseRateId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				baseRateUUID, _ := uuid.FromString(baseRateId)

				markups := &db.Markups{OwnerId: ownerUUID, Markup: TestMarkup10}
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(markups, nil)

				baserateClient := ts.setupBaseRateClient()
				baserateClient.On("GetBaseRatesById", mock.Anything, &bpb.GetBaseRatesByIdRequest{
					Uuid: baseRateUUID.String(),
				}).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: true,
			errorContains: ErrRecordNotFound,
		},
		{
			name:       "GetRateById_UserMarkupError",
			ownerId:    uuid.NewV4().String(),
			baseRateId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string, baseRateId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(nil, errors.New(ErrMarkupDatabase))
			},
			expectedError: true,
			errorContains: ErrMarkupDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestSetup()
			tt.setupMocks(ts, tt.ownerId, tt.baseRateId)

			req := &pb.GetRateByIdRequest{
				OwnerId:  tt.ownerId,
				BaseRate: tt.baseRateId,
			}
			rateRes, err := ts.rateService.GetRateById(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, rateRes)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, rateRes)
				assert.NotNil(t, rateRes.Rate)
				// Verify markup calculation
				expectedData := MarkupRate(TestDataRate, TestMarkup10)
				expectedSmsMo := MarkupRate(TestSmsMoRate, TestMarkup10)
				expectedSmsMt := MarkupRate(TestSmsMtRate, TestMarkup10)

				assert.InDelta(t, expectedData, rateRes.Rate.Data, 1e-8)
				assert.InDelta(t, expectedSmsMo, rateRes.Rate.SmsMo, 1e-8)
				assert.InDelta(t, expectedSmsMt, rateRes.Rate.SmsMt, 1e-8)
			}

			ts.assertAllExpectations(t)
		})
	}
}
