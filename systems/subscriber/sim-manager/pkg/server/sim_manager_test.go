package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/mocks"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/providers"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
	sims "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"

	subspb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	subsmocks "github.com/ukama/ukama/systems/subscriber/registry/pb/gen/mocks"

	splpb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	splmocks "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen/mocks"

	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
)

const testIccid = "890000-this-is-a-test-iccid"

func TestSimManagerServer_GetSim(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		var simID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		sim := simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				Iccid:      testIccid,
				Status:     db.SimStatusInactive,
				Type:       db.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*db.Sim)

		agentAdapter := agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("GetSim", mock.Anything,
			sim.Iccid).Return(nil, nil).Once()

		s := NewSimManagerServer(simRepo, nil, agentFactory, nil, nil, nil, "", nil, "", "", nil, nil)
		simResp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			SimId: simID.String()})

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

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "", nil, "", "", nil, nil)
		simResp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			SimId: simID.String()})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimUUIDInvalid", func(t *testing.T) {
		var simID = "1"

		simRepo := &mocks.SimRepo{}

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "", nil, "", "", nil, nil)
		simResp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			SimId: simID})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_GetSimsBySubscriber(t *testing.T) {
	t.Run("SubscriberFound", func(t *testing.T) {
		var simID = uuid.NewV4()
		var subscriberID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}

		simRepo.On("GetBySubscriber", subscriberID).Return(
			[]db.Sim{
				{Id: simID,
					SubscriberId: subscriberID,
					IsPhysical:   false,
				}}, nil).Once()

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "", nil, "", "", nil, nil)
		simResp, err := s.GetSimsBySubscriber(context.TODO(),
			&pb.GetSimsBySubscriberRequest{SubscriberId: subscriberID.String()})

		assert.NoError(t, err)
		assert.NotNil(t, simResp)
		assert.Equal(t, simID.String(), simResp.GetSims()[0].GetId())
		assert.Equal(t, subscriberID.String(), simResp.SubscriberId)
		simRepo.AssertExpectations(t)
	})

	t.Run("SubscriberNotFound", func(t *testing.T) {
		var subscriberID = uuid.Nil

		simRepo := &mocks.SimRepo{}

		simRepo.On("GetBySubscriber", subscriberID).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "", nil, "", "", nil, nil)
		simResp, err := s.GetSimsBySubscriber(context.TODO(), &pb.GetSimsBySubscriberRequest{
			SubscriberId: subscriberID.String()})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SubscriberUUIDInvalid", func(t *testing.T) {
		var subscriberID = "1"

		simRepo := &mocks.SimRepo{}

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "", nil, "", "", nil, nil)
		simResp, err := s.GetSimsBySubscriber(context.TODO(), &pb.GetSimsBySubscriberRequest{
			SubscriberId: subscriberID})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})

}

