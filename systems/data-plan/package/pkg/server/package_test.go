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
	"log"
	"testing"
	"time"

	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/data-plan/package/mocks"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	ukama "github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	rpb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	splmocks "github.com/ukama/ukama/systems/data-plan/rate/pb/gen/mocks"
)

const OrgName = "testorg"
const OrgId = "8c6c2bec-5f90-4fee-8ffd-ee6456abf4fc"

// Get Packages //
// Success case
type UkamaDbMock struct {
	GormDb *gorm.DB
}

func (u UkamaDbMock) Init(model ...interface{}) error {
	panic("implement me: Init()")
}

func (u UkamaDbMock) Connect() error {
	panic("implement me: Connect()")
}

func (u UkamaDbMock) GetGormDb() *gorm.DB {
	return u.GormDb
}

func (u UkamaDbMock) InitDB() error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func() error) error {
	log.Fatal("implement me: ExecuteInTransaction()")
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func(tx *gorm.DB) error) error {
	log.Fatal("implement me: ExecuteInTransaction2()")
	return nil
}

const testSim = "ukama_data"

func TestPackageServer_GetPackages_Success(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	packageUUID := uuid.NewV4()

	var mockFilters = &pb.GetPackageRequest{
		Uuid: packageUUID.String(),
	}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageRepo.On("Get", packageUUID).Return(&db.Package{
		Name: "Daily-pack",
	}, nil)
	pkg, err := s.Get(context.TODO(), mockFilters)
	assert.NoError(t, err)
	assert.Equal(t, "Daily-pack", pkg.Package.Name)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_GetPackages_Error1(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageUUID := uuid.NewV4()
	var mockFilters = &pb.GetPackageRequest{
		Uuid: packageUUID.String(),
	}
	packageRepo.On("Get", packageUUID).Return(nil, grpc.SqlErrorToGrpc(errors.New("SQL error while fetching records"), "packages"))
	pkg2, err := s.Get(context.TODO(), mockFilters)
	assert.Error(t, err)
	assert.Nil(t, pkg2)
}

func TestPackageServer_GetPackage_Error(t *testing.T) {
	var orgId = uuid.NewV4()

	packageRepo := &mocks.PackageRepo{}
	var mockFilters = &pb.GetAllRequest{}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, orgId.String())

	packageRepo.On("GetAll").
		Return(nil, grpc.SqlErrorToGrpc(errors.New("SQL error while fetching records"), "packages"))
	pkg, err := s.GetAll(context.TODO(), mockFilters)
	assert.Error(t, err)
	assert.Nil(t, pkg)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_AddPackage(t *testing.T) {
	ownerId := uuid.NewV4().String()
	baserate := uuid.NewV4().String()
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}
	packageRepo.On("Add", mock.MatchedBy(func(p *db.Package) bool {
		return p.Active == true && p.Name == "daily-pack"
	}), mock.Anything).Return(nil).Once()

	rateClient := rate.On("GetClient").
		Return(&splmocks.RateServiceClient{}, nil).
		Once().
		ReturnArguments.Get(0).(*splmocks.RateServiceClient)

	rateClient.On("GetRateById", &rpb.GetRateByIdRequest{
		OwnerId:  ownerId,
		BaseRate: baserate,
	}).Return(&rpb.GetRateByIdResponse{
		Rate: &bpb.Rate{
			SmsMo:    1,
			SmsMt:    1,
			Data:     10,
			Country:  "USA",
			Provider: "ukama",
		},
	}, nil).Once()

	rateResp := rateClient.On("GetRateById", mock.Anything,
		&rpb.GetRateByIdRequest{
			OwnerId:  ownerId,
			BaseRate: baserate,
		}).
		Return(&rpb.GetRateByIdResponse{
			Rate: &bpb.Rate{
				SmsMo:    1,
				SmsMt:    1,
				Data:     10,
				Country:  "USA",
				Provider: "ukama",
			},
		}, nil).Once().
		ReturnArguments.Get(0).(*rpb.GetRateByIdResponse)

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	ActPackage, err := s.Add(context.TODO(), &pb.AddPackageRequest{
		Active:     true,
		Name:       "daily-pack",
		SimType:    testSim,
		OwnerId:    ownerId,
		BaserateId: baserate,
		Country:    rateResp.Rate.Country,
		Provider:   rateResp.Rate.Provider,
		From:       time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:         time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, ActPackage.Package.Active)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_UpdatePackage(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	s := NewPackageServer(OrgName, packageRepo, nil, msgbusClient, OrgId)
	packageUUID := uuid.NewV4()
	mockPackage := &pb.UpdatePackageRequest{
		Name: "Daily-pack-updated",
	}
	packageRepo.On("Update", packageUUID, mock.MatchedBy(func(p *db.Package) bool {
		return p.Active == true && p.Name == "Daily-pack-updated"
	})).Return(nil).Once()

	// Mock the Get call to return a complete package
	packageRepo.On("Get", packageUUID).Return(&db.Package{
		Uuid:   packageUUID,
		Name:   "Daily-pack-updated",
		Active: true,
		From:   time.Now(),
		To:     time.Now().Add(time.Hour * 24 * 30),
	}, nil).Once()

	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
	pkg, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
		Uuid:   packageUUID.String(),
		Name:   "Daily-pack-updated",
		Active: true,
	})
	assert.NoError(t, err)
	assert.Equal(t, mockPackage.Name, pkg.Package.Name)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_DeletePackage_Error1(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageUUID := uuid.NewV4()
	var mockFilters = &pb.DeletePackageRequest{
		Uuid: packageUUID.String(),
	}
	packageRepo.On("Delete", packageUUID).
		Return(status.Errorf(codes.InvalidArgument, "OrgId is required."))
	pkg1, err := s.Delete(context.TODO(), mockFilters)
	assert.Error(t, err)
	assert.Nil(t, pkg1)
}

func TestPackageServer_DeletePackage_Success_Error2(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageUUID := uuid.NewV4()
	var mockFilters = &pb.DeletePackageRequest{
		Uuid: packageUUID.String(),
	}
	packageRepo.On("Delete", packageUUID).
		Return(status.Errorf(codes.InvalidArgument, "Id is required."))
	pkg2, err := s.Delete(context.TODO(), mockFilters)
	assert.Error(t, err)
	assert.Nil(t, pkg2)
}

func TestPackageServer_Get_InvalidUUID(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)

	req := &pb.GetPackageRequest{
		Uuid: "invalid-uuid",
	}

	pkg, err := s.Get(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, pkg)
	assert.Contains(t, err.Error(), "invalid format of package uuid")
}

