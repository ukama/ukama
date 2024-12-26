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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/billing/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/billing/report/pb/gen/mocks"
	"github.com/ukama/ukama/systems/common/uuid"

	pb "github.com/ukama/ukama/systems/billing/report/pb/gen"
)

func TestReportClient_Add(t *testing.T) {
	t.Run("OwnerIdValid", func(t *testing.T) {
		var bc = &mocks.ReportServiceClient{}

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
	"file_url": "https://getlago.com/report/file",
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

		reportReq := &pb.AddRequest{
			RawReport: raw,
		}

		reportResp := &pb.ReportResponse{Report: &pb.Report{
			Id:        uuid.NewV4().String(),
			OwnerId:   "5eb02857-a71e-4ea2-bcf9-57d3a41bc6ba",
			NetworkId: uuid.NewV4().String(),
		}}

		bc.On("Add", mock.Anything, reportReq).Return(reportResp, nil)

		i := client.NewReportFromClient(bc)

		resp, err := i.Add(raw)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		bc.AssertExpectations(t)
	})

	t.Run("OwnerIdNotValid", func(t *testing.T) {
		var bc = &mocks.ReportServiceClient{}

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
	"file_url": "https://getlago.com/report/file",
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

		reportReq := &pb.AddRequest{
			RawReport: raw,
		}

		bc.On("Add", mock.Anything, reportReq).Return(nil,
			status.Errorf(codes.InvalidArgument, "invalid ownerId"))

		i := client.NewReportFromClient(bc)

		resp, err := i.Add(raw)

		assert.Error(t, err)
		assert.Nil(t, resp)
		bc.AssertExpectations(t)
	})
}

func TestReportClient_Get(t *testing.T) {
	var ic = &mocks.ReportServiceClient{}

	t.Run("ReportFound", func(t *testing.T) {
		reportId := uuid.NewV4()

		reportReq := &pb.GetRequest{
			ReportId: reportId.String(),
		}

		reportResp := &pb.ReportResponse{Report: &pb.Report{
			Id:      reportId.String(),
			OwnerId: uuid.NewV4().String(),
		}}

		ic.On("Get", mock.Anything, reportReq).Return(reportResp, nil)

		i := client.NewReportFromClient(ic)

		resp, err := i.Get(reportId.String())

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, resp.Report.Id, reportId.String())
		ic.AssertExpectations(t)
	})

	t.Run("ReportNotFound", func(t *testing.T) {
		reportId := uuid.NewV4()

		reportReq := &pb.GetRequest{
			ReportId: reportId.String(),
		}

		ic.On("Get", mock.Anything, reportReq).Return(nil,
			status.Errorf(codes.NotFound, "report not found"))

		i := client.NewReportFromClient(ic)

		resp, err := i.Get(reportId.String())

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "not found")
		ic.AssertExpectations(t)
	})
}
func TestReportClient_List(t *testing.T) {
	var (
		ic = &mocks.ReportServiceClient{}

		reportId            = uuid.NewV4().String()
		ownerId             = uuid.NewV4().String()
		networkId           = uuid.NewV4().String()
		ownerTypeSubscriber = "subscriber"
		// ownerTypeOrg        = "org"
		reportTypeInvoice = "invoice"

		isPaid = true
	)

	listReq := &pb.ListRequest{
		OwnerId:    ownerId,
		OwnerType:  ownerTypeSubscriber,
		NetworkId:  networkId,
		ReportType: reportTypeInvoice,
		IsPaid:     isPaid,
		Count:      uint32(1),
		Sort:       true}

	listResp := &pb.ListResponse{Reports: []*pb.Report{
		&pb.Report{
			Id:        reportId,
			OwnerId:   ownerId,
			OwnerType: ownerTypeSubscriber,
			NetworkId: networkId,
			IsPaid:    isPaid,
		}}}

	ic.On("List", mock.Anything, listReq).Return(listResp, nil)

	n := client.NewReportFromClient(ic)

	resp, err := n.List(ownerId, ownerTypeSubscriber, networkId,
		reportTypeInvoice, isPaid, uint32(1), true)

	assert.NoError(t, err)
	assert.Equal(t, resp.Reports[0].Id, reportId)
	ic.AssertExpectations(t)
}

func TestReportClient_Remove(t *testing.T) {
	var bc = &mocks.ReportServiceClient{}

	t.Run("ReportFound", func(t *testing.T) {
		reportId := uuid.NewV4()

		reportReq := &pb.DeleteRequest{
			ReportId: reportId.String(),
		}

		bc.On("Delete", mock.Anything, reportReq).Return(nil, nil)

		i := client.NewReportFromClient(bc)

		err := i.Remove(reportId.String())

		assert.NoError(t, err)
		bc.AssertExpectations(t)
	})

	t.Run("ReportNotFound", func(t *testing.T) {
		reportId := uuid.NewV4()

		reportReq := &pb.DeleteRequest{
			ReportId: reportId.String(),
		}

		bc.On("Delete", mock.Anything, reportReq).Return(nil,
			status.Errorf(codes.NotFound, "report not found"))

		i := client.NewReportFromClient(bc)

		err := i.Remove(reportId.String())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		bc.AssertExpectations(t)
	})
}
