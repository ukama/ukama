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
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/mailer/mocks"
	"github.com/ukama/ukama/systems/notification/mailer/pkg"
	"github.com/ukama/ukama/systems/notification/mailer/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
)

// Test data constants
const (
	testSMTPHost     = "smtp.example.com"
	testSMTPPort     = 587
	testSMTPUsername = "test@example.com"
	testSMTPPassword = "password"
	testSMTPFrom     = "from@example.com"
	testTemplateDir  = "../../templates"

	testTemplateName = "test-template"
	testEmail1       = "test@example.com"
	testEmail2       = "user1@example.com"
	testEmail3       = "user2@example.com"
	testEmail4       = "user3@example.com"
	testEmail5       = "user@example.com"
	testEmail6       = "admin@example.com"
	testEmail7       = "invalid-email"
	testEmail8       = "user@"
	testEmail9       = "userexample.com"

	testName1 = "John Doe"
	testName2 = "Team"
	testName3 = "User"
	testName4 = "Administrator"
	testName5 = "José María"
	testName6 = "Test User"

	testMessage1 = "Welcome to Ukama!"
	testMessage2 = "Meeting reminder"
	testMessage3 = "Please find attached"
	testMessage4 = "System notification"
	testMessage5 = "¡Hola! ¿Cómo estás?"

	testCode    = "ABC123"
	testDate    = "2024-01-15"
	testStatus  = "Active"
	testSymbols = "!@#$%^&*()_+-=[]{}|;':\",./<>?"

	testFilename1    = "document.pdf"
	testFilename2    = "image.jpg"
	testContentType1 = "application/pdf"
	testContentType2 = "image/jpeg"
	testPdfContent   = "fake pdf content"
	testImageContent = "fake image content"

	testSuccessMessage = "Email queued for sending!"
	testInvalidUUID    = "invalid-uuid"
	testEmptyString    = ""

	testRetryCount     = 2
	testRetryDuration  = 10 * time.Minute
	testSleepDuration1 = 300 * time.Millisecond
	testSleepDuration2 = 200 * time.Millisecond
	testSleepDuration3 = 100 * time.Millisecond
)

// Test error messages
var (
	errTestNotFound           = errors.New("not found")
	errTestDatabaseConnection = errors.New("database connection failed")
	errTestDatabaseError      = errors.New("database error")
	errTestUpdateFailed       = errors.New("update failed")
)