func TestPackageServer_GetDetails_Success(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	packageUUID := uuid.NewV4()

	req := &pb.GetPackageRequest{
		Uuid: packageUUID.String(),
	}

	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageRepo.On("GetDetails", packageUUID).Return(&db.Package{
		Name: "Test Package Details",
	}, nil)

	pkg, err := s.GetDetails(context.TODO(), req)
	assert.NoError(t, err)
	assert.Equal(t, "Test Package Details", pkg.Package.Name)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_GetDetails_InvalidUUID(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)

	req := &pb.GetPackageRequest{
		Uuid: "invalid-uuid",
	}

	pkg, err := s.GetDetails(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, pkg)
	assert.Contains(t, err.Error(), "invalid format of package uuid")
}

func TestPackageServer_GetDetails_DatabaseError(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	packageUUID := uuid.NewV4()
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)

	req := &pb.GetPackageRequest{
		Uuid: packageUUID.String(),
	}

	// Mock database error
	packageRepo.On("GetDetails", packageUUID).Return(nil, errors.New("database connection failed"))

	pkg, err := s.GetDetails(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, pkg)
	assert.Contains(t, err.Error(), "database connection failed")
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_GetDetails_EmptyUUID(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)

	req := &pb.GetPackageRequest{
		Uuid: "",
	}

	pkg, err := s.GetDetails(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, pkg)
	assert.Contains(t, err.Error(), "invalid format of package uuid")
}

func TestPackageServer_GetDetails_CompletePackageData(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	packageUUID := uuid.NewV4()
	ownerUUID := uuid.NewV4()
	baseRateUUID := uuid.NewV4()
	now := time.Now()

	req := &pb.GetPackageRequest{
		Uuid: packageUUID.String(),
	}

	completePackage := &db.Package{
		Uuid:          packageUUID,
		OwnerId:       ownerUUID,
		Name:          "Premium Data Package",
		SimType:       ukama.ParseSimType("ukama_data"),
		Active:        true,
		Duration:      30,
		SmsVolume:     1000,
		DataVolume:    5000,
		VoiceVolume:   100,
		Type:          ukama.ParsePackageType("prepaid"),
		DataUnits:     ukama.ParseDataUnitType("MegaBytes"),
		VoiceUnits:    ukama.ParseCallUnitType("minutes"),
		MessageUnits:  ukama.ParseMessageType("int"),
		Flatrate:      false,
		Currency:      "USD",
		From:          now,
		To:            now.AddDate(0, 1, 0),
		Country:       "USA",
		Provider:      "ukama",
		Overdraft:     10.5,
		TrafficPolicy: 1,
		Networks:      []string{"network1", "network2"},
		SyncStatus:    ukama.ParseStatusType("pending"),
		PackageRate: db.PackageRate{
			Amount: 99.99,
			SmsMo:  0.05,
			SmsMt:  0.05,
			Data:   0.02,
		},
		PackageMarkup: db.PackageMarkup{
			BaseRateId: baseRateUUID,
			Markup:     15.0,
		},
		PackageDetails: db.PackageDetails{
			Apn:  "internet.ukama.com",
			Dlbr: 10240000,
			Ulbr: 10240000,
		},
	}

	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageRepo.On("GetDetails", packageUUID).Return(completePackage, nil)

	pkg, err := s.GetDetails(context.TODO(), req)
	assert.NoError(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, packageUUID.String(), pkg.Package.Uuid)
	assert.Equal(t, "Premium Data Package", pkg.Package.Name)
	assert.Equal(t, true, pkg.Package.Active)
	assert.Equal(t, uint64(30), pkg.Package.Duration)
	assert.Equal(t, int64(1000), pkg.Package.SmsVolume)
	assert.Equal(t, int64(5000), pkg.Package.DataVolume)
	assert.Equal(t, int64(100), pkg.Package.VoiceVolume)
	assert.Equal(t, "prepaid", pkg.Package.Type)
	assert.Equal(t, "MegaBytes", pkg.Package.DataUnit)
	assert.Equal(t, "minutes", pkg.Package.VoiceUnit)
	assert.Equal(t, "int", pkg.Package.MessageUnit)
	assert.Equal(t, false, pkg.Package.Flatrate)
	assert.Equal(t, "USD", pkg.Package.Currency)
	assert.Equal(t, "USA", pkg.Package.Country)
	assert.Equal(t, "ukama", pkg.Package.Provider)
	assert.Equal(t, float64(10.5), pkg.Package.Overdraft)
	assert.Equal(t, uint32(1), pkg.Package.TrafficPolicy)
	assert.Equal(t, []string{"network1", "network2"}, pkg.Package.Networks)
	assert.Equal(t, "pending", pkg.Package.SyncStatus)
	assert.Equal(t, "internet.ukama.com", pkg.Package.Apn)

	assert.Equal(t, float64(99.99), pkg.Package.Rate.Amount)
	assert.Equal(t, float64(0.05), pkg.Package.Rate.SmsMo)
	assert.Equal(t, float64(0.05), pkg.Package.Rate.SmsMt)
	assert.Equal(t, float64(0.02), pkg.Package.Rate.Data)

	assert.Equal(t, baseRateUUID.String(), pkg.Package.Markup.Baserate)
	assert.Equal(t, float64(15.0), pkg.Package.Markup.Markup)

	packageRepo.AssertExpectations(t)
}

