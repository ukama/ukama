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

func TestRateService_GetMarkup(t *testing.T) {

	t.Run("MarkupforOwnerIdExists", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

		markups := &db.Markups{
			OwnerId: uuid.NewV4(),
			Markup:  10,
		}

		req := &pb.GetMarkupRequest{
			OwnerId: markups.OwnerId.String(),
		}

		markupRepo.On("GetMarkupRate", markups.OwnerId).Return(markups, nil)
		rateRes, err := rateService.GetMarkup(context.Background(), req)
		assert.NoError(t, err)

		assert.Equal(t, rateRes.Markup, markups.Markup)
		assert.Equal(t, rateRes.OwnerId, markups.OwnerId.String())

		markupRepo.AssertExpectations(t)
	})

	t.Run("MarkupforOwnerIdDoesn'tExists", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

		markups := &db.Markups{
			OwnerId: uuid.NewV4(),
			Markup:  10,
		}

		defMarkup := &db.DefaultMarkup{
			Markup: 5,
		}

		req := &pb.GetMarkupRequest{
			OwnerId: markups.OwnerId.String(),
		}

		markupRepo.On("GetMarkupRate", markups.OwnerId).Return(nil, gorm.ErrRecordNotFound)
		defMarkupRepo.On("GetDefaultMarkupRate").Return(defMarkup, nil)

		rateRes, err := rateService.GetMarkup(context.Background(), req)
		assert.NoError(t, err)

		assert.Equal(t, rateRes.Markup, defMarkup.Markup)
		assert.Equal(t, rateRes.OwnerId, markups.OwnerId.String())

		markupRepo.AssertExpectations(t)
		defMarkupRepo.AssertExpectations(t)

	})

}
func TestRateService_UpdateDefaultMarkup(t *testing.T) {

	t.Run("UpdateDefaultMarkupSuccess", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

		defMarkup := &db.DefaultMarkup{
			Markup: 5,
		}

		req := &pb.UpdateDefaultMarkupRequest{
			Markup: defMarkup.Markup,
		}

		defMarkupRepo.On("UpdateDefaultMarkupRate", defMarkup.Markup).Return(nil)
		msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)
		_, err := rateService.UpdateDefaultMarkup(context.Background(), req)
		assert.NoError(t, err)

		defMarkupRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("UpdateDefaultMarkup_DatabaseError", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

		req := &pb.UpdateDefaultMarkupRequest{
			Markup: 5,
		}

		defMarkupRepo.On("UpdateDefaultMarkupRate", req.Markup).Return(errors.New("database error"))

		_, err := rateService.UpdateDefaultMarkup(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		defMarkupRepo.AssertExpectations(t)
	})

	t.Run("UpdateDefaultMarkup_MessageBusError", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

		req := &pb.UpdateDefaultMarkupRequest{
			Markup: 5,
		}

		defMarkupRepo.On("UpdateDefaultMarkupRate", req.Markup).Return(nil)
		msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("message bus error"))

		_, err := rateService.UpdateDefaultMarkup(context.Background(), req)
		assert.NoError(t, err) // Message bus errors are logged but don't fail the operation

		defMarkupRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("UpdateDefaultMarkup_RecordNotFound", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

		req := &pb.UpdateDefaultMarkupRequest{
			Markup: 5,
		}

		defMarkupRepo.On("UpdateDefaultMarkupRate", req.Markup).Return(gorm.ErrRecordNotFound)

		_, err := rateService.UpdateDefaultMarkup(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "record not found")

		defMarkupRepo.AssertExpectations(t)
	})
}

func TestRateService_GetDefaultMarkup(t *testing.T) {

	t.Run("GetDefaultMarkupSuccess", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

		defMarkup := &db.DefaultMarkup{
			Markup: 5,
		}

		req := &pb.GetDefaultMarkupRequest{}

		defMarkupRepo.On("GetDefaultMarkupRate").Return(defMarkup, nil)
		rateRes, err := rateService.GetDefaultMarkup(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, rateRes.Markup, defMarkup.Markup)
		defMarkupRepo.AssertExpectations(t)
	})

}

