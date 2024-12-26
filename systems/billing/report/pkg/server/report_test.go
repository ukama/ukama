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
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/billing/report/mocks"
	"github.com/ukama/ukama/systems/billing/report/pkg/db"
	"github.com/ukama/ukama/systems/billing/report/pkg/server"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	pb "github.com/ukama/ukama/systems/billing/report/pb/gen"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	csub "github.com/ukama/ukama/systems/common/rest/client/subscriber"
)

const (
	ownerTypeOrg        = "org"
	reportTypeInvoice   = "invoice"
	ownerTypeSubscriber = "subscriber"
)

func TestReportServer_Add(t *testing.T) {
	t.Run("OwnerIsValid", func(t *testing.T) {
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

		var ownerIdString = "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba"

		ownerId, err := uuid.FromString(ownerIdString)
		if err != nil {
			t.Fatalf("invalid OwnerId input: %s", ownerIdString)
		}

		reportRepo := &mocks.ReportRepo{}
		subscriberClient := &cmocks.SubscriberClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		reportRepo.On("Add", mock.Anything, mock.Anything).Return(nil).Once()

		subscriberClient.On("Get", ownerId.String()).Return(&csub.SubscriberInfo{
			SubscriberId: ownerId,
			NetworkId:    uuid.NewV4(),
		}, nil).Once()

		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := server.NewReportServer(OrgName, OrgId, reportRepo, subscriberClient, msgbusClient)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			RawReport: raw,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, ownerId.String(), res.Report.OwnerId)
		reportRepo.AssertExpectations(t)
	})

	t.Run("OwnerIdIsNotValid", func(t *testing.T) {
		// Arrange
		reportRepo := &mocks.ReportRepo{}
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

		s := server.NewReportServer(OrgName, OrgId, reportRepo, subscriberClient, nil)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			RawReport: raw,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		reportRepo.AssertExpectations(t)
	})

	t.Run("OwnerIsNotValid", func(t *testing.T) {
		// Arrange
		reportRepo := &mocks.ReportRepo{}
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

		s := server.NewReportServer(OrgName, OrgId, reportRepo, subscriberClient, nil)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			RawReport: raw,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		reportRepo.AssertExpectations(t)
	})

	t.Run("RawReportIsNotValid", func(t *testing.T) {
		// Arrange
		var raw = "+{}"

		reportRepo := &mocks.ReportRepo{}
		subscriberClient := &cmocks.SubscriberClient{}

		s := server.NewReportServer(OrgName, OrgId, reportRepo, subscriberClient, nil)

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			RawReport: raw,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		reportRepo.AssertExpectations(t)
	})
}

