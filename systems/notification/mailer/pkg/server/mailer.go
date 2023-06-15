package server

import (
	"context"

	"net/smtp"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/common/config"
	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	"github.com/ukama/ukama/systems/notification/mailer/pkg/db"
)

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
	body := req.GetBody()

	// Compose the email message
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", s.mailer.Username, s.mailer.Password, s.mailer.Host)

	port := strconv.Itoa(s.mailer.Port)

	err := smtp.SendMail(s.mailer.Host+":"+port, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Errorf("Failed to send email: %v", err.Error())
		return nil, err
	}
	log.Infof("Email sent successfully sent to %s", to)
	
	response := &pb.SendEmailResponse{
		Message:   "Email sent successfully",
	}


	return response, nil

}
