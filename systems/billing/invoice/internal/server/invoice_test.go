package server_test

import (
	"context"
	"testing"
	"time"

	pb "github.com/ukama/ukama/systems/billing/invoice/pb/gen"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/billing/invoice/internal/db"
	"github.com/ukama/ukama/systems/billing/invoice/internal/server"
	"github.com/ukama/ukama/systems/billing/invoice/mocks"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func TestInvoiceServer_Add(t *testing.T) {
	t.Run("SubscriberFound", func(t *testing.T) {
		// Arrange
		var subscriberId = uuid.NewV4()
		var period = time.Now().UTC()
		var raw = "{}"

		invoiceRepo := &mocks.InvoiceRepo{}

		invoice := &db.Invoice{
			SubscriberId: subscriberId,
			Period:       period,
			RawInvoice:   datatypes.JSON([]byte(raw)),
		}

		invoiceRepo.On("Add", invoice, mock.Anything).Return(nil).Once()

		s := server.NewInvoiceServer(invoiceRepo)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			SubscriberId: subscriberId.String(),
			Period:       timestamppb.New(period),
			RawInvoice:   raw,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, subscriberId.String(), res.Invoice.SubscriberId)
		invoiceRepo.AssertExpectations(t)
	})

	t.Run("SubscriberNotFound", func(t *testing.T) {
		// Arrange
		var period = time.Now().UTC()
		var raw = "{}"

		invoiceRepo := &mocks.InvoiceRepo{}

		s := server.NewInvoiceServer(invoiceRepo)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			SubscriberId: "lol",
			Period:       timestamppb.New(period),
			RawInvoice:   raw,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		invoiceRepo.AssertExpectations(t)
	})
}

func TestInvoiceServer_Get(t *testing.T) {
	t.Run("InvoiceFound", func(t *testing.T) {
		var invoiceId = uuid.NewV4()
		var subscriberId = uuid.NewV4()
		var period = time.Now().UTC()
		var raw = "{}"

		invoiceRepo := &mocks.InvoiceRepo{}

		invoice := invoiceRepo.On("Get", invoiceId).
			Return(&db.Invoice{
				Id:           invoiceId,
				SubscriberId: subscriberId,
				Period:       period,
				RawInvoice:   datatypes.JSON([]byte(raw)),
				IsPaid:       false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*db.Invoice)

		s := server.NewInvoiceServer(invoiceRepo)
		res, err := s.Get(context.TODO(), &pb.GetRequest{
			InvoiceId: invoiceId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, invoice.Id.String(), res.GetInvoice().GetId())
		assert.Equal(t, false, res.GetInvoice().IsPaid)
		invoiceRepo.AssertExpectations(t)
	})

	t.Run("InvoiceNotFound", func(t *testing.T) {
		var invoiceId = uuid.Nil

		invoiceRepo := &mocks.InvoiceRepo{}

		invoiceRepo.On("Get", invoiceId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewInvoiceServer(invoiceRepo)
		resp, err := s.Get(context.TODO(), &pb.GetRequest{
			InvoiceId: invoiceId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		invoiceRepo.AssertExpectations(t)
	})

	t.Run("InvoiceUUIDInvalid", func(t *testing.T) {
		var invoiceId = "1"

		invoiceRepo := &mocks.InvoiceRepo{}

		s := server.NewInvoiceServer(invoiceRepo)
		res, err := s.Get(context.TODO(), &pb.GetRequest{
			InvoiceId: invoiceId})

		assert.Error(t, err)
		assert.Nil(t, res)
		invoiceRepo.AssertExpectations(t)
	})
}

func TestInvoiceServer_GetInvoiceBySubscriber(t *testing.T) {
	t.Run("SubscriberFound", func(t *testing.T) {
		var invoiceId = uuid.NewV4()
		var subscriberId = uuid.NewV4()

		invoiceRepo := &mocks.InvoiceRepo{}

		invoiceRepo.On("GetBySubscriber", subscriberId).Return(
			[]db.Invoice{
				{Id: invoiceId,
					SubscriberId: subscriberId,
					IsPaid:       false,
				}}, nil).Once()

		s := server.NewInvoiceServer(invoiceRepo)

		res, err := s.GetBySubscriber(context.TODO(),
			&pb.GetBySubscriberRequest{SubscriberId: subscriberId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, invoiceId.String(), res.GetInvoices()[0].GetId())
		assert.Equal(t, subscriberId.String(), res.SubscriberId)
		invoiceRepo.AssertExpectations(t)
	})

	t.Run("SubscriberNotFound", func(t *testing.T) {
		var subscriberId = uuid.Nil

		invoiceRepo := &mocks.InvoiceRepo{}

		invoiceRepo.On("GetBySubscriber", subscriberId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewInvoiceServer(invoiceRepo)

		res, err := s.GetBySubscriber(context.TODO(), &pb.GetBySubscriberRequest{
			SubscriberId: subscriberId.String()})

		assert.Error(t, err)
		assert.Nil(t, res)
		invoiceRepo.AssertExpectations(t)
	})

	t.Run("SubscriberUUIDInvalid", func(t *testing.T) {
		var subscriberID = "1"

		invoiceRepo := &mocks.InvoiceRepo{}

		s := server.NewInvoiceServer(invoiceRepo)

		res, err := s.GetBySubscriber(context.TODO(), &pb.GetBySubscriberRequest{
			SubscriberId: subscriberID})

		assert.Error(t, err)
		assert.Nil(t, res)
		invoiceRepo.AssertExpectations(t)
	})
}

func TestInvoiceServer_Delete(t *testing.T) {
	t.Run("InvoiceFound", func(t *testing.T) {
		var invoiceId = uuid.NewV4()

		invoiceRepo := &mocks.InvoiceRepo{}

		invoiceRepo.On("Delete", invoiceId, mock.Anything).Return(nil).Once()

		s := server.NewInvoiceServer(invoiceRepo)

		res, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			InvoiceId: invoiceId.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		invoiceRepo.AssertExpectations(t)
	})

	t.Run("InvoiceNotFound", func(t *testing.T) {
		var invoiceId = uuid.NewV4()

		invoiceRepo := &mocks.InvoiceRepo{}

		invoiceRepo.On("Delete", invoiceId, mock.Anything).Return(gorm.ErrRecordNotFound).Once()

		s := server.NewInvoiceServer(invoiceRepo)

		res, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			InvoiceId: invoiceId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		invoiceRepo.AssertExpectations(t)
	})

	t.Run("InvoiceUUIDInvalid", func(t *testing.T) {
		var invoiceId = "1"

		invoiceRepo := &mocks.InvoiceRepo{}

		s := server.NewInvoiceServer(invoiceRepo)

		res, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			InvoiceId: invoiceId,
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		invoiceRepo.AssertExpectations(t)
	})

}