func TestPackageServer_GetDetails_InactivePackage(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	packageUUID := uuid.NewV4()

	req := &pb.GetPackageRequest{
		Uuid: packageUUID.String(),
	}

	inactivePackage := &db.Package{
		Uuid:   packageUUID,
		Name:   "Inactive Package",
		Active: false,
		From:   time.Now(),
		To:     time.Now().AddDate(0, 1, 0),
	}

	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageRepo.On("GetDetails", packageUUID).Return(inactivePackage, nil)

	pkg, err := s.GetDetails(context.TODO(), req)
	assert.NoError(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, false, pkg.Package.Active)
	assert.Equal(t, "Inactive Package", pkg.Package.Name)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_GetDetails_FlatratePackage(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	packageUUID := uuid.NewV4()

	req := &pb.GetPackageRequest{
		Uuid: packageUUID.String(),
	}

	flatratePackage := &db.Package{
		Uuid:     packageUUID,
		Name:     "Flatrate Package",
		Active:   true,
		Flatrate: true,
		From:     time.Now(),
		To:       time.Now().AddDate(0, 1, 0),
		PackageRate: db.PackageRate{
			Amount: 29.99,
			SmsMo:  0.0,
			SmsMt:  0.0,
			Data:   0.0,
		},
	}

	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageRepo.On("GetDetails", packageUUID).Return(flatratePackage, nil)

	pkg, err := s.GetDetails(context.TODO(), req)
	assert.NoError(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, true, pkg.Package.Flatrate)
	assert.Equal(t, float64(29.99), pkg.Package.Rate.Amount)
	assert.Equal(t, float64(0.0), pkg.Package.Rate.SmsMo)
	assert.Equal(t, float64(0.0), pkg.Package.Rate.SmsMt)
	assert.Equal(t, float64(0.0), pkg.Package.Rate.Data)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_GetDetails_SpecialCharactersInName(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	packageUUID := uuid.NewV4()

	req := &pb.GetPackageRequest{
		Uuid: packageUUID.String(),
	}

	specialNamePackage := &db.Package{
		Uuid:   packageUUID,
		Name:   "Package with special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?",
		Active: true,
		From:   time.Now(),
		To:     time.Now().AddDate(0, 1, 0),
	}

	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageRepo.On("GetDetails", packageUUID).Return(specialNamePackage, nil)

	pkg, err := s.GetDetails(context.TODO(), req)
	assert.NoError(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, "Package with special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?", pkg.Package.Name)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_GetDetails_ZeroVolumes(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	packageUUID := uuid.NewV4()

	req := &pb.GetPackageRequest{
		Uuid: packageUUID.String(),
	}

	zeroVolumePackage := &db.Package{
		Uuid:        packageUUID,
		Name:        "Zero Volume Package",
		Active:      true,
		SmsVolume:   0,
		DataVolume:  0,
		VoiceVolume: 0,
		From:        time.Now(),
		To:          time.Now().AddDate(0, 1, 0),
	}

	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageRepo.On("GetDetails", packageUUID).Return(zeroVolumePackage, nil)

	pkg, err := s.GetDetails(context.TODO(), req)
	assert.NoError(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, int64(0), pkg.Package.SmsVolume)
	assert.Equal(t, int64(0), pkg.Package.DataVolume)
	assert.Equal(t, int64(0), pkg.Package.VoiceVolume)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_GetDetails_ContextCancellation(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	packageUUID := uuid.NewV4()
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)

	req := &pb.GetPackageRequest{
		Uuid: packageUUID.String(),
	}

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Mock the GetDetails call to avoid panic
	packageRepo.On("GetDetails", packageUUID).Return(nil, context.Canceled)

	// The method should handle context cancellation gracefully
	// The actual behavior depends on the gRPC framework implementation
	// We're testing that it doesn't panic
	assert.NotPanics(t, func() {
		_, _ = s.GetDetails(ctx, req)
	})
}

func TestPackageServer_GetDetails_RecordNotFound(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	packageUUID := uuid.NewV4()
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)

	req := &pb.GetPackageRequest{
		Uuid: packageUUID.String(),
	}

	// Mock record not found error
	packageRepo.On("GetDetails", packageUUID).Return(nil, gorm.ErrRecordNotFound)

	pkg, err := s.GetDetails(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, pkg)
	assert.Contains(t, err.Error(), "record not found")
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_GetDetails_DifferentSimTypes(t *testing.T) {
	testCases := []struct {
		name     string
		simType  string
		expected string
	}{
		{"ukama_data", "ukama_data", "ukama_data"},
		{"operator_data", "operator_data", "operator_data"},
		{"test", "test", "test"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			packageRepo := &mocks.PackageRepo{}
			packageUUID := uuid.NewV4()
			s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)

			req := &pb.GetPackageRequest{
				Uuid: packageUUID.String(),
			}

			testPackage := &db.Package{
				Uuid:    packageUUID,
				Name:    tc.name,
				SimType: ukama.ParseSimType(tc.simType),
				Active:  true,
				From:    time.Now(),
				To:      time.Now().AddDate(0, 1, 0),
			}

			packageRepo.On("GetDetails", packageUUID).Return(testPackage, nil)

			pkg, err := s.GetDetails(context.TODO(), req)
			assert.NoError(t, err)
			assert.NotNil(t, pkg)
			assert.Equal(t, tc.expected, pkg.Package.SimType)
			packageRepo.AssertExpectations(t)
		})
	}
}

func TestPackageServer_GetAll_Success(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}

	req := &pb.GetAllRequest{}

	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packages := []db.Package{
		{Name: "Package 1"},
		{Name: "Package 2"},
	}
	packageRepo.On("GetAll").Return(packages, nil)

	resp, err := s.GetAll(context.TODO(), req)
	assert.NoError(t, err)
	assert.Len(t, resp.Packages, 2)
	assert.Equal(t, "Package 1", resp.Packages[0].Name)
	assert.Equal(t, "Package 2", resp.Packages[1].Name)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_Add_InvalidOwnerUUID(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:    "invalid-uuid",
		BaserateId: uuid.NewV4().String(),
		From:       time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:         time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	}

	resp, err := s.Add(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid format of owner uuid")
}

func TestPackageServer_Add_InvalidBaserateUUID(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:    uuid.NewV4().String(),
		BaserateId: "invalid-uuid",
		From:       time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:         time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	}

	resp, err := s.Add(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid format of base rate")
}

func TestPackageServer_Add_InvalidFromDate(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:    uuid.NewV4().String(),
		BaserateId: uuid.NewV4().String(),
		From:       "invalid-date",
		To:         time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	}

	resp, err := s.Add(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Error:")
}

func TestPackageServer_Add_InvalidToDate(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:    uuid.NewV4().String(),
		BaserateId: uuid.NewV4().String(),
		From:       time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:         "invalid-date",
	}

	resp, err := s.Add(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Error:")
}

func TestPackageServer_Add_PastFromDate(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:    uuid.NewV4().String(),
		BaserateId: uuid.NewV4().String(),
		From:       time.Now().Add(-time.Hour * 24).Format(time.RFC3339), // Past date
		To:         time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	}

	resp, err := s.Add(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Error:")
}

func TestPackageServer_Add_PastToDate(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:    uuid.NewV4().String(),
		BaserateId: uuid.NewV4().String(),
		From:       time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:         time.Now().Add(-time.Hour * 24).Format(time.RFC3339), // Past date
	}

	resp, err := s.Add(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Error:")
}

func TestPackageServer_Add_InvalidDateRange(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:    uuid.NewV4().String(),
		BaserateId: uuid.NewV4().String(),
		From:       time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339), // Later date
		To:         time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339), // Earlier date
	}

	resp, err := s.Add(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Error:")
}

