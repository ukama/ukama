package utils

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

type Resty struct {
	C *resty.Client
}

func NewResty() *Resty {

	c := resty.New()
	c.SetDebug(true)

	return &Resty{
		C: c,
	}
}

func (r *Resty) Put(url string, b []byte) (*resty.Response, error) {

	resp := &resty.Response{}
	var err error
	errStatus := &rest.ErrorResponse{}
	if b != nil {
		resp, err = r.C.R().
			SetError(errStatus).
			SetBody(b).
			Put(url)
	} else {
		resp, err = r.C.R().
			SetError(errStatus).
			Put(url)
	}
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if resp.StatusCode() != http.StatusCreated {
		log.Errorf("Failed to perform PUT operation on %s HTTP resp code %d and Error message is %s", url, resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return resp, nil
}

func (r *Resty) Post(url string, b []byte) (*resty.Response, error) {

	resp := &resty.Response{}
	var err error
	errStatus := &rest.ErrorResponse{}
	if b != nil {
		resp, err = r.C.R().
			SetError(errStatus).
			SetBody(b).
			Post(url)
	} else {
		resp, err = r.C.R().
			SetError(errStatus).
			Post(url)
	}

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if !((resp.StatusCode() >= http.StatusOK) && resp.StatusCode() < http.StatusBadRequest) {
		log.Errorf("Failed to perform POST operation on %s HTTP resp code %d and Error message is %s", url, resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return resp, nil
}

func (r *Resty) Get(url string) (*resty.Response, error) {

	errStatus := &rest.ErrorResponse{}

	resp, err := r.C.R().
		SetError(errStatus).
		Get(url)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Failed to perform GET on %s operation HTTP resp code %d and Error message is %s", url, resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return resp, nil
}

func (r *Resty) Patch(url string, b []byte) (*resty.Response, error) {

	errStatus := &rest.ErrorResponse{}

	resp, err := r.C.R().
		SetError(errStatus).
		SetBody(b).
		Patch(url)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Failed to perform PATCH operation on %s HTTP resp code %d and Error message is %s", url, resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return resp, nil
}

func (r *Resty) Delete(url string) (*resty.Response, error) {

	errStatus := &rest.ErrorResponse{}

	resp, err := r.C.R().
		SetError(errStatus).
		Delete(url)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Failed to perform Delete operation on %s HTTP resp code %d and Error message is %s", url, resp.StatusCode(), errStatus.Error)
		return nil, fmt.Errorf("rest api failure. error : %s", errStatus.Error)
	}

	return resp, nil
}
