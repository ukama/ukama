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

	t.Run("successful email processing", func(t *testing.T) {
		mailId := uuid.NewV4()
		payload := &EmailPayload{
			To:           []string{"test@example.com"},
			TemplateName: "test-template",
			Values:       map[string]interface{}{"Name": "Test User"},
			MailId:       mailId,
		}

		mockRepo.On("GetEmailById", mailId).Return(&db.Mailing{
			Status: ukama.Pending,
		}, nil).Once()
		mockRepo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
			return m.MailId == mailId && m.Status == ukama.Process
		})).Return(nil).Once()
		mockRepo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
			return m.MailId == mailId && m.Status == ukama.Retry && m.RetryCount == 1
		})).Return(nil).Once()

		server.emailQueue <- payload

		time.Sleep(300 * time.Millisecond)

		mockRepo.AssertExpectations(t)
	})

	t.Run("email already sent - skip processing", func(t *testing.T) {
		mailId := uuid.NewV4()
		payload := &EmailPayload{
			To:           []string{"test@example.com"},
			TemplateName: "test-template",
			Values:       map[string]interface{}{"Name": "Test User"},
			MailId:       mailId,
		}

		mockRepo.On("GetEmailById", mailId).Return(&db.Mailing{
			Status: ukama.Success,
		}, nil).Once()

		server.emailQueue <- payload

		time.Sleep(200 * time.Millisecond)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failed to fetch email - skip processing", func(t *testing.T) {
		mailId := uuid.NewV4()
		payload := &EmailPayload{
			To:           []string{"test@example.com"},
			TemplateName: "test-template",
			Values:       map[string]interface{}{"Name": "Test User"},
			MailId:       mailId,
		}

		mockRepo.On("GetEmailById", mailId).Return(nil, errors.New("database error")).Once()

		server.emailQueue <- payload

		time.Sleep(200 * time.Millisecond)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failed to update status to process - skip processing", func(t *testing.T) {
		mailId := uuid.NewV4()
		payload := &EmailPayload{
			To:           []string{"test@example.com"},
			TemplateName: "test-template",
			Values:       map[string]interface{}{"Name": "Test User"},
			MailId:       mailId,
		}

		mockRepo.On("GetEmailById", mailId).Return(&db.Mailing{
			Status: ukama.Pending,
		}, nil).Once()
		mockRepo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
			return m.MailId == mailId && m.Status == ukama.Process
		})).Return(errors.New("update failed")).Once()

		server.emailQueue <- payload

		time.Sleep(200 * time.Millisecond)

		mockRepo.AssertExpectations(t)
	})
}

