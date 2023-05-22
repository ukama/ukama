package server

import (
	"context"
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
)

func TestRateService_GetMarkup(t *testing.T) {

	t.Run("MarkupforOwnerIdExists", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baseRate := &mocks.BaseRateSrvc{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(markupRepo, defMarkupRepo, baseRate, msgbusClient)

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
		baseRate := &mocks.BaseRateSrvc{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(markupRepo, defMarkupRepo, baseRate, msgbusClient)

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
		baseRate := &mocks.BaseRateSrvc{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(markupRepo, defMarkupRepo, baseRate, msgbusClient)

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
	})

}

func TestRateService_GetDefaultMarkup(t *testing.T) {

	t.Run("GetDefaultMarkupSuccess", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baseRate := &mocks.BaseRateSrvc{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(markupRepo, defMarkupRepo, baseRate, msgbusClient)

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
		baseRate := &mocks.BaseRateSrvc{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(markupRepo, defMarkupRepo, baseRate, msgbusClient)
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
		baseRate := &mocks.BaseRateSrvc{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(markupRepo, defMarkupRepo, baseRate, msgbusClient)
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
		baseRate.On("GetBaseRates", &bpb.GetBaseRatesByPeriodRequest{
			Country:  req.Country,
			Provider: req.Provider,
			To:       to.Format(time.RFC3339),
			From:     from.Format(time.RFC3339),
			SimType:  req.SimType,
		}).Return(rates, nil)

		rateRes, err := rateService.GetRate(context.Background(), req)
		assert.NoError(t, err)
		if assert.NotNil(t, rateRes) {
			for _, r := range rateRes.Rates {
				assert.Equal(t, req.Country, r.Country)
				assert.Equal(t, req.SimType, r.SimType)
			}
		}

		markupRepo.AssertExpectations(t)
		baseRate.AssertExpectations(t)
		defMarkupRepo.AssertExpectations(t)
	})
}

func TestRateService_UpdateMarkup(t *testing.T) {

	t.Run("UpdateMarkupSuccess", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baseRate := &mocks.BaseRateSrvc{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(markupRepo, defMarkupRepo, baseRate, msgbusClient)

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

		defMarkupRepo.AssertExpectations(t)
	})

}

func TestRateService_DeleteMarkup(t *testing.T) {

	t.Run("DeleteMarkupSuccess", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baseRate := &mocks.BaseRateSrvc{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(markupRepo, defMarkupRepo, baseRate, msgbusClient)

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

		defMarkupRepo.AssertExpectations(t)
	})

}

func TestRateService_GetMarkupVal(t *testing.T) {

	t.Run("GetMarkupSuccess", func(t *testing.T) {
		markupRepo := &mocks.MarkupsRepo{}
		defMarkupRepo := &mocks.DefaultMarkupRepo{}
		baseRate := &mocks.BaseRateSrvc{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		rateService := NewRateServer(markupRepo, defMarkupRepo, baseRate, msgbusClient)

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
		baseRate := &mocks.BaseRateSrvc{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		OwnerId := uuid.NewV4()

		rateService := NewRateServer(markupRepo, defMarkupRepo, baseRate, msgbusClient)
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
