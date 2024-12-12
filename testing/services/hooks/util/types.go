/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package util

import stripelib "github.com/stripe/stripe-go/v78"

type PaymentIntent struct {
	*stripelib.PaymentIntent
}

type Deposit struct {
	DepositId                string                               `json:"depositId,omitempty"`
	Status                   string                               `json:"status,omitempty"`
	Amount                   string                               `json:"amount,omitempty"`
	RequestedAmount          string                               `json:"requestedAmount,omitempty"`
	Currency                 string                               `json:"currency,omitempty"`
	Country                  string                               `json:"country,omitempty"`
	Correspondent            string                               `json:"correspondent,omitempty"`
	Payer                    Payer                                `json:"payer,omitempty"`
	CustomerTimestamp        string                               `json:"customerTimestamp,omitempty"`
	StatementDescription     string                               `json:"statementDescription,omitempty"`
	Created                  string                               `json:"created,omitempty"`
	DepositedAmount          string                               `json:"depositedAmount,omitempty"`
	RespondedByPayer         string                               `json:"respondedByPayer,omitempty"`
	CorrespondentIds         map[string]interface{}               `json:"correspondentIds,omitempty"`
	SuspiciousActivityReport []SuspiciousDepositTransactionReport `json:"SuspiciousActivityReport,omitempty"`
	RejectionReason          DepositRejectionReason               `json:"rejectionReason,omitempty"`
	FailureReason            DepositFailureReason                 `json:"failureReason,omitempty"`
}

type Payer struct {
	Type    string  `json:"type" validate:"required"`
	Address Address `json:"address" validate:"required"`
}

type Address struct {
	Value string `json:"value" validate:"required"`
}

type DepositRejectionReason struct {
	RejectionCode    string `json:"rejectionCode,omitempty"`
	RejectionMessage string `json:"rejectionMessage,omitempty"`
}

type DepositFailureReason struct {
	FailureCode    string `json:"failureCode,omitempty"`
	FailureMessage string `json:"failureMessage,omitempty"`
}

type SuspiciousDepositTransactionReport struct {
	ActivityType string `json:"activityType,omitempty"`
	Comment      string `json:"comment,omitempty"`
}
