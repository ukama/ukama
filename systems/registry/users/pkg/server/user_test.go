package server

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/systems/registry/users/mocks"
	pb "github.com/ukama/ukama/systems/registry/users/pb/gen"

	"github.com/ukama/ukama/systems/registry/users/pkg/db"
)

func TestUserService_Add(t *testing.T) {
	userRepo := &mocks.UserRepo{}
	userUUID := uuid.New()

	userRequest := &pb.User{
		Uuid:  userUUID.String(),
		Name:  "Joe",
		Email: "test@example.com",
		Phone: "12324",
	}

	userRepo.On("Add", mock.Anything, mock.Anything).Return(&db.User{Uuid: userUUID,
		Email: userRequest.Email, Phone: userRequest.Phone, Name: userRequest.Name}, nil)

	t.Run("NewUser", func(tt *testing.T) {
		srv := NewUserService(userRepo)
		addResp, err := srv.Add(context.TODO(), &pb.AddRequest{User: userRequest})

		assert.NoError(t, err)
		assert.NotEmpty(t, addResp.User.Uuid)
		assert.Equal(t, userUUID.String(), addResp.User.Uuid)
		assert.Equal(t, userRequest.Name, addResp.User.Name)
		assert.Equal(t, userRequest.Phone, addResp.User.Phone)
		assert.Equal(t, userRequest.Email, addResp.User.Email)
	})
}

//TestGet
//UserNotFound

func TestUserService_Get(t *testing.T) {
	userRepo := &mocks.UserRepo{}
	userUUID := uuid.NewString()

	userRepo.On("Get", uuid.MustParse(userUUID)).Return(&db.User{
		Uuid: uuid.MustParse(userUUID),
	}, nil)

	t.Run("UserFound", func(tt *testing.T) {
		srv := NewUserService(userRepo)

		user, err := srv.Get(context.TODO(), &pb.GetRequest{UserId: userUUID})

		assert.NoError(t, err)
		assert.Equal(t, userUUID, user.GetUser().Uuid)
		userRepo.AssertExpectations(t)
	})
}

//TestUpdate

//TestDactivate
//UserAlreadyDeactivated
//UserNotAlreadyDeactivate

//TestDelete
//UserAlreadyDeactivated
//UserNotAlreadyDeactivate

func TestUserService_Deactivate(t *testing.T) {
	userRepo := &mocks.UserRepo{}
	userUUID := uuid.NewString()

	userRepo.On("Get", uuid.MustParse(userUUID)).Return(&db.User{
		Uuid: uuid.MustParse(userUUID),
	}, nil)

	userRepo.On("Update", mock.MatchedBy(func(u *db.User) bool {
		return u.Uuid.String() == userUUID
	})).Return(&db.User{}, nil)

	t.Run("UserNotAlreadyDeactivated", func(tt *testing.T) {
		srv := NewUserService(userRepo)

		_, err := srv.Deactivate(context.Background(), &pb.DeactivateRequest{
			UserId: userUUID,
		})

		assert.NoError(t, err, "Error deactivating user")
		userRepo.AssertExpectations(t)
	})
}

func TestUserService_Validation_Add(t *testing.T) {
	const name = "nn"

	tests := []struct {
		name        string
		user        *pb.User
		expectErr   bool
		errContains string
	}{
		{
			name:        "emptyName",
			user:        &pb.User{},
			expectErr:   true,
			errContains: "Name",
		},
		{
			name:        "email",
			user:        &pb.User{Email: "test_example.com", Name: name},
			expectErr:   true,
			errContains: "must be an email format",
		},
		{
			name:        "emailNoTopLevelDomain",
			user:        &pb.User{Email: "test@example", Name: name},
			expectErr:   true,
			errContains: "must be an email format",
		},
		{
			name:      "emailNotRequired",
			user:      &pb.User{Name: name},
			expectErr: false,
		},
		{
			name:        "emailIsEmpty",
			user:        &pb.User{Email: "@example.com", Name: name},
			expectErr:   true,
			errContains: "must be an email format",
		},

		{
			name:      "phone1",
			user:      &pb.User{Phone: "(+351) 282 43 50 50", Name: name},
			expectErr: false,
		},
		{
			name:      "phone2",
			user:      &pb.User{Phone: "90191919908", Name: name},
			expectErr: false,
		},

		{
			name:      "phone3",
			user:      &pb.User{Phone: "555-8909", Name: name},
			expectErr: false,
		},
		{
			name:      "phone4",
			user:      &pb.User{Phone: "001 6867684", Name: name},
			expectErr: false,
		},
		{
			name:      "phone5",
			user:      &pb.User{Phone: "1 (234) 567-8901", Name: name},
			expectErr: false,
		},
		{
			name:      "phone6",
			user:      &pb.User{Phone: "+1 34 567-8901", Name: name},
			expectErr: false,
		},
		{
			name:      "phoneEmpty",
			user:      &pb.User{Name: name},
			expectErr: false,
		},

		{
			name:        "phoneErr",
			user:        &pb.User{Phone: "sdfewr", Name: name},
			expectErr:   true,
			errContains: "phone number",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {

			// test add requeset
			r := &pb.AddRequest{
				User: test.user,
			}

			err := r.Validate()
			assertValidationErr(tt, err, test.expectErr, test.errContains)
		})
	}
}

func TestUserService_Validation_Update(t *testing.T) {
	const name = "nn"

	tests := []struct {
		name        string
		user        *pb.User
		expectErr   bool
		errContains string
	}{
		{
			name:      "emptyName",
			user:      &pb.User{},
			expectErr: false,
		},
		{
			name:        "email",
			user:        &pb.User{Email: "test_example.com", Name: name},
			expectErr:   true,
			errContains: "must be an email format",
		},
		{
			name:        "emailNoTopLevelDomain",
			user:        &pb.User{Email: "test@example", Name: name},
			expectErr:   true,
			errContains: "must be an email format",
		},
		{
			name:      "emailNotRequired",
			user:      &pb.User{Name: name},
			expectErr: false,
		},
		{
			name:        "emailIsEmpty",
			user:        &pb.User{Email: "@example.com", Name: name},
			expectErr:   true,
			errContains: "must be an email format",
		},

		{
			name:      "phone1",
			user:      &pb.User{Phone: "(+351) 282 43 50 50", Name: name},
			expectErr: false,
		},
		{
			name:      "phone2",
			user:      &pb.User{Phone: "90191919908", Name: name},
			expectErr: false,
		},

		{
			name:      "phone3",
			user:      &pb.User{Phone: "555-8909", Name: name},
			expectErr: false,
		},
		{
			name:      "phone4",
			user:      &pb.User{Phone: "001 6867684", Name: name},
			expectErr: false,
		},
		{
			name:      "phone5",
			user:      &pb.User{Phone: "1 (234) 567-8901", Name: name},
			expectErr: false,
		},
		{
			name:      "phone6",
			user:      &pb.User{Phone: "+1 34 567-8901", Name: name},
			expectErr: false,
		},
		{
			name:      "phoneEmpty",
			user:      &pb.User{Name: name},
			expectErr: false,
		},

		{
			name:        "phoneErr",
			user:        &pb.User{Phone: "sdfewr", Name: name},
			expectErr:   true,
			errContains: "phone number",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {

			// test update request
			ru := &pb.UpdateRequest{
				UserId: uuid.NewString(),
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
	t.Helper()

	if expectErr {
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), errContains)
		}
	} else {
		assert.NoError(t, err)
	}
}
