package server

import (
	"bytes"
	"context"
	"html/template"
	"net/smtp"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/common/config"
	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	"github.com/ukama/ukama/systems/notification/mailer/pkg/db"
)

type EmailData struct {
	To      string
	Subject string
	Body    string
	Values  map[string]interface{}
}

type MaillingServer struct {
	maillingRepoRepo db.MaillingRepo
	pb.UnimplementedMaillingServiceServer
	mailer *config.Mailer
}

func NewMaillingServer(maillingRepoRepo db.MaillingRepo, mail *config.Mailer) *MaillingServer {
	return &MaillingServer{
		maillingRepoRepo: maillingRepoRepo,
		mailer:           mail,
	}
}

func (s *MaillingServer) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {
	from := s.mailer.From
	to := req.GetTo()
	subject := req.GetSubject()
	bodyTemplate := req.GetBody()
	values := req.Values

	// Create a new email data struct
	emailData := &EmailData{
		To:      to,
		Subject: subject,
		Body:    bodyTemplate,
		Values:  make(map[string]interface{}),
	}

	// Add values to the email data
	for key, value := range values {
		emailData.Values[key] = value
	}

	tmpl, err := template.New("email").Parse(emailData.Body)
	if err != nil {
		log.Errorf("Failed to parse email template: %v", err)
		return nil, err
	}

	var bodyBuffer bytes.Buffer
	err = tmpl.Execute(&bodyBuffer, emailData.Values)
	if err != nil {
		log.Errorf("Failed to render email template: %v", err)
		return nil, err
	}

	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" + 
		"\r\n" +
		bodyBuffer.String()

	auth := smtp.PlainAuth("", s.mailer.Username, s.mailer.Password, s.mailer.Host)
	port := strconv.Itoa(s.mailer.Port)
	err = smtp.SendMail(s.mailer.Host+":"+port, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Errorf("Failed to send email: %v", err.Error())
		return nil, err
	}

	log.Infof("Email sent successfully to %s", to)

	response := &pb.SendEmailResponse{
		Message: "Email sent successfully",
	}

	return response, nil
}