func TestSendEmail(t *testing.T) {
	server, mockRepo := setupServer(t)

	tests := []struct {
		name        string
		request     *pb.SendEmailRequest
		setup       func(*mocks.MailerRepo)
		wantMessage string
		wantErr     bool
		errCode     codes.Code
	}{
		{
			name: "successful email queueing with single recipient",
			request: &pb.SendEmailRequest{
				To:           []string{"test@example.com"},
				TemplateName: "test-template",
				Values: map[string]string{
					"Name":    "John Doe",
					"Message": "Welcome to Ukama!",
				},
			},
			setup: func(repo *mocks.MailerRepo) {
				repo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).Return(nil).Once()
				// Setup expectation for background processing
				repo.On("GetEmailById", mock.AnythingOfType("uuid.UUID")).Return(&db.Mailing{
					Status: ukama.Pending,
				}, nil).Maybe()
				repo.On("UpdateEmailStatus", mock.AnythingOfType("*db.Mailing")).Return(nil).Maybe()
			},
			wantMessage: "Email queued for sending!",
			wantErr:     false,
		},
		{
			name: "successful email queueing with multiple recipients",
			request: &pb.SendEmailRequest{
				To:           []string{"user1@example.com", "user2@example.com", "user3@example.com"},
				TemplateName: "test-template",
				Values: map[string]string{
					"Name":    "Team",
					"Message": "Meeting reminder",
				},
			},
			setup: func(repo *mocks.MailerRepo) {
				repo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).Return(nil).Once()
				repo.On("GetEmailById", mock.AnythingOfType("uuid.UUID")).Return(&db.Mailing{
					Status: ukama.Pending,
				}, nil).Maybe()
				repo.On("UpdateEmailStatus", mock.AnythingOfType("*db.Mailing")).Return(nil).Maybe()
			},
			wantMessage: "Email queued for sending!",
			wantErr:     false,
		},
		{
			name: "successful email queueing with attachments",
			request: &pb.SendEmailRequest{
				To:           []string{"user@example.com"},
				TemplateName: "test-template",
				Values: map[string]string{
					"Name":    "User",
					"Message": "Please find attached",
				},
				Attachments: []*pb.Attachment{
					{
						Filename:    "document.pdf",
						ContentType: "application/pdf",
						Content:     []byte("fake pdf content"),
					},
					{
						Filename:    "image.jpg",
						ContentType: "image/jpeg",
						Content:     []byte("fake image content"),
					},
				},
			},
			setup: func(repo *mocks.MailerRepo) {
				repo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).Return(nil).Once()
				// Setup expectation for background processing
				repo.On("GetEmailById", mock.AnythingOfType("uuid.UUID")).Return(&db.Mailing{
					Status: ukama.Pending,
				}, nil).Maybe()
				repo.On("UpdateEmailStatus", mock.AnythingOfType("*db.Mailing")).Return(nil).Maybe()
			},
			wantMessage: "Email queued for sending!",
			wantErr:     false,
		},
		{
			name: "successful email queueing with complex values",
			request: &pb.SendEmailRequest{
				To:           []string{"admin@example.com"},
				TemplateName: "test-template",
				Values: map[string]string{
					"Name":    "Administrator",
					"Message": "System notification",
					"Code":    "ABC123",
					"Date":    "2024-01-15",
					"Status":  "Active",
				},
			},
			setup: func(repo *mocks.MailerRepo) {
				repo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).Return(nil).Once()

				repo.On("GetEmailById", mock.AnythingOfType("uuid.UUID")).Return(&db.Mailing{
					Status: ukama.Pending,
				}, nil).Maybe()
				repo.On("UpdateEmailStatus", mock.AnythingOfType("*db.Mailing")).Return(nil).Maybe()
			},
			wantMessage: "Email queued for sending!",
			wantErr:     false,
		},
		{
			name: "successful email queueing with empty values",
			request: &pb.SendEmailRequest{
				To:           []string{"user@example.com"},
				TemplateName: "test-template",
				Values:       map[string]string{},
			},
			setup: func(repo *mocks.MailerRepo) {
				repo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).Return(nil).Once()
				repo.On("GetEmailById", mock.AnythingOfType("uuid.UUID")).Return(&db.Mailing{
					Status: ukama.Pending,
				}, nil).Maybe()
				repo.On("UpdateEmailStatus", mock.AnythingOfType("*db.Mailing")).Return(nil).Maybe()
			},
			wantMessage: "Email queued for sending!",
			wantErr:     false,
		},
		{
			name: "successful email queueing with special characters in values",
			request: &pb.SendEmailRequest{
				To:           []string{"user@example.com"},
				TemplateName: "test-template",
				Values: map[string]string{
					"Name":    "José María",
					"Message": "¡Hola! ¿Cómo estás?",
					"Symbols": "!@#$%^&*()_+-=[]{}|;':\",./<>?",
				},
			},
			setup: func(repo *mocks.MailerRepo) {
				repo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).Return(nil).Once()
				// Setup expectation for background processing
				repo.On("GetEmailById", mock.AnythingOfType("uuid.UUID")).Return(&db.Mailing{
					Status: ukama.Pending,
				}, nil).Maybe()
				repo.On("UpdateEmailStatus", mock.AnythingOfType("*db.Mailing")).Return(nil).Maybe()
			},
			wantMessage: "Email queued for sending!",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil

			tt.setup(mockRepo)

			resp, err := server.SendEmail(context.Background(), tt.request)

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.errCode, st.Code())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.wantMessage, resp.Message)
				assert.NotEmpty(t, resp.MailId)

				_, uuidErr := uuid.FromString(resp.MailId)
				assert.NoError(t, uuidErr)

				time.Sleep(100 * time.Millisecond)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSendEmail_ValidationErrors(t *testing.T) {
	server, mockRepo := setupServer(t)

	tests := []struct {
		name    string
		request *pb.SendEmailRequest
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "empty recipients list",
			request: &pb.SendEmailRequest{
				To:           []string{},
				TemplateName: "test-template",
				Values:       map[string]string{"Name": "Test"},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "invalid email address",
			request: &pb.SendEmailRequest{
				To:           []string{"invalid-email"},
				TemplateName: "test-template",
				Values:       map[string]string{"Name": "Test"},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "missing template name",
			request: &pb.SendEmailRequest{
				To:     []string{"test@example.com"},
				Values: map[string]string{"Name": "Test"},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "email without domain",
			request: &pb.SendEmailRequest{
				To:           []string{"user@"},
				TemplateName: "test-template",
				Values:       map[string]string{"Name": "Test"},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "email without @ symbol",
			request: &pb.SendEmailRequest{
				To:           []string{"userexample.com"},
				TemplateName: "test-template",
				Values:       map[string]string{"Name": "Test"},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock expectations
			mockRepo.ExpectedCalls = nil

			// Execute the method
			resp, err := server.SendEmail(context.Background(), tt.request)

			// Assertions
			assert.Error(t, err)
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.errCode, st.Code())
			assert.Nil(t, resp)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSendEmail_DatabaseError(t *testing.T) {
	server, mockRepo := setupServer(t)

	request := &pb.SendEmailRequest{
		To:           []string{"test@example.com"},
		TemplateName: "test-template",
		Values:       map[string]string{"Name": "Test"},
	}

	mockRepo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).Return(errors.New("database connection failed")).Once()

	resp, err := server.SendEmail(context.Background(), request)

	// Assertions
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Nil(t, resp)

	mockRepo.AssertExpectations(t)
}

func TestUpdateStatus(t *testing.T) {
	server, mockRepo := setupServer(t)
	mailId := uuid.NewV4()

	tests := []struct {
		name       string
		status     ukama.Status
		setupMocks func(*mocks.MailerRepo)
		wantErr    bool
	}{
		{
			name:   "successful status update",
			status: ukama.Success,
			setupMocks: func(repo *mocks.MailerRepo) {
				repo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
					return m.MailId == mailId && m.Status == ukama.Success
				})).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name:   "database error during status update",
			status: ukama.Failed,
			setupMocks: func(repo *mocks.MailerRepo) {
				repo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
					return m.MailId == mailId && m.Status == ukama.Failed
				})).Return(errors.New("database error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMocks(mockRepo)

			err := server.updateStatus(mailId, tt.status)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateRetryStatus(t *testing.T) {
	server, mockRepo := setupServer(t)
	mailId := uuid.NewV4()
	retryCount := 2
	nextRetryTime := time.Now().Add(10 * time.Minute)

	tests := []struct {
		name          string
		retryCount    int
		nextRetryTime *time.Time
		setupMocks    func(*mocks.MailerRepo)
		wantErr       bool
	}{
		{
			name:          "successful retry status update",
			retryCount:    retryCount,
			nextRetryTime: &nextRetryTime,
			setupMocks: func(repo *mocks.MailerRepo) {
				repo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
					return m.MailId == mailId &&
						m.Status == ukama.Retry &&
						m.RetryCount == retryCount &&
						m.NextRetryTime != nil
				})).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name:          "retry status update with nil next retry time",
			retryCount:    retryCount,
			nextRetryTime: nil,
			setupMocks: func(repo *mocks.MailerRepo) {
				repo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
					return m.MailId == mailId &&
						m.Status == ukama.Retry &&
						m.RetryCount == retryCount &&
						m.NextRetryTime == nil
				})).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name:          "database error during retry status update",
			retryCount:    retryCount,
			nextRetryTime: &nextRetryTime,
			setupMocks: func(repo *mocks.MailerRepo) {
				repo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
					return m.MailId == mailId &&
						m.Status == ukama.Retry &&
						m.RetryCount == retryCount
				})).Return(errors.New("database error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMocks(mockRepo)

			err := server.updateRetryStatus(mailId, tt.retryCount, tt.nextRetryTime)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
