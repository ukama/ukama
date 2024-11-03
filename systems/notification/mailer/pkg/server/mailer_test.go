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

	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/mailer/mocks"
	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	"github.com/ukama/ukama/systems/notification/mailer/pkg"
	"github.com/ukama/ukama/systems/notification/mailer/pkg/db"
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

	return server, mockRepo
}

func TestSendEmail(t *testing.T) {
	server, mockRepo := setupServer(t)

	tests := []struct {
		name    string
		req     *pb.SendEmailRequest
		setup   func(*mocks.MailerRepo, *uuid.UUID)
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "successful email queueing",
			req: &pb.SendEmailRequest{
				To:           []string{"test@example.com"},
				TemplateName: "test-template",
				Values:       map[string]string{"name": "John"},
			},
			setup: func(repo *mocks.MailerRepo, mailId *uuid.UUID) {
				// Setup CreateEmail expectation
				repo.On("CreateEmail", mock.MatchedBy(func(m *db.Mailing) bool {
					*mailId = m.MailId // Capture the mailId for subsequent calls
					return m.Email == "test@example.com" &&
						m.TemplateName == "test-template" &&
						m.Status == db.Pending &&
						m.Values["name"] == "John"
				})).Return(nil).Once()

				// Setup GetEmailById expectation
				repo.On("GetEmailById", mock.MatchedBy(func(id uuid.UUID) bool {
					return id == *mailId
				})).Return(&db.Mailing{
					MailId:       *mailId,
					Email:        "test@example.com",
					TemplateName: "test-template",
					Status:       db.Pending,
					Values:       db.JSONMap{"name": "John"},
				}, nil).Once()

				// Setup UpdateEmailStatus expectations for processing and success
				repo.On("UpdateEmailStatus", mock.MatchedBy(func(m *db.Mailing) bool {
					return m.MailId == *mailId && m.Status == db.Process
				})).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "invalid email address",
			req: &pb.SendEmailRequest{
				To:           []string{"invalid-email"},
				TemplateName: "test-template",
				Values:       map[string]string{"name": "John"},
			},
			setup:   func(repo *mocks.MailerRepo, mailId *uuid.UUID) {},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "database error",
			req: &pb.SendEmailRequest{
				To:           []string{"test@example.com"},
				TemplateName: "test-template",
				Values:       map[string]string{"name": "John"},
			},
			setup: func(repo *mocks.MailerRepo, mailId *uuid.UUID) {
				repo.On("CreateEmail", mock.Anything).Return(errors.New("database error")).Once()
			},
			wantErr: true,
			errCode: codes.Internal,
		},
		{
			name: "empty template name",
			req: &pb.SendEmailRequest{
				To:           []string{"test@example.com"},
				TemplateName: "",
				Values:       map[string]string{"name": "John"},
			},
			setup:   func(repo *mocks.MailerRepo, mailId *uuid.UUID) {},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mailId uuid.UUID
			tt.setup(mockRepo, &mailId)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := server.SendEmail(ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.errCode, st.Code())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.MailId)
				assert.Contains(t, resp.Message, "Email queued")

				// Wait for async processing
				time.Sleep(200 * time.Millisecond)
			}

			mockRepo.AssertExpectations(t)
		})
	}
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
					Status:       db.Success,
					SentAt:       &testTime,
				}, nil).Once()
			},
			want: &pb.GetEmailByIdResponse{
				MailId:       testMailId.String(),
				TemplateName: "test-template",
				Status:       pb.Status(pb.Status_value[db.Success.String()]),
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
