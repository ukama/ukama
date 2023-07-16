package server_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/registry/users/mocks"

	"github.com/ukama/ukama/systems/registry/users/pkg/db"
	"github.com/ukama/ukama/systems/registry/users/pkg/server"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	orgpb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	omocks "github.com/ukama/ukama/systems/registry/org/pb/gen/mocks"

	pb "github.com/ukama/ukama/systems/registry/users/pb/gen"
)

func TestUserService_Add(t *testing.T) {
	name := "Joe"
	email := "test@example.com"
	phone := "12324"
	authId := uuid.NewV4()

	userRepo := &mocks.UserRepo{}
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	user := &db.User{
		Name:   name,
		Email:  email,
		Phone:  phone,
		AuthId: authId,
	}

	userRequest := &pb.User{
		Name:   name,
		Email:  email,
		Phone:  phone,
		AuthId: authId.String(),
	}

	userRepo.On("Add", user, mock.Anything).Return(nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, &pb.AddRequest{User: userRequest}).Return(nil).Once()
	userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

	s := server.NewUserService(userRepo, nil, msgclientRepo, "")

	t.Run("AddValidUser", func(tt *testing.T) {
		aResp, err := s.Add(context.Background(), &pb.AddRequest{User: userRequest})

		assert.NoError(t, err)
		assert.NotEmpty(t, aResp.User.Id)

		assert.Equal(t, userRequest.Name, aResp.User.Name)
		assert.Equal(t, userRequest.Phone, aResp.User.Phone)
		assert.Equal(t, userRequest.Email, aResp.User.Email)
	})

	t.Run("AddNonValidUser", func(tt *testing.T) {
		userRequest.AuthId = "df7d48f9-9ca0-4f0d-89f1-42df51ea2f6z"

		aResp, err := s.Add(context.Background(), &pb.AddRequest{User: userRequest})

		assert.Error(t, err)
		assert.Nil(t, aResp)
	})
}

