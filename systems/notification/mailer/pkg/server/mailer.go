package server

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"net"
	"net/smtp"
	"path/filepath"
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

const (
	maxRetries        = 3
	retryDelay        = 5 * time.Second
	templateExtension = ".tmpl"
	defaultTimeout    = 80 * time.Second
	dialTimeout       = 80 * time.Second
	smtpTimeout       = 60 * time.Second
)

type EmailPayload struct {
	To           []string               `json:"to"`
	TemplateName string                 `json:"template_name"`
	Values       map[string]interface{} `json:"values"`
}

type MailerServer struct {
	mailerRepo db.MailerRepo
	pb.UnimplementedMailerServiceServer
	mailer        *pkg.Mailer
	templatesPath string
	templates     *template.Template
}

func NewMailerServer(mailerRepo db.MailerRepo, mail *pkg.Mailer, templatesPath string) (*MailerServer, error) {
	templates, err := template.ParseGlob(filepath.Join(templatesPath, "*"+templateExtension))
	if err != nil {
		return nil, fmt.Errorf("failed to load email templates: %v", err)
	}

	return &MailerServer{
		mailerRepo:    mailerRepo,
		mailer:        mail,
		templatesPath: templatesPath,
		templates:     templates,
	}, nil
}

func (s *MailerServer) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	mailId := uuid.NewV4()
	logger := log.WithField("mail_id", mailId)
	logger.Infof("Sending email to %v", req.To)

	payload := &EmailPayload{
		To:           req.To,
		TemplateName: req.TemplateName,
		Values:       s.convertValues(req.Values),
	}

	var lastError error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := s.attemptSendEmail(ctx, payload)
		if err == nil {
			if err := s.saveEmailStatus(mailId, req.To[0], req.TemplateName, "sent"); err != nil {
				logger.WithError(err).Error("Failed to save successful email status")
			}
			return &pb.SendEmailResponse{
				Message: "Email Sent!",
				MailId:  mailId.String(),
			}, nil
		}

		lastError = err
		logger.WithFields(log.Fields{
			"attempt": attempt,
			"error":   err,
		}).Error("Email sending failed")

		if attempt < maxRetries {
			backoff := retryDelay * time.Duration(attempt)
			logger.WithField("backoff", backoff).Info("Retrying after backoff")

			timer := time.NewTimer(backoff)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, status.Errorf(codes.Canceled, "operation cancelled: %v", ctx.Err())
			case <-timer.C:
				continue
			}
		}
	}

	if err := s.saveEmailStatus(mailId, req.To[0], req.TemplateName, "failed"); err != nil {
		logger.WithError(err).Error("Failed to save failed email status")
	}

	return nil, status.Errorf(codes.Internal, "failed to send email after %d attempts: %v", maxRetries, lastError)
}

func (s *MailerServer) validateRequest(req *pb.SendEmailRequest) error {
	if len(req.To) == 0 {
		return status.Error(codes.InvalidArgument, "recipient email address required")
	}
	for _, email := range req.To {
		if !isValidEmail(email) {
			return status.Errorf(codes.InvalidArgument, "invalid email address: %s", email)
		}
	}
	if req.TemplateName == "" {
		return status.Error(codes.InvalidArgument, "template name required")
	}
	return nil
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func (s *MailerServer) convertValues(reqValues map[string]string) map[string]interface{} {
	values := make(map[string]interface{}, len(reqValues))
	for key, value := range reqValues {
		values[key] = value
	}
	return values
}

func (s *MailerServer) attemptSendEmail(ctx context.Context, payload *EmailPayload) error {
	body, err := s.prepareMsg(payload)
	if err != nil {
		return fmt.Errorf("failed to prepare email body: %v", err)
	}

	client, err := s.createSMTPClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Close()

	return s.sendWithClient(ctx, client, payload, body)
}

func (s *MailerServer) createSMTPClient(ctx context.Context) (*smtp.Client, error) {
	dialer := &net.Dialer{
		Timeout: dialTimeout,
	}

	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", s.mailer.Host, s.mailer.Port))
	if err != nil {
		return nil, fmt.Errorf("SMTP connection failed: %v", err)
	}

	client, err := smtp.NewClient(conn, s.mailer.Host)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("SMTP client creation failed: %v", err)
	}

	config := &tls.Config{
		ServerName:         s.mailer.Host,
		InsecureSkipVerify: false,
	}
	if err := client.StartTLS(config); err != nil {
		client.Close()
		return nil, fmt.Errorf("TLS setup failed: %v", err)
	}

	auth := smtp.PlainAuth("", s.mailer.Username, s.mailer.Password, s.mailer.Host)
	if err := client.Auth(auth); err != nil {
		client.Close()
		if strings.Contains(err.Error(), "authentication failed") {
			return nil, fmt.Errorf("SMTP authentication failed, check credentials: %v", err)
		}
		return nil, fmt.Errorf("SMTP authentication error: %v", err)
	}

	return client, nil
}