func TestRateService_GetDefaultMarkupHistory(t *testing.T) {

	t.Run("GetDefaultMarkupHistorySuccess", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)
		cTime, err := time.Parse(time.RFC3339, "2021-11-12T11:45:26.371Z")
		assert.NoError(t, err)
		uTime, err := time.Parse(time.RFC3339, "2022-10-12T11:45:26.371Z")
		assert.NoError(t, err)
		dTime, err := time.Parse(time.RFC3339, "2022-11-12T11:45:26.371Z")
		DeleteAt := gorm.DeletedAt{
			Time:  dTime,
			Valid: true,
		}
		assert.NoError(t, err)

		defMarkup := []db.DefaultMarkup{
			{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: cTime,
					DeletedAt: DeleteAt,
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

		req := &pb.GetDefaultMarkupHistoryRequest{}

		defMarkupRepo.On("GetDefaultMarkupRateHistory").Return(defMarkup, nil)
		rateRes, err := rateService.GetDefaultMarkupHistory(context.Background(), req)
		assert.NoError(t, err)
		if assert.NotNil(t, rateRes) {
			for i, rate := range rateRes.MarkupRates {
				assert.Equal(t, defMarkup[i].Markup, rate.Markup)
				assert.Equal(t, defMarkup[i].CreatedAt.Format(time.RFC3339), rate.CreatedAt)
			}
		}
		defMarkupRepo.AssertExpectations(t)
	})

}

func TestRateService_GetRate(t *testing.T) {
	t.Run("GetRate_Success", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)
		ownerId := uuid.NewV4()

		req := &pb.GetRateRequest{
			OwnerId:  ownerId.String(),
			Country:  "USA",
			Provider: "Ukama",
			SimType:  "ukama_data",
			From:     "2033-04-20T20:31:24-00:00",
			To:       "2043-04-20T20:31:24-00:00",
		}

		to, err := validation.FromString(req.To)
		assert.NoError(t, err)
		from, err := validation.FromString(req.From)
		assert.NoError(t, err)
		markups := &db.Markups{
			OwnerId: ownerId,
			Markup:  10,
		}

		rates := &bpb.GetBaseRatesResponse{
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
		}

		markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)
		baserateClient := baserateSvc.On("GetClient").
			Return(&splmocks.BaseRatesServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*splmocks.BaseRatesServiceClient)

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
					SimType:     rates.Rates[0].SimType,
					SmsMo:       0.0100,
					SmsMt:       0.0001,
					Vpmn:        "TTC",
				},
			},
		}, nil).Once()

		rateRes, err := rateService.GetRate(context.Background(), req)
		assert.NoError(t, err)
		if assert.NotNil(t, rateRes) {
			for _, r := range rateRes.Rates {
				assert.Equal(t, req.Country, r.Country)
				assert.Equal(t, req.SimType, r.SimType)
			}
		}

		markupRepo.AssertExpectations(t)
		baserateClient.AssertExpectations(t)
		defMarkupRepo.AssertExpectations(t)
	})
}

func TestRateService_UpdateMarkup(t *testing.T) {

	t.Run("UpdateMarkupSuccess", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

		markup := &db.Markups{
			OwnerId: uuid.NewV4(),
			Markup:  10,
		}

		req := &pb.UpdateMarkupRequest{
			OwnerId: markup.OwnerId.String(),
			Markup:  markup.Markup,
		}

		markupRepo.On("UpdateMarkupRate", markup.OwnerId, markup.Markup).Return(nil)
		msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)
		_, err := rateService.UpdateMarkup(context.Background(), req)
		assert.NoError(t, err)

		markupRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("UpdateMarkup_InvalidUUID", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

		req := &pb.UpdateMarkupRequest{
			OwnerId: "invalid-uuid",
			Markup:  10,
		}

		_, err := rateService.UpdateMarkup(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("UpdateMarkup_DatabaseError", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

		ownerId := uuid.NewV4()
		req := &pb.UpdateMarkupRequest{
			OwnerId: ownerId.String(),
			Markup:  10,
		}

		markupRepo.On("UpdateMarkupRate", ownerId, req.Markup).Return(errors.New("database error"))

		_, err := rateService.UpdateMarkup(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		markupRepo.AssertExpectations(t)
	})

	t.Run("UpdateMarkup_MessageBusError", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

		markup := &db.Markups{
			OwnerId: uuid.NewV4(),
			Markup:  10,
		}

		req := &pb.UpdateMarkupRequest{
			OwnerId: markup.OwnerId.String(),
			Markup:  markup.Markup,
		}

		markupRepo.On("UpdateMarkupRate", markup.OwnerId, markup.Markup).Return(nil)
		msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("message bus error"))

		_, err := rateService.UpdateMarkup(context.Background(), req)
		assert.NoError(t, err) // Message bus errors are logged but don't fail the operation

		markupRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})
}

