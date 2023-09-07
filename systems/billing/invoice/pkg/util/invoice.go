package util

import (
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
)

type RawInvoice struct {
	LagoID               uuid.UUID                 `json:"lago_id,omitempty"`
	SequentialID         int                       `json:"sequential_id,omitempty"`
	Number               string                    `json:"number,omitempty"`
	IssuingDate          string                    `json:"issuing_date,omitempty"`
	Status               string                    `json:"status,omitempty"`
	PaymentStatus        string                    `json:"payment_status,omitempty"`
	AmountCents          int                       `json:"amount_cents,omitempty"`
	AmountCurrency       string                    `json:"amount_currency,omitempty"`
	VatAmountCents       int                       `json:"vat_amount_cents,omitempty"`
	VatAmountCurrency    string                    `json:"vat_amount_currency,omitempty"`
	CreditAmountCents    int                       `json:"credit_amount_cents,omitempty"`
	CreditAmountCurrency string                    `json:"credit_amount_currency,omitempty"`
	TotalAmountCents     int                       `json:"total_amount_cents,omitempty"`
	TotalAmountCurrency  string                    `json:"total_amount_currency,omitempty"`
	FileURL              string                    `json:"file_url,omitempty"`
	Legacy               bool                      `json:"legacy,omitempty"`
	Customer             *Customer                 `json:"customer,omitempty"`
	Subscriptions        []Subscription            `json:"subscriptions,omitempty"`
	Fees                 []Fee                     `json:"fees,omitempty"`
	Credits              []InvoiceCredit           `json:"credits,omitempty"`
	Metadata             []InvoiceMetadataResponse `json:"metadata,omitempty"`
}

type InvoiceCreditItem struct {
	LagoID uuid.UUID `json:"lago_id,omitempty"`
	Type   string    `json:"type,omitempty"`
	Code   string    `json:"code,omitempty"`
	Name   string    `json:"name,omitempty"`
}

type InvoiceCredit struct {
	Item           InvoiceCreditItem `json:"item,omitempty"`
	LagoID         uuid.UUID         `json:"lago_id,omitempty"`
	AmountCents    int               `json:"amount_cents,omitempty"`
	AmountCurrency string            `json:"amount_currency,omitempty"`
}

type InvoiceMetadataResponse struct {
	LagoID    uuid.UUID `json:"lago_id,omitempty"`
	Key       string    `json:"key,omitempty"`
	Value     string    `json:"value,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type Fee struct {
	LagoID              uuid.UUID `json:"lago_id,omitempty"`
	LagoGroupID         uuid.UUID `json:"lago_group_id,omitempty"`
	AmountCents         int       `json:"amount_cents,omitempty"`
	AmountCurrency      string    `json:"amount_currenty,omitempty"`
	VatAmountCents      int       `json:"vat_amount_cents,omitempty"`
	VatAmountCurrency   string    `json:"vat_amount_currency,omitempty"`
	TotalAmountCents    int       `json:"total_amount_cents,omitempty"`
	TotalAmountCurrency string    `json:"total_amount_currency,omitempty"`
	Units               string    `json:"units,omitempty"`
	EventsCount         int       `json:"events_count,omitempty"`
	Item                FeeItem   `json:"item,omitempty"`
	// Units               float32   `json:"units,omitempty"`
}

type FeeItem struct {
	Type string `json:"type,omitempty"`
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

type Customer struct {
	LagoID       uuid.UUID `json:"lago_id,omitempty"`
	ExternalID   string    `json:"external_id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Email        string    `json:"email,omitempty"`
	AddressLine1 string    `json:"address_line1,omitempty"`
	AddressLine2 string    `json:"address_line2,omitempty"`
	City         string    `json:"city,omitempty"`
	State        string    `json:"state,omitempty"`
	Zipcode      string    `json:"zipcode,omitempty"`
	Country      string    `json:"country,omitempty"`
	LegalName    string    `json:"legal_name,omitempty"`
	LegalNumber  string    `json:"legal_number,omitempty"`
	LogoURL      string    `json:"logo_url,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	URL          string    `json:"url,omitempty"`
	VatRate      float32   `json:"vat_rate,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

type Subscription struct {
	LagoID             uuid.UUID  `json:"lago_id"`
	LagoCustomerID     uuid.UUID  `json:"lago_customer_id"`
	ExternalCustomerID string     `json:"external_customer_id"`
	ExternalID         string     `json:"external_id"`
	PlanCode           string     `json:"plan_code"`
	Status             string     `json:"status"`
	CreatedAt          *time.Time `json:"created_at"`
	StartedAt          *time.Time `json:"started_at"`
	CanceledAt         *time.Time `json:"canceled_at"`
	TerminatedAt       *time.Time `json:"terminated_at"`
}