func TestUserService_Get(t *testing.T) {
	userRepo := &mocks.UserRepo{}
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	s := server.NewUserService(userRepo, nil, msgclientRepo, "")

	t.Run("UserFound", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("Get", userId).Return(&db.User{
			Id: userId,
		}, nil)

		uResp, err := s.Get(context.TODO(), &pb.GetRequest{UserId: userId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, uResp)

		assert.NoError(t, err)
		assert.Equal(t, userId.String(), uResp.GetUser().Id)
		userRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("Get", userId).Return(nil, gorm.ErrRecordNotFound).Once()

		uResp, err := s.Get(context.TODO(), &pb.GetRequest{UserId: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, uResp)
		userRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByAuthId(t *testing.T) {
	userRepo := &mocks.UserRepo{}
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	s := server.NewUserService(userRepo, nil, msgclientRepo, "")

	t.Run("UserFound", func(t *testing.T) {
		authId := uuid.NewV4()

		userRepo.On("GetByAuthId", authId).Return(&db.User{
			AuthId: authId,
		}, nil)

		uResp, err := s.GetByAuthId(context.TODO(), &pb.GetByAuthIdRequest{AuthId: authId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, uResp)

		assert.NoError(t, err)
		assert.Equal(t, authId.String(), uResp.GetUser().AuthId)
		userRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		authId := uuid.NewV4()

		userRepo.On("GetByAuthId", authId).Return(nil, gorm.ErrRecordNotFound).Once()

		uResp, err := s.GetByAuthId(context.TODO(), &pb.GetByAuthIdRequest{AuthId: authId.String()})

		assert.Error(t, err)
		assert.Nil(t, uResp)
		userRepo.AssertExpectations(t)
	})
}

//TestUpdate

//TestDactivate
//UserAlreadyDeactivated
//UserNotAlreadyDeactivate

func TestUserService_Deactivate(t *testing.T) {
	userRepo := &mocks.UserRepo{}
	userUUID := uuid.NewV4()
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	userRepo.On("Get", userUUID).Return(&db.User{
		Id: userUUID,
	}, nil)

	userRepo.On("Update", mock.MatchedBy(func(u *db.User) bool {
		return u.Id.String() == userUUID.String()
	}), mock.Anything).Return(nil)

	msgclientRepo.On("PublishRequest", mock.Anything, &pb.DeactivateRequest{UserId: userUUID.String()}).Return(nil).Once()
	userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

	s := server.NewUserService(userRepo, nil, msgclientRepo, "")

	t.Run("UserNotAlreadyDeactivated", func(tt *testing.T) {

		res, err := s.Deactivate(context.Background(), &pb.DeactivateRequest{
			UserId: userUUID.String(),
		})

		assert.NoError(t, err, "Error deactivating user")
		assert.NotNil(t, res, "Response should not be nil")

		userRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})
}

func TestUserService_Delete(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	userRepo := &mocks.UserRepo{}

	s := server.NewUserService(userRepo, nil, msgclientRepo, "")

	t.Run("UserFoundAndInactive", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("Get", userId).Return(&db.User{Id: userId, Deactivated: true}, nil).Once()
		userRepo.On("Delete", userId, mock.Anything).Return(nil).Once()

		msgclientRepo.On("PublishRequest", mock.Anything, &pb.DeleteRequest{
			UserId: userId.String(),
		}).Return(nil).Once()

		userRepo.On("GetUserCount").Return(int64(1), int64(0), nil).Once()

		resp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			UserId: userId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		userRepo.AssertExpectations(t)
	})

	t.Run("UserFoundAndActive", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("Get", userId).Return(&db.User{Id: userId, Deactivated: false}, nil).Once()

		resp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			UserId: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		userRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		userId := uuid.NewV4()

		userRepo.On("Get", userId).Return(nil, gorm.ErrRecordNotFound).Once()

		resp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			UserId: userId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		userRepo.AssertExpectations(t)
	})
}

func TestUserService_Whoami(t *testing.T) {
	name := "Joe"
	email := "test@example.com"
	phone := "12324"
	authId := uuid.NewV4()

	userRepo := &mocks.UserRepo{}
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	orgService := &mocks.OrgClientProvider{}

	user := &db.User{
		Name:   name,
		Email:  email,
		Phone:  phone,
		AuthId: authId,
	}

	s := server.NewUserService(userRepo, orgService, msgclientRepo, "")

	t.Run("NonValidUser", func(tt *testing.T) {
		userId := "df7d48f9-9ca0-4f0d-89f1-42df51ea2f6z"

		aResp, err := s.Whoami(context.Background(), &pb.GetRequest{UserId: userId})

		assert.Error(t, err)
		assert.Nil(t, aResp)
	})

	t.Run("UserNotFound", func(tt *testing.T) {
		userRepo.On("Get", user.Id).Return(nil, gorm.ErrRecordNotFound).Once()

		uResp, err := s.Whoami(context.TODO(), &pb.GetRequest{UserId: user.Id.String()})

		assert.Error(t, err)
		assert.Nil(t, uResp)
		userRepo.AssertExpectations(t)
	})

	t.Run("OrgServiceNotFound", func(tt *testing.T) {
		userRepo.On("Get", user.Id).Return(user, nil)

		orgService.On("GetClient").
			Return(nil, errors.New("Internal")).Once()

		uResp, err := s.Whoami(context.TODO(), &pb.GetRequest{UserId: user.Id.String()})

		assert.Error(t, err)
		assert.Nil(t, uResp)
		userRepo.AssertExpectations(t)
	})

	t.Run("OrgServiceUserNotFound", func(tt *testing.T) {
		orgService := &mocks.OrgClientProvider{}
		userRepo := &mocks.UserRepo{}

		userRepo.On("Get", user.Id).Return(user, nil)

		orgClient := orgService.On("GetClient").
			Return(&omocks.OrgServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*omocks.OrgServiceClient)

		orgClient.On("GetByUser", mock.Anything,
			&orgpb.GetByOwnerRequest{UserUuid: user.Id.String()}).
			Return(nil, errors.New("Not Found")).Once()

		s := server.NewUserService(userRepo, orgService, msgclientRepo, "")

		uResp, err := s.Whoami(context.TODO(), &pb.GetRequest{UserId: user.Id.String()})

		assert.Error(t, err)
		assert.Nil(t, uResp)
		userRepo.AssertExpectations(t)
	})

	t.Run("OrgServiceUserFound", func(tt *testing.T) {
		userRepo.On("Get", user.Id).Return(user, nil)

		orgClient := orgService.On("GetClient").
			Return(&omocks.OrgServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*omocks.OrgServiceClient)

		orgClient.On("GetByUser", mock.Anything,
			&orgpb.GetByOwnerRequest{UserUuid: user.Id.String()}).
			Return(&orgpb.GetByUserResponse{
				User: user.Id.String(),
				OwnerOf: []*orgpb.Organization{
					&orgpb.Organization{},
				},

				MemberOf: []*orgpb.OrgUser{
					&orgpb.OrgUser{},
					&orgpb.OrgUser{},
				},
			}, nil).Once()

		uResp, err := s.Whoami(context.TODO(), &pb.GetRequest{UserId: user.Id.String()})

		assert.NoError(t, err)
		assert.NotNil(t, uResp)

		assert.Equal(t, user.Id.String(), uResp.User.Id)
		assert.Equal(t, user.Name, uResp.User.Name)
		assert.Equal(t, user.Phone, uResp.User.Phone)
		assert.Equal(t, user.Email, uResp.User.Email)
		assert.Equal(t, 2, len(uResp.MemberOf))
		assert.Equal(t, 1, len(uResp.OwnerOf))
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
			test.user.Id = uuid.NewV4().String()
			test.user.AuthId = uuid.NewV4().String()

			// test add requeset
			r := &pb.AddRequest{
				User: &pb.UserAttributes{
					Name:  test.user.Name,
					Email: test.user.Email,
					Phone: test.user.Phone,
				},
			}

			err := r.Validate()
			assertValidationErr(tt, err, test.expectErr, test.errContains)
		})
	}
}

func TestUserService_Validation_Update(t *testing.T) {
	const name = "nn"
	var userId = uuid.NewV4()

	tests := []struct {
		name        string
		user        *pb.User
		expectErr   bool
		errContains string
	}{
		{
			name: "emptyName",
			user: &pb.User{
				Name: "",
			},
			expectErr: true,
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
		// { TODO: fix this test. Regex has a bug
		// 	name:        "phoneErr",
		// 	user:        &pb.User{Phone: "sdfewr", Name: name},
		// 	expectErr:   true,
		// 	errContains: "phone number",
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {

			// test update request
			ru := &pb.UpdateRequest{
				UserId: userId.String(),
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
