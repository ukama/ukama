/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client_test

import (
	"testing"

	"github.com/ukama/ukama/systems/billing/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/billing/invoice/pb/gen/mocks"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pmocks "github.com/ukama/ukama/systems/billing/api-gateway/mocks"
	pb "github.com/ukama/ukama/systems/billing/invoice/pb/gen"
)

func TestBillingClient_AddInvoice(t *testing.T) {
	t.Run("SubscriberIdValid", func(t *testing.T) {
		var bc = &mocks.InvoiceServiceClient{}

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

		invoiceReq := &pb.AddRequest{
			RawInvoice: raw,
		}

		invoiceResp := &pb.AddResponse{Invoice: &pb.Invoice{
			Id:           uuid.NewV4().String(),
			SubscriberId: "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
			NetworkId:    uuid.NewV4().String(),
		}}

		bc.On("Add", mock.Anything, invoiceReq).Return(invoiceResp, nil)

		b := client.NewBillingFromClient(bc, nil)

		resp, err := b.AddInvoice(raw)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		bc.AssertExpectations(t)
	})

	t.Run("SubscriberIdNotValid", func(t *testing.T) {
		var bc = &mocks.InvoiceServiceClient{}

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
	"external_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6bX",
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

		invoiceReq := &pb.AddRequest{
			RawInvoice: raw,
		}

		bc.On("Add", mock.Anything, invoiceReq).Return(nil,
			status.Errorf(codes.InvalidArgument, "invalid subscriberId"))

		b := client.NewBillingFromClient(bc, nil)

		resp, err := b.AddInvoice(raw)

		assert.Error(t, err)
		assert.Nil(t, resp)
		bc.AssertExpectations(t)
	})
}

func TestBillingClient_GetInvoice(t *testing.T) {
	var bc = &mocks.InvoiceServiceClient{}

	t.Run("InvoiceFound", func(t *testing.T) {
		invoiceId := uuid.NewV4()

		invoiceReq := &pb.GetRequest{
			InvoiceId: invoiceId.String(),
			AsPdf:     false,
		}

		invoiceResp := &pb.GetResponse{Invoice: &pb.Invoice{
			Id:           invoiceId.String(),
			SubscriberId: uuid.NewV4().String(),
		}}

		bc.On("Get", mock.Anything, invoiceReq).Return(invoiceResp, nil)

		b := client.NewBillingFromClient(bc, nil)

		resp, err := b.GetInvoice(invoiceId.String(), false)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, resp.Invoice.Id, invoiceId.String())
		bc.AssertExpectations(t)
	})

	t.Run("InvoiceNotFound", func(t *testing.T) {
		invoiceId := uuid.NewV4()

		invoiceReq := &pb.GetRequest{
			InvoiceId: invoiceId.String(),
			AsPdf:     false,
		}

		bc.On("Get", mock.Anything, invoiceReq).Return(nil,
			status.Errorf(codes.NotFound, "invoice not found"))

		b := client.NewBillingFromClient(bc, nil)

		resp, err := b.GetInvoice(invoiceId.String(), false)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "not found")
		bc.AssertExpectations(t)
	})
}

func TestBillingClient_GetInvoicesBySubscriber(t *testing.T) {
	var bc = &mocks.InvoiceServiceClient{}

	t.Run("SubscriberFound", func(t *testing.T) {
		subscriberId := uuid.NewV4()
		invoiceId := uuid.NewV4()

		req := &pb.GetBySubscriberRequest{
			SubscriberId: subscriberId.String(),
		}

		invoiceResp := &pb.GetBySubscriberResponse{Invoices: []*pb.Invoice{
			&pb.Invoice{
				Id:           invoiceId.String(),
				SubscriberId: subscriberId.String(),
			}}}

		bc.On("GetBySubscriber", mock.Anything, req).Return(invoiceResp, nil)

		b := client.NewBillingFromClient(bc, nil)

		resp, err := b.GetInvoicesBySubscriber(subscriberId.String())

		assert.NoError(t, err)
		assert.Equal(t, resp.Invoices[0].Id, invoiceId.String())
		bc.AssertExpectations(t)
	})

	t.Run("SubscriberNotFound", func(t *testing.T) {
		subscriberId := uuid.NewV4()

		req := &pb.GetBySubscriberRequest{
			SubscriberId: subscriberId.String(),
		}

		bc.On("GetBySubscriber", mock.Anything, req).Return(nil,
			status.Errorf(codes.NotFound, "subscriber not found"))

		b := client.NewBillingFromClient(bc, nil)

		resp, err := b.GetInvoicesBySubscriber(subscriberId.String())

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "not found")
		bc.AssertExpectations(t)
	})
}

