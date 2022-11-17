package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/data-plan/package/mocks"
	"github.com/ukama/ukama/systems/data-plan/package/pb"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"
)

func TestPackageServer_GetPackages(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	var mockFilters = &pb.GetPackagesRequest{
		Id:    1,
		OrgId: 2323,
	}
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

	s := NewPackageServer(packageRepo)

	_package, err := s.GetPackages(context.TODO(), mockFilters)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), _package.Packages[0].DataVolume)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_AddPackage(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	packageRepo.On("Add", mock.MatchedBy(func(p *db.Package) bool {
		return p.Active == true && p.Name == "daily-pack"
	})).Return(nil).Once()

	s := NewPackageServer(packageRepo)

	ActPackage, err := s.AddPackage(context.TODO(), &pb.AddPackageRequest{
		Active: true,
		Name:   "daily-pack",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, ActPackage.Package.Active)
	packageRepo.AssertExpectations(t)
}

func TestPackageServer_UpdatePackage(t *testing.T) {
	packageRepo := &mocks.PackageRepo{}
	packageRepo.On("Update", uint64(1), mock.Anything).Return(&db.Package{
		Active: false,
	}, nil)

	s := NewPackageServer(packageRepo)

	ap, err := s.UpdatePackage(context.TODO(), &pb.UpdatePackageRequest{
		Id:     uint64(1),
		Active: false,
	})
	assert.NoError(t, err)
	assert.Equal(t, false, ap.Package.Active)
	packageRepo.AssertExpectations(t)
}
