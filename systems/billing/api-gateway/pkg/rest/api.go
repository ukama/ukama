package rest

import "github.com/ukama/ukama/systems/billing/invoice/pkg/util"

type GetInvoicesRequest struct {
	SubscriberId string `example:"{{SubscriberUUID}}" form:"subscriber_id" json:"subscriber" query:"subscriber" binding:"required" validate:"required"`
}

type GetInvoiceRequest struct {
	InvoiceId string `example:"{{InvoiceUUID}}" path:"invoice_id" validate:"required"`
}

type HandleWebHookRequest struct {
	WebhookType string          `example:"webhook-type" json:"webhook_type" validate:"required"`
	ObjectType  string          `example:"object-type" json:"object_type" validate:"required"`
	Invoice     util.RawInvoice `example:"{}" json:"invoice" validate:"required"`
}
