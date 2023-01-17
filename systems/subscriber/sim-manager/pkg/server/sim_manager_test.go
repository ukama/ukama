package server

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/mocks"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
	"gorm.io/gorm"

	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
)

func TestSimManagerServer_GetSim(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		var simID = uuid.New()

		simRepo := &mocks.SimRepo{}

		simRepo.On("Get", simID).Return(
			&db.Sim{ID: simID,
				SubscriberID: uuid.New(),
				NetworkID:    uuid.New(),
				OrgID:        uuid.New(),
				IsPhysical:   false,
			}, nil).Once()

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "")
		simResp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			SimID: simID.String()})

		assert.NoError(t, err)
		assert.NotNil(t, simResp)
		assert.Equal(t, simID.String(), simResp.GetSim().GetId())
		assert.Equal(t, false, simResp.GetSim().IsPhysical)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		var simID = uuid.Nil

		simRepo := &mocks.SimRepo{}

		simRepo.On("Get", simID).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "")
		simResp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			SimID: simID.String()})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimUUIDInvalid", func(t *testing.T) {
		var simID = "1"

		simRepo := &mocks.SimRepo{}

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "")
		simResp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			SimID: simID})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_GetSimsBySubscriber(t *testing.T) {
	t.Run("SubscriberFound", func(t *testing.T) {
		var simID = uuid.New()
		var subscriberID = uuid.New()

		simRepo := &mocks.SimRepo{}

		simRepo.On("GetBySubscriber", subscriberID).Return(
			[]db.Sim{
				db.Sim{ID: simID,
					SubscriberID: subscriberID,
					IsPhysical:   false,
				}}, nil).Once()

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "")
		simResp, err := s.GetSimsBySubscriber(context.TODO(),
			&pb.GetSimsBySubscriberRequest{SubscriberID: subscriberID.String()})

		assert.NoError(t, err)
		assert.NotNil(t, simResp)
		assert.Equal(t, simID.String(), simResp.GetSims()[0].GetId())
		assert.Equal(t, subscriberID.String(), simResp.SubscriberID)
		simRepo.AssertExpectations(t)
	})

	t.Run("SubscriberNotFound", func(t *testing.T) {
		var subscriberID = uuid.Nil

		simRepo := &mocks.SimRepo{}

		simRepo.On("GetBySubscriber", subscriberID).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "")
		simResp, err := s.GetSimsBySubscriber(context.TODO(), &pb.GetSimsBySubscriberRequest{
			SubscriberID: subscriberID.String()})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SubscriberUUIDInvalid", func(t *testing.T) {
		var subscriberID = "1"

		simRepo := &mocks.SimRepo{}

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "")
		simResp, err := s.GetSimsBySubscriber(context.TODO(), &pb.GetSimsBySubscriberRequest{
			SubscriberID: subscriberID})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})

}

func TestSimManagerServer_GetSimsByNetwork(t *testing.T) {
	t.Run("NetworkFound", func(t *testing.T) {
		var simID = uuid.New()
		var networkID = uuid.New()

		simRepo := &mocks.SimRepo{}

		simRepo.On("GetByNetwork", networkID).Return(
			[]db.Sim{
				db.Sim{ID: simID,
					NetworkID:  networkID,
					IsPhysical: false,
				}}, nil).Once()

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "")
		simResp, err := s.GetSimsByNetwork(context.TODO(),
			&pb.GetSimsByNetworkRequest{NetworkID: networkID.String()})

		assert.NoError(t, err)
		assert.NotNil(t, simResp)
		assert.Equal(t, simID.String(), simResp.GetSims()[0].GetId())
		assert.Equal(t, networkID.String(), simResp.NetworkID)
		simRepo.AssertExpectations(t)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		var networkID = uuid.Nil

		simRepo := &mocks.SimRepo{}

		simRepo.On("GetByNetwork", networkID).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "")
		simResp, err := s.GetSimsByNetwork(context.TODO(), &pb.GetSimsByNetworkRequest{
			NetworkID: networkID.String()})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("NetworkUUIDInvalid", func(t *testing.T) {
		var networkID = "1"

		simRepo := &mocks.SimRepo{}

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "")
		simResp, err := s.GetSimsByNetwork(context.TODO(), &pb.GetSimsByNetworkRequest{
			NetworkID: networkID})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_GetPackagesBySim(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		var simID = uuid.New()
		var packageID = uuid.New()

		packageRepo := &mocks.PackageRepo{}

		packageRepo.On("GetBySim", simID).Return(
			[]db.Package{
				db.Package{ID: packageID,
					SimID:    simID,
					IsActive: false,
				}}, nil).Once()

		s := NewSimManagerServer(nil, packageRepo, nil, nil, nil, nil, "")

		resp, err := s.GetPackagesBySim(context.TODO(),
			&pb.GetPackagesBySimRequest{SimID: simID.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, packageID.String(), resp.GetPackages()[0].GetId())
		assert.Equal(t, simID.String(), resp.SimID)
		packageRepo.AssertExpectations(t)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		var simID = uuid.Nil

		packageRepo := &mocks.PackageRepo{}

		packageRepo.On("GetBySim", simID).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSimManagerServer(nil, packageRepo, nil, nil, nil, nil, "")

		resp, err := s.GetPackagesBySim(context.TODO(), &pb.GetPackagesBySimRequest{
			SimID: simID.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})

	t.Run("SimUUIDInvalid", func(t *testing.T) {
		var simID = "1"

		packageRepo := &mocks.PackageRepo{}

		s := NewSimManagerServer(nil, packageRepo, nil, nil, nil, nil, "")

		resp, err := s.GetPackagesBySim(context.TODO(), &pb.GetPackagesBySimRequest{
			SimID: simID})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_RemovePackageForSim(t *testing.T) {
	t.Run("PackageFound", func(t *testing.T) {
		var packageID = uuid.New()

		packageRepo := &mocks.PackageRepo{}

		packageRepo.On("Delete", packageID,
			mock.AnythingOfType("func(uuid.UUID, *gorm.DB) error")).Return(nil).Once()

		s := NewSimManagerServer(nil, packageRepo, nil, nil, nil, nil, "")

		resp, err := s.RemovePackageForSim(context.TODO(), &pb.RemovePackageRequest{
			PackageID: packageID.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		packageRepo.AssertExpectations(t)
	})

	t.Run("PackageNotFound", func(t *testing.T) {
		var packageID = uuid.Nil

		packageRepo := &mocks.PackageRepo{}

		packageRepo.On("Delete", packageID,
			mock.AnythingOfType("func(uuid.UUID, *gorm.DB) error")).Return(gorm.ErrRecordNotFound).Once()

		s := NewSimManagerServer(nil, packageRepo, nil, nil, nil, nil, "")

		resp, err := s.RemovePackageForSim(context.TODO(), &pb.RemovePackageRequest{
			PackageID: packageID.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})

	t.Run("PackageUUIDInvalid", func(t *testing.T) {
		var packageID = "1"

		packageRepo := &mocks.PackageRepo{}

		s := NewSimManagerServer(nil, packageRepo, nil, nil, nil, nil, "")

		resp, err := s.RemovePackageForSim(context.TODO(), &pb.RemovePackageRequest{
			PackageID: packageID})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})
}
