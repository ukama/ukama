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
	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	"github.com/ukama/ukama/systems/data-plan/package/mocks"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	ukama "github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	rpb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	splmocks "github.com/ukama/ukama/systems/data-plan/rate/pb/gen/mocks"
)

const OrgName = "testorg"
const OrgId = "8c6c2bec-5f90-4fee-8ffd-ee6456abf4fc"

// Fixed timestamps for consistent testing
var (
	fixedBaseTime    = time.Date(2025, 7, 15, 12, 0, 0, 0, time.UTC)
	fixedFromTime    = time.Date(2025, 8, 15, 12, 0, 0, 0, time.UTC) // 30 days from base
	fixedToTime      = time.Date(2025, 9, 15, 12, 0, 0, 0, time.UTC) // 60 days from base
	fixedPastTime    = time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC) // 1 day before base
	fixedLaterTime   = time.Date(2025, 9, 15, 12, 0, 0, 0, time.UTC) // 60 days from base (for invalid range tests)
	fixedEarlierTime = time.Date(2025, 8, 15, 12, 0, 0, 0, time.UTC) // 30 days from base (for invalid range tests)
)

// Helper types
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

// ============================================================================
// GET TESTS
// ============================================================================

func TestPackageServer_Get(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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
	})

	t.Run("Error_Database", func(t *testing.T) {
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
	})

	t.Run("Error_InvalidUUID", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)

		req := &pb.GetPackageRequest{
			Uuid: "invalid-uuid",
		}

		pkg, err := s.Get(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, pkg)
		assert.Contains(t, err.Error(), "invalid format of package uuid")
	})
}

func TestPackageServer_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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
	})

	t.Run("Error_Database", func(t *testing.T) {
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
	})
}

func TestPackageServer_GetDetails(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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
	})

	t.Run("Error_InvalidUUID", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)

		req := &pb.GetPackageRequest{
			Uuid: "invalid-uuid",
		}

		pkg, err := s.GetDetails(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, pkg)
		assert.Contains(t, err.Error(), "invalid format of package uuid")
	})

	t.Run("Error_Database", func(t *testing.T) {
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
	})

	t.Run("Error_EmptyUUID", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)

		req := &pb.GetPackageRequest{
			Uuid: "",
		}

		pkg, err := s.GetDetails(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, pkg)
		assert.Contains(t, err.Error(), "invalid format of package uuid")
	})

	t.Run("Success_CompletePackageData", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		packageUUID := uuid.NewV4()
		ownerUUID := uuid.NewV4()
		baseRateUUID := uuid.NewV4()
		now := fixedBaseTime

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
	})

	t.Run("Success_InactivePackage", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		packageUUID := uuid.NewV4()

		req := &pb.GetPackageRequest{
			Uuid: packageUUID.String(),
		}

		inactivePackage := &db.Package{
			Uuid:   packageUUID,
			Name:   "Inactive Package",
			Active: false,
			From:   fixedBaseTime,
			To:     fixedToTime,
		}

		s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
		packageRepo.On("GetDetails", packageUUID).Return(inactivePackage, nil)

		pkg, err := s.GetDetails(context.TODO(), req)
		assert.NoError(t, err)
		assert.NotNil(t, pkg)
		assert.Equal(t, false, pkg.Package.Active)
		assert.Equal(t, "Inactive Package", pkg.Package.Name)
		packageRepo.AssertExpectations(t)
	})

	t.Run("Success_FlatratePackage", func(t *testing.T) {
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
			From:     fixedBaseTime,
			To:       fixedToTime,
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
	})

	t.Run("Success_SpecialCharactersInName", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		packageUUID := uuid.NewV4()

		req := &pb.GetPackageRequest{
			Uuid: packageUUID.String(),
		}

		specialNamePackage := &db.Package{
			Uuid:   packageUUID,
			Name:   "Package with special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?",
			Active: true,
			From:   fixedBaseTime,
			To:     fixedToTime,
		}

		s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
		packageRepo.On("GetDetails", packageUUID).Return(specialNamePackage, nil)

		pkg, err := s.GetDetails(context.TODO(), req)
		assert.NoError(t, err)
		assert.NotNil(t, pkg)
		assert.Equal(t, "Package with special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?", pkg.Package.Name)
		packageRepo.AssertExpectations(t)
	})

	t.Run("Success_ZeroVolumes", func(t *testing.T) {
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
			From:        fixedBaseTime,
			To:          fixedToTime,
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
	})

	t.Run("Error_ContextCancellation", func(t *testing.T) {
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
	})

	t.Run("Error_RecordNotFound", func(t *testing.T) {
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
	})

	t.Run("Success_DifferentSimTypes", func(t *testing.T) {
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
					From:    fixedBaseTime,
					To:      fixedToTime,
				}

				packageRepo.On("GetDetails", packageUUID).Return(testPackage, nil)

				pkg, err := s.GetDetails(context.TODO(), req)
				assert.NoError(t, err)
				assert.NotNil(t, pkg)
				assert.Equal(t, tc.expected, pkg.Package.SimType)
				packageRepo.AssertExpectations(t)
			})
		}
	})
}