func TestBillingClient_GetInvoicesByNetwork(t *testing.T) {
	var bc = &mocks.InvoiceServiceClient{}

	t.Run("NetworkFound", func(t *testing.T) {
		networkId := uuid.NewV4()
		invoiceId := uuid.NewV4()

		req := &pb.GetByNetworkRequest{
			NetworkId: networkId.String(),
		}

		invoiceResp := &pb.GetByNetworkResponse{Invoices: []*pb.Invoice{
			&pb.Invoice{
				Id:        invoiceId.String(),
				NetworkId: networkId.String(),
			}}}

		bc.On("GetByNetwork", mock.Anything, req).Return(invoiceResp, nil)

		b := client.NewBillingFromClient(bc, nil)

		resp, err := b.GetInvoicesByNetwork(networkId.String())

		assert.NoError(t, err)
		assert.Equal(t, resp.Invoices[0].Id, invoiceId.String())
		bc.AssertExpectations(t)
	})

	t.Run("InvoiceNotFound", func(t *testing.T) {
		networkId := uuid.NewV4()

		req := &pb.GetByNetworkRequest{
			NetworkId: networkId.String(),
		}

		bc.On("GetByNetwork", mock.Anything, req).Return(nil,
			status.Errorf(codes.NotFound, "network not found"))

		b := client.NewBillingFromClient(bc, nil)

		resp, err := b.GetInvoicesByNetwork(networkId.String())

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "not found")
		bc.AssertExpectations(t)
	})
}

func TestBillingClient_RemoveInvoice(t *testing.T) {
	var bc = &mocks.InvoiceServiceClient{}

	t.Run("InvoiceFound", func(t *testing.T) {
		invoiceId := uuid.NewV4()

		invoiceReq := &pb.DeleteRequest{
			InvoiceId: invoiceId.String(),
		}

		bc.On("Delete", mock.Anything, invoiceReq).Return(nil, nil)

		b := client.NewBillingFromClient(bc, nil)

		err := b.RemoveInvoice(invoiceId.String())

		assert.NoError(t, err)
		bc.AssertExpectations(t)
	})

	t.Run("InvoiceNotFound", func(t *testing.T) {
		invoiceId := uuid.NewV4()

		invoiceReq := &pb.DeleteRequest{
			InvoiceId: invoiceId.String(),
		}

		bc.On("Delete", mock.Anything, invoiceReq).Return(nil,
			status.Errorf(codes.NotFound, "invoice not found"))

		b := client.NewBillingFromClient(bc, nil)

		err := b.RemoveInvoice(invoiceId.String())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		bc.AssertExpectations(t)
	})
}

func TestBillingClient_GetInvoicePDF(t *testing.T) {
	var pc = &pmocks.PdfClient{}
	var bc = &mocks.InvoiceServiceClient{}

	t.Run("InvoiceFound", func(t *testing.T) {
		invoiceId := uuid.NewV4()

		pc.On("GetPdf", invoiceId.String()).Return([]byte("some fake pdf data"), nil)

		b := client.NewBillingFromClient(bc, pc)

		resp, err := b.GetInvoicePDF(invoiceId.String())

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		pc.AssertExpectations(t)
	})

	t.Run("InvoiceNotFound", func(t *testing.T) {
		invoiceId := uuid.NewV4()

		pc.On("GetPdf", invoiceId.String()).Return(nil,
			status.Errorf(codes.NotFound, "invoice not found"))

		b := client.NewBillingFromClient(bc, pc)

		resp, err := b.GetInvoicePDF(invoiceId.String())

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "not found")
		pc.AssertExpectations(t)
	})
}
