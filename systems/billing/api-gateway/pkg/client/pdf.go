package client

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

const FileServerEndpoint = "/pdf/"

var ErrInvoicePDFNotFound = errors.New("invoice PDF file not found")

type PdfClient interface {
	GetPdf(invoiceId string) ([]byte, error)
}

type pdfClient struct {
	R *rest.RestClient
}

func NewPdfClient(url string, debug bool) (*pdfClient, error) {
	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		log.Errorf("Can't conncet to %s url. Error %s", url, err.Error())

		return nil, err
	}

	N := &pdfClient{
		R: f,
	}

	return N, nil
}

func (p *pdfClient) GetPdf(invoiceId string) ([]byte, error) {
	errStatus := &rest.ErrorMessage{}

	resp, err := p.R.C.R().
		SetError(errStatus).
		Get(p.R.URL.String() + FileServerEndpoint + invoiceId + ".pdf")

	if err != nil {
		log.Errorf("Failed to send request to billing/invoice. Error %s", err.Error())

		return nil, fmt.Errorf("api request to billing system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch invoice file. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		if resp.StatusCode() == http.StatusNotFound {
			return nil, fmt.Errorf("%w: %s", ErrInvoicePDFNotFound, errStatus.Message)
		}

		return nil, fmt.Errorf("error while retrieving invoice file: %s", errStatus.Message)
	}

	return resp.Body(), nil
}
