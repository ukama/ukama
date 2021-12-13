package server

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukamaX/cloud/hss/mocks"
	pb "github.com/ukama/ukamaX/cloud/hss/pb/gen"
	"github.com/ukama/ukamaX/cloud/hss/pkg/db"
	"testing"
)

const testOrg = "org"
const testImis = "1"

func Test_AddUser(t *testing.T) {
	// Arrange
	userRepo := &mocks.UserRepo{}
	imsiRepo := &mocks.ImsiRepo{}

	userRequest := &pb.User{
		Imsi:      "1",
		FirstName: "Joe",
		LastName:  "Doe",
		Email:     "test@example.com",
		Phone:     "12324",
	}

	userUuid := uuid.NewV4()
	userRepo.On("Add", mock.Anything).Return(&db.User{Uuid: userUuid,
		Email: userRequest.Email, Phone: userRequest.Phone,
		LastName: userRequest.LastName, FirstName: userRequest.FirstName}, nil)

	imsiRepo.On("Add", testOrg, mock.MatchedBy(func(n *db.Imsi) bool {
		return n.Imsi == testImis
	})).Return(nil)

	srv := NewUserService(userRepo, imsiRepo)

	// Act
	addResp, err := srv.Add(context.Background(), &pb.AddUserRequest{
		Org:  testOrg,
		User: userRequest,
	})

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, addResp.User.Uuid)
	assert.Equal(t, userUuid.String(), addResp.User.Uuid)
	assert.Equal(t, userRequest.FirstName, addResp.User.FirstName)
	assert.Equal(t, userRequest.LastName, addResp.User.LastName)
	assert.Equal(t, userRequest.Phone, addResp.User.Phone)
	assert.Equal(t, userRequest.Email, addResp.User.Email)
}
