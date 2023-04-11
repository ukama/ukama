package rest

type GetInvoicesRequest struct {
	SubscriberId string `example:"{{SubscriberUUID}}" form:"subscriber_id" json:"subscriber" query:"subscriber" binding:"required" validate:"required"`
}

type GetInvoiceRequest struct {
	InvoiceId string `example:"{{InvoiceUUID}}" path:"invoice_id" validate:"required"`
}

type AddInvoiceRequest struct {
	SubscriberId string `example:"SubscriberUUID"  json:"subscriber_id" validate:"required"`
	RawInvoice   string `example:"mesh-network" json:"raw_invoice" validate:"required"`
}
