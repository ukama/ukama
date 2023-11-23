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
	SubscriberId string `example:"{{SubscriberUUID}}" form:"subscriber_id" json:"subscriber" query:"subscriber" binding:"required"`
	NetworkId    string `example:"{{NetworkUUID}}" form:"network_id" json:"network" query:"network" binding:"required"`
}

type GetInvoiceRequest struct {
	InvoiceId string `example:"{{InvoiceUUID}}" path:"invoice_id" validate:"required"`
}

type WebHookRequest struct {
	WebhookType string          `example:"webhook-type" json:"webhook_type" validate:"required"`
	ObjectType  string          `example:"object-type" json:"object_type" validate:"required"`
	Invoice     util.RawInvoice `example:"{}" json:"invoice" validate:"required"`
}
