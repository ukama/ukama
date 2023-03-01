package server

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"testing"

// 	uuid "github.com/ukama/ukama/systems/common/uuid"
// 	"gorm.io/gorm"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/ukama/ukama/systems/common/grpc"
// 	"github.com/ukama/ukama/systems/data-plan/package/mocks"
// 	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
// 	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// // Get Packages //
// // Success case
// type UkamaDbMock struct {
// 	GormDb *gorm.DB
// }

// func (u UkamaDbMock) Init(model ...interface{}) error {
// 	panic("implement me: Init()")
// }

// func (u UkamaDbMock) Connect() error {
// 	panic("implement me: Connect()")
// }

// func (u UkamaDbMock) GetGormDb() *gorm.DB {
// 	return u.GormDb
// }

// func (u UkamaDbMock) InitDB() error {
// 	return nil
// }

// func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB,
// 	nestedFuncs ...func() error) error {
// 	log.Fatal("implement me: ExecuteInTransaction()")
// 	return nil
// }

// func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB,
// 	nestedFuncs ...func(tx *gorm.DB) error) error {
// 	log.Fatal("implement me: ExecuteInTransaction2()")
// 	return nil
// }

// func TestPackageServer_GetPackages_Success(t *testing.T) {
// 	packageRepo := &mocks.PackageRepo{}
// 	packageUUID := uuid.NewV4()

// 	var mockFilters = &pb.GetPackageRequest{
// 		PackageID: packageUUID.String(),
// 	}
// 	s := NewPackageServer(packageRepo,nil)
// 	packageRepo.On("Get", packageUUID).Return(&db.Package{
// 		Name: "Daily-pack",
// 	}, nil)
// 	pkg, err := s.Get(context.TODO(), mockFilters)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "Daily-pack", pkg.Package.Name)
// 	packageRepo.AssertExpectations(t)
// }

// // Error case SQL error
// func TestPackageServer_GetPackages_Error1(t *testing.T) {
// 	packageRepo := &mocks.PackageRepo{}
// 	s := NewPackageServer(packageRepo,nil)
// 	packageUUID := uuid.NewV4()
// 	var mockFilters = &pb.GetPackageRequest{
// 		PackageID: packageUUID.String(),
// 	}
// 	packageRepo.On("Get", packageUUID).Return(nil, grpc.SqlErrorToGrpc(errors.New("SQL error while fetching records"), "packages"))
// 	pkg2, err := s.Get(context.TODO(), mockFilters)
// 	assert.Error(t, err)
// 	assert.Nil(t, pkg2)
// }

// // End Get packages //

// // Get Package by org //

// func TestPackageServer_GetPackageByOrg_Success(t *testing.T) {
// 	var orgID = uuid.NewV4()

// 	packageRepo := &mocks.PackageRepo{}
// 	var mockFilters = &pb.GetByOrgPackageRequest{
// 		OrgID: orgID.String(),
// 	}
// 	s := NewPackageServer(packageRepo,nil)

// 	packageRepo.On("GetByOrg", orgID).Return([]db.Package{{
// 		SimType:     db.ParseType("INTER_MNO_ALL"),
// 		Name:         "Daily-pack",
// 		OrgID:       orgID,
// 		Active:       true,
// 		Duration:     1,
// 		SmsVolume:   20,
// 		DataVolume:  12,
// 		VoiceVolume: 34,
// 		OrgRatesID: 00,
// 	}}, nil)
// 	pkg, err := s.GetByOrg(context.TODO(), mockFilters)
// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(12), pkg.Packages[0].DataVolume)
// 	packageRepo.AssertExpectations(t)
// }

// // Error cases

// func TestPackageServer_GetPackageByOrg_Error(t *testing.T) {
// 	var orgID = uuid.NewV4()

// 	packageRepo := &mocks.PackageRepo{}
// 	var mockFilters = &pb.GetByOrgPackageRequest{
// 		OrgID: orgID.String(),
// 	}
// 	s := NewPackageServer(packageRepo,nil)

// 	packageRepo.On("GetByOrg", orgID).
// 		Return(nil, grpc.SqlErrorToGrpc(errors.New("SQL error while fetching records"), "packages"))
// 	pkg, err := s.GetByOrg(context.TODO(), mockFilters)
// 	assert.Error(t, err)
// 	assert.Nil(t, pkg)
// 	packageRepo.AssertExpectations(t)
// }

// // End Get package org //

// // Add packages //
// func TestPackageServer_AddPackage(t *testing.T) {
// 	packageRepo := &mocks.PackageRepo{}
// 	packageRepo.On("Add", mock.MatchedBy(func(p *db.Package) bool {
// 		return p.Active == true && p.Name == "daily-pack"
// 	})).Return(nil).Once()

// 	s := NewPackageServer(packageRepo,nil)

