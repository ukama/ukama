package server

import (
	"context"
	"testing"

	"github.com/ukama/ukama/systems/notification/mailer/mocks"
	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	"github.com/ukama/ukama/systems/notification/mailer/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mailId string

func TestMailer(t *testing.T) {
	// Test case 1: Add subscriber successfully
	mailerRepo := &mocks.MailerRepo{}
	m := NewMailerServer(mailerRepo, &pkg.Mailer{
		Host:     "sandbox.smtp.mailtrap.io",
		Port:     587,
		Username: "a7ee775f59cebc",
		Password: "939047730fb6ea",
		From:     "from@example.com",
	}, "../../templates")
	t.Run("Send email", func(t *testing.T) {

		mailerRepo.On("", mock.AnythingOfType("*db.Mailing")).Return(nil)
		mailerRepo.On("SendEmail", mock.Anything, mock.Anything).Return(nil)

		req := &pb.SendEmailRequest{
			To:           []string{" brackley@ukama.com"},
			TemplateName: "test-template",
			Values: map[string]string{
				"Name":    "John",
				"Message": "Hello World!",
			},
		}
		res, err := m.SendEmail(context.Background(), req)
		mailId = res.MailId

		assert.NoError(t, err)
	})

}