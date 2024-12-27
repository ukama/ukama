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
	"gorm.io/datatypes"

	"github.com/ukama/ukama/systems/billing/report/mocks"
	"github.com/ukama/ukama/systems/billing/report/pkg/db"
	"github.com/ukama/ukama/systems/billing/report/pkg/server"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/protobuf/types/known/anypb"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

const (
	OrgName = "testOrg"
	OrgId   = "592f7a8e-f318-4d3a-aab8-8d4187cde7f9"
)

func TestReportEventServer_HandlePaymentSuccessEvent(t *testing.T) {
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.payments.processor.payment.success")

	t.Run("ReportIdNotValid", func(t *testing.T) {
		reportRepo := &mocks.ReportRepo{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		payment := epb.Payment{
			ItemId: "lol",
		}

		anyE, err := anypb.New(&payment)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s, err := server.NewReportEventServer(OrgName, OrgId, reportRepo, msgbusClient)
		assert.NoError(t, err)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("ReportNotFound", func(t *testing.T) {
		reportRepo := &mocks.ReportRepo{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		reportId := uuid.NewV4()

		reportRepo.On("Get", reportId, mock.Anything).
			Return(nil, errors.New("not found")).Once()

		payment := epb.Payment{
			ItemId: reportId.String(),
		}

		anyE, err := anypb.New(&payment)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s, err := server.NewReportEventServer(OrgName, OrgId, reportRepo, msgbusClient)
		assert.NoError(t, err)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("ErrroOnUpdate", func(t *testing.T) {
		reportRepo := &mocks.ReportRepo{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		reportId := uuid.NewV4()

		report := &db.Report{
			Id:     reportId,
			IsPaid: false,
		}

		reportRepo.On("Get", reportId, mock.Anything).
			Return(report, nil).Once()

		reportRepo.On("Update", report, mock.Anything).
			Return(errors.New("Error on update")).Once()

		payment := epb.Payment{
			ItemId: reportId.String(),
		}

		anyE, err := anypb.New(&payment)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s, err := server.NewReportEventServer(OrgName, OrgId, reportRepo, msgbusClient)
		assert.NoError(t, err)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("PaymentSuccessEventSent", func(t *testing.T) {

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

		reportRepo := &mocks.ReportRepo{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		reportId := uuid.NewV4()

		report := &db.Report{
			Id:        reportId,
			IsPaid:    false,
			RawReport: datatypes.JSON([]byte(raw)),
		}

		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).
			Return(nil).Once()

		reportRepo.On("Get", reportId, mock.Anything).
			Return(report, nil).Once()

		reportRepo.On("Update", report, mock.Anything).
			Return(nil).Once()

		payment := epb.Payment{
			ItemId: reportId.String(),
		}

		anyE, err := anypb.New(&payment)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s, err := server.NewReportEventServer(OrgName, OrgId, reportRepo, msgbusClient)
		assert.NoError(t, err)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("PaymentSuccessPayloadNotSent", func(t *testing.T) {
		reportRepo := &mocks.ReportRepo{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		payment := epb.Notification{}

		anyE, err := anypb.New(&payment)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s, err := server.NewReportEventServer(OrgName, OrgId, reportRepo, msgbusClient)
		assert.NoError(t, err)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("PaymentSuccessEventNotSent", func(t *testing.T) {
		reportRepo := &mocks.ReportRepo{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		payment := epb.Notification{}

		anyE, err := anypb.New(&payment)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s, err := server.NewReportEventServer(OrgName, OrgId, reportRepo, msgbusClient)
		assert.NoError(t, err)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}
