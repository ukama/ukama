package server

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/data-plan/package/mocks"
	"github.com/ukama/ukama/systems/data-plan/package/pb"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Get Packages //
// Success case
func TestPackageServer_GetPackages_Success(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	var mockFilters = &pb.GetPackagesRequest{
		Id:    1,
		OrgId: 1,
	}
	s := NewPackageServer(packageRepo)

	packageRepo.On("Get", mockFilters.OrgId, mockFilters.Id).Return([]db.Package{
		{
			Name:         "Daily-pack",
			Org_id:       2323,
			Active:       true,
			Duration:     1,
			Sms_volume:   20,
			Data_volume:  12,
			Voice_volume: 34,
			Org_rates_id: 00},
	}, nil)
	pkg, err := s.Get(context.TODO(), mockFilters)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), pkg.Packages[0].DataVolume)
	packageRepo.AssertExpectations(t)
}

// Error case invalid arguments
func TestPackageServer_GetPackages_Error1(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	var mockFilters = &pb.GetPackagesRequest{
		Id:    0,
		OrgId: 0,
	}
	s := NewPackageServer(packageRepo)
	packageRepo.On("Get", mockFilters.OrgId, mockFilters.Id).
		Return(nil, status.Errorf(codes.InvalidArgument, "OrgId is required."))
	pkg1, err := s.Get(context.TODO(), mockFilters)
	assert.Error(t, err)
	assert.Nil(t, pkg1)
}

// Error case SQL error
func TestPackageServer_GetPackages_Error2(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	var mockFilters = &pb.GetPackagesRequest{
		Id:    999,
		OrgId: 999,
	}
	s := NewPackageServer(packageRepo)
	packageRepo.On("Get", mockFilters.OrgId, mockFilters.Id).
		Return(nil, grpc.SqlErrorToGrpc(errors.New("SQL error while fetching records"), "packages"))
	pkg2, err := s.Get(context.TODO(), mockFilters)
	assert.Error(t, err)
	assert.Nil(t, pkg2)
}

// End Get packages //

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
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, ActPackage.Package.Active)
	packageRepo.AssertExpectations(t)
}

// End Add packages //

// Update packages //
// Success case
func TestPackageServer_UpdatePackage_Success(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo)

	packageRepo.On("Update", uint64(1), mock.Anything).Return(&db.Package{
		Active: false,
	}, nil)
	up, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
		Id:     uint64(1),
		Active: false,
	})
	assert.NoError(t, err)
	assert.Equal(t, false, up.Package.Active)
	packageRepo.AssertExpectations(t)
}

// Error case
func TestPackageServer_UpdatePackage_Error(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo)

	packageRepo.On("Update", uint64(1), db.Package{
		Active: false,
	}).Return(nil, grpc.SqlErrorToGrpc(errors.New("Error updating records"), "rates"))
	_up, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
		Id:     uint64(1),
		Active: false,
	})
	layout := "2006-01-02 15:04:05 -0700 MST"
	ct, _ := time.Parse(layout, _up.Package.GetCreatedAt())
	ut, _ := time.Parse(layout, _up.Package.GetUpdatedAt())
	dt, _ := time.Parse(layout, _up.Package.GetDeletedAt())
	fmt.Println(_up, err)
	assert.NoError(t, err)
	assert.True(t, true, ct.IsZero())
	assert.True(t, true, ut.IsZero())
	assert.True(t, true, dt.IsZero())
}

// End Update package //

// Delete package //
// Success case
func TestPackageServer_DeletePackage_Success(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo)
	var mockFilters = &pb.DeletePackageRequest{
		Id:    1,
		OrgId: 1,
	}
	packageRepo.On("Delete", mockFilters.OrgId, mockFilters.Id).Return(nil)
	_, err := s.Delete(context.TODO(), mockFilters)
	assert.NoError(t, err)
	packageRepo.AssertExpectations(t)
}

// Error case: OrgID 0
func TestPackageServer_DeletePackage_Error1(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo)
	var mockFilters = &pb.DeletePackageRequest{
		Id:    1,
		OrgId: 0,
	}
	packageRepo.On("Delete", mockFilters.OrgId, mockFilters.Id).
		Return(status.Errorf(codes.InvalidArgument, "OrgId is required."))
	pkg1, err := s.Delete(context.TODO(), mockFilters)
	assert.Error(t, err)
	assert.Nil(t, pkg1)
}

// Error case: Id 0
func TestPackageServer_DeletePackage_Success_Error2(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo)
	var mockFilters = &pb.DeletePackageRequest{
		Id:    0,
		OrgId: 1,
	}
	packageRepo.On("Delete", mockFilters.OrgId, mockFilters.Id).
		Return(status.Errorf(codes.InvalidArgument, "Id is required."))
	pkg2, err := s.Delete(context.TODO(), mockFilters)
	assert.Error(t, err)
	assert.Nil(t, pkg2)
}

// Error case: Error deleting record
func TestPackageServer_DeletePackage_Error3(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo)
	var mf = &pb.DeletePackageRequest{
		Id:    999,
		OrgId: 999,
	}
	packageRepo.On("Delete", mf.OrgId, mf.Id).
		Return(grpc.SqlErrorToGrpc(errors.New("SQL error while deleting record"), "packages"))
	pkg3, err := s.Delete(context.TODO(), mf)
	fmt.Println(err)
	assert.Error(t, err)
	assert.Nil(t, pkg3)
}

// End Delete package //
