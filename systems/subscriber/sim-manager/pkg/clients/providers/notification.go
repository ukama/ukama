/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package providers

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	res "github.com/ukama/ukama/systems/common/rest"
)


type SendEmailReq struct {
	To      []string `json:"to" validate:"required"`
	TemplateName string `json:"template_name" validate:"required"`
	Values  map[string]interface{}
}


type NotificationClient interface {
	SendEmail(body SendEmailReq) error
}

type notificationClient struct {
	RestClient *res.RestClient
}

type Notification struct {
	Client NotificationClient `json:"notification"`
}

type notificationResponse struct {
	Message string `json:"message"`
	MailId  string `json:"mail_id"`
}

func NewNotificationClient(url string, debug bool) (*notificationClient, error) {
	restClient, err := res.NewRestClient(url, debug)
	if err != nil {
		logrus.Errorf("Failed to connect to %s. Error: %s", url, err.Error())
		return nil, err
	}

	notificationClient := &notificationClient{
		RestClient: restClient,
	}

	return notificationClient, nil
}

func (nc *notificationClient) SendEmail(emailBody SendEmailReq) error {
	errStatus := &res.ErrorMessage{}
	notificationRes := &notificationResponse{}

	resp, err := nc.RestClient.C.R().
		SetError(errStatus).
		SetBody(emailBody).
		Post(nc.RestClient.URL.String() + "/v1/mailer/sendEmail")
	if err != nil {
		logrus.Errorf("Failed to send API request to the notification system. Error: %s", err.Error())
		return err
	}
	err = json.Unmarshal(resp.Body(), &notificationRes)
	if err != nil {
		logrus.Tracef("Failed to deserialize. Error message: %s", err.Error())
		return fmt.Errorf("Failed to deserialization failure: %s", err.Error())
	}

	return nil
}