func TestRateService_DeleteMarkup(t *testing.T) {

	t.Run("DeleteMarkupSuccess", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

		markup := &db.Markups{
			OwnerId: uuid.NewV4(),
			Markup:  10,
		}

		req := &pb.DeleteMarkupRequest{
			OwnerId: markup.OwnerId.String(),
		}

		markupRepo.On("DeleteMarkupRate", markup.OwnerId).Return(nil)
		msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)
		_, err := rateService.DeleteMarkup(context.Background(), req)
		assert.NoError(t, err)

		markupRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("DeleteMarkup_InvalidUUID", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

		req := &pb.DeleteMarkupRequest{
			OwnerId: "invalid-uuid",
		}

		_, err := rateService.DeleteMarkup(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("DeleteMarkup_DatabaseError", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

		ownerId := uuid.NewV4()
		req := &pb.DeleteMarkupRequest{
			OwnerId: ownerId.String(),
		}

		markupRepo.On("DeleteMarkupRate", ownerId).Return(errors.New("database error"))

		_, err := rateService.DeleteMarkup(context.Background(), req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		markupRepo.AssertExpectations(t)
	})
}

func TestRateService_GetMarkupVal(t *testing.T) {

	t.Run("GetMarkupSuccess", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)

		markup := &db.Markups{
			OwnerId: uuid.NewV4(),
			Markup:  10,
		}

		req := &pb.GetMarkupRequest{
			OwnerId: markup.OwnerId.String(),
		}

		markupRepo.On("GetMarkupRate", markup.OwnerId).Return(markup, nil)
		rateRes, err := rateService.GetMarkup(context.Background(), req)
		assert.NoError(t, err)
		if assert.NotNil(t, rateRes) {
			assert.Equal(t, markup.Markup, rateRes.Markup)
			assert.Equal(t, markup.OwnerId.String(), rateRes.OwnerId)
		}
		defMarkupRepo.AssertExpectations(t)
	})

}

func TestRateService_GetMarkupHistory(t *testing.T) {

	t.Run("GetMarkupHistorySuccess", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		OwnerId := uuid.NewV4()

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)
		cTime, err := time.Parse(time.RFC3339, "2021-11-12T11:45:26.371Z")
		assert.NoError(t, err)
		uTime, err := time.Parse(time.RFC3339, "2022-10-12T11:45:26.371Z")
		assert.NoError(t, err)
		dTime, err := time.Parse(time.RFC3339, "2022-11-12T11:45:26.371Z")
		DeleteAt := gorm.DeletedAt{
			Time:  dTime,
			Valid: true,
		}
		assert.NoError(t, err)

		markup := []db.Markups{
			{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: cTime,
					DeletedAt: DeleteAt,
					UpdatedAt: uTime,
				},
				OwnerId: OwnerId,
				Markup:  5,
			},
			{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: cTime,
					UpdatedAt: uTime,
				},
				OwnerId: OwnerId,
				Markup:  10,
			},
		}

		req := &pb.GetMarkupHistoryRequest{
			OwnerId: OwnerId.String(),
		}

		markupRepo.On("GetMarkupRateHistory", OwnerId).Return(markup, nil)
		rateRes, err := rateService.GetMarkupHistory(context.Background(), req)
		assert.NoError(t, err)
		if assert.NotNil(t, rateRes) {
			for i, rate := range rateRes.MarkupRates {
				assert.Equal(t, markup[i].Markup, rate.Markup)
				assert.Equal(t, markup[i].CreatedAt.Format(time.RFC3339), rate.CreatedAt)
			}
		}
		defMarkupRepo.AssertExpectations(t)
	})

}

