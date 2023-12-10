/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/ukama/ukama/systems/billing/invoice/mocks"
	"github.com/ukama/ukama/systems/billing/invoice/pkg/db"
	"github.com/ukama/ukama/systems/billing/invoice/pkg/server"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	pb "github.com/ukama/ukama/systems/billing/invoice/pb/gen"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	csub "github.com/ukama/ukama/systems/common/rest/client/subscriber"
)

const OrgName = "testorg"

func TestInvoiceServer_Add(t *testing.T) {
	t.Run("SubscriberIsValid", func(t *testing.T) {
		// Arrange
		var raw = `{
	"lago_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
	"sequential_id": 2,
	"number": "LAG-1234-001-002",
	"issuing_date": "2022-04-30",
	"status": "finalized",
	"payment_status": "succeeded",
	"amount_cents": 100,
	"amount_currency": "EUR",
	"vat_amount_cents": 20,
	"vat_amount_currency": "EUR",
	"credit_amount_cents": 10,
	"credit_amount_currency": "EUR",
	"total_amount_cents": 110,
	"total_amount_currency": "EUR",
	"file_url": "https://getlago.com/invoice/file",
	"legacy": false,
	"customer": {
	"lago_id": "99a6094e-199b-4101-896a-54e927ce7bd7",
	"external_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
	"address_line1": "5230 Penfield Ave",
	"address_line2": null,
	"city": "Woodland Hills",
	"country": "US",
	"created_at": "2022-04-29T08:59:51Z",
	"email": "dinesh@piedpiper.test",
	"legal_name": "Coleman-Blair",
	"legal_number": "49-008-2965",
	"logo_url": "http://hooli.com/logo.png",
	"name": "Gavin Belson",
	"phone": "1-171-883-3711 x245",
	"state": "CA",
	"url": "http://hooli.com",
	"vat_rate": 20.0,
	"zipcode": "91364"
	},
	"subscriptions": [
	{
	"lago_id": "b7ab2926-1de8-4428-9bcd-779314ac129b",
	"external_id": "susbcription_external_id",
	"lago_customer_id": "99a6094e-199b-4101-896a-54e927ce7bd7",
	"external_customer_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
	"canceled_at": "2022-04-29T08:59:51Z",
	"created_at": "2022-04-29T08:59:51Z",
	"plan_code": "new_code",
	"started_at": "2022-04-29T08:59:51Z",
	"status": "active",
	"terminated_at": null
	}
	],
	"fees": [
	{
	"lago_id": "6be23c42-47d2-45a3-9770-5b3572f225c3",
	"lago_group_id": null,
	"item": {
	"type": "subscription",
	"code": "plan_code",
	"name": "Plan"
	},
	"amount_cents": 100,
	"amount_currency": "EUR",
	"vat_amount_cents": 20,
	"vat_amount_currency": "EUR",
	"total_amount_cents": 120,
	"total_amount_currency": "EUR",
	"units": "0.32",
	"events_count": 23
	}
	],
	"credits": [
	{
	"lago_id": "b7ab2926-1de8-4428-9bcd-779314ac129b",
	"item": {
	"lago_id": "b7ab2926-1de8-4428-9bcd-779314ac129b",
	"type": "coupon",
	"code": "coupon_code",
	"name": "Coupon"
	},
	"amount_cents": 100,
	"amount_currency": "EUR"
	}
	],
	"metadata": [
	{
	"lago_id": "27f12d13-4ae0-437b-b822-8771bcd62e3a",
	"key": "digital_ref_id",
	"value": "INV-0123456-98765",
	"created_at": "2022-04-29T08:59:51Z"
	}
	]
	}`

		var sId = "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba"

		subscriberId, err := uuid.FromString(sId)
		if err != nil {
			t.Fatalf("invalid subscriberId input: %s", sId)
		}

		invoiceRepo := &mocks.InvoiceRepo{}
		subscriberClient := &cmocks.SubscriberClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		invoiceRepo.On("Add", mock.Anything, mock.Anything).Return(nil).Once()

		subscriberClient.On("Get", subscriberId.String()).Return(&csub.SubscriberInfo{
			SubscriberId: subscriberId,
			NetworkId:    uuid.NewV4(),
		}, nil).Once()

		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := server.NewInvoiceServer(OrgName, invoiceRepo, subscriberClient, msgbusClient)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			RawInvoice: raw,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, subscriberId.String(), res.Invoice.SubscriberId)
		invoiceRepo.AssertExpectations(t)
	})

	t.Run("SubscriberIsNotValid", func(t *testing.T) {
		// Arrange
		invoiceRepo := &mocks.InvoiceRepo{}
		subscriberClient := &cmocks.SubscriberClient{}

		var raw = `{
	"lago_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
	"sequential_id": 2,
	"number": "LAG-1234-001-002",
	"issuing_date": "2022-04-30",
	"status": "finalized",
	"payment_status": "succeeded",
	"amount_cents": 100,
	"amount_currency": "EUR",
	"vat_amount_cents": 20,
	"vat_amount_currency": "EUR",
	"credit_amount_cents": 10,
	"credit_amount_currency": "EUR",
	"total_amount_cents": 110,
	"total_amount_currency": "EUR",
	"file_url": "https://getlago.com/invoice/file",
	"legacy": false,
	"customer": {
	"lago_id": "99a6094e-199b-4101-896a-54e927ce7bd7",
	"external_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6baX",
	"address_line1": "5230 Penfield Ave",
	"address_line2": null,
	"city": "Woodland Hills",
	"country": "US",
	"created_at": "2022-04-29T08:59:51Z",
	"email": "dinesh@piedpiper.test",
	"legal_name": "Coleman-Blair",
	"legal_number": "49-008-2965",
	"logo_url": "http://hooli.com/logo.png",
	"name": "Gavin Belson",
	"phone": "1-171-883-3711 x245",
	"state": "CA",
	"url": "http://hooli.com",
	"vat_rate": 20.0,
	"zipcode": "91364"
	},
	"subscriptions": [
	{
	"lago_id": "b7ab2926-1de8-4428-9bcd-779314ac129b",
	"external_id": "susbcription_external_id",
	"lago_customer_id": "99a6094e-199b-4101-896a-54e927ce7bd7",
	"external_customer_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
	"canceled_at": "2022-04-29T08:59:51Z",
	"created_at": "2022-04-29T08:59:51Z",
	"plan_code": "new_code",
	"started_at": "2022-04-29T08:59:51Z",
	"status": "active",
	"terminated_at": null
	}
	],
	"fees": [
	{
	"lago_id": "6be23c42-47d2-45a3-9770-5b3572f225c3",
	"lago_group_id": null,
	"item": {
	"type": "subscription",
	"code": "plan_code",
	"name": "Plan"
	},
	"amount_cents": 100,
	"amount_currency": "EUR",
	"vat_amount_cents": 20,
	"vat_amount_currency": "EUR",
	"total_amount_cents": 120,
	"total_amount_currency": "EUR",
	"units": "0.32",
	"events_count": 23
	}
	],
	"credits": [
	{
	"lago_id": "b7ab2926-1de8-4428-9bcd-779314ac129b",
	"item": {
	"lago_id": "b7ab2926-1de8-4428-9bcd-779314ac129b",
	"type": "coupon",
	"code": "coupon_code",
	"name": "Coupon"
	},
	"amount_cents": 100,
	"amount_currency": "EUR"
	}
	],
	"metadata": [
	{
	"lago_id": "27f12d13-4ae0-437b-b822-8771bcd62e3a",
	"key": "digital_ref_id",
	"value": "INV-0123456-98765",
	"created_at": "2022-04-29T08:59:51Z"
	}
	]
	}`

		s := server.NewInvoiceServer(OrgName, invoiceRepo, subscriberClient, nil)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			RawInvoice: raw,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		invoiceRepo.AssertExpectations(t)
	})

	t.Run("RawInvoiceIsNotValid", func(t *testing.T) {
		// Arrange
		var raw = "+{}"

		invoiceRepo := &mocks.InvoiceRepo{}
		subscriberClient := &cmocks.SubscriberClient{}

		s := server.NewInvoiceServer(OrgName, invoiceRepo, subscriberClient, nil)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			RawInvoice: raw,
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

		var raw = `{
	"lago_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
	"sequential_id": 2,
	"number": "LAG-1234-001-002",
	"issuing_date": "2022-04-30",
	"status": "finalized",
	"payment_status": "succeeded",
	"amount_cents": 100,
	"amount_currency": "EUR",
	"vat_amount_cents": 20,
	"vat_amount_currency": "EUR",
	"credit_amount_cents": 10,
	"credit_amount_currency": "EUR",
	"total_amount_cents": 110,
	"total_amount_currency": "EUR",
	"file_url": "https://getlago.com/invoice/file",
	"legacy": false,
	"customer": {
	"lago_id": "99a6094e-199b-4101-896a-54e927ce7bd7",
	"external_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
	"address_line1": "5230 Penfield Ave",
	"address_line2": null,
	"city": "Woodland Hills",
	"country": "US",
	"created_at": "2022-04-29T08:59:51Z",
	"email": "dinesh@piedpiper.test",
	"legal_name": "Coleman-Blair",
	"legal_number": "49-008-2965",
	"logo_url": "http://hooli.com/logo.png",
	"name": "Gavin Belson",
	"phone": "1-171-883-3711 x245",
	"state": "CA",
	"url": "http://hooli.com",
	"vat_rate": 20.0,
	"zipcode": "91364"
	},
	"subscriptions": [
	{
	"lago_id": "b7ab2926-1de8-4428-9bcd-779314ac129b",
	"external_id": "susbcription_external_id",
	"lago_customer_id": "99a6094e-199b-4101-896a-54e927ce7bd7",
	"external_customer_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
	"canceled_at": "2022-04-29T08:59:51Z",
	"created_at": "2022-04-29T08:59:51Z",
	"plan_code": "new_code",
	"started_at": "2022-04-29T08:59:51Z",
	"status": "active",
	"terminated_at": null
	}
	],
	"fees": [
	{
	"lago_id": "6be23c42-47d2-45a3-9770-5b3572f225c3",
	"lago_group_id": null,
	"item": {
	"type": "subscription",
	"code": "plan_code",
	"name": "Plan"
	},
	"amount_cents": 100,
	"amount_currency": "EUR",
	"vat_amount_cents": 20,
	"vat_amount_currency": "EUR",
	"total_amount_cents": 120,
	"total_amount_currency": "EUR",
	"units": "0.32",
	"events_count": 23
	}
	],
	"credits": [
	{
	"lago_id": "b7ab2926-1de8-4428-9bcd-779314ac129b",
	"item": {
	"lago_id": "b7ab2926-1de8-4428-9bcd-779314ac129b",
	"type": "coupon",
	"code": "coupon_code",
	"name": "Coupon"
	},
	"amount_cents": 100,
	"amount_currency": "EUR"
	}
	],
	"metadata": [
	{
	"lago_id": "27f12d13-4ae0-437b-b822-8771bcd62e3a",
	"key": "digital_ref_id",
	"value": "INV-0123456-98765",
	"created_at": "2022-04-29T08:59:51Z"
	}
	]
	}`

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

		s := server.NewInvoiceServer(OrgName, invoiceRepo, nil, nil)
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

		s := server.NewInvoiceServer(OrgName, invoiceRepo, nil, nil)
		resp, err := s.Get(context.TODO(), &pb.GetRequest{
			InvoiceId: invoiceId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		invoiceRepo.AssertExpectations(t)
	})

	t.Run("InvoiceUUIDInvalid", func(t *testing.T) {
		var invoiceId = "1"

		invoiceRepo := &mocks.InvoiceRepo{}

		s := server.NewInvoiceServer(OrgName, invoiceRepo, nil, nil)
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

		s := server.NewInvoiceServer(OrgName, invoiceRepo, nil, nil)

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

		s := server.NewInvoiceServer(OrgName, invoiceRepo, nil, nil)

		res, err := s.GetBySubscriber(context.TODO(), &pb.GetBySubscriberRequest{
			SubscriberId: subscriberId.String()})

		assert.Error(t, err)
		assert.Nil(t, res)
		invoiceRepo.AssertExpectations(t)
	})

	t.Run("SubscriberUUIDInvalid", func(t *testing.T) {
		var subscriberID = "1"

		invoiceRepo := &mocks.InvoiceRepo{}

		s := server.NewInvoiceServer(OrgName, invoiceRepo, nil, nil)

		res, err := s.GetBySubscriber(context.TODO(), &pb.GetBySubscriberRequest{
			SubscriberId: subscriberID})

		assert.Error(t, err)
		assert.Nil(t, res)
		invoiceRepo.AssertExpectations(t)
	})
}

func TestInvoiceServer_GetInvoiceByNetwork(t *testing.T) {
	t.Run("NetworkFound", func(t *testing.T) {
		var invoiceId = uuid.NewV4()
		var networkId = uuid.NewV4()

		invoiceRepo := &mocks.InvoiceRepo{}

		invoiceRepo.On("GetByNetwork", networkId).Return(
			[]db.Invoice{
				{Id: invoiceId,
					NetworkId: networkId,
					IsPaid:    false,
				}}, nil).Once()

		s := server.NewInvoiceServer(OrgName, invoiceRepo, nil, nil)

		res, err := s.GetByNetwork(context.TODO(),
			&pb.GetByNetworkRequest{NetworkId: networkId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, invoiceId.String(), res.GetInvoices()[0].GetId())
		assert.Equal(t, networkId.String(), res.NetworkId)
		invoiceRepo.AssertExpectations(t)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		var networkId = uuid.Nil

		invoiceRepo := &mocks.InvoiceRepo{}

		invoiceRepo.On("GetByNetwork", networkId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewInvoiceServer(OrgName, invoiceRepo, nil, nil)

		res, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{
			NetworkId: networkId.String()})

		assert.Error(t, err)
		assert.Nil(t, res)
		invoiceRepo.AssertExpectations(t)
	})

	t.Run("NetworkUUIDInvalid", func(t *testing.T) {
		var networkId = "1"

		invoiceRepo := &mocks.InvoiceRepo{}

		s := server.NewInvoiceServer(OrgName, invoiceRepo, nil, nil)

		res, err := s.GetByNetwork(context.TODO(), &pb.GetByNetworkRequest{
			NetworkId: networkId})

		assert.Error(t, err)
		assert.Nil(t, res)
		invoiceRepo.AssertExpectations(t)
	})
}

func TestInvoiceServer_Delete(t *testing.T) {
	t.Run("InvoiceFound", func(t *testing.T) {
		var invoiceId = uuid.NewV4()

		invoiceRepo := &mocks.InvoiceRepo{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		invoiceRepo.On("Delete", invoiceId, mock.Anything).Return(nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := server.NewInvoiceServer(OrgName, invoiceRepo, nil, msgbusClient)

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

		s := server.NewInvoiceServer(OrgName, invoiceRepo, nil, nil)

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

		s := server.NewInvoiceServer(OrgName, invoiceRepo, nil, nil)

		res, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			InvoiceId: invoiceId,
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		invoiceRepo.AssertExpectations(t)
	})
}