func setupServer(t *testing.T) (*MailerServer, *mocks.MailerRepo) {
	mockRepo := mocks.NewMailerRepo(t)
	mailer := &pkg.MailerConfig{
		Host:     testSMTPHost,
		Port:     testSMTPPort,
		Username: testSMTPUsername,
		Password: testSMTPPassword,
		From:     testSMTPFrom,
	}

	server, err := NewMailerServer(mockRepo, mailer, testTemplateDir)
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
					Email:        testEmail1,
					TemplateName: testTemplateName,
					Status:       ukama.MailStatusSuccess,
					SentAt:       &testTime,
				}, nil).Once()
			},
			want: &pb.GetEmailByIdResponse{
				MailId:       testMailId.String(),
				TemplateName: testTemplateName,
				Status:       pb.Status(pb.Status_value[ukama.MailStatusSuccess.String()]),
				SentAt:       testTime.String(),
			},
			wantErr: false,
		},
		{
			name:    "invalid UUID",
			mailId:  testInvalidUUID,
			setup:   func(repo *mocks.MailerRepo) {},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name:   "email not found",
			mailId: testMailId.String(),
			setup: func(repo *mocks.MailerRepo) {
				repo.On("GetEmailById", testMailId).Return(nil, errTestNotFound).Once()
			},
			wantErr: true,
			errCode: codes.Internal,
		},
		{
			name:    "empty mail ID",
			mailId:  testEmptyString,
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
			To:           []string{testEmail1},
			TemplateName: testTemplateName,
			Values:       map[string]interface{}{"Name": testName6},
			MailId:       mailId,
		}

		// Note: Email sending will fail due to fake SMTP credentials,
		// so we expect the retry path, not success path
		mockRepo.On("GetEmailById", mailId).Return(&db.Mailing{
			Status: ukama.MailStatusPending,
		}, nil).Once()
		mockRepo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
			return m.MailId == mailId && m.Status == ukama.MailStatusProcess
		})).Return(nil).Once()
		mockRepo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
			return m.MailId == mailId && m.Status == ukama.MailStatusRetry && m.RetryCount == 1
		})).Return(nil).Once()

		server.emailQueue <- payload

		time.Sleep(testSleepDuration1)

		mockRepo.AssertExpectations(t)
	})

	t.Run("email already sent - skip processing", func(t *testing.T) {
		mailId := uuid.NewV4()
		payload := &EmailPayload{
			To:           []string{testEmail1},
			TemplateName: testTemplateName,
			Values:       map[string]interface{}{"Name": testName6},
			MailId:       mailId,
		}

		mockRepo.On("GetEmailById", mailId).Return(&db.Mailing{
			Status: ukama.MailStatusSuccess,
		}, nil).Once()

		server.emailQueue <- payload

		time.Sleep(testSleepDuration2)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failed to fetch email - skip processing", func(t *testing.T) {
		mailId := uuid.NewV4()
		payload := &EmailPayload{
			To:           []string{testEmail1},
			TemplateName: testTemplateName,
			Values:       map[string]interface{}{"Name": testName6},
			MailId:       mailId,
		}

		mockRepo.On("GetEmailById", mailId).Return(nil, errTestDatabaseError).Once()

		server.emailQueue <- payload

		time.Sleep(testSleepDuration2)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failed to update status to process - skip processing", func(t *testing.T) {
		mailId := uuid.NewV4()
		payload := &EmailPayload{
			To:           []string{testEmail1},
			TemplateName: testTemplateName,
			Values:       map[string]interface{}{"Name": testName6},
			MailId:       mailId,
		}

		mockRepo.On("GetEmailById", mailId).Return(&db.Mailing{
			Status: ukama.MailStatusPending,
		}, nil).Once()
		mockRepo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
			return m.MailId == mailId && m.Status == ukama.MailStatusProcess
		})).Return(errTestUpdateFailed).Once()

		server.emailQueue <- payload

		time.Sleep(testSleepDuration2)

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
				To:           []string{testEmail1},
				TemplateName: testTemplateName,
				Values: map[string]string{
					"Name":    testName1,
					"Message": testMessage1,
				},
			},
			setup: func(repo *mocks.MailerRepo) {
				repo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).Return(nil).Once()
				// Setup expectation for background processing
				repo.On("GetEmailById", mock.AnythingOfType("uuid.UUID")).Return(&db.Mailing{
					Status: ukama.MailStatusPending,
				}, nil).Maybe()
				repo.On("UpdateEmailStatus", mock.AnythingOfType("*db.Mailing")).Return(nil).Maybe()
			},
			wantMessage: testSuccessMessage,
			wantErr:     false,
		},
		{
			name: "successful email queueing with multiple recipients",
			request: &pb.SendEmailRequest{
				To:           []string{testEmail2, testEmail3, testEmail4},
				TemplateName: testTemplateName,
				Values: map[string]string{
					"Name":    testName2,
					"Message": testMessage2,
				},
			},
			setup: func(repo *mocks.MailerRepo) {
				repo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).Return(nil).Once()
				repo.On("GetEmailById", mock.AnythingOfType("uuid.UUID")).Return(&db.Mailing{
					Status: ukama.MailStatusPending,
				}, nil).Maybe()
				repo.On("UpdateEmailStatus", mock.AnythingOfType("*db.Mailing")).Return(nil).Maybe()
			},
			wantMessage: testSuccessMessage,
			wantErr:     false,
		},
		{
			name: "successful email queueing with attachments",
			request: &pb.SendEmailRequest{
				To:           []string{testEmail5},
				TemplateName: testTemplateName,
				Values: map[string]string{
					"Name":    testName3,
					"Message": testMessage3,
				},
				Attachments: []*pb.Attachment{
					{
						Filename:    testFilename1,
						ContentType: testContentType1,
						Content:     []byte(testPdfContent),
					},
					{
						Filename:    testFilename2,
						ContentType: testContentType2,
						Content:     []byte(testImageContent),
					},
				},
			},
			setup: func(repo *mocks.MailerRepo) {
				repo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).Return(nil).Once()
				// Setup expectation for background processing
				repo.On("GetEmailById", mock.AnythingOfType("uuid.UUID")).Return(&db.Mailing{
					Status: ukama.MailStatusPending,
				}, nil).Maybe()
				repo.On("UpdateEmailStatus", mock.AnythingOfType("*db.Mailing")).Return(nil).Maybe()
			},
			wantMessage: testSuccessMessage,
			wantErr:     false,
		},
		{
			name: "successful email queueing with complex values",
			request: &pb.SendEmailRequest{
				To:           []string{testEmail6},
				TemplateName: testTemplateName,
				Values: map[string]string{
					"Name":    testName4,
					"Message": testMessage4,
					"Code":    testCode,
					"Date":    testDate,
					"Status":  testStatus,
				},
			},
			setup: func(repo *mocks.MailerRepo) {
				repo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).Return(nil).Once()

				repo.On("GetEmailById", mock.AnythingOfType("uuid.UUID")).Return(&db.Mailing{
					Status: ukama.MailStatusPending,
				}, nil).Maybe()
				repo.On("UpdateEmailStatus", mock.AnythingOfType("*db.Mailing")).Return(nil).Maybe()
			},
			wantMessage: testSuccessMessage,
			wantErr:     false,
		},
		{
			name: "successful email queueing with empty values",
			request: &pb.SendEmailRequest{
				To:           []string{testEmail5},
				TemplateName: testTemplateName,
				Values:       map[string]string{},
			},
			setup: func(repo *mocks.MailerRepo) {
				repo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).Return(nil).Once()
				repo.On("GetEmailById", mock.AnythingOfType("uuid.UUID")).Return(&db.Mailing{
					Status: ukama.MailStatusPending,
				}, nil).Maybe()
				repo.On("UpdateEmailStatus", mock.AnythingOfType("*db.Mailing")).Return(nil).Maybe()
			},
			wantMessage: testSuccessMessage,
			wantErr:     false,
		},
		{
			name: "successful email queueing with special characters in values",
			request: &pb.SendEmailRequest{
				To:           []string{testEmail5},
				TemplateName: testTemplateName,
				Values: map[string]string{
					"Name":    testName5,
					"Message": testMessage5,
					"Symbols": testSymbols,
				},
			},
			setup: func(repo *mocks.MailerRepo) {
				repo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).Return(nil).Once()
				// Setup expectation for background processing
				repo.On("GetEmailById", mock.AnythingOfType("uuid.UUID")).Return(&db.Mailing{
					Status: ukama.MailStatusPending,
				}, nil).Maybe()
				repo.On("UpdateEmailStatus", mock.AnythingOfType("*db.Mailing")).Return(nil).Maybe()
			},
			wantMessage: testSuccessMessage,
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

				time.Sleep(testSleepDuration3)
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
				TemplateName: testTemplateName,
				Values:       map[string]string{"Name": "Test"},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "invalid email address",
			request: &pb.SendEmailRequest{
				To:           []string{testEmail7},
				TemplateName: testTemplateName,
				Values:       map[string]string{"Name": "Test"},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "missing template name",
			request: &pb.SendEmailRequest{
				To:     []string{testEmail1},
				Values: map[string]string{"Name": "Test"},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "email without domain",
			request: &pb.SendEmailRequest{
				To:           []string{testEmail8},
				TemplateName: testTemplateName,
				Values:       map[string]string{"Name": "Test"},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "email without @ symbol",
			request: &pb.SendEmailRequest{
				To:           []string{testEmail9},
				TemplateName: testTemplateName,
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
		To:           []string{testEmail1},
		TemplateName: testTemplateName,
		Values:       map[string]string{"Name": "Test"},
	}

	mockRepo.On("CreateEmail", mock.AnythingOfType("*db.Mailing")).Return(errTestDatabaseConnection).Once()

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
		status     ukama.MailStatus
		setupMocks func(*mocks.MailerRepo)
		wantErr    bool
	}{
		{
			name:   "successful status update",
			status: ukama.MailStatusSuccess,
			setupMocks: func(repo *mocks.MailerRepo) {
				repo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
					return m.MailId == mailId && m.Status == ukama.MailStatusSuccess
				})).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name:   "database error during status update",
			status: ukama.MailStatusFailed,
			setupMocks: func(repo *mocks.MailerRepo) {
				repo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
					return m.MailId == mailId && m.Status == ukama.MailStatusFailed
				})).Return(errTestDatabaseError).Once()
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
	nextRetryTime := time.Now().Add(testRetryDuration)

	tests := []struct {
		name          string
		retryCount    int
		nextRetryTime *time.Time
		setupMocks    func(*mocks.MailerRepo)
		wantErr       bool
	}{
		{
			name:          "successful retry status update",
			retryCount:    testRetryCount,
			nextRetryTime: &nextRetryTime,
			setupMocks: func(repo *mocks.MailerRepo) {
				repo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
					return m.MailId == mailId &&
						m.Status == ukama.MailStatusRetry &&
						m.RetryCount == testRetryCount &&
						m.NextRetryTime != nil
				})).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name:          "retry status update with nil next retry time",
			retryCount:    testRetryCount,
			nextRetryTime: nil,
			setupMocks: func(repo *mocks.MailerRepo) {
				repo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
					return m.MailId == mailId &&
						m.Status == ukama.MailStatusRetry &&
						m.RetryCount == testRetryCount &&
						m.NextRetryTime == nil
				})).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name:          "database error during retry status update",
			retryCount:    testRetryCount,
			nextRetryTime: &nextRetryTime,
			setupMocks: func(repo *mocks.MailerRepo) {
				repo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
					return m.MailId == mailId &&
						m.Status == ukama.MailStatusRetry &&
						m.RetryCount == testRetryCount
				})).Return(errTestDatabaseError).Once()
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
