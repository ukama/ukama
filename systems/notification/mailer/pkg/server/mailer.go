package server

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"

	"github.com/ukama/ukama/systems/notification/mailer/pkg"
	"github.com/ukama/ukama/systems/notification/mailer/pkg/db"
)

type EmailData struct {
	To      []string
	Subject string
	Body    string
	Values  map[string]interface{}
}

type MailerServer struct {
	mailerRepoRepo db.MailerRepo
	pb.UnimplementedMailerServiceServer
	mailer *pkg.Mailer
}

func NewMailerServer(mailerRepoRepo db.MailerRepo, mail *pkg.Mailer) *MailerServer {
	return &MailerServer{
		mailerRepoRepo: mailerRepoRepo,
		mailer:         mail,
	}
}

func (s *MailerServer) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {

	if req.To == nil || req.Subject == "" || req.Body == "" {
		return nil, errors.New("missing required fields in SendEmailRequest")
	}
	mailId := uuid.NewV4()
	currentTime := time.Now()
	sentAt := &currentTime

	from := s.mailer.From
	to := req.GetTo()
	subject := req.GetSubject()
	bodyTemplate := req.GetBody()
	values := req.GetValues()
	emailData := &EmailData{
		Subject: subject,
		Body:    bodyTemplate,
		Values:  make(map[string]interface{}),
	}

	for key, value := range values {
		emailData.Values[key] = value
	}
	emailData.Values["EmailID"] = mailId.String()
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

	var recipientList []string
	if len(to) > 0 {
		recipientList = make([]string, len(to))
		copy(recipientList, to)
	}

	msg := "From: " + from + "\r\n" +
		"To: " + strings.Join(recipientList, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" +
		"\r\n" +
		bodyBuffer.String()

	auth := smtp.PlainAuth("", s.mailer.Username, s.mailer.Password, s.mailer.Host)
	port := strconv.Itoa(s.mailer.Port)
	err = smtp.SendMail(s.mailer.Host+":"+port, auth, from, recipientList, []byte(msg))
	if err != nil {
		log.Errorf("Failed to send email: %v", err.Error())
		err = s.mailerRepoRepo.SendEmail(&db.Mailing{
			MailId:  mailId,
			Email:   recipientList[0], // Use the first email if only one is provided
			Subject: subject,
			Body:    bodyBuffer.String(),
			SentAt:  sentAt,
			Status:  "failed",
		})
		return nil, err
	}

	log.Infof("Email sent successfully to %v", recipientList)

	for _, recipient := range recipientList {
		err = s.mailerRepoRepo.SendEmail(&db.Mailing{
			MailId:  mailId,
			Email:   recipient,
			Subject: subject,
			Body:    bodyBuffer.String(),
			SentAt:  sentAt,
			Status:  "sent",
		})
		if err != nil {
			log.Error("Error while sending email" + err.Error())
			return nil, grpc.SqlErrorToGrpc(err, "Failed to send email")
		}
	}

	response := &pb.SendEmailResponse{
		Message: "Email sent successfully",
		MailId:  mailId.String(),
	}

	return response, nil
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