// ============================================================================
// ADD TESTS
// ============================================================================

func TestPackageServer_Add(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ownerId := uuid.NewV4().String()
		baserate := uuid.NewV4().String()
		packageRepo := &mocks.PackageRepo{}
		rate := &mocks.RateClientProvider{}
		packageRepo.On("Add", mock.MatchedBy(func(p *db.Package) bool {
			return p.Active == true && p.Name == "daily-pack"
		}), mock.Anything).Return(nil).Once()

		rateClient := &splmocks.RateServiceClient{}
		rate.On("GetClient").Return(rateClient, nil).Once()

		rateResponse := &rpb.GetRateByIdResponse{
			Rate: &bpb.Rate{
				SmsMo:    1,
				SmsMt:    1,
				Data:     10,
				Country:  "USA",
				Provider: "ukama",
			},
		}

		rateClient.On("GetRateById", mock.Anything, &rpb.GetRateByIdRequest{
			OwnerId:  ownerId,
			BaseRate: baserate,
		}).Return(rateResponse, nil).Once()

		s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

		ActPackage, err := s.Add(context.TODO(), &pb.AddPackageRequest{
			Active:     true,
			Name:       "daily-pack",
			SimType:    testSim,
			OwnerId:    ownerId,
			BaserateId: baserate,
			Country:    rateResponse.Rate.Country,
			Provider:   rateResponse.Rate.Provider,
			From:       fixedFromTime.Format(time.RFC3339),
			To:         fixedToTime.Format(time.RFC3339),
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, ActPackage.Package.Active)
		packageRepo.AssertExpectations(t)
	})

	t.Run("Error_InvalidOwnerUUID", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		rate := &mocks.RateClientProvider{}

		s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

		req := &pb.AddPackageRequest{
			OwnerId:    "invalid-uuid",
			BaserateId: uuid.NewV4().String(),
			From:       fixedFromTime.Format(time.RFC3339),
			To:         fixedToTime.Format(time.RFC3339),
		}

		resp, err := s.Add(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid format of owner uuid")
	})

	t.Run("Error_InvalidBaserateUUID", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		rate := &mocks.RateClientProvider{}

		s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

		req := &pb.AddPackageRequest{
			OwnerId:    uuid.NewV4().String(),
			BaserateId: "invalid-uuid",
			From:       fixedFromTime.Format(time.RFC3339),
			To:         fixedToTime.Format(time.RFC3339),
		}

		resp, err := s.Add(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid format of base rate")
	})

	t.Run("Error_InvalidFromDate", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		rate := &mocks.RateClientProvider{}

		s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

		req := &pb.AddPackageRequest{
			OwnerId:    uuid.NewV4().String(),
			BaserateId: uuid.NewV4().String(),
			From:       "invalid-date",
			To:         fixedToTime.Format(time.RFC3339),
		}

		resp, err := s.Add(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error:")
	})

	t.Run("Error_InvalidToDate", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		rate := &mocks.RateClientProvider{}

		s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

		req := &pb.AddPackageRequest{
			OwnerId:    uuid.NewV4().String(),
			BaserateId: uuid.NewV4().String(),
			From:       fixedFromTime.Format(time.RFC3339),
			To:         "invalid-date",
		}

		resp, err := s.Add(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error:")
	})

	t.Run("Error_PastFromDate", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		rate := &mocks.RateClientProvider{}

		s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

		req := &pb.AddPackageRequest{
			OwnerId:    uuid.NewV4().String(),
			BaserateId: uuid.NewV4().String(),
			From:       fixedPastTime.Format(time.RFC3339), // Past date
			To:         fixedToTime.Format(time.RFC3339),
		}

		resp, err := s.Add(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error:")
	})

	t.Run("Error_PastToDate", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		rate := &mocks.RateClientProvider{}

		s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

		req := &pb.AddPackageRequest{
			OwnerId:    uuid.NewV4().String(),
			BaserateId: uuid.NewV4().String(),
			From:       fixedFromTime.Format(time.RFC3339),
			To:         fixedPastTime.Format(time.RFC3339), // Past date
		}

		resp, err := s.Add(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error:")
	})

	t.Run("Error_InvalidDateRange", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		rate := &mocks.RateClientProvider{}

		s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

		req := &pb.AddPackageRequest{
			OwnerId:    uuid.NewV4().String(),
			BaserateId: uuid.NewV4().String(),
			From:       fixedToTime.Format(time.RFC3339),   // Later date
			To:         fixedFromTime.Format(time.RFC3339), // Earlier date
		}

		resp, err := s.Add(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error:")
	})

	t.Run("Error_RateClientError", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		rate := &mocks.RateClientProvider{}

		s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

		req := &pb.AddPackageRequest{
			OwnerId:    uuid.NewV4().String(),
			BaserateId: uuid.NewV4().String(),
			From:       fixedFromTime.Format(time.RFC3339),
			To:         fixedToTime.Format(time.RFC3339),
		}

		rate.On("GetClient").Return(nil, errors.New("rate client error"))

		resp, err := s.Add(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "rate client error")
	})

	t.Run("Error_GetRateByIdError", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		rate := &mocks.RateClientProvider{}
		ownerId := uuid.NewV4().String()
		baserate := uuid.NewV4().String()

		s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

		req := &pb.AddPackageRequest{
			OwnerId:    ownerId,
			BaserateId: baserate,
			From:       fixedFromTime.Format(time.RFC3339),
			To:         fixedToTime.Format(time.RFC3339),
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
	})

	t.Run("Error_PackageRepoError", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		rate := &mocks.RateClientProvider{}
		ownerId := uuid.NewV4().String()
		baserate := uuid.NewV4().String()

		s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

		req := &pb.AddPackageRequest{
			OwnerId:    ownerId,
			BaserateId: baserate,
			From:       fixedFromTime.Format(time.RFC3339),
			To:         fixedToTime.Format(time.RFC3339),
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
	})

	t.Run("Success_FlatratePackage", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		rate := &mocks.RateClientProvider{}
		ownerId := uuid.NewV4().String()
		baserate := uuid.NewV4().String()

		s := NewPackageServer(OrgName, packageRepo, rate, nil, OrgId)

		req := &pb.AddPackageRequest{
			OwnerId:    ownerId,
			BaserateId: baserate,
			Flatrate:   true, // Test flatrate package
			From:       fixedFromTime.Format(time.RFC3339),
			To:         fixedToTime.Format(time.RFC3339),
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
	})

	t.Run("Success_WithMessageBus", func(t *testing.T) {
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
			From:          fixedFromTime.Format(time.RFC3339),
			To:            fixedToTime.Format(time.RFC3339),
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
	})

	t.Run("Success_MessageBusPublishError", func(t *testing.T) {
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
			From:       fixedFromTime.Format(time.RFC3339),
			To:         fixedToTime.Format(time.RFC3339),
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
	})

	t.Run("Success_WithEmptyApn", func(t *testing.T) {
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
			From:       fixedFromTime.Format(time.RFC3339),
			To:         fixedToTime.Format(time.RFC3339),
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
	})

	t.Run("Success_NonFlatratePackageWithCalculations", func(t *testing.T) {
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
			From:        fixedFromTime.Format(time.RFC3339),
			To:          fixedToTime.Format(time.RFC3339),
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
		assert.Equal(t, 532.0, resp.Package.Amount)

		packageRepo.AssertExpectations(t)
	})

	t.Run("Success_WithDifferentSimTypes", func(t *testing.T) {
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
					From:       fixedFromTime.Format(time.RFC3339),
					To:         fixedToTime.Format(time.RFC3339),
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
	})

	t.Run("Success_WithDifferentPackageTypes", func(t *testing.T) {
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
					From:       fixedFromTime.Format(time.RFC3339),
					To:         fixedToTime.Format(time.RFC3339),
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
	})

	t.Run("Success_WithDifferentUnits", func(t *testing.T) {
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
					From:        fixedFromTime.Format(time.RFC3339),
					To:          fixedToTime.Format(time.RFC3339),
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
	})

	t.Run("Success_WithZeroVolumes", func(t *testing.T) {
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
			From:        fixedFromTime.Format(time.RFC3339),
			To:          fixedToTime.Format(time.RFC3339),
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
	})

	t.Run("Success_WithLargeVolumes", func(t *testing.T) {
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
			From:        fixedFromTime.Format(time.RFC3339),
			To:          fixedToTime.Format(time.RFC3339),
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
		assert.Equal(t, 10200.0, resp.Package.Amount)

		packageRepo.AssertExpectations(t)
	})

	t.Run("Success_WithNegativeMarkup", func(t *testing.T) {
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
			From:       fixedFromTime.Format(time.RFC3339),
			To:         fixedToTime.Format(time.RFC3339),
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
	})
}

