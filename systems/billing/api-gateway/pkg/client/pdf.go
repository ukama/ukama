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

var ErrReportPDFNotFound = errors.New("report PDF file not found")

type Pdf interface {
	GetPdf(reportId string) ([]byte, error)
}

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

func (p *pdf) GetPdf(reportId string) ([]byte, error) {
	errStatus := &rest.ErrorMessage{}

	resp, err := p.R.C.R().
		SetError(errStatus).
		Get(p.R.URL.String() + FileServerEndpoint + reportId + ".pdf")

	if err != nil {
		log.Errorf("Failed to send request to report/pdf. Error %s", err.Error())

		return nil, fmt.Errorf("api request to report service failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch report file. HTTP resp code %d and Error message is %s",
			resp.StatusCode(), errStatus.Message)

		if resp.StatusCode() == http.StatusNotFound {
			return nil, fmt.Errorf("%w: %s", ErrReportPDFNotFound, errStatus.Message)
		}

		return nil, fmt.Errorf("error while retrieving report file: %s", errStatus.Message)
	}

	return resp.Body(), nil
}
