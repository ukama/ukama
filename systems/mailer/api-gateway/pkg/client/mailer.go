package client

import (
	"fmt"
	"net/http"
	"time"

	res "github.com/ukama/ukama/systems/mailer/api-gateway/pkg/rest"
)

type MailerClient struct {
	client  *http.Client
	host    string
	port    int
	username string
	password string
	
}

func NewMailerClient(host string,port int ,timeout time.Duration,username string ,password string) *MailerClient {
	return &MailerClient{
		client: &http.Client{
			Timeout: 10 * time.Second, // Set a timeout for HTTP requests
		},
		host: host,
		port: port,
		username: username,
		password: password,

	}

}

func (c *MailerClient) sendEmail(to string, message string) (res.SendEmailRes, error) {
	url := c.baseURL + "/v1/sendEmail"

	// Send POST request to the mailer API GW
	resp, err := c.client.Post(url, "application/json", nil)
	if err != nil {
		return res.SendEmailRes{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return res.SendEmailRes{}, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Handle the response body if needed
	// ...

	// Return the response
	return res.SendEmailRes{
		Message: "Email sent successfully",
	}, nil
}
