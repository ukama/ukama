package server

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"

	"github.com/ukama/ukama/systems/notification/mailer/pkg"
	"github.com/ukama/ukama/systems/notification/mailer/pkg/db"
)

type EmailPayload struct {
	To      []string
	TemplateName string `json:"template_name"`
	Values  map[string]interface{}
}


type MailerServer struct {
	mailerRepoRepo db.MailerRepo
	pb.UnimplementedMailerServiceServer
	mailer *pkg.Mailer
	templatesPath string
}

func NewMailerServer(mailerRepoRepo db.MailerRepo, mail *pkg.Mailer , templatesPath string) *MailerServer {
	return &MailerServer{
		mailerRepoRepo: mailerRepoRepo,
		mailer:         mail,
		templatesPath:  templatesPath,
		
	}
}
func (s *MailerServer) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {
	log.Infof("Sending email to %v", req.To)
	values := make(map[string]interface{})

	for key, value := range req.Values {
		values[key] = value
	}

	payload := &EmailPayload{
		To:           req.To,
		TemplateName: req.TemplateName,
		Values:       values,
	}

	body, err := s.prepareMsg(payload)
	if err != nil {
		return nil, err
	}

	
	c, err := smtp.Dial(fmt.Sprintf("%s:%d", s.mailer.Host, s.mailer.Port))
	if err != nil {
		return nil, err
	}
	defer c.Close()
	config := &tls.Config{
		ServerName: s.mailer.Host,
	}
	if err = c.StartTLS(config); err != nil {
		return nil, err
	}

	auth := smtp.PlainAuth("", s.mailer.Username, s.mailer.Password, s.mailer.Host)

	// Authenticate with the SMTP server
	if err = c.Auth(auth); err != nil {
		return nil, err
	}

	// Set the sender email address
	if err = c.Mail(s.mailer.From); err != nil {
		return nil, err
	}

	// Add the recipient email addresses
	for _, recipient := range req.To {
		if err = c.Rcpt(recipient); err != nil {
			return nil, err
		}
	}

	w, err := c.Data()
	if err != nil {
		return nil, err
	}
	defer w.Close()

	_, err = w.Write(body.Bytes())
	if err != nil {
		return nil, err
	}

	err = c.Quit()
	if err != nil {
		log.Errorln("Email not sent ", err.Error())
	}

	mailId := uuid.NewV4()
	return &pb.SendEmailResponse{
		Message: "Email Sent!",
		MailId:  mailId.String(),
	}, nil
}

func (s *MailerServer) prepareMsg(data *EmailPayload) (bytes.Buffer, error) {
	tmplName := data.TemplateName
	if filepath.Ext(tmplName) == "" {
		tmplName += ".tmpl"
	}

	t, err := template.ParseFiles(filepath.Join(s.templatesPath, tmplName))
	if err != nil {
		return bytes.Buffer{}, err
	}

	var body bytes.Buffer

	err = t.Execute(&body, data)
	if err != nil {
		return bytes.Buffer{}, err
	}

	if pkg.IsDebugMode {
		log.Printf("%s", body.String())
	}
	return body, nil
}

func (s *MailerServer) GetEmailById(ctx context.Context, req *pb.GetEmailByIdRequest) (*pb.GetEmailByIdResponse, error) {
	if req.MailId == "" {
		return nil, errors.New("missing required fields in GetEmailByIdRequest")
	}
	mailerId, err := uuid.FromString(req.GetMailId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of mailer uuid. Error %s", err.Error())
	}
	mail, err := s.mailerRepoRepo.GetEmailById(mailerId)
	if err != nil {
		log.Error("Error while getting email" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "Failed to get email")
	}
	log.Infof("getting email with id %v", mail.MailId.String())
	response := &pb.GetEmailByIdResponse{
		MailId:  mail.MailId.String(),
		Subject: mail.Subject,
		Body:    mail.Body,
		SentAt:  mail.SentAt.String(),
		Status:  mail.Status,
	}

	return response, nil
}
