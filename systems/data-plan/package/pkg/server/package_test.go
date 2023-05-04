package server

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/grpc"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	"github.com/ukama/ukama/systems/data-plan/package/mocks"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"
	rpb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	s := NewPackageServer(packageRepo, nil, nil)
	packageRepo.On("Get", packageUUID).Return(&db.Package{
		Name: "Daily-pack",
	}, nil)
	pkg, err := s.Get(context.TODO(), mockFilters)
	assert.NoError(t, err)
	assert.Equal(t, "Daily-pack", pkg.Package.Name)
	packageRepo.AssertExpectations(t)
}

// Error case SQL error
func TestPackageServer_GetPackages_Error1(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo, nil, nil)
	packageUUID := uuid.NewV4()
	var mockFilters = &pb.GetPackageRequest{
		Uuid: packageUUID.String(),
	}
	packageRepo.On("Get", packageUUID).Return(nil, grpc.SqlErrorToGrpc(errors.New("SQL error while fetching records"), "packages"))
	pkg2, err := s.Get(context.TODO(), mockFilters)
	assert.Error(t, err)
	assert.Nil(t, pkg2)
}

func TestPackageServer_GetPackageByOrg_Success(t *testing.T) {
	var orgId = uuid.NewV4()

	packageRepo := &mocks.PackageRepo{}
	var mockFilters = &pb.GetByOrgPackageRequest{
		OrgId: orgId.String(),
	}
	s := NewPackageServer(packageRepo, nil, nil)

	packageRepo.On("GetByOrg", orgId).Return([]db.Package{{
		Uuid:        uuid.NewV4(),
		Name:        "Silver Plan",
		SimType:     ukama.SimTypeTest,
		OrgId:       orgId,
		OwnerId:     uuid.NewV4(),
		Active:      true,
		Duration:    30,
		SmsVolume:   1000,
		DataVolume:  5000000,
		VoiceVolume: 500,
		Type:        ukama.PackageTypePostpaid,

		DataUnits:    ukama.DataUnitTypeMB,
		VoiceUnits:   ukama.CallUnitTypeSec,
		MessageUnits: ukama.MessageUnitTypeInt,
		Flatrate:     false,
		Currency:     "Dollar",
		From:         time.Now(),
		To:           time.Now().Add(time.Hour * 24 * 30),
		Country:      "USA",
		Provider:     "ukama",
	}}, nil)
	pkg, err := s.GetByOrg(context.TODO(), mockFilters)
	assert.NoError(t, err)
	assert.Equal(t, int64(5000000), pkg.Packages[0].DataVolume)
	packageRepo.AssertExpectations(t)
}

// Error cases

func TestPackageServer_GetPackageByOrg_Error(t *testing.T) {
	var orgId = uuid.NewV4()

	packageRepo := &mocks.PackageRepo{}
	var mockFilters = &pb.GetByOrgPackageRequest{
		OrgId: orgId.String(),
	}
	s := NewPackageServer(packageRepo, nil, nil)

	packageRepo.On("GetByOrg", orgId).
		Return(nil, grpc.SqlErrorToGrpc(errors.New("SQL error while fetching records"), "packages"))
	pkg, err := s.GetByOrg(context.TODO(), mockFilters)
	assert.Error(t, err)
	assert.Nil(t, pkg)
	packageRepo.AssertExpectations(t)
}

// End Get package org //

// Add packages //
func TestPackageServer_AddPackage(t *testing.T) {
	ownerId := uuid.NewV4().String()
	baserate := uuid.NewV4().String()
	packageRepo := &mocks.PackageRepo{}
	rate := &mocks.RateService{}
	packageRepo.On("Add", mock.MatchedBy(func(p *db.Package) bool {
		return p.Active == true && p.Name == "daily-pack"
	})).Return(nil).Once()

	rate.On("GetRateById", &rpb.GetRateByIdRequest{
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
	s := NewPackageServer(packageRepo, rate, nil)

	ActPackage, err := s.Add(context.TODO(), &pb.AddPackageRequest{
		Active:   true,
		Name:     "daily-pack",
		OrgId:    uuid.NewV4().String(),
		SimType:  testSim,
		OwnerId:  ownerId,
		Baserate: baserate,
		From:     time.Now().Add(time.Hour * 24 * 30).Format(time.RFC3339),
		To:       time.Now().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, ActPackage.Package.Active)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_UpdatePackage(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	s := NewPackageServer(packageRepo, nil, msgbusClient)
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
	s := NewPackageServer(packageRepo, nil, nil)
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

// Error case: Id 0
func TestPackageServer_DeletePackage_Success_Error2(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo, nil, nil)
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
