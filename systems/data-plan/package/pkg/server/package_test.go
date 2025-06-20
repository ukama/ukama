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
