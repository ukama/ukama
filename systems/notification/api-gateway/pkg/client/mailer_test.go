/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/notification/mailer/pb/gen/mocks"

	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
)

var mc = &mocks.MailerServiceClient{}
var mailId = uuid.NewV4().String()

func TestMailer_SendEmail(t *testing.T) {
	emailReq := &pb.SendEmailRequest{
		To:           []string{"test@example.com", "user@example.com"},
		TemplateName: "welcome_template",
		Values: map[string]string{
			"name":    "John Doe",
			"company": "Ukama Inc",
		},
		Status: pb.Status_Pending,
		Attachments: []*pb.Attachment{
			{
				Filename:    "document.pdf",
				ContentType: "application/pdf",
				Content:     []byte("fake pdf content"),
			},
		},
	}

	emailResp := &pb.SendEmailResponse{
		Message: "Email sent successfully",
		MailId:  mailId,
	}

	mc.On("SendEmail", mock.Anything, emailReq).Return(emailResp, nil)

	c := client.NewMailerFromClient(mc)

	resp, err := c.SendEmail(emailReq)
	assert.NoError(t, err)

	if assert.NotNil(t, resp) {
		assert.Equal(t, emailResp.Message, resp.Message)
		assert.Equal(t, emailResp.MailId, resp.MailId)
	}

	mc.AssertExpectations(t)
}

func TestMailer_SendEmail_WithMinimalData(t *testing.T) {
	emailReq := &pb.SendEmailRequest{
		To:           []string{"simple@example.com"},
		TemplateName: "simple_template",
		Values:       map[string]string{},
		Status:       pb.Status_Pending,
	}

	emailResp := &pb.SendEmailResponse{
		Message: "Email queued for sending",
		MailId:  mailId,
	}

	mc.On("SendEmail", mock.Anything, emailReq).Return(emailResp, nil)

	c := client.NewMailerFromClient(mc)

	resp, err := c.SendEmail(emailReq)
	assert.NoError(t, err)

	if assert.NotNil(t, resp) {
		assert.Equal(t, emailResp.Message, resp.Message)
		assert.Equal(t, emailResp.MailId, resp.MailId)
	}

	mc.AssertExpectations(t)
}

func TestMailer_GetEmailById(t *testing.T) {
	emailReq := &pb.GetEmailByIdRequest{
		MailId: mailId,
	}

	emailResp := &pb.GetEmailByIdResponse{
		MailId:       mailId,
		To:           "test@example.com",
		TemplateName: "welcome_template",
		Status:       pb.Status_Success,
		SentAt:       "2023-12-01T10:00:00Z",
		Values: map[string]string{
			"name":    "John Doe",
			"company": "Ukama Inc",
		},
	}

	mc.On("GetEmailById", mock.Anything, emailReq).Return(emailResp, nil)

	c := client.NewMailerFromClient(mc)

	resp, err := c.GetEmailById(mailId)
	assert.NoError(t, err)

	if assert.NotNil(t, resp) {
		assert.Equal(t, emailResp.MailId, resp.MailId)
		assert.Equal(t, emailResp.To, resp.To)
		assert.Equal(t, emailResp.TemplateName, resp.TemplateName)
		assert.Equal(t, emailResp.Status, resp.Status)
		assert.Equal(t, emailResp.SentAt, resp.SentAt)
		assert.Equal(t, emailResp.Values, resp.Values)
	}

	mc.AssertExpectations(t)
}

func TestMailer_GetEmailById_WithFailedStatus(t *testing.T) {
	emailReq := &pb.GetEmailByIdRequest{
		MailId: mailId,
	}

	emailResp := &pb.GetEmailByIdResponse{
		MailId:       mailId,
		To:           "failed@example.com",
		TemplateName: "error_template",
		Status:       pb.Status_Failed,
		SentAt:       "",
		Values:       map[string]string{},
	}

	mc.On("GetEmailById", mock.Anything, emailReq).Return(emailResp, nil)

	c := client.NewMailerFromClient(mc)

	resp, err := c.GetEmailById(mailId)
	assert.NoError(t, err)

	if assert.NotNil(t, resp) {
		assert.Equal(t, emailResp.MailId, resp.MailId)
		assert.Equal(t, emailResp.Status, pb.Status_Failed)
		assert.Empty(t, emailResp.SentAt)
	}

	mc.AssertExpectations(t)
}

func TestMailer_GetEmailById_WithRetryStatus(t *testing.T) {
	emailReq := &pb.GetEmailByIdRequest{
		MailId: mailId,
	}

	emailResp := &pb.GetEmailByIdResponse{
		MailId:       mailId,
		To:           "retry@example.com",
		TemplateName: "retry_template",
		Status:       pb.Status_Retry,
		SentAt:       "",
		Values:       map[string]string{},
	}

	mc.On("GetEmailById", mock.Anything, emailReq).Return(emailResp, nil)

	c := client.NewMailerFromClient(mc)

	resp, err := c.GetEmailById(mailId)
	assert.NoError(t, err)

	if assert.NotNil(t, resp) {
		assert.Equal(t, emailResp.MailId, resp.MailId)
		assert.Equal(t, emailResp.Status, pb.Status_Retry)
	}

	mc.AssertExpectations(t)
}