func TestPackageServer_Add_RateClientError(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:    uuid.NewV4().String(),
		BaserateId: uuid.NewV4().String(),
		From:       time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:         time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	}

	rate.On("GetClient").Return(nil, errors.New("rate client error"))

	resp, err := s.Add(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "rate client error")
}

func TestPackageServer_Add_GetRateByIdError(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}
	ownerId := uuid.NewV4().String()
	baserate := uuid.NewV4().String()

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:    ownerId,
		BaserateId: baserate,
		From:       time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:         time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	}

	rateClient := &splmocks.RateServiceClient{}
	rate.On("GetClient").Return(rateClient, nil)
	rateClient.On("GetRateById", mock.Anything, &rpb.GetRateByIdRequest{
		OwnerId:  ownerId,
		BaseRate: baserate,
	}).Return(nil, errors.New("rate not found"))

	resp, err := s.Add(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid base id")
}

func TestPackageServer_Add_PackageRepoError(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}
	ownerId := uuid.NewV4().String()
	baserate := uuid.NewV4().String()

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:    ownerId,
		BaserateId: baserate,
		From:       time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:         time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	}

	rateClient := &splmocks.RateServiceClient{}
	rate.On("GetClient").Return(rateClient, nil)
	rateClient.On("GetRateById", mock.Anything, &rpb.GetRateByIdRequest{
		OwnerId:  ownerId,
		BaseRate: baserate,
	}).Return(&rpb.GetRateByIdResponse{
		Rate: &bpb.Rate{
			SmsMo:    1,
			SmsMt:    1,
			Data:     10,
			Country:  "USA",
			Provider: "ukama",
		},
	}, nil)

	packageRepo.On("Add", mock.Anything, mock.Anything).Return(errors.New("package repo error"))

	resp, err := s.Add(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Error while adding a package")
}

