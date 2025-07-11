/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"html/template"
	"net"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/mailer/pkg"
	"github.com/ukama/ukama/systems/notification/mailer/pkg/db"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
)

const (
	maxRetries         = 3
	retryDelay         = 5 * time.Second
	templateExtension  = ".tmpl"
	defaultTimeout     = 80 * time.Second
	dialTimeout        = 80 * time.Second
	smtpTimeout        = 60 * time.Second
	emailQueueCapacity = 100
	MaxRetryCount      = 3
	InitialBackoff     = 5 * time.Minute
)

type EmailPayload struct {
	To           []string               `json:"to"`
	TemplateName string                 `json:"template_name"`
	Values       map[string]interface{} `json:"values"`
	MailId       uuid.UUID
	Attachments  []struct {
		Filename    string
		ContentType string
		Content     []byte
	}
}

type MailerServer struct {
	mailerRepo db.MailerRepo
	pb.UnimplementedMailerServiceServer
	mailer        *pkg.MailerConfig
	templatesPath string
	templates     *template.Template
	emailQueue    chan *EmailPayload
	retryTicker   *time.Ticker
}

func NewMailerServer(mailerRepo db.MailerRepo, mail *pkg.MailerConfig, templatesPath string) (*MailerServer, error) {
	templates, err := template.ParseGlob(filepath.Join(templatesPath, "*"+templateExtension))
	if err != nil {
		return nil, fmt.Errorf("failed to load email templates: %w", err)
	}

	server := &MailerServer{
		mailerRepo:    mailerRepo,
		mailer:        mail,
		templatesPath: templatesPath,
		templates:     templates,
		emailQueue:    make(chan *EmailPayload, emailQueueCapacity),
		retryTicker:   time.NewTicker(1 * time.Minute),
	}

	go server.processEmailQueue()
	go server.processRetries()

	return server, nil
}

func (s *MailerServer) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	mailId := uuid.NewV4()
	log.Infof("Queueing email to %v", req.To)

	payload := &EmailPayload{
		To:           req.To,
		TemplateName: req.TemplateName,
		Values:       s.convertValues(req.Values),
		MailId:       mailId,
		Attachments: make([]struct {
			Filename    string
			ContentType string
			Content     []byte
		}, len(req.Attachments)),
	}

	// Convert attachments
	for i, att := range req.Attachments {
		payload.Attachments[i] = struct {
			Filename    string
			ContentType string
			Content     []byte
		}{
			Filename:    att.Filename,
			ContentType: att.ContentType,
			Content:     att.Content,
		}
	}

	if err := s.saveEmailStatus(mailId, strings.Join(req.To, ","), req.TemplateName, ukama.MailStatusPending, req.Values); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save email status: %v", err)
	}

	select {
	case s.emailQueue <- payload:
		return &pb.SendEmailResponse{
			Message: "Email queued for sending!",
			MailId:  mailId.String(),
		}, nil
	case <-timeoutCtx.Done():
		return nil, status.Errorf(codes.Canceled, "operation canceled: %v", timeoutCtx.Err())
	}
}

