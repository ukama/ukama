/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import "github.com/ukama/ukama/systems/billing/report/pkg/util"

type GetReportsRequest struct {
	OwnerId   string `form:"owner_id" json:"owner_id" query:"owner_id" binding:"required"`
	OwnerType string `form:"owner_type" json:"owner_type" query:"owner_type" binding:"required"`
	NetworkId string `form:"network_id" json:"network_id" query:"network_id" binding:"required"`
	ReortType string `form:"report_type" json:"report_type" query:"report_type" binding:"required"`
	IsPaid    bool   `form:"is_paid" json:"is_paid" query:"is_paid" binding:"required"`
	Count     uint32 `form:"count" json:"count" query:"count" binding:"required"`
	Sort      bool   `form:"sort" json:"sort" query:"sort" binding:"required"`
}

type GetReportRequest struct {
	ReportId string `example:"{{ReportUUID}}" path:"report_id" validate:"required"`
	AsPdf    bool   `form:"as_pdf" json:"as_pdf" query:"as_pdf" binding:"required"`
}

type WebHookRequest struct {
	WebhookType string          `example:"webhook-type" json:"webhook_type" validate:"required"`
	ObjectType  string          `example:"object-type" json:"object_type" validate:"required"`
	Invoice     util.RawInvoice `example:"{}" json:"invoice" validate:"required"`
}