func TestSimManagerServer_GetSimsByNetwork(t *testing.T) {
	t.Run("NetworkFound", func(t *testing.T) {
		var simID = uuid.NewV4()
		var networkID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}

		simRepo.On("GetByNetwork", networkID).Return(
			[]db.Sim{
				{Id: simID,
					NetworkId:  networkID,
					IsPhysical: false,
				}}, nil).Once()

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "", nil, "", "", nil, nil)
		simResp, err := s.GetSimsByNetwork(context.TODO(),
			&pb.GetSimsByNetworkRequest{NetworkId: networkID.String()})

		assert.NoError(t, err)
		assert.NotNil(t, simResp)
		assert.Equal(t, simID.String(), simResp.GetSims()[0].GetId())
		assert.Equal(t, networkID.String(), simResp.NetworkId)
		simRepo.AssertExpectations(t)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		var networkID = uuid.Nil

		simRepo := &mocks.SimRepo{}

		simRepo.On("GetByNetwork", networkID).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "", nil, "", "", nil, nil)
		simResp, err := s.GetSimsByNetwork(context.TODO(), &pb.GetSimsByNetworkRequest{
			NetworkId: networkID.String()})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("NetworkUUIDInvalid", func(t *testing.T) {
		var networkID = "1"

		simRepo := &mocks.SimRepo{}

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "", nil, "", "", nil, nil)
		simResp, err := s.GetSimsByNetwork(context.TODO(), &pb.GetSimsByNetworkRequest{
			NetworkId: networkID})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_GetPackagesBySim(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		var simID = uuid.NewV4()
		var packageID = uuid.NewV4()

		packageRepo := &mocks.PackageRepo{}

		packageRepo.On("GetBySim", simID).Return(
			[]db.Package{
				{Id: packageID,
					SimId:    simID,
					IsActive: false,
				}}, nil).Once()

		s := NewSimManagerServer(nil, packageRepo, nil, nil, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.GetPackagesBySim(context.TODO(),
			&pb.GetPackagesBySimRequest{SimId: simID.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, packageID.String(), resp.GetPackages()[0].GetId())
		assert.Equal(t, simID.String(), resp.SimId)
		packageRepo.AssertExpectations(t)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		var simID = uuid.Nil

		packageRepo := &mocks.PackageRepo{}

		packageRepo.On("GetBySim", simID).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSimManagerServer(nil, packageRepo, nil, nil, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.GetPackagesBySim(context.TODO(), &pb.GetPackagesBySimRequest{
			SimId: simID.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})

	t.Run("SimUUIDInvalid", func(t *testing.T) {
		var simID = "1"

		packageRepo := &mocks.PackageRepo{}

		s := NewSimManagerServer(nil, packageRepo, nil, nil, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.GetPackagesBySim(context.TODO(), &pb.GetPackagesBySimRequest{
			SimId: simID})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_AllocateSim(t *testing.T) {
	t.Run("SimAndPackageFound", func(t *testing.T) {
		var subscriberID = uuid.NewV4()
		var networkID = uuid.NewV4()
		var packageID = uuid.NewV4()
		var orgID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}

		msgbusClient := &mbmocks.MsgBusServiceClient{}

		subscriberService := &mocks.SubscriberRegistryClientProvider{}
		packageClient := &mocks.PackageClient{}
		simPoolService := &mocks.SimPoolClientProvider{}

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.SubscriberRegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.SubscriberRegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberID.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &subspb.Subscriber{
					SubscriberId: subscriberID.String(),
					NetworkId:    networkID.String(),
					OrgId:        orgID.String(),
				},
			}, nil).Once()

		packageClient.On("GetPackageInfo", packageID.String()).
			Return(&providers.PackageInfo{
				OrgId:    orgID.String(),
				IsActive: true,
				Duration: 3600,
				SimType:  db.SimTypeTest.String(),
			}, nil).Once()

		simPoolClient := simPoolService.On("GetClient").
			Return(&splmocks.SimServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*splmocks.SimServiceClient)

		simPoolResp := simPoolClient.On("Get", mock.Anything,
			&splpb.GetRequest{IsPhysicalSim: false,
				SimType: db.SimTypeTest.String(),
			}).
			Return(&splpb.GetResponse{
				Sim: &splpb.Sim{
					IsPhysical: false,
					SimType:    db.SimTypeTest.String(),
				},
			}, nil).Once().
			ReturnArguments.Get(0).(*splpb.GetResponse)

		sim := &db.Sim{
			SubscriberId: subscriberID,
			NetworkId:    networkID,
			OrgId:        orgID,
			Type:         db.SimTypeTest,
			Status:       sims.SimStatusInactive,
			IsPhysical:   simPoolResp.Sim.IsPhysical,
		}

		simRepo.On("Add", sim,
			mock.Anything).Return(nil).Once()

		pkg := &sims.Package{
			SimId:     sim.Id,
			PackageId: packageID,
			IsActive:  false,
		}

		packageRepo.On("Add", pkg,
			mock.Anything).Return(nil).Once()

		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
		simRepo.On("GetSimMetrics").Return(int64(0), int64(0), int64(0), int64(0), nil).Once()

		s := NewSimManagerServer(simRepo, packageRepo, nil,
			packageClient, subscriberService, simPoolService, "", msgbusClient, "", "", nil, nil)
		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberID.String(),
			NetworkId:    networkID.String(),
			PackageId:    packageID.String(),
			SimType:      sims.SimTypeTest.String(),
			SimToken:     "",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		simRepo.AssertExpectations(t)

		subscriberService.AssertExpectations(t)
		subscriberClient.AssertExpectations(t)

		simPoolService.AssertExpectations(t)
		simPoolClient.AssertExpectations(t)

		packageRepo.AssertExpectations(t)
		packageClient.AssertExpectations(t)
	})

	t.Run("SubscriberNotRegisteredOnProvidedNetwork", func(t *testing.T) {
		var subscriberID = uuid.NewV4()
		var networkID = uuid.NewV4()
		var packageID = uuid.NewV4()

		subscriberService := &mocks.SubscriberRegistryClientProvider{}

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.SubscriberRegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.SubscriberRegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberID.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &subspb.Subscriber{
					SubscriberId: subscriberID.String(),
					NetworkId:    uuid.NewV4().String(),
				},
			}, nil).Once()

		s := NewSimManagerServer(nil, nil, nil, nil, subscriberService, nil, "", nil, "", "", nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberID.String(),
			NetworkId:    networkID.String(),
			PackageId:    packageID.String(),
			SimType:      sims.SimTypeTest.String(),
			SimToken:     "",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		subscriberService.AssertExpectations(t)
		subscriberClient.AssertExpectations(t)
	})

	t.Run("SubscriberOrgAndPackageOrgMismatch", func(t *testing.T) {
		var subscriberID = uuid.NewV4()
		var networkID = uuid.NewV4()
		var packageID = uuid.NewV4()
		var orgID = uuid.NewV4()

		subscriberService := &mocks.SubscriberRegistryClientProvider{}
		packageClient := &mocks.PackageClient{}

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.SubscriberRegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.SubscriberRegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberID.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &subspb.Subscriber{
					SubscriberId: subscriberID.String(),
					NetworkId:    networkID.String(),
					OrgId:        orgID.String(),
				},
			}, nil).Once()

		packageClient.On("GetPackageInfo", packageID.String()).
			Return(
				&providers.PackageInfo{
					OrgId:    uuid.NewV4().String(),
					IsActive: true,
					Duration: 3600,
					SimType:  db.SimTypeTest.String(),
				}, nil).Once()

		s := NewSimManagerServer(nil, nil, nil,
			packageClient, subscriberService, nil, "", nil, "", "", nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberID.String(),
			NetworkId:    networkID.String(),
			PackageId:    packageID.String(),
			SimType:      sims.SimTypeTest.String(),
			SimToken:     "",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		subscriberService.AssertExpectations(t)
		subscriberClient.AssertExpectations(t)

		packageClient.AssertExpectations(t)
	})

	t.Run("OrgPackageNoMoreActive", func(t *testing.T) {
		var subscriberID = uuid.NewV4()
		var networkID = uuid.NewV4()
		var packageID = uuid.NewV4()
		var orgID = uuid.NewV4()

		subscriberService := &mocks.SubscriberRegistryClientProvider{}
		packageClient := &mocks.PackageClient{}

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.SubscriberRegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.SubscriberRegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberID.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &subspb.Subscriber{
					SubscriberId: subscriberID.String(),
					NetworkId:    networkID.String(),
					OrgId:        orgID.String(),
				},
			}, nil).Once()

		packageClient.On("GetPackageInfo", packageID.String()).
			Return(&providers.PackageInfo{
				OrgId:    orgID.String(),
				IsActive: false,
				Duration: 3600,
				SimType:  db.SimTypeTest.String(),
			}, nil).Once()

		s := NewSimManagerServer(nil, nil, nil,
			packageClient, subscriberService, nil, "", nil, "", "", nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberID.String(),
			NetworkId:    networkID.String(),
			PackageId:    packageID.String(),
			SimType:      sims.SimTypeTest.String(),
			SimToken:     "",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		subscriberService.AssertExpectations(t)
		subscriberClient.AssertExpectations(t)

		packageClient.AssertExpectations(t)
	})

	t.Run("PackageSimtypeAndProvidedSimtypeMismatch", func(t *testing.T) {
		var subscriberID = uuid.NewV4()
		var networkID = uuid.NewV4()
		var packageID = uuid.NewV4()
		var orgID = uuid.NewV4()

		subscriberService := &mocks.SubscriberRegistryClientProvider{}
		packageClient := &mocks.PackageClient{}

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.SubscriberRegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.SubscriberRegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberID.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &subspb.Subscriber{
					SubscriberId: subscriberID.String(),
					NetworkId:    networkID.String(),
					OrgId:        orgID.String(),
				},
			}, nil).Once()

		packageClient.On("GetPackageInfo", packageID.String()).
			Return(
				&providers.PackageInfo{
					OrgId:    orgID.String(),
					IsActive: true,
					Duration: 3600,
					SimType:  db.SimTypeUnknown.String(),
				}, nil).Once()

		s := NewSimManagerServer(nil, nil, nil,
			packageClient, subscriberService, nil, "", nil, "", "", nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberID.String(),
			NetworkId:    networkID.String(),
			PackageId:    packageID.String(),
			SimType:      sims.SimTypeTest.String(),
			SimToken:     "",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		subscriberService.AssertExpectations(t)
		subscriberClient.AssertExpectations(t)

		packageClient.AssertExpectations(t)
	})

}

func TestSimManagerServer_SetActivePackageForSim(t *testing.T) {
	t.Run("SimAndPackageFound", func(t *testing.T) {
		var packageID = uuid.NewV4()
		var simID = uuid.NewV4()

		packageRepo := &mocks.PackageRepo{}
		simRepo := &mocks.SimRepo{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				IsPhysical: false,
				Status:     db.SimStatusActive,
			}, nil).
			Once()

		packageRepo.On("Get", packageID).Return(
			&db.Package{Id: packageID,
				SimId:    simID,
				EndDate:  time.Now().AddDate(0, 1, 0), // next month
				IsActive: false,
			}, nil).Once()

		packageRepo.On("Update",
			&sims.Package{
				Id:       packageID,
				IsActive: true,
			},
			mock.Anything).Return(nil).Once()

		s := NewSimManagerServer(simRepo, packageRepo, nil, nil, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
	})

	t.Run("SimStatusInvalid", func(t *testing.T) {
		var packageID = uuid.NewV4()
		var simID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				IsPhysical: false,
				Status:     db.SimStatusUnknown,
			}, nil).
			Once()

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimIdAndPackageSimIdMismatch", func(t *testing.T) {
		var packageID = uuid.NewV4()
		var simID = uuid.NewV4()

		packageRepo := &mocks.PackageRepo{}
		simRepo := &mocks.SimRepo{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				IsPhysical: false,
				Status:     db.SimStatusActive,
			}, nil).
			Once()

		packageRepo.On("Get", packageID).Return(
			&db.Package{Id: packageID,
				SimId:    uuid.NewV4(),
				EndDate:  time.Now().AddDate(0, 1, 0), // next month
				IsActive: false,
			}, nil).Once()

		s := NewSimManagerServer(simRepo, packageRepo, nil, nil, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
	})

	t.Run("PackageAlreadyExpired", func(t *testing.T) {
		var packageID = uuid.NewV4()
		var simID = uuid.NewV4()

		packageRepo := &mocks.PackageRepo{}
		simRepo := &mocks.SimRepo{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				IsPhysical: false,
				Status:     db.SimStatusActive,
			}, nil).
			Once()

		packageRepo.On("Get", packageID).Return(
			&db.Package{Id: packageID,
				SimId:    simID,
				EndDate:  time.Now().AddDate(0, -1, 0), // one month ago
				IsActive: false,
			}, nil).Once()

		s := NewSimManagerServer(simRepo, packageRepo, nil, nil, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_RemovePackageForSim(t *testing.T) {
	t.Run("PackageFound", func(t *testing.T) {
		var packageID = uuid.NewV4()
		var simID = uuid.NewV4()
		simRepo := &mocks.SimRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}

		packageRepo := &mocks.PackageRepo{}

		packageRepo.On("Get", packageID).Return(
			&db.Package{Id: packageID,
				SimId:    simID,
				IsActive: false,
			}, nil).Once()

		packageRepo.On("Delete", packageID,
			mock.Anything).Return(nil).Once()

		s := NewSimManagerServer(simRepo, packageRepo, nil, nil, nil, nil, "", msgbusClient, "", "", nil, nil)
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		resp, err := s.RemovePackageForSim(context.TODO(), &pb.RemovePackageRequest{
			PackageId: packageID.String(),
			SimId:     simID.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		packageRepo.AssertExpectations(t)
	})

	t.Run("PackageDeleteError", func(t *testing.T) {
		var packageID = uuid.Nil
		var simID = uuid.NewV4()

		packageRepo := &mocks.PackageRepo{}

		packageRepo.On("Get", packageID).Return(
			&db.Package{Id: packageID,
				SimId:    simID,
				IsActive: false,
			}, nil).Once()

		packageRepo.On("Delete", packageID,
			mock.Anything).Return(gorm.ErrRecordNotFound).Once()

		s := NewSimManagerServer(nil, packageRepo, nil, nil, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.RemovePackageForSim(context.TODO(), &pb.RemovePackageRequest{
			PackageId: packageID.String(),
			SimId:     simID.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})

	t.Run("PackageUUIDInvalid", func(t *testing.T) {
		var packageID = "1"

		packageRepo := &mocks.PackageRepo{}

		s := NewSimManagerServer(nil, packageRepo, nil, nil, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.RemovePackageForSim(context.TODO(), &pb.RemovePackageRequest{
			PackageId: packageID})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})

	t.Run("SimUUIDMismatch", func(t *testing.T) {
		var packageID = uuid.Nil
		var simID = uuid.NewV4()

		packageRepo := &mocks.PackageRepo{}

		packageRepo.On("Get", packageID).Return(
			&db.Package{Id: packageID,
				SimId:    simID,
				IsActive: false,
			}, nil).Once()

		s := NewSimManagerServer(nil, packageRepo, nil, nil, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.RemovePackageForSim(context.TODO(), &pb.RemovePackageRequest{
			PackageId: packageID.String(),
			SimId:     uuid.NewV4().String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_AddPackageForSim(t *testing.T) {
	t.Run("SimAndPackageFound", func(t *testing.T) {
		var simID = uuid.NewV4()
		var packageID = uuid.NewV4()
		msgbusClient := &mbmocks.MsgBusServiceClient{}

		var orgID = uuid.NewV4()
		startDate := time.Now().UTC().AddDate(0, 0, 1) // tomorrow

		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}
		packageClient := &mocks.PackageClient{}

		sim := simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				OrgId:      orgID,
				Iccid:      testIccid,
				Status:     db.SimStatusInactive,
				Type:       db.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*db.Sim)

		pkgInfo := packageClient.On("GetPackageInfo", packageID.String()).
			Return(&providers.PackageInfo{
				Id:       packageID.String(),
				OrgId:    orgID.String(),
				IsActive: true,
				Duration: 3600,
				SimType:  db.SimTypeTest.String(),
			}, nil).
			Once().
			ReturnArguments.Get(0).(*providers.PackageInfo)

		pkg := &sims.Package{
			SimId:     sim.Id,
			StartDate: startDate,
			EndDate:   startDate.Add(time.Duration(pkgInfo.Duration)),
			PackageId: packageID,
			IsActive:  false,
		}

		packageRepo.On("GetOverlap", pkg).
			Return(nil, nil).Once()

		packageRepo.On("Add", pkg,
			mock.Anything).Return(nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewSimManagerServer(simRepo, packageRepo, nil, packageClient, nil, nil, "", msgbusClient, "", "", nil, nil)

		resp, err := s.AddPackageForSim(context.TODO(), &pb.AddPackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
			StartDate: timestamppb.New(startDate),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		packageClient.AssertExpectations(t)
	})

	t.Run("PackageStartDateNotValid", func(t *testing.T) {
		var simID = uuid.NewV4()
		var packageID = uuid.NewV4()
		startDate := time.Now().UTC()

		s := NewSimManagerServer(nil, nil, nil, nil, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.AddPackageForSim(context.TODO(), &pb.AddPackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
			StartDate: timestamppb.New(startDate),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("OrgPackageNoMoreActive", func(t *testing.T) {
		var simID = uuid.NewV4()
		var packageID = uuid.NewV4()
		var orgID = uuid.NewV4()
		startDate := time.Now().UTC().AddDate(0, 0, 1) // tomorrow

		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}
		packageClient := &mocks.PackageClient{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				OrgId:      orgID,
				Iccid:      testIccid,
				Status:     db.SimStatusInactive,
				Type:       db.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once()

		packageClient.On("GetPackageInfo", packageID.String()).
			Return(&providers.PackageInfo{
				Id:       packageID.String(),
				OrgId:    orgID.String(),
				IsActive: false,
				Duration: 3600,
				SimType:  db.SimTypeTest.String(),
			}, nil).Once()

		s := NewSimManagerServer(simRepo, packageRepo, nil, packageClient, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.AddPackageForSim(context.TODO(), &pb.AddPackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
			StartDate: timestamppb.New(startDate),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		packageClient.AssertExpectations(t)
	})

	t.Run("SimOrgAndPackageOrgMismatch", func(t *testing.T) {
		var simID = uuid.NewV4()
		var packageID = uuid.NewV4()
		var orgID = uuid.NewV4()
		startDate := time.Now().UTC().AddDate(0, 0, 1) // tomorrow

		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}
		packageClient := &mocks.PackageClient{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				OrgId:      orgID,
				Iccid:      testIccid,
				Status:     db.SimStatusInactive,
				Type:       db.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once()

		packageClient.On("GetPackageInfo", packageID.String()).
			Return(&providers.PackageInfo{
				Id:       packageID.String(),
				OrgId:    uuid.NewV4().String(),
				IsActive: true,
				Duration: 3600,
				SimType:  db.SimTypeTest.String(),
			}, nil).Once()

		s := NewSimManagerServer(simRepo, packageRepo, nil, packageClient, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.AddPackageForSim(context.TODO(), &pb.AddPackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
			StartDate: timestamppb.New(startDate),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		packageClient.AssertExpectations(t)
	})

	t.Run("SimSimtypeAndPackageSimtypeMismatch", func(t *testing.T) {
		var simID = uuid.NewV4()
		var packageID = uuid.NewV4()
		var orgID = uuid.NewV4()
		startDate := time.Now().UTC().AddDate(0, 0, 1) // tomorrow

		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}
		packageClient := &mocks.PackageClient{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				OrgId:      orgID,
				Iccid:      testIccid,
				Status:     db.SimStatusInactive,
				Type:       db.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once()

		packageClient.On("GetPackageInfo", packageID.String()).
			Return(&providers.PackageInfo{
				Id:       packageID.String(),
				OrgId:    orgID.String(),
				IsActive: true,
				Duration: 3600,
				SimType:  db.SimTypeUnknown.String(),
			}, nil).Once()

		s := NewSimManagerServer(simRepo, packageRepo, nil, packageClient, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.AddPackageForSim(context.TODO(), &pb.AddPackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
			StartDate: timestamppb.New(startDate),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		packageClient.AssertExpectations(t)
	})

	t.Run("PackageValidityPeriodOverlapsWithExistingPackages", func(t *testing.T) {
		var simID = uuid.NewV4()
		var packageID = uuid.NewV4()
		var orgID = uuid.NewV4()
		startDate := time.Now().UTC().AddDate(0, 0, 1) // tomorrow

		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}
		packageClient := &mocks.PackageClient{}

		sim := simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				OrgId:      orgID,
				Iccid:      testIccid,
				Status:     db.SimStatusInactive,
				Type:       db.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*db.Sim)

		pkgInfo := packageClient.On("GetPackageInfo", packageID.String()).
			Return(&providers.PackageInfo{
				Id:       packageID.String(),
				OrgId:    orgID.String(),
				IsActive: true,
				Duration: 3600,
				SimType:  db.SimTypeTest.String(),
			}, nil).Once().
			ReturnArguments.Get(0).(*providers.PackageInfo)

		pkg := &sims.Package{
			SimId:     sim.Id,
			StartDate: startDate,
			EndDate:   startDate.Add(time.Duration(pkgInfo.Duration)),
			PackageId: packageID,
			IsActive:  false,
		}

		packageRepo.On("GetOverlap", pkg).
			Return([]db.Package{
				{}, {},
			}, nil).Once()

		s := NewSimManagerServer(simRepo, packageRepo, nil, packageClient, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.AddPackageForSim(context.TODO(), &pb.AddPackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
			StartDate: timestamppb.New(startDate),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		packageClient.AssertExpectations(t)
	})

}

func TestSimManagerServer_DeleteSim(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		var simID = uuid.NewV4()
		msgbusClient := &mbmocks.MsgBusServiceClient{}

		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		sim := simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				Iccid:      testIccid,
				Status:     db.SimStatusInactive,
				Type:       db.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*db.Sim)

		agentAdapter := agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("TerminateSim", mock.Anything,
			sim.Iccid).Return(nil).Once()

		simRepo.On("Update",
			&sims.Sim{
				Id:     sim.Id,
				Status: sims.SimStatusTerminated,
			},
			mock.Anything).Return(nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
		simRepo.On("GetSimMetrics").Return(int64(0), int64(0), int64(0), int64(0), nil).Once()

		s := NewSimManagerServer(simRepo, nil, agentFactory, nil, nil, nil, "", msgbusClient, "", "", nil, nil)

		resp, err := s.DeleteSim(context.TODO(), &pb.DeleteSimRequest{
			SimId: simID.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		simRepo.AssertExpectations(t)
		agentFactory.AssertExpectations(t)
	})

	t.Run("SimStatusInvalid", func(t *testing.T) {
		var simID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				Iccid:      testIccid,
				Status:     db.SimStatusActive,
				Type:       db.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once()

		s := NewSimManagerServer(simRepo, nil, nil, nil, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.DeleteSim(context.TODO(), &pb.DeleteSimRequest{
			SimId: simID.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
	})

	t.Run("SimTypeNotSupported", func(t *testing.T) {
		var simID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		sim := simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				Iccid:      testIccid,
				Status:     db.SimStatusInactive,
				Type:       100,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*db.Sim)

		agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, false).
			Once()

		s := NewSimManagerServer(simRepo, nil, agentFactory, nil, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.DeleteSim(context.TODO(), &pb.DeleteSimRequest{
			SimId: simID.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
	})

	t.Run("SimAgentFailToTerminate", func(t *testing.T) {
		var simID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		sim := simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				Iccid:      testIccid,
				Status:     db.SimStatusInactive,
				Type:       db.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*db.Sim)

		agentAdapter := agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("TerminateSim", mock.Anything,
			sim.Iccid).Return(errors.New("anyError")).Once()

		s := NewSimManagerServer(simRepo, nil, agentFactory, nil, nil, nil, "", nil, "", "", nil, nil)

		resp, err := s.DeleteSim(context.TODO(), &pb.DeleteSimRequest{
			SimId: simID.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
	})
}