// 	ActPackage, err := s.Add(context.TODO(), &pb.AddPackageRequest{
// 		Active: true,
// 		Name:   "daily-pack",
// 		OrgID:  uuid.NewV4().String(),
// 	})
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, ActPackage.Package.Active)
// 	packageRepo.AssertExpectations(t)
// }

// // Error case adding package
// func TestPackageServer_AddPackage_Error(t *testing.T) {
// 	packageRepo := &mocks.PackageRepo{}
// 	packageRepo.On("Add", mock.MatchedBy(func(p *db.Package) bool {
// 		return p.Active == true && p.Name == "daily-pack"
// 	})).Return(status.Errorf(codes.Internal, "error adding a package"))

// 	s := NewPackageServer(packageRepo,nil)

// 	ActPackage, err := s.Add(context.TODO(), &pb.AddPackageRequest{
// 		Active: true,
// 		Name:   "daily-pack",
// 		OrgID:  uuid.NewV4().String(),
// 	})
// 	assert.Error(t, err)
// 	assert.Nil(t, ActPackage)
// 	packageRepo.AssertExpectations(t)
// }

// func TestPackageServer_UpdatePackage_Error(t *testing.T) {
// 	packageRepo := &mocks.PackageRepo{}
// 	s := NewPackageServer(packageRepo,nil)
// 	packageUUID := uuid.NewV4()
// 	packageRepo.On("Update", packageUUID, mock.Anything).Return(nil, grpc.SqlErrorToGrpc(errors.New("Error updating records"), "rates"))
// 	pkg, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
// 		PackageID: packageUUID.String(),
// 	})
// 	assert.Error(t, err)
// 	assert.Nil(t, pkg)
// }

// func TestPackageServer_DeletePackage_Error1(t *testing.T) {
// 	packageRepo := &mocks.PackageRepo{}
// 	s := NewPackageServer(packageRepo,nil)
// 	packageUUID := uuid.NewV4()
// 	var mockFilters = &pb.DeletePackageRequest{
// 		PackageID: packageUUID.String(),
// 	}
// 	packageRepo.On("Delete", packageUUID).
// 		Return(status.Errorf(codes.InvalidArgument, "OrgId is required."))
// 	pkg1, err := s.Delete(context.TODO(), mockFilters)
// 	assert.Error(t, err)
// 	assert.Nil(t, pkg1)
// }

// // Error case: Id 0
// func TestPackageServer_DeletePackage_Success_Error2(t *testing.T) {
// 	packageRepo := &mocks.PackageRepo{}
// 	s := NewPackageServer(packageRepo,nil)
// 	packageUUID := uuid.NewV4()
// 	var mockFilters = &pb.DeletePackageRequest{
// 		PackageID: packageUUID.String(),
// 	}
// 	packageRepo.On("Delete", packageUUID).
// 		Return(status.Errorf(codes.InvalidArgument, "Id is required."))
// 	pkg2, err := s.Delete(context.TODO(), mockFilters)
// 	assert.Error(t, err)
// 	assert.Nil(t, pkg2)
// }

// // Error case: Error deleting record
// func TestPackageServer_DeletePackage_Error3(t *testing.T) {
// 	packageRepo := &mocks.PackageRepo{}
// 	s := NewPackageServer(packageRepo,nil)
// 	packageUUID := uuid.NewV4()
// 	var mf = &pb.DeletePackageRequest{
// 		PackageID: packageUUID.String(),
// 	}
// 	packageRepo.On("Delete", packageUUID).
// 		Return(grpc.SqlErrorToGrpc(errors.New("SQL error while deleting record"), "packages"))
// 	pkg3, err := s.Delete(context.TODO(), mf)
// 	fmt.Println(err)
// 	assert.Error(t, err)
// 	assert.Nil(t, pkg3)
// }

// // End Delete package //

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/grpc"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/data-plan/package/mocks"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Get Packages //
// Success case
func TestPackageServer_GetPackages_Success(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	packageUUID := uuid.NewV4()

	var mockFilters = &pb.GetPackageRequest{
		Uuid: packageUUID.String(),
	}
	s := NewPackageServer(packageRepo, nil)
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
	s := NewPackageServer(packageRepo, nil)
	packageUUID := uuid.NewV4()
	var mockFilters = &pb.GetPackageRequest{
		Uuid: packageUUID.String(),
	}
	packageRepo.On("Get", packageUUID).Return(nil, grpc.SqlErrorToGrpc(errors.New("SQL error while fetching records"), "packages"))
	pkg2, err := s.Get(context.TODO(), mockFilters)
	assert.Error(t, err)
	assert.Nil(t, pkg2)
}

// End Get packages //

// Get Package by org //

