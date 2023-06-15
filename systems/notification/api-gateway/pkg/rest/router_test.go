package rest_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/mailer/api-gateway/pkg/rest"
)

type MockMailerManager struct {
	SendEmailFunc      func(to, message, subject string) (rest.SendEmailRes, error)
	SendEmailCalled    bool
	SendEmailTo        string
	SendEmailMessage   string
	SendEmailSubject   string
}

func (m *MockMailerManager) SendEmail(to, message, subject string) (rest.SendEmailRes, error) {
	m.SendEmailCalled = true
	m.SendEmailTo = to
	m.SendEmailMessage = message
	m.SendEmailSubject = subject
	if m.SendEmailFunc != nil {
		return m.SendEmailFunc(to, message, subject)
	}
	return rest.SendEmailRes{}, nil
}

func TestSendEmail(t *testing.T) {
	// Create a new instance of the Router
	router := rest.NewRouter(nil, nil)

	// Create a Gin context and request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a mock MailerManager
	mockMailerManager := &MockMailerManager{
		SendEmailFunc: func(to, message, subject string) (rest.SendEmailRes, error) {
			// Define the expected behavior of the mock MailerManager
			if to == "test@example.com" && message == "Hello, World!" && subject == "Test Subject" {
				return rest.SendEmailRes{
					Message: "Email sent successfully",
				}, nil
			}
			return rest.SendEmailRes{}, errors.New("Failed to send email")
		},
	}

	// Assign the mock MailerManager to the Router
	router.Client = &rest.Clients{
		Ma: mockMailerManager,
	}

	// Create a SendEmailReq
	req := &rest.SendEmailReq{
		To:      "test@example.com",
		Message: "Hello, World!",
		Subject: "Test Subject",
	}

	// Call the sendEmail handler
	res, err := router.SendEmail(c, req)

	// Assert the response and error
	assert.NoError(t, err)
	assert.Equal(t, &rest.SendEmailRes{Message: "Email sent successfully"}, res)

	// Assert that the SendEmail function of the mock MailerManager was called
	assert.True(t, mockMailerManager.SendEmailCalled)
	assert.Equal(t, "test@example.com", mockMailerManager.SendEmailTo)
	assert.Equal(t, "Hello, World!", mockMailerManager.SendEmailMessage)
	assert.Equal(t, "Test Subject", mockMailerManager.SendEmailSubject)
}
