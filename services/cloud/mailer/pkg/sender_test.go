//go:build localtest

package pkg

import (
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/services/common/msgbus"
	"testing"
)

func TestMail_SendEmail(t *testing.T) {
	IsDebugMode = true
	mail := NewMail(&SmtpConfig{
		Host:     "localhost",
		Port:     587,
		Username: "username",
		Password: "password",
		From:     "hello@dev.ukama.com",
	}, "../templates/")

	err := mail.SendEmail(
		&msgbus.MailMessage{
			To:           "denis@ukama.com",
			TemplateName: "test-template",
			Values: map[string]any{
				"mail":    "",
				"Name":    "Denis",
				"Message": "test",
			},
		})
	assert.NoError(t, err)
}