func (s *MailerServer) sendWithClient(ctx context.Context, client *smtp.Client, payload *EmailPayload, body bytes.Buffer) error {
	sendCtx, cancel := context.WithTimeout(ctx, smtpTimeout)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		if err := client.Mail(s.mailer.From); err != nil {
			errCh <- fmt.Errorf("failed to set sender: %v", err)
			return
		}

		for _, recipient := range payload.To {
			if err := client.Rcpt(recipient); err != nil {
				errCh <- fmt.Errorf("failed to add recipient %s: %v", recipient, err)
				return
			}
		}

		writer, err := client.Data()
		if err != nil {
			errCh <- fmt.Errorf("failed to create message writer: %v", err)
			return
		}
		defer writer.Close()

		if _, err := writer.Write(body.Bytes()); err != nil {
			errCh <- fmt.Errorf("failed to write message body: %v", err)
			return
		}

		if err := client.Quit(); err != nil {
			errCh <- fmt.Errorf("failed to close SMTP connection: %v", err)
			return
		}

		errCh <- nil
	}()

	select {
	case <-sendCtx.Done():
		return fmt.Errorf("send operation timed out: %v", sendCtx.Err())
	case err := <-errCh:
		return err
	}
}

func (s *MailerServer) prepareMsg(data *EmailPayload) (bytes.Buffer, error) {
	var body bytes.Buffer
	tmplName := data.TemplateName
	if filepath.Ext(tmplName) == "" {
		tmplName += templateExtension
	}

	t := s.templates.Lookup(tmplName)
	if t == nil {
		return body, fmt.Errorf("template %s not found", tmplName)
	}

	if err := t.Execute(&body, data); err != nil {
		return body, fmt.Errorf("failed to execute template: %v", err)
	}

	if pkg.IsDebugMode {
		log.WithField("body", body.String()).Debug("Email body prepared")
	}

	return body, nil
}

func (s *MailerServer) saveEmailStatus(mailId uuid.UUID, email, templateName, status string) error {
	return s.mailerRepo.SendEmail(&db.Mailing{
		MailId:       mailId,
		Email:        email,
		TemplateName: templateName,
		Status:       status,
	})
}

func (s *MailerServer) GetEmailById(ctx context.Context, req *pb.GetEmailByIdRequest) (*pb.GetEmailByIdResponse, error) {
	if req.MailId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing mail ID")
	}

	mailerId, err := uuid.FromString(req.GetMailId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid mail ID")
	}

	mail, err := s.mailerRepo.GetEmailById(mailerId)
	if err != nil {
		log.WithError(err).Error("Error while getting email")
		return nil, grpc.SqlErrorToGrpc(err, "failed to get email")
	}

	log.WithField("mail_id", mail.MailId).Info("Retrieved email")

	return &pb.GetEmailByIdResponse{
		MailId:       mail.MailId.String(),
		TemplateName: mail.TemplateName,
		SentAt:       mail.SentAt.String(),
		Status:       mail.Status,
	}, nil
}