func TestPackageServer_Delete_Success(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	packageUUID := uuid.NewV4()

	s := NewPackageServer(OrgName, packageRepo, nil, msgbusClient, OrgId)

	req := &pb.DeletePackageRequest{
		Uuid: packageUUID.String(),
	}

	packageRepo.On("Delete", packageUUID).Return(nil)
	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

	resp, err := s.Delete(context.TODO(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, packageUUID.String(), resp.Uuid)
	packageRepo.AssertExpectations(t)
	msgbusClient.AssertExpectations(t)
}

func TestPackageServer_Delete_InvalidUUID(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)

	req := &pb.DeletePackageRequest{
		Uuid: "invalid-uuid",
	}

	resp, err := s.Delete(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid format of package uuid")
}

func TestPackageServer_Update_InvalidUUID(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)

	req := &pb.UpdatePackageRequest{
		Uuid: "invalid-uuid",
		Name: "Updated Package",
	}

	resp, err := s.Update(context.TODO(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid format of package uuid")
}

func TestPackageServer_Add_FlatratePackage(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}
	ownerId := uuid.NewV4().String()
	baserate := uuid.NewV4().String()

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:    ownerId,
		BaserateId: baserate,
		Flatrate:   true, // Test flatrate package
		From:       time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:         time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	}

	rateClient := &splmocks.RateServiceClient{}
	rate.On("GetClient").Return(rateClient, nil)
	rateClient.On("GetRateById", mock.Anything, &rpb.GetRateByIdRequest{
		OwnerId:  ownerId,
		BaseRate: baserate,
	}).Return(&rpb.GetRateByIdResponse{
		Rate: &bpb.Rate{
			SmsMo:    1,
			SmsMt:    1,
			Data:     10,
			Country:  "USA",
			Provider: "ukama",
		},
	}, nil)

	packageRepo.On("Add", mock.Anything, mock.Anything).Return(nil)

	resp, err := s.Add(context.TODO(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Package.Flatrate)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_Add_WithMessageBus(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	ownerId := uuid.NewV4().String()
	baserate := uuid.NewV4().String()

	s := NewPackageServer(OrgName, packageRepo, rate, msgbusClient, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:       ownerId,
		BaserateId:    baserate,
		Name:          "test-package",
		SimType:       "ukama_data",
		Active:        true,
		Duration:      30,
		SmsVolume:     100,
		DataVolume:    1024,
		VoiceVolume:   60,
		MessageUnit:   "message",
		VoiceUnit:     "minute",
		DataUnit:      "mb",
		Type:          "prepaid",
		Flatrate:      false,
		Amount:        50.0,
		Markup:        10.0,
		Currency:      "USD",
		Overdraft:     5.0,
		TrafficPolicy: 1,
		Networks:      []string{"network1", "network2"},
		Country:       "USA",
		Provider:      "ukama",
		Apn:           "test.apn",
		From:          time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:            time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	}

	rateClient := &splmocks.RateServiceClient{}
	rate.On("GetClient").Return(rateClient, nil)
	rateClient.On("GetRateById", mock.Anything, &rpb.GetRateByIdRequest{
		OwnerId:  ownerId,
		BaseRate: baserate,
	}).Return(&rpb.GetRateByIdResponse{
		Rate: &bpb.Rate{
			SmsMo:    0.1,
			SmsMt:    0.1,
			Data:     0.5,
			Country:  "USA",
			Provider: "ukama",
			Apn:      "default.apn",
		},
	}, nil)

	packageRepo.On("Add", mock.Anything, mock.Anything).Return(nil)
	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

	resp, err := s.Add(context.TODO(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Name, resp.Package.Name)
	assert.Equal(t, req.SimType, resp.Package.SimType)
	assert.Equal(t, req.Active, resp.Package.Active)
	assert.Equal(t, req.Flatrate, resp.Package.Flatrate)
	assert.Equal(t, req.Currency, resp.Package.Currency)
	assert.Equal(t, req.Country, resp.Package.Country)
	assert.Equal(t, req.Provider, resp.Package.Provider)
	assert.Equal(t, req.Apn, resp.Package.Apn)
	assert.Equal(t, req.Overdraft, resp.Package.Overdraft)
	assert.Equal(t, req.TrafficPolicy, resp.Package.TrafficPolicy)
	assert.Equal(t, req.Networks, resp.Package.Networks)

	packageRepo.AssertExpectations(t)
	msgbusClient.AssertExpectations(t)
}

func TestPackageServer_Add_MessageBusPublishError(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	ownerId := uuid.NewV4().String()
	baserate := uuid.NewV4().String()

	s := NewPackageServer(OrgName, packageRepo, rate, msgbusClient, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:    ownerId,
		BaserateId: baserate,
		Name:       "test-package",
		SimType:    "ukama_data",
		Active:     true,
		From:       time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:         time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	}

	rateClient := &splmocks.RateServiceClient{}
	rate.On("GetClient").Return(rateClient, nil)
	rateClient.On("GetRateById", mock.Anything, &rpb.GetRateByIdRequest{
		OwnerId:  ownerId,
		BaseRate: baserate,
	}).Return(&rpb.GetRateByIdResponse{
		Rate: &bpb.Rate{
			SmsMo:    0.1,
			SmsMt:    0.1,
			Data:     0.5,
			Country:  "USA",
			Provider: "ukama",
		},
	}, nil)

	packageRepo.On("Add", mock.Anything, mock.Anything).Return(nil)
	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(errors.New("message bus error"))

	resp, err := s.Add(context.TODO(), req)
	// Should still succeed even if message bus fails
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	packageRepo.AssertExpectations(t)
	msgbusClient.AssertExpectations(t)
}

func TestPackageServer_Add_WithEmptyApn(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}
	ownerId := uuid.NewV4().String()
	baserate := uuid.NewV4().String()

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:    ownerId,
		BaserateId: baserate,
		Name:       "test-package",
		SimType:    "ukama_data",
		Active:     true,
		Apn:        "", // Empty APN should use rate's APN
		From:       time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:         time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	}

	rateClient := &splmocks.RateServiceClient{}
	rate.On("GetClient").Return(rateClient, nil)
	rateClient.On("GetRateById", mock.Anything, &rpb.GetRateByIdRequest{
		OwnerId:  ownerId,
		BaseRate: baserate,
	}).Return(&rpb.GetRateByIdResponse{
		Rate: &bpb.Rate{
			SmsMo:    0.1,
			SmsMt:    0.1,
			Data:     0.5,
			Country:  "USA",
			Provider: "ukama",
			Apn:      "default.apn",
		},
	}, nil)

	packageRepo.On("Add", mock.Anything, mock.Anything).Return(nil)

	resp, err := s.Add(context.TODO(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "default.apn", resp.Package.Apn)

	packageRepo.AssertExpectations(t)
}

func TestPackageServer_Add_NonFlatratePackageWithCalculations(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}
	ownerId := uuid.NewV4().String()
	baserate := uuid.NewV4().String()

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:     ownerId,
		BaserateId:  baserate,
		Name:        "test-package",
		SimType:     "ukama_data",
		Active:      true,
		Duration:    30,
		SmsVolume:   100,
		DataVolume:  1024,
		VoiceVolume: 60,
		MessageUnit: "int",
		VoiceUnit:   "minutes",
		DataUnit:    "MegaBytes",
		Type:        "prepaid",
		Flatrate:    false, // Non-flatrate should calculate rates
		Amount:      0.0,   // Should be calculated
		From:        time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:          time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	}

	rateClient := &splmocks.RateServiceClient{}
	rate.On("GetClient").Return(rateClient, nil)
	rateClient.On("GetRateById", mock.Anything, &rpb.GetRateByIdRequest{
		OwnerId:  ownerId,
		BaseRate: baserate,
	}).Return(&rpb.GetRateByIdResponse{
		Rate: &bpb.Rate{
			SmsMo:    0.1,
			SmsMt:    0.1,
			Data:     0.5,
			Country:  "USA",
			Provider: "ukama",
		},
	}, nil)

	packageRepo.On("Add", mock.Anything, mock.Anything).Return(nil)

	resp, err := s.Add(context.TODO(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Package.Flatrate)
	// Amount should be calculated: (0.1 + 0.1) * 100 + 0.5 * 1024 = 20 + 512 = 532
	assert.Greater(t, resp.Package.Amount, 0.0)

	packageRepo.AssertExpectations(t)
}

func TestPackageServer_Add_WithDifferentSimTypes(t *testing.T) {
	testCases := []struct {
		name     string
		simType  string
		expected bool
	}{
		{"ukama_data", "ukama_data", true},
		{"test", "test", true},
		{"operator_data", "operator_data", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			packageRepo := &mocks.PackageRepo{}
			rate := &mocks.RateClientProvider{}
			ownerId := uuid.NewV4().String()
			baserate := uuid.NewV4().String()

			s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

			req := &pb.AddPackageRequest{
				OwnerId:    ownerId,
				BaserateId: baserate,
				Name:       "test-package",
				SimType:    tc.simType,
				Active:     true,
				From:       time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
				To:         time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
			}

			rateClient := &splmocks.RateServiceClient{}
			rate.On("GetClient").Return(rateClient, nil)
			rateClient.On("GetRateById", mock.Anything, &rpb.GetRateByIdRequest{
				OwnerId:  ownerId,
				BaseRate: baserate,
			}).Return(&rpb.GetRateByIdResponse{
				Rate: &bpb.Rate{
					SmsMo:    0.1,
					SmsMt:    0.1,
					Data:     0.5,
					Country:  "USA",
					Provider: "ukama",
				},
			}, nil)

			packageRepo.On("Add", mock.Anything, mock.Anything).Return(nil)

			resp, err := s.Add(context.TODO(), req)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tc.simType, resp.Package.SimType)

			packageRepo.AssertExpectations(t)
		})
	}
}

