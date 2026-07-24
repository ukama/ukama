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
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/report/generator/internal/server"
	"github.com/ukama/ukama/systems/report/generator/mocks"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

const (
	OrgName = "testorg"
)

var raw = `{
	"lago_id": "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
	"sequential_id": 2,
	"number": "LAG-1234-001-002",
	"issuing_date": "2022-04-30",
	"status": "failed",
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

var rawReceipt = `{
	"number": "RCPT-58C1E087",
	"issuing_date": "2026-07-24T10:00:00Z",
	"invoice_type": "receipt",
	"status": "finalized",
	"payment_status": "succeeded",
	"fees_amount_cents": 1050,
	"sub_total_excluding_taxes_amount_cents": 1050,
	"sub_total_including_taxes_amount_cents": 1050,
	"total_amount_cents": 1050,
	"currency": "USD",
	"file_url": "http://{API_ENDPOINT}/pdf/58c1e087-9575-4bb5-bf8a-44fdbad01adf.pdf",
	"customer": {
	"external_id": "adb6a5b9-2f92-4c72-884b-d408b283d6e4",
	"name": "Gavin Belson",
	"email": "gavin@hooli.test",
	"phone": "1-171-883-3711",
	"country": "US",
	"currency": "USD"
	},
	"fees": [
	{
	"external_subscription_id": "txn-123",
	"amount_cents": 1050,
	"total_amount_cents": 1050,
	"total_amount_currency": "USD",
	"description": "Payment for package",
	"item": {
	"type": "cash",
	"code": "ff8e5b35-8947-4e9f-8168-c798bdf3e326",
	"name": "package"
	}
	}
	]
	}`

func TestGeneratorEventServer_HandleReceiptGenerateEvent(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	receiptEvent := func(t *testing.T) *epb.Report {
		t.Helper()

		val := &epb.RawReport{}

		m := protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}

		err := m.Unmarshal([]byte(rawReceipt), val)
		assert.NoError(t, err)

		return &epb.Report{
			Id:            uuid.NewV4().String(),
			OwnerId:       uuid.NewV4().String(),
			OwnerType:     ukama.OwnerTypeOrg.String(),
			Type:          ukama.ReportTypeReceipt.String(),
			RawReport:     val,
			IsPaid:        true,
			TransactionId: uuid.NewV4().String(),
		}
	}

	eventFor := func(t *testing.T, routingKey string) *epb.Event {
		t.Helper()

		anyE, err := anypb.New(receiptEvent(t))
		assert.NoError(t, err)

		return &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}
	}

	t.Run("ReceiptGeneratedEventSent", func(t *testing.T) {
		pdfEngine := &mocks.PdfEngine{}

		pdfEngine.On("Configure", mock.Anything, mock.Anything).
			Return(nil).Once()

		pdfEngine.On("Generate", mock.Anything).
			Return(nil).Once()

		routingKey := msgbus.PrepareRoute(OrgName,
			"event.cloud.local.{{ .Org}}.billing.report.receipt.generate")

		s := server.NewGeneratorEventServer(OrgName, pdfEngine, msgbusClient)

		_, err := s.EventNotification(context.TODO(), eventFor(t, routingKey))

		assert.NoError(t, err)
		pdfEngine.AssertExpectations(t)
	})

	t.Run("ReceiptUpdatedEventSent", func(t *testing.T) {
		pdfEngine := &mocks.PdfEngine{}

		pdfEngine.On("Configure", mock.Anything, mock.Anything).
			Return(nil).Once()

		pdfEngine.On("Generate", mock.Anything).
			Return(nil).Once()

		routingKey := msgbus.PrepareRoute(OrgName,
			"event.cloud.local.{{ .Org}}.billing.report.receipt.update")

		s := server.NewGeneratorEventServer(OrgName, pdfEngine, msgbusClient)

		_, err := s.EventNotification(context.TODO(), eventFor(t, routingKey))

		assert.NoError(t, err)
		pdfEngine.AssertExpectations(t)
	})

	t.Run("ErrorOnGenerateReceiptPDF", func(t *testing.T) {
		pdfEngine := &mocks.PdfEngine{}

		pdfEngine.On("Configure", mock.Anything, mock.Anything).
			Return(nil).Once()

		pdfEngine.On("Generate", mock.Anything).
			Return(errors.New("fail to generate file")).Once()

		routingKey := msgbus.PrepareRoute(OrgName,
			"event.cloud.local.{{ .Org}}.billing.report.receipt.generate")

		s := server.NewGeneratorEventServer(OrgName, pdfEngine, msgbusClient)

		_, err := s.EventNotification(context.TODO(), eventFor(t, routingKey))

		assert.Error(t, err)
		pdfEngine.AssertExpectations(t)
	})

	t.Run("InvalidReceiptEventTypeSent", func(t *testing.T) {
		pdfEngine := &mocks.PdfEngine{}

		anyE, err := anypb.New(&epb.Notification{Id: uuid.NewV4().String()})
		assert.NoError(t, err)

		routingKey := msgbus.PrepareRoute(OrgName,
			"event.cloud.local.{{ .Org}}.billing.report.receipt.generate")

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewGeneratorEventServer(OrgName, pdfEngine, msgbusClient)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestGeneratorEventServer_HandleInvoiceGenerateEvent(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	routingKey := msgbus.PrepareRoute(OrgName,
		"event.cloud.local.{{ .Org}}.billing.report.invoice.generate")

	t.Run("InvoiceGeneratedEventSent", func(t *testing.T) {
		pdfEngine := &mocks.PdfEngine{}

		pdfEngine.On("Configure", mock.Anything, mock.Anything).
			Return(nil).Once()

		pdfEngine.On("Generate", mock.Anything).
			Return(nil).Once()

		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).
			Return(nil).Once()

		val := &epb.RawReport{}

		m := protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}

		err := m.Unmarshal([]byte(raw), val)
		assert.NoError(t, err)

		evt := &epb.Report{
			Id:        uuid.NewV4().String(),
			OwnerId:   uuid.NewV4().String(),
			OwnerType: ukama.OwnerTypeOrg.String(),
			RawReport: val,
			IsPaid:    false,
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewGeneratorEventServer(OrgName, pdfEngine, msgbusClient)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("ErrorOnConfigurePDF", func(t *testing.T) {
		pdfEngine := &mocks.PdfEngine{}
		val := &epb.RawReport{}

		pdfEngine.On("Configure", mock.Anything, mock.Anything).
			Return(errors.New("fail to generate file")).Once()

		m := protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}

		err := m.Unmarshal([]byte(raw), val)
		assert.NoError(t, err)

		evt := &epb.Report{
			Id:        uuid.NewV4().String(),
			OwnerId:   uuid.NewV4().String(),
			OwnerType: ukama.OwnerTypeSubscriber.String(),
			NetworkId: uuid.NewV4().String(),
			RawReport: val,
			IsPaid:    false,
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewGeneratorEventServer(OrgName, pdfEngine, msgbusClient)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("ErrorOnGeneratePDF", func(t *testing.T) {
		pdfEngine := &mocks.PdfEngine{}
		val := &epb.RawReport{}

		pdfEngine.On("Configure", mock.Anything, mock.Anything).
			Return(nil).Once()

		pdfEngine.On("Generate", mock.Anything).
			Return(errors.New("fail to generate file")).Once()

		m := protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}

		err := m.Unmarshal([]byte(raw), val)
		assert.NoError(t, err)

		evt := &epb.Report{
			Id:        uuid.NewV4().String(),
			OwnerId:   uuid.NewV4().String(),
			OwnerType: ukama.OwnerTypeOrg.String(),
			RawReport: val,
			IsPaid:    false,
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewGeneratorEventServer(OrgName, pdfEngine, msgbusClient)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("InvalidEventTypeSent", func(t *testing.T) {
		pdfEngine := &mocks.PdfEngine{}
		evt := &epb.Notification{
			Id: uuid.NewV4().String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewGeneratorEventServer(OrgName, pdfEngine, msgbusClient)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}
