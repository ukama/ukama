package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/configurator/mocks"
	"github.com/ukama/ukama/systems/node/configurator/pkg"
	"google.golang.org/protobuf/types/known/anypb"
)

var testNode = ukama.NewVirtualNodeId("HomeNode")
var orgId = uuid.NewV4()

func TestConfiguratorServer_EventNotification(t *testing.T) {
	// Arrange
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	commitRepo := &mocks.CommitRepo{}
	configRepo := &mocks.ConfigRepo{}
	configStore := &mocks.ConfigStoreProvider{}
	registry := &mocks.RegistryProvider{}

	s := NewConfiguratorServer(msgbusClient, registry, configRepo, commitRepo, configStore, testOrgName, pkg.IsDebugMode)

	eventServer := NewConfiguratorEventServer(testOrgName, s)

	/* Node Cretaed event */
	evt := epb.NodeCreatedEvent{
		NodeId: testNode.String(),
		Name:   "testnode",
		Type:   "hnode",
		Org:    orgId.String(),
	}

	any, err := anypb.New(&evt)
	assert.NoError(t, err)

	configRepo.On("Add", testNode.String()).Return(nil).Once()

	_, err = eventServer.EventNotification(context.Background(), &epb.Event{
		RoutingKey: "event.cloud.local.testorg.registry.node.node.create",
		Msg:        any,
	})

	assert.NoError(t, err)
	// Assert
	configStore.AssertExpectations(t)
	assert.NoError(t, err)

}