func TestPackageServer_Add_WithDifferentPackageTypes(t *testing.T) {
	testCases := []struct {
		name     string
		pkgType  string
		expected bool
	}{
		{"prepaid", "prepaid", true},
		{"postpaid", "postpaid", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			packageRepo := &mocks.PackageRepo{}
			rate := &mocks.RateClientProvider{}
			ownerId := uuid.NewV4().String()
			baserate := uuid.NewV4().String()

			s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

			req := &pb.AddPackageRequest{
				OwnerId:    ownerId,
				BaserateId: baserate,
				Name:       "test-package",
				SimType:    "ukama_data",
				Type:       tc.pkgType,
				Active:     true,
				From:       time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
				To:         time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
			}

			rateClient := &splmocks.RateServiceClient{}
			rate.On("GetClient").Return(rateClient, nil)
			rateClient.On("GetRateById", mock.Anything, &rpb.GetRateByIdRequest{
				OwnerId:  ownerId,
				BaseRate: baserate,
			}).Return(&rpb.GetRateByIdResponse{
				Rate: &bpb.Rate{
					SmsMo:    0.1,
					SmsMt:    0.1,
					Data:     0.5,
					Country:  "USA",
					Provider: "ukama",
				},
			}, nil)

			packageRepo.On("Add", mock.Anything, mock.Anything).Return(nil)

			resp, err := s.Add(context.TODO(), req)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tc.pkgType, resp.Package.Type)

			packageRepo.AssertExpectations(t)
		})
	}
}