func (s *MailerServer) processEmailQueue() {
	for payload := range s.emailQueue {
		email, err := s.mailerRepo.GetEmailById(payload.MailId)
		if err != nil {
			log.WithError(err).Error("Failed to fetch email")
			continue
		}

		if email.Status == ukama.MailStatusSuccess {
			log.Warnf("Skipping email with mailId %v as it has already been sent", payload.MailId)
			continue
		}

		if err := s.mailerRepo.UpdateEmailStatus(&db.Mailing{
			MailId:        payload.MailId,
			Status:        ukama.MailStatusProcess,
			RetryCount:    0,
			NextRetryTime: nil,
			UpdatedAt:     time.Now(),
		}); err != nil {
			log.Errorf("Failed to update email status: %v", err)
			continue
		}
		now := time.Now()

		err = s.attemptSendEmail(payload)
		if err != nil {
			log.Errorf("Failed to send email: %v", err)

			nextRetry := time.Now().Add(InitialBackoff)
			if updateErr := s.mailerRepo.UpdateEmailStatus(&db.Mailing{
				MailId:        payload.MailId,
				Status:        ukama.MailStatusRetry,
				RetryCount:    1,
				NextRetryTime: &nextRetry,
				UpdatedAt:     time.Now(),
			}); updateErr != nil {
				log.Errorf("Failed to update email status: %v", updateErr)
			}
		} else {
			log.Infof("Email sent successfully to %v with template %s and mail ID %v",
				payload.To, payload.TemplateName, payload.MailId)

			if updateErr := s.mailerRepo.UpdateEmailStatus(&db.Mailing{
				MailId:        payload.MailId,
				Status:        ukama.MailStatusSuccess,
				SentAt:        &now,
				RetryCount:    0,
				NextRetryTime: nil,
				UpdatedAt:     time.Now(),
			}); updateErr != nil {
				log.Errorf("Failed to update email status: %v", updateErr)

			}
		}
	}
}

func (s *MailerServer) updateStatus(mailId uuid.UUID, status ukama.MailStatus) error {
	return s.mailerRepo.UpdateEmailStatus(&db.Mailing{
		MailId:    mailId,
		Status:    status,
		UpdatedAt: time.Now(),
	})
}

func (s *MailerServer) updateRetryStatus(mailId uuid.UUID, retryCount int, nextRetryTime *time.Time) error {
	return s.mailerRepo.UpdateEmailStatus(&db.Mailing{
		MailId:        mailId,
		Status:        ukama.MailStatusRetry,
		RetryCount:    retryCount,
		NextRetryTime: nextRetryTime,
		UpdatedAt:     time.Now(),
	})
}

func (s *MailerServer) processRetries() {
	for range s.retryTicker.C {
		emails, err := s.mailerRepo.GetFailedEmails()
		if err != nil {
			log.WithError(err).Error("Failed to fetch failed emails")
			continue
		}

		for _, email := range emails {
			if email.Status == ukama.MailStatusSuccess || email.Status == ukama.MailStatusProcess {
				continue
			}

			if email.NextRetryTime != nil && time.Now().Before(*email.NextRetryTime) {
				continue
			}

			if email.RetryCount >= MaxRetryCount {
				log.WithField("mailId", email.MailId).Info("Max retries reached, marking as permanently failed")

				if err := s.updateStatus(email.MailId, ukama.MailStatusFailed); err != nil {
					log.WithError(err).Error("Failed to update email status")
				}
				continue
			}

			if err := s.updateStatus(email.MailId, ukama.MailStatusProcess); err != nil {
				log.WithError(err).Error("Failed to update status to processing")
				continue
			}

			payload := &EmailPayload{
				To:           strings.Split(email.Email, ","),
				TemplateName: email.TemplateName,
				Values:       email.Values,
				MailId:       email.MailId,
			}

			if err := s.attemptSendEmail(payload); err != nil {
				log.WithError(err).WithField("mailId", email.MailId).Error("Retry attempt failed")
				nextRetry := time.Now().Add(InitialBackoff * time.Duration(1<<uint(email.RetryCount)))
				if err := s.updateRetryStatus(email.MailId, email.RetryCount+1, &nextRetry); err != nil {
					log.WithError(err).Error("Failed to update retry status")
				}
			} else {
				log.WithField("mailId", email.MailId).Info("Retry successful")
				if err := s.updateStatus(email.MailId, ukama.MailStatusSuccess); err != nil {
					log.WithError(err).Error("Failed to update email status")
				}
			}
		}
	}
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

func (s *MailerServer) saveEmailStatus(mailId uuid.UUID, email, templateName string, status ukama.MailStatus, values map[string]string) error {
	var nextRetryTime *time.Time
	if status == ukama.MailStatusFailed {
		t := time.Now().Add(InitialBackoff)
		nextRetryTime = &t
	}

	jsonMap := make(db.JSONMap)
	for k, v := range values {
		jsonMap[k] = v
	}

	return s.mailerRepo.CreateEmail(&db.Mailing{
		MailId:        mailId,
		Email:         email,
		TemplateName:  templateName,
		Status:        status,
		RetryCount:    0,
		NextRetryTime: nextRetryTime,
		Values:        jsonMap,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	})
}

func (s *MailerServer) attemptSendEmail(payload *EmailPayload) error {
	body, err := s.prepareMsg(payload)
	if err != nil {
		return fmt.Errorf("failed to prepare email body: %w", err)
	}

	client, err := s.createSMTPClient(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Warnf("failed to close smtp client connection: %v", err)
		}
	}()

	return s.sendWithClient(client, payload, body)
}