func TestReportServer_Get(t *testing.T) {
	t.Run("ReportFound", func(t *testing.T) {
		var reportId = uuid.NewV4()
		var OwnerId = uuid.NewV4()
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

		reportRepo := &mocks.ReportRepo{}

		report := reportRepo.On("Get", reportId).
			Return(&db.Report{
				Id:        reportId,
				OwnerId:   OwnerId,
				Period:    period,
				RawReport: datatypes.JSON([]byte(raw)),
				IsPaid:    false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*db.Report)

		s := server.NewReportServer(OrgName, OrgId, reportRepo, nil, nil)
		res, err := s.Get(context.TODO(), &pb.GetRequest{
			ReportId: reportId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, report.Id.String(), res.GetReport().GetId())
		assert.Equal(t, false, res.GetReport().IsPaid)
		reportRepo.AssertExpectations(t)
	})

	t.Run("ReportNotFound", func(t *testing.T) {
		var reportId = uuid.Nil

		reportRepo := &mocks.ReportRepo{}

		reportRepo.On("Get", reportId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewReportServer(OrgName, OrgId, reportRepo, nil, nil)
		resp, err := s.Get(context.TODO(), &pb.GetRequest{
			ReportId: reportId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		reportRepo.AssertExpectations(t)
	})

	t.Run("ReportUUIDInvalid", func(t *testing.T) {
		var reportId = "1"

		reportRepo := &mocks.ReportRepo{}

		s := server.NewReportServer(OrgName, OrgId, reportRepo, nil, nil)
		res, err := s.Get(context.TODO(), &pb.GetRequest{
			ReportId: reportId})

		assert.Error(t, err)
		assert.Nil(t, res)
		reportRepo.AssertExpectations(t)
	})
}

func TestReportServer_List(t *testing.T) {
	resp := make([]db.Report, 1)
	var reportId = uuid.NewV4()
	var OwnerId = uuid.NewV4()
	var networkId = uuid.NewV4()
	var isPaid = true

	t.Run("ListAll", func(t *testing.T) {
		repo := &mocks.ReportRepo{}
		reportResp := db.Report{
			Id:        reportId,
			OwnerId:   OwnerId,
			OwnerType: ukama.OwnerTypeOrg,
		}

		resp[0] = reportResp

		repo.On("List", "", ukama.OwnerTypeUnknown, "", ukama.ReportTypeUnknown, false,
			uint32(0), false).Return(resp, nil)

		s := server.NewReportServer(OrgName, OrgId, repo, nil, nil)
		list, err := s.List(context.TODO(), &pb.ListRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assertList(t, list, resp)
	})

	t.Run("ListSpecificOwner", func(t *testing.T) {
		repo := &mocks.ReportRepo{}
		reportResp := db.Report{
			Id:        reportId,
			OwnerId:   OwnerId,
			OwnerType: ukama.OwnerTypeSubscriber,
		}

		resp[0] = reportResp

		repo.On("List", OwnerId.String(), ukama.OwnerTypeUnknown, "",
			ukama.ReportTypeUnknown, false, uint32(0), false).Return(resp, nil)

		s := server.NewReportServer(OrgName, OrgId, repo, nil, nil)
		list, err := s.List(context.TODO(), &pb.ListRequest{
			OwnerId: OwnerId.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assertList(t, list, resp)
	})

	t.Run("SpecificOwnerNotFound", func(t *testing.T) {
		repo := &mocks.ReportRepo{}
		notFoundId := uuid.NewV4()

		repo.On("List", notFoundId.String(), ukama.OwnerTypeUnknown, "",
			ukama.ReportTypeUnknown, false, uint32(0), false).
			Return(nil, errors.New("not found"))

		s := server.NewReportServer(OrgName, OrgId, repo, nil, nil)
		list, err := s.List(context.TODO(), &pb.ListRequest{
			OwnerId: notFoundId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("ListSpecificNetwork", func(t *testing.T) {
		repo := &mocks.ReportRepo{}
		reportResp := db.Report{
			Id:        reportId,
			OwnerId:   OwnerId,
			OwnerType: ukama.OwnerTypeSubscriber,
			NetworkId: networkId,
		}

		resp[0] = reportResp

		repo.On("List", "", ukama.OwnerTypeUnknown, networkId.String(),
			ukama.ReportTypeUnknown, false, uint32(0), false).Return(resp, nil)

		s := server.NewReportServer(OrgName, OrgId, repo, nil, nil)
		list, err := s.List(context.TODO(), &pb.ListRequest{
			NetworkId: networkId.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assertList(t, list, resp)
	})

	t.Run("InvalidReportNetwork", func(t *testing.T) {
		repo := &mocks.ReportRepo{}
		repo.On("List", "", "", "lol", "", false, uint32(0), false).Return(nil, errors.New("invalid"))

		s := server.NewReportServer(OrgName, OrgId, repo, nil, nil)
		list, err := s.List(context.TODO(), &pb.ListRequest{
			NetworkId: "lol",
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("ListSortedUnpaidOrgInvoiceReports", func(t *testing.T) {
		repo := &mocks.ReportRepo{}
		reportResp := db.Report{
			Id:        reportId,
			OwnerId:   OwnerId,
			OwnerType: ukama.OwnerTypeOrg,
			Type:      ukama.ReportTypeInvoice,
		}

		resp[0] = reportResp

		repo.On("List", "", ukama.OwnerTypeOrg, "", ukama.ReportTypeInvoice,
			false, uint32(0), true).
			Return(resp, nil)

		s := server.NewReportServer(OrgName, OrgId, repo, nil, nil)
		list, err := s.List(context.TODO(), &pb.ListRequest{
			OwnerType:  ownerTypeOrg,
			ReportType: reportTypeInvoice,
			IsPaid:     false,
			Sort:       true,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assertList(t, list, resp)
	})

	t.Run("ListPaidSubscriberReports", func(t *testing.T) {
		repo := &mocks.ReportRepo{}
		reportResp := db.Report{
			Id:        reportId,
			OwnerId:   OwnerId,
			OwnerType: ukama.OwnerTypeSubscriber,
			IsPaid:    true,
		}

		resp[0] = reportResp

		repo.On("List", "", ukama.OwnerTypeSubscriber, "", ukama.ReportTypeUnknown,
			isPaid, uint32(0), false).Return(resp, nil)

		s := server.NewReportServer(OrgName, OrgId, repo, nil, nil)
		list, err := s.List(context.TODO(), &pb.ListRequest{
			OwnerType: ownerTypeSubscriber,
			IsPaid:    true,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assertList(t, list, resp)
	})

	t.Run("InvalidOwnerId", func(t *testing.T) {
		repo := &mocks.ReportRepo{}
		repo.On("List", "lol", "", "", "", false, uint32(0), false).Return(nil, errors.New("invalid"))

		s := server.NewReportServer(OrgName, OrgId, repo, nil, nil)
		list, err := s.List(context.TODO(), &pb.ListRequest{
			OwnerId: "lol",
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("InvalidOwnerType", func(t *testing.T) {
		repo := &mocks.ReportRepo{}
		repo.On("List", "", "lol", "", "", uint32(0), false).Return(nil, errors.New("invalid"))

		s := server.NewReportServer(OrgName, OrgId, repo, nil, nil)
		list, err := s.List(context.TODO(), &pb.ListRequest{
			OwnerType: "lol",
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("InvalidReprtType", func(t *testing.T) {
		repo := &mocks.ReportRepo{}
		repo.On("List", "", "", "", "lol", uint32(0), false).Return(nil, errors.New("invalid"))

		s := server.NewReportServer(OrgName, OrgId, repo, nil, nil)
		list, err := s.List(context.TODO(), &pb.ListRequest{
			ReportType: "lol",
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})
}

func TestReportServer_Delete(t *testing.T) {
	t.Run("ReportFound", func(t *testing.T) {
		var reportId = uuid.NewV4()

		reportRepo := &mocks.ReportRepo{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		reportRepo.On("Delete", reportId, mock.Anything).Return(nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := server.NewReportServer(OrgName, OrgId, reportRepo, nil, msgbusClient)

		res, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			ReportId: reportId.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		reportRepo.AssertExpectations(t)
	})

	t.Run("ReportNotFound", func(t *testing.T) {
		var repoortId = uuid.NewV4()

		reportRepo := &mocks.ReportRepo{}

		reportRepo.On("Delete", repoortId, mock.Anything).Return(gorm.ErrRecordNotFound).Once()

		s := server.NewReportServer(OrgName, OrgId, reportRepo, nil, nil)

		res, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			ReportId: repoortId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		reportRepo.AssertExpectations(t)
	})

	t.Run("ReportUUIDInvalid", func(t *testing.T) {
		var reportId = "1"

		reportRepo := &mocks.ReportRepo{}

		s := server.NewReportServer(OrgName, OrgId, reportRepo, nil, nil)

		res, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			ReportId: reportId,
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		reportRepo.AssertExpectations(t)
	})
}

func assertList(t *testing.T, list *pb.ListResponse, resp []db.Report) {
	for idx, paymt := range list.Reports {
		assert.Equal(t, paymt.OwnerId, resp[idx].OwnerId.String())
		assert.Equal(t, paymt.OwnerType, resp[idx].OwnerType.String())
	}
}
