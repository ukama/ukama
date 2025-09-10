/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/mocks"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/server"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	cgenukama "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	cdplan "github.com/ukama/ukama/systems/common/rest/client/dataplan"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	subregpb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	subregpbmocks "github.com/ukama/ukama/systems/subscriber/registry/pb/gen/mocks"
	sims "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
)

func TestSimManagerEventServer_HandleSimManagerSimAllocateEvent(t *testing.T) {
	msgbusClient := &cmocks.MsgBusServiceClient{}

	routingKey := msgbus.PrepareRoute(OrgName,
		"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate")

	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	t.Run("AllocatedSimNotFound", func(t *testing.T) {
		repo := mocks.SimRepo{}
		repo.On("Get", mock.Anything).Return(nil, errors.New("sim not found"))

		allocatedSim := epb.EventSimAllocation{
			Id: simId,
		}

		anyE, err := anypb.New(&allocatedSim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("ActiveSimFound", func(t *testing.T) {
		repo := mocks.SimRepo{}
		simId := uuid.NewV4()

		repo.On("Get", mock.Anything).Return(&sims.Sim{
			Id:     simId,
			Status: ukama.SimStatusActive,
		}, nil)

		allocatedSim := epb.EventSimAllocation{
			Id:     simId.String(),
			Status: ukama.SimStatusActive.String(),
		}

		anyE, err := anypb.New(&allocatedSim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("ActiveSimUpdateFailure", func(t *testing.T) {
		repo := mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		repo.On("Get", mock.Anything).Return(&sims.Sim{
			Id:     simId,
			Status: ukama.SimStatusInactive,
			Type:   ukama.SimTypeUkamaData,
		}, nil)

		repo.On("Update", mock.Anything, mock.Anything).
			Return(errors.New("failed to update sim"))

		allocatedSim := epb.EventSimAllocation{
			Id:     simId.String(),
			Status: ukama.SimStatusInactive.String(),
			Type:   ukama.SimTypeUkamaData.String(),
		}

		anyE, err := anypb.New(&allocatedSim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, agentFactory, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("ActiveSimAgentFailure", func(t *testing.T) {
		repo := mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		repo.On("Get", mock.Anything).Return(&sims.Sim{
			Id:     simId,
			Status: ukama.SimStatusInactive,
		}, nil)

		repo.On("Update", mock.Anything, mock.Anything).Return(nil)

		agentFactory.On("GetAgentAdapter", mock.Anything).
			Return(nil, false).Once()

		allocatedSim := epb.EventSimAllocation{
			Id:     simId.String(),
			Status: ukama.SimStatusInactive.String(),
		}

		anyE, err := anypb.New(&allocatedSim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, agentFactory, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("ActiveSimAgentActivateSimFailure", func(t *testing.T) {
		repo := mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		repo.On("Get", mock.Anything).Return(&sims.Sim{
			Id:     simId,
			Status: ukama.SimStatusInactive,
			Type:   ukama.SimTypeUkamaData,
		}, nil)

		repo.On("Update", mock.Anything, mock.Anything).Return(nil)

		agentAdapter := agentFactory.On("GetAgentAdapter", ukama.SimTypeUkamaData).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("ActivateSim", mock.Anything, mock.Anything).
			Return(errors.New("fail to activate sim on remote agent")).Once()

		allocatedSim := epb.EventSimAllocation{
			Id:     simId.String(),
			Status: ukama.SimStatusInactive.String(),
			Type:   ukama.SimTypeUkamaData.String(),
		}

		anyE, err := anypb.New(&allocatedSim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, agentFactory, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("ActiveSimAgentActivateSimSuccess", func(t *testing.T) {
		repo := mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		repo.On("Get", mock.Anything).Return(&sims.Sim{
			Id:     simId,
			Status: ukama.SimStatusInactive,
			Type:   ukama.SimTypeUkamaData,
		}, nil)

		repo.On("Update", mock.Anything, mock.Anything).Return(nil)

		agentAdapter := agentFactory.On("GetAgentAdapter", ukama.SimTypeUkamaData).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("ActivateSim", mock.Anything, mock.Anything).
			Return(nil).Once()

		repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{}, nil)

		allocatedSim := epb.EventSimAllocation{
			Id:     simId.String(),
			Status: ukama.SimStatusInactive.String(),
			Type:   ukama.SimTypeUkamaData.String(),
		}

		anyE, err := anypb.New(&allocatedSim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, agentFactory, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("InvalidEventTypeSent", func(t *testing.T) {
		repo := mocks.SimRepo{}
		evt := &epb.Payment{
			Id: uuid.NewV4().String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestSimManagerEventServer_HandleProcessorPaymentSuccessEvent(t *testing.T) {
	msgbusClient := &cmocks.MsgBusServiceClient{}

	routingKey := msgbus.PrepareRoute(OrgName,
		"event.cloud.local.{{ .Org}}.payments.processor.payment.success")

	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	t.Run("AddPackageSuccess", func(t *testing.T) {
		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}
		packageClient := &cmocks.PackageClient{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		networkClient := &cmocks.NetworkClient{}
		mailerClient := &cmocks.MailerClient{}
		subscriberRegistryProvider := &mocks.SubscriberRegistryClientProvider{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		simRepo.On("Get", mock.Anything).Return(&sims.Sim{
			Id: simId,
			Package: sims.Package{
				Id: packageId,
			},
			Type: ukama.SimTypeUkamaData,
		}, nil)

		packageRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Package{
				sims.Package{
					Id:    packageId,
					SimId: simId,
				},
			}, nil)

		packageRepo.On("Add", mock.Anything, mock.Anything).Return(nil)

		packageClient.On("Get", mock.Anything).
			Return(&cdplan.PackageInfo{
				IsActive: true,
				SimType:  ukama.SimTypeUkamaData.String(),
			}, nil)

		orgClient.On("Get", mock.Anything).
			Return(&cnucl.OrgInfo{}, nil)

		userClient.On("GetById", mock.Anything).
			Return(&cnucl.UserInfo{}, nil)

		subscriberRegistryClient := subscriberRegistryProvider.On("GetClient", mock.Anything).
			Return(&subregpbmocks.RegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subregpbmocks.RegistryServiceClient)

		subscriberRegistryClient.On("Get", mock.Anything, mock.Anything).
			Return(&subregpb.GetSubscriberResponse{
				Subscriber: &cgenukama.Subscriber{},
			}, nil)

		networkClient.On("Get", mock.Anything).
			Return(&creg.NetworkInfo{}, nil)

		mailerClient.On("SendEmail", mock.Anything).
			Return(nil)

		evt := &epb.Payment{
			Id:       uuid.NewV4().String(),
			ItemId:   packageId.String(),
			Status:   ukama.StatusTypeCompleted.String(),
			ItemType: ukama.ItemTypePackage.String(),
			Metadata: []byte(fmt.Sprintf(`{"sim": "%s"}`, simId.String())),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, packageClient,
			subscriberRegistryProvider, networkClient, mailerClient, orgClient, userClient, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("AddPackageError", func(t *testing.T) {
		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}
		packageClient := &cmocks.PackageClient{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		simRepo.On("Get", mock.Anything).Return(&sims.Sim{
			Id: simId,
			Package: sims.Package{
				Id: packageId,
			},
			Type: ukama.SimTypeUkamaData,
		}, nil)

		packageRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Package{}, nil)

		packageRepo.On("Add", mock.Anything, mock.Anything).Return(errors.New("failed to add package to sim"))

		packageClient.On("Get", mock.Anything).
			Return(&cdplan.PackageInfo{
				IsActive: true,
				SimType:  ukama.SimTypeUkamaData.String(),
			}, nil)

		evt := &epb.Payment{
			Id:       uuid.NewV4().String(),
			ItemId:   packageId.String(),
			Status:   ukama.StatusTypeCompleted.String(),
			ItemType: ukama.ItemTypePackage.String(),
			Metadata: []byte(fmt.Sprintf(`{"sim": "%s"}`, simId.String())),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, packageClient, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("SimPackagesListError", func(t *testing.T) {
		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}
		packageClient := &cmocks.PackageClient{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		simRepo.On("Get", mock.Anything).Return(&sims.Sim{
			Id: simId,
			Package: sims.Package{
				Id: packageId,
			},
			Type: ukama.SimTypeUkamaData,
		}, nil)

		packageRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("failed to list current packages on sim"))

		packageClient.On("Get", mock.Anything).
			Return(&cdplan.PackageInfo{
				IsActive: true,
				SimType:  ukama.SimTypeUkamaData.String(),
			}, nil)

		evt := &epb.Payment{
			Id:       uuid.NewV4().String(),
			ItemId:   packageId.String(),
			Status:   ukama.StatusTypeCompleted.String(),
			ItemType: ukama.ItemTypePackage.String(),
			Metadata: []byte(fmt.Sprintf(`{"sim": "%s"}`, simId.String())),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, packageClient, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("SimTypeAndPackageTypeMismatch", func(t *testing.T) {
		repo := mocks.SimRepo{}
		packageClient := &cmocks.PackageClient{}
		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		repo.On("Get", mock.Anything).Return(&sims.Sim{
			Id: simId,
			Package: sims.Package{
				Id: packageId,
			},
			Type: ukama.SimTypeUkamaData,
		}, nil)

		packageClient.On("Get", mock.Anything).
			Return(&cdplan.PackageInfo{
				IsActive: true,
				SimType:  ukama.SimTypeOperatorData.String(),
			}, nil)

		evt := &epb.Payment{
			Id:       uuid.NewV4().String(),
			ItemId:   packageId.String(),
			Status:   ukama.StatusTypeCompleted.String(),
			ItemType: ukama.ItemTypePackage.String(),
			Metadata: []byte(fmt.Sprintf(`{"sim": "%s"}`, simId.String())),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, packageClient, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("PackageNotActive", func(t *testing.T) {
		repo := mocks.SimRepo{}
		packageClient := &cmocks.PackageClient{}
		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		repo.On("Get", mock.Anything).Return(&sims.Sim{
			Id: simId,
			Package: sims.Package{
				Id: packageId,
			},
		}, nil)

		packageClient.On("Get", mock.Anything).
			Return(&cdplan.PackageInfo{
				IsActive: false,
			}, nil)

		evt := &epb.Payment{
			Id:       uuid.NewV4().String(),
			ItemId:   packageId.String(),
			Status:   ukama.StatusTypeCompleted.String(),
			ItemType: ukama.ItemTypePackage.String(),
			Metadata: []byte(fmt.Sprintf(`{"sim": "%s"}`, simId.String())),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, packageClient, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("FailedToGetPackageClient", func(t *testing.T) {
		repo := mocks.SimRepo{}
		packageClient := &cmocks.PackageClient{}
		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		repo.On("Get", mock.Anything).Return(&sims.Sim{
			Id: simId,
			Package: sims.Package{
				Id: packageId,
			},
		}, nil)

		packageClient.On("Get", mock.Anything).
			Return(nil, errors.New("failed to get package client"))

		evt := &epb.Payment{
			Id:       uuid.NewV4().String(),
			ItemId:   packageId.String(),
			Status:   ukama.StatusTypeCompleted.String(),
			ItemType: ukama.ItemTypePackage.String(),
			Metadata: []byte(fmt.Sprintf(`{"sim": "%s"}`, simId.String())),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, packageClient, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("InvalidPackageId", func(t *testing.T) {
		repo := mocks.SimRepo{}
		simId := uuid.NewV4()
		repo.On("Get", mock.Anything).Return(&sims.Sim{
			Id:      simId,
			Package: sims.Package{},
		}, nil)

		evt := &epb.Payment{
			Id:       uuid.NewV4().String(),
			Status:   ukama.StatusTypeCompleted.String(),
			ItemType: ukama.ItemTypePackage.String(),
			Metadata: []byte(fmt.Sprintf(`{"sim": "%s"}`, simId.String())),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		repo := mocks.SimRepo{}
		repo.On("Get", mock.Anything).Return(nil, errors.New("sim not found"))

		evt := &epb.Payment{
			Id:       uuid.NewV4().String(),
			Status:   ukama.StatusTypeCompleted.String(),
			ItemType: ukama.ItemTypePackage.String(),
			Metadata: []byte(`{"sim": "03cb753f-5e03-4c97-8e47-625115476c72"}`),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("PaymentMetadataSimKeyMissing", func(t *testing.T) {
		repo := mocks.SimRepo{}
		evt := &epb.Payment{
			Id:       uuid.NewV4().String(),
			Status:   ukama.StatusTypeCompleted.String(),
			ItemType: ukama.ItemTypePackage.String(),
			Metadata: []byte("{}"),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("InvalidPaymentMetadata", func(t *testing.T) {
		repo := mocks.SimRepo{}
		evt := &epb.Payment{
			Id:       uuid.NewV4().String(),
			Status:   ukama.StatusTypeCompleted.String(),
			ItemType: ukama.ItemTypePackage.String(),
			Metadata: []byte("+++"),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("InvalidPaymentTypeOrStatus", func(t *testing.T) {
		repo := mocks.SimRepo{}
		evt := &epb.Payment{
			Id: uuid.NewV4().String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("InvalidEventTypeSent", func(t *testing.T) {
		repo := mocks.SimRepo{}
		evt := &epb.EventSimAllocation{
			Id: uuid.NewV4().String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestSimManagerEventServer_HandleOperatorCdrCreateEvent(t *testing.T) {
	msgbusClient := &cmocks.MsgBusServiceClient{}

	routingKey := msgbus.PrepareRoute(OrgName,
		"event.cloud.local.{{ .Org}}.operator.cdr.cdr.create")

	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	t.Run("MultipleSimIccidFound", func(t *testing.T) {
		repo := mocks.SimRepo{}

		simId := uuid.NewV4()
		subscriberId := uuid.NewV4()

		repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{
					Id:           simId,
					SubscriberId: subscriberId,
				},
			}, nil)

		evt := &epb.EventOperatorCdrReport{
			Iccid: testIccid,
			Type:  ukama.CdrTypeData.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("MultipleSimIccidFound", func(t *testing.T) {
		repo := mocks.SimRepo{}
		repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{},
				sims.Sim{},
			}, nil)

		evt := &epb.EventOperatorCdrReport{
			Iccid: testIccid,
			Type:  ukama.CdrTypeData.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("SimIccidNotFound", func(t *testing.T) {
		repo := mocks.SimRepo{}
		repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{}, nil)

		evt := &epb.EventOperatorCdrReport{
			Iccid: testIccid,
			Type:  ukama.CdrTypeData.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("SimIccidListError", func(t *testing.T) {
		repo := mocks.SimRepo{}
		repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("failed to list sim by Iccid"))

		evt := &epb.EventOperatorCdrReport{
			Iccid: testIccid,
			Type:  ukama.CdrTypeData.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("UnsupportableCDRType", func(t *testing.T) {
		repo := mocks.SimRepo{}
		evt := &epb.EventOperatorCdrReport{
			Type: ukama.CdrTypeSms.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("InvalidEventTypeSent", func(t *testing.T) {
		repo := mocks.SimRepo{}
		evt := &epb.EventArtifactChunkReady{}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}