func TestPackageServer_Add_WithDifferentUnits(t *testing.T) {
	testCases := []struct {
		name        string
		messageUnit string
		voiceUnit   string
		dataUnit    string
	}{
		{"int_minutes_MegaBytes", "int", "minutes", "MegaBytes"},
		{"int_minutes_GigaBytes", "int", "minutes", "GigaBytes"},
		{"int_hours_MegaBytes", "int", "hours", "MegaBytes"},
		{"int_hours_GigaBytes", "int", "hours", "GigaBytes"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			packageRepo := &mocks.PackageRepo{}
			rate := &mocks.RateClientProvider{}
			ownerId := uuid.NewV4().String()
			baserate := uuid.NewV4().String()

			s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

			req := &pb.AddPackageRequest{
				OwnerId:     ownerId,
				BaserateId:  baserate,
				Name:        "test-package",
				SimType:     "ukama_data",
				MessageUnit: tc.messageUnit,
				VoiceUnit:   tc.voiceUnit,
				DataUnit:    tc.dataUnit,
				Active:      true,
				From:        time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
				To:          time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
			}

			rateClient := &splmocks.RateServiceClient{}
			rate.On("GetClient").Return(rateClient, nil)
			rateClient.On("GetRateById", mock.Anything, &rpb.GetRateByIdRequest{
				OwnerId:  ownerId,
				BaseRate: baserate,
			}).Return(&rpb.GetRateByIdResponse{
				Rate: &bpb.Rate{
					SmsMo:    0.1,
					SmsMt:    0.1,
					Data:     0.5,
					Country:  "USA",
					Provider: "ukama",
				},
			}, nil)

			packageRepo.On("Add", mock.Anything, mock.Anything).Return(nil)

			resp, err := s.Add(context.TODO(), req)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tc.messageUnit, resp.Package.MessageUnit)
			assert.Equal(t, tc.voiceUnit, resp.Package.VoiceUnit)
			assert.Equal(t, tc.dataUnit, resp.Package.DataUnit)

			packageRepo.AssertExpectations(t)
		})
	}
}

func TestPackageServer_Add_WithZeroVolumes(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}
	ownerId := uuid.NewV4().String()
	baserate := uuid.NewV4().String()

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:     ownerId,
		BaserateId:  baserate,
		Name:        "test-package",
		SimType:     "ukama_data",
		Active:      true,
		SmsVolume:   0,
		DataVolume:  0,
		VoiceVolume: 0,
		Flatrate:    false,
		From:        time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:          time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	}

	rateClient := &splmocks.RateServiceClient{}
	rate.On("GetClient").Return(rateClient, nil)
	rateClient.On("GetRateById", mock.Anything, &rpb.GetRateByIdRequest{
		OwnerId:  ownerId,
		BaseRate: baserate,
	}).Return(&rpb.GetRateByIdResponse{
		Rate: &bpb.Rate{
			SmsMo:    0.1,
			SmsMt:    0.1,
			Data:     0.5,
			Country:  "USA",
			Provider: "ukama",
		},
	}, nil)

	packageRepo.On("Add", mock.Anything, mock.Anything).Return(nil)

	resp, err := s.Add(context.TODO(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.Package.SmsVolume)
	assert.Equal(t, int64(0), resp.Package.DataVolume)
	assert.Equal(t, int64(0), resp.Package.VoiceVolume)
	// Amount should be 0 since volumes are 0
	assert.Equal(t, 0.0, resp.Package.Amount)

	packageRepo.AssertExpectations(t)
}

func TestPackageServer_Add_WithLargeVolumes(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}
	ownerId := uuid.NewV4().String()
	baserate := uuid.NewV4().String()

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:     ownerId,
		BaserateId:  baserate,
		Name:        "test-package",
		SimType:     "ukama_data",
		Active:      true,
		SmsVolume:   10000,
		DataVolume:  100000,
		VoiceVolume: 1000,
		MessageUnit: "int",
		VoiceUnit:   "minutes",
		DataUnit:    "MegaBytes",
		Flatrate:    false,
		From:        time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:          time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	}

	rateClient := &splmocks.RateServiceClient{}
	rate.On("GetClient").Return(rateClient, nil)
	rateClient.On("GetRateById", mock.Anything, &rpb.GetRateByIdRequest{
		OwnerId:  ownerId,
		BaseRate: baserate,
	}).Return(&rpb.GetRateByIdResponse{
		Rate: &bpb.Rate{
			SmsMo:    0.01,
			SmsMt:    0.01,
			Data:     0.1,
			Country:  "USA",
			Provider: "ukama",
		},
	}, nil)

	packageRepo.On("Add", mock.Anything, mock.Anything).Return(nil)

	resp, err := s.Add(context.TODO(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(10000), resp.Package.SmsVolume)
	assert.Equal(t, int64(100000), resp.Package.DataVolume)
	assert.Equal(t, int64(1000), resp.Package.VoiceVolume)
	// Amount should be calculated: (0.01 + 0.01) * 10000 + 0.1 * 100000 = 200 + 10000 = 10200
	assert.Greater(t, resp.Package.Amount, 0.0)

	packageRepo.AssertExpectations(t)
}