func TestPackageServer_GetPackageByOrg_Success(t *testing.T) {
	var orgId = uuid.NewV4()

	packageRepo := &mocks.PackageRepo{}
	var mockFilters = &pb.GetByOrgPackageRequest{
		OrgID: orgId.String(),
	}
	s := NewPackageServer(packageRepo, nil)

	packageRepo.On("GetByOrg", orgId).Return([]db.Package{{
		SimType:     db.SimTypeTest,
		Name:        "Daily-pack",
		OrgID:       orgId,
		Active:      true,
		Duration:    1,
		SmsVolume:   20,
		DataVolume:  12,
		VoiceVolume: 34,
		OrgRatesID:  00,
	}}, nil)
	pkg, err := s.GetByOrg(context.TODO(), mockFilters)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), pkg.Packages[0].DataVolume)
	packageRepo.AssertExpectations(t)
}

// Error cases

func TestPackageServer_GetPackageByOrg_Error(t *testing.T) {
	var orgId = uuid.NewV4()

	packageRepo := &mocks.PackageRepo{}
	var mockFilters = &pb.GetByOrgPackageRequest{
		OrgID: orgId.String(),
	}
	s := NewPackageServer(packageRepo, nil)

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
	packageRepo := &mocks.PackageRepo{}
	reqMock := &pb.AddPackageRequest{
		SimType:     "ukama_data",
		Name:        "Daily-pack",
		OrgID:       uuid.NewV4().String(),
		Active:      true,
		Duration:    1,
		SmsVolume:   20,
		DataVolume:  12,
		VoiceVolume: 34,
	}
	packageRepo.On("Add", mock.Anything, reqMock).Return(nil)

	s := NewPackageServer(packageRepo, nil)

	ActPackage, err := s.Add(context.TODO(), &pb.AddPackageRequest{
		Active: true,
		Name:   "daily-pack",
		OrgID:  uuid.NewV4().String(),
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, ActPackage.Package.Active)
	packageRepo.AssertExpectations(t)
}

// Error case adding package
func TestPackageServer_AddPackage_Error(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	packageRepo.On("Add", mock.MatchedBy(func(p *db.Package) bool {
		return p.Active == true && p.Name == "daily-pack"
	})).Return(status.Errorf(codes.Internal, "error adding a package"))

	s := NewPackageServer(packageRepo, nil)

	ActPackage, err := s.Add(context.TODO(), &pb.AddPackageRequest{
		Active: true,
		Name:   "daily-pack",
		OrgID:  uuid.NewV4().String(),
	})
	assert.Error(t, err)
	assert.Nil(t, ActPackage)
	packageRepo.AssertExpectations(t)
}

// End Add packages //

// Update packages //
// Success case
// func TestPackageServer_UpdatePackage_Success(t *testing.T) {
// packageRepo := &mocks.PackageRepo{}
// s := NewPackageServer(packageRepo)
// packageUUID := uuid.NewString()
// packageRepo.On("Update", uuid.MustParse(packageUUID), mock.Anything).Return(&db.Package{
// Active: false,
// }, nil)
// up, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
// Uuid:   packageUUID,
// Active: false,
// })
// assert.NoError(t, err)
// assert.Equal(t, false, up.Package.Active)
// packageRepo.AssertExpectations(t)
// }

// Error case
func TestPackageServer_UpdatePackage_Error(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo, nil)
	packageUUID := uuid.NewV4()
	packageRepo.On("Update", packageUUID, mock.Anything).Return(nil, grpc.SqlErrorToGrpc(errors.New("Error updating records"), "rates"))
	pkg, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
		Uuid: packageUUID.String(),
	})
	assert.Error(t, err)
	assert.Nil(t, pkg)
}

// // End Update package //

// Delete package //
// Success case
// func TestPackageServer_DeletePackage_Success(t *testing.T) {
// 	packageRepo := &mocks.PackageRepo{}
// 	s := NewPackageServer(packageRepo)
// 	packageUUID := uuid.NewV4()
// 	var mockFilters = &pb.DeletePackageRequest{
// 		Uuid: packageUUID.String(),
// 	}
// 	packageRepo.On("Delete", packageUUID).Return(db.Package{
// 		Uuid: packageUUID,
// 	}, nil)
// 	_, err := s.Delete(context.TODO(), mockFilters)
// 	assert.NoError(t, err)
// 	packageRepo.AssertExpectations(t)
// }

// Error case: OrgID 0
func TestPackageServer_DeletePackage_Error1(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo, nil)
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
	s := NewPackageServer(packageRepo, nil)
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

// Error case: Error deleting record
func TestPackageServer_DeletePackage_Error3(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo, nil)
	packageUUID := uuid.NewV4()
	var mf = &pb.DeletePackageRequest{
		Uuid: packageUUID.String(),
	}
	packageRepo.On("Delete", packageUUID).
		Return(grpc.SqlErrorToGrpc(errors.New("SQL error while deleting record"), "packages"))
	pkg3, err := s.Delete(context.TODO(), mf)
	assert.Error(t, err)
	assert.Nil(t, pkg3)
}

// End Delete package //
