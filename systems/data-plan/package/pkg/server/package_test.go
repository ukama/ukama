package server

import (
	"context"
	"errors"
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

func TestPackageServer_GetPackages(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	var mockFilters = &pb.GetPackagesRequest{
		Id:    1,
		OrgId: 1,
	}
	var failMockFilters1 = &pb.GetPackagesRequest{
		Id:    0,
		OrgId: 0,
	}
	var failMockFilters2 = &pb.GetPackagesRequest{
		Id:    999,
		OrgId: 999,
	}
	s := NewPackageServer(packageRepo)

	// Success case
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

	// Error case: Invalid orgID
	packageRepo.On("Get", failMockFilters1.OrgId, failMockFilters1.Id).
		Return(nil, status.Errorf(codes.InvalidArgument, "OrgId is required."))
	pkg1, err := s.Get(context.TODO(), failMockFilters1)
	assert.Error(t, err)
	assert.Nil(t, pkg1)

	// Error case: Error fetching records
	packageRepo.On("Get", failMockFilters2.OrgId, failMockFilters2.Id).
		Return(nil, grpc.SqlErrorToGrpc(errors.New("SQL error while fetching records"), "packages"))
	pkg2, err := s.Get(context.TODO(), failMockFilters2)
	assert.Error(t, err)
	assert.Nil(t, pkg2)
}

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

func TestPackageServer_UpdatePackage(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo)

	// Success case
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

	//Error case
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
	assert.NoError(t, err)
	assert.True(t, true, ct.IsZero())
	assert.True(t, true, ut.IsZero())
	assert.True(t, true, dt.IsZero())
}

func TestPackageServer_DeletePackage(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	s := NewPackageServer(packageRepo)
	var mockFilters = &pb.DeletePackageRequest{
		Id:    1,
		OrgId: 1,
	}
	var failMockFilters1 = &pb.DeletePackageRequest{
		Id:    1,
		OrgId: 0,
	}
	var failMockFilters2 = &pb.DeletePackageRequest{
		Id:    0,
		OrgId: 1,
	}
	var failMockFilters3 = &pb.DeletePackageRequest{
		Id:    999,
		OrgId: 999,
	}

	// Success case
	packageRepo.On("Delete", mockFilters.OrgId, mockFilters.Id).Return(nil)
	_, err := s.Delete(context.TODO(), mockFilters)
	assert.NoError(t, err)
	packageRepo.AssertExpectations(t)

	// Error case: OrgID 0
	packageRepo.On("Delete", failMockFilters1.OrgId, failMockFilters1.Id).
		Return(status.Errorf(codes.InvalidArgument, "OrgId is required."))
	pkg1, err := s.Delete(context.TODO(), failMockFilters1)
	assert.Error(t, err)
	assert.Nil(t, pkg1)

	// Error case: Id 0
	packageRepo.On("Delete", failMockFilters2.OrgId, failMockFilters2.Id).
		Return(status.Errorf(codes.InvalidArgument, "Id is required."))
	pkg2, err := s.Delete(context.TODO(), failMockFilters2)
	assert.Error(t, err)
	assert.Nil(t, pkg2)

	// Error case: Error deleting record
	packageRepo.On("Delete", failMockFilters3.OrgId, failMockFilters3.Id).
		Return(grpc.SqlErrorToGrpc(errors.New("SQL error while deleting record"), "packages"))
	pkg3, err := s.Delete(context.TODO(), failMockFilters3)
	assert.Error(t, err)
	assert.Nil(t, pkg3)
}
