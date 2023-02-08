package server

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/data_plan/package/mocks"
	pb "github.com/ukama/ukama/systems/data_plan/package/pb/gen"
	"github.com/ukama/ukama/systems/data_plan/package/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Get Packages //
// Success case
func TestPackageServer_GetPackages_Success(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	packageUUID := uuid.NewString()

	var mockFilters = &pb.GetPackageRequest{
		PackageUuid: packageUUID,
	}
	s := NewPackageServer(packageRepo)
	packageRepo.On("Get", uuid.MustParse(packageUUID)).Return(&db.Package{
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
	s := NewPackageServer(packageRepo)
	packageUUID := uuid.NewString()
	var mockFilters = &pb.GetPackageRequest{
		PackageUuid: packageUUID,
	}
	packageRepo.On("Get", uuid.MustParse(packageUUID)).Return(nil, grpc.SqlErrorToGrpc(errors.New("SQL error while fetching records"), "packages"))
	pkg2, err := s.Get(context.TODO(), mockFilters)
	assert.Error(t, err)
	assert.Nil(t, pkg2)
}

// End Get packages //

// Get Package by org //

func TestPackageServer_GetPackageByOrg_Success(t *testing.T) {
	var orgId = uuid.New()

	packageRepo := &mocks.PackageRepo{}
	var mockFilters = &pb.GetByOrgPackageRequest{
		OrgId: orgId.String(),
	}
	s := NewPackageServer(packageRepo)

	packageRepo.On("GetByOrg", orgId).Return([]db.Package{{
		Sim_type:     pb.SimType_INTER_MNO_ALL.String(),
		Name:         "Daily-pack",
		Org_id:       orgId,
		Active:       true,
		Duration:     1,
		Sms_volume:   20,
		Data_volume:  12,
		Voice_volume: 34,
		Org_rates_id: 00,
	}}, nil)
	pkg, err := s.GetByOrg(context.TODO(), mockFilters)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), pkg.Packages[0].DataVolume)
	packageRepo.AssertExpectations(t)
}

// Error cases

func TestPackageServer_GetPackageByOrg_Error(t *testing.T) {
	var orgId = uuid.New()

	packageRepo := &mocks.PackageRepo{}
	var mockFilters = &pb.GetByOrgPackageRequest{
		OrgId: orgId.String(),
	}
	s := NewPackageServer(packageRepo)

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
	packageRepo.On("Add", mock.MatchedBy(func(p *db.Package) bool {
		return p.Active == true && p.Name == "daily-pack"
	})).Return(nil).Once()

	s := NewPackageServer(packageRepo)

	ActPackage, err := s.Add(context.TODO(), &pb.AddPackageRequest{
		Active: true,
		Name:   "daily-pack",
		OrgId:  uuid.New().String(),
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

	s := NewPackageServer(packageRepo)

	ActPackage, err := s.Add(context.TODO(), &pb.AddPackageRequest{
		Active: true,
		Name:   "daily-pack",
		OrgId:  uuid.New().String(),
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
	s := NewPackageServer(packageRepo)
	packageUUID := uuid.NewString()
	packageRepo.On("Update", uuid.MustParse(packageUUID), mock.Anything).Return(nil, grpc.SqlErrorToGrpc(errors.New("Error updating records"), "rates"))
	pkg, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
		Uuid: packageUUID,
	})
	assert.Error(t, err)
	assert.Nil(t, pkg)
}

// // End Update package //

// // Delete package //
// // Success case
// func TestPackageServer_DeletePackage_Success(t *testing.T) {
// packageRepo := &mocks.PackageRepo{}
// s := NewPackageServer(packageRepo)
// packageUUID := uuid.NewString()
// var mockFilters = &pb.DeletePackageRequest{
// PackageUuid: packageUUID,
// }
// packageRepo.On("Delete", uuid.MustParse(packageUUID)).Return(nil)
// _, err := s.Delete(context.TODO(), mockFilters)
// assert.NoError(t, err)
// packageRepo.AssertExpectations(t)
// }

// Error case: OrgID 0
func TestPackageServer_DeletePackage_Error1(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo)
	packageUUID := uuid.NewString()
	var mockFilters = &pb.DeletePackageRequest{
		PackageUuid: packageUUID,
	}
	packageRepo.On("Delete", uuid.MustParse(packageUUID)).
		Return(status.Errorf(codes.InvalidArgument, "OrgId is required."))
	pkg1, err := s.Delete(context.TODO(), mockFilters)
	assert.Error(t, err)
	assert.Nil(t, pkg1)
}

// Error case: Id 0
func TestPackageServer_DeletePackage_Success_Error2(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo)
	packageUUID := uuid.NewString()
	var mockFilters = &pb.DeletePackageRequest{
		PackageUuid: packageUUID,
	}
	packageRepo.On("Delete", uuid.MustParse(packageUUID)).
		Return(status.Errorf(codes.InvalidArgument, "Id is required."))
	pkg2, err := s.Delete(context.TODO(), mockFilters)
	assert.Error(t, err)
	assert.Nil(t, pkg2)
}

// Error case: Error deleting record
func TestPackageServer_DeletePackage_Error3(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo)
	packageUUID := uuid.NewString()
	var mf = &pb.DeletePackageRequest{
		PackageUuid: packageUUID,
	}
	packageRepo.On("Delete", uuid.MustParse(packageUUID)).
		Return(grpc.SqlErrorToGrpc(errors.New("SQL error while deleting record"), "packages"))
	pkg3, err := s.Delete(context.TODO(), mf)
	fmt.Println(err)
	assert.Error(t, err)
	assert.Nil(t, pkg3)
}

// End Delete package //
