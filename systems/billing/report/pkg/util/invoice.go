/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package util

import (
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
)

type Currency string

type RawInvoice struct {
	LagoID                            uuid.UUID                 `json:"lago_id,omitempty"`
	SequentialID                      int                       `json:"sequential_id,omitempty"`
	Number                            string                    `json:"number,omitempty"`
	IssuingDate                       string                    `json:"issuing_date,omitempty"`
	PaymentDueDate                    string                    `json:"payment_due_date,omitempty"`
	PaymentOverdue                    bool                      `json:"payment_overdue,omitempty"`
	InvoiceType                       string                    `json:"invoice_type,omitempty"`
	Status                            string                    `json:"status,omitempty"`
	PaymentStatus                     string                    `json:"payment_status,omitempty"`
	Currency                          string                    `json:"currency,omitempty"`
	FeesAmountCents                   int                       `json:"fees_amount_cents,omitempty"`
	TaxesAmountCents                  int                       `json:"taxes_amount_cents,omitempty"`
	SubTotalExcludingTaxesAmountCents int                       `json:"sub_total_excluding_taxes_amount_cents,omitempty"`
	SubTotalIncludingTaxesAmountCents int                       `json:"sub_total_including_taxes_amount_cents,omitempty"`
	TotalAmountCents                  int                       `json:"total_amount_cents,omitempty"`
	FileURL                           string                    `json:"file_url,omitempty"`
	Legacy                            bool                      `json:"legacy,omitempty"`
	Customer                          *Customer                 `json:"customer,omitempty"`
	Subscriptions                     []Subscription            `json:"subscriptions,omitempty"`
	Fees                              []Fee                     `json:"fees,omitempty"`
	Credits                           []InvoiceCredit           `json:"credits,omitempty"`
	Metadata                          []InvoiceMetadataResponse `json:"metadata,omitempty"`
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
	TaxesAmountCents    int       `json:"taxes_amount_cents,omitempty"`
	TaxesPreciseAmount  string    `json:"taxes_precise_amount,omitempty"`
	TotalAmountCents    int       `json:"total_amount_cents,omitempty"`
	TotalAmountCurrency string    `json:"total_amount_currency,omitempty"`
	Units               string    `json:"units,omitempty"`
	EventsCount         int       `json:"events_count,omitempty"`
	Item                FeeItem   `json:"item,omitempty"`
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
	Currency     string    `json:"currency,omitempty"`
	Timezone     string    `json:"timezone,omitempty"`
	VatRate      float32   `json:"vat_rate,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

type Subscription struct {
	LagoID             uuid.UUID  `json:"lago_id"`
	LagoCustomerID     uuid.UUID  `json:"lago_customer_id"`
	ExternalCustomerID string     `json:"external_customer_id"`
	ExternalID         string     `json:"external_id"`
	PlanCode           string     `json:"plan_code"`
	Name               string     `json:"name"`
	Status             string     `json:"status"`
	CreatedAt          *time.Time `json:"created_at"`
	StartedAt          *time.Time `json:"started_at"`
	CanceledAt         *time.Time `json:"canceled_at"`
	TerminatedAt       *time.Time `json:"terminated_at"`
}

type PlanInterval string

type Plan struct {
	LagoID                   uuid.UUID    `json:"lago_id"`
	Name                     string       `json:"name,omitempty"`
	Code                     string       `json:"code,omitempty"`
	Interval                 PlanInterval `json:"interval,omitempty"`
	Description              string       `json:"description,omitempty"`
	AmountCents              int          `json:"amount_cents,omitempty"`
	AmountCurrency           Currency     `json:"amount_currency,omitempty"`
	PayInAdvance             bool         `json:"pay_in_advance,omitempty"`
	BillChargeMonthly        bool         `json:"bill_charge_monthly,omitempty"`
	ActiveSubscriptionsCount int          `json:"active_subscriptions_count,omitempty"`
	DraftInvoicesCount       int          `json:"draft_invoices_count,omitempty"`
	Charges                  []Charge     `json:"charges,omitempty"`
}

type ChargeModel string

type Charge struct {
	LagoID               uuid.UUID              `json:"lago_id,omitempty"`
	LagoBillableMetricID uuid.UUID              `json:"lago_billable_metric_id,omitempty"`
	BillableMetricCode   string                 `json:"billable_metric_code,omitempty"`
	ChargeModel          ChargeModel            `json:"charge_model,omitempty"`
	CreatedAt            time.Time              `json:"created_at,omitempty"`
	PayInAdvance         bool                   `json:"pay_in_advance,omitempty"`
	MinAmountCents       int                    `json:"min_amount_cents,omitempty"`
	Properties           map[string]interface{} `json:"properties,omitempty"`
	GroupProperties      []GroupProperties      `json:"group_properties,omitempty"`
}

type GroupProperties struct {
	GroupID uuid.UUID              `json:"group_id"`
	Values  map[string]interface{} `json:"values"`
}