func (s *MailerServer) createSMTPClient(ctx context.Context) (*smtp.Client, error) {
	dialer := &net.Dialer{
		Timeout: dialTimeout,
	}

	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", s.mailer.Host, s.mailer.Port))
	if err != nil {
		return nil, fmt.Errorf("SMTP connection failed: %w", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Warnf("failed to close net client connection: %v", err)
		}
	}()

	client, err := smtp.NewClient(conn, s.mailer.Host)
	if err != nil {
		return nil, fmt.Errorf("SMTP client creation failed: %w", err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Warnf("failed to close smtp client connection: %V", err)
		}
	}()

	config := &tls.Config{
		ServerName:         s.mailer.Host,
		InsecureSkipVerify: false,
	}

	if err := client.StartTLS(config); err != nil {
		return nil, fmt.Errorf("TLS setup failed: %w", err)
	}

	auth := smtp.PlainAuth("", s.mailer.Username, s.mailer.Password, s.mailer.Host)
	if err := client.Auth(auth); err != nil {
		if strings.Contains(err.Error(), "authentication failed") {
			return nil, fmt.Errorf("SMTP authentication failed, check credentials: %w", err)
		}
		return nil, fmt.Errorf("SMTP authentication error: %w", err)
	}

	return client, nil
}

func (s *MailerServer) sendWithClient(client *smtp.Client, payload *EmailPayload, body bytes.Buffer) error {
	sendCtx, cancel := context.WithTimeout(context.Background(), smtpTimeout)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		if err := client.Mail(s.mailer.From); err != nil {
			errCh <- fmt.Errorf("failed to set sender: %w", err)
			return
		}

		for _, recipient := range payload.To {
			if err := client.Rcpt(recipient); err != nil {
				errCh <- fmt.Errorf("failed to add recipient %s: %w", recipient, err)
				return
			}
		}

		writer, err := client.Data()
		if err != nil {
			errCh <- fmt.Errorf("failed to create message writer: %w", err)
			return
		}

		defer func() {
			if err := writer.Close(); err != nil {
				log.Warnf("failed to close mail writer: %V", err)
			}
		}()

		if _, err := writer.Write(body.Bytes()); err != nil {
			errCh <- fmt.Errorf("failed to write message body: %w", err)
			return
		}

		if err := client.Quit(); err != nil {
			if !strings.Contains(err.Error(), "250 Ok") {
				errCh <- fmt.Errorf("failed to close SMTP connection: %w", err)
				return
			}
		}

		errCh <- nil
	}()

	select {
	case <-sendCtx.Done():
		return fmt.Errorf("send operation timed out: %w", sendCtx.Err())
	case err := <-errCh:
		return err
	}
}

