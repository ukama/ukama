/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package clients

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/testing/services/hooks/util"

	log "github.com/sirupsen/logrus"
)

const (
	invalidArgCode  = 3
	DepositEndpoint = "/deposits"

	DepositStatusAccepted = "ACCEPTED"
	DepositStatusRejected = "REJECTED"

	deserializeLogMsg      = "Failed to deserialize %s info. Error message is: %s"
	deserializeErrorMsg    = "%s info deserialization failure: %w"
	resourceLogMsg         = "%s info: %v"
	requestMarshalErrorMsg = "request marshal error. error: %w"
)

type PawapayClient interface {
	GetDeposit(string) (*util.Deposit, error)
	AddDeposit(AddDepositRequest) (*util.Deposit, error)
	ResendDepositCallback(CallbackRequest) (*util.Deposit, error)
	PredictMno(MsisdnRequest) (*Operator, error)
	GetMnosAvailability() ([]CountryOperators, error)
}

type pawapayClient struct {
	u *url.URL
	R *client.Resty
}

func NewPawapayClient(h string, key string) *pawapayClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &pawapayClient{
		u: u,
		R: client.NewResty(client.WithBearer(key), client.WithError(&Err{}),
			client.WithDebug(), client.WithContentTypeJSON()),
	}
}

func (p *pawapayClient) GetDeposit(id string) (*util.Deposit, error) {
	log.Debugf("Getting deposit: %v", id)

	dep := []util.Deposit{}

	resp, err := p.R.Get(p.u.String() + DepositEndpoint + "/" + id)
	if err != nil {
		log.Errorf("GetDeposit failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetDeposit failure: %w", err)
	}
	err = json.Unmarshal(resp.Body(), &dep)
	if err != nil {
		log.Tracef(deserializeLogMsg, "deposit", err.Error())

		return nil, fmt.Errorf(deserializeErrorMsg, "deposit", err)
	}

	log.Infof(resourceLogMsg, "Deposit", dep)

	if len(dep) == 0 {
		return nil, fmt.Errorf("deposit not found")
	}

	return &dep[0], nil
}

func (p *pawapayClient) AddDeposit(req AddDepositRequest) (*util.Deposit, error) {
	log.Debugf("Adding deposit: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf(requestMarshalErrorMsg, err)
	}

	dep := util.Deposit{}

	resp, err := p.R.Post(p.u.String()+DepositEndpoint, b)
	if err != nil {
		log.Errorf("AddDeposit failure. error: %s", err.Error())

		return nil, fmt.Errorf("AddDeposit failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &dep)
	if err != nil {
		log.Tracef(deserializeLogMsg, "deposit", err.Error())

		return nil, fmt.Errorf(deserializeErrorMsg, "deposit", err)
	}

	if dep.Status == DepositStatusRejected {
		log.Errorf("Failed to submit a deposit. Status: %s, Reason: %s (%s)",
			dep.Status, dep.RejectionReason.RejectionMessage, dep.RejectionReason.RejectionCode)

		err = fmt.Errorf("status: %s, reason: %s (%s)", dep.Status,
			dep.RejectionReason.RejectionMessage, dep.RejectionReason.RejectionCode)

		return nil, fmt.Errorf("failed to submit a deposit. %w",
			client.ErrorStatus{
				StatusCode: invalidArgCode,
				RawError:   err,
			})
	}

	log.Infof(resourceLogMsg, "Deposit", dep)

	return &dep, nil
}

func (p *pawapayClient) ResendDepositCallback(req CallbackRequest) (*util.Deposit, error) {
	log.Debugf("Re-requesting callback for deposit: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf(requestMarshalErrorMsg, err)
	}

	dep := util.Deposit{}

	resp, err := p.R.Post(p.u.String()+DepositEndpoint+"/resend-callback", b)
	if err != nil {
		log.Errorf("ResendDepositCallback failure. error: %s", err.Error())

		return nil, fmt.Errorf("ResendDepositCallback failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &dep)
	if err != nil {
		log.Tracef(deserializeLogMsg, "deposit", err.Error())

		return nil, fmt.Errorf(deserializeErrorMsg, "deposit", err)
	}

	log.Infof(resourceLogMsg, "Deposit", dep)

	return &dep, nil
}

func (p *pawapayClient) PredictMno(req MsisdnRequest) (*Operator, error) {
	log.Debugf("Predicting MNO correspondent: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf(requestMarshalErrorMsg, err)
	}

	op := Operator{}

	resp, err := p.R.Post(p.u.String()+"/v1/predict-correspondent", b)
	if err != nil {
		log.Errorf("PredictMno failure. error: %s", err.Error())

		return nil, fmt.Errorf("PredictMno failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &op)
	if err != nil {
		log.Tracef(deserializeLogMsg, "operator", err.Error())

		return nil, fmt.Errorf(deserializeErrorMsg, "operator", err)
	}

	log.Infof(resourceLogMsg, "Operator", op)

	return &op, nil
}

func (p *pawapayClient) GetMnosAvailability() ([]CountryOperators, error) {
	log.Debug("Getting MNOs availability")

	ops := []CountryOperators{}

	resp, err := p.R.Get(p.u.String() + "/availability")
	if err != nil {
		log.Errorf("GetMnoAvailability failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetMnoAvailability failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &ops)
	if err != nil {
		log.Tracef(deserializeLogMsg, "country operators", err.Error())

		return nil, fmt.Errorf(deserializeErrorMsg, "country operators", err)
	}

	log.Infof(resourceLogMsg, "Country operators", ops)

	return ops, nil
}

type AddDepositRequest struct {
	DepositId            string     `json:"depositId" validate:"required"`
	Amount               string     `json:"amount" validate:"required"`
	Currency             string     `json:"currency" validate:"required"`
	Country              string     `json:"country" validate:"required"`
	Correspondent        string     `json:"correspondent" validate:"required"`
	Payer                util.Payer `json:"payer" validate:"required"`
	CustomerTimestamp    string     `json:"customerTimestamp" validate:"required"`
	StatementDescription string     `json:"statementDescription" validate:"required"`
	PreAuthorisationCode string     `json:"preAuthorisationCode"`
}

type CallbackRequest struct {
	DepositId string `json:"depositId" validate:"required"`
}

type MsisdnRequest struct {
	Msisdn string `json:"msisdn" validate:"required"`
}

type Operator struct {
	Country       string `json:"country,omitempty"`
	Operator      string `json:"operator,omitempty"`
	Correspondent string `json:"correspondent,omitempty"`
	Msisdn        string `json:"msisdn,omitempty"`
}

type CountryOperators struct {
	Country        string          `json:"country,omitempty"`
	Correspondents []Correspondent `json:"correspondents,omitempty"`
}

type Correspondent struct {
	Correspondent  string      `json:"correspondent,omitempty"`
	OperationTypes []Operation `json:"operationTypes,omitempty"`
}

type Operation struct {
	OperationType string `json:"operationType,omitempty"`
	Status        string `json:"status,omitempty"`
}

type Err struct {
	ErrorId      string `json:"errorId,omitempty"`
	ErrorCode    int    `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func (e Err) Error() string {
	return fmt.Sprintf("%s (provider statuscode: %d, provider errorid: %s)",
		e.ErrorMessage, e.ErrorCode, e.ErrorId)
}
