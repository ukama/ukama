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
		subscriberRegistryProvider := &mocks.SubscriberRegistryClientProvider{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		simRepo.On("Get", mock.Anything).Return(&sims.Sim{
			Id: simId,
			Package: sims.Package{
				Id:              packageId,
				DefaultDuration: 1,
			},
			Type: ukama.SimTypeUkamaData,
		}, nil)

		packageRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Package{
				sims.Package{
					Id:              packageId,
					DefaultDuration: 1,
					SimId:           simId,
				},
			}, nil)

		packageRepo.On("Add", mock.Anything, mock.Anything).Return(nil)

		packageClient.On("Get", mock.Anything).
			Return(&cdplan.PackageInfo{
				Duration:   1,
				IsActive:   true,
				SimType:    ukama.SimTypeUkamaData.String(),
				DataVolume: 10,
				DataUnit:   "GB",
				Name:       "Ukama Package",
				Amount:     100,
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
			subscriberRegistryProvider, networkClient, orgClient, userClient, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("AddFirstPackageIsQueuedNotSetInUseDirectly", func(t *testing.T) {
		localMsgbusClient := &cmocks.MsgBusServiceClient{}
		localMsgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}
		packageClient := &cmocks.PackageClient{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		networkClient := &cmocks.NetworkClient{}
		subscriberRegistryProvider := &mocks.SubscriberRegistryClientProvider{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		simRepo.On("Get", mock.Anything).Return(&sims.Sim{
			Id:   simId,
			Type: ukama.SimTypeUkamaData,
		}, nil)

		packageRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Package{}, nil)

		packageRepo.On("Add", mock.MatchedBy(func(p *sims.Package) bool {
			return !p.IsCurrentlyInUse
		}), mock.Anything).Return(nil)

		packageClient.On("Get", mock.Anything).
			Return(&cdplan.PackageInfo{
				Duration:   1,
				IsActive:   true,
				SimType:    ukama.SimTypeUkamaData.String(),
				DataVolume: 10,
				DataUnit:   "GB",
				Name:       "Ukama Package",
				Amount:     100,
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
			subscriberRegistryProvider, networkClient, orgClient, userClient, localMsgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)

		packageRepo.AssertExpectations(t)
		localMsgbusClient.AssertExpectations(t)
	})

	t.Run("ProvisionedPaymentSkipsPackageAdd", func(t *testing.T) {
		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}
		packageClient := &cmocks.PackageClient{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		// A payment recorded by sim allocation is marked provisioned; the handler
		// must skip adding the package (it was already added at allocation) and
		// must not touch the repos/clients.
		evt := &epb.Payment{
			Id:       uuid.NewV4().String(),
			ItemId:   packageId.String(),
			Status:   ukama.StatusTypeCompleted.String(),
			ItemType: ukama.ItemTypePackage.String(),
			Metadata: []byte(fmt.Sprintf(`{"sim": "%s", "provisioned": "true"}`, simId.String())),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, packageClient,
			nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
		packageRepo.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
		packageClient.AssertNotCalled(t, "Get", mock.Anything)
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, packageClient, nil, nil, nil, nil, msgbusClient, "")
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, packageClient, nil, nil, nil, nil, msgbusClient, "")
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, packageClient, nil, nil, nil, nil, msgbusClient, "")
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, packageClient, nil, nil, nil, nil, msgbusClient, "")
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, packageClient, nil, nil, nil, nil, msgbusClient, "")
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
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
					Type:         ukama.SimTypeOperatorData,
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("UnsupportedCDRType", func(t *testing.T) {
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
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

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestSimManagerEventServer_HandleUkamaAgentCdrCreateEvent(t *testing.T) {
	msgbusClient := &cmocks.MsgBusServiceClient{}
	routingKey := msgbus.PrepareRoute(OrgName,
		"event.cloud.local.{{ .Org}}.ukamaagent.cdr.cdr.create")

	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	t.Run("SimFound", func(t *testing.T) {
		repo := mocks.SimRepo{}

		repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{
					Type: ukama.SimTypeUkamaData,
				},
			}, nil)

		evt := &epb.CDRReported{}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("MultipleSimsFound", func(t *testing.T) {
		repo := mocks.SimRepo{}

		repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{},
				sims.Sim{},
			}, nil)

		evt := &epb.CDRReported{}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		repo := mocks.SimRepo{}

		repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{}, nil)

		evt := &epb.CDRReported{}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("SimListError", func(t *testing.T) {
		repo := mocks.SimRepo{}

		repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("failed to list sim by Iccid"))

		evt := &epb.CDRReported{}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("InvalidEventTypeSent", func(t *testing.T) {
		repo := mocks.SimRepo{}
		evt := &epb.Customer{}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestSimManagerEventServer_HandleUkamaAgentAsrProfileDeleteEvent(t *testing.T) {
	routingKey := msgbus.PrepareRoute(OrgName,
		"event.cloud.local.{{ .Org}}.ukamaagent.asr.activesubscriber.delete")

	t.Run("NextPackagesFound", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{
					Id:     simId,
					Status: ukama.SimStatusServiceOn,
					Type:   ukama.SimTypeUkamaData,
				},
			}, nil)

		simRepo.On("Get", mock.Anything).
			Return(&sims.Sim{
				Id:     simId,
				Status: ukama.SimStatusServiceOn,
				Type:   ukama.SimTypeUkamaData,
			}, nil)

		packageRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		packageRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Package{
				sims.Package{
					Id:               packageId,
					SimId:            simId,
					IsCurrentlyInUse: false,
					IsExpired:        true,
				},
				sims.Package{
					Id:               uuid.NewV4(),
					SimId:            simId,
					IsCurrentlyInUse: true,
					IsExpired:        false,
				},
			}, nil)

		simRepo.On("Update",
			&sims.Sim{
				Id:     simId,
				Status: ukama.SimStatusTerminated,
			},
			mock.Anything).Return(nil).Once()

		evt := &epb.AsrInactivated{
			Subscriber: &epb.Subscriber{
				SimPackage: packageId.String(),
			}}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("NextPackagesListError", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{
					Id:     simId,
					Status: ukama.SimStatusServiceOn,
					Type:   ukama.SimTypeUkamaData,
				},
			}, nil)

		packageRepo.On("Get", mock.Anything).
			Return(&sims.Package{
				SimId:            simId,
				IsCurrentlyInUse: true,
				IsExpired:        false,
			}, nil)

		simRepo.On("Get", mock.Anything).
			Return(&sims.Sim{
				Id:     simId,
				Status: ukama.SimStatusServiceOn,
				Type:   ukama.SimTypeUkamaData,
			}, nil)

		packageRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		packageRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("failed to list next packages"))

		evt := &epb.AsrInactivated{
			Subscriber: &epb.Subscriber{
				SimPackage: packageId.String(),
			}}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("ExpirePackagesError", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{
					Id:     simId,
					Iccid:  testIccid,
					Status: ukama.SimStatusServiceOn,
					Type:   ukama.SimTypeUkamaData,
				},
			}, nil)

		packageRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Package{
				sims.Package{
					Id:               packageId,
					SimId:            simId,
					IsCurrentlyInUse: false,
					IsExpired:        true,
				},
				sims.Package{
					Id:               uuid.NewV4(),
					SimId:            simId,
					IsCurrentlyInUse: true,
					IsExpired:        false,
				},
			}, nil)

		packageRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("package expired update failure"))

		evt := &epb.AsrInactivated{
			Subscriber: &epb.Subscriber{
				Iccid:      testIccid,
				SimPackage: packageId.String(),
			}}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("SimGetError", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{
					Id: simId,
				},
			}, nil)

		packageRepo.On("Get", mock.Anything).
			Return(&sims.Package{
				SimId:            simId,
				IsCurrentlyInUse: true,
				IsExpired:        false,
			}, nil)

		simRepo.On("Get", mock.Anything).
			Return(nil, errors.New("failed to get sim"))

		evt := &epb.AsrInactivated{
			Subscriber: &epb.Subscriber{
				SimPackage: packageId.String(),
			}}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("SimAndPackageIdsMismatch", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{
					Id: simId,
				},
			}, nil)

		packageRepo.On("Get", mock.Anything).
			Return(&sims.Package{
				SimId: uuid.NewV4(),
			}, nil)

		evt := &epb.AsrInactivated{
			Subscriber: &epb.Subscriber{
				SimPackage: packageId.String(),
			}}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("PackageNotFound", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}

		packageId := uuid.NewV4()

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{},
			}, nil)

		packageRepo.On("Get", mock.Anything).Return(nil, errors.New("error while looking up package"))

		evt := &epb.AsrInactivated{
			Subscriber: &epb.Subscriber{
				SimPackage: packageId.String(),
			}}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("PackageIdNotValid", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		repo := mocks.SimRepo{}

		repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{},
			}, nil)

		evt := &epb.AsrInactivated{
			Subscriber: &epb.Subscriber{
				SimPackage: "lol",
			}}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("MultipleSimsFound", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		repo := mocks.SimRepo{}

		repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{},
				sims.Sim{},
			}, nil)

		evt := &epb.AsrInactivated{Subscriber: &epb.Subscriber{}}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		repo := mocks.SimRepo{}

		repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{}, nil)

		evt := &epb.AsrInactivated{Subscriber: &epb.Subscriber{}}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("SimListError", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		repo := mocks.SimRepo{}

		repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("failed to list sim by Iccid"))

		evt := &epb.AsrInactivated{Subscriber: &epb.Subscriber{}}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("InvalidEventTypeSent", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		repo := mocks.SimRepo{}
		evt := &epb.AsrActivated{Subscriber: &epb.Subscriber{}}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestSimManagerEventServer_HandleUkamaAgentAsrPolicyViolationEvent(t *testing.T) {
	routingKey := msgbus.PrepareRoute(OrgName,
		"event.cloud.local.{{ .Org}}.ukamaagent.asr.policy.violation")

	t.Run("DataCapExceededOnlyRecordsUsage", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{
					Id:      simId,
					Type:    ukama.SimTypeUkamaData,
					Package: sims.Package{Id: packageId},
				},
			}, nil)

		packageRepo.On("Update", []uuid.UUID{packageId}, &sims.Package{UsedDataAtExpiry: uint64(500)}, mock.Anything).
			Return(nil).Once()

		evt := &epb.PolicyViolation{
			Profile: &epb.Profile{
				SimPackage:     packageId.String(),
				TotalDataBytes: 500,
			},
			Reason: ukama.PolicyViolationReasonDataCapExceeded.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)

		packageRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("PackageExpiredRollsOverToNextPackage", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Twice()

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()
		nextPackageId := uuid.NewV4()
		dataPlanId := uuid.NewV4()

		simd := &sims.Sim{
			Id:      simId,
			Type:    ukama.SimTypeUkamaData,
			Package: sims.Package{Id: packageId},
		}

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{*simd}, nil)

		simRepo.On("Get", simId).Return(simd, nil).Once()

		packageRepo.On("Update", []uuid.UUID{packageId}, &sims.Package{UsedDataAtExpiry: uint64(1000)}, mock.Anything).
			Return(nil).Once()

		packageRepo.On("Get", packageId).Return(
			&sims.Package{
				Id:               packageId,
				SimId:            simId,
				IsCurrentlyInUse: true,
				IsExpired:        false,
			}, nil)

		packageRepo.On("Update", []uuid.UUID{packageId},
			&sims.Package{IsCurrentlyInUse: false, IsExpired: true}, mock.Anything).
			Return(nil).Once()

		simRepo.On("Get", simId).Return(&sims.Sim{
			Id:   simId,
			Type: ukama.SimTypeUkamaData,
		}, nil).Once()

		simRepo.On("Get", simId).Return(&sims.Sim{
			Id:   simId,
			Type: ukama.SimTypeUkamaData,
		}, nil).Once()

		packageRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Package{
				sims.Package{
					Id:               packageId,
					SimId:            simId,
					IsExpired:        true,
					IsCurrentlyInUse: false,
				},
				sims.Package{
					Id:              nextPackageId,
					SimId:           simId,
					PackageId:       dataPlanId,
					DefaultDuration: 1,
				},
			}, nil)

		packageRepo.On("Get", nextPackageId).Return(
			&sims.Package{
				Id:              nextPackageId,
				SimId:           simId,
				PackageId:       dataPlanId,
				DefaultDuration: 1,
			}, nil)

		packageRepo.On("Update", []uuid.UUID{nextPackageId},
			&sims.Package{IsCurrentlyInUse: true}, mock.Anything).
			Return(nil).Once()

		agentAdapter := agentFactory.On("GetAgentAdapter", ukama.SimTypeUkamaData).
			Return(&mocks.AgentAdapter{}, true).
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("UpdatePackage", mock.Anything, mock.Anything).Return(nil).Once()

		evt := &epb.PolicyViolation{
			Profile: &epb.Profile{
				SimPackage:     packageId.String(),
				TotalDataBytes: 1000,
			},
			Reason: ukama.PolicyViolationReasonPackageExpired.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, agentFactory, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)

		packageRepo.AssertExpectations(t)
		agentAdapter.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("PackageExpiredNoNextPackageQueued", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		simd := &sims.Sim{
			Id:      simId,
			Type:    ukama.SimTypeUkamaData,
			Package: sims.Package{Id: packageId},
		}

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{*simd}, nil)

		simRepo.On("Get", simId).Return(simd, nil).Once()

		packageRepo.On("Update", []uuid.UUID{packageId}, &sims.Package{UsedDataAtExpiry: uint64(1000)}, mock.Anything).
			Return(nil).Once()

		packageRepo.On("Get", packageId).Return(
			&sims.Package{
				Id:               packageId,
				SimId:            simId,
				IsCurrentlyInUse: true,
				IsExpired:        false,
			}, nil)

		packageRepo.On("Update", []uuid.UUID{packageId},
			&sims.Package{IsCurrentlyInUse: false, IsExpired: true}, mock.Anything).
			Return(nil).Once()

		simRepo.On("Get", simId).Return(&sims.Sim{
			Id:   simId,
			Type: ukama.SimTypeUkamaData,
		}, nil).Once()

		packageRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Package{
				sims.Package{
					Id:               packageId,
					SimId:            simId,
					IsExpired:        true,
					IsCurrentlyInUse: false,
				},
			}, nil)

		evt := &epb.PolicyViolation{
			Profile: &epb.Profile{
				SimPackage:     packageId.String(),
				TotalDataBytes: 1000,
			},
			Reason: ukama.PolicyViolationReasonPackageExpired.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)

		packageRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("PackageExpiredWithNoQueuedPackageTurnsSimOff", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Twice()

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		simd := &sims.Sim{
			Id:      simId,
			Type:    ukama.SimTypeUkamaData,
			Status:  ukama.SimStatusServiceOn,
			Package: sims.Package{Id: packageId},
		}

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{*simd}, nil).Once()

		simRepo.On("Get", simId).Return(simd, nil).Once()

		packageRepo.On("Update", []uuid.UUID{packageId}, &sims.Package{UsedDataAtExpiry: uint64(1000)}, mock.Anything).
			Return(nil).Once()

		packageRepo.On("Get", packageId).Return(
			&sims.Package{
				Id:               packageId,
				SimId:            simId,
				IsCurrentlyInUse: true,
				IsExpired:        false,
			}, nil)

		packageRepo.On("Update", []uuid.UUID{packageId},
			&sims.Package{IsCurrentlyInUse: false, IsExpired: true}, mock.Anything).
			Return(nil).Once()

		simRepo.On("Get", simId).Return(&sims.Sim{
			Id:     simId,
			Type:   ukama.SimTypeUkamaData,
			Status: ukama.SimStatusServiceOn,
		}, nil).Once()

		packageRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Package{
				sims.Package{
					Id:               packageId,
					SimId:            simId,
					IsExpired:        true,
					IsCurrentlyInUse: false,
				},
			}, nil)

		simRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{}, nil).Twice()

		evt := &epb.PolicyViolation{
			Profile: &epb.Profile{
				SimPackage:     packageId.String(),
				TotalDataBytes: 1000,
			},
			Reason: ukama.PolicyViolationReasonPackageExpired.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("StaleViolationForNoLongerCurrentPackageIsNoOp", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}

		simId := uuid.NewV4()
		stalePackageId := uuid.NewV4()
		currentPackageId := uuid.NewV4()

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{
					Id:      simId,
					Type:    ukama.SimTypeUkamaData,
					Package: sims.Package{Id: currentPackageId},
				},
			}, nil)

		simRepo.On("Get", simId).Return(&sims.Sim{
			Id:      simId,
			Type:    ukama.SimTypeUkamaData,
			Package: sims.Package{Id: currentPackageId},
		}, nil).Once()

		evt := &epb.PolicyViolation{
			Profile: &epb.Profile{
				SimPackage:     stalePackageId.String(),
				TotalDataBytes: 500,
			},
			Reason: ukama.PolicyViolationReasonPackageExpired.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("DuplicateViolationAfterRolloverHasNoFurtherEffect", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Twice()

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()
		nextPackageId := uuid.NewV4()
		dataPlanId := uuid.NewV4()

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{
					Id:      simId,
					Type:    ukama.SimTypeUkamaData,
					Package: sims.Package{Id: packageId},
				},
			}, nil).Once()

		simRepo.On("Get", simId).Return(&sims.Sim{
			Id:      simId,
			Type:    ukama.SimTypeUkamaData,
			Package: sims.Package{Id: packageId},
		}, nil).Once()

		packageRepo.On("Update", []uuid.UUID{packageId}, &sims.Package{UsedDataAtExpiry: uint64(1000)}, mock.Anything).
			Return(nil).Once()

		packageRepo.On("Get", packageId).Return(
			&sims.Package{
				Id:               packageId,
				SimId:            simId,
				IsCurrentlyInUse: true,
				IsExpired:        false,
			}, nil)

		packageRepo.On("Update", []uuid.UUID{packageId},
			&sims.Package{IsCurrentlyInUse: false, IsExpired: true}, mock.Anything).
			Return(nil).Once()

		simRepo.On("Get", simId).Return(&sims.Sim{
			Id:   simId,
			Type: ukama.SimTypeUkamaData,
		}, nil).Once()

		simRepo.On("Get", simId).Return(&sims.Sim{
			Id:   simId,
			Type: ukama.SimTypeUkamaData,
		}, nil).Once()

		packageRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Package{
				sims.Package{
					Id:               packageId,
					SimId:            simId,
					IsExpired:        true,
					IsCurrentlyInUse: false,
				},
				sims.Package{
					Id:              nextPackageId,
					SimId:           simId,
					PackageId:       dataPlanId,
					DefaultDuration: 1,
				},
			}, nil)

		packageRepo.On("Get", nextPackageId).Return(
			&sims.Package{
				Id:              nextPackageId,
				SimId:           simId,
				PackageId:       dataPlanId,
				DefaultDuration: 1,
			}, nil)

		packageRepo.On("Update", []uuid.UUID{nextPackageId},
			&sims.Package{IsCurrentlyInUse: true}, mock.Anything).
			Return(nil).Once()

		agentAdapter := agentFactory.On("GetAgentAdapter", ukama.SimTypeUkamaData).
			Return(&mocks.AgentAdapter{}, true).
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("UpdatePackage", mock.Anything, mock.Anything).Return(nil).Once()

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{
				sims.Sim{
					Id:      simId,
					Type:    ukama.SimTypeUkamaData,
					Package: sims.Package{Id: nextPackageId},
				},
			}, nil).Once()

		simRepo.On("Get", simId).Return(&sims.Sim{
			Id:      simId,
			Type:    ukama.SimTypeUkamaData,
			Package: sims.Package{Id: nextPackageId},
		}, nil).Once()

		evt := &epb.PolicyViolation{
			Profile: &epb.Profile{
				SimPackage:     packageId.String(),
				TotalDataBytes: 1000,
			},
			Reason: ukama.PolicyViolationReasonPackageExpired.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, agentFactory, nil, nil, nil, nil, nil, msgbusClient, "")

		_, err = s.EventNotification(context.TODO(), msg)
		assert.NoError(t, err)

		_, err = s.EventNotification(context.TODO(), msg)
		assert.NoError(t, err)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		agentAdapter.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})
}

