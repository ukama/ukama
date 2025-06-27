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

// testSetup contains all the mocks and service instance for testing
type testSetup struct {
	markupRepo     *mocks.MarkupsRepo
	defMarkupRepo  *mocks.DefaultMarkupRepo
	baserateSvc    *mocks.BaserateClientProvider
	msgbusClient   *mbmocks.MsgBusServiceClient
	rateService    *RateServer
	baserateClient *splmocks.BaseRatesServiceClient
}

// newTestSetup creates a new test setup with all mocks initialized
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

// setupBaseRateClient sets up the base rate client mock
func (ts *testSetup) setupBaseRateClient() *splmocks.BaseRatesServiceClient {
	baserateClient := &splmocks.BaseRatesServiceClient{}
	ts.baserateSvc.On("GetClient").Return(baserateClient, nil)
	ts.baserateClient = baserateClient
	return baserateClient
}

// assertAllExpectations asserts expectations on all mocks
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
					Markup:  10,
				}
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(markups, nil)
			},
			expectedError:  false,
			expectedMarkup: 10,
		},
		{
			name:    "MarkupforOwnerIdDoesn'tExists",
			ownerId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				defMarkup := &db.DefaultMarkup{Markup: 5}
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(nil, gorm.ErrRecordNotFound)
				ts.defMarkupRepo.On("GetDefaultMarkupRate").Return(defMarkup, nil)
			},
			expectedError:  false,
			expectedMarkup: 5,
		},
		{
			name:    "InvalidUUID",
			ownerId: "invalid-uuid",
			setupMocks: func(ts *testSetup, ownerId string) {
				// No mocks needed for invalid UUID
			},
			expectedError: true,
			errorContains: "Error parsing UUID",
		},
		{
			name:    "DatabaseErrorWhenGettingMarkup",
			ownerId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				dbError := errors.New("database connection failed")
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(nil, dbError)
			},
			expectedError: true,
			errorContains: "database connection failed",
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
			markup: 5,
			setupMocks: func(ts *testSetup) {
				ts.defMarkupRepo.On("UpdateDefaultMarkupRate", 5.0).Return(nil)
				ts.msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "UpdateDefaultMarkup_DatabaseError",
			markup: 5,
			setupMocks: func(ts *testSetup) {
				ts.defMarkupRepo.On("UpdateDefaultMarkupRate", 5.0).Return(errors.New("database error"))
			},
			expectedError: true,
			errorContains: "database error",
		},
		{
			name:   "UpdateDefaultMarkup_MessageBusError",
			markup: 5,
			setupMocks: func(ts *testSetup) {
				ts.defMarkupRepo.On("UpdateDefaultMarkupRate", 5.0).Return(nil)
				ts.msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("message bus error"))
			},
			expectedError: false, // Message bus errors are logged but don't fail the operation
		},
		{
			name:   "UpdateDefaultMarkup_RecordNotFound",
			markup: 5,
			setupMocks: func(ts *testSetup) {
				ts.defMarkupRepo.On("UpdateDefaultMarkupRate", 5.0).Return(gorm.ErrRecordNotFound)
			},
			expectedError: true,
			errorContains: "record not found",
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
				defMarkup := &db.DefaultMarkup{Markup: 5}
				ts.defMarkupRepo.On("GetDefaultMarkupRate").Return(defMarkup, nil)
			},
			expectedError:  false,
			expectedMarkup: 5,
		},
		{
			name: "GetDefaultMarkup_DatabaseError",
			setupMocks: func(ts *testSetup) {
				ts.defMarkupRepo.On("GetDefaultMarkupRate").Return(nil, errors.New("database error"))
			},
			expectedError: true,
			errorContains: "database error",
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
	cTime, _ := time.Parse(time.RFC3339, "2021-11-12T11:45:26.371Z")
	uTime, _ := time.Parse(time.RFC3339, "2022-10-12T11:45:26.371Z")
	dTime, _ := time.Parse(time.RFC3339, "2022-11-12T11:45:26.371Z")
	deleteAt := gorm.DeletedAt{Time: dTime, Valid: true}

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
							CreatedAt: cTime,
							DeletedAt: deleteAt,
							UpdatedAt: uTime,
						},
						Markup: 5,
					},
					{
						Model: gorm.Model{
							ID:        2,
							CreatedAt: cTime,
							UpdatedAt: uTime,
						},
						Markup: 5,
					},
				}
				ts.defMarkupRepo.On("GetDefaultMarkupRateHistory").Return(defMarkup, nil)
			},
			expectedError: false,
		},
		{
			name: "GetDefaultMarkupHistory_DatabaseError",
			setupMocks: func(ts *testSetup) {
				ts.defMarkupRepo.On("GetDefaultMarkupRateHistory").Return(nil, errors.New("database error"))
			},
			expectedError: true,
			errorContains: "database error",
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
					assert.Equal(t, float64(5), rate.Markup)
					assert.Equal(t, cTime.Format(time.RFC3339), rate.CreatedAt)
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
				Country:  "USA",
				Provider: "Ukama",
				SimType:  "ukama_data",
				From:     "2033-04-20T20:31:24-00:00",
				To:       "2043-04-20T20:31:24-00:00",
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				markups := &db.Markups{OwnerId: ownerId, Markup: 10}
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
							Apn:         "Manual entry required",
							Country:     req.Country,
							Data:        0.0014,
							EffectiveAt: "2033-04-20T20:31:24+00:00",
							Imsi:        1,
							Lte:         true,
							Provider:    "Multi Tel",
							SimType:     req.SimType,
							SmsMo:       0.0100,
							SmsMt:       0.0001,
							Vpmn:        "TTC",
						},
					},
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "GetRate_InvalidUUID",
			req: &pb.GetRateRequest{
				OwnerId:  "invalid-uuid",
				Country:  "USA",
				Provider: "Ukama",
				SimType:  "ukama_data",
				From:     "2033-04-20T20:31:24-00:00",
				To:       "2043-04-20T20:31:24-00:00",
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				// No mocks needed for invalid UUID
			},
			expectedError: true,
			errorContains: "Error parsing UUID",
		},
		{
			name: "GetRate_InvalidDateFormatTo",
			req: &pb.GetRateRequest{
				OwnerId:  uuid.NewV4().String(),
				Country:  "USA",
				Provider: "Ukama",
				SimType:  "ukama_data",
				From:     "2033-04-20T20:31:24-00:00",
				To:       "invalid-date-format",
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				markups := &db.Markups{OwnerId: ownerId, Markup: 10}
				ts.markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)
			},
			expectedError: true,
			errorContains: "invalid date format for to",
		},
		{
			name: "GetRate_InvalidDateFormatFrom",
			req: &pb.GetRateRequest{
				OwnerId:  uuid.NewV4().String(),
				Country:  "USA",
				Provider: "Ukama",
				SimType:  "ukama_data",
				From:     "invalid-date-format",
				To:       "2043-04-20T20:31:24-00:00",
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				markups := &db.Markups{OwnerId: ownerId, Markup: 10}
				ts.markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)
			},
			expectedError: true,
			errorContains: "invalid date format for from",
		},
		{
			name: "GetRate_BaseRateClientError",
			req: &pb.GetRateRequest{
				OwnerId:  uuid.NewV4().String(),
				Country:  "USA",
				Provider: "Ukama",
				SimType:  "ukama_data",
				From:     "2033-04-20T20:31:24-00:00",
				To:       "2043-04-20T20:31:24-00:00",
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				markups := &db.Markups{OwnerId: ownerId, Markup: 10}
				ts.markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)
				ts.baserateSvc.On("GetClient").Return(nil, errors.New("base rate client error"))
			},
			expectedError: true,
			errorContains: "base rate client error",
		},
		{
			name: "GetRate_BaseRatesError",
			req: &pb.GetRateRequest{
				OwnerId:  uuid.NewV4().String(),
				Country:  "USA",
				Provider: "Ukama",
				SimType:  "ukama_data",
				From:     "2033-04-20T20:31:24-00:00",
				To:       "2043-04-20T20:31:24-00:00",
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				markups := &db.Markups{OwnerId: ownerId, Markup: 10}
				ts.markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)

				baserateClient := ts.setupBaseRateClient()
				baserateClient.On("GetBaseRatesForPeriod", mock.Anything, mock.Anything).Return(nil, errors.New("base rates error"))
			},
			expectedError: true,
			errorContains: "base rates error",
		},
		{
			name: "GetRate_NoValidBaseRates",
			req: &pb.GetRateRequest{
				OwnerId:  uuid.NewV4().String(),
				Country:  "USA",
				Provider: "Ukama",
				SimType:  "ukama_data",
				From:     "2033-04-20T20:31:24-00:00",
				To:       "2043-04-20T20:31:24-00:00",
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				markups := &db.Markups{OwnerId: ownerId, Markup: 10}
				ts.markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)

				baserateClient := ts.setupBaseRateClient()
				baserateClient.On("GetBaseRatesForPeriod", mock.Anything, mock.Anything).Return(&bpb.GetBaseRatesResponse{
					Rates: []*bpb.Rate{},
				}, nil)
			},
			expectedError: true,
			errorContains: "no valid base rates found",
		},
		{
			name: "GetRate_NilBaseRates",
			req: &pb.GetRateRequest{
				OwnerId:  uuid.NewV4().String(),
				Country:  "USA",
				Provider: "Ukama",
				SimType:  "ukama_data",
				From:     "2033-04-20T20:31:24-00:00",
				To:       "2043-04-20T20:31:24-00:00",
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				markups := &db.Markups{OwnerId: ownerId, Markup: 10}
				ts.markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)

				baserateClient := ts.setupBaseRateClient()
				baserateClient.On("GetBaseRatesForPeriod", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedError: true,
			errorContains: "no valid base rates found",
		},
		{
			name: "GetRate_UserMarkupError",
			req: &pb.GetRateRequest{
				OwnerId:  uuid.NewV4().String(),
				Country:  "USA",
				Provider: "Ukama",
				SimType:  "ukama_data",
				From:     "2033-04-20T20:31:24-00:00",
				To:       "2043-04-20T20:31:24-00:00",
			},
			setupMocks: func(ts *testSetup, req *pb.GetRateRequest) {
				ownerId, _ := uuid.FromString(req.OwnerId)
				ts.markupRepo.On("GetMarkupRate", ownerId).Return(nil, errors.New("markup database error"))
			},
			expectedError: true,
			errorContains: "markup database error",
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
			markup:  10,
			setupMocks: func(ts *testSetup, ownerId string, markup float64) {
				ownerUUID, _ := uuid.FromString(ownerId)
				ts.markupRepo.On("UpdateMarkupRate", ownerUUID, markup).Return(nil)
				ts.msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "UpdateMarkup_InvalidUUID",
			ownerId: "invalid-uuid",
			markup:  10,
			setupMocks: func(ts *testSetup, ownerId string, markup float64) {
				// No mocks needed for invalid UUID
			},
			expectedError: true,
			errorContains: "Error parsing UUID",
		},
		{
			name:    "UpdateMarkup_DatabaseError",
			ownerId: uuid.NewV4().String(),
			markup:  10,
			setupMocks: func(ts *testSetup, ownerId string, markup float64) {
				ownerUUID, _ := uuid.FromString(ownerId)
				ts.markupRepo.On("UpdateMarkupRate", ownerUUID, markup).Return(errors.New("database error"))
			},
			expectedError: true,
			errorContains: "database error",
		},
		{
			name:    "UpdateMarkup_MessageBusError",
			ownerId: uuid.NewV4().String(),
			markup:  10,
			setupMocks: func(ts *testSetup, ownerId string, markup float64) {
				ownerUUID, _ := uuid.FromString(ownerId)
				ts.markupRepo.On("UpdateMarkupRate", ownerUUID, markup).Return(nil)
				ts.msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("message bus error"))
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
			ownerId: "invalid-uuid",
			setupMocks: func(ts *testSetup, ownerId string) {
				// No mocks needed for invalid UUID
			},
			expectedError: true,
			errorContains: "Error parsing UUID",
		},
		{
			name:    "DeleteMarkup_DatabaseError",
			ownerId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				ts.markupRepo.On("DeleteMarkupRate", ownerUUID).Return(errors.New("database error"))
			},
			expectedError: true,
			errorContains: "database error",
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
			Markup:  10,
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
	cTime, _ := time.Parse(time.RFC3339, "2021-11-12T11:45:26.371Z")
	uTime, _ := time.Parse(time.RFC3339, "2022-10-12T11:45:26.371Z")
	dTime, _ := time.Parse(time.RFC3339, "2022-11-12T11:45:26.371Z")
	deleteAt := gorm.DeletedAt{Time: dTime, Valid: true}

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
							CreatedAt: cTime,
							DeletedAt: deleteAt,
							UpdatedAt: uTime,
						},
						OwnerId: ownerUUID,
						Markup:  5,
					},
					{
						Model: gorm.Model{
							ID:        2,
							CreatedAt: cTime,
							UpdatedAt: uTime,
						},
						OwnerId: ownerUUID,
						Markup:  10,
					},
				}
				ts.markupRepo.On("GetMarkupRateHistory", ownerUUID).Return(markup, nil)
			},
			expectedError: false,
		},
		{
			name:    "GetMarkupHistory_InvalidUUID",
			ownerId: "invalid-uuid",
			setupMocks: func(ts *testSetup, ownerId string) {
				// No mocks needed for invalid UUID
			},
			expectedError: true,
			errorContains: "Error parsing UUID",
		},
		{
			name:    "GetMarkupHistory_DatabaseError",
			ownerId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				ts.markupRepo.On("GetMarkupRateHistory", ownerUUID).Return(nil, errors.New("database error"))
			},
			expectedError: true,
			errorContains: "database error",
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
					assert.Equal(t, float64(5+i*5), rate.Markup) // 5, 10
					assert.Equal(t, cTime.Format(time.RFC3339), rate.CreatedAt)
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

				markups := &db.Markups{OwnerId: ownerUUID, Markup: 10}
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(markups, nil)

				baserateClient := ts.setupBaseRateClient()
				origData := 0.0014
				origSmsMo := 0.0100
				origSmsMt := 0.0001

				baseRate := &bpb.Rate{
					X2G:         true,
					X3G:         true,
					Apn:         "Manual entry required",
					Country:     "USA",
					Data:        origData,
					EffectiveAt: "2033-04-20T20:31:24+00:00",
					Imsi:        1,
					Lte:         true,
					Provider:    "Multi Tel",
					SimType:     "ukama_data",
					SmsMo:       origSmsMo,
					SmsMt:       origSmsMt,
					Vpmn:        "TTC",
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
			ownerId:    "invalid-uuid",
			baseRateId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string, baseRateId string) {
				// No mocks needed for invalid UUID
			},
			expectedError: true,
			errorContains: "Error parsing UUID",
		},
		{
			name:       "GetRateById_BaseRateClientError",
			ownerId:    uuid.NewV4().String(),
			baseRateId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string, baseRateId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				markups := &db.Markups{OwnerId: ownerUUID, Markup: 10}
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(markups, nil)
				ts.baserateSvc.On("GetClient").Return(nil, errors.New("base rate client error"))
			},
			expectedError: true,
			errorContains: "base rate client error",
		},
		{
			name:       "GetRateById_BaseRateServiceError",
			ownerId:    uuid.NewV4().String(),
			baseRateId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string, baseRateId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				baseRateUUID, _ := uuid.FromString(baseRateId)

				markups := &db.Markups{OwnerId: ownerUUID, Markup: 10}
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(markups, nil)

				baserateClient := ts.setupBaseRateClient()
				baserateClient.On("GetBaseRatesById", mock.Anything, &bpb.GetBaseRatesByIdRequest{
					Uuid: baseRateUUID.String(),
				}).Return(nil, errors.New("base rate service error"))
			},
			expectedError: true,
			errorContains: "base rate service error",
		},
		{
			name:       "GetRateById_BaseRateNotFound",
			ownerId:    uuid.NewV4().String(),
			baseRateId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string, baseRateId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				baseRateUUID, _ := uuid.FromString(baseRateId)

				markups := &db.Markups{OwnerId: ownerUUID, Markup: 10}
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(markups, nil)

				baserateClient := ts.setupBaseRateClient()
				baserateClient.On("GetBaseRatesById", mock.Anything, &bpb.GetBaseRatesByIdRequest{
					Uuid: baseRateUUID.String(),
				}).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: true,
			errorContains: "record not found",
		},
		{
			name:       "GetRateById_UserMarkupError",
			ownerId:    uuid.NewV4().String(),
			baseRateId: uuid.NewV4().String(),
			setupMocks: func(ts *testSetup, ownerId string, baseRateId string) {
				ownerUUID, _ := uuid.FromString(ownerId)
				ts.markupRepo.On("GetMarkupRate", ownerUUID).Return(nil, errors.New("markup database error"))
			},
			expectedError: true,
			errorContains: "markup database error",
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
				origData := 0.0014
				origSmsMo := 0.0100
				origSmsMt := 0.0001
				expectedData := MarkupRate(origData, 10)
				expectedSmsMo := MarkupRate(origSmsMo, 10)
				expectedSmsMt := MarkupRate(origSmsMt, 10)

				assert.InDelta(t, expectedData, rateRes.Rate.Data, 1e-8)
				assert.InDelta(t, expectedSmsMo, rateRes.Rate.SmsMo, 1e-8)
				assert.InDelta(t, expectedSmsMt, rateRes.Rate.SmsMt, 1e-8)
			}

			ts.assertAllExpectations(t)
		})
	}
}
