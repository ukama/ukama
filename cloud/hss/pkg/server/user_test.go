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

	t.Run("WithDebugSimToken", func(tt *testing.T) {
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

func Test_Validation(t *testing.T) {
	const name = "nn"
	tests := []struct {
		name        string
		user        *pb.User
		expectErr   bool
		errContains string
	}{
		{name: "emptyName",
			user:        &pb.User{},
			expectErr:   true,
			errContains: "Name",
		},
		{name: "email",
			user:        &pb.User{Email: "test_example.com", Name: name},
			expectErr:   true,
			errContains: "must be an email format",
		},
		{name: "emailNoTopLevelDomain",
			user:        &pb.User{Email: "test@example", Name: name},
			expectErr:   true,
			errContains: "must be an email format",
		},
		{name: "emailNotRequired",
			user:      &pb.User{Name: name},
			expectErr: false,
		},
		{name: "emailIsEmpty",
			user:        &pb.User{Email: "@example.com", Name: name},
			expectErr:   true,
			errContains: "must be an email format",
		},

		{name: "phone1",
			user:      &pb.User{Phone: "(+351) 282 43 50 50", Name: name},
			expectErr: false,
		},
		{name: "phone2",
			user:      &pb.User{Phone: "90191919908", Name: name},
			expectErr: false,
		},

		{name: "phone3",
			user:      &pb.User{Phone: "555-8909", Name: name},
			expectErr: false,
		},
		{name: "phone4",
			user:      &pb.User{Phone: "001 6867684", Name: name},
			expectErr: false,
		},
		{name: "phone5",
			user:      &pb.User{Phone: "1 (234) 567-8901", Name: name},
			expectErr: false,
		},
		{name: "phone6",
			user:      &pb.User{Phone: "+1 34 567-8901", Name: name},
			expectErr: false,
		},
		{name: "phoneEmpty",
			user:      &pb.User{Name: name},
			expectErr: false,
		},

		{name: "phoneErr",
			user:        &pb.User{Phone: "sdfewr", Name: name},
			expectErr:   true,
			errContains: "phone number",
		},
	}

	for _, test := range tests {
		t.Run("addRequest_"+test.name, func(tt *testing.T) {

			// test add requeset
			r := &pb.AddRequest{
				Org:  testOrg,
				User: test.user,
			}
			err := r.Validate()
			assertValidationErr(tt, err, test.expectErr, test.errContains)
		})
	}

	for _, test := range tests {
		t.Run("updateRequest_"+test.name, func(tt *testing.T) {
			// test update request
			ru := &pb.UpdateRequest{
				Uuid: uuid.NewString(),
				User: &pb.UserAttributes{
					Phone: test.user.Phone,
					Email: test.user.Email,
					Name:  test.user.Name,
				},
			}
			err := ru.Validate()
			assertValidationErr(tt, err, test.expectErr, test.errContains)
		})
	}

}

func assertValidationErr(t *testing.T, err error, expectErr bool, errContains string) {
	if expectErr {
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), errContains)
		}
	} else {
		assert.NoError(t, err)
	}
}
