package configstore

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/node/configurator/mocks"
	"github.com/ukama/ukama/systems/node/configurator/pkg/db"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
)

var testNode1 = "uk-000000-hnode-0000"
var testNode2 = "uk-000000-hnode-0001"
var Service = "node/configurator"
var TestData = "node/configurator/test/integration/data"

const OrgName = "testorg"

func TestConfigStore_HandleConfigStoreEvent(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	commitRepo := &mocks.CommitRepo{}
	configRepo := &mocks.ConfigRepo{}
	store := &mocks.StoreProvider{}
	registry := &mocks.RegistryProvider{}

	cS := NewConfigStore(msgbusClient, registry, configRepo, commitRepo, OrgName, store, (10 * time.Second))
	t.Run("SameVersion", func(t *testing.T) {
		store.On("GetLatestRemoteConfigs", mock.Anything).Return("000", nil).Once()
		commitRepo.On("GetLatest").Return(&db.Commit{Hash: "000"}, nil).Once()
		store.On("GetRemoteConfigVersion", mock.Anything, mock.Anything).Return(nil).Once()
		err := cS.HandleConfigStoreEvent(context.Background())
		assert.NoError(t, err)
		commitRepo.AssertExpectations(t)
	})

	t.Run("DifferentVersionButNoChanges", func(t *testing.T) {
		store.On("GetLatestRemoteConfigs", mock.Anything).Return("000", nil).Once()
		commitRepo.On("GetLatest").Return(&db.Commit{Hash: "001"}, nil).Once()
		store.On("GetRemoteConfigVersion", mock.Anything, mock.Anything).Return(nil).Once()
		store.On("GetDiff", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
		err := cS.HandleConfigStoreEvent(context.Background())
		assert.NoError(t, err)
		commitRepo.AssertExpectations(t)
	})

}

func TestConfigStore_ProcessConfigStoreEvent(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	commitRepo := &mocks.CommitRepo{}
	configRepo := &mocks.ConfigRepo{}
	store := &mocks.StoreProvider{}
	registry := &mocks.RegistryProvider{}
	cVer := "0.0.0."
	rVer := "0.0.1"
	path, err := os.Getwd()
	assert.NoError(t, err)
	p := strings.Split(path, Service)
	dir := p[0] + TestData
	cS := NewConfigStore(msgbusClient, registry, configRepo, commitRepo, OrgName, store, (10 * time.Second))

	t.Run("DifferentVersionWithChanges", func(t *testing.T) {
		var node string
		store.On("GetDiff", mock.Anything, mock.Anything, mock.Anything).Return([]string{"networkABC/siteXYZ/uk-000000-hnode-0000/epc/epc.json", "networkABC/siteXYZ/uk-000000-hnode-0000/deviced/deviced.json", "networkABC/siteXYZ/uk-000000-hnode-0001/epc/epc.json"}, nil).Once()
		configRepo.On("Get", mock.MatchedBy(func(n string) bool {
			if n == testNode1 {
				node = testNode1
				return true
			} else if n == testNode2 {
				node = testNode2
				return true
			}
			return false
		})).Return(&db.Configuration{NodeId: node}, nil)

		msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)
		configRepo.On("UpdateLastCommit", mock.Anything, mock.MatchedBy(func(a *db.CommitState) bool { return a != nil && *a == db.Published })).Return(nil)
		files, ldir, err := cS.LookingForChanges(dir, cVer, rVer)
		assert.NoError(t, err)
		err = cS.ProcessConfigStoreEvent(files, cVer, ldir)
		assert.NoError(t, err)
		configRepo.AssertExpectations(t)
	})

}
