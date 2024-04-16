/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import "github.com/ukama/ukama/systems/billing/invoice/pkg/util"

type GetInvoicesRequest struct {
	InvoiceeId   string `form:"invoicee_id" json:"invoicee_id" query:"invoicee_id" binding:"required"`
	InvoiceeType string `form:"invoicee_type" json:"invoicee_type" query:"invoicee_type" binding:"required"`
	NetworkId    string `form:"network_id" json:"network_id" query:"network_id" binding:"required"`
	IsPaid       bool   `form:"is_paid" json:"is_paid" query:"is_paid" binding:"required"`
	Count        uint32 `form:"count" json:"count" query:"count" binding:"required"`
	Sort         bool   `form:"sort" json:"sort" query:"sort" binding:"required"`
}

type GetInvoiceRequest struct {
	InvoiceId string `example:"{{InvoiceUUID}}" path:"invoice_id" validate:"required"`
	AsPdf     bool   `form:"as_pdf" json:"as_pdf" query:"as_pdf" binding:"required"`
}

type WebHookRequest struct {
	WebhookType string          `example:"webhook-type" json:"webhook_type" validate:"required"`
	ObjectType  string          `example:"object-type" json:"object_type" validate:"required"`
	Invoice     util.RawInvoice `example:"{}" json:"invoice" validate:"required"`
}
