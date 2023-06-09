package client

import (
	"net/http"
	"time"

	"github.com/go-playground/validator"
	res "github.com/ukama/ukama/systems/mailer/api-gateway/pkg/rest"

	"gopkg.in/gomail.v2"
)

type SmtpConfig struct {
	From     string `default:"hello@dev.ukama.com" validate:"required"`
	Host     string `default:"localhost" validate:"required"`
	Port     int    `default:"25" validate:"required"`
	Password string
	Username string
}

type MailerClient struct {
	client   *http.Client
	host     string
	port     int
	username string
	password string
}

func NewMailerClient(host string, port int, timeout time.Duration, username string, password string) *MailerClient {
	return &MailerClient{
		client:   &http.Client{Timeout: timeout},
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (c *MailerClient) SendEmail(to string, message string) (res.SendEmailRes, error) {
	var smtpConfig SmtpConfig
	smtpConfig.From = c.username
	smtpConfig.Host = c.host
	smtpConfig.Port = c.port
	smtpConfig.Password = c.password
	smtpConfig.Username = c.username

	validate := validator.New()
	if err := validate.Struct(smtpConfig); err != nil {
		return res.SendEmailRes{}, err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpConfig.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Hello")
	m.SetBody("text/plain", message)

	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.Password)

	if err := d.DialAndSend(m); err != nil {
		return res.SendEmailRes{}, err
	}

	return res.SendEmailRes{
		Message: "Email sent successfully",
	}, nil
}
