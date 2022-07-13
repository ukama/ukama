package server

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	hsspb "github.com/ukama/ukama/services/cloud/hss/pb/gen"
	hssmocks "github.com/ukama/ukama/services/cloud/hss/pb/gen/mocks"
	"github.com/ukama/ukama/services/cloud/users/mocks"
	pb "github.com/ukama/ukama/services/cloud/users/pb/gen"
	mocks2 "github.com/ukama/ukama/services/cloud/users/pb/gen/mocks"
	pbclient "github.com/ukama/ukama/services/cloud/users/pb/gen/simmgr"
	"github.com/ukama/ukama/services/cloud/users/pkg/db"
	"github.com/ukama/ukama/services/cloud/users/pkg/sims"
	commock "github.com/ukama/ukama/services/common/mocks"
)

const testOrg = "org"
const testImis = "1"

func Test_AddInternal(t *testing.T) {
	// Arrange
	userRepo := &mocks.UserRepo{}
	hssClient := &hssmocks.ImsiServiceClient{}
<<<<<<< HEAD
	kratosClient:=&mocks.KratosClient{}
=======
	kratosClient := &mocks.KratosClient{}
>>>>>>> f2070c6a81 (update user service test)
	simRepo := &mocks.SimcardRepo{}
	simManager := &mocks2.SimManagerServiceClient{}
	simProvider := &mocks.SimProvider{}
	hssProv := &mocks.ImsiClientProvider{}
	hssProv.On("GetClient").Return(hssClient, nil)

	userRequest := &pb.User{
		Name:  "Joe",
		Email: "test@example.com",
		Phone: "12324",
	}

	userUuid := uuid.New()
	userRepo.On("Add", mock.Anything, testOrg, mock.Anything).Return(&db.User{Uuid: userUuid,
		Email: userRequest.Email, Phone: userRequest.Phone,
		Name: userRequest.Name}, nil)

	hssClient.On("Add", mock.MatchedBy(func(n *hsspb.AddImsiRequest) bool {
		return n.Imsi.Imsi == testImis && n.Imsi.UserId == userUuid.String()
	})).Return(&hsspb.AddImsiResponse{}, nil)

	simManager.On("GetQrCode", mock.Anything, mock.Anything).Return(&pbclient.GetQrCodeResponse{
		QrCode: "qr",
	}, nil)

	pub := &commock.QPub{}
	pub.On("PublishToQueue", "mailer", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	srv := NewUserService(userRepo, hssProv, simRepo, simProvider, simManager, "simManager", pub, kratosClient)
<<<<<<< HEAD
	md := metadata.Pairs("x-requester", "89273897297392")
	ctx := metadata.NewOutgoingContext(context.TODO(), md)
=======
	md := metadata.Pairs("x-requester", tesRequesterId)
	ctx := metadata.NewIncomingContext(context.TODO(), md)
>>>>>>> f2070c6a81 (update user service test)
	// Act
	addResp, err := srv.AddInternal(ctx, &pb.AddInternalRequest{
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
	hssClient := &hssmocks.ImsiServiceClient{}
	simRepo := &mocks.SimcardRepo{}
	simManager := &mocks2.SimManagerServiceClient{}
	hssProv := &mocks.ImsiClientProvider{}
	hssProv.On("GetClient").Return(hssClient, nil)

	userRequest := &pb.User{
		Name:  "Joe",
		Email: "test@example.com",
		Phone: "12324",
	}

	userUuid := uuid.New()
	userRepo.On("Add", mock.Anything, testOrg, mock.Anything).Return(&db.User{Uuid: userUuid,
		Email: userRequest.Email, Phone: userRequest.Phone,
		Name: userRequest.Name}, nil)

	hssClient.On("Add", mock.MatchedBy(func(n *hsspb.AddImsiRequest) bool {
		return n.Imsi.Imsi == testImis && n.Imsi.UserId == userUuid.String()
	})).Return(&hsspb.AddImsiResponse{}, nil)

	simManager.On("GetQrCode", mock.Anything, mock.Anything).Return(&pbclient.GetQrCodeResponse{
		QrCode: "qr",
	}, nil).Maybe()

	pub := &commock.QPub{}

	t.Run("WithSimToken", func(tt *testing.T) {
		simProvider := &mocks.SimProvider{}
		simProvider.On("GetICCIDWithCode", TEST_SIM_TOKEN).Return(sims.GetDubugIccid(), nil)
<<<<<<< HEAD
kratosClient:=&mocks.KratosClient{}
=======
		kratosClient := &mocks.KratosClient{}
>>>>>>> f2070c6a81 (update user service test)
		srv := NewUserService(userRepo, hssProv, simRepo, simProvider, simManager, "simManager", pub, kratosClient)
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
<<<<<<< HEAD
kratosClient :=&mocks.KratosClient{}
=======
		kratosClient := &mocks.KratosClient{}
>>>>>>> f2070c6a81 (update user service test)
		pub := &commock.QPub{}
		pub.On("PublishToQueue", "mailer", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		srv := NewUserService(userRepo, hssProv, simRepo, simProvider, simManager, "simManager", pub, kratosClient)
<<<<<<< HEAD
=======

		md := metadata.Pairs("x-requester", tesRequesterId)
		ctx := metadata.NewIncomingContext(context.TODO(), md)
		kratosClient.On("GetAccountName", tesRequesterId).Return("TestNO", nil)

>>>>>>> f2070c6a81 (update user service test)
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
		pub := &commock.QPub{}
<<<<<<< HEAD
		kratosClient :=&mocks.KratosClient{}
		pub.On("PublishToQueue", "mailer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		srv := NewUserService(userRepo, hssProv, simRepo, simProvider, simManager, "simManager", pub, kratosClient)
=======
		kratosClient := &mocks.KratosClient{}
		pub.On("PublishToQueue", "mailer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		srv := NewUserService(userRepo, hssProv, simRepo, simProvider, simManager, "simManager", pub, kratosClient)

		kratosClient.On("GetAccountName", tesRequesterId).Return("TestNO", nil)
		md := metadata.Pairs("x-requester", tesRequesterId)
		ctx := metadata.NewIncomingContext(context.TODO(), md)

>>>>>>> f2070c6a81 (update user service test)
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

func Test_Deactivate(t *testing.T) {
	// Arrange
	userRepo := &mocks.UserRepo{}
	hssClient := &hssmocks.ImsiServiceClient{}
	simRepo := &mocks.SimcardRepo{}
	simManager := &mocks2.SimManagerServiceClient{}
	simProvider := &mocks.SimProvider{}
<<<<<<< HEAD
	kratosClient :=&mocks.KratosClient{}
=======
	kratosClient := &mocks.KratosClient{}
>>>>>>> f2070c6a81 (update user service test)
	hssProv := &mocks.ImsiClientProvider{}
	hssProv.On("GetClient").Return(hssClient, nil)
	pub := &commock.QPub{}

	iccid := sims.GetDubugIccid()
	userId := uuid.NewString()
	userRepo.On("Get", uuid.MustParse(userId)).Return(&db.User{
		Uuid: uuid.MustParse(userId),
		Simcard: db.Simcard{
			Iccid: iccid,
		},
	}, nil)
	userRepo.On("Update", mock.MatchedBy(func(u *db.User) bool {
		return u.Uuid.String() == userId
	})).Return(&db.User{}, nil)

	simManager.On("TerminateSim", mock.Anything, mock.MatchedBy(func(t *pbclient.TerminateSimRequest) bool {
		return t.Iccid == iccid
	})).Return(&pbclient.TerminateSimResponse{}, nil)

	hssClient.On("Delete", mock.Anything, mock.Anything).Return(&hsspb.DeleteImsiResponse{}, nil)

	srv := NewUserService(userRepo, hssProv, simRepo, simProvider, simManager, "simManager", pub, kratosClient)

	_, err := srv.DeactivateUser(context.Background(), &pb.DeactivateUserRequest{
		UserId: userId,
	})
	if assert.NoError(t, err, "Error deactivating user") {
		hssClient.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	}
}

func Test_AddValidation(t *testing.T) {
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
		t.Run(test.name, func(tt *testing.T) {

			// test add requeset
			r := &pb.AddRequest{
				Org:  testOrg,
				User: test.user,
			}
			err := r.Validate()
			assertValidationErr(tt, err, test.expectErr, test.errContains)
		})
	}
}

func Test_UpdateValidation(t *testing.T) {
	const name = "nn"
	tests := []struct {
		name        string
		user        *pb.User
		expectErr   bool
		errContains string
	}{
		{name: "emptyName",
			user:      &pb.User{},
			expectErr: false,
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

func Test_UpdateServices(t *testing.T) {
	// Arrange
	userRepo := &mocks.UserRepo{}
	simProvider := &mocks.SimProvider{}
	hssProv := &mocks.ImsiClientProvider{}
	testIccid := "890000000000000001"

	simRepo := &mocks.SimcardRepo{}
	simRepo.On("Get", testIccid).Return(&db.Simcard{Iccid: testIccid, Services: []*db.Service{
		{Data: false, Type: db.ServiceTypeUkama},
		{Data: true, Type: db.ServiceTypeCarrier},
	}}, nil)
	simRepo.On("UpdateServices", mock.Anything, mock.Anything, mock.MatchedBy(func(f func() error) bool {
		// call the function that is passed as nestec func
		err := f()
		if err != nil {
			t.Errorf("Error calling nested simmanager.updateServices. Error: %v", err)
			t.Fail()
		}
		return true
	})).Return(nil)

	pub := &commock.QPub{}

	t.Run("UpdateUkama", func(tt *testing.T) {
		simManager := &mocks2.SimManagerServiceClient{}
<<<<<<< HEAD
		kratosClient :=&mocks.KratosClient{}

		srv := NewUserService(userRepo, hssProv, simRepo, simProvider, simManager, "simManager", pub,kratosClient)
=======
		kratosClient := &mocks.KratosClient{}

		srv := NewUserService(userRepo, hssProv, simRepo, simProvider, simManager, "simManager", pub, kratosClient)
>>>>>>> f2070c6a81 (update user service test)
		// Act
		resp, err := srv.SetSimStatus(context.TODO(), &pb.SetSimStatusRequest{
			Iccid: testIccid,
			Ukama: &pb.SetSimStatusRequest_SetServices{
				Data: wrapperspb.Bool(true),
			},
		})

		// Assert
		if assert.NoError(t, err) {
			simRepo.AssertExpectations(tt)
			simManager.AssertExpectations(tt)
			assert.NotNil(tt, resp)
		}
	})

	t.Run("UpdateCarrier", func(tt *testing.T) {
		simManager := &mocks2.SimManagerServiceClient{}
<<<<<<< HEAD
		kratosClient :=&mocks.KratosClient{}
=======
		kratosClient := &mocks.KratosClient{}
>>>>>>> f2070c6a81 (update user service test)
		simManager.On("SetServiceStatus", mock.Anything, mock.MatchedBy(func(p *pbclient.SetServiceStatusRequest) bool {
			return p.Services.Data.GetValue()
		})).Return(nil, nil)

		srv := NewUserService(userRepo, hssProv, simRepo, simProvider, simManager, "simManager", pub, kratosClient)
		// Act
		resp, err := srv.SetSimStatus(context.TODO(), &pb.SetSimStatusRequest{
			Iccid: testIccid,
			Ukama: &pb.SetSimStatusRequest_SetServices{
				Data: wrapperspb.Bool(true),
			},
			Carrier: &pb.SetSimStatusRequest_SetServices{
				Data: wrapperspb.Bool(true),
			},
		})

		// Assert
		if assert.NoError(t, err) {
			simRepo.AssertExpectations(tt)
			simManager.AssertExpectations(tt)
			assert.NotNil(tt, resp)
		}
	})

	t.Run("DisableAllServicesButKeepCarrier", func(tt *testing.T) {
		simManager := &mocks2.SimManagerServiceClient{}
<<<<<<< HEAD
		kratosClient :=&mocks.KratosClient{}
=======
		kratosClient := &mocks.KratosClient{}
>>>>>>> f2070c6a81 (update user service test)
		simManager.On("SetServiceStatus", mock.Anything, mock.MatchedBy(func(p *pbclient.SetServiceStatusRequest) bool {
			return p.Services.Data != nil && p.Services.Data.GetValue() == false
		})).Return(nil, nil)

		srv := NewUserService(userRepo, hssProv, simRepo, simProvider, simManager, "simManager", pub, kratosClient)
		// Act
		resp, err := srv.SetSimStatus(context.TODO(), &pb.SetSimStatusRequest{
			Iccid: testIccid,
			Ukama: &pb.SetSimStatusRequest_SetServices{
				Data: wrapperspb.Bool(false),
			},
			Carrier: &pb.SetSimStatusRequest_SetServices{
				Data: wrapperspb.Bool(true),
			},
		})

		// Assert
		if assert.NoError(t, err) {
			simRepo.AssertExpectations(tt)
			simManager.AssertExpectations(tt)
			assert.NotNil(tt, resp)
		}
	})
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