func TestRateService_GetRateById(t *testing.T) {
	t.Run("GetRateById_Success", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)
		ownerId := uuid.NewV4()
		baseRateId := uuid.NewV4()

		req := &pb.GetRateByIdRequest{
			OwnerId:  ownerId.String(),
			BaseRate: baseRateId.String(),
		}

		markups := &db.Markups{
			OwnerId: ownerId,
			Markup:  10,
		}

		// Use original values for expected calculations
		origData := 0.0014
		origSmsMo := 0.0100
		origSmsMt := 0.0001

		expectedData := MarkupRate(origData, markups.Markup)
		expectedSmsMo := MarkupRate(origSmsMo, markups.Markup)
		expectedSmsMt := MarkupRate(origSmsMt, markups.Markup)

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

		markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)
		baserateClient := baserateSvc.On("GetClient").
			Return(&splmocks.BaseRatesServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*splmocks.BaseRatesServiceClient)

		baserateClient.On("GetBaseRatesById", mock.Anything, &bpb.GetBaseRatesByIdRequest{
			Uuid: baseRateId.String(),
		}).Return(&bpb.GetBaseRatesByIdResponse{
			Rate: baseRate,
		}, nil).Once()

		rateRes, err := rateService.GetRateById(context.Background(), req)
		assert.NoError(t, err)
		if assert.NotNil(t, rateRes) && assert.NotNil(t, rateRes.Rate) {
			// Use InDelta for floating point comparison
			assert.InDelta(t, expectedData, rateRes.Rate.Data, 1e-8)
			assert.InDelta(t, expectedSmsMo, rateRes.Rate.SmsMo, 1e-8)
			assert.InDelta(t, expectedSmsMt, rateRes.Rate.SmsMt, 1e-8)
			assert.Equal(t, baseRate.Country, rateRes.Rate.Country)
			assert.Equal(t, baseRate.SimType, rateRes.Rate.SimType)
		}

		markupRepo.AssertExpectations(t)
		baserateClient.AssertExpectations(t)
	})

	t.Run("GetRateById_InvalidOwnerId", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)
		baseRateId := uuid.NewV4()

		req := &pb.GetRateByIdRequest{
			OwnerId:  "invalid-uuid",
			BaseRate: baseRateId.String(),
		}

		rateRes, err := rateService.GetRateById(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, rateRes)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("GetRateById_BaseRateClientError", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)
		ownerId := uuid.NewV4()
		baseRateId := uuid.NewV4()

		req := &pb.GetRateByIdRequest{
			OwnerId:  ownerId.String(),
			BaseRate: baseRateId.String(),
		}

		markups := &db.Markups{
			OwnerId: ownerId,
			Markup:  10,
		}

		markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)
		baserateSvc.On("GetClient").Return(nil, errors.New("base rate client error"))

		rateRes, err := rateService.GetRateById(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, rateRes)
		assert.Contains(t, err.Error(), "base rate client error")

		markupRepo.AssertExpectations(t)
		baserateSvc.AssertExpectations(t)
	})

	t.Run("GetRateById_BaseRateServiceError", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)
		ownerId := uuid.NewV4()
		baseRateId := uuid.NewV4()

		req := &pb.GetRateByIdRequest{
			OwnerId:  ownerId.String(),
			BaseRate: baseRateId.String(),
		}

		markups := &db.Markups{
			OwnerId: ownerId,
			Markup:  10,
		}

		markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)
		baserateClient := baserateSvc.On("GetClient").
			Return(&splmocks.BaseRatesServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*splmocks.BaseRatesServiceClient)

		baserateClient.On("GetBaseRatesById", mock.Anything, &bpb.GetBaseRatesByIdRequest{
			Uuid: baseRateId.String(),
		}).Return(nil, errors.New("base rate service error")).Once()

		rateRes, err := rateService.GetRateById(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, rateRes)
		assert.Contains(t, err.Error(), "base rate service error")

		markupRepo.AssertExpectations(t)
		baserateClient.AssertExpectations(t)
	})

	t.Run("GetRateById_BaseRateNotFound", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baserateSvc := &mocks.BaserateClientProvider{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(OrgName, markupRepo, defMarkupRepo, baserateSvc, msgbusClient)
		ownerId := uuid.NewV4()
		baseRateId := uuid.NewV4()

		req := &pb.GetRateByIdRequest{
			OwnerId:  ownerId.String(),
			BaseRate: baseRateId.String(),
		}

		markups := &db.Markups{
			OwnerId: ownerId,
			Markup:  10,
		}

		markupRepo.On("GetMarkupRate", ownerId).Return(markups, nil)
		baserateClient := baserateSvc.On("GetClient").
			Return(&splmocks.BaseRatesServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*splmocks.BaseRatesServiceClient)

		baserateClient.On("GetBaseRatesById", mock.Anything, &bpb.GetBaseRatesByIdRequest{
			Uuid: baseRateId.String(),
		}).Return(nil, gorm.ErrRecordNotFound).Once()

		rateRes, err := rateService.GetRateById(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, rateRes)
		assert.Contains(t, err.Error(), "record not found")

		markupRepo.AssertExpectations(t)
		baserateClient.AssertExpectations(t)
	})
}