func (s *MailerServer) prepareMsg(data *EmailPayload) (bytes.Buffer, error) {
	var body bytes.Buffer

	// First, execute the template to get the email content
	tmplName := data.TemplateName
	if filepath.Ext(tmplName) == "" {
		tmplName += templateExtension
	}

	t := s.templates.Lookup(tmplName)
	if t == nil {
		return body, fmt.Errorf("template %s not found", tmplName)
	}

	templateData := struct {
		Values map[string]interface{}
	}{
		Values: data.Values,
	}

	var templateBuffer bytes.Buffer
	if err := t.Execute(&templateBuffer, templateData); err != nil {
		log.WithError(err).Error("Template execution failed")
		return body, fmt.Errorf("failed to execute template: %w", err)
	}

	templateContent := templateBuffer.String()
	parts := strings.SplitN(templateContent, "\n\n", 2)

	headers := make(map[string]string)
	if len(parts) > 1 {
		headerLines := strings.Split(parts[0], "\n")
		for _, line := range headerLines {
			if colonIdx := strings.Index(line, ":"); colonIdx != -1 {
				key := strings.TrimSpace(line[:colonIdx])
				value := strings.TrimSpace(line[colonIdx+1:])
				headers[strings.ToLower(key)] = value
			}
		}
	}

	htmlContent := ""
	if len(parts) > 1 {
		htmlContent = parts[1]
	} else {
		htmlContent = parts[0]
	}

	fmt.Fprintf(&body, "From: %s\r\n", s.mailer.From)
	fmt.Fprintf(&body, "To: %s\r\n", strings.Join(data.To, ", "))

	subjectTemplate := template.Must(template.New("subject").Parse(headers["subject"]))
	var processedSubject bytes.Buffer
	if err := subjectTemplate.Execute(&processedSubject, templateData); err != nil {
		log.WithError(err).Error("Subject template execution failed")
		processedSubject.WriteString("No Subject")
	}

	fmt.Fprintf(&body, "Subject: %s\r\n", processedSubject.String())
	fmt.Fprintf(&body, "MIME-Version: 1.0\r\n")

	boundary := "UkamaMailBoundary" + uuid.NewV4().String()

	if len(data.Attachments) > 0 {
		fmt.Fprintf(&body, "Content-Type: multipart/mixed; boundary=%s\r\n\r\n", boundary)
		fmt.Fprintf(&body, "--%s\r\n", boundary)

		fmt.Fprintf(&body, "Content-Type: text/html; charset=UTF-8\r\n")
		fmt.Fprintf(&body, "Content-Transfer-Encoding: 7bit\r\n\r\n")
		fmt.Fprintf(&body, "%s\r\n", htmlContent)

		for _, att := range data.Attachments {
			fmt.Fprintf(&body, "\r\n--%s\r\n", boundary)
			fmt.Fprintf(&body, "Content-Type: %s; name=\"%s\"\r\n", att.ContentType, att.Filename)
			fmt.Fprintf(&body, "Content-Disposition: attachment; filename=\"%s\"\r\n", att.Filename)
			fmt.Fprintf(&body, "Content-Transfer-Encoding: base64\r\n\r\n")

			content := att.Content
			encoder := base64.StdEncoding
			encoded := make([]byte, encoder.EncodedLen(len(content)))
			encoder.Encode(encoded, content)

			lineLength := 76
			for i := 0; i < len(encoded); i += lineLength {
				end := i + lineLength
				if end > len(encoded) {
					end = len(encoded)
				}
				fmt.Fprintf(&body, "%s\r\n", encoded[i:end])
			}
		}

		fmt.Fprintf(&body, "\r\n--%s--\r\n", boundary)
	} else {
		fmt.Fprintf(&body, "Content-Type: text/html; charset=UTF-8\r\n")
		fmt.Fprintf(&body, "Content-Transfer-Encoding: 7bit\r\n\r\n")
		fmt.Fprintf(&body, "%s\r\n", htmlContent)
	}

	if pkg.IsDebugMode {
		log.WithFields(log.Fields{
			"finalBody": body.String(),
			"values":    data.Values,
		}).Debug("Email body prepared")
	}

	return body, nil
}

func (s *MailerServer) GetEmailById(ctx context.Context, req *pb.GetEmailByIdRequest) (*pb.GetEmailByIdResponse, error) {
	log.Infof("GetEmailById: %v", req)
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

	return &pb.GetEmailByIdResponse{
		MailId:       mail.MailId.String(),
		TemplateName: mail.TemplateName,
		SentAt:       mail.CreatedAt.String(),
		Status:       pb.Status(pb.Status_value[ukama.MailStatus(mail.Status).String()]),
	}, nil
}
