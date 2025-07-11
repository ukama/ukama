/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/mailer/mocks"
	"github.com/ukama/ukama/systems/notification/mailer/pkg"
	"github.com/ukama/ukama/systems/notification/mailer/pkg/db"

	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
)

func setupServer(t *testing.T) (*MailerServer, *mocks.MailerRepo) {
	mockRepo := mocks.NewMailerRepo(t)
	mailer := &pkg.MailerConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "test@example.com",
		Password: "password",
		From:     "from@example.com",
	}

	server, err := NewMailerServer(mockRepo, mailer, "../../templates")
	require.NoError(t, err)

	t.Cleanup(func() {
		server.Stop()
	})

	return server, mockRepo
}

func TestGetEmailById(t *testing.T) {
	server, mockRepo := setupServer(t)
	testMailId := uuid.NewV4()
	testTime := time.Now()

	tests := []struct {
		name    string
		mailId  string
		setup   func(*mocks.MailerRepo)
		want    *pb.GetEmailByIdResponse
		wantErr bool
		errCode codes.Code
	}{
		{
			name:   "successful retrieval",
			mailId: testMailId.String(),
			setup: func(repo *mocks.MailerRepo) {
				repo.On("GetEmailById", testMailId).Return(&db.Mailing{
					MailId:       testMailId,
					Email:        "test@example.com",
					TemplateName: "test-template",
					Status:       ukama.MailStatusSuccess,
					SentAt:       &testTime,
				}, nil).Once()
			},
			want: &pb.GetEmailByIdResponse{
				MailId:       testMailId.String(),
				TemplateName: "test-template",
				Status:       pb.Status(pb.Status_value[ukama.MailStatusSuccess.String()]),
				SentAt:       testTime.String(),
			},
			wantErr: false,
		},
		{
			name:    "invalid UUID",
			mailId:  "invalid-uuid",
			setup:   func(repo *mocks.MailerRepo) {},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name:   "email not found",
			mailId: testMailId.String(),
			setup: func(repo *mocks.MailerRepo) {
				repo.On("GetEmailById", testMailId).Return(nil, errors.New("not found")).Once()
			},
			wantErr: true,
			errCode: codes.Internal,
		},
		{
			name:    "empty mail ID",
			mailId:  "",
			setup:   func(repo *mocks.MailerRepo) {},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockRepo)

			resp, err := server.GetEmailById(context.Background(), &pb.GetEmailByIdRequest{
				MailId: tt.mailId,
			})

			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.errCode, st.Code())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.want.MailId, resp.MailId)
				assert.Equal(t, tt.want.TemplateName, resp.TemplateName)
			}
		})
	}

	mockRepo.AssertExpectations(t)
}

func (s *MailerServer) Stop() {
	close(s.emailQueue)
	s.retryTicker.Stop()
}

func TestProcessEmailQueue(t *testing.T) {
	server, mockRepo := setupServer(t)
	mailId := uuid.NewV4()
	mockRepo.On("GetEmailById", mailId).Return(&db.Mailing{Status: ukama.MailStatusPending}, nil)
	mockRepo.On("UpdateEmailStatus", mock.Anything).Return(nil)

	go server.processEmailQueue()

	server.emailQueue <- &EmailPayload{
		To:           []string{"recipient@test.com"},
		TemplateName: "test-template",
		Values:       map[string]interface{}{"key": "value"},
		MailId:       mailId,
	}

	time.Sleep(1 * time.Second)

	mockRepo.AssertExpectations(t)
}
