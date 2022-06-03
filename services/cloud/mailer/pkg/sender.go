package pkg

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/msgbus"
	"net/smtp"
	"path/filepath"
	"text/template"
)

type Sender struct {
	smtpConf *SmtpConfig
	// path to template directory
	tmplPath string
}

func NewMail(smtpConf *SmtpConfig, tmplPath string) *Sender {
	return &Sender{
		smtpConf: smtpConf,
		tmplPath: tmplPath}
}

func (m *Sender) SendEmail(data *msgbus.MailMessage) error {
	c, err := smtp.Dial(fmt.Sprintf("%s:%d", m.smtpConf.Host, m.smtpConf.Port))
	if err != nil {
		return err
	}
	defer c.Close()

	body, err := m.prepareMsg(data)
	if err != nil {
		return err
	}

	// Sending email.
	err = c.Mail(m.smtpConf.From)
	if err != nil {
		return err
	}
	if err = c.Rcpt(data.To); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.Write(body.Bytes())
	if err != nil {
		return err
	}

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = c.Quit()
	if err != nil {
		logrus.Warningf("Error sending quit command to smtp: %s", err)
	}

	fmt.Println("Email Sent!")
	return nil
}

func (m *Sender) prepareMsg(data *msgbus.MailMessage) (bytes.Buffer, error) {
	tmplName := data.TemplateName
	if filepath.Ext(tmplName) == "" {
		tmplName += ".tmpl"
	}

	t, err := template.ParseFiles(filepath.Join(m.tmplPath, tmplName))
	if err != nil {
		return bytes.Buffer{}, err
	}

	var body bytes.Buffer

	err = t.Execute(&body, data)
	if err != nil {
		return bytes.Buffer{}, err
	}

	if IsDebugMode {
		logrus.Printf("%s", body.String())
	}
	return body, nil
}
