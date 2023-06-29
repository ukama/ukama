package providers

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

type SendEmailReq struct {
	To      []string          `json:"to" validate:"required"`
	Subject string            `json:"subject" validate:"required"`
	Body    string            `json:"body" validate:"required"`
	Values  map[string]string `json:"values"`
}

type NotificationClient interface {
	SendEmail(body SendEmailReq) error
}

type notificationClient struct {
	RestClient *RestClient
}

type Notification struct {
	Client NotificationClient `json:"notification"`
}

type notificationResponse struct {
	Message string `json:"message"`
	MailId  string `json:"mail_id"`
}

func NewNotificationClient(url string, debug bool) (*notificationClient, error) {
	restClient, err := NewRestClient(url, debug)
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
	errStatus := &rest.ErrorMessage{}
	notificationRes := &notificationResponse{}

	resp, err := nc.RestClient.Client.R().
		SetError(errStatus).
		SetBody(emailBody).
		Post(nc.RestClient.URL.String() + "/v1/mail/sendEmail")
	if err != nil {
		logrus.Errorf("Failed to send API request to the notification system. Error: %s", err.Error())
		return err
	}

	if !resp.IsSuccess() {
		logrus.Tracef("Failed to fetch network info. HTTP response code: %d, Error message: %s", resp.StatusCode(), errStatus.Message)
		return fmt.Errorf("Network Info failure: %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &notificationRes)
	if err != nil {
		logrus.Tracef("Failed to deserialize network info. Error message: %s", err.Error())
		return fmt.Errorf("Network info deserialization failure: %s", err.Error())
	}

	return nil
}

type RestClient struct {
	Client *resty.Client
	URL    *url.URL
}

func NewRestClient(path string, debug bool) (*RestClient, error) {
	parsedURL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	client := resty.New()
	client.SetDebug(debug)

	restClient := &RestClient{
		Client: client,
		URL:    parsedURL,
	}

	logrus.Tracef("Client created %+v for %s", restClient, restClient.URL.String())
	return restClient, nil
}
