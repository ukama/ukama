/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ukama/ukama/systems/common/rest"

	log "github.com/sirupsen/logrus"
)

const FileServerEndpoint = "/pdf/"

var ErrInvoicePDFNotFound = errors.New("invoice PDF file not found")

type Pdf interface {
	GetPdf(invoiceId string) ([]byte, error)
}

// TODO: try again with common/rest/client resty
type pdf struct {
	R *rest.RestClient
}

func NewPdfClient(fileHost string, debug bool) *pdf {
	f, err := rest.NewRestClient(fileHost, debug)
	if err != nil {
		log.Fatalf("Can't conncet to pdf host  %s url. Error %s", fileHost, err.Error())
	}

	return &pdf{
		R: f,
	}
}

func (p *pdf) GetPdf(invoiceId string) ([]byte, error) {
	errStatus := &rest.ErrorMessage{}

	resp, err := p.R.C.R().
		SetError(errStatus).
		Get(p.R.URL.String() + FileServerEndpoint + invoiceId + ".pdf")

	if err != nil {
		log.Errorf("Failed to send request to invoice/pdf. Error %s", err.Error())

		return nil, fmt.Errorf("api request to invoice service failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch invoice file. HTTP resp code %d and Error message is %s",
			resp.StatusCode(), errStatus.Message)

		if resp.StatusCode() == http.StatusNotFound {
			return nil, fmt.Errorf("%w: %s", ErrInvoicePDFNotFound, errStatus.Message)
		}

		return nil, fmt.Errorf("error while retrieving invoice file: %s", errStatus.Message)
	}

	return resp.Body(), nil
}
