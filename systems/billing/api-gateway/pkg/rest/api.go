package rest

import (
	"github.com/ukama/ukama/systems/billing/invoice/pkg/util"
	"github.com/ukama/ukama/systems/common/rest"
)

type GetInvoicesRequest struct {
	rest.BaseRequest
	SubscriberId string `example:"{{SubscriberUUID}}" form:"subscriber_id" json:"subscriber" query:"subscriber" binding:"required" validate:"required"`
}

type GetInvoiceRequest struct {
	rest.BaseRequest
	InvoiceId string `example:"{{InvoiceUUID}}" path:"invoice_id" validate:"required"`
}

type HandleWebHookRequest struct {
	rest.BaseRequest
	WebhookType string          `example:"webhook-type" json:"webhook_type" validate:"required"`
	ObjectType  string          `example:"object-type" json:"object_type" validate:"required"`
	Invoice     util.RawInvoice `example:"{}" json:"invoice" validate:"required"`
}