func TestSimManagerEventServer_HandleSimManagerSimAddPackageEvent(t *testing.T) {
	routingKey := msgbus.PrepareRoute(OrgName,
		"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.addpackage")

	t.Run("IdleSimSetsQueuedPackageInUse", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()
		dataPlanId := uuid.NewV4()

		simRepo.On("Get", simId).Return(&sims.Sim{
			Id:   simId,
			Type: ukama.SimTypeTest,
		}, nil).Once()

		simRepo.On("Get", simId).Return(&sims.Sim{
			Id:   simId,
			Type: ukama.SimTypeTest,
		}, nil).Once()

		packageRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Package{
				sims.Package{
					Id:              packageId,
					SimId:           simId,
					PackageId:       dataPlanId,
					DefaultDuration: 1,
				},
			}, nil)

		packageRepo.On("Get", packageId).Return(
			&sims.Package{
				Id:              packageId,
				SimId:           simId,
				PackageId:       dataPlanId,
				DefaultDuration: 1,
			}, nil)

		packageRepo.On("Update", []uuid.UUID{packageId},
			&sims.Package{IsCurrentlyInUse: true}, mock.Anything).
			Return(nil).Once()

		agentAdapter := agentFactory.On("GetAgentAdapter", ukama.SimTypeTest).
			Return(&mocks.AgentAdapter{}, true).
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("UpdatePackage", mock.Anything, mock.Anything).Return(nil).Once()

		evt := &epb.EventSimAddPackage{
			Id:        simId.String(),
			PackageId: dataPlanId.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, agentFactory, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		agentAdapter.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("IdleSimResumesServiceWhenPackageActivated", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Twice()

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()
		dataPlanId := uuid.NewV4()

		simRepo.On("Get", simId).Return(&sims.Sim{
			Id:     simId,
			Type:   ukama.SimTypeTest,
			Status: ukama.SimStatusServiceOff,
		}, nil).Once()

		simRepo.On("Get", simId).Return(&sims.Sim{
			Id:     simId,
			Type:   ukama.SimTypeTest,
			Status: ukama.SimStatusServiceOff,
		}, nil).Once()

		packageRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Package{
				sims.Package{
					Id:              packageId,
					SimId:           simId,
					PackageId:       dataPlanId,
					DefaultDuration: 1,
				},
			}, nil)

		packageRepo.On("Get", packageId).Return(
			&sims.Package{
				Id:              packageId,
				SimId:           simId,
				PackageId:       dataPlanId,
				DefaultDuration: 1,
			}, nil)

		packageRepo.On("Update", []uuid.UUID{packageId},
			&sims.Package{IsCurrentlyInUse: true}, mock.Anything).
			Return(nil).Once()

		agentAdapter := agentFactory.On("GetAgentAdapter", ukama.SimTypeTest).
			Return(&mocks.AgentAdapter{}, true).
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("UpdatePackage", mock.Anything, mock.Anything).Return(nil).Once()

		simRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{}, nil).Twice()

		evt := &epb.EventSimAddPackage{
			Id:        simId.String(),
			PackageId: dataPlanId.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, agentFactory, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		agentAdapter.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("NonIdleSimIsNoOp", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}

		simRepo := mocks.SimRepo{}
		packageRepo := mocks.PackageRepo{}

		simId := uuid.NewV4()
		currentPackageId := uuid.NewV4()

		simRepo.On("Get", simId).Return(&sims.Sim{
			Id:      simId,
			Type:    ukama.SimTypeUkamaData,
			Package: sims.Package{Id: currentPackageId},
		}, nil).Once()

		evt := &epb.EventSimAddPackage{
			Id: simId.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &simRepo, &packageRepo, nil, nil, nil, nil, nil, nil, msgbusClient, "")
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})
}
