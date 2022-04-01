package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukamaX/cloud/hss/mocks"
	pb "github.com/ukama/ukamaX/cloud/hss/pb/gen"
	mocks2 "github.com/ukama/ukamaX/cloud/hss/pb/gen/mocks"
	"github.com/ukama/ukamaX/cloud/hss/pkg/db"
	"github.com/ukama/ukamaX/cloud/hss/pkg/sims"
	"testing"
)

const testOrg = "org"
const testImis = "1"

func Test_AddInternal(t *testing.T) {
	// Arrange
	userRepo := &mocks.UserRepo{}
	imsiRepo := &mocks.ImsiRepo{}
	simRepo := &mocks.SimcardRepo{}
	simManager := &mocks2.SimManagerServiceClient{}
	simProvider := &mocks.SimProvider{}

	userRequest := &pb.User{
		Name:  "Joe",
		Email: "test@example.com",
		Phone: "12324",
	}

	userUuid := uuid.New()
	userRepo.On("Add", mock.Anything, testOrg, mock.Anything).Return(&db.User{Uuid: userUuid,
		Email: userRequest.Email, Phone: userRequest.Phone,
		Name: userRequest.Name}, nil)

	imsiRepo.On("Add", testOrg, mock.MatchedBy(func(n *db.Imsi) bool {
		return n.Imsi == testImis
	})).Return(nil)

	srv := NewUserService(userRepo, imsiRepo, simRepo, simProvider, simManager, "simManager")

	// Act
	addResp, err := srv.AddInternal(context.Background(), &pb.AddInternalRequest{
		Org:  testOrg,
		User: userRequest,
	})

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, addResp.User.Uuid)
	assert.Equal(t, userUuid.String(), addResp.User.Uuid)
	assert.Equal(t, userRequest.Name, addResp.User.Name)
	assert.Equal(t, userRequest.Phone, addResp.User.Phone)
	assert.Equal(t, userRequest.Email, addResp.User.Email)
}

const TEST_SIM_TOKEN = "QQQQQQQQQQQ"

func Test_Add(t *testing.T) {
	// Arrange
	userRepo := &mocks.UserRepo{}
	imsiRepo := &mocks.ImsiRepo{}
	simRepo := &mocks.SimcardRepo{}
	simManager := &mocks2.SimManagerServiceClient{}

	userRequest := &pb.User{
		Name:  "Joe",
		Email: "test@example.com",
		Phone: "12324",
	}

	userUuid := uuid.New()
	userRepo.On("Add", mock.Anything, testOrg, mock.Anything).Return(&db.User{Uuid: userUuid,
		Email: userRequest.Email, Phone: userRequest.Phone,
		Name: userRequest.Name}, nil)

	imsiRepo.On("Add", testOrg, mock.MatchedBy(func(n *db.Imsi) bool {
		return n.Imsi == testImis
	})).Return(nil)

	t.Run("WithSimToken", func(tt *testing.T) {
		simProvider := &mocks.SimProvider{}
		simProvider.On("GetICCIDWithCode", TEST_SIM_TOKEN).Return(sims.GetDubugIccid(), nil)

		srv := NewUserService(userRepo, imsiRepo, simRepo, simProvider, simManager, "simManager")
		// Act
		addResp, err := srv.Add(context.Background(), &pb.AddRequest{
			Org:      testOrg,
			User:     userRequest,
			SimToken: TEST_SIM_TOKEN,
		})

		// Assert
		if assert.NoError(t, err) {
			assert.NotEmpty(t, addResp.User.Uuid)
			assert.Equal(t, userUuid.String(), addResp.User.Uuid)
			simProvider.AssertExpectations(tt)
			simManager.AssertExpectations(tt)
		}
	})

	t.Run("WithoutSimToken", func(tt *testing.T) {
		simProvider := &mocks.SimProvider{}
		simProvider.On("GetICCIDFromPool").Return(sims.GetDubugIccid(), nil)

		srv := NewUserService(userRepo, imsiRepo, simRepo, simProvider, simManager, "simManager")
		// Act
		addResp, err := srv.Add(context.Background(), &pb.AddRequest{
			Org:  testOrg,
			User: userRequest,
		})

		// Assert
		if assert.NoError(t, err) {
			assert.NotEmpty(t, addResp.User.Uuid)
			assert.Equal(t, userUuid.String(), addResp.User.Uuid)
			simProvider.AssertExpectations(tt)
		}
	})

	t.Run("WutDebugSimToken", func(tt *testing.T) {
		simProvider := &mocks.SimProvider{}
		srv := NewUserService(userRepo, imsiRepo, simRepo, simProvider, simManager, "simManager")
		// Act
		addResp, err := srv.Add(context.Background(), &pb.AddRequest{
			Org:      testOrg,
			User:     userRequest,
			SimToken: "I_DO_NOT_NEED_A_SIM",
		})

		// Assert
		if assert.NoError(t, err) {
			assert.NotEmpty(t, addResp.User.Uuid)
			assert.Equal(t, userUuid.String(), addResp.User.Uuid)
			simProvider.AssertExpectations(tt)
		}
	})
}