func TestPackageServer_Add_WithNegativeMarkup(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateClientProvider{}
	ownerId := uuid.NewV4().String()
	baserate := uuid.NewV4().String()

	s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

	req := &pb.AddPackageRequest{
		OwnerId:    ownerId,
		BaserateId: baserate,
		Name:       "test-package",
		SimType:    "ukama_data",
		Active:     true,
		Markup:     -10.0, // Negative markup
		From:       time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:         time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	}

	rateClient := &splmocks.RateServiceClient{}
	rate.On("GetClient").Return(rateClient, nil)
	rateClient.On("GetRateById", mock.Anything, &rpb.GetRateByIdRequest{
		OwnerId:  ownerId,
		BaseRate: baserate,
	}).Return(&rpb.GetRateByIdResponse{
		Rate: &bpb.Rate{
			SmsMo:    0.1,
			SmsMt:    0.1,
			Data:     0.5,
			Country:  "USA",
			Provider: "ukama",
		},
	}, nil)

	packageRepo.On("Add", mock.Anything, mock.Anything).Return(nil)

	resp, err := s.Add(context.TODO(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, -10.0, resp.Package.Markup.Markup)

	packageRepo.AssertExpectations(t)
}

func TestPackageServer_Update_RepoError(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	s := NewPackageServer(OrgName, packageRepo, nil, msgbusClient, OrgId)
	packageUUID := uuid.NewV4()
	packageRepo.On("Update", packageUUID, mock.Anything).Return(errors.New("db error")).Once()
	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Maybe()
	resp, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
		Uuid:   packageUUID.String(),
		Name:   "fail-update",
		Active: true,
	})
	assert.Error(t, err)
	assert.Nil(t, resp)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_Update_MessageBusError(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	s := NewPackageServer(OrgName, packageRepo, nil, msgbusClient, OrgId)
	packageUUID := uuid.NewV4()
	packageRepo.On("Update", packageUUID, mock.Anything).Return(nil).Once()

	packageRepo.On("Get", packageUUID).Return(&db.Package{
		Uuid:   packageUUID,
		Name:   "msgbus-fail",
		Active: true,
		From:   time.Now(),
		To:     time.Now().Add(time.Hour * 24 * 30),
	}, nil).Once()

	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(errors.New("msgbus error")).Once()
	resp, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
		Uuid:   packageUUID.String(),
		Name:   "msgbus-fail",
		Active: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "msgbus-fail", resp.Package.Name)
	packageRepo.AssertExpectations(t)
	msgbusClient.AssertExpectations(t)
}

func TestPackageServer_Update_OnlyName(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageUUID := uuid.NewV4()

	packageRepo.On("Update", packageUUID, mock.MatchedBy(func(p *db.Package) bool {
		return p.Name == "NameOnly" && p.Active == false
	})).Return(nil).Once()

	packageRepo.On("Get", packageUUID).Return(&db.Package{
		Uuid:   packageUUID,
		Name:   "NameOnly",
		Active: false,
		From:   time.Now(),
		To:     time.Now().Add(time.Hour * 24 * 30),
	}, nil).Once()

	resp, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
		Uuid: packageUUID.String(),
		Name: "NameOnly",
	})
	assert.NoError(t, err)
	assert.Equal(t, "NameOnly", resp.Package.Name)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_Update_OnlyActive(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageUUID := uuid.NewV4()

	packageRepo.On("Update", packageUUID, mock.MatchedBy(func(p *db.Package) bool {
		return p.Name == "" && p.Active == true
	})).Return(nil).Once()

	packageRepo.On("Get", packageUUID).Return(&db.Package{
		Uuid:   packageUUID,
		Name:   "Original Name",
		Active: true,
		From:   time.Now(),
		To:     time.Now().Add(time.Hour * 24 * 30),
	}, nil).Once()

	resp, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
		Uuid:   packageUUID.String(),
		Active: true,
	})
	assert.NoError(t, err)
	assert.True(t, resp.Package.Active)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_Update_Inactive(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageUUID := uuid.NewV4()

	packageRepo.On("Update", packageUUID, mock.MatchedBy(func(p *db.Package) bool {
		return p.Active == false
	})).Return(nil).Once()

	packageRepo.On("Get", packageUUID).Return(&db.Package{
		Uuid:   packageUUID,
		Name:   "Inactive Package",
		Active: false,
		From:   time.Now(),
		To:     time.Now().Add(time.Hour * 24 * 30),
	}, nil).Once()

	resp, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
		Uuid:   packageUUID.String(),
		Active: false,
	})
	assert.NoError(t, err)
	assert.False(t, resp.Package.Active)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_Update_NilMsgBus(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageUUID := uuid.NewV4()

	packageRepo.On("Update", packageUUID, mock.Anything).Return(nil).Once()

	packageRepo.On("Get", packageUUID).Return(&db.Package{
		Uuid:   packageUUID,
		Name:   "NoMsgBus",
		Active: true,
		From:   time.Now(),
		To:     time.Now().Add(time.Hour * 24 * 30),
	}, nil).Once()

	resp, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
		Uuid:   packageUUID.String(),
		Name:   "NoMsgBus",
		Active: true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "NoMsgBus", resp.Package.Name)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_Update_GetAfterUpdateError(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageUUID := uuid.NewV4()

	packageRepo.On("Update", packageUUID, mock.Anything).Return(nil).Once()

	packageRepo.On("Get", packageUUID).Return(nil, errors.New("get error")).Once()

	resp, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
		Uuid:   packageUUID.String(),
		Name:   "GetFail",
		Active: true,
	})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "get error")
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_Update_EmptyName(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
	packageUUID := uuid.NewV4()

	packageRepo.On("Update", packageUUID, mock.MatchedBy(func(p *db.Package) bool {
		return p.Name == "" && p.Active == true
	})).Return(nil).Once()

	packageRepo.On("Get", packageUUID).Return(&db.Package{
		Uuid:   packageUUID,
		Name:   "",
		Active: true,
		From:   time.Now(),
		To:     time.Now().Add(time.Hour * 24 * 30),
	}, nil).Once()

	resp, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
		Uuid:   packageUUID.String(),
		Name:   "",
		Active: true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "", resp.Package.Name)
	packageRepo.AssertExpectations(t)
}
