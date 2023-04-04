package server

import (
	"context"
	"testing"
	"time"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	mocks "github.com/ukama/ukama/systems/data-plan/rate/mocks"
	pb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	"github.com/ukama/ukama/systems/data-plan/rate/pkg/db"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
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
	})

}
