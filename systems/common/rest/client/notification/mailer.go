/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package notification

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const MailerEndpoint = "/v1/mailer"

type SendEmailReq struct {
	To           []string `json:"to" validate:"required"`
	TemplateName string   `json:"template_name" validate:"required"`
	Values       map[string]interface{}
}

type MailerClient interface {
	SendEmail(body SendEmailReq) error
}

type mailerClient struct {
	u *url.URL
	R *client.Resty
}

func NewMailerClient(h string, options ...client.Option) *mailerClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &mailerClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (m *mailerClient) SendEmail(body SendEmailReq) error {
	log.Debugf("Sending email req: %v", body)

	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("email body marshal error. error: %w", err)
	}

	_, err = m.R.Post(m.u.String()+MailerEndpoint+"/sendEmail", b)
	if err != nil {
		log.Errorf("SendEmail failure. error: %s", err.Error())

		return fmt.Errorf("SendEmail failure: %w", err)
	}

	log.Debugf("Email sent successfully")
	return nil
}