// ============================================================================
// DELETE TESTS
// ============================================================================

func TestPackageServer_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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
	})

	t.Run("Error_InvalidUUID", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)

		req := &pb.DeletePackageRequest{
			Uuid: "invalid-uuid",
		}

		resp, err := s.Delete(context.TODO(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid format of package uuid")
	})

	t.Run("Error_DatabaseError1", func(t *testing.T) {
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
	})

	t.Run("Error_DatabaseError2", func(t *testing.T) {
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
	})

	t.Run("Success_MessageBusPublishError", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		packageUUID := uuid.NewV4()

		s := NewPackageServer(OrgName, packageRepo, nil, msgbusClient, OrgId)

		req := &pb.DeletePackageRequest{
			Uuid: packageUUID.String(),
		}

		// Mock successful package deletion
		packageRepo.On("Delete", packageUUID).Return(nil)

		// Mock message bus publish failure
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(errors.New("message bus publish failed"))

		resp, err := s.Delete(context.TODO(), req)
		// Should still succeed even if message bus fails
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, packageUUID.String(), resp.Uuid)

		packageRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})
}

// ============================================================================
// UPDATE TESTS
// ============================================================================

func TestPackageServer_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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
			From:   fixedBaseTime,
			To:     fixedToTime,
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
	})

	t.Run("Error_InvalidUUID", func(t *testing.T) {
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
	})

	t.Run("Error_RepoError", func(t *testing.T) {
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
	})

	t.Run("Success_MessageBusError", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		s := NewPackageServer(OrgName, packageRepo, nil, msgbusClient, OrgId)
		packageUUID := uuid.NewV4()
		packageRepo.On("Update", packageUUID, mock.Anything).Return(nil).Once()

		packageRepo.On("Get", packageUUID).Return(&db.Package{
			Uuid:   packageUUID,
			Name:   "msgbus-fail",
			Active: true,
			From:   fixedBaseTime,
			To:     fixedToTime,
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
	})

	t.Run("Success_OnlyName", func(t *testing.T) {
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
			From:   fixedBaseTime,
			To:     fixedToTime,
		}, nil).Once()

		resp, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
			Uuid: packageUUID.String(),
			Name: "NameOnly",
		})
		assert.NoError(t, err)
		assert.Equal(t, "NameOnly", resp.Package.Name)
		packageRepo.AssertExpectations(t)
	})

	t.Run("Success_Inactive", func(t *testing.T) {
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
			From:   fixedBaseTime,
			To:     fixedToTime,
		}, nil).Once()

		resp, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
			Uuid:   packageUUID.String(),
			Active: false,
		})
		assert.NoError(t, err)
		assert.False(t, resp.Package.Active)
		packageRepo.AssertExpectations(t)
	})

	t.Run("Success_NilMsgBus", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		s := NewPackageServer(OrgName, packageRepo, nil, nil, OrgId)
		packageUUID := uuid.NewV4()

		packageRepo.On("Update", packageUUID, mock.Anything).Return(nil).Once()

		packageRepo.On("Get", packageUUID).Return(&db.Package{
			Uuid:   packageUUID,
			Name:   "NoMsgBus",
			Active: true,
			From:   fixedBaseTime,
			To:     fixedToTime,
		}, nil).Once()

		resp, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
			Uuid:   packageUUID.String(),
			Name:   "NoMsgBus",
			Active: true,
		})
		assert.NoError(t, err)
		assert.Equal(t, "NoMsgBus", resp.Package.Name)
		packageRepo.AssertExpectations(t)
	})

	t.Run("Error_GetAfterUpdateError", func(t *testing.T) {
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
	})

	t.Run("Success_EmptyName", func(t *testing.T) {
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
			From:   fixedBaseTime,
			To:     fixedToTime,
		}, nil).Once()

		resp, err := s.Update(context.TODO(), &pb.UpdatePackageRequest{
			Uuid:   packageUUID.String(),
			Name:   "",
			Active: true,
		})
		assert.NoError(t, err)
		assert.Equal(t, "", resp.Package.Name)
		packageRepo.AssertExpectations(t)
	})
